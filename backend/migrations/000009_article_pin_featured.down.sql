-- File: 000009_article_pin_featured.down.sql
-- Purpose: Revert is_pinned and is_featured columns from articles.
-- Module: backend/migrations, schema evolution layer.
-- Related: 000009_article_pin_featured.up.sql

DROP INDEX IF EXISTS idx_articles_pinned;
DROP INDEX IF EXISTS idx_articles_featured;

ALTER TABLE articles
    DROP COLUMN IF EXISTS is_pinned,
    DROP COLUMN IF EXISTS is_featured;
