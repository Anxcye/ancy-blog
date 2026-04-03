-- File: 000023_gallery.up.sql
-- Purpose: Create gallery module tables: photos, gallery tags, photo-tag junction, and photo assets.
-- Module: backend/migrations, schema evolution layer.
-- Related: gallery domain, gallery_repo.go, gallery_service.go, gallery handlers.

-- =========================
-- Gallery Tags (separate from article tags for independent taxonomy)
-- =========================

CREATE TABLE IF NOT EXISTS gallery_tags (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(128) NOT NULL UNIQUE,
    slug        VARCHAR(128) NOT NULL UNIQUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_gallery_tags_slug ON gallery_tags (slug) WHERE deleted_at IS NULL;

-- =========================
-- Gallery Photos
-- =========================

CREATE TABLE IF NOT EXISTS gallery_photos (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title           VARCHAR(256) NOT NULL DEFAULT '',
    slug            VARCHAR(256) NOT NULL UNIQUE,
    description     TEXT NOT NULL DEFAULT '',
    status          VARCHAR(16)  NOT NULL DEFAULT 'draft'
                        CHECK (status IN ('draft', 'published', 'hidden')),

    -- Location (city-level, no raw GPS stored)
    location_name   VARCHAR(256) NOT NULL DEFAULT '',
    location_city   VARCHAR(128) NOT NULL DEFAULT '',
    location_state  VARCHAR(128) NOT NULL DEFAULT '',
    location_country VARCHAR(128) NOT NULL DEFAULT '',

    -- Shooting time
    taken_at        TIMESTAMPTZ,

    -- EXIF metadata (whitelisted fields only)
    camera_make     VARCHAR(128) NOT NULL DEFAULT '',
    camera_model    VARCHAR(128) NOT NULL DEFAULT '',
    lens_model      VARCHAR(256) NOT NULL DEFAULT '',
    focal_length    VARCHAR(32)  NOT NULL DEFAULT '',
    aperture        VARCHAR(32)  NOT NULL DEFAULT '',
    shutter_speed   VARCHAR(32)  NOT NULL DEFAULT '',
    iso             VARCHAR(16)  NOT NULL DEFAULT '',

    -- Image dimensions
    width           INTEGER NOT NULL DEFAULT 0,
    height          INTEGER NOT NULL DEFAULT 0,

    -- Per-photo metadata display switches
    taken_at_display  BOOLEAN NOT NULL DEFAULT TRUE,
    camera_display    BOOLEAN NOT NULL DEFAULT TRUE,
    location_display  BOOLEAN NOT NULL DEFAULT TRUE,
    exif_display      BOOLEAN NOT NULL DEFAULT TRUE,
    tags_display      BOOLEAN NOT NULL DEFAULT TRUE,

    -- Derived asset URLs
    placeholder_data  TEXT    NOT NULL DEFAULT '',
    display_url       TEXT    NOT NULL DEFAULT '',
    large_url         TEXT    NOT NULL DEFAULT '',

    -- Processing status
    processing_status VARCHAR(32) NOT NULL DEFAULT 'pending'
                        CHECK (processing_status IN ('pending', 'processing', 'completed', 'failed')),
    processing_error  TEXT    NOT NULL DEFAULT '',

    -- Ordering
    sort_order        INTEGER NOT NULL DEFAULT 0,

    -- Article reference count (denormalized for deletion safety)
    article_ref_count INTEGER NOT NULL DEFAULT 0,

    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_gallery_photos_status_sort
    ON gallery_photos (status, sort_order DESC, taken_at DESC NULLS LAST, created_at DESC)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_gallery_photos_slug
    ON gallery_photos (slug) WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_gallery_photos_processing_status
    ON gallery_photos (processing_status) WHERE deleted_at IS NULL;

-- =========================
-- Gallery Photo Tags (many-to-many)
-- =========================

CREATE TABLE IF NOT EXISTS gallery_photo_tags (
    photo_id    UUID NOT NULL REFERENCES gallery_photos(id) ON DELETE CASCADE,
    tag_id      UUID NOT NULL REFERENCES gallery_tags(id) ON DELETE CASCADE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (photo_id, tag_id)
);

CREATE INDEX IF NOT EXISTS idx_gallery_photo_tags_tag_photo
    ON gallery_photo_tags (tag_id, photo_id);
