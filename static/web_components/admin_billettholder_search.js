if (!customElements.get("admin-billettholder-search")) {
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
    const EMPTY_SEARCH_RESULTS_TEXT = "Ingen billettholdere funnet"

    const globalStyleSheetsPromise = (async () => {
        const supportsConstructableStyleSheets =
            !!(Document.prototype && "adoptedStyleSheets" in Document.prototype) &&
            !!(CSSStyleSheet.prototype && "replace" in CSSStyleSheet.prototype)

        if (!supportsConstructableStyleSheets) return null

        const styleSheets = []
        for (const url of GLOBAL_STYLE_URLS) {
            const response = await fetch(url, { credentials: "same-origin" })
            const cssText = await response.text()
            const styleSheet = new CSSStyleSheet()
            await styleSheet.replace(cssText)
            styleSheets.push(styleSheet)
        }
        return styleSheets
    })()

    const ADMIN_BILLETTHOLDER_SEARCH_SHADOW_STYLES = `
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
            color: var(--bg-base);
            margin-inline-end: 0.3ch;
        }
    }
`

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
    const normalizeQuery = (value) => value.trim().toLowerCase()

    /**
     * Score how well a name matches a query.
     * @param {string} candidateName
     * @param {string} normalizedQuery
     * @returns {number}
     */
    const calculateMatchScore = (candidateName, normalizedQuery) => {
        if (!normalizedQuery) return 0
        if (candidateName === normalizedQuery) return 3
        if (candidateName.startsWith(normalizedQuery)) return 2
        if (candidateName.includes(normalizedQuery)) return 1
        return 0
    }

    /**
     * Merge overlapping [start, end] ranges.
     * @param {Array<[number, number]>} ranges
     * @returns {Array<[number, number]>}
     */
    const mergeOverlappingRanges = (ranges) => {
        if (ranges.length === 0) return []

        ranges.sort((leftRange, rightRange) => leftRange[0] - rightRange[0])
        const mergedRanges = [ranges[0]]
        for (let i = 1; i < ranges.length; i += 1) {
            const currentRange = ranges[i]
            const previousMergedRange = mergedRanges[mergedRanges.length - 1]
            if (currentRange[0] <= previousMergedRange[1]) {
                previousMergedRange[1] = Math.max(previousMergedRange[1], currentRange[1])
            } else {
                mergedRanges.push(currentRange)
            }
        }
        return mergedRanges
    }

    /**
     * Find and merge all highlight ranges in a visible label.
     * @param {string} label
     * @param {string} normalizedQuery
     * @returns {Array<[number, number]>}
     */
    const collectHighlightRanges = (label, normalizedQuery) => {
        if (!normalizedQuery) return []
        const queryParts = normalizedQuery.split(/\s+/).filter(Boolean).map(escapeHighlightPart)
        /** @type {Array<[number, number]>} */
        const highlightRanges = []

        for (const queryPart of queryParts) {
            const partMatcher = new RegExp(queryPart, "ig")
            let partMatch = partMatcher.exec(label)
            while (partMatch) {
                highlightRanges.push([partMatch.index, partMatch.index + partMatch[0].length])
                partMatch = partMatcher.exec(label)
            }
        }

        return mergeOverlappingRanges(highlightRanges)
    }

    /**
     * Build a DOM fragment with highlighted <mark> ranges.
     * @param {string} label
     * @param {string} normalizedQuery
     * @returns {DocumentFragment}
     */
    const renderHighlightedLabelFragment = (label, normalizedQuery) => {
        const fragment = document.createDocumentFragment()
        const highlightRanges = collectHighlightRanges(label, normalizedQuery)

        if (highlightRanges.length === 0) {
            fragment.append(document.createTextNode(label))
            return fragment
        }

        let currentCursor = 0
        for (const [startIndex, endIndex] of highlightRanges) {
            if (currentCursor < startIndex) {
                fragment.append(document.createTextNode(label.slice(currentCursor, startIndex)))
            }
            const markElement = document.createElement("mark")
            markElement.append(document.createTextNode(label.slice(startIndex, endIndex)))
            fragment.append(markElement)
            currentCursor = endIndex
        }
        if (currentCursor < label.length) {
            fragment.append(document.createTextNode(label.slice(currentCursor)))
        }
        return fragment
    }

    /**
     * Find and rank the best visible billettholder matches for a normalized query.
     * @param {Array<{id:number, label:string, normalizedLabel:string}>} searchableBillettholderOptions
     * @param {string} normalizedQuery
     * @returns {Array<{id:number, label:string, normalizedLabel:string, score:number}>}
     */
    const getTopMatchingBillettholderOptions = (searchableBillettholderOptions, normalizedQuery) => {
        return searchableBillettholderOptions
            .map((option) => ({ ...option, score: calculateMatchScore(option.normalizedLabel, normalizedQuery) }))
            .filter((option) => option.score > 0)
            // @ts-ignore
            .toSorted((leftOption, rightOption) => rightOption.score - leftOption.score || leftOption.label.localeCompare(rightOption.label))
            .toSpliced(8)
    }

    /**
     * Searchable admin billettholder picker web component.
     *
     * Attributes:
     * - `data-billettholdere`: JSON array of `{ Id, FirstName, LastName }`
     * - `data-clear-input`: changing the value clears the current search input and results
     * - `placeholder`: optional input placeholder text
     * - `input-tippy`: optional tooltip text for the search input
     *
     * Events:
     * - `billettholder-select`: emitted when a result is selected
     *   detail: `{ id: number, label: string }`
     *
     * Datastar usage:
     * - `data-attr:clear-input="$clearInput"`
     * - `data-on:billettholder-select="$assignmentBillettholderId = evt.detail.id"`
     *
     * Renders its internal UI inside Shadow DOM so Datastar light-DOM patching
     * does not wipe the search input after the component has connected.
     *
     * @extends HTMLElement
     */
    class AdminBillettholderSearch extends HTMLElement {
        static observedAttributes = [COMPONENT_ATTRIBUTES.billettholdereJson, COMPONENT_ATTRIBUTES.clearInputVersion]

        /** @type {Array<{id:number, label:string, normalizedLabel:string}>} */
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
        /** @type {boolean} */
        #hasConnected = false
        /** @type {string} */
        #clearInputVersion = ""

        constructor() {
            super()
            this.#clearInputVersion = this.getAttribute(COMPONENT_ATTRIBUTES.clearInputVersion) ?? ""
            this.#shadowContentRoot = document.createElement("div")
            this.#shadowContentRoot.className = "admin-billettholder-search-root"

            if (!this.shadowRoot) {
                const shadowRoot = this.attachShadow({ mode: "open" })
                globalStyleSheetsPromise.then((styleSheets) => {
                    if (styleSheets && this.shadowRoot) {
                        this.shadowRoot.adoptedStyleSheets = [...this.shadowRoot.adoptedStyleSheets, ...styleSheets]
                    } else if (this.shadowRoot) {
                        for (const url of GLOBAL_STYLE_URLS) {
                            const linkElement = document.createElement("link")
                            linkElement.rel = "stylesheet"
                            linkElement.href = url
                            this.shadowRoot.appendChild(linkElement)
                        }
                    }
                })
                shadowRoot.append(this.#createShadowStyleElement(), this.#shadowContentRoot)
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
         * Public API for clearing the current search value and visible results.
         * @returns {void}
         */
        clearSearch() {
            this.#clearSearchInputAndResults()
        }

        /**
         * React to Datastar-driven attribute updates.
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

        /**
         * Create the component-local shadow DOM style element.
         * @returns {HTMLStyleElement}
         */
        #createShadowStyleElement() {
            const styleElement = document.createElement("style")
            styleElement.textContent = ADMIN_BILLETTHOLDER_SEARCH_SHADOW_STYLES
            return styleElement
        }

        /**
         * Read billettholder JSON from the host element attribute into component state.
         * @returns {void}
         */
        #syncBillettholdereFromAttribute() {
            const billettholdereJson = this.getAttribute(COMPONENT_ATTRIBUTES.billettholdereJson)
            if (!billettholdereJson) return
            try {
                const parsedBillettholdere = JSON.parse(billettholdereJson)
                this.#availableBillettholdere = Array.isArray(parsedBillettholdere) ? parsedBillettholdere : []
            } catch (error) {
                console.warn("billettholder-search: invalid JSON data", error)
            }
        }

        /**
         * Build the normalized option list used by the in-memory search.
         * @returns {void}
         */
        #rebuildSearchableOptions() {
            this.#searchableBillettholderOptions = this.#availableBillettholdere.map((billettholder) => {
                const label = `${ billettholder.FirstName } ${ billettholder.LastName }`
                return {
                    id: billettholder.Id,
                    label,
                    normalizedLabel: label.toLowerCase(),
                }
            })
        }

        /**
         * Render the search input and result container inside shadow DOM.
         * @returns {void}
         */
        #renderSearchInterface() {
            this.#shadowContentRoot.replaceChildren()
            const searchInputElement = this.#createSearchInputElement()
            const searchResultsContainer = this.#createSearchResultsContainer()
            this.#shadowContentRoot.append(searchInputElement, searchResultsContainer)
            this.#searchInputElement = searchInputElement
            this.#searchResultsElement = searchResultsContainer
        }

        /**
         * Create the search input element from current host attributes.
         * @returns {HTMLInputElement}
         */
        #createSearchInputElement() {
            const searchInputElement = document.createElement("input")
            searchInputElement.type = "search"
            searchInputElement.autocomplete = "off"
            searchInputElement.placeholder = this.getAttribute("placeholder") ?? "søk etter spiller"
            searchInputElement.className = "input"
            searchInputElement.required = true
            searchInputElement.title = ""
            searchInputElement.setAttribute("data-tippy-content", this.getAttribute("input-tippy") ?? "")
            return searchInputElement
        }

        /**
         * Create the search result container.
         * @returns {HTMLDivElement}
         */
        #createSearchResultsContainer() {
            const searchResultsContainer = document.createElement("div")
            searchResultsContainer.className = "gm-search-results"
            searchResultsContainer.setAttribute("aria-live", "polite")
            return searchResultsContainer
        }

        /**
         * Attach DOM listeners using an AbortController so rebinds stay cheap and safe.
         * @returns {void}
         */
        #bindEventListeners() {
            if (!this.#searchInputElement || !this.#searchResultsElement) return

            this.#eventListenerController?.abort()
            this.#eventListenerController = new AbortController()
            const { signal } = this.#eventListenerController

            this.#searchInputElement.addEventListener("input", this.handleInput, { signal })
            this.#searchInputElement.addEventListener("keydown", this.handleInputKeydown, { signal })
            this.#searchResultsElement.addEventListener("click", this.handleClick, { signal })
        }

        /**
         * Render the current search result list for a query string.
         * @param {string} query
         * @returns {void}
         */
        #renderSearchResults(query) {
            const normalizedQuery = normalizeQuery(query)

            this.#searchResultsElement?.replaceChildren()
            if (!normalizedQuery) return

            const matchingBillettholderOptions = getTopMatchingBillettholderOptions(
                this.#searchableBillettholderOptions,
                normalizedQuery,
            )

            if (matchingBillettholderOptions.length === 0) {
                const emptyStateElement = document.createElement("div")
                emptyStateElement.classList.add("gm-search-empty")
                emptyStateElement.append(document.createTextNode(EMPTY_SEARCH_RESULTS_TEXT))
                this.#searchResultsElement?.append(emptyStateElement)
                return
            }

            const resultButtonsFragment = document.createDocumentFragment()
            for (const option of matchingBillettholderOptions) {
                resultButtonsFragment.append(this.#createSearchResultButton(option, normalizedQuery))
            }

            this.#searchResultsElement?.append(resultButtonsFragment)
        }

        /**
         * Create a selectable result button for one billettholder option.
         * @param {{id:number, label:string, normalizedLabel:string}} option
         * @param {string} normalizedQuery
         * @returns {HTMLButtonElement}
         */
        #createSearchResultButton(option, normalizedQuery) {
            const resultButtonElement = document.createElement("button")
            resultButtonElement.type = "button"
            resultButtonElement.classList.add("btn", "btn--outline", "gm-search-item")
            resultButtonElement.dataset.value = option.label
            resultButtonElement.dataset.id = String(option.id)
            resultButtonElement.append(renderHighlightedLabelFragment(option.label, normalizedQuery))
            return resultButtonElement
        }

        /**
         * Handle text input updates from the search field.
         * @returns {void}
         */
        handleInput() {
            this.#renderSearchResults(this.#searchInputElement?.value ?? "")
        }

        /**
         * Pressing Enter selects the first visible search result.
         * @param {KeyboardEvent} event
         * @returns {void}
         */
        handleInputKeydown(event) {
            if (event.key !== "Enter") return
            const firstVisibleResultButton = this.#searchResultsElement?.querySelector(".gm-search-item")
            if (!(firstVisibleResultButton instanceof HTMLButtonElement)) return
            event.preventDefault()
            this.#selectSearchResultButton(firstVisibleResultButton)
        }

        /**
         * Handle click selection from the rendered search results.
         * @param {MouseEvent} event
         * @returns {void}
         */
        handleClick(event) {
            const eventTarget = event.target
            if (!(eventTarget instanceof HTMLElement)) return
            const resultButton = eventTarget.closest(".gm-search-item")
            if (!(resultButton instanceof HTMLButtonElement)) return
            this.#selectSearchResultButton(resultButton)
        }

        /**
         * Select a result button, write the label into the input, and emit the public selection event.
         * @param {HTMLButtonElement} selectedResultButton
         * @returns {void}
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
                    bubbles: true,
                    composed: true,
                }),
            )
            this.#restoreInputFocusAfterSelection()
        }

        /**
         * Clear the current search input value and any rendered search results.
         * @returns {void}
         */
        #clearSearchInputAndResults() {
            if (this.#searchInputElement) this.#searchInputElement.value = ""
            this.#searchResultsElement?.replaceChildren()
        }

        /**
         * Restore focus to the current search input after Datastar updates settle.
         * @returns {void}
         */
        #restoreInputFocusAfterSelection() {
            requestAnimationFrame(() => {
                const currentInputElement = this.shadowRoot?.querySelector(".input")
                if (currentInputElement instanceof HTMLInputElement) {
                    currentInputElement.focus()
                }
            })
        }

    }

    customElements.define(ADMIN_BILLETTHOLDER_SEARCH_TAG, AdminBillettholderSearch)
}
