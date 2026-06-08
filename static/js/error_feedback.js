import { mergePatch } from "../datastar.js"

const feedbackSignalName = "feedbackErrors"
const localFeedbackPatchStartEventName = "feedback-errors-local-patch-start"
const localFeedbackPatchEndEventName = "feedback-errors-local-patch-end"
const fallbackFeedbackMessage = "Klarte ikkje å lagre endringa. Prøv igjen."
const feedbackRootSelector = "form, dialog"

/**
 * @typedef {"started" | "finished" | "error" | "retrying" | "retries-failed"} DatastarFetchType
 * @typedef {{
 *     readonly type: DatastarFetchType,
 *     readonly el: Element | null,
 *     readonly argsRaw?: Record<string, string>,
 * }} DatastarFetchEventDetail
 * @typedef {Record<string, string>} FeedbackErrors
 * @typedef {{ feedbackErrors?: FeedbackErrors }} DatastarSignalPatchDetail
 */

/** @type {Map<HTMLElement, string>} */
const activeRequestTriggers = new Map()

/** @type {number} */
let localFeedbackPatchDepth = 0

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
document.addEventListener("datastar-signal-patch", handleDatastarSignalPatch)

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
        activeRequestTriggers.set(trigger, key)
        failedRequestTriggers.delete(trigger)
        patchFeedbackError(key, "")
        return
    }

    if (detail.type === "retrying") {
        patchFeedbackError(key, feedbackMessage(trigger))
        return
    }

    if (lastingFailureTypes.has(detail.type)) {
        failedRequestTriggers.add(trigger)
        patchFeedbackError(key, feedbackMessage(trigger))
        return
    }

    if (detail.type === "finished") {
        activeRequestTriggers.delete(trigger)
        if (!failedRequestTriggers.has(trigger)) {
            patchFeedbackError(key, "")
        }
    }
}

/**
 * Datastar custom backend feedback arrives as signal patches:
 * `{ feedbackErrors: { [key]: message } }`.
 * The message elements render the signal directly; this listener only marks active
 * requests with backend-patched errors so the later `finished` fetch event does
 * not clear the backend message. Local fallback patches are ignored here.
 *
 * @param {Event} event
 * @returns {void}
 */
function handleDatastarSignalPatch(event) {
    if (localFeedbackPatchDepth > 0) {
        return
    }

    const detail = getDatastarSignalPatchDetail(event)
    const feedbackErrors = detail?.[feedbackSignalName]
    if (!feedbackErrors) {
        return
    }

    markActiveTriggersWithPatchedErrors(feedbackErrors)
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
 * @param {Event} event
 * @returns {DatastarSignalPatchDetail | null}
 */
function getDatastarSignalPatchDetail(event) {
    if (!(event instanceof CustomEvent)) {
        return null
    }

    /** @type {unknown} */
    const detail = event.detail
    return isDatastarSignalPatchDetail(detail)
        ? /** @type {DatastarSignalPatchDetail} */ (detail)
        : null
}

/**
 * @param {unknown} detail
 * @returns {detail is DatastarSignalPatchDetail}
 */
function isDatastarSignalPatchDetail(detail) {
    if (!detail || typeof detail !== "object") {
        return false
    }

    /** @type {{ feedbackErrors?: unknown }} */
    const candidate = detail
    return candidate.feedbackErrors === undefined || isFeedbackErrors(candidate.feedbackErrors)
}

/**
 * @param {unknown} feedbackErrors
 * @returns {feedbackErrors is FeedbackErrors}
 */
function isFeedbackErrors(feedbackErrors) {
    return !!feedbackErrors
        && typeof feedbackErrors === "object"
        && Object.values(feedbackErrors).every((message) => typeof message === "string")
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
 * @param {FeedbackErrors} feedbackErrors
 * @returns {void}
 */
function markActiveTriggersWithPatchedErrors(feedbackErrors) {
    for (const [trigger, key] of activeRequestTriggers.entries()) {
        const message = feedbackErrors[key]
        if (message && message.trim() !== "") {
            failedRequestTriggers.add(trigger)
        }
    }
}

/**
 * @param {string} key
 * @param {string} message
 * @returns {void}
 */
function patchFeedbackError(key, message) {
    localFeedbackPatchDepth += 1
    document.dispatchEvent(new CustomEvent(localFeedbackPatchStartEventName))
    try {
        mergePatch({
            [feedbackSignalName]: {
                [key]: message,
            },
        })
    } finally {
        document.dispatchEvent(new CustomEvent(localFeedbackPatchEndEventName))
        localFeedbackPatchDepth = Math.max(0, localFeedbackPatchDepth - 1)
    }
}

/**
 * @param {HTMLElement} trigger
 * @returns {Document | HTMLElement}
 */
function feedbackRoot(trigger) {
    return trigger.closest(feedbackRootSelector) || document
}
