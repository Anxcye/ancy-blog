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
- Admin article editing must allow authors to set and update the AI disclosure level explicitly.
- Public article detail metadata should display publish time as the primary timestamp.
- Public article detail should only show update time when the article has materially changed after publication.
- Public article detail should expose AI disclosure with a hover explanation for the selected disclosure level.

## Content Types
- Long-form content: `articles`.
- Short updates/status posts: `moments` (independent content type).

## Comment Rules
- Comments support threaded replies.
- Articles and moments both support the same public comment system.
- Public pages should lazy-load per-item comment threads for moments instead of expanding all comment sections by default.
- Public article comment lists must return recursive reply threads instead of flat reply pages.
- Public comment payloads must exclude visitor privacy fields such as `email`, `ip`, and `user_agent`.
- Comments support admin pinning on root comments and replies.
- Admins can reply from the moderation center; admin replies are stored with `source=admin` and are immediately approved.
- Public comment items with `isAuthor=true` should display the site-level avatar and site name instead of the raw stored commenter nickname/avatar.
- Store `ip` in plaintext by project decision.
- Record `user_agent` for moderation and abuse analysis.
- Comment reactions are tracked via generic `reactions` records.

## Link (Friend Link) Rules
- Visitors can submit link applications.
- Site settings must support a global `link_submission_enabled` toggle for public friend-link applications.
- Admin reviews each submission (`pending/approved/rejected`).
- Approved links can optionally associate one article (`related_article_id`).
- Public friends page should prioritize editorial context and the approved directory before the submission form.
- Public friend-link submission form should default to collapsed and expand on explicit user intent.
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
  - `favicon_url`
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
- Admin translation management must expose translated rows independently from job history and provide source context (`source title/slug`) for review.
- Admin translation editing must load the full persisted translation detail before saving overrides, instead of relying on truncated list payloads.
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
- Editor image insertion must support both direct remote image URLs and uploaded image assets, storing final URL references in content JSON.
- Admin article editor must support table authoring with persisted TipTap table nodes that can also render on the public blog.
- Editor must support extensible block embeds with explicit node types and attrs.
  - Initial custom blocks: `x_post` and `tmdb_card`.
  - Embed blocks should be stored in the same rich-text JSON document as normal content nodes.

## Navigation Rules
- Top navigation is dynamically managed by admin.
- Default menu items:
  - Home
  - Articles
  - Gallery
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
- When an internal page is backed by an article record that allows comments, the public page may reuse the shared article comment system.

## Moments Interaction Rules
- Timeline should behave as a month-grouped archive index:
  - newest month expanded by default
  - older months collapsed by default
  - article rows navigate to article detail pages
  - moment rows open the shared moment detail modal
  - compact row layout should prefer `date + text + tail label/category` over card styling
- Public `moments` feed is for browsing only; detailed interaction happens in a route-driven modal (`/moments/:id`).
- Each moment detail modal must support previous/next navigation when adjacent moments are available in the loaded feed.
- Closing the moment detail modal must preserve the current feed state and scroll context instead of refetching the list.
- Moment comments reuse the same public threaded comment components as article comments.
- Moment content is authored and stored as Markdown source, then rendered as Markdown on the public blog and in admin preview.

## Gallery Rules
- Gallery is a built-in public module under `/gallery`, not a separately deployed site.
- Public photo URLs use `/gallery/{slug}`.
- Auto-generated photo slugs must not be derived from original local filenames, because filenames may expose private information.
- When uploaded photos have no manually edited title, the system should generate a random-looking unique `IMGxxxxxxxx` display title from internal photo identity, and must not fall back to original filenames or sequential counters in public UI.
- Gallery uses a global masonry photo stream with tag filtering; albums are out of scope for v1.
- Gallery cards and viewer images must preserve original photo aspect ratios.
- Opening a photo from the gallery stream should use a card-to-viewer transition; direct entry to `/gallery/{slug}` should render the viewer immediately without relying on the source-card animation.
- Desktop stream cards show only concise hover metadata; full metadata appears in the viewer side panel.
- Mobile viewer metadata is opened from a top-left details button; the top-right close button returns to the gallery stream.
- Gallery metadata fields are public only when the field has a value and its per-photo display switch is enabled; hidden fields should be omitted by public APIs.
- Recommended per-photo display switches:
  - `taken_at_display`
  - `camera_display`
  - `location_display`
  - `exif_display`
  - `tags_display`
- Photo statuses:
  - `draft`: not publicly visible
  - `published`: listed in gallery and usable by articles
  - `hidden`: excluded from gallery stream but still usable by explicit article references
- Admin gallery upload supports batch file/folder upload, draft review, per-photo metadata editing, batch publish, and retry/removal for failed items.
- Upload processing must extract a whitelist of EXIF fields, reverse-geocode GPS into city-level location data, then remove GPS and non-whitelisted EXIF data from stored image assets.
- Admin photo upload should accept common camera/phone image formats including JPEG, PNG, WebP, HEIC, and HEIF; generated public `display`/`large` assets are normalized to JPEG.
- Precise original GPS coordinates should not be stored by default; public UI should show city-level or manually edited location text only.
- Public gallery and article image rendering should use a blurred placeholder plus optimized display/large assets stored on Cloudflare R2, not raw original camera files.
- Articles should reference gallery photos through stable gallery photo identity rather than hardcoded final image URLs whenever possible.
- Hard deletion should be blocked or safely degraded when a gallery photo is already referenced by articles.
- Detailed gallery requirements are tracked in `docs/GALLERY_REQUIREMENTS.md`.

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

- When deployed behind a reverse proxy or CDN, moderation IP capture should prefer the real visitor IP headers (for example `CF-Connecting-IP`) over proxy hop addresses.

## Analytics Rules
- Visitor analytics uses explicit event ingestion, not generic API access counting.
- Public page-visit analytics must be reported by the blog frontend from the browser after route entry so SSR data prefetch does not inflate page views.
- Initial analytics event types:
  - `page_view` for page entry and route transitions
  - `page_ping` for ongoing active-page heartbeat
    - backend updates the latest `page_view` for the same `session_id + path` with a newer active timestamp and duration instead of creating another raw row
- The analytics ingest API accepts page context from the frontend and enriches events on the backend with:
  - visitor IP
  - user agent
  - inferred device type
  - inferred browser name
  - inferred operating system
  - bot flag
- Analytics may enrich visitor IPs into cached geographic metadata (`country`, `region`, `city`, `ISP`) through an offline IP database.
- Analytics stores visitor IP in plaintext by project decision.
- Analytics should persist each accepted `page_view` as a raw visit record, and use `page_ping` to update that record's `last_engaged_at` and active duration instead of inserting unbounded heartbeat rows.
- Admin analytics must support:
  - overview metrics (`PV`, `UV`, unique IPs, sessions)
  - path-level aggregation
  - raw visit/event browsing
- Geographic filtering depends on cached offline IP lookup data and is best-effort rather than guaranteed accurate.
- Analytics is a product/statistics domain; HTTP request logs may still exist for operations or debugging, but they are not the source of truth for visitor metrics.
