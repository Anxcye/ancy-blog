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

## Environment Variables
- `APP_NAME` (default: `ancy-blog-api`)
- `APP_ENV` (default: `dev`)
- `HTTP_HOST` (default: `0.0.0.0`)
- `HTTP_PORT` (default: `8080`)
- `AUTH_ADMIN_USERNAME` (default: `admin`)
- `AUTH_ADMIN_PASSWORD` (default: `123456`)
- `AUTH_ACCESS_TOKEN_TTL_SECONDS` (default: `3600`)
- `AUTH_REFRESH_TOKEN_TTL_SECONDS` (default: `604800`)

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
