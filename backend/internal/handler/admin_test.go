// File: admin_test.go
// Purpose: Verify key admin handler responses for integration and translation APIs.
// Module: backend/internal/handler, admin HTTP test layer.
// Related: admin.go, middleware auth context, and content service.
package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

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
	h := NewAdminHandler(service.NewContentService(repo, nil))
	r := adminRouter(h)
	r.GET("/integrations", h.ListIntegrations)

	req := httptest.NewRequest(http.MethodGet, "/integrations", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestAdminUpdateIntegrationProviderNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &handlerRepoStub{updateIntegrationProvider: func(string, bool, []byte, []byte) (domain.IntegrationProvider, error) {
		return domain.IntegrationProvider{}, errors.New("provider not found")
	}}
	h := NewAdminHandler(service.NewContentService(repo, nil))
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
	h := NewAdminHandler(service.NewContentService(repo, nil))
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
	h := NewAdminHandler(service.NewContentService(&handlerRepoStub{}, nil))
	r := adminRouter(h)
	r.GET("/translations/jobs/:id", h.TranslationJobDetail)

	req := httptest.NewRequest(http.MethodGet, "/translations/jobs/job-x", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}
