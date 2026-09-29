# Automated tests

This page describes how the Go test suite is organized and the conventions we follow when writing tests. The step-by-step test structure (BDD metadata, `// Given` / `// When` / `// Then`, production helpers vs. test helpers) is defined in [AGENTS.md](../AGENTS.md#go-test-structure). Manual test checklists live in [testing/](testing/index.md).

## Running tests and the behavior report

Tests create their SQLite databases from `schema.sql` in the project root, so that file must exist before `go test` works. `go tool task test` regenerates it from `database/events.db` and then runs `go test ./...`. After that, plain `go test ./...` is enough. See [Run tests in the README](../README.md#run-tests).

`cmd/testreport` produces a Markdown "Automated Behavior Test Report":

```bash
go tool task test:report   # same as: go run ./cmd/testreport
```

It runs `go list -json ./...` to find test files, reads the BDD metadata at the top of each `Test...` function, then runs `go test -json ./...`. The report has a summary (packages with tests, tests run, failed tests, failed packages, skipped, tests missing BDD metadata) and, per package, each test's PASS/FAIL/SKIP status followed by its Given/When/Then text. Tests without metadata show `BDD-metadata mangler.` The command exits non-zero when `go test` fails.

CI runs the same command in the "Run tests with behavior report" step of `.github/workflows/buildAndTest.yml`. The report is only printed to the CI log; it is not saved as an artifact.

## BDD metadata

Declare a test's behavior with the `testutil/bdd` package:

```go
bdd.Behavior(t, bdd.BDD{
	Given: "Gitt at brukeren ikke er innlogget.",
	When:  "Når hovednavigasjonen vises.",
	Then:  "Så skal brukeren bare få navigasjonslenker til forsiden og innlogging.",
})
```

`bdd.Behavior` fails the test if any of the three fields is blank. The text may be Norwegian or English.

`testutil/bdd` is a leaf package on purpose: the root `testutil` package imports `service` (to create test databases), so tests in `service` packages cannot import `testutil` without an import cycle. They can always import `testutil/bdd`.

How `cmd/testreport` finds the metadata (it parses the Go AST, it does not run the call):

- It looks for a call named `Behavior` (`bdd.Behavior(...)` or `Behavior(...)`) whose argument is a composite literal of type `BDD` (`bdd.BDD{...}` or `BDD{...}`). Other struct literals with Given/When/Then fields are ignored.
- Only string-literal `Given`, `When` and `Then` values are read, and all three must be non-empty. Constants or concatenated strings are not picked up.
- Fallback for tests that have not been migrated: the first comment group at the top of the function body (before the first statement) is used if it contains given/gitt, when/når and then/så. The bare section markers `// Given`, `// When` and `// Then` are skipped.

## Test infrastructure

Generic infrastructure lives in `testutil` (`testutil/createTmpDbLogger.go`):

- `CreateTestDB(t, name)` creates a SQLite database in `t.TempDir()` from `schema.sql` via `service.InitTestDBFrom` and closes it with `t.Cleanup`.
- `CreateTestDBAndLogger(t, name)` also returns a stub `*slog.Logger`.
- `NewTestLogger()`, `MustExec(t, db, query, args...)` and `QueryInt(t, db, query, args...)`.

`testutil` does not seed lookup rows or domain data. Most tests use it; `service/puljefordeling/commit_consistency_test.go` still calls `service.InitTestDBFrom` directly.

Behavior-specific fixture helpers stay package-local and are named after what they set up, for example `insertRootPageEvent` and `insertRootPageEventPulje` in `pages/root`, or `seedEventInterestUpdateFixture` in `pages/event`. Do not centralize them into one shared fixture factory: a mega-factory saves lines but makes tests harder to read. Only generic infrastructure belongs in `testutil`, and component render/query helpers belong in `testutil/templtest`.

## Component tests with templ and goquery

templ components implement `Render(ctx, io.Writer)`, so conditional and role-based rendering in `.templ` files can be tested in Go without a browser, as recommended in the [templ testing guide](https://templ.guide/core-concepts/testing/).

Use `testutil/templtest` instead of writing render or parse code in each test:

| Helper | What it does |
| --- | --- |
| `Render(t, component)` | Renders a `templ.Component` and parses it into a `*goquery.Document`. |
| `HasSelector(doc, selector)` | Reports whether the selector matches anything. |
| `CollectTexts(doc, selector)` | Returns the text of each match, with whitespace normalized. |
| `CollectUniqueHrefs(doc)` | Returns the sorted, unique, non-empty `href` values of all `a[href]`. |
| `AssertSameHrefs(t, expected, actual)` | Order-independent exact comparison of two href sets. |

`templtest` imports only templ and goquery. It is deliberately separate from the root `testutil` package, which imports `service`. Component tests that also need a database (for example `components/header/menu_test.go`) import both.

For QR codes, `testutil/qrtest.ScanSVG(t, svg)` paints the dark modules of an SVG drawn by `components/qrcode` into an image and decodes it with gozxing, so a test checks the drawn code, not only the encoder. `components/qrcode/qrcode_test.go` and `pages/admin/feedback_admin_qr_test.go` use it. Like `templtest`, it does not import `service`.

Query the rendered HTML with CSS selectors through [goquery](https://github.com/PuerkitoBio/goquery), for example `` doc.Find(`a[href="https://www.regncon.no/vanlege-sporsmal/"] .inline-icon`) ``. Prefer goquery over `golang.org/x/net/html` (more verbose) and over `strings.Contains` on raw HTML. Plain string checks are fine for smoke tests and for output that is not a DOM, such as the sanitized Markdown fragments in `service/eventService/event_helpers_test.go`.

## Conventions

These add to the rules in [AGENTS.md](../AGENTS.md#go-test-structure).

- **One behavior per test**, with a name specific to that behavior. Declare expected values as named variables (for example `expectedHrefs`) at the top of `// Given`, before other setup.
- **Allow-list href tests per role.** The header menu tests (`components/header/menu_test.go`) have one test per role: anonymous (`/`, `/auth`), logged-in non-admin (`/`, `/profile`, `/tilbakemelding`, `/auth/logout`, FAQ) and admin (those plus `/admin`, `/admin/billettholder/`, `/admin/approval/`). Each collects the unique hrefs and compares them to the exact expected set with `templtest.AssertSameHrefs`. One exact-set comparison covers both "expected links exist" and "no unexpected links exist", so these tests need no separate absence checks. Only the placement of the feedback link has its own tests (`TestMenu_LoggedInUserSeesFeedbackLinkOnlyInBurgerMenus`, `TestMenu_AnonymousUserDoesNotSeeFeedbackLink`), and they select it by `href`. Absence checks break when links are renamed and are unreliable when several links share the same text (for example dropdown "Admin" entries). Assert on `href` rather than link text, since it is more stable. The external FAQ link `https://www.regncon.no/vanlege-sporsmal/` is part of the expected set, and the logged-in and admin tests also check that it shows the external-link icon.
- **Render tests assert visible text.** Assert on what the user sees, not on internal struct fields. `components/profile/my_program_render_test.go` renders `MyProgram` with `templtest.Render` and checks the visible text (`profileProgramVisibleText`, which joins `strings.Fields(doc.Text())`); `templtest.CollectTexts` on a CSS hook works too. Check both that the expected text is present and that text that should be hidden is absent.
- **Test first for behavior changes.** Write the failing behavior test, confirm it fails for the expected reason, then change the production query or template. Mitt festivalprogram is pinned at two layers in `components/profile`: query tests on `GetAllEventsForUser` (`my_program_test.go`) and render tests on `MyProgram` (`my_program_render_test.go`). Current behavior:
  - A player assignment made by the solver in a pulje that is not Completed (open or locked) is hidden, and the billettholder's interesser are shown instead.
  - In a Completed pulje the assigned arrangement is shown.
  - GM assignments are always shown, and win over interesser.
  - Manual player assignments are shown immediately, also in open or locked puljer.
  - Player and GM on the same arrangement gives one card, marked GM.
  - If the program is not published, the player result is hidden and an info text is shown.
- **Split large test files by behavior**, not by line count. Move long fixture, helper and fake code into a companion `<topic>_test_helpers_test.go` in the same package. Examples:
  - `pages/event/event_interest_test.go`, `event_interest_update_test.go`, `event_interest_test_helpers_test.go`
  - `pages/root/root_page_test.go`, `root_page_program_published_test.go`, `root_page_program_unpublished_test.go`, `root_page_test_helpers_test.go`
  - `components/profile/my_program_test.go`, `my_program_render_test.go`, `my_program_test_helpers_test.go`
  - `service/live/live_test.go`, `live_test_helpers_test.go` (fake JetStream KV)
  - `service/eventService/previous_next_root_list_test.go`, `previous_next_root_list_test_helpers_test.go`
  - `service/rooms/*_test.go`, `rooms_test_helpers_test.go`
- **No production helpers that only tests call.** golangci-lint runs with `tests: false` (`.golangci.yml`), so a function used only from `_test.go` files is reported by the `unused` linter. Test the real entry point instead. For example, the event visibility tests call `decideEventView(...)` and assert on the returned decision (`.CanView` and other fields); the `TestCanViewEvent_` name prefix is historical. Put setup helpers in `_test.go` files or `testutil`.

## Notable covered behaviors and deliberate gaps

- **Friendly errors on the front page.** `pages/root/root.go` defines `rootPageLoadErrorMessage` and `rootEventsLoadErrorMessage`. `root_page.templ` and `event_list.templ` render these constants instead of raw errors, and `root_page_test.go` uses the same constants to check that the friendly message is shown and `Error fetching` is not.
- **Creating an arrangement from Min side.** `POST /profile/api/create` returns 401 without a logged-in user, and HTTP 500 with `createEventFailureMessage` if the `user_id` lookup or the insert fails, never an empty 200. On success it redirects to `/profile/new/{id}`. The logic lives in `createNewEventFormSubmissionForUser` so tests can pass an explicit user. `TestCreateNewEventFormSubmission_WhenInsertFails_ReturnsFriendlyError` drops the `events` table to force the failure.
- **Post-login sync.** `GET /auth/post-login` reads email and user ID from the Descope token and calls `syncPostLoginUser(db, userID, email, isAdmin, logger)`. If no local user with that email exists, one is inserted (`external_id`, `email`, `is_admin`); `is_admin` is then updated from the token. The return target comes from the `neste` query parameter, which `safeReturnPath` only accepts as a same-site relative path and otherwise replaces with `/`. A sync error redirects to `/auth` (keeping `neste`), otherwise to the return target. `pages/login/post_login_user_test.go` calls `syncPostLoginUser` directly, so it covers this app-owned behavior without Descope; the `neste` handling is covered by `TestPostLogin_*Neste*` and `TestSafeReturnPath_*` in `pages/login/login_test.go`.
- **Breadcrumbs.** `components/breadcrumbs_test.go` tests the shared `Breadcrumbs` component: parent crumbs are links, the current crumb is `.breadcrumb-end` text, separators appear, a `.breadcrumb-mobile-return` link points to the parent, and a current-only path renders no navigation links. Each page still needs its own small breadcrumb assertion (for example `TestRootPageContent_RendersHomeBreadcrumb` and `TestAdminPage_RendersBreadcrumb`), because the component test cannot prove which crumbs a route passes in.
- **No tests for an empty pulje or an empty front page.** The app is never expected to run with zero arrangementer in a pulje or zero announced arrangementer.
- **No guard tests against placeholder text** (for example "Game System" or "Arrangørnavn") on public event cards. The event form already requires those fields, and an admin who overrides validation is responsible for the content.

## Coverage gaps and known issues

Components with render tests include `components/event_card_test.go`, `components/breadcrumbs_test.go`, `components/event_components/programpulje_interests_test.go`, `components/profile/my_events_test.go`, `components/profile/my_tickets_test.go` and `EventInterestPanel` in `pages/event/event_interest_test.go`.

Still without direct component tests:

- `components/previous_next.templ`
- `components/breadcrumb_nav_buttons.templ`
- `components/ui/button/eventStatus.templ`
- `components/event_components/programpulje_event.templ`
- `components/formsubmission/statusCard.templ`

Known issue: in `statusCard.templ` the admin branch renders every status `<option>` with `selected`, so the browser initially shows the last one (Archived) regardless of the arrangement's status, until Datastar's `data-bind` overrides it.

Templates that query the database or read the request during render are harder to component-test. Extract a data-only partial first and test that.
