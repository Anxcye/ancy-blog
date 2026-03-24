-- Add engagement fields so page_ping can update page_view duration instead of inserting new rows.
-- Related: visit_events analytics ingest, admin analytics duration display, and storage-growth control.
ALTER TABLE visit_events
    ADD COLUMN IF NOT EXISTS last_engaged_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS active_duration_seconds INTEGER NOT NULL DEFAULT 0;

UPDATE visit_events
SET last_engaged_at = COALESCE(last_engaged_at, occurred_at),
    active_duration_seconds = COALESCE(active_duration_seconds, 0);

CREATE INDEX IF NOT EXISTS idx_visit_events_page_view_session_path_occurred_at
    ON visit_events (session_id, path, occurred_at DESC)
    WHERE event_type = 'page_view';

DROP INDEX IF EXISTS idx_visit_events_ping_session_path_occurred_at;
