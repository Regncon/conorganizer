# Puljefordeling Bordfordeling Visual Parity Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Restyle the puljefordeling tab page to match the bordfordeling (rooms-assignment) page — a tab-connected panel, responsive grid of room-style game boxes with bordered player tiles, and removal of the emulator-era chrome (legend, page stats line, soft-grey flatness).

**Architecture:** CSS + minor markup edits confined to `pages/admin/puljefordeling_tab.templ`. The bordfordeling style values are copied from `pages/admin/rooms/rooms_index.templ` and `rooms_assignment_page.templ` into this file's two `<style>` blocks (the `puljefordelingTabIndex` layout block and the `puljefordelingTabStyles` block), matching the codebase's per-templ `<style>` convention. No behavior, routing, endpoint, or data changes.

**Tech Stack:** templ (`go tool templ generate`), CSS custom properties + native CSS nesting (already used in the rooms `<style>` blocks).

## Global Constraints

- Generate with `go tool templ generate -path pages/admin -log-level error` (templ is a Go tool, not on PATH).
- `*_templ.go` generated files are gitignored repo-wide — commit only the `.templ` source.
- Visual/markup-tidy change only: no behavior, routing, endpoint, copy-meaning, or data changes.
- Preserve exactly: "+"/"×" (and open/locked gating), the shared picker dialog, the rerun endpoint + button, the status-toggle controls, the live SSE re-render, the level emoji, 📌 pin, red moved accent (`.puljefordeling-tab-moved`), turquoise SL (`.puljefordeling-tab-dm`), the error and "Ingen deltakere" states.
- Remove: the legend line, the page-level stats line, and the soft-grey flat styling.
- Norwegian (Bokmål) copy meaning unchanged.
- Verification gate: `go test ./...` and `golangci-lint run` both clean.

---

## Task 1: Bordfordeling restyle of the puljefordeling tab

**Files:**
- Modify: `pages/admin/puljefordeling_tab.templ`

**Interfaces:**
- Consumes: existing templ funcs `PuljefordelingTabContent`, `puljefordelingTabEventCard`, `puljefordelingTabIndex`, `puljefordelingTabStyles`, `puljefordelingTabLegend`.
- Produces: nothing new; restyles and tidies existing markup. `puljefordelingTabLegend` is deleted.

- [ ] **Step 1: Add the connected-panel CSS to the tabs layout**

In `puljefordelingTabIndex`'s `<style>` block, the last rule is `.tabs a.active { ... }`. Immediately after its closing `}` (and before `</style>`), add the `.rooms-section` panel rule (values from `rooms_index.templ:101-109`):

```css
            .rooms-section {
                display: flex;
                flex-direction: column;
                gap: var(--responsive-page-section-spacing);
                border: 2px solid var(--color-primary);
                border-radius: 0 0 var(--border-radius-2x) var(--border-radius-2x);
                padding: var(--responsive-pane-padding);
                background-color: var(--bg-surface);
            }
```

(The content section already has `class="rooms-section"`, so it now picks this up and joins the active tab — which already renders `background: var(--bg-surface); border-color: var(--color-primary); border-bottom: none`.)

- [ ] **Step 2: Remove the page-level stats line from the header**

In `PuljefordelingTabContent`, replace this header heading block:

```go
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
```

with just the heading:

```go
				<div class="puljefordeling-tab-heading">
					<h1 class="page-heading">{ em.Name }</h1>
				</div>
```

> `meta` is now unused in this templ. After this edit, change the binding block at the top of `PuljefordelingTabContent` from `em, meta, emErr := pf.EmulatePulje(db, pulje)` to `em, _, emErr := pf.EmulatePulje(db, pulje)` so the build stays clean.

- [ ] **Step 3: Remove the legend (call + templ)**

In `PuljefordelingTabContent`, delete this line:

```go
			@puljefordelingTabLegend()
```

Then delete the entire `puljefordelingTabLegend` templ function:

```go
templ puljefordelingTabLegend() {
	<p class="puljefordeling-tab-legend">
		🔥 Veldig interessert · 👍 Middels · 🤷 Litt · 📌 Manuelt plassert ·
		<span class="puljefordeling-tab-dm">Turkis</span> = spilleder et annet sted ·
		<span class="puljefordeling-tab-moved-legend">Rød strek</span> = flyttet ned for å gi plass til andre
	</p>
}
```

- [ ] **Step 4: Restyle the game-card top into room-style rows**

In `puljefordelingTabEventCard`, replace this block:

```go
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
```

with room-style centered title + label/value `span` rows:

```go
		<h3>{ ev.Title }</h3>
		<span>
			<p>Spillere</p>
			<p>{ fmt.Sprintf("%d / %d", len(ev.AssignedPlayers), ev.Capacity) }</p>
		</span>
		if ev.GMName != "" {
			<span>
				<p>Spilleder</p>
				<p>{ ev.GMName }</p>
			</span>
		}
		if ev.Undersubscribed {
			<span class="puljefordeling-tab-badge">Få deltakere</span>
		}
```

- [ ] **Step 5: Update the grid, card, and player-tile styles; drop dead selectors**

In `puljefordelingTabStyles`, make the following changes.

(5a) Replace the grid rule:

```css
		.puljefordeling-tab-grid {
			display: grid;
			grid-template-columns: repeat(auto-fit, minmax(var(--mobile-min-width), 1fr));
			gap: var(--spacing-4x);
		}
```

with bordfordeling's responsive columns (from `rooms_assignment_page.templ:262-312`):

```css
		.puljefordeling-tab-grid {
			display: grid;
			grid-template-columns: 1fr;
			gap: var(--responsive-page-section-spacing);
		}
		@media screen and (width > 1200px) {
			.puljefordeling-tab-grid {
				grid-template-columns: repeat(2, 1fr);
			}
		}
		@media screen and (width > 1800px) {
			.puljefordeling-tab-grid {
				grid-template-columns: repeat(3, 1fr);
			}
		}
```

(5b) Replace the game-card rule (add centered `h3` and `span` row layout via nesting, mirroring `.room`):

```css
		.puljefordeling-tab-event {
			display: flex;
			flex-direction: column;
			gap: var(--spacing-2x);
			padding: var(--spacing-2x);
			border: 1px solid var(--bg-item-border);
			border-radius: var(--border-radius-1x);
			background-color: var(--bg-item);
			transition: border-color 300ms ease-in-out;
		}
		.puljefordeling-tab-event:hover {
			border-color: var(--bg-item-border-hover);
		}
```

with:

```css
		.puljefordeling-tab-event {
			display: flex;
			flex-flow: column nowrap;
			gap: var(--spacing-1x);
			padding: var(--spacing-2x);
			border: 1px solid var(--bg-item-border);
			border-radius: var(--border-radius-1x);
			background-color: var(--bg-item);
			transition: border-color 300ms ease-in-out;
		}
		.puljefordeling-tab-event h3 {
			text-align: center;
			color: var(--color-text-strong);
			font-size: var(--text-body);
		}
		.puljefordeling-tab-event > span {
			display: flex;
			justify-content: space-between;
			gap: var(--spacing-1x);
		}
		.puljefordeling-tab-event:hover {
			border-color: var(--bg-item-border-hover);
		}
```

(5c) Delete these now-dead rules (their markup/classes were removed in Steps 2-4):

```css
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
```

```css
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
```

```css
		.puljefordeling-tab-moved-legend {
			border-left: 3px solid #e05050;
			padding-left: var(--spacing-1x);
		}
```

(5d) Drop the soft-grey default on the player list — replace:

```css
		.puljefordeling-tab-players {
			margin: 0;
			padding: 0;
			list-style: none;
			display: flex;
			flex-direction: column;
			gap: var(--spacing-1x);
			color: var(--color-text-soft);
		}
```

with (no `color` override; inherits the normal text color):

```css
		.puljefordeling-tab-players {
			margin: 0;
			padding: 0;
			list-style: none;
			display: flex;
			flex-direction: column;
			gap: var(--spacing-1x);
		}
```

- [ ] **Step 6: Generate templates and run the existing render tests**

The restyle attaches to the same structural classes/buttons the tests assert on (`.puljefordeling-tab-remove`, `.puljefordeling-tab-add`, `.puljefordeling-tab-rerun`, the `h3` title, the data-init URL, the status-control label), so they must still pass.

Run:
```bash
go tool templ generate -path pages/admin -log-level error
go test ./pages/admin/ -run 'TestPuljefordelingTabContent|TestRerun' -v
```
Expected: PASS. If a compile error reports `puljefordelingTabLegend` still referenced, an instance of the call in Step 3 was missed — remove it. If `meta declared and not used`, apply the `em, _, emErr` change from Step 2.

- [ ] **Step 7: Run the full package suite**

Run: `go test ./pages/admin/...`
Expected: all PASS.

- [ ] **Step 8: Commit**

```bash
git add pages/admin/puljefordeling_tab.templ
git commit -m "style(puljefordeling): match bordfordeling page — connected panel, room boxes, drop emulator chrome"
```

---

## Task 2: Verification gate and visual smoke test

**Files:** none (verification only)

- [ ] **Step 1: Full suite + linter**

Run:
```bash
go test ./...
go tool templ generate -path pages/admin -log-level error && golangci-lint run
```
Expected: all tests PASS; `golangci-lint` → `0 issues` (watch for an unused-variable finding if Step 2's `meta` change was missed).

- [ ] **Step 2: Visual smoke test**

```bash
go tool task start
```

At `http://localhost:7331/admin/puljefordeling/FredagKveld`, compare side-by-side with `http://localhost:7331/admin/rooms/assignment/FredagKveld`:
- The active tab joins a primary-bordered surface panel (one connected surface), like the rooms page.
- A responsive grid (1 col, 2 above 1200px, 3 above 1800px) of filled game boxes; each box has a centered title and right-aligned count rows ("Spillere N / cap", "Spilleder …").
- Player tiles are bordered with a hover highlight, normal text color (no grey flatness); level emoji, 📌, red moved accent, turquoise SL all intact.
- No legend line, no page-level stats line.
- "+/×", the picker dialog, "Rerun fordeling" (when locked), and the status toggles all still work; the live re-render after an add/remove still fires.

---

## Self-Review Notes

- **Spec coverage:** connected panel (Step 1), responsive grid (5a), room-style boxes (Step 4 + 5b), player tile color drop (5d), legend removal (Step 3 + 5c), page stats removal (Step 2 + 5c), dead-selector cleanup (5c), preserved markers/behavior (constraints + untouched `.puljefordeling-tab-moved`/`-dm`/`-players li`/`-remove`/`-add`/`-rerun`), verification gate (Task 2).
- **Placeholder scan:** every CSS/markup value is concrete, copied from the rooms source or the current file.
- **Consistency:** all selector names (`.rooms-section`, `.puljefordeling-tab-grid`, `.puljefordeling-tab-event`, `.puljefordeling-tab-players`) match the current file, verified against the live `puljefordelingTabStyles` and `puljefordelingTabIndex` blocks. The `meta` → `_` change (Step 2) is the only Go-side adjustment and is called out in Steps 2 and 6.
