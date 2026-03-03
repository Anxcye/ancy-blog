-- File: 000005_i18n_publish_control.down.sql
-- Purpose: Roll back translation publish control and extra translation fields.
-- Module: backend/migrations, schema rollback layer.
-- Related: 000005_i18n_publish_control.up.sql.

DROP INDEX IF EXISTS idx_moment_translations_status_published;
DROP INDEX IF EXISTS idx_article_translations_status_published;

ALTER TABLE moment_translations
    DROP COLUMN IF EXISTS published_at,
    DROP COLUMN IF EXISTS status;

ALTER TABLE article_translations
    DROP COLUMN IF EXISTS published_at,
    DROP COLUMN IF EXISTS status,
    DROP COLUMN IF EXISTS summary,
    DROP COLUMN IF EXISTS title;

ALTER TABLE translation_jobs
    DROP COLUMN IF EXISTS publish_at,
    DROP COLUMN IF EXISTS auto_publish;
