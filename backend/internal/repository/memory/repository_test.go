// File: repository_test.go
// Purpose: Verify in-memory repository core CRUD and list behavior.
// Module: backend/internal/repository/memory, repository unit test layer.
// Related: repository.go.
package memory

import (
	"testing"
	"time"

	"github.com/anxcye/ancy-blog/backend/internal/domain"
)

func TestCreateAndGetArticle(t *testing.T) {
	repo := NewRepository()
	created, err := repo.CreateArticle(domain.Article{Title: "T1", Slug: "t1", ContentKind: "post", Status: "published"})
	if err != nil {
		t.Fatalf("expected create article success, got error: %v", err)
	}
	if created.ID == "" {
		t.Fatalf("expected generated article id")
	}
	got, ok := repo.GetArticleByID(created.ID)
	if !ok {
		t.Fatalf("expected article found by id")
	}
	if got.Slug != "t1" {
		t.Fatalf("unexpected article slug: %s", got.Slug)
	}
}

func TestListPublishedArticlesFiltering(t *testing.T) {
	repo := NewRepository()
	_, _ = repo.CreateArticle(domain.Article{Title: "Draft", Slug: "draft", ContentKind: "post", Status: "draft", CategorySlug: "tech"})
	_, _ = repo.CreateArticle(domain.Article{Title: "Pub", Slug: "pub", ContentKind: "post", Status: "published", CategorySlug: "tech", TagSlugs: []string{"go"}})

	rows, total := repo.ListPublishedArticles(1, 20, "tech", "go", "post")
	if total < 1 {
		t.Fatalf("expected at least one published article, got total=%d", total)
	}
	if len(rows) == 0 {
		t.Fatalf("expected non-empty rows")
	}
}

func TestListArticlesIncludesDraftAndKeywordFilter(t *testing.T) {
	repo := NewRepository()
	_, _ = repo.CreateArticle(domain.Article{Title: "Draft Note", Slug: "draft-note", ContentKind: "post", Status: "draft"})
	_, _ = repo.CreateArticle(domain.Article{Title: "Published Note", Slug: "pub-note", ContentKind: "post", Status: "published"})

	rows, total := repo.ListArticles(1, 20, "draft", "post", "draft")
	if total == 0 || len(rows) == 0 {
		t.Fatalf("expected draft rows, total=%d len=%d", total, len(rows))
	}
	if rows[0].Status != "draft" {
		t.Fatalf("expected draft status, got %s", rows[0].Status)
	}
}

func TestCreateSlotItemRequiresSlot(t *testing.T) {
	repo := NewRepository()
	if _, err := repo.CreateSlotItem("not_exists", domain.SlotItem{ContentType: "article", ContentID: "a1"}); err == nil {
		t.Fatalf("expected slot not found error")
	}
}

func TestLocaleTranslationForMomentAndTimeline(t *testing.T) {
	repo := NewRepository()

	moments, total := repo.ListPublishedMoments(1, 10, "")
	if total == 0 || len(moments) == 0 {
		t.Fatalf("expected seeded published moments")
	}
	momentID := moments[0].ID
	originalContent := moments[0].Content

	if err := repo.UpsertMomentTranslation(momentID, "en-US", "translated moment", "published", time.Now().UTC(), "job-1"); err != nil {
		t.Fatalf("upsert moment translation failed: %v", err)
	}

	localizedMoments, _ := repo.ListPublishedMoments(1, 10, "en-US")
	if localizedMoments[0].Content != "translated moment" {
		t.Fatalf("expected localized moment content, got %s", localizedMoments[0].Content)
	}

	defaultTimeline, _ := repo.ListTimeline(1, 50, "")
	foundDefault := false
	for _, item := range defaultTimeline {
		if item.ContentType == "moment" && item.ID == momentID {
			foundDefault = true
			if item.Content != originalContent {
				t.Fatalf("expected default timeline content, got %s", item.Content)
			}
		}
	}
	if !foundDefault {
		t.Fatalf("expected moment in timeline")
	}

	localizedTimeline, _ := repo.ListTimeline(1, 50, "en-US")
	foundLocalized := false
	for _, item := range localizedTimeline {
		if item.ContentType == "moment" && item.ID == momentID {
			foundLocalized = true
			if item.Content != "translated moment" {
				t.Fatalf("expected localized timeline content, got %s", item.Content)
			}
		}
	}
	if !foundLocalized {
		t.Fatalf("expected localized moment in timeline")
	}
}

func TestTranslationContentCRUD(t *testing.T) {
	repo := NewRepository()
	articles, total := repo.ListPublishedArticles(1, 10, "", "", "post")
	if total == 0 || len(articles) == 0 {
		t.Fatalf("expected seeded article")
	}
	articleID := articles[0].ID

	row, err := repo.UpsertTranslationContent("article", articleID, "en-US", "Translated Title", "Translated Summary", "translated article", "published", time.Now().UTC(), "job-1")
	if err != nil {
		t.Fatalf("upsert translation content failed: %v", err)
	}
	if row.Content != "translated article" {
		t.Fatalf("unexpected upserted content: %s", row.Content)
	}

	got, ok := repo.GetTranslationContent("article", articleID, "en-US")
	if !ok {
		t.Fatalf("expected translation content detail")
	}
	if got.Content != "translated article" {
		t.Fatalf("unexpected translation content detail: %s", got.Content)
	}

	items, count := repo.ListTranslationContents(1, 10, "article", articleID, "en-US")
	if count != 1 || len(items) != 1 {
		t.Fatalf("expected one translation content row, count=%d len=%d", count, len(items))
	}
}
