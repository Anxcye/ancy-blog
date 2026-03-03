//go:build integration

// File: e2e_smoke_integration_test.go
// Purpose: Validate end-to-end API smoke flows with real PostgreSQL and full HTTP routing.
// Module: backend/internal/server, integration test layer.
// Related: app wiring, auth/public/admin handlers, and migrations.
package server

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/anxcye/ancy-blog/backend/internal/config"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type envelope struct {
	Code    string          `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

func TestAPISmokeFlow(t *testing.T) {
	dsn := strings.TrimSpace(os.Getenv("TEST_DB_DSN"))
	if dsn == "" {
		t.Skip("TEST_DB_DSN is not set")
	}
	dbCfg := parseDBConfigFromDSN(t, dsn)
	resetDBForE2E(t, dsn)

	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	cfg := &config.Config{
		App:  config.AppConfig{Name: "ancy-blog-api", Env: "test"},
		HTTP: config.HTTPConfig{Host: "127.0.0.1", Port: 0},
		Auth: config.AuthConfig{
			AdminUsername:          "admin",
			AdminPassword:          "123456",
			AccessTokenTTLSeconds:  3600,
			RefreshTokenTTLSeconds: 604800,
		},
		DB:    dbCfg,
		Redis: config.RedisConfig{Enabled: false},
		Translation: config.TranslationConfig{
			WorkerEnabled:  false,
			PollIntervalMS: 3000,
		},
	}

	srv, err := NewHTTPServer(cfg, logger)
	if err != nil {
		t.Fatalf("new http server failed: %v", err)
	}
	ts := httptest.NewServer(srv.server.Handler)
	defer ts.Close()

	token := loginAndGetToken(t, ts.URL)
	articleID := createArticle(t, ts.URL, token)
	readArticleBySlug(t, ts.URL)
	createComment(t, ts.URL, articleID)
	listCommentsByAdmin(t, ts.URL, token)
	enableLLMIntegration(t, ts.URL, token)
	jobID := createTranslationJob(t, ts.URL, token, articleID)
	getTranslationJobDetail(t, ts.URL, token, jobID)
}

func loginAndGetToken(t *testing.T, baseURL string) string {
	t.Helper()
	body := `{"username":"admin","password":"123456"}`
	env := doJSON(t, http.MethodPost, baseURL+"/api/v1/auth/login", "", body)
	if env.Code != "OK" {
		t.Fatalf("login failed code=%s msg=%s", env.Code, env.Message)
	}
	var data struct {
		AccessToken string `json:"accessToken"`
	}
	if err := json.Unmarshal(env.Data, &data); err != nil {
		t.Fatalf("unmarshal login data failed: %v", err)
	}
	if data.AccessToken == "" {
		t.Fatalf("access token is empty")
	}
	return data.AccessToken
}

func createArticle(t *testing.T, baseURL, token string) string {
	t.Helper()
	body := `{"title":"Smoke Post","slug":"smoke-post","contentKind":"post","summary":"smoke","content":"hello smoke","status":"published","visibility":"public","allowComment":true,"originType":"original","aiAssistLevel":"none"}`
	env := doJSON(t, http.MethodPost, baseURL+"/api/v1/admin/articles", token, body)
	if env.Code != "OK" {
		t.Fatalf("create article failed code=%s msg=%s", env.Code, env.Message)
	}
	var data map[string]string
	if err := json.Unmarshal(env.Data, &data); err != nil {
		t.Fatalf("unmarshal article data failed: %v", err)
	}
	id := data["id"]
	if id == "" {
		t.Fatalf("article id is empty")
	}
	return id
}

func readArticleBySlug(t *testing.T, baseURL string) {
	t.Helper()
	env := doJSON(t, http.MethodGet, baseURL+"/api/v1/public/articles/smoke-post", "", "")
	if env.Code != "OK" {
		t.Fatalf("read article failed code=%s msg=%s", env.Code, env.Message)
	}
}

func createComment(t *testing.T, baseURL, articleID string) {
	t.Helper()
	body := fmt.Sprintf(`{"articleId":"%s","nickname":"tester","content":"nice"}`, articleID)
	env := doJSON(t, http.MethodPost, baseURL+"/api/v1/public/comments", "", body)
	if env.Code != "OK" {
		t.Fatalf("create comment failed code=%s msg=%s", env.Code, env.Message)
	}
}

func listCommentsByAdmin(t *testing.T, baseURL, token string) {
	t.Helper()
	env := doJSON(t, http.MethodGet, baseURL+"/api/v1/admin/comments?page=1&pageSize=10", token, "")
	if env.Code != "OK" {
		t.Fatalf("list comments failed code=%s msg=%s", env.Code, env.Message)
	}
}

func enableLLMIntegration(t *testing.T, baseURL, token string) {
	t.Helper()
	body := `{"enabled":true,"configJson":{"base_url":"https://example.com/v1","api_key":"k","model":"gpt-4.1-mini"},"metaJson":{"health":"ok"}}`
	env := doJSON(t, http.MethodPut, baseURL+"/api/v1/admin/integrations/openai_compatible", token, body)
	if env.Code != "OK" {
		t.Fatalf("enable llm integration failed code=%s msg=%s", env.Code, env.Message)
	}
}

func createTranslationJob(t *testing.T, baseURL, token, articleID string) string {
	t.Helper()
	body := fmt.Sprintf(`{"sourceType":"article","sourceId":"%s","sourceLocale":"zh-CN","targetLocale":"en-US","providerKey":"openai_compatible","modelName":"gpt-4.1-mini"}`, articleID)
	env := doJSON(t, http.MethodPost, baseURL+"/api/v1/admin/translations/jobs", token, body)
	if env.Code != "OK" {
		t.Fatalf("create translation job failed code=%s msg=%s", env.Code, env.Message)
	}
	var data map[string]string
	if err := json.Unmarshal(env.Data, &data); err != nil {
		t.Fatalf("unmarshal translation job create data failed: %v", err)
	}
	jobID := data["id"]
	if jobID == "" {
		t.Fatalf("translation job id is empty")
	}
	return jobID
}

func getTranslationJobDetail(t *testing.T, baseURL, token, jobID string) {
	t.Helper()
	env := doJSON(t, http.MethodGet, baseURL+"/api/v1/admin/translations/jobs/"+jobID, token, "")
	if env.Code != "OK" {
		t.Fatalf("translation job detail failed code=%s msg=%s", env.Code, env.Message)
	}
}

func doJSON(t *testing.T, method, endpoint, token, body string) envelope {
	t.Helper()
	var reqBody io.Reader
	if body != "" {
		reqBody = bytes.NewBufferString(body)
	}
	req, err := http.NewRequest(method, endpoint, reqBody)
	if err != nil {
		t.Fatalf("new request failed: %v", err)
	}
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("do request failed: %v", err)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	var env envelope
	if err := json.Unmarshal(raw, &env); err != nil {
		t.Fatalf("unmarshal envelope failed status=%d body=%s err=%v", resp.StatusCode, string(raw), err)
	}
	if resp.StatusCode >= 400 {
		t.Fatalf("unexpected status=%d code=%s message=%s body=%s", resp.StatusCode, env.Code, env.Message, string(raw))
	}
	return env
}

func parseDBConfigFromDSN(t *testing.T, dsn string) config.DBConfig {
	t.Helper()
	u, err := url.Parse(dsn)
	if err != nil {
		t.Fatalf("parse TEST_DB_DSN failed: %v", err)
	}
	password, _ := u.User.Password()
	port := 5432
	if strings.Contains(u.Host, ":") {
		parts := strings.Split(u.Host, ":")
		u.Host = parts[0]
		fmt.Sscanf(parts[1], "%d", &port)
	}
	sslMode := u.Query().Get("sslmode")
	if sslMode == "" {
		sslMode = "disable"
	}
	return config.DBConfig{
		Host:         u.Host,
		Port:         port,
		Name:         strings.TrimPrefix(u.Path, "/"),
		User:         u.User.Username(),
		Password:     password,
		SSLMode:      sslMode,
		MaxOpenConns: 20,
		MaxIdleConns: 5,
	}
}

func resetDBForE2E(t *testing.T, dsn string) {
	t.Helper()
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatalf("open test db failed: %v", err)
	}
	defer db.Close()
	execSQLFile(t, db, migrationFilePath(t, "000001_init.down.sql"))
	execSQLFile(t, db, migrationFilePath(t, "000001_init.up.sql"))
	execSQLFile(t, db, migrationFilePath(t, "000002_translation_job_result.up.sql"))
	execSQLFile(t, db, migrationFilePath(t, "000003_content_translations.up.sql"))
}

func migrationFilePath(t *testing.T, name string) string {
	t.Helper()
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatalf("runtime caller unavailable")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(currentFile), "../../migrations", name))
}

func execSQLFile(t *testing.T, db *sql.DB, path string) {
	t.Helper()
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read sql file failed: %s err=%v", path, err)
	}
	if _, err := db.Exec(string(body)); err != nil {
		t.Fatalf("exec sql file failed: %s err=%v", path, err)
	}
}
