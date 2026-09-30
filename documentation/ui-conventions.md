# UI conventions

These are shared conventions for templ components, CSS and user-facing HTML. The rules for writing CSS are in [AGENTS.md](../AGENTS.md#css-rules): native nesting, container queries for component responsiveness, range syntax, logical properties and so on. This page covers conventions that are specific to this codebase.

## CSS conventions

### Prefer `.surface-pane` over `.card`

`static/css/index.css` defines two panel styles:

- `.surface-pane` / `.surface-pane-title`: the current panel style. Use it, or a specific local class, in new code.
- `.form-card` / `.form-card-title`: marked deprecated in favour of `.surface-pane`. The two share one rule, and `.form-card` is still used in `components/formsubmission/`.

The global `.card` class (`--level-1-background`, `border-radius-1x`, a shadow and `1rem` padding) is legacy and should not be used in new code. It is still used in:

- `pages/event/event_interest_panel.templ` (`card event-interest-picker`)
- `pages/profile/tickets/tickets_page.templ` (the info text and `profile-ticket-empty-state`)
- `pages/profile/tickets/billettholder_profile_card.templ`
- `pages/admin/billettholder_admin/billettholder_card.templ`
- `pages/admin/billettholder_admin/add/ticket_card.templ`

Remove `.card` from `index.css` only after all of these have been migrated.

### Shared form and choice classes

- `.label`, `.label-bold` and `.form-group` (the label + control + message stack) are global rules in `static/css/index.css`, so they work outside `components/formsubmission/` too.
- `static/css/choice.css` (loaded by `layouts/base.templ` after `buttons.css`) holds controls that switch an option on or off. `.choice-toggle-group` with `.choice-toggle` buttons (used together with `.btn .btn--secondary`) picks one of a few options. The buttons stack, and sit side by side in equal columns when the nearest container is wider than 25rem. `.choice-chip-group` with `.choice-chip` buttons gives wrapping pill-shaped chips for picking any number of options. The selected look follows `aria-pressed="true"`, so the markup only binds `aria-pressed`. The feedback form (`pages/feedback/feedback_page.templ`) uses both.

## Focus-visible and accessibility

Every interactive element needs a visible `:focus-visible` state, not only a `:hover` state (GitHub issue Regncon/conorganizer#407). Use the soft focus colour (`--color-text-soft-50`, or `--btn-secondary-focus-shadow` / `--btn-ghost-focus-shadow`, which point to it) for focus rings. Primary elements use `--color-primary-focus-visible`.

### Reset and global fallback

`static/css/index.css` has a universal reset that removes outlines:

```css
*, *::before, *::after { box-sizing: inherit; outline: none; }
```

A global fallback comes right after it:

```css
:where(a[href], button, summary, [role="button"], input:not([type="hidden"]), select, textarea):focus-visible {
    outline: 2px solid var(--color-text-soft-50);
    outline-offset: 2px;
}
```

- Keep the fallback directly after the reset, so the two rules are read together. The fallback wins over the reset through specificity, not file order: `*` has specificity (0,0,0), and the `:focus-visible` outside `:where()` gives the fallback (0,1,0).
- The element list is wrapped in `:where()` so the fallback stays at (0,1,0). Any component rule with a class and `:focus-visible`, such as `.btn--*`, `.input`, the cards and the header, has higher specificity and overrides it.

### Per-component focus states

| Element | Where | Focus-visible style |
| --- | --- | --- |
| Header logo link `.logo-link` | `components/header/menu.templ` | `border-radius: 50%`; on focus, `outline: none`, the ghost-button hover background and `box-shadow: var(--btn-ghost-shadow)` (a `3px` `--btn-ghost-focus-shadow` ring). `box-shadow` follows the border radius, so the ring is round. |
| Header dropdown links and buttons | `components/header/menu.templ` | Shared hover/focus background, plus an inset `box-shadow` ring on focused links. |
| `.event-card-container`, `.event-bar-container` | `static/css/card.css`, used by `components/event_card.templ` and the admin event bars | The Level 2 card states: `--level-2-background-hover` and `--level-2-border-hover` on hover and focus, `--level-2-background-active` on `:active`, plus a `3px` `--color-text-soft-50` outline on focus. |
| `.profile-event-bar` | `static/css/card.css` | The Level 3 card states: `--level-3-background-hover` and `--level-3-border-hover` on hover and focus, plus a `3px` `--color-text-soft-50` outline with `outline-offset: 2px` on focus. |
| `.card-clickable--level-2`, `.card-clickable--level-3` | `static/css/card.css`, used by the ticket holder pickers and the header menu | The level's hover background and border on hover and focus, plus a `3px` `--color-text-soft-50` outline on focus when the card is not `.selected`. |
| `a.inline-link` | global, `static/css/index.css` | `2px` `--color-text-soft-50` outline with offset and radius. Used for example by the help link in `event_interest_panel.templ`. |
| `.pulje-interests-summary` (`<summary>`) | `components/event_components/programpulje_interests.templ` | Same background as hover, plus a `3px` `--btn-ghost-focus-shadow` ring on the inner text span. |
| `.choice-toggle`, `.choice-chip` | `static/css/choice.css` | `box-shadow: var(--choice-shadow)` (a `3px` `--color-text-soft-50` ring) on keyboard focus only (a mouse click gets no ring). Selected buttons (`aria-pressed="true"`) use `--choice-selected-shadow` (`--color-primary-focus-visible`) instead. `.choice-chip` also gets the Level 3 hover background and border on focus. |
| `a.admin-tool-card` | `pages/admin/admin_card.templ` | The Level 2 card states: the same background and border as hover, `--level-2-background-active` on `:active`, plus the global `2px` outline fallback (the card rule sets no outline of its own). |

The logo image's `alt` is "Regncon forside", because the image is the link's only accessible name. Give an image-only link an `alt` that describes where the link goes.

A `<details>` element that live updates re-render should have `data-preserve-attr="open"` (as `.pulje-interests-collapse` and the room assignment notes do), so a Datastar morph doesn't close it while the user is reading it.

## External links

A link that opens another site in a new tab uses `target="_blank" rel="noopener noreferrer"` and ends with the external-link icon:

```templ
<a href="https://www.regncon.no/vanlege-sporsmal/" target="_blank" rel="noopener noreferrer">
	Vanlige spørsmål
	@icons.Icon(icons.ExternalLink, icons.Size24)
</a>
```

`icons.ExternalLink` maps to `components/icons/assets/icon-external-link.svg`. It is currently used by:

- "Kjøp billetter" (the CheckIn link) on `/profile/tickets`,
- both "Vanlige spørsmål" links in the header menu (desktop and mobile),
- "Hvordan fungerer interessevalget?" on the event page (`Size16`),
- "Hent utskriftsvennlig versjon av programmet" on the Publiser program card on `/admin` (`Size16`). This is an exception: the link goes to the same-site `/print` and opens in the same tab, without `target` or `rel`.

The room map link on `/admin/rooms` is same-site. It opens the map SVG in a new tab with `rel="noopener"` and no icon. The link text says "i ny fane" instead.

Links in Markdown event descriptions do not open in a new tab. See [Markdown rendering](#markdown-rendering-and-sanitization).

## Event form and preview layout

The event form with its live preview is built from two components.

- **`formsubmission.EventFormPageLayout(title string, ultrawide bool)`** (`components/formsubmission/event_form_page_layout.templ`) owns the page shell: `page-content-container` (plus `ultrawide-content` when `ultrawide` is true), the `<h1>`, and the `.event-form-with-preview` grid with its CSS. It renders its children inside the grid.
- **`event.EventFormPreview(eventData, eventImageDir)`** (`pages/event/event_form_preview.templ`) renders the preview. It stays in `pages/event` because it depends on `Event_mobile`, so `components/formsubmission` never imports `pages/event`.

The pages compose the two explicitly:

| Page | Title | `ultrawide` | Preview |
| --- | --- | --- | --- |
| `pages/profile/newevent/new_page.templ` | "Nytt arrangement" | `false` | only for admins |
| `pages/admin/approval/editForm/edit_form_page.templ` | "Rediger arrangement" | `true` | always |

`EventFormPreview` sets the preview policy in one place. It wraps `<div class="event-container event-form-preview">` around `Event_mobile(eventData, components.PreviousNext{IsRemoved: true}, nil, eventImageDir)`. This removes the previous/next navigation, and because `interestPanel` is `nil`, no interest or ticket panel is shown. The profile and admin previews therefore look the same.

### Container query breakpoints

The pages wrap their content in `.formsubmission-css-container`, defined in `static/css/index.css` (`container-type: inline-size; container-name: formsubmission-css-container`). All breakpoints query that container:

| Breakpoint | Where | Effect |
| --- | --- | --- |
| `width > 600px` | `components/formsubmission/form_body.templ` | `.submit-section` becomes a two-column grid (`1fr 200px`). |
| `width > 920px` | `components/formsubmission/form_body.templ` | The wider form section layout, for example the organizer fields side by side. |
| `width > 1360px` | `event_form_page_layout.templ` | When a preview is present (`.event-form-with-preview:has(.event-form-preview)`), the grid becomes `minmax(880px, 1fr) minmax(var(--mobile-min-width), 1fr)`. |

The form + preview grid is one column by default. The 1360px breakpoint and the `880px` minimum make sure the form always keeps at least 880px before the preview moves beside it. Don't lower the breakpoint without checking the form for overlap.

## Markdown rendering and sanitization

Event descriptions are Markdown written by the organizer. They are rendered in `service/eventService/event_helpers.go`:

- `MdToHTML` renders with the gomarkdown HTML renderer, using the flags `html.CommonFlags | html.HrefTargetBlank`.
- `SanitizeMdToHTML` passes that output through `bluemonday.UGCPolicy()`.

Always render user Markdown through `SanitizeMdToHTML`, never through `MdToHTML` alone.

### Security contract

`service/eventService/event_helpers_test.go` covers this contract:

- **Kept on purpose:** safe Markdown and safe raw HTML, such as headings (with auto heading IDs), `<strong>`, `https` links, `<div>…</div>` and `<img src alt>`.
- **Stripped:** executable content, which means `<script>` elements, `javascript:` link URLs and `on*` event-handler attributes such as `onerror`.

### Link attributes

The renderer adds `target="_blank"` to links, but `UGCPolicy` doesn't allow `target`, so it is stripped. The policy then adds `rel="nofollow"` to links. As a result, links in descriptions open in the same tab.

Because the final attributes differ from the renderer output, tests for `SanitizeMdToHTML` do not compare the full output byte for byte. They assert on selected fragments (for example `<a href="https://example.com"` and `>event page</a>`), plus a list of forbidden substrings that are matched case-insensitively (`assertStringContainsAll` / `assertStringExcludesAll`). Follow the same pattern when you add cases.

### No `<p>` wrapper

The sanitized description can contain block elements (`<p>`, headings, lists), so it must not be wrapped in a `<p>`. `pages/event/event_mobile.templ` renders it directly in a `<div>`:

```templ
<div class="description">
	@templ.Raw(string(eventservice.SanitizeMdToHTML([]byte(event.Description))))
</div>
```

Paragraph spacing comes from nested `.description p` rules in the same template.
