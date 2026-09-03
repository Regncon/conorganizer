# Access Pages and User Context Boundary Design

**Date:** 2026-08-30
**Revised:** 2026-09-01
**Status:** Approved

## Scope

This design moves complete 401 and 403 responses out of service packages, simplifies the access middleware contracts, and removes the menu route's temporary user-info callback.

The change intentionally keeps the current `layouts.Base(title, userInfo, db, logger, children)` signature and full server-side menu rendering.

## Problem

`service/userctx` currently does two unrelated jobs:

1. It derives `requestctx.UserRequestInfo` and gates authenticated routes.
2. It renders complete 401 and 403 HTML responses through `layouts.Base`.

Because `layouts.Base` renders `components/header`, directly importing `userctx` from the header currently creates this cycle:

```text
header
↓
userctx
↓
layouts
↓
header
```

`SetupMenuRoute` avoids the cycle by receiving `GetUserRequestInfo` as a function argument. That callback is not a genuine configurable dependency; it exists only to work around package ownership.

The access presentation is also spread across `service/userctx` and `service/authctx`. Both pages duplicate the same styles, while `authctx.RequireAdmin` uses a functional-option configuration even though production has one call site and always supplies the HTML forbidden handler.

## Decision

Create one focused `pages/access` package:

```text
pages/access/
├── access.go
├── access.templ
└── access_test.go
```

- `access.templ` keeps the shared styles and both small components together: `Unauthorized` and `Forbidden`.
- `access.go` owns the complete HTTP 401 and 403 handlers.
- `access_test.go` covers status codes, copy, and links for both responses.

Keep the access decisions in the service packages:

- `userctx.UserMiddleware` decides whether an authenticated route may continue and delegates the complete 401 response to a required `http.HandlerFunc`.
- `authctx.RequireAdmin` decides whether an administrator route may continue and delegates the complete 403 response to a required `http.HandlerFunc`.

Use direct required arguments rather than optional configuration:

```go
func UserMiddleware(
	logger *slog.Logger,
	unauthorizedHandler http.HandlerFunc,
) func(http.Handler) http.Handler

func RequireAdmin(
	logger *slog.Logger,
	forbiddenHandler http.HandlerFunc,
) func(http.Handler) http.Handler
```

Delete `requireAdminConfig`, `RequireAdminOption`, `WithForbiddenHandler`, and the unused default plain-text forbidden response. The application has one production `RequireAdmin` call and always provides its response handler.

## Package Responsibilities

### `service/userctx`

Owns:

- deriving `requestctx.UserRequestInfo` from authentication context;
- allowing or rejecting authenticated routes;
- user-ID lookup helpers.

Does not own:

- Templ access pages;
- layout rendering;
- database access for rendering the header;
- the complete 403 response.

The middleware keeps the current request logging, calls the unauthorized handler when the user is logged out, and never writes a second response itself.

### `service/authctx`

Owns:

- reading authentication claims;
- checking the administrator role;
- allowing or rejecting administrator routes.

It does not own the forbidden page or choose a default presentation.

### `pages/access`

Owns:

- the shared access-denied presentation;
- the 401 page content and complete response;
- the 403 page content and complete response;
- render-failure logging for those responses.

The handlers retain the database argument required by the current layout:

```go
func UnauthorizedHandler(db *sql.DB, logger *slog.Logger) http.HandlerFunc
func ForbiddenHandler(db *sql.DB, logger *slog.Logger) http.HandlerFunc
```

## Accepted Layout Dependency

The header renders authoritative billettholder data in the first server response. `layouts.Base` therefore receives the database and logger used by `header.Menu`, and complete page handlers pass those dependencies through to the layout.

This explicit dependency is retained deliberately:

- SQLite access is inexpensive for the expected event traffic.
- The first response contains the complete avatar, ticket type, and switcher without a loading state.
- `/menu/api` remains responsible for live updates after the initial render.
- Request context and package globals are not used to hide infrastructure dependencies.
- A one-off renderer would replace the dependency rather than remove it and would be inconsistent with the rest of the application.

If page rendering is redesigned later, it should be an application-wide change covering data loading, view models, layout rendering, statuses, and render-error handling consistently. That broader redesign is outside this plan.

## Router Wiring

`router.go` remains the composition boundary:

```go
isLoggedInRouter := router.With(
	userctx.UserMiddleware(logger, access.UnauthorizedHandler(db, logger)),
)

header.SetupMenuRoute(isLoggedInRouter, liveManager, db, logger)

routerAdmin := isLoggedInRouter.With(
	authctx.RequireAdmin(logger, access.ForbiddenHandler(db, logger)),
)
```

Once `userctx` no longer imports `layouts`, `components/header` can directly use `userctx.GetUserRequestInfo`:

```go
func SetupMenuRoute(
	router chi.Router,
	liveManager *live.Manager,
	db *sql.DB,
	logger *slog.Logger,
)
```

The injected `getUserRequestInfo func(context.Context) requestctx.UserRequestInfo` parameter is removed.

## Request Flows

### HTTP 401

```text
request
↓
userctx.UserMiddleware checks request context
↓
pages/access.UnauthorizedHandler
↓
HTTP 401 + layouts.Base + Unauthorized
```

### HTTP 403

```text
authenticated administrator route request
↓
authctx.RequireAdmin checks role
↓
pages/access.ForbiddenHandler
↓
HTTP 403 + layouts.Base + Forbidden
```

### Live menu route

```text
/menu/api
↓
header.SetupMenuRoute
↓
userctx.GetUserRequestInfo(ctx)
↓
MenuBillettholderLive
```

## Final Dependency Direction

```text
pages/access
├── layouts
└── service/userctx

layouts
└── components/header

components/header
└── service/userctx

service/userctx
└── service/authctx
```

There is no path from `service/userctx` back to `layouts`, so the header can depend directly on `userctx` without a cycle.

## Behavior Preserved

- Logged-out requests to protected routes return HTTP 401.
- Non-administrators requesting administrator routes return HTTP 403.
- Existing Norwegian text and links remain unchanged.
- The 401 response renders a logged-out header.
- The 403 response renders using the current request's user information.
- Access-page render failures retain structured request ID and path logging.
- The menu stream still watches `live.BucketBillettholders`.

## Removed Code

- `service/userctx/unauthenticated.templ`
- `service/userctx/unauthenticated_test.go`
- `service/authctx/forbidden.templ`
- `service/authctx/forbidden_test.go`
- `userctx.AdminForbiddenHandler`
- `authctx.requireAdminConfig`
- `authctx.RequireAdminOption`
- `authctx.WithForbiddenHandler`
- the default plain-text forbidden handler
- the redundant `if !userInfo.IsLoggedIn` branch
- the no-op `r.WithContext(ctx)` call
- the `SetupMenuRoute` user-info callback parameter

## Constraints

- Keep `requestctx.UserRequestInfo` in `service/requestctx`.
- Keep `layouts.Base` and `header.Menu` database/logger parameters to preserve complete first-response rendering.
- Keep all access-page text, links, titles, and status codes unchanged.
- Do not change billettholder persistence, Datastar signals, or live-bucket behavior.
- Follow `domeneordbok.md`; this change introduces no new domain terms.
- Do not manually edit generated `*_templ.go` files.
- Templ generation and tests remain user-owned.
- Keep this implementation as one commit.

## Acceptance Criteria

- `pages/access` contains exactly one Go source, one Templ source, and one test source.
- `service/userctx` no longer imports `layouts` or owns presentation components.
- `service/authctx` no longer owns the forbidden component or optional forbidden-handler configuration.
- Both middleware functions receive required response handlers directly.
- `router.go` composes both handlers.
- `components/header/menu.templ` directly calls `userctx.GetUserRequestInfo`.
- `SetupMenuRoute` has four parameters and no user-info callback.
- 401/403 behavior and menu live behavior remain unchanged.
- The code compiles without an import cycle after user-owned generation.
