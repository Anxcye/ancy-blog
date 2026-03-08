#!/usr/bin/env bash
# File: update.sh
# Purpose: Perform a production upgrade from git fetch to tag-aligned image release and smoke checks.
# Module: deploy automation, upgrade workflow layer.
# Related: deploy/release.sh, deploy/backup-postgres.sh, docs/DEPLOYMENT.md.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"
TARGET_REF="${1:-}"
ENV_FILE="${SCRIPT_DIR}/.env"

cd "${SCRIPT_DIR}"

if [[ ! -f "${ENV_FILE}" ]]; then
  echo "deploy/.env is required. Copy deploy/.env.example first." >&2
  exit 1
fi

set -a
source "${ENV_FILE}"
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

if [[ -n "${TARGET_REF}" ]]; then
  python3 - "${ENV_FILE}" "${TARGET_REF}" <<'PY'
from pathlib import Path
import sys

env_path = Path(sys.argv[1])
target = sys.argv[2]
content = env_path.read_text()
lines = content.splitlines()
updated = False
for i, line in enumerate(lines):
    if line.startswith("APP_IMAGE_TAG="):
        lines[i] = f"APP_IMAGE_TAG={target}"
        updated = True
        break
if not updated:
    lines.append(f"APP_IMAGE_TAG={target}")
env_path.write_text("\\n".join(lines) + "\\n")
PY
  export APP_IMAGE_TAG="${TARGET_REF}"
fi

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
