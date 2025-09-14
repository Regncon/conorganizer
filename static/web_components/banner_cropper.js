class BannerCropper extends HTMLElement {
    static get observedAttributes() {
        return ['width', 'height', 'image-url'];
    }

    constructor() {
        super();

        // Defaults (don’t read attributes here)
        this.bannerWidth = 430;
        this.bannerHeight = 180;

        // State
        this.image = new Image();
        this.imageLoaded = false;
        this.scale = 1;
        this.minScale = 1;
        this.drawX = 0;
        this.drawY = 0;
        this.isDragging = false;
        this.dragStartX = 0;
        this.dragStartY = 0;
        this.startDrawX = 0;
        this.startDrawY = 0;

        // Shadow DOM
        const root = this.attachShadow({ mode: 'open' });
        root.innerHTML = `
      <div>
        <div>
          <span id="cameraIcon" style="display:none">
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" width="32" height="32">
              <path fill="currentColor" d="M12 7a5 5 0 1 1 0 10 5 5 0 0 1 0-10Zm0 2a3 3 0 1 0 0 6 3 3 0 0 0 0-6ZM4 4h3.2l1.6-2h6.4l1.6 2H20a2 2 0 0 1 2 2v12a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2Zm0 2v12h16V6h-3.2l-1.6 2H8.8L7.2 6H4Z"/>
            </svg>
          </span>
        </div>

        <div>
          <button id="exportButton" type="button">Export PNG</button>
          <a id="downloadLink"></a>
        </div>

        <div>
          <label for="zoom">Zoom:</label>
          <input id="zoom" type="range" min="1" max="3" step="0.01" value="1" disabled>
        </div>

        <canvas id="canvas" style="cursor:move" aria-label="Banner canvas"></canvas>
      </div>
    `;

        // Elements
        this.canvas = root.getElementById('canvas');
        this.ctx = this.canvas.getContext('2d');
        this.cameraIcon = root.getElementById('cameraIcon');
        this.zoom = root.getElementById('zoom');
        this.exportButton = root.getElementById('exportButton');
        this.downloadLink = root.getElementById('downloadLink');

        // Bind handlers once
        this.handleZoomInput = this.handleZoomInput.bind(this);
        this.onPointerDown = this.onPointerDown.bind(this);
        this.onPointerMove = this.onPointerMove.bind(this);
        this.onPointerUp = this.onPointerUp.bind(this);
        this.handleExport = this.handleExport.bind(this);
    }

    connectedCallback() {
        this._applyInitialAttributes();

        // Listeners
        this.zoom.addEventListener('input', this.handleZoomInput);
        this.canvas.addEventListener('pointerdown', this.onPointerDown);
        window.addEventListener('pointermove', this.onPointerMove);
        window.addEventListener('pointerup', this.onPointerUp);
        this.exportButton.addEventListener('click', this.handleExport);

        this.redraw();
    }

    disconnectedCallback() {
        this.zoom.removeEventListener('input', this.handleZoomInput);
        this.canvas.removeEventListener('pointerdown', this.onPointerDown);
        window.removeEventListener('pointermove', this.onPointerMove);
        window.removeEventListener('pointerup', this.onPointerUp);
        this.exportButton.removeEventListener('click', this.handleExport);
    }

    attributeChangedCallback(name, oldValue, newValue) {
        if (oldValue === newValue) return;

        if (name === 'width' || name === 'height') {
            const w = Number(this.getAttribute('width')) || this.bannerWidth;
            const h = Number(this.getAttribute('height')) || this.bannerHeight;
            this.setCanvasSize(w, h);
        }

        if (name === 'image-url') {
            const url = this.getAttribute('image-url');
            if (url) this._loadImage(url);
        }
    }

    // --- UI handlers ---
    handleZoomInput(e) {
        const newScale = parseFloat(e.target.value);
        this.setScale(newScale);
    }

    onPointerDown(e) {
        if (!this.imageLoaded) return;
        this.isDragging = true;
        this.canvas.setPointerCapture?.(e.pointerId);
        this.dragStartX = e.clientX;
        this.dragStartY = e.clientY;
        this.startDrawX = this.drawX;
        this.startDrawY = this.drawY;
    }

    onPointerMove(e) {
        if (!this.isDragging) return;
        const dx = e.clientX - this.dragStartX;
        const dy = e.clientY - this.dragStartY;
        this.drawX = this.startDrawX + dx;
        this.drawY = this.startDrawY + dy;
        this.redraw();
    }

    onPointerUp(e) {
        this.isDragging = false;
        try { this.canvas.releasePointerCapture?.(e.pointerId); } catch { }
    }

    handleExport() {
        const dataURL = this.canvas.toDataURL('image/png');
        this.downloadLink.href = dataURL;
        this.downloadLink.download = 'banner.png';
        this.downloadLink.textContent = 'Download banner.png';
    }

    // --- Helpers ---
    _applyInitialAttributes() {
        if (this.hasAttribute('width')) {
            const w = Number(this.getAttribute('width'));
            if (!Number.isNaN(w) && w > 0) this.bannerWidth = w;
        }
        if (this.hasAttribute('height')) {
            const h = Number(this.getAttribute('height'));
            if (!Number.isNaN(h) && h > 0) this.bannerHeight = h;
        }
        this.setCanvasSize(this.bannerWidth, this.bannerHeight);

        const url = this.getAttribute('image-url');
        if (url) {
            this._loadImage(url);
        } else {
            this.cameraIcon.style.display = 'block';
            this.zoom.disabled = true;
        }
    }

    setCanvasSize(w, h) {
        this.bannerWidth = w;
        this.bannerHeight = h;
        this.canvas.width = w;
        this.canvas.height = h;
        if (this.imageLoaded) this.setInitialView();
        this.redraw();
    }

    _loadImage(url) {
        this.cameraIcon.style.display = 'none';
        this.image.onload = () => {
            this.imageLoaded = true;
            this.setInitialView();
        };
        this.image.onerror = () => {
            this.imageLoaded = false;
            this.cameraIcon.style.display = 'block';
            this.zoom.disabled = true;
            this.redraw();
        };
        this.image.src = url;
    }

    setInitialView() {
        const coverScaleX = this.canvas.width / this.image.width;
        const coverScaleY = this.canvas.height / this.image.height;
        this.minScale = Math.max(coverScaleX, coverScaleY);
        this.scale = this.minScale;

        this.drawX = (this.canvas.width - this.image.width * this.scale) / 2;
        this.drawY = (this.canvas.height - this.image.height * this.scale) / 2;

        this.zoom.min = this.minScale.toFixed(3);
        this.zoom.max = (this.minScale * 3).toFixed(3);
        this.zoom.step = (this.minScale / 100).toFixed(4);
        this.zoom.value = this.scale.toFixed(3);
        this.zoom.disabled = false;

        this.ctx.imageSmoothingEnabled = true;
        this.ctx.imageSmoothingQuality = 'high';

        this.redraw();
    }

    setScale(newScale) {
        if (!this.imageLoaded) return;
        const oldScale = this.scale;
        const cx = this.canvas.width / 2;
        const cy = this.canvas.height / 2;
        const imgXAtCenter = (cx - this.drawX) / oldScale;
        const imgYAtCenter = (cy - this.drawY) / oldScale;

        this.scale = Math.max(this.minScale, newScale);
        this.drawX = cx - imgXAtCenter * this.scale;
        this.drawY = cy - imgYAtCenter * this.scale;
        this.redraw();
    }

    clampPosition() {
        const maxX = 0;
        const maxY = 0;
        const minX = this.canvas.width - this.image.width * this.scale;
        const minY = this.canvas.height - this.image.height * this.scale;
        this.drawX = Math.min(maxX, Math.max(minX, this.drawX));
        this.drawY = Math.min(maxY, Math.max(minY, this.drawY));
    }

    redraw() {
        this.ctx.clearRect(0, 0, this.canvas.width, this.canvas.height);
        if (!this.imageLoaded) return;
        this.clampPosition();
        this.ctx.drawImage(
            this.image,
            this.drawX,
            this.drawY,
            this.image.width * this.scale,
            this.image.height * this.scale
        );
    }
}

customElements.define('banner-cropper', BannerCropper);
