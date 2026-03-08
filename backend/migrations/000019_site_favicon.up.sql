-- File: 000019_site_favicon.up.sql
-- Purpose: Add favicon_url to site_settings for configurable public favicon metadata.
-- Module: backend/migrations, schema evolution layer.
-- Related: site_settings and public/admin site settings flows.

ALTER TABLE site_settings
ADD COLUMN IF NOT EXISTS favicon_url TEXT;
