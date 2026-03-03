// File: content_test.go
// Purpose: Verify content service business rules, including validation, defaults, and cache behavior.
// Module: backend/internal/service, content unit test layer.
// Related: content.go and repository contracts.
package service

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/anxcye/ancy-blog/backend/internal/domain"
	"github.com/anxcye/ancy-blog/backend/internal/repository"
)

type contentRepoStub struct {
	repository.ContentRepository

	createArticleFunc             func(article domain.Article) (domain.Article, error)
	updateCommentAdminFunc        func(id, status, isPinned string) (domain.Comment, error)
	submitLinkFunc                func(link domain.Link) (domain.Link, error)
	reviewLinkFunc                func(id, reviewStatus, reviewNote, relatedArticleID string) (domain.Link, error)
	getArticleByIDFunc            func(id string) (domain.Article, bool)
	getPublishedArticleBySlugFunc func(slug string) (domain.Article, bool)
	createFooterItemFunc          func(item domain.FooterItem) (domain.FooterItem, error)
	createSlotItemFunc            func(slotKey string, item domain.SlotItem) (domain.SlotItem, error)
	getSiteSettingsFunc           func() domain.SiteSettings
	listIntegrationProvidersFunc  func(providerType string) []domain.IntegrationProvider

	getIntegrationProviderFunc    func(providerKey string) (domain.IntegrationProvider, bool)
	updateIntegrationProviderFunc func(providerKey string, enabled bool, configJSON, metaJSON []byte) (domain.IntegrationProvider, error)
	createTranslationJobFunc      func(job domain.TranslationJob) (domain.TranslationJob, error)
}

func (s *contentRepoStub) CreateArticle(article domain.Article) (domain.Article, error) {
	if s.createArticleFunc != nil {
		return s.createArticleFunc(article)
	}
	return article, nil
}

func (s *contentRepoStub) UpdateCommentAdmin(id, status, isPinned string) (domain.Comment, error) {
	if s.updateCommentAdminFunc != nil {
		return s.updateCommentAdminFunc(id, status, isPinned)
	}
	return domain.Comment{}, nil
}

func (s *contentRepoStub) SubmitLink(link domain.Link) (domain.Link, error) {
	if s.submitLinkFunc != nil {
		return s.submitLinkFunc(link)
	}
	return link, nil
}

func (s *contentRepoStub) ReviewLink(id, reviewStatus, reviewNote, relatedArticleID string) (domain.Link, error) {
	if s.reviewLinkFunc != nil {
		return s.reviewLinkFunc(id, reviewStatus, reviewNote, relatedArticleID)
	}
	return domain.Link{}, nil
}

func (s *contentRepoStub) GetArticleByID(id string) (domain.Article, bool) {
	if s.getArticleByIDFunc != nil {
		return s.getArticleByIDFunc(id)
	}
	return domain.Article{}, false
}

func (s *contentRepoStub) GetPublishedArticleBySlug(slug string) (domain.Article, bool) {
	if s.getPublishedArticleBySlugFunc != nil {
		return s.getPublishedArticleBySlugFunc(slug)
	}
	return domain.Article{}, false
}

func (s *contentRepoStub) CreateFooterItem(item domain.FooterItem) (domain.FooterItem, error) {
	if s.createFooterItemFunc != nil {
		return s.createFooterItemFunc(item)
	}
	return item, nil
}

func (s *contentRepoStub) CreateSlotItem(slotKey string, item domain.SlotItem) (domain.SlotItem, error) {
	if s.createSlotItemFunc != nil {
		return s.createSlotItemFunc(slotKey, item)
	}
	return item, nil
}

func (s *contentRepoStub) GetSiteSettings() domain.SiteSettings {
	if s.getSiteSettingsFunc != nil {
		return s.getSiteSettingsFunc()
	}
	return domain.SiteSettings{}
}

func (s *contentRepoStub) ListIntegrationProviders(providerType string) []domain.IntegrationProvider {
	if s.listIntegrationProvidersFunc != nil {
		return s.listIntegrationProvidersFunc(providerType)
	}
	return nil
}

func (s *contentRepoStub) GetIntegrationProvider(providerKey string) (domain.IntegrationProvider, bool) {
	if s.getIntegrationProviderFunc != nil {
		return s.getIntegrationProviderFunc(providerKey)
	}
	return domain.IntegrationProvider{}, false
}

func (s *contentRepoStub) UpdateIntegrationProvider(providerKey string, enabled bool, configJSON, metaJSON []byte) (domain.IntegrationProvider, error) {
	if s.updateIntegrationProviderFunc != nil {
		return s.updateIntegrationProviderFunc(providerKey, enabled, configJSON, metaJSON)
	}
	return domain.IntegrationProvider{}, nil
}

func (s *contentRepoStub) CreateTranslationJob(job domain.TranslationJob) (domain.TranslationJob, error) {
	if s.createTranslationJobFunc != nil {
		return s.createTranslationJobFunc(job)
	}
	return job, nil
}

type fakeCache struct {
	store map[string]string
}

func (c *fakeCache) Get(_ context.Context, key string) (string, bool, error) {
	v, ok := c.store[key]
	if !ok {
		return "", false, nil
	}
	return v, true, nil
}

func (c *fakeCache) Set(_ context.Context, key, value string, _ time.Duration) error {
	if c.store == nil {
		c.store = map[string]string{}
	}
	c.store[key] = value
	return nil
}

func (c *fakeCache) Del(_ context.Context, keys ...string) error {
	for _, key := range keys {
		delete(c.store, key)
	}
	return nil
}

func TestCreateArticleSetsDefaults(t *testing.T) {
	repo := &contentRepoStub{createArticleFunc: func(article domain.Article) (domain.Article, error) {
		if article.ContentKind != "post" || article.Status != "draft" || article.Visibility != "public" || article.OriginType != "original" || article.AIAssistLevel != "none" {
			t.Fatalf("unexpected defaults: %#v", article)
		}
		article.ID = "a1"
		return article, nil
	}}
	svc := NewContentService(repo, nil)

	created, err := svc.CreateArticle(domain.Article{Title: "Hello", Slug: "hello"})
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if created.ID != "a1" {
		t.Fatalf("unexpected article id: %s", created.ID)
	}
}

func TestCreateArticleValidation(t *testing.T) {
	svc := NewContentService(&contentRepoStub{}, nil)
	if _, err := svc.CreateArticle(domain.Article{Slug: "x"}); err == nil {
		t.Fatalf("expected error for missing title")
	}
}

func TestUpdateCommentAdminDefaultStatus(t *testing.T) {
	repo := &contentRepoStub{updateCommentAdminFunc: func(_ string, status, _ string) (domain.Comment, error) {
		if status != "approved" {
			t.Fatalf("expected default status approved, got %s", status)
		}
		return domain.Comment{ID: "c1", Status: status}, nil
	}}
	svc := NewContentService(repo, nil)
	if _, err := svc.UpdateCommentAdmin("c1", "", "0"); err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
}

func TestSubmitLinkValidation(t *testing.T) {
	svc := NewContentService(&contentRepoStub{}, nil)
	if _, err := svc.SubmitLink(domain.Link{Name: "x", URL: "not-a-url"}); err == nil {
		t.Fatalf("expected invalid url error")
	}
}

func TestReviewLinkValidation(t *testing.T) {
	svc := NewContentService(&contentRepoStub{}, nil)
	if _, err := svc.ReviewLink("id", "bad", "", ""); err == nil {
		t.Fatalf("expected invalid review status error")
	}

	svc = NewContentService(&contentRepoStub{getArticleByIDFunc: func(string) (domain.Article, bool) {
		return domain.Article{}, false
	}}, nil)
	if _, err := svc.ReviewLink("id", "approved", "", "missing"); err == nil {
		t.Fatalf("expected related article not found error")
	}
}

func TestCreateFooterItemValidation(t *testing.T) {
	svc := NewContentService(&contentRepoStub{}, nil)
	if _, err := svc.CreateFooterItem(domain.FooterItem{Label: "A", RowNum: 4}); err == nil {
		t.Fatalf("expected rowNum validation error")
	}

	svc = NewContentService(&contentRepoStub{getPublishedArticleBySlugFunc: func(string) (domain.Article, bool) {
		return domain.Article{ContentKind: "post"}, true
	}}, nil)
	if _, err := svc.CreateFooterItem(domain.FooterItem{Label: "A", RowNum: 1, LinkType: "internal", InternalArticleSlug: "about"}); err == nil {
		t.Fatalf("expected internal page validation error")
	}
}

func TestCreateSlotItemValidation(t *testing.T) {
	svc := NewContentService(&contentRepoStub{}, nil)
	if _, err := svc.CreateSlotItem("home", domain.SlotItem{ContentType: "bad"}); err == nil {
		t.Fatalf("expected contentType validation error")
	}

	svc = NewContentService(&contentRepoStub{getArticleByIDFunc: func(string) (domain.Article, bool) {
		return domain.Article{}, false
	}}, nil)
	if _, err := svc.CreateSlotItem("home", domain.SlotItem{ContentType: "article", ContentID: "a1"}); err == nil {
		t.Fatalf("expected article not found error")
	}
}

func TestGetSiteSettingsReadsFromCache(t *testing.T) {
	calls := 0
	repo := &contentRepoStub{getSiteSettingsFunc: func() domain.SiteSettings {
		calls++
		return domain.SiteSettings{SiteName: "Ancy"}
	}}
	cache := &fakeCache{store: map[string]string{}}
	svc := NewContentService(repo, cache)

	_ = svc.GetSiteSettings()
	_ = svc.GetSiteSettings()
	if calls != 1 {
		t.Fatalf("expected one repository call due to cache, got %d", calls)
	}
}

func TestListIntegrationProvidersMasksSecrets(t *testing.T) {
	repo := &contentRepoStub{listIntegrationProvidersFunc: func(string) []domain.IntegrationProvider {
		return []domain.IntegrationProvider{{ProviderKey: "cloudflare_r2", ConfigJSON: []byte(`{"secret_access_key":"abc","public_base_url":"x"}`)}}
	}}
	svc := NewContentService(repo, nil)
	rows := svc.ListIntegrationProviders("")
	if len(rows) != 1 {
		t.Fatalf("expected one provider")
	}
	var payload map[string]any
	if err := json.Unmarshal(rows[0].ConfigJSON, &payload); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if payload["secret_access_key"] != "******" {
		t.Fatalf("expected secret masked, got %#v", payload)
	}
}

func TestUpdateIntegrationProviderMasksSecrets(t *testing.T) {
	repo := &contentRepoStub{
		updateIntegrationProviderFunc: func(providerKey string, enabled bool, configJSON, metaJSON []byte) (domain.IntegrationProvider, error) {
			return domain.IntegrationProvider{
				ProviderKey: providerKey,
				Enabled:     enabled,
				ConfigJSON:  configJSON,
				MetaJSON:    metaJSON,
			}, nil
		},
	}
	svc := NewContentService(repo, nil)

	config := []byte(`{"access_key_id":"abc","secret_access_key":"def","public_base_url":"https://cdn.example.com"}`)
	meta := []byte(`{"health":"ok"}`)
	got, err := svc.UpdateIntegrationProvider("cloudflare_r2", true, config, meta)
	if err != nil {
		t.Fatalf("expected update success, got error: %v", err)
	}

	var payload map[string]any
	if err := json.Unmarshal(got.ConfigJSON, &payload); err != nil {
		t.Fatalf("failed to parse masked config json: %v", err)
	}
	if payload["access_key_id"] != "******" || payload["secret_access_key"] != "******" {
		t.Fatalf("expected secret keys to be masked, got: %#v", payload)
	}
}

func TestUpdateIntegrationProviderInvalidConfigJSON(t *testing.T) {
	svc := NewContentService(&contentRepoStub{}, nil)
	if _, err := svc.UpdateIntegrationProvider("cloudflare_r2", true, []byte("{"), nil); err == nil {
		t.Fatalf("expected validation error for invalid config json")
	}
}

func TestTestIntegrationProviderValidation(t *testing.T) {
	cases := []struct {
		name    string
		provide func(string) (domain.IntegrationProvider, bool)
		wantErr bool
	}{
		{
			name: "provider not found",
			provide: func(string) (domain.IntegrationProvider, bool) {
				return domain.IntegrationProvider{}, false
			},
			wantErr: true,
		},
		{
			name: "provider disabled",
			provide: func(string) (domain.IntegrationProvider, bool) {
				return domain.IntegrationProvider{ProviderKey: "openai_compatible", Enabled: false, ConfigJSON: []byte(`{}`)}, true
			},
			wantErr: true,
		},
		{
			name: "missing required config",
			provide: func(string) (domain.IntegrationProvider, bool) {
				return domain.IntegrationProvider{ProviderKey: "openai_compatible", Enabled: true, ConfigJSON: []byte(`{"base_url":"https://example.com"}`)}, true
			},
			wantErr: true,
		},
		{
			name: "valid openai compatible config",
			provide: func(string) (domain.IntegrationProvider, bool) {
				return domain.IntegrationProvider{
					ProviderKey: "openai_compatible",
					Enabled:     true,
					ConfigJSON:  []byte(`{"base_url":"https://example.com","api_key":"k","model":"gpt-4.1-mini"}`),
				}, true
			},
			wantErr: false,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			svc := NewContentService(&contentRepoStub{getIntegrationProviderFunc: tc.provide}, nil)
			_, err := svc.TestIntegrationProvider("openai_compatible")
			if tc.wantErr && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("expected success, got error: %v", err)
			}
		})
	}
}

func TestCreateTranslationJobValidation(t *testing.T) {
	svc := NewContentService(&contentRepoStub{
		getIntegrationProviderFunc: func(string) (domain.IntegrationProvider, bool) {
			return domain.IntegrationProvider{ProviderKey: "openai_compatible", ProviderType: "llm", Enabled: true}, true
		},
		createTranslationJobFunc: func(job domain.TranslationJob) (domain.TranslationJob, error) {
			job.ID = "job-1"
			return job, nil
		},
	}, nil)

	job, err := svc.CreateTranslationJob(domain.TranslationJob{
		SourceType:   "article",
		SourceID:     "a1",
		SourceLocale: "zh-CN",
		TargetLocale: "en-US",
		ProviderKey:  "openai_compatible",
		ModelName:    "gpt-4.1-mini",
	})
	if err != nil {
		t.Fatalf("expected translation job creation success, got error: %v", err)
	}
	if job.Status != "queued" {
		t.Fatalf("expected queued status, got %s", job.Status)
	}
}

func TestCreateTranslationJobRejectsSameLocale(t *testing.T) {
	svc := NewContentService(&contentRepoStub{}, nil)
	_, err := svc.CreateTranslationJob(domain.TranslationJob{
		SourceType:   "article",
		SourceID:     "a1",
		SourceLocale: "en-US",
		TargetLocale: "en-US",
		ProviderKey:  "openai_compatible",
		ModelName:    "gpt-4.1-mini",
	})
	if err == nil {
		t.Fatalf("expected error when sourceLocale equals targetLocale")
	}
}
