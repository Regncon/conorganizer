// @ts-check

/**
 * Encode the full source image using the same browser API and quality as
 * banner_cropper.js. Cropping happens later, after the source has been saved.
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

        /** @type {Blob | null} */
        const webpBlob = await new Promise((resolve) => canvas.toBlob(resolve, "image/webp", 0.9))
        // Browsers may return null or fall back to PNG if WebP encoding fails.
        if (!webpBlob || webpBlob.type !== "image/webp") {
            throw new Error("Nettleseren kunne ikke lage et WebP-bilde.")
        }

        const filename = sourceFile.name.replace(/\.[^.]*$/, "") + ".webp"
        return new File([webpBlob], filename, { type: webpBlob.type })
    } finally {
        URL.revokeObjectURL(sourceImageUrl)
        sourceImage.src = ""
        canvas.width = 0
        canvas.height = 0
    }
}
