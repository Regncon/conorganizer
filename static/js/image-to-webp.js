// @ts-check

/**
 * Shared WebP encoder for source uploads and cropped images.
 * Keeps encoding quality and output validation in one place.
 *
 * @param {HTMLCanvasElement} canvas
 * @returns {Promise<Blob>} Rejects if the browser cannot produce WebP.
 */
export async function canvasToWebp(canvas) {
    /** @type {Blob | null} */
    const webpBlob = await new Promise((resolve) => canvas.toBlob(resolve, "image/webp", 0.9))
    // Browsers may return null or fall back to PNG if WebP encoding fails.
    if (!webpBlob || webpBlob.type !== "image/webp") {
        throw new Error("Nettleseren kunne ikke lage et WebP-bilde.")
    }
    return webpBlob
}

/**
 * Load the full source image onto a canvas and encode it as a WebP file.
 * The cropper uses canvasToWebp directly because it already has a canvas.
 *
 * @param {File} sourceFile
 * @returns {Promise<File>}
 */
export async function imageToWebp(sourceFile) {
    const sourceImageUrl = URL.createObjectURL(sourceFile)
    const sourceImage = new Image()
    const canvas = document.createElement("canvas")
    try {
        // Use the browser's image decoder, including orientation and SVG support.
        sourceImage.src = sourceImageUrl
        await sourceImage.decode()
        canvas.width = sourceImage.naturalWidth
        canvas.height = sourceImage.naturalHeight
        const context = canvas.getContext("2d")
        if (!context) throw new Error("Nettleseren kunne ikke klargjøre bildet.")
        context.drawImage(sourceImage, 0, 0)

        const webpBlob = await canvasToWebp(canvas)
        const filename = sourceFile.name.replace(/\.[^.]*$/, "") + ".webp"
        return new File([webpBlob], filename, { type: webpBlob.type })
    } finally {
        URL.revokeObjectURL(sourceImageUrl)
        sourceImage.src = ""
        canvas.width = 0
        canvas.height = 0
    }
}
