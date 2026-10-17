## Names
- Code is in English, except the domain words in domeneordbok.md, which stay Norwegian: `getBillettholder`, not `getTicketHolder`.
- Write domain words without æøå in identifiers, as domeneordbok.md does: `Forstevalg`, `Lordag`.
- Some existing names translate a domain word: `event` (arrangement), `interest` (interesse), `room` (rom) and the `ticketholder` package (billettholder). They will be renamed in a later, separate change. Don't rename them as part of other work. New names use the domain word.

## Billettholder access
- A user is linked to a billettholder only through `relation_billettholdere_users` (user id ↔ billettholder id). Every read, picker and access check must use that link.
- Never match an email (`relation_billettholder_emails.email` or `users.email`) against the logged-in user to decide what they can see or do. An email is text a user typed, not a key.
- Emails may only be used to *create* links, in the CheckIn flows: "Hent billetter" and the login sync (`checkIn.ImportUserTickets`), admin order conversion, and adding or deleting a manual email. See documentation/billettholdere.md.
- Look users up by `users.external_id` (the Descope id), never by email.

## Read the docs first
Before changing an area, read its doc in documentation/ and update it in the same PR as the behavior change. When a doc disagrees with the code, the code wins.
- pulje status, publishing and the interest UI: pulje-status-and-publishing.md
- billettholdere and CheckIn: billettholdere.md
- rooms and room maps: room-assignment.md
- live pages and Broadcast: live-update-lifecycle.md
- signals and dialogs: datastar-signals.md
- access and the 401/403/404 pages: access-control-and-error-pages.md
- tests: automated-tests.md
- migrations: migrations.md
- CSS classes and focus states: ui-conventions.md
- deploy, CI and servers: deployment.md
- SQLite, startup checks and health endpoints: sqlite-and-health.md

Skills live in `.agents/skills/<name>/` (SKILL.md, references/, agents/openai.yaml). `go tool task setup:skills` links them into .claude/skills.

## Text and names
- User-facing text is bokmål, never nynorsk. The one nynorsk exception is the external link https://www.regncon.no/vanlege-sporsmal/.
- An event is "annonsert", never "publisert". "Publisert" is for the program, the puljefordeling and the romfordeling (domeneordbok.md). Forkastet is `EventStatusArchived`; there is no Avvist.
- When you change a user-facing string, grep the tests and the «»-quoted text in documentation/testing/*.md for the old text, and add a negative `!strings.Contains` check for it.
- Use made-up names in test fixtures (Kari Nordmann, Ola Nordmann, @example.com), never real people's names.

## Dev setup
- `*_templ.go` files are gitignored; commit only `.templ`. Never hand-edit `*_templ.go`.
- In the main checkout the dev server (`go tool task start`, which Docker also runs) regenerates them. Don't run `go tool templ generate`, `go tool task build`, `make generate` or `make check` next to it.
- A .claude/worktrees checkout has no watcher: run `go tool templ generate` there before `go build` or `go test`.
- Never open database/events.db with the host sqlite3 while Docker runs; the container then fails with "unable to open database file (14)". Use `docker compose exec -T webserver sqlite3 database/events.db "..."`. Copy a database with `.backup`, never cp.
- Run tests with `go test ./<pkg>/...`. Don't run `go tool task test` or `make test` on the Windows host: they run the host sqlite3 on events.db and rewrite schema.sql, and in a worktree without database/events.db they empty it. Never commit a schema.sql change without a migration; restore it with `git checkout -- schema.sql`.
- .claude/worktrees/ holds other sessions' checkouts. Scope repo-wide searches and replacements to `git ls-files` or exclude `.claude/`, and stage only your own files.
- Docker serves https://localhost:7331; port 7332 is only reachable inside compose. http://localhost:8080 is the native default.
- Put throwaway Go programs in tmp/ (gitignored).
- Never print .env values.

## Routes and access
- The router a Setup*Route receives in router.go sets its tier: publicRouter, authenticatedRouter, isLoggedInRouter or routerAdmin. Admin endpoints hang off routerAdmin, in the `/admin` block of SetupAdminRoute (pages/admin/admin.go) or in SetupBillettholderAdminRoute. A second `Route("/admin")` or `Route("/profile")` panics at startup, because chi refuses to mount an existing path.
- Route handlers often live in .templ files. When you look for a route or handler, grep `*.go` and `*.templ`, excluding `*_templ.go` and `.claude/`.
- Menu links exist twice in components/header/menu.templ (desktop dropdown and mobile dialog), and menu_test.go asserts the exact hrefs per role. A new admin tool also needs an `adminToolCard` in pages/admin/admin_page.templ, a graphic in static/admin/ and an entry in pages/admin/admin_landing_integration_test.go.
- Every SSE, API and mutation endpoint enforces the same owner-or-admin rule as its page GET. Hiding a control in templ is not a check.
- `userctx.GetUserRequestInfo(ctx).Id` is the Descope id (users.external_id), not users.id. Get users.id with `userctx.GetUserIDFromExternalIDFromContext`.
- Admin means the Descope `Admin` role (`userInfo.IsAdmin` or `authctx.GetAdminFromUserToken`). users.is_admin is only synced at login; never use it for authorization.

## Live updates and Datastar
- After a durable write, broadcast every live bucket whose subscribers render the changed data. Room assignment broadcasts rooms and events; tildelinger broadcast events, interests and rooms. The subscription table in live-update-lifecycle.md drifts, so trust the code, and update the table when you add a subscriber or broadcaster.
- Validate input and send `http.Error` before `datastar.NewSSE(w, r)`. NewSSE flushes the 200 at once, so after it report errors with an errorfeedback signal patch or `sse.ConsoleError`.
- Datastar is client v1.0.4 (static/datastar.js) with datastar-go v1.2.2. Use colon syntax (`data-on:click`, `data-init`). The old hyphen form (`data-on-click`) is silently ignored, and `data-on-load` does not exist.
- Signals starting with `_` stay in the browser and never reach `ReadSignals`.
- Pass server values into data-signals with `components.DataSignals(map)`, not `fmt.Sprintf("'%s'")`. Never leave `<pre data-json-signals>` debug dumps in templates.

## Go code
- Import edges that exist today: layouts imports components/header; userctx imports layouts and authctx; testutil imports service; pages import service packages. So header and layouts cannot import userctx, authctx cannot import userctx, package service cannot import testutil (testutil/bdd is fine), and service packages never import pages.
- New service packages take `db *sql.DB` like their neighbours, return wrapped errors without logging them, and export sentinel errors (like `feedback.ErrNotFound`) so handlers can map them to 404.
- golangci-lint runs with `tests: false`, so a production function that only tests call is reported as unused. Error strings start lowercase (ST1005).
- Don't reformat untouched code. Keep the blank line between stdlib and module imports.
- Don't add code comments unless they are needed. Keep the `// Given`, `// When`, `// Then` markers and notes on magic numbers.

## Database
- Goose is not a go.mod tool. Run `go run github.com/pressly/goose/v3/cmd/goose@v3.28.0 -env /dev/null -dir migrations ...` (documentation/migrations.md).
- A migration PR also commits the regenerated schema.sql. Databases are migrated in place with their data (main and demo by hand after the merge, previews by CI), so added columns need `NOT NULL DEFAULT`. Goose never re-runs an applied migration, so don't edit one that is already on main; add a new one.
- relation_events_players.source is lowercase `manual`/`solver` on purpose. Use `SourceManual` and `SourceSolver`, and capitalized values for new enum-like columns.
- Removing an event from a pulje keeps its relation_event_puljer row and room_id with is_in_pulje=0, so readers of assignments and rooms filter on is_in_pulje=1.
- Don't add age back into the puljefordeling solver (board decision, #618 reversed #524). Age is handled by `FinnAldersvarsler` and the "Plasser likevel" confirmation.

# CSS rules

### Required CSS style
- Use native CSS nesting.
- Use container queries for responsive component styling.
- Prefer component-based responsiveness over viewport-based media queries.
- Use modern CSS features where appropriate:
  - custom properties
  - logical properties
  - `:has()`
  - `:is()` and `:where()` to group selectors, like `&:is(:hover, :focus-visible)` instead of separate `&:hover` and `&:focus-visible` rules
  - `clamp()`
  - `min()`
  - `max()`
  - `color-mix()`
  - container query units like `cqw`, `cqh`, `cqi`, `cqb`
  - Use Range Syntax: Use modern comparison operators like @container (width > 500px)
  - place-content, place-items, place-self instead of align and justify
- Use `&` to attach to the parent selector, like `&:hover`, `&::after`, `&.selected` or `&[data-open]`, or to add context around it, like `.card:focus-visible &`.
- Do not use `&` for Sass-style name joining, like `&-title` or `&__item`. It does not work in native CSS.
- Keep selectors shallow and readable.
- Avoid Sass-style CSS that does not work natively in browsers.
- Avoid viewport media queries for components. Use `@container` instead. Media queries about the device, like `(hover: hover)`, are fine.
- In container and media queries write `(width > 600px)`, not `(min-width: 600px)`. The `min-width` property itself is fine.
- Use viewport media queries only for page-level/global layout concerns.

### Backgrounds and tokens
- Pick backgrounds by Level (the `--level-*` tokens in static/css/index.css): containers on base are Level 1; non-clickable boxes inside Level 1 are Level 2; clickable items on base use the Level 2 states; clickable items and form fields on Level 1 use Level 3; anything inside Level 2 is Level 3; Level 1 edges use `--level-3-border`. Ask when a case isn't covered.
- Keep the Figma design tokens in index.css even when unused; remove only legacy ones. Retired names must not return: `--bg-surface*`, `--bg-item*`, `--color-secondary*`, `--color-disabled`, `--item-flat-active`, `btn--outline`, `label-small`.
- Use `.surface-pane` for panels, not the legacy `.card` or the deprecated `.form-card` (documentation/ui-conventions.md).
- `style="display: none;"` next to `data-show` is intentional: it stops a flash before Datastar loads. Don't move it into CSS.
- Shadow-DOM web components adopt only index.css and buttons.css (`BASE_STYLE_URLS` in static/js/conorganizer.js), so card, chip, choice and layout classes don't apply inside them. Link a new global stylesheet in layouts/base.templ.
- Custom properties don't work in `@container` or `@media` conditions; write the literal `600px` or `1200px`.

## Go test structure
When writing or cleaning up Go tests, prefer behavior-focused tests that read from intent to verification.

- Add short BDD metadata at the top of each test when it improves readability. Prefer
  `bdd.Behavior(t, bdd.BDD{Given: "...", When: "...", Then: "..."})`
  for tests that should appear in the automated behavior report. Existing BDD-style
  comments are still supported while older tests are migrated:
  - Given some initial context,
  - When an event occurs,
  - Then ensure the expected outcome.
- Split each test into three visible sections:
  - `// Given`
  - `// When`
  - `// Then`
- Start the `// Given` section with the expected output or expected result. This makes the test's purpose clear before setup details.
- Keep test names specific to one behavior. Prefer splitting creation, duplicate/idempotency, validation, and error cases into separate tests.
- Move noisy setup SQL and repeated assertions into local test helpers when that makes the test read like behavior.
- `bdd.Behavior` must be a top-level statement in the Test function with plain string literals; cmd/testreport skips calls inside `t.Run` and concatenated strings.

Example shape:

```go
func TestService_EventProducesExpectedResult(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given an existing domain object.",
		When:  "When the service reconciles related data.",
		Then:  "Then the expected durable relation is stored.",
	})

	// Given
	expectedResult := models.SomeResult{ID: 123}
	db := createTestDB(t)
	insertRequiredFixture(t, db, expectedResult.ID)

	// When
	err := ServiceFunction(expectedResult.ID, db)

	// Then
	if err != nil {
		t.Fatalf("expected service call to succeed: %v", err)
	}
	assertExpectedResult(t, db, expectedResult)
}
```

### Production helpers vs test helpers
Reuse production functions for the behavior under test. Do not add production service functions only to make test setup nicer.

- Use production helpers when they represent real application behavior that the test is exercising.
- Use test helpers for arranging prerequisite database state, such as inserting fixture users, billettholdere, emails, or association rows.
- If setup helpers become useful across packages, move them to `testutil`, not to production services.
- Avoid arranging a test through the exact code path being verified unless the test is explicitly an integration or route test.
- Route-level tests can exercise handlers and their production wiring, but keep service-level regression tests focused on the service behavior.

## When creating custom components that use signals

Datastar interacts with custom web components using a "props down, events up" model:
- Data flows into the component via HTML attributes.
- Changes inside the component are sent back via custom events.

### Interaction methods

#### Passing data in (Datastar -> Web Component)
- Use the `data-attr` attribute to bind a Datastar signal to a custom element attribute.
- The web component must observe that attribute (via `attributeChangedCallback`).

Example:
```html
<my-component data-attr:name="$_name"></my-component>
```

#### Getting data out (Web Component -> Datastar)
- Dispatch a custom event from the web component when internal state changes.
- Datastar listens with `data-on` and updates a signal.
- `evt.detail` or `$_eventDetails` provides the payload.

Example:
```html
<my-component data-on:custom-change="$_signal = evt.detail.value"></my-component>
```

### Key attributes and variables
- `data-attr:[attribute-name]`: Binds a Datastar signal to a component attribute.
- `data-on:[event-name]`: Listens for a (custom) event from the component.
- `evt`: Event object available in Datastar expressions.
- `evt.detail` or `$_eventDetails`: Event payload.

### Example flow
1. A user interacts with a `<reverse-component>` custom element.
2. The component reverses a string internally.
3. The component dispatches a `reverse` event with the result in `event.detail`.

JavaScript inside the component:
```javascript
// Inside your web component class
const event = new CustomEvent('reverse', {
  detail: { value: reversedValue }
});
this.dispatchEvent(event);
```

Datastar binding:
```html
<reverse-component
  data-on:reverse="$_reversed = evt.detail.value"
  data-attr:name="$_name">
</reverse-component>
```

This sets the `$_reversed` signal to the value from the event payload.

For more details and live examples, refer to the Datastar documentation on web components.

### Sub Agents
Use subagents when this task would be faster or more reliable with them.

## Festival year
Year-bound values change every festival:
- the CheckIn event id in the GraphQL query and `regnConDate` (the 18+ cutoff) in service/checkIn/cache.go
- `TicketTypeMiddag` in service/checkIn/ticket.go
- `checkInTicketsURL` in pages/profile/tickets/tickets_page.templ
- "Regncon 2026" in pages/login/login.go

## Git
- Commit subjects are imperative English and describe the user-visible outcome, keeping Norwegian domain words. No Co-Authored-By trailer.
- Bring a branch up to date by merging main into it, not by rebasing.
- Don't edit the preview block CI adds to the PR body (between the `preview-deployment-url` markers).

## finishing instructions
When done reading this you will output "Done reading AGENTS.MD file."
