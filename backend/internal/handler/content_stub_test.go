// File: content_stub_test.go
// Purpose: Provide reusable repository stubs for handler-level tests.
// Module: backend/internal/handler, test support layer.
// Related: public/admin handlers and content service wiring.
package handler

import (
	"errors"

	"github.com/anxcye/ancy-blog/backend/internal/domain"
	"github.com/anxcye/ancy-blog/backend/internal/repository"
)

type handlerRepoStub struct {
	repository.ContentRepository

	createCommentFunc            func(comment domain.Comment) (domain.Comment, error)
	getPublishedArticleBySlug    func(slug string) (domain.Article, bool)
	listPublishedArticlesFunc    func(page, pageSize int, category, tag, contentKind string) ([]domain.Article, int)
	listPublishedMomentsFunc     func(page, pageSize int, locale string) ([]domain.Moment, int)
	listTimelineFunc             func(page, pageSize int, locale string) ([]domain.TimelineItem, int)
	listFooterItemsFunc          func() []domain.FooterItem
	countArticleCommentsFunc     func(articleID string) (int, error)
	slugExistsFunc               func(slug string) bool
	listIntegrationProvidersFunc func(providerType string) []domain.IntegrationProvider
	updateIntegrationProvider    func(providerKey string, enabled bool, configJSON, metaJSON []byte) (domain.IntegrationProvider, error)
	getIntegrationProviderFunc   func(providerKey string) (domain.IntegrationProvider, bool)
	createTranslationJobFunc     func(job domain.TranslationJob) (domain.TranslationJob, error)
	getTranslationJobByIDFunc    func(id string) (domain.TranslationJob, bool)
}

func (s *handlerRepoStub) CreateComment(comment domain.Comment) (domain.Comment, error) {
	if s.createCommentFunc != nil {
		return s.createCommentFunc(comment)
	}
	return domain.Comment{}, errors.New("not implemented")
}

func (s *handlerRepoStub) GetPublishedArticleBySlug(slug string) (domain.Article, bool) {
	if s.getPublishedArticleBySlug != nil {
		return s.getPublishedArticleBySlug(slug)
	}
	return domain.Article{}, false
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
