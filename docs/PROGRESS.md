# Progress Board

## Todo
- [ ] Initialize `frontend-blog/` Vue3 scaffold.
- [ ] Set up lint/test/build scripts for frontend workspaces.

## In Progress
- [ ] Add API rate limiting for auth/comment/AI/translation-create endpoints.
- [ ] Add frontend i18n routing and SEO output (`/` default zh, `/en/*` english).

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
- [x] Added translation worker unit tests with mocked LLM responses (success/failure/empty output/provider disabled).
- [x] Added non-streaming admin AI assist endpoints for article summary and slug suggestion (`/api/v1/admin/ai/summary`, `/api/v1/admin/ai/slug`) with deterministic fallback and slug de-duplication.
- [x] Added locale-aware read support for public moments and timeline APIs (`locale` query passthrough + repository translation readback).
- [x] Added admin translation-content management APIs for browse/detail/manual override (`/api/v1/admin/translations/contents*`).
- [x] Added translation worker retry/backoff and manual retry endpoint for failed jobs (`/api/v1/admin/translations/jobs/{id}/retry`).
- [x] Added i18n publish-control schema and runtime: translation `autoPublish/publishAt`, article translation `title/summary/content`, and locale read visibility based on translation publish state.
- [x] Updated frontend style baseline with hero/nav/motion/particle constraints and mobile+SEO requirements.
- [x] Added `frontend-admin` responsive shell baseline (desktop top nav + mobile 5-tab full-feature entry) and fresh teal design tokens.
- [x] Added admin article management baseline: backend list/detail endpoints and frontend list+editor pages (content/meta/preview modes).
- [x] Added admin site/system operation pages: site settings/footer/social/nav management and integration-center config/test UI.
- [x] Added admin interaction center page with comment moderation and friend-link review operations.
- [x] Added translation operations center (job create/list/retry + translation content override) under admin system workspace.
- [x] Added AI assist actions in article editor (summary generation and slug suggestion).
- [x] Added admin moments management (backend list/update APIs and frontend moments CRUD page baseline).
- [x] Added actionable dashboard metrics and quick entries for content/moderation/system workflows.
- [x] Completed admin site CRUD loop (footer/social/nav create+edit+delete) and added global 401 redirect handling plus translation-center pagination.

- [x] Added pagination to interaction center lists (comments and link submissions).

- [x] Added slot management end-to-end (slot list/item list APIs + admin UI create/list/delete slot items).

- [x] Added admin article bulk operations (batch status switch + single/batch delete).

- [x] Added quick jump from article rows to translation center with prefilled source context.
- [x] Upgraded admin UI to Naive UI design system with enterprise shell (responsive sider/header/drawer, theme + locale switch, refined login/dashboard/articles pages).
- [x] Replaced raw article textarea editor with TipTap rich-text editor (toolbar, JSON content persistence, and inline image insertion).
- [x] Added structured integration-center forms for `openai_compatible` and `cloudflare_r2` with optional advanced JSON section.
- [x] Enabled runtime-refresh R2 upload configuration (no backend restart needed after integration save) and protected masked secrets from accidental overwrite.
- [x] Improved article editor usability with rich preview rendering and local auto-draft restore/clear workflow.
- [x] Added post-save automatic provider test in integration center for immediate initialization feedback.
- [x] Added admin moment deletion and batch operations (batch status switch and batch delete) with API contract updates.
- [x] Unified admin pages (`site/interaction/system/moments`) onto consistent Naive UI layout and theme behavior.
- [x] Upgraded TipTap editor UX with grouped sticky toolbar, link popover actions, undo/redo, and richer typography styling.
- [x] Added extensible embed-block baseline in TipTap editor/preview (`x_post`, `tmdb_card`) for future custom content modules.
- [x] Migrated TipTap toolbar to component-library-first interactions (button groups, dropdown command menus, and modal-based embed insertion).
- [x] Integrated TipTap simple-editor template into React admin ArticleEditorPage (full toolbar, image upload, form-controlled value/onChange binding).
- [x] Initialized `frontend-admin-react` rewrite baseline (`React + TypeScript + Ant Design + React Router + Zustand + React Query + Axios`) with login/auth-guard/layout/dashboard skeleton.

- [x] Fixed article editor gaps in `frontend-admin-react`: category selector, tag multi-select, conditional scheduled-publish DatePicker, cover image upload, category filter on list page.
- [x] Implemented M2 content module for `frontend-admin-react`: shared API types, articles and moments API modules, ArticlesPage (filters/batch ops/pagination), ArticleEditorPage (create/edit with AI slug+summary assist), MomentsPage (list+drawer CRUD+batch ops), updated router, and expanded sider navigation.

- [x] Implemented categories & tags full CRUD: migration 000007 (`category_id` FK on articles), backend repo write/read helpers (`resolveCategoryID`, `resolveTagIDs`, `syncArticleTags`), admin `POST/DELETE /admin/categories` and `/admin/tags` endpoints, TaxonomyPage management UI under content nav, and inline create in article editor dropdowns for both category and tag.

## Blocked
- None.

## Notes
- Keep each task small enough to complete in one session when possible.
- Move items across sections instead of duplicating them.
