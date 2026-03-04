-- File: 000007_article_category_tags.down.sql
-- Purpose: Revert category_id column from articles.
-- Module: backend/migrations, schema evolution layer.
-- Related: 000007_article_category_tags.up.sql

DROP INDEX IF EXISTS idx_articles_category_id;

ALTER TABLE articles
    DROP COLUMN IF EXISTS category_id;
