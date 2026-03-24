DROP INDEX IF EXISTS idx_visit_events_page_view_session_path_occurred_at;

ALTER TABLE visit_events
    DROP COLUMN IF EXISTS active_duration_seconds,
    DROP COLUMN IF EXISTS last_engaged_at;
