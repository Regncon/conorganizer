// @ts-check

(() => {
    /**
     * @typedef {Object} StoredBillettholder
     * @property {number} Id
     * @property {string} Name
     * @property {string} Email
     * @property {string} [Color]
     */
    /**
     * @typedef {Object} SharedStyles
     * @property {(urls: string[]) => Promise<CSSStyleSheet[] | null>} getStyleSheets
     * @property {(extraUrls?: string[]) => string[]} getStyleUrls
     * @property {(shadowRoot: ShadowRoot, urls: string[]) => Promise<void>} applyStyleUrlsToShadowRoot
     */
    /**
     * @typedef {Object} BillettholderSelectionStyle
     * @property {(color: string) => string} backgroundColor
     * @property {(color: string) => string} border
     */
    /**
     * @typedef {Object} BillettholderSelection
     * @property {() => void} clear
     * @property {(name: string) => string} colorFromName
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

        const supportsConstructableStyleSheets =
            !!(Document.prototype && "adoptedStyleSheets" in Document.prototype) &&
            !!(CSSStyleSheet.prototype && "replace" in CSSStyleSheet.prototype)

        /** @type {Map<string, Promise<CSSStyleSheet[] | null>>} */
        const sharedStyleSheetPromises = new Map()

        /**
         * Load shared stylesheets once per URL combination.
         * @param {string[]} urls
         * @returns {Promise<CSSStyleSheet[] | null>}
         */
        function getStyleSheets(urls) {
            if (!supportsConstructableStyleSheets) return Promise.resolve(null)

            const cacheKey = urls.join("|")
            const existingPromise = sharedStyleSheetPromises.get(cacheKey)
            if (existingPromise) return existingPromise

            const styleSheetsPromise = (async () => {
                const styleSheets = []
                for (const url of urls) {
                    const response = await fetch(url, { credentials: "same-origin" })
                    const cssText = await response.text()
                    const styleSheet = new CSSStyleSheet()
                    await styleSheet.replace(cssText)
                    styleSheets.push(styleSheet)
                }
                return styleSheets
            })()

            sharedStyleSheetPromises.set(cacheKey, styleSheetsPromise)
            return styleSheetsPromise
        }

        /**
         * Combine the shared base styles with optional component-local styles.
         * @param {string[]} [extraUrls=[]]
         * @returns {string[]}
         */
        function getStyleUrls(extraUrls = []) {
            return [...BASE_STYLE_URLS, ...extraUrls]
        }

        /**
         * Apply shared styles to a shadow root via adoptedStyleSheets when supported,
         * otherwise fall back to injecting <link> elements.
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

        /**
         * @param {string} color
         * @returns {boolean}
         */
        function isBillettholderColor(color) {
            return ACCENT_COLORS.includes(color)
        }

        /**
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

            if (selectedBillettholder.Name.length > 0) {
                return { ...selectedBillettholder, Color: colorFromName(selectedBillettholder.Name) }
            }

            return selectedBillettholder
        }

        /**
         * @param {StoredBillettholder | null} billettholder
         * @returns {void}
         */
        function dispatchSelectionChange(billettholder) {
            window.dispatchEvent(new CustomEvent(SELECTION_CHANGE_EVENT, { detail: billettholder }))
        }

        /**
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
         * @param {unknown[]} associatedBillettholdere
         * @param {unknown} yourBillettholder
         * @returns {StoredBillettholder | null}
         */
        function initialize(associatedBillettholdere, yourBillettholder) {
            const candidates = Array.isArray(associatedBillettholdere) ? associatedBillettholdere : []
            /**
             * @param {unknown} billettholderId
             * @returns {unknown | undefined}
             */
            const findAssociatedBillettholder = (billettholderId) => candidates.find((billettholder) => {
                const candidate = /** @type {{Id?: unknown}} */ (billettholder)
                return Number(candidate.Id ?? 0) === Number(billettholderId)
            })

            const storedBillettholder = get()
            if (storedBillettholder) {
                const selectedBillettholder = findAssociatedBillettholder(storedBillettholder.Id)
                if (selectedBillettholder) {
                    return set(selectedBillettholder)
                }

                clear()
            }

            const normalizedYourBillettholder = normalizeBillettholder(yourBillettholder)
            if (!normalizedYourBillettholder) {
                return null
            }

            const selectedBillettholder = findAssociatedBillettholder(normalizedYourBillettholder.Id)
            if (!selectedBillettholder) {
                return null
            }

            return set(selectedBillettholder)
        }

        /**
         * @param {(billettholder: StoredBillettholder | null) => void} callback
         * @returns {() => void}
         */
        function onChange(callback) {
            const handleSelectionChange = (event) => {
                const selectionEvent = /** @type {CustomEvent<unknown>} */ (event)
                callback(normalizeBillettholder(selectionEvent.detail))
            }

            window.addEventListener(SELECTION_CHANGE_EVENT, handleSelectionChange)
            return () => window.removeEventListener(SELECTION_CHANGE_EVENT, handleSelectionChange)
        }

        /**
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
         * @param {string} name
         * @returns {string}
         */
        function getInitials(name) {
            const words = name.trim().split(/\s+/).filter(Boolean)
            if (words.length === 0) {
                return "TT"
            }

            const firstInitial = words[0]?.[0] ?? "T"
            const lastInitial = words[words.length - 1]?.[0] ?? "T"

            return `${letterOrFallback(firstInitial)}${letterOrFallback(lastInitial)}`.toUpperCase()
        }

        /**
         * @param {string} value
         * @returns {string}
         */
        function letterOrFallback(value) {
            return /\p{L}/u.test(value) ? value : "T"
        }

        /**
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
