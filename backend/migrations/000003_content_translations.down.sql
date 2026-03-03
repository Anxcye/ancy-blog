-- File: 000003_content_translations.down.sql
-- Purpose: Roll back locale translation tables.
-- Module: backend/migrations, schema rollback layer.
-- Related: 000003_content_translations.up.sql.

DROP TABLE IF EXISTS moment_translations;
DROP TABLE IF EXISTS article_translations;
