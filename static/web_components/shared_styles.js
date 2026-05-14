// @ts-check

(() => {
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

    window.conorganizerSharedStyles = Object.freeze({
        applyStyleUrlsToShadowRoot,
        getStyleUrls,
        getStyleSheets,
    })
})()
