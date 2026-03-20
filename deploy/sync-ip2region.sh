#!/usr/bin/env bash
# File: sync-ip2region.sh
# Purpose: Download or refresh ip2region xdb runtime data files for local and deployment use.
# Module: deploy automation, runtime data sync workflow layer.
# Related: backend runtime-data/ip, deploy/release.sh, deploy/update.sh, and backend analytics IP enrichment.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"
ENV_FILE="${SCRIPT_DIR}/.env"
TARGET_DIR="${REPO_DIR}/backend/runtime-data/ip"
SOURCE_REF="${IP2REGION_SOURCE_REF:-master}"
AUTO_SYNC="${IP2REGION_AUTO_SYNC:-true}"
FORCE_SYNC="${1:-}"

if [[ -f "${ENV_FILE}" ]]; then
  set -a
  source "${ENV_FILE}"
  set +a
fi

if [[ "${AUTO_SYNC}" != "true" && "${AUTO_SYNC}" != "1" && "${AUTO_SYNC}" != "yes" ]]; then
  echo "Skipping ip2region sync because IP2REGION_AUTO_SYNC=${AUTO_SYNC}."
  exit 0
fi

mkdir -p "${TARGET_DIR}"

V4_URL="${IP2REGION_V4_URL:-https://raw.githubusercontent.com/lionsoul2014/ip2region/${SOURCE_REF}/data/ip2region_v4.xdb}"
V6_URL="${IP2REGION_V6_URL:-https://raw.githubusercontent.com/lionsoul2014/ip2region/${SOURCE_REF}/data/ip2region_v6.xdb}"

download_file() {
  local url="$1"
  local output="$2"
  local label="$3"
  local tmp="${output}.tmp"

  echo "Syncing ${label} from ${url}"
  curl --fail --location --silent --show-error "${url}" -o "${tmp}"
  mv "${tmp}" "${output}"
}

if [[ "${FORCE_SYNC}" == "--force" ]]; then
  download_file "${V4_URL}" "${TARGET_DIR}/ip2region_v4.xdb" "ip2region IPv4 xdb"
  download_file "${V6_URL}" "${TARGET_DIR}/ip2region_v6.xdb" "ip2region IPv6 xdb"
  echo "ip2region sync completed."
  exit 0
fi

download_file "${V4_URL}" "${TARGET_DIR}/ip2region_v4.xdb" "ip2region IPv4 xdb"
download_file "${V6_URL}" "${TARGET_DIR}/ip2region_v6.xdb" "ip2region IPv6 xdb"

echo "ip2region sync completed."
