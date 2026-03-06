-- File: 000018_link_submission_toggle.down.sql
-- Purpose: Remove the global toggle for public friend-link submissions.
-- Module: backend/migrations, schema rollback layer.
-- Related: 000018_link_submission_toggle.up.sql and site_settings.

ALTER TABLE site_settings
DROP COLUMN IF EXISTS link_submission_enabled;
