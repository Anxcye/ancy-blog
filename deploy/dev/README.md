# Local Development Infrastructure

This directory owns local Docker infrastructure only: PostgreSQL and Redis. It is intentionally separate from `deploy/docker-compose.yml`, which represents the production deployment topology.

## Commands

```bash
./deploy/dev/dev-env.sh up
./deploy/dev/dev-env.sh reset
./deploy/dev/dev-env.sh stop
./deploy/dev/dev-env.sh status
```

`up` starts PostgreSQL and Redis, waits for PostgreSQL, and runs backend migrations.

`reset` removes the local development containers and named volumes, recreates them, and runs migrations again. This deletes local development database and Redis data.

The first run creates `deploy/dev/.env` from `.env.example`.

## Backend Settings

The default values match `backend/.env.example`:

```env
DB_HOST=127.0.0.1
DB_PORT=5432
DB_NAME=ancy_blog
DB_USER=ancy
DB_PASSWORD=ancy_dev_password
REDIS_ADDR=127.0.0.1:6379
```

## Practice

Keep local infrastructure, production deployment, and application runtime commands separate:

- `deploy/dev/compose.yml` runs local dependencies with disposable volumes.
- `deploy/docker-compose.yml` runs the production host topology.
- Backend/frontend dev servers still run from their own workspaces.
