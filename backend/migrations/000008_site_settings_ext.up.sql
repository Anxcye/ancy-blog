-- File: 000008_site_settings_ext.up.sql
-- Purpose: Extend site_settings with comment policy and SEO fields.
-- Module: backend/migrations, schema evolution layer.
-- Related: site_repo.go, domain/models.go, handler admin site endpoints.

ALTER TABLE site_settings
    ADD COLUMN IF NOT EXISTS comment_enabled          BOOLEAN NOT NULL DEFAULT TRUE,
    ADD COLUMN IF NOT EXISTS comment_require_approval BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS site_description         TEXT,
    ADD COLUMN IF NOT EXISTS seo_keywords             TEXT,
    ADD COLUMN IF NOT EXISTS og_image_url             TEXT;
