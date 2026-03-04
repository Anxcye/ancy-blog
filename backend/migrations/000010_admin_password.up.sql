-- File: 000010_admin_password.up.sql
-- Purpose: Add admin_password_hash column to site_settings for persistent credential storage.
-- Module: backend/migrations, schema evolution layer.
-- Related: auth service, site_repo.go.

ALTER TABLE site_settings
    ADD COLUMN IF NOT EXISTS admin_password_hash TEXT;
