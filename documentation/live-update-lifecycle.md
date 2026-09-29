# Live Update Lifecycle

## Human-readable summary

Conorganizer uses Datastar Server-Sent Events and embedded NATS KeyValue buckets to refresh open pages when server-side content changes.

For browser state, signal responses, and signals carried by live HTML, see
[Datastar signals and examples](datastar-signals.md).

Every live page must render full content during the normal HTTP request (see [Full content first](#full-content-first) for the known exceptions). After the page loads, a Datastar `data-init` request opens a live SSE endpoint. The endpoint ensures the browser has a Gorilla session cookie named `connections`, sends one full Datastar patch immediately, then writes the connection id from that session as a key in each relevant NATS KeyValue bucket and waits for updates to that key.

When a mutation changes content, the server broadcasts to the affected bucket by looping through all keys in that bucket and writing a new timestamp value to each key. Each open SSE watcher sees its key change and re-renders the full page fragment from the database.

The live update KV data is runtime-only state. NATS does not need to persist these connection keys across process restarts. After a restart, clients reconnect through Datastar, the server recreates the KV key from the existing Gorilla session cookie when possible, and the SSE endpoint sends a full content patch again.

## Authentication before streaming

Application startup creates one Descope session validator shared by page requests, mutations, live endpoints, login session establishment, and the not-found handler. The SDK retrieves signing keys automatically when needed and retains its key cache and HTTP connection pools for the application lifetime. Requests without authentication cookies do not call session validation or refresh.

`POST /auth/session` and `GET /auth/logout` are recovery routes registered outside the authentication middleware. Session establishment validates only the tokens submitted in the request body, so stale cookies cannot trigger an upstream refresh before a new login completes. Logout clears local authentication cookies without contacting Descope. The remaining authentication routes use the application middleware's authentication result rather than running that middleware again.

Authentication finishes before a live handler opens its SSE response, so a refreshed session cookie is sent before the stream headers are flushed. No authentication deadline is attached to the downstream request context.

Authentication timing logs use `component=auth`, `operation=session_validation` or `operation=session_refresh`, `duration_ms`, `succeeded`, and `request_id`. Normal timings are debug-level; operations taking at least one second produce a warning, including slow operations that eventually succeed. Timing logs do not include tokens or cookies. `session_validation` includes signing-key retrieval on a cache miss; it is not necessarily local-only work.

Client reuse mitigates repeated key retrieval but does not remove the upstream dependency for refresh or an uncached signing key. The Descope Go SDK (v1.33.0 in `go.mod`) uses a background context for key retrieval, so a caller context timeout alone would not bound the complete authentication operation. Authentication failure responses and Datastar retry behavior are unchanged by the shared validator.

## Decisions

- Keep the existing Gorilla session cookie named `connections`.
- The session value key is `id`; this id is the live connection id.
- Do not introduce a literal `connection` cookie unless a future migration explicitly changes this.
- Live KV values are timestamp values only (`time.RFC3339Nano`, UTC).
- Live KV values are not page state, not rendered content, and not an application model.
- Do not use inherited Northstar placeholder-state names in live update code.
- Live KV TTL is `26h` (`live.DefaultTTL`), giving a buffer over the current `24h` Gorilla session max age.
- NATS live connection state is ephemeral and does not need persistence across restarts.
- Conorganizer has roughly 200 users in total, so the number of concurrent live connections stays low. The simple broadcast loop and the absence of per-connection heartbeats rely on this.

## Terminology

- **Connection session**: the Gorilla `connections` session cookie. The `id` value inside this session is the live connection id.
- **Live bucket**: a NATS KeyValue bucket used to notify a class of pages that related content changed.
- **Live key**: a KV key named with the live connection id.
- **Live value**: a timestamp written to a live key. Its only purpose is to trigger NATS watchers.
- **Live endpoint**: a Datastar SSE route opened from `data-init`.
- **Page renderer**: a server-side function that renders the full live fragment from durable application state, usually SQLite plus request auth context.
- **Broadcast**: an operation that loops all live keys in a bucket and writes a new timestamp to each key.

## Naming and public API

The live update code lives in `service/live`. It has no MVC or Todo naming from the Northstar template. The public API is:

- `NewManager`, `WithTTL`, `WithLogger`
- `Manager.EnsureConnection`, `Manager.Stream`, `Manager.Broadcast`
- `DatastarInit`, `DatastarInitExpression`
- `Page`, `Bucket` and the four `Bucket*` constants, `DefaultTTL`

Internal helpers include `ensureSession`, `touchConnection` and `forwardWatcherUpdates`.

## Lifecycle

1. A normal HTTP route renders the full page content.
2. The live wrapper includes a Datastar `data-init` request to the page's live endpoint, built with `live.DatastarInit` or `live.DatastarInitExpression` (see [Datastar init and retries](#datastar-init-and-retries)).
3. `Manager.Stream` ensures the `connections` session before creating the Datastar SSE generator.
4. If the session is missing or expired, it creates a new session id (a UUID) and saves the session cookie.
5. It creates `datastar.NewSSE(w, r)`.
6. It immediately sends a full Datastar patch for the live fragment. If this fails, it logs the error, sends it to the browser console, and ends the stream.
7. For each bucket the page subscribes to, it writes (touches) the KV key named by the connection id. If this fails, it logs the error, sends it to the browser console, and ends the stream. The user still has the content from step 6.
8. It starts one NATS watcher per subscribed bucket for the connection key, using `jetstream.UpdatesOnly()`.
9. When a watched key changes, it re-renders the full fragment from durable state and patches it through Datastar. Several updates arriving close together are coalesced into one re-render.
10. When the browser disconnects, the request context is cancelled, the stream returns, and the watchers are stopped.
11. If a NATS watcher closes while the browser request is still active, the endpoint logs a warning and closes the SSE stream instead of staying open without updates (see [Troubleshooting](#troubleshooting)).

Important ordering rule: only the Gorilla session must be created or saved before `datastar.NewSSE(w, r)`. `NewSSE` flushes response headers, so a handler that creates the session after that point may fail to send `Set-Cookie` reliably.

The initial patch deliberately comes before the KV touch and watcher setup. If NATS is slow or unavailable, the user still gets current content instead of a blank container or a failed request (covered by `TestManager_Stream_WhenTouchConnectionFails_SendsInitialPatch`). The trade-off is that a broadcast in the short window between the initial patch and watcher start is missed. This is accepted; do not change the order.

`EnsureConnection`, used for per-connection interactions such as search requests, ensures the session and then touches the key immediately.

## Broadcast Lifecycle

1. A mutation handler updates durable state, usually SQLite.
2. After the durable update succeeds, the handler calls `liveManager.Broadcast(r.Context(), buckets...)`.
3. The service lists all keys in each bucket with `kv.Keys()`.
4. For each key, the service writes a fresh timestamp value with `kv.Put`.
5. Open live endpoints watching those keys receive the update and re-render.
6. Missing or expired keys are not listed, so they get no update; the next stream or `EnsureConnection` call recreates them. An empty bucket is not an error.

At the expected Conorganizer scale (about 200 users), looping all keys on every broadcast is acceptable and preferred for simplicity over a more scalable pub/sub design.

## Buckets

The bucket list should stay small. `service/live` defines exactly four buckets. Pages may subscribe to multiple buckets when they render data from multiple domains. Authorization is enforced by HTTP middleware and render logic, not by bucket names.

| Bucket | Purpose | Typical broadcasters | Typical subscribers |
| --- | --- | --- | --- |
| `events` | Event, program, pulje, publishing, and event-form data. | Event form updates, event creation and submission, approval status changes, program publishing, puljeoppsett changes, pulje status and closing-warning changes, puljefordeling seat changes. | Root page, event details, profile, profile event form, admin dashboard, admin approval, admin event edit, room assignment, puljeoppsett, puljefordeling. |
| `interests` | Interest choices, first-choice data, player/GM assignment state, and views that show who is interested in an event. | User interest updates on the event page, admin interest changes in puljefordeling, tildelinger. | Event details, profile, admin approval, puljefordeling, admin billettholder overview. |
| `billettholders` | Billettholder data and billettholder emails. | Add/remove billettholder emails (admin and profile), ticket conversion, ticket fetch on the profile tickets page. | Header menu, profile, profile tickets, admin approval, admin billettholder overview, add billettholder page. |
| `rooms` | Room data and room assignment choices. | Create, update, delete room; assign or remove a room for an event in a pulje. | Event details, admin rooms, room assignment, profile event form, admin event edit. |

Dev hot reload does not use a live bucket. It is a separate in-process hub in `dev_reload.go` (build tag `dev`), with a `/reload` SSE endpoint and a `/hotreload` trigger.

### Cross-bucket broadcasts

Room assignment is cross-cutting: it affects event rendering and room-dependent forms, so these handlers broadcast more than one bucket:

- Admin room assignment (`POST /admin/rooms/api/assignment/{pulje}/{event}/{room}`) and removal (`DELETE` on the same path) broadcast `rooms` and `events`.
- Assigning a room from the event form (`PUT /profile/api/new/{id}/assign-room/{puljeId}`) broadcasts `events` and `rooms`.
- Saving a puljefordeling (`POST /admin/api/puljefordeling/{pulje}/commit`) broadcasts `events` and `rooms`.
- Saving tildelinger and removing a GM in puljefordeling broadcast `events`, `interests` and `rooms`.
- Creating, updating or deleting a room broadcasts only `rooms`.

A pulje status or closing-warning change broadcasts `events`, so open event pages (`/event/api/{id}` watches `events`) re-render with the new pulje state.

### Page subscriptions

| Page | Live endpoint | Buckets | Notes |
| --- | --- | --- | --- |
| `/` | `/root/api` | `events` | Keeps the selected program day through the `date` query parameter. |
| `/event/{id}` | `/event/api/{id}` | `events`, `interests`, `rooms` | The page's query string is passed to the live endpoint. |
| `/profile` | `/profile/api` | `events`, `interests`, `billettholders` | Uses `DatastarInitExpression` to pass `window.location.search`. Patches only `#profile-main-column`. |
| `/profile/new/{id}` | `/profile/api/new/{id}` | `events`, `rooms` | Event form with pulje and room choices. |
| `/profile/tickets` | `/profile/tickets/api` | `billettholders` | |
| Header menu on every page for logged-in users | `/menu/api` | `billettholders` | Billettholder selector (`#main-menu-billettholder-live`). |
| `/admin/` | `/admin/api/` | `events` | Renders a placeholder first; see [Full content first](#full-content-first). |
| `/admin/approval/` | `/admin/approval/api/` | `events`, `interests`, `billettholders` | |
| `/admin/approval/edit/{id}` | `/admin/approval/edit/api/{id}` | `events`, `rooms` | Admin event edit form. |
| `/admin/rooms/` | `/admin/rooms/api/` | `rooms` | |
| `/admin/rooms/assignment/{pulje}` | `/admin/rooms/api/assignment/{pulje}` | `rooms`, `events` | Patches only `#room-assignment`; the dialog is outside the live fragment. |
| `/admin/puljeoppsett/` | `/admin/api/puljeoppsett` | `events` | |
| `/admin/puljefordeling/{pulje}` | `/admin/api/puljefordeling/{pulje}` | `events`, `interests` | |
| `/admin/billettholder/` | `/admin/billettholder/api/` | `billettholders`, `interests` | Renders a placeholder first. Search uses `EnsureConnection` with `billettholders`. |
| `/admin/billettholder/add/` | `/admin/billettholder/add/api/` | `billettholders` | Renders a placeholder first. Ticket search uses `EnsureConnection` with `billettholders`. |
| `/auth` | None | None | No live updates. |
| `/print` | None | None | Static render only. |

## Full content first

A normal GET must render real content in the first HTTP response. The SSE live patch later replaces the same container. If Datastar starts late or the SSE connection fails, the user still sees current content.

Known exceptions still in the code:

- `/admin/` (`adminIndex` in `pages/admin/admin_index.templ`)
- `/admin/billettholder/` (`pages/admin/billettholder_admin/billettholder_admin_index.templ`)
- `/admin/billettholder/add/` (`pages/admin/billettholder_admin/add/add_billettholder_index.templ`)

These render only a placeholder container with `data-init` ("If you are seeing this message, please clear your cookies and refresh the page.") and rely fully on the first SSE patch. Any delay in Datastar startup or the first patch shows the placeholder, and a failed SSE connection leaves the page without content. They should be migrated to render full content first.

## Datastar init and retries

Every live `data-init` uses `live.DatastarInit(url)` or `live.DatastarInitExpression(expr)`. Both produce:

```js
@get('/some/live/api', {
  requestCancellation: 'disabled',
  retryMaxCount: Infinity,
  retryInterval: 1000,
  retryMaxWaitMs: 30000
})
```

`DatastarInit` quotes a static URL and escapes single quotes. Use `DatastarInitExpression` when the URL is a JavaScript expression, for example the profile page:

```templ
data-init={ live.DatastarInitExpression("'/profile/api' + window.location.search") }
```

The query string carries the selected billettholder to the live stream. `/profile/api` re-renders `ProfileMainColumn` (`#profile-main-column`, with `MyEvents` and `MyProgram`), and Datastar morphs it over the same element. The account administration column is not part of the live fragment.

Retry behavior of the bundled client (`static/datastar.js`, Datastar v1.0.4) with these options:

- Datastar retries when the fetch throws, for example when the server process is down or the connection breaks mid-stream. The wait starts at `retryInterval` (1000 ms) and doubles up to 30 s. `retryMaxCount: Infinity` keeps retrying, and the reconnect gets a full patch.
- The default `retry: 'auto'` does not retry after a non-2xx response (for example a proxy `502` while the service restarts) or after the server ends the stream cleanly.
- Datastar v1.0.4 reads `retryMaxWait`, not `retryMaxWaitMs`. The 30 s cap therefore comes from Datastar's default.
- GET requests are aborted while the tab is hidden and reopened when it becomes visible, because `openWhenHidden` defaults to `false` for GET. Each reopen is a new `Stream` call with a new initial patch and key touch.

## Shared session keys across tabs

Live keys are named by the `connections` session id, and that cookie is shared by all tabs and windows in the same browser session. So every open live page in that browser watches the same key in each bucket it subscribes to.

`Stream` and `EnsureConnection` touch that key in their buckets when a live page opens, when a tab becomes visible again, or when a search request runs. Any other tab watching the same key in an overlapping bucket then re-renders even though no data changed. For example, `/admin/billettholder/` watches `billettholders` and `interests`. Opening `/event/{id}` (`events`, `interests`, `rooms`) in another tab re-renders it through `interests`. Because the header menu streams `billettholders` on every page for logged-in users, opening any page touches the `billettholders` key and re-renders every tab that watches `billettholders`.

These extra re-renders are harmless when renderers are idempotent and browser-owned state survives morphs (see [Datastar signals](datastar-signals.md#keep-browser-owned-state-through-morphs)).

## Native dialogs and live fragments

Keep a native modal `<dialog>` that is opened imperatively with `showModal()` outside the live-patched fragment. If a live refresh morphs or replaces an open modal dialog, the modal and inert state can stay attached to a detached element and leave the whole document inert with no visible modal. A tab switch reopens the stream, so it can trigger this too.

- Room assignment: a stable `#room-assignment-page` wrapper owns the shared signals (`room`, `_roomNumber`, `_roomSaving`, `_roomAssignmentStatus`, `draggedEventId`, `dragOverRoom`) and the `room-saved` window handler. It contains the live `#room-assignment` section (with the `data-init`) and `<dialog id="assignment-dialog">`. The SSE endpoint renders only `RoomsAssignmentPageContent`, which is the `#room-assignment` section.
- Puljefordeling: the SSE endpoint renders `puljefordelingPage` (`#puljefordeling-page`), and `#tildeling-dialog` is rendered outside it in `puljefordelingIndex`.

A dialog that must stay inside a live fragment must be opened from a signal through `data-effect` and keep `data-preserve-attr="open"`. Examples are the puljefordeling assign dialog, the billettholder interest dialog and the event room map dialog. See [Datastar signals](datastar-signals.md#native-modal-dialogs-in-live-fragments).

## NATS storage and TTL

The embedded NATS server is started in `setupRoutes` (`router.go`) with JetStream and an explicit `StoreDir` from the `-nats-store-dir` flag. The flag defaults to `data/nats` for local development and is passed through `main` -> `run` -> `startServer` -> `setupRoutes` into `natsserver.Options.StoreDir`. `setupRoutes` returns an error if it is empty.

Every app instance starts its own embedded NATS server, so each instance needs its own store directory. Each systemd unit sets `RuntimeDirectory=conorganizer-<safe name>` and passes `-nats-store-dir /run/conorganizer-<safe name>/nats`. This applies to main, demo, restored and PR previews. The CI deploy fills `%SAFE_NAME%` into `deploy/conorganizer.service.tmpl` with `sed`. Instances therefore never share NATS state, and systemd removes the directory when the service stops, so it does not survive a restart. The path is not configured through an environment variable, and it is not derived from the random NATS port or the working directory.

The live buckets use `jetstream.MemoryStorage` with a 16 MiB `MaxBytes` limit. Their contents are never persisted: live keys are rebuildable runtime state. The app has no other JetStream stream or consumer. JetStream still needs a `StoreDir` for its own runtime bookkeeping, so the per-instance directory holds only NATS server metadata, not app data. After a restart, Datastar retries reconnect, `Stream` recreates the keys from the `connections` session, and the initial full patch renders the latest database state.

The bucket TTL is `live.DefaultTTL` (`26h`), longer than the `24h` Gorilla session `MaxAge` set in `router.go`. `Broadcast` only writes to keys returned by `kv.Keys()`, so an expired key is skipped. A connection whose key expired before its session would silently miss broadcasts. A shorter TTL than the cookie lifetime caused exactly that for idle open pages.

## Targeted Updates

The current live value is intentionally just a timestamp. That does not permanently limit the architecture to bucket-wide broadcasts, but the current service should avoid new per-user or per-session key namespaces until a concrete feature needs them.

The current key shape is:

```text
<connection-id>
```

That shape means a bucket broadcast updates every active connection subscribed to that bucket.

For transient per-connection UI interactions, such as search/filter endpoints that should refresh only the current open page, reuse the existing connection key instead of adding a new key namespace. The handler may call `EnsureConnection(w, r, bucket)` for the relevant bucket. Because `EnsureConnection` writes a fresh timestamp for the current `connections` session id, only the current connection key is touched. Other tabs in the same browser share that key (see [Shared session keys across tabs](#shared-session-keys-across-tabs)).

If a future feature needs targeted updates, add it deliberately with tests and a clear use case. Two reasonable extensions are:

```text
connection.<connection-id>
user.<user-id>.<connection-id>
```

A targeted per-user broadcast would list or watch keys matching `user.<user-id>.*` in the relevant bucket and write a fresh timestamp only to those keys. A targeted per-session broadcast would write only `connection.<connection-id>`.

Do not add these key shapes before a real feature needs them. The global connection-id key is easier to reason about and is sufficient for the current app.

## Future JSON KV Values

The current implementation uses a plain timestamp value because it is enough to wake NATS watchers and keeps the broadcast path simple.

JSON values are acceptable in the future if a concrete feature needs structured metadata, but they should not be introduced speculatively. JSON values may be useful for debugging, observability, schema evolution, or preserving small connection metadata alongside the nonce.

Example future value shape:

```json
{
  "version": 1,
  "nonce": "01JZ...",
  "updated_at": "2026-06-07T12:00:00Z",
  "page": "admin.approval",
  "user_id": "auth-provider-user-id"
}
```

Rules for future JSON values:

- Keep the value small.
- Include a `version` field.
- Include a fresh `nonce` or `updated_at` on every broadcast so watchers receive an update.
- Do not store rendered HTML.
- Do not store form state.
- Do not store secrets or sensitive personal data.
- Do not treat the KV value as durable application state.
- Use typed Go structs and `encoding/json`, not string concatenation.
- Handle corrupt or old JSON values defensively; a bad value should not break broadcasts for other connections.
- Add tests for create, read, broadcast update, old-version handling, and corrupt-value handling before adopting JSON values.

JSON values should not be used for recipient selection unless there is a strong reason. NATS can list and watch keys efficiently, but it cannot query JSON value contents. If a feature needs targeted broadcasts, prefer key namespaces such as `user.<user-id>.<connection-id>` over scanning every value and filtering decoded JSON.

## Service API

All live pages use the shared service in `service/live` instead of open-coded NATS/session logic.

```go
type Bucket string

const (
	BucketEvents         Bucket = "events"
	BucketInterests      Bucket = "interests"
	BucketBillettholders Bucket = "billettholders"
	BucketRooms          Bucket = "rooms"
)

type Manager struct {
	// Owns the NATS connection, bucket handles, session handling, and broadcast helpers.
}

type Page struct {
	Buckets []Bucket
	Render  func(ctx context.Context, r *http.Request) templ.Component
}

func NewManager(ctx context.Context, ns *embeddednats.Server, store sessions.Store, opts ...Option) (*Manager, error)
func (m *Manager) EnsureConnection(w http.ResponseWriter, r *http.Request, buckets ...Bucket) (string, error)
func (m *Manager) Stream(w http.ResponseWriter, r *http.Request, page Page)
func (m *Manager) Broadcast(ctx context.Context, buckets ...Bucket) error
func DatastarInit(url string) string
func DatastarInitExpression(urlExpression string) string
```

`NewManager` creates or updates all four buckets with `CreateOrUpdateKeyValue`. `router.go` creates one manager with `live.WithLogger(logger)` and passes it to the page setup functions.

Usage in page setup:

```go
router.Get("/root/api", func(w http.ResponseWriter, r *http.Request) {
	requestedDate := r.URL.Query().Get(programDateQueryParam)
	liveManager.Stream(w, r, live.Page{
		Buckets: []live.Bucket{live.BucketEvents},
		Render: func(ctx context.Context, r *http.Request) templ.Component {
			return rootPage(db, eventImageDir, requestedDate)
		},
	})
})
```

Usage after mutations:

```go
if err := liveManager.Broadcast(r.Context(), live.BucketEvents); err != nil {
	logger.Error("failed to broadcast live update", "error", err)
	http.Error(w, "Failed to broadcast update", http.StatusInternalServerError)
	return
}
```

Some handlers, such as the puljefordeling handlers, treat the broadcast as best effort: they log the error and still return success.

Usage for current-connection UI refreshes:

```go
if _, err := liveManager.EnsureConnection(w, r, live.BucketBillettholders); err != nil {
	http.Error(w, err.Error(), http.StatusInternalServerError)
	return
}
```

## Logging

Live logs come from the manager with `component=live` and carry `method`, `path`, `request_id` and `buckets`. Error messages have the form `<message>: <error>`. The live connection id (session-derived key) is never logged. Broadcast errors are returned to the caller and logged by the mutation handler.

| Level | Message | Where |
| --- | --- | --- |
| Error | `failed to ensure live session` | `Stream`, `EnsureConnection` |
| Error | `failed to send initial live patch` | `Stream` |
| Error | `failed to touch live key before watching` | `Stream` |
| Error | `failed to touch live key` | `EnsureConnection` |
| Error | `failed to prepare live watcher`, `failed to start live watcher` | `Stream` |
| Error | `failed to send live patch after update` | `Stream` |
| Warn | `live watcher closed; closing stream for reconnect` | `Stream`, when a watcher closes while the request is active |
| Warn | `failed to stop live watcher` | `Stream` cleanup |
| Info | `live watcher already stopped` | `Stream` cleanup, when `Stop()` returns `nats.ErrBadSubscription` |

`live watcher already stopped` is expected cleanup after navigation, reconnect, client disconnect or shutdown. It is not a fault.

Touch failures also include `nats_touch_duration_ms`. When a NATS connection exists, touch failures and the watcher-closed warning include `nats_status` and, if set, `nats_last_error`. These fields tell a slow JetStream write apart from a disconnected or stale NATS client.

## Troubleshooting

### Open pages stop refreshing

Symptom: a mutation succeeds and is saved, but open pages do not refresh. It looks random, and a service restart fixes it.

Known root cause: the JetStream `KeyWatcher.Updates()` channel can close while the browser's SSE request is still open. `forwardWatcherUpdates` now reports the closed channel to `Stream`, which logs `live watcher closed; closing stream for reconnect` and returns, closing the SSE response. If the request context is already cancelled, it returns without the warning. The regression test is `TestManager_Stream_WhenWatcherCloses_ExitsStreamForReconnect`.

Ruled out: cancellation of the mutating request, KV TTL or persistence (the TTL is 26h, and every timestamp write creates a new revision, so updates are not suppressed), and the race between the initial patch and watcher start (the bug also happened long after the page was idle).

No SSE heartbeat or keepalive is implemented. Recovery relies on closing the stream so Datastar reconnects with a fresh watcher. Note the retry rules in [Datastar init and retries](#datastar-init-and-retries): Datastar v1.0.4 with `retry: 'auto'` does not retry after a clean stream end, so a page whose stream was closed this way may need a reload or tab switch to reconnect.

If missed reloads come back, search the logs for `live watcher closed; closing stream for reconnect`.

### `failed to touch live key ... context deadline exceeded`

A value of about 5000 in `nats_touch_duration_ms` means the JetStream call hit the nats.go jetstream client's default API timeout (`defaultAPITimeout = 5 * time.Second`). The client applies that timeout when the context has no deadline. Conorganizer creates its JetStream context with a plain `jetstream.New(nc)`, so the default is not changed with `jetstream.WithDefaultTimeout`, and request contexts have no deadline. `kv.Put` is a synchronous JetStream publish that waits for the server's ack, so when JetStream is stuck the touch fails after about 5 s.

Known cause: several instances sharing one JetStream store directory. Before `-nats-store-dir` existed, main, demo and every PR preview on the host used the default `/tmp/nats/jetstream` and corrupted each other's files. The errors were intermittent, with `nats_touch_duration_ms` around 5000 and `nats_status` `CONNECTED`. `sudo fuser -vm /tmp/nats/jetstream` showed several services holding files. Check that each unit passes its own `-nats-store-dir` (see [NATS storage and TTL](#nats-storage-and-ttl)).

## Testing Strategy

Tests follow the repository's Given/When/Then structure:

```go
func TestManager_EnsureConnection_WhenCookieExistsAndKeyMissing_RecreatesLiveKey(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given an existing live session without a matching KV key.",
		When:  "When the manager ensures the connection.",
		Then:  "Then the same connection id is reused and the KV key is recreated.",
	})

	// Given
	expectedConnectionID := "..."

	// When

	// Then
}
```

Current service tests in `service/live/live_test.go` cover:

- `EnsureConnection` creates a `connections` session cookie and live KV key when missing.
- `EnsureConnection` reuses an existing `connections` session id and recreates a missing KV key.
- Every live bucket uses memory storage and the `26h` TTL.
- `Broadcast` writes a new timestamp to every key in the bucket.
- `Broadcast` succeeds when a bucket has no keys.
- A watcher receives an update after `Broadcast`.
- `Stream` sends the initial patch when touching the live key fails, and logs the touch error with `nats_touch_duration_ms`.
- `Stream` logs an already stopped watcher at Info.
- `Stream` returns when a watcher closes while the request is active.
- `DatastarInit` and `DatastarInitExpression` include the retry settings.

The service tests do not use a real NATS server. `newTestManager` wires every bucket to an in-memory `fakeKeyValue` (`Get`, `Put`, `Purge`, `Keys`, `Watch`) with a `fakeWatcher`, and has hooks such as `onWatch`, `onWatchWatcher`, `watcherStopErr` and `putErr`. This makes lifecycle edge cases deterministic, such as a watcher closing, a stop error or a failed touch. The production `keyValue` interface only needs `Put`, `Keys` and `Watch`.

Known gap: fake-based tests do not exercise real JetStream KV `Watch()`/`Keys()` semantics, and the dead-watcher bug lived in those semantics. Real embedded NATS (`embeddednats` plus `live.NewManager`) is used only in route-level tests such as `pages/admin/rooms_assignment_route_test.go` and `pages/event/event_rooms_browser_test.go`. A service-level embedded-NATS integration test for the watch/broadcast pattern would still be valuable.

## LLM Implementation Contract

This section is intentionally explicit for AI coding agents.

When implementing or modifying live update code:

- Do not introduce inherited Northstar placeholder-state or Todo terminology into live update code.
- Treat SQLite and request context as the source of truth for rendered content.
- Do not store rendered page content in NATS KV.
- Do not store form state in NATS KV for this lifecycle.
- Store only a timestamp as the live KV value.
- If JSON KV values are introduced later, keep them small, schema-versioned, and metadata-only.
- Use global connection-id keys for now. Do not implement per-user or per-session key namespaces until a concrete feature needs targeted updates.
- Use the Gorilla `connections` session cookie and the session value key `id`.
- Ensure the session before calling `datastar.NewSSE(w, r)`. Touch the KV key after the initial patch, not before.
- Every live SSE stream must send one full patch immediately after opening, before touching KV keys or starting watchers.
- Every normal GET of a live page must render full content before Datastar connects.
- Use `live.DatastarInit` or `live.DatastarInitExpression` for live `data-init` so the retry settings stay consistent.
- Broadcast by looping every key in the target bucket and writing a fresh timestamp.
- Set live KV bucket TTL to `26h`.
- Keep live NATS connection state ephemeral; the live buckets use `jetstream.MemoryStorage`. Do not add persistent NATS storage for live update buckets.
- Each app instance must use its own `-nats-store-dir`.
- Keep bucket definitions centralized in `service/live`.
- Prefer broad buckets over premature fine-grained splitting unless the page subscriptions show a real correctness issue.
- Pages that render data from multiple domains should subscribe to multiple buckets.
- Broadcast every bucket whose subscribers render the changed data, for example `rooms` and `events` for room assignment.
- Keep native modal dialogs opened with `showModal()` outside live fragments, or drive them from a signal with `data-preserve-attr="open"`.
- Security belongs in HTTP middleware and render logic, not in bucket names.
- New tests must use behavior-focused names and Given/When/Then sections.
