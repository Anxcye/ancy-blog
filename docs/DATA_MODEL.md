# Data Model (v1)

## Scope
This document defines database table design, field semantics, indexes, and constraints.  
API request/response shapes remain in `docs/API_CONTRACT.md`.

## Conventions
- Primary key: `id` (UUID across all business tables).
- Timestamp fields: `created_at`, `updated_at` (`TIMESTAMPTZ`).
- Soft delete: use `deleted_at` where needed.
- Enum fields: use `VARCHAR` + check constraint (or PostgreSQL enum type).

## Table: `articles`
Purpose: Core blog content entity.

Fields:
- `id`
- `title` (required)
- `slug` (required, unique)
- `content_kind` (`post | page`)
- `summary`
- `content` (Markdown source)
- `status` (`draft | published | archived`)
- `visibility` (`public | unlisted | private`)
- `is_pinned` (bool, default false)
- `pin_order` (int, default 0)
- `pinned_at` (nullable)
- `allow_comment` (bool, default true)
- `origin_type` (`original | repost | translation`)
- `source_url` (nullable)
- `ai_assist_level` (`none | polish | dictation | assisted | generated | translated`)
- `cover_image` (nullable)
- `published_at` (nullable)
- `created_at`
- `updated_at`
- `deleted_at` (nullable)

Indexes/Constraints:
- Unique: `slug`
- Index: `(content_kind, status, published_at DESC)`
- Index: `(status, published_at DESC)`
- Index: `(is_pinned DESC, pin_order DESC, published_at DESC)`

## Table: `categories`
Purpose: Article category taxonomy.

Fields:
- `id`
- `name` (required, unique)
- `slug` (required, unique)
- `description` (nullable)
- `created_at`
- `updated_at`
- `deleted_at` (nullable)

Indexes/Constraints:
- Unique: `name`
- Unique: `slug`

## Table: `tags`
Purpose: Article tag taxonomy.

Fields:
- `id`
- `name` (required, unique)
- `slug` (required, unique)
- `description` (nullable)
- `created_at`
- `updated_at`
- `deleted_at` (nullable)

Indexes/Constraints:
- Unique: `name`
- Unique: `slug`

## Table: `article_tags`
Purpose: Many-to-many relation between articles and tags.

Fields:
- `article_id` (required, FK to `articles.id`)
- `tag_id` (required, FK to `tags.id`)
- `created_at`

Indexes/Constraints:
- Unique: `(article_id, tag_id)`
- Index: `(tag_id, article_id)`

## Table: `comments`
Purpose: Threaded comments for articles (issue-like discussion style).

Fields:
- `id`
- `article_id` (required, FK to `articles.id`)
- `parent_id` (nullable, FK to `comments.id`)
- `root_id` (nullable, FK to `comments.id`)
- `content` (required)
- `status` (`pending | approved | rejected | spam | deleted`)
- `is_pinned` (bool, default false)
- `like_count` (int, default 0)
- `reply_count` (int, default 0)
- `nickname` (required)
- `email` (nullable, admin-visible)
- `website` (nullable)
- `avatar_url` (nullable)
- `source` (`web | admin | api`)
- `ip` (required, plaintext storage by project decision)
- `user_agent` (nullable)
- `risk_score` (int, default 0)
- `approved_at` (nullable)
- `approved_by` (nullable, admin user id)
- `created_at`
- `updated_at`
- `deleted_at` (nullable)

Indexes/Constraints:
- Index: `(article_id, status, created_at DESC)`
- Index: `(parent_id, created_at ASC)`
- Index: `(root_id, created_at ASC)`
- Index: `(ip, created_at DESC)`
- Check: `parent_id != id`

## Table: `moments`
Purpose: Short-form status/log posts independent from long-form articles.

Fields:
- `id`
- `content` (required, short text/markdown)
- `status` (`draft | published | archived`)
- `is_pinned` (bool, default false)
- `pin_order` (int, default 0)
- `allow_comment` (bool, default true)
- `source` (`web | admin | api`)
- `published_at` (nullable)
- `created_at`
- `updated_at`
- `deleted_at` (nullable)

Indexes/Constraints:
- Index: `(status, published_at DESC)`
- Index: `(is_pinned DESC, pin_order DESC, published_at DESC)`

## Table: `links`
Purpose: Friend link management with moderation workflow.

Fields:
- `id`
- `name` (required)
- `url` (required)
- `avatar_url` (nullable)
- `description` (nullable)
- `contact_email` (nullable)
- `review_status` (`pending | approved | rejected`)
- `review_note` (nullable)
- `submitted_ip` (nullable, plaintext)
- `submitted_user_agent` (nullable)
- `related_article_id` (nullable, FK to `articles.id`)
- `approved_at` (nullable)
- `approved_by` (nullable, admin user id)
- `created_at`
- `updated_at`
- `deleted_at` (nullable)

Indexes/Constraints:
- Unique: `url`
- Index: `(review_status, created_at DESC)`
- Index: `(related_article_id)`

## Table: `site_settings`
Purpose: Global singleton-like key/value settings for site rendering.

Fields:
- `id`
- `site_name` (required)
- `avatar_url` (nullable)
- `hero_intro_md` (nullable, short markdown text for homepage center intro)
- `default_locale` (required, default `en`)
- `created_at`
- `updated_at`

Indexes/Constraints:
- Unique: `id` (single active row policy can be enforced in service layer)

## Table: `nav_items`
Purpose: Dynamic top navigation configuration.

Fields:
- `id`
- `name` (required)
- `key` (required, unique)
- `type` (`menu | dropdown | link`)
- `target_type` (`route | category | slot | external`)
- `target_value` (nullable)
- `order_num` (required, default 0)
- `enabled` (bool, default true)
- `created_at`
- `updated_at`
- `deleted_at` (nullable)

Indexes/Constraints:
- Unique: `key`
- Index: `(enabled, order_num ASC)`

## Table: `content_slots`
Purpose: Define reusable slot positions for content mounting.

Fields:
- `id`
- `slot_key` (required, unique)
- `name` (required)
- `description` (nullable)
- `enabled` (bool, default true)
- `created_at`
- `updated_at`

Indexes/Constraints:
- Unique: `slot_key`

## Table: `content_slot_items`
Purpose: Many-to-many mapping between slots and content.

Fields:
- `id`
- `slot_key` (required, FK to `content_slots.slot_key`)
- `content_type` (`article | moment`)
- `content_id` (required)
- `order_num` (required, default 0)
- `enabled` (bool, default true)
- `created_at`
- `updated_at`

Indexes/Constraints:
- Unique: `(slot_key, content_type, content_id)`
- Index: `(slot_key, enabled, order_num ASC)`

## Table: `footer_items`
Purpose: Configurable footer display items with row placement and ordering.

Fields:
- `id`
- `label` (required)
- `link_type` (`none | internal | external`)
- `internal_article_slug` (nullable, for internal links to `articles.slug` where `content_kind=page`)
- `external_url` (nullable, for external links)
- `row_num` (required, 1..3)
- `order_num` (required, default 0)
- `enabled` (bool, default true)
- `created_at`
- `updated_at`
- `deleted_at` (nullable)

Indexes/Constraints:
- Index: `(row_num, order_num ASC)`
- Index: `(enabled, row_num, order_num ASC)`
- Check: `row_num BETWEEN 1 AND 3`
- Check: `link_type='internal' -> internal_article_slug IS NOT NULL`
- Check: `link_type='external' -> external_url IS NOT NULL`

## Table: `social_links`
Purpose: Managed homepage social links shown in center section.

Fields:
- `id`
- `platform` (`github | mail | x | linkedin | custom`)
- `title` (required)
- `url` (required)
- `icon_key` (nullable)
- `order_num` (required, default 0)
- `enabled` (bool, default true)
- `created_at`
- `updated_at`
- `deleted_at` (nullable)

Indexes/Constraints:
- Index: `(enabled, order_num ASC)`

## Table: `integration_providers`
Purpose: Unified provider configuration center for external integrations (R2, LLM, etc.).

Fields:
- `id`
- `provider_type` (`object_storage | llm`)
- `provider_key` (required, unique, e.g. `cloudflare_r2`, `openai_compatible`)
- `name` (required)
- `enabled` (bool, default false)
- `config_json` (required, encrypted/masked for secret fields)
- `meta_json` (nullable, for health check info and non-secret runtime metadata)
- `created_at`
- `updated_at`

Indexes/Constraints:
- Unique: `provider_key`
- Index: `(provider_type, enabled)`

## Table: `translation_jobs`
Purpose: Track async translation workflow for articles/moments/pages.

Fields:
- `id`
- `source_type` (`article | moment`)
- `source_id` (required)
- `source_locale` (required)
- `target_locale` (required)
- `provider_key` (required, FK to `integration_providers.provider_key`)
- `model_name` (required)
- `status` (`queued | running | succeeded | failed`)
- `error_message` (nullable)
- `result_text` (nullable, translated output from worker)
- `requested_by` (nullable, admin user id)
- `retry_count` (required, default 0)
- `max_retries` (required, default 3)
- `next_retry_at` (required, default now)
- `auto_publish` (required, default false)
- `publish_at` (nullable, scheduled publish time for localized content)
- `created_at`
- `updated_at`
- `finished_at` (nullable)

Indexes/Constraints:
- Index: `(source_type, source_id, target_locale)`
- Index: `(status, created_at DESC)`
- Index: `(status, next_retry_at ASC, created_at ASC)`

## Table: `article_translations`
Purpose: Store locale-specific article content generated by translation worker or manual edits.

Fields:
- `id`
- `article_id` (required, FK to `articles.id`)
- `locale` (required)
- `title` (nullable, translated title)
- `summary` (nullable, translated summary)
- `content` (required, translated body)
- `status` (`draft | published | archived`)
- `published_at` (nullable)
- `translated_by_job_id` (nullable, FK to `translation_jobs.id`)
- `created_at`
- `updated_at`

Indexes/Constraints:
- Unique: `(article_id, locale)`
- Index: `(locale, updated_at DESC)`
- Index: `(status, published_at DESC, updated_at DESC)`

## Table: `moment_translations`
Purpose: Store locale-specific moment content generated by translation worker or manual edits.

Fields:
- `id`
- `moment_id` (required, FK to `moments.id`)
- `locale` (required)
- `content` (required, translated body)
- `status` (`draft | published | archived`)
- `published_at` (nullable)
- `translated_by_job_id` (nullable, FK to `translation_jobs.id`)
- `created_at`
- `updated_at`

Indexes/Constraints:
- Unique: `(moment_id, locale)`
- Index: `(locale, updated_at DESC)`
- Index: `(status, published_at DESC, updated_at DESC)`

## Table: `reactions`
Purpose: Generic reaction records for article/comment (like/upvote/love/etc.).

Fields:
- `id`
- `target_type` (`article | comment`)
- `target_id` (required)
- `reaction_type` (`like | love | clap | insightful`)
- `actor_type` (`admin | visitor`)
- `actor_id` (nullable, for admin/user account)
- `visitor_key` (nullable, for anonymous visitor fingerprint/cookie key)
- `ip` (nullable, plaintext for anti-abuse checks)
- `user_agent` (nullable)
- `created_at`

Indexes/Constraints:
- Unique: `(target_type, target_id, reaction_type, actor_type, actor_id, visitor_key)`
- Index: `(target_type, target_id, reaction_type)`
- Check: at least one of `actor_id` or `visitor_key` is not null.

## Notes
- `ip` is intentionally stored in plaintext for moderation/abuse tracing (project decision).
- If privacy requirements change later, move to hashed or masked storage and add retention policy.
