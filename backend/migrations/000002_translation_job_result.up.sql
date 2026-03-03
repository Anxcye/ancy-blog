-- File: 000002_translation_job_result.up.sql
-- Purpose: Add translation result payload storage to translation jobs.
-- Module: backend/migrations, schema evolution layer.
-- Related: 000001_init.up.sql and translation worker runtime.

ALTER TABLE translation_jobs
ADD COLUMN IF NOT EXISTS result_text TEXT;
