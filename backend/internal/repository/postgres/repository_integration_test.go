//go:build integration

// File: repository_integration_test.go
// Purpose: Validate PostgreSQL repository behavior with a real database.
// Module: backend/internal/repository/postgres, integration test layer.
// Related: repository domain files and SQL migrations.
package postgres

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/anxcye/ancy-blog/backend/internal/domain"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func newIntegrationRepo(t *testing.T) *Repository {
	t.Helper()
	dsn := strings.TrimSpace(os.Getenv("TEST_DB_DSN"))
	if dsn == "" {
		t.Skip("TEST_DB_DSN is not set")
	}
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatalf("open test db failed: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		t.Fatalf("ping test db failed: %v", err)
	}
	repo := &Repository{db: db}
	resetTestDatabase(t, db)
	return repo
}

func resetTestDatabase(t *testing.T, db *sql.DB) {
	t.Helper()
	execSQLFile(t, db, migrationFilePath(t, "000001_init.down.sql"))
	execSQLFile(t, db, migrationFilePath(t, "000001_init.up.sql"))
	execSQLFile(t, db, migrationFilePath(t, "000002_translation_job_result.up.sql"))
}

func migrationFilePath(t *testing.T, name string) string {
	t.Helper()
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatalf("runtime caller unavailable")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(currentFile), "../../../migrations", name))
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

func TestRepositoryIntegration_CreateAndListPublishedArticles(t *testing.T) {
	repo := newIntegrationRepo(t)
	now := time.Now().UTC()

	created, err := repo.CreateArticle(domain.Article{
		Title:         "Integration Article",
		Slug:          "integration-article",
		ContentKind:   "post",
		Summary:       "summary",
		Content:       "content",
		Status:        "published",
		Visibility:    "public",
		AllowComment:  true,
		OriginType:    "original",
		AIAssistLevel: "none",
		PublishedAt:   now,
	})
	if err != nil {
		t.Fatalf("create article failed: %v", err)
	}
	if created.ID == "" {
		t.Fatalf("expected article id")
	}

	rows, total := repo.ListPublishedArticles(1, 10, "", "", "post")
	if total < 1 {
		t.Fatalf("expected total >=1, got %d", total)
	}
	if len(rows) == 0 {
		t.Fatalf("expected at least one row")
	}

	got, ok := repo.GetPublishedArticleBySlug("integration-article")
	if !ok {
		t.Fatalf("expected published article by slug")
	}
	if got.Title != "Integration Article" {
		t.Fatalf("unexpected article title: %s", got.Title)
	}
}

func TestRepositoryIntegration_CommentFlow(t *testing.T) {
	repo := newIntegrationRepo(t)
	now := time.Now().UTC()

	article, err := repo.CreateArticle(domain.Article{
		Title:         "Comment Target",
		Slug:          "comment-target",
		ContentKind:   "post",
		Summary:       "summary",
		Content:       "content",
		Status:        "published",
		Visibility:    "public",
		AllowComment:  true,
		OriginType:    "original",
		AIAssistLevel: "none",
		PublishedAt:   now,
	})
	if err != nil {
		t.Fatalf("create article failed: %v", err)
	}

	comment, err := repo.CreateComment(domain.Comment{
		ArticleID: article.ID,
		Content:   "hello",
		Status:    "approved",
		Nickname:  "tester",
		Source:    "web",
		IP:        "127.0.0.1",
		UserAgent: "integration-test",
	})
	if err != nil {
		t.Fatalf("create comment failed: %v", err)
	}
	if comment.ID == "" {
		t.Fatalf("expected comment id")
	}

	rows, total := repo.ListArticleComments(article.ID, 1, 10)
	if total != 1 || len(rows) != 1 {
		t.Fatalf("expected one comment, total=%d len=%d", total, len(rows))
	}

	count, err := repo.CountArticleComments(article.ID)
	if err != nil {
		t.Fatalf("count comments failed: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected count=1, got %d", count)
	}
}

func TestRepositoryIntegration_TranslationJobLifecycle(t *testing.T) {
	repo := newIntegrationRepo(t)
	now := time.Now().UTC()

	article, err := repo.CreateArticle(domain.Article{
		Title:         "Translate Me",
		Slug:          "translate-me",
		ContentKind:   "post",
		Summary:       "summary",
		Content:       "content",
		Status:        "published",
		Visibility:    "public",
		AllowComment:  true,
		OriginType:    "original",
		AIAssistLevel: "none",
		PublishedAt:   now,
	})
	if err != nil {
		t.Fatalf("create article failed: %v", err)
	}

	_, err = repo.UpdateIntegrationProvider("openai_compatible", true, []byte(`{"base_url":"https://example.com/v1","api_key":"k","model":"gpt-4.1-mini"}`), []byte(`{"health":"ok"}`))
	if err != nil {
		t.Fatalf("update provider failed: %v", err)
	}

	job, err := repo.CreateTranslationJob(domain.TranslationJob{
		SourceType:   "article",
		SourceID:     article.ID,
		SourceLocale: "zh-CN",
		TargetLocale: "en-US",
		ProviderKey:  "openai_compatible",
		ModelName:    "gpt-4.1-mini",
		Status:       "queued",
	})
	if err != nil {
		t.Fatalf("create job failed: %v", err)
	}

	claimed, ok, err := repo.ClaimNextQueuedTranslationJob()
	if err != nil {
		t.Fatalf("claim queued job failed: %v", err)
	}
	if !ok {
		t.Fatalf("expected claimed job")
	}
	if claimed.ID != job.ID || claimed.Status != "running" {
		t.Fatalf("unexpected claimed job: %#v", claimed)
	}

	if err := repo.MarkTranslationJobSucceeded(job.ID, "translated content"); err != nil {
		t.Fatalf("mark succeeded failed: %v", err)
	}

	stored, ok := repo.GetTranslationJobByID(job.ID)
	if !ok {
		t.Fatalf("expected stored job")
	}
	if stored.Status != "succeeded" || stored.ResultText != "translated content" {
		t.Fatalf("unexpected stored job: %#v", stored)
	}

	sourceText, ok, err := repo.GetTranslationSourceText("article", article.ID)
	if err != nil {
		t.Fatalf("get source text failed: %v", err)
	}
	if !ok || strings.TrimSpace(sourceText) == "" {
		t.Fatalf("expected non-empty source text")
	}
}
