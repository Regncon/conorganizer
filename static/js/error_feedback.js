// @ts-check

const fallbackMessage = "Klarte ikkje å lagre endringa. Prøv igjen."
const rootSelector = "form, dialog"

/**
 * @typedef {"started" | "finished" | "error" | "retrying" | "retries-failed"} DatastarFetchType
 * @typedef {{ type: DatastarFetchType, el: Element | null }} DatastarFetchEventDetail
 */

document.addEventListener("datastar-fetch", handleDatastarFetch)

/**
 * @param {Event} event
 * @returns {void}
 */
function handleDatastarFetch(event) {
    const detail = getDatastarFetchDetail(event)
    if (!detail || !(detail.el instanceof HTMLElement)) {
        return
    }

    const key = detail.el.dataset.feedbackKey?.trim()
    if (!key) {
        return
    }

    if (detail.type === "started") {
        updateFeedback(detail.el, key, "")
        return
    }

    if (detail.type === "error" || detail.type === "retrying" || detail.type === "retries-failed") {
        updateFeedback(detail.el, key, feedbackMessage(detail.el))
    }
}

/**
 * @param {Event} event
 * @returns {DatastarFetchEventDetail | null}
 */
function getDatastarFetchDetail(event) {
    if (!(event instanceof CustomEvent)) {
        return null
    }

    /** @type {unknown} */
    const detail = event.detail
    return isDatastarFetchDetail(detail) ? detail : null
}

/**
 * @param {unknown} detail
 * @returns {detail is DatastarFetchEventDetail}
 */
function isDatastarFetchDetail(detail) {
    if (!detail || typeof detail !== "object") {
        return false
    }

    /** @type {{ type?: unknown, el?: unknown }} */
    const candidate = detail
    return isDatastarFetchType(candidate.type)
        && (candidate.el === null || candidate.el instanceof Element)
}

/**
 * @param {unknown} type
 * @returns {type is DatastarFetchType}
 */
function isDatastarFetchType(type) {
    return type === "started"
        || type === "finished"
        || type === "error"
        || type === "retrying"
        || type === "retries-failed"
}

/**
 * @param {HTMLElement} trigger
 * @returns {string}
 */
function feedbackMessage(trigger) {
    const root = feedbackRoot(trigger)
    return trigger.dataset.feedbackMessage?.trim()
        || (root instanceof HTMLElement ? root.dataset.feedbackDefaultMessage?.trim() : "")
        || fallbackMessage
}

/**
 * @param {HTMLElement} trigger
 * @param {string} key
 * @param {string} message
 * @returns {void}
 */
function updateFeedback(trigger, key, message) {
    for (const target of feedbackTargets(trigger, key)) {
        target.textContent = message
        target.hidden = message === ""
        target.classList.toggle("is-visible", message !== "")
    }
}

/**
 * @param {HTMLElement} trigger
 * @param {string} key
 * @returns {HTMLElement[]}
 */
function feedbackTargets(trigger, key) {
    const root = feedbackRoot(trigger)
    return [...root.querySelectorAll("[data-feedback-for]")].filter((target) => (
        target instanceof HTMLElement && target.dataset.feedbackFor === key
    ))
}

/**
 * @param {HTMLElement} trigger
 * @returns {Document | HTMLElement}
 */
function feedbackRoot(trigger) {
    return trigger.closest(rootSelector) || document
}
