-- Migration 000011: Add parent_id to nav_items for hierarchical (dropdown) navigation support.
ALTER TABLE nav_items ADD COLUMN IF NOT EXISTS parent_id UUID REFERENCES nav_items(id) ON DELETE CASCADE;
