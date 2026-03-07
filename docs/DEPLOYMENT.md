# Deployment Guide

## Target Topology
- Public blog: `https://example.com`
- Public API: `https://example.com/api/v1`
- Admin app: `https://admin.example.com`
- CDN / edge proxy: Cloudflare
- Origin runtime: single Linux host with Docker Compose
- Image delivery: GitHub Actions builds to GHCR, server pulls prebuilt images

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
4. In GitHub repository settings, add variable `APP_DOMAIN=example.com` so frontend images are built with the correct API base.
5. Publish GHCR packages as public, or run `docker login ghcr.io` on the server before the first deploy.
6. Copy `deploy/.env.example` to `deploy/.env` and fill all secrets. The default `IMAGE_NAMESPACE` is already `anxcye/ancy-blog`.
7. If the server needs custom redirects or environment-specific Caddy rules, place them in `deploy/caddy/local/*.caddy`.

## Initial Deploy
```bash
cd deploy
cp .env.example .env
# edit .env
./release.sh
```

Recommended production value in `deploy/.env`:
```env
IMAGE_REGISTRY=ghcr.io
IMAGE_NAMESPACE=anxcye/ancy-blog
APP_IMAGE_TAG=latest
```

## Upgrade Flow
For normal production upgrades, use the wrapper script:
```bash
cd deploy
./update.sh
```

To deploy a specific git tag or commit:
```bash
cd deploy
./update.sh v1.0.3
# or
./update.sh 5d3024f
```

`deploy/update.sh` performs this order:
1. `git fetch --tags --prune`
2. `git pull --ff-only` or `git checkout <ref>`
3. PostgreSQL backup
4. Image pull from GHCR
5. Database migration
6. Service restart
7. Basic smoke checks for blog, admin, and the public site API

`deploy/release.sh` remains the lower-level script that only handles the container release itself.

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

Run the lower-level release script without pulling git:
```bash
cd deploy
./release.sh
```

Add server-local Caddy overrides without changing tracked files:
```bash
cd deploy
cat > caddy/local/legacy-redirects.caddy <<'EOF'
redir /home/7 /articles/example-article-slug 301
EOF
./release.sh
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
Current recommendation: keep upgrades manual and deterministic.
- Tag a release in git.
- Wait for the `build-images` workflow to finish successfully in GitHub Actions.
- The workflow only publishes on tag pushes (plus manual dispatch), so normal branch pushes do not produce production images.
- Let GitHub Actions build and publish images to GHCR.
- Keep secrets in `deploy/.env`.
- Keep server-specific Caddy rules in `deploy/caddy/local/*.caddy`.
- Run `deploy/update.sh` on the server.
- Let the script backup, pull, migrate, restart, and smoke-check the stack.

Later, this can evolve into:
- CI builds images
- registry push
- server pull + `docker compose up -d`

## Backup Policy
- At minimum, run one PostgreSQL dump per day.
- Store dumps outside the VM when possible.
- Keep at least 7 daily copies before pruning.
