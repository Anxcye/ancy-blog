-- File: 000006_translation_policy.up.sql
-- Purpose: Add translation_policy JSONB column to site_settings for global auto-translate config.
-- Module: backend/migrations, schema evolution layer.
-- Related: site_repo.go, domain/models.go, handler admin translation policy endpoints.

ALTER TABLE site_settings
    ADD COLUMN IF NOT EXISTS translation_policy JSONB;
