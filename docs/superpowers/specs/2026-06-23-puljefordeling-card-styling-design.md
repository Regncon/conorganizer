# Puljefordeling: romfordeling-style box styling for game cards

**Date:** 2026-06-23
**Status:** Approved design, ready for implementation plan
**Builds on:** `2026-06-23-puljefordeling-interactive-assignment-design.md`

## Summary

Restyle the puljefordeling tab's game cards to match the romfordeling
(rooms-assignment) "table selection" boxes: each game card becomes a filled,
rounded box like `.room`, and each assigned player becomes a bordered tile like
`.room-event` (without the image banner, since players have no image). Pure CSS
plus a minor markup tidy in `pages/admin/puljefordeling_tab.templ`
(`puljefordelingTabStyles`); no behavior, copy, routing, or data changes.

## Motivation

The interactive game cards work, but visually they are flat — a thin border, no
background — and the player rows are plain text rows. The team wants them to read
as the same kind of UI as the romfordeling assignment page (which has filled room
boxes containing bordered game tiles), so the two admin pages feel consistent.

## Visual target (copied from `pages/admin/rooms/rooms_assignment_page.templ`)

**Game card** — mirror `.room`:
- `background-color: var(--bg-item)`
- `padding: var(--spacing-2x)`
- `border-radius: var(--border-radius-1x)`
- keep `border: 1px solid var(--bg-item-border)`
- `transition: border-color 300ms ease-in-out`
- `:hover { border-color: var(--bg-item-border-hover); }`

**Player tile** — mirror `.room-event`, minus the `<img>` and `::before` gradient:
- `border: 2px solid var(--bg-item-border)`
- `border-radius: var(--border-radius-2x)`
- `padding: var(--spacing-1x)`
- `transition: border-color 300ms ease-in-out`
- `:hover { border-color: var(--color-primary); }`
- **background: transparent** (border-only — lets the card's `--bg-item` show
  through, matching how room-events are defined by their border, not a fill)

## Scope of change

In `pages/admin/puljefordeling_tab.templ`, `puljefordelingTabStyles`:

- Extend `.puljefordeling-tab-event` with the background / padding / radius /
  transition / hover above.
- Extend `.puljefordeling-tab-players li` (the player tile) with the border /
  radius / padding / transition / hover above. The existing rules on that
  selector stay: `display: flex; align-items: center; gap; justify-content:
  flex-start`.

Preserved exactly (must not regress):
- The `.puljefordeling-tab-moved` red left-border for relocated players.
- The `.puljefordeling-tab-dm` turquoise color for SL-elsewhere players.
- The 📌 pin marker and the `×` remove button (and its open/locked gating).
- The "+" add button placement, the legend, the header, the rerun button.
- The existing `.puljefordeling-tab-grid` (auto-fit columns) — unchanged.

## Out of scope

- No image banners on player tiles (players have no event image).
- No grid-column / responsive-breakpoint changes (grid stays auto-fit).
- No markup restructuring beyond what the tile styling needs; no behavior,
  routing, endpoint, copy, or data changes.

## Testing

- This is a visual-only change. The existing render tests
  (`TestPuljefordelingTabContent_*`) must still pass — in particular the "×"
  gating and "+" presence assertions, confirming the markup the styles attach to
  is unchanged structurally.
- Manual smoke test: the game cards show a filled background with rounded corners
  and a hover border; each player renders as a bordered tile with a hover
  highlight; moved (red border) and SL (turquoise) markers still render.
- Verification gate: `go test ./...` and `golangci-lint run` both clean.
