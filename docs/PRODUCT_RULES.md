# Product Rules (v1)

## Purpose
Define business behavior and content policy for the new blog platform.  
Technical schema details belong to `docs/DATA_MODEL.md`.

## Identifier and URL Strategy
- Internal primary keys use UUID.
- Public article URLs use `slug`, not UUID.
- URL format: `/article/{slug}`.
- Slug must be unique, lowercase, and hyphen-separated words.

## Article Publishing Rules
- Articles support `draft`, `published`, `archived`.
- Articles are unified content records with `content_kind`:
  - `post` for normal blog posts
  - `page` for internal site pages (e.g. About Me, About Site)
- Per-article comment switch is required: `allow_comment`.
- Content placement is configurable through slot relations (not a single article field).
- Recommended slot keys:
  - `home_about`
  - `home_featured`
  - `home_top_story`
  - `nav_home_hover`
  - `nav_articles_hover`

## AI Disclosure Rules
- Every article must include an AI disclosure level:
  - `none`: no AI involved
  - `polish`: AI used for wording/polishing
  - `dictation`: AI used to convert speech/dictation to draft
  - `assisted`: AI helped with structure/partial generation
  - `generated`: AI generated most content
  - `translated`: AI used for translation

## Content Types
- Long-form content: `articles`.
- Short updates/status posts: `moments` (independent content type).

## Comment Rules
- Comments support threaded replies.
- Store `ip` in plaintext by project decision.
- Record `user_agent` for moderation and abuse analysis.
- Comment reactions are tracked via generic `reactions` records.

## Link (Friend Link) Rules
- Visitors can submit link applications.
- Admin reviews each submission (`pending/approved/rejected`).
- Approved links can optionally associate one article (`related_article_id`).
- Suggested submission fields:
  - site name
  - site URL
  - avatar/logo URL
  - short description
  - contact email

## Taxonomy Rules
- Articles can bind one category and multiple tags.
- Categories and tags are independent managed entities.

## Site Configuration Rules
- Site-level settings must support:
  - `site_name`
  - `avatar_url`
  - `hero_intro_md` (short intro on homepage center, Markdown enabled)
- Footer content is fully configurable by admin.
  - Supports plain text (`link_type=none`)
  - Supports internal links (`link_type=internal`, target article by slug)
  - Supports external links (`link_type=external`)
  - Supports ordering and row placement
  - Default UI layout uses 3 rows (`row_num` from 1 to 3)
- Homepage social links (e.g. GitHub, email) are managed items with ordering and enable flags.

## Navigation Rules
- Top navigation is dynamically managed by admin.
- Default menu items:
  - Home
  - Articles
  - Moments
  - Timeline
  - Links
- Navigation supports:
  - static route links
  - external links
  - dropdown by category
  - dropdown by content slot

## Internal Pages Rules
- Internal static-like pages (e.g. "About Me", "About This Site") are stored in `articles` with `content_kind=page`.
- Internal link items should reference article slug, not hardcoded URL.

## Caching Rules (Redis)
- Cache read-heavy site data:
  - site settings
  - footer items
  - social links
  - top navigation
  - homepage/hover slot content
- Key examples:
  - `site:settings:{locale}`
  - `site:footer:{locale}`
  - `site:social:{locale}`
  - `site:nav:{locale}`
  - `site:slot:{slot_key}:{locale}`
- Pattern: cache-aside.
  - Public read APIs: read cache first, fallback to DB, then set cache.
  - Admin update APIs: write DB first, then invalidate related cache keys.
- Redis is acceleration only; database remains source of truth.
