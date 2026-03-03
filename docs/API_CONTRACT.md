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
