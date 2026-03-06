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

## ADR-006: Singleton Settings Update Semantics
- Date: 2026-03-06
- Decision: Treat singleton-style settings resources as patch-merge updates even when exposed through `PUT` endpoints.
- Rationale: Admin UIs may save one settings section at a time, so omitted fields must preserve existing values instead of being cleared by whole-object replacement.
- Consequences:
  - Handler DTOs for singleton settings should use optional/pointer fields.
  - Handlers should read current state first, then merge only fields present in the request.
  - Frontend section forms may still merge local cached settings before submit, but backend merge behavior is the required safety boundary.
