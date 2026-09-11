import { action, mergePatch, mergePaths } from "../datastar.js"

// Fetch does not expose upload progress. Keep XHR at the transport boundary,
// while using Datastar signals and request events for all page state.
action({
    name: "upload",
    apply({ el, cleanups }, url, { progressSignal } = {}) {
        const form = el.closest("form")
        if (!form || (!form.noValidate && !form.reportValidity())) return

        const cleanupKey = "@upload"
        if (cleanups.has(cleanupKey)) return

        // Collect files before data-indicator disables the file input.
        const body = new FormData(form)
        const request = new XMLHttpRequest()
        const dispatch = (type, argsRaw = {}) => document.dispatchEvent(
            new CustomEvent("datastar-fetch", { detail: { type, el, argsRaw } }),
        )
        const progress = (value) => {
            if (progressSignal) mergePaths([[progressSignal, value]])
        }

        return new Promise((resolve) => {
            let finished = false
            const finish = (message) => {
                if (finished) return
                finished = true
                if (message) dispatch("error", { message, status: String(request.status) })
                cleanups.delete(cleanupKey)
                dispatch("finished")
                resolve()
            }

            cleanups.set(cleanupKey, () => request.abort())
            request.upload.addEventListener("progress", (event) => {
                if (event.lengthComputable) {
                    progress(Math.min(100, Math.round(event.loaded / event.total * 100)))
                }
            })
            request.addEventListener("load", () => {
                if (request.status < 200 || request.status >= 300) {
                    finish(request.responseText || "Kunne ikke laste opp bildet. Prøv igjen.")
                    return
                }
                try {
                    const signals = JSON.parse(request.responseText)
                    if (!signals || typeof signals.sourceImageUrl !== "string" || !signals.sourceImageUrl) {
                        throw new Error("Missing sourceImageUrl")
                    }
                    progress(100)
                    mergePatch(signals)
                    finish()
                } catch {
                    finish("Kunne ikke lese svaret fra opplastingen. Prøv igjen.")
                }
            })
            request.addEventListener("error", () => finish("Kunne ikke laste opp bildet. Prøv igjen."))
            request.addEventListener("abort", () => finish("Opplastingen ble avbrutt."))
            request.addEventListener("timeout", () => finish("Opplastingen tok for lang tid. Prøv igjen."))

            progress(0)
            dispatch("started")
            try {
                request.open("POST", url)
                request.setRequestHeader("Accept", "application/json")
                request.setRequestHeader("Datastar-Request", "true")
                request.setRequestHeader("Upload-Progress-Request", "true")
                request.send(body)
            } catch {
                finish("Kunne ikke starte opplastingen. Prøv igjen.")
            }
        })
    },
})
