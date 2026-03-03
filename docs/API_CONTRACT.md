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
- Request: None
- Response: article detail (includes `allowComment`, AI disclosure, source info, `contentKind`)
- Error Codes: ARTICLE_NOT_FOUND

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

## Public/Admin - Moments
- ID: PUB-MOMENT-001
- Method: GET
- Path: /api/v1/public/moments
- Auth Required: No
- Request: query `page`, `pageSize`
- Response: paginated moments
- Error Codes: None

- ID: ADM-MOMENT-001
- Method: POST
- Path: /api/v1/admin/moments
- Auth Required: Yes
- Request: moment create payload
- Response: created moment id
- Error Codes: VALIDATION_ERROR, AUTH_UNAUTHORIZED

## Public/Admin - Links
- ID: PUB-LINK-001
- Method: POST
- Path: /api/v1/public/links/submissions
- Auth Required: No
- Request: link submission payload
- Response: submission id
- Error Codes: VALIDATION_ERROR, RATE_LIMITED

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
  "heroIntroMd": "Hi, I build things.",
  "defaultLocale": "en"
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
- Request: query `page`, `pageSize`
- Response: mixed timeline feed (`article` + `moment`) sorted by publish time desc.
- Error Codes: None

## Admin - Site
- ID: ADM-SITE-001
- Method: PUT
- Path: /api/v1/admin/site/settings
- Auth Required: Yes
- Request: site settings payload
- Response: success boolean
- Error Codes: VALIDATION_ERROR, AUTH_UNAUTHORIZED
- Notes: Invalidate `site:settings:{locale}` cache after successful update.

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
- Response: created nav item id
- Error Codes: VALIDATION_ERROR, AUTH_UNAUTHORIZED

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
- Method: POST
- Path: /api/v1/admin/site/slots
- Auth Required: Yes
- Request: slot definition payload
- Response: created slot id
- Error Codes: VALIDATION_ERROR, AUTH_UNAUTHORIZED

- ID: ADM-SITE-012
- Method: POST
- Path: /api/v1/admin/site/slots/{slotKey}/items
- Auth Required: Yes
- Request: slot item payload (`contentType`, `contentId`, `orderNum`, `enabled`)
- Response: created slot item id
- Error Codes: SLOT_NOT_FOUND, VALIDATION_ERROR

- ID: ADM-SITE-013
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
