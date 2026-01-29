const HIGHLIGHT_ESCAPE_PATTERN = /[.*+?^${}()|[\]\\]/g

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
const scoreMatch = (name, query) => {
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
            /** @type {[number, number]} */
            const range = [match.index, match.index + match[0].length]
            ranges.push(range)
            match = partMatch.exec(label)
        }
    }

    return mergeRanges(ranges)
}

/**
 * Build a DOM fragment in memory, then insert once into the live DOM.
 * This avoids repeated reflows while adding <mark> nodes.
 * @param {string} label
 * @param {string} query
 * @returns {DocumentFragment}
 */
const buildHighlight = (label, query) => {
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

/**
 * Admin GM select web component.
 *
 * Events:
 * - name: "gm-select"
 *   detail: {
 *     id: number,
 *     label: string
 *   }
 *
 *   Use in Datastar:
 *   data-on:gm-select="$gmSearchBillettholderId = evt.detail.id"
 *
 * Attributes:
 * - data-billettholdere: JSON array of { Id, FirstName, LastName }
 */
class AdminGmSelect extends HTMLElement {
    /**
     * Datastar reads updates via attributes; observe changes to re-render matches.
     * @returns {string[]}
     */
    static get observedAttributes() {
        return ["data-billettholdere"]
    }

    constructor() {
        super()

        /** @type {Array<{id:number, label:string, norm:string}>} */
        this.options = []

        /** @type {Array<{Id:number, FirstName:string, LastName:string}>} */
        this._billettholdere = []
        this._initialized = false

        this.handleInput = this.handleInput.bind(this)
        this.handleClick = this.handleClick.bind(this)
    }

    /**
     * Build UI and bind listeners when inserted into the DOM.
     */
    connectedCallback() {
        if (this._initialized) return
        this._initialized = true

        this._loadDataFromAttribute()
        this._render()
        this._setOptions()
        this._bind()

        if (this.input) {
            this._renderMatches(this.input.value)
        }
    }

    /**
     * Cleanup listeners when removed from the DOM.
     */
    disconnectedCallback() {
        if (!this.input || !this.results) return
        this.input.removeEventListener("input", this.handleInput)
        this.results.removeEventListener("click", this.handleClick)
    }

    /**
     * Update billettholder data via property assignment.
     * @param {Array<{Id:number, FirstName:string, LastName:string}>} value
     */
    set billettholdere(value) {
        this._billettholdere = Array.isArray(value) ? value : []
        this._setOptions()
        if (this.input && this.results) {
            this._renderMatches(this.input.value)
        }
    }

    /**
     * Current billettholder list.
     * @returns {Array<{Id:number, FirstName:string, LastName:string}>}
     */
    get billettholdere() {
        return this._billettholdere
    }

    /**
     * React to Datastar-driven attribute updates.
     * @param {string} name
     * @param {string|null} oldValue
     * @param {string|null} newValue
     */

    attributeChangedCallback(name, oldValue, newValue) {
        if (oldValue === newValue) return
        if (name !== "data-billettholdere") return
        this._loadDataFromAttribute()
        this._setOptions()
        if (this.input && this.results) {
            this._renderMatches(this.input.value)
        }
    }

    /**
     * Parse billettholder data from data-billettholdere.
     */
    _loadDataFromAttribute() {
        const raw = this.getAttribute("data-billettholdere")
        if (!raw) return
        try {
            const data = JSON.parse(raw)
            this._billettholdere = Array.isArray(data) ? data : []
        } catch (error) {
            console.warn("gm-picker: invalid JSON data", error)
        }
    }

    /**
     * Build searchable options from billettholder data.
     */
    _setOptions() {
        this.options = (this._billettholdere || []).map((billettholder) => {
            const label = `${ billettholder.FirstName } ${ billettholder.LastName }`
            return {
                id: billettholder.Id,
                label,
                norm: label.toLowerCase(),
            }
        })
    }

    /**
     * Render the light DOM structure so page styles apply.
     */
    _render() {
        // Light DOM so global styles from the page keep working.
        const labelText = this.getAttribute("label") || "Søk etter spiller som skal være spilleder"
        const placeholder = this.getAttribute("placeholder") || "søk etter spiller"
        const submitLabel = this.getAttribute("submit-label") || "Lagre"
        const inputId = this.getAttribute("input-id") || `gm-search-${ Math.random().toString(36).slice(2, 8) }`

        this.replaceChildren()

        const label = document.createElement("label")
        label.setAttribute("for", inputId)
        label.append(document.createTextNode(labelText))

        const input = document.createElement("input")
        input.id = inputId
        input.placeholder = placeholder
        input.className = "input"
        input.required = true

        const button = document.createElement("button")
        button.type = "submit"
        button.className = "btn btn--primary"
        button.append(document.createTextNode(submitLabel))

        const results = document.createElement("div")
        results.className = "gm-search-results"
        results.setAttribute("aria-live", "polite")

        this.append(label, input, button, results)

        /** @type {HTMLInputElement} */
        this.input = input
        /** @type {HTMLDivElement} */
        this.results = results
    }

    /**
     * Wire input + click handlers.
     */
    _bind() {

        this.input.addEventListener("input", this.handleInput)

        this.results.addEventListener("click", this.handleClick)
    }

    /**
     * Render search result buttons for the current query.
     * @param {string} query
     */
    _renderMatches(query) {
        const norm = normalize(query || "")

        this.results.replaceChildren()
        if (!norm) return

        const matches = this.options
            .map((opt) => ({ ...opt, score: scoreMatch(opt.norm, norm) }))
            .filter((opt) => opt.score > 0)
            .sort((a, b) => b.score - a.score || a.label.localeCompare(b.label))
            .slice(0, 8)

        if (matches.length === 0) {
            const empty = document.createElement("div")
            empty.classList.add("gm-search-empty")
            empty.append(document.createTextNode("Ingen billettholdere funnet"))

            this.results.append(empty)
            return
        }

        const fragment = document.createDocumentFragment()
        for (const opt of matches) {
            const button = document.createElement("button")
            button.type = "button"
            button.classList.add("btn", "btn--outline", "gm-search-item")
            button.dataset.value = opt.label
            button.dataset.id = String(opt.id)
            button.append(buildHighlight(opt.label, norm))
            fragment.append(button)
        }

        this.results.append(fragment)
    }

    /**
     * Handle input changes.
     */
    handleInput() {

        this._renderMatches(this.input.value)
    }

    /**
     * Handle result selection and dispatch gm-select.
     * @param {MouseEvent} event
     */

    handleClick(event) {
        const target = event.target
        if (!(target instanceof HTMLElement)) return
        const button = target.closest(".gm-search-item")
        if (!button) return
        const value = button.getAttribute("data-value")
        if (!value) return

        this.input.value = value

        this.results.replaceChildren()
        const id = button.getAttribute("data-id")
        if (id) {
            // Datastar listens to this event and updates signals via data-on:gm-select.
            this.dispatchEvent(
                new CustomEvent("gm-select", {
                    detail: {
                        id: Number(id),
                        label: value,
                    },
                })
            )
        }
    }

}

customElements.define("admin-gm-select", AdminGmSelect)
