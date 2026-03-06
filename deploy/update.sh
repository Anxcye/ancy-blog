#!/usr/bin/env bash
# File: update.sh
# Purpose: Perform a production upgrade from git fetch to backup, release, and smoke checks.
# Module: deploy automation, upgrade workflow layer.
# Related: deploy/release.sh, deploy/backup-postgres.sh, docs/DEPLOYMENT.md.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"
TARGET_REF="${1:-}"

cd "${SCRIPT_DIR}"

if [[ ! -f .env ]]; then
  echo "deploy/.env is required. Copy deploy/.env.example first." >&2
  exit 1
fi

set -a
source ./.env
set +a

if [[ -z "${APP_DOMAIN:-}" || -z "${ADMIN_DOMAIN:-}" ]]; then
  echo "APP_DOMAIN and ADMIN_DOMAIN must be set in deploy/.env." >&2
  exit 1
fi

cd "${REPO_DIR}"

git fetch --tags --prune

if [[ -n "${TARGET_REF}" ]]; then
  git checkout "${TARGET_REF}"
else
  git pull --ff-only
fi

cd "${SCRIPT_DIR}"

./backup-postgres.sh
./release.sh

BLOG_URL="https://${APP_DOMAIN}"
ADMIN_URL="https://${ADMIN_DOMAIN}"
API_HEALTH_URL="https://${APP_DOMAIN}/api/v1/public/site/settings"

echo "Running smoke checks..."
curl --fail --silent --show-error --location "${BLOG_URL}" >/dev/null
curl --fail --silent --show-error --location "${ADMIN_URL}" >/dev/null
curl --fail --silent --show-error --location "${API_HEALTH_URL}" >/dev/null

echo "Upgrade completed successfully."
