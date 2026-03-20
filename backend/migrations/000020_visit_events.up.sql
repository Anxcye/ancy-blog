-- Track raw visitor analytics events for page paths, referrers, and device breakdowns.
-- Keeps each accepted browser event so admin analytics can inspect exact visit records.
CREATE TABLE IF NOT EXISTS visit_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id VARCHAR(128) NOT NULL UNIQUE,
    event_type VARCHAR(32) NOT NULL CHECK (event_type IN ('page_view', 'page_ping')),
    occurred_at TIMESTAMPTZ NOT NULL,
    received_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    visitor_id VARCHAR(128) NOT NULL,
    session_id VARCHAR(128) NOT NULL,
    path VARCHAR(512) NOT NULL,
    route_name VARCHAR(128),
    page_title VARCHAR(256),
    referrer TEXT,
    referrer_host VARCHAR(255),
    content_type VARCHAR(32),
    content_id VARCHAR(128),
    content_slug VARCHAR(255),
    locale VARCHAR(32),
    screen_width INTEGER,
    screen_height INTEGER,
    viewport_width INTEGER,
    viewport_height INTEGER,
    timezone VARCHAR(64),
    ip VARCHAR(64) NOT NULL,
    user_agent TEXT,
    device_type VARCHAR(16) NOT NULL DEFAULT 'unknown',
    browser_name VARCHAR(64),
    os_name VARCHAR(64),
    is_bot BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_visit_events_occurred_at_desc
    ON visit_events (occurred_at DESC);
CREATE INDEX IF NOT EXISTS idx_visit_events_event_type_occurred_at
    ON visit_events (event_type, occurred_at DESC);
CREATE INDEX IF NOT EXISTS idx_visit_events_path_occurred_at
    ON visit_events (path, occurred_at DESC);
CREATE INDEX IF NOT EXISTS idx_visit_events_visitor_occurred_at
    ON visit_events (visitor_id, occurred_at DESC);
CREATE INDEX IF NOT EXISTS idx_visit_events_session_occurred_at
    ON visit_events (session_id, occurred_at DESC);
CREATE INDEX IF NOT EXISTS idx_visit_events_content_occurred_at
    ON visit_events (content_type, content_id, occurred_at DESC);
CREATE INDEX IF NOT EXISTS idx_visit_events_ip_occurred_at
    ON visit_events (ip, occurred_at DESC);
CREATE INDEX IF NOT EXISTS idx_visit_events_referrer_host_occurred_at
    ON visit_events (referrer_host, occurred_at DESC);
