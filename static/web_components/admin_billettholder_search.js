const HIGHLIGHT_ESCAPE_PATTERN = /[.*+?^${}()|[\]\\]/g
const ATTRS = Object.freeze({
    billettholdere: "data-billettholdere",
    clearInput: "data-clear-input",
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
    static observedAttributes = [ATTRS.billettholdere, ATTRS.clearInput]

    /** @type {Array<{id:number, label:string, norm:string}>} */
    #searchOptions = []
    /** @type {Array<{Id:number, FirstName:string, LastName:string}>} */
    #billettholdere = []
    /** @type {HTMLInputElement | null} */
    #inputEl = null
    /** @type {HTMLDivElement | null} */
    #searchResultsEl = null
    /** @type {HTMLDivElement} */
    #mountEl
    /** @type {AbortController | null} */
    #events = null
    #initialized = false
    #clearInput = ""

    constructor() {
        super()
        this.#clearInput = this.getAttribute(ATTRS.clearInput) ?? ""
        this.#mountEl = document.createElement("div")
        this.#mountEl.className = "admin-billettholder-search-root"

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
            shadowRoot.append(style, this.#mountEl)
        }

        this.handleInput = this.handleInput.bind(this)
        this.handleClick = this.handleClick.bind(this)
        this.handleInputKeydown = this.handleInputKeydown.bind(this)
    }

    connectedCallback() {
        if (this.#initialized) return
        this.#initialized = true

        this.#loadDataFromAttribute()
        this.#render()
        this.#setOptions()
        this.#bind()
        this.#renderMatches(this.#inputEl?.value ?? "")
    }

    disconnectedCallback() {
        this.#events?.abort()
        this.#events = null
    }

    /**
     * Update billettholder data via property assignment.
     * @param {Array<{Id:number, FirstName:string, LastName:string}>} value
     */
    set billettholdere(value) {
        this.#billettholdere = Array.isArray(value) ? value : []
        this.#setOptions()
        this.#renderMatches(this.#inputEl?.value ?? "")
    }

    /**
     * Current billettholder list.
     * @returns {Array<{Id:number, FirstName:string, LastName:string}>}
     */
    get billettholdere() {
        return this.#billettholdere
    }

    /**
     * @param {`data-billettholdere`|`data-clear-input`} name
     * @param {string|null} oldValue
     * @param {string|null} newValue
     */
    attributeChangedCallback(name, oldValue, newValue) {
        if (oldValue === newValue) return
        if (name === ATTRS.billettholdere) {
            this.#loadDataFromAttribute()
            this.#setOptions()
            this.#renderMatches(this.#inputEl?.value ?? "")
            return
        }
        if (name === ATTRS.clearInput) {
            this.#clearInput = newValue ?? ""
            this.#clearSearch()
        }
    }

    #loadDataFromAttribute() {
        const raw = this.getAttribute(ATTRS.billettholdere)
        if (!raw) return
        try {
            const data = JSON.parse(raw)
            this.#billettholdere = Array.isArray(data) ? data : []
        } catch (error) {
            console.warn("billettholder-search: invalid JSON data", error)
        }
    }

    #setOptions() {
        this.#searchOptions = this.#billettholdere.map((billettholder) => {
            const label = `${ billettholder.FirstName } ${ billettholder.LastName }`
            return {
                id: billettholder.Id,
                label,
                norm: label.toLowerCase(),
            }
        })
    }

    #render() {
        const placeholder = this.getAttribute("placeholder") ?? "s›k etter spiller"
        const inputId = this.getAttribute("input-id") ?? `gm-search-${ Math.random().toString(36).substring(2, 8) }`
        const inputTippy = this.getAttribute("input-tippy") ?? ""

        this.#mountEl.replaceChildren()

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

        this.#mountEl.append(input, results)
        this.#inputEl = input
        this.#searchResultsEl = results
    }

    #bind() {
        if (!this.#inputEl || !this.#searchResultsEl) return

        this.#events?.abort()
        this.#events = new AbortController()
        const signal = this.#events.signal

        this.#inputEl.addEventListener("input", this.handleInput, { signal })
        this.#inputEl.addEventListener("keydown", this.handleInputKeydown, { signal })
        this.#searchResultsEl.addEventListener("click", this.handleClick, { signal })
    }

    /**
     * @param {string} query
     */
    #renderMatches(query) {
        const norm = normalize(query)

        this.#searchResultsEl?.replaceChildren()
        if (!norm) return

        const matches = this.#searchOptions
            .map((opt) => ({ ...opt, score: matchScore(opt.norm, norm) }))
            .filter((opt) => opt.score > 0)
            .sort((a, b) => b.score - a.score || a.label.localeCompare(b.label))
            .slice(0, 8)

        if (matches.length === 0) {
            const empty = document.createElement("div")
            empty.classList.add("gm-search-empty")
            empty.append(document.createTextNode("Ingen billettholdere funnet"))
            this.#searchResultsEl?.append(empty)
            return
        }

        const fragment = document.createDocumentFragment()
        for (const opt of matches) {
            const button = document.createElement("button")
            button.type = "button"
            button.classList.add("btn", "btn--outline", "gm-search-item")
            button.dataset.value = opt.label
            button.dataset.id = String(opt.id)
            button.append(renderHighlightFragment(opt.label, norm))
            fragment.append(button)
        }

        this.#searchResultsEl?.append(fragment)
    }

    handleInput() {
        this.#renderMatches(this.#inputEl?.value ?? "")
    }

    /**
     * @param {KeyboardEvent} event
     */
    handleInputKeydown(event) {
        if (event.key !== "Enter") return
        const firstResult = this.#searchResultsEl?.querySelector(".gm-search-item")
        if (!(firstResult instanceof HTMLButtonElement)) return
        event.preventDefault()
        this.#selectMatchButton(firstResult)
    }

    /**
     * @param {MouseEvent} event
     */
    handleClick(event) {
        const target = event.target
        if (!(target instanceof HTMLElement)) return
        const button = target.closest(".gm-search-item")
        if (!(button instanceof HTMLButtonElement)) return
        this.#selectMatchButton(button)
    }

    /**
     * @param {HTMLButtonElement} button
     */
    #selectMatchButton(button) {
        const value = button.getAttribute("data-value")
        if (!value) return
        const id = button.getAttribute("data-id")
        if (!id) return

        if (this.#inputEl) this.#inputEl.value = value
        this.#searchResultsEl?.replaceChildren()

        this.dispatchEvent(
            new CustomEvent("billettholder-select", {
                detail: {
                    id: Number(id),
                    label: value,
                },
            }),
        )
    }

    #clearSearch() {
        if (this.#inputEl) this.#inputEl.value = ""
        this.#searchResultsEl?.replaceChildren()
    }
}

if (!customElements.get(ADMIN_BILLETTHOLDER_SEARCH_TAG)) {
    customElements.define(ADMIN_BILLETTHOLDER_SEARCH_TAG, AdminBillettholderSearch)
}
