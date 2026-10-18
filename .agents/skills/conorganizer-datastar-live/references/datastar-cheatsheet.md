# Datastar cheatsheet for Conorganizer

Client: Datastar v1.0.4, bundled as `static/datastar.js` (first line of the file). Go SDK: `github.com/starfederation/datastar-go v1.2.2` (`go.mod`). Check both before trusting syntax from the web; Datastar changed its attribute syntax before 1.0.

The client-side behaviour below was checked in a browser against the bundled file on 2026-10-18, with a check page and a small server sending real SSE patches. Re-check it the same way after a client upgrade; reading the minified bundle is not enough.

## How the client reads an attribute

`data-<plugin>:<key>__<modifier>.<arg>`. The client splits the name at `__` first, then at the first `:`, and looks the plugin up by exact name. So the key goes before the modifiers: `data-signals:open-id__ifmissing="''"` declares `$openId`. In `data-signals__ifmissing:open-id` the key lands inside the modifier: the attribute is read as `data-signals` without a key, declares nothing, is not if-missing, and spreads a string value into signals named `0`, `1` and so on. The object form `data-signals__ifmissing="{openId: ''}"` has no key and works.

Registered attribute plugins in the bundle: `attr`, `bind`, `class`, `computed`, `effect`, `indicator`, `init`, `json-signals`, `on`, `on-intersect`, `on-interval`, `on-signal-patch`, `ref`, `show`, `signals`, `style`, `text`. `peek`, `setAll` and `toggleAll` are actions, written `@peek(...)`; a `data-peek` attribute does nothing.

So `data-on-click` is read as plugin `on-click`, which does not exist, and is dropped silently. The live map's last section lists hyphen-form and modifier-before-key attributes still in the repo. Do not copy a hit.

Keys are kebab-case in HTML and camelCase in expressions: `data-signals:room-count="3"` creates `$roomCount`. HTML lowercases attribute names, so `data-signals:roomCount` creates `$roomcount`.

## Attributes in use

| Attribute | Repo example | Note |
| --- | --- | --- |
| `data-signals:name="expr"` | `data-signals:selected-interest-level={ fmt.Sprintf("%q", interest) }` in `components/event_components/event_interests.templ` | One signal. Applies on load and again when a patch changes the attribute value. A patch that renders the same value keeps the browser's value. |
| `data-signals={ components.DataSignals(map) }` | `components/event_components/event_interests.templ` dialog | Several signals, JSON-escaped, keys verbatim (`"puljeCanEdit"`). Nested maps become `$crop.ready`. A keyless `data-signals` needs an object: a string value, such as the `"errorMessage-42"` the billettholder cards write, is spread into signals named `0`, `1` and so on. |
| `data-signals:name__ifmissing="expr"` or `data-signals__ifmissing="{...}"` | object form in `components/header/menu.templ` with `templ.JSONString(...)` and in `pages/event/event_rooms.templ` | Only sets signals that do not exist yet, so a morph keeps the browser's value. The key goes before `__`. |
| `data-bind:name` or `data-bind="name"` | `data-bind="title"` in `components/formsubmission/about_event.templ` | Two-way binding for inputs. The quoted form is the same thing. |
| `data-text="$x"` | `data-text="$getTicketsErrorMessage"` in `pages/profile/tickets/tickets_page.templ` | Also render the initial text as content so the page has it before Datastar starts. |
| `data-show="$x"` | `data-show="$_fetching"` in the fetch button | Pair with `style="display: none;"` for content that starts hidden. |
| `data-attr:name="expr"` | `data-attr:disabled="$_fetching"` | Boolean false removes the attribute. |
| `data-class:name="expr"` | `data-class:billettholder-pending="!$_menuBillettholderReady"` in `menu.templ` | |
| `data-computed:name="expr"` | `data-computed:_image-busy="$_uploading || ..."` in `components/formsubmission/event_img_upload/event_img_upload.templ` | Derived read-only signal. |
| `data-style:property="expr"` | `data-style:background-color="..."` in `components/header/menu.templ` | |
| `data-on:event="expr"` | `data-on:click="@post('/profile/tickets/api/get-tickets')"` | `evt` and `el` are available in the expression. |
| `data-on:datastar-fetch__document` | `event_interests.templ` | Global Datastar fetch events, see below. |
| `data-indicator="name"` or `data-indicator:_name` | `data-indicator="toastTitle"`, `data-indicator:_fetching` | Signal is true while the request from that element runs. Only the value form with prefix `toast` triggers the saved toast: `toast.js` reads `el.dataset.indicator`, so `data-indicator:toast-title` shows nothing. |
| `data-init="expr"` | `data-init={ live.DatastarInit("/profile/tickets/api") }` | Runs when the element appears, including after a patch that inserts it, and again when a patch changes the expression. A patch with the same expression does not run it. |
| `data-effect="expr"` | `data-effect="$tildelingOpen ? (!el.open && el.showModal()) : (el.open && el.close())"` in `pages/admin/tildelingsdialog.templ` | Re-runs when any signal it reads changes. |
| `data-preserve-attr="open"` | `components/event_components/programpulje_interests.templ` | Morphs keep the listed attributes from the live DOM. |
| `data-ignore-morph` | `InterestContent` in `event_interests.templ`, two spots in `menu.templ` | Region is left alone by parent patches; it needs its own endpoint to update. |

Modifiers seen in the repo: `__ifmissing`, `__prevent`, `__stop`, `__document`, `__window` and `__outside` on `data-on`. Others from the 1.0 docs, such as `__debounce.<ms>`, `__once`, `__throttle.<ms>`, `__delay.<ms>` and `__case.kebab`, work the same way: after the key and `__`, arguments after `.`.

## Actions

`@get(url)`, `@post(url)`, `@put(url)`, `@delete(url)`, `@patch(url)`, with an optional options object. Form posts add `{contentType: 'form'}`, as `pages/profile/newevent/new_page.templ` does, and the cropped image upload in `components/formsubmission/event_img_upload/event_img_upload.templ` also turns retries off. The live `data-init` gets its options from `live.DatastarInit`, which writes:

```js
@get('/profile/tickets/api', {requestCancellation: 'disabled', retryMaxCount: Infinity, retryInterval: 1000, retryMaxWaitMs: 30000})
```

The 1.0.4 client reads `retryMaxWait`, not `retryMaxWaitMs`, so the 30 s cap is Datastar's default, as the lifecycle doc notes. Keep using the helper; do not hand-write this object.

Datastar sends every signal whose name does not start with `_` with each request. GET and DELETE carry them in the `datastar` query parameter, other methods in the JSON body. The request also carries the header `Datastar-Request: true`, which handlers use to choose `sse.Redirect` over a 303. A 303 answer is followed by `fetch`, and Datastar morphs the target page's HTML into the current page by matching ids without changing the URL.

## Go SDK calls in use

| Call | Repo example | Note |
| --- | --- | --- |
| `datastar.ReadSignals(r, &store)` | `components/formsubmission/about_event.templ`, `pages/admin/puljefordeling.go` | Struct with `json:"camelCaseName"` tags. Reads query on GET/DELETE, body otherwise. Call it before `NewSSE`. |
| `sse := datastar.NewSSE(w, r)` | everywhere | Sets `text/event-stream` and flushes headers. Nothing that sets a status may follow, and Go's HTTP/1 server then discards an unread request body. |
| `sse.PatchElementTempl(component, opts...)` | `service/live/live.go`, the feedback form reset in `pages/feedback/feedback.go` | Default: outer morph matched by the top element's id. |
| `datastar.WithModeReplace()`, `WithSelector(sel)`, `WithModeInner()`, `WithModeAppend()`... | `WithModeReplace()` in the selected-interest endpoint in `pages/event/event.go` | Use replace for `data-ignore-morph` regions owned by the endpoint; `WithSelector` when the patch root has no id. |
| `sse.MarshalAndPatchSignals(any)` | `errorfeedback.Patch`, `pages/admin/tildelinger.go` | Marshals for you. Prefer it over hand-built JSON. The billettholder cards still marshal by hand and call `sse.PatchSignals`. |
| `sse.PatchSignals([]byte)` | `pages/profile/tickets/tickets_page.templ` | When you already have JSON. |
| `sse.MarshalAndPatchSignalsIfMissing(any)` | not yet used | Server-side equivalent of `__ifmissing`. |
| `sse.ExecuteScript(js)` | `pages/admin/admin.go` closes the room dialog and dispatches `room-saved` | Appends a `<script>` to `body` through an element patch. |
| `sse.Redirect(url)` | `components/formsubmission/submit_section.templ` | Browser navigation from a Datastar request. |
| `sse.ConsoleError(err)` | `service/live/live.go` | Surfaces a server error in the browser console once headers are sent. |
| `sse.RemoveElement(selector)` | not used | Element patch with mode remove. |

Not every method of `ServerSentEventGenerator` is listed; `go doc github.com/starfederation/datastar-go/datastar ServerSentEventGenerator` prints the rest.

## DOM events from the client

Dispatched on `document`:

- `datastar-fetch` with `evt.detail.type` one of `started`, `finished`, `error`, `retrying`, `retries-failed`, `datastar-patch-elements`, `datastar-patch-signals`, plus `evt.detail.el` (the element that started the request) and `argsRaw`. A non-200 answer gives `started`, `error`, `finished` with the body dropped; a fetch that throws (server down) goes through `retrying` to `retries-failed`.
- `datastar-signal-patch` with the patched signals as `evt.detail`.
- `datastar-ready`, `datastar-prop-change`, `datastar-scope-children`.

`static/js/error_feedback.js` and `static/web_components/toast.js` are built on `datastar-fetch` and `datastar-signal-patch`. There is no before-morph hook; react after a patch with `data-on:datastar-fetch__document` and a check on `evt.detail.type === 'datastar-patch-elements'` and `evt.detail.el?.id`.

## Writing signals from Go

- Values in `data-signals:name` are JavaScript expressions. Use `templ.JSONString(value)` for a string and `strconv.Itoa(n)` for a number, or `components.DataSignals` for several. Go's `%q` and hand-built `'...'` quoting break on apostrophes, backslashes and newlines in user text.
- Signal names with hyphens (`errorMessage-42`) work, because the client reads `$errorMessage-42` as one signal name. The same rule reads `$count-1` as a signal named `count-1`, so put spaces around a minus that follows a signal. New code uses camelCase names in expressions and puts the id in the value or in a nested map.
- Initialize value signals with `data-signals` or `components.DataSignals`, not with `errorfeedback`; that package owns only `feedbackErrors`.
