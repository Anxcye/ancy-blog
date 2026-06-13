-- File: 000026_site_hero_quotes.down.sql
-- Purpose: Remove configurable localized homepage quote candidates.
-- Module: backend/migrations, database schema rollback.
-- Related: 000026_site_hero_quotes.up.sql and site_settings.

ALTER TABLE site_settings
  DROP COLUMN IF EXISTS hero_quotes;
