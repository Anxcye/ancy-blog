// File: content_stub_test.go
// Purpose: Provide reusable repository stubs for handler-level tests.
// Module: backend/internal/handler, test support layer.
// Related: public/admin handlers and content service wiring.
package handler

import (
	"errors"
	"time"

	"github.com/anxcye/ancy-blog/backend/internal/domain"
	"github.com/anxcye/ancy-blog/backend/internal/repository"
)

type handlerRepoStub struct {
	repository.ContentRepository

	createCommentFunc            func(comment domain.Comment) (domain.Comment, error)
	listArticlesFunc             func(page, pageSize int, status, contentKind, keyword string) ([]domain.Article, int)
	getArticleByIDFunc           func(id string) (domain.Article, bool)
	listMomentsFunc              func(page, pageSize int, status string) ([]domain.Moment, int)
	updateMomentFunc             func(id string, moment domain.Moment) (domain.Moment, error)
	getPublishedArticleBySlug    func(slug string) (domain.Article, bool)
	getPublishedMomentByIDFunc   func(id, locale string) (domain.Moment, bool)
	listPublishedArticlesFunc    func(page, pageSize int, category, tag, contentKind string) ([]domain.Article, int)
	listPublishedMomentsFunc     func(page, pageSize int, locale string) ([]domain.Moment, int)
	listArticleCommentsFunc      func(articleID string, page, pageSize int) ([]domain.Comment, int)
	listContentCommentsFunc      func(contentType, contentID string, page, pageSize int) ([]domain.Comment, int)
	listTimelineFunc             func(page, pageSize int, locale string) ([]domain.TimelineItem, int)
	listFooterItemsFunc          func() []domain.FooterItem
	countArticleCommentsFunc     func(articleID string) (int, error)
	countContentCommentsFunc     func(contentType, contentID string) (int, error)
	listCommentDescendantsFunc   func(rootIDs []string) []domain.Comment
	slugExistsFunc               func(slug string) bool
	listIntegrationProvidersFunc func(providerType string) []domain.IntegrationProvider
	updateIntegrationProvider    func(providerKey string, enabled bool, configJSON, metaJSON []byte) (domain.IntegrationProvider, error)
	getIntegrationProviderFunc   func(providerKey string) (domain.IntegrationProvider, bool)
	createTranslationJobFunc     func(job domain.TranslationJob) (domain.TranslationJob, error)
	getTranslationJobByIDFunc    func(id string) (domain.TranslationJob, bool)
	retryTranslationJobFunc      func(id string) (domain.TranslationJob, error)
	listTranslationContentsFunc  func(page, pageSize int, sourceType, sourceID, locale string) ([]domain.TranslationContent, int)
	getTranslationContentFunc    func(sourceType, sourceID, locale string) (domain.TranslationContent, bool)
	upsertTranslationContentFunc func(sourceType, sourceID, locale, title, summary, content, status string, publishedAt time.Time, translatedByJobID string) (domain.TranslationContent, error)
}

func (s *handlerRepoStub) CreateComment(comment domain.Comment) (domain.Comment, error) {
	if s.createCommentFunc != nil {
		return s.createCommentFunc(comment)
	}
	return domain.Comment{}, errors.New("not implemented")
}

// GetSiteSettings returns a permissive default so comment tests pass without extra setup.
func (s *handlerRepoStub) GetSiteSettings() domain.SiteSettings {
	return domain.SiteSettings{SiteName: "Test", DefaultLocale: "en", CommentEnabled: true}
}

func (s *handlerRepoStub) ListArticles(page, pageSize int, status, contentKind, keyword string) ([]domain.Article, int) {
	if s.listArticlesFunc != nil {
		return s.listArticlesFunc(page, pageSize, status, contentKind, keyword)
	}
	return nil, 0
}

func (s *handlerRepoStub) GetArticleByID(id string) (domain.Article, bool) {
	if s.getArticleByIDFunc != nil {
		return s.getArticleByIDFunc(id)
	}
	return domain.Article{}, false
}

func (s *handlerRepoStub) ListMoments(page, pageSize int, status string) ([]domain.Moment, int) {
	if s.listMomentsFunc != nil {
		return s.listMomentsFunc(page, pageSize, status)
	}
	return nil, 0
}

func (s *handlerRepoStub) UpdateMoment(id string, moment domain.Moment) (domain.Moment, error) {
	if s.updateMomentFunc != nil {
		return s.updateMomentFunc(id, moment)
	}
	return domain.Moment{}, errors.New("not implemented")
}

func (s *handlerRepoStub) GetPublishedArticleBySlug(slug string) (domain.Article, bool) {
	if s.getPublishedArticleBySlug != nil {
		return s.getPublishedArticleBySlug(slug)
	}
	return domain.Article{}, false
}

func (s *handlerRepoStub) GetPublishedMomentByID(id, locale string) (domain.Moment, bool) {
	if s.getPublishedMomentByIDFunc != nil {
		return s.getPublishedMomentByIDFunc(id, locale)
	}
	return domain.Moment{}, false
}

func (s *handlerRepoStub) ListPublishedArticles(page, pageSize int, category, tag, contentKind string) ([]domain.Article, int) {
	if s.listPublishedArticlesFunc != nil {
		return s.listPublishedArticlesFunc(page, pageSize, category, tag, contentKind)
	}
	return nil, 0
}

func (s *handlerRepoStub) ListPublishedMoments(page, pageSize int, locale string) ([]domain.Moment, int) {
	if s.listPublishedMomentsFunc != nil {
		return s.listPublishedMomentsFunc(page, pageSize, locale)
	}
	return nil, 0
}

func (s *handlerRepoStub) ListArticleComments(articleID string, page, pageSize int) ([]domain.Comment, int) {
	if s.listArticleCommentsFunc != nil {
		return s.listArticleCommentsFunc(articleID, page, pageSize)
	}
	return nil, 0
}

func (s *handlerRepoStub) ListContentComments(contentType, contentID string, page, pageSize int) ([]domain.Comment, int) {
	if s.listContentCommentsFunc != nil {
		return s.listContentCommentsFunc(contentType, contentID, page, pageSize)
	}
	return nil, 0
}

func (s *handlerRepoStub) ListTimeline(page, pageSize int, locale string) ([]domain.TimelineItem, int) {
	if s.listTimelineFunc != nil {
		return s.listTimelineFunc(page, pageSize, locale)
	}
	return nil, 0
}

func (s *handlerRepoStub) ListFooterItems() []domain.FooterItem {
	if s.listFooterItemsFunc != nil {
		return s.listFooterItemsFunc()
	}
	return nil
}

func (s *handlerRepoStub) CountArticleComments(articleID string) (int, error) {
	if s.countArticleCommentsFunc != nil {
		return s.countArticleCommentsFunc(articleID)
	}
	return 0, nil
}

func (s *handlerRepoStub) CountContentComments(contentType, contentID string) (int, error) {
	if s.countContentCommentsFunc != nil {
		return s.countContentCommentsFunc(contentType, contentID)
	}
	return 0, nil
}

func (s *handlerRepoStub) ListCommentDescendants(rootIDs []string) []domain.Comment {
	if s.listCommentDescendantsFunc != nil {
		return s.listCommentDescendantsFunc(rootIDs)
	}
	return nil
}

func (s *handlerRepoStub) SlugExists(slug string) bool {
	if s.slugExistsFunc != nil {
		return s.slugExistsFunc(slug)
	}
	return false
}

func (s *handlerRepoStub) ListIntegrationProviders(providerType string) []domain.IntegrationProvider {
	if s.listIntegrationProvidersFunc != nil {
		return s.listIntegrationProvidersFunc(providerType)
	}
	return nil
}

func (s *handlerRepoStub) UpdateIntegrationProvider(providerKey string, enabled bool, configJSON, metaJSON []byte) (domain.IntegrationProvider, error) {
	if s.updateIntegrationProvider != nil {
		return s.updateIntegrationProvider(providerKey, enabled, configJSON, metaJSON)
	}
	return domain.IntegrationProvider{}, errors.New("provider not found")
}

func (s *handlerRepoStub) GetIntegrationProvider(providerKey string) (domain.IntegrationProvider, bool) {
	if s.getIntegrationProviderFunc != nil {
		return s.getIntegrationProviderFunc(providerKey)
	}
	return domain.IntegrationProvider{}, false
}

func (s *handlerRepoStub) CreateTranslationJob(job domain.TranslationJob) (domain.TranslationJob, error) {
	if s.createTranslationJobFunc != nil {
		return s.createTranslationJobFunc(job)
	}
	job.ID = "job-1"
	return job, nil
}

func (s *handlerRepoStub) GetTranslationJobByID(id string) (domain.TranslationJob, bool) {
	if s.getTranslationJobByIDFunc != nil {
		return s.getTranslationJobByIDFunc(id)
	}
	return domain.TranslationJob{}, false
}

func (s *handlerRepoStub) RetryTranslationJob(id string) (domain.TranslationJob, error) {
	if s.retryTranslationJobFunc != nil {
		return s.retryTranslationJobFunc(id)
	}
	return domain.TranslationJob{}, errors.New("not implemented")
}

func (s *handlerRepoStub) ListTranslationContents(page, pageSize int, sourceType, sourceID, locale string) ([]domain.TranslationContent, int) {
	if s.listTranslationContentsFunc != nil {
		return s.listTranslationContentsFunc(page, pageSize, sourceType, sourceID, locale)
	}
	return nil, 0
}

func (s *handlerRepoStub) GetTranslationContent(sourceType, sourceID, locale string) (domain.TranslationContent, bool) {
	if s.getTranslationContentFunc != nil {
		return s.getTranslationContentFunc(sourceType, sourceID, locale)
	}
	return domain.TranslationContent{}, false
}

func (s *handlerRepoStub) UpsertTranslationContent(sourceType, sourceID, locale, title, summary, content, status string, publishedAt time.Time, translatedByJobID string) (domain.TranslationContent, error) {
	if s.upsertTranslationContentFunc != nil {
		return s.upsertTranslationContentFunc(sourceType, sourceID, locale, title, summary, content, status, publishedAt, translatedByJobID)
	}
	return domain.TranslationContent{}, errors.New("not implemented")
}
