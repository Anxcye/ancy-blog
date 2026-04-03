-- File: 000023_gallery.down.sql
-- Purpose: Rollback gallery module tables.
-- Module: backend/migrations, schema evolution layer.
-- Related: 000023_gallery.up.sql

DROP TABLE IF EXISTS gallery_photo_tags;
DROP TABLE IF EXISTS gallery_photos;
DROP TABLE IF EXISTS gallery_tags;
