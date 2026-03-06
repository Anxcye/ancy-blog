-- File: schema_v1.sql
-- Purpose: Initialize PostgreSQL schema for the blog rewrite (v1) based on DATA_MODEL.
-- Module: backend/sql, database initialization layer.
-- Related: docs/DATA_MODEL.md and docs/API_CONTRACT.md.

BEGIN;

-- Enable UUID generation.
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- =========================
-- Core Content
-- =========================

CREATE TABLE IF NOT EXISTS articles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(256) NOT NULL,
    slug VARCHAR(256) NOT NULL UNIQUE,
    content_kind VARCHAR(16) NOT NULL CHECK (content_kind IN ('post', 'page')),
    summary TEXT,
    content TEXT,
    status VARCHAR(16) NOT NULL CHECK (status IN ('draft', 'published', 'archived')),
    visibility VARCHAR(16) NOT NULL DEFAULT 'public' CHECK (visibility IN ('public', 'unlisted', 'private')),
    is_pinned BOOLEAN NOT NULL DEFAULT FALSE,
    pin_order INTEGER NOT NULL DEFAULT 0,
    pinned_at TIMESTAMPTZ,
    allow_comment BOOLEAN NOT NULL DEFAULT TRUE,
    origin_type VARCHAR(16) NOT NULL DEFAULT 'original' CHECK (origin_type IN ('original', 'repost', 'translation')),
    source_url TEXT,
    ai_assist_level VARCHAR(16) NOT NULL DEFAULT 'none' CHECK (ai_assist_level IN ('none', 'polish', 'dictation', 'assisted', 'generated', 'translated')),
    cover_image TEXT,
    published_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_articles_content_kind_status_published_at ON articles (content_kind, status, published_at DESC);
CREATE INDEX IF NOT EXISTS idx_articles_status_published_at ON articles (status, published_at DESC);
CREATE INDEX IF NOT EXISTS idx_articles_pinned_order_published_at ON articles (is_pinned DESC, pin_order DESC, published_at DESC);

CREATE TABLE IF NOT EXISTS categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(128) NOT NULL UNIQUE,
    slug VARCHAR(128) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS tags (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(128) NOT NULL UNIQUE,
    slug VARCHAR(128) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS article_tags (
    article_id UUID NOT NULL REFERENCES articles(id) ON DELETE CASCADE,
    tag_id UUID NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (article_id, tag_id)
);

CREATE INDEX IF NOT EXISTS idx_article_tags_tag_article ON article_tags (tag_id, article_id);

CREATE TABLE IF NOT EXISTS moments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    content TEXT NOT NULL,
    status VARCHAR(16) NOT NULL CHECK (status IN ('draft', 'published', 'archived')),
    is_pinned BOOLEAN NOT NULL DEFAULT FALSE,
    pin_order INTEGER NOT NULL DEFAULT 0,
    allow_comment BOOLEAN NOT NULL DEFAULT TRUE,
    source VARCHAR(16) NOT NULL DEFAULT 'admin' CHECK (source IN ('web', 'admin', 'api')),
    published_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_moments_status_published_at ON moments (status, published_at DESC);
CREATE INDEX IF NOT EXISTS idx_moments_pinned_order_published_at ON moments (is_pinned DESC, pin_order DESC, published_at DESC);

-- =========================
-- Comments & Reactions
-- =========================

CREATE TABLE IF NOT EXISTS comments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    article_id UUID REFERENCES articles(id) ON DELETE CASCADE,
    content_type VARCHAR(16) NOT NULL CHECK (content_type IN ('article', 'moment')),
    content_id UUID NOT NULL,
    parent_id UUID REFERENCES comments(id) ON DELETE SET NULL,
    root_id UUID REFERENCES comments(id) ON DELETE SET NULL,
    content TEXT NOT NULL,
    status VARCHAR(16) NOT NULL CHECK (status IN ('pending', 'approved', 'rejected', 'spam', 'deleted')),
    is_pinned BOOLEAN NOT NULL DEFAULT FALSE,
    like_count INTEGER NOT NULL DEFAULT 0,
    reply_count INTEGER NOT NULL DEFAULT 0,
    nickname VARCHAR(64) NOT NULL,
    email VARCHAR(128),
    website TEXT,
    avatar_url TEXT,
    source VARCHAR(16) NOT NULL DEFAULT 'web' CHECK (source IN ('web', 'admin', 'api')),
    ip VARCHAR(64) NOT NULL,
    user_agent TEXT,
    risk_score INTEGER NOT NULL DEFAULT 0,
    approved_at TIMESTAMPTZ,
    approved_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT chk_comments_parent_not_self CHECK (parent_id IS NULL OR parent_id <> id)
);

CREATE INDEX IF NOT EXISTS idx_comments_content_status_created_at ON comments (content_type, content_id, status, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_comments_article_status_created_at ON comments (article_id, status, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_comments_parent_created_at ON comments (parent_id, created_at ASC);
CREATE INDEX IF NOT EXISTS idx_comments_root_created_at ON comments (root_id, created_at ASC);
CREATE INDEX IF NOT EXISTS idx_comments_ip_created_at ON comments (ip, created_at DESC);

CREATE TABLE IF NOT EXISTS reactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    target_type VARCHAR(16) NOT NULL CHECK (target_type IN ('article', 'comment')),
    target_id UUID NOT NULL,
    reaction_type VARCHAR(16) NOT NULL CHECK (reaction_type IN ('like', 'love', 'clap', 'insightful')),
    actor_type VARCHAR(16) NOT NULL CHECK (actor_type IN ('admin', 'visitor')),
    actor_id UUID,
    visitor_key VARCHAR(256),
    ip VARCHAR(64),
    user_agent TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_reactions_actor_or_visitor CHECK (actor_id IS NOT NULL OR visitor_key IS NOT NULL)
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_reactions_target_actor
    ON reactions (target_type, target_id, reaction_type, actor_type, actor_id, visitor_key);
CREATE INDEX IF NOT EXISTS idx_reactions_target_type_id_reaction
    ON reactions (target_type, target_id, reaction_type);

-- =========================
-- Friend Links
-- =========================

CREATE TABLE IF NOT EXISTS links (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(128) NOT NULL,
    url TEXT NOT NULL UNIQUE,
    avatar_url TEXT,
    description TEXT,
    contact_email VARCHAR(128),
    review_status VARCHAR(16) NOT NULL DEFAULT 'pending' CHECK (review_status IN ('pending', 'approved', 'rejected')),
    review_note TEXT,
    submitted_ip VARCHAR(64),
    submitted_user_agent TEXT,
    related_article_id UUID REFERENCES articles(id) ON DELETE SET NULL,
    approved_at TIMESTAMPTZ,
    approved_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_links_review_status_created_at ON links (review_status, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_links_related_article_id ON links (related_article_id);

-- =========================
-- Site Configuration
-- =========================

CREATE TABLE IF NOT EXISTS site_settings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    site_name VARCHAR(128) NOT NULL,
    avatar_url TEXT,
    hero_intro_md TEXT,
    default_locale VARCHAR(16) NOT NULL DEFAULT 'en',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS nav_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(128) NOT NULL,
    key VARCHAR(128) NOT NULL UNIQUE,
    type VARCHAR(16) NOT NULL CHECK (type IN ('menu', 'dropdown', 'link')),
    target_type VARCHAR(16) NOT NULL CHECK (target_type IN ('route', 'category', 'slot', 'external')),
    target_value TEXT,
    order_num INTEGER NOT NULL DEFAULT 0,
    enabled BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_nav_items_enabled_order_num ON nav_items (enabled, order_num ASC);

CREATE TABLE IF NOT EXISTS content_slots (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    slot_key VARCHAR(128) NOT NULL UNIQUE,
    name VARCHAR(128) NOT NULL,
    description TEXT,
    enabled BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS content_slot_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    slot_key VARCHAR(128) NOT NULL REFERENCES content_slots(slot_key) ON DELETE CASCADE,
    content_type VARCHAR(16) NOT NULL CHECK (content_type IN ('article', 'moment')),
    content_id UUID NOT NULL,
    order_num INTEGER NOT NULL DEFAULT 0,
    enabled BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_content_slot_items_unique ON content_slot_items (slot_key, content_type, content_id);
CREATE INDEX IF NOT EXISTS idx_content_slot_items_slot_enabled_order ON content_slot_items (slot_key, enabled, order_num ASC);

CREATE TABLE IF NOT EXISTS footer_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    label VARCHAR(256) NOT NULL,
    link_type VARCHAR(16) NOT NULL CHECK (link_type IN ('none', 'internal', 'external')),
    internal_article_slug VARCHAR(256),
    external_url TEXT,
    row_num INTEGER NOT NULL CHECK (row_num BETWEEN 1 AND 3),
    order_num INTEGER NOT NULL DEFAULT 0,
    enabled BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT chk_footer_internal_slug_required CHECK (
        link_type <> 'internal' OR internal_article_slug IS NOT NULL
    ),
    CONSTRAINT chk_footer_external_url_required CHECK (
        link_type <> 'external' OR external_url IS NOT NULL
    )
);

CREATE INDEX IF NOT EXISTS idx_footer_items_row_order ON footer_items (row_num, order_num ASC);
CREATE INDEX IF NOT EXISTS idx_footer_items_enabled_row_order ON footer_items (enabled, row_num, order_num ASC);

CREATE TABLE IF NOT EXISTS social_links (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    platform VARCHAR(32) NOT NULL CHECK (platform IN ('github', 'mail', 'x', 'linkedin', 'custom')),
    title VARCHAR(128) NOT NULL,
    url TEXT NOT NULL,
    icon_key VARCHAR(128),
    order_num INTEGER NOT NULL DEFAULT 0,
    enabled BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_social_links_enabled_order ON social_links (enabled, order_num ASC);

-- =========================
-- Integration Center
-- =========================

CREATE TABLE IF NOT EXISTS integration_providers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    provider_type VARCHAR(32) NOT NULL CHECK (provider_type IN ('object_storage', 'llm')),
    provider_key VARCHAR(64) NOT NULL UNIQUE,
    name VARCHAR(128) NOT NULL,
    enabled BOOLEAN NOT NULL DEFAULT FALSE,
    config_json JSONB NOT NULL,
    meta_json JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_integration_providers_type_enabled ON integration_providers (provider_type, enabled);

CREATE TABLE IF NOT EXISTS translation_jobs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    source_type VARCHAR(16) NOT NULL CHECK (source_type IN ('article', 'moment')),
    source_id UUID NOT NULL,
    source_locale VARCHAR(16) NOT NULL,
    target_locale VARCHAR(16) NOT NULL,
    provider_key VARCHAR(64) NOT NULL REFERENCES integration_providers(provider_key) ON DELETE RESTRICT,
    model_name VARCHAR(128) NOT NULL,
    status VARCHAR(16) NOT NULL CHECK (status IN ('queued', 'running', 'succeeded', 'failed')),
    error_message TEXT,
    result_text TEXT,
    requested_by UUID,
    retry_count INTEGER NOT NULL DEFAULT 0,
    max_retries INTEGER NOT NULL DEFAULT 3,
    next_retry_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    auto_publish BOOLEAN NOT NULL DEFAULT FALSE,
    publish_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    finished_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_translation_jobs_source_target
    ON translation_jobs (source_type, source_id, target_locale);
CREATE INDEX IF NOT EXISTS idx_translation_jobs_status_created_at
    ON translation_jobs (status, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_translation_jobs_status_next_retry
    ON translation_jobs (status, next_retry_at ASC, created_at ASC);

CREATE TABLE IF NOT EXISTS article_translations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    article_id UUID NOT NULL REFERENCES articles(id) ON DELETE CASCADE,
    locale VARCHAR(16) NOT NULL,
    title VARCHAR(512),
    summary TEXT,
    content TEXT NOT NULL,
    status VARCHAR(16) NOT NULL DEFAULT 'draft',
    published_at TIMESTAMPTZ,
    translated_by_job_id UUID REFERENCES translation_jobs(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (article_id, locale)
);

CREATE INDEX IF NOT EXISTS idx_article_translations_locale
    ON article_translations (locale, updated_at DESC);
CREATE INDEX IF NOT EXISTS idx_article_translations_status_published
    ON article_translations (status, published_at DESC, updated_at DESC);

CREATE TABLE IF NOT EXISTS moment_translations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    moment_id UUID NOT NULL REFERENCES moments(id) ON DELETE CASCADE,
    locale VARCHAR(16) NOT NULL,
    content TEXT NOT NULL,
    status VARCHAR(16) NOT NULL DEFAULT 'draft',
    published_at TIMESTAMPTZ,
    translated_by_job_id UUID REFERENCES translation_jobs(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (moment_id, locale)
);

CREATE INDEX IF NOT EXISTS idx_moment_translations_locale
    ON moment_translations (locale, updated_at DESC);
CREATE INDEX IF NOT EXISTS idx_moment_translations_status_published
    ON moment_translations (status, published_at DESC, updated_at DESC);

-- =========================
-- Optional seed data
-- =========================

INSERT INTO site_settings (site_name, avatar_url, hero_intro_md, default_locale)
SELECT 'Ancy Blog', '', 'Hi, I build things.', 'en'
WHERE NOT EXISTS (SELECT 1 FROM site_settings);

INSERT INTO integration_providers (provider_type, provider_key, name, enabled, config_json, meta_json)
SELECT 'object_storage', 'cloudflare_r2', 'Cloudflare R2', FALSE, '{}'::jsonb, '{}'::jsonb
WHERE NOT EXISTS (SELECT 1 FROM integration_providers WHERE provider_key = 'cloudflare_r2');

INSERT INTO integration_providers (provider_type, provider_key, name, enabled, config_json, meta_json)
SELECT 'llm', 'openai_compatible', 'OpenAI Compatible', FALSE, '{}'::jsonb, '{}'::jsonb
WHERE NOT EXISTS (SELECT 1 FROM integration_providers WHERE provider_key = 'openai_compatible');

COMMIT;
