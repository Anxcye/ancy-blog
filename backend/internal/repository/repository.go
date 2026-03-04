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
	ListArticles(page, pageSize int, status, contentKind, keyword string) ([]domain.Article, int)
	ListPublishedArticles(page, pageSize int, category, tag, contentKind string) ([]domain.Article, int)
	GetPublishedArticleBySlug(slug string) (domain.Article, bool)
	GetPublishedArticleBySlugWithLocale(slug, locale string) (domain.Article, bool)
	SlugExists(slug string) bool
	GetArticleByID(id string) (domain.Article, bool)

	CreateMoment(moment domain.Moment) (domain.Moment, error)
	UpdateMoment(id string, moment domain.Moment) (domain.Moment, error)
	ListMoments(page, pageSize int, status string) ([]domain.Moment, int)
	ListPublishedMoments(page, pageSize int, locale string) ([]domain.Moment, int)

	CreateComment(comment domain.Comment) (domain.Comment, error)
	ListArticleComments(articleID string, page, pageSize int) ([]domain.Comment, int)
	ListCommentChildren(parentID string, page, pageSize int) ([]domain.Comment, int)
	CountArticleComments(articleID string) (int, error)
	ListCommentPage(page, pageSize int, status string) ([]domain.Comment, int)
	UpdateCommentAdmin(id string, status, isPinned string) (domain.Comment, error)

	SubmitLink(link domain.Link) (domain.Link, error)
	ListApprovedLinks() []domain.Link
	ListLinkSubmissions(page, pageSize int, reviewStatus string) ([]domain.Link, int)
	ReviewLink(id, reviewStatus, reviewNote, relatedArticleID string) (domain.Link, error)

	ListCategories() []domain.Category
	ListTags() []domain.Tag

	GetSiteSettings() domain.SiteSettings
	UpdateSiteSettings(settings domain.SiteSettings) domain.SiteSettings

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
