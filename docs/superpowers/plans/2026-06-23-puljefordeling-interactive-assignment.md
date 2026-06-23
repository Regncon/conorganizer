# Puljefordeling Interactive Per-Game Player Assignment — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Make each game card in the puljefordeling tab interactive — add players via a shared searchable picker ("+"), remove seats ("×"), and re-solve a locked pulje ("Rerun fordeling") — mirroring the romfordeling assignment UX.

**Architecture:** Reuse existing manual-assignment backend (`event-players` endpoints, which write `source='manual'` pins the solver already honors by reducing capacity) and the existing `<admin-billettholder-search>` picker. Add a billettholder id to the emulation's player rows so the "×" can target a player, export the picker from `formsubmission`, extend the tab's event card with "+"/"×", add one shared `<dialog>` per tab, and add one new `rerun` endpoint that calls the existing `CommitPuljeAssignments`.

**Tech Stack:** Go 1.26, templ (`go tool templ generate`), Datastar (`datastar-go`), chi router, SQLite, `service/live` SSE manager.

## Global Constraints

- Templates generate with `go tool templ generate -path <dir> -log-level error` (templ/air/task are Go tools, not on PATH).
- `*_templ.go` generated files are gitignored repo-wide — commit only `.templ`/`.go` source; never `git add` a `_templ.go`.
- All user-facing copy is Norwegian (Bokmål).
- Reused endpoints (do not re-path, do not reimplement): `POST /admin/approval/api/event-players/post/add_first_choice`, `POST /admin/approval/api/event-players/post/add_gm`, `PUT /admin/approval/api/event-players/update_status`. All three already broadcast `live.BucketInterests`.
- The solver already reduces a game's effective capacity by the number of manual pins (`solver.SolveSlotFixed`) — no solver change.
- "×" visibility: when the pulje is **open**, only on pinned (`source='manual'`) players; when **locked/completed**, on every seated player.
- "Rerun fordeling" button + endpoint apply only to **locked/completed** puljer; open puljer 409 (already live).
- Lock/publish controls are unchanged (tab header only; dashboard untouched).
- Datastar signals are page-global; kebab `data-signals:assignment-event-id` maps to `$assignmentEventId`.
- Verification gate: `go test ./...` and `golangci-lint run` both clean.

---

## File Structure

- `service/puljefordeling/emulate.go` (modify) — add `BillettholderID int` to `AssignedPlayer`; populate it in `assignedPlayers`.
- `service/puljefordeling/emulate_pulje_test.go` (modify) — assert the new field is populated.
- `components/formsubmission/who_is_interested.templ` (modify) — add exported `PuljeAssignmentSearch(db, logger, pulje)` templ wrapping the existing picker, and exported `GetBillettholdereForPulje` helper if needed.
- `components/formsubmission/assignment_search_test.go` (create) — render test for `PuljeAssignmentSearch`.
- `pages/admin/puljefordeling_tab.templ` (modify) — interactive event card ("+"/"×"), shared `<dialog>`, signal wrapper on the shell, rerun button; pass pulje+status to the card.
- `pages/admin/puljefordeling_tab_test.go` (modify) — card render tests for "×" gating, "+", rerun button visibility.
- `pages/admin/puljefordeling_tab_rerun_test.go` (create) — rerun endpoint behavior test.

---

## Task 1: Add billettholder id to emulated player rows

**Files:**
- Modify: `service/puljefordeling/emulate.go`
- Test: `service/puljefordeling/emulate_pulje_test.go`

**Interfaces:**
- Consumes: existing `assignedPlayers(...)` (parses each solver player id with `strconv.Atoi(id)`), `AssignedPlayer`.
- Produces: `AssignedPlayer.BillettholderID int` — the billettholder id of the seated player (0 if the id was non-numeric).

- [ ] **Step 1: Write the failing test**

Append to `service/puljefordeling/emulate_pulje_test.go`:

```go
func TestEmulatePulje_AssignedPlayerHasBillettholderID(t *testing.T) {
	db, _ := testutil.CreateTestDBAndLogger(t, "test_emulate_pulje_bhid")

	const fredag = models.PuljeFredagKveld
	seedPulje(t, db, fredag, "Fredag Kveld", "2026-09-04T18:00:00Z")
	seedEvent(t, db, "evF", "Fredagsspill", 4, fredag)
	seedParticipant(t, db, 7, "Anna", "A")
	seedInterest(t, db, 7, "evF", fredag, models.InterestLevelHigh)

	pulje, _, err := EmulatePulje(db, fredag)
	if err != nil {
		t.Fatalf("EmulatePulje: %v", err)
	}
	ev, ok := findEvent(pulje, "evF")
	if !ok {
		t.Fatal("evF missing")
	}
	if len(ev.AssignedPlayers) != 1 {
		t.Fatalf("expected 1 assigned player, got %d", len(ev.AssignedPlayers))
	}
	if ev.AssignedPlayers[0].BillettholderID != 7 {
		t.Errorf("expected BillettholderID 7, got %d", ev.AssignedPlayers[0].BillettholderID)
	}
}
```

- [ ] **Step 2: Run the test to verify it fails**

Run: `go test ./service/puljefordeling/ -run TestEmulatePulje_AssignedPlayerHasBillettholderID -v`
Expected: FAIL — `ev.AssignedPlayers[0].BillettholderID undefined`.

- [ ] **Step 3: Add the field and populate it**

In `service/puljefordeling/emulate.go`, add the field to `AssignedPlayer` (after `Name`):

```go
type AssignedPlayer struct {
	Name           string
	BillettholderID int                  // billettholder id of the seated player (0 if non-numeric)
	IsDM           bool                 // runs at least one game in the weekend (DM bump)
	Level          models.InterestLevel // their interest in the game they got
	Moved          bool                 // relocated off a higher-scoring event by the solver to make room for others
	Pinned         bool                 // manually placed (source=manual); honored by the solver, not chosen by it
}
```

In `assignedPlayers`, set it inside the loop where `bh` is parsed:

```go
		ap := AssignedPlayer{
			Name:            names[bh],
			BillettholderID: bh,
			IsDM:            dmSet[bh],
			Moved:           moved[id],
			Pinned:          pinned[seatKey(eventID, id)],
		}
```

(The non-numeric fallback row `AssignedPlayer{Name: id}` keeps `BillettholderID: 0`.)

- [ ] **Step 4: Run the test to verify it passes**

Run: `go test ./service/puljefordeling/ -run TestEmulatePulje -v`
Expected: PASS (all `TestEmulatePulje_*`).

- [ ] **Step 5: Commit**

```bash
git add service/puljefordeling/emulate.go service/puljefordeling/emulate_pulje_test.go
git commit -m "feat(puljefordeling): expose billettholder id on emulated player rows"
```

---

## Task 2: Export the assignment-search picker from formsubmission

**Files:**
- Modify: `components/formsubmission/who_is_interested.templ`
- Test: `components/formsubmission/assignment_search_test.go`

**Interfaces:**
- Consumes: existing unexported `getBillettholdere(db, logger) ([]Billettholder, error)` and `billettholderAssignmentActions(eventId string, puljeId models.Pulje, billettholdere []Billettholder)` (its add buttons read `$assignmentEventId`/`$assignmentBillettholderId` from signals and bake only `puljeId`; the `eventId` arg is unused by the buttons).
- Produces: `templ PuljeAssignmentSearch(db *sql.DB, logger *slog.Logger, pulje models.Pulje)` — fetches all billettholdere and renders the searchable picker + "Legg til" buttons for the given pulje. Exported for use from `package admin`.

- [ ] **Step 1: Write the failing test**

Create `components/formsubmission/assignment_search_test.go`:

```go
package formsubmission

import (
	"database/sql"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/templtest"
)

func TestPuljeAssignmentSearch_RendersPickerAndButtons(t *testing.T) {
	db, logger := testutil.CreateTestDBAndLogger(t, "test_pulje_assignment_search")

	if _, err := db.Exec(
		`INSERT INTO billettholdere (id, first_name, last_name, ticket_type_id, ticket_type, order_id, ticket_id)
		 VALUES (1, 'Anna', 'A', 0, '', 0, 1)`,
	); err != nil {
		t.Fatalf("seed billettholder: %v", err)
	}

	doc := templtest.Render(t, PuljeAssignmentSearch(db, logger, models.PuljeFredagKveld))

	if doc.Find("admin-billettholder-search").Length() == 0 {
		t.Errorf("expected the search web component in the picker")
	}
	if doc.Find("button:contains('Legg til som førsteval')").Length() == 0 {
		t.Errorf("expected the 'add as first choice' button")
	}
}
```

- [ ] **Step 2: Run the test to verify it fails**

Run: `go tool templ generate -path components/formsubmission -log-level error && go test ./components/formsubmission/ -run TestPuljeAssignmentSearch -v`
Expected: FAIL — `undefined: PuljeAssignmentSearch`.

- [ ] **Step 3: Add the exported wrapper templ**

In `components/formsubmission/who_is_interested.templ`, add (near `billettholderAssignmentActions`):

```go
// PuljeAssignmentSearch renders the searchable billettholder picker and the
// "add as player / add as GM" buttons for one pulje, for reuse outside the
// event-edit page (e.g. the puljefordeling tab). The target event is taken from
// the $assignmentEventId signal set by the caller before opening the dialog.
templ PuljeAssignmentSearch(db *sql.DB, logger *slog.Logger, pulje models.Pulje) {
	{{ billettholdere, err := getBillettholdere(db, logger) }}
	if err != nil {
		<p style="color: var(--color-error);">Kunne ikke hente billettholdere</p>
	} else {
		@billettholderAssignmentActions("", pulje, billettholdere)
	}
}
```

(`database/sql`, `log/slog`, `models` are already imported in this file.)

- [ ] **Step 4: Run the test to verify it passes**

Run: `go tool templ generate -path components/formsubmission -log-level error && go test ./components/formsubmission/ -run TestPuljeAssignmentSearch -v`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add components/formsubmission/who_is_interested.templ components/formsubmission/assignment_search_test.go
git commit -m "feat(formsubmission): export PuljeAssignmentSearch picker for reuse"
```

---

## Task 3: Interactive event card, shared dialog, and signal wrapper

**Files:**
- Modify: `pages/admin/puljefordeling_tab.templ`
- Test: `pages/admin/puljefordeling_tab_test.go`

**Interfaces:**
- Consumes: `pf.EmulatedEvent` (`.EventID`, `.AssignedPlayers` with `.BillettholderID`/`.Pinned` from Task 1); `formsubmission.PuljeAssignmentSearch` (Task 2); existing `puljeIsLocked`/`puljeIsCompleted` (in `pages/admin/puljefordeling.go`); the reused `update_status` endpoint.
- Produces: updated `PuljefordelingTabPage` (signal wrapper + dialog), `PuljefordelingTabContent` (passes pulje+removable flag to the card), `puljefordelingTabEventCard(ev, pulje, removableAll)`, new `puljefordelingAssignmentDialog(db, logger, pulje)`.

- [ ] **Step 1: Restructure the shell to wrap signals + dialog around the content**

In `pages/admin/puljefordeling_tab.templ`, replace `PuljefordelingTabPage`:

```go
templ PuljefordelingTabPage(db *sql.DB, logger *slog.Logger, pulje models.Pulje) {
	<div
		class="puljefordeling-tab-shell"
		data-signals:assignment-event-id="''"
		data-signals:assignment-pulje-id={ fmt.Sprintf("'%s'", string(pulje)) }
		data-signals:assignment-billettholder-id="0"
		data-signals:assignment-is-player="false"
		data-signals:assignment-is-gm="false"
		data-signals:clear-input="0"
	>
		@puljefordelingAssignmentDialog(db, logger, pulje)
		@PuljefordelingTabContent(db, logger, pulje)
	</div>
}

templ puljefordelingAssignmentDialog(db *sql.DB, logger *slog.Logger, pulje models.Pulje) {
	<dialog id="pulje-assignment-dialog" closedBy="any">
		@formsubmission.PuljeAssignmentSearch(db, logger, pulje)
		<button type="button" class="btn btn--outline" data-on:click="document.getElementById('pulje-assignment-dialog').close()">Lukk</button>
	</dialog>
}
```

Add to the file's import block: `"github.com/Regncon/conorganizer/components/formsubmission"` (a plain-text close button avoids any icons dependency).

- [ ] **Step 2: Pass pulje + removable flag from content to the card**

In `PuljefordelingTabContent`, change the event-grid loop to compute the removable flag once and pass it down. Replace the grid block:

```go
			{{ removableAll := puljeIsLocked(status) || puljeIsCompleted(status) }}
			<div class="puljefordeling-tab-grid">
				for _, ev := range em.Events {
					@puljefordelingTabEventCard(ev, pulje, removableAll)
				}
			</div>
```

- [ ] **Step 3: Extend the event card with "×" and "+"**

Replace `puljefordelingTabEventCard` with:

```go
templ puljefordelingTabEventCard(ev pf.EmulatedEvent, pulje models.Pulje, removableAll bool) {
	<div class="item-card puljefordeling-tab-event">
		<div class="puljefordeling-tab-event-top">
			<h3>{ ev.Title }</h3>
			<span class="puljefordeling-tab-cap">{ fmt.Sprintf("%d / %d", len(ev.AssignedPlayers), ev.Capacity) }</span>
		</div>
		if ev.GMName != "" {
			<p class="puljefordeling-tab-gm">Spilleder: { ev.GMName }</p>
		}
		if ev.Undersubscribed {
			<span class="puljefordeling-tab-badge">Få deltakere</span>
		}
		if len(ev.AssignedPlayers) > 0 {
			<ul class="puljefordeling-tab-players">
				for _, pl := range ev.AssignedPlayers {
					<li class={ templ.KV("puljefordeling-tab-moved", pl.Moved) }>
						<span class="puljefordeling-tab-emoji">{ pl.Level.Emoji() }</span>
						<span class={ templ.KV("puljefordeling-tab-dm", pl.IsDM) }>{ pl.Name }</span>
						if pl.Pinned {
							<span class="puljefordeling-tab-pinned" title="Manuelt plassert">📌</span>
						}
						if (removableAll || pl.Pinned) && pl.BillettholderID > 0 {
							<button
								type="button"
								class="puljefordeling-tab-remove"
								title="Fjern fra spillet"
								data-on:click={ fmt.Sprintf(
									"$assignmentEventId = '%s'; $assignmentPuljeId = '%s'; $assignmentBillettholderId = %d; $assignmentIsPlayer = false; $assignmentIsGm = false; @put('/admin/approval/api/event-players/update_status')",
									ev.EventID, string(pulje), pl.BillettholderID,
								) }
							>×</button>
						}
					</li>
				}
			</ul>
		} else {
			<p class="puljefordeling-tab-empty">Ingen deltakere</p>
		}
		<button
			type="button"
			class="btn btn--outline puljefordeling-tab-add"
			data-on:click={ fmt.Sprintf("$assignmentEventId = '%s'; document.getElementById('pulje-assignment-dialog').showModal()", ev.EventID) }
		>+</button>
	</div>
}
```

- [ ] **Step 4: Add styles for the "×" / "+" affordances**

Inside the existing `puljefordelingTabStyles` `<style>` block (append before `</style>`):

```css
		.puljefordeling-tab-players li {
			justify-content: flex-start;
		}
		.puljefordeling-tab-remove {
			margin-left: auto;
			background: transparent;
			border: none;
			color: var(--color-text-error, #e05050);
			cursor: pointer;
			font-size: 1.1rem;
			line-height: 1;
			padding: 0 var(--spacing-1x);
		}
		.puljefordeling-tab-add {
			align-self: center;
			margin-top: var(--spacing-2x);
			min-width: 3rem;
		}
```

- [ ] **Step 5: Write the card render tests**

Append to `pages/admin/puljefordeling_tab_test.go` (reusing the `seedTab*` helpers already in that file):

```go
func TestPuljefordelingTabContent_RemoveGatedByState(t *testing.T) {
	db, logger := testutil.CreateTestDBAndLogger(t, "test_tab_remove_gating")

	const fredag = models.PuljeFredagKveld
	seedTabPulje(t, db, fredag, "Fredag Kveld", "2026-09-04T18:00:00Z")
	seedTabEvent(t, db, "evF", "Fredagsspill", 4, fredag)
	seedTabParticipant(t, db, 1, "Anna", "A")
	seedTabInterest(t, db, 1, "evF", fredag, models.InterestLevelHigh)

	// Open pulje: Anna is solver-placed (not pinned) -> no remove button.
	doc := templtest.Render(t, PuljefordelingTabContent(db, logger, fredag))
	if doc.Find(".puljefordeling-tab-remove").Length() != 0 {
		t.Errorf("open pulje: solver-placed player must not have a remove button")
	}
	if doc.Find(".puljefordeling-tab-add").Length() == 0 {
		t.Errorf("expected an add (+) button on the game card")
	}

	// Lock the pulje: every seated player becomes removable.
	if _, err := db.Exec(`UPDATE puljer SET status = 'Locked' WHERE id = ?`, string(fredag)); err != nil {
		t.Fatalf("lock pulje: %v", err)
	}
	// Persist Anna as a committed seat so the locked view shows her.
	if _, err := db.Exec(
		`INSERT INTO relation_events_players (event_id, pulje_id, billettholder_id, role, source)
		 VALUES ('evF', ?, 1, 'Player', 'solver')`, string(fredag),
	); err != nil {
		t.Fatalf("seed committed seat: %v", err)
	}
	docLocked := templtest.Render(t, PuljefordelingTabContent(db, logger, fredag))
	if docLocked.Find(".puljefordeling-tab-remove").Length() == 0 {
		t.Errorf("locked pulje: seated player must have a remove button")
	}
}
```

- [ ] **Step 6: Generate and run the tests**

Run:
```bash
go tool templ generate -path pages/admin -log-level error
go test ./pages/admin/ -run TestPuljefordelingTabContent -v
```
Expected: PASS (existing render test + the two new ones).

- [ ] **Step 7: Commit**

```bash
git add pages/admin/puljefordeling_tab.templ pages/admin/puljefordeling_tab_test.go
git commit -m "feat(puljefordeling): interactive game cards with add/remove and shared picker dialog"
```

---

## Task 4: Rerun endpoint and button for locked puljer

**Files:**
- Modify: `pages/admin/puljefordeling_tab.templ` (add the `rerun` route to `SetupPuljefordelingTabRoute`; add the button to the tab header)
- Test: `pages/admin/puljefordeling_tab_rerun_test.go`

**Interfaces:**
- Consumes: `pf.CommitPuljeAssignments(db, pulje, logger) error`; `models.ParsePulje`; `models.PuljeStatusOpen/Locked/Completed`; `liveManager.Broadcast`; existing route group in `SetupPuljefordelingTabRoute`.
- Produces: `PUT /admin/puljefordeling/api/{pulje}/rerun`; a "Rerun fordeling" button rendered in the tab header only when locked/completed.

- [ ] **Step 1: Write the failing endpoint test**

Create `pages/admin/puljefordeling_tab_rerun_test.go`:

```go
package admin

import (
	"database/sql"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/service/live"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/go-chi/chi/v5"
)

func setupRerunRouter(db *sql.DB, logger *slog.Logger) *chi.Mux {
	r := chi.NewRouter()
	r.Route("/admin/puljefordeling", func(pf chi.Router) {
		SetupPuljefordelingTabRoute(pf, db, &live.Manager{}, logger)
	})
	return r
}

func TestRerun_LockedPuljeReSolves(t *testing.T) {
	db, logger := testutil.CreateTestDBAndLogger(t, "test_rerun_locked")
	const fredag = models.PuljeFredagKveld
	seedTabPulje(t, db, fredag, "Fredag Kveld", "2026-09-04T18:00:00Z")
	seedTabEvent(t, db, "evF", "Fredagsspill", 4, fredag)
	seedTabParticipant(t, db, 1, "Anna", "A")
	seedTabInterest(t, db, 1, "evF", fredag, models.InterestLevelHigh)
	if _, err := db.Exec(`UPDATE puljer SET status = 'Locked' WHERE id = ?`, string(fredag)); err != nil {
		t.Fatalf("lock pulje: %v", err)
	}

	router := setupRerunRouter(db, logger)
	req := httptest.NewRequest(http.MethodPut, "/admin/puljefordeling/api/FredagKveld/rerun", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rec.Code)
	}
	var n int
	if err := db.QueryRow(
		`SELECT COUNT(*) FROM relation_events_players WHERE pulje_id = ? AND source = 'solver'`,
		string(fredag),
	).Scan(&n); err != nil {
		t.Fatalf("count solver seats: %v", err)
	}
	if n == 0 {
		t.Errorf("expected rerun to commit solver seats, found none")
	}
}

func TestRerun_OpenPuljeConflicts(t *testing.T) {
	db, logger := testutil.CreateTestDBAndLogger(t, "test_rerun_open")
	const fredag = models.PuljeFredagKveld
	seedTabPulje(t, db, fredag, "Fredag Kveld", "2026-09-04T18:00:00Z")

	router := setupRerunRouter(db, logger)
	req := httptest.NewRequest(http.MethodPut, "/admin/puljefordeling/api/FredagKveld/rerun", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusConflict {
		t.Errorf("expected 409 for open pulje, got %d", rec.Code)
	}
}

func TestRerun_InvalidPuljeBadRequest(t *testing.T) {
	db, logger := testutil.CreateTestDBAndLogger(t, "test_rerun_invalid")
	router := setupRerunRouter(db, logger)
	req := httptest.NewRequest(http.MethodPut, "/admin/puljefordeling/api/Nope/rerun", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid pulje, got %d", rec.Code)
	}
}
```

- [ ] **Step 2: Run the test to verify it fails**

Run: `go test ./pages/admin/ -run TestRerun -v`
Expected: FAIL — the `/api/{pulje}/rerun` route is not registered (404, not 204/409/400).

- [ ] **Step 3: Register the rerun route**

In `pages/admin/puljefordeling_tab.templ`, inside `SetupPuljefordelingTabRoute`, after the existing `router.Get("/api/{pulje}", ...)` block, add:

```go
	router.Put("/api/{pulje}/rerun", func(w http.ResponseWriter, r *http.Request) {
		puljeQuery := chi.URLParam(r, "pulje")
		puljeID, ok := models.ParsePulje(puljeQuery)
		if !ok {
			http.Error(w, "Expected a valid pulje ID, got: "+puljeQuery, http.StatusBadRequest)
			return
		}

		var raw string
		if err := db.QueryRow(`SELECT status FROM puljer WHERE id = ?`, string(puljeID)).Scan(&raw); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Error(w, "Pulje not found", http.StatusNotFound)
				return
			}
			logger.Error(fmt.Errorf("read pulje status for rerun: %w", err).Error(), "pulje_id", puljeID)
			http.Error(w, "Failed to read pulje status", http.StatusInternalServerError)
			return
		}
		status := models.PuljeStatus(raw)
		if status == models.PuljeStatusOpen {
			http.Error(w, "Open puljer re-solve live; rerun applies to locked puljer", http.StatusConflict)
			return
		}

		if err := pf.CommitPuljeAssignments(db, puljeID, logger); err != nil {
			logger.Error(fmt.Errorf("rerun commit for %s: %w", puljeID, err).Error(), "pulje_id", puljeID)
			http.Error(w, "Failed to re-run fordeling", http.StatusInternalServerError)
			return
		}

		if err := liveManager.Broadcast(r.Context(), live.BucketInterests); err != nil {
			logger.Error(fmt.Errorf("broadcast interests after rerun: %w", err).Error(), "pulje_id", puljeID)
		}
		if err := liveManager.Broadcast(r.Context(), live.BucketEvents); err != nil {
			logger.Error(fmt.Errorf("broadcast events after rerun: %w", err).Error(), "pulje_id", puljeID)
		}
		w.WriteHeader(http.StatusNoContent)
	})
```

Add `"errors"` to the file's import block if not already present (`database/sql`, `fmt`, `net/http`, `chi`, `live`, `pf`, `models` are already imported from Tasks in `puljefordeling_tab.templ`; verify and add only what's missing).

- [ ] **Step 4: Run the endpoint test to verify it passes**

Run: `go test ./pages/admin/ -run TestRerun -v`
Expected: PASS (locked→204, open→409, invalid→400).

- [ ] **Step 5: Add the "Rerun fordeling" button to the tab header**

In `PuljefordelingTabContent`, inside the `<header class="puljefordeling-tab-header">`, after `@puljefordelingTabStatusControls(...)`, add:

```go
				if puljeIsLocked(status) || puljeIsCompleted(status) {
					<button
						type="button"
						class="btn btn--outline puljefordeling-tab-rerun"
						data-on:click={ fmt.Sprintf("@put('/admin/puljefordeling/api/%s/rerun')", pulje) }
					>Rerun fordeling</button>
				}
```

- [ ] **Step 6: Add a render assertion for the button and generate**

Append to `pages/admin/puljefordeling_tab_test.go`:

```go
func TestPuljefordelingTabContent_RerunButtonOnlyWhenLocked(t *testing.T) {
	db, logger := testutil.CreateTestDBAndLogger(t, "test_tab_rerun_button")
	const fredag = models.PuljeFredagKveld
	seedTabPulje(t, db, fredag, "Fredag Kveld", "2026-09-04T18:00:00Z")
	seedTabEvent(t, db, "evF", "Fredagsspill", 4, fredag)

	if templtest.Render(t, PuljefordelingTabContent(db, logger, fredag)).
		Find(".puljefordeling-tab-rerun").Length() != 0 {
		t.Errorf("open pulje must not show the rerun button")
	}

	if _, err := db.Exec(`UPDATE puljer SET status = 'Locked' WHERE id = ?`, string(fredag)); err != nil {
		t.Fatalf("lock pulje: %v", err)
	}
	if templtest.Render(t, PuljefordelingTabContent(db, logger, fredag)).
		Find(".puljefordeling-tab-rerun").Length() == 0 {
		t.Errorf("locked pulje must show the rerun button")
	}
}
```

Run:
```bash
go tool templ generate -path pages/admin -log-level error
go test ./pages/admin/ -run 'TestRerun|TestPuljefordelingTabContent' -v
```
Expected: PASS.

- [ ] **Step 7: Commit**

```bash
git add pages/admin/puljefordeling_tab.templ pages/admin/puljefordeling_tab_test.go pages/admin/puljefordeling_tab_rerun_test.go
git commit -m "feat(puljefordeling): rerun endpoint and button for locked puljer"
```

---

## Task 5: Verification gate and manual smoke test

**Files:** none (verification only)

- [ ] **Step 1: Full test suite**

Run: `go test ./...`
Expected: all packages PASS.

- [ ] **Step 2: Linter**

Run: `go tool templ generate -path pages/admin -log-level error && go tool templ generate -path components/formsubmission -log-level error && golangci-lint run`
Expected: `0 issues`. Fix any unused-symbol / errcheck findings introduced by the changes (e.g., an unused import after the edits).

- [ ] **Step 3: Manual smoke test**

```bash
go tool task start
```

In a browser at `http://localhost:7331/admin/puljefordeling/FredagKveld`:
- Each game card shows a **+** button; clicking it opens the shared picker dialog; searching and clicking "Legg til som førsteval" adds the player, the card re-renders with the new 📌 pin, and capacity reflects it.
- On an **open** pulje, only 📌 (manual) players show a **×**; clicking it removes the pin and the card re-renders.
- **Lock** the pulje (tab header control): every seated player now shows **×**, and a **Rerun fordeling** button appears. Remove a player, then click Rerun fordeling — the freed seat is re-solved.
- The dashboard at `/admin/` is unchanged (single "Gå til puljefordeling" link; no inline status grid).

- [ ] **Step 4: Final commit (if smoke-test fixes were needed)**

```bash
git add -A
git commit -m "fix(puljefordeling): smoke-test adjustments"
```

---

## Self-Review Notes

- **Spec coverage:** behavior-by-state (Task 3 "×" gating + Task 4 rerun), interactive card + shared picker modal (Tasks 2-3), reuse of existing endpoints (Tasks 3-4 use them directly), new rerun endpoint with open→409/invalid→400 guards (Task 4), billettholder id for "×" targeting (Task 1), lock/publish unchanged + dashboard untouched (no task modifies `admin_page.templ` or the status route — by design), verification gate incl. `golangci-lint` (Task 5).
- **Type consistency:** `AssignedPlayer.BillettholderID int` defined in Task 1, consumed in Task 3's card. `PuljeAssignmentSearch(db, logger, pulje)` defined in Task 2, called in Task 3's dialog. `puljefordelingTabEventCard(ev, pulje, removableAll)` signature consistent between Task 3's definition and `PuljefordelingTabContent`'s call. The rerun path `/admin/puljefordeling/api/{pulje}/rerun` matches between Task 4's route, the header button, and the tests.
- **Resolved during planning:** picker reuse is signal-driven (no hard-coded event id — confirmed `who_is_interested.templ:639-650`); the reused endpoints live at `/admin/approval/api/event-players/...`; the solver already reduces capacity for pins (no solver change).
- **Implementer judgement calls flagged inline:** exact import additions to `puljefordeling_tab.templ` (Task 3 Step 1, Task 4 Step 3) — add only what's missing, the compiler/`golangci-lint` is the check.
