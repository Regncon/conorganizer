# Puljefordeling: per-pulje tabbed view

**Date:** 2026-06-23
**Status:** Approved design, ready for implementation plan

## Summary

Replace the current two-part puljefordeling UI (an inline status-card grid on
the `/admin/` dashboard plus a separate all-puljer emulation page) with a single
tabbed page that mirrors the rooms-assignment page (`/admin/rooms/assignment/{pulje}`).
Each tab shows one pulje's live distribution together with its lock/publish
control. The solver only computes up to the pulje being viewed.

This is a front-end / routing redesign. The backend from the closed PR #451
(solver, emulate/commit split, `source` column migration) is kept as-is; only a
small, focused service helper is added.

## Motivation

- The current emulation page stacks all puljer vertically and exposes a single
  global rerun, which does not scale visually and recomputes everything.
- Lock/publish controls live on the dashboard, separate from the distribution
  they affect, so the admin has to switch context.
- The rooms-assignment page already establishes a tabbed, per-pulje, live
  pattern the team likes. Reusing it keeps the admin UI consistent.

## Routing & entry point

| Route | Type | Purpose |
|---|---|---|
| `/admin/puljefordeling/` | GET | Redirect → `/admin/puljefordeling/{first pulje}` (`models.PuljeFredagKveld`) |
| `/admin/puljefordeling/{pulje}` | GET | HTML shell: tabs bar + `data-init` subscribing to the SSE API below |
| `/admin/puljefordeling/api/{pulje}` | GET (SSE) | `liveManager.Stream(...)`; renders the selected pulje's distribution + lock/publish control; re-renders on broadcasts |
| `/admin/puljefordeling/api/puljer/{puljeId}/status` | PUT | The existing lock/publish transition handler, relocated under this route group |

- The tabs bar is one `<a href="/admin/puljefordeling/{pulje}">` per
  `models.AllPuljer()`, with the `active` class on the current pulje — identical
  markup pattern to `pages/admin/rooms/rooms_index.templ`.
- Invalid `{pulje}` values return `400` (same as the rooms-assignment handler)
  or redirect to the first pulje; follow whichever the rooms handler does for
  consistency.

### Dashboard changes

- The inline status-card grid (`@puljefordeling(db)` in
  `pages/admin/admin_page.templ`) and the separate "Emulér puljefordeling" card
  collapse into **one** admin card: "Puljefordeling" with a
  "Gå til puljefordeling" button linking to `/admin/puljefordeling/`.
- The old `/admin/puljefordeling-emulate/` route and the
  `pages/admin/puljefordeling_emulate/` page are removed. Their rendering logic
  (legend, `puljeResult`, `eventCard`, styles) moves into the new per-pulje tab
  content.

## Per-pulje tab content

The SSE `Render` callback renders, for the selected pulje:

1. **Header row:** pulje name + summary stats (deltakere med ønsker, fikk
   førstevalg, uten plass) on the left; the **status control**
   (Åpen / Låst / Fullført) on the right.
2. **Legend:** the existing emoji/colour legend.
3. **Event grid:** the existing `eventCard` markup, one card per event in the
   pulje, showing assigned players with level emoji, DM colour, pinned marker,
   and the moved-down indicator.
4. **Uten plass:** list of unassigned participants, if any.

Reuse the existing `puljeResult` / `eventCard` / legend / styles from the
emulate page, scoped to a single pulje instead of looped over all.

### Live behaviour

- The shell's `data-init` is `live.DatastarInit("/admin/puljefordeling/api/{pulje}")`.
- The SSE stream subscribes to `live.BucketInterests` and `live.BucketEvents`
  (the buckets the status PUT already broadcasts, and the ones that change a
  distribution).
- **Open** pulje → live simulation: the solver runs on each render.
- **Locked / Completed** pulje → renders the persisted seats
  (`EmulateSeatings`/`EmulatePulje` already return persisted data for frozen
  puljes via `puljeFrozen`).
- Locking a pulje re-renders the tab to show the now-saved distribution;
  unlocking reverts and resumes live simulation. Backend transition behaviour
  (`CommitPuljeAssignments` / `RevertPuljeAssignments`) is unchanged.
- No explicit "Kjør på nytt" button — live updates make it redundant.

### Status control

The Åpen / Låst / Fullført control is the existing lock/publish control from the
status card, moved into the tab header. It PUTs to
`/admin/puljefordeling/api/puljer/{puljeId}/status`. Locking commits the
distribution; unlocking reverts; locked↔completed is status-only. No backend
change to the transition handler beyond its route path.

## Backend addition: scoped emulation

Add a focused helper next to `EmulateSeatings` in
`service/puljefordeling/emulate.go`:

```go
// EmulatePulje solves chronologically only up to (and including) puljeID and
// returns just that pulje's distribution. Earlier puljer are solved as needed
// (locked ones contribute their persisted seats); later puljer are never
// computed.
func EmulatePulje(db *sql.DB, puljeID models.Pulje) (EmulatedPulje, EmulationMeta, error)
```

- Loads seating data, finds the chronological index of `puljeID`, calls
  `d.solveChronological(idx)`, and shapes only `puljer[idx]`.
- `EmulationMeta` carries the per-pulje summary numbers the header needs
  (player count with wishes, satisfied count for this pulje). Exact shape to be
  finalised in the plan; may reuse fields already on `EmulatedPulje`.
- The existing `EmulateSeatings` (all-puljer) stays for any other caller; the
  tab page does not use it.

### Why "don't emulate the rest until locked"

`solveChronological(upTo)` already accepts a stop index; `EmulateSeatings` just
always passes the last pulje. Solving only up to the selected pulje means
opening e.g. the Fredag kveld tab never runs the solver for the three later
puljer. A later pulje only becomes relevant once an earlier one is locked, at
which point its persisted seats feed forward as fixed input to later solves.
This also bounds the per-render cost of the live view.

## Testing

- **Service:** unit test `EmulatePulje` — (a) returns only the requested pulje;
  (b) for a locked pulje returns persisted seats unchanged; (c) solving up to an
  earlier pulje does not depend on or compute later puljer; (d) equivalence:
  `EmulatePulje(db, lastPulje)` matches the corresponding entry from
  `EmulateSeatings(db)`.
- **Route:** the SSE/page handler returns `400` for an invalid pulje and renders
  the expected tab for a valid one; the status PUT still performs lock/commit and
  unlock/revert under the new path.
- Keep the existing puljefordeling solver/commit/emulate tests passing; update
  any test that referenced the removed emulate route or the moved status route.

## Out of scope

- Manual seat editing within the tab (handled by the rooms-assignment page).
- Changes to the solver algorithm or the commit/revert semantics.
- Any change to the `source` column migration (already applied).

## Risks / watch-items

- **Live solver cost:** even scoped to one pulje, an open pulje re-runs the
  solver on every `BucketInterests` / `BucketEvents` broadcast for every admin
  with the tab open. Bounded by the per-pulje scope above; revisit with caching
  if it shows up in practice.
- **Route relocation:** moving the status PUT path means updating its callers
  (the status control markup) and any tests that hit the old path.
