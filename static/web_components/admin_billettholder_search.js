const HIGHLIGHT_ESCAPE_PATTERN = /[.*+?^${}()|[\]\\]/g
const COMPONENT_ATTRIBUTES = Object.freeze({
    billettholdereJson: "data-billettholdere",
    clearInputVersion: "data-clear-input",
})
const ADMIN_BILLETTHOLDER_SEARCH_TAG = "admin-billettholder-search"
const GLOBAL_STYLE_URLS = [
    "/static/index.css",
    "/static/buttons.css",
]

/**
 * Escape RegExp meta characters for a safe literal match.
 * @param {string} value
 * @returns {string}
 */
const escapeHighlightPart = (value) => value.replace(HIGHLIGHT_ESCAPE_PATTERN, "\\$&")

/**
 * Normalize a user query for matching.
 * @param {string} value
 * @returns {string}
 */
const normalize = (value) => value.trim().toLowerCase()

/**
 * Score how well a name matches a query.
 * @param {string} name
 * @param {string} query
 * @returns {number}
 */
const matchScore = (name, query) => {
    if (!query) return 0
    if (name === query) return 3
    if (name.startsWith(query)) return 2
    if (name.includes(query)) return 1
    return 0
}

/**
 * Merge overlapping [start, end] ranges.
 * @param {Array<[number, number]>} ranges
 * @returns {Array<[number, number]>}
 */
const mergeRanges = (ranges) => {
    if (ranges.length === 0) return []

    ranges.sort((a, b) => a[0] - b[0])
    const merged = [ranges[0]]
    for (let i = 1; i < ranges.length; i += 1) {
        const current = ranges[i]
        const last = merged[merged.length - 1]
        if (current[0] <= last[1]) {
            last[1] = Math.max(last[1], current[1])
        } else {
            merged.push(current)
        }
    }
    return merged
}

/**
 * Find and merge all match ranges in a label.
 * @param {string} label
 * @param {string} query
 * @returns {Array<[number, number]>}
 */
const collectMatchRanges = (label, query) => {
    if (!query) return []
    const parts = query.split(/\s+/).filter(Boolean).map(escapeHighlightPart)
    /** @type {Array<[number, number]>} */
    const ranges = []

    for (const part of parts) {
        const partMatch = new RegExp(part, "ig")
        let match = partMatch.exec(label)
        while (match) {
            ranges.push([match.index, match.index + match[0].length])
            match = partMatch.exec(label)
        }
    }

    return mergeRanges(ranges)
}

/**
 * Build a DOM fragment in memory, then insert once into the live DOM.
 * @param {string} label
 * @param {string} query
 * @returns {DocumentFragment}
 */
const renderHighlightFragment = (label, query) => {
    const fragment = document.createDocumentFragment()
    const ranges = collectMatchRanges(label, query)

    if (ranges.length === 0) {
        fragment.append(document.createTextNode(label))
        return fragment
    }

    let cursor = 0
    for (const [start, end] of ranges) {
        if (cursor < start) {
            fragment.append(document.createTextNode(label.slice(cursor, start)))
        }
        const mark = document.createElement("mark")
        mark.append(document.createTextNode(label.slice(start, end)))
        fragment.append(mark)
        cursor = end
    }
    if (cursor < label.length) {
        fragment.append(document.createTextNode(label.slice(cursor)))
    }
    return fragment
}

class AdminBillettholderSearch extends HTMLElement {
    static observedAttributes = [COMPONENT_ATTRIBUTES.billettholdereJson, COMPONENT_ATTRIBUTES.clearInputVersion]

    /** @type {Array<{id:number, label:string, norm:string}>} */
    #searchableBillettholderOptions = []
    /** @type {Array<{Id:number, FirstName:string, LastName:string}>} */
    #availableBillettholdere = []
    /** @type {HTMLInputElement | null} */
    #searchInputElement = null
    /** @type {HTMLDivElement | null} */
    #searchResultsElement = null
    /** @type {HTMLDivElement} */
    #shadowContentRoot
    /** @type {AbortController | null} */
    #eventListenerController = null
    #hasConnected = false
    #clearInputVersion = ""

    constructor() {
        super()
        this.#clearInputVersion = this.getAttribute(COMPONENT_ATTRIBUTES.clearInputVersion) ?? ""
        this.#shadowContentRoot = document.createElement("div")
        this.#shadowContentRoot.className = "admin-billettholder-search-root"

        if (!this.shadowRoot) {
            const shadowRoot = this.attachShadow({ mode: "open" })
            for (const url of GLOBAL_STYLE_URLS) {
                const link = document.createElement("link")
                link.rel = "stylesheet"
                link.href = url
                shadowRoot.appendChild(link)
            }
            const style = document.createElement("style")
            style.textContent = `
                :host {
                    display: block;
                }

                .admin-billettholder-search-root {
                    display: block;

                    input {
                        font-family: var(--font-monospace);
                    }

                    .input {
                        background-color: var(--bg-item);
                        color: var(--color-text-primary);
                        border-radius: var(--border-radius-2x);
                        min-height: 2.6rem;
                        border: 1px solid var(--bg-item-border);
                        font-size: 1rem;
                        padding-inline: 1rem;
                        margin: 0;
                        transition: border-color 0.15s ease-in-out, box-shadow 0.15s ease-in-out;
                        box-sizing: border-box;

                        &::placeholder {
                            color: var(--color-text-soft-50);
                        }

                        &:focus-visible {
                            outline: 0;
                            border: 1px solid var(--color-primary-hover);
                            box-shadow: 0 0 0 0.25rem hsla(from var(--color-primary-hover) h s l / 0.25);
                        }
                    }

                    .gm-search-results {
                        margin-top: var(--spacing-4x);
                        display: block;
                    }

                    .gm-search-empty {
                        color: var(--color-text-soft);
                    }

                    .gm-search-item {
                        display: inline-flex;
                        align-items: center;
                        justify-content: flex-start;
                        height: var(--btn-height);
                        padding: 0 var(--btn-padding-x);
                        font-size: var(--btn-font-size);
                        font-weight: var(--btn-font-weight);
                        line-height: 1;
                        border-radius: var(--btn-border-radius);
                        border-style: solid;
                        border-width: var(--btn-border-width);
                        cursor: pointer;
                        user-select: none;
                        transition:
                            background-color var(--btn-transition-duration) ease,
                            color var(--btn-transition-duration) ease,
                            border-color var(--btn-transition-duration) ease,
                            box-shadow var(--btn-transition-duration) ease,
                            transform var(--btn-transition-duration) ease;
                        inline-size: auto;
                        max-inline-size: 100%;
                        text-align: left;
                        border-color: var(--btn-outline-border);
                        color: var(--color-secondary);
                        background-color: transparent;
                        margin: 0 var(--spacing-2x) var(--spacing-2x) 0;
                        vertical-align: top;

                        &:hover {
                            background-color: var(--btn-outline-hover-bg);
                            color: var(--btn-outline-active-text);
                        }

                        &:focus-visible {
                            outline: none;
                            background-color: var(--btn-outline-hover-bg);
                            color: var(--btn-outline-active-text);
                            box-shadow: 0 0 0 3px var(--btn-outline-focus-shadow);
                        }
                    }

                    mark {
                        background: var(--color-warning);
                        color: inherit;
                    }
                }
            `
            shadowRoot.append(style, this.#shadowContentRoot)
        }

        this.handleInput = this.handleInput.bind(this)
        this.handleClick = this.handleClick.bind(this)
        this.handleInputKeydown = this.handleInputKeydown.bind(this)
    }

    connectedCallback() {
        if (this.#hasConnected) return
        this.#hasConnected = true

        this.#syncBillettholdereFromAttribute()
        this.#renderSearchInterface()
        this.#rebuildSearchableOptions()
        this.#bindEventListeners()
        this.#renderSearchResults(this.#searchInputElement?.value ?? "")
    }

    disconnectedCallback() {
        this.#eventListenerController?.abort()
        this.#eventListenerController = null
    }

    /**
     * Update billettholder data via property assignment.
     * @param {Array<{Id:number, FirstName:string, LastName:string}>} value
     */
    set billettholdere(value) {
        this.#availableBillettholdere = Array.isArray(value) ? value : []
        this.#rebuildSearchableOptions()
        this.#renderSearchResults(this.#searchInputElement?.value ?? "")
    }

    /**
     * Current billettholder list.
     * @returns {Array<{Id:number, FirstName:string, LastName:string}>}
     */
    get billettholdere() {
        return this.#availableBillettholdere
    }

    /**
     * @param {`data-billettholdere`|`data-clear-input`} name
     * @param {string|null} oldValue
     * @param {string|null} newValue
     */
    attributeChangedCallback(name, oldValue, newValue) {
        if (oldValue === newValue) return
        if (name === COMPONENT_ATTRIBUTES.billettholdereJson) {
            this.#syncBillettholdereFromAttribute()
            this.#rebuildSearchableOptions()
            this.#renderSearchResults(this.#searchInputElement?.value ?? "")
            return
        }
        if (name === COMPONENT_ATTRIBUTES.clearInputVersion) {
            this.#clearInputVersion = newValue ?? ""
            this.#clearSearchInputAndResults()
        }
    }

    #syncBillettholdereFromAttribute() {
        const raw = this.getAttribute(COMPONENT_ATTRIBUTES.billettholdereJson)
        if (!raw) return
        try {
            const data = JSON.parse(raw)
            this.#availableBillettholdere = Array.isArray(data) ? data : []
        } catch (error) {
            console.warn("billettholder-search: invalid JSON data", error)
        }
    }

    #rebuildSearchableOptions() {
        this.#searchableBillettholderOptions = this.#availableBillettholdere.map((billettholder) => {
            const label = `${ billettholder.FirstName } ${ billettholder.LastName }`
            return {
                id: billettholder.Id,
                label,
                norm: label.toLowerCase(),
            }
        })
    }

    #renderSearchInterface() {
        const placeholder = this.getAttribute("placeholder") ?? "s›k etter spiller"
        const inputId = this.getAttribute("input-id") ?? `gm-search-${ Math.random().toString(36).substring(2, 8) }`
        const inputTippy = this.getAttribute("input-tippy") ?? ""

        this.#shadowContentRoot.replaceChildren()

        const input = document.createElement("input")
        input.id = inputId
        input.type = "search"
        input.autocomplete = "off"
        input.placeholder = placeholder
        input.className = "input"
        input.required = true
        input.title = ""
        input.setAttribute("data-tippy-content", inputTippy)

        const results = document.createElement("div")
        results.className = "gm-search-results"
        results.setAttribute("aria-live", "polite")

        this.#shadowContentRoot.append(input, results)
        this.#searchInputElement = input
        this.#searchResultsElement = results
    }

    #bindEventListeners() {
        if (!this.#searchInputElement || !this.#searchResultsElement) return

        this.#eventListenerController?.abort()
        this.#eventListenerController = new AbortController()
        const signal = this.#eventListenerController.signal

        this.#searchInputElement.addEventListener("input", this.handleInput, { signal })
        this.#searchInputElement.addEventListener("keydown", this.handleInputKeydown, { signal })
        this.#searchResultsElement.addEventListener("click", this.handleClick, { signal })
    }

    /**
     * @param {string} query
     */
    #renderSearchResults(query) {
        const normalizedQuery = normalize(query)

        this.#searchResultsElement?.replaceChildren()
        if (!normalizedQuery) return

        const matchingBillettholderOptions = this.#searchableBillettholderOptions
            .map((option) => ({ ...option, score: matchScore(option.norm, normalizedQuery) }))
            .filter((option) => option.score > 0)
            .sort((leftOption, rightOption) => rightOption.score - leftOption.score || leftOption.label.localeCompare(rightOption.label))
            .slice(0, 8)

        if (matchingBillettholderOptions.length === 0) {
            const emptyStateElement = document.createElement("div")
            emptyStateElement.classList.add("gm-search-empty")
            emptyStateElement.append(document.createTextNode("Ingen billettholdere funnet"))
            this.#searchResultsElement?.append(emptyStateElement)
            return
        }

        const resultButtonsFragment = document.createDocumentFragment()
        for (const option of matchingBillettholderOptions) {
            const resultButton = document.createElement("button")
            resultButton.type = "button"
            resultButton.classList.add("btn", "btn--outline", "gm-search-item")
            resultButton.dataset.value = option.label
            resultButton.dataset.id = String(option.id)
            resultButton.append(renderHighlightFragment(option.label, normalizedQuery))
            resultButtonsFragment.append(resultButton)
        }

        this.#searchResultsElement?.append(resultButtonsFragment)
    }

    handleInput() {
        this.#renderSearchResults(this.#searchInputElement?.value ?? "")
    }

    /**
     * @param {KeyboardEvent} event
     */
    handleInputKeydown(event) {
        if (event.key !== "Enter") return
        const firstVisibleResultButton = this.#searchResultsElement?.querySelector(".gm-search-item")
        if (!(firstVisibleResultButton instanceof HTMLButtonElement)) return
        event.preventDefault()
        this.#selectSearchResultButton(firstVisibleResultButton)
    }

    /**
     * @param {MouseEvent} event
     */
    handleClick(event) {
        const target = event.target
        if (!(target instanceof HTMLElement)) return
        const resultButton = target.closest(".gm-search-item")
        if (!(resultButton instanceof HTMLButtonElement)) return
        this.#selectSearchResultButton(resultButton)
    }

    /**
     * @param {HTMLButtonElement} button
     */
    #selectSearchResultButton(selectedResultButton) {
        const selectedLabel = selectedResultButton.getAttribute("data-value")
        if (!selectedLabel) return
        const selectedId = selectedResultButton.getAttribute("data-id")
        if (!selectedId) return

        if (this.#searchInputElement) this.#searchInputElement.value = selectedLabel
        this.#searchResultsElement?.replaceChildren()

        this.dispatchEvent(
            new CustomEvent("billettholder-select", {
                detail: {
                    id: Number(selectedId),
                    label: selectedLabel,
                },
            }),
        )
    }

    #clearSearchInputAndResults() {
        if (this.#searchInputElement) this.#searchInputElement.value = ""
        this.#searchResultsElement?.replaceChildren()
    }
}

if (!customElements.get(ADMIN_BILLETTHOLDER_SEARCH_TAG)) {
    customElements.define(ADMIN_BILLETTHOLDER_SEARCH_TAG, AdminBillettholderSearch)
}
