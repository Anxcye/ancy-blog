-- File: 000026_site_hero_quotes.up.sql
-- Purpose: Add configurable localized homepage quote candidates.
-- Module: backend/migrations, database schema evolution.
-- Related: site_settings, public site settings API, admin site settings UI.

ALTER TABLE site_settings
  ADD COLUMN IF NOT EXISTS hero_quotes JSONB NOT NULL DEFAULT '[]'::jsonb;
