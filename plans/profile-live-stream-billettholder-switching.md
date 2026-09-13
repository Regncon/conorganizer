# Profile Cookie Selection and Live Billettholder Switching Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Status:** Deferred. Implement after the current branch is finished. Saving this plan does not start implementation.

**Goal:** Open `/profile` with the saved billettholder already rendered, then update `#profile-main-column` when `$billettHolderId` changes without reloading the document or using `b_id` in the URL.

**Architecture:** The existing `selectedBillettholderId` cookie supplies the initial server-rendered selection through `requestctx.SelectedBillettholderID(ctx)`. A profile-owned Datastar effect observes the menu's `$billettHolderId`, replaces the previous `GET /profile/api` stream, and lets that stream patch the main column. Each stream captures the signal from its request and validates it against the user's current billettholdere on every render.

**Tech Stack:** Go, templ, the bundled `static/datastar.js`, `datastar-go v1.2.2`, existing `service/live.Manager`, and existing NATS buckets.

**Spec:** The Requirements and Selection Rules sections in this document. Updated from the earlier switching plan on 2026-09-11 to include the cookie middleware and current menu wiring.

## Global Constraints

- This is a plan for later implementation; change no application code or run application tests while saving it.
- At execution time, reread the current branch's instructions and relevant files. The current branch may change before this plan is started.
- Use the names in `domeneordbok.md`, including billettholder, billettholdere, interesse, and pulje.
- Keep `window.conorganizer.billettholderSelection` as the existing browser selection store. Add no cookie, localStorage key, session selection store, or NATS subject.
- JavaScript already writes the cookie; the middleware reads it. Do not add a request merely to update the cookie or refresh DevTools.
- Keep the global menu independent of profile URLs and requests. The profile owns its stream.
- Validate every cookie and signal ID against the authenticated user's associated billettholdere. Neither the cookie nor the signal grants access.
- Billettholder selection scopes only Mitt festivalprogram. The profile, account controls, and Mine arrangementer remain available based on the authenticated user, including when no billettholder exists.
- Patch `#profile-main-column` only. Preserve the menu/dialog, account controls, ticket summary, document, and scroll position during selection changes.
- Keep the current buckets: `live.BucketEvents`, `live.BucketInterests`, and `live.BucketBillettholders`.
- Use the fixed action URI `/profile/api`. Do not put the selected ID into the action URI; Datastar sends the public signal automatically.
- Remove all profile selection reads and writes of `b_id`. Do not add URL synchronization, redirects, or history API calls. Existing links containing `b_id` still open the profile, but that parameter is ignored and need not be stripped from the address bar.
- Selection changes leave the visible URL, query parameters, fragment, and history state untouched. The `datastar` query parameter on the API request remains the SDK's signal transport; it is not a visible profile URL selection parameter.
- Use the existing templ watcher when it is running. Do not hand-edit generated `*_templ.go` files or run competing generators.
- During later implementation, follow AGENTS.md's behavior-focused Go test structure. Keep fixture setup in test helpers.
- No shared live-service refactor, dependency upgrade, or CSS change is needed.

## Requirements

1. Opening `/profile` with a valid selection cookie renders that billettholder's program in the first HTML response.
2. Neither opening the profile nor switching billettholder triggers the current `window.location.replace` navigation.
3. Switching billettholder updates the main column using the newly selected ID. A changing cookie alone does not change an already-open stream.
4. At most one profile stream remains active after initialization or a switch; cancelled streams cannot later restore an earlier selection.
5. Existing NATS notifications continue updating the active stream. Revalidate its captured ID on each render in case the user's relation is removed.
6. Missing, malformed, stale, or unrelated IDs never expose another user's profile data.
7. A user with no billettholdere still receives a working profile and live updates for account-owned content; only Mitt festivalprogram shows its empty state. A missing cookie or missing/zero signal is not an error.
8. A transient connection failure reconnects the current selection. Clean server-side stream closure also reconnects.
9. Menu selection and displayed program agree after normal initialization and switches, with no `b_id` dependency or URL mutation. Cookies remain browser-wide, while an open stream follows its own request's signal.
10. Switching on `/` and `/event/{id}` continues working after the small menu fallback change.

## Selection Rules

| Situation | Selection order |
| --- | --- |
| Initial `GET /profile` | Associated cookie ID; otherwise existing email match / first billettholder / zero fallback |
| An old link includes `b_id` | Ignore it, whether or not a selection cookie exists |
| New `GET /profile/api` | Associated `billettHolderId` from Datastar signals; otherwise existing email match / first billettholder / zero fallback |
| Malformed JSON, or a string, boolean, object, array, or fractional signal value | HTTP 400 before opening an SSE stream |
| Later NATS update | Recheck the captured signal ID against freshly loaded associations; use the normal fallback if it is no longer associated |

The cookie is a hint for the first render. The live request's signal is authoritative for selecting among authorized billettholdere; do not let a cookie override that signal. Neither handler reads `b_id`. This also avoids depending on whether a request starts before or after the browser writes the cookie during a click.

The menu already validates localStorage during initialization. If cookies are unavailable, or localStorage differs from the cookie, the first live patch reconciles the profile to the menu without navigation. No query parameter selects the billettholder for the initial page render.

For example, with associated IDs 101 and 202:

- Cookie 202 and URL `?b_id=101`: first render uses 202.
- No cookie and URL `?b_id=202`: ignore the query and use the existing authorized default.
- Live signal 101 and cookie 202: that stream renders 101.
- Live signal 999: render the authorized default, never 999.
- No associated billettholdere: selection is zero.

## Files and Responsibilities

| File | Planned change |
| --- | --- |
| `pages/profile/profile.go` | Cookie/default initial selection; remove `b_id` parsing; decode live signals once; validate captured ID on each render; remove client-only valid-ID plumbing |
| `pages/profile/profile_index.templ` | Replace reload script and initial stream owner with one reactive profile stream; remove all selection URL handling |
| `components/header/menu.templ` | Use numeric zero in the menu's three empty-selection assignments |
| `pages/profile/profile_selection_test.go` | Cookie/default selection, ignored legacy query, and live selection authorization regressions |
| `pages/profile/profile_page_test.go` | Update the existing rendering test for the changed template contract |
| `documentation/testing/profile.md` | Add the browser acceptance checklist |
| `service/requestctx/billettholder_selection.go` | Read only: existing middleware and accessor |
| `static/js/conorganizer.js` | Read only: storage, cookie writes, and selection event |
| `service/live/live.go` | Read only: stream cancellation, watcher cleanup, and current request options |
| `static/datastar.js` | Read only: verify the shipped action and retry behavior |

## Task 1: Render the Initial Profile From the Cookie

**Files:** `pages/profile/profile.go`, `pages/profile/profile_selection_test.go`.

**Interfaces:**
- Consumes `requestctx.SelectedBillettholderID(r.Context()) int`.
- Preserves `selectedBillettholderIDFromRequest(r *http.Request, user requestctx.UserRequestInfo, billettholdere []models.Billettholder, logger *slog.Logger) int`.
- Preserves only the existing email match / first billettholder / zero fallback after cookie validation; removes the query fallback.

- [ ] **Step 1: Add a regression test exercising the real middleware and initial selector.**

Replace the three existing query-specific tests in `profile_selection_test.go` with this test matrix, reusing its current imports and fixture helpers. Retain the existing email-match, first-billettholder, and no-billettholdere tests.

```go
func TestSelectedBillettholderIDFromRequest_UsesCookieOrDefaultAndIgnoresQuery(t *testing.T) {
    cases := []struct {
        name string
        cookie string
        target string
        expectedID int
    }{
        {"cookie without query", "202", "/profile", 202},
        {"old query ignored with cookie", "202", "/profile?b_id=101", 202},
        {"invalid query ignored with cookie", "202", "/profile?b_id=invalid", 202},
        {"missing cookie ignores old query", "", "/profile?b_id=202", 101},
        {"unrelated cookie ignores old query", "999", "/profile?b_id=202", 101},
        {"invalid cookie ignores old query", "invalid", "/profile?b_id=202", 101},
        {"negative cookie uses default", "-1", "/profile", 101},
        {"unrelated cookie uses default", "999", "/profile", 101},
    }
    for _, tc := range cases {
        t.Run(tc.name, func(t *testing.T) {
            bdd.Behavior(t, bdd.BDD{
                Given: "An authenticated user with two associated billettholdere and an optional selection cookie.",
                When: "The profile resolves its initial selection.",
                Then: "The cookie or authorized default selects the billettholder; the URL never does.",
            })

            // Given
            expectedID := tc.expectedID
            user := profileSelectionUser("owner@example.com")
            billettholdere := []models.Billettholder{
                profileSelectionBillettholder(101, user.Email),
                profileSelectionBillettholder(202, "other@example.com"),
            }
            request := profileSelectionRequest(t, tc.target)
            if tc.cookie != "" {
                request.AddCookie(&http.Cookie{
                    Name: requestctx.SelectedBillettholderCookieName,
                    Value: tc.cookie,
                })
            }
            var actualID int
            handler := requestctx.BillettholderSelectionMiddleware(
                http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                    actualID = selectedBillettholderIDFromRequest(
                        r, user, billettholdere, testutil.NewTestLogger(),
                    )
                }),
            )

            // When
            handler.ServeHTTP(httptest.NewRecorder(), request)

            // Then
            if actualID != expectedID {
                t.Fatalf("expected billettholder %d; got %d", expectedID, actualID)
            }
        })
    }
}
```

- [ ] **Step 2: Run the focused test and confirm that cookie selection and ignored-query cases fail for the expected reasons.**

```powershell
go test ./pages/profile -run '^TestSelectedBillettholderIDFromRequest_' -count=1
```

- [ ] **Step 3: Replace `selectedBillettholderIDFromRequest` with cookie/default selection only.**

```go
func selectedBillettholderIDFromRequest(r *http.Request, user requestctx.UserRequestInfo, billettholdere []models.Billettholder, logger *slog.Logger) int {
    cookieID := requestctx.SelectedBillettholderID(r.Context())
    if cookieID > 0 && hasBillettholderID(billettholdere, cookieID) {
        return cookieID
    }
    return defaultSelectedBillettholderID(user, billettholdere, logger)
}
```

Delete the previous `b_id` parsing and query-specific logging. Remove the now-unused `strconv` import from `profile.go`; keep `strings`, which is still used elsewhere. Keep `defaultSelectedBillettholderID` unchanged.

- [ ] **Step 4: Rerun the same focused test command.** Existing default-selection tests and the new cookie/ignored-query tests must pass.

## Task 2: Capture and Validate the Live Selection

**Files:** `components/header/menu.templ`, `pages/profile/profile.go`, `pages/profile/profile_selection_test.go`.

**Interfaces:**
- Consumes the existing public `$billettHolderId` signal.
- Produces a captured integer per `/profile/api` request.
- Adds `selectedBillettholderIDFromSignal(id int, user requestctx.UserRequestInfo, billettholdere []models.Billettholder, logger *slog.Logger) int`.

- [ ] **Step 1: Add the signal authorization test.**

```go
func TestSelectedBillettholderIDFromSignal_OnlySelectsAssociatedBillettholdere(t *testing.T) {
    cases := []struct {
        name string
        signalID int
        expectedID int
        noBillettholdere bool
    }{
        {"associated selection", 202, 202, false},
        {"unrelated selection", 999, 101, false},
        {"cleared selection", 0, 101, false},
        {"negative selection", -1, 101, false},
        {"removed association", 202, 0, true},
        {"no billettholdere", 0, 0, true},
    }
    for _, tc := range cases {
        t.Run(tc.name, func(t *testing.T) {
            bdd.Behavior(t, bdd.BDD{
                Given: "A live profile selection and the user's current associated billettholdere.",
                When: "The profile resolves the live selection.",
                Then: "Only an associated billettholder or the authorized fallback is rendered.",
            })

            // Given
            expectedID := tc.expectedID
            user := profileSelectionUser("owner@example.com")
            billettholdere := []models.Billettholder{
                profileSelectionBillettholder(101, user.Email),
                profileSelectionBillettholder(202, "other@example.com"),
            }
            if tc.noBillettholdere {
                billettholdere = nil
            }

            // When
            actualID := selectedBillettholderIDFromSignal(
                tc.signalID, user, billettholdere, testutil.NewTestLogger(),
            )

            // Then
            if actualID != expectedID {
                t.Fatalf("expected billettholder %d; got %d", expectedID, actualID)
            }
        })
    }
}
```

- [ ] **Step 2: Run the new test and confirm the missing helper fails compilation.**

```powershell
go test ./pages/profile -run '^TestSelectedBillettholderIDFromSignal_' -count=1
```

- [ ] **Step 3: Add the live selection helper to `profile.go`.**

```go
func selectedBillettholderIDFromSignal(id int, user requestctx.UserRequestInfo, billettholdere []models.Billettholder, logger *slog.Logger) int {
    if id > 0 && hasBillettholderID(billettholdere, id) {
        return id
    }
    return defaultSelectedBillettholderID(user, billettholdere, logger)
}
```

- [ ] **Step 4: Change the menu's three empty fallbacks to numeric zero.**

Use the current window event, not the old plan's removed `menu-billettholder-change` handler:

```templ
data-signals:billett-holder-id="0"
data-on:billettholder-selection-change__window="$_menuBillettholder = evt.detail; $billettHolderId = evt.detail?.Id ?? 0; $_menuBillettholderReady = true"
```

In `menuBillettholderSelectionEffect`, change only the assignment to:

```javascript
$billettHolderId = selected?.Id ?? 0;
```

Keep the existing cookie, pending/skeleton behavior, initialization, and selection event. These menu assignments supply an integer on the profile, where there is no event interest picker. Do not expand this into an unrelated rewrite of every selection component.

- [ ] **Step 5: Read signals once, before `liveManager.Stream` in the profile API GET handler.**

Add the same Go SDK import already used by `pages/event/event.go`:

```go
datastar "github.com/starfederation/datastar-go/datastar"
```

Insert:

```go
signals := struct {
    BillettHolderID int `json:"billettHolderId"`
}{}
if err := datastar.ReadSignals(r, &signals); err != nil {
    http.Error(w, "Ugyldig profilvalg.", http.StatusBadRequest)
    return
}
requestedBillettholderID := signals.BillettHolderID
```

Inside the existing `Render` callback, after fetching current billettholdere, replace the call to `selectedBillettholderIDFromRequest` with:

```go
selectedBillettholderID := selectedBillettholderIDFromSignal(
    requestedBillettholderID, user, billettholdere, requestLogger,
)
```

Keep the existing `ProfileMainColumn(...)` return and buckets. Only the raw requested ID is captured; association checks remain inside `Render`. Missing or null signal values decode as zero and use the authorized default. Malformed JSON and incompatible JSON types receive 400.

- [ ] **Step 6: Rerun the selection tests and existing menu tests.**

After the existing watcher regenerates changed templates:

```powershell
go test ./pages/profile ./components/header -count=1
```

## Task 3: Replace Navigation With One Profile Stream Effect

**Files:** `pages/profile/profile_index.templ`, `pages/profile/profile.go`, `pages/profile/profile_page_test.go`.

**Interfaces:**
- Consumes `$billettHolderId` and fixed action URI `GET /profile/api`.
- Keeps `ProfileMainColumn` and its root ID `profile-main-column`.
- Changes `ProfilePage` to remove only the `validBillettholderIDs []int` parameter.

- [ ] **Step 1: Replace the profile container's selection metadata and `data-init` with this effect.**

```templ
data-effect="
    $billettHolderId;
    @get('/profile/api', {
        requestCancellation: 'auto',
        retry: 'always',
        retryMaxCount: Infinity,
        retryInterval: 1000,
        retryMaxWait: 30000,
    })
"
```

The standalone `$billettHolderId;` read makes the effect depend on that signal; `@get` sends its value automatically. Keep that read even though it does not assign a variable. This element stays outside `#profile-main-column`, so incoming patches do not recreate the effect. Opening the menu, closing the dialog, and unrelated signal changes must not replace this stream. There is no URL or history synchronization.

Use `retry: 'always'` deliberately: the shared manager can close a stream cleanly when a watcher closes, and the bundled client's default `auto` retries network errors but not a clean EOF. Keep the option spelled `retryMaxWait`; the shared helper currently uses `retryMaxWaitMs`, which is not the shipped action option. Do not modify that shared helper in this task.

- [ ] **Step 2: Delete the old inline reload script and unused plumbing.**

Remove:

- The entire profile selection `<script>` with `selection.onChange`, `selection.get`, and `window.location.replace`.
- Both `data-profile-*` attributes used by that script.
- The now-unused `fmt` and `service/live` imports from `profile_index.templ`.
- `validBillettholderIDs := billettholderIDs(billettholdere)` and that argument from the GET handler's `ProfilePage` call.
- The `billettholderIDs` helper from `profile.go` and its template parameter.

Retain `selectedBillettholderID int` for the initial server render. The resulting call is:

```go
ProfilePage(user, events, tickets, selectedBillettholderID, db, requestLogger, eventImageDir)
```

- [ ] **Step 3: Update the existing profile rendering test to the new template contract.**

Replace `TestProfilePage_RendersBreadcrumbAndBillettholderSelectionMetadata` with the following test. Remove the obsolete `slices` import; `strings`, `testing`, and the existing project imports remain.

```go
func TestProfilePage_RendersOverviewWithLiveMainColumn(t *testing.T) {
    bdd.Behavior(t, bdd.BDD{
        Given: "An authenticated user opening the profile.",
        When: "The profile overview renders.",
        Then: "The overview has a stable main column and an initialized reactive live update.",
    })

    // Given
    expectedMainColumns := 1
    db, logger := testutil.CreateTestDBAndLogger(t, "profile_page")
    user := requestctx.UserRequestInfo{
        IsLoggedIn: true,
        Id: "profile-page-user",
        Email: "profile-page-user@example.com",
    }

    // When
    doc := templtest.Render(t, ProfilePage(user, nil, nil, 22, db, logger, nil))

    // Then
    if actual := doc.Find("#profile-main-column").Length(); actual != expectedMainColumns {
        t.Fatalf("expected %d main column; got %d", expectedMainColumns, actual)
    }
    if doc.Find(".breadcrumb-end").Text() != "Min Side" {
        t.Fatal("expected Min Side breadcrumb")
    }
    effect, exists := doc.Find(".profile-container").Attr("data-effect")
    if !exists || !strings.Contains(effect, "@get('/profile/api'") {
        t.Fatal("expected the profile container to own its reactive live request")
    }
}
```

This maintains the existing render-level coverage; the browser checks below verify actual reactivity and cancellation.

- [ ] **Step 4: Let the templ watcher regenerate, then run the focused suite and source checks.**

```powershell
go test ./pages/profile ./components/header ./service/requestctx ./service/live -count=1
git diff --check
rg -n 'b_id|window.location|history\.|searchParams|selection.onChange|DatastarInitExpression|validBillettholderIDs|billettholderIDs' pages/profile/profile.go pages/profile/profile_index.templ
```

Expected: tests pass, diff check is clean, and the final search has no matches (rg exit code 1 means no matches).

## Task 4: Verify the Full Browser Flow and Record It

**Files:** `documentation/testing/profile.md`; verify the implementation files without broadening scope.

**Interfaces:**
- Consumes a running app and an authenticated user with two associated billettholdere whose programs differ.
- Produces evidence for initial HTML selection, live updates, cancellation, retries, and access checks.

- [ ] **Step 1: Check first render from the cookie and ignored legacy URLs.**

Select billettholder B on `/`, then navigate normally to `/profile`. Inspect the initial document response, not only the patched DOM. It must already show B's program. There must be one document navigation for entering the page and no selection-driven second document request.

Repeat with a URL containing A's old `b_id`, another query parameter, and a fragment. The B cookie determines the first render and the address stays untouched. Repeat the initial request without a cookie: the old `b_id` must not override the existing authorized default. A later patch may reconcile to the menu's saved localStorage selection.

- [ ] **Step 2: Check switching and unrelated interactions.**

With Network filtered to `/profile/api`, switch from B to A and back:

- Each new request carries the selected numeric `billettHolderId` in its `datastar` query parameter.
- The preceding request is cancelled; one profile request remains active.
- The main column shows the final selected billettholder's program.
- Menu/dialog, account controls, ticket summary, and scroll position remain stable.
- Opening/closing the menu without changing the ID creates no replacement profile request.
- The entire visible URL and history state stay unchanged. Opening `/profile` and switching must never add `b_id`; an old link's existing `b_id` stays inert.
- Back does not cycle through billettholder choices.

- [ ] **Step 3: Check rapid switching and NATS updates.**

Throttle the browser connection and switch A/B several times quickly. After settling on B, only B may remain displayed. Trigger an ordinary interesse or event update and confirm B's stream receives it. No response from a cancelled A request may overwrite B.

Inspect the request-context cancellation and deferred watcher cleanup in `service/live/live.go` if a superseded stream remains active. Do not create additional selection state or NATS buckets.

- [ ] **Step 4: Check reconnects and tab isolation.**

Interrupt connectivity briefly, restore it, and confirm the current selection reconnects. Also restart the development server or close the active watcher connection to exercise clean EOF handling.

Open another profile tab and select a different billettholder there. An existing stream must still follow its own captured signal on subsequent updates, even though the shared browser cookie has changed.

- [ ] **Step 5: Check invalid and missing selections.**

Exercise each of these cases using test accounts/data:

- Missing cookie with valid localStorage: the live patch reconciles to the menu without document navigation.
- Cookie containing malformed, negative, zero, or unrelated ID: authorized initial fallback.
- A numeric but unrelated live signal: authorized live fallback, never another user's program.
- Malformed `datastar` JSON or a string/object instead of integer `billettHolderId`: HTTP 400 with no SSE stream.
- No associated billettholdere: numeric zero on the profile, an empty Mitt festivalprogram, and a working profile/main-column stream. No cookie or ID is required to use the profile.
- Remove the selected association while the profile is open and broadcast the change: the next render must not expose that billettholder's data.

Use the Cookies table refresh or a read of `document.cookie` to inspect actual stored values. A stale DevTools display is not evidence that the cookie write failed.

- [ ] **Step 6: Check adjacent pages.**

Switch billettholder on `/` and `/event/{id}`. Verify the menu, interest picker, selected interests, and cookie still agree. The home page should not gain any profile or cookie-update request.

- [ ] **Step 7: Update the manual checklist and record verification.**

Add entries under `documentation/testing/profile.md` for:

- Initial render from the selection cookie.
- Legacy `b_id` parameters have no effect, and selection changes do not modify the URL.
- Live billettholder switching without navigation.
- Rapid switching with only the final stream active.
- Reconnect behavior and missing/invalid selections.

Use its existing Given/When/Then format and Bokmål for user-facing text. Record the commands run, browser scenarios checked, and any remaining limitation in the implementation handoff. Commit the verified implementation on its later branch using the normal repository workflow.

## Execution Handoff

When this branch is finished, start the implementation from the completed branch state, using a fresh `codex/` branch or worktree as appropriate. First confirm that the cookie middleware, current menu wiring, profile handlers, and bundled Datastar options still match this plan. Do not apply the earlier plan's obsolete menu event names.

Suggested later task prompt:

> Implement plans/profile-live-stream-billettholder-switching.md. Use the existing cookie for the profile's initial selection and the menu's billettHolderId signal for live updates. Remove b_id parsing and URL synchronization. Keep the profile working without a billettholder; the selection scopes Mitt festivalprogram only. Complete the focused tests and browser acceptance checks, including cancellation of superseded streams.

## References

The shipped `static/datastar.js` is the compatibility reference for this repository. The public [Datastar actions reference](https://data-star.dev/reference/actions#request-cancellation) documents automatic same-method/same-URI cancellation, signal transport, and retry options. Verify these against the bundled implementation again when executing the plan.

The existing `service/live/live.go` captures its request context, patches the component returned by `Page.Render`, and stops watchers when the request is cancelled. Reuse that lifecycle.
