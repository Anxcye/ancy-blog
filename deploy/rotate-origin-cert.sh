#!/usr/bin/env bash
# File: rotate-origin-cert.sh
# Purpose: Reload Caddy after replacing Cloudflare Origin Certificate files on the server.
# Module: deploy automation, certificate maintenance workflow layer.
# Related: deploy/docker-compose.yml, deploy/.env.example, deploy/caddy/certs/.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

cd "${SCRIPT_DIR}"

if [[ ! -f .env ]]; then
  echo "deploy/.env is required. Copy deploy/.env.example first." >&2
  exit 1
fi

set -a
source ./.env
set +a

if [[ -z "${ORIGIN_CERT_FILE:-}" || -z "${ORIGIN_KEY_FILE:-}" ]]; then
  echo "ORIGIN_CERT_FILE and ORIGIN_KEY_FILE must be set in deploy/.env." >&2
  exit 1
fi

HOST_CERT_PATH="${SCRIPT_DIR}/caddy/certs/$(basename "${ORIGIN_CERT_FILE}")"
HOST_KEY_PATH="${SCRIPT_DIR}/caddy/certs/$(basename "${ORIGIN_KEY_FILE}")"

if [[ ! -f "${HOST_CERT_PATH}" ]]; then
  echo "Certificate file not found: ${HOST_CERT_PATH}" >&2
  exit 1
fi

if [[ ! -f "${HOST_KEY_PATH}" ]]; then
  echo "Private key file not found: ${HOST_KEY_PATH}" >&2
  exit 1
fi

docker compose --env-file .env up -d caddy

echo "Origin certificate rotation completed."
