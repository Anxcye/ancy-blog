-- File: 000025_gallery_file_size.up.sql
-- Purpose: Store original uploaded photo file size for gallery metadata display.
-- Module: backend/migrations, schema evolution layer.
-- Related: gallery_photos, gallery_repo.go, public gallery viewer metadata.

ALTER TABLE gallery_photos
ADD COLUMN IF NOT EXISTS file_size_bytes BIGINT NOT NULL DEFAULT 0;
