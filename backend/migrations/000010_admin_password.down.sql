-- File: 000010_admin_password.down.sql
-- Purpose: Revert admin_password_hash column from site_settings.
-- Module: backend/migrations, schema evolution layer.
-- Related: 000010_admin_password.up.sql

ALTER TABLE site_settings
    DROP COLUMN IF EXISTS admin_password_hash;
