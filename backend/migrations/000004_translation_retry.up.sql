-- File: 000004_translation_retry.up.sql
-- Purpose: Add retry scheduling columns for translation jobs.
-- Module: backend/migrations, schema evolution layer.
-- Related: translation worker retry/backoff and manual retry APIs.

ALTER TABLE translation_jobs
    ADD COLUMN IF NOT EXISTS retry_count INTEGER NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS max_retries INTEGER NOT NULL DEFAULT 3,
    ADD COLUMN IF NOT EXISTS next_retry_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

CREATE INDEX IF NOT EXISTS idx_translation_jobs_status_next_retry
    ON translation_jobs (status, next_retry_at ASC, created_at ASC);
