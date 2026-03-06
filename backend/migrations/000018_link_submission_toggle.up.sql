-- File: 000018_link_submission_toggle.up.sql
-- Purpose: Add a global toggle for public friend-link submissions.
-- Module: backend/migrations, schema evolution layer.
-- Related: site_settings and public link submission flow.

ALTER TABLE site_settings
ADD COLUMN IF NOT EXISTS link_submission_enabled BOOLEAN NOT NULL DEFAULT TRUE;
