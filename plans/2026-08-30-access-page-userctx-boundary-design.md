# Access Pages and User Context Boundary Design

**Date:** 2026-08-30  
**Status:** Approved

## Problem

`service/userctx` currently has two responsibilities:

1. It derives `requestctx.UserRequestInfo` from authentication data stored in the request context.
2. It renders complete 401 and 403 HTML responses through `layouts.Base`.

The second responsibility makes `userctx` import `layouts`. Because `layouts.Base` renders `components/header`, importing `userctx` from the menu creates this cycle:

```text
header
↓
userctx
↓
layouts
↓
header
```

`SetupMenuRoute` currently avoids the cycle by receiving `GetUserRequestInfo` as a function parameter. That injection is a temporary package-boundary workaround rather than a dependency the menu genuinely needs to abstract.

## Decision

Move the complete 401 and 403 response rendering into a new `pages/access` package. Keep authentication and user-context interpretation in `service/userctx` and `service/authctx`.

Keep `requestctx.UserRequestInfo` in `service/requestctx` in this migration. Moving the type to `userctx` would touch many unrelated pages, components, and tests and is explicitly outside this design.

Once `userctx` no longer imports `layouts`, restore direct use of `userctx.GetUserRequestInfo` inside `SetupMenuRoute` at the bottom of `components/header/menu.templ`.

## Package Responsibilities

### `service/userctx`

`userctx` owns:

- deriving `requestctx.UserRequestInfo` from the authenticated request context;
- deciding whether a request may continue through `UserMiddleware`;
- user-ID lookup helpers.

`userctx` does not own:

- Templ access-denied components;
- `layouts.Base` rendering;
- database or layout dependencies for the unauthenticated middleware response;
- the full 403 response.

The middleware receives the response behavior through this interface:

```go
func UserMiddleware(
	logger *slog.Logger,
	unauthorizedHandler http.HandlerFunc,
) func(http.Handler) http.Handler
```

When `GetUserRequestInfo` reports an unauthenticated user, the middleware logs the existing request information, invokes `unauthorizedHandler`, and does not call the protected handler.

### `service/authctx`

`authctx` continues to own session interpretation and administrator-role checks. `RequireAdmin` and `WithForbiddenHandler` remain unchanged.

The `Forbidden` Templ component moves out of `authctx`, because presentation is not authentication-context logic.

### `pages/access`

The new package owns complete 401 and 403 responses:

```text
pages/access/
├── access.go
├── access.templ
├── unauthorized.templ
├── forbidden.templ
├── access_test.go
├── unauthorized_test.go
└── forbidden_test.go
```

- `access.go` exposes `UnauthorizedHandler` and `ForbiddenHandler`.
- `access.templ` contains shared access-denied presentation styles.
- `unauthorized.templ` contains the HTTP 401 page content.
- `forbidden.templ` contains the HTTP 403 page content.
- The tests cover status codes, response content, links, and the moved components.

The public handler interfaces are:

```go
func UnauthorizedHandler(db *sql.DB, logger *slog.Logger) http.HandlerFunc
func ForbiddenHandler(db *sql.DB, logger *slog.Logger) http.HandlerFunc
```

`UnauthorizedHandler` renders `layouts.Base` with an empty `requestctx.UserRequestInfo`, preserving the existing logged-out header behavior.

`ForbiddenHandler` obtains the authenticated request information with `userctx.GetUserRequestInfo(r.Context())` and passes it to `layouts.Base`, preserving the existing logged-in header behavior.

Both handlers set the HTTP status before rendering and include the request ID and path if rendering fails.

## Router Wiring

`router.go` composes the middleware and presentation handlers:

```go
isLoggedInRouter := router.With(
	userctx.UserMiddleware(logger, access.UnauthorizedHandler(db, logger)),
)

header.SetupMenuRoute(isLoggedInRouter, liveManager, db, logger)

routerAdmin := isLoggedInRouter.With(
	authctx.RequireAdmin(
		logger,
		authctx.WithForbiddenHandler(access.ForbiddenHandler(db, logger)),
	),
)
```

This keeps dependency composition at the application boundary instead of placing page rendering inside a service package.

## Request Flows

### HTTP 401

```text
Request
↓
AuthMiddleware populates authentication context
↓
UserMiddleware calls GetUserRequestInfo
↓
User is not logged in
↓
pages/access.UnauthorizedHandler
↓
HTTP 401 + layouts.Base + Unauthorized component
```

### HTTP 403

```text
Authenticated request to an administrator route
↓
RequireAdmin checks the role
↓
User is not an administrator
↓
pages/access.ForbiddenHandler
↓
HTTP 403 + layouts.Base + Forbidden component
```

### Live Menu Route

```text
/menu/api
↓
SetupMenuRoute
↓
userctx.GetUserRequestInfo(ctx)
↓
MenuBillettholderLive
```

The final `SetupMenuRoute` signature is:

```go
func SetupMenuRoute(
	router chi.Router,
	liveManager *live.Manager,
	db *sql.DB,
	logger *slog.Logger,
)
```

The temporary `getUserRequestInfo func(context.Context) requestctx.UserRequestInfo` parameter is removed.

## Final Dependency Direction

```text
pages/access
↓
layouts
↓
components/header
↓
service/userctx
↓
service/authctx
```

There is no dependency from `userctx` back to `layouts`, so importing `userctx` from the menu is safe.

## Behavior Preserved

- Unauthenticated protected requests return HTTP 401.
- Non-administrator requests to administrator routes return HTTP 403.
- The existing Norwegian access-denied text and links remain unchanged.
- 401 renders the logged-out header.
- 403 renders using the current request's user information.
- Existing request logging remains, with access-page render failures owned by the `access` component.
- The menu live stream still watches `live.BucketBillettholders` and renders `MenuBillettholderLive`.

## Constraints

- Do not move `requestctx.UserRequestInfo` in this migration.
- Do not change the billettholder persistence or live-bucket behavior.
- Do not manually edit generated `*_templ.go` files.
- Codex does not run Templ generation; the user performs the project's normal generation workflow.
- Codex does not run tests; the user performs test verification.
- Do not create commits unless the user explicitly asks for them later.

## Acceptance Criteria

- `service/userctx` no longer imports `layouts`.
- `service/userctx` no longer contains access-denied Templ components or a forbidden page handler.
- `pages/access` owns the 401 and 403 components and complete response handlers.
- `UserMiddleware` receives and invokes an `http.HandlerFunc` for unauthenticated requests.
- `router.go` wires both access handlers into the existing middleware chain.
- `components/header/menu.templ` directly calls `userctx.GetUserRequestInfo` in `SetupMenuRoute`.
- `SetupMenuRoute` no longer receives a user-info function parameter.
- The import cycle is removed without moving `UserRequestInfo` out of `requestctx`.
