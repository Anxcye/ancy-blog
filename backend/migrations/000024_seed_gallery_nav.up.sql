-- File: 000024_seed_gallery_nav.up.sql
-- Purpose: Seed the public Gallery navigation entry so the gallery page is discoverable from the blog header.
-- Module: backend/migrations, navigation seed data layer.
-- Related: nav_items, frontend-blog default layout navigation, gallery public routes.

INSERT INTO nav_items (name, key, type, target_type, target_value, order_num, enabled)
SELECT '画廊', 'gallery', 'link', 'route', '/gallery', 25, true
WHERE NOT EXISTS (SELECT 1 FROM nav_items WHERE key = 'gallery');
