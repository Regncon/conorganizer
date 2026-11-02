---
name: conorganizer-datastar-live
description: Conorganizer's rules for Datastar and live updates (NATS buckets, SSE) in .go and .templ files. Use whenever you add or change a POST/PUT/DELETE handler that writes events, interests, billettholdere, puljer, tildelinger or rooms; make a page live or add an SSE endpoint (live.Page, liveManager.Stream, Broadcast, *_index/*_page templ pairs); call datastar.NewSSE, ReadSignals, PatchSignals or PatchElementTempl; or put data-signals, data-on:, data-effect, data-init, a native dialog, a details element or a tooltip inside live HTML. Also use when the user says a page doesn't update until reload, other tabs or admins don't see a change, a dialog closes or state resets on update, the page becomes unclickable, a signal is empty on the server or an error message doesn't clear, and when reviewing a PR or diff that touches any of this. Use it even when live updates are not the main task, since most mutation handlers in this repo must broadcast.
---

# Conorganizer Datastar Live

One model runs the whole app: a mutation writes to SQLite, broadcasts a NATS bucket, and every open page subscribed to that bucket re-renders its live fragment over SSE from the database. Most bugs in this area are a missing or mis-targeted broadcast, a response written in the wrong order, or browser state that a morph wipes out. This skill keeps new code on that model. AGENTS.md has the short version under "Live updates and Datastar"; this is the depth.

## Read first

The documentation is the contract. The code wins when they differ.

1. `documentation/live-update-lifecycle.md`: "Lifecycle", "Broadcast Lifecycle", "Buckets" with "Cross-bucket broadcasts" and "Page subscriptions", "Full content first", "Native dialogs and live fragments", "Testing Strategy", "LLM Implementation Contract".
2. `documentation/datastar-signals.md`: "Choose who owns the value", "Effects and requests", "Keep browser-owned state through morphs", "Native modal dialogs in live fragments", "Datastar DOM events".
3. `service/live/live.go`: the bucket constants, `Page`, `Stream`, `Broadcast`, `EnsureConnection`, `DatastarInit`.
4. Run the live map from the repo root before you trust the doc tables, because they drift:

```bash
bash .agents/skills/conorganizer-datastar-live/scripts/live-map.sh
```

It prints, per bucket, every `Stream` subscriber with its local route (relative to where its setup function is mounted, so `/tickets/api/` is `/profile/tickets/api`) and every `Broadcast` call site. Then it lists the POST/PUT/DELETE/PATCH routes that broadcast only through a function they call, with the call chain, and the routes that never broadcast. Open each route in that last list before you call it a missed broadcast. The last two sections list handlers that call `http.Error` or `ReadSignals` after `NewSSE`, and Datastar attributes the client ignores or misreads. Nothing in those two is a pattern to copy.

Versions: the bundled client is Datastar v1.0.4 (`static/datastar.js`) and the Go SDK is `datastar-go v1.2.2` (`go.mod`). Syntax from older Datastar releases or blog posts does not apply; see `references/datastar-cheatsheet.md`.

## Add or change a mutation handler

Handler order, top to bottom:

1. When the handler needs signals, `datastar.ReadSignals(r, &store)`. JSON tags are the camelCase signal names. GET and DELETE read the `datastar` query parameter, other methods read the body.
2. Check access, then validate, and reject with `http.Error(...)` here, while the response is untouched. Input the user can fix is the exception: it gets a 200 `errorfeedback` patch (see below).
3. Write to SQLite.
4. `liveManager.Broadcast(r.Context(), buckets...)` for every bucket whose subscribers render the changed data.
5. Only now `sse := datastar.NewSSE(w, r)`, and only if the request needs a signal patch, field feedback or a redirect. When nothing has to reach the requester because the live stream re-renders the page, answer `w.WriteHeader(http.StatusNoContent)`, as the pulje status routes do.

`NewSSE` sets `Content-Type: text/event-stream` and flushes the headers, so the status is already 200. A later `http.Error` appends plain text that Datastar ignores, and Go logs a superfluous `WriteHeader`. After `NewSSE`, report problems through `errorfeedback`, a signal patch or `sse.ConsoleError(err)`.

`ReadSignals` must come before `NewSSE` too. When the headers are flushed, Go's HTTP/1 server discards the unread request body (`net/http`, `maxPostHandlerReadBytes`), so a POST or PUT that reads signals after `NewSSE` gets nothing on a real server. `httptest.NewRecorder` does not discard the body, so a recorder test hides this. Some older handlers break one of these orders; the live map lists them, so check a handler there before you copy it.

An `http.Error` before `NewSSE` never shows its text. For a non-200 status the 1.0.4 client fires `started`, `error`, then `finished` and drops the body. A trigger with `errorfeedback.Attrs(key)` and a matching `RenderMessage` then shows a fallback message (its own, the form's default, or the generic one), and a toast indicator on it stays quiet; a trigger without a feedback key shows nothing at all. Give every mutation trigger that can fail a feedback key.

When the user can fix the problem and needs the specific reason ("Epostadressen finnes allerede"), answer 200 with an `errorfeedback` patch instead. Keep the order above and add one branch: check access first, then validate before the write, and when there are errors, open the SSE, patch them and return without writing. A field error sent before the access check answers a user without access with 200, and a check that reads the database tells them about data they may not see. On success, write and broadcast, then open the SSE and patch `feedback.Signals()`, which sends every key the response owns as `""`. `http.Error` stays available for everything before an SSE is opened, a failed broadcast included. Read the signals into `map[string]any` or a struct: the browser also sends the `feedbackErrors` object, which `map[string]string` cannot hold. `references/live-page-skeleton.md` sections 7 and 8 have the handler and its tests.

Choosing buckets:

| You changed | Broadcast |
| --- | --- |
| Events, event form fields, status, submission, program publishing, puljer, pulje status or closing warning | `live.BucketEvents` |
| Interests, førstevalg, tildelinger, who is interested in what | `live.BucketInterests`; tildelinger also `events` and `rooms` |
| Billettholdere, their e-mails, their links to users, ticket fetch or conversion | `live.BucketBillettholders` |
| Rooms | `live.BucketRooms`; room assignment and its publishing also `events` |

Some existing handlers broadcast fewer buckets than this table. Do not copy a bucket list from a neighbouring handler; check it against the map. Data outside these domains, such as feedback, gets its own bucket when a page streams it (see "Make a page live").

Pages subscribe to several buckets at once (`/` streams `events`, `interests` and `billettholders`; `/event/{id}` streams `events`, `interests` and `rooms`; the header menu streams `billettholders` on every page). The live map shows the current list, so check it rather than guessing from the page name.

The map lists subscriptions, not what each renderer reads. To know which pages show the data you changed, grep for the service function or table behind it (for example `GetTicketHolders` or `relation_billettholdere_users`) in `*.templ` and `*.go` outside `*_templ.go`, and compare the pages you find with the map. A page that renders the data without subscribing to the bucket is a gap: add the bucket to its `Stream` call when the page is yours to touch, otherwise name the gap in your notes.

Required or best effort: by default a failed broadcast is a failed request. Log it and answer `http.Error(w, "Failed to broadcast update", http.StatusInternalServerError)`, as the pulje status routes in `pages/admin/puljefordeling.go` do. When the requester's only feedback is the live re-render (an item disappears, a list grows), the broadcast must be required, or a failure looks like nothing happened; the user then sees an error for a change that is already saved, and a reload shows the truth. Treat it as best effort (log and still succeed) only when the SSE response has already started or when the mutation must not look failed because the notification failed, as the puljefordeling tildeling handlers and the billettholder e-mail responders do.

A broadcast that outlives the request must not use `r.Context()`: pass `context.WithoutCancel(r.Context())` into the goroutine, as `pages/login/login.go` does for the post-login ticket sync.

Redirect after a mutation: a Datastar request runs through `fetch`, which follows a 303 on its own. Datastar then morphs the target page's HTML into the current page by matching ids, and the URL does not change. Branch on the header and send `sse.Redirect` instead, as `components/formsubmission/submit_section.templ` does:

```go
if r.Header.Get("Datastar-Request") == "true" {
	sse := datastar.NewSSE(w, r)
	if err := sse.Redirect("/profile/"); err != nil { ... }
	return
}
http.Redirect(w, r, "/profile/", http.StatusSeeOther)
```

Do not patch elements the live stream already re-renders. If the page is subscribed, the broadcast refreshes it on every open tab; an element patch in the mutation response only reaches the requester and gives two renders that fight. Patch elements from a mutation only for regions no live page streams, such as the feedback form after a submit in `pages/feedback/feedback.go`. Check the live map first: a region that has become live has lost that exemption.

Closing the requester's dialog: close it from the response, after the write succeeded. Patch the signal that drives a signal-driven dialog to its closed value, or send `sse.ExecuteScript` with `close()` for a dialog opened directly, as the room dialog in `pages/admin/admin.go` does. Render the `errorfeedback` message inside the dialog. A trigger that writes `$open = false; @post(...)` closes before the result is known, so a failure message in the dialog is never seen.

Update the doc tables in the same change: "Typical broadcasters" under "Buckets", and "Page subscriptions".

## Make a page live

Copy `pages/profile/tickets`. Do not copy `admin_index.templ`, `billettholder_admin_index.templ` or `add_billettholder_index.templ`; they render a placeholder and are listed under "Full content first" as exceptions to migrate.

1. `<name>_index.templ`: the normal GET renders `layouts.Base(...)` with a wrapper that holds the full content, for example `<div id="profile-tickets-page-wrapper" data-init={ live.DatastarInit("/profile/tickets/api") }>`. If SSE never connects, the user still has the page.
2. `<name>_page.templ`: a `templ` component that renders the same wrapper id, `<div id="profile-tickets-page-wrapper">`, around the same content component, without `data-init`. The live patch is an outer morph matched by id; with a different id the patch lands nowhere.
3. The `/api` GET route calls `liveManager.Stream(w, r, live.Page{Buckets: ..., Render: func(ctx context.Context, r *http.Request) templ.Component { ... }})`. Validate URL parameters and `http.Error` before calling `Stream`, as the room assignment route in `pages/admin/admin.go` does. `Stream` ensures the `connections` session, creates the SSE, sends one full patch, then watches the buckets.
4. Carry query state into the live URL. The root page passes `date`, the event page passes `r.URL.RawQuery`, and the profile page uses `live.DatastarInitExpression("'/profile/api' + window.location.search")`. The renderer reads it from `r` on every re-render.
5. Access: the router middleware for login and admin guards the `/api` route too, but per-resource checks the page GET does, such as event ownership, are not inherited. Guard the stream and every mutation under the same `/api` router. A chi middleware added with `Use` on that sub-router covers them all, as `newevent.RequireEventOwnerOrAdmin` does in `pages/profile/profile.go`.
6. The renderer is idempotent and reads SQLite every time. It runs on every broadcast, every tab-visibility change, and whenever another tab in the same browser touches the shared connection key ("Shared session keys across tabs").
7. A native `<dialog>` opened with `showModal()` belongs outside the live wrapper (room assignment, puljefordeling) or follows the signal pattern in `references/morph-safety.md`.
8. Add the page to "Page subscriptions".

A page that stays in one `.templ` file can hold both components; `domeneordbok.md` ("Kode") names the files `_index` (sets up the NATS integration) and `_page` (renders the entry id). Splitting an existing file renames it, so ask before you do it.

Add a bucket only when no existing bucket's subscribers render the data, and broadcasting an existing one would re-render many unrelated pages. A page whose data no bucket covers is that case; it is the "real correctness issue" the lifecycle doc allows. Add the constant and the `allBuckets` entry in `service/live/live.go`, a row in the "Buckets" table, and update every sentence in the lifecycle doc that counts the buckets (grep it for "four": "Naming and public API", "Buckets", "Service API" with its code block). The doc prefers a few broad buckets over fine-grained ones, so say in the change why the new one is needed.

Splitting or renaming a `.templ` file leaves the old `*_templ.go` behind until templ generates again. `go tool templ generate` deletes a generated file whose `.templ` is gone, so run it before you build or vet in a worktree.

`references/live-page-skeleton.md` has the index route, the wrapper, the page renderer, the `/api` route, a mutation handler in the right order, a handler that returns field errors, route tests for both, and an end-to-end stream test, all based on `pages/profile/tickets`.

## Datastar markup inside live HTML

- Colon syntax only: `data-on:click`, `data-bind:name`, `data-attr:disabled`, `data-class:selected`, `data-signals:name`, `data-indicator:_saving`. Modifiers go after the key and `__`: `data-on:input__debounce.300ms`, `data-on:click__window`, `data-signals:open-id__ifmissing`. The 1.0.4 client splits the name at `__` first, then at the first `:`, and looks the plugin up by that exact name. So `data-on-click` or `data-bind-search` is ignored without an error, and in `data-signals__ifmissing:open-id` the key lands inside the modifier and nothing is declared.
- Keys are kebab-case. HTML lowercases attribute names, so `data-signals:roomCount` declares `$roomcount`; write `data-signals:room-count` for `$roomCount`.
- Signals are global to the page. Markup rendered once per item (a card, a row) must not declare or bind a fixed signal name, or the items overwrite each other and every copy shows the last one's value. Put the item id in the name, as the billettholder card does with `newEmail-42`, or keep one shared dialog outside the list and pass the item id in a signal.
- `components.DataSignals(map[string]any{...})` writes several signals in one attribute with proper JSON escaping and keeps keys verbatim, so write `"puljeCanEdit"`. A `_` prefix marks browser-only state that Datastar leaves out of requests, which is why a `_`-signal "is empty on the server".
- `data-signals:name__ifmissing="..."` or the object form `data-signals__ifmissing="{...}"` for browser-owned state that must survive a live patch, such as open panels and selections. Plain `data-signals` applies its value on load and again only when a patch changes the attribute value; a patch that renders the same value keeps whatever the browser set.
- The live `data-init` is always `live.DatastarInit(url)` or `live.DatastarInitExpression(expr)`, never a hand-written `@get`, so retry settings stay consistent.
- `data-effect` runs initially and whenever a signal it reads changes. Keep its request from changing the signals it reads, or it loops. For one-shot setup inside live HTML use `data-init`: when a patch changes a `data-effect` expression, Datastar starts the new effect before disposing the old one, and two effects that talk through an event overflow the stack (`components/header/menu.templ`). `data-init` runs again only when a patch changes its expression, so put the data the setup depends on into the expression.
- `data-preserve-attr="open"` on `<details>` and on signal-driven `<dialog>`; `data-ignore-morph` only on regions a widget or a separate endpoint owns.
- Adding a `<dialog>` to an existing live page: a dialog opened with `showModal()`, or with `commandfor` and `command="show-modal"`, must sit outside the live wrapper (a sibling in the `_index` templ is enough; the trigger may stay inside) or follow the signal pattern. `references/morph-safety.md` has both.
- User text never goes into a Datastar expression through `fmt.Sprintf`; Go's `%q` is not JavaScript quoting. Put it in a `data-*` attribute and read `el.dataset.name` in the expression, or JSON-encode it with `templ.JSONString` or `components.DataSignals`.
- Tooltips: only `data-tippy-content="..."`. `layouts/base.templ` delegates from `document.body`, so morphed elements keep their tooltips; per-element `tippy(...)` instances die with the morph.
- Style blocks: a component rendered in a loop, or by a parent that repeats it, wraps its `<style>` in a package-level `templ.NewOnceHandle()`, as "Looped components render their style once" in `documentation/ui-conventions.md` describes. The handle is scoped to one `Render` call, and each live patch is its own call, so every patch carries the style once. Without the handle, every broadcast resends one style block per item. A component rendered once keeps a plain `<style>` block.

## Feedback to the user

- Field errors: `feedback := errorfeedback.New("title", "intro")`, `feedback.Set(key, msg)`, `feedback.Patch(sse)`, with the SSE opened only in the branch that patches (see "Add or change a mutation handler"). Markup: `{ errorfeedback.Attrs("title")... }` on the trigger and `@errorfeedback.RenderMessage("title")` beside it. Keys you do not pass to `New` are never cleared; that is how old messages get stuck.
- Saved toast: `@toast.Render("Lagret", icons.ProgressComplete)` once per page, outside the live fragment, then the value form `data-indicator="toastTitle"` on each trigger. `toast.js` reads `el.dataset.indicator`, so `data-indicator:toast-title` never shows a toast. It shows when the request finishes without `error` or `retries-failed`. A backend `feedbackErrors` patch suppresses it only when it has a non-empty message for the trigger's own `errorfeedback.Attrs` key.
- Per-request messages that are not field errors: patch plain signals, as `getTicketsErrorMessage` in `pages/profile/tickets/tickets_page.templ` does. Declare them kebab-case, `data-signals:get-tickets-error-message="''"`; that file writes `data-signals:getTicketsErrorMessage`, which declares `$getticketserrormessage`.

Details and the client-side lifecycle are in `references/errorfeedback-and-toast.md`.

## Tests

- A handler whose broadcast is required needs a real manager in its route test: embedded NATS plus `live.NewManager`, as in `pages/admin/rooms_assignment_route_test.go` and the `newMaxPlayersTestLiveManager` helper in `components/formsubmission/other_details_test.go`. `&live.Manager{}` has no buckets, so `Broadcast` returns `unknown live bucket`; it only fits handlers that treat the broadcast as best effort, tests that stop before it, or a test that checks the broadcast-failure path on purpose.
- A route test that checks status and database still passes when the `Broadcast` line is deleted. To prove that another open page re-renders, open the stream over `httptest.NewServer`, read the first patch, wait until the stream's NATS watcher exists, then run the mutation and read the next patch. A broadcast between the first patch and the watcher start is missed by design, so a test without the wait is flaky. `references/live-page-skeleton.md` has the helpers.
- Server-side feedback is asserted from the response body: each `data: signals` line carries a JSON patch, and `feedbackErrors` is nested in it. Test the error path (message present, nothing written) and the success path (the key patched to `""`). `references/live-page-skeleton.md` section 8 has the helpers; the flat `profilePatchedSignals` helper cannot read the nested map.
- Markup rules are asserted with `templtest.Render`, as `TestBillettholderCard_UsesDatastarStateForInterestDialogOpenState` does for the dialog attributes. `templtest.SignalsObject` parses a JSON `data-signals` attribute and `templtest.SignalsAttrsWithModifierBeforeKey` finds `__ifmissing` written before the key; `pages/admin/puljefordeling_signals_test.go` uses both.
- `service/live/live_test.go` uses an in-memory fake KV. Keep NATS out of it.
- Use `bdd.Behavior` and Given/When/Then sections as AGENTS.md describes.

## Troubleshooting

| Symptom | Look for |
| --- | --- |
| Page updates only after reload | No `Broadcast` in the handler, wrong bucket, page not subscribed to that bucket, or `_index` and `_page` wrapper ids differ |
| Other tabs or admins do not see a change | Mutation patched elements for the requester instead of broadcasting |
| Dialog closes, or page goes unclickable after an update | `<dialog>` opened with `showModal()` inside the live fragment |
| Open panel or selection resets on update | Plain `data-signals` whose rendered value changed since the last patch, `__ifmissing` written before the key, or missing `data-preserve-attr` |
| Every card shows the same value, or saves another card's | One fixed signal name declared or bound in markup that repeats per item |
| A button or input does nothing, with no console error | Hyphen-form attribute such as `data-on-click`; the live map lists them |
| Signal is empty on the server | `_` prefix, kebab/camel name mismatch with the JSON tag, camelCase after the colon that HTML lowercased, or the signal was never declared |
| Signal has no initial value, or signals named `0`, `1` and so on appear | `data-signals__ifmissing:name` with the modifier before the key; write `data-signals:name__ifmissing` |
| Error message does not clear, or toast shows next to an error | Key missing from `errorfeedback.New`, a plain-signal message next to `feedbackErrors`, or a trigger whose `errorfeedback.Attrs` key is not the key the backend patched |
| Saved toast never shows | `data-indicator:toast-name` instead of `data-indicator="toastName"`, or no `toast.Render` on the page |
| A rejected request shows nothing, or only a generic message | `http.Error` drops its text; the trigger needs `errorfeedback.Attrs` for the fallback, or the handler a 200 `errorfeedback` patch for the specific reason |
| A mutation replaces parts of the page with another page | `http.Redirect` answered a Datastar request; use `sse.Redirect` |
| Save "works" in a recorder test but nothing is saved on the real server | `ReadSignals` after `NewSSE`; the flushed headers made Go discard the body |
| A live page never updates locally, and the stream answers 500 `failed to ensure live session` | A `connections` cookie signed by another local server on the same host; Datastar does not retry a non-2xx. Clear the cookie or use a distinct `*.localhost` host |
| Live pages stop refreshing until restart | `live watcher closed; closing stream for reconnect` in the logs; see "Troubleshooting" in the lifecycle doc |
| `failed to touch live key ... context deadline exceeded` | Several instances sharing one `-nats-store-dir` |

Live logs carry `component=live`; the message table is under "Logging" in the lifecycle doc.

## PR review checklist

1. Every handler that writes events, interests, billettholdere, puljer, tildelinger or rooms broadcasts, and the buckets cover every subscriber that renders the data (run the live map).
2. `ReadSignals` and `http.Error` only before `datastar.NewSSE`; feedback after it goes through signals.
3. Broadcast failure handling is deliberate: 500 by default, best effort only with a reason.
4. New live page: full content on GET, same wrapper id in `_index` and `_page`, `live.DatastarInit`, query state carried, access checked on the stream and every mutation.
5. No `showModal()` or `command="show-modal"` dialog inside a live fragment without the signal pattern; `data-preserve-attr` on `<details>`; `data-ignore-morph` justified.
6. Colon syntax with modifiers after the key, kebab-case keys, `_` prefix for browser-only signals, `__ifmissing` where state must survive a patch, `components.DataSignals` instead of hand-quoted JSON.
7. No user text formatted into a Datastar expression; per-item signal names in markup that repeats, and its `<style>` rendered once through `templ.NewOnceHandle()`.
8. `errorfeedback.New` lists every key the response owns; every trigger that can fail has a feedback key; input the user can fix is rejected with a 200 `errorfeedback` patch that says why, not with `http.Error`; toasts use the value form of `data-indicator`.
9. Route tests use embedded NATS when the broadcast is required; a test that passes with `&live.Manager{}` proves nothing about the broadcast.
10. Doc tables updated: "Buckets", "Cross-bucket broadcasts", "Page subscriptions".

## References

- `references/live-page-skeleton.md`: copy-ready index, page, `/api` and mutation code based on `pages/profile/tickets`, a handler that returns field errors, the embedded-NATS route tests for both, and the end-to-end stream test.
- `references/datastar-cheatsheet.md`: attributes, modifiers, actions and Go SDK calls for Datastar v1.0.4 and datastar-go v1.2.2, with one repo example each.
- `references/morph-safety.md`: dialogs, `<details>`, `data-ignore-morph`, effects, the `datastar-fetch` hook, tooltips and user text in expressions.
- `references/errorfeedback-and-toast.md`: the feedback and toast contract, server and browser side.
