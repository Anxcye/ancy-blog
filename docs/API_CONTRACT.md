# API Contract

## Conventions
- Base path: `/api/v1`
- Response envelope:
```json
{
  "code": "OK",
  "message": "success",
  "data": {}
}
```
- Auth: `Authorization: Bearer <access_token>`
- Time format: ISO-8601 UTC.

## Route Groups
- Public read APIs: `/api/v1/public/*`
- Admin management APIs: `/api/v1/admin/*`
- Auth APIs: `/api/v1/auth/*`

## Endpoint Template
- ID:
- Method:
- Path:
- Auth Required: Yes/No
- Request:
- Response:
- Error Codes:
- Notes:

## Auth
- ID: AUTH-001
- Method: POST
- Path: /api/v1/auth/login
- Auth Required: No
- Request:
```json
{ "username": "admin", "password": "string" }
```
- Response:
```json
{
  "accessToken": "jwt",
  "refreshToken": "token",
  "expiresIn": 3600
}
```
- Error Codes: AUTH_INVALID_CREDENTIALS, AUTH_RATE_LIMITED
- Notes: Single admin account for MVP.

- ID: AUTH-002
- Method: POST
- Path: /api/v1/auth/refresh
- Auth Required: No
- Request:
```json
{ "refreshToken": "token" }
```
- Response:
```json
{
  "accessToken": "jwt",
  "refreshToken": "token",
  "expiresIn": 3600
}
```
- Error Codes: AUTH_REFRESH_INVALID

- ID: AUTH-003
- Method: GET
- Path: /api/v1/auth/me
- Auth Required: Yes
- Request: None
- Response:
```json
{
  "id": "uuid",
  "username": "admin",
  "displayName": "Ancy"
}
```
- Error Codes: AUTH_UNAUTHORIZED

## Public - Articles
- ID: PUB-ARTICLE-001
- Method: GET
- Path: /api/v1/public/articles
- Auth Required: No
- Request: query `page`, `pageSize`, `category`, `tag`, `contentKind` (`post` by default)
- Response: paginated article cards
- Error Codes: None

- ID: PUB-ARTICLE-002
- Method: GET
- Path: /api/v1/public/articles/{slug}
- Auth Required: No
- Request: query `locale` (optional, e.g. `en-US`)
- Response: article detail (includes `allowComment`, AI disclosure, source info, `contentKind`)
- Error Codes: ARTICLE_NOT_FOUND
- Notes: when `locale` is provided and translation exists, response content returns localized variant.

- ID: PUB-ARTICLE-003
- Method: GET
- Path: /api/v1/public/articles/by-category/{categorySlug}
- Auth Required: No
- Request: query `page`, `pageSize`
- Response: paginated article cards
- Error Codes: CATEGORY_NOT_FOUND

## Admin - Articles
- ID: ADM-ARTICLE-001
- Method: POST
- Path: /api/v1/admin/articles
- Auth Required: Yes
- Request: article create payload (`contentKind` supports `post | page`)
- Response: created article id
- Error Codes: VALIDATION_ERROR, AUTH_UNAUTHORIZED

- ID: ADM-ARTICLE-002
- Method: PUT
- Path: /api/v1/admin/articles/{id}
- Auth Required: Yes
- Request: article update payload
- Response: updated article id
- Error Codes: ARTICLE_NOT_FOUND, VALIDATION_ERROR
- Notes: Content mounting is handled by slot APIs under `/api/v1/admin/site/slots/*`.

- ID: ADM-ARTICLE-003
- Method: GET
- Path: /api/v1/admin/articles
- Auth Required: Yes
- Request: query `page`, `pageSize`, `status`, `contentKind`, `keyword`
- Response: paginated article list (includes `draft/published/scheduled` records)
- Error Codes: AUTH_UNAUTHORIZED

- ID: ADM-ARTICLE-004
- Method: GET
- Path: /api/v1/admin/articles/{id}
- Auth Required: Yes
- Request: None
- Response: article detail for editing
- Error Codes: ARTICLE_NOT_FOUND, AUTH_UNAUTHORIZED

- ID: ADM-ARTICLE-005
- Method: DELETE
- Path: /api/v1/admin/articles/{id}
- Auth Required: Yes
- Request: None
- Response: success boolean
- Error Codes: ARTICLE_NOT_FOUND, AUTH_UNAUTHORIZED

- ID: ADM-ARTICLE-006
- Method: POST
- Path: /api/v1/admin/articles/batch-status
- Auth Required: Yes
- Request: `{ "ids": ["uuid"], "status": "draft|published|scheduled" }`
- Response: affected count
- Error Codes: VALIDATION_ERROR, AUTH_UNAUTHORIZED

- ID: ADM-ARTICLE-007
- Method: POST
- Path: /api/v1/admin/articles/batch-delete
- Auth Required: Yes
- Request: `{ "ids": ["uuid"] }`
- Response: affected count
- Error Codes: VALIDATION_ERROR, AUTH_UNAUTHORIZED

## Public/Admin - Moments
- ID: PUB-MOMENT-001
- Method: GET
- Path: /api/v1/public/moments
- Auth Required: No
- Request: query `page`, `pageSize`, `locale` (optional)
- Response: paginated moments
- Error Codes: None

- ID: ADM-MOMENT-001
- Method: POST
- Path: /api/v1/admin/moments
- Auth Required: Yes
- Request: moment create payload
- Response: created moment id
- Error Codes: VALIDATION_ERROR, AUTH_UNAUTHORIZED

- ID: ADM-MOMENT-002
- Method: GET
- Path: /api/v1/admin/moments
- Auth Required: Yes
- Request: query `page`, `pageSize`, `status`
- Response: paginated moments (includes non-published records)
- Error Codes: AUTH_UNAUTHORIZED

- ID: ADM-MOMENT-003
- Method: PUT
- Path: /api/v1/admin/moments/{id}
- Auth Required: Yes
- Request: moment update payload
- Response: updated moment object
- Error Codes: MOMENT_NOT_FOUND, VALIDATION_ERROR

- ID: ADM-MOMENT-004
- Method: DELETE
- Path: /api/v1/admin/moments/{id}
- Auth Required: Yes
- Request: None
- Response: success boolean
- Error Codes: MOMENT_NOT_FOUND, AUTH_UNAUTHORIZED

- ID: ADM-MOMENT-005
- Method: POST
- Path: /api/v1/admin/moments/batch-status
- Auth Required: Yes
- Request: `{ "ids": ["uuid"], "status": "draft|published|scheduled" }`
- Response: affected count
- Error Codes: VALIDATION_ERROR, AUTH_UNAUTHORIZED

- ID: ADM-MOMENT-006
- Method: POST
- Path: /api/v1/admin/moments/batch-delete
- Auth Required: Yes
- Request: `{ "ids": ["uuid"] }`
- Response: affected count
- Error Codes: VALIDATION_ERROR, AUTH_UNAUTHORIZED

## Public/Admin - Comments
- ID: PUB-COMMENT-001
- Method: GET
- Path: /api/v1/public/comments/content/{contentType}/{contentId}
- Auth Required: No
- Request: query `page`, `pageSize`
- Response: paginated threaded comments, each row is a root comment node with recursive `children`
- Error Codes: None

- ID: PUB-COMMENT-002
- Method: GET
- Path: /api/v1/public/comments/{id}/children
- Auth Required: No
- Request: query `page`, `pageSize`
- Response: paginated child comments
- Error Codes: None

- ID: PUB-COMMENT-003
- Method: GET
- Path: /api/v1/public/comments/content/{contentType}/{contentId}/total
- Auth Required: No
- Request: None
- Response: total approved comments count
- Error Codes: None

- ID: PUB-COMMENT-004
- Method: POST
- Path: /api/v1/public/comments
- Auth Required: No
- Request: comment create payload including `contentType` and `contentId`
- Response: created comment id
- Error Codes: VALIDATION_ERROR

- ID: PUB-MOMENT-002
- Method: GET
- Path: /api/v1/public/moments/{id}
- Auth Required: No
- Request: query `locale`
- Response: published moment detail
- Error Codes: MOMENT_NOT_FOUND

- ID: ADM-COMMENT-001
- Method: GET
- Path: /api/v1/admin/comments
- Auth Required: Yes
- Request: query `page`, `pageSize`, `status`
- Response: paginated comments
- Error Codes: AUTH_UNAUTHORIZED

- ID: ADM-COMMENT-002
- Method: PUT
- Path: /api/v1/admin/comments/{id}
- Auth Required: Yes
- Request:
```json
{
  "status": "approved",
  "isPinned": "1"
}
```
- Response: updated comment object
- Error Codes: COMMENT_NOT_FOUND, VALIDATION_ERROR

## Public/Admin - Links
- ID: PUB-LINK-001
- Method: POST
- Path: /api/v1/public/links/submissions
- Auth Required: No
- Request: link submission payload
- Response: submission id
- Error Codes: VALIDATION_ERROR, RATE_LIMITED, LINK_SUBMISSION_DISABLED

- ID: PUB-LINK-002
- Method: GET
- Path: /api/v1/public/links
- Auth Required: No
- Request: None
- Response: approved links list
- Error Codes: None

- ID: ADM-LINK-001
- Method: GET
- Path: /api/v1/admin/links
- Auth Required: Yes
- Request: query `reviewStatus`, `page`, `pageSize`
- Response: paginated link submissions
- Error Codes: AUTH_UNAUTHORIZED

- ID: ADM-LINK-002
- Method: PATCH
- Path: /api/v1/admin/links/{id}/review
- Auth Required: Yes
- Request:
```json
{
  "reviewStatus": "approved",
  "reviewNote": "Looks good",
  "relatedArticleId": "uuid"
}
```
- Response: success boolean
- Error Codes: LINK_NOT_FOUND, VALIDATION_ERROR

## Public - Taxonomy
- ID: PUB-TAX-001
- Method: GET
- Path: /api/v1/public/categories
- Auth Required: No
- Request: None
- Response: category list
- Error Codes: None

- ID: PUB-TAX-002
- Method: GET
- Path: /api/v1/public/tags
- Auth Required: No
- Request: None
- Response: tag list
- Error Codes: None

## Public - Site
- ID: PUB-SITE-001
- Method: GET
- Path: /api/v1/public/site/settings
- Auth Required: No
- Request: None
- Response:
```json
{
  "siteName": "Ancy Blog",
  "avatarUrl": "https://...",
  "faviconUrl": "https://...",
  "heroIntroMd": "Hi, I build things.",
  "defaultLocale": "en",
  "commentEnabled": true,
  "commentRequireApproval": false,
  "linkSubmissionEnabled": true,
  "siteDescription": "...",
  "seoKeywords": "...",
  "ogImageUrl": "https://..."
}
```
- Error Codes: None
- Notes: Redis cache recommended (`site:settings:{locale}`).

- ID: PUB-SITE-002
- Method: GET
- Path: /api/v1/public/site/footer
- Auth Required: No
- Request: None
- Response: grouped footer items by `rowNum` with ordering.
- Error Codes: None
- Notes: Redis cache recommended (`site:footer:{locale}`).

- ID: PUB-SITE-003
- Method: GET
- Path: /api/v1/public/site/social-links
- Auth Required: No
- Request: None
- Response: enabled social links ordered by `orderNum`.
- Error Codes: None
- Notes: Redis cache recommended (`site:social:{locale}`).

- ID: PUB-SITE-004
- Method: GET
- Path: /api/v1/public/site/nav
- Auth Required: No
- Request: None
- Response: top navigation config with dropdown metadata.
- Error Codes: None
- Notes: Redis cache recommended (`site:nav:{locale}`).

- ID: PUB-SITE-005
- Method: GET
- Path: /api/v1/public/site/slots/{slotKey}
- Auth Required: No
- Request: query `limit` (optional)
- Response: ordered slot content list (articles/moments)
- Error Codes: SLOT_NOT_FOUND
- Notes: Redis cache recommended (`site:slot:{slotKey}:{locale}`).

## Public - Timeline
- ID: PUB-TL-001
- Method: GET
- Path: /api/v1/public/timeline
- Auth Required: No
- Request: query `page`, `pageSize`, `locale` (optional)
- Response: mixed timeline feed (`article` + `moment`) sorted by publish time desc.
- Item Fields:
  - `contentType`: `article | moment`
  - `id`
  - `title`, `summary`, `slug` for article rows
  - `categorySlug`, `categoryName` for article rows
  - `content` for moment rows
  - `publishedAt`
- Error Codes: None

## Admin - Site
- ID: ADM-SITE-001
- Method: PUT
- Path: /api/v1/admin/site/settings
- Auth Required: Yes
- Request: site settings payload
- Response: success boolean
- Error Codes: VALIDATION_ERROR, AUTH_UNAUTHORIZED
- Notes:
  - Supports partial update payloads; omitted fields keep their previous values.
  - Supported fields include `siteName`, `avatarUrl`, `faviconUrl`, `heroIntroMd`, locale, comment policy, friend-link submission policy, and SEO metadata.
  - Invalidate `site:settings:{locale}` cache after successful update.

- ID: ADM-SITE-002
- Method: POST
- Path: /api/v1/admin/site/footer-items
- Auth Required: Yes
- Request: footer item create payload
- Response: created item id
- Error Codes: VALIDATION_ERROR, AUTH_UNAUTHORIZED

- ID: ADM-SITE-003
- Method: PUT
- Path: /api/v1/admin/site/footer-items/{id}
- Auth Required: Yes
- Request: footer item update payload
- Response: success boolean
- Error Codes: FOOTER_ITEM_NOT_FOUND, VALIDATION_ERROR

- ID: ADM-SITE-004
- Method: DELETE
- Path: /api/v1/admin/site/footer-items/{id}
- Auth Required: Yes
- Request: None
- Response: success boolean
- Error Codes: FOOTER_ITEM_NOT_FOUND
- Notes: Invalidate `site:footer:{locale}` cache on create/update/delete.

- ID: ADM-SITE-005
- Method: POST
- Path: /api/v1/admin/site/social-links
- Auth Required: Yes
- Request: social link create payload
- Response: created item id
- Error Codes: VALIDATION_ERROR, AUTH_UNAUTHORIZED

- ID: ADM-SITE-006
- Method: PUT
- Path: /api/v1/admin/site/social-links/{id}
- Auth Required: Yes
- Request: social link update payload
- Response: success boolean
- Error Codes: SOCIAL_LINK_NOT_FOUND, VALIDATION_ERROR

- ID: ADM-SITE-007
- Method: DELETE
- Path: /api/v1/admin/site/social-links/{id}
- Auth Required: Yes
- Request: None
- Response: success boolean
- Error Codes: SOCIAL_LINK_NOT_FOUND
- Notes: Invalidate `site:social:{locale}` cache on create/update/delete.

- ID: ADM-SITE-008
- Method: POST
- Path: /api/v1/admin/site/nav-items
- Auth Required: Yes
- Request: nav item create payload
```json
{
  "parentId": "uuid (optional, for child nav items)",
  "name": "string",
  "key": "string",
  "type": "menu | dropdown | link",
  "targetType": "route | category | slot | external | article",
  "targetValue": "string (optional)",
  "orderNum": 1,
  "enabled": true
}
```
- Response: created nav item id
- Error Codes: VALIDATION_ERROR, AUTH_UNAUTHORIZED
- Notes: Set `parentId` to create child nav items. Frontend auto-injects category list when `targetType=category` with no children.

- ID: ADM-SITE-009
- Method: PUT
- Path: /api/v1/admin/site/nav-items/{id}
- Auth Required: Yes
- Request: nav item update payload
- Response: success boolean
- Error Codes: NAV_ITEM_NOT_FOUND, VALIDATION_ERROR

- ID: ADM-SITE-010
- Method: DELETE
- Path: /api/v1/admin/site/nav-items/{id}
- Auth Required: Yes
- Request: None
- Response: success boolean
- Error Codes: NAV_ITEM_NOT_FOUND
- Notes: Invalidate `site:nav:{locale}` cache on create/update/delete.

- ID: ADM-SITE-011
- Method: GET
- Path: /api/v1/admin/site/slots
- Auth Required: Yes
- Request: None
- Response: slot definition list
- Error Codes: AUTH_UNAUTHORIZED

- ID: ADM-SITE-012
- Method: POST
- Path: /api/v1/admin/site/slots
- Auth Required: Yes
- Request: slot definition payload
- Response: created slot id
- Error Codes: VALIDATION_ERROR, AUTH_UNAUTHORIZED

- ID: ADM-SITE-013
- Method: GET
- Path: /api/v1/admin/site/slots/{slotKey}/items
- Auth Required: Yes
- Request: None
- Response: slot item list
- Error Codes: SLOT_NOT_FOUND

- ID: ADM-SITE-014
- Method: POST
- Path: /api/v1/admin/site/slots/{slotKey}/items
- Auth Required: Yes
- Request: slot item payload (`contentType`, `contentId`, `orderNum`, `enabled`)
- Response: created slot item id
- Error Codes: SLOT_NOT_FOUND, VALIDATION_ERROR

- ID: ADM-SITE-015
- Method: DELETE
- Path: /api/v1/admin/site/slots/{slotKey}/items/{id}
- Auth Required: Yes
- Request: None
- Response: success boolean
- Error Codes: SLOT_ITEM_NOT_FOUND
- Notes: Invalidate `site:slot:{slotKey}:{locale}` cache on slot updates.

## Admin - Upload
- ID: ADM-UPLOAD-001
- Method: POST
- Path: /api/v1/admin/upload/image
- Auth Required: Yes
- Request: `multipart/form-data` with `file`
- Response:
```json
{
  "key": "uploads/images/202603/uuid.png",
  "url": "https://cdn.example.com/uploads/images/202603/uuid.png"
}
```
- Error Codes: UPLOAD_NOT_CONFIGURED, VALIDATION_ERROR, UPLOAD_FAILED
- Notes: Uses Cloudflare R2 when configured.

## Admin - Integrations
- ID: ADM-INT-001
- Method: GET
- Path: /api/v1/admin/integrations
- Auth Required: Yes
- Request: query `providerType` (optional)
- Response: integration provider list (secret fields masked)
- Error Codes: AUTH_UNAUTHORIZED

- ID: ADM-INT-002
- Method: PUT
- Path: /api/v1/admin/integrations/{providerKey}
- Auth Required: Yes
- Request: provider config payload (`enabled`, `configJson`, `metaJson`)
- Response: success boolean
- Error Codes: PROVIDER_NOT_FOUND, VALIDATION_ERROR
- Notes: Supports both `cloudflare_r2` and `openai_compatible` in one management flow.

- ID: ADM-INT-003
- Method: POST
- Path: /api/v1/admin/integrations/{providerKey}/test
- Auth Required: Yes
- Request: optional test payload
- Response:
```json
{
  "ok": true,
  "message": "connection ok",
  "latencyMs": 120
}
```
- Error Codes: PROVIDER_NOT_FOUND, PROVIDER_TEST_FAILED

## Admin - Translation
- ID: ADM-TR-001
- Method: POST
- Path: /api/v1/admin/translations/jobs
- Auth Required: Yes
- Request:
```json
{
  "sourceType": "article",
  "sourceId": "uuid",
  "sourceLocale": "zh-CN",
  "targetLocale": "en-US",
  "providerKey": "openai_compatible",
  "modelName": "gpt-4.1-mini",
  "maxRetries": 3,
  "autoPublish": true,
  "publishAt": "2026-03-05T09:00:00Z"
}
```
- Response: created job id
- Error Codes: VALIDATION_ERROR, PROVIDER_NOT_FOUND

- ID: ADM-TR-002
- Method: GET
- Path: /api/v1/admin/translations/jobs
- Auth Required: Yes
- Request: query `status`, `sourceType`, `sourceId`, `page`, `pageSize`
- Response: paginated translation jobs
- Error Codes: AUTH_UNAUTHORIZED

- ID: ADM-TR-003
- Method: GET
- Path: /api/v1/admin/translations/jobs/{id}
- Auth Required: Yes
- Request: None
- Response: translation job detail
- Error Codes: TRANSLATION_JOB_NOT_FOUND

- ID: ADM-TR-004
- Method: POST
- Path: /api/v1/admin/translations/jobs/{id}/retry
- Auth Required: Yes
- Request: None
- Response: queued translation job object
- Error Codes: TRANSLATION_JOB_NOT_FOUND, VALIDATION_ERROR
- Notes: Manual retry resets `retryCount` to `0` and re-queues job immediately.

- ID: ADM-TR-005
- Method: GET
- Path: /api/v1/admin/translations/contents
- Auth Required: Yes
- Request: query `sourceType` (`article|moment`, required), `sourceId` (optional), `locale` (optional), `page`, `pageSize`
- Response: paginated translation content rows
- Error Codes: VALIDATION_ERROR, AUTH_UNAUTHORIZED

- ID: ADM-TR-006
- Method: GET
- Path: /api/v1/admin/translations/contents/{sourceType}/{sourceId}/{locale}
- Auth Required: Yes
- Request: None
- Response: translation content detail
- Error Codes: TRANSLATION_CONTENT_NOT_FOUND, VALIDATION_ERROR

- ID: ADM-TR-007
- Method: PUT
- Path: /api/v1/admin/translations/contents
- Auth Required: Yes
- Request:
```json
{
  "sourceType": "article",
  "sourceId": "uuid",
  "locale": "en-US",
  "title": "Translated title",
  "summary": "Translated summary",
  "content": "manual translation override",
  "status": "draft",
  "publishedAt": "2026-03-05T09:00:00Z",
  "translatedByJobId": "uuid"
}
```
- Response: saved translation content row
- Error Codes: VALIDATION_ERROR
- Notes: Supports manual correction/override after machine translation.

## Admin - AI Assist
- ID: ADM-AI-001
- Method: POST
- Path: /api/v1/admin/ai/summary
- Auth Required: Yes
- Request:
```json
{
  "title": "Go Testing Notes",
  "content": "long article content...",
  "providerKey": "openai_compatible",
  "modelName": "gpt-4.1-mini",
  "maxLength": 180
}
```
- Response:
```json
{
  "summary": "short summary text",
  "fallbackUsed": false
}
```
- Error Codes: VALIDATION_ERROR, INTERNAL_ERROR
- Notes: Non-streaming API. Falls back to local truncation when provider is unavailable or invocation fails.

- ID: ADM-AI-002
- Method: POST
- Path: /api/v1/admin/ai/slug
- Auth Required: Yes
- Request:
```json
{
  "title": "Go Testing Notes",
  "providerKey": "openai_compatible",
  "modelName": "gpt-4.1-mini"
}
```
- Response:
```json
{
  "slug": "go-testing-notes",
  "fallbackUsed": false
}
```
- Error Codes: VALIDATION_ERROR, INTERNAL_ERROR
- Notes: Slug output is normalized and de-duplicated (`-2`, `-3`, ...).

## System
- ID: SYS-001
- Method: GET
- Path: /healthz
- Auth Required: No
- Request: None
- Response:
```json
{
  "code": "OK",
  "message": "service is healthy",
  "data": { "status": "up" }
}
```
- Error Codes: None
- Notes: Used by local and deployment health probes.

## Translation Policy (Admin)

- ID: ADM-TRANS-POLICY-001
- Method: GET
- Path: /api/v1/admin/site/translation-policy
- Auth Required: Yes
- Request: None
- Response:
```json
{
  "code": "OK",
  "message": "success",
  "data": {
    "enabled": true,
    "targetLocales": ["en-US", "ja-JP"],
    "providerKey": "openai_compatible",
    "autoPublish": false
  }
}
```
- Error Codes: AUTH_REQUIRED
- Notes: Returns defaults if no policy has been saved yet.

- ID: ADM-TRANS-POLICY-002
- Method: PUT
- Path: /api/v1/admin/site/translation-policy
- Auth Required: Yes
- Request:
```json
{
  "enabled": true,
  "targetLocales": ["en-US", "ja-JP"],
  "providerKey": "openai_compatible",
  "autoPublish": false
}
```
- Response: Same shape as GET.
- Error Codes: AUTH_REQUIRED, VALIDATION_ERROR, UPDATE_FAILED
- Notes:
  - `targetLocales` must use BCP-47 locale tags (e.g. `en-US`, `ja-JP`, `ko-KR`, `fr-FR`).
  - `providerKey` must match an existing integration provider key (e.g. `openai_compatible`).
  - When `enabled=true`, saving a published article (create or update) automatically enqueues
    translation jobs for each `targetLocale`. Jobs use `autoPublish` from this policy.
  - For articles with TipTap JSON content, the worker extracts text leaf nodes, translates them
    in a single LLM batch call, then reconstructs the JSON with translated text.
  - Policy is stored as JSONB in the `site_settings` table (migration 000006).

---

## Category & Tag Management (Admin)

### ADM-CAT-001 — Create Category
- Method: POST
- Path: `/api/v1/admin/categories`
- Auth Required: Yes
- Request: `{ "name": "Tech", "slug": "tech" }`
- Response: `{ "id": "uuid", "name": "Tech", "slug": "tech" }`
- Error Codes: AUTH_REQUIRED, VALIDATION_ERROR

### ADM-CAT-002 — Delete Category
- Method: DELETE
- Path: `/api/v1/admin/categories/:id`
- Auth Required: Yes
- Response: `{}`
- Error Codes: AUTH_REQUIRED, NOT_FOUND

### ADM-TAG-001 — Create Tag
- Method: POST
- Path: `/api/v1/admin/tags`
- Auth Required: Yes
- Request: `{ "name": "Go", "slug": "go" }`
- Response: `{ "id": "uuid", "name": "Go", "slug": "go" }`
- Error Codes: AUTH_REQUIRED, VALIDATION_ERROR

### ADM-TAG-002 — Delete Tag
- Method: DELETE
- Path: `/api/v1/admin/tags/:id`
- Auth Required: Yes
- Response: `{}`
- Error Codes: AUTH_REQUIRED, NOT_FOUND

---

## Account Management (Admin)

### ADM-AUTH-PWD-001 — Change Admin Password
- Method: PUT
- Path: `/api/v1/admin/auth/password`
- Auth Required: Yes
- Request: `{ "oldPassword": "current", "newPassword": "newpass123" }`
- Response: `{}`
- Error Codes: AUTH_REQUIRED, VALIDATION_ERROR (wrong old password / too short)
- Notes:
  - Minimum new password length: 6 characters.
  - All existing sessions are invalidated immediately; re-login required.
  - Password is stored as bcrypt hash in `site_settings.admin_password_hash` (migration 000010).

---

## Site Settings Extensions

### ADM-SITE-SET-001 (Updated) — Update Site Settings
- New fields added (migration 000008):
  - `commentEnabled` (bool) — global switch to allow/disallow new comments.
  - `commentRequireApproval` (bool) — if true, new comments default to `pending` status.
  - `siteDescription` (string) — used for SEO meta description.
  - `seoKeywords` (string) — comma-separated SEO keywords.
  - `ogImageUrl` (string) — default Open Graph image URL.
- New field added (migration 000018):
  - `linkSubmissionEnabled` (bool) — global switch to allow/disallow public friend-link submissions.
- Update semantics:
  - This singleton settings endpoint uses patch-style merge semantics.
  - Omitted fields keep their previous values instead of being reset.

### Article Pin / Featured (migration 000009)
- `isPinned` (bool) — pinned articles sort first in `ListArticles`; shown as badge in admin.
- `isFeatured` (bool) — featured flag for frontend showcase use; shown as badge in admin.
