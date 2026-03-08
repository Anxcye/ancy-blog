-- File: 000019_site_favicon.down.sql
-- Purpose: Remove favicon_url from site_settings.
-- Module: backend/migrations, schema rollback layer.
-- Related: 000019_site_favicon.up.sql and site_settings.

ALTER TABLE site_settings
DROP COLUMN IF EXISTS favicon_url;
