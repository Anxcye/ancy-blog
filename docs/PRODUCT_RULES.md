# Product Rules (v1)

## Purpose
Define business behavior and content policy for the new blog platform.  
Technical schema details belong to `docs/DATA_MODEL.md`.

## Identifier and URL Strategy
- Internal primary keys use UUID.
- Public article URLs use `slug`, not UUID.
- URL format: `/article/{slug}`.
- Slug must be unique, lowercase, and hyphen-separated words.
- Locale route strategy:
  - Chinese (default): no locale prefix (e.g. `/article/{slug}`).
  - English: `/en/*` prefix (e.g. `/en/article/{slug}`).
  - No `/zh` prefix in default route set.

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

## Integration Configuration Rules
- External integrations must be managed in a unified admin configuration center.
- R2 image storage and LLM provider config should share one management domain and one UI module.
- Integration config fields are editable in admin UI and stored in DB (not hardcoded).
- Integration config updates must take effect immediately for runtime actions (e.g., image upload) without requiring service restart.
- Secret fields (API keys/tokens) are write-only in UI and must be masked in read responses.
- Integrations must support:
  - `enabled` switch
  - provider type
  - connection settings
  - health-check status metadata
- Initial integration scopes:
  - Object storage (`cloudflare_r2`)
  - LLM translation (`openai_compatible`)

## AI Translation Rules
- Translation should be asynchronous (job-based), not blocking publish flow.
- Translation task uses configured LLM provider from integration center.
- Translation output should be stored in locale-specific content records.
- Translation result must include metadata: provider, model, status, and error message (if failed).
- Translation scope:
  - Article: `title`, `summary`, `content`.
  - Moment: `content`.
- Locale fallback:
  - English route falls back to Chinese content when localized content is unavailable.
  - Fallback pages should be marked `noindex` on frontend to avoid SEO confusion.
- Publish control:
  - Chinese source content is the default publish baseline.
  - Translation job supports `auto_publish` switch.
  - Translation job supports `publish_at` for scheduled publish.
  - When `auto_publish=false`, localized content remains draft until manual publish.

## Editor Rules
- Admin article editor uses rich-text JSON as persisted content format.
- Editor image insertion must upload through configured object storage integration (`cloudflare_r2`) and store URL references in content JSON.

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
  - integration public capabilities (optional)
- Key examples:
  - `site:settings:{locale}`
  - `site:footer:{locale}`
  - `site:social:{locale}`
  - `site:nav:{locale}`
  - `site:slot:{slot_key}:{locale}`
  - `integration:capabilities`
- Pattern: cache-aside.
  - Public read APIs: read cache first, fallback to DB, then set cache.
  - Admin update APIs: write DB first, then invalidate related cache keys.
- Redis is acceleration only; database remains source of truth.
