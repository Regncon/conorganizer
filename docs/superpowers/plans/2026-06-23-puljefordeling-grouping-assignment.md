# Puljefordeling Grouping Assignment Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Turn the read-only puljefordeling min-cost-flow preview into a real, persisted grouping assignment integrated into the `Open → Locked → Completed` pulje lifecycle.

**Architecture:** The solver (`solver/flow.go`, `solver/solver.go`) stays as the assignment engine; we add (1) pinned-seat support so manual placements are honored, (2) an `ApplyActual` replay so each frozen pulje seeds the next pulje's fairness from real seats, (3) a `source` column on `relation_events_players` to separate solver-written from hand-placed seats, and (4) orchestration that commits on Lock / reverts on Unlock plus a re-run button on the preview page.

**Tech Stack:** Go, SQLite (`modernc.org/sqlite`, STRICT tables), goose migrations, templ v0.3.x, chi router, datastar.

## Global Constraints

- `solver/flow.go` is NOT modified.
- Tests load the schema from `schema.sql` via `testutil.CreateTestDBAndLogger`; production uses goose migrations; `initialize.sql` is the fresh-bootstrap copy. A schema column change must touch **all three**.
- After editing any `.templ` file, run `go tool templ generate` before building. Do NOT hand-edit generated `*_templ.go` files for templ markup.
- Pulje statuses are exactly `Open`, `Locked`, `Completed` (`models.PuljeStatus`). "Frozen" = `Locked` OR `Completed`.
- `relation_events_players` role values are exactly `Player`, `GM` (`models.EventPlayerRole`). Solver writes `role='Player', source='solver'`.
- The build must stay green after every task.
- Commit after every task.

---

### Task 1: Add `source` column to `relation_events_players`

**Files:**
- Modify: `schema.sql` (the `relation_events_players` table, ~line 68)
- Modify: `initialize.sql` (the `relation_events_players` table, ~line 150)
- Create: `migrations/20260623120000_add_source_to_relation_events_players.sql`
- Test: `service/puljefordeling/schema_source_test.go`

**Interfaces:**
- Produces: a `source TEXT NOT NULL DEFAULT 'manual' CHECK (source IN ('manual','solver'))` column on `relation_events_players`.

- [ ] **Step 1: Write the failing test**

Create `service/puljefordeling/schema_source_test.go`:

```go
package puljefordeling

import (
	"testing"

	"github.com/Regncon/conorganizer/testutil"
)

func TestRelationEventsPlayersHasSourceColumn(t *testing.T) {
	db, _ := testutil.CreateTestDBAndLogger(t, "test_schema_source")

	rows, err := db.Query(`PRAGMA table_info(relation_events_players)`)
	if err != nil {
		t.Fatalf("pragma table_info: %v", err)
	}
	defer rows.Close()

	found := false
	for rows.Next() {
		var (
			cid        int
			name       string
			ctype      string
			notnull    int
			dflt       any
			primaryKey int
		)
		if err := rows.Scan(&cid, &name, &ctype, &notnull, &dflt, &primaryKey); err != nil {
			t.Fatalf("scan column: %v", err)
		}
		if name == "source" {
			found = true
		}
	}
	if !found {
		t.Error("relation_events_players is missing the 'source' column")
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./service/puljefordeling/ -run TestRelationEventsPlayersHasSourceColumn -v`
Expected: FAIL — "relation_events_players is missing the 'source' column".

- [ ] **Step 3: Add the column to `schema.sql`**

In `schema.sql`, inside `CREATE TABLE relation_events_players(...)`, add a line immediately after the `role` column:

```sql
    role TEXT NOT NULL DEFAULT 'Player' CHECK (role IN ('Player', 'GM')),
    source TEXT NOT NULL DEFAULT 'manual' CHECK (source IN ('manual','solver')),
```

(Read the file first to match its exact `role` line, then insert the `source` line right after it.)

- [ ] **Step 4: Add the column to `initialize.sql`**

Make the identical insertion in `initialize.sql`'s `relation_events_players` definition (after its `role` column line).

- [ ] **Step 5: Create the goose migration**

Create `migrations/20260623120000_add_source_to_relation_events_players.sql`:

```sql
-- +goose Up
ALTER TABLE relation_events_players
    ADD COLUMN source TEXT NOT NULL DEFAULT 'manual' CHECK (source IN ('manual','solver'));

-- +goose Down
ALTER TABLE relation_events_players DROP COLUMN source;
```

- [ ] **Step 6: Run test to verify it passes**

Run: `go test ./service/puljefordeling/ -run TestRelationEventsPlayersHasSourceColumn -v`
Expected: PASS.

- [ ] **Step 7: Commit**

```bash
git add schema.sql initialize.sql migrations/20260623120000_add_source_to_relation_events_players.sql service/puljefordeling/schema_source_test.go
git commit -m "feat(puljefordeling): add source column to relation_events_players"
```

---

### Task 2: Solver — pinned-seat support + `applyResult` extraction

**Files:**
- Modify: `service/puljefordeling/solver/solver.go`
- Test: `service/puljefordeling/solver/solver_test.go`

**Interfaces:**
- Consumes: `model.Slot`, `model.Player`, `model.Event`, `model.SlotResult`, `model.MaxScore`.
- Produces:
  - `func (s *State) SolveSlot(slot model.Slot, players []model.Player) model.SlotResult` — unchanged signature, now delegates to `SolveSlotFixed(slot, players, nil)`.
  - `func (s *State) SolveSlotFixed(slot model.Slot, players []model.Player, fixed map[string]string) model.SlotResult` — `fixed` maps `playerID → eventID` for pinned manual placements.
  - `func (s *State) applyResult(result *model.SlotResult, slotID string, interested []model.Player, playerByID map[string]model.Player, assignments map[string][]string) map[string]bool` (private).

- [ ] **Step 1: Write the failing tests**

Append to `service/puljefordeling/solver/solver_test.go`:

```go
func TestSolveSlotFixed_PinnedSeatHonoredWithoutPreference(t *testing.T) {
	// "kid" has NO preference for A, but is manually pinned there. The pin must be
	// honored, and the seat must reduce A's effective capacity.
	sl := slot("s1", event("A", 2))
	players := []model.Player{
		player("alice", prefs("s1", map[string]model.Score{"A": 5})),
		player("bob", prefs("s1", map[string]model.Score{"A": 5})),
		// "kid" intentionally absent from players: a manual placement with no interest.
	}

	result := NewState(2026, weekendOf(sl)).SolveSlotFixed(sl, players, map[string]string{"kid": "A"})

	if !slices.Contains(assigned(result, "A"), "kid") {
		t.Errorf("pinned kid must be seated in A, got %v", assigned(result, "A"))
	}
	if len(assigned(result, "A")) != 2 {
		t.Errorf("A cap 2: pin takes one seat, solver fills one more, want 2 total, got %d", len(assigned(result, "A")))
	}
	// Exactly one of alice/bob got the remaining seat; the other is unassigned.
	if len(result.Unassigned) != 1 {
		t.Errorf("want 1 unassigned (only 1 free seat after the pin), got %v", result.Unassigned)
	}
}

func TestSolveSlotFixed_PinnedSeatCountsAsSeated(t *testing.T) {
	// A pinned player who got their top choice is satisfied; one with no/low
	// interest is seated but not satisfied.
	sl := slot("s1", event("A", 4))
	players := []model.Player{
		player("top", prefs("s1", map[string]model.Score{"A": 5})),
	}

	st := NewState(2026, weekendOf(sl))
	result := st.SolveSlotFixed(sl, players, map[string]string{"top": "A", "nopref": "A"})

	if !st.IsSatisfied("top") {
		t.Error("pinned player whose pinned event is their top choice should be satisfied")
	}
	if st.IsSatisfied("nopref") {
		t.Error("pinned player with no top-choice preference must not be satisfied")
	}
	if !slices.Contains(result.NewlySatisfied, "top") {
		t.Errorf("top should be newly satisfied, got %v", result.NewlySatisfied)
	}
}

func TestSolveSlot_NilFixedRegression(t *testing.T) {
	// The 2-arg SolveSlot must behave exactly as before (delegates with nil).
	sl := slot("s1", event("A", 1))
	players := []model.Player{
		player("low", prefs("s1", map[string]model.Score{"A": 3})),
		player("high", prefs("s1", map[string]model.Score{"A": 5})),
	}
	result := NewState(2026, weekendOf(sl)).SolveSlot(sl, players)
	if !slices.Contains(assigned(result, "A"), "high") {
		t.Error("high scorer should still win with nil fixed")
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./service/puljefordeling/solver/ -run 'TestSolveSlotFixed|TestSolveSlot_NilFixedRegression' -v`
Expected: FAIL — `SolveSlotFixed` undefined (compile error).

- [ ] **Step 3: Extract `applyResult` and refactor `SolveSlot`**

In `service/puljefordeling/solver/solver.go`, replace the existing `SolveSlot` method (currently `solver.go:147-259`) with the following two methods plus the shared helper:

```go
// SolveSlot assigns players to events for one slot with no pinned placements.
func (s *State) SolveSlot(slot model.Slot, players []model.Player) model.SlotResult {
	return s.SolveSlotFixed(slot, players, nil)
}

// SolveSlotFixed assigns players to events for one slot, honoring pinned manual
// placements (fixed maps playerID → eventID), updates the fairness state, and
// returns the result.
//
// A pinned player is reserved into their event (consuming a seat and reducing the
// event's effective capacity) and is removed from the free assignment pool. Pins
// are honored even when the player expressed no interest in that event. Players
// DMing any event in this slot are excluded from the player pool.
func (s *State) SolveSlotFixed(slot model.Slot, players []model.Player, fixed map[string]string) model.SlotResult {
	currentIndex := s.slotIndex
	seed := int64(s.year)*1000 + int64(currentIndex)
	s.slotIndex++

	result := model.SlotResult{
		SlotID:      slot.ID,
		Assignments: make(map[string][]string),
		Seed:        seed,
	}

	// Players DMing in this slot are unavailable as players.
	dmingHere := make(map[string]bool)
	for _, ev := range slot.Events {
		if ev.DMID != "" {
			dmingHere[ev.DMID] = true
		}
	}

	// Validate pins: keep only those whose event is in this slot and whose player
	// is not DMing here.
	eventInSlot := make(map[string]bool, len(slot.Events))
	for _, ev := range slot.Events {
		eventInSlot[ev.ID] = true
	}
	pinnedByEvent := make(map[string][]string)
	pinned := make(map[string]bool)
	for pid, evID := range fixed {
		if !eventInSlot[evID] || dmingHere[pid] {
			continue
		}
		pinnedByEvent[evID] = append(pinnedByEvent[evID], pid)
		pinned[pid] = true
	}

	// Free pool: interested players who are neither DMing here nor pinned.
	interested := make([]model.Player, 0, len(players))
	for _, p := range players {
		if dmingHere[p.ID] || pinned[p.ID] {
			continue
		}
		if len(p.Prefs[slot.ID]) > 0 {
			interested = append(interested, p)
		}
	}

	// Index ALL players so pinned players' prefs are available for satisfaction.
	playerByID := make(map[string]model.Player, len(players))
	for _, p := range players {
		playerByID[p.ID] = p
	}

	// Reduce each event's capacity by the number of players pinned into it.
	events := make([]model.Event, len(slot.Events))
	copy(events, slot.Events)
	for i := range events {
		if n := len(pinnedByEvent[events[i].ID]); n > 0 {
			events[i].Capacity -= n
			if events[i].Capacity < 0 {
				events[i].Capacity = 0
			}
		}
	}

	// Solve the free pool over the reduced-capacity events.
	assignments := make(map[string][]string)
	var moved map[string]bool
	if len(interested) > 0 {
		rng := rand.New(rand.NewPCG(uint64(seed), 0)) //nolint:gosec
		rng.Shuffle(len(interested), func(i, j int) {
			interested[i], interested[j] = interested[j], interested[i]
		})
		assignments, moved = s.runMCMF(slot.ID, events, interested)
	}

	// Merge pinned placements into the assignment.
	for evID, pids := range pinnedByEvent {
		assignments[evID] = append(assignments[evID], pids...)
	}
	result.Assignments = assignments

	// Flag events with fewer than minViablePlayers (against final counts).
	for _, ev := range slot.Events {
		if len(assignments[ev.ID]) < minViablePlayers {
			result.UndersubscribedEvents = append(result.UndersubscribedEvents, ev.ID)
		}
	}

	// Update fairness/totals/misses/unassigned; returns the seated set.
	assigned := s.applyResult(&result, slot.ID, interested, playerByID, assignments)

	// Players bumped off a higher-scoring event by a residual augmentation and
	// still holding a seat.
	for pid := range moved {
		if assigned[pid] {
			result.MovedPlayers = append(result.MovedPlayers, pid)
		}
	}

	sortSlotResult(&result)
	return result
}

// applyResult updates fairness state (seated/satisfied/misses) from a final
// assignment, fills the result's per-assignment fields (TotalScore, NewlySatisfied,
// Unassigned), and returns the set of seated player IDs. assignments is
// eventID → []playerID. Scores for players absent from playerByID are treated as 0.
func (s *State) applyResult(
	result *model.SlotResult,
	slotID string,
	interested []model.Player,
	playerByID map[string]model.Player,
	assignments map[string][]string,
) map[string]bool {
	assigned := make(map[string]bool)
	for evID, playerIDs := range assignments {
		for _, pid := range playerIDs {
			score := playerByID[pid].Prefs[slotID][evID]
			result.TotalScore += int(score)
			s.seated[pid] = true
			if score == model.MaxScore && !s.satisfied[pid] {
				s.satisfied[pid] = true
				result.NewlySatisfied = append(result.NewlySatisfied, pid)
			}
			assigned[pid] = true
		}
	}

	// Record a miss for every still-unsatisfied free-pool player who wanted a top
	// choice this slot but did not get one.
	for _, p := range interested {
		if s.satisfied[p.ID] {
			continue
		}
		if wantedTopChoice(p, slotID) {
			s.misses[p.ID]++
		}
	}

	// Unassigned: interested (free-pool) players who got no seat.
	for _, p := range interested {
		if !assigned[p.ID] {
			result.Unassigned = append(result.Unassigned, p.ID)
		}
	}

	return assigned
}

// sortSlotResult sorts all output slices for deterministic results.
func sortSlotResult(result *model.SlotResult) {
	sort.Strings(result.NewlySatisfied)
	sort.Strings(result.Unassigned)
	sort.Strings(result.UndersubscribedEvents)
	sort.Strings(result.MovedPlayers)
	for evID := range result.Assignments {
		sort.Strings(result.Assignments[evID])
	}
}
```

(Leave `runMCMF`, `adjustScore`, `bandBase`, `wantedTopChoice`, and all constants unchanged.)

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./service/puljefordeling/solver/ -v`
Expected: PASS — the new tests plus all pre-existing solver tests (they call the unchanged 2-arg `SolveSlot`).

- [ ] **Step 5: Commit**

```bash
git add service/puljefordeling/solver/solver.go service/puljefordeling/solver/solver_test.go
git commit -m "feat(puljefordeling): pinned-seat support in solver (SolveSlotFixed)"
```

---

### Task 3: Solver — `ApplyActual` fairness replay

**Files:**
- Modify: `service/puljefordeling/solver/solver.go`
- Test: `service/puljefordeling/solver/solver_test.go`

**Interfaces:**
- Consumes: `applyResult`, `sortSlotResult`, `wantedTopChoice` (from Task 2).
- Produces: `func (s *State) ApplyActual(slot model.Slot, players []model.Player, assignments map[string][]string) model.SlotResult` — seeds fairness from a known persisted assignment without solving.

- [ ] **Step 1: Write the failing test**

Append to `service/puljefordeling/solver/solver_test.go`:

```go
func TestApplyActual_SeedsFairnessFromRealSeats(t *testing.T) {
	// Replay slot 1 where alice actually got her top choice (A) and bob actually
	// missed his top choice. Then in slot 2 (same single seat) bob — now carrying a
	// miss and never satisfied — should beat the already-satisfied alice.
	sl1 := slot("s1", event("A", 1))
	sl2 := slot("s2", event("A", 1))

	alice := model.Player{ID: "alice", Name: "alice", Prefs: map[string]map[string]model.Score{
		"s1": {"A": 5}, "s2": {"A": 5},
	}}
	bob := model.Player{ID: "bob", Name: "bob", Prefs: map[string]map[string]model.Score{
		"s1": {"A": 5}, "s2": {"A": 5},
	}}

	st := NewState(2026, weekendOf(sl1, sl2))

	// Replay the actual slot-1 result: alice seated in A, bob unseated.
	r1 := st.ApplyActual(sl1, []model.Player{alice, bob}, map[string][]string{"A": {"alice"}})
	if !slices.Contains(r1.Assignments["A"], "alice") {
		t.Fatalf("ApplyActual should echo the actual assignment, got %v", r1.Assignments["A"])
	}
	if !st.IsSatisfied("alice") {
		t.Error("alice got her top choice in the replayed slot → satisfied")
	}
	if st.IsSatisfied("bob") {
		t.Error("bob was not seated → not satisfied")
	}

	// Now solve slot 2: bob (unsatisfied + a miss) should win over satisfied alice.
	r2 := st.SolveSlot(sl2, []model.Player{alice, bob})
	if !slices.Contains(assigned(r2, "A"), "bob") {
		t.Errorf("bob (missed + unsatisfied) should win slot 2 over satisfied alice, got %v", assigned(r2, "A"))
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./service/puljefordeling/solver/ -run TestApplyActual_SeedsFairnessFromRealSeats -v`
Expected: FAIL — `ApplyActual` undefined.

- [ ] **Step 3: Implement `ApplyActual`**

Add to `service/puljefordeling/solver/solver.go`:

```go
// ApplyActual seeds the fairness state from a known persisted assignment for a slot
// without running the solver. It advances the slot index (so subsequent seeds stay
// aligned), echoes the assignment into the returned result, and updates
// seated/satisfied/misses exactly as a solved slot would. Used to replay frozen
// puljer so later puljer are seeded from what actually happened.
func (s *State) ApplyActual(slot model.Slot, players []model.Player, assignments map[string][]string) model.SlotResult {
	currentIndex := s.slotIndex
	seed := int64(s.year)*1000 + int64(currentIndex)
	s.slotIndex++

	// Copy the assignment so sorting/appends never mutate the caller's map.
	clone := make(map[string][]string, len(assignments))
	for evID, pids := range assignments {
		clone[evID] = append([]string(nil), pids...)
	}

	result := model.SlotResult{
		SlotID:      slot.ID,
		Assignments: clone,
		Seed:        seed,
	}

	dmingHere := make(map[string]bool)
	for _, ev := range slot.Events {
		if ev.DMID != "" {
			dmingHere[ev.DMID] = true
		}
	}

	interested := make([]model.Player, 0, len(players))
	for _, p := range players {
		if dmingHere[p.ID] {
			continue
		}
		if len(p.Prefs[slot.ID]) > 0 {
			interested = append(interested, p)
		}
	}

	playerByID := make(map[string]model.Player, len(players))
	for _, p := range players {
		playerByID[p.ID] = p
	}

	s.applyResult(&result, slot.ID, interested, playerByID, result.Assignments)

	for _, ev := range slot.Events {
		if len(result.Assignments[ev.ID]) < minViablePlayers {
			result.UndersubscribedEvents = append(result.UndersubscribedEvents, ev.ID)
		}
	}

	sortSlotResult(&result)
	return result
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./service/puljefordeling/solver/ -v`
Expected: PASS (all solver tests).

- [ ] **Step 5: Commit**

```bash
git add service/puljefordeling/solver/solver.go service/puljefordeling/solver/solver_test.go
git commit -m "feat(puljefordeling): ApplyActual fairness replay from persisted seats"
```

---

### Task 4: Orchestration — seeded chronological run + manual-seat marker

**Files:**
- Modify: `service/puljefordeling/emulate.go`
- Test: `service/puljefordeling/emulate_test.go`

**Interfaces:**
- Consumes: `solver.NewState`, `solver.State.SolveSlotFixed`, `solver.State.ApplyActual`, `solver.State.SatisfiedCount`.
- Produces (used by Task 5):
  - `type seatingData struct { ... }` with fields `puljer []models.PuljeRow`, `weekend smodel.Weekend`, `players []smodel.Player`, `gms map[string]int`, `names map[int]string`, `prefs map[int]map[string]map[string]smodel.Score`, `dmSet map[int]bool`, `actual map[string]map[string][]string`, `manualFixed map[string]map[string]string`, `pinnedSet map[string]map[string]bool`, `year int`.
  - `func loadSeatingData(db *sql.DB) (*seatingData, error)`.
  - `func (d *seatingData) solveChronological(upTo int) (*solver.State, []smodel.SlotResult)`.
  - `func puljeFrozen(status models.PuljeStatus) bool`.
  - `AssignedPlayer` gains a `Pinned bool` field.

- [ ] **Step 1: Write the failing test**

Append to `service/puljefordeling/emulate_test.go`:

```go
func seedAssignment(t *testing.T, db *sql.DB, eventID string, pulje models.Pulje, bhID int, source string) {
	t.Helper()
	_, err := db.Exec(
		`INSERT INTO relation_events_players (event_id, pulje_id, billettholder_id, role, source) VALUES (?, ?, ?, 'Player', ?)`,
		eventID, string(pulje), bhID, source,
	)
	if err != nil {
		t.Fatalf("seed assignment bh=%d ev=%s: %v", bhID, eventID, err)
	}
}

func setPuljeStatus(t *testing.T, db *sql.DB, pulje models.Pulje, status models.PuljeStatus) {
	t.Helper()
	if _, err := db.Exec(`UPDATE puljer SET status = ? WHERE id = ?`, string(status), string(pulje)); err != nil {
		t.Fatalf("set pulje %s status: %v", pulje, err)
	}
}

func TestEmulateSeatings_ManualPlacementPinnedAndMarked(t *testing.T) {
	db, _ := testutil.CreateTestDBAndLogger(t, "test_emulate_pinned")

	const fredag = models.PuljeFredagKveld
	seedPulje(t, db, fredag, "Fredag Kveld", "2026-09-04T18:00:00Z")
	seedEvent(t, db, "evA", "Alpha", 1, fredag) // capacity 1

	seedParticipant(t, db, 1, "Anna", "A")
	seedParticipant(t, db, 2, "Kid", "K")

	// Anna wants evA highly. Kid has NO interest but is manually placed in evA.
	seedInterest(t, db, 1, "evA", fredag, models.InterestLevelHigh)
	seedAssignment(t, db, "evA", fredag, 2, "manual")

	em, err := EmulateSeatings(db)
	if err != nil {
		t.Fatalf("EmulateSeatings: %v", err)
	}

	evA, ok := findEvent(em.Puljer[0], "evA")
	if !ok {
		t.Fatal("evA missing")
	}
	names := playerNames(evA.AssignedPlayers)
	if !slices.Contains(names, "Kid K") {
		t.Errorf("manually placed Kid must be seated in evA, got %v", names)
	}
	// evA capacity is 1 and the pin took it, so Anna cannot also be seated there.
	if slices.Contains(names, "Anna A") {
		t.Errorf("Anna should not fit (pin took the only seat), got %v", names)
	}
	// The pinned seat must be marked.
	for _, ap := range evA.AssignedPlayers {
		if ap.Name == "Kid K" && !ap.Pinned {
			t.Error("Kid's seat should be marked Pinned")
		}
	}
}

func TestEmulateSeatings_FrozenPuljeSeedsNext(t *testing.T) {
	db, _ := testutil.CreateTestDBAndLogger(t, "test_emulate_seed")

	const p1 = models.PuljeFredagKveld
	const p2 = models.PuljeLordagMorgen
	seedPulje(t, db, p1, "Pulje 1", "2026-09-04T18:00:00Z")
	seedPulje(t, db, p2, "Pulje 2", "2026-09-05T09:00:00Z")

	seedEvent(t, db, "e1", "Event 1", 1, p1)
	seedEvent(t, db, "e2", "Event 2", 1, p2)

	seedParticipant(t, db, 1, "Alice", "A")
	seedParticipant(t, db, 2, "Bob", "B")

	// Both want their pulje-1 event (cap 1) and the pulje-2 event (cap 1).
	seedInterest(t, db, 1, "e1", p1, models.InterestLevelHigh)
	seedInterest(t, db, 2, "e1", p1, models.InterestLevelHigh)
	seedInterest(t, db, 1, "e2", p2, models.InterestLevelHigh)
	seedInterest(t, db, 2, "e2", p2, models.InterestLevelHigh)

	// Pulje 1 is frozen with Alice actually seated in e1 (she's now satisfied).
	setPuljeStatus(t, db, p1, models.PuljeStatusCompleted)
	seedAssignment(t, db, "e1", p1, 1, "solver")

	em, err := EmulateSeatings(db)
	if err != nil {
		t.Fatalf("EmulateSeatings: %v", err)
	}

	// In pulje 2, Bob (unsatisfied) should win e1's single seat over satisfied Alice.
	var p2res EmulatedPulje
	for _, p := range em.Puljer {
		if p.PuljeID == p2 {
			p2res = p
		}
	}
	e2, ok := findEvent(p2res, "e2")
	if !ok {
		t.Fatal("e2 missing")
	}
	if !slices.Contains(playerNames(e2.AssignedPlayers), "Bob B") {
		t.Errorf("unsatisfied Bob should win pulje-2 seat (Alice satisfied in frozen pulje 1), got %v", playerNames(e2.AssignedPlayers))
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./service/puljefordeling/ -run 'TestEmulateSeatings_ManualPlacementPinnedAndMarked|TestEmulateSeatings_FrozenPuljeSeedsNext' -v`
Expected: FAIL — `seedAssignment` references the `source` column (exists after Task 1) but `AssignedPlayer.Pinned` is undefined and seeding/pinning is not wired, so assertions fail/compile-fail.

- [ ] **Step 3: Add the `Pinned` field**

In `service/puljefordeling/emulate.go`, add to the `AssignedPlayer` struct:

```go
type AssignedPlayer struct {
	Name   string
	IsDM   bool                 // runs at least one game in the weekend (DM bump)
	Level  models.InterestLevel // their interest in the game they got
	Moved  bool                 // relocated off a higher-scoring event by the solver to make room for others
	Pinned bool                 // manually placed (source=manual); honored by the solver, not chosen by it
}
```

- [ ] **Step 4: Add the seating-data loader, frozen helper, and chronological runner**

Add to `service/puljefordeling/emulate.go` (a `persistedSeat` loader + assembly). Place near the other loaders:

```go
// puljeFrozen reports whether a pulje's seats are committed and no longer
// auto-solved (Locked or Completed).
func puljeFrozen(status models.PuljeStatus) bool {
	return status == models.PuljeStatusLocked || status == models.PuljeStatusCompleted
}

func seatKey(eventID, playerID string) string {
	return eventID + "\x00" + playerID
}

type persistedSeat struct {
	pulje    string
	event    string
	playerID string
	source   string
}

// loadPersistedPlayerSeats returns every role=Player row (manual and solver).
func loadPersistedPlayerSeats(db *sql.DB) ([]persistedSeat, error) {
	const query = `
		SELECT event_id, pulje_id, billettholder_id, source
		FROM relation_events_players
		WHERE role = ?
	`
	rows, err := db.Query(query, models.EventPlayerRolePlayer)
	if err != nil {
		return nil, fmt.Errorf("query player seats: %w", err)
	}
	defer rows.Close()

	var seats []persistedSeat
	for rows.Next() {
		var s persistedSeat
		var bhID int
		if err := rows.Scan(&s.event, &s.pulje, &bhID, &s.source); err != nil {
			return nil, fmt.Errorf("scan player seat: %w", err)
		}
		s.playerID = strconv.Itoa(bhID)
		seats = append(seats, s)
	}
	return seats, rows.Err()
}

// seatingData holds everything needed to run the seeded chronological solve.
type seatingData struct {
	puljer      []models.PuljeRow
	weekend     smodel.Weekend
	players     []smodel.Player
	gms         map[string]int
	names       map[int]string
	prefs       map[int]map[string]map[string]smodel.Score
	dmSet       map[int]bool
	actual      map[string]map[string][]string // puljeID -> eventID -> []playerID (all Player rows)
	manualFixed map[string]map[string]string   // puljeID -> playerID -> eventID (source=manual)
	pinnedSet   map[string]map[string]bool      // puljeID -> seatKey(event,player) (source=manual)
	year        int
}

// loadSeatingData loads puljer, events, GMs, names, prefs, and persisted Player
// seats, and assembles the solver model. Players include every participant with an
// interest plus anyone holding a manual placement (even without an interest).
func loadSeatingData(db *sql.DB) (*seatingData, error) {
	puljer, err := loadPuljer(db)
	if err != nil {
		return nil, err
	}
	if len(puljer) == 0 {
		return &seatingData{}, nil
	}
	events, err := loadEligibleEvents(db)
	if err != nil {
		return nil, err
	}
	gms, err := loadGMs(db)
	if err != nil {
		return nil, err
	}
	names, err := loadParticipantNames(db)
	if err != nil {
		return nil, err
	}
	prefs, err := loadPrefs(db, events)
	if err != nil {
		return nil, err
	}
	seats, err := loadPersistedPlayerSeats(db)
	if err != nil {
		return nil, err
	}

	d := &seatingData{
		puljer:      puljer,
		gms:         gms,
		names:       names,
		prefs:       prefs,
		actual:      make(map[string]map[string][]string),
		manualFixed: make(map[string]map[string]string),
		pinnedSet:   make(map[string]map[string]bool),
	}

	for _, s := range seats {
		if d.actual[s.pulje] == nil {
			d.actual[s.pulje] = make(map[string][]string)
		}
		d.actual[s.pulje][s.event] = append(d.actual[s.pulje][s.event], s.playerID)
		if s.source == "manual" {
			if d.manualFixed[s.pulje] == nil {
				d.manualFixed[s.pulje] = make(map[string]string)
			}
			d.manualFixed[s.pulje][s.playerID] = s.event
			if d.pinnedSet[s.pulje] == nil {
				d.pinnedSet[s.pulje] = make(map[string]bool)
			}
			d.pinnedSet[s.pulje][seatKey(s.event, s.playerID)] = true
		}
	}

	// Build the solver's Weekend in chronological pulje order.
	d.weekend = smodel.Weekend{Slots: make([]smodel.Slot, 0, len(puljer))}
	for _, p := range puljer {
		slot := smodel.Slot{ID: string(p.ID), Name: p.Name}
		for _, eid := range sortedEventIDs(events[p.ID]) {
			e := events[p.ID][eid]
			ev := smodel.Event{ID: eid, Name: e.title, Capacity: e.capacity}
			if gmID, ok := gms[eventPuljeKey(eid, p.ID)]; ok {
				ev.DMID = strconv.Itoa(gmID)
			}
			slot.Events = append(slot.Events, ev)
		}
		d.weekend.Slots = append(d.weekend.Slots, slot)
	}

	// Players: everyone with an interest, plus anyone holding a manual seat.
	playerIDs := make(map[int]bool, len(prefs))
	for bh := range prefs {
		playerIDs[bh] = true
	}
	for _, s := range seats {
		if s.source != "manual" {
			continue
		}
		if bh, err := strconv.Atoi(s.playerID); err == nil {
			playerIDs[bh] = true
		}
	}
	for _, bh := range sortedIntKeys(playerIDs) { // sortedIntKeys is generic over map[int]V
		d.players = append(d.players, smodel.Player{
			ID:    strconv.Itoa(bh),
			Name:  names[bh],
			Prefs: prefs[bh],
		})
	}

	d.dmSet = make(map[int]bool, len(gms))
	for _, bhID := range gms {
		d.dmSet[bhID] = true
	}

	d.year = puljer[0].StartAt.TimeOrZero().Year()
	return d, nil
}

// solveChronological threads one State across slots[0..upTo] inclusive: frozen
// puljer are replayed from their persisted seats (ApplyActual); open puljer are
// solved with their manual placements pinned (SolveSlotFixed). Returns the State
// and per-slot results index-aligned with d.puljer[0..upTo].
func (d *seatingData) solveChronological(upTo int) (*solver.State, []smodel.SlotResult) {
	state := solver.NewState(d.year, d.weekend)
	results := make([]smodel.SlotResult, 0, upTo+1)
	for i := 0; i <= upTo; i++ {
		slot := d.weekend.Slots[i]
		pid := string(d.puljer[i].ID)
		var res smodel.SlotResult
		if puljeFrozen(d.puljer[i].Status) {
			res = state.ApplyActual(slot, d.players, d.actual[pid])
		} else {
			res = state.SolveSlotFixed(slot, d.players, d.manualFixed[pid])
		}
		results = append(results, res)
	}
	return state, results
}
```

- [ ] **Step 5: Rewrite `EmulateSeatings` to use the seeded run**

Replace the body of `EmulateSeatings` (`emulate.go:63-132`) with:

```go
func EmulateSeatings(db *sql.DB) (Emulation, error) {
	d, err := loadSeatingData(db)
	if err != nil {
		return Emulation{}, err
	}
	if len(d.puljer) == 0 {
		return Emulation{}, nil
	}

	state, results := d.solveChronological(len(d.puljer) - 1)

	// PlayerCount is distinct participants with at least one interest (unchanged
	// semantics — excludes manual-only placements with no interest).
	emulation := Emulation{Year: d.year, PlayerCount: len(d.prefs)}
	for i := range d.puljer {
		pid := string(d.puljer[i].ID)
		emulation.Puljer = append(emulation.Puljer, shapePulje(
			d.puljer[i], d.weekend.Slots[i], results[i],
			d.gms, d.names, d.prefs, d.dmSet, d.pinnedSet[pid],
		))
	}
	emulation.SatisfiedTotal = state.SatisfiedCount()
	return emulation, nil
}
```

- [ ] **Step 6: Thread the pinned set through shaping**

Update `shapePulje` to accept `pinned map[string]bool` and pass it to `assignedPlayers`; update `assignedPlayers` to accept `pinned map[string]bool` and set `ap.Pinned`. Change the two signatures and call sites:

In `shapePulje`, add the parameter and forward it:

```go
func shapePulje(
	pulje models.PuljeRow,
	slot smodel.Slot,
	res smodel.SlotResult,
	gms map[string]int,
	names map[int]string,
	prefs map[int]map[string]map[string]smodel.Score,
	dmSet map[int]bool,
	pinned map[string]bool,
) EmulatedPulje {
```

and in its event loop change the `AssignedPlayers:` line to:

```go
			AssignedPlayers: assignedPlayers(res.Assignments[ev.ID], ev.ID, string(pulje.ID), names, prefs, dmSet, moved, pinned),
```

In `assignedPlayers`, add the parameter and set the flag:

```go
func assignedPlayers(
	ids []string,
	eventID, puljeID string,
	names map[int]string,
	prefs map[int]map[string]map[string]smodel.Score,
	dmSet map[int]bool,
	moved map[string]bool,
	pinned map[string]bool,
) []AssignedPlayer {
	if len(ids) == 0 {
		return nil
	}
	out := make([]AssignedPlayer, 0, len(ids))
	for _, id := range ids {
		bh, err := strconv.Atoi(id)
		if err != nil {
			out = append(out, AssignedPlayer{Name: id})
			continue
		}
		ap := AssignedPlayer{
			Name:   names[bh],
			IsDM:   dmSet[bh],
			Moved:  moved[id],
			Pinned: pinned[seatKey(eventID, id)],
		}
		if byPulje, ok := prefs[bh]; ok {
			got := byPulje[puljeID][eventID]
			ap.Level = models.InterestLevelFromScore(int(got))
		}
		out = append(out, ap)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}
```

- [ ] **Step 7: Run tests to verify they pass**

Run: `go test ./service/puljefordeling/ -v`
Expected: PASS — the two new tests plus the existing `TestEmulateSeatings` (the `pinned` map is empty there, so behavior is unchanged).

- [ ] **Step 8: Commit**

```bash
git add service/puljefordeling/emulate.go service/puljefordeling/emulate_test.go
git commit -m "feat(puljefordeling): seed preview from persisted seats + mark manual placements"
```

---

### Task 5: Commit & revert persisted assignments

**Files:**
- Create: `service/puljefordeling/commit.go`
- Test: `service/puljefordeling/commit_test.go`

**Interfaces:**
- Consumes: `loadSeatingData`, `seatingData.solveChronological`, `seatingData.manualFixed`, `models.EventPlayerRolePlayer`, `models.Pulje`.
- Produces:
  - `func CommitPuljeAssignments(db *sql.DB, target models.Pulje, logger *slog.Logger) error` — solves the target (seeded from frozen priors, manual pinned), and in one transaction replaces `source='solver'` Player rows for the pulje and sets its status to `Locked`.
  - `func RevertPuljeAssignments(db *sql.DB, target models.Pulje) error` — in one transaction deletes `source='solver'` Player rows for the pulje and sets its status to `Open`.

- [ ] **Step 1: Write the failing tests**

Create `service/puljefordeling/commit_test.go`:

```go
package puljefordeling

import (
	"slices"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil"
)

func TestCommitPuljeAssignments_WritesSolverSeatsPreservesManual(t *testing.T) {
	db, logger := testutil.CreateTestDBAndLogger(t, "test_commit")

	const fredag = models.PuljeFredagKveld
	seedPulje(t, db, fredag, "Fredag Kveld", "2026-09-04T18:00:00Z")
	seedEvent(t, db, "evA", "Alpha", 4, fredag)

	seedParticipant(t, db, 1, "Anna", "A")
	seedParticipant(t, db, 2, "Bob", "B")
	seedParticipant(t, db, 3, "Kid", "K")

	seedInterest(t, db, 1, "evA", fredag, models.InterestLevelHigh)
	seedInterest(t, db, 2, "evA", fredag, models.InterestLevelHigh)
	// Kid is manually pinned, no interest.
	seedAssignment(t, db, "evA", fredag, 3, "manual")

	if err := CommitPuljeAssignments(db, fredag, logger); err != nil {
		t.Fatalf("CommitPuljeAssignments: %v", err)
	}

	// Status flipped to Locked.
	var status string
	if err := db.QueryRow(`SELECT status FROM puljer WHERE id = ?`, string(fredag)).Scan(&status); err != nil {
		t.Fatalf("read status: %v", err)
	}
	if status != string(models.PuljeStatusLocked) {
		t.Errorf("status: want Locked, got %s", status)
	}

	// Anna & Bob written as solver; Kid preserved as manual.
	rows, err := db.Query(`SELECT billettholder_id, source FROM relation_events_players WHERE pulje_id = ? AND role = 'Player' ORDER BY billettholder_id`, string(fredag))
	if err != nil {
		t.Fatalf("query seats: %v", err)
	}
	defer rows.Close()
	type seat struct {
		bh     int
		source string
	}
	var seats []seat
	for rows.Next() {
		var s seat
		if err := rows.Scan(&s.bh, &s.source); err != nil {
			t.Fatalf("scan: %v", err)
		}
		seats = append(seats, s)
	}
	if len(seats) != 3 {
		t.Fatalf("want 3 Player rows, got %d (%v)", len(seats), seats)
	}
	bySource := map[int]string{}
	for _, s := range seats {
		bySource[s.bh] = s.source
	}
	if bySource[1] != "solver" || bySource[2] != "solver" {
		t.Errorf("Anna/Bob should be solver-written, got %v", bySource)
	}
	if bySource[3] != "manual" {
		t.Errorf("Kid should remain manual, got %q", bySource[3])
	}
}

func TestRevertPuljeAssignments_RemovesOnlySolverSeats(t *testing.T) {
	db, logger := testutil.CreateTestDBAndLogger(t, "test_revert")

	const fredag = models.PuljeFredagKveld
	seedPulje(t, db, fredag, "Fredag Kveld", "2026-09-04T18:00:00Z")
	seedEvent(t, db, "evA", "Alpha", 4, fredag)
	seedParticipant(t, db, 1, "Anna", "A")
	seedParticipant(t, db, 2, "Kid", "K")
	seedInterest(t, db, 1, "evA", fredag, models.InterestLevelHigh)
	seedAssignment(t, db, "evA", fredag, 2, "manual")

	if err := CommitPuljeAssignments(db, fredag, logger); err != nil {
		t.Fatalf("commit: %v", err)
	}
	if err := RevertPuljeAssignments(db, fredag); err != nil {
		t.Fatalf("revert: %v", err)
	}

	var status string
	if err := db.QueryRow(`SELECT status FROM puljer WHERE id = ?`, string(fredag)).Scan(&status); err != nil {
		t.Fatalf("read status: %v", err)
	}
	if status != string(models.PuljeStatusOpen) {
		t.Errorf("status after revert: want Open, got %s", status)
	}

	var solverCount, manualCount int
	if err := db.QueryRow(`SELECT COUNT(*) FROM relation_events_players WHERE pulje_id = ? AND role='Player' AND source='solver'`, string(fredag)).Scan(&solverCount); err != nil {
		t.Fatalf("count solver: %v", err)
	}
	if err := db.QueryRow(`SELECT COUNT(*) FROM relation_events_players WHERE pulje_id = ? AND role='Player' AND source='manual'`, string(fredag)).Scan(&manualCount); err != nil {
		t.Fatalf("count manual: %v", err)
	}
	if solverCount != 0 {
		t.Errorf("solver seats should be gone, got %d", solverCount)
	}
	if manualCount != 1 {
		t.Errorf("manual seat must survive revert, got %d", manualCount)
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./service/puljefordeling/ -run 'TestCommitPuljeAssignments_WritesSolverSeatsPreservesManual|TestRevertPuljeAssignments_RemovesOnlySolverSeats' -v`
Expected: FAIL — `CommitPuljeAssignments`/`RevertPuljeAssignments` undefined.

- [ ] **Step 3: Implement commit & revert**

Create `service/puljefordeling/commit.go`:

```go
package puljefordeling

import (
	"database/sql"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/Regncon/conorganizer/models"
)

// CommitPuljeAssignments solves the target pulje — seeded by the actual seats of
// earlier frozen puljer and pinning the pulje's existing manual placements — and
// persists the solver-chosen seats. In one transaction it removes any prior
// source='solver' Player rows for the pulje, inserts the fresh ones, and sets the
// pulje status to Locked. Manual and GM rows are left untouched. Idempotent.
func CommitPuljeAssignments(db *sql.DB, target models.Pulje, logger *slog.Logger) error {
	d, err := loadSeatingData(db)
	if err != nil {
		return fmt.Errorf("load seating data: %w", err)
	}
	idx := -1
	for i := range d.puljer {
		if d.puljer[i].ID == target {
			idx = i
			break
		}
	}
	if idx == -1 {
		return fmt.Errorf("pulje %s not found", target)
	}

	_, results := d.solveChronological(idx)
	res := results[idx]
	manual := d.manualFixed[string(target)]

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("begin commit tx for %s: %w", target, err)
	}
	defer tx.Rollback()

	if _, err := tx.Exec(
		`DELETE FROM relation_events_players WHERE pulje_id = ? AND role = ? AND source = 'solver'`,
		string(target), models.EventPlayerRolePlayer,
	); err != nil {
		return fmt.Errorf("clear solver seats for %s: %w", target, err)
	}

	const insert = `
		INSERT INTO relation_events_players (event_id, pulje_id, billettholder_id, role, source)
		VALUES (?, ?, ?, ?, 'solver')
		ON CONFLICT(billettholder_id, event_id, pulje_id) DO NOTHING
	`
	inserted := 0
	for evID, pids := range res.Assignments {
		for _, pid := range pids {
			if _, isManual := manual[pid]; isManual {
				continue // already persisted as a manual row
			}
			bh, convErr := strconv.Atoi(pid)
			if convErr != nil {
				continue
			}
			if _, err := tx.Exec(insert, evID, string(target), bh, models.EventPlayerRolePlayer); err != nil {
				return fmt.Errorf("insert solver seat (pulje=%s event=%s bh=%d): %w", target, evID, bh, err)
			}
			inserted++
		}
	}

	if _, err := tx.Exec(`UPDATE puljer SET status = ? WHERE id = ?`, string(models.PuljeStatusLocked), string(target)); err != nil {
		return fmt.Errorf("lock pulje %s: %w", target, err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit seats for %s: %w", target, err)
	}
	if logger != nil {
		logger.Info("committed puljefordeling", "pulje_id", target, "solver_seats", inserted)
	}
	return nil
}

// RevertPuljeAssignments removes the solver-written seats for a pulje (leaving manual
// and GM rows intact) and reopens it, in one transaction.
func RevertPuljeAssignments(db *sql.DB, target models.Pulje) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("begin revert tx for %s: %w", target, err)
	}
	defer tx.Rollback()

	if _, err := tx.Exec(
		`DELETE FROM relation_events_players WHERE pulje_id = ? AND role = ? AND source = 'solver'`,
		string(target), models.EventPlayerRolePlayer,
	); err != nil {
		return fmt.Errorf("revert solver seats for %s: %w", target, err)
	}
	if _, err := tx.Exec(`UPDATE puljer SET status = ? WHERE id = ?`, string(models.PuljeStatusOpen), string(target)); err != nil {
		return fmt.Errorf("reopen pulje %s: %w", target, err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit revert for %s: %w", target, err)
	}
	return nil
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./service/puljefordeling/ -v`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add service/puljefordeling/commit.go service/puljefordeling/commit_test.go
git commit -m "feat(puljefordeling): commit/revert persisted assignments on lock/unlock"
```

---

### Task 6: Wire commit/revert into the pulje status handler

**Files:**
- Modify: `pages/admin/puljefordeling.templ` (the `puljefordelingStatusRoute` handler)
- Regenerate: `pages/admin/puljefordeling_templ.go` (via `go tool templ generate`)
- Test: `pages/admin/puljefordeling_status_test.go`

**Interfaces:**
- Consumes: `puljefordeling.CommitPuljeAssignments`, `puljefordeling.RevertPuljeAssignments`, existing `updatePuljeStatus`, `models.ParsePulje`.

- [ ] **Step 1: Write the failing test**

This drives the real HTTP handler through `httptest` (modeled on
`pages/admin/billettholder_admin/billettholder_email_routes_test.go:180-203`), using a
zero-value `&live.Manager{}` and asserting DB side-effects (not the broadcast). Create
`pages/admin/puljefordeling_status_test.go`:

```go
package admin

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/service/live"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/go-chi/chi/v5"
)

func putPuljeStatus(t *testing.T, router http.Handler, pulje models.Pulje, status models.PuljeStatus) {
	t.Helper()
	body, _ := json.Marshal(map[string]string{"puljeStatus": string(status)})
	req := httptest.NewRequest(http.MethodPut, "/api/puljer/"+string(pulje)+"/status", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
}

func TestPuljeStatusHandler_LockCommitsUnlockReverts(t *testing.T) {
	db, logger := testutil.CreateTestDBAndLogger(t, "test_admin_lock")
	router := chi.NewRouter()
	puljefordelingStatusRoute(router, db, &live.Manager{}, logger)

	const fredag = models.PuljeFredagKveld
	mustExec(t, db, `INSERT INTO puljer (id, name, status, start_at, end_at) VALUES (?, 'F', 'Open', ?, ?)`,
		string(fredag), "2026-09-04T18:00:00Z", "2026-09-04T22:00:00Z")
	mustExec(t, db, `INSERT INTO events (id, title, intro, description, host_name, email, phone_number, max_players)
		VALUES ('evA','Alpha','','','','','',4)`)
	mustExec(t, db, `INSERT INTO relation_event_puljer (event_id, pulje_id, is_in_pulje) VALUES ('evA', ?, 1)`, string(fredag))
	mustExec(t, db, `INSERT INTO billettholdere (id, first_name, last_name, ticket_type_id, ticket_type, order_id, ticket_id) VALUES (1,'Anna','A',0,'',0,1)`)
	mustExec(t, db, `INSERT INTO interests (billettholder_id, event_id, pulje_id, interest_level) VALUES (1,'evA',?, ?)`,
		string(fredag), string(models.InterestLevelHigh))

	// Lock → commit writes a solver seat and status becomes Locked.
	putPuljeStatus(t, router, fredag, models.PuljeStatusLocked)

	var solverCount int
	mustQueryRow(t, db, `SELECT COUNT(*) FROM relation_events_players WHERE pulje_id=? AND source='solver'`, string(fredag)).Scan(&solverCount)
	if solverCount != 1 {
		t.Errorf("after lock: want 1 solver seat, got %d", solverCount)
	}
	var status string
	mustQueryRow(t, db, `SELECT status FROM puljer WHERE id=?`, string(fredag)).Scan(&status)
	if status != string(models.PuljeStatusLocked) {
		t.Errorf("after lock: want status Locked, got %s", status)
	}

	// Unlock → revert removes solver seats and status becomes Open.
	putPuljeStatus(t, router, fredag, models.PuljeStatusOpen)

	mustQueryRow(t, db, `SELECT COUNT(*) FROM relation_events_players WHERE pulje_id=? AND source='solver'`, string(fredag)).Scan(&solverCount)
	if solverCount != 0 {
		t.Errorf("after unlock: want 0 solver seats, got %d", solverCount)
	}
	mustQueryRow(t, db, `SELECT status FROM puljer WHERE id=?`, string(fredag)).Scan(&status)
	if status != string(models.PuljeStatusOpen) {
		t.Errorf("after unlock: want status Open, got %s", status)
	}
}
```

Add these small DB helpers in the same file (only if the `admin` package test files do not already define them — search first with `grep -rn "func mustExec\|func mustQueryRow" pages/admin/`):

```go
func mustExec(t *testing.T, db *sql.DB, query string, args ...any) {
	t.Helper()
	if _, err := db.Exec(query, args...); err != nil {
		t.Fatalf("exec %q: %v", query, err)
	}
}

func mustQueryRow(t *testing.T, db *sql.DB, query string, args ...any) *sql.Row {
	t.Helper()
	return db.QueryRow(query, args...)
}
```

(If you add the helpers, also add `"database/sql"` to the test file's imports. If the
package already provides `testutil.MustExec`, prefer that and drop the local `mustExec`.)

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./pages/admin/ -run TestPuljeStatusHandler_LockCommitsUnlockReverts -v`
Expected: FAIL — before the wiring (Step 3) the handler only calls `updatePuljeStatus`, so no solver seat is written: "after lock: want 1 solver seat, got 0".

- [ ] **Step 3: Wire the handler**

In `pages/admin/puljefordeling.templ`, add the import:

```go
	"github.com/Regncon/conorganizer/service/puljefordeling"
```

Then replace the body of the `router.Put("/api/puljer/{puljeId}/status", ...)` handler (currently `puljefordeling.templ:108-146`, after the `isValidPuljeStatus` check) so the transition drives commit/revert. Replace the block from `if err := updatePuljeStatus(...)` through the broadcast with:

```go
		// Read the current status to decide the transition side-effects.
		var currentRaw string
		if err := db.QueryRow(`SELECT status FROM puljer WHERE id = ?`, string(puljeID)).Scan(&currentRaw); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Error(w, "Pulje not found", http.StatusNotFound)
				return
			}
			logger.Error(fmt.Errorf("read current pulje status: %w", err).Error(), "pulje_id", puljeID)
			http.Error(w, "Failed to read pulje status", http.StatusInternalServerError)
			return
		}
		current := models.PuljeStatus(currentRaw)
		next := store.PuljeStatus

		var transitionErr error
		switch {
		case current == models.PuljeStatusOpen && next == models.PuljeStatusLocked:
			// Lock: commit the distribution and set status (atomic, inside the service).
			transitionErr = puljefordeling.CommitPuljeAssignments(db, puljeID, logger)
		case current == models.PuljeStatusLocked && next == models.PuljeStatusOpen:
			// Unlock: drop solver seats and reopen (atomic, inside the service).
			transitionErr = puljefordeling.RevertPuljeAssignments(db, puljeID)
		default:
			// Locked↔Completed and any no-op: status only.
			transitionErr = updatePuljeStatus(db, puljeID, next)
		}
		if transitionErr != nil {
			if errors.Is(transitionErr, errPuljeNotFound) {
				http.Error(w, "Pulje not found", http.StatusNotFound)
				return
			}
			logger.Error(transitionErr.Error(), "pulje_id", puljeID, "pulje_status", next)
			http.Error(w, "Failed to update pulje status", http.StatusInternalServerError)
			return
		}

		if err := liveManager.Broadcast(r.Context(), live.BucketEvents); err != nil {
			logger.Error(fmt.Errorf("failed to broadcast pulje status update: %w", err).Error(), "pulje_id", puljeID, "pulje_status", next)
			http.Error(w, "Failed to broadcast update", http.StatusInternalServerError)
			return
		}
		if err := liveManager.Broadcast(r.Context(), live.BucketInterests); err != nil {
			logger.Error(fmt.Errorf("failed to broadcast interests update: %w", err).Error(), "pulje_id", puljeID, "pulje_status", next)
			http.Error(w, "Failed to broadcast update", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
```

(Confirm `database/sql` is imported in `puljefordeling.templ` — it already is. `errors`, `fmt` are already imported.)

- [ ] **Step 4: Regenerate templ and build**

Run: `go tool templ generate`
Run: `go build ./...`
Expected: builds cleanly.

- [ ] **Step 5: Run tests to verify they pass**

Run: `go test ./pages/admin/ -run TestPuljeStatusHandler_LockCommitsUnlockReverts -v`
Expected: PASS.

- [ ] **Step 6: Commit**

```bash
git add pages/admin/puljefordeling.templ pages/admin/puljefordeling_templ.go pages/admin/puljefordeling_status_test.go
git commit -m "feat(puljefordeling): lock commits assignment, unlock reverts it"
```

---

### Task 7: Emulate page — re-run button, manual marker, copy

**Files:**
- Modify: `pages/admin/puljefordeling_emulate/puljefordeling_emulate.templ`
- Regenerate: `pages/admin/puljefordeling_emulate/puljefordeling_emulate_templ.go` (via `go tool templ generate`)

**Interfaces:**
- Consumes: `puljefordeling.EmulatedEvent` / `AssignedPlayer.Pinned` (from Task 4).

- [ ] **Step 1: Render the pinned marker**

In `puljefordeling_emulate.templ`, in `templ eventCard`, replace the player `<li>` (lines 85-90) with one that shows a 📌 for pinned seats:

```templ
					for _, pl := range ev.AssignedPlayers {
						<li class={ templ.KV("emulate-moved", pl.Moved) }>
							<span class="emulate-emoji">{ pl.Level.Emoji() }</span>
							<span class={ templ.KV("emulate-dm", pl.IsDM) }>{ pl.Name }</span>
							if pl.Pinned {
								<span class="emulate-pinned" title="Manuelt plassert">📌</span>
							}
						</li>
					}
```

- [ ] **Step 2: Add the re-run button and fix the intro copy**

In `templ emulatePage`, replace the intro/legend block (lines 21-27) with copy that reflects persistence, and add a re-run button. The button reloads the page, which recomputes the distribution (`EmulateSeatings` runs on every GET):

```templ
			<h1 class="page-heading">Emulér puljefordeling</h1>
			<p class="emulate-intro">
				Forhåndsvis fordelingen for hver pulje. Åpne puljer er en live simulering du kan
				kjøre på nytt så ofte du vil; når en pulje låses, lagres fordelingen og videre
				endringer gjøres manuelt. Låste og fullførte puljer viser den faktiske fordelingen.
			</p>
			<p class="emulate-legend">
				🔥 Veldig interessert · 👍 Middels · 🤷 Litt · 📌 Manuelt plassert ·
				<span class="emulate-dm">Turkis</span> = spilleder et annet sted ·
				<span class="emulate-moved-legend">Rød strek</span> = flyttet ned for å gi plass til andre
			</p>
			<a role="button" href="/admin/puljefordeling-emulate/" class="btn btn--outline emulate-rerun">
				Kjør fordeling på nytt
			</a>
```

- [ ] **Step 3: Add a style for the marker**

In `templ emulateStyles`, add inside the `<style>` block (e.g. after the `.emulate-moved` rule):

```css
			.emulate-pinned {
				flex: 0 0 auto;
			}
			.emulate-rerun {
				align-self: flex-start;
				margin-bottom: var(--spacing-4x);
			}
```

- [ ] **Step 4: Regenerate templ and build**

Run: `go tool templ generate`
Run: `go build ./...`
Expected: builds cleanly.

- [ ] **Step 5: Run the package tests**

Run: `go test ./... 2>&1 | tail -20`
Expected: PASS across the repo.

- [ ] **Step 6: Commit**

```bash
git add pages/admin/puljefordeling_emulate/puljefordeling_emulate.templ pages/admin/puljefordeling_emulate/puljefordeling_emulate_templ.go
git commit -m "feat(puljefordeling): emulate page re-run button, manual marker, persistence copy"
```

---

## Final verification

- [ ] `go tool templ generate && go build ./...`
- [ ] `go test ./service/puljefordeling/... ./pages/admin/...`
- [ ] Manual smoke test (via `/run` or `go run .`): on `/admin/puljefordeling`, pin a participant to an event (existing admin flow); open the emulate page and confirm the seat shows 📌 and the solver fills around it; click "Kjør fordeling på nytt" and watch it recompute; lock the pulje → confirm `relation_events_players` gains `role=Player, source=solver` rows and the manual row persists; move a player manually; unlock → solver rows gone, manual kept, status Open; lock again, then publish (Completed) → confirm the next pulje's preview reflects the prior pulje's actual seats.

## Notes / deferred

- The "re-run" button reloads the page (the solve is recomputed on every GET). A future enhancement could stream live updates via `live.Manager` + `live.BucketInterests` instead of a manual reload.
- Manual placements are **not** counted as "misses" in fairness (a deliberate placement is not a missed top choice). If organisers want pinned players to still earn the scarcity bonus, revisit `applyResult`'s miss loop.
- If too many manual placements exceed an event's capacity, all pins are still honored (a hand-placement is never dropped); the event renders over capacity. Surfacing that as a warning is out of scope here.
