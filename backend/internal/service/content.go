// File: content.go
// Purpose: Orchestrate content, site configuration, and timeline business operations.
// Module: backend/internal/service, content service layer.
// Related: content repository implementations, cache layer, and HTTP handlers.
package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/anxcye/ancy-blog/backend/internal/apperr"
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
		return domain.Article{}, fmt.Errorf("%w: title and slug are required", apperr.ErrValidation)
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
	result, err := s.repo.CreateArticle(article)
	if err == nil && result.Status == "published" {
		go s.triggerAutoTranslation(result.ID, s.repo.GetTranslationPolicy())
	}
	return result, err
}

func (s *ContentService) UpdateArticle(id string, article domain.Article) (domain.Article, error) {
	if strings.TrimSpace(id) == "" {
		return domain.Article{}, fmt.Errorf("%w: article id is required", apperr.ErrValidation)
	}
	if strings.TrimSpace(article.Title) == "" || strings.TrimSpace(article.Slug) == "" {
		return domain.Article{}, fmt.Errorf("%w: title and slug are required", apperr.ErrValidation)
	}
	if article.ContentKind == "" {
		article.ContentKind = "post"
	}
	result, err := s.repo.UpdateArticle(id, article)
	if err == nil && result.Status == "published" {
		go s.triggerAutoTranslation(result.ID, s.repo.GetTranslationPolicy())
	}
	return result, err
}

func (s *ContentService) ListArticles(page, pageSize int, status, contentKind, keyword string) ([]domain.Article, int) {
	return s.repo.ListArticles(page, pageSize, status, contentKind, keyword)
}

func (s *ContentService) DeleteArticle(id string) bool {
	return s.repo.DeleteArticle(id)
}

func (s *ContentService) BatchUpdateArticleStatus(ids []string, status string) (int, error) {
	if len(ids) == 0 {
		return 0, fmt.Errorf("%w: ids are required", apperr.ErrValidation)
	}
	switch status {
	case "draft", "published", "scheduled":
	default:
		return 0, fmt.Errorf("%w: invalid status", apperr.ErrValidation)
	}
	return s.repo.BatchUpdateArticleStatus(ids, status), nil
}

func (s *ContentService) ListPublishedArticles(page, pageSize int, category, tag, contentKind string) ([]domain.Article, int) {
	if contentKind == "" {
		contentKind = "post"
	}
	return s.repo.ListPublishedArticles(page, pageSize, category, tag, contentKind)
}

func (s *ContentService) GetArticleByID(id string) (domain.Article, bool) {
	return s.repo.GetArticleByID(id)
}

func (s *ContentService) GetPublishedArticleBySlug(slug string) (domain.Article, bool) {
	return s.repo.GetPublishedArticleBySlug(slug)
}

func (s *ContentService) GetPublishedArticleBySlugWithLocale(slug, locale string) (domain.Article, bool) {
	if strings.TrimSpace(locale) == "" {
		return s.repo.GetPublishedArticleBySlug(slug)
	}
	return s.repo.GetPublishedArticleBySlugWithLocale(slug, locale)
}

func (s *ContentService) SlugExists(slug string) bool {
	return s.repo.SlugExists(slug)
}

func (s *ContentService) RecordView(articleID, visitorKey string) (int64, error) {
	return s.repo.RecordView(articleID, visitorKey)
}

func (s *ContentService) CreateMoment(moment domain.Moment) (domain.Moment, error) {
	if strings.TrimSpace(moment.Content) == "" {
		return domain.Moment{}, fmt.Errorf("%w: content is required", apperr.ErrValidation)
	}
	if moment.Status == "" {
		moment.Status = "draft"
	}
	return s.repo.CreateMoment(moment)
}

func (s *ContentService) UpdateMoment(id string, moment domain.Moment) (domain.Moment, error) {
	if strings.TrimSpace(id) == "" {
		return domain.Moment{}, fmt.Errorf("%w: moment id is required", apperr.ErrValidation)
	}
	if strings.TrimSpace(moment.Content) == "" {
		return domain.Moment{}, fmt.Errorf("%w: content is required", apperr.ErrValidation)
	}
	if moment.Status == "" {
		moment.Status = "draft"
	}
	return s.repo.UpdateMoment(id, moment)
}

func (s *ContentService) ListMoments(page, pageSize int, status string) ([]domain.Moment, int) {
	return s.repo.ListMoments(page, pageSize, status)
}

func (s *ContentService) GetMomentByID(id string) (domain.Moment, bool) {
	return s.repo.GetMomentByID(id)
}

func (s *ContentService) DeleteMoment(id string) bool {
	return s.repo.DeleteMoment(id)
}

func (s *ContentService) BatchUpdateMomentStatus(ids []string, status string) (int, error) {
	if len(ids) == 0 {
		return 0, fmt.Errorf("%w: ids are required", apperr.ErrValidation)
	}
	switch status {
	case "draft", "published", "scheduled":
	default:
		return 0, fmt.Errorf("%w: invalid status", apperr.ErrValidation)
	}
	return s.repo.BatchUpdateMomentStatus(ids, status), nil
}

func (s *ContentService) ListPublishedMoments(page, pageSize int, locale string) ([]domain.Moment, int) {
	return s.repo.ListPublishedMoments(page, pageSize, locale)
}

func (s *ContentService) GetPublishedMomentByID(id, locale string) (domain.Moment, bool) {
	return s.repo.GetPublishedMomentByID(id, locale)
}

func (s *ContentService) CreateComment(comment domain.Comment) (domain.Comment, error) {
	if strings.TrimSpace(comment.ContentType) == "" {
		if strings.TrimSpace(comment.ArticleID) != "" {
			comment.ContentType = "article"
		}
	}
	if strings.TrimSpace(comment.ContentID) == "" && strings.TrimSpace(comment.ArticleID) != "" {
		comment.ContentID = comment.ArticleID
	}
	if comment.ContentType == "article" && strings.TrimSpace(comment.ArticleID) == "" {
		comment.ArticleID = comment.ContentID
	}
	if comment.ContentType != "article" && comment.ContentType != "moment" {
		return domain.Comment{}, fmt.Errorf("%w: contentType is invalid", apperr.ErrValidation)
	}
	if strings.TrimSpace(comment.ContentID) == "" {
		return domain.Comment{}, fmt.Errorf("%w: contentId is required", apperr.ErrValidation)
	}
	if strings.TrimSpace(comment.Content) == "" {
		return domain.Comment{}, fmt.Errorf("%w: content is required", apperr.ErrValidation)
	}
	if strings.TrimSpace(comment.Nickname) == "" {
		return domain.Comment{}, fmt.Errorf("%w: nickname is required", apperr.ErrValidation)
	}
	if strings.TrimSpace(comment.ParentID) != "" && strings.TrimSpace(comment.RootID) == "" {
		comment.RootID = comment.ParentID
	}

	// Enforce site-level comment policy
	settings := s.repo.GetSiteSettings()
	if !settings.CommentEnabled {
		return domain.Comment{}, fmt.Errorf("%w: commenting is disabled", apperr.ErrValidation)
	}
	if comment.Status == "" {
		if settings.CommentRequireApproval {
			comment.Status = "pending"
		} else {
			comment.Status = "approved"
		}
	}
	return s.repo.CreateComment(comment)
}

func (s *ContentService) ListArticleComments(articleID string, page, pageSize int) ([]domain.Comment, int) {
	return s.repo.ListArticleComments(articleID, page, pageSize)
}

func (s *ContentService) ListContentComments(contentType, contentID string, page, pageSize int) ([]domain.Comment, int) {
	return s.repo.ListContentComments(contentType, contentID, page, pageSize)
}

func (s *ContentService) ListArticleCommentThreads(articleID string, page, pageSize int) ([]domain.CommentNode, int) {
	return s.ListContentCommentThreads("article", articleID, page, pageSize)
}

func (s *ContentService) ListContentCommentThreads(contentType, contentID string, page, pageSize int) ([]domain.CommentNode, int) {
	roots, total := s.repo.ListContentComments(contentType, contentID, page, pageSize)
	if len(roots) == 0 {
		return []domain.CommentNode{}, total
	}

	rootIDs := make([]string, 0, len(roots))
	for _, root := range roots {
		rootIDs = append(rootIDs, root.ID)
	}

	descendants := s.repo.ListCommentDescendants(rootIDs)
	rootSet := make(map[string]struct{}, len(roots))
	commentByID := make(map[string]domain.Comment, len(roots)+len(descendants))
	childrenByParent := make(map[string][]domain.Comment)
	childrenByRoot := make(map[string][]domain.Comment)
	normalizedDescendants := make([]domain.Comment, 0, len(descendants))

	for _, root := range roots {
		rootSet[root.ID] = struct{}{}
		commentByID[root.ID] = root
	}

	for _, comment := range descendants {
		parent, hasParent := commentByID[comment.ParentID]
		if hasParent {
			comment.ToCommentID = parent.ID
			comment.ToCommentNickname = parent.Nickname
			comment.ToCommentIsAuthor = strings.EqualFold(parent.Source, "admin") || strings.EqualFold(parent.Nickname, "admin")
		}
		commentByID[comment.ID] = comment
		normalizedDescendants = append(normalizedDescendants, comment)
	}

	for _, comment := range normalizedDescendants {
		parentID := comment.ParentID
		if parentID != "" {
			if _, ok := commentByID[parentID]; !ok {
				if _, rootOK := rootSet[comment.RootID]; rootOK {
					childrenByRoot[comment.RootID] = append(childrenByRoot[comment.RootID], comment)
				}
				continue
			}
			childrenByParent[parentID] = append(childrenByParent[parentID], comment)
			continue
		}
		if _, ok := rootSet[comment.RootID]; ok {
			childrenByRoot[comment.RootID] = append(childrenByRoot[comment.RootID], comment)
		}
	}

	for parentID := range childrenByParent {
		sortCommentsByThreadOrder(childrenByParent[parentID])
	}
	for rootID := range childrenByRoot {
		sortCommentsByThreadOrder(childrenByRoot[rootID])
	}

	items := make([]domain.CommentNode, 0, len(roots))
	for _, root := range roots {
		items = append(items, buildCommentNode(root, childrenByParent, childrenByRoot, map[string]bool{}))
	}

	return items, total
}

func (s *ContentService) ListCommentChildren(parentID string, page, pageSize int) ([]domain.Comment, int) {
	return s.repo.ListCommentChildren(parentID, page, pageSize)
}

func (s *ContentService) CountArticleComments(articleID string) (int, error) {
	return s.repo.CountArticleComments(articleID)
}

func (s *ContentService) CountContentComments(contentType, contentID string) (int, error) {
	return s.repo.CountContentComments(contentType, contentID)
}

func (s *ContentService) ListCommentPage(page, pageSize int, status string) ([]domain.Comment, int) {
	return s.repo.ListCommentPage(page, pageSize, status)
}

func (s *ContentService) ReplyToCommentAsAdmin(id, content, authorName, approvedBy, ip, userAgent string) (domain.Comment, error) {
	if strings.TrimSpace(id) == "" {
		return domain.Comment{}, fmt.Errorf("%w: id is required", apperr.ErrValidation)
	}
	if strings.TrimSpace(content) == "" {
		return domain.Comment{}, fmt.Errorf("%w: content is required", apperr.ErrValidation)
	}
	target, ok := s.repo.GetCommentByID(id)
	if !ok || target.Status == "deleted" {
		return domain.Comment{}, apperr.ErrCommentNotFound
	}
	if !s.repo.GetSiteSettings().CommentEnabled {
		return domain.Comment{}, fmt.Errorf("%w: commenting is disabled", apperr.ErrValidation)
	}
	if strings.TrimSpace(authorName) == "" {
		authorName = "Admin"
	}
	rootID := target.RootID
	if strings.TrimSpace(rootID) == "" {
		rootID = target.ID
	}

	reply := domain.Comment{
		ContentType: target.ContentType,
		ContentID:   target.ContentID,
		ParentID:    target.ID,
		RootID:      rootID,
		Content:     content,
		Status:      "approved",
		Nickname:    authorName,
		Source:      "admin",
		ApprovedBy:  approvedBy,
		IP:          ip,
		UserAgent:   userAgent,
	}
	if target.ContentType == "article" {
		reply.ArticleID = target.ContentID
	}
	return s.repo.CreateComment(reply)
}

func (s *ContentService) UpdateCommentAdmin(id, status, isPinned string) (domain.Comment, error) {
	if strings.TrimSpace(id) == "" {
		return domain.Comment{}, fmt.Errorf("%w: id is required", apperr.ErrValidation)
	}
	if status == "" {
		status = "approved"
	}
	return s.repo.UpdateCommentAdmin(id, status, isPinned)
}

func (s *ContentService) SubmitLink(link domain.Link) (domain.Link, error) {
	if strings.TrimSpace(link.Name) == "" || strings.TrimSpace(link.URL) == "" {
		return domain.Link{}, fmt.Errorf("%w: name and url are required", apperr.ErrValidation)
	}
	if !s.repo.GetSiteSettings().LinkSubmissionEnabled {
		return domain.Link{}, apperr.ErrLinkSubmissionDisabled
	}
	if _, err := url.ParseRequestURI(link.URL); err != nil {
		return domain.Link{}, fmt.Errorf("%w: url is invalid", apperr.ErrValidation)
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
		return domain.Link{}, fmt.Errorf("%w: invalid review status", apperr.ErrValidation)
	}
	if relatedArticleID != "" {
		if _, ok := s.repo.GetArticleByID(relatedArticleID); !ok {
			return domain.Link{}, apperr.ErrArticleNotFound
		}
	}
	return s.repo.ReviewLink(id, reviewStatus, reviewNote, relatedArticleID)
}

func (s *ContentService) ListCategories() []domain.Category {
	return s.repo.ListCategories()
}

func (s *ContentService) CreateCategory(category domain.Category) (domain.Category, error) {
	if strings.TrimSpace(category.Name) == "" || strings.TrimSpace(category.Slug) == "" {
		return domain.Category{}, fmt.Errorf("%w: name and slug are required", apperr.ErrValidation)
	}
	return s.repo.CreateCategory(category)
}

func (s *ContentService) DeleteCategory(id string) bool {
	return s.repo.DeleteCategory(id)
}

func (s *ContentService) ListTags() []domain.Tag {
	return s.repo.ListTags()
}

func (s *ContentService) CreateTag(tag domain.Tag) (domain.Tag, error) {
	if strings.TrimSpace(tag.Name) == "" || strings.TrimSpace(tag.Slug) == "" {
		return domain.Tag{}, fmt.Errorf("%w: name and slug are required", apperr.ErrValidation)
	}
	return s.repo.CreateTag(tag)
}

func (s *ContentService) DeleteTag(id string) bool {
	return s.repo.DeleteTag(id)
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
		return domain.FooterItem{}, fmt.Errorf("%w: label is required", apperr.ErrValidation)
	}
	if item.RowNum < 1 || item.RowNum > 3 {
		return domain.FooterItem{}, fmt.Errorf("%w: rowNum must be between 1 and 3", apperr.ErrValidation)
	}
	if item.LinkType == "internal" {
		a, ok := s.repo.GetPublishedArticleBySlug(item.InternalArticleSlug)
		if !ok || a.ContentKind != "page" {
			return domain.FooterItem{}, fmt.Errorf("%w: internal article slug must point to a published page", apperr.ErrValidation)
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
		return domain.FooterItem{}, fmt.Errorf("%w: rowNum must be between 1 and 3", apperr.ErrValidation)
	}
	if item.LinkType == "internal" {
		a, ok := s.repo.GetPublishedArticleBySlug(item.InternalArticleSlug)
		if !ok || a.ContentKind != "page" {
			return domain.FooterItem{}, fmt.Errorf("%w: internal article slug must point to a published page", apperr.ErrValidation)
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
		return domain.SocialLink{}, fmt.Errorf("%w: title and url are required", apperr.ErrValidation)
	}
	created, err := s.repo.CreateSocialLink(item)
	if err == nil {
		s.delCache(cacheSiteSocialKey)
	}
	return created, err
}

func (s *ContentService) UpdateSocialLink(id string, item domain.SocialLink) (domain.SocialLink, error) {
	if strings.TrimSpace(item.Title) == "" || strings.TrimSpace(item.URL) == "" {
		return domain.SocialLink{}, fmt.Errorf("%w: title and url are required", apperr.ErrValidation)
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
		return domain.NavItem{}, fmt.Errorf("%w: name and key are required", apperr.ErrValidation)
	}
	created, err := s.repo.CreateNavItem(item)
	if err == nil {
		s.delCache(cacheSiteNavKey)
	}
	return created, err
}

func (s *ContentService) UpdateNavItem(id string, item domain.NavItem) (domain.NavItem, error) {
	if strings.TrimSpace(item.Name) == "" || strings.TrimSpace(item.Key) == "" {
		return domain.NavItem{}, fmt.Errorf("%w: name and key are required", apperr.ErrValidation)
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
		return domain.ContentSlot{}, fmt.Errorf("%w: slotKey and name are required", apperr.ErrValidation)
	}
	created, err := s.repo.CreateContentSlot(slot)
	if err == nil {
		s.delCache(fmt.Sprintf(cacheSiteSlotKeyFmt, slot.SlotKey))
	}
	return created, err
}

func (s *ContentService) ListContentSlots() []domain.ContentSlot {
	return s.repo.ListContentSlots()
}

func (s *ContentService) CreateSlotItem(slotKey string, item domain.SlotItem) (domain.SlotItem, error) {
	if item.ContentType != "article" && item.ContentType != "moment" {
		return domain.SlotItem{}, fmt.Errorf("%w: contentType must be article or moment", apperr.ErrValidation)
	}
	if item.ContentType == "article" {
		if _, ok := s.repo.GetArticleByID(item.ContentID); !ok {
			return domain.SlotItem{}, apperr.ErrArticleNotFound
		}
	}
	created, err := s.repo.CreateSlotItem(slotKey, item)
	if err == nil {
		s.delCache(fmt.Sprintf(cacheSiteSlotKeyFmt, slotKey))
	}
	return created, err
}

func (s *ContentService) ListSlotItems(slotKey string) ([]domain.SlotItem, bool) {
	return s.repo.ListSlotItems(slotKey)
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
		return domain.IntegrationProvider{}, fmt.Errorf("%w: provider key is required", apperr.ErrValidation)
	}
	if !json.Valid(configJSON) {
		return domain.IntegrationProvider{}, fmt.Errorf("%w: configJson must be valid JSON", apperr.ErrValidation)
	}
	if len(metaJSON) > 0 && !json.Valid(metaJSON) {
		return domain.IntegrationProvider{}, fmt.Errorf("%w: metaJson must be valid JSON", apperr.ErrValidation)
	}
	if len(metaJSON) == 0 {
		metaJSON = []byte(defaultJSONObj)
	}
	if existingProvider, exists := s.repo.GetIntegrationProvider(providerKey); exists {
		configJSON = preserveMaskedSecrets(configJSON, existingProvider.ConfigJSON)
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
		return domain.ProviderTestResult{}, apperr.ErrProviderNotFound
	}
	if !provider.Enabled {
		return domain.ProviderTestResult{}, fmt.Errorf("%w: provider is disabled", apperr.ErrValidation)
	}

	var config map[string]any
	if err := json.Unmarshal(provider.ConfigJSON, &config); err != nil {
		return domain.ProviderTestResult{}, fmt.Errorf("%w: provider config is invalid JSON", apperr.ErrValidation)
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
		return domain.TranslationJob{}, fmt.Errorf("%w: sourceType must be article or moment", apperr.ErrValidation)
	}
	if strings.TrimSpace(job.SourceID) == "" {
		return domain.TranslationJob{}, fmt.Errorf("%w: sourceId is required", apperr.ErrValidation)
	}
	if strings.TrimSpace(job.SourceLocale) == "" || strings.TrimSpace(job.TargetLocale) == "" {
		return domain.TranslationJob{}, fmt.Errorf("%w: sourceLocale and targetLocale are required", apperr.ErrValidation)
	}
	if strings.EqualFold(job.SourceLocale, job.TargetLocale) {
		return domain.TranslationJob{}, fmt.Errorf("%w: sourceLocale and targetLocale must be different", apperr.ErrValidation)
	}
	if strings.TrimSpace(job.ProviderKey) == "" || strings.TrimSpace(job.ModelName) == "" {
		return domain.TranslationJob{}, fmt.Errorf("%w: providerKey and modelName are required", apperr.ErrValidation)
	}
	provider, ok := s.repo.GetIntegrationProvider(job.ProviderKey)
	if !ok {
		return domain.TranslationJob{}, apperr.ErrProviderNotFound
	}
	if provider.ProviderType != "llm" {
		return domain.TranslationJob{}, fmt.Errorf("%w: provider is not llm", apperr.ErrValidation)
	}
	if !provider.Enabled {
		return domain.TranslationJob{}, fmt.Errorf("%w: provider is disabled", apperr.ErrValidation)
	}
	if job.MaxRetries <= 0 {
		job.MaxRetries = 3
	}
	if job.MaxRetries > 10 {
		return domain.TranslationJob{}, fmt.Errorf("%w: maxRetries must be <= 10", apperr.ErrValidation)
	}
	job.RetryCount = 0
	if job.NextRetryAt.IsZero() {
		job.NextRetryAt = time.Now().UTC()
	}
	if !job.PublishAt.IsZero() {
		job.PublishAt = job.PublishAt.UTC()
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

func (s *ContentService) ClaimNextQueuedTranslationJob() (domain.TranslationJob, bool, error) {
	return s.repo.ClaimNextQueuedTranslationJob()
}

func (s *ContentService) MarkTranslationJobSucceeded(id, resultText string) error {
	return s.repo.MarkTranslationJobSucceeded(id, resultText)
}

func (s *ContentService) MarkTranslationJobFailed(id, errorMessage string) error {
	return s.repo.MarkTranslationJobFailed(id, errorMessage)
}

func (s *ContentService) ScheduleTranslationJobRetry(id, errorMessage string, nextRetryAt time.Time) error {
	if strings.TrimSpace(id) == "" {
		return fmt.Errorf("%w: id is required", apperr.ErrValidation)
	}
	if nextRetryAt.IsZero() {
		nextRetryAt = time.Now().UTC()
	}
	return s.repo.ScheduleTranslationJobRetry(id, errorMessage, nextRetryAt)
}

func (s *ContentService) RetryTranslationJob(id string) (domain.TranslationJob, error) {
	if strings.TrimSpace(id) == "" {
		return domain.TranslationJob{}, fmt.Errorf("%w: id is required", apperr.ErrValidation)
	}
	return s.repo.RetryTranslationJob(id)
}

func (s *ContentService) GetTranslationSourceText(sourceType, sourceID string) (string, bool, error) {
	return s.repo.GetTranslationSourceText(sourceType, sourceID)
}

func (s *ContentService) UpsertTranslationResult(sourceType, sourceID, targetLocale, title, summary, content, status string, publishedAt time.Time, translatedByJobID string) error {
	if status == "" {
		status = "draft"
	}
	switch sourceType {
	case "article":
		return s.repo.UpsertArticleTranslation(sourceID, targetLocale, title, summary, content, status, publishedAt, translatedByJobID)
	case "moment":
		return s.repo.UpsertMomentTranslation(sourceID, targetLocale, content, status, publishedAt, translatedByJobID)
	default:
		return fmt.Errorf("%w: unsupported sourceType", apperr.ErrValidation)
	}
}

func (s *ContentService) ListTranslationContents(page, pageSize int, sourceType, sourceID, locale string) ([]domain.TranslationContent, int, error) {
	if sourceType != "article" && sourceType != "moment" {
		return nil, 0, fmt.Errorf("%w: sourceType must be article or moment", apperr.ErrValidation)
	}
	rows, total := s.repo.ListTranslationContents(page, pageSize, sourceType, sourceID, locale)
	return rows, total, nil
}

func (s *ContentService) GetTranslationContent(sourceType, sourceID, locale string) (domain.TranslationContent, bool, error) {
	if sourceType != "article" && sourceType != "moment" {
		return domain.TranslationContent{}, false, fmt.Errorf("%w: sourceType must be article or moment", apperr.ErrValidation)
	}
	if strings.TrimSpace(sourceID) == "" || strings.TrimSpace(locale) == "" {
		return domain.TranslationContent{}, false, fmt.Errorf("%w: sourceId and locale are required", apperr.ErrValidation)
	}
	row, ok := s.repo.GetTranslationContent(sourceType, sourceID, locale)
	return row, ok, nil
}

func (s *ContentService) UpsertTranslationContent(sourceType, sourceID, locale, title, summary, content, status string, publishedAt time.Time, translatedByJobID string) (domain.TranslationContent, error) {
	if sourceType != "article" && sourceType != "moment" {
		return domain.TranslationContent{}, fmt.Errorf("%w: sourceType must be article or moment", apperr.ErrValidation)
	}
	if strings.TrimSpace(sourceID) == "" || strings.TrimSpace(locale) == "" {
		return domain.TranslationContent{}, fmt.Errorf("%w: sourceId and locale are required", apperr.ErrValidation)
	}
	if strings.TrimSpace(content) == "" {
		return domain.TranslationContent{}, fmt.Errorf("%w: content is required", apperr.ErrValidation)
	}
	if status == "" {
		status = "draft"
	}
	if status != "draft" && status != "published" && status != "archived" {
		return domain.TranslationContent{}, fmt.Errorf("%w: status must be draft, published, or archived", apperr.ErrValidation)
	}
	if status == "published" && publishedAt.IsZero() {
		publishedAt = time.Now().UTC()
	}
	if !publishedAt.IsZero() {
		publishedAt = publishedAt.UTC()
	}
	return s.repo.UpsertTranslationContent(sourceType, sourceID, locale, title, summary, content, status, publishedAt, translatedByJobID)
}

func (s *ContentService) GetTranslationPolicy() domain.TranslationPolicy {
	return s.repo.GetTranslationPolicy()
}

func (s *ContentService) UpdateTranslationPolicy(policy domain.TranslationPolicy) error {
	return s.repo.UpdateTranslationPolicy(policy)
}

// triggerAutoTranslation enqueues translation jobs for all locales in the policy when an
// article is published (or re-published after an update). It runs best-effort and never
// blocks the caller.
func (s *ContentService) triggerAutoTranslation(articleID string, policy domain.TranslationPolicy) {
	if !policy.Enabled || len(policy.TargetLocales) == 0 || policy.ProviderKey == "" {
		return
	}
	sourceLocale := s.repo.GetSiteSettings().DefaultLocale
	if sourceLocale == "" {
		sourceLocale = "en"
	}
	for _, locale := range policy.TargetLocales {
		if locale == "" || locale == sourceLocale {
			continue
		}
		job := domain.TranslationJob{
			SourceType:   "article",
			SourceID:     articleID,
			SourceLocale: sourceLocale,
			TargetLocale: locale,
			ProviderKey:  policy.ProviderKey,
			Status:       "queued",
			AutoPublish:  policy.AutoPublish,
			MaxRetries:   3,
		}
		_, _ = s.repo.CreateTranslationJob(job)
	}
}

func (s *ContentService) GetIntegrationProviderForRuntime(providerKey string) (domain.IntegrationProvider, bool) {
	return s.repo.GetIntegrationProvider(providerKey)
}

func (s *ContentService) ListTimeline(page, pageSize int, locale string) ([]domain.TimelineItem, int) {
	return s.repo.ListTimeline(page, pageSize, locale)
}

func (s *ContentService) RecordVisitEvents(events []domain.VisitEvent) (domain.AnalyticsIngestResult, error) {
	if len(events) == 0 {
		return domain.AnalyticsIngestResult{}, fmt.Errorf("%w: events are required", apperr.ErrValidation)
	}
	normalized := make([]domain.VisitEvent, 0, len(events))
	now := time.Now().UTC()
	for _, event := range events {
		if strings.TrimSpace(event.EventID) == "" || strings.TrimSpace(event.VisitorID) == "" || strings.TrimSpace(event.SessionID) == "" {
			return domain.AnalyticsIngestResult{}, fmt.Errorf("%w: eventId, visitorId, and sessionId are required", apperr.ErrValidation)
		}
		switch event.EventType {
		case "page_view", "page_ping":
		default:
			return domain.AnalyticsIngestResult{}, fmt.Errorf("%w: invalid analytics event type", apperr.ErrValidation)
		}
		event.Path = normalizeAnalyticsPath(event.Path)
		if event.Path == "" {
			return domain.AnalyticsIngestResult{}, fmt.Errorf("%w: path is required", apperr.ErrValidation)
		}
		if event.OccurredAt.IsZero() {
			event.OccurredAt = now
		}
		event.ReferrerHost = analyticsReferrerHost(event.Referrer)
		if event.DeviceType == "" {
			event.DeviceType = "unknown"
		}
		event.OccurredAt = event.OccurredAt.UTC()
		normalized = append(normalized, event)
	}
	return s.repo.CreateVisitEvents(normalized)
}

func (s *ContentService) GetAnalyticsOverview(days int) (domain.AnalyticsOverview, error) {
	return s.repo.GetAnalyticsOverview(normalizeAnalyticsDays(days))
}

func (s *ContentService) ListAnalyticsPages(page, pageSize, days int, path, contentType string) ([]domain.AnalyticsPathStat, int, error) {
	return s.repo.ListAnalyticsPages(page, pageSize, normalizeAnalyticsDays(days), strings.TrimSpace(path), strings.TrimSpace(contentType))
}

func (s *ContentService) ListAnalyticsVisits(page, pageSize, days int, path, eventType, visitorID, sessionID, contentType, ip, deviceType, browserName, osName, isBot string) ([]domain.VisitEvent, int, error) {
	return s.repo.ListAnalyticsVisits(
		page,
		pageSize,
		normalizeAnalyticsDays(days),
		strings.TrimSpace(path),
		strings.TrimSpace(eventType),
		strings.TrimSpace(visitorID),
		strings.TrimSpace(sessionID),
		strings.TrimSpace(contentType),
		strings.TrimSpace(ip),
		strings.TrimSpace(deviceType),
		strings.TrimSpace(browserName),
		strings.TrimSpace(osName),
		strings.TrimSpace(isBot),
	)
}

func normalizeAnalyticsDays(days int) int {
	if days <= 0 {
		return 7
	}
	if days > 90 {
		return 90
	}
	return days
}

func normalizeAnalyticsPath(raw string) string {
	value := strings.TrimSpace(raw)
	if value == "" {
		return ""
	}
	parsed, err := url.Parse(value)
	if err == nil {
		if parsed.Path != "" {
			value = parsed.Path
		}
	}
	if !strings.HasPrefix(value, "/") {
		value = "/" + value
	}
	return value
}

func analyticsReferrerHost(raw string) string {
	value := strings.TrimSpace(raw)
	if value == "" {
		return ""
	}
	parsed, err := url.Parse(value)
	if err != nil {
		return ""
	}
	return strings.ToLower(parsed.Hostname())
}

func buildCommentNode(
	comment domain.Comment,
	childrenByParent map[string][]domain.Comment,
	childrenByRoot map[string][]domain.Comment,
	visited map[string]bool,
) domain.CommentNode {
	if visited[comment.ID] {
		return domain.CommentNode{Comment: comment}
	}
	visited[comment.ID] = true

	children := childrenByParent[comment.ID]
	if comment.ParentID == "" && len(children) == 0 {
		children = childrenByRoot[comment.ID]
	}

	node := domain.CommentNode{Comment: comment}
	for _, child := range children {
		node.Children = append(node.Children, buildCommentNode(child, childrenByParent, childrenByRoot, visited))
	}
	return node
}

func sortCommentsByThreadOrder(items []domain.Comment) {
	sort.SliceStable(items, func(i, j int) bool {
		leftPinned := items[i].IsPinned == "1"
		rightPinned := items[j].IsPinned == "1"
		if leftPinned != rightPinned {
			return leftPinned
		}
		return items[i].CreatedAt.Before(items[j].CreatedAt)
	})
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

func preserveMaskedSecrets(newRaw, oldRaw []byte) []byte {
	if len(newRaw) == 0 || len(oldRaw) == 0 || !json.Valid(newRaw) || !json.Valid(oldRaw) {
		return newRaw
	}
	var newPayload map[string]any
	if err := json.Unmarshal(newRaw, &newPayload); err != nil {
		return newRaw
	}
	var oldPayload map[string]any
	if err := json.Unmarshal(oldRaw, &oldPayload); err != nil {
		return newRaw
	}
	secretKeys := map[string]struct{}{
		"access_key_id":     {},
		"secret_access_key": {},
		"api_key":           {},
		"token":             {},
		"secret":            {},
	}
	for key, value := range newPayload {
		if _, isSecret := secretKeys[strings.ToLower(key)]; !isSecret {
			continue
		}
		if strings.TrimSpace(fmt.Sprintf("%v", value)) != "******" {
			continue
		}
		if oldValue, ok := oldPayload[key]; ok {
			newPayload[key] = oldValue
		}
	}
	merged, err := json.Marshal(newPayload)
	if err != nil {
		return newRaw
	}
	return merged
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
