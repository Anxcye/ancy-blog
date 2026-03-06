-- File: 000017_comment_content_target.down.sql
-- Purpose: Revert comments polymorphic targets back to article-only linkage.
-- Module: backend/migrations, schema migration layer.
-- Related: comments table and public comment APIs.
DELETE FROM comments WHERE content_type = 'moment';

UPDATE comments
SET article_id = content_id
WHERE content_type = 'article' AND article_id IS NULL;

DROP INDEX IF EXISTS idx_comments_content_status_created_at;

ALTER TABLE comments DROP CONSTRAINT IF EXISTS chk_comments_content_type;
ALTER TABLE comments DROP COLUMN IF EXISTS content_id;
ALTER TABLE comments DROP COLUMN IF EXISTS content_type;
ALTER TABLE comments ALTER COLUMN article_id SET NOT NULL;
