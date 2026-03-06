# Deployment Guide

## Target Topology
- Public blog: `https://example.com`
- Public API: `https://example.com/api/v1`
- Admin app: `https://admin.example.com`
- CDN / edge proxy: Cloudflare
- Origin runtime: single Linux host with Docker Compose

## Why This Topology
- `frontend-blog` runs as Nuxt SSR, so it should stay as a Node service instead of a static export.
- `frontend-admin-react` is a Vite SPA, so it can be built once and served as static files.
- `backend` is a standalone Go API and migration runner.
- A single-host `docker compose` deployment keeps the release path simple and debuggable while the product is still evolving.

## Services
- `caddy`: TLS termination and reverse proxy
- `frontend-blog`: Nuxt SSR service
- `backend`: Go API service
- `admin`: static React admin site served by nginx
- `postgres`: primary database
- `redis`: optional cache / worker acceleration

## Request Routing
- `example.com/*` -> `frontend-blog`
- `example.com/api/*` -> `backend`
- `admin.example.com/*` -> `admin`

## Cloudflare Setup
1. Create DNS records for `example.com` and `admin.example.com`.
2. Point both records to the origin server public IP.
3. Enable the Cloudflare proxy for both records.
4. Set SSL mode to `Full (strict)`.
5. Disable caching for:
   - `example.com/api/*`
   - `admin.example.com/*`
6. Strongly recommend adding Cloudflare Access or an IP allowlist in front of `admin.example.com`.

## Server Preparation
1. Provision an Ubuntu 22.04+ host.
2. Install Docker Engine and Docker Compose plugin.
3. Clone this repository to the server.
4. Copy `deploy/.env.example` to `deploy/.env` and fill all secrets.

## Initial Deploy
```bash
cd deploy
cp .env.example .env
# edit .env
./release.sh
```

## Release Flow
Each release should follow this order:
1. Build fresh images.
2. Start `postgres` and `redis`.
3. Run database migrations.
4. Restart `backend`.
5. Restart `frontend-blog`.
6. Restart `admin` and `caddy`.

The provided `deploy/release.sh` already follows this order.

## Manual Operations
Run migrations:
```bash
cd deploy
docker compose --env-file .env run --rm backend /app/migrate -action up
```

View logs:
```bash
cd deploy
docker compose --env-file .env logs -f backend
docker compose --env-file .env logs -f frontend-blog
docker compose --env-file .env logs -f caddy
```

Create a database backup:
```bash
cd deploy
./backup-postgres.sh
```

## Environment Variables
### Public domains
- `APP_DOMAIN`
- `ADMIN_DOMAIN`
- `ACME_EMAIL`

### Backend auth
- `AUTH_ADMIN_USERNAME`
- `AUTH_ADMIN_PASSWORD`
- `AUTH_ACCESS_TOKEN_TTL_SECONDS`
- `AUTH_REFRESH_TOKEN_TTL_SECONDS`

### Database
- `DB_NAME`
- `DB_USER`
- `DB_PASSWORD`
- `DB_SSLMODE`
- `DB_MAX_OPEN_CONNS`
- `DB_MAX_IDLE_CONNS`

### Redis
- `REDIS_ENABLED`
- `REDIS_PASSWORD`
- `REDIS_DB`
- `REDIS_POOL_SIZE`
- `REDIS_MIN_IDLE_CONNS`

### Runtime
- `CORS_ALLOWED_ORIGINS`
- `TRANSLATION_WORKER_ENABLED`
- `TRANSLATION_WORKER_POLL_INTERVAL_MS`
- `TRANSLATION_WORKER_BACKOFF_BASE_MS`
- `TRANSLATION_WORKER_BACKOFF_MAX_MS`

## Update Strategy
Current recommendation: keep releases manual and deterministic.
- Tag a release in git.
- Pull the target revision on the server.
- Run `deploy/release.sh`.
- Validate `/healthz`, homepage, and admin login.

Later, this can evolve into:
- CI builds images
- registry push
- server pull + `docker compose up -d`

## Backup Policy
- At minimum, run one PostgreSQL dump per day.
- Store dumps outside the VM when possible.
- Keep at least 7 daily copies before pruning.
