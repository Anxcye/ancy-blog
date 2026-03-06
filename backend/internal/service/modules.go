// File: modules.go
// Purpose: Expose module-oriented service facades to decouple handlers from monolithic service dependencies.
// Module: backend/internal/service, service composition layer.
// Related: content.go core logic and handler constructors.
package service

import (
	"time"

	"github.com/anxcye/ancy-blog/backend/internal/domain"
)

type ArticleService struct{ core *ContentService }

type CommentService struct{ core *ContentService }

type LinkService struct{ core *ContentService }

type SiteService struct{ core *ContentService }

type IntegrationService struct{ core *ContentService }

type TranslationService struct{ core *ContentService }

type TimelineService struct{ core *ContentService }

func NewArticleService(core *ContentService) *ArticleService { return &ArticleService{core: core} }
func NewCommentService(core *ContentService) *CommentService { return &CommentService{core: core} }
func NewLinkService(core *ContentService) *LinkService       { return &LinkService{core: core} }
func NewSiteService(core *ContentService) *SiteService       { return &SiteService{core: core} }
func NewIntegrationService(core *ContentService) *IntegrationService {
	return &IntegrationService{core: core}
}
func NewTranslationService(core *ContentService) *TranslationService {
	return &TranslationService{core: core}
}
func NewTimelineService(core *ContentService) *TimelineService { return &TimelineService{core: core} }

func (s *ArticleService) CreateArticle(article domain.Article) (domain.Article, error) {
	return s.core.CreateArticle(article)
}
func (s *ArticleService) UpdateArticle(id string, article domain.Article) (domain.Article, error) {
	return s.core.UpdateArticle(id, article)
}
func (s *ArticleService) ListArticles(page, pageSize int, status, contentKind, keyword string) ([]domain.Article, int) {
	return s.core.ListArticles(page, pageSize, status, contentKind, keyword)
}
func (s *ArticleService) DeleteArticle(id string) bool { return s.core.DeleteArticle(id) }
func (s *ArticleService) BatchUpdateArticleStatus(ids []string, status string) (int, error) {
	return s.core.BatchUpdateArticleStatus(ids, status)
}
func (s *ArticleService) ListPublishedArticles(page, pageSize int, category, tag, contentKind string) ([]domain.Article, int) {
	return s.core.ListPublishedArticles(page, pageSize, category, tag, contentKind)
}
func (s *ArticleService) GetArticleByID(id string) (domain.Article, bool) {
	return s.core.GetArticleByID(id)
}
func (s *ArticleService) GetPublishedArticleBySlug(slug string) (domain.Article, bool) {
	return s.core.GetPublishedArticleBySlug(slug)
}
func (s *ArticleService) GetPublishedArticleBySlugWithLocale(slug, locale string) (domain.Article, bool) {
	return s.core.GetPublishedArticleBySlugWithLocale(slug, locale)
}
func (s *ArticleService) SlugExists(slug string) bool { return s.core.SlugExists(slug) }
func (s *ArticleService) CreateMoment(moment domain.Moment) (domain.Moment, error) {
	return s.core.CreateMoment(moment)
}
func (s *ArticleService) UpdateMoment(id string, moment domain.Moment) (domain.Moment, error) {
	return s.core.UpdateMoment(id, moment)
}
func (s *ArticleService) ListMoments(page, pageSize int, status string) ([]domain.Moment, int) {
	return s.core.ListMoments(page, pageSize, status)
}
func (s *ArticleService) DeleteMoment(id string) bool { return s.core.DeleteMoment(id) }
func (s *ArticleService) BatchUpdateMomentStatus(ids []string, status string) (int, error) {
	return s.core.BatchUpdateMomentStatus(ids, status)
}
func (s *ArticleService) ListPublishedMoments(page, pageSize int, locale string) ([]domain.Moment, int) {
	return s.core.ListPublishedMoments(page, pageSize, locale)
}
func (s *ArticleService) ListCategories() []domain.Category { return s.core.ListCategories() }
func (s *ArticleService) CreateCategory(category domain.Category) (domain.Category, error) {
	return s.core.CreateCategory(category)
}
func (s *ArticleService) DeleteCategory(id string) bool { return s.core.DeleteCategory(id) }
func (s *ArticleService) ListTags() []domain.Tag        { return s.core.ListTags() }
func (s *ArticleService) CreateTag(tag domain.Tag) (domain.Tag, error) {
	return s.core.CreateTag(tag)
}
func (s *ArticleService) DeleteTag(id string) bool { return s.core.DeleteTag(id) }
func (s *ArticleService) RecordView(articleID, visitorKey string) (int64, error) {
	return s.core.RecordView(articleID, visitorKey)
}

func (s *CommentService) CreateComment(comment domain.Comment) (domain.Comment, error) {
	return s.core.CreateComment(comment)
}
func (s *CommentService) ListArticleComments(articleID string, page, pageSize int) ([]domain.Comment, int) {
	return s.core.ListArticleComments(articleID, page, pageSize)
}
func (s *CommentService) ListCommentChildren(parentID string, page, pageSize int) ([]domain.Comment, int) {
	return s.core.ListCommentChildren(parentID, page, pageSize)
}
func (s *CommentService) CountArticleComments(articleID string) (int, error) {
	return s.core.CountArticleComments(articleID)
}
func (s *CommentService) ListCommentPage(page, pageSize int, status string) ([]domain.Comment, int) {
	return s.core.ListCommentPage(page, pageSize, status)
}
func (s *CommentService) UpdateCommentAdmin(id, status, isPinned string) (domain.Comment, error) {
	return s.core.UpdateCommentAdmin(id, status, isPinned)
}

func (s *LinkService) SubmitLink(link domain.Link) (domain.Link, error) {
	return s.core.SubmitLink(link)
}
func (s *LinkService) ListApprovedLinks() []domain.Link { return s.core.ListApprovedLinks() }
func (s *LinkService) ListLinkSubmissions(page, pageSize int, reviewStatus string) ([]domain.Link, int) {
	return s.core.ListLinkSubmissions(page, pageSize, reviewStatus)
}
func (s *LinkService) ReviewLink(id, reviewStatus, reviewNote, relatedArticleID string) (domain.Link, error) {
	return s.core.ReviewLink(id, reviewStatus, reviewNote, relatedArticleID)
}

func (s *SiteService) GetSiteSettings() domain.SiteSettings { return s.core.GetSiteSettings() }
func (s *SiteService) UpdateSiteSettings(settings domain.SiteSettings) domain.SiteSettings {
	return s.core.UpdateSiteSettings(settings)
}
func (s *SiteService) CreateFooterItem(item domain.FooterItem) (domain.FooterItem, error) {
	return s.core.CreateFooterItem(item)
}
func (s *SiteService) UpdateFooterItem(id string, item domain.FooterItem) (domain.FooterItem, error) {
	return s.core.UpdateFooterItem(id, item)
}
func (s *SiteService) DeleteFooterItem(id string) bool { return s.core.DeleteFooterItem(id) }
func (s *SiteService) ListFooterItems() []domain.FooterItem {
	return s.core.ListFooterItems()
}
func (s *SiteService) CreateSocialLink(item domain.SocialLink) (domain.SocialLink, error) {
	return s.core.CreateSocialLink(item)
}
func (s *SiteService) UpdateSocialLink(id string, item domain.SocialLink) (domain.SocialLink, error) {
	return s.core.UpdateSocialLink(id, item)
}
func (s *SiteService) DeleteSocialLink(id string) bool { return s.core.DeleteSocialLink(id) }
func (s *SiteService) ListSocialLinks() []domain.SocialLink {
	return s.core.ListSocialLinks()
}
func (s *SiteService) CreateNavItem(item domain.NavItem) (domain.NavItem, error) {
	return s.core.CreateNavItem(item)
}
func (s *SiteService) UpdateNavItem(id string, item domain.NavItem) (domain.NavItem, error) {
	return s.core.UpdateNavItem(id, item)
}
func (s *SiteService) DeleteNavItem(id string) bool { return s.core.DeleteNavItem(id) }
func (s *SiteService) ListNavItems() []domain.NavItem {
	return s.core.ListNavItems()
}
func (s *SiteService) CreateContentSlot(slot domain.ContentSlot) (domain.ContentSlot, error) {
	return s.core.CreateContentSlot(slot)
}
func (s *SiteService) ListContentSlots() []domain.ContentSlot { return s.core.ListContentSlots() }
func (s *SiteService) CreateSlotItem(slotKey string, item domain.SlotItem) (domain.SlotItem, error) {
	return s.core.CreateSlotItem(slotKey, item)
}
func (s *SiteService) ListSlotItems(slotKey string) ([]domain.SlotItem, bool) {
	return s.core.ListSlotItems(slotKey)
}
func (s *SiteService) DeleteSlotItem(slotKey, itemID string) bool {
	return s.core.DeleteSlotItem(slotKey, itemID)
}
func (s *SiteService) ListSlotContent(slotKey string, limit int) ([]domain.SlotContentItem, bool) {
	return s.core.ListSlotContent(slotKey, limit)
}
func (s *SiteService) GetTranslationPolicy() domain.TranslationPolicy {
	return s.core.GetTranslationPolicy()
}
func (s *SiteService) UpdateTranslationPolicy(policy domain.TranslationPolicy) error {
	return s.core.UpdateTranslationPolicy(policy)
}

func (s *IntegrationService) ListIntegrationProviders(providerType string) []domain.IntegrationProvider {
	return s.core.ListIntegrationProviders(providerType)
}
func (s *IntegrationService) GetIntegrationProviderForRuntime(providerKey string) (domain.IntegrationProvider, bool) {
	return s.core.GetIntegrationProviderForRuntime(providerKey)
}
func (s *IntegrationService) UpdateIntegrationProvider(providerKey string, enabled bool, configJSON, metaJSON []byte) (domain.IntegrationProvider, error) {
	return s.core.UpdateIntegrationProvider(providerKey, enabled, configJSON, metaJSON)
}
func (s *IntegrationService) TestIntegrationProvider(providerKey string) (domain.ProviderTestResult, error) {
	return s.core.TestIntegrationProvider(providerKey)
}

func (s *TranslationService) CreateTranslationJob(job domain.TranslationJob) (domain.TranslationJob, error) {
	return s.core.CreateTranslationJob(job)
}
func (s *TranslationService) ListTranslationJobs(page, pageSize int, status, sourceType, sourceID string) ([]domain.TranslationJob, int) {
	return s.core.ListTranslationJobs(page, pageSize, status, sourceType, sourceID)
}
func (s *TranslationService) GetTranslationJobByID(id string) (domain.TranslationJob, bool) {
	return s.core.GetTranslationJobByID(id)
}
func (s *TranslationService) ClaimNextQueuedTranslationJob() (domain.TranslationJob, bool, error) {
	return s.core.ClaimNextQueuedTranslationJob()
}
func (s *TranslationService) MarkTranslationJobSucceeded(id, resultText string) error {
	return s.core.MarkTranslationJobSucceeded(id, resultText)
}
func (s *TranslationService) MarkTranslationJobFailed(id, errorMessage string) error {
	return s.core.MarkTranslationJobFailed(id, errorMessage)
}
func (s *TranslationService) ScheduleTranslationJobRetry(id, errorMessage string, nextRetryAt time.Time) error {
	return s.core.ScheduleTranslationJobRetry(id, errorMessage, nextRetryAt)
}
func (s *TranslationService) RetryTranslationJob(id string) (domain.TranslationJob, error) {
	return s.core.RetryTranslationJob(id)
}
func (s *TranslationService) GetTranslationSourceText(sourceType, sourceID string) (string, bool, error) {
	return s.core.GetTranslationSourceText(sourceType, sourceID)
}
func (s *TranslationService) UpsertTranslationResult(sourceType, sourceID, targetLocale, title, summary, content, status string, publishedAt time.Time, translatedByJobID string) error {
	return s.core.UpsertTranslationResult(sourceType, sourceID, targetLocale, title, summary, content, status, publishedAt, translatedByJobID)
}
func (s *TranslationService) ListTranslationContents(page, pageSize int, sourceType, sourceID, locale string) ([]domain.TranslationContent, int, error) {
	return s.core.ListTranslationContents(page, pageSize, sourceType, sourceID, locale)
}
func (s *TranslationService) GetTranslationContent(sourceType, sourceID, locale string) (domain.TranslationContent, bool, error) {
	return s.core.GetTranslationContent(sourceType, sourceID, locale)
}
func (s *TranslationService) UpsertTranslationContent(sourceType, sourceID, locale, title, summary, content, status string, publishedAt time.Time, translatedByJobID string) (domain.TranslationContent, error) {
	return s.core.UpsertTranslationContent(sourceType, sourceID, locale, title, summary, content, status, publishedAt, translatedByJobID)
}

func (s *TimelineService) ListTimeline(page, pageSize int, locale string) ([]domain.TimelineItem, int) {
	return s.core.ListTimeline(page, pageSize, locale)
}
