#!/usr/bin/env bash
# File: release.sh
# Purpose: Run the production release flow with image pull, migration, and service restart.
# Module: deploy automation, release workflow layer.
# Related: deploy/docker-compose.yml, backend/cmd/migrate, deploy/.env.example.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

cd "${SCRIPT_DIR}"

if [[ ! -f .env ]]; then
  echo "deploy/.env is required. Copy deploy/.env.example first." >&2
  exit 1
fi

./sync-ip2region.sh

docker compose --env-file .env pull backend frontend-blog admin
docker compose --env-file .env up -d postgres redis
docker compose --env-file .env run --rm backend /app/migrate -action up
docker compose --env-file .env up -d backend frontend-blog admin caddy

echo "Release completed."
