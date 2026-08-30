# Access Pages and User Context Boundary Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Remove the `header → userctx → layouts → header` import cycle by moving complete 401/403 responses to `pages/access`, then restore direct `userctx.GetUserRequestInfo` use in the live menu route.

**Architecture:** `service/userctx` remains responsible for interpreting authenticated request context and gating protected requests, while `pages/access` owns HTML response rendering. `router.go` injects an `UnauthorizedHandler` into `UserMiddleware` and a `ForbiddenHandler` into `RequireAdmin`; the menu can then import `userctx` directly because `userctx` no longer imports `layouts`.

**Tech Stack:** Go, Chi middleware, Templ, `log/slog`, Datastar live streams.

**Spec:** `plans/2026-08-30-access-page-userctx-boundary-design.md`

## Global Constraints

- Keep `requestctx.UserRequestInfo` in `service/requestctx`.
- Keep all existing Norwegian access-page text and links unchanged.
- Do not change billettholder persistence, Datastar signals, or live-bucket behavior.
- Follow `domeneordbok.md`; this plan introduces no new billettholder domain terms.
- Do not manually edit generated `*_templ.go` files.
- Do not run Templ generation; pause for the user to perform it.
- Do not run tests; provide the exact relevant commands to the user instead.
- Do not create commits unless the user explicitly requests one.
- Preserve unrelated staged and unstaged changes in `components/header`, `router.go`, and `static/js/conorganizer.js`.

---

## File Map

### Create

- `pages/access/access.go` — complete HTTP 401 and 403 response handlers.
- `pages/access/access.templ` — shared access-denied styles.
- `pages/access/unauthorized.templ` — unauthenticated page content.
- `pages/access/forbidden.templ` — administrator-access-denied page content.
- `pages/access/access_test.go` — handler status and rendered-response coverage.
- `pages/access/unauthorized_test.go` — moved 401 component coverage.
- `pages/access/forbidden_test.go` — moved 403 component coverage.
- `service/userctx/userctx_test.go` — verifies that unauthenticated middleware delegates to the injected handler.

### Modify

- `service/userctx/userctx.go:3-53` — remove layout rendering and accept an unauthorized handler.
- `router.go:11-21,65-69` — import `pages/access` and wire both access handlers.
- `components/header/menu.templ:3-15,694-710` — restore the direct `userctx` import and simplify `SetupMenuRoute`.

### Delete

- `service/userctx/unauthenticated.templ` — presentation moves to `pages/access`.
- `service/userctx/unauthenticated_test.go` — replacement test lives with the moved component.
- `service/authctx/forbidden.templ` — presentation moves to `pages/access`.
- `service/authctx/forbidden_test.go` — replacement test lives with the moved component.

### Generated artifacts

- Do not edit `service/userctx/unauthenticated_templ.go`, `service/authctx/forbidden_templ.go`, `pages/access/*_templ.go`, or `components/header/menu_templ.go` manually.
- After the source changes are complete, the user runs the normal Templ workflow. Confirm afterward that obsolete generated files from deleted templates are gone; if the generator leaves orphans, ask the user to remove those generated files through their normal cleanup workflow.

---

### Task 1: Create the access-page presentation package

**Files:**
- Create: `pages/access/access.templ`
- Create: `pages/access/unauthorized.templ`
- Create: `pages/access/forbidden.templ`
- Create: `pages/access/unauthorized_test.go`
- Create: `pages/access/forbidden_test.go`
- Delete: `service/userctx/unauthenticated.templ`
- Delete: `service/userctx/unauthenticated_test.go`
- Delete: `service/authctx/forbidden.templ`
- Delete: `service/authctx/forbidden_test.go`

**Interfaces:**
- Consumes: Templ components and the existing global button/color variables.
- Produces: `access.Unauthorized() templ.Component` and `access.Forbidden() templ.Component` for the handlers in Task 2.

- [ ] **Step 1: Add the shared access-page styles**

Create `pages/access/access.templ` with the shared styles currently duplicated by the two components. Use native CSS nesting and the repository's preferred place shorthands:

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
```

- [ ] **Step 2: Move the unauthenticated component**

Create `pages/access/unauthorized.templ`. Preserve the existing text, hrefs, CSS classes, and accessible heading relationship:

```templ
package access

templ Unauthorized() {
	@accessDeniedStyles()
	<section class="access-denied" aria-labelledby="access-denied-heading">
		<h1 id="access-denied-heading">Du må logge inn</h1>
		<p>Logg inn for å se denne siden.</p>
		<div class="access-denied-actions">
			<a href="/auth" class="btn btn--primary btn-login">
				Logg inn
			</a>
			<a href="/" class="btn btn--outline">
				Gå til arrangementslisten
			</a>
		</div>
	</section>
}
```

- [ ] **Step 3: Move the forbidden component**

Create `pages/access/forbidden.templ` with the existing administrator-access message:

```templ
package access

templ Forbidden() {
	@accessDeniedStyles()
	<section class="access-denied" aria-labelledby="access-denied-heading">
		<h1 id="access-denied-heading">Du har ikke tilgang</h1>
		<p>Du er logget inn, men denne siden krever administratortilgang.</p>
		<div class="access-denied-actions">
			<a href="/" class="btn btn--primary">
				Gå til arrangementslisten
			</a>
		</div>
	</section>
}
```

- [ ] **Step 4: Move the component tests into the new package**

Create `pages/access/unauthorized_test.go` with the moved behavior and update the component name from `Unauthenticated` to `Unauthorized`:

```go
package access

import (
	"strings"
	"testing"

	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/Regncon/conorganizer/testutil/templtest"
)

func TestUnauthorized_RendersClearLoginAndHomeLinks(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at en bruker ikke er logget inn.",
		When:  "Når innloggingsfeilsiden vises.",
		Then:  "Så skal brukeren få tydelige veier til innlogging og arrangementslisten.",
	})

	// Given
	expectedHrefs := []string{"/", "/auth"}
	expectedTextParts := []string{
		"Du må logge inn",
		"Logg inn for å se denne siden.",
		"Logg inn",
		"Gå til arrangementslisten",
	}

	// When
	doc := templtest.Render(t, Unauthorized())
	actualHrefs := templtest.CollectUniqueHrefs(doc)
	actualText := strings.Join(templtest.CollectTexts(doc, ".access-denied"), " ")

	// Then
	templtest.AssertSameHrefs(t, expectedHrefs, actualHrefs)
	for _, expectedTextPart := range expectedTextParts {
		if !strings.Contains(actualText, expectedTextPart) {
			t.Fatalf("unauthorized page text mismatch\nexpected text to contain: %q\nactual text:              %q", expectedTextPart, actualText)
		}
	}
}
```

Create `pages/access/forbidden_test.go` with the existing behavior in the new package:

```go
package access

import (
	"strings"
	"testing"

	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/Regncon/conorganizer/testutil/templtest"
)

func TestForbidden_RendersAdminAccessMessageAndHomeLink(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at en innlogget bruker mangler adminrolle.",
		When:  "Når tilgangsfeilsiden vises.",
		Then:  "Så skal brukeren få en tydelig forklaring og lenke til arrangementslisten.",
	})

	// Given
	expectedHrefs := []string{"/"}
	expectedTextParts := []string{
		"Du har ikke tilgang",
		"Du er logget inn, men denne siden krever administratortilgang.",
		"Gå til arrangementslisten",
	}

	// When
	doc := templtest.Render(t, Forbidden())
	actualHrefs := templtest.CollectUniqueHrefs(doc)
	actualText := strings.Join(templtest.CollectTexts(doc, ".access-denied"), " ")

	// Then
	templtest.AssertSameHrefs(t, expectedHrefs, actualHrefs)
	for _, expectedTextPart := range expectedTextParts {
		if !strings.Contains(actualText, expectedTextPart) {
			t.Fatalf("forbidden page text mismatch\nexpected text to contain: %q\nactual text:              %q", expectedTextPart, actualText)
		}
	}
}
```

- [ ] **Step 5: Remove the old presentation sources and tests**

Delete the four old source/test files listed for this task. Do not delete or edit ignored generated Go files; the user owns the generation/cleanup step.

- [ ] **Step 6: Review Task 1 without generating or testing**

Inspect the diff and confirm:

```text
Unauthorized component package:    pages/access
Forbidden component package:       pages/access
401 links:                          / and /auth
403 links:                          /
Old source templates:               removed
Generated *_templ.go files:         untouched
```

---

### Task 2: Add complete 401 and 403 handlers to `pages/access`

**Files:**
- Create: `pages/access/access.go`
- Create: `pages/access/access_test.go`

**Interfaces:**
- Consumes: `access.Unauthorized`, `access.Forbidden`, `layouts.Base`, `requestctx.UserRequestInfo`, and `userctx.GetUserRequestInfo`.
- Produces: `func UnauthorizedHandler(db *sql.DB, logger *slog.Logger) http.HandlerFunc` and `func ForbiddenHandler(db *sql.DB, logger *slog.Logger) http.HandlerFunc` for `router.go`.

- [ ] **Step 1: Write handler behavior tests**

Create `pages/access/access_test.go` with these two behavior-focused tests. The empty context in the 403 test intentionally makes the header render logged out and avoids database access; component-specific copy and links remain covered by the moved tests.

```go
package access

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Regncon/conorganizer/testutil/bdd"
)

func TestUnauthorizedHandler_RendersUnauthorizedResponse(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at en beskyttet side avviser en bruker som ikke er innlogget.",
		When:  "Når 401-handleren renderer svaret.",
		Then:  "Så skal svaret ha status 401 og forklare at brukeren må logge inn.",
	})

	// Given
	expectedStatusCode := http.StatusUnauthorized
	expectedText := "Du må logge inn"
	request := httptest.NewRequest(http.MethodGet, "/profile", nil)
	recorder := httptest.NewRecorder()
	handler := UnauthorizedHandler(nil, discardAccessLogger())

	// When
	handler.ServeHTTP(recorder, request)

	// Then
	if recorder.Code != expectedStatusCode {
		t.Fatalf("HTTP status mismatch\nexpected: %d\nactual:   %d", expectedStatusCode, recorder.Code)
	}
	if !strings.Contains(recorder.Body.String(), expectedText) {
		t.Fatalf("response body mismatch\nexpected body to contain: %q\nactual body:              %q", expectedText, recorder.Body.String())
	}
}

func TestForbiddenHandler_RendersForbiddenResponse(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at en adminside avviser en bruker uten administratorrolle.",
		When:  "Når 403-handleren renderer svaret.",
		Then:  "Så skal svaret ha status 403 og forklare at brukeren ikke har tilgang.",
	})

	// Given
	expectedStatusCode := http.StatusForbidden
	expectedText := "Du har ikke tilgang"
	request := httptest.NewRequest(http.MethodGet, "/admin", nil)
	recorder := httptest.NewRecorder()
	handler := ForbiddenHandler(nil, discardAccessLogger())

	// When
	handler.ServeHTTP(recorder, request)

	// Then
	if recorder.Code != expectedStatusCode {
		t.Fatalf("HTTP status mismatch\nexpected: %d\nactual:   %d", expectedStatusCode, recorder.Code)
	}
	if !strings.Contains(recorder.Body.String(), expectedText) {
		t.Fatalf("response body mismatch\nexpected body to contain: %q\nactual body:              %q", expectedText, recorder.Body.String())
	}
}

func discardAccessLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}
```

- [ ] **Step 2: Implement both complete response handlers**

Create `pages/access/access.go` with the complete implementation:

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
		requestID := middleware.GetReqID(r.Context())

		w.WriteHeader(http.StatusUnauthorized)
		if err := layouts.Base("Logg inn", requestctx.UserRequestInfo{}, db, logger, Unauthorized()).Render(r.Context(), w); err != nil {
			logger.Error(
				fmt.Errorf("failed to render unauthenticated page: %w", err).Error(),
				"request_id", requestID,
				"path", r.URL.Path,
			)
		}
	}
}

func ForbiddenHandler(db *sql.DB, logger *slog.Logger) http.HandlerFunc {
	logger = logger.With("component", "access")
	return func(w http.ResponseWriter, r *http.Request) {
		requestID := middleware.GetReqID(r.Context())
		userInfo := userctx.GetUserRequestInfo(r.Context())

		w.WriteHeader(http.StatusForbidden)
		if err := layouts.Base("Ingen tilgang", userInfo, db, logger, Forbidden()).Render(r.Context(), w); err != nil {
			logger.Error(
				fmt.Errorf("failed to render forbidden page: %w", err).Error(),
				"request_id", requestID,
				"path", r.URL.Path,
			)
		}
	}
}
```

This implementation preserves the status codes, titles, and layout inputs while changing ownership from the service packages to `pages/access`.

- [ ] **Step 3: Check the handler contracts against router composition**

Confirm that both public functions return the exact type accepted by their consumers:

```go
func UnauthorizedHandler(db *sql.DB, logger *slog.Logger) http.HandlerFunc
func ForbiddenHandler(db *sql.DB, logger *slog.Logger) http.HandlerFunc
```

- [ ] **Step 4: Review Task 2 without running tests**

Confirm the tests and implementation agree on these interfaces and statuses:

```text
UnauthorizedHandler(*sql.DB, *slog.Logger) http.HandlerFunc → 401
ForbiddenHandler(*sql.DB, *slog.Logger) http.HandlerFunc    → 403
```

Do not run the tests because generated access components do not exist until the user performs Templ generation.

---

### Task 3: Make `UserMiddleware` presentation-independent

**Files:**
- Modify: `service/userctx/userctx.go:3-53`
- Create: `service/userctx/userctx_test.go`

**Interfaces:**
- Consumes: an `http.HandlerFunc` supplied by the application router.
- Produces: `func UserMiddleware(logger *slog.Logger, unauthorizedHandler http.HandlerFunc) func(http.Handler) http.Handler` without any `layouts` dependency.

- [ ] **Step 1: Add the unauthenticated delegation test**

Create `service/userctx/userctx_test.go` with this behavior-focused middleware test:

```go
package userctx

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Regncon/conorganizer/testutil/bdd"
)

func TestUserMiddleware_WhenUserIsNotLoggedIn_UsesUnauthorizedHandler(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at en bruker ikke er innlogget og en egen 401-handler er satt.",
		When:  "Når brukeren åpner en beskyttet rute.",
		Then:  "Så skal 401-handleren behandle svaret uten at den beskyttede handleren kjøres.",
	})

	// Given
	expectedStatusCode := http.StatusUnauthorized
	expectedBody := "custom unauthorized"
	expectedUnauthorizedHandlerCalled := true
	expectedProtectedHandlerCalled := false

	unauthorizedHandlerCalled := false
	protectedHandlerCalled := false
	unauthorizedHandler := func(w http.ResponseWriter, _ *http.Request) {
		unauthorizedHandlerCalled = true
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(expectedBody))
	}
	protectedHandler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		protectedHandlerCalled = true
		w.WriteHeader(http.StatusOK)
	})
	handler := UserMiddleware(discardUserctxLogger(), unauthorizedHandler)(protectedHandler)
	request := httptest.NewRequest(http.MethodGet, "/profile", nil)
	recorder := httptest.NewRecorder()

	// When
	handler.ServeHTTP(recorder, request)

	// Then
	if recorder.Code != expectedStatusCode {
		t.Fatalf("HTTP status mismatch\nexpected: %d\nactual:   %d", expectedStatusCode, recorder.Code)
	}
	if recorder.Body.String() != expectedBody {
		t.Fatalf("HTTP body mismatch\nexpected: %q\nactual:   %q", expectedBody, recorder.Body.String())
	}
	if unauthorizedHandlerCalled != expectedUnauthorizedHandlerCalled {
		t.Fatalf("unauthorized handler call mismatch\nexpected: %v\nactual:   %v", expectedUnauthorizedHandlerCalled, unauthorizedHandlerCalled)
	}
	if protectedHandlerCalled != expectedProtectedHandlerCalled {
		t.Fatalf("protected handler call mismatch\nexpected: %v\nactual:   %v", expectedProtectedHandlerCalled, protectedHandlerCalled)
	}
}

func discardUserctxLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}
```

- [ ] **Step 2: Change the middleware signature**

Replace:

```go
func UserMiddleware(logger *slog.Logger, db *sql.DB) func(http.Handler) http.Handler
```

with:

```go
func UserMiddleware(logger *slog.Logger, unauthorizedHandler http.HandlerFunc) func(http.Handler) http.Handler
```

Keep the existing component-enriched logger and authenticated path. Simplify the unauthenticated branch so it logs the existing request ID and path, invokes `unauthorizedHandler(w, r)`, and returns:

```go
if userInfo.IsLoggedIn {
	logger.Debug("User is logged in", "request_id", requestID)
	next.ServeHTTP(w, r.WithContext(ctx))
	return
}

logger.Debug("User is not logged in", "request_id", requestID, "path", r.URL.Path)
unauthorizedHandler(w, r)
```

Do not set the status in `UserMiddleware`; the injected handler owns the complete response.

- [ ] **Step 3: Remove presentation dependencies from `userctx.go`**

Remove:

- the `layouts` import;
- `AdminForbiddenHandler`;
- the `w.WriteHeader` and `layouts.Base` rendering previously inside `UserMiddleware`.

Keep `database/sql`, `fmt`, `authctx`, `requestctx`, and Chi middleware because the remaining context and user-ID helpers still require them.

- [ ] **Step 4: Review the package boundary without running tests**

Run only read-only searches or inspect the diff. The expected source relationship is:

```text
service/userctx imports layouts:                         no
service/userctx contains AdminForbiddenHandler:          no
service/userctx contains Unauthenticated Templ source:   no
UserMiddleware accepts unauthorizedHandler:              yes
```

---

### Task 4: Wire access handlers and restore the direct menu dependency

**Files:**
- Modify: `router.go:11-21,65-69`
- Modify: `components/header/menu.templ:3-15,694-710`

**Interfaces:**
- Consumes: `access.UnauthorizedHandler`, `access.ForbiddenHandler`, the new `userctx.UserMiddleware` signature, and `userctx.GetUserRequestInfo`.
- Produces: final route composition with no injected user-info function in `SetupMenuRoute`.

- [ ] **Step 1: Wire the access package in `router.go`**

Add:

```go
"github.com/Regncon/conorganizer/pages/access"
```

Replace the current middleware/menu/admin block with:

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

This removes both `userctx.GetUserRequestInfo` function injection and `userctx.AdminForbiddenHandler` from router composition.

- [ ] **Step 2: Restore the `userctx` import in `menu.templ`**

Add:

```go
"github.com/Regncon/conorganizer/service/userctx"
```

Keep the `requestctx` import because the menu component signatures still use `requestctx.UserRequestInfo`.

- [ ] **Step 3: Simplify `SetupMenuRoute`**

Change the function signature to:

```go
func SetupMenuRoute(
	router chi.Router,
	liveManager *live.Manager,
	db *sql.DB,
	logger *slog.Logger,
)
```

Change the live render callback to:

```go
Render: func(ctx context.Context, r *http.Request) templ.Component {
	userInfo := userctx.GetUserRequestInfo(ctx)
	return MenuBillettholderLive(userInfo, db, logger)
},
```

Do not change the `/menu/api` path, `live.BucketBillettholders`, or `MenuBillettholderLive` arguments.

- [ ] **Step 4: Inspect the final dependency direction**

Confirm from source imports:

```text
pages/access → layouts
pages/access → userctx
layouts → header
header → userctx
userctx → authctx
userctx ⇏ layouts
```

The direct `header → userctx` dependency is now valid because there is no path from `userctx` back to `header`.

---

### Task 5: User-owned generation and verification handoff

**Files:**
- Generated by user: `pages/access/*_templ.go`
- Regenerated by user: `components/header/menu_templ.go`
- Removed through user's normal generated-file cleanup if orphaned: old access component `*_templ.go` files.

**Interfaces:**
- Consumes: all approved source changes from Tasks 1–4.
- Produces: generated Go matching the moved Templ components and the restored menu route signature.

- [ ] **Step 1: Stop and ask the user to run their normal Templ generation workflow**

Do not run a generation command. Tell the user which source templates moved and that the generated menu file must reflect the simplified `SetupMenuRoute` signature.

- [ ] **Step 2: Inspect generated results after the user confirms generation**

Confirm without editing generated files:

```text
pages/access generated components exist
components/header/menu_templ.go imports userctx
generated SetupMenuRoute has four parameters
obsolete generated Unauthenticated/Forbidden components are absent from service packages
```

- [ ] **Step 3: Perform allowed static checks only**

The implementation agent may run:

```powershell
git diff --check
rg -n "conorganizer/layouts|AdminForbiddenHandler|func Unauthenticated|func Forbidden" service/userctx service/authctx
rg -n "func SetupMenuRoute|GetUserRequestInfo" components/header/menu.templ router.go
```

Expected results:

- `git diff --check` exits successfully;
- no `layouts`, access component, or forbidden-handler ownership remains in `service/userctx` or `service/authctx` source;
- `menu.templ` directly calls `userctx.GetUserRequestInfo`;
- `router.go` calls `SetupMenuRoute` with four arguments.

- [ ] **Step 4: Give the user the test commands without executing them**

Provide these commands for user-owned verification:

```powershell
go test ./pages/access ./service/userctx ./service/authctx ./components/header
go test ./...
```

Expected result: all packages compile without an import cycle and all tests pass.

- [ ] **Step 5: Report the completed boundary change**

The handoff must state:

- 401 and 403 page rendering now lives in `pages/access`;
- `userctx` no longer imports `layouts`;
- `UserMiddleware` delegates 401 responses to the injected handler;
- `SetupMenuRoute` directly obtains request user information again;
- Templ generation and tests were not run by Codex;
- no commit was created.

---

## Self-Review Record

- **Spec coverage:** Every acceptance criterion maps to Tasks 1–5.
- **Placeholder scan:** The plan contains no deferred implementation choices; public signatures, ownership, statuses, copy, and wiring are specified.
- **Type consistency:** `UnauthorizedHandler` and `ForbiddenHandler` both return `http.HandlerFunc`; this matches `UserMiddleware` and `authctx.WithForbiddenHandler`. `SetupMenuRoute` has the same four-argument signature in `menu.templ` and `router.go`.
- **Scope check:** The plan changes one package boundary spanning access presentation, middleware composition, and the menu workaround. Moving `UserRequestInfo`, billettholder behavior, and broader error-page consolidation remain outside scope.
