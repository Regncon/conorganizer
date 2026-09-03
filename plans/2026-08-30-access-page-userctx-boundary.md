# Access Pages and User Context Boundary Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Move complete 401/403 responses out of service packages, remove one-use middleware configuration, and remove the menu route's import-cycle callback.

**Architecture:** `service/userctx` and `service/authctx` keep access decisions and delegate rejected responses to required handlers. A compact `pages/access` package keeps both related pages, their shared styles, their handlers, and their tests together. `router.go` composes the dependencies, allowing the header to call `userctx.GetUserRequestInfo` directly.

**Tech Stack:** Go 1.27, Chi middleware, Templ, `log/slog`.

**Spec:** `plans/2026-08-30-access-page-userctx-boundary-design.md`

## Global Constraints

- Keep `requestctx.UserRequestInfo` in `service/requestctx`.
- Keep the current `layouts.Base(title, userInfo, db, logger, children)` signature and complete server-side menu rendering.
- Do not introduce asynchronous initial menu loading, hidden request-context dependencies, or a one-off layout renderer.
- Preserve HTTP 401/403 statuses, page titles, Norwegian copy, links, and header login state.
- Do not change billettholder persistence, Datastar signals, or live buckets.
- Follow `domeneordbok.md`; this change introduces no new domain terms.
- Use native CSS nesting and `place-*` shorthands in the moved styles.
- Do not manually edit generated `*_templ.go` files.
- Do not run Templ generation or tests; pause for the user to run them.
- Preserve unrelated worktree changes.
- Deliver all tasks as one commit: `refactor(access): separate access responses from context services`.

---

## File Map

### Create

- `pages/access/access.templ` — shared styles plus `Unauthorized` and `Forbidden`.
- `pages/access/access.go` — complete 401 and 403 response handlers.
- `pages/access/access_test.go` — component and handler behavior tests.
- `service/userctx/userctx_test.go` — middleware delegation test.

### Modify

- `service/userctx/userctx.go` — remove presentation and accept a required unauthorized handler.
- `service/authctx/utils.go` — replace optional forbidden-handler configuration with one required handler.
- `service/authctx/require_admin_test.go` — test the direct handler contract.
- `router.go` — compose `pages/access` handlers.
- `components/header/menu.templ` — directly call `userctx.GetUserRequestInfo`.

### Delete

- `service/userctx/unauthenticated.templ`
- `service/userctx/unauthenticated_test.go`
- `service/authctx/forbidden.templ`
- `service/authctx/forbidden_test.go`

### Generated files

Do not manually edit generated files. After the source move, the user's generation workflow must create `pages/access/access_templ.go`, refresh `components/header/menu_templ.go`, and remove obsolete generated access components from the service packages.

---

### Task 1: Consolidate the access presentation

**Files:**

- Create: `pages/access/access.templ`
- Create: `pages/access/access_test.go`
- Delete: `service/userctx/unauthenticated.templ`
- Delete: `service/userctx/unauthenticated_test.go`
- Delete: `service/authctx/forbidden.templ`
- Delete: `service/authctx/forbidden_test.go`

**Interfaces:**

- Produces: `access.Unauthorized() templ.Component`
- Produces: `access.Forbidden() templ.Component`

- [ ] **Step 1: Create one Templ source for both related pages**

Create `pages/access/access.templ`:

```templ
package access

templ accessDeniedStyles() {
	<style>
		.access-denied {
			display: flex;
			flex-direction: column;
			gap: 1rem;
			place-items: center;
			place-content: center;
			max-width: 42rem;
			margin-inline: auto;
			padding: 5rem 1rem 2rem;
			background-color: var(--color-background);
			color: var(--color-text);
			font-family: var(--font-sans-serif);
			text-align: center;

			h1 {
				margin: 0;
				font-size: 2rem;
			}

			p {
				margin: 0;
				font-size: 1.2rem;
				line-height: 1.5;
			}
		}

		.access-denied-actions {
			display: flex;
			flex-wrap: wrap;
			gap: 0.75rem;
			place-content: center;
		}
	</style>
}

templ Unauthorized() {
	@accessDeniedStyles()
	<section class="access-denied" aria-labelledby="access-denied-heading">
		<h1 id="access-denied-heading">Du må logge inn</h1>
		<p>Logg inn for å se denne siden.</p>
		<div class="access-denied-actions">
			<a href="/auth" class="btn btn--primary btn-login">Logg inn</a>
			<a href="/" class="btn btn--outline">Gå til arrangementslisten</a>
		</div>
	</section>
}

templ Forbidden() {
	@accessDeniedStyles()
	<section class="access-denied" aria-labelledby="access-denied-heading">
		<h1 id="access-denied-heading">Du har ikke tilgang</h1>
		<p>Du er logget inn, men denne siden krever administratortilgang.</p>
		<div class="access-denied-actions">
			<a href="/" class="btn btn--primary">Gå til arrangementslisten</a>
		</div>
	</section>
}
```

- [ ] **Step 2: Move both component behaviors into one test source**

Create `pages/access/access_test.go`. Move the existing component assertions into these tests without changing expected copy or hrefs:

```go
func TestUnauthorized_RendersLoginAndHomeLinks(t *testing.T)
func TestForbidden_RendersAdminAccessMessageAndHomeLink(t *testing.T)
```

Use `templtest.Render`, `templtest.CollectUniqueHrefs`, and `templtest.AssertSameHrefs`. The exact expectations are:

```go
unauthorizedHrefs := []string{"/", "/auth"}
unauthorizedText := []string{
	"Du må logge inn",
	"Logg inn for å se denne siden.",
	"Logg inn",
	"Gå til arrangementslisten",
}

forbiddenHrefs := []string{"/"}
forbiddenText := []string{
	"Du har ikke tilgang",
	"Du er logget inn, men denne siden krever administratortilgang.",
	"Gå til arrangementslisten",
}
```

Keep BDD metadata and visible `// Given`, `// When`, and `// Then` sections.

- [ ] **Step 3: Delete the four replaced sources**

Delete the two old Templ files and their two package-local tests listed above. Do not delete generated Go files manually.

- [ ] **Step 4: Inspect the Task 1 diff**

Confirm that both components and their shared CSS have one owner and that the rendered copy, hrefs, classes, and accessible heading relationships are unchanged.

---

### Task 2: Add complete 401 and 403 response handlers

**Files:**

- Create: `pages/access/access.go`
- Modify: `pages/access/access_test.go`

**Interfaces:**

- Consumes: `layouts.Base`, `requestctx.UserRequestInfo`, and `userctx.GetUserRequestInfo`.
- Produces: `func UnauthorizedHandler(*sql.DB, *slog.Logger) http.HandlerFunc`.
- Produces: `func ForbiddenHandler(*sql.DB, *slog.Logger) http.HandlerFunc`.

- [ ] **Step 1: Add failing handler behavior tests**

Append two tests to `pages/access/access_test.go`:

```go
func TestUnauthorizedHandler_ReturnsCompleteUnauthorizedResponse(t *testing.T) {
	// Given
	expectedStatus := http.StatusUnauthorized
	request := httptest.NewRequest(http.MethodGet, "/profile", nil)
	recorder := httptest.NewRecorder()
	handler := UnauthorizedHandler(nil, discardAccessLogger())

	// When
	handler.ServeHTTP(recorder, request)

	// Then
	if recorder.Code != expectedStatus {
		t.Fatalf("status mismatch: expected %d, got %d", expectedStatus, recorder.Code)
	}
	if !strings.Contains(recorder.Body.String(), "Du må logge inn") {
		t.Fatalf("expected unauthorized page, got %q", recorder.Body.String())
	}
}

func TestForbiddenHandler_ReturnsCompleteForbiddenResponse(t *testing.T) {
	// Given
	expectedStatus := http.StatusForbidden
	request := httptest.NewRequest(http.MethodGet, "/admin", nil)
	recorder := httptest.NewRecorder()
	handler := ForbiddenHandler(nil, discardAccessLogger())

	// When
	handler.ServeHTTP(recorder, request)

	// Then
	if recorder.Code != expectedStatus {
		t.Fatalf("status mismatch: expected %d, got %d", expectedStatus, recorder.Code)
	}
	if !strings.Contains(recorder.Body.String(), "Du har ikke tilgang") {
		t.Fatalf("expected forbidden page, got %q", recorder.Body.String())
	}
}
```

Add BDD metadata to both tests. Add this local helper:

```go
func discardAccessLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}
```

- [ ] **Step 2: Implement both handlers**

Create `pages/access/access.go`:

```go
package access

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Regncon/conorganizer/layouts"
	"github.com/Regncon/conorganizer/service/requestctx"
	"github.com/Regncon/conorganizer/service/userctx"
	"github.com/go-chi/chi/v5/middleware"
)

func UnauthorizedHandler(db *sql.DB, logger *slog.Logger) http.HandlerFunc {
	logger = logger.With("component", "access")
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		if err := layouts.Base("Logg inn", requestctx.UserRequestInfo{}, db, logger, Unauthorized()).Render(r.Context(), w); err != nil {
			logger.Error(
				fmt.Errorf("failed to render unauthorized page: %w", err).Error(),
				"request_id", middleware.GetReqID(r.Context()),
				"path", r.URL.Path,
			)
		}
	}
}

func ForbiddenHandler(db *sql.DB, logger *slog.Logger) http.HandlerFunc {
	logger = logger.With("component", "access")
	return func(w http.ResponseWriter, r *http.Request) {
		userInfo := userctx.GetUserRequestInfo(r.Context())
		w.WriteHeader(http.StatusForbidden)
		if err := layouts.Base("Ingen tilgang", userInfo, db, logger, Forbidden()).Render(r.Context(), w); err != nil {
			logger.Error(
				fmt.Errorf("failed to render forbidden page: %w", err).Error(),
				"request_id", middleware.GetReqID(r.Context()),
				"path", r.URL.Path,
			)
		}
	}
}
```

- [ ] **Step 3: Inspect the handler contracts**

Confirm that each handler sets its status before rendering, uses the correct header login state, and logs render failures with component, request ID, and path.

---

### Task 3: Make `UserMiddleware` presentation-independent

**Files:**

- Modify: `service/userctx/userctx.go`
- Create: `service/userctx/userctx_test.go`

**Interfaces:**

- Consumes: a required `http.HandlerFunc` for rejected requests.
- Produces: `func UserMiddleware(*slog.Logger, http.HandlerFunc) func(http.Handler) http.Handler`.

- [ ] **Step 1: Add the delegation behavior test**

Create `service/userctx/userctx_test.go` with a test that supplies a custom unauthorized handler, sends an unauthenticated request, and verifies all four outcomes:

```go
expectedStatus := http.StatusUnauthorized
expectedBody := "custom unauthorized"
expectedUnauthorizedHandlerCalled := true
expectedProtectedHandlerCalled := false
```

Construct the middleware exactly as production will use it:

```go
handler := UserMiddleware(discardUserctxLogger(), unauthorizedHandler)(protectedHandler)
```

Use BDD metadata and `// Given`, `// When`, `// Then` sections. Add a local discard logger helper.

- [ ] **Step 2: Change the middleware signature and branch**

Replace the database parameter with the required handler:

```go
func UserMiddleware(logger *slog.Logger, unauthorizedHandler http.HandlerFunc) func(http.Handler) http.Handler
```

Use this complete request branch:

```go
userInfo := GetUserRequestInfo(r.Context())
requestID := middleware.GetReqID(r.Context())
if userInfo.IsLoggedIn {
	logger.Debug("User is logged in", "request_id", requestID)
	next.ServeHTTP(w, r)
	return
}

logger.Debug("User is not logged in", "request_id", requestID, "path", r.URL.Path)
unauthorizedHandler(w, r)
```

- [ ] **Step 3: Remove obsolete user-context presentation code**

Remove the `layouts` import, all inline 401 rendering, and `AdminForbiddenHandler`. Keep `database/sql` because the user-ID lookup helpers still use it.

- [ ] **Step 4: Inspect the package boundary**

Confirm `service/userctx` contains no layout import, Templ component, status write, or forbidden response handler.

---

### Task 4: Simplify the administrator middleware contract

**Files:**

- Modify: `service/authctx/utils.go`
- Modify: `service/authctx/require_admin_test.go`

**Interfaces:**

- Consumes: one required `http.HandlerFunc` for rejected requests.
- Produces: `func RequireAdmin(*slog.Logger, http.HandlerFunc) func(http.Handler) http.Handler`.

- [ ] **Step 1: Reduce the tests to the two production behaviors**

Keep and update these behaviors:

```go
func TestRequireAdmin_WhenUserIsNotAdmin_UsesForbiddenHandler(t *testing.T)
func TestRequireAdmin_WhenUserIsAdmin_AllowsRequest(t *testing.T)
```

Both construct the middleware with a handler argument:

```go
handler := RequireAdmin(discardLogger(), forbiddenHandler)(protectedHandler)
```

Delete the test for the default plain-text response because that behavior is being removed.

- [ ] **Step 2: Replace functional options with a direct dependency**

Use this signature and rejection branch:

```go
func RequireAdmin(logger *slog.Logger, forbiddenHandler http.HandlerFunc) func(http.Handler) http.Handler {
	logger = logger.With("component", "auth")
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !GetAdminFromUserToken(r.Context()) {
				logger.Warn(
					"User is not an admin",
					"request_id", middleware.GetReqID(r.Context()),
					"path", r.URL.Path,
				)
				forbiddenHandler(w, r)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
```

Delete `requireAdminConfig`, `RequireAdminOption`, `WithForbiddenHandler`, and the default handler.

- [ ] **Step 3: Inspect all call sites**

Confirm every `RequireAdmin` call supplies exactly one forbidden handler. Production currently has one call in `router.go`.

---

### Task 5: Compose access responses and remove the menu callback

**Files:**

- Modify: `router.go`
- Modify: `components/header/menu.templ`

**Interfaces:**

- Consumes: the handlers and middleware contracts from Tasks 2–4.
- Produces: `func SetupMenuRoute(chi.Router, *live.Manager, *sql.DB, *slog.Logger)`.

- [ ] **Step 1: Update router composition**

Import `pages/access` and replace the current access/menu block with:

```go
isLoggedInRouter := router.With(
	userctx.UserMiddleware(logger, access.UnauthorizedHandler(db, logger)),
)
header.SetupMenuRoute(isLoggedInRouter, liveManager, db, logger)
routerAdmin := isLoggedInRouter.With(
	authctx.RequireAdmin(logger, access.ForbiddenHandler(db, logger)),
)
```

- [ ] **Step 2: Restore the direct header dependency**

In `components/header/menu.templ`, import:

```go
"github.com/Regncon/conorganizer/service/userctx"
```

Keep the `requestctx` import because menu component signatures still use `requestctx.UserRequestInfo`.

- [ ] **Step 3: Remove the callback parameter**

Use this route signature:

```go
func SetupMenuRoute(
	router chi.Router,
	liveManager *live.Manager,
	db *sql.DB,
	logger *slog.Logger,
)
```

Use the direct context lookup in the live renderer:

```go
Render: func(ctx context.Context, r *http.Request) templ.Component {
	userInfo := userctx.GetUserRequestInfo(ctx)
	return MenuBillettholderLive(userInfo, db, logger)
},
```

- [ ] **Step 4: Inspect the final imports**

The expected dependency chain is:

```text
pages/access → layouts → header → userctx → authctx
pages/access → userctx
```

There must be no `userctx → layouts` edge.

---

### Task 6: User-owned generation, verification, and single commit

**Files:**

- Generated by user: `pages/access/access_templ.go`
- Regenerated by user: `components/header/menu_templ.go`
- Removed by normal generated-file cleanup: obsolete service-package access components.

- [ ] **Step 1: Stop for Templ generation**

Ask the user to run their normal Templ generation workflow. Do not run it for them.

- [ ] **Step 2: Inspect generated ownership**

After generation, confirm:

```text
pages/access/access_templ.go exists
components/header/menu_templ.go imports userctx
service/userctx has no generated Unauthenticated component
service/authctx has no generated Forbidden component
```

- [ ] **Step 3: Run allowed static checks**

```powershell
git diff --check
rg -n "conorganizer/layouts|AdminForbiddenHandler|func Unauthenticated" service/userctx
rg -n "requireAdminConfig|RequireAdminOption|WithForbiddenHandler|func Forbidden" service/authctx
rg -n "func SetupMenuRoute|GetUserRequestInfo" components/header/menu.templ router.go
```

Expected: `git diff --check` succeeds; the first two searches return no obsolete ownership; the final search shows direct user-context lookup and four-argument route setup.

- [ ] **Step 4: Give the user the verification commands**

```powershell
go test ./pages/access ./service/userctx ./service/authctx ./components/header
go test ./...
```

Expected: all packages compile without an import cycle and all tests pass.

- [ ] **Step 5: Create one commit after user verification**

```powershell
git add pages/access service/userctx service/authctx components/header/menu.templ components/header/menu_templ.go router.go
git commit -m "refactor(access): separate access responses from context services"
```

Do not commit task-by-task; this entire access-boundary change is one review point.

---

## Self-Review Record

- **Spec coverage:** Tasks 1–5 cover every ownership, API, behavior, and routing decision; Task 6 covers the user-owned generation boundary.
- **Placeholder scan:** No implementation choice is deferred.
- **Type consistency:** Both access handlers return `http.HandlerFunc`; both middleware functions require that exact type; `SetupMenuRoute` has four arguments in source and router wiring.
- **Scope check:** This plan intentionally retains `db` and `logger` in `layouts.Base` as the explicit cost of complete first-response menu rendering. Any replacement belongs to a future application-wide rendering design.
