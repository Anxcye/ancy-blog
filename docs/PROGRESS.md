# Progress Board

## Todo
- [ ] Initialize `frontend-admin/` Vue3 scaffold.
- [ ] Initialize `frontend-blog/` Vue3 scaffold.
- [ ] Implement PostgreSQL schema in code (migration not required in current plan).
- [ ] Set up lint/test/build scripts.
- [ ] Replace in-memory repository with PostgreSQL + Redis cache-aside.

## In Progress
- [ ] Replace in-memory repository with PostgreSQL + Redis cache-aside.

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

## Blocked
- None.

## Notes
- Keep each task small enough to complete in one session when possible.
- Move items across sections instead of duplicating them.
