# Access control and error pages

This document covers who can see an arrangement, what the page shows when they cannot, the shared 401, 403 and 404 pages, and how login returns the user to the page they tried to open. User-facing texts are quoted verbatim in Bokmål.

## Event visibility

`decideEventView` in `pages/event/event_visibility.go` decides whether the viewer may see an arrangement. It is used by:

- the full page `GET /event/{idx}` (`pages/event/event_index.templ`),
- the live SSE endpoint `GET /event/api/{idx}`, which renders `event_page` / `event_page_content` in `pages/event/event_page.templ`,
- `PUT /event/api/{idx}/interest/selected-interest` (`pages/event/event.go`). When the viewer may not see the arrangement, this endpoint returns plain text: `Arrangementet er ikke tilgjengelig.`

### Rules

Visibility depends on the viewer's login and admin status. The roles a person has on the arrangement itself do not count.

| Status | Admin | Owner | Other logged-in users and anonymous visitors |
| --- | --- | --- | --- |
| `Annonsert` | Full arrangement | Full arrangement | Full arrangement |
| `Kladd`, `Innsendt`, `Godkjent` | Full arrangement and the "ikke annonsert" banner | Full arrangement and the "ikke annonsert" banner | Hidden notice, HTTP 200 |
| `Forkastet` | Full arrangement and the "forkastet" banner | Full arrangement and the "forkastet" banner | Hidden notice, HTTP 410 Gone |

- **Admin** means `userInfo.IsAdmin` or `authctx.GetAdminFromUserToken(ctx)` (the `Admin` role in the Descope token).
- **Owner** means the logged-in user's token ID matches `users.external_id` on the `users` row whose `id` equals `events.user_id` (`eventOwnerMatchesUser`). Anonymous viewers, and arrangements with no `user_id`, never match.
- Being assigned to the arrangement as Player or GM does not make it visible. The person must also be the owner or an admin.

### What the viewer sees

Admin and owner get the normal arrangement page with a warning banner above it:

- Not announced: title `Arrangementet er ikke annonsert`, message `Dette arrangementet er ikke annonsert ennå. Vanlige brukere kan ikke se det.`
- Forkastet: title `Arrangementet er forkastet`, message `Dette arrangementet er forkastet og er ikke tilgjengelig for vanlige brukere.`

Everyone else gets the normal layout with a notice in place of the arrangement:

| Case | Page title and heading | Message | Status on `/event/{idx}` |
| --- | --- | --- | --- |
| Not announced | `Arrangementet er ikke annonsert ennå` | `Dette arrangementet er ikke annonsert ennå. Kom tilbake senere, så får du se hva som venter.` | 200 |
| Forkastet | `Arrangementet er ikke tilgjengelig` | `Dette arrangementet er ikke tilgjengelig lenger.` | 410 Gone |

The messages are the constants `eventNotAnnouncedMessage` and `eventArchivedMessage`, and the page title comes from `eventPageTitle`. The SSE endpoint renders the same notice but does not set these status codes. The status codes apply only to the full-page route.

A hidden arrangement never gets the generic 404 page. `/event/{idx}` calls `notfound.RenderEvent` only when no arrangement has that ID. See [404 page](#404-page).

### Front page

The front page only lists `Annonsert` arrangementer:

- Program not published: `program.GetAnnouncedEvents` lists them alphabetically (`title COLLATE NOCASE`, then `id`).
- Program published: the lists come from the view `v_events_by_pulje_active`, which keeps rows with `e.status = 'Annonsert'` and `relation_event_puljer.is_in_pulje = 1`. The legacy `is_published` column is selected but not filtered on. See `migrations/20260919140000_add_event_puljefordeling.sql` and `service/program/events.go`.

## Visibility UI in `event_page.templ`

There are two separate components.

**Hidden notice** (`eventHiddenNotice`, used by `eventNotAnnouncedHiddenNotice` and `eventArchivedHiddenNotice`). It is rendered instead of the arrangement content and shows a centered warning icon, an `h1` title and the message. It carries these CSS hooks:

- the state class on the section: `.event-not-announced` or `.event-archived`
- the message paragraph: `.event-not-announced-message` or `.event-archived-message`

`pages/event/event_visibility_test.go` selects the message hooks, so do not rename them.

**Warning banner** (`eventVisibilityWarning(title, message)`), shown only to admin and owner. It is a two-column grid with `role="status"`. The left column is a 4.75rem icon rail that shows `icons.Warning` (from `components/icons`) scaled to 3rem, so the warning stands out. The right column holds the title (`.event-visibility-warning-title`, which the tests select) and the message (`.event-visibility-warning-message`).

## 401 and 403 pages

Both pages are in Bokmål and render inside `layouts.Base`.

| | 401 Unauthorized | 403 Forbidden |
| --- | --- | --- |
| Meaning | Not logged in | Logged in, but without the `Admin` role |
| Component | `userctx.Unauthenticated(loginHref)` (`service/userctx/unauthenticated.templ`) | `authctx.Forbidden()` (`service/authctx/forbidden.templ`) |
| Rendered by | `userctx.UserMiddleware` | `userctx.AdminForbiddenHandler(db, logger)` |
| Page title | `Logg inn` | `Ingen tilgang` |
| Heading | `Du må logge inn` | `Du har ikke tilgang` |
| Text | `Logg inn for å se denne siden.` | `Du er logget inn, men denne siden krever administratortilgang.` |
| Links | Primary `Logg inn` to `/auth?neste=<requested path and query>` (plain `/auth` for `/`), outline `Gå til arrangementslisten` to `/` | Primary `Gå til arrangementslisten` to `/` |

`userctx.UserMiddleware` builds the `Logg inn` link with `loginHrefWithNeste(r)` (`service/userctx/unauthenticated.go`). It takes the request's path and raw query and URL-encodes them into the `neste` parameter, so `/tilbakemelding?om=festivalen` gives `/auth?neste=%2Ftilbakemelding%3Fom%3Dfestivalen`. See [Login return target](#login-return-target-neste).

### Wiring and the import cycle

`service/userctx` imports `service/authctx`. That means `authctx` cannot import `userctx` or the layout helpers, because that would create an import cycle. Instead, `authctx.RequireAdmin(logger, opts ...RequireAdminOption)` accepts an optional `authctx.WithForbiddenHandler(http.HandlerFunc)`:

- A nil handler is ignored.
- Without the option, a non-admin gets a plain-text 403: `Du har ikke administratortilgang`.

The `Forbidden()` component lives in `authctx`. The handler that wraps it in the layout lives in `userctx`.

`router.go` builds the admin router like this:

```go
isLoggedInRouter := authenticatedRouter.With(userctx.UserMiddleware(logger, db))
routerAdmin := isLoggedInRouter.With(
	authctx.RequireAdmin(logger, authctx.WithForbiddenHandler(userctx.AdminForbiddenHandler(db, logger))),
)
```

`admin.SetupAdminRoute` and `billettholderadmin.SetupBillettholderAdminRoute` are mounted on `routerAdmin`. An anonymous visitor therefore gets the 401 page first, and a logged-in non-admin gets the HTML 403 page. This includes the feedback admin routes under `/admin/tilbakemeldinger/`, which `admin.SetupAdminRoute` registers. `profilepage.SetupProfileRoute` and `feedbackpage.SetupFeedbackRoute` (`/tilbakemelding`) are mounted on `isLoggedInRouter`, so they only require login.

Some nested routes call `authctx.RequireAdmin(baseLogger)` again without the option. Examples are the `/admin/billettholder/api/` and `/admin/billettholder/add/api/` streams in `pages/admin/billettholder_admin/billettholder_admin.go`. These inner checks would fall back to the plain-text 403, but the outer `routerAdmin` check stops non-admins before they get that far.

## Login return target (`neste`)

`neste` ("next") is the query parameter that carries where login should send the user afterwards. It is handled in `pages/login/login.go` and `pages/login/login.templ`:

- `GET /auth` reads `neste` and validates it with `safeReturnPath`. The login form (`loginForm(neste)`) puts the value on `.descope-login-wrapper` as `data-login-neste`. After a successful Descope login, the script stores the session and goes to `/auth/post-login?neste=<encoded value>`.
- `GET /auth/post-login` validates `neste` again and redirects there with 303 See Other. If the user token is missing or the local user sync fails, it redirects back to `/auth?neste=<encoded value>` (`authPathWithNeste`), so a retry keeps the return target. For `/` or an empty value it redirects to plain `/auth`.
- A user who is already logged in and opens `/auth` gets `alreadyLogedIn(neste)`: `Velkommen tilbake!`, a five-second countdown and then `window.location.replace` to the return target. The fallback link `denne lenken` points to the same target.

`safeReturnPath` only accepts a same-site relative path. It returns `/` for an empty value, anything with a control character or whitespace, a value that does not start with `/`, a scheme-relative `//host` path (backslashes are treated as `/` first, as browsers do), and anything that parses with a scheme or host. The query string of an accepted path is kept.

Only the 401 page adds `neste`. The `Logg inn` button in the header links to plain `/auth`, so logging in from there lands on the front page.

Tests: `pages/login/login_test.go` (`TestPostLogin_ValidNesteRedirectsThereAfterLogin`, `TestPostLogin_UnsafeNesteFallsBackToFrontPage`, `TestSafeReturnPath_*`), `pages/login/login_form_test.go` (`TestLoginForm_CarriesNesteReturnTargetToPostLogin`), `service/userctx/unauthenticated_route_test.go` (`loginHrefWithNeste`) and `service/userctx/unauthenticated_test.go` (`TestUnauthenticated_LoginLinkCarriesGivenHref`).

## 404 page

`pages/notfound/notfound.templ` contains both the Go helpers and the templ component:

- `Render(w, r, db, logger, message)` writes 404 and renders `Page(message)` inside `layouts.Base` with the title `Siden finnes ikke`. An empty message falls back to `defaultMessage`: `Vi fant ikke siden du leter etter. Lenken kan være feil, eller siden kan ha flyttet.`
- `RenderEvent(w, r, db, logger)` calls `Render` with `Vi fant ikke arrangementet du leter etter. Lenken kan være feil, eller arrangementet finnes ikke lenger.`

### Wiring

`main.go` registers the global handler after `setupRoutes` succeeds:

```go
router.NotFound(authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	notfound.Render(w, r, db, baseLogger, "")
})).ServeHTTP)
```

It is wrapped in `authMiddleware` so that the header shows whether the user is logged in. In degraded mode, when the image directory check fails, `db` is nil or `setupRoutes` fails, `mountDegradedRoutes` sets its own plain 503 page as `NotFound` instead.

These full-page routes call `notfound.RenderEvent` when the arrangement does not exist:

- `/event/{idx}` (`pages/event/event_index.templ`)
- `/admin/approval/edit/{id}` (`pages/admin/approval/editForm/edit_form_index.templ`)
- `/profile/new/{id}` (`pages/profile/newevent/new_index.templ`). This route also returns 404 when the arrangement exists but belongs to a different user.

API, SSE and admin action endpoints still return plain `http.Error` text with 404, so Datastar and other client callers do not get HTML. Examples are `pages/admin/puljefordeling.go` (`Pulje not found`), `pages/admin/tildelinger.go` (`Fant ikke billettholder, arrangement eller pulje`) and `DELETE /admin/tilbakemeldinger/{id}` in `pages/admin/feedback_admin.templ` (`Fant ikke tilbakemeldingen.`).

### Content

The page is one centered column with:

- the heading `404 Not Found`
- an `h2` with `Siden finnes ikke`
- the message
- a `Til forsiden` button that links to `/`
- the image `/static/404.webp` (750x600, alt `HTTP Cats 404 Not Found`)
- a visible `figcaption` credit: `Image from HTTP Cats` links to https://http.cat/status/404, and `Original images by Tomomi Imura` links to https://girliemac.com/

The image is served locally. There is no remote fallback to http.cat and no separate attribution file. The credit links open in the same tab, so they have no `target` or `rel` attributes. The CSS is natively nested under `.not-found-page`.

### Kept simple on purpose

The page follows KISS/YAGNI. There are no option or data structs. Callers pass a custom message only when they need one, and `RenderEvent` is currently the only one that does. The package deliberately has no automated tests, so check the page visually by hand after changing it.
