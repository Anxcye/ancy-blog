INSERT INTO nav_items (name, key, type, target_type, target_value, order_num, enabled)
SELECT '友人帐', 'links', 'link', 'route', '/friends', 50, true
WHERE NOT EXISTS (SELECT 1 FROM nav_items WHERE key = 'links');
