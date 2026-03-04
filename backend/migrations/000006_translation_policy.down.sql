-- File: 000006_translation_policy.down.sql
-- Purpose: Revert translation_policy column from site_settings.
-- Module: backend/migrations, schema evolution layer.
-- Related: 000006_translation_policy.up.sql

ALTER TABLE site_settings
    DROP COLUMN IF EXISTS translation_policy;
