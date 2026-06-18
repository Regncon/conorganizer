# Live Update Lifecycle

## Human-readable summary

Conorganizer uses Datastar Server-Sent Events and embedded NATS KeyValue buckets to refresh open pages when server-side content changes.

Every live page must render full content during the normal HTTP request. After the page loads, a Datastar `data-init` request opens a live SSE endpoint. The endpoint ensures the browser has a Gorilla session cookie named `connections`, ensures the connection id from that session exists as a key in the relevant NATS KeyValue bucket, sends one full Datastar patch immediately, and then waits for bucket updates.

When a mutation changes content, the server broadcasts to the affected bucket by looping through all keys in that bucket and writing a new timestamp/nonce value to each key. Each open SSE watcher sees its key change and re-renders the full page fragment from the database.

The live update KV data is runtime-only state. NATS does not need to persist these connection keys across process restarts. After a restart, clients reconnect through Datastar, the server recreates the KV key from the existing Gorilla session cookie when possible, and the SSE endpoint sends a full content patch again.

## Decisions

- Keep the existing Gorilla session cookie named `connections`.
- The session value key is `id`; this id is the live connection id.
- Do not introduce a literal `connection` cookie unless a future migration explicitly changes this.
- Live KV values are timestamp/nonce values only.
- Live KV values are not page state, not rendered content, and not an application model.
- Do not use inherited Northstar placeholder-state names in live update code.
- Live KV TTL is `26h`, giving a buffer over the current `24h` Gorilla session max age.
- NATS live connection state is ephemeral and does not need persistence across restarts.
- Scheduled NATS messages are rebuilt from the database on startup. Missed scheduled thresholds during downtime do not need catch-up broadcasts.

## Terminology

- **Connection session**: the Gorilla `connections` session cookie. The `id` value inside this session is the live connection id.
- **Live bucket**: a NATS KeyValue bucket used to notify a class of pages that related content changed.
- **Live key**: a KV key named with the live connection id.
- **Live value**: a timestamp/nonce written to a live key. Its only purpose is to trigger NATS watchers.
- **Live endpoint**: a Datastar SSE route opened from `data-init`.
- **Page renderer**: a server-side function that renders the full live fragment from durable application state, usually SQLite plus request auth context.
- **Broadcast**: an operation that loops all live keys in a bucket and writes a new timestamp/nonce to each key.

## Lifecycle

1. A normal HTTP route renders the full page content.
2. The live wrapper includes a Datastar `data-init` request to the page's live endpoint.
3. The `data-init` request uses retry settings suitable for process restarts:

   ```js
   @get('/some/live/api', {
     requestCancellation: 'disabled',
     retryMaxCount: Infinity,
     retryInterval: 1000,
     retryMaxWaitMs: 30000
   })
   ```

4. The live endpoint ensures the `connections` session before creating the Datastar SSE generator.
5. If the session is missing or expired, the endpoint creates a new session id and saves the session cookie.
6. For each bucket the page subscribes to, the endpoint creates or refreshes the KV key named by the session id.
7. The endpoint creates `datastar.NewSSE(w, r)`.
8. The endpoint immediately sends a full Datastar patch for the live fragment.
9. The endpoint starts NATS watchers for the connection key in each subscribed bucket.
10. When a watched key changes, the endpoint re-renders the full fragment from durable state and patches it through Datastar.
11. When the browser disconnects, the request context is cancelled and the watcher stops.
12. If a NATS watcher closes while the browser request is still active, the endpoint closes the SSE stream so Datastar can reconnect and establish a fresh watcher.

Important ordering rule: create or save the Gorilla session before `datastar.NewSSE(w, r)`. `NewSSE` flushes response headers, so a handler that creates the session after that point may fail to send `Set-Cookie` reliably.

## Broadcast Lifecycle

1. A mutation handler updates durable state, usually SQLite.
2. After the durable update succeeds, the handler calls the live update service to broadcast one or more buckets.
3. The service gets all keys in each bucket.
4. For each key, the service writes a fresh timestamp/nonce value.
5. Open live endpoints watching those keys receive the update and re-render.
6. Missing or expired keys are ignored; the next client request or reconnect recreates them.

At the expected Conorganizer scale, looping all keys on every broadcast is acceptable and preferred for simplicity.

## Buckets

The bucket list should stay small. Pages may subscribe to multiple buckets when they render data from multiple domains. Authorization is enforced by HTTP middleware and render logic, not by bucket names.

| Bucket | Purpose | Typical broadcasters | Typical subscribers |
| --- | --- | --- | --- |
| `events` | Event, program, pulje, publishing, and event-form data. | Event form updates, event submission, approval changes, program publishing, pulje status updates, scheduled pulje threshold broadcasts. | Root page, event details, profile event list, profile event form, admin dashboard, admin approval, admin event edit. |
| `interests` | Interest choices, first-choice data, player/GM assignment state, and views that show who is interested in an event. | User interest updates, admin approval player assignment updates, first-choice and GM changes. | Event details, admin approval, admin billettholder filters, profile views that show interest or assignment state. |
| `billettholders` | Ticket holder and billettholder data. | Add/remove billettholder emails, ticket conversion, ticket fetch/check-in flows, billettholder admin updates. | Profile tickets, profile overview where ticket holders are shown, admin billettholder overview, add billettholder page, possibly event details if ticket holder choices are displayed. |
| `rooms` | Room data and room assignment choices. | Create, update, delete room; assign room to an event pulje. | Admin rooms, event form pages that show room assignment choices, admin event edit. |

### Page Subscription Matrix

This matrix is the starting point for the refactor. Confirm each row while migrating the page.

| Page | Live endpoint | Buckets | Notes |
| --- | --- | --- | --- |
| `/` | `/root/api` | `events` | Must render full root content before Datastar connects. |
| `/event/{id}` | `/event/api/{id}` | `events`, `interests`, possibly `billettholders` | Event detail renders event state, interest state, and user ticket holder choices. |
| `/profile` | `/profile/api` | `events`, `interests`, `billettholders` | Shows user's submitted events and ticket holder context. Include `interests` if assignment or interest state is rendered. |
| `/profile/new/{id}` | `/profile/api/new/{id}` | `events`, `rooms` | Event form renders event fields, pulje choices, and room assignment choices. |
| `/profile/tickets` | `/profile/tickets/api` | `billettholders` | Ticket holder profile data. |
| `/admin` | `/admin/api` | `events` | Admin dashboard currently focuses on event/program controls. |
| `/admin/approval` | `/admin/approval/api` | `events`, `interests`, `billettholders` | Approval views render event data, interest data, assignment state, and interested ticket holders. |
| `/admin/approval/edit/{id}` | `/admin/approval/edit/api/{id}` | `events`, `rooms` | Admin event edit form. |
| `/admin/rooms` | To be added when live updates are introduced | `rooms` | Currently not standardized with the live page lifecycle. |
| `/admin/billettholder` | `/admin/billettholder/api` | `billettholders`, `interests`, possibly `events` | Some filters depend on first-choice/event interest data. Confirm whether event metadata is also rendered during migration. |
| `/admin/billettholder/add` | `/admin/billettholder/add/api` | `billettholders` | Add/convert ticket holder workflows. |
| `/login` | None | None | No live updates. |
| Print-friendly pages | None | None | Static render only. |

## Targeted Updates

The current live value is intentionally just a timestamp/nonce. That does not permanently limit the architecture to bucket-wide broadcasts, but the current service should avoid new per-user or per-session key namespaces until a concrete feature needs them.

The current key shape is:

```text
<connection-id>
```

That shape means a bucket broadcast updates every active connection subscribed to that bucket.

For transient per-connection UI interactions, such as search/filter endpoints that should refresh only the current open page, reuse the existing connection key instead of adding a new key namespace. The handler may call `EnsureConnection(w, r, bucket)` for the relevant bucket. Because `EnsureConnection` writes a fresh timestamp/nonce for the current `connections` session id, only the current connection key is touched.

If a future feature needs targeted updates, add it deliberately with tests and a clear use case. Two reasonable extensions are:

```text
connection.<connection-id>
user.<user-id>.<connection-id>
```

A targeted per-user broadcast would list or watch keys matching `user.<user-id>.*` in the relevant bucket and write a fresh timestamp/nonce only to those keys. A targeted per-session broadcast would write only `connection.<connection-id>`.

Do not add these key shapes before a real feature needs them. The global connection-id key is easier to reason about and is sufficient for the current app.

## Future JSON KV Values

The current implementation should use a plain timestamp/nonce value because it is enough to wake NATS watchers and keeps the broadcast path simple.

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

## Target Service Shape

The exact API can change during implementation, but all live pages should use one shared service instead of open-coded NATS/session logic.

```go
type Bucket string

const (
	BucketEvents         Bucket = "events"
	BucketInterests      Bucket = "interests"
	BucketBillettholders Bucket = "billettholders"
	BucketRooms          Bucket = "rooms"
)

type Manager struct {
	// Owns NATS, JetStream, bucket creation, session handling, and broadcast helpers.
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

Expected usage in page setup:

```go
router.Get("/root/api", func(w http.ResponseWriter, r *http.Request) {
	liveManager.Stream(w, r, live.Page{
		Buckets: []live.Bucket{live.BucketEvents},
		Render: func(ctx context.Context, r *http.Request) templ.Component {
			return rootPage(db, isAdmin, eventImageDir)
		},
	})
})
```

Expected usage after mutations:

```go
if err := liveManager.Broadcast(r.Context(), live.BucketEvents); err != nil {
	logger.Error("failed to broadcast live update", "error", err)
	http.Error(w, "Failed to broadcast update", http.StatusInternalServerError)
	return
}
```

Expected usage for current-connection UI refreshes:

```go
if _, err := liveManager.EnsureConnection(w, r, live.BucketBillettholders); err != nil {
	http.Error(w, err.Error(), http.StatusInternalServerError)
	return
}
```

## Testing Strategy

Implement the live update service with behavior-focused tests before migrating pages.

Tests should follow the repository's Given/When/Then structure:

```go
func TestManager_EnsureConnection_WhenCookieExistsAndKeyExpired_RecreatesLiveKey(t *testing.T) {
	// Given an existing live session without a matching KV key,
	// when the manager ensures the connection,
	// then the same connection id is reused and the KV key is recreated.

	// Given
	expectedConnectionID := "..."

	// When

	// Then
}
```

Recommended service tests:

- `EnsureConnection` creates a `connections` session cookie and live KV key when missing.
- `EnsureConnection` reuses an existing `connections` session id.
- `EnsureConnection` recreates a missing KV key for an existing session id.
- The live KV bucket TTL is `26h`.
- `Broadcast` writes a new timestamp/nonce to every key in the bucket.
- `Broadcast` succeeds when a bucket has no keys.
- A watcher receives an update after `Broadcast`.
- `Stream` sends an initial full Datastar patch before waiting for broadcasts.

Use embedded NATS in service integration tests. Keep page rendering, database setup, and auth out of the initial service tests by using a tiny test renderer.

## Scheduled NATS Messages

Scheduled pulje broadcasts are not live connection state. They are derived jobs built from database pulje start times.

Current behavior to preserve:

- Create a scheduled stream for pulje warning thresholds.
- Rebuild future schedules from the database on startup.
- Ignore thresholds that are already in the past when startup runs.
- Do not replay missed threshold broadcasts after downtime.
- When a scheduled threshold fires while the app is running, broadcast the `events` bucket.

This is sufficient because clients reconnect after a restart and receive a full page patch from the live endpoint.

## LLM Implementation Contract

This section is intentionally explicit for AI coding agents.

When implementing or modifying live update code:

- Do not introduce inherited Northstar placeholder-state or Todo terminology into live update code.
- Treat SQLite and request context as the source of truth for rendered content.
- Do not store rendered page content in NATS KV.
- Do not store form state in NATS KV for this lifecycle.
- Store only a timestamp/nonce as the live KV value.
- If JSON KV values are introduced later, keep them small, schema-versioned, and metadata-only.
- Use global connection-id keys for now. Do not implement per-user or per-session key namespaces until a concrete feature needs targeted updates.
- Use the Gorilla `connections` session cookie and the session value key `id`.
- Ensure the session and KV key before calling `datastar.NewSSE(w, r)`.
- Every live SSE stream must send one full patch immediately after opening.
- Use Datastar retry settings that survive server restarts.
- Broadcast by looping every key in the target bucket and writing a fresh timestamp/nonce.
- Set live KV bucket TTL to `26h`.
- Keep live NATS connection state ephemeral; do not add persistent NATS storage for live update buckets.
- Scheduled NATS messages may use JetStream scheduling, but schedules must be rebuildable from durable database state.
- Keep bucket definitions centralized.
- Prefer broad buckets over premature fine-grained splitting unless the page/bucket matrix shows a real correctness issue.
- Pages that render data from multiple domains should subscribe to multiple buckets.
- Security belongs in HTTP middleware and render logic, not in bucket names.
- New tests must use behavior-focused names and Given/When/Then sections.
