// File: repository_test.go
// Purpose: Verify in-memory repository core CRUD and list behavior.
// Module: backend/internal/repository/memory, repository unit test layer.
// Related: repository.go.
package memory

import (
	"testing"

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

func TestCreateSlotItemRequiresSlot(t *testing.T) {
	repo := NewRepository()
	if _, err := repo.CreateSlotItem("not_exists", domain.SlotItem{ContentType: "article", ContentID: "a1"}); err == nil {
		t.Fatalf("expected slot not found error")
	}
}
