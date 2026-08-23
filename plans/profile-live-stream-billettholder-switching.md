# Profile Live-Stream Billettholder Switching Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Update the profile page for a newly selected billettholder without reloading the document, while ensuring only the selected billettholder's live stream can patch the profile.

**Architecture:** Keep the menu-owned `$billettHolderId` signal as the only live selection input, using numeric `0` when no billettholder exists. A small profile Datastar effect mirrors that ID into the visible `b_id`, cancels the old SSE request, and opens a new `/profile/api` stream without reloading the document. The backend validates the automatically sent signal, while the existing `service/live.Manager` and NATS buckets continue updating the replacement stream unchanged.

**Tech Stack:** Go, templ, Datastar signals/actions, server-sent events, and the existing `service/live.Manager`.

**Spec:** `plans/profile-live-stream-billettholder-switching.md#requirements`

## Global Constraints

- Do not add a profile endpoint to the global menu.
- Do not persist the selected billettholder in a second client-side store; continue using `window.conorganizer.billettholderSelection` and `$billettHolderId`.
- Do not add server-side session or connection state for the selection; local storage persists it in the browser and the Datastar signal carries it into each new SSE request.
- Do not publish the selection change through NATS or add a new NATS subject/bucket; NATS continues notifying the replacement stream about ordinary application-data changes.
- Change only the menu's empty signal fallback from `""` to numeric `0`; `billettholderSelection.initialize(...)` already normalizes and validates selected IDs.
- Do not override the Datastar request payload; `/profile/api` must consume the public `$billettHolderId` signal sent by default.
- Do not duplicate the menu's associated-billettholder validation in the profile template; authorization remains mandatory in `/profile/api`.
- Preserve every query parameter except the `b_id` value being updated.
- Use `history.replaceState`; billettholder switches must not add browser-history entries.
- Keep the Datastar action URI exactly `/profile/api` and use `requestCancellation: "auto"` so a new request cancels its predecessor.
- Keep the current infinite-retry behavior for an active profile stream.
- Do not manually run templ generation; the existing server watcher owns generation.
- Do not create or run automated tests unless the user removes that constraint. Use source inspection and the manual verification checklist in this plan.

---

## Requirements

1. Selecting a billettholder from the global menu while `/profile` is open must update `#profile-main-column` without a document navigation.
2. The browser URL must reflect the active billettholder as `b_id=<id>` so an ordinary refresh retains the same server-side selection.
3. Starting a new `/profile/api` request must automatically cancel the previous same-URL request. An old stream must not patch the profile after a newer selection.
4. Rapid selection changes must converge on the final billettholder with at most one active profile stream.
5. Existing query parameters and URL fragments must survive a selection change.
6. Users without billettholdere must use numeric signal value `0`, retain a live `/profile/api` connection, and have no visible `b_id` parameter.
7. The backend must validate query `b_id` for full-page requests and signal `billettHolderId` for live requests against the authenticated user's billettholdere; client validation is a UX guard, not authorization.
8. A server restart or transient network failure must reconnect the currently selected billettholder's stream.
9. After switching, existing NATS broadcasts must rerender the replacement stream with the new billettholder, never revive the cancelled stream or its old selection.

## File Structure

- Modify `components/header/menu.templ`: replace only the signal's empty fallback with numeric `0`.
- Modify `pages/profile/profile_index.templ`: remove the duplicate client validation and reload script, then inline the URL/SSE effect on the profile container.
- Modify `pages/profile/profile.go`: read and validate `billettHolderId` for `/profile/api`, retain query-based `b_id` for full-page rendering, and remove the valid-ID list used only by the deleted script.
- Read only `service/live/live.go`: rely on `Manager.Stream` stopping the cancelled stream's NATS bucket watchers and creating fresh watchers for the replacement stream.

## Selected Design

The menu already calls `billettholderSelection.initialize(...)`, which normalizes stored IDs and rejects selections that are no longer associated with the user. The profile therefore consumes `$billettHolderId` directly instead of receiving another JSON list and repeating the same validation.

```text
menu selection
    -> billettholderSelection.set(...) persists local storage
    -> menu event updates numeric $billettHolderId
    -> history.replaceState updates visible b_id
    -> Datastar cancels the old SSE and opens a new /profile/api SSE
    -> live.Manager attaches the replacement stream to existing NATS buckets
    -> later NATS broadcasts rerender the replacement stream
```

Only the SSE connection is replaced. There is no document navigation, no `window.location.replace(...)`, and no attempt to make NATS read browser local storage.

The profile container needs only this effect:

```templ
data-effect="
    const selectedID = Number($billettHolderId ?? 0);
    const profileURL = new URL(window.location.href);
    if (selectedID > 0) {
        profileURL.searchParams.set('b_id', String(selectedID));
    } else {
        profileURL.searchParams.delete('b_id');
    }
    window.history.replaceState(null, '', profileURL);
    @get('/profile/api', {
        requestCancellation: 'auto',
        retryMaxCount: Infinity,
        retryInterval: 1000,
        retryMaxWait: 30000,
    })
"
```

The client produces only associated IDs during normal use. `/profile/api` remains authoritative: a manually altered or stale signal falls back to a billettholder related to the authenticated user rather than bypassing authorization.

The fixed action URI is intentional. Datastar sends public signals automatically, so `/profile/api` receives `billettHolderId` without a payload override. Repeating the same action URI with `requestCancellation: "auto"` lets Datastar cancel the previous request before opening the replacement stream: [Datastar actions reference](https://data-star.dev/reference/actions#request-cancellation).

---

### Task 1: Use a Numeric Empty Selection in the Menu

**Files:**
- Modify: `components/header/menu.templ:314-318`

**Interfaces:**
- Consumes: `window.conorganizer.billettholderSelection.initialize(...)`.
- Produces: menu-owned signal `$billettHolderId`, using numeric `0` only when initialization returns no billettholder.

- [ ] **Step 1: Replace the menu's empty fallback**

Keep the existing event handler and initialization logic; change only both empty-string fallbacks to zero:

```templ
data-signals:billett-holder-id="0"
data-on:menu-billettholder-change="$billettHolderId = evt.detail.id"
data-init={ fmt.Sprintf("const selectedBillettholder = window.conorganizer.billettholderSelection.initialize(%s, %s); $billettHolderId = selectedBillettholder?.Id ?? 0;", associatedBillettholdereJSON, currentBillettholderJSON) }
```

- [ ] **Step 2: Inspect the menu signal assignment**

```powershell
rg -n "data-signals:billett-holder-id|menu-billettholder-change|selectedBillettholder\?\.Id" components/header/menu.templ
git diff --check -- components/header/menu.templ
```

Expected: the menu uses `0` for both empty fallbacks; the event continues assigning the already-normalized numeric `evt.detail.id`.

---

### Task 2: Read and Validate the Signal in the Profile Live Endpoint

**Files:**
- Modify: `pages/profile/profile.go:3-23`
- Modify: `pages/profile/profile.go:35-65`
- Modify: `pages/profile/profile.go:84-100`
- Modify: `pages/profile/profile.go:281-287`

**Interfaces:**
- Consumes: Datastar request signal `billettHolderId: number`.
- Produces: `/profile/api` renders with a captured, relation-validated billettholder ID; the full-page route continues using query `b_id`.

- [ ] **Step 1: Import the Datastar Go SDK**

Add:

```go
datastar "github.com/starfederation/datastar-go/datastar"
```

- [ ] **Step 2: Read the selection before opening the live stream**

At the beginning of the `/profile/api` GET handler, read the public signal once. The value is then fixed for that stream; a later selection opens a replacement request.

```go
var signals struct {
	BillettHolderID int `json:"billettHolderId"`
}

if err := datastar.ReadSignals(r, &signals); err != nil {
	http.Error(w, "Ugyldig profilvalg.", http.StatusBadRequest)
	return
}
```

Keep the existing `requestLogger`; signal decoding does not require a logging refactor. A malformed Datastar request returns `400`, which the existing HTTP middleware records.

- [ ] **Step 3: Validate the captured signal on every live render**

Replace the live endpoint's call to `selectedBillettholderIDFromRequest` with the existing relation and fallback helpers:

```go
selectedBillettholderID := signals.BillettHolderID
if !hasBillettholderID(billettholdere, selectedBillettholderID) {
	selectedBillettholderID = defaultSelectedBillettholderID(user, billettholdere, requestLogger)
}
```

The captured signal stays fixed for that SSE connection. A selection change opens a replacement request with a new captured value. Keep `selectedBillettholderIDFromRequest` unchanged for full-page requests and refreshes that use visible `b_id`.

- [ ] **Step 4: Remove the client-only valid-ID plumbing**

Delete this full-page handler line:

```go
validBillettholderIDs := billettholderIDs(billettholdere)
```

Remove `validBillettholderIDs` from the `ProfilePage(...)` call, and delete the now-unused `billettholderIDs(...)` helper. Task 3 removes the matching template parameter and data attribute.

- [ ] **Step 5: Inspect the two server selection entry points**

```powershell
rg -n -C 4 "ReadSignals|signals.BillettHolderID|selectedBillettholderIDFromRequest|hasBillettholderID" pages/profile/profile.go
rg -n "validBillettholderIDs|billettholderIDs" pages/profile/profile.go
git diff --check -- pages/profile/profile.go
```

Expected: `/profile/api` reads and validates the signal inline; the full `/profile` handler continues using the query helper; the second `rg` prints no matches.

---

### Task 3: Replace Reload-Based Profile Synchronization With One Fixed-URL Effect

**Files:**
- Modify: `pages/profile/profile_index.templ:3-13`
- Modify: `pages/profile/profile_index.templ:23-119`

**Interfaces:**
- Consumes: menu-owned Datastar signal `$billettHolderId`; fixed action URI `GET /profile/api`.
- Produces: a URL whose visible `b_id` mirrors the signal and one automatically managed same-URL SSE request.

- [ ] **Step 1: Remove parameters and imports used only by duplicate client validation**

Remove `validBillettholderIDs []int` from the `ProfilePage` signature. Keep `selectedBillettholderID int`, because the server still uses it for the initial `ProfileMainColumn` render.

Remove both now-unused imports:

```go
"fmt"
"github.com/Regncon/conorganizer/service/live"
```

- [ ] **Step 2: Replace the profile attributes with one inline effect**

Replace these section attributes:

```templ
data-profile-selected-billettholder-id={ fmt.Sprintf("%d", selectedBillettholderID) }
data-profile-valid-billettholder-ids={ templ.JSONString(validBillettholderIDs) }
data-init={ live.DatastarInitExpression("'/profile/api' + window.location.search") }
```

with:

```templ
data-effect="
	const selectedID = Number($billettHolderId ?? 0);
	const profileURL = new URL(window.location.href);
	if (selectedID > 0) {
		profileURL.searchParams.set('b_id', String(selectedID));
	} else {
		profileURL.searchParams.delete('b_id');
	}
	window.history.replaceState(null, '', profileURL);
	@get('/profile/api', {
		requestCancellation: 'auto',
		retryMaxCount: Infinity,
		retryInterval: 1000,
		retryMaxWait: 30000,
	})
"
```

The effect intentionally has one reactive input and no helper function, duplicate signal declaration, data attributes, JSON parsing, or payload override.

- [ ] **Step 3: Remove the old reload synchronization**

Delete the entire inline `<script>` immediately after the profile `</section>`. It currently reads local storage, mutates `b_id`, subscribes with `selection.onChange`, and calls `window.location.replace`.

The menu continues owning local storage and `$billettHolderId`; the profile owns only URL synchronization and its SSE request.

- [ ] **Step 4: Inspect the completed source for competing stream owners**

Run these source-only checks:

```powershell
rg -n "window\.location\.replace|DatastarInitExpression|selection\.onChange" pages/profile/profile_index.templ
rg -n "data-effect|@get|requestCancellation" pages/profile/profile_index.templ
rg -n "AbortController|profileStreamController|payload:" pages/profile/profile_index.templ
git diff --check -- pages/profile/profile_index.templ
```

Expected results:

- The first `rg` prints no matches.
- The second `rg` shows one effect, one fixed `/profile/api` action, and `requestCancellation: "auto"`.
- The third `rg` prints no matches; the effect has no manual request controller or payload override.
- `git diff --check` prints nothing.

- [ ] **Step 5: Review the server cancellation path without changing it**

Confirm these existing behaviors remain intact:

```powershell
rg -n -C 3 "case <-ctx.Done\(\)|watcher\.watcher\.Stop" service/live/live.go
rg -n -C 3 "ReadSignals|signals.BillettHolderID|selectedBillettholderIDFromRequest|hasBillettholderID" pages/profile/profile.go
```

Expected results:

- `Manager.Stream` exits on request-context cancellation and stops every bucket watcher in its deferred cleanup.
- The replacement `Manager.Stream` call creates new watchers for the existing events, interests, and billettholdere buckets; no NATS subject or selection-state change is required.
- `/profile/api` validates the signal ID, while the full-page handler validates query `b_id`; both require a billettholder related to the authenticated user.

---

### Task 4: Verify No-Reload Switching and Stream Ownership in the Browser

**Files:**
- Verify only: `pages/profile/profile_index.templ`
- Verify only: `pages/profile/profile.go`
- Verify only: `service/live/live.go`

**Interfaces:**
- Consumes: the managed profile effect from Task 3 and two billettholdere with distinguishable profile data.
- Produces: evidence that URL state, rendered data, connection cancellation, retries, and browser history behave as required.

- [ ] **Step 1: Establish the initial connection**

Open `/profile` while signed in as a user related to at least two billettholdere. In browser developer tools, filter Network requests by `/profile/api`.

Expected:

- The address bar contains the active `b_id` after initialization.
- There is one pending `/profile/api` SSE request whose Datastar signal envelope contains the numeric active `billettHolderId`.
- The profile main column shows that billettholder's program and interests.

- [ ] **Step 2: Switch billettholder without a document navigation**

Keep the Network panel open, switch to a second billettholder from the menu modal, and observe the request list and address bar.

Expected:

- No new `document` request occurs.
- The old `/profile/api` request becomes cancelled.
- Exactly one new `/profile/api` request remains pending, and its Datastar signal envelope contains the second `billettHolderId`.
- The visible URL contains the second billettholder as `b_id=<id>`.
- `#profile-main-column` changes to the second billettholder's data.
- Other query parameters and the URL fragment are unchanged.

- [ ] **Step 3: Verify rapid changes cannot leave stale streams**

Switch between the two billettholdere several times quickly and stop on the second one.

Expected:

- Every superseded request is cancelled.
- At most one `/profile/api` request remains pending after interaction settles.
- The address bar and profile data both represent the final selection.
- A subsequent interest or event NATS broadcast patches the replacement stream without a document request and does not restore data from an earlier selection.

- [ ] **Step 4: Verify history and refresh behavior**

After switching, press Back once, then return to the profile and refresh normally.

Expected:

- Back does not cycle through prior billettholder selections because the URL was changed with `replaceState`.
- Refresh renders the same billettholder on the server because the current `b_id` is present in the URL.
- Local storage, the menu highlight, and profile data agree after refresh.

- [ ] **Step 5: Verify retry behavior**

With one billettholder selected, temporarily restart the development server and leave the profile tab open.

Expected:

- The active request retries after the connection failure.
- It reconnects with the same selected `b_id`.
- Only the final selected billettholder's stream resumes patching the profile.

- [ ] **Step 6: Verify the zero-billettholder fallback**

Open `/profile` as a valid user with no related billettholdere.

Expected:

- The URL has no `b_id` parameter.
- `$billettHolderId` is numeric `0`.
- One `/profile/api` request remains active with `billettHolderId: 0` in its Datastar signal envelope.
- The rest of the profile, including data not scoped to a billettholder, continues receiving live updates.

---

## Completion Criteria

- Changing billettholder on `/profile` performs no document reload.
- During normal switcher use, URL `b_id`, local storage, menu selection, rendered profile data, and active SSE request agree.
- Only one profile stream remains active after initialization or any number of selection changes.
- Automatically cancelled streams stop their server-side bucket watchers through request-context cancellation.
- Refresh, retry, browser history, backend fallback for an invalid signal, and zero-billettholder behavior match this plan.
