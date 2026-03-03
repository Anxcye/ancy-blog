// File: content.go
// Purpose: Orchestrate content, site configuration, and timeline business operations.
// Module: backend/internal/service, content service layer.
// Related: content repository implementations, cache layer, and HTTP handlers.
package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/anxcye/ancy-blog/backend/internal/cache"
	"github.com/anxcye/ancy-blog/backend/internal/domain"
	"github.com/anxcye/ancy-blog/backend/internal/repository"
)

type ContentService struct {
	repo  repository.ContentRepository
	cache cache.Cache
}

func NewContentService(repo repository.ContentRepository, cacheClient cache.Cache) *ContentService {
	return &ContentService{repo: repo, cache: cacheClient}
}

const (
	cacheTTL             = 30 * time.Minute
	cacheSiteSettingsKey = "site:settings:default"
	cacheSiteFooterKey   = "site:footer:default"
	cacheSiteSocialKey   = "site:social:default"
	cacheSiteNavKey      = "site:nav:default"
	cacheSiteSlotKeyFmt  = "site:slot:%s:default"
	defaultJSONObj       = "{}"
)

func (s *ContentService) CreateArticle(article domain.Article) (domain.Article, error) {
	if strings.TrimSpace(article.Title) == "" || strings.TrimSpace(article.Slug) == "" {
		return domain.Article{}, errors.New("title and slug are required")
	}
	if article.ContentKind == "" {
		article.ContentKind = "post"
	}
	if article.Status == "" {
		article.Status = "draft"
	}
	if article.Visibility == "" {
		article.Visibility = "public"
	}
	if article.OriginType == "" {
		article.OriginType = "original"
	}
	if article.AIAssistLevel == "" {
		article.AIAssistLevel = "none"
	}
	return s.repo.CreateArticle(article)
}

func (s *ContentService) UpdateArticle(id string, article domain.Article) (domain.Article, error) {
	if strings.TrimSpace(id) == "" {
		return domain.Article{}, errors.New("article id is required")
	}
	if strings.TrimSpace(article.Title) == "" || strings.TrimSpace(article.Slug) == "" {
		return domain.Article{}, errors.New("title and slug are required")
	}
	if article.ContentKind == "" {
		article.ContentKind = "post"
	}
	return s.repo.UpdateArticle(id, article)
}

func (s *ContentService) ListPublishedArticles(page, pageSize int, category, tag, contentKind string) ([]domain.Article, int) {
	if contentKind == "" {
		contentKind = "post"
	}
	return s.repo.ListPublishedArticles(page, pageSize, category, tag, contentKind)
}

func (s *ContentService) GetPublishedArticleBySlug(slug string) (domain.Article, bool) {
	return s.repo.GetPublishedArticleBySlug(slug)
}

func (s *ContentService) CreateMoment(moment domain.Moment) (domain.Moment, error) {
	if strings.TrimSpace(moment.Content) == "" {
		return domain.Moment{}, errors.New("content is required")
	}
	if moment.Status == "" {
		moment.Status = "draft"
	}
	return s.repo.CreateMoment(moment)
}

func (s *ContentService) ListPublishedMoments(page, pageSize int) ([]domain.Moment, int) {
	return s.repo.ListPublishedMoments(page, pageSize)
}

func (s *ContentService) CreateComment(comment domain.Comment) (domain.Comment, error) {
	if strings.TrimSpace(comment.ArticleID) == "" {
		return domain.Comment{}, errors.New("articleId is required")
	}
	if strings.TrimSpace(comment.Content) == "" {
		return domain.Comment{}, errors.New("content is required")
	}
	if strings.TrimSpace(comment.Nickname) == "" {
		return domain.Comment{}, errors.New("nickname is required")
	}
	if comment.Status == "" {
		comment.Status = "approved"
	}
	return s.repo.CreateComment(comment)
}

func (s *ContentService) ListArticleComments(articleID string, page, pageSize int) ([]domain.Comment, int) {
	return s.repo.ListArticleComments(articleID, page, pageSize)
}

func (s *ContentService) ListCommentChildren(parentID string, page, pageSize int) ([]domain.Comment, int) {
	return s.repo.ListCommentChildren(parentID, page, pageSize)
}

func (s *ContentService) CountArticleComments(articleID string) (int, error) {
	return s.repo.CountArticleComments(articleID)
}

func (s *ContentService) ListCommentPage(page, pageSize int, status string) ([]domain.Comment, int) {
	return s.repo.ListCommentPage(page, pageSize, status)
}

func (s *ContentService) UpdateCommentAdmin(id, status, isPinned string) (domain.Comment, error) {
	if strings.TrimSpace(id) == "" {
		return domain.Comment{}, errors.New("id is required")
	}
	if status == "" {
		status = "approved"
	}
	return s.repo.UpdateCommentAdmin(id, status, isPinned)
}

func (s *ContentService) SubmitLink(link domain.Link) (domain.Link, error) {
	if strings.TrimSpace(link.Name) == "" || strings.TrimSpace(link.URL) == "" {
		return domain.Link{}, errors.New("name and url are required")
	}
	if _, err := url.ParseRequestURI(link.URL); err != nil {
		return domain.Link{}, errors.New("url is invalid")
	}
	return s.repo.SubmitLink(link)
}

func (s *ContentService) ListApprovedLinks() []domain.Link {
	return s.repo.ListApprovedLinks()
}

func (s *ContentService) ListLinkSubmissions(page, pageSize int, reviewStatus string) ([]domain.Link, int) {
	return s.repo.ListLinkSubmissions(page, pageSize, reviewStatus)
}

func (s *ContentService) ReviewLink(id, reviewStatus, reviewNote, relatedArticleID string) (domain.Link, error) {
	if reviewStatus != "pending" && reviewStatus != "approved" && reviewStatus != "rejected" {
		return domain.Link{}, errors.New("invalid review status")
	}
	if relatedArticleID != "" {
		if _, ok := s.repo.GetArticleByID(relatedArticleID); !ok {
			return domain.Link{}, errors.New("related article not found")
		}
	}
	return s.repo.ReviewLink(id, reviewStatus, reviewNote, relatedArticleID)
}

func (s *ContentService) ListCategories() []domain.Category {
	return s.repo.ListCategories()
}

func (s *ContentService) ListTags() []domain.Tag {
	return s.repo.ListTags()
}

func (s *ContentService) GetSiteSettings() domain.SiteSettings {
	if s.cache != nil {
		var cached domain.SiteSettings
		if ok := s.getCache(cacheSiteSettingsKey, &cached); ok {
			return cached
		}
	}
	settings := s.repo.GetSiteSettings()
	s.setCache(cacheSiteSettingsKey, settings)
	return settings
}

func (s *ContentService) UpdateSiteSettings(settings domain.SiteSettings) domain.SiteSettings {
	if settings.DefaultLocale == "" {
		settings.DefaultLocale = "en"
	}
	updated := s.repo.UpdateSiteSettings(settings)
	s.delCache(cacheSiteSettingsKey)
	return updated
}

func (s *ContentService) CreateFooterItem(item domain.FooterItem) (domain.FooterItem, error) {
	if strings.TrimSpace(item.Label) == "" {
		return domain.FooterItem{}, errors.New("label is required")
	}
	if item.RowNum < 1 || item.RowNum > 3 {
		return domain.FooterItem{}, errors.New("rowNum must be between 1 and 3")
	}
	if item.LinkType == "internal" {
		a, ok := s.repo.GetPublishedArticleBySlug(item.InternalArticleSlug)
		if !ok || a.ContentKind != "page" {
			return domain.FooterItem{}, errors.New("internal article slug must point to a published page")
		}
	}
	if !item.Enabled && item.OrderNum == 0 {
		item.OrderNum = 1
	}
	created, err := s.repo.CreateFooterItem(item)
	if err == nil {
		s.delCache(cacheSiteFooterKey)
	}
	return created, err
}

func (s *ContentService) UpdateFooterItem(id string, item domain.FooterItem) (domain.FooterItem, error) {
	if item.RowNum < 1 || item.RowNum > 3 {
		return domain.FooterItem{}, errors.New("rowNum must be between 1 and 3")
	}
	if item.LinkType == "internal" {
		a, ok := s.repo.GetPublishedArticleBySlug(item.InternalArticleSlug)
		if !ok || a.ContentKind != "page" {
			return domain.FooterItem{}, errors.New("internal article slug must point to a published page")
		}
	}
	updated, err := s.repo.UpdateFooterItem(id, item)
	if err == nil {
		s.delCache(cacheSiteFooterKey)
	}
	return updated, err
}

func (s *ContentService) DeleteFooterItem(id string) bool {
	ok := s.repo.DeleteFooterItem(id)
	if ok {
		s.delCache(cacheSiteFooterKey)
	}
	return ok
}

func (s *ContentService) ListFooterItems() []domain.FooterItem {
	if s.cache != nil {
		var cached []domain.FooterItem
		if ok := s.getCache(cacheSiteFooterKey, &cached); ok {
			return cached
		}
	}
	items := s.repo.ListFooterItems()
	s.setCache(cacheSiteFooterKey, items)
	return items
}

func (s *ContentService) CreateSocialLink(item domain.SocialLink) (domain.SocialLink, error) {
	if strings.TrimSpace(item.Title) == "" || strings.TrimSpace(item.URL) == "" {
		return domain.SocialLink{}, errors.New("title and url are required")
	}
	created, err := s.repo.CreateSocialLink(item)
	if err == nil {
		s.delCache(cacheSiteSocialKey)
	}
	return created, err
}

func (s *ContentService) UpdateSocialLink(id string, item domain.SocialLink) (domain.SocialLink, error) {
	if strings.TrimSpace(item.Title) == "" || strings.TrimSpace(item.URL) == "" {
		return domain.SocialLink{}, errors.New("title and url are required")
	}
	updated, err := s.repo.UpdateSocialLink(id, item)
	if err == nil {
		s.delCache(cacheSiteSocialKey)
	}
	return updated, err
}

func (s *ContentService) DeleteSocialLink(id string) bool {
	ok := s.repo.DeleteSocialLink(id)
	if ok {
		s.delCache(cacheSiteSocialKey)
	}
	return ok
}

func (s *ContentService) ListSocialLinks() []domain.SocialLink {
	if s.cache != nil {
		var cached []domain.SocialLink
		if ok := s.getCache(cacheSiteSocialKey, &cached); ok {
			return cached
		}
	}
	items := s.repo.ListSocialLinks()
	s.setCache(cacheSiteSocialKey, items)
	return items
}

func (s *ContentService) CreateNavItem(item domain.NavItem) (domain.NavItem, error) {
	if strings.TrimSpace(item.Name) == "" || strings.TrimSpace(item.Key) == "" {
		return domain.NavItem{}, errors.New("name and key are required")
	}
	created, err := s.repo.CreateNavItem(item)
	if err == nil {
		s.delCache(cacheSiteNavKey)
	}
	return created, err
}

func (s *ContentService) UpdateNavItem(id string, item domain.NavItem) (domain.NavItem, error) {
	if strings.TrimSpace(item.Name) == "" || strings.TrimSpace(item.Key) == "" {
		return domain.NavItem{}, errors.New("name and key are required")
	}
	updated, err := s.repo.UpdateNavItem(id, item)
	if err == nil {
		s.delCache(cacheSiteNavKey)
	}
	return updated, err
}

func (s *ContentService) DeleteNavItem(id string) bool {
	ok := s.repo.DeleteNavItem(id)
	if ok {
		s.delCache(cacheSiteNavKey)
	}
	return ok
}

func (s *ContentService) ListNavItems() []domain.NavItem {
	if s.cache != nil {
		var cached []domain.NavItem
		if ok := s.getCache(cacheSiteNavKey, &cached); ok {
			return cached
		}
	}
	items := s.repo.ListNavItems()
	s.setCache(cacheSiteNavKey, items)
	return items
}

func (s *ContentService) CreateContentSlot(slot domain.ContentSlot) (domain.ContentSlot, error) {
	if strings.TrimSpace(slot.SlotKey) == "" || strings.TrimSpace(slot.Name) == "" {
		return domain.ContentSlot{}, errors.New("slotKey and name are required")
	}
	created, err := s.repo.CreateContentSlot(slot)
	if err == nil {
		s.delCache(fmt.Sprintf(cacheSiteSlotKeyFmt, slot.SlotKey))
	}
	return created, err
}

func (s *ContentService) CreateSlotItem(slotKey string, item domain.SlotItem) (domain.SlotItem, error) {
	if item.ContentType != "article" && item.ContentType != "moment" {
		return domain.SlotItem{}, errors.New("contentType must be article or moment")
	}
	if item.ContentType == "article" {
		if _, ok := s.repo.GetArticleByID(item.ContentID); !ok {
			return domain.SlotItem{}, errors.New("article not found")
		}
	}
	created, err := s.repo.CreateSlotItem(slotKey, item)
	if err == nil {
		s.delCache(fmt.Sprintf(cacheSiteSlotKeyFmt, slotKey))
	}
	return created, err
}

func (s *ContentService) DeleteSlotItem(slotKey, itemID string) bool {
	ok := s.repo.DeleteSlotItem(slotKey, itemID)
	if ok {
		s.delCache(fmt.Sprintf(cacheSiteSlotKeyFmt, slotKey))
	}
	return ok
}

func (s *ContentService) ListSlotContent(slotKey string, limit int) ([]domain.SlotContentItem, bool) {
	cacheKey := fmt.Sprintf(cacheSiteSlotKeyFmt, slotKey)
	if s.cache != nil {
		var cached []domain.SlotContentItem
		if ok := s.getCache(cacheKey, &cached); ok {
			if limit > 0 && len(cached) > limit {
				return cached[:limit], true
			}
			return cached, true
		}
	}
	items, exists := s.repo.ListSlotContent(slotKey, limit)
	if exists {
		s.setCache(cacheKey, items)
	}
	return items, exists
}

func (s *ContentService) ListIntegrationProviders(providerType string) []domain.IntegrationProvider {
	rows := s.repo.ListIntegrationProviders(providerType)
	for i := range rows {
		rows[i].ConfigJSON = maskSecretJSON(rows[i].ConfigJSON)
	}
	return rows
}

func (s *ContentService) UpdateIntegrationProvider(providerKey string, enabled bool, configJSON, metaJSON []byte) (domain.IntegrationProvider, error) {
	if strings.TrimSpace(providerKey) == "" {
		return domain.IntegrationProvider{}, errors.New("provider key is required")
	}
	if !json.Valid(configJSON) {
		return domain.IntegrationProvider{}, errors.New("configJson must be valid JSON")
	}
	if len(metaJSON) > 0 && !json.Valid(metaJSON) {
		return domain.IntegrationProvider{}, errors.New("metaJson must be valid JSON")
	}
	if len(metaJSON) == 0 {
		metaJSON = []byte(defaultJSONObj)
	}
	provider, err := s.repo.UpdateIntegrationProvider(providerKey, enabled, configJSON, metaJSON)
	if err != nil {
		return domain.IntegrationProvider{}, err
	}
	provider.ConfigJSON = maskSecretJSON(provider.ConfigJSON)
	return provider, nil
}

func (s *ContentService) TestIntegrationProvider(providerKey string) (domain.ProviderTestResult, error) {
	start := time.Now()
	provider, ok := s.repo.GetIntegrationProvider(providerKey)
	if !ok {
		return domain.ProviderTestResult{}, errors.New("provider not found")
	}
	if !provider.Enabled {
		return domain.ProviderTestResult{}, errors.New("provider is disabled")
	}

	var config map[string]any
	if err := json.Unmarshal(provider.ConfigJSON, &config); err != nil {
		return domain.ProviderTestResult{}, errors.New("provider config is invalid JSON")
	}

	requiredByProvider := map[string][]string{
		"cloudflare_r2":     {"account_id", "bucket", "access_key_id", "secret_access_key", "public_base_url"},
		"openai_compatible": {"base_url", "api_key", "model"},
	}
	required := requiredByProvider[provider.ProviderKey]
	for _, key := range required {
		v, ok := config[key]
		if !ok || strings.TrimSpace(fmt.Sprintf("%v", v)) == "" {
			return domain.ProviderTestResult{}, fmt.Errorf("missing required config key: %s", key)
		}
	}

	return domain.ProviderTestResult{
		OK:        true,
		Message:   "configuration looks valid",
		LatencyMS: time.Since(start).Milliseconds(),
	}, nil
}

func (s *ContentService) CreateTranslationJob(job domain.TranslationJob) (domain.TranslationJob, error) {
	if job.SourceType != "article" && job.SourceType != "moment" {
		return domain.TranslationJob{}, errors.New("sourceType must be article or moment")
	}
	if strings.TrimSpace(job.SourceID) == "" {
		return domain.TranslationJob{}, errors.New("sourceId is required")
	}
	if strings.TrimSpace(job.SourceLocale) == "" || strings.TrimSpace(job.TargetLocale) == "" {
		return domain.TranslationJob{}, errors.New("sourceLocale and targetLocale are required")
	}
	if strings.EqualFold(job.SourceLocale, job.TargetLocale) {
		return domain.TranslationJob{}, errors.New("sourceLocale and targetLocale must be different")
	}
	if strings.TrimSpace(job.ProviderKey) == "" || strings.TrimSpace(job.ModelName) == "" {
		return domain.TranslationJob{}, errors.New("providerKey and modelName are required")
	}
	provider, ok := s.repo.GetIntegrationProvider(job.ProviderKey)
	if !ok {
		return domain.TranslationJob{}, errors.New("provider not found")
	}
	if provider.ProviderType != "llm" {
		return domain.TranslationJob{}, errors.New("provider is not llm")
	}
	if !provider.Enabled {
		return domain.TranslationJob{}, errors.New("provider is disabled")
	}
	job.Status = "queued"
	return s.repo.CreateTranslationJob(job)
}

func (s *ContentService) ListTranslationJobs(page, pageSize int, status, sourceType, sourceID string) ([]domain.TranslationJob, int) {
	return s.repo.ListTranslationJobs(page, pageSize, status, sourceType, sourceID)
}

func (s *ContentService) GetTranslationJobByID(id string) (domain.TranslationJob, bool) {
	return s.repo.GetTranslationJobByID(id)
}

func (s *ContentService) ListTimeline(page, pageSize int) ([]domain.TimelineItem, int) {
	return s.repo.ListTimeline(page, pageSize)
}

func maskSecretJSON(raw []byte) []byte {
	if len(raw) == 0 || !json.Valid(raw) {
		return raw
	}
	var payload map[string]any
	if err := json.Unmarshal(raw, &payload); err != nil {
		return raw
	}
	secretKeys := map[string]struct{}{
		"access_key_id":     {},
		"secret_access_key": {},
		"api_key":           {},
		"token":             {},
		"secret":            {},
	}
	for k, v := range payload {
		if _, ok := secretKeys[strings.ToLower(k)]; ok {
			if s, ok2 := v.(string); ok2 && strings.TrimSpace(s) != "" {
				payload[k] = "******"
			}
		}
	}
	masked, err := json.Marshal(payload)
	if err != nil {
		return raw
	}
	return masked
}

func (s *ContentService) getCache(key string, out any) bool {
	if s.cache == nil {
		return false
	}
	raw, ok, err := s.cache.Get(context.Background(), key)
	if err != nil || !ok {
		return false
	}
	if err := json.Unmarshal([]byte(raw), out); err != nil {
		return false
	}
	return true
}

func (s *ContentService) setCache(key string, value any) {
	if s.cache == nil {
		return
	}
	b, err := json.Marshal(value)
	if err != nil {
		return
	}
	_ = s.cache.Set(context.Background(), key, string(b), cacheTTL)
}

func (s *ContentService) delCache(keys ...string) {
	if s.cache == nil {
		return
	}
	_ = s.cache.Del(context.Background(), keys...)
}
