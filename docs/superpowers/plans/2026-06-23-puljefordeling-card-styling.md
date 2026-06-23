# Puljefordeling Card Styling (romfordeling boxes) Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Restyle the puljefordeling tab's game cards and player rows to match the romfordeling "table selection" boxes — filled `.room`-style game cards and bordered `.room-event`-style player tiles.

**Architecture:** Pure CSS edits to the `puljefordelingTabStyles` `<style>` block in `pages/admin/puljefordeling_tab.templ`, copying the box treatment values verbatim from `pages/admin/rooms/rooms_assignment_page.templ`. No markup, behavior, routing, copy, or data changes.

**Tech Stack:** templ (`go tool templ generate`), CSS custom properties.

## Global Constraints

- `go tool templ generate -path pages/admin -log-level error` to regenerate (templ is a Go tool, not on PATH).
- `*_templ.go` generated files are gitignored repo-wide — commit only the `.templ` source.
- Visual-only change: no markup restructuring, no behavior/routing/copy/data changes.
- Preserve exactly: the `.puljefordeling-tab-moved` red left-border, `.puljefordeling-tab-dm` turquoise, the 📌 pin, the `×` remove button and its gating, the `+` button, the legend/header/rerun button, and the `.puljefordeling-tab-grid` (auto-fit columns).
- Player tiles use a transparent background (border-only), letting the card's `var(--bg-item)` show through.
- Verification gate: `go test ./...` and `golangci-lint run` both clean.

---

## Task 1: Box styling for game cards and player tiles

**Files:**
- Modify: `pages/admin/puljefordeling_tab.templ` (the `puljefordelingTabStyles` `<style>` block)

**Interfaces:**
- Consumes: existing class names `.puljefordeling-tab-event`, `.puljefordeling-tab-players li` (rendered by `puljefordelingTabEventCard`).
- Produces: nothing new — restyles existing selectors and adds two `:hover` rules.

- [ ] **Step 1: Restyle the game card (`.puljefordeling-tab-event`)**

In `pages/admin/puljefordeling_tab.templ`, replace this rule:

```css
		.puljefordeling-tab-event {
			display: flex;
			flex-direction: column;
			gap: var(--spacing-2x);
			border: 1px solid var(--bg-item-border);
		}
```

with (adds background / padding / radius / transition, plus a hover rule — values mirror `.room`):

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

- [ ] **Step 2: Restyle the player tile (`.puljefordeling-tab-players li`)**

Replace this rule:

```css
		.puljefordeling-tab-players li {
			display: flex;
			align-items: center;
			gap: var(--spacing-2x);
			justify-content: flex-start;
		}
```

with (adds border / radius / padding / transition, plus a hover rule — values mirror `.room-event`, background stays transparent):

```css
		.puljefordeling-tab-players li {
			display: flex;
			align-items: center;
			gap: var(--spacing-2x);
			justify-content: flex-start;
			padding: var(--spacing-1x);
			border: 2px solid var(--bg-item-border);
			border-radius: var(--border-radius-2x);
			transition: border-color 300ms ease-in-out;
		}
		.puljefordeling-tab-players li:hover {
			border-color: var(--color-primary);
		}
```

(The existing `.puljefordeling-tab-moved` rule still applies its `border-left: 3px solid #e05050` and `padding-left` on top of this tile border — the red "moved" accent is preserved by cascade.)

- [ ] **Step 3: Generate templates and run the existing render tests**

The styles attach to unchanged markup, so the existing structural tests must still pass.

Run:
```bash
go tool templ generate -path pages/admin -log-level error
go test ./pages/admin/ -run TestPuljefordelingTabContent -v
```
Expected: PASS (the existing `TestPuljefordelingTabContent_*` tests — "×" gating, "+" presence, rerun button visibility — unchanged).

- [ ] **Step 4: Run the full package suite**

Run: `go test ./pages/admin/...`
Expected: all PASS.

- [ ] **Step 5: Commit**

```bash
git add pages/admin/puljefordeling_tab.templ
git commit -m "style(puljefordeling): romfordeling-style boxes for game cards and player tiles"
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
Expected: all tests PASS; `golangci-lint` → `0 issues`.

- [ ] **Step 2: Visual smoke test**

```bash
go tool task start
```

At `http://localhost:7331/admin/puljefordeling/FredagKveld`:
- Each game card has a filled background (`var(--bg-item)`), rounded corners, and a hover border highlight.
- Each assigned player is a bordered tile with a hover highlight; the emoji · name · 📌 · × layout is intact.
- A relocated player still shows the red left accent; an SL-elsewhere player still shows turquoise; the legend, "+", and (when locked) "Rerun fordeling" are unchanged.

---

## Self-Review Notes

- **Spec coverage:** game-card `.room` treatment (Task 1 Step 1), player-tile `.room-event` treatment with transparent background (Task 1 Step 2), preserved markers/behavior (constraints + the cascade note), existing tests still pass (Task 1 Steps 3-4), verification gate (Task 2).
- **No placeholders:** every CSS value is concrete and copied from the romfordeling source or the spec.
- **Consistency:** selectors `.puljefordeling-tab-event` and `.puljefordeling-tab-players li` match the names already in the file (verified against the current `puljefordelingTabStyles` block).
