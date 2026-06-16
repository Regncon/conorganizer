// @ts-check

// Shared browser helpers for templates and web components.
// Keep this file build-free and browser-native.
(() => {
    /**
     * @typedef {Object} StoredBillettholder
     * @property {number} Id
     * @property {string} Name
     * @property {string} Email
     * @property {string} [Color]
     */
    /**
     * Shared style loading API for custom elements.
     * @typedef {Object} SharedStyles
     * @property {(urls: string[]) => Promise<CSSStyleSheet[] | null>} getStyleSheets
     * @property {(extraUrls?: string[]) => string[]} getStyleUrls
     * @property {(shadowRoot: ShadowRoot, urls: string[]) => Promise<void>} applyStyleUrlsToShadowRoot
     */
    /**
     * Style value helpers for billettholder initials.
     * @typedef {Object} BillettholderSelectionStyle
     * @property {(color: string) => string} backgroundColor
     * @property {(color: string) => string} border
     */
    /**
     * Shared selected-billettholder API used by templates and web components.
     * @typedef {Object} BillettholderSelection
     * @property {() => void} clear
     * @property {(name: string) => string} colorFromName
     * @property {(associatedBillettholdere: unknown[], billettholderId: unknown) => StoredBillettholder | null} findAssociated
     * @property {() => StoredBillettholder | null} get
     * @property {(name: string) => string} getInitials
     * @property {(associatedBillettholdere: unknown[], yourBillettholder: unknown) => StoredBillettholder | null} initialize
     * @property {(callback: (billettholder: StoredBillettholder | null) => void) => () => void} onChange
     * @property {(billettholder: unknown) => StoredBillettholder | null} set
     * @property {BillettholderSelectionStyle} style
     * @property {string} storageKey
     */
    /**
     * @typedef {Object} ConorganizerGlobal
     * @property {BillettholderSelection} billettholderSelection
     * @property {SharedStyles} sharedStyles
     */
    /**
     * @typedef {Window & typeof globalThis & {
     *   conorganizer?: Partial<ConorganizerGlobal>,
     * }} ConorganizerWindow
     */

    /** @type {ConorganizerWindow} */
    const typedWindow = window
    const existingConorganizer = typedWindow.conorganizer ?? {}

    const sharedStyles = existingConorganizer.sharedStyles ?? createSharedStyles()
    const billettholderSelection = existingConorganizer.billettholderSelection ?? createBillettholderSelection()

    /** @type {ConorganizerGlobal} */
    const conorganizer = Object.freeze({
        ...existingConorganizer,
        billettholderSelection,
        sharedStyles,
    })

    typedWindow.conorganizer = conorganizer

    function createSharedStyles() {
        const BASE_STYLE_URLS = Object.freeze([
            "/static/css/index.css",
            "/static/css/buttons.css",
        ])

        // Constructable stylesheets keep shared CSS deduplicated across shadow roots.
        const supportsConstructableStyleSheets =
            !!(Document.prototype && "adoptedStyleSheets" in Document.prototype) &&
            !!(CSSStyleSheet.prototype && "replace" in CSSStyleSheet.prototype)

        /** @type {Map<string, Promise<CSSStyleSheet[] | null>>} */
        const sharedStyleSheetPromises = new Map()

        /**
         * Public: load shared stylesheets once per URL combination.
         * @param {string[]} urls
         * @returns {Promise<CSSStyleSheet[] | null>}
         */
        function getStyleSheets(urls) {
            if (!supportsConstructableStyleSheets) return Promise.resolve(null)

            const cacheKey = urls.join("|")
            const existingPromise = sharedStyleSheetPromises.get(cacheKey)
            if (existingPromise) return existingPromise

            const styleSheetsPromise = Promise.all(
                urls.map(async (url) => {
                    const response = await fetch(url, { credentials: "same-origin" })
                    const cssText = await response.text()
                    const styleSheet = new CSSStyleSheet()
                    await styleSheet.replace(cssText)
                    return styleSheet
                })
            )

            sharedStyleSheetPromises.set(cacheKey, styleSheetsPromise)
            return styleSheetsPromise
        }

        /**
         * Public: combine shared base stylesheet URLs with component-local stylesheet URLs.
         * @param {string[]} [extraUrls=[]]
         * @returns {string[]}
         */
        function getStyleUrls(extraUrls = []) {
            return [...BASE_STYLE_URLS, ...extraUrls]
        }

        /**
         * Public: apply shared styles to a shadow root.
         * Uses adoptedStyleSheets when available, and falls back to <link> elements.
         * @param {ShadowRoot} shadowRoot
         * @param {string[]} urls
         * @returns {Promise<void>}
         */
        async function applyStyleUrlsToShadowRoot(shadowRoot, urls) {
            const styleSheets = await getStyleSheets(urls)

            if (styleSheets) {
                shadowRoot.adoptedStyleSheets = [...shadowRoot.adoptedStyleSheets, ...styleSheets]
                return
            }

            for (const url of urls) {
                const linkElement = document.createElement("link")
                linkElement.rel = "stylesheet"
                linkElement.href = url
                shadowRoot.appendChild(linkElement)
            }
        }

        return Object.freeze({
            applyStyleUrlsToShadowRoot,
            getStyleSheets,
            getStyleUrls,
        })
    }

    function createBillettholderSelection() {
        const SELECTED_BILLETTHOLDER_STORAGE_KEY = "selectedBillettHolder"
        const SELECTION_CHANGE_EVENT = "billettholder-selection-change"
        const ACCENT_COLORS = Object.freeze([
            "var(--color-accent-blue)",
            "var(--color-accent-purple)",
            "var(--color-accent-pink)",
            "var(--color-accent-red)",
            "var(--color-accent-orange)",
            "var(--color-accent-yellow)",
            "var(--color-accent-green)",
            "var(--color-accent-teal)",
            "var(--color-accent-cyan)",
        ])
        const ACCENT_COLOR_SET = new Set(ACCENT_COLORS)

        /**
         * Internal: accepts only color tokens produced by this helper.
         * @param {string} color
         * @returns {boolean}
         */
        function isBillettholderColor(color) {
            return ACCENT_COLOR_SET.has(color)
        }

        /**
         * Internal: normalize untrusted server/localStorage values into the shared shape.
         * @param {unknown} value
         * @returns {StoredBillettholder | null}
         */
        function normalizeBillettholder(value) {
            if (!value || typeof value !== "object") {
                return null
            }

            const billettholder = /** @type {{Id?: unknown, Name?: unknown, Email?: unknown, Color?: unknown}} */ (value)
            const id = Number(billettholder.Id ?? 0)
            if (!Number.isInteger(id) || id <= 0) {
                return null
            }

            const selectedBillettholder = {
                Id: id,
                Name: typeof billettholder.Name === "string" ? billettholder.Name : "",
                Email: typeof billettholder.Email === "string" ? billettholder.Email : "",
            }

            const color = typeof billettholder.Color === "string" ? billettholder.Color.trim() : ""
            if (isBillettholderColor(color)) {
                return { ...selectedBillettholder, Color: color }
            }

            // Repair older localStorage entries that have no color, or have a name stored as Color.
            if (selectedBillettholder.Name.length > 0) {
                return { ...selectedBillettholder, Color: colorFromName(selectedBillettholder.Name) }
            }

            return selectedBillettholder
        }

        /**
         * Internal: notify all selection subscribers after storage changes.
         * @param {StoredBillettholder | null} billettholder
         * @returns {void}
         */
        function dispatchSelectionChange(billettholder) {
            window.dispatchEvent(new CustomEvent(SELECTION_CHANGE_EVENT, { detail: billettholder }))
        }

        /**
         * Public: read the current selected billettholder from localStorage.
         * @returns {StoredBillettholder | null}
         */
        function get() {
            let storedBillettholderString = null
            try {
                storedBillettholderString = localStorage.getItem(SELECTED_BILLETTHOLDER_STORAGE_KEY)
            } catch {
                return null
            }
            if (!storedBillettholderString) {
                return null
            }

            try {
                return normalizeBillettholder(JSON.parse(storedBillettholderString))
            } catch {
                clear()
                return null
            }
        }

        /**
         * Public: store a selected billettholder and notify subscribers.
         * @param {unknown} billettholder
         * @returns {StoredBillettholder | null}
         */
        function set(billettholder) {
            const selectedBillettholder = normalizeBillettholder(billettholder)
            if (!selectedBillettholder) {
                return null
            }

            try {
                localStorage.setItem(SELECTED_BILLETTHOLDER_STORAGE_KEY, JSON.stringify(selectedBillettholder))
            } catch {
                return null
            }
            dispatchSelectionChange(selectedBillettholder)
            return selectedBillettholder
        }

        /**
         * Public: remove the selected billettholder and notify subscribers.
         * @returns {void}
         */
        function clear() {
            try {
                localStorage.removeItem(SELECTED_BILLETTHOLDER_STORAGE_KEY)
            } catch {
                return
            }
            dispatchSelectionChange(null)
        }

        /**
         * Public: find the canonical associated billettholder for a selected id.
         * @param {unknown[]} associatedBillettholdere
         * @param {unknown} billettholderId
         * @returns {StoredBillettholder | null}
         */
        function findAssociated(associatedBillettholdere, billettholderId) {
            if (!Array.isArray(associatedBillettholdere)) {
                return null
            }

            const associatedBillettholder = associatedBillettholdere.find((billettholder) => {
                const candidate = /** @type {{Id?: unknown}} */ (billettholder)
                return Number(candidate.Id ?? 0) === Number(billettholderId)
            })

            return normalizeBillettholder(associatedBillettholder)
        }

        /**
         * Public: choose the initial billettholder for the current page.
         * Keeps localStorage only when it still matches one of the associated billettholdere.
         * @param {unknown[]} associatedBillettholdere
         * @param {unknown} yourBillettholder
         * @returns {StoredBillettholder | null}
         */
        function initialize(associatedBillettholdere, yourBillettholder) {
            const storedBillettholder = get()
            if (storedBillettholder) {
                const selectedBillettholder = findAssociated(associatedBillettholdere, storedBillettholder.Id)
                if (selectedBillettholder) {
                    return set(selectedBillettholder)
                }

                clear()
            }

            const normalizedYourBillettholder = normalizeBillettholder(yourBillettholder)
            if (!normalizedYourBillettholder) {
                return null
            }

            const selectedBillettholder = findAssociated(associatedBillettholdere, normalizedYourBillettholder.Id)
            if (!selectedBillettholder) {
                return null
            }

            return set(selectedBillettholder)
        }

        /**
         * Public: subscribe to selected-billettholder changes.
         * Returns a cleanup function that removes this subscriber.
         * @param {(billettholder: StoredBillettholder | null) => void} callback
         * @returns {() => void}
         */
        function onChange(callback) {
            /** @type {(event: Event) => void} */
            const handleSelectionChange = (event) => {
                const selectionEvent = /** @type {CustomEvent<unknown>} */ (event)
                callback(normalizeBillettholder(selectionEvent.detail))
            }

            window.addEventListener(SELECTION_CHANGE_EVENT, handleSelectionChange)
            return () => window.removeEventListener(SELECTION_CHANGE_EVENT, handleSelectionChange)
        }

        /**
         * Public: build the initials badge background value for a CSSStyleDeclaration.
         * @param {string} color
         * @returns {string}
         */
        function initialsBackgroundColor(color) {
            if (!isBillettholderColor(color)) {
                return ""
            }

            return `hsl(from ${color} h s l / 0.5)`
        }

        /**
         * Public: build the initials badge border value for a CSSStyleDeclaration.
         * @param {string} color
         * @returns {string}
         */
        function initialsBorder(color) {
            if (!isBillettholderColor(color)) {
                return ""
            }

            return `1px solid ${color}`
        }

        /** @type {BillettholderSelectionStyle} */
        const style = Object.freeze({
            backgroundColor: initialsBackgroundColor,
            border: initialsBorder,
        })

        /**
         * Public: derive initials from a full name using first and last word.
         * @param {string} name
         * @returns {string}
         */
        function getInitials(name) {
            const words = name.trim().split(/\s+/).filter(Boolean)
            if (words.length === 0) {
                return "TT"
            }

            const firstInitial = words[0]?.[0] ?? "T"
            const lastInitial = words.at(-1)?.[0] ?? "T"

            return `${letterOrFallback(firstInitial)}${letterOrFallback(lastInitial)}`.toUpperCase()
        }

        /**
         * Internal: keep initials stable when a name part starts with punctuation or a digit.
         * @param {string} value
         * @returns {string}
         */
        function letterOrFallback(value) {
            return /\p{L}/u.test(value) ? value : "T"
        }

        /**
         * Public: derive the stable accent color token for a billettholder name.
         * @param {string} name
         * @returns {string}
         */
        function colorFromName(name) {
            let hash = 2166136261
            for (let index = 0; index < name.length; index += 1) {
                hash ^= name.charCodeAt(index)
                hash = Math.imul(hash, 16777619)
            }

            return ACCENT_COLORS[(hash >>> 0) % ACCENT_COLORS.length]
        }

        return Object.freeze({
            clear,
            colorFromName,
            findAssociated,
            get,
            getInitials,
            initialize,
            onChange,
            set,
            style,
            storageKey: SELECTED_BILLETTHOLDER_STORAGE_KEY,
        })
    }
})()
