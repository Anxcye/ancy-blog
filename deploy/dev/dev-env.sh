#!/usr/bin/env bash
# File: dev-env.sh
# Purpose: Manage local development infrastructure lifecycle.
# Module: deploy/dev command automation layer.
# Related: compose.yml, .env.example, backend/cmd/migrate.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "${SCRIPT_DIR}/../.." && pwd)"
ENV_FILE="${SCRIPT_DIR}/.env"
ENV_EXAMPLE_FILE="${SCRIPT_DIR}/.env.example"
COMPOSE_FILE="${SCRIPT_DIR}/compose.yml"

if [[ ! -f "${ENV_FILE}" ]]; then
  cp "${ENV_EXAMPLE_FILE}" "${ENV_FILE}"
  echo "Created ${ENV_FILE} from .env.example"
fi

compose() {
  docker compose --env-file "${ENV_FILE}" -f "${COMPOSE_FILE}" "$@"
}

wait_for_postgres() {
  compose up -d postgres
  echo "Waiting for PostgreSQL to become healthy..."
  for _ in $(seq 1 60); do
    if compose exec -T postgres pg_isready -U "${POSTGRES_USER:-ancy}" -d "${POSTGRES_DB:-ancy_blog}" >/dev/null 2>&1; then
      return 0
    fi
    sleep 1
  done
  echo "PostgreSQL did not become ready in time." >&2
  return 1
}

run_migrations() {
  wait_for_postgres
  (
    cd "${ROOT_DIR}/backend"
    DB_HOST=127.0.0.1 \
      DB_PORT="${POSTGRES_PORT:-5432}" \
      DB_NAME="${POSTGRES_DB:-ancy_blog}" \
      DB_USER="${POSTGRES_USER:-ancy}" \
      DB_PASSWORD="${POSTGRES_PASSWORD:-ancy_dev_password}" \
      DB_SSLMODE=disable \
      run_go ./cmd/migrate -action up
  )
}

run_go() {
  if command -v go >/dev/null 2>&1; then
    go run "$@"
    return
  fi

  if command -v mise >/dev/null 2>&1; then
    mise exec -- go run "$@"
    return
  fi

  echo "Go is not installed and mise is unavailable. Install Go or run 'mise install'." >&2
  return 127
}

usage() {
  cat <<'EOF'
Usage: ./deploy/dev/dev-env.sh <command>

Commands:
  up       Start PostgreSQL and Redis, then run database migrations.
  start    Alias for up.
  migrate  Run database migrations against the local PostgreSQL container.
  reset    Stop containers, delete dev volumes, recreate services, and migrate.
  stop     Stop local development containers without deleting volumes.
  down     Stop and remove local development containers without deleting volumes.
  status   Show local development container status.
  logs     Follow local development service logs.

The script creates deploy/dev/.env from .env.example on first run.
EOF
}

load_env() {
  set -a
  # shellcheck disable=SC1090
  source "${ENV_FILE}"
  set +a
}

command="${1:-}"
load_env

case "${command}" in
  up | start)
    compose up -d
    run_migrations
    compose ps
    ;;
  migrate)
    run_migrations
    ;;
  reset)
    compose down --volumes --remove-orphans
    compose up -d
    run_migrations
    compose ps
    ;;
  stop)
    compose stop
    ;;
  down)
    compose down --remove-orphans
    ;;
  status)
    compose ps
    ;;
  logs)
    compose logs -f
    ;;
  -h | --help | help)
    usage
    ;;
  *)
    usage
    exit 1
    ;;
esac
