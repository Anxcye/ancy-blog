# Repository Guidelines

## Purpose
This repository is a full rewrite of the blog system using `Go + PostgreSQL + Redis + Vue3`.  
Development is AI-assisted and document-driven.

## Project Structure
- `backend/`: Go API service.
- `frontend-admin/`: Vue3 admin app.
- `frontend-blog/`: Vue3 public blog app.
- `deploy/`: deployment scripts and runtime manifests.
- `docs/`: roadmap, progress, API contracts, decisions, and session logs.

## Mandatory File Header Rule
Every source code file must start with a short header comment (3-6 lines) describing:
1. File purpose.
2. Module and layer position.
3. Main dependencies or related modules.

Example:
```go
// File: auth_handler.go
// Purpose: Handles login and token refresh HTTP endpoints.
// Module: backend/auth, presentation layer.
// Related: auth_service, token_store, user_repository.
```

## Comment Language
All inline code comments must be in English.

## Workflow Rules
1. Update `docs/PRODUCT_RULES.md` when business behavior changes.
2. Update `docs/DATA_MODEL.md` when entity fields/relations change.
3. Update `docs/API_CONTRACT.md` before implementing or changing endpoints.
4. Update `docs/FRONTEND_STYLE_GUIDE.md` when frontend visual/motion direction changes.
5. Update `docs/PROGRESS.md` when a task status changes.
6. Record architecture choices in `docs/DECISIONS.md`.
7. At the end of each work session, append to `docs/SESSION_LOG.md`.

## Content Modeling Rules
1. Use unified `articles` with `content_kind` (`post | page`) instead of a separate `pages` table.
2. Use slot mapping (`content_slots` + `content_slot_items`) for content placement, not hardcoded placement fields.
3. Use dynamic navigation (`nav_items`) for top menu and hover/dropdown behavior.

## Engineering Baseline (Enterprise-Oriented)
1. Backend framework: `Gin` with structured logging and recovery middleware.
2. Configuration: environment-driven config only; no hardcoded secrets.
3. API standard: versioned routes (`/api/v1`), unified response envelope, stable error codes.
4. Layering: `handler -> service -> repository`, with clear dependency direction.
5. Database change management: versioned SQL migrations via `golang-migrate` (`backend/migrations` + `cmd/migrate`).
6. Quality gates: run lint/test/build before merge; keep code production-ready by default.

## Definition of Done
A task is done only if:
1. Code is implemented.
2. Lint/test/build pass locally.
3. Contract/docs are updated.
4. Progress and session log are updated.
