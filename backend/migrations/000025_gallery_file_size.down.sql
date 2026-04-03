-- File: 000025_gallery_file_size.down.sql
-- Purpose: Roll back gallery original file size storage.
-- Module: backend/migrations, schema rollback layer.
-- Related: 000025_gallery_file_size.up.sql.

ALTER TABLE gallery_photos
DROP COLUMN IF EXISTS file_size_bytes;
