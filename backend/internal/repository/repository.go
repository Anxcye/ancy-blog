// File: repository.go
// Purpose: Declare persistence contracts for content-related domains.
// Module: backend/internal/repository, repository abstraction layer.
// Related: repository implementations and content service orchestration.
package repository

import (
	"time"

	"github.com/anxcye/ancy-blog/backend/internal/domain"
)

type ContentRepository interface {
	CreateArticle(article domain.Article) (domain.Article, error)
	UpdateArticle(id string, article domain.Article) (domain.Article, error)
	DeleteArticle(id string) bool
	BatchUpdateArticleStatus(ids []string, status string) int
	ListArticles(page, pageSize int, status, contentKind, keyword string) ([]domain.Article, int)
	ListPublishedArticles(page, pageSize int, category, tag, contentKind string) ([]domain.Article, int)
	GetPublishedArticleBySlug(slug string) (domain.Article, bool)
	GetPublishedArticleBySlugWithLocale(slug, locale string) (domain.Article, bool)
	SlugExists(slug string) bool
	GetArticleByID(id string) (domain.Article, bool)
	// RecordView inserts a deduplicated view row and bumps the denormalized counter.
	// visitorKey = SHA-256(ip + user-agent + date); returns the new total.
	RecordView(articleID, visitorKey string) (int64, error)

	CreateMoment(moment domain.Moment) (domain.Moment, error)
	UpdateMoment(id string, moment domain.Moment) (domain.Moment, error)
	DeleteMoment(id string) bool
	BatchUpdateMomentStatus(ids []string, status string) int
	ListMoments(page, pageSize int, status string) ([]domain.Moment, int)
	GetMomentByID(id string) (domain.Moment, bool)
	ListPublishedMoments(page, pageSize int, locale string) ([]domain.Moment, int)
	GetPublishedMomentByID(id, locale string) (domain.Moment, bool)

	CreateComment(comment domain.Comment) (domain.Comment, error)
	ListArticleComments(articleID string, page, pageSize int) ([]domain.Comment, int)
	ListContentComments(contentType, contentID string, page, pageSize int) ([]domain.Comment, int)
	ListCommentChildren(parentID string, page, pageSize int) ([]domain.Comment, int)
	ListCommentDescendants(rootIDs []string) []domain.Comment
	CountArticleComments(articleID string) (int, error)
	CountContentComments(contentType, contentID string) (int, error)
	ListCommentPage(page, pageSize int, status string) ([]domain.Comment, int)
	GetCommentByID(id string) (domain.Comment, bool)
	UpdateCommentAdmin(id string, status, isPinned string) (domain.Comment, error)

	SubmitLink(link domain.Link) (domain.Link, error)
	ListApprovedLinks() []domain.Link
	ListLinkSubmissions(page, pageSize int, reviewStatus string) ([]domain.Link, int)
	ReviewLink(id, reviewStatus, reviewNote, relatedArticleID string) (domain.Link, error)

	ListCategories() []domain.Category
	CreateCategory(category domain.Category) (domain.Category, error)
	DeleteCategory(id string) bool
	ListTags() []domain.Tag
	CreateTag(tag domain.Tag) (domain.Tag, error)
	DeleteTag(id string) bool

	GetSiteSettings() domain.SiteSettings
	UpdateSiteSettings(settings domain.SiteSettings) domain.SiteSettings
	GetTranslationPolicy() domain.TranslationPolicy
	UpdateTranslationPolicy(policy domain.TranslationPolicy) error

	CreateFooterItem(item domain.FooterItem) (domain.FooterItem, error)
	UpdateFooterItem(id string, item domain.FooterItem) (domain.FooterItem, error)
	DeleteFooterItem(id string) bool
	ListFooterItems() []domain.FooterItem

	CreateSocialLink(item domain.SocialLink) (domain.SocialLink, error)
	UpdateSocialLink(id string, item domain.SocialLink) (domain.SocialLink, error)
	DeleteSocialLink(id string) bool
	ListSocialLinks() []domain.SocialLink

	CreateNavItem(item domain.NavItem) (domain.NavItem, error)
	UpdateNavItem(id string, item domain.NavItem) (domain.NavItem, error)
	DeleteNavItem(id string) bool
	ListNavItems() []domain.NavItem

	CreateContentSlot(slot domain.ContentSlot) (domain.ContentSlot, error)
	ListContentSlots() []domain.ContentSlot
	CreateSlotItem(slotKey string, item domain.SlotItem) (domain.SlotItem, error)
	ListSlotItems(slotKey string) ([]domain.SlotItem, bool)
	DeleteSlotItem(slotKey, itemID string) bool
	ListSlotContent(slotKey string, limit int) ([]domain.SlotContentItem, bool)

	ListIntegrationProviders(providerType string) []domain.IntegrationProvider
	GetIntegrationProvider(providerKey string) (domain.IntegrationProvider, bool)
	UpdateIntegrationProvider(providerKey string, enabled bool, configJSON, metaJSON []byte) (domain.IntegrationProvider, error)

	CreateTranslationJob(job domain.TranslationJob) (domain.TranslationJob, error)
	ListTranslationJobs(page, pageSize int, status, sourceType, sourceID string) ([]domain.TranslationJob, int)
	GetTranslationJobByID(id string) (domain.TranslationJob, bool)
	ClaimNextQueuedTranslationJob() (domain.TranslationJob, bool, error)
	MarkTranslationJobRunning(id string) error
	MarkTranslationJobSucceeded(id, resultText string) error
	MarkTranslationJobFailed(id, errorMessage string) error
	ScheduleTranslationJobRetry(id, errorMessage string, nextRetryAt time.Time) error
	RetryTranslationJob(id string) (domain.TranslationJob, error)
	GetTranslationSourceText(sourceType, sourceID string) (string, bool, error)
	UpsertArticleTranslation(articleID, locale, title, summary, content, status string, publishedAt time.Time, translatedByJobID string) error
	UpsertMomentTranslation(momentID, locale, content, status string, publishedAt time.Time, translatedByJobID string) error
	ListTranslationContents(page, pageSize int, sourceType, sourceID, locale string) ([]domain.TranslationContent, int)
	GetTranslationContent(sourceType, sourceID, locale string) (domain.TranslationContent, bool)
	UpsertTranslationContent(sourceType, sourceID, locale, title, summary, content, status string, publishedAt time.Time, translatedByJobID string) (domain.TranslationContent, error)

	ListTimeline(page, pageSize int, locale string) ([]domain.TimelineItem, int)
}

// CredentialStore persists admin credentials independently of the in-memory session layer.
type CredentialStore interface {
	GetAdminPasswordHash() (string, bool)
	SetAdminPasswordHash(hash string) error
}
