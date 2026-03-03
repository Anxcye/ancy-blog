# Progress Board

## Todo
- [ ] Initialize `frontend-admin/` Vue3 scaffold.
- [ ] Initialize `frontend-blog/` Vue3 scaffold.
- [ ] Set up lint/test/build scripts.

## In Progress
- [ ] Add worker-oriented tests for translation execution with mocked LLM responses.
- [ ] Add localized content read APIs for moments and timeline (`locale` aware).

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
- [x] Added backend unit tests for auth service, auth middleware, and integration/translation business validation.
- [x] Added handler/config/logger/response/postgres-helper test suites and wired them into `go test ./...`.
- [x] Added tests for `app`, `server`, `cache/redis`, and `repository/memory`.
- [x] Refactored handlers to module-oriented service dependencies (`article/comment/link/site/integration/translation/timeline`).
- [x] Introduced handler DTO layer (`internal/handler/dto`) for transport payload isolation.
- [x] Introduced typed application error taxonomy (`internal/apperr`) and replaced string-based handler matching.
- [x] Added versioned migration mechanism with `golang-migrate` (`cmd/migrate` + `migrations/*.sql` + `Makefile` shortcuts).
- [x] Implemented translation worker runtime (`queued -> running -> succeeded/failed`) with OpenAI-compatible execution path.
- [x] Added PostgreSQL integration tests (`-tags=integration`) for article/comment/translation lifecycle.
- [x] Added API e2e smoke test (`-tags=integration`) covering auth/admin/public/translation flows.
- [x] Added locale translation persistence (`article_translations`, `moment_translations`) and worker writeback.
- [x] Added article detail `locale` query support with translation fallback.

## Blocked
- None.

## Notes
- Keep each task small enough to complete in one session when possible.
- Move items across sections instead of duplicating them.
