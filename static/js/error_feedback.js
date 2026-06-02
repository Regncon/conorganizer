// @ts-check

const fallbackFeedbackMessage = "Klarte ikkje å lagre endringa. Prøv igjen."
const feedbackRootSelector = "form, dialog"
const feedbackTargetSelector = "[data-feedback-for]"

/**
 * @typedef {"started" | "finished" | "error" | "retrying" | "retries-failed"} DatastarFetchType
 * @typedef {{
 *     readonly type: DatastarFetchType,
 *     readonly el: Element | null,
 *     readonly argsRaw?: Record<string, string>,
 * }} DatastarFetchEventDetail
 */

// Datastar emits `finished` after both success and failure; these triggers keep their error visible.
/** @type {WeakSet<HTMLElement>} */
const failedRequestTriggers = new WeakSet()

/** @type {ReadonlySet<DatastarFetchType>} */
const datastarFetchTypes = new Set([
    "started",
    "finished",
    "error",
    "retrying",
    "retries-failed",
])

/** @type {ReadonlySet<DatastarFetchType>} */
const lastingFailureTypes = new Set(["error", "retries-failed"])

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

    const trigger = detail.el
    const key = trigger.dataset.feedbackKey?.trim()
    if (!key) {
        return
    }

    if (detail.type === "started") {
        failedRequestTriggers.delete(trigger)
        updateFeedback(trigger, key, "")
        return
    }

    if (detail.type === "retrying") {
        updateFeedback(trigger, key, feedbackMessage(trigger))
        return
    }

    if (lastingFailureTypes.has(detail.type)) {
        failedRequestTriggers.add(trigger)
        updateFeedback(trigger, key, feedbackMessage(trigger))
        return
    }

    if (detail.type === "finished" && !failedRequestTriggers.has(trigger)) {
        updateFeedback(trigger, key, "")
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
    return typeof type === "string" && datastarFetchTypes.has(/** @type {DatastarFetchType} */(type))
}

/**
 * @param {HTMLElement} trigger
 * @returns {string}
 */
function feedbackMessage(trigger) {
    const root = feedbackRoot(trigger)
    return trigger.dataset.feedbackMessage?.trim()
        || (root instanceof HTMLElement ? root.dataset.feedbackDefaultMessage?.trim() : "")
        || fallbackFeedbackMessage
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
    return [...root.querySelectorAll(feedbackTargetSelector)].filter((target) => (
        isFeedbackTargetForKey(target, key)
    ))
}

/**
 * @param {Element} target
 * @param {string} key
 * @returns {target is HTMLElement}
 */
function isFeedbackTargetForKey(target, key) {
    return target instanceof HTMLElement && target.dataset.feedbackFor === key
}

/**
 * @param {HTMLElement} trigger
 * @returns {Document | HTMLElement}
 */
function feedbackRoot(trigger) {
    return trigger.closest(feedbackRootSelector) || document
}
