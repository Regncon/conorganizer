# Puljefordeling: full bordfordeling visual parity

**Date:** 2026-06-23
**Status:** Approved design, ready for implementation plan
**Builds on:** `2026-06-23-puljefordeling-card-styling-design.md`

## Summary

Restyle the puljefordeling tab page to look as close to the bordfordeling
(rooms-assignment) page as possible: a connected panel that joins the tabs, a
responsive grid of filled "room"-style game boxes holding bordered
"room-event"-style player tiles, and removal of the emulator-era chrome (legend,
page-level stats line, soft-grey flat styling). Functional per-player markers
(interest-level emoji, 📌 pinned, red moved accent, turquoise SL) are kept,
restyled to fit the boxes. All changes are CSS plus minor markup in
`pages/admin/puljefordeling_tab.templ`; no behavior, routing, copy meaning,
endpoint, or data changes.

## Motivation

The puljefordeling content section already carries `class="rooms-section"`, but
the panel CSS that class depends on lives only on the rooms page, so the
puljefordeling page renders flat and visually disconnected from its tabs. The
team wants the two admin assignment pages to look like one consistent UI. Rather
than reference styles that are not loaded, this work brings the bordfordeling
visual vocabulary into the puljefordeling page.

## Visual targets (source: `pages/admin/rooms/rooms_index.templ` and `rooms_assignment_page.templ`)

### Connected panel — mirror `.rooms-section` (rooms_index.templ:101-109)

Bring this into the puljefordeling tabs layout (`puljefordelingTabIndex`, next to
the `.tabs` CSS it already has):

```css
display: flex;
flex-direction: column;
gap: var(--responsive-page-section-spacing);
border: 2px solid var(--color-primary);
border-radius: 0 0 var(--border-radius-2x) var(--border-radius-2x);
padding: var(--responsive-pane-padding);
background-color: var(--bg-surface);
```

The `.tabs` rules already copied include the active-tab merge
(`.tabs a.active { background: var(--bg-surface); border-color: var(--color-primary);
border-bottom: none; }` and `.tabs { margin-bottom: -2px; z-index: 2; }`), so the
active tab joins this panel as one surface. Keep the puljefordeling content
section's `class="rooms-section"` so it picks up this rule.

### Grid — mirror `.rooms-container` (rooms_assignment_page.templ:262-312)

Replace the current `auto-fit` grid with bordfordeling's responsive columns:

```css
display: grid;
grid-template-columns: 1fr;
gap: var(--responsive-page-section-spacing);
```
```css
@media screen and (width > 1200px) { grid-template-columns: repeat(2, 1fr); }
@media screen and (width > 1800px) { grid-template-columns: repeat(3, 1fr); }
```

### Game box — mirror `.room` (rooms_assignment_page.templ:272-300)

```css
display: flex;
flex-flow: column nowrap;
gap: var(--spacing-1x);
padding: var(--spacing-2x);
border: 1px solid var(--bg-item-border);
border-radius: var(--border-radius-1x);
background-color: var(--bg-item);
transition: border-color 300ms ease-in-out;
```
plus `h3 { text-align: center; }`, `span { display: flex; justify-content:
space-between; gap: var(--spacing-1x); }`, and `:hover { border-color:
var(--bg-item-border-hover); }`. Per-box counts (assigned/capacity, GM) render as
these right-aligned `span` rows like the room's "Tildelte spill / N" rows.

### Player tile — mirror `.room-event` (rooms_assignment_page.templ:321-380), minus the image

```css
position: relative;
display: flex;
flex-flow: row nowrap;
align-items: center;
gap: var(--spacing-2x);
padding: var(--spacing-1x);
border: 2px solid var(--bg-item-border);
border-radius: var(--border-radius-2x);
transition: border-color 300ms ease-in-out;
```
plus `:hover { border-color: var(--color-primary); }`. Keep the inner content:
level emoji · name · 📌 · ×. No `<img>`, no `::before` gradient. Use the normal
text color (drop the `--color-text-soft` soft-grey). Preserve the
`.puljefordeling-tab-moved` red left-accent and `.puljefordeling-tab-dm`
turquoise.

## Removed (emulator chrome)

- The legend line (🔥 Veldig interessert · 👍 … explanation) — deleted from
  `PuljefordelingTabContent` markup and its `.puljefordeling-tab-legend` style.
- The page-level stats line ("N deltakere med ønsker · N fikk førstevalg · N uten
  plass") at the top — deleted; per-box counts carry the information instead.
- The soft-grey flat styling: `color: var(--color-text-soft)` defaults and the
  flat (background-less, chrome-less) treatment on the section and lists.

## Header and "uten plass" inside the panel

- **Header:** mirror bordfordeling's section heading — a heading for the pulje
  (the pulje name) inside the panel. Keep the **status toggles** (Puljefordeling
  lukket / publisert) and the **Rerun fordeling** button (locked/completed only),
  restyled to sit cleanly in the panel header rather than the old emulator header
  block. No change to their behavior or endpoints.
- **Uten plass:** restyle the unseated-participants section to bordfordeling's
  "Eventer uten rom" treatment — a heading plus the list, inside the panel.

## Preserved exactly (must not regress)

- All behavior: "+", "×" (and open/locked gating), the shared picker dialog, the
  rerun endpoint + button, the status-toggle endpoint, the live SSE re-render.
- Functional markers: level emoji, 📌 pinned, red moved accent, turquoise SL.
- The `error` and `empty` ("Ingen deltakere") states.
- Norwegian (Bokmål) copy meaning.

## Out of scope

- No image banners on player tiles (players have no event image).
- No refactor to share CSS across the rooms and puljefordeling pages — the
  bordfordeling CSS values are copied into the puljefordeling page, matching the
  codebase's existing per-templ `<style>` convention. (A future shared-stylesheet
  extraction is a separate effort.)
- No behavior, routing, endpoint, or data changes.

## Testing

- Visual-only / markup-tidy change. The existing render tests
  (`TestPuljefordelingTabContent_*`: "×" gating, "+" presence, rerun button
  visibility) must still pass against the restyled markup. If removing the legend
  / stats line changes any text those tests assert on, update only the affected
  assertion (the tests target structural classes and buttons, not the legend).
- Manual smoke test: the page shows tabs joined to a primary-bordered surface
  panel; a responsive grid (1/2/3 columns) of filled game boxes with centered
  titles and right-aligned count rows; bordered player tiles with hover; no
  legend, no top stats line, no grey flatness; moved/SL/pin markers intact;
  "+/×", rerun, and status toggles all still work.
- Verification gate: `go test ./...` and `golangci-lint run` both clean.
