# Backend Service

## Overview
This is the Go API backend for the rewritten blog platform, built with Gin.

## Requirements
- Go `1.22+`

## Run Locally
```bash
cd backend
go run ./cmd/server
```

## Database Migrations (golang-migrate)
```bash
cd backend
go run ./cmd/migrate -action up
go run ./cmd/migrate -action version
```

Quick shortcuts:
```bash
cd backend
make migrate-up
make migrate-version
```

## Environment Variables
- `APP_NAME` (default: `ancy-blog-api`)
- `APP_ENV` (default: `dev`)
- `HTTP_HOST` (default: `0.0.0.0`)
- `HTTP_PORT` (default: `8080`)
- `AUTH_ADMIN_USERNAME` (default: `admin`)
- `AUTH_ADMIN_PASSWORD` (default: `123456`)
- `AUTH_ACCESS_TOKEN_TTL_SECONDS` (default: `3600`)
- `AUTH_REFRESH_TOKEN_TTL_SECONDS` (default: `604800`)
- `DB_HOST` (default: `127.0.0.1`)
- `DB_PORT` (default: `5432`)
- `DB_NAME` (default: `ancy_blog`)
- `DB_USER` (default: `ancy`)
- `DB_PASSWORD` (default: `ancy_dev_password`)
- `DB_SSLMODE` (default: `disable`)
- `DB_MAX_OPEN_CONNS` (default: `20`)
- `DB_MAX_IDLE_CONNS` (default: `5`)
- `REDIS_ENABLED` (default: `false`)
- `REDIS_ADDR` (default: `127.0.0.1:6379`)
- `REDIS_PASSWORD` (default: empty)
- `REDIS_DB` (default: `0`)
- `REDIS_POOL_SIZE` (default: `10`)
- `REDIS_MIN_IDLE_CONNS` (default: `2`)
- `TRANSLATION_WORKER_ENABLED` (default: `true`)
- `TRANSLATION_WORKER_POLL_INTERVAL_MS` (default: `3000`)

## Health Check
- `GET /healthz`

## API Groups
- Auth: `/api/v1/auth/*`
- Public: `/api/v1/public/*`
- Admin: `/api/v1/admin/*` (Bearer token required)

## Image Upload
- Endpoint: `POST /api/v1/admin/upload/image`
- Auth: Bearer token required
- Content-Type: `multipart/form-data`
- Form field: `file`
- Current limit: `10MB`
- Note: provider configuration is planned via integration center (DB-backed).

## Repository Mode
- The app requires PostgreSQL repository at startup.
- If PostgreSQL connection fails, server initialization fails and the process exits.
- When Redis is enabled and reachable, site-related reads use cache-aside:
  - `site:settings:default`
  - `site:footer:default`
  - `site:social:default`
  - `site:nav:default`
  - `site:slot:{slotKey}:default`

## Comment APIs (Implemented)
- Public:
  - `GET /api/v1/public/comments/article/:articleId`
  - `GET /api/v1/public/comments/:id/children`
  - `GET /api/v1/public/comments/article/:articleId/total`
  - `POST /api/v1/public/comments`
- Admin:
  - `GET /api/v1/admin/comments`
  - `PUT /api/v1/admin/comments/:id`

## Demo Login
```bash
curl -X POST http://127.0.0.1:8080/api/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"admin","password":"123456"}'
```

## Notes
- The HTTP layer uses `gin-gonic/gin` with:
  - `gin.Recovery()` for panic recovery.
  - A custom structured request logging middleware.
- Handler dependencies are module-oriented (`article/comment/link/site/integration/translation/timeline` services).
- Request payloads use DTO structs under `internal/handler/dto` to decouple transport schema from domain models.
- Translation worker:
  - polls `translation_jobs` in `queued` status
  - claims one job with DB lock (`FOR UPDATE SKIP LOCKED`)
  - calls OpenAI-compatible `/chat/completions`
  - updates job to `succeeded/failed` and writes `result_text`
