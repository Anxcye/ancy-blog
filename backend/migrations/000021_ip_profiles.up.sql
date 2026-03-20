-- File: 000021_ip_profiles.up.sql
-- Purpose: Cache offline ip2region lookups so analytics queries can filter by geographic fields.
-- Module: backend/migrations, analytics persistence change set.
-- Related: visit_events analytics queries, ip2region resolver, and admin analytics filters.
CREATE TABLE IF NOT EXISTS ip_profiles (
    ip VARCHAR(64) PRIMARY KEY,
    country_code VARCHAR(16),
    country_name VARCHAR(128),
    region_name VARCHAR(128),
    city_name VARCHAR(128),
    isp VARCHAR(128),
    raw_region VARCHAR(255),
    source VARCHAR(32) NOT NULL DEFAULT 'ip2region',
    resolved_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_ip_profiles_country_name
    ON ip_profiles (country_name);
CREATE INDEX IF NOT EXISTS idx_ip_profiles_region_name
    ON ip_profiles (region_name);
CREATE INDEX IF NOT EXISTS idx_ip_profiles_city_name
    ON ip_profiles (city_name);
CREATE INDEX IF NOT EXISTS idx_ip_profiles_isp
    ON ip_profiles (isp);
