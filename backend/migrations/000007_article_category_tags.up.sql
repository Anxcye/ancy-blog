-- File: 000007_article_category_tags.up.sql
-- Purpose: Add category_id FK to articles and enable admin CRUD for categories/tags.
-- Module: backend/migrations, schema evolution layer.
-- Related: content_repo.go, domain/models.go, handler admin category/tag endpoints.

ALTER TABLE articles
    ADD COLUMN IF NOT EXISTS category_id UUID REFERENCES categories(id) ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS idx_articles_category_id ON articles (category_id) WHERE category_id IS NOT NULL;
