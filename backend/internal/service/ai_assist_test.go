// File: ai_assist_test.go
// Purpose: Verify AI assist service fallback behavior for summary and slug generation.
// Module: backend/internal/service, AI assist unit test layer.
// Related: ai_assist.go.
package service

import (
	"context"
	"testing"

	"github.com/anxcye/ancy-blog/backend/internal/domain"
	"github.com/anxcye/ancy-blog/backend/internal/repository"
)

type aiRepoStub struct {
	repository.ContentRepository
	slugExistsFunc            func(slug string) bool
	getIntegrationProviderFun func(providerKey string) (domain.IntegrationProvider, bool)
}

func (s *aiRepoStub) SlugExists(slug string) bool {
	if s.slugExistsFunc != nil {
		return s.slugExistsFunc(slug)
	}
	return false
}

func (s *aiRepoStub) GetIntegrationProvider(providerKey string) (domain.IntegrationProvider, bool) {
	if s.getIntegrationProviderFun != nil {
		return s.getIntegrationProviderFun(providerKey)
	}
	return domain.IntegrationProvider{}, false
}

func TestAIAssistSuggestSlugFallbackAndUnique(t *testing.T) {
	repo := &aiRepoStub{slugExistsFunc: func(slug string) bool { return slug == "hello-world" }}
	core := NewContentService(repo, nil)
	ai := NewAIAssistService(NewArticleService(core), NewIntegrationService(core))

	slug, fallback, _, err := ai.SuggestSlug(context.Background(), "Hello World", "", "")
	if err != nil {
		t.Fatalf("suggest slug failed: %v", err)
	}
	if !fallback {
		t.Fatalf("expected fallback mode when provider is missing")
	}
	if slug != "hello-world-2" {
		t.Fatalf("unexpected slug: %s", slug)
	}
}

func TestAIAssistGenerateSummaryFallback(t *testing.T) {
	repo := &aiRepoStub{}
	core := NewContentService(repo, nil)
	ai := NewAIAssistService(NewArticleService(core), NewIntegrationService(core))

	summary, fallback, _, err := ai.GenerateSummary(context.Background(), "", "abcdefg", "", "", 4)
	if err != nil {
		t.Fatalf("generate summary failed: %v", err)
	}
	if !fallback {
		t.Fatalf("expected fallback mode")
	}
	if summary != "abcd" {
		t.Fatalf("unexpected summary: %s", summary)
	}
}
