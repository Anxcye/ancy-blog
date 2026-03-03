# Session Log

## 2026-03-03
### Summary
- Created collaboration baseline documents:
  - `AGENTS.md`
  - `docs/ROADMAP.md`
  - `docs/PROGRESS.md`
  - `docs/API_CONTRACT.md`
  - `docs/DECISIONS.md`
  - `docs/SESSION_LOG.md`
- Added mandatory source file header rule.
- Added English-only inline comment rule.
- Initialized backend Go scaffold:
  - `backend/go.mod`
  - `backend/cmd/server/main.go`
  - `backend/internal/config/config.go`
  - `backend/internal/logger/logger.go`
  - `backend/internal/server/http.go`
  - `backend/internal/handler/health.go`
  - `backend/internal/response/response.go`
  - `backend/README.md`
- Noted local environment gap: Go is not installed on this machine yet.
- Switched HTTP layer from `net/http` mux to `gin-gonic/gin`.
- Added request logging middleware and Gin recovery middleware.
- Added `docs/DATA_MODEL.md` with v1 table design for:
  - `articles`
  - `comments`
  - `reactions`
- Recorded explicit decision to store `ip` in plaintext for moderation use.
- Added `docs/PRODUCT_RULES.md` and clarified:
  - UUID internal IDs + slug-based public article URLs
  - homepage placement slots (e.g. `home_about`)
  - per-article comment switch
  - moments as independent short content
  - friend-link submission and review workflow
  - AI disclosure level taxonomy
- Expanded `docs/DATA_MODEL.md` with:
  - taxonomy tables (`categories`, `tags`, `article_tags`)
  - `moments`
  - `links` (with review and optional related article)
- Expanded `docs/API_CONTRACT.md` with:
  - route group conventions (`public/admin/auth`)
  - article, moment, link, and taxonomy endpoints
- Expanded site configuration scope:
  - global settings (`site_name`, avatar, homepage intro markdown)
  - footer item management (3-row placement + ordering, internal/external/plain text)
  - homepage social links management
  - internal page-like content managed by `articles` (`content_kind=page`) and slug
- Added Redis cache policy for read-heavy site data in product rules and API notes.
- Merged standalone `pages` concept into `articles`:
  - removed `pages` table contract
  - added `content_kind` in `articles`
  - footer internal link target switched to article slug
- Added frontend-oriented dynamic navigation and slot architecture:
  - dynamic top nav via `nav_items`
  - reusable content mounting via `content_slots` + `content_slot_items`
  - timeline API for mixed article/moment feed
  - hover-related slot/category retrieval APIs
- Removed legacy article `placement` API from contract and unified mounting via slot APIs.
- Implemented backend code based on contracts:
  - layered structure: `handler -> service -> repository`
  - in-memory repository with seed data
  - auth APIs (`login`, `refresh`, `me`) + bearer middleware
  - public APIs for articles, moments, links, taxonomy, site, slots, timeline
  - admin APIs for content and site management
- Verified compile and smoke test:
  - `go test ./...` passed
  - `/healthz`, `/api/v1/auth/login`, and `/api/v1/public/timeline` returned expected results
- Added Cloudflare R2 image upload capability:
  - new storage abstraction and R2 uploader
  - new admin endpoint `POST /api/v1/admin/upload/image`
  - new R2 environment configuration fields
  - updated backend usage docs and API contract
- Refined documentation-only design for unified integration center:
  - grouped R2 and LLM config into one admin integration domain
  - added data model tables: `integration_providers`, `translation_jobs`
  - expanded API contract with integration management and translation job endpoints
  - updated roadmap/progress to include integration center milestone

### Next Suggested Tasks
1. Install Go `1.22+`, run `go mod tidy`, then `go run ./cmd/server` under `backend/`.
2. Implement domain structs and DTOs from `docs/DATA_MODEL.md` and `docs/API_CONTRACT.md`.
3. Implement first admin/public APIs for site settings, footer items, and social links with cache-aside pattern.
