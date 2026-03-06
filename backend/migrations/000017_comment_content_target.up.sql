-- File: 000017_comment_content_target.up.sql
-- Purpose: Expand comments to support both articles and moments through polymorphic content targets.
-- Module: backend/migrations, schema migration layer.
-- Related: comments table, public comment APIs, and moments page interactions.
ALTER TABLE comments ADD COLUMN IF NOT EXISTS content_type VARCHAR(16);
ALTER TABLE comments ADD COLUMN IF NOT EXISTS content_id UUID;

UPDATE comments
SET content_type = 'article',
    content_id = article_id
WHERE content_type IS NULL OR content_id IS NULL;

ALTER TABLE comments
    ALTER COLUMN article_id DROP NOT NULL,
    ALTER COLUMN content_type SET NOT NULL,
    ALTER COLUMN content_id SET NOT NULL;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'chk_comments_content_type'
    ) THEN
        ALTER TABLE comments
            ADD CONSTRAINT chk_comments_content_type CHECK (content_type IN ('article', 'moment'));
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_comments_content_status_created_at
    ON comments (content_type, content_id, status, created_at DESC);
