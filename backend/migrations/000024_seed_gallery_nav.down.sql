-- File: 000024_seed_gallery_nav.down.sql
-- Purpose: Remove the seeded Gallery navigation entry.
-- Module: backend/migrations, navigation seed rollback layer.
-- Related: 000024_seed_gallery_nav.up.sql.

DELETE FROM nav_items WHERE key = 'gallery';
