# Architecture Decisions

## ADR-001: Backend Stack
- Date: 2026-03-03
- Decision: Use Go as the backend language.
- Rationale: Lower resource usage on VPS, faster startup, simpler runtime model.
- Consequences: Need new service scaffolding and backend reimplementation.

## ADR-002: Database
- Date: 2026-03-03
- Decision: Use PostgreSQL as primary database.
- Rationale: Strong SQL capabilities, good transactional behavior, future-friendly querying.
- Consequences: New schema and migration plan required.

## ADR-003: Cache
- Date: 2026-03-03
- Decision: Use Redis as optional acceleration layer.
- Rationale: Improve hot-read performance and support token/session workflows.
- Consequences: System must degrade safely when Redis is unavailable.

## ADR-004: Auth Model
- Date: 2026-03-03
- Decision: Start with single-admin authentication (no RBAC in MVP).
- Rationale: Faster delivery and lower maintenance overhead.
- Consequences: Keep extension points for future role support.

## ADR-005: HTTP Framework
- Date: 2026-03-03
- Decision: Use Gin as the HTTP framework.
- Rationale: Widely adopted in Go projects, mature middleware ecosystem, clear routing model.
- Consequences: Standardize handler/middleware patterns around Gin context.
