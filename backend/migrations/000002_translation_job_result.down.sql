-- File: 000002_translation_job_result.down.sql
-- Purpose: Roll back translation result payload column.
-- Module: backend/migrations, schema rollback layer.
-- Related: 000002_translation_job_result.up.sql.

ALTER TABLE translation_jobs
DROP COLUMN IF EXISTS result_text;
