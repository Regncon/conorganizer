// @ts-check

import { action, mergePatch, mergePaths } from "../../datastar.js"
import { imageToWebp } from "../image-to-webp.js"

/** @typedef {"started" | "finished" | "error"} UploadRequestEventType */
/** @typedef {{ message?: string, status?: string }} UploadRequestEventDetails */

/**
 * Upload a source image with byte progress exposed as a Datastar signal.
 * XMLHttpRequest supplies the progress events that fetch does not expose.
 * Datastar still owns the signals, request indicators, and page feedback.
 *
 * Usage on a form containing an input named "image":
 *   data-on:submit__prevent="@uploadWithProgress('/profile/api/new/123/upload',
 *       { progressSignal: '_uploadProgress', convertToWebp: true })"
 *   data-indicator="_uploading"
 *
 * progressSignal is an optional signal path without the "$" prefix. Its value
 * is a percentage from 0 to 100 describing transfer of the request body.
 * Even at 100, the request may still be waiting for the server to save the image.
 * convertToWebp optionally converts the "image" file before sending it. The
 * request indicator covers conversion too; byte progress starts when sending.
 *
 * A successful response must be a JSON signal patch with a nonempty sourceImageUrl,
 * for example { "sourceImageUrl": "/event-images/123_source.webp?v=456" }.
 * This adapter handles JSON responses, not Datastar SSE streams.
 */
action({
    name: "uploadWithProgress",
    /**
     * @param {{ el: Element, cleanups: Map<string, () => void> }} context Datastar's action context.
     * @param {string} uploadUrl URL of the source-image upload endpoint.
     * @param {{ progressSignal?: string, convertToWebp?: boolean }} [options] Progress signal and optional source-image conversion.
     * @returns {Promise<void> | undefined} Resolves after success or failure; errors are reported through request events.
     */
    apply({ el: triggerElement, cleanups: actionCleanups }, uploadUrl, { progressSignal: progressSignalPath, convertToWebp = false } = {}) {
        const uploadForm = triggerElement.closest("form")
        if (!uploadForm || (!uploadForm.noValidate && !uploadForm.reportValidity())) return

        // Datastar provides a cleanup map per expression. The entry also guards
        // against starting another upload from that expression before this one ends.
        const uploadCleanupKey = "@uploadWithProgress"
        if (actionCleanups.has(uploadCleanupKey)) return

        // Capture the files before "started" activates data-indicator: disabled
        // inputs are excluded when FormData reads a form.
        const formData = new FormData(uploadForm)
        const uploadRequest = new XMLHttpRequest()

        // Keep Datastar's event field names. Indicators use detail.el to match
        // the request to its trigger; page handlers read errors from argsRaw.
        /**
         * @param {UploadRequestEventType} eventType
         * @param {UploadRequestEventDetails} [eventDetails]
         * @returns {boolean}
         */
        const dispatchRequestEvent = (eventType, eventDetails = {}) => document.dispatchEvent(
            new CustomEvent("datastar-fetch", {
                detail: { type: eventType, el: triggerElement, argsRaw: eventDetails },
            }),
        )
        /**
         * @param {number} percentComplete
         * @returns {void}
         */
        const updateProgressSignal = (percentComplete) => {
            if (progressSignalPath) mergePaths([[progressSignalPath, percentComplete]])
        }

        return new Promise((resolveUpload) => {
            let hasFinished = false
            /**
             * @param {string} [errorMessage] Omit when the upload succeeds.
             * @returns {void}
             */
            const finishUpload = (errorMessage) => {
                // Every outcome must release the indicator exactly once. On
                // failure, "error" precedes "finished" so handlers can suppress success feedback.
                if (hasFinished) return
                hasFinished = true
                if (errorMessage) {
                    dispatchRequestEvent("error", { message: errorMessage, status: String(uploadRequest.status) })
                }
                actionCleanups.delete(uploadCleanupKey)
                dispatchRequestEvent("finished")
                resolveUpload()
            }

            // Cleanup must also finish an upload that is still converting its
            // image: abort() alone does not emit an event before XHR is opened.
            actionCleanups.set(uploadCleanupKey, () => {
                uploadRequest.abort()
                finishUpload("Opplastingen ble avbrutt.")
            })
            uploadRequest.upload.addEventListener("progress", (progressEvent) => {
                if (progressEvent.lengthComputable) {
                    updateProgressSignal(Math.min(100, Math.round(progressEvent.loaded / progressEvent.total * 100)))
                }
            })
            uploadRequest.addEventListener("load", () => {
                // XHR also fires "load" for HTTP errors such as 400 and 500.
                if (uploadRequest.status < 200 || uploadRequest.status >= 300) {
                    finishUpload(uploadRequest.responseText || "Kunne ikke laste opp bildet. Prøv igjen.")
                    return
                }
                try {
                    /** @type {unknown} */
                    const responseSignals = JSON.parse(uploadRequest.responseText)
                    if (
                        typeof responseSignals !== "object" ||
                        responseSignals === null ||
                        !("sourceImageUrl" in responseSignals) ||
                        typeof responseSignals.sourceImageUrl !== "string" ||
                        !responseSignals.sourceImageUrl
                    ) {
                        throw new Error("Missing sourceImageUrl")
                    }
                    updateProgressSignal(100)
                    // XHR bypasses @post's response handling. Apply the server's
                    // signal patch before "finished" lets the page react to success.
                    mergePatch(responseSignals)
                    finishUpload()
                } catch {
                    finishUpload("Kunne ikke lese svaret fra opplastingen. Prøv igjen.")
                }
            })
            uploadRequest.addEventListener("error", () => finishUpload("Kunne ikke laste opp bildet. Prøv igjen."))
            uploadRequest.addEventListener("abort", () => finishUpload("Opplastingen ble avbrutt."))
            uploadRequest.addEventListener("timeout", () => finishUpload("Opplastingen tok for lang tid. Prøv igjen."))

            updateProgressSignal(0)
            dispatchRequestEvent("started")
            const sendUpload = async () => {
                if (convertToWebp) {
                    try {
                        const sourceFile = formData.get("image")
                        if (!(sourceFile instanceof File) || sourceFile.size === 0) {
                            throw new Error("Missing source image")
                        }
                        formData.set("image", await imageToWebp(sourceFile))
                    } catch {
                        finishUpload("Kunne ikke konvertere bildet til WebP. Prøv et annet bilde.")
                        return
                    }
                }
                // Conversion is asynchronous; the binding may have been removed.
                if (hasFinished) return
                try {
                    uploadRequest.open("POST", uploadUrl)
                    uploadRequest.setRequestHeader("Accept", "application/json")
                    uploadRequest.setRequestHeader("Datastar-Request", "true")
                    // The source-image handler uses this header to return JSON
                    // instead of its normal Datastar SSE response or form redirect.
                    uploadRequest.setRequestHeader("Upload-Progress-Request", "true")
                    uploadRequest.send(formData)
                } catch {
                    finishUpload("Kunne ikke starte opplastingen. Prøv igjen.")
                }
            }
            void sendUpload()
        })
    },
})
