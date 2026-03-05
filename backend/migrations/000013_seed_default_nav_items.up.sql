INSERT INTO nav_items (name, key, type, target_type, target_value, order_num, enabled)
SELECT '首页', 'home', 'link', 'route', '/', 10, true
WHERE NOT EXISTS (SELECT 1 FROM nav_items WHERE key = 'home');

INSERT INTO nav_items (name, key, type, target_type, target_value, order_num, enabled)
SELECT '文章', 'articles', 'link', 'route', '/articles', 20, true
WHERE NOT EXISTS (SELECT 1 FROM nav_items WHERE key = 'articles');

INSERT INTO nav_items (name, key, type, target_type, target_value, order_num, enabled)
SELECT '说说', 'moments', 'link', 'route', '/moments', 30, true
WHERE NOT EXISTS (SELECT 1 FROM nav_items WHERE key = 'moments');

INSERT INTO nav_items (name, key, type, target_type, target_value, order_num, enabled)
SELECT '时间线', 'timeline', 'link', 'route', '/timeline', 40, true
WHERE NOT EXISTS (SELECT 1 FROM nav_items WHERE key = 'timeline');
