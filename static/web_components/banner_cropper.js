// @ts-check

import { canvasToWebp } from "../js/image-to-webp.js"

/**
 * Source of truth: [SharedStyles](../js/conorganizer.js#L16) and
 * [ConorganizerGlobal / ConorganizerWindow](../js/conorganizer.js#L42).
 * Keep this local subset in sync with those definitions.
 *
 * @typedef {Window & typeof globalThis & {
 *   conorganizer: {
 *     sharedStyles: {
 *       getStyleUrls: (names?: string[]) => string[],
 *       applyStyleUrlsToShadowRoot: (shadowRoot: ShadowRoot, styleUrls: string[]) => Promise<void>,
 *     },
 *   },
 * }} BannerCropperWindow
 */

/**
 * Payloads dispatched by the cropper and consumed by Datastar bindings.
 * @typedef {{
 *   "crop-started": Record<string, never>,
 *   "crop-ready": { files: FileList },
 *   "crop-error": { message: string },
 *   "image-loading": { url: string },
 *   "image-ready": { url: string },
 *   "image-error": { message: string, url: string },
 * }} BannerCropperEventDetails
 */

// The page loads conorganizer.js before this component.
const typedWindow = /** @type {BannerCropperWindow} */ (window)
const STYLE_URLS = typedWindow.conorganizer.sharedStyles.getStyleUrls()
const UPLOAD_ERROR_MESSAGE = 'Klarte ikkje å lagre endringa. Prøv igjen. Kontakt styret dersom problemet held fram.'
// ---- component --------------------------------------------------------------
class BannerCropper extends HTMLElement {
    static get observedAttributes() {
        return ["width", "height", "preview-width", "preview-height", "image-url", "label"]
    }

    constructor() {
        super()

        // Defaults (don’t read attributes here)
        this.bannerWidth = 430
        this.bannerHeight = 180

        // State
        this.image = new Image()
        this.imageLoaded = false
        this.imageLoadVersion = 0
        this.exportVersion = 0
        this.exporting = false
        this.scale = 1
        this.minScale = 1
        this.drawX = 0
        this.drawY = 0
        this.isDragging = false
        this.dragStartX = 0
        this.dragStartY = 0
        this.startDrawX = 0
        this.startDrawY = 0
        this.dragScaleX = 1
        this.dragScaleY = 1

        // Shadow DOM
        const root = this.attachShadow({ mode: "open" })

        typedWindow.conorganizer.sharedStyles.applyStyleUrlsToShadowRoot(root, STYLE_URLS)

        root.innerHTML = `
        <style>
            .banner-cropper-wrapper {
                --left-side-padding: calc((var(--spacing-2x) * 2));
                --right-side-padding: calc((var(--spacing-4x) * 2));
                --both-sides: calc(var(--left-side-padding) + var(--right-side-padding));

                display: flex;
                flex-direction: column;
                gap: var(--spacing-4x);
                inline-size: min-content;

                .banner-cropper-header-label {
                    color: var(--color-text-soft);
                }

                .banner-cropper-image-slider {
                    display: grid;
                    inline-size: min(calc(100cqi - var(--both-sides)), var(--banner-cropper-preview-width));
                }

                .banner-cropper-canvas {
                    inline-size: 100%;
                    block-size: var(--banner-cropper-preview-height);
                    cursor: move;
                    touch-action: none;
                }

                input[type="range"] {
                    --range-thumb-border: #FFC483;
                    --range-thumb-size: 20px;
                    --range-thumb-background: var(--color-primary);

                    --range-track-size: 12px;
                    --range-track-border: var(--bg-item-border);
                    --range-track-border-size: 1px;
                    --range-track-background: var(--bg-item);

                    --range-progress: 0%;
                    --range-progress-background: var(--color-primary-focus-visible);
                    /* firefox makes color a bit darker */
                    --range-progress-background-chromium: #997759;
                    --range-progress-border: var(--bg-item-hover);

                    --range-focus-ring: var(--bg-item-border-hover);

                    appearance: none;
                    -webkit-appearance: none;
                    inline-size: 100%;
                    block-size: var(--range-track-size);
                    margin-block: 10px;
                    background: transparent;
                    cursor: pointer;
                    margin: 0;

                    &:focus {
                        outline: none;
                    }

                    /* Chrome, Edge and Safari: complete track */
                    &::-webkit-slider-runnable-track {
                        block-size: var(--range-track-size);
                        background: linear-gradient(
                            to right,
                            var(--range-progress-background-chromium) 0 var(--range-progress),
                            var(--range-track-background) var(--range-progress) 100%
                        );
                        outline: 1px solid var(--range-track-border);
                        outline-offset: -1px;
                        border-bottom-right-radius: var(--border-radius-2x);
                        border-bottom-left-radius: var(--border-radius-2x);
                    }

                    /* Chrome, Edge and Safari: thumb */
                    &::-webkit-slider-thumb {
                        appearance: none;
                        -webkit-appearance: none;
                        inline-size: var(--range-thumb-size);
                        block-size: var(--range-thumb-size);
                        margin-top: calc(
                            (var(--range-track-size) - var(--range-thumb-size)) / 2
                        );
                        background: var(--range-thumb-background);
                        border: 1px solid var(--range-thumb-border);
                        border-radius: 50%;
                        cursor: grab;
                        /* take thumb above outline */
                        position: relative;
                        translate: calc(var(--range-progress) - 50%) 0;
                    }

                    &:active::-webkit-slider-thumb {
                        cursor: grabbing;
                    }

                    &:focus-visible::-webkit-slider-runnable-track {
                        box-shadow: 0 0 0 2px var(--range-focus-ring);
                    }

                    /* Firefox: right-side track */
                    &::-moz-range-track {
                        block-size: var(--range-track-size);
                        box-sizing: border-box;
                        background: var(--range-track-background);
                        border: 1px solid var(--range-track-border);
                        border-bottom-right-radius: var(--border-radius-2x);
                        border-bottom-left-radius: var(--border-radius-2x);
                    }

                    /* Firefox: left-side progress */
                    &::-moz-range-progress {
                        block-size: var(--range-track-size);
                        background: var(--range-progress-background);
                        outline: 1px var(--bg-item-hover) solid;
                        outline-offset: -1px;
                        border-bottom-right-radius: var(--border-radius-2x);
                        border-bottom-left-radius: var(--border-radius-2x);
                    }

                    /* Firefox: thumb */
                    &::-moz-range-thumb {
                        inline-size: var(--range-thumb-size);
                        block-size: var(--range-thumb-size);
                        box-sizing: border-box;
                        background: var(--range-thumb-background);
                        border: 1px solid var(--range-thumb-border);
                        border-radius: 50%;
                        cursor: grab;
                        translate: calc(var(--range-progress) - 50%) 0;
                    }

                    &:active::-moz-range-thumb {
                        cursor: grabbing;
                    }

                    &:focus-visible::-moz-range-track {
                        box-shadow: 0 0 0 2px var(--range-focus-ring);
                    }
                }
            }

            @container main (width > 600px) {
                .banner-cropper-wrapper {
                    --both-sides: calc(var(--spacing-10x) * 2)
                }
            }
        </style>
        <div class="banner-cropper-wrapper">
            <h4 class="banner-cropper-header-label">Plassholder tekst</h4>
            <div class="banner-cropper-image-slider">
                <canvas id="canvas" class="banner-cropper-canvas" aria-label="Banner canvas"></canvas>
                <input id="zoom" class="slider" type="range" min="1" max="3" step="0.01" value="1" disabled>
            </div>

            <slot name="actions"></slot>
        </div>
        `

        // Validate the owned template once so handlers can use concrete DOM types.
        const canvas = root.getElementById("canvas")
        const zoom = root.getElementById("zoom")
        const headerLabel = root.querySelector(".banner-cropper-header-label")
        if (!(canvas instanceof HTMLCanvasElement) || !(zoom instanceof HTMLInputElement) || !(headerLabel instanceof HTMLHeadingElement)) {
            throw new Error("Banner cropper template is missing its canvas, zoom input, or heading.")
        }
        const context = canvas.getContext("2d")
        if (!context) {
            throw new Error("Banner cropper requires a 2D canvas context.")
        }
        this.canvas = canvas
        this.ctx = context
        this.zoom = zoom
        this.headerLabelEl = headerLabel

        // Bind handlers once
        this.handleZoomInput = this.handleZoomInput.bind(this)
        this.onPointerDown = this.onPointerDown.bind(this)
        this.onPointerMove = this.onPointerMove.bind(this)
        this.onPointerUp = this.onPointerUp.bind(this)
    }

    connectedCallback() {
        this._applyInitialAttributes()

        // Listeners
        this.zoom.addEventListener("input", this.handleZoomInput)
        this.canvas.addEventListener("pointerdown", this.onPointerDown)
        window.addEventListener("pointermove", this.onPointerMove)
        window.addEventListener("pointerup", this.onPointerUp)
        window.addEventListener("pointercancel", this.onPointerUp)

        this.redraw()
    }

    disconnectedCallback() {
        this.zoom.removeEventListener("input", this.handleZoomInput)
        this.canvas.removeEventListener("pointerdown", this.onPointerDown)
        window.removeEventListener("pointermove", this.onPointerMove)
        window.removeEventListener("pointerup", this.onPointerUp)
        window.removeEventListener("pointercancel", this.onPointerUp)
        this._clearImage()
    }

    /**
     * @param {string} name
     * @param {string | null} oldValue
     * @param {string | null} newValue
     * @returns {void}
     */
    attributeChangedCallback(name, oldValue, newValue) {
        if (oldValue === newValue) return

        if (name === "width" || name === "height") {
            const w = Number(this.getAttribute("width")) || this.bannerWidth
            const h = Number(this.getAttribute("height")) || this.bannerHeight
            this.setCanvasSize(w, h)
        }

        if (name === "preview-width" || name === "preview-height") {
            this._applyPreviewSize()
        }

        if (name === "label" && this.headerLabelEl) {
            this.headerLabelEl.textContent = newValue
        }

        if (name === "image-url" && this.isConnected) {
            if (newValue) {
                this._loadImage(newValue)
            } else {
                this._clearImage()
            }
        }
    }

    // --- UI handlers ---
    /**
     * @param {Event} event
     * @returns {void}
     */
    handleZoomInput(event) {
        const zoomInput = event.currentTarget
        if (!(zoomInput instanceof HTMLInputElement)) return
        const newScale = parseFloat(zoomInput.value)
        this.setScale(newScale)
        this.#updateCssForZoom()
    }

    /**
     * @param {PointerEvent} e
     * @returns {void}
     */
    onPointerDown(e) {
        if (!this.imageLoaded) return
        this.isDragging = true
        this.canvas.setPointerCapture?.(e.pointerId)
        this.dragStartX = e.clientX
        this.dragStartY = e.clientY
        this.startDrawX = this.drawX
        this.startDrawY = this.drawY

        const canvasBounds = this.canvas.getBoundingClientRect()
        this.dragScaleX = canvasBounds.width > 0 ? this.canvas.width / canvasBounds.width : 1
        this.dragScaleY = canvasBounds.height > 0 ? this.canvas.height / canvasBounds.height : 1
    }

    /**
     * @param {PointerEvent} e
     * @returns {void}
     */
    onPointerMove(e) {
        if (!this.isDragging) return
        const dx = (e.clientX - this.dragStartX) * this.dragScaleX
        const dy = (e.clientY - this.dragStartY) * this.dragScaleY
        this.drawX = this.startDrawX + dx
        this.drawY = this.startDrawY + dy
        this.redraw()
    }

    /**
     * @param {PointerEvent} e
     * @returns {void}
     */
    onPointerUp(e) {
        this.isDragging = false
        try {
            this.canvas.releasePointerCapture?.(e.pointerId)
        } catch { }
    }

    /**
     * The page handles crop-started/crop-ready/crop-error and uploads the exported FileList.
     * @returns {Promise<void>}
     */
    async exportImage() {
        if (this.exporting || !this.isConnected) return

        this.exporting = true
        const version = ++this.exportVersion
        const isCurrentExport = () => this.isConnected && version === this.exportVersion
        this._emit("crop-started", {})

        try {
            if (!isCurrentExport()) return
            if (!this.imageLoaded) {
                this._emit("crop-error", { message: "Det er ikke noe bilde å lagre." })
                return
            }

            const blob = await canvasToWebp(this.canvas)
            if (!isCurrentExport()) return

            const transfer = new DataTransfer()
            transfer.items.add(new File([blob], "crop.webp", { type: blob.type }))
            this._emit("crop-ready", { files: transfer.files })
        } catch (error) {
            if (isCurrentExport()) {
                this._emit("crop-error", {
                    message: error instanceof Error ? error.message : "Klarte ikke å klargjøre bildet for lagring. Prøv igjen.",
                })
            }
        } finally {
            if (version === this.exportVersion) this.exporting = false
        }
    }

    // --- Helpers ---
    _applyInitialAttributes() {
        if (this.hasAttribute("width")) {
            const w = Number(this.getAttribute("width"))
            if (!Number.isNaN(w) && w > 0) this.bannerWidth = w
        }
        if (this.hasAttribute("height")) {
            const h = Number(this.getAttribute("height"))
            if (!Number.isNaN(h) && h > 0) this.bannerHeight = h
        }
        this.setCanvasSize(this.bannerWidth, this.bannerHeight)

        const url = this.getAttribute("image-url")
        if (url) {
            this._loadImage(url)
        } else {
            this._clearImage()
        }
    }

    /**
     * @param {number} w
     * @param {number} h
     * @returns {void}
     */
    setCanvasSize(w, h) {
        this.bannerWidth = w
        this.bannerHeight = h
        this.canvas.width = w
        this.canvas.height = h
        this._applyPreviewSize()
        if (this.imageLoaded) this.setInitialView()
        this.redraw()
    }

    _applyPreviewSize() {
        const previewWidth = Number(this.getAttribute("preview-width")) || this.bannerWidth
        const previewHeight = Number(this.getAttribute("preview-height")) || this.bannerHeight

        this.style.setProperty("--banner-cropper-preview-width", `${ previewWidth }px`)
        this.style.setProperty("--banner-cropper-preview-height", `${ previewHeight }px`)
    }

    /**
     * @param {string} url
     * @returns {void}
     */
    _loadImage(url) {
        this._clearImage()
        const version = this.imageLoadVersion
        const image = new Image()
        this.image = image
        const isCurrentImage = () => this.isConnected && version === this.imageLoadVersion

        image.onload = () => {
            if (!isCurrentImage()) return
            this.imageLoaded = true
            this.setInitialView()
            this._emit("image-ready", { url })
        }
        image.onerror = () => {
            if (!isCurrentImage()) return
            this.imageLoaded = false
            this.zoom.disabled = true
            this.redraw()
            this._emit("image-error", { message: "Klarte ikke å laste bildet. Prøv igjen.", url })
        }
        this._emit("image-loading", { url })
        if (isCurrentImage()) image.src = url
    }

    _clearImage() {
        // Retire pending image loads and exports before replacing or disconnecting the image.
        this.imageLoadVersion += 1
        this.exportVersion += 1
        this.exporting = false
        this.imageLoaded = false
        this.image.onload = null
        this.image.onerror = null
        this.image.removeAttribute("src")
        this.isDragging = false
        this.scale = 1
        this.minScale = 1
        this.drawX = 0
        this.drawY = 0
        this.zoom.disabled = true
        this.redraw()
    }

    /**
     * @template {keyof BannerCropperEventDetails} EventType
     * @param {EventType} type
     * @param {BannerCropperEventDetails[EventType]} detail
     * @returns {void}
     */
    _emit(type, detail) {
        this.dispatchEvent(new CustomEvent(type, { bubbles: true, composed: true, detail }))
    }

    setInitialView() {
        const coverScaleX = this.canvas.width / this.image.width
        const coverScaleY = this.canvas.height / this.image.height
        this.minScale = Math.max(coverScaleX, coverScaleY)
        this.scale = this.minScale

        this.drawX = (this.canvas.width - this.image.width * this.scale) / 2
        this.drawY = (this.canvas.height - this.image.height * this.scale) / 2

        this.zoom.min = this.minScale.toFixed(3)
        this.zoom.max = (this.minScale * 3).toFixed(3)
        this.zoom.step = (this.minScale / 100).toFixed(4)
        this.zoom.value = this.scale.toFixed(3)
        this.zoom.disabled = false
        this.#updateCssForZoom()

        this.ctx.imageSmoothingEnabled = true
        this.ctx.imageSmoothingQuality = "high"

        this.redraw()
    }

    #updateCssForZoom() {
        const min = Number(this.zoom.min)
        const max = Number(this.zoom.max)
        const value = this.zoom.valueAsNumber

        const progress = max === min ? 0 : ((value - min) / (max - min)) * 100

        const clampedProgress = Math.min(100, Math.max(0, progress))

        this.zoom.style.setProperty("--range-progress", `${ clampedProgress }%`)
    }

    /**
     * @param {number} newScale
     * @returns {void}
     */
    setScale(newScale) {
        if (!this.imageLoaded) return
        const oldScale = this.scale
        const cx = this.canvas.width / 2
        const cy = this.canvas.height / 2
        const imgXAtCenter = (cx - this.drawX) / oldScale
        const imgYAtCenter = (cy - this.drawY) / oldScale

        this.scale = Math.max(this.minScale, newScale)
        this.drawX = cx - imgXAtCenter * this.scale
        this.drawY = cy - imgYAtCenter * this.scale
        this.redraw()
    }

    clampPosition() {
        const maxX = 0
        const maxY = 0
        const minX = this.canvas.width - this.image.width * this.scale
        const minY = this.canvas.height - this.image.height * this.scale
        this.drawX = Math.min(maxX, Math.max(minX, this.drawX))
        this.drawY = Math.min(maxY, Math.max(minY, this.drawY))
    }

    redraw() {
        this.ctx.clearRect(0, 0, this.canvas.width, this.canvas.height)
        if (!this.imageLoaded) return
        this.clampPosition()
        this.ctx.drawImage(
            this.image,
            this.drawX,
            this.drawY,
            this.image.width * this.scale,
            this.image.height * this.scale,
        )
    }
}

customElements.define("banner-cropper", BannerCropper)
