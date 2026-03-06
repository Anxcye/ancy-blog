-- Track unique page views per article per visitor (fingerprinted by IP + UA hash).
-- A single row per (article_id, visitor_key) so repeat visits within a day are deduplicated.
CREATE TABLE IF NOT EXISTS article_views (
    id          UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    article_id  UUID         NOT NULL REFERENCES articles(id) ON DELETE CASCADE,
    visitor_key VARCHAR(128) NOT NULL,  -- SHA-256(ip + ua + date), rotates daily
    viewed_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_article_views UNIQUE (article_id, visitor_key)
);

CREATE INDEX IF NOT EXISTS idx_article_views_article_id ON article_views (article_id);

-- Denormalised counter on articles for fast list queries.
ALTER TABLE articles ADD COLUMN IF NOT EXISTS view_count BIGINT NOT NULL DEFAULT 0;
