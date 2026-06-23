# Puljefordeling Per-Pulje Tabbed View Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Replace the split puljefordeling UI (inline status cards on `/admin/` + a separate all-puljer emulation page) with one tabbed page, one tab per pulje, each showing that pulje's live distribution plus its lock/publish control.

**Architecture:** Mirror the existing rooms-assignment pattern: a server-rendered shell with a `.tabs` bar and a `data-init` section that subscribes to a Datastar SSE stream; the stream's `Render` re-renders the selected pulje's distribution on `BucketInterests`/`BucketEvents` broadcasts. A new scoped service function solves the solver chronologically only up to the viewed pulje, so later puljer are never computed until they themselves are reached.

**Tech Stack:** Go 1.26, templ (`go tool templ generate`), Datastar (`datastar-go`), chi router, SQLite, `service/live` SSE manager.

## Global Constraints

- Templates are generated with `go tool templ generate` (templ/air/task are Go tools, not on PATH).
- All user-facing copy is Norwegian (Bokmål), matching existing puljefordeling text.
- The status-PUT endpoint stays at its current path `/admin/api/puljer/{puljeId}/status` (registered by the unchanged `puljefordelingStatusRoute`). The `puljeStatusUpdateAction` markup already targets it — do not move it.
- Backend solver/commit/emulate semantics are unchanged; only `EmulatePulje` is added.
- No manual seat editing in the tab (that lives on the rooms-assignment page).
- Verification gate (must both pass clean before "done"): `go test ./...` and `golangci-lint run`.
- New templ symbols in `package admin` must not collide with existing ones (`puljefordeling`, `puljeStatusCard`, `puljeStatusUpdateAction`, `puljeIsLocked`, `puljeIsCompleted`, `getPuljer`). Prefix new ones with `puljefordelingTab`.

---

## File Structure

- `service/puljefordeling/emulate.go` (modify) — add `EmulationMeta` struct + `EmulatePulje` function.
- `service/puljefordeling/emulate_pulje_test.go` (create) — tests for `EmulatePulje`.
- `pages/admin/puljefordeling_tab.templ` (create) — tabbed page: `puljefordelingTabIndex` (breadcrumbs+tabs layout), `PuljefordelingTabPage` (shell), `PuljefordelingTabContent` (SSE-rendered body), `puljefordelingTabEventCard`, `puljefordelingTabLegend`, `puljefordelingTabStatusControls`, `puljefordelingTabStyles`, and `SetupPuljefordelingTabRoute` (route registration).
- `pages/admin/puljefordeling_tab_test.go` (create) — render test for `PuljefordelingTabContent`.
- `pages/admin/admin.go` (modify) — register the new routes; remove the `/puljefordeling-emulate` route block and its import.
- `pages/admin/admin_page.templ` (modify) — remove inline `@puljefordeling(db)` section and the "Emulér puljefordeling" card; add one "Puljefordeling" card linking to `/admin/puljefordeling/`.
- `pages/admin/puljefordeling.templ` (modify) — delete the now-unused `puljefordeling` and `puljeStatusCard` templ funcs (keep `getPuljer`, `puljefordelingStatusRoute`, the status helpers, and `puljeStatusUpdateAction`).
- `pages/admin/puljefordeling_emulate/` (delete) — whole package removed.

---

## Task 1: Scoped per-pulje emulation (`EmulatePulje`)

**Files:**
- Modify: `service/puljefordeling/emulate.go`
- Test: `service/puljefordeling/emulate_pulje_test.go`

**Interfaces:**
- Consumes: existing unexported `seatingData` (`loadSeatingData`, fields `puljer`, `weekend`, `gms`, `names`, `prefs`, `dmSet`, `pinnedSet`, `year`), `solveChronological(upTo int)`, `shapePulje(...)`, `EmulatedPulje`.
- Produces:
  - `type EmulationMeta struct { Year int; PlayerCount int }`
  - `func EmulatePulje(db *sql.DB, puljeID models.Pulje) (EmulatedPulje, EmulationMeta, error)` — solves only up to `puljeID`'s chronological index and returns just that pulje. Errors if the pulje id is unknown.

- [ ] **Step 1: Write the failing test**

Create `service/puljefordeling/emulate_pulje_test.go`:

```go
package puljefordeling

import (
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil"
)

func TestEmulatePulje_ReturnsOnlyRequestedPulje(t *testing.T) {
	db, _ := testutil.CreateTestDBAndLogger(t, "test_emulate_pulje_single")

	const fredag = models.PuljeFredagKveld
	const lordag = models.PuljeLordagMorgen
	seedPulje(t, db, fredag, "Fredag Kveld", "2026-09-04T18:00:00Z")
	seedPulje(t, db, lordag, "Lørdag Morgen", "2026-09-05T10:00:00Z")
	seedEvent(t, db, "evF", "Fredagsspill", 4, fredag)
	seedEvent(t, db, "evL", "Lørdagsspill", 4, lordag)
	seedParticipant(t, db, 1, "Anna", "A")
	seedInterest(t, db, 1, "evF", fredag, models.InterestLevelHigh)
	seedInterest(t, db, 1, "evL", lordag, models.InterestLevelHigh)

	pulje, meta, err := EmulatePulje(db, fredag)
	if err != nil {
		t.Fatalf("EmulatePulje: %v", err)
	}
	if pulje.PuljeID != fredag {
		t.Errorf("expected pulje %s, got %s", fredag, pulje.PuljeID)
	}
	if _, ok := findEvent(pulje, "evF"); !ok {
		t.Errorf("expected fredag event evF in result")
	}
	if _, ok := findEvent(pulje, "evL"); ok {
		t.Errorf("lørdag event evL must NOT appear in the fredag pulje result")
	}
	if meta.PlayerCount != 1 {
		t.Errorf("expected PlayerCount 1, got %d", meta.PlayerCount)
	}
}

func TestEmulatePulje_UnknownPuljeErrors(t *testing.T) {
	db, _ := testutil.CreateTestDBAndLogger(t, "test_emulate_pulje_unknown")
	seedPulje(t, db, models.PuljeFredagKveld, "Fredag Kveld", "2026-09-04T18:00:00Z")

	if _, _, err := EmulatePulje(db, models.Pulje("DoesNotExist")); err == nil {
		t.Fatal("expected error for unknown pulje, got nil")
	}
}

func TestEmulatePulje_MatchesFullEmulationForLastPulje(t *testing.T) {
	db, _ := testutil.CreateTestDBAndLogger(t, "test_emulate_pulje_equiv")

	const fredag = models.PuljeFredagKveld
	const lordag = models.PuljeLordagMorgen
	seedPulje(t, db, fredag, "Fredag Kveld", "2026-09-04T18:00:00Z")
	seedPulje(t, db, lordag, "Lørdag Morgen", "2026-09-05T10:00:00Z")
	seedEvent(t, db, "evL", "Lørdagsspill", 4, lordag)
	seedParticipant(t, db, 1, "Anna", "A")
	seedInterest(t, db, 1, "evL", lordag, models.InterestLevelHigh)

	scoped, _, err := EmulatePulje(db, lordag)
	if err != nil {
		t.Fatalf("EmulatePulje: %v", err)
	}
	full, err := EmulateSeatings(db)
	if err != nil {
		t.Fatalf("EmulateSeatings: %v", err)
	}
	// lørdag is the last pulje chronologically -> index 1 in the full result.
	want := full.Puljer[1]
	if scoped.NewlySatisfied != want.NewlySatisfied || scoped.TotalScore != want.TotalScore {
		t.Errorf("scoped=%+v want=%+v", scoped, want)
	}
}
```

- [ ] **Step 2: Run the test to verify it fails**

Run: `go test ./service/puljefordeling/ -run TestEmulatePulje -v`
Expected: FAIL — `undefined: EmulatePulje` (and `EmulationMeta`).

- [ ] **Step 3: Implement `EmulationMeta` and `EmulatePulje`**

In `service/puljefordeling/emulate.go`, add after the `Emulation` struct (around line 57):

```go
// EmulationMeta carries the page-level numbers a per-pulje view needs that are
// not specific to one pulje.
type EmulationMeta struct {
	Year        int
	PlayerCount int // distinct participants with at least one interest
}

// EmulatePulje solves chronologically only up to (and including) puljeID and
// returns just that pulje's distribution. Earlier puljer are solved as needed
// (locked ones contribute their persisted seats); later puljer are never
// computed.
func EmulatePulje(db *sql.DB, puljeID models.Pulje) (EmulatedPulje, EmulationMeta, error) {
	d, err := loadSeatingData(db)
	if err != nil {
		return EmulatedPulje{}, EmulationMeta{}, err
	}

	idx := -1
	for i := range d.puljer {
		if d.puljer[i].ID == puljeID {
			idx = i
			break
		}
	}
	if idx < 0 {
		return EmulatedPulje{}, EmulationMeta{}, fmt.Errorf("pulje %q not found", puljeID)
	}

	_, results := d.solveChronological(idx)
	pid := string(d.puljer[idx].ID)
	shaped := shapePulje(
		d.puljer[idx], d.weekend.Slots[idx], results[idx],
		d.gms, d.names, d.prefs, d.dmSet, d.pinnedSet[pid],
	)

	return shaped, EmulationMeta{Year: d.year, PlayerCount: len(d.prefs)}, nil
}
```

(`fmt` and `models` are already imported in this file.)

- [ ] **Step 4: Run the test to verify it passes**

Run: `go test ./service/puljefordeling/ -run TestEmulatePulje -v`
Expected: PASS (all three).

- [ ] **Step 5: Commit**

```bash
git add service/puljefordeling/emulate.go service/puljefordeling/emulate_pulje_test.go
git commit -m "feat(puljefordeling): add EmulatePulje for scoped per-pulje emulation"
```

---

## Task 2: Per-pulje tab content templates

**Files:**
- Create: `pages/admin/puljefordeling_tab.templ`
- Test: `pages/admin/puljefordeling_tab_test.go`

**Interfaces:**
- Consumes: `puljefordeling.EmulatePulje` + `EmulationMeta`/`EmulatedPulje`/`EmulatedEvent`/`AssignedPlayer` (Task 1); existing `package admin` helpers `puljeStatusUpdateAction`, `puljeIsLocked`, `puljeIsCompleted`; `models.PuljeRow`/`getPuljer`; `service/live.DatastarInit`.
- Produces:
  - `templ PuljefordelingTabContent(db *sql.DB, logger *slog.Logger, pulje models.Pulje)` — the SSE-rendered body (used by Task 3's stream and by this task's test).
  - unexported `puljefordelingTabEventCard`, `puljefordelingTabLegend`, `puljefordelingTabStatusControls`, `puljefordelingTabStyles`.

> Note: this task creates the content templates and a render test that uses them, so nothing is orphaned. The shell, tabs layout, and routes come in Task 3.

- [ ] **Step 1: Write the content templates**

Create `pages/admin/puljefordeling_tab.templ`:

```go
package admin

import (
	"database/sql"
	"fmt"
	"log/slog"
	"strings"

	"github.com/Regncon/conorganizer/models"
	pf "github.com/Regncon/conorganizer/service/puljefordeling"
	"github.com/Regncon/conorganizer/service/live"
)

templ PuljefordelingTabContent(db *sql.DB, logger *slog.Logger, pulje models.Pulje) {
	{{
		em, meta, emErr := pf.EmulatePulje(db, pulje)
		status := puljeStatusFor(db, pulje)
	}}
	<section
		id="puljefordeling-tab"
		class="rooms-section"
		data-signals:pulje-status={ fmt.Sprintf("'%s'", string(status)) }
		data-init={ live.DatastarInit(fmt.Sprintf("/admin/puljefordeling/api/%s", pulje)) }
	>
		@puljefordelingTabStyles()
		if emErr != nil {
			<p class="puljefordeling-tab-error">Kunne ikke beregne fordeling: { emErr.Error() }</p>
		} else {
			<header class="puljefordeling-tab-header">
				<div class="puljefordeling-tab-heading">
					<h1 class="page-heading">{ em.Name }</h1>
					<div class="puljefordeling-tab-stats">
						<span>{ fmt.Sprintf("%d deltakere med ønsker", meta.PlayerCount) }</span>
						<span>{ fmt.Sprintf("%d fikk førstevalg", em.NewlySatisfied) }</span>
						if len(em.Unassigned) > 0 {
							<span class="puljefordeling-tab-warn">{ fmt.Sprintf("%d uten plass", len(em.Unassigned)) }</span>
						}
					</div>
				</div>
				@puljefordelingTabStatusControls(models.PuljeRow{ID: pulje, Name: em.Name, Status: status})
			</header>
			@puljefordelingTabLegend()
			<div class="puljefordeling-tab-grid">
				for _, ev := range em.Events {
					@puljefordelingTabEventCard(ev)
				}
			</div>
			if len(em.Unassigned) > 0 {
				<div class="puljefordeling-tab-unassigned">
					<h4>{ fmt.Sprintf("Uten plass (%d)", len(em.Unassigned)) }</h4>
					<p>{ strings.Join(em.Unassigned, ", ") }</p>
				</div>
			}
		}
	</section>
}

templ puljefordelingTabStatusControls(pulje models.PuljeRow) {
	<div class="puljefordeling-status-controls">
		<label class="puljefordeling-status-control">
			<input
				if puljeIsLocked(pulje.Status) {
					checked
				}
				data-on:click={ puljeStatusUpdateAction(
					pulje,
					fmt.Sprintf("Vil du endre lukking for %s?", pulje.Name),
					models.PuljeStatusLocked,
					models.PuljeStatusOpen,
				) }
				type="checkbox"
				class="checkbox input"
			/>
			<span>Puljefordeling lukket</span>
		</label>
		<label class="puljefordeling-status-control">
			<input
				if puljeIsCompleted(pulje.Status) {
					checked
				}
				data-on:click={ puljeStatusUpdateAction(
					pulje,
					fmt.Sprintf("Vil du endre publisering av puljefordeling for %s?", pulje.Name),
					models.PuljeStatusCompleted,
					models.PuljeStatusLocked,
				) }
				type="checkbox"
				class="checkbox input"
			/>
			<span>Puljefordeling publisert</span>
		</label>
	</div>
}

templ puljefordelingTabLegend() {
	<p class="puljefordeling-tab-legend">
		🔥 Veldig interessert · 👍 Middels · 🤷 Litt · 📌 Manuelt plassert ·
		<span class="puljefordeling-tab-dm">Turkis</span> = spilleder et annet sted ·
		<span class="puljefordeling-tab-moved-legend">Rød strek</span> = flyttet ned for å gi plass til andre
	</p>
}

templ puljefordelingTabEventCard(ev pf.EmulatedEvent) {
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
					</li>
				}
			</ul>
		} else {
			<p class="puljefordeling-tab-empty">Ingen deltakere</p>
		}
	</div>
}
```

- [ ] **Step 2: Add the styles template and the `puljeStatusFor` helper**

Append to `pages/admin/puljefordeling_tab.templ`:

```go
// puljeStatusFor reads a single pulje's current status, defaulting to Open on
// any error (the content render is best-effort; a status read failure should
// not blank the page).
func puljeStatusFor(db *sql.DB, pulje models.Pulje) models.PuljeStatus {
	var raw string
	if err := db.QueryRow(`SELECT status FROM puljer WHERE id = ?`, string(pulje)).Scan(&raw); err != nil {
		return models.PuljeStatusOpen
	}
	return models.PuljeStatus(raw)
}

templ puljefordelingTabStyles() {
	<style>
		.puljefordeling-tab-header {
			display: flex;
			align-items: flex-start;
			justify-content: space-between;
			gap: var(--spacing-4x);
			flex-wrap: wrap;
			border-bottom: 1px solid var(--bg-item-border);
			padding-bottom: var(--spacing-2x);
		}
		.puljefordeling-tab-stats {
			display: flex;
			gap: var(--spacing-4x);
			color: var(--color-text-soft);
			flex-wrap: wrap;
		}
		.puljefordeling-tab-warn {
			color: var(--color-text-warning, #e0a030);
		}
		.puljefordeling-tab-legend {
			color: var(--color-text-soft);
			font-size: var(--text-small, 0.85rem);
		}
		.puljefordeling-tab-grid {
			display: grid;
			grid-template-columns: repeat(auto-fit, minmax(var(--mobile-min-width), 1fr));
			gap: var(--spacing-4x);
		}
		.puljefordeling-tab-event {
			display: flex;
			flex-direction: column;
			gap: var(--spacing-2x);
			border: 1px solid var(--bg-item-border);
		}
		.puljefordeling-tab-event-top {
			display: flex;
			align-items: baseline;
			justify-content: space-between;
			gap: var(--spacing-2x);
		}
		.puljefordeling-tab-event-top h3 {
			color: var(--color-text-strong);
			font-size: var(--text-body);
		}
		.puljefordeling-tab-cap {
			color: var(--color-text-soft);
			white-space: nowrap;
		}
		.puljefordeling-tab-gm {
			color: var(--color-text-soft);
			font-style: italic;
		}
		.puljefordeling-tab-badge {
			align-self: flex-start;
			padding: 0 var(--spacing-2x);
			border-radius: var(--border-radius-1x);
			background-color: var(--bg-item);
			color: var(--color-text-warning, #e0a030);
			font-size: var(--text-small, 0.85rem);
		}
		.puljefordeling-tab-players {
			margin: 0;
			padding: 0;
			list-style: none;
			display: flex;
			flex-direction: column;
			gap: var(--spacing-1x);
			color: var(--color-text-soft);
		}
		.puljefordeling-tab-players li {
			display: flex;
			align-items: center;
			gap: var(--spacing-2x);
		}
		.puljefordeling-tab-emoji {
			flex: 0 0 auto;
			width: 1.4em;
			text-align: center;
		}
		.puljefordeling-tab-dm {
			color: turquoise;
		}
		.puljefordeling-tab-pinned {
			flex: 0 0 auto;
		}
		.puljefordeling-tab-moved {
			border-left: 3px solid #e05050;
			padding-left: var(--spacing-2x);
		}
		.puljefordeling-tab-moved-legend {
			border-left: 3px solid #e05050;
			padding-left: var(--spacing-1x);
		}
		.puljefordeling-tab-empty {
			color: var(--color-text-soft-50);
			font-style: italic;
		}
		.puljefordeling-tab-unassigned {
			color: var(--color-text-soft);
		}
		.puljefordeling-tab-error {
			color: var(--color-text-error, #e05050);
		}
	</style>
}
```

- [ ] **Step 3: Write the render test (with fully-specified seed helpers)**

Create `pages/admin/puljefordeling_tab_test.go`. The render helper is `testutil/templtest.Render`, which returns a `*goquery.Document` (same helper `pages/admin/rooms/rooms_page_test.go` uses). The seed SQL is copied verbatim from `service/puljefordeling/emulate_test.go`:

```go
package admin

import (
	"database/sql"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/templtest"
)

func TestPuljefordelingTabContent_RendersPuljeAndControls(t *testing.T) {
	db, logger := testutil.CreateTestDBAndLogger(t, "test_puljefordeling_tab_content")

	const fredag = models.PuljeFredagKveld
	seedTabPulje(t, db, fredag, "Fredag Kveld", "2026-09-04T18:00:00Z")
	seedTabEvent(t, db, "evF", "Fredagsspill", 4, fredag)
	seedTabParticipant(t, db, 1, "Anna", "A")
	seedTabInterest(t, db, 1, "evF", fredag, models.InterestLevelHigh)

	doc := templtest.Render(t, PuljefordelingTabContent(db, logger, fredag))

	if doc.Find("h3:contains('Fredagsspill')").Length() == 0 {
		html, _ := doc.Html()
		t.Errorf("expected event title in rendered tab, got:\n%s", html)
	}
	if doc.Find("span:contains('Puljefordeling lukket')").Length() == 0 {
		t.Errorf("expected lock control in rendered tab")
	}
	dataInit := doc.Find("#puljefordeling-tab").AttrOr("data-init", "")
	if !strings.Contains(dataInit, "/admin/puljefordeling/api/FredagKveld") {
		t.Errorf("expected data-init SSE url, got: %q", dataInit)
	}
}

func seedTabPulje(t *testing.T, db *sql.DB, id models.Pulje, name, startAt string) {
	t.Helper()
	if _, err := db.Exec(
		`INSERT INTO puljer (id, name, status, start_at, end_at) VALUES (?, ?, 'Open', ?, ?)`,
		string(id), name, startAt, startAt,
	); err != nil {
		t.Fatalf("seed pulje %s: %v", id, err)
	}
}

func seedTabEvent(t *testing.T, db *sql.DB, id, title string, maxPlayers int, pulje models.Pulje) {
	t.Helper()
	if _, err := db.Exec(
		`INSERT INTO events (id, title, intro, description, host_name, email, phone_number, max_players)
		 VALUES (?, ?, '', '', '', '', '', ?)`,
		id, title, maxPlayers,
	); err != nil {
		t.Fatalf("seed event %s: %v", id, err)
	}
	if _, err := db.Exec(
		`INSERT INTO relation_event_puljer (event_id, pulje_id, is_in_pulje) VALUES (?, ?, 1)`,
		id, string(pulje),
	); err != nil {
		t.Fatalf("place event %s in %s: %v", id, pulje, err)
	}
}

func seedTabParticipant(t *testing.T, db *sql.DB, id int, first, last string) {
	t.Helper()
	if _, err := db.Exec(
		`INSERT INTO billettholdere (id, first_name, last_name, ticket_type_id, ticket_type, order_id, ticket_id)
		 VALUES (?, ?, ?, 0, '', 0, ?)`,
		id, first, last, id,
	); err != nil {
		t.Fatalf("seed participant %d: %v", id, err)
	}
}

func seedTabInterest(t *testing.T, db *sql.DB, bhID int, eventID string, pulje models.Pulje, level models.InterestLevel) {
	t.Helper()
	if _, err := db.Exec(
		`INSERT INTO interests (billettholder_id, event_id, pulje_id, interest_level) VALUES (?, ?, ?, ?)`,
		bhID, eventID, string(pulje), string(level),
	); err != nil {
		t.Fatalf("seed interest bh=%d ev=%s: %v", bhID, eventID, err)
	}
}
```

Add `"strings"` to the import block (used by the `data-init` assertion).

- [ ] **Step 4: Generate templates and run the test**

Run:
```bash
go tool templ generate -path pages/admin -log-level error
go test ./pages/admin/ -run TestPuljefordelingTabContent -v
```
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add pages/admin/puljefordeling_tab.templ pages/admin/puljefordeling_tab_templ.go pages/admin/puljefordeling_tab_test.go
git commit -m "feat(puljefordeling): add per-pulje tab content templates"
```

---

## Task 3: Tabbed shell, tabs layout, routes; remove emulate page

**Files:**
- Modify: `pages/admin/puljefordeling_tab.templ` (add shell, tabs layout, route setup)
- Modify: `pages/admin/admin.go` (register new routes; remove emulate route + import)
- Delete: `pages/admin/puljefordeling_emulate/` (whole directory)

**Interfaces:**
- Consumes: `PuljefordelingTabContent` (Task 2); `getPuljer` (existing); `models.AllPuljer`, `models.ParsePulje`, `models.PuljeFredagKveld`; `layouts.Base`, `userctx.GetUserRequestInfo`, `components.Breadcrumbs`; `live.Manager.Stream`, `live.BucketInterests`, `live.BucketEvents`.
- Produces: `func SetupPuljefordelingTabRoute(router chi.Router, db *sql.DB, liveManager *live.Manager, logger *slog.Logger)` registering `/`, `/{pulje}`, and `/api/{pulje}` on the passed sub-router.

- [ ] **Step 1: Add the shell + tabs layout to `puljefordeling_tab.templ`**

Add these imports to the file's import block: `"net/http"`, `"github.com/Regncon/conorganizer/components"`, `"github.com/Regncon/conorganizer/layouts"`, `"github.com/Regncon/conorganizer/service/userctx"`, `"github.com/go-chi/chi/v5"`, `"context"`, `"github.com/a-h/templ"`.

Append:

```go
templ puljefordelingTabIndex(children templ.Component, puljeQuery models.Pulje) {
	@components.Breadcrumbs([]components.BreadcrumbPath{
		{Name: "Hjem", Url: "/"},
		{Name: "Admin", Url: "/admin/"},
		{Name: "Puljefordeling", Url: ""},
	})
	<div class="page-content-container flex flex-col">
		<style>
            .tabs {
                display: flex;
                overflow-x: auto;
                margin-top: var(--responsive-page-section-spacing);
                margin-bottom: -2px;
                z-index: 2;
            }
            .tabs a {
                white-space: nowrap;
                text-decoration: none;
                padding: 0.75rem 1rem;
                border: 2px solid transparent;
                border-radius: var(--border-radius-2x) var(--border-radius-2x) 0 0;
                color: var(--color-primary-text);
                font-weight: 500;
            }
            .tabs a:hover {
                color: var(--color-primary-hover);
            }
            .tabs a.active {
                background: var(--bg-surface);
                border-color: var(--color-primary);
                border-bottom: none;
                color: var(--color-primary);
            }
		</style>
		<div class="tabs">
			for _, pulje := range models.AllPuljer() {
				<a
					role="button"
					href={ templ.SafeURL("/admin/puljefordeling/" + string(pulje)) }
					class={ templ.KV("active", puljeQuery == pulje) }
				>{ string(pulje) }</a>
			}
		</div>
		@children
	</div>
}

templ PuljefordelingTabPage(db *sql.DB, logger *slog.Logger, pulje models.Pulje) {
	@PuljefordelingTabContent(db, logger, pulje)
}
```

> Tab labels use the pulje id (e.g. `FredagKveld`), matching the rooms-assignment tabs which render the raw pulje value. If readable names are wanted later, swap to `models.PuljeRow.Name`; out of scope here.

- [ ] **Step 2: Add the route setup function**

Append to `puljefordeling_tab.templ`:

```go
func SetupPuljefordelingTabRoute(router chi.Router, db *sql.DB, liveManager *live.Manager, logger *slog.Logger) {
	logger = logger.With("component", "puljefordeling_tab")

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, fmt.Sprintf("/admin/puljefordeling/%s", models.PuljeFredagKveld), http.StatusSeeOther)
	})

	router.Get("/{pulje}", func(w http.ResponseWriter, r *http.Request) {
		userInfo := userctx.GetUserRequestInfo(r.Context())
		puljeQuery := chi.URLParam(r, "pulje")
		puljeID, ok := models.ParsePulje(puljeQuery)
		if !ok {
			http.Redirect(w, r, fmt.Sprintf("/admin/puljefordeling/%s", models.PuljeFredagKveld), http.StatusSeeOther)
			return
		}
		if err := layouts.Base(
			"Puljefordeling",
			userInfo,
			puljefordelingTabIndex(PuljefordelingTabPage(db, logger, puljeID), puljeID),
		).Render(r.Context(), w); err != nil {
			logger.Error(fmt.Errorf("render puljefordeling tab: %w", err).Error(), "user_id", userInfo.Id)
		}
	})

	router.Get("/api/{pulje}", func(w http.ResponseWriter, r *http.Request) {
		puljeQuery := chi.URLParam(r, "pulje")
		puljeID, ok := models.ParsePulje(puljeQuery)
		if !ok {
			http.Error(w, "Expected a valid pulje ID, got: "+puljeQuery, http.StatusBadRequest)
			return
		}
		liveManager.Stream(w, r, live.Page{
			Buckets: []live.Bucket{live.BucketInterests, live.BucketEvents},
			Render: func(ctx context.Context, r *http.Request) templ.Component {
				return PuljefordelingTabContent(db, logger, puljeID)
			},
		})
	})
}
```

> `models.ParsePulje(s string) (models.Pulje, bool)` returns a `models.Pulje` directly — no conversion needed, `puljeID` is already a `models.Pulje`.

- [ ] **Step 3: Wire routes into `admin.go` and remove the emulate route**

In `pages/admin/admin.go`:
- Remove the import line `"github.com/Regncon/conorganizer/pages/admin/puljefordeling_emulate"`.
- Replace the emulate route block (lines ~34-36):

```go
		adminRouter.Route("/puljefordeling-emulate", func(emulateRouter chi.Router) {
			puljefordeling_emulate.SetupPuljefordelingEmulateRoute(emulateRouter, db, baseLogger)
		})
```

with:

```go
		adminRouter.Route("/puljefordeling", func(pfRouter chi.Router) {
			SetupPuljefordelingTabRoute(pfRouter, db, liveManager, logger)
		})
```

(If `baseLogger` becomes unused after this, leave it — it is used elsewhere in the function; verify with the compiler.)

- [ ] **Step 4: Delete the emulate package**

```bash
git rm -r pages/admin/puljefordeling_emulate/
```

- [ ] **Step 5: Generate, build, and run the whole admin/service test set**

Run:
```bash
go tool templ generate -path pages/admin -log-level error
go build ./...
go test ./pages/admin/... ./service/puljefordeling/...
```
Expected: build succeeds (no references to the deleted package remain), tests PASS.

- [ ] **Step 6: Commit**

```bash
git add -A
git commit -m "feat(puljefordeling): tabbed per-pulje page with live SSE; remove emulate page"
```

---

## Task 4: Consolidate the admin dashboard entry point

**Files:**
- Modify: `pages/admin/admin_page.templ`
- Modify: `pages/admin/puljefordeling.templ` (delete unused `puljefordeling` + `puljeStatusCard` templ funcs)

**Interfaces:**
- Consumes: existing `adminCard` templ helper in `admin_page.templ`.
- Produces: nothing new; removes the inline status grid and the emulate card, adds one link card.

- [ ] **Step 1: Replace the inline status section with a single card**

In `pages/admin/admin_page.templ`, remove this block (lines ~64-66):

```go
					<section class="admin-panel-card">
						@puljefordeling(db)
					</section>
```

and remove the "Emulér puljefordeling" card block (lines ~75-82):

```go
					@adminCard(
						"Emulér puljefordeling",
						"Forhåndsvis hvordan deltakere ville blitt fordelt på arrangementer i puljene, basert på påmeldte ønsker. Ingenting lagres.",
						"/static/call-to-action-avatar.webp",
						"Emulate seatings",
					) {
						<a role="button" href="/admin/puljefordeling-emulate/" class="btn btn--outline">Emulér fordeling</a>
					}
```

Add one card in their place (where the status section was):

```go
					@adminCard(
						"Puljefordeling",
						"Se og styr fordelingen av deltakere på arrangementer per pulje. Lås og publiser fordelingen når den er klar.",
						"/static/call-to-action-avatar.webp",
						"Puljefordeling",
					) {
						<a role="button" href="/admin/puljefordeling/" class="btn btn--outline">Gå til puljefordeling</a>
					}
```

- [ ] **Step 2: Delete the now-unused templ funcs**

In `pages/admin/puljefordeling.templ`, delete the entire `templ puljefordeling(db *sql.DB) { ... }` function (the one starting around line 179, including its `<style>` block) and the entire `templ puljeStatusCard(pulje models.PuljeRow) { ... }` function. Keep `getPuljer`, `puljefordelingStatusRoute`, `isValidPuljeStatus`, `puljeIsLocked`, `puljeIsCompleted`, `updatePuljeStatus`, and `puljeStatusUpdateAction` (still used by Task 2's controls).

- [ ] **Step 3: Generate and build**

Run:
```bash
go tool templ generate -path pages/admin -log-level error
go build ./...
```
Expected: builds clean. If the compiler reports `getPuljer` or any kept helper is now unused, leave `getPuljer` only if still referenced; otherwise delete it too and re-run.

- [ ] **Step 4: Commit**

```bash
git add -A
git commit -m "refactor(admin): collapse puljefordeling status + emulate cards into one tab entry"
```

---

## Task 5: Verification gate and manual smoke test

**Files:** none (verification only)

- [ ] **Step 1: Full test suite**

Run: `go test ./...`
Expected: all packages PASS. Fix any test that still references the removed emulate route/package or the deleted templ funcs.

- [ ] **Step 2: Linter**

Run: `golangci-lint run`
Expected: no findings. Address any unused-symbol / error-check findings introduced by the changes.

- [ ] **Step 3: Manual smoke test**

Run the dev server and verify in a browser:

```bash
go tool task start
```

- Open `http://localhost:7331/admin/puljefordeling/` → redirects to the first pulje tab.
- The `.tabs` bar shows one tab per pulje; clicking a tab navigates and marks it active.
- Each tab shows that pulje's event grid, stats, legend, and the lock/publish controls.
- Toggling "Puljefordeling lukket" prompts a confirm, commits, and the view re-renders to the saved distribution; unlocking reverts.
- The dashboard at `/admin/` shows a single "Puljefordeling" card linking here; the old `/admin/puljefordeling-emulate/` URL is gone (404).

- [ ] **Step 4: Final commit (if any smoke-test fixes were needed)**

```bash
git add -A
git commit -m "fix(puljefordeling): smoke-test adjustments"
```

---

## Self-Review Notes

- **Spec coverage:** routing/entry-point (Tasks 3-4), per-pulje tab content + status control (Task 2), scoped emulation / "don't emulate the rest until locked" (Task 1), live SSE on `BucketInterests`+`BucketEvents` (Task 3), removal of old emulate page (Task 3), dashboard collapse (Task 4), verification gate incl. `golangci-lint` (Task 5). The spec's "relocate status PUT" was intentionally simplified to "leave it in place" — see Global Constraints — to avoid touching the working markup.
- **Type consistency:** `EmulatePulje` returns `(EmulatedPulje, EmulationMeta, error)` (Task 1) and is consumed exactly that way in Task 2. `PuljefordelingTabContent(db, logger, pulje)` signature is identical in Task 2 (definition + test) and Task 3 (shell + stream).
- **Resolved during planning (no placeholders remain):** render helper is `testutil/templtest.Render` (returns `*goquery.Document`); `models.ParsePulje(s) (models.Pulje, bool)` returns a `models.Pulje` directly; seed-helper SQL is copied verbatim from `service/puljefordeling/emulate_test.go`. The one remaining judgement call left to the implementer is whether `getPuljer` becomes unused after Task 4 (delete it if the compiler flags it) — Task 4 Step 3 covers this.
