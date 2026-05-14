# Centralize Billettholder LocalStorage In `TicketHolderPicker`

## Summary
- Move all `selectedBillettHolder` initialization, validation, and persistence into [ticket_holder_picker.templ](/e:/programmering/Workspace/Web%20devlopment/conorganizer/components/ticket_holder/ticket_holder_picker.templ#L38).
- Make the button and dropdown pure selector children: props down, events up.
- Normalize the serialized selected-billettholder payload to `Id`, `Name`, and `Email` only.
- Keep `Color` available in render-only input data where the dropdown needs it for UI, but do not emit or persist it.

## Implementation Changes
- In [ticket_holder_picker.templ](/e:/programmering/Workspace/Web%20devlopment/conorganizer/components/ticket_holder/ticket_holder_picker.templ#L38) lines 38-115, replace the current split ownership:
  - remove `data-on:billettholder-selected="$billettHolderId = evt.detail;"`
  - replace `data-init` so the parent resolves the selected billettholder from:
    - valid `localStorage`
    - else valid `yourBillettHolder`
    - else first available associated billettholder
  - persist that resolved object once
  - set `$billettHolderId`
  - on later selection changes, persist from the parent

```templ
<div
	class="ticket-holder-container"
	data-signals:billett-holder-id=""
	data-on:billettholder-selected="
		$billettHolderId = evt.detail.Id;
		setSelectedBillettholderInLocalStorage(evt.detail);
	"
	data-init={ fmt.Sprintf("
		const selectedBillettholder = initializeSelectedBillettholder(%s, %s);
		$billettHolderId = selectedBillettholder?.Id ?? '';
	",
		templ.JSONString(associatedTicketHolders),
		templ.JSONString(yourBillettHolder),
	) }
>
```

- Keep the helper functions in the same file so `TicketHolderPicker` remains the only owner of selected billettholder persistence:

```javascript
const LSKey = {
	SelectedBilletHolder: "selectedBillettHolder",
};

function getSelectedBillettholderFromLocalStorage() {
	const raw = localStorage.getItem(LSKey.SelectedBilletHolder);
	if (!raw) return null;
	try {
		return JSON.parse(raw);
	} catch {
		return null;
	}
}

function setSelectedBillettholderInLocalStorage(billettholder) {
	if (!billettholder?.Id) return;
	localStorage.setItem(
		LSKey.SelectedBilletHolder,
		JSON.stringify({
			Id: billettholder.Id,
			Name: billettholder.Name ?? "",
			Email: billettholder.Email ?? "",
		}),
	);
}

function initializeSelectedBillettholder(associatedBillettholdere, yourBillettholder) {
	const stored = getSelectedBillettholderFromLocalStorage();
	const validStored = associatedBillettholdere.find((bh) => bh.Id === stored?.Id);
	if (validStored) {
		setSelectedBillettholderInLocalStorage(validStored);
		return validStored;
	}

	const validDefault = associatedBillettholdere.find((bh) => bh.Id === yourBillettholder?.Id);
	if (validDefault) {
		setSelectedBillettholderInLocalStorage(validDefault);
		return validDefault;
	}

	const fallback = associatedBillettholdere[0] ?? null;
	if (fallback) {
		setSelectedBillettholderInLocalStorage(fallback);
	}
	return fallback;
}
```

- In the same file, replace the current button-specific query/listener block with parent-side delegated persistence, or remove it entirely if the new parent `data-on:billettholder-selected` covers both selector variants. The parent must remain the only writer to `localStorage`.

- In [ticket_holder_dropdown.templ](/e:/programmering/Workspace/Web%20devlopment/conorganizer/components/ticket_holder/ticket_holder_dropdown.templ#L8) lines 8-16:
  - remove unused params `eventId` and `userInfo`
  - pass selected state down from the parent with `data-attr:selected-id="$billettHolderId"`

```templ
templ TicketHolderDropDown(associatedBillettholder []BillettHolder) {
	<billettholder-dropdown
		class="custom-select"
		data-billettholdere={ templ.JSONString(associatedBillettholder) }
		data-attr:selected-id="$billettHolderId"
	>
		<template data-arrow-icon>
			@icons.Icon(icons.ChevronDown, icons.Size24)
		</template>
	</billettholder-dropdown>
}
```

- Update the call site in [ticket_holder_picker.templ](/e:/programmering/Workspace/Web%20devlopment/conorganizer/components/ticket_holder/ticket_holder_picker.templ#L60) so it only passes `associatedTicketHolders`.

- In [ticket_holder_button.templ](/e:/programmering/Workspace/Web%20devlopment/conorganizer/components/ticket_holder/ticket_holder_button.templ#L29) lines 29-33, normalize the embedded payload to the same serialized shape used everywhere else:

```templ
<button
	class="btn btn--secondary no-marking billett-holder-picker-button"
	data-current-billett-holder={ fmt.Sprintf(
		"{\"Id\": %d, \"Name\": \"%s\", \"Email\": \"%s\"}",
		currentBillettHolder.Id,
		currentBillettHolder.Name,
		currentBillettHolder.Email,
	) }
	data-class:selected={ fmt.Sprintf("$billettHolderId == %d", currentBillettHolder.Id) }
	data-on:click={ fmt.Sprintf("$billettHolderId = %d", currentBillettHolder.Id) }
>
```

- In [ticket_holder_dropdown.js](/e:/programmering/Workspace/Web%20devlopment/conorganizer/static/web_components/ticket_holder_dropdown.js#L38), remove dropdown-owned persistence and make the component controlled by the parent:
  - add `selected-id` as an observed attribute
  - remove `#localStorageKey`
  - remove `saveSelectedToLocalStorage`
  - remove `hydrateSelection`
  - add `syncSelectedFromAttribute`
  - emit a selected payload object with `Id`, `Name`, and `Email`
  - keep `Color` only in `data-billettholdere` parsing and internal rendering, since the UI still uses it

```javascript
const DATA_BILLETTHOLDERE_ATTR = "data-billettholdere";
const SELECTED_ID_ATTR = "selected-id";

static get observedAttributes() {
	return [DATA_BILLETTHOLDERE_ATTR, SELECTED_ID_ATTR];
}

syncFromAttribute() {
	this.#billettholdere = this.parseBillettholdere();
	if (this.#billettholdere.length === 0) {
		this.teardownInteractiveElements();
		this.shadowRoot?.replaceChildren();
		return;
	}

	this.#arrowIconTemplateEle ??= this.querySelector("template[data-arrow-icon]");

	this.render();
	if (!this.setupInteractiveElements()) {
		return;
	}
	this.syncSelectedFromAttribute();
}

syncSelectedFromAttribute() {
	const selectedId = Number(this.getAttribute(SELECTED_ID_ATTR) ?? "0");
	const optionEles = this.getOptionElements();
	const selectedOptionEle = optionEles.find(
		(optionEle) => Number(optionEle.dataset.Id ?? "0") === selectedId,
	);
	if (selectedOptionEle) {
		this.renderSelected(selectedOptionEle);
	}
}

toSelectedPayload(optionEle) {
	return {
		Id: Number(optionEle.dataset.Id ?? "0"),
		Name: optionEle.dataset.Name ?? "",
		Email: optionEle.dataset.Email ?? "",
	};
}

handleOptionSelect(optionEle) {
	this.renderSelected(optionEle);
	this.emitBillettholderSelected(this.toSelectedPayload(optionEle));
}

emitBillettholderSelected(billettholder) {
	setTimeout(() => {
		this.dispatchEvent(
			new CustomEvent("billettholder-selected", {
				detail: billettholder,
				bubbles: true,
				composed: true,
			}),
		);
	}, 150);
}
```

- Also remove the inline `liEle.onclick = ...` path in the same JS file so there is only one selection path through the dropdown’s click/keyboard handlers.

- In [event_interests.templ](/e:/programmering/Workspace/Web%20devlopment/conorganizer/components/event_components/event_interests.templ#L144) lines 144-157, remove the duplicate `LSKey` script. Keep only the dialog click handler if still needed.

```templ
<script>
	const dialogElement = document.querySelector(`.interest-dialog`);

	dialogElement.addEventListener("click", (ev) => {
		ev.stopPropagation();
	});
</script>
```

## Public Interface Changes
- `<billettholder-dropdown>` gains a new input attribute: `selected-id`.
- `billettholder-selected` changes from `detail: number` to:

```javascript
detail: {
	Id,
	Name,
	Email,
}
```

- `data-current-billett-holder` uses the same canonical serialized shape:

```javascript
{
	Id,
	Name,
	Email,
}
```

- `data-billettholdere` remains unchanged and may still include `Color`, because the dropdown needs it for styling.

## Test Plan
- No stored billettholder:
  - parent resolves a valid default
  - parent writes one canonical payload to `localStorage`
  - `$billettHolderId` matches the resolved billettholder
  - the active button or dropdown selection reflects that value
- Valid stored billettholder:
  - parent keeps it selected
  - dropdown highlights it through `selected-id`
  - no child overwrites it on init
- Invalid or stale stored billettholder:
  - parent falls back to `yourBillettHolder`
  - if `yourBillettHolder` is not present, parent falls back to the first associated billettholder
  - storage is rewritten with a valid canonical payload
- Button mode (`2-3` billettholdere):
  - click updates `$billettHolderId`
  - parent writes only `Id`, `Name`, and `Email`
- Dropdown mode (`>3` billettholdere):
  - click and keyboard selection emit the canonical selected payload
  - parent writes only `Id`, `Name`, and `Email`
  - visual selection tracks `selected-id`
- Regression checks:
  - no remaining `localStorage` ownership in `ticket_holder_dropdown.js`
  - no duplicate `LSKey` definition outside `TicketHolderPicker`
  - no persistence code writes `Color` into serialized selected-billettholder payloads

## Assumptions
- `TicketHolderPicker` is the correct ownership boundary because it is the shared parent of both selector variants and already owns `$billettHolderId`.
- Backward-compatible reads only require `Id` from older stored JSON; missing `Name` or `Email` can be normalized when rewriting storage.
- The preferred fallback order is valid stored selection, then valid `yourBillettHolder`, then first available billettholder.
