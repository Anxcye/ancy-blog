-- File: 000001_init.down.sql
-- Purpose: Roll back the initial schema for local/dev reset.
-- Module: backend/migrations, migration rollback layer.
-- Related: 000001_init.up.sql.

DROP TABLE IF EXISTS translation_jobs;
DROP TABLE IF EXISTS integration_providers;
DROP TABLE IF EXISTS social_links;
DROP TABLE IF EXISTS footer_items;
DROP TABLE IF EXISTS content_slot_items;
DROP TABLE IF EXISTS content_slots;
DROP TABLE IF EXISTS nav_items;
DROP TABLE IF EXISTS site_settings;
DROP TABLE IF EXISTS links;
DROP TABLE IF EXISTS reactions;
DROP TABLE IF EXISTS comments;
DROP TABLE IF EXISTS moments;
DROP TABLE IF EXISTS article_tags;
DROP TABLE IF EXISTS tags;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS articles;
