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


## ADR-007: Production Deployment Topology
- Date: 2026-03-06
- Decision: Deploy the system on a single Linux host using Docker Compose, with Cloudflare in front, Caddy as the origin reverse proxy, `frontend-blog` as a Nuxt SSR service on `example.com`, `backend` routed through `example.com/api`, and `frontend-admin-react` served on `admin.example.com`.
- Rationale: This project currently benefits more from operational simplicity, deterministic releases, and easy rollback/debugging than from orchestration complexity. The chosen split also keeps the public site same-origin with the API while isolating the admin surface on its own domain.
- Consequences:
  - Production releases should run migrations before application restarts.
  - Frontend and backend images should be built in CI and published to GHCR, while the server only pulls and runs them.
  - `frontend-blog` remains an SSR Node service instead of static export.
  - `frontend-admin-react` is built as static assets and served behind the same reverse proxy layer.
  - Cloudflare caching and access policies should treat `example.com/api/*` and `admin.example.com/*` as dynamic/non-cacheable paths.

## ADR-008: Visitor Analytics Collection Model
- Date: 2026-03-20
- Decision: Collect visitor analytics through an explicit browser-reported events API instead of deriving page metrics from generic backend request logs.
- Rationale: The public blog runs with Nuxt SSR and server-side data prefetch, so counting backend content API calls would overstate real page visits and blur browser traffic with server/internal fetches.
- Consequences:
  - The blog frontend must report page events from the browser after route entry.
  - The backend remains responsible for IP capture, user-agent parsing, bot detection, and persistent storage.
  - Raw analytics events are queryable by admins as first-class application data.
  - HTTP access logs remain optional for debugging/ops, but they are not the product analytics source of truth.

## ADR-009: Offline IP Geography Enrichment
- Date: 2026-03-20
- Decision: Use `ip2region` xdb files as an offline IP-enrichment source, cache lookup results in `ip_profiles`, and keep the xdb file out of git as runtime data.
- Rationale: The project already stores plaintext visitor IPs by decision, and admin analytics needs country/region/city/ISP filtering without adding synchronous online API dependencies or large in-memory caches.
- Consequences:
  - Runtime configuration should provide `IP2REGION_V4_XDB_PATH` and optionally `IP2REGION_V6_XDB_PATH`.
  - Analytics ingest best-effort enriches new IPs into cached region records.
  - Admin analytics visit queries may filter by `countryName`, `regionName`, `cityName`, and `isp`.
  - Deployment automation should keep the xdb files refreshed outside git before container release.
  - Geographic metadata quality depends on the deployed xdb data file and may lag real-world IP ownership changes.

## ADR-010: CGO Runtime Dependencies in Backend Images
- Date: 2026-04-03
- Decision: Keep `CGO_ENABLED=1` for backend builds that depend on HEIF decoding, and install the matching C/C++ runtime libraries (`libstdc++`, `libgcc`) in the final Alpine runtime image.
- Rationale: The gallery HEIC/HEIF pipeline depends on CGO-backed libraries; without the C++ runtime shared objects, the server binary fails at container startup even though the build stage succeeds.
- Consequences:
  - Backend runtime images must include the shared libraries required by CGO-linked binaries.
  - Future dependency additions that introduce CGO should be validated against the final runtime image, not only the build stage.
  - If a feature can remain pure Go, `CGO_ENABLED=0` is still preferable for simpler runtime portability.
