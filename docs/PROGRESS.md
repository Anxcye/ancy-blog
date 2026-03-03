# Progress Board

## Todo
- [ ] Initialize `frontend-admin/` Vue3 scaffold.
- [ ] Initialize `frontend-blog/` Vue3 scaffold.
- [ ] Set up lint/test/build scripts.

## In Progress
- [ ] Design and implement translation job worker execution flow.

## Done
- [x] Created tracking docs and collaboration guidelines.
- [x] Established file header and English comment rules.
- [x] Initialized `backend/` Go service scaffold (`cmd`, `internal`, `healthz`, config, logger).
- [x] Switched HTTP stack to Gin and added request logging middleware.
- [x] Defined v1 data model for `articles/comments/reactions` in `docs/DATA_MODEL.md`.
- [x] Defined product-level business rules in `docs/PRODUCT_RULES.md`.
- [x] Expanded data model for `categories/tags/article_tags/moments/links`.
- [x] Expanded API contract for `public/admin/auth` route groups and core resources.
- [x] Defined site config domain (`site_settings/footer_items/social_links`) and Redis cache policy.
- [x] Defined dynamic nav/slot model and timeline API for frontend interaction design.
- [x] Implemented backend v1 APIs (`auth/public/admin`) with Gin, in-memory repository, and auth middleware.
- [x] Added Cloudflare R2 image upload support (`/api/v1/admin/upload/image`).
- [x] Defined unified integration-center docs for R2 + LLM config and translation jobs.
- [x] Added `backend/sql/schema_v1.sql` for PostgreSQL schema initialization.
- [x] Implemented PostgreSQL repository and runtime fallback to in-memory when DB unavailable.
- [x] Implemented Redis cache-aside for `site/settings/footer/social/nav/slot`.
- [x] Removed in-memory fallback: backend now fails fast when PostgreSQL init fails.
- [x] Implemented comments APIs for public and admin flows.
- [x] Implemented integration-center APIs (`integration_providers`, `translation_jobs`).

## Blocked
- None.

## Notes
- Keep each task small enough to complete in one session when possible.
- Move items across sections instead of duplicating them.
