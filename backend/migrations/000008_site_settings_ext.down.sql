-- File: 000008_site_settings_ext.down.sql
-- Purpose: Revert site_settings comment policy and SEO field additions.
-- Module: backend/migrations, schema evolution layer.
-- Related: 000008_site_settings_ext.up.sql

ALTER TABLE site_settings
    DROP COLUMN IF EXISTS comment_enabled,
    DROP COLUMN IF EXISTS comment_require_approval,
    DROP COLUMN IF EXISTS site_description,
    DROP COLUMN IF EXISTS seo_keywords,
    DROP COLUMN IF EXISTS og_image_url;
