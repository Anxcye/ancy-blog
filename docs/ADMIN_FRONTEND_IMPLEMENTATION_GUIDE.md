# Admin Frontend Implementation Guide

## 1. Goal and Scope
This document is a handoff guide for rebuilding the admin app with a production-ready React stack while keeping the existing Go backend APIs (`/api/v1`) unchanged.

Target workspace: `frontend-admin-react/`

Out of scope:
- Backend API redesign
- Public blog frontend (`frontend-blog/`)

## 2. Recommended Stack (Enterprise-Oriented)
- Framework: `React 19 + TypeScript + Vite`
- UI system: `Ant Design (antd)`
- Routing: `react-router-dom`
- Data fetching/cache: `@tanstack/react-query`
- HTTP client: `axios` (single shared client with interceptors)
- Local app state: `zustand` (auth/session/theme)
- Form system: `antd` Form + schema validation rules
- Optional editor layer: TipTap (headless core) wrapped by UI-library components

## 3. UI/Design System Baseline
Use Ant Design as the default component system. Avoid hand-built primitive UI unless a component does not exist in the library.

Required UI principles:
- Unified spacing scale: 8px grid (`8, 12, 16, 24`)
- Global radius: `10-12px`
- Primary color: `#1f8f8a`
- Desktop: left sider + top header
- Mobile: bottom navigation with full feature access
- Dark mode support must be token-based, not ad-hoc CSS overrides

Core components to prioritize:
- Layout: `Layout`, `Menu`, `Drawer`, `Grid`
- Data: `Table`, `Pagination`, `Tag`, `Badge`, `Empty`, `Skeleton`
- Forms: `Form`, `Input`, `Select`, `Switch`, `DatePicker`, `Upload`, `Modal`, `Drawer`
- Feedback: `Message`, `Notification`, `Result`, `Spin`

## 4. Current React Baseline (Already Done)
Implemented in `frontend-admin-react/`:
- Auth store (`src/store/auth.ts`)
- HTTP client + 401 interceptor (`src/lib/http.ts`)
- Protected routing (`src/App.tsx`)
- Login page (`src/pages/LoginPage.tsx`)
- Admin shell (`src/layouts/AdminLayout.tsx`)
- Dashboard placeholder (`src/pages/DashboardPage.tsx`)

## 5. Migration Modules (Execution Order)
1. Content module
- Article list, filters, pagination, batch actions
- Article editor: TipTap + AI assist (summary/slug) + image upload
- Moments list/editor + batch status/delete

2. Interaction module
- Comment moderation table
- Link submission review table
- Unified status actions and audit feedback

3. Site module
- Site settings
- Footer/social/nav/slot CRUD
- Stable responsive form layout

4. System module
- Integration center (`openai_compatible`, `cloudflare_r2`)
- Structured provider forms (no raw JSON-only UX)
- Save -> auto test feedback
- Translation jobs and translation content override pages

## 6. API Integration Rules
- Keep response envelope parsing centralized:
  - `{ code, message, data }`
- Handle 401 globally in axios interceptor.
- Never duplicate base URL logic in page components.
- API modules should be split by domain, e.g.:
  - `src/api/auth.ts`
  - `src/api/articles.ts`
  - `src/api/moments.ts`
  - `src/api/interactions.ts`
  - `src/api/site.ts`
  - `src/api/system.ts`

## 7. Editor Strategy
- Keep TipTap core for long-term custom block extensibility.
- Required block baseline:
  - `x_post`
  - `tmdb_card`
- Replace prompt-based insertion with typed form modals/drawers.
- Ensure preview renderer uses the same extension set as editor.

## 8. Quality Gates
Before each merge:
- `pnpm run build` passes
- Type errors = zero
- No direct API calls inside presentational components
- Responsive check on `390px` and `>=1280px`
- Update docs:
  - `docs/PROGRESS.md`
  - `docs/SESSION_LOG.md`

## 9. Suggested Milestone Plan
- M1: Shell + Auth + Dashboard + shared API/query infra
- M2: Content module complete
- M3: Interaction + Site modules complete
- M4: System/translation module complete
- M5: UI polish, accessibility, performance split chunks

## 10. Handoff Notes for Next Process
- Use `frontend-admin-react/` as the only active admin workspace.
- Keep backend contract-first approach; update `docs/API_CONTRACT.md` before endpoint changes.
