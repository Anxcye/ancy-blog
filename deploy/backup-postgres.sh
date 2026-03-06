#!/usr/bin/env bash
# File: backup-postgres.sh
# Purpose: Create a timestamped PostgreSQL dump from the running production container.
# Module: deploy automation, backup workflow layer.
# Related: deploy/docker-compose.yml, deploy/.env.example.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BACKUP_DIR="${SCRIPT_DIR}/backups"
TIMESTAMP="$(date +%Y%m%d-%H%M%S)"

cd "${SCRIPT_DIR}"

if [[ ! -f .env ]]; then
  echo "deploy/.env is required. Copy deploy/.env.example first." >&2
  exit 1
fi

mkdir -p "${BACKUP_DIR}"

set -a
source ./.env
set +a

docker compose --env-file .env exec -T postgres \
  pg_dump -U "${DB_USER}" -d "${DB_NAME}" \
  > "${BACKUP_DIR}/${DB_NAME}-${TIMESTAMP}.sql"

echo "Backup written to ${BACKUP_DIR}/${DB_NAME}-${TIMESTAMP}.sql"
