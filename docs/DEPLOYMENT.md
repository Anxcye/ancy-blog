# Deployment Guide

## Target Topology
- Public blog: `https://example.com`
- Public API: `https://example.com/api/v1`
- Admin app: `https://admin.example.com`
- CDN / edge proxy: Cloudflare
- Origin runtime: single Linux host with Docker Compose
- Image delivery: GitHub Actions builds to GHCR, server pulls prebuilt images
- Origin TLS: Cloudflare Origin Certificate mounted into Caddy

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
7. Create a Cloudflare Origin Certificate that covers both the apex domain and wildcard subdomain, for example:
   - `example.com`
   - `*.example.com`

## Server Preparation
1. Provision an Ubuntu 22.04+ host.
2. Install Docker Engine and Docker Compose plugin.
3. Clone this repository to the server.
4. Publish GHCR packages as public, or run `docker login ghcr.io` on the server before the first deploy.
5. Create `deploy/caddy/certs/origin.pem` and `deploy/caddy/certs/origin.key` from the Cloudflare Origin Certificate.
6. Copy `deploy/.env.example` to `deploy/.env` and fill all secrets. The default `IMAGE_NAMESPACE` is already `anxcye/ancy-blog`.
7. If the server needs custom redirects or environment-specific Caddy rules, place them in `deploy/caddy/local/*.caddy`.
8. Review the `ip2region` sync settings in `deploy/.env`. The default setup auto-downloads `ip2region_v4.xdb` and `ip2region_v6.xdb` into `backend/runtime-data/ip/` during each release.

## Initial Deploy
```bash
cd deploy
cp .env.example .env
# edit .env
./sync-ip2region.sh
./release.sh
```

Recommended production value in `deploy/.env`:
```env
IMAGE_REGISTRY=ghcr.io
IMAGE_NAMESPACE=anxcye/ancy-blog
APP_IMAGE_TAG=v1.0.0
ORIGIN_CERT_FILE=/etc/caddy/certs/origin.pem
ORIGIN_KEY_FILE=/etc/caddy/certs/origin.key
IP2REGION_AUTO_SYNC=true
IP2REGION_SOURCE_REF=master
IP2REGION_V4_XDB_PATH=/app/runtime-data/ip/ip2region_v4.xdb
IP2REGION_V6_XDB_PATH=/app/runtime-data/ip/ip2region_v6.xdb
```

## Upgrade Flow
For normal production upgrades, use the wrapper script:
```bash
cd deploy
./update.sh v1.0.3
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
3. If `<ref>` is provided, update `deploy/.env` so `APP_IMAGE_TAG=<ref>`
4. PostgreSQL backup
5. Sync `ip2region` runtime data files into `backend/runtime-data/ip/`
6. Image pull from GHCR
7. Database migration
8. Service restart
9. Basic smoke checks for blog, admin, and the public site API

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

Rotate the Cloudflare Origin Certificate after replacing `deploy/caddy/certs/origin.pem` and `deploy/caddy/certs/origin.key`:
```bash
cd deploy
./rotate-origin-cert.sh
```

Run the lower-level release script without pulling git:
```bash
cd deploy
./release.sh
```

Refresh the offline IP databases manually:
```bash
cd deploy
./sync-ip2region.sh
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
- `ORIGIN_CERT_FILE`
- `ORIGIN_KEY_FILE`

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
- `IP2REGION_AUTO_SYNC`
- `IP2REGION_SOURCE_REF`
- `IP2REGION_V4_URL`
- `IP2REGION_V6_URL`
- `IP2REGION_V4_XDB_PATH`
- `IP2REGION_V6_XDB_PATH`
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
- Keep `APP_IMAGE_TAG` aligned to a published git tag instead of relying on `latest`.
- Keep Cloudflare origin certificate files in `deploy/caddy/certs/`.
- Keep server-specific Caddy rules in `deploy/caddy/local/*.caddy`.
- Run `deploy/update.sh <tag>` on the server.
- Let the script backup, pull, migrate, restart, and smoke-check the stack.

Later, this can evolve into:
- CI builds images
- registry push
- server pull + `docker compose up -d`

## Backup Policy
- At minimum, run one PostgreSQL dump per day.
- Store dumps outside the VM when possible.
- Keep at least 7 daily copies before pruning.

## Origin Certificate Rotation
Cloudflare Origin Certificates are long-lived, but they are not auto-renewed by this stack.

Recommended rotation flow:
1. Generate a new Origin Certificate in Cloudflare.
2. Replace:
   - `deploy/caddy/certs/origin.pem`
   - `deploy/caddy/certs/origin.key`
3. Run:
   ```bash
   cd deploy
   ./rotate-origin-cert.sh
   ```
4. Verify:
   - `https://example.com`
   - `https://admin.example.com`
   - `docker compose --env-file .env logs --tail=100 caddy`
