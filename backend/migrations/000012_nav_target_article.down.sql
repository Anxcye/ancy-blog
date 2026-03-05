ALTER TABLE nav_items DROP CONSTRAINT IF EXISTS nav_items_target_type_check;
ALTER TABLE nav_items ADD CONSTRAINT nav_items_target_type_check CHECK (target_type IN ('route', 'category', 'slot', 'external'));
