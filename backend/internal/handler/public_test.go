// File: public_test.go
// Purpose: Verify key public handler responses and request-to-service integration.
// Module: backend/internal/handler, public HTTP test layer.
// Related: public.go and service.ContentService.
package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/anxcye/ancy-blog/backend/internal/domain"
	"github.com/anxcye/ancy-blog/backend/internal/response"
	"github.com/anxcye/ancy-blog/backend/internal/service"
	"github.com/gin-gonic/gin"
)

func TestPublicArticleBySlugNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &handlerRepoStub{}
	core := service.NewContentService(repo, nil)
	h := NewPublicHandler(
		service.NewArticleService(core),
		service.NewCommentService(core),
		service.NewLinkService(core),
		service.NewSiteService(core),
		service.NewTimelineService(core),
	)

	r := gin.New()
	r.GET("/articles/:slug", h.ArticleBySlug)

	req := httptest.NewRequest(http.MethodGet, "/articles/not-exist", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestPublicAddCommentSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	captured := domain.Comment{}
	repo := &handlerRepoStub{createCommentFunc: func(comment domain.Comment) (domain.Comment, error) {
		captured = comment
		comment.ID = "c1"
		return comment, nil
	}}
	core := service.NewContentService(repo, nil)
	h := NewPublicHandler(
		service.NewArticleService(core),
		service.NewCommentService(core),
		service.NewLinkService(core),
		service.NewSiteService(core),
		service.NewTimelineService(core),
	)

	r := gin.New()
	r.POST("/comments", h.AddComment)

	body := bytes.NewBufferString(`{"articleId":"a1","nickname":"ancy","content":"hello"}`)
	req := httptest.NewRequest(http.MethodPost, "/comments", body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "test-agent")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	if captured.IP == "" {
		t.Fatalf("expected client IP captured")
	}
	if captured.UserAgent != "test-agent" {
		t.Fatalf("expected user agent captured, got %s", captured.UserAgent)
	}
}

func TestPublicSiteFooterGroupedByRow(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &handlerRepoStub{listFooterItemsFunc: func() []domain.FooterItem {
		return []domain.FooterItem{{ID: "1", Label: "About", RowNum: 1, OrderNum: 1}, {ID: "2", Label: "ICP", RowNum: 2, OrderNum: 1}}
	}}
	core := service.NewContentService(repo, nil)
	h := NewPublicHandler(
		service.NewArticleService(core),
		service.NewCommentService(core),
		service.NewLinkService(core),
		service.NewSiteService(core),
		service.NewTimelineService(core),
	)

	r := gin.New()
	r.GET("/footer", h.SiteFooter)

	req := httptest.NewRequest(http.MethodGet, "/footer", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var env response.Envelope
	if err := json.Unmarshal(w.Body.Bytes(), &env); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if env.Code != "OK" {
		t.Fatalf("unexpected code: %s", env.Code)
	}
}

func TestPublicMomentsPassesLocale(t *testing.T) {
	gin.SetMode(gin.TestMode)
	capturedLocale := ""
	repo := &handlerRepoStub{
		listPublishedMomentsFunc: func(_ int, _ int, locale string) ([]domain.Moment, int) {
			capturedLocale = locale
			return []domain.Moment{{ID: "m1", Content: "hello"}}, 1
		},
	}
	core := service.NewContentService(repo, nil)
	h := NewPublicHandler(
		service.NewArticleService(core),
		service.NewCommentService(core),
		service.NewLinkService(core),
		service.NewSiteService(core),
		service.NewTimelineService(core),
	)

	r := gin.New()
	r.GET("/moments", h.Moments)

	req := httptest.NewRequest(http.MethodGet, "/moments?locale=en-US", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if capturedLocale != "en-US" {
		t.Fatalf("expected locale en-US, got %s", capturedLocale)
	}
}

func TestPublicTimelinePassesLocale(t *testing.T) {
	gin.SetMode(gin.TestMode)
	capturedLocale := ""
	repo := &handlerRepoStub{
		listTimelineFunc: func(_ int, _ int, locale string) ([]domain.TimelineItem, int) {
			capturedLocale = locale
			return []domain.TimelineItem{{ContentType: "moment", ID: "m1", Content: "hello"}}, 1
		},
	}
	core := service.NewContentService(repo, nil)
	h := NewPublicHandler(
		service.NewArticleService(core),
		service.NewCommentService(core),
		service.NewLinkService(core),
		service.NewSiteService(core),
		service.NewTimelineService(core),
	)

	r := gin.New()
	r.GET("/timeline", h.Timeline)

	req := httptest.NewRequest(http.MethodGet, "/timeline?locale=en-US", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if capturedLocale != "en-US" {
		t.Fatalf("expected locale en-US, got %s", capturedLocale)
	}
}
