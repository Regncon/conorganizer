# Morph safety

A live fragment is re-rendered from the database and morphed over the DOM many times while a page is open: on every broadcast to a subscribed bucket, when the tab becomes visible again (the GET stream is aborted while hidden and reopened), and when any other tab in the same browser touches the shared connection key. Server-rendered HTML never knows what the browser did in between. Everything the browser owns must therefore survive a morph by design, not by luck.

## How a patch lands

`Stream` sends `sse.PatchElementTempl(render())`. The default mode is an outer morph: the top-level element of the patch is matched by id against the DOM, and the subtree is reconciled attribute by attribute and node by node. Attributes present in the DOM but absent from the new HTML are removed, which is what strips `open` from a dialog or details element.

## Native `<dialog>`

A `<dialog>` opened with `showModal()` sits in the top layer with the rest of the document inert. If a morph removes its `open` attribute, the browser does not run `close()`: the page can stay inert with no visible dialog, and nothing is clickable until a reload. The lifecycle doc describes this under "Native dialogs and live fragments".

An invoker button with `commandfor="x"` and `command="show-modal"` is `showModal()` too and has the same problem, with narrower browser support on top. `command="close"` is harmless.

Three options, in order of preference:

1. Render the dialog outside the live wrapper. Room assignment keeps `#room-assignment-page` as a stable outer wrapper with the signals and the `<dialog id="assignment-dialog">`, and streams only the `#room-assignment` section inside it (`pages/admin/rooms/rooms_assignment_page.templ`). Puljefordeling renders `#tildeling-dialog` in `puljefordelingIndex`, outside the streamed `#puljefordeling-page`. The rooms page renders `formModal()` with both dialogs outside `RoomsPageContent`. A page without such an outer wrapper does not need one: render the dialog as a sibling right after the live `<div id="..." data-init=...>` in the `_index` templ. Signals are global, so the dialog can declare the ones it needs itself with `data-signals:_name__ifmissing="''"` or `data-signals__ifmissing="{_name: ''}"`. The trigger button may stay inside the live fragment and call `showModal()` on the outside dialog; what matters is where the `<dialog>` element sits. Pass the item the dialog acts on through signals set by the trigger, and keep user text out of the expression (`data-*` attribute plus `el.dataset`, see "User text in expressions" below). Have the trigger overwrite every value the dialog uses each time it opens it. Clearing them in `data-on:close` races with a quick reopen, because the `close` event can arrive after the dialog was opened again; if you clear them there, guard it with `if (!el.open) { ... }`.

   Close the dialog from the mutation's response, after the write succeeded: `sse.ExecuteScript` with `close()`, as the room dialog in `pages/admin/admin.go` does, or a patch of the signal that drives it. Put the trigger's `errorfeedback.RenderMessage` inside the dialog, so a failure stays visible while the dialog is open.
2. Keep it inside the fragment but drive it from a signal (`datastar-signals.md`, "Native modal dialogs in live fragments"):

```html
<div data-signals:billettholder-interest-open-dialog-id__ifmissing="''">
    <button data-on:click="$billettholderInterestOpenDialogId = 'billettholder-interests-42'">Interesser</button>
    <dialog
        id="billettholder-interests-42"
        data-preserve-attr="open"
        data-effect="if ($billettholderInterestOpenDialogId == 'billettholder-interests-42' && !el.open) { el.showModal(); } if ($billettholderInterestOpenDialogId != 'billettholder-interests-42' && el.open) { el.close(); }"
        data-on:close="if ($billettholderInterestOpenDialogId == 'billettholder-interests-42') { $billettholderInterestOpenDialogId = ''; }"
    >
        ...
    </dialog>
</div>
```

   The signal is declared where a morph keeps it: `data-signals:name__ifmissing` on a container, the object form `data-signals__ifmissing="{_name: ''}"` as `pages/event/event_rooms.templ` writes it, or an element outside the live wrapper, as the `<header>` in `components/header/menu.templ` does. `data-signals__ifmissing:name`, with the modifier before the key, declares nothing; a dialog written that way only works because its open button creates the signal by assigning it. Buttons only set the signal. The effect opens and closes. `data-on:close` catches Escape and `form method="dialog"`. `data-preserve-attr="open"` stops the morph from stripping the attribute between the patch and the effect. Real uses: `pages/admin/billettholder_admin/billettholder_interest_dialog.templ`, `pages/event/event_rooms.templ`, `pages/admin/puljefordeling_assignment_dialog.templ`, `pages/admin/tildelingsdialog.templ`, the billettholder switcher in `components/header/menu.templ`.
3. The "Meld interesse" dialog in `components/event_components/event_interests.templ` is the exception: it binds `data-attr:open="$_interestIsOpen ? true : false"` and never calls `showModal()`, so it has no top layer, no inert background and no focus trap. Do not copy it for new dialogs.

The bug pattern is a `<dialog>` inside a live fragment that is opened with `showModal()` or `command="show-modal"`, from anywhere, without the signal pattern. Grep a changed templ file for `showModal` and `show-modal`, find the dialog each one targets, and check which wrapper the dialog sits in.

## User text in expressions

Datastar attributes are JavaScript expressions. A name or e-mail address written into one with `fmt.Sprintf("... '%s' ...", value)` breaks on an apostrophe and lets user input run as code; Go's `%q` quotes for Go, not for JavaScript. Keep user text out of expressions: render it into a `data-*` attribute on the trigger and read it as `el.dataset.name` in the expression, or JSON-encode it with `templ.JSONString` or `components.DataSignals` when it is a signal value. Ids and other values you control may be formatted in.

## `<details>` and other browser-owned attributes

List the attribute in `data-preserve-attr`:

```html
<details class="pulje-interests-collapse" data-preserve-attr="open">
```

Several attributes are space separated: the header menu keeps `data-preserve-attr="class aria-busy"`. This only protects state across morphs; a reload resets it. State that must survive a reload belongs in the URL, `sessionStorage` or the database.

Selections and toggles kept in signals need `data-signals:name__ifmissing` or the object form `data-signals__ifmissing="{...}"` on an element that survives the patch. With plain `data-signals`, the server's value wins whenever a patch renders a different value. The billettholder interest dialog keeps its active pulje tab in `$billettholderInterestActivePulje` this way, and its open button sets the signal explicitly to the billettholder's first pulje.

## `data-ignore-morph`

Marks a region the parent patch must leave alone. Text and markup inside it never receive parent updates either, so it is only right when something else owns the region: a widget, or an endpoint that replaces the wrapper explicitly with `datastar.WithModeReplace()` (and `WithSelector` when the id is not on the patch root). `InterestContent` in `event_interests.templ` is the model: the page morph skips it, and the selected-interest endpoint replaces it after validating the current selection. Keep the controls that trigger that endpoint outside the replaced region, and render every initial value the region needs in the replacement HTML.

## Effects

`data-effect` runs once when the element appears and again whenever a signal it reads changes.

- A request inside an effect whose response changes the signals the effect reads re-triggers the effect. Limit the explicit dependencies to the values that should trigger a request, and keep the response from touching those values. The event page's selected-interest effect reads only `$billettHolderId` and `$puljeId`.
- A live patch that changes the text of a `data-effect` attribute starts the new effect before the old one is disposed. In the header menu the two effects talked to each other through `billettholder-selection-change` until the stack overflowed; the fix was `data-init` on a hidden span for the one-shot setup. Rule: one-shot setup goes in `data-init`; reactions to signal changes go in `data-effect`. `data-init` runs when the element is inserted and again when a patch changes its expression. A patch with the same expression does not run it, so put the data the setup depends on into the expression, as `menuBillettholderSelectionInit` in `components/header/menu.templ` does.
- Effects that touch the DOM (`el.showModal()`, `document.body.classList`) must be idempotent, because the next patch reruns them.

## Reacting after a patch

There is no before-morph hook. To do something after the live fragment was replaced, listen on the document and filter on the fragment id, as the Meld interesse dialog does to refresh its selection:

```html
data-on:datastar-fetch__document="if (el.isConnected && evt.detail.type === 'datastar-patch-elements' && evt.detail.el?.id === 'event-container') { ... }"
```

`el.isConnected` matters because the element that registered the listener may itself have been morphed away.

## Tooltips

`layouts/base.templ` creates one delegated Tippy instance on `document.body` with `target: "[data-tippy-content]"` and reads the content in `onShow` from the attribute. Morphed elements keep their tooltips because nothing is attached to the element itself, and a changed `data-tippy-content` shows the new text on the next hover. Write `data-tippy-content="..."` (plus `data-tippy-theme` or `data-attr:data-tippy-trigger` as `components/interest_indicator.templ` does) and never create per-element `tippy(el)` instances in page scripts inside a live fragment; the morph replaces the element and leaves the instance orphaned.

## Scripts and web components inside fragments

- A `<script>` node the morph keeps does not run again; a `<script>` node a patch inserts runs once (that is how `sse.ExecuteScript` works: it appends a script element to `body`). Do not depend on either for behaviour inside a live fragment; put it in `data-on`, `data-init` or a web component.
- Web components keep their JS state through a morph as long as the element is kept (same tag, same id). Attributes they observe are re-applied from the new HTML, so pass state down through attributes bound with `data-attr` and send changes up with custom events, as AGENTS.md describes.

## Shared session keys across tabs

All tabs in one browser share the `connections` cookie, so they watch the same key per bucket. Opening or refocusing any live page touches that key in its buckets and re-renders every other tab that watches an overlapping bucket, even though no data changed. Because the header menu streams `billettholders` on every page, every tab that watches `billettholders` re-renders when any page opens. This is why renderers must be idempotent and why the patterns above are not optional for rarely-changing data.
