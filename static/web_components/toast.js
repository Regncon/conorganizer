// @ts-check

const CUSTOM_ELEMENT_TAG_NAME = "app-toast"

/**
 * @typedef {{ message?: unknown }} ToastEventDetail
 * @typedef {"started" | "finished" | "error" | "retrying" | "retries-failed"} DatastarFetchType
 * @typedef {{ type: DatastarFetchType, el: Element | null }} DatastarFetchEventDetail
 * @typedef {Record<string, string>} FeedbackErrors
 * @typedef {{ feedbackErrors?: FeedbackErrors }} DatastarSignalPatchDetail
 * @typedef {{ activeCount: number, failed: boolean }} IndicatorState
 */

/** @type {AppToast | null} */
let activeAppToast = null

class AppToast extends HTMLElement {
    static #defaultIndicatorPrefix = "toast"
    static #missingMessage = "Mangler toast-melding"
    static #toastDurationMs = 3000
    static #toastExitMs = 180

    static {
        if (!customElements.get(CUSTOM_ELEMENT_TAG_NAME)) {
            customElements.define(CUSTOM_ELEMENT_TAG_NAME, AppToast)
        }
    }

    /** @type {Map<string, IndicatorState>} */
    #indicatorStates = new Map()

    /** @type {(event: Event) => void} */
    #onDatastarFetch = (event) => this.#handleDatastarFetch(event)

    /** @type {(event: Event) => void} */
    #onDatastarSignalPatch = (event) => this.#handleDatastarSignalPatch(event)

    /** @type {(event: Event) => void} */
    #onToastEvent = (event) => this.#handleToastEvent(event)

    /** @returns {void} */
    connectedCallback() {
        // Keep the fixed-position region outside page containers that may use containment.
        if (this.parentElement !== document.body) {
            document.body.append(this)
            return
        }

        // A page can render several hosts through partials; only one should listen.
        if (activeAppToast) {
            return
        }

        activeAppToast = this
        document.addEventListener("datastar-fetch", this.#onDatastarFetch)
        document.addEventListener("datastar-signal-patch", this.#onDatastarSignalPatch)
        document.addEventListener("toast", this.#onToastEvent)
    }

    /** @returns {void} */
    disconnectedCallback() {
        if (activeAppToast !== this) {
            return
        }

        document.removeEventListener("datastar-fetch", this.#onDatastarFetch)
        document.removeEventListener("datastar-signal-patch", this.#onDatastarSignalPatch)
        document.removeEventListener("toast", this.#onToastEvent)
        activeAppToast = null
    }

    /** @returns {string} */
    get #indicatorPrefix() {
        return this.getAttribute("indicator-prefix")?.trim() || AppToast.#defaultIndicatorPrefix
    }

    /** @returns {string} */
    get #indicatorMessage() {
        return this.getAttribute("indicator-message")?.trim() || AppToast.#missingMessage
    }

    /**
     * Shows a toast when a tracked Datastar request finishes without an error.
     * @param {Event} event
     * @returns {void}
     */
    #handleDatastarFetch(event) {
        const detail = AppToast.#getDatastarFetchDetail(event)
        if (!detail) {
            return
        }

        const indicator = detail.el instanceof HTMLElement
            ? detail.el.dataset.indicator?.trim()
            : ""

        if (!indicator?.startsWith(this.#indicatorPrefix)) {
            return
        }

        if (detail.type === "started") {
            const state = this.#getIndicatorState(indicator)
            if (state.activeCount === 0) {
                state.failed = false
            }
            state.activeCount += 1
            return
        }

        if (detail.type === "error" || detail.type === "retries-failed") {
            this.#getIndicatorState(indicator).failed = true
            return
        }

        if (detail.type !== "finished") {
            return
        }

        const state = this.#indicatorStates.get(indicator)
        if (!state) {
            return
        }

        state.activeCount = Math.max(0, state.activeCount - 1)
        if (state.activeCount > 0) {
            return
        }

        this.#indicatorStates.delete(indicator)
        if (state.failed) {
            return
        }

        this.#showToast(this.#indicatorMessage)
    }

    /**
     * Datastar custom backend feedback arrives as signal patches. Non-empty
     * `feedbackErrors` marks active requests as failed so their later `finished`
     * fetch event does not show a saved toast.
     * @param {Event} event
     * @returns {void}
     */
    #handleDatastarSignalPatch(event) {
        const detail = AppToast.#getDatastarSignalPatchDetail(event)
        if (!detail || !AppToast.#hasFeedbackErrors(detail.feedbackErrors)) {
            return
        }

        for (const state of this.#indicatorStates.values()) {
            state.failed = true
        }
    }

    /**
     * @param {string} indicator
     * @returns {IndicatorState}
     */
    #getIndicatorState(indicator) {
        const state = this.#indicatorStates.get(indicator)
        if (state) {
            return state
        }

        const newState = { activeCount: 0, failed: false }
        this.#indicatorStates.set(indicator, newState)
        return newState
    }

    /**
     * @param {Event} event
     * @returns {DatastarFetchEventDetail | null}
     */
    static #getDatastarFetchDetail(event) {
        if (!(event instanceof CustomEvent)) {
            return null
        }

        /** @type {unknown} */
        const detail = event.detail
        return AppToast.#isDatastarFetchDetail(detail) ? detail : null
    }

    /**
     * @param {unknown} detail
     * @returns {detail is DatastarFetchEventDetail}
     */
    static #isDatastarFetchDetail(detail) {
        if (!detail || typeof detail !== "object") {
            return false
        }

        /** @type {{ type?: unknown, el?: unknown }} */
        const candidate = detail
        return AppToast.#isDatastarFetchType(candidate.type)
            && (candidate.el === null || candidate.el instanceof Element)
    }

    /**
     * @param {unknown} type
     * @returns {type is DatastarFetchType}
     */
    static #isDatastarFetchType(type) {
        return type === "started"
            || type === "finished"
            || type === "error"
            || type === "retrying"
            || type === "retries-failed"
    }

    /**
     * @param {Event} event
     * @returns {DatastarSignalPatchDetail | null}
     */
    static #getDatastarSignalPatchDetail(event) {
        if (!(event instanceof CustomEvent)) {
            return null
        }

        /** @type {unknown} */
        const detail = event.detail
        return AppToast.#isDatastarSignalPatchDetail(detail)
            ? /** @type {DatastarSignalPatchDetail} */ (detail)
            : null
    }

    /**
     * @param {unknown} detail
     * @returns {detail is DatastarSignalPatchDetail}
     */
    static #isDatastarSignalPatchDetail(detail) {
        if (!detail || typeof detail !== "object") {
            return false
        }

        /** @type {{ feedbackErrors?: unknown }} */
        const candidate = detail
        return candidate.feedbackErrors === undefined || AppToast.#isFeedbackErrors(candidate.feedbackErrors)
    }

    /**
     * @param {unknown} feedbackErrors
     * @returns {feedbackErrors is FeedbackErrors}
     */
    static #isFeedbackErrors(feedbackErrors) {
        return !!feedbackErrors
            && typeof feedbackErrors === "object"
            && Object.values(feedbackErrors).every((message) => typeof message === "string")
    }

    /**
     * @param {FeedbackErrors | undefined} feedbackErrors
     * @returns {boolean}
     */
    static #hasFeedbackErrors(feedbackErrors) {
        if (!feedbackErrors) {
            return false
        }

        return Object.values(feedbackErrors).some((message) => (
            typeof message === "string" && message.trim() !== ""
        ))
    }

    /**
     * Public event API for manual toasts from other web components or page scripts.
     * @param {Event} event
     * @returns {void}
     */
    #handleToastEvent(event) {
        const detail = AppToast.#getToastEventDetail(event)
        const message = detail?.message
        if (typeof message !== "string" || message.trim() === "") {
            console.warn("toast event requires detail.message")
            return
        }

        this.#showToast(message)
    }

    /**
     * @param {Event} event
     * @returns {ToastEventDetail | null}
     */
    static #getToastEventDetail(event) {
        if (!(event instanceof CustomEvent)) {
            return null
        }

        /** @type {unknown} */
        const detail = event.detail
        return detail && typeof detail === "object" ? detail : null
    }

    /**
     * @param {string} message
     * @returns {void}
     */
    #showToast(message) {
        /** @type {HTMLTemplateElement | null} */
        const template = this.querySelector("template")
        const templateToast = template?.content.firstElementChild?.cloneNode(true)
        const toast = templateToast instanceof HTMLElement ? templateToast : AppToast.#createFallbackToastElement()
        const text = toast.querySelector("[data-toast-text]")
        if (text) {
            text.textContent = message
        }

        this.append(toast)
        window.setTimeout(() => AppToast.#removeToastElement(toast), AppToast.#toastDurationMs)
    }

    /** @returns {HTMLElement} */
    static #createFallbackToastElement() {
        const toast = document.createElement("div")
        toast.className = "toast"
        const text = document.createElement("span")
        text.setAttribute("data-toast-text", "")
        toast.append(text)
        return toast
    }

    /**
     * Wait for the exit animation before removing the toast from the DOM.
     * @param {HTMLElement} toast
     * @returns {void}
     */
    static #removeToastElement(toast) {
        toast.classList.add("is-removing")
        window.setTimeout(() => toast.remove(), AppToast.#toastExitMs)
    }

}
