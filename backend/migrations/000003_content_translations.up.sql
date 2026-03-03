-- File: 000003_content_translations.up.sql
-- Purpose: Add locale-specific translation tables for articles and moments.
-- Module: backend/migrations, schema evolution layer.
-- Related: translation worker writeback and locale-aware reads.

CREATE TABLE IF NOT EXISTS article_translations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    article_id UUID NOT NULL REFERENCES articles(id) ON DELETE CASCADE,
    locale VARCHAR(16) NOT NULL,
    content TEXT NOT NULL,
    translated_by_job_id UUID REFERENCES translation_jobs(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (article_id, locale)
);

CREATE INDEX IF NOT EXISTS idx_article_translations_locale
    ON article_translations (locale, updated_at DESC);

CREATE TABLE IF NOT EXISTS moment_translations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    moment_id UUID NOT NULL REFERENCES moments(id) ON DELETE CASCADE,
    locale VARCHAR(16) NOT NULL,
    content TEXT NOT NULL,
    translated_by_job_id UUID REFERENCES translation_jobs(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (moment_id, locale)
);

CREATE INDEX IF NOT EXISTS idx_moment_translations_locale
    ON moment_translations (locale, updated_at DESC);
