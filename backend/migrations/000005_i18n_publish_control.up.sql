-- File: 000005_i18n_publish_control.up.sql
-- Purpose: Add translation publish control and richer article translation fields.
-- Module: backend/migrations, schema evolution layer.
-- Related: translation job auto-publish/schedule and locale-aware public reads.

ALTER TABLE translation_jobs
    ADD COLUMN IF NOT EXISTS auto_publish BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS publish_at TIMESTAMPTZ;

ALTER TABLE article_translations
    ADD COLUMN IF NOT EXISTS title VARCHAR(512),
    ADD COLUMN IF NOT EXISTS summary TEXT,
    ADD COLUMN IF NOT EXISTS status VARCHAR(16) NOT NULL DEFAULT 'draft',
    ADD COLUMN IF NOT EXISTS published_at TIMESTAMPTZ;

ALTER TABLE moment_translations
    ADD COLUMN IF NOT EXISTS status VARCHAR(16) NOT NULL DEFAULT 'draft',
    ADD COLUMN IF NOT EXISTS published_at TIMESTAMPTZ;

UPDATE article_translations
SET status = 'published',
    published_at = COALESCE(published_at, created_at)
WHERE status = 'draft';

UPDATE moment_translations
SET status = 'published',
    published_at = COALESCE(published_at, created_at)
WHERE status = 'draft';

CREATE INDEX IF NOT EXISTS idx_article_translations_status_published
    ON article_translations (status, published_at DESC, updated_at DESC);

CREATE INDEX IF NOT EXISTS idx_moment_translations_status_published
    ON moment_translations (status, published_at DESC, updated_at DESC);
