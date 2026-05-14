# Plan: Move Auth Cookies to Backend (`HttpOnly`, `SameSite=Lax`) Without Descope Custom Domain

## Why We Are Doing This
- Current login stores auth cookies from browser JavaScript (`document.cookie`), so cookies are not `HttpOnly`.
- Non-`HttpOnly` auth cookies are exposed to XSS theft risk.
- With `SameSite=Strict`, browser recovery flows (for example after server downtime/offline error pages) can cause cookies to be omitted on the next navigation, creating “looks logged out until another click” behavior.
- Moving cookie issuance to backend gives `HttpOnly` cookies and consistent policy (`SameSite=Lax`) while keeping Descope UI.


## Summary
Replace browser-side auth cookie writes with a backend session-establishment endpoint so cookies become `HttpOnly` and resilient to the Chromium offline-refresh flow. Keep Descope widget login UI, but let Go own cookie issuance and policy.

## Key Implementation Changes
- Add `POST /auth/session` under existing auth routes (same subsystem as current `/auth/post-login`).
- New endpoint contract:
  - Request JSON: `{"sessionJwt":"...","refreshJwt":"..."}`
  - Server validates token pair with Descope SDK before setting cookies.
  - Response: `204 No Content` (or `401` on invalid tokens).
- Server sets both cookies (`session_token`, `refresh_token`) with:
  - `HttpOnly=true`
  - `Secure=true`
  - `SameSite=Lax`
  - `Path=/`
  - 1-year expiry/max-age (match existing behavior).
- Refactor cookie writing to one shared helper in auth package so policy is consistent across:
  - new `/auth/session`
  - middleware refresh write path
  - logout clear path
- Update login page script to:
  - stop using `document.cookie`
  - on Descope success, `fetch('/auth/session', { method: 'POST', credentials: 'same-origin' ... })`
  - then navigate to `/auth/post-login` only after successful `/auth/session`.
- Keep current middleware + `GetUserRequestInfo` flow unchanged (only source of cookies changes).

## Public/API & Interface Changes
- New backend endpoint: `POST /auth/session`.
- New client-server payload interface: `{sessionJwt, refreshJwt}` from Descope success event.
- Cookie behavior change: auth cookies become `HttpOnly` and `SameSite=Lax` (from JS-set non-HttpOnly Strict).

## Test Plan
- Unit tests for `/auth/session`:
  - valid token pair => `204` and both `Set-Cookie` headers include `HttpOnly`, `Secure`, `SameSite=Lax`.
  - invalid/missing token(s) => `401`/`400` and no auth cookies set.
- Middleware regression test:
  - existing authenticated requests still resolve `IsLoggedIn=true`.
- Logout regression test:
  - both cookies are cleared with same `Path`/policy compatibility.
- Manual scenario:
  - start app with Docker/Air, stop during request, get Chromium offline page, refresh when app is back, verify session remains usable without “magic” recovery via logo click.

## Assumptions
- Keep `Secure=true` in all environments (works on `localhost`; non-local plain HTTP envs are not supported by default).
- Keep 1-year cookie lifetime to match current behavior.
- Keep `/auth/post-login` logic unchanged; only precondition changes (backend now establishes cookies first).
- No Descope paid custom-domain features required.
