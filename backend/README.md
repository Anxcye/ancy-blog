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

## Health Check
- `GET /healthz`

## Notes
- The HTTP layer uses `gin-gonic/gin` with:
  - `gin.Recovery()` for panic recovery.
  - A custom structured request logging middleware.
