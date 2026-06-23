# Puljefordeling: real grouping assignment integrated into the pulje lifecycle

**Date:** 2026-06-23
**Status:** Design — pending review

## Problem

`service/puljefordeling/` already contains a correct min-cost max-flow solver
(`solver/flow.go` — SPFA augmentation; `solver/solver.go` — weight bands and
cross-slot fairness). Today it is **read-only**: `EmulateSeatings` (`emulate.go`)
loads data, runs the solver across every pulje in one in-memory pass, and renders a
preview at `/admin/puljefordeling-emulate/`. Its package doc states it "never writes to
the database." Nothing turns the proposal into real groups.

We want **proper integration** that matches how organisers actually run the weekend —
one pulje at a time, in chronological order:

1. **While signup is Open** — re-run the distribution often as a live preview, seeded by
   what already happened in earlier (frozen) puljer.
2. **On Lock** — closing signup also commits that pulje's solver result into the
   database as real groups, preserving any manual placements. The solver then stops for
   that pulje.
3. **After Lock** — admins move/shuffle people manually; no re-running.
4. **On Publish (Completed)** — the pulje is final; later puljer are seeded from these
   actual seats.

## Key facts established during exploration

- `Open → Locked → Completed` statuses already exist (`models.PuljeStatus`;
  `pages/admin/puljefordeling_templ.go` has the status card UI and the
  `PUT /admin/api/puljer/{puljeId}/status` handler).
- `relation_events_players` rows with `role='Player'` are **admin-created only** — public
  signup writes only `interests`. So any existing `Player` row is a deliberate manual
  placement.
- **Manual placements happen throughout Open**, often long before lock (e.g. an admin
  seats a DM's child at a specific table). They are authoritative at all times; the
  solver always works *around* them.
- A manually placed player **may have no `interests` row** for that event. The pin must be
  honored regardless of whether a preference exists.
- Tests bootstrap the schema from `schema.sql` (`service/testdb.go`); production existing
  DBs use goose migrations (`migrations/`); `initialize.sql` is the fresh-bootstrap copy.
- The solver itself is correct; the work is integration, pinned-seat support, and
  fairness-seeding from persisted reality — not an algorithm change. `flow.go` is
  untouched.

## Lifecycle (state machine, per pulje)

| Transition | What happens |
|---|---|
| **Open** (default) | Signup active. Live preview; "re-run" button. Nothing persisted by the solver. Preview seeded from earlier frozen puljer; existing manual placements pinned. |
| **Open → Locked** | Signup closes **and** commit: solver-chosen seats written to `relation_events_players` as `role='Player', source='solver'`; existing manual rows preserved. Solver stops for this pulje. Transaction also flips status to `Locked`. |
| **Locked** | Admins move/shuffle via existing manual endpoints (those write `source='manual'`). No solver re-runs. |
| **Locked → Completed** | Publish. Status only — seats already persisted. Downstream puljer now seed from these seats. |
| **Locked → Open** (unlock) | Remove only `role='Player' AND source='solver'` rows for this pulje (manual rows survive); flip status to `Open`; signup reopens. |
| **Completed → Locked** | Status only (un-publish). |

**Confirmed decisions:**
- **Lock = close signup + commit in one action** (no separate commit step).
- **Fairness seeds from puljer that are Locked OR Completed** (both have real persisted
  seats and reflect "what actually happened so far").

## Architecture

Three layers, each independently testable.

### 1. Schema — a `source` discriminator

Add to `relation_events_players`:

```sql
source TEXT NOT NULL DEFAULT 'manual' CHECK (source IN ('manual','solver'))
```

Why: lets the solver (a) pin manual placements during a solve, and (b) on unlock remove
*only* the seats it created, leaving hand-placements intact. Existing rows default to
`'manual'`, which is correct — they are all admin placements/GMs.

Touch all three schema sources (matching the precedent set by the pulje-status changes):
- `schema.sql` (canonical; loaded by tests)
- `initialize.sql` (fresh bootstrap)
- `migrations/<timestamp>_add_source_to_relation_events_players.sql` — goose
  (`-- +goose Up` / `-- +goose Down`). Prefer `ALTER TABLE ... ADD COLUMN`; fall back to
  the create-new/copy/drop/rename pattern (see
  `migrations/20260522120000_pulje_status_open_locked_completed.sql`) if the CHECK-on-add
  is rejected on the STRICT table.

### 2. Solver — pinned seats + fairness replay (`solver/solver.go`)

`flow.go` is unchanged. Three changes in `solver.go`:

1. **Extract `applyResult(slot, players, assignments)`** — pull the
   seated/satisfied/misses + `NewlySatisfied`/`TotalScore` update block out of
   `SolveSlot` (currently `solver.go:204-247`) into one private method, so the fairness
   rules live in exactly one place.

2. **Pinned seats** — change the signature to
   `SolveSlot(slot model.Slot, players []model.Player, fixed map[string]string)` where
   `fixed` maps `playerID → eventID` (the manual placements to honor):
   - Operate on a copy of `slot.Events` whose `Capacity` is reduced by the number of
     pinned players in each event (clamp ≥ 0).
   - Remove pinned players (and DMs-here) from the pool passed to `runMCMF` — a pinned
     player is reserved, not re-contested, and needs **no preference edge**.
   - Merge pinned players back into `result.Assignments` before `applyResult`, so they
     count as seated (and as satisfied iff their pinned event is their top choice).
   - Ignore a `fixed` entry whose event is not in this slot, or whose player DMs here.
   - Existing callers pass `nil`.

3. **`ApplyActual(slot, players, assignments map[string][]string) model.SlotResult`** —
   seed `State` from a known persisted assignment **without solving** (reuses
   `applyResult`). Used to replay each frozen prior pulje so satisfied/seated/misses
   reflect reality, with the exact same rules as a solved slot.

### 3. Orchestration + wiring (`emulate.go`, `puljefordeling_templ.go`)

Refactor `emulate.go` from "one in-memory pass" to "chronological, seeded by persisted
reality":

- **New loaders:** `loadPersistedAssignments(db)` → `map[puljeID]map[eventID][]playerID`
  for `role='Player'`, tagging each row manual/solver via `source`; pulje statuses come
  from the existing `loadPuljer`.
- **Shared chronological core** threading one `State`:
  - frozen prior (Locked/Completed) → `state.ApplyActual(slot, players, persisted)`;
  - open prior → `state.SolveSlot(slot, players, manualFixed)` (preview chain);
  - target/open pulje → `state.SolveSlot(slot, players, manualFixed)`.
  `manualFixed` for a pulje = its existing `role='Player'` rows.
- **`EmulateSeatings(db)`** reuses the core: frozen puljer render their **actuals**, open
  puljer render the **live preview**. Display shaping (`shapePulje`/`assignedPlayers`) is
  unchanged — it consumes the resulting `SlotResult`.
- **`CommitPuljeAssignments(db, puljeID, logger) error`** (Lock path), in one tx:
  compute the target solve (seeded from frozen priors, manual pinned) → delete existing
  `role='Player' AND source='solver'` rows for the pulje → insert solver-chosen seats as
  `role='Player', source='solver'`. Manual + GM rows untouched. Idempotent.
- **`RevertPuljeAssignments(db, puljeID) error`** (Unlock path): delete only
  `role='Player' AND source='solver'` rows for the pulje.
- Set `source='manual'` explicitly in the admin insert paths in
  `components/formsubmission/who_is_interested_templ.go` (`AddPlayersFirstChoice`,
  `UpdatePlayerStatus`) — the default covers it, but be explicit.

**Lifecycle wiring** — in the existing `PUT /admin/api/puljer/{puljeId}/status` handler
(`puljefordeling_templ.go:108`), branch on the requested status within a transaction:
- Open → Locked: `CommitPuljeAssignments` + `updatePuljeStatus(Locked)`; broadcast
  `live.BucketInterests` + `live.BucketEvents`.
- Locked → Open: `RevertPuljeAssignments` + `updatePuljeStatus(Open)`.
- Locked → Completed / Completed → Locked: status only.

### UI — re-run button, manual-seat marker, copy (`pages/admin/puljefordeling_emulate/`)

- Add a **"Kjør fordeling på nytt"** button that recomputes and swaps the
  `#puljefordeling-emulate` fragment via a sibling route (e.g.
  `GET /admin/puljefordeling-emulate/api/`) wired with datastar `@get`, consistent with
  other admin live fragments (`admin.go:198-211`).
- **Visually distinguish a pinned/manual seat** from a solver-proposed one in the preview
  (e.g. a small "📌 manuelt plassert" marker), so organisers see what is hand-placed vs
  computed. This needs `AssignedPlayer` to carry a `Pinned`/`Manual` flag, populated from
  the `source` of the underlying row.
- Update the intro copy (`puljefordeling_emulate.templ:21-24`): currently "ingenting
  lagres" — clarify that Open puljer are a live preview while Locked/Completed puljer
  show the committed seating, and that locking persists.

## Data flow (Open → Lock example)

```
Admin pins a DM's child  ──► relation_events_players (role=Player, source=manual)
                              (capacity for that event effectively −1)

Admin opens emulate page / clicks re-run
   EmulateSeatings walks puljer chronologically:
     frozen priors  → ApplyActual(persisted seats)  → seeds satisfied/seated/misses
     this pulje      → SolveSlot(players, fixed={child→event})
                        ├─ child reserved, capacity reduced, no preference needed
                        └─ solver fills remaining seats by preference + fairness
   → preview rendered (manual seats marked 📌)

Admin flips "Locked"
   PUT /api/puljer/{id}/status  (tx):
     CommitPuljeAssignments → delete old source=solver rows
                            → insert solver seats (source=solver)
     updatePuljeStatus(Locked)
   manual rows preserved; solver stops for this pulje
```

## Edge cases & decisions

- **Pinned player with no interest** — honored; seated, not satisfied (no top-choice
  match). Counts as seated for future-pulje fairness.
- **Pinned seats exceed capacity** (too many manual placements) — all pins honored
  (we never drop a hand-placement); event renders over capacity and should be surfaced as
  a warning rather than silently truncated.
- **Locking a pulje while an earlier pulje is still Open** — allowed; it seeds from
  whatever is currently frozen. (Optional: warn in the UI. Not blocking.)
- **Re-lock after unlock** — `CommitPuljeAssignments` is idempotent: it clears prior
  `source=solver` rows first, so a fresh solve replaces them while manual rows persist.

## Testing

Reuse existing helpers (`seedPulje/seedEvent/seedParticipant/seedInterest/seedGM`);
add `seedAssignment(..., source)` and a pulje-status seeder.

- **`solver_test.go`:** pinned seat honored without a preference edge; capacity reduced by
  pins; pinned player counted in fairness; `SolveSlot(..., nil)` regression;
  `ApplyActual` reconstructs satisfied/seated/misses correctly.
- **`emulate_test.go`:** target pulje seeded from a frozen prior's actuals; manual rows
  pinned in preview and at lock; `CommitPuljeAssignments` writes `source='solver'` rows
  and preserves manual/GM; `RevertPuljeAssignments` removes only solver rows; whole-
  weekend `EmulateSeatings` shows actuals for frozen + preview for open.
- **schema check:** assert the `source` column exists (tests load `schema.sql`).

## Critical files

- `service/puljefordeling/solver/solver.go` — pinned seats, `applyResult`, `ApplyActual`
- `service/puljefordeling/emulate.go` — seeded orchestration, commit/revert, loaders, pinned-seat marker data
- `pages/admin/puljefordeling_templ.go` — lock/unlock → commit/revert wiring
- `pages/admin/puljefordeling_emulate/puljefordeling_emulate.{templ,go}` — re-run button, manual marker, copy
- `components/formsubmission/who_is_interested_templ.go` — explicit `source='manual'`
- `schema.sql`, `initialize.sql`, `migrations/<ts>_add_source_to_relation_events_players.sql`

## Verification

1. `go tool templ generate` (regenerate `_templ.go` after `.templ` edits).
2. `go build ./...` and `go test ./service/puljefordeling/... ./pages/admin/...`.
3. Manual (via `/run` or `go run .`): pin a player on `/admin/...`; open the emulate page
   and confirm the seat shows as manual and the solver fills around it; click re-run and
   watch the preview update; lock the pulje → confirm `relation_events_players` gains
   `role=Player, source=solver` rows and the manual row persists; move a player manually;
   unlock → solver rows gone, manual kept; publish → confirm the next pulje's preview
   reflects the prior pulje's actual seats.
