/**
 * Datastar signal patches can be nested objects. The toast component flattens
 * them into dot paths so both `toastName` and nested signal paths can be handled.
 * @typedef {Record<string, unknown>} SignalPatch
 */

const defaultIndicatorPrefix = "toast"
const defaultIndicatorMessage = "Dette er en toast!!!"
const toastDurationMs = 3000
const toastExitMs = 180

/** @type {ToastMessage | null} */
let activeHost = null

/**
 * @param {unknown} signals
 * @param {string} [prefix]
 * @returns {Array<[string, unknown]>}
 */
function flattenSignals(signals, prefix = "") {
    if (!signals || typeof signals !== "object") {
        return []
    }

    return Object.entries(signals).flatMap(([key, value]) => {
        const path = prefix ? `${ prefix }.${ key }` : key
        if (value && typeof value === "object" && !Array.isArray(value)) {
            return flattenSignals(value, path)
        }
        return [[path, value]]
    })
}

class ToastMessage extends HTMLElement {
    constructor() {
        super()

        /** @type {Set<string>} */
        this.activeIndicators = new Set()

        /** @type {(event: Event) => void} */
        this.onSignalPatch = (event) => this.handleSignalPatch(event)

        /** @type {(event: Event) => void} */
        this.onToast = (event) => this.handleToast(event)
    }

    /** @returns {void} */
    connectedCallback() {
        // Keep the fixed-position region outside page containers that may use containment.
        if (this.parentElement !== document.body) {
            document.body.append(this)
            return
        }

        // A page can render multiple toast hosts through partials; only one should listen.
        if (activeHost) {
            return
        }

        activeHost = this
        document.addEventListener("datastar-signal-patch", this.onSignalPatch)
        document.addEventListener("toast", this.onToast)
    }

    /** @returns {void} */
    disconnectedCallback() {
        if (activeHost !== this) {
            return
        }

        document.removeEventListener("datastar-signal-patch", this.onSignalPatch)
        document.removeEventListener("toast", this.onToast)
        activeHost = null
    }

    /** @returns {string} */
    get indicatorPrefix() {
        return this.getAttribute("indicator-prefix")?.trim() || defaultIndicatorPrefix
    }

    /** @returns {string} */
    get indicatorMessage() {
        return this.getAttribute("indicator-message")?.trim() || defaultIndicatorMessage
    }

    /**
     * Shows a toast when a tracked Datastar indicator goes from active to idle.
     * @param {Event} event
     * @returns {void}
     */
    handleSignalPatch(event) {
        const detail = event instanceof CustomEvent ? event.detail : null

        for (const [path, value] of flattenSignals(detail)) {
            if (!path.startsWith(this.indicatorPrefix)) {
                continue
            }

            if (value === true) {
                this.activeIndicators.add(path)
                continue
            }

            if (value === false && this.activeIndicators.has(path)) {
                this.activeIndicators.delete(path)
                this.show(this.indicatorMessage)
            }
        }
    }

    /**
     * Public event API for manual toasts from other web components or page scripts.
     * @param {Event} event
     * @returns {void}
     */
    handleToast(event) {
        /** @type {{ message?: unknown } | null} */
        const detail = event instanceof CustomEvent && event.detail && typeof event.detail === "object"
            ? event.detail
            : null
        const message = detail?.message
        if (typeof message !== "string" || message.trim() === "") {
            console.warn("toast event requires detail.message")
            return
        }

        this.show(message)
    }

    /**
     * @param {string} message
     * @returns {void}
     */
    show(message) {
        /** @type {HTMLTemplateElement | null} */
        const template = this.querySelector("template")
        const templateToast = template?.content.firstElementChild?.cloneNode(true)
        const toast = templateToast instanceof HTMLElement ? templateToast : this.createFallbackToast()
        const text = toast.querySelector("[data-toast-message-text]")
        if (text) {
            text.textContent = message
        }

        this.append(toast)
        window.setTimeout(() => this.removeToast(toast), toastDurationMs)
    }

    /** @returns {HTMLElement} */
    createFallbackToast() {
        const toast = document.createElement("div")
        toast.className = "toast"
        const text = document.createElement("span")
        text.setAttribute("data-toast-message-text", "")
        toast.append(text)
        return toast
    }

    /**
     * Wait for the exit animation before removing the toast from the DOM.
     * @param {HTMLElement} toast
     * @returns {void}
     */
    removeToast(toast) {
        toast.classList.add("is-removing")
        window.setTimeout(() => toast.remove(), toastExitMs)
    }
}

if (!customElements.get("toast-message")) {
    customElements.define("toast-message", ToastMessage)
}
