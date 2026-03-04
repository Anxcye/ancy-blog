// File: admin_test.go
// Purpose: Verify key admin handler responses for integration and translation APIs.
// Module: backend/internal/handler, admin HTTP test layer.
// Related: admin.go, middleware auth context, and content service.
package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/anxcye/ancy-blog/backend/internal/apperr"
	"github.com/anxcye/ancy-blog/backend/internal/domain"
	"github.com/anxcye/ancy-blog/backend/internal/middleware"
	"github.com/anxcye/ancy-blog/backend/internal/response"
	"github.com/anxcye/ancy-blog/backend/internal/service"
	"github.com/gin-gonic/gin"
)

func adminRouter(h *AdminHandler) *gin.Engine {
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set(middleware.ContextUserKey, domain.User{ID: "admin-1", Username: "admin", IsAdmin: true})
		c.Next()
	})
	return r
}

func TestAdminListIntegrations(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &handlerRepoStub{listIntegrationProvidersFunc: func(providerType string) []domain.IntegrationProvider {
		return []domain.IntegrationProvider{{ProviderKey: "cloudflare_r2", ProviderType: "object_storage", Enabled: true, ConfigJSON: []byte(`{"access_key_id":"x"}`)}}
	}}
	core := service.NewContentService(repo, nil)
	h := NewAdminHandler(
		service.NewArticleService(core),
		service.NewCommentService(core),
		service.NewLinkService(core),
		service.NewSiteService(core),
		service.NewIntegrationService(core),
		service.NewTranslationService(core),
		service.NewAIAssistService(service.NewArticleService(core), service.NewIntegrationService(core)),
	)
	r := adminRouter(h)
	r.GET("/integrations", h.ListIntegrations)

	req := httptest.NewRequest(http.MethodGet, "/integrations", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestAdminListArticles(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &handlerRepoStub{
		listArticlesFunc: func(page, pageSize int, status, contentKind, keyword string) ([]domain.Article, int) {
			if page != 1 || pageSize != 10 || status != "draft" || contentKind != "post" || keyword != "go" {
				t.Fatalf("unexpected query params page=%d pageSize=%d status=%s contentKind=%s keyword=%s", page, pageSize, status, contentKind, keyword)
			}
			return []domain.Article{{ID: "a1", Title: "Go Draft", Slug: "go-draft", Status: "draft", ContentKind: "post"}}, 1
		},
	}
	core := service.NewContentService(repo, nil)
	h := NewAdminHandler(
		service.NewArticleService(core),
		service.NewCommentService(core),
		service.NewLinkService(core),
		service.NewSiteService(core),
		service.NewIntegrationService(core),
		service.NewTranslationService(core),
		service.NewAIAssistService(service.NewArticleService(core), service.NewIntegrationService(core)),
	)
	r := adminRouter(h)
	r.GET("/articles", h.ListArticles)

	req := httptest.NewRequest(http.MethodGet, "/articles?status=draft&contentKind=post&keyword=go", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestAdminArticleDetailNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &handlerRepoStub{}
	core := service.NewContentService(repo, nil)
	h := NewAdminHandler(
		service.NewArticleService(core),
		service.NewCommentService(core),
		service.NewLinkService(core),
		service.NewSiteService(core),
		service.NewIntegrationService(core),
		service.NewTranslationService(core),
		service.NewAIAssistService(service.NewArticleService(core), service.NewIntegrationService(core)),
	)
	r := adminRouter(h)
	r.GET("/articles/:id", h.ArticleDetail)

	req := httptest.NewRequest(http.MethodGet, "/articles/missing", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestAdminListMoments(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &handlerRepoStub{
		listMomentsFunc: func(page, pageSize int, status string) ([]domain.Moment, int) {
			if page != 1 || pageSize != 10 || status != "draft" {
				t.Fatalf("unexpected query params page=%d pageSize=%d status=%s", page, pageSize, status)
			}
			return []domain.Moment{{ID: "m1", Content: "draft", Status: "draft"}}, 1
		},
	}
	core := service.NewContentService(repo, nil)
	h := NewAdminHandler(
		service.NewArticleService(core),
		service.NewCommentService(core),
		service.NewLinkService(core),
		service.NewSiteService(core),
		service.NewIntegrationService(core),
		service.NewTranslationService(core),
		service.NewAIAssistService(service.NewArticleService(core), service.NewIntegrationService(core)),
	)
	r := adminRouter(h)
	r.GET("/moments", h.ListMoments)

	req := httptest.NewRequest(http.MethodGet, "/moments?status=draft", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestAdminUpdateIntegrationProviderNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &handlerRepoStub{updateIntegrationProvider: func(string, bool, []byte, []byte) (domain.IntegrationProvider, error) {
		return domain.IntegrationProvider{}, apperr.ErrProviderNotFound
	}}
	core := service.NewContentService(repo, nil)
	h := NewAdminHandler(
		service.NewArticleService(core),
		service.NewCommentService(core),
		service.NewLinkService(core),
		service.NewSiteService(core),
		service.NewIntegrationService(core),
		service.NewTranslationService(core),
		service.NewAIAssistService(service.NewArticleService(core), service.NewIntegrationService(core)),
	)
	r := adminRouter(h)
	r.PUT("/integrations/:providerKey", h.UpdateIntegration)

	body := bytes.NewBufferString(`{"enabled":true,"configJson":{},"metaJson":{}}`)
	req := httptest.NewRequest(http.MethodPut, "/integrations/missing", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
	var env response.Envelope
	if err := json.Unmarshal(w.Body.Bytes(), &env); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if env.Code != "PROVIDER_NOT_FOUND" {
		t.Fatalf("unexpected code: %s", env.Code)
	}
}

func TestAdminCreateTranslationJobSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &handlerRepoStub{
		getIntegrationProviderFunc: func(string) (domain.IntegrationProvider, bool) {
			return domain.IntegrationProvider{ProviderKey: "openai_compatible", ProviderType: "llm", Enabled: true}, true
		},
		createTranslationJobFunc: func(job domain.TranslationJob) (domain.TranslationJob, error) {
			job.ID = "job-1"
			return job, nil
		},
	}
	core := service.NewContentService(repo, nil)
	h := NewAdminHandler(
		service.NewArticleService(core),
		service.NewCommentService(core),
		service.NewLinkService(core),
		service.NewSiteService(core),
		service.NewIntegrationService(core),
		service.NewTranslationService(core),
		service.NewAIAssistService(service.NewArticleService(core), service.NewIntegrationService(core)),
	)
	r := adminRouter(h)
	r.POST("/translations/jobs", h.CreateTranslationJob)

	body := bytes.NewBufferString(`{"sourceType":"article","sourceId":"a1","sourceLocale":"zh-CN","targetLocale":"en-US","providerKey":"openai_compatible","modelName":"gpt-4.1-mini"}`)
	req := httptest.NewRequest(http.MethodPost, "/translations/jobs", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestAdminTranslationJobDetailNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	core := service.NewContentService(&handlerRepoStub{}, nil)
	h := NewAdminHandler(
		service.NewArticleService(core),
		service.NewCommentService(core),
		service.NewLinkService(core),
		service.NewSiteService(core),
		service.NewIntegrationService(core),
		service.NewTranslationService(core),
		service.NewAIAssistService(service.NewArticleService(core), service.NewIntegrationService(core)),
	)
	r := adminRouter(h)
	r.GET("/translations/jobs/:id", h.TranslationJobDetail)

	req := httptest.NewRequest(http.MethodGet, "/translations/jobs/job-x", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestAdminRetryTranslationJobSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &handlerRepoStub{
		retryTranslationJobFunc: func(id string) (domain.TranslationJob, error) {
			return domain.TranslationJob{ID: id, Status: "queued"}, nil
		},
	}
	core := service.NewContentService(repo, nil)
	h := NewAdminHandler(
		service.NewArticleService(core),
		service.NewCommentService(core),
		service.NewLinkService(core),
		service.NewSiteService(core),
		service.NewIntegrationService(core),
		service.NewTranslationService(core),
		service.NewAIAssistService(service.NewArticleService(core), service.NewIntegrationService(core)),
	)
	r := adminRouter(h)
	r.POST("/translations/jobs/:id/retry", h.RetryTranslationJob)

	req := httptest.NewRequest(http.MethodPost, "/translations/jobs/job-1/retry", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestAdminRetryTranslationJobNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &handlerRepoStub{
		retryTranslationJobFunc: func(id string) (domain.TranslationJob, error) {
			return domain.TranslationJob{}, apperr.ErrTranslationJobNotFound
		},
	}
	core := service.NewContentService(repo, nil)
	h := NewAdminHandler(
		service.NewArticleService(core),
		service.NewCommentService(core),
		service.NewLinkService(core),
		service.NewSiteService(core),
		service.NewIntegrationService(core),
		service.NewTranslationService(core),
		service.NewAIAssistService(service.NewArticleService(core), service.NewIntegrationService(core)),
	)
	r := adminRouter(h)
	r.POST("/translations/jobs/:id/retry", h.RetryTranslationJob)

	req := httptest.NewRequest(http.MethodPost, "/translations/jobs/job-1/retry", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestAdminListTranslationContents(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &handlerRepoStub{
		listTranslationContentsFunc: func(_ int, _ int, sourceType, sourceID, locale string) ([]domain.TranslationContent, int) {
			if sourceType != "article" || sourceID != "a1" || locale != "en-US" {
				t.Fatalf("unexpected filters sourceType=%s sourceID=%s locale=%s", sourceType, sourceID, locale)
			}
			return []domain.TranslationContent{{SourceType: "article", SourceID: "a1", Locale: "en-US", Content: "hello"}}, 1
		},
	}
	core := service.NewContentService(repo, nil)
	h := NewAdminHandler(
		service.NewArticleService(core),
		service.NewCommentService(core),
		service.NewLinkService(core),
		service.NewSiteService(core),
		service.NewIntegrationService(core),
		service.NewTranslationService(core),
		service.NewAIAssistService(service.NewArticleService(core), service.NewIntegrationService(core)),
	)
	r := adminRouter(h)
	r.GET("/translations/contents", h.ListTranslationContents)

	req := httptest.NewRequest(http.MethodGet, "/translations/contents?sourceType=article&sourceId=a1&locale=en-US", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestAdminUpsertTranslationContent(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &handlerRepoStub{
		upsertTranslationContentFunc: func(sourceType, sourceID, locale, title, summary, content, status string, publishedAt time.Time, translatedByJobID string) (domain.TranslationContent, error) {
			return domain.TranslationContent{
				SourceType:        sourceType,
				SourceID:          sourceID,
				Locale:            locale,
				Title:             title,
				Summary:           summary,
				Content:           content,
				Status:            status,
				PublishedAt:       publishedAt,
				TranslatedByJobID: translatedByJobID,
			}, nil
		},
	}
	core := service.NewContentService(repo, nil)
	h := NewAdminHandler(
		service.NewArticleService(core),
		service.NewCommentService(core),
		service.NewLinkService(core),
		service.NewSiteService(core),
		service.NewIntegrationService(core),
		service.NewTranslationService(core),
		service.NewAIAssistService(service.NewArticleService(core), service.NewIntegrationService(core)),
	)
	r := adminRouter(h)
	r.PUT("/translations/contents", h.UpsertTranslationContent)

	body := bytes.NewBufferString(`{"sourceType":"article","sourceId":"a1","locale":"en-US","content":"translated","translatedByJobId":"job-1"}`)
	req := httptest.NewRequest(http.MethodPut, "/translations/contents", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestAdminTranslationContentDetailNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	core := service.NewContentService(&handlerRepoStub{}, nil)
	h := NewAdminHandler(
		service.NewArticleService(core),
		service.NewCommentService(core),
		service.NewLinkService(core),
		service.NewSiteService(core),
		service.NewIntegrationService(core),
		service.NewTranslationService(core),
		service.NewAIAssistService(service.NewArticleService(core), service.NewIntegrationService(core)),
	)
	r := adminRouter(h)
	r.GET("/translations/contents/:sourceType/:sourceId/:locale", h.TranslationContentDetail)

	req := httptest.NewRequest(http.MethodGet, "/translations/contents/article/a1/en-US", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestAdminAISuggestSlugFallback(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &handlerRepoStub{slugExistsFunc: func(slug string) bool { return slug == "my-title" }}
	core := service.NewContentService(repo, nil)
	h := NewAdminHandler(
		service.NewArticleService(core),
		service.NewCommentService(core),
		service.NewLinkService(core),
		service.NewSiteService(core),
		service.NewIntegrationService(core),
		service.NewTranslationService(core),
		service.NewAIAssistService(service.NewArticleService(core), service.NewIntegrationService(core)),
	)
	r := adminRouter(h)
	r.POST("/ai/slug", h.SuggestSlug)

	body := bytes.NewBufferString(`{"title":"My Title"}`)
	req := httptest.NewRequest(http.MethodPost, "/ai/slug", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	var env response.Envelope
	if err := json.Unmarshal(w.Body.Bytes(), &env); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	data, ok := env.Data.(map[string]any)
	if !ok {
		t.Fatalf("unexpected data type: %T", env.Data)
	}
	if data["slug"] != "my-title-2" {
		t.Fatalf("unexpected slug: %#v", data)
	}
}

func TestAdminAIGenerateSummaryFallback(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &handlerRepoStub{}
	core := service.NewContentService(repo, nil)
	h := NewAdminHandler(
		service.NewArticleService(core),
		service.NewCommentService(core),
		service.NewLinkService(core),
		service.NewSiteService(core),
		service.NewIntegrationService(core),
		service.NewTranslationService(core),
		service.NewAIAssistService(service.NewArticleService(core), service.NewIntegrationService(core)),
	)
	r := adminRouter(h)
	r.POST("/ai/summary", h.GenerateSummary)

	body := bytes.NewBufferString(`{"title":"T","content":"This is a long body for fallback summary.","maxLength":10}`)
	req := httptest.NewRequest(http.MethodPost, "/ai/summary", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
}
