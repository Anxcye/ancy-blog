# ancy-blog (Rewrite)

A clean rewrite of the blog platform with:
- Backend: Go + Gin + PostgreSQL + Redis
- Frontend: React (admin) + Vue 3 (public site)

## Current Status
- Repository is initialized for the rewrite phase.
- Core planning and contracts are defined under `docs/`.
- Backend skeleton is available in `backend/`.

## Project Structure
- `backend/`: Go API service.
- `frontend-admin-react/`: React admin app.
- `frontend-blog/`: Vue 3 public app.
- `docs/`: product rules, data model, API contracts, progress logs.
- `AGENTS.md`: contributor and AI-collaboration rules.

## Quick Start (Backend)
```bash
cd backend
go mod tidy
go run ./cmd/server
```

Health check:
```bash
curl http://127.0.0.1:8080/healthz
```

## Core Docs
- `docs/PRODUCT_RULES.md`
- `docs/DATA_MODEL.md`
- `docs/API_CONTRACT.md`
- `docs/PROGRESS.md`
- `docs/SESSION_LOG.md`

## Collaboration Rules
- Follow `AGENTS.md`.
- Keep code comments in English.
- Add a top-of-file header comment for each source file:
  - file purpose
  - module/layer position
  - related modules
