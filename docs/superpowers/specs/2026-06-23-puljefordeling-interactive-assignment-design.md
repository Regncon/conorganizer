# Puljefordeling: interactive per-game player assignment

**Date:** 2026-06-23
**Status:** Approved design, ready for implementation plan
**Builds on:** `2026-06-23-puljefordeling-tabs-design.md` (the per-pulje tabbed view)

## Summary

Make each game (event) card in the puljefordeling tab interactive, mirroring the
romfordeling (rooms-assignment) page: every game lists its assigned players, has a
**"+"** that opens a shared searchable picker modal to add a player, and a **"×"**
to remove a seat. Manual additions are pins (`source='manual'`) that the solver
already honors by reducing the game's effective capacity. A locked pulje gains a
**"Rerun fordeling"** button that re-solves around the remaining pins.

Lock/publish controls stay where the tabbed-view work put them — in the tab
header only; the dashboard remains a single "Gå til puljefordeling" link.

Almost all backend already exists; the only new endpoint is the rerun action.

## Motivation

The tabbed view (previous spec) is read-only. The team wants the same hands-on
assignment workflow romfordeling already provides — assign participants to games
directly — without destabilising the solver. The existing `event-players`
endpoints and the `<admin-billettholder-search>` picker already implement manual
assignment for the event-edit page; this spec brings that interaction into the
puljefordeling tab and adds the small pieces that are missing (a locked-pulje
re-solve, and the dashboard control restoration).

## Behavior by pulje state

Each game card is interactive, but what is editable depends on the pulje's status.

| | **Open** (live solver) | **Locked / Completed** (frozen, persisted rows) |
|---|---|---|
| Player list shown | Live solver emulation + manual pins (📌) | Committed seats (solver + manual rows), via `ApplyActual` |
| **"+"** (add) | Adds a manual pin → solver reduces that game's capacity → live re-render | Adds a manual pin row |
| **"×"** (remove) | **Only on 📌 manual pins** → removes the pin; solver reclaims the seat next render | On **any** seated player → unseats (deletes the row) |
| **Rerun** | Not shown (open puljer auto re-render live) | **"Rerun fordeling"** button → re-solve, reseating freed players around remaining pins |

Rationale: open puljer are a live simulation, so the only persisted thing an admin
can remove is a manual pin they themselves created; the solver's own placements are
not rows yet. Once locked, every seat is a real row, so removal and an explicit
re-solve become meaningful.

## Card and modal UI

Mirror `roomPuljeContainer` + `assignEventModal` from
`pages/admin/rooms/rooms_assignment_page.templ`, swapping room→game and
game→player. Keep the existing **compact emoji player rows** (a player is not an
event with a banner image), not romfordeling's image-banner cards.

Per game card (extends the existing `puljefordelingTabEventCard`):
- Existing header: title, GM name, `assigned / capacity`, undersubscribed badge.
- Existing player rows: level emoji, name (turquoise if SL elsewhere), 📌 if pinned,
  red left-border if moved.
- **New:** a conditional **"×"** on each player row — shown when
  `pulje is open && player.Pinned`, or when `pulje is locked/completed` (any player).
- **New:** a **"+"** button at the bottom of the card (same placement as the room
  card's "+").

Per tab:
- **New:** a **"Rerun fordeling"** button in the tab header, shown only when the
  pulje is locked or completed.
- **New:** **one shared `<dialog>`** (romfordeling's single-dialog pattern). The
  card "+" sets `$assignmentEventId = '<eventID>'` and `$assignmentPuljeId =
  '<pulje>'`, then opens the dialog.

Dialog body: reuse the existing
`components/formsubmission.billettholderAssignmentActions(eventId, puljeId,
billettholdere)` — the `<admin-billettholder-search>` searchable picker over all
billettholdere plus "Legg til som spiller" / "Legg til som SL" buttons. Because
the card sets `$assignmentEventId` before opening, the dialog targets the right
game. (If `billettholderAssignmentActions` hard-codes its own `eventId`, the plan
will either parameterise the dialog per render or bind it to the
`$assignmentEventId` signal — to be resolved in the plan against the component's
actual markup.)

## Backend

### Reused unchanged

- `POST /admin/approval/api/event-players/post/add_first_choice` →
  `formsubmission.AddPlayersFirstChoice`: upserts a high interest **and** a
  `relation_events_players` row with `role='Player', source='manual'`. The solver
  already reduces each event's effective capacity by the number of pinned players
  (`solver.SolveSlotFixed`, "Reduce each event's capacity by the number of players
  pinned into it"), so "manual seat reduces capacity" needs **no** solver change.
- `POST /admin/approval/api/event-players/post/add_gm` → manual GM pin.
- `PUT /admin/approval/api/event-players/update_status` with
  `isPlayer=false, isGm=false` → DELETEs the `relation_events_players` row (the "×"
  remove).
- `getBillettholdere`, `billettholderAssignmentActions`,
  `/static/web_components/admin_billettholder_search.js` — the picker.
- `service/puljefordeling.CommitPuljeAssignments` — deletes `source='solver'` rows,
  re-solves chronologically around manual pins, re-inserts solver rows.

All three event-players POST/PUT handlers already broadcast `live.BucketInterests`,
which the puljefordeling tab subscribes to — so adds/removes re-render the tab with
no extra wiring.

### New: rerun endpoint

```
PUT /admin/puljefordeling/api/{pulje}/rerun
```

- Parse and validate `{pulje}` (`models.ParsePulje`); 400 on invalid.
- Read the pulje's current status. Only act when **locked or completed**; if open,
  no-op with HTTP 409 (open puljer are already live — there is nothing to re-solve
  on demand).
- Call `CommitPuljeAssignments(db, pulje, logger)`.
- On success, broadcast `live.BucketInterests` and `live.BucketEvents` (best-effort,
  matching the status-PUT handler) and return 204.
- Register it in the puljefordeling route group alongside the page/SSE routes.

### Render data

The interactive card needs, per render:
- Each player's `Pinned` flag — already on `puljefordeling.AssignedPlayer`.
- The pulje status — already read in the tab content (`puljeStatusFor`) to drive
  state-dependent "×" visibility and the rerun button.
- The billettholder list for the modal — `getBillettholdere(db, logger)`.
- `event_id`, `pulje_id`, `billettholder_id` for each "+"/"×" action — all present
  in the render.

## Lock/publish controls (unchanged)

Lock/publish stays in the tab header only, exactly as the tabbed-view work left it.
The `/admin/` dashboard remains a single "Gå til puljefordeling" link card — the
inline status grid is **not** restored. No change to the status endpoint or to
`pages/admin/admin_page.templ` for this work.

## Testing

- **Rerun endpoint:** locked pulje → re-solves and reseats freed players around
  remaining pins (assert solver rows regenerated, manual pins preserved); open pulje
  → 409 no-op; invalid pulje → 400.
- **Card render:** "×" appears only on pinned players when the pulje is open, and on
  all seated players when locked; the "+" wires the correct `assignmentEventId` /
  `assignmentPuljeId`; the "Rerun fordeling" button shows only when locked/completed.
- Reuse the existing `event-players` and `CommitPuljeAssignments` test coverage; do
  not duplicate it.
- **Verification gate:** `go test ./...` and `golangci-lint run` both clean.

## Out of scope

- **Solver exclusion** ("×" on a *solver-placed* player in an open pulje to forbid
  the solver from seating them there). Explicitly dropped; instead, open-pulje
  removal is manual-pin-only, and you influence open distributions by pinning.
- Changes to the solver algorithm, the commit/revert semantics, or the
  `source` column / migration.
- Manual seat editing for puljer that have no events.

## Risks / watch-items

- **Dialog targeting:** the reused `billettholderAssignmentActions` was written for
  a single event-edit page where `eventId` is fixed. Driving one shared dialog from
  many game cards relies on the `$assignmentEventId` signal being set before the
  dialog opens; the plan must confirm the component reads the signal (not a
  hard-coded id) or adapt it. This is the main integration risk.
- **Live solver cost** (inherited from the tabbed-view spec): an open pulje re-runs
  the scoped solver on every `BucketInterests` broadcast; each "+"/"×" triggers one.
  Bounded by the per-pulje scope; revisit with caching only if it shows up.
- **Rerun on a busy locked pulje:** `CommitPuljeAssignments` re-solves the whole
  pulje; fine for an explicit button press, but it should not be wired to fire
  automatically on every broadcast.
