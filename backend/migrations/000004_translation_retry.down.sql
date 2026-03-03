-- File: 000004_translation_retry.down.sql
-- Purpose: Remove retry scheduling columns from translation jobs.
-- Module: backend/migrations, schema rollback layer.
-- Related: 000004_translation_retry.up.sql.

DROP INDEX IF EXISTS idx_translation_jobs_status_next_retry;

ALTER TABLE translation_jobs
    DROP COLUMN IF EXISTS next_retry_at,
    DROP COLUMN IF EXISTS max_retries,
    DROP COLUMN IF EXISTS retry_count;
