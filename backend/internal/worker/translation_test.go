// File: translation_test.go
// Purpose: Verify translation worker job processing behavior with mocked LLM responses.
// Module: backend/internal/worker, worker unit test layer.
// Related: translation.go and service translation/integration facades.
package worker

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/anxcye/ancy-blog/backend/internal/domain"
	"github.com/anxcye/ancy-blog/backend/internal/repository"
	"github.com/anxcye/ancy-blog/backend/internal/service"
)

type workerRepoStub struct {
	repository.ContentRepository
	mu sync.Mutex

	jobToClaim  domain.TranslationJob
	hasJob      bool
	sourceText  string
	sourceFound bool
	provider    domain.IntegrationProvider
	providerOK  bool

	upsertCalls []upsertCall
	succeeded   map[string]string
	failed      map[string]string
}

type upsertCall struct {
	sourceType string
	sourceID   string
	locale     string
	content    string
	jobID      string
}

func newWorkerRepoStub() *workerRepoStub {
	return &workerRepoStub{
		succeeded: map[string]string{},
		failed:    map[string]string{},
	}
}

func (s *workerRepoStub) ClaimNextQueuedTranslationJob() (domain.TranslationJob, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.hasJob {
		return domain.TranslationJob{}, false, nil
	}
	s.hasJob = false
	job := s.jobToClaim
	job.Status = "running"
	return job, true, nil
}

func (s *workerRepoStub) GetTranslationSourceText(sourceType, sourceID string) (string, bool, error) {
	_ = sourceType
	_ = sourceID
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.sourceText, s.sourceFound, nil
}

func (s *workerRepoStub) GetIntegrationProvider(providerKey string) (domain.IntegrationProvider, bool) {
	_ = providerKey
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.provider, s.providerOK
}

func (s *workerRepoStub) UpsertArticleTranslation(articleID, locale, content, translatedByJobID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.upsertCalls = append(s.upsertCalls, upsertCall{sourceType: "article", sourceID: articleID, locale: locale, content: content, jobID: translatedByJobID})
	return nil
}

func (s *workerRepoStub) UpsertMomentTranslation(momentID, locale, content, translatedByJobID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.upsertCalls = append(s.upsertCalls, upsertCall{sourceType: "moment", sourceID: momentID, locale: locale, content: content, jobID: translatedByJobID})
	return nil
}

func (s *workerRepoStub) MarkTranslationJobSucceeded(id, resultText string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.succeeded[id] = resultText
	return nil
}

func (s *workerRepoStub) MarkTranslationJobFailed(id, errorMessage string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.failed[id] = errorMessage
	return nil
}

func buildWorkerForTest(repo *workerRepoStub, pollInterval time.Duration) *TranslationWorker {
	core := service.NewContentService(repo, nil)
	w := NewTranslationWorker(
		slog.New(slog.NewJSONHandler(io.Discard, nil)),
		service.NewTranslationService(core),
		service.NewIntegrationService(core),
		pollInterval,
	)
	return w
}

func TestTranslationWorkerProcessOnceSuccess(t *testing.T) {
	llmServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/chat/completions" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"choices":[{"message":{"content":"translated text"}}]}`))
	}))
	defer llmServer.Close()

	cfg := map[string]any{
		"base_url": llmServer.URL,
		"api_key":  "test-key",
		"model":    "gpt-4.1-mini",
	}
	rawCfg, _ := json.Marshal(cfg)

	repo := newWorkerRepoStub()
	repo.hasJob = true
	repo.jobToClaim = domain.TranslationJob{
		ID:           "job-1",
		SourceType:   "article",
		SourceID:     "article-1",
		SourceLocale: "zh-CN",
		TargetLocale: "en-US",
		ProviderKey:  "openai_compatible",
		ModelName:    "gpt-4.1-mini",
	}
	repo.sourceFound = true
	repo.sourceText = "原文"
	repo.providerOK = true
	repo.provider = domain.IntegrationProvider{ProviderKey: "openai_compatible", Enabled: true, ConfigJSON: rawCfg}

	worker := buildWorkerForTest(repo, time.Second)
	if err := worker.processOnce(context.Background()); err != nil {
		t.Fatalf("processOnce failed: %v", err)
	}

	if got := repo.succeeded["job-1"]; got != "translated text" {
		t.Fatalf("expected succeeded result, got %q", got)
	}
	if len(repo.upsertCalls) != 1 {
		t.Fatalf("expected one upsert call, got %d", len(repo.upsertCalls))
	}
	if repo.upsertCalls[0].locale != "en-US" {
		t.Fatalf("expected locale en-US, got %s", repo.upsertCalls[0].locale)
	}
	if len(repo.failed) != 0 {
		t.Fatalf("expected no failed jobs, got %+v", repo.failed)
	}
}

func TestTranslationWorkerProcessOnceProviderDisabled(t *testing.T) {
	repo := newWorkerRepoStub()
	repo.hasJob = true
	repo.jobToClaim = domain.TranslationJob{ID: "job-2", SourceType: "article", SourceID: "article-1", TargetLocale: "en-US", ProviderKey: "openai_compatible"}
	repo.sourceFound = true
	repo.sourceText = "content"
	repo.providerOK = true
	repo.provider = domain.IntegrationProvider{ProviderKey: "openai_compatible", Enabled: false, ConfigJSON: []byte(`{"base_url":"https://example.com","api_key":"k","model":"m"}`)}

	worker := buildWorkerForTest(repo, time.Second)
	if err := worker.processOnce(context.Background()); err != nil {
		t.Fatalf("processOnce failed: %v", err)
	}
	msg := repo.failed["job-2"]
	if !strings.Contains(msg, "disabled") {
		t.Fatalf("expected disabled error, got %q", msg)
	}
	if len(repo.succeeded) != 0 {
		t.Fatalf("expected no succeeded jobs, got %+v", repo.succeeded)
	}
}

func TestTranslationWorkerProcessOnceLLMFailure(t *testing.T) {
	llmServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "bad gateway", http.StatusBadGateway)
	}))
	defer llmServer.Close()

	rawCfg, _ := json.Marshal(map[string]any{"base_url": llmServer.URL, "api_key": "k", "model": "m"})
	repo := newWorkerRepoStub()
	repo.hasJob = true
	repo.jobToClaim = domain.TranslationJob{ID: "job-3", SourceType: "article", SourceID: "article-1", TargetLocale: "en-US", ProviderKey: "openai_compatible", ModelName: "m"}
	repo.sourceFound = true
	repo.sourceText = "content"
	repo.providerOK = true
	repo.provider = domain.IntegrationProvider{ProviderKey: "openai_compatible", Enabled: true, ConfigJSON: rawCfg}

	worker := buildWorkerForTest(repo, time.Second)
	if err := worker.processOnce(context.Background()); err != nil {
		t.Fatalf("processOnce failed: %v", err)
	}
	msg := repo.failed["job-3"]
	if !strings.Contains(msg, "status 502") {
		t.Fatalf("expected llm status error, got %q", msg)
	}
}

func TestTranslationWorkerProcessOnceEmptyOutput(t *testing.T) {
	llmServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"choices":[{"message":{"content":"   "}}]}`))
	}))
	defer llmServer.Close()

	rawCfg, _ := json.Marshal(map[string]any{"base_url": llmServer.URL, "api_key": "k", "model": "m"})
	repo := newWorkerRepoStub()
	repo.hasJob = true
	repo.jobToClaim = domain.TranslationJob{ID: "job-4", SourceType: "article", SourceID: "article-1", TargetLocale: "en-US", ProviderKey: "openai_compatible", ModelName: "m"}
	repo.sourceFound = true
	repo.sourceText = "content"
	repo.providerOK = true
	repo.provider = domain.IntegrationProvider{ProviderKey: "openai_compatible", Enabled: true, ConfigJSON: rawCfg}

	worker := buildWorkerForTest(repo, time.Second)
	if err := worker.processOnce(context.Background()); err != nil {
		t.Fatalf("processOnce failed: %v", err)
	}
	msg := repo.failed["job-4"]
	if !strings.Contains(msg, "empty translation output") {
		t.Fatalf("expected empty output error, got %q", msg)
	}
	if len(repo.upsertCalls) != 0 {
		t.Fatalf("expected no upsert when output empty")
	}
}
