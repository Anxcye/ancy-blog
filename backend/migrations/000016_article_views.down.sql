DROP INDEX IF EXISTS idx_article_views_article_id;
DROP TABLE  IF EXISTS article_views;
ALTER TABLE articles DROP COLUMN IF EXISTS view_count;
