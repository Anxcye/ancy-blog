-- File: 000009_article_pin_featured.up.sql
-- Purpose: Add is_pinned and is_featured flags to articles table.
-- Module: backend/migrations, schema evolution layer.
-- Related: content_repo.go, domain/models.go, handler admin article endpoints.

ALTER TABLE articles
    ADD COLUMN IF NOT EXISTS is_pinned   BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS is_featured BOOLEAN NOT NULL DEFAULT FALSE;

CREATE INDEX IF NOT EXISTS idx_articles_pinned   ON articles (is_pinned)   WHERE is_pinned = TRUE AND deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_articles_featured ON articles (is_featured) WHERE is_featured = TRUE AND deleted_at IS NULL;
