class BannerCropper extends HTMLElement {
  constructor() {
    super();

    // Read desired banner size from attributes (fallback to defaults)
    this.bannerWidth = parseInt(this.getAttribute('width')) || 1200;
    this.bannerHeight = parseInt(this.getAttribute('height')) || 400;

    // Shadow DOM with simple markup (no styling)
    this.shadow = this.attachShadow({ mode: 'open' });
    this.shadow.innerHTML = `
      <div>
        <div>
          <input id="fileInput" type="file" accept="image/*">
          <button id="loadButton">Load Image</button>
        </div>
        <div>
          <label for="zoom">Zoom:</label>
          <input id="zoom" type="range" min="1" max="3" step="0.01" value="1" disabled>
        </div>
        <canvas id="canvas" width="${this.bannerWidth}" height="${this.bannerHeight}" aria-label="Banner canvas"></canvas>
        <div>
          <button id="exportButton">Export PNG</button>
          <a id="downloadLink"></a>
        </div>
      </div>
    `;

    // Elements
    this.canvas = this.shadow.getElementById('canvas');
    this.ctx = this.canvas.getContext('2d');
    this.fileInput = this.shadow.getElementById('fileInput');
    this.loadButton = this.shadow.getElementById('loadButton');
    this.zoom = this.shadow.getElementById('zoom');
    this.exportButton = this.shadow.getElementById('exportButton');
    this.downloadLink = this.shadow.getElementById('downloadLink');

    // Image state
    this.image = new Image();
    this.imageLoaded = false;
    this.scale = 1;
    this.minScale = 1;
    this.drawX = 0;
    this.drawY = 0;

    // Drag state
    this.isDragging = false;
    this.dragStartX = 0;
    this.dragStartY = 0;
    this.startDrawX = 0;
    this.startDrawY = 0;

    // Bind handlers
    this.handleLoadClick = this.handleLoadClick.bind(this);
    this.handleZoomInput = this.handleZoomInput.bind(this);
    this.onPointerDown = this.onPointerDown.bind(this);
    this.onPointerMove = this.onPointerMove.bind(this);
    this.onPointerUp = this.onPointerUp.bind(this);
    this.handleExport = this.handleExport.bind(this);
  }

  connectedCallback() {
    this.loadButton.addEventListener('click', this.handleLoadClick);
    this.zoom.addEventListener('input', this.handleZoomInput);
    this.canvas.addEventListener('pointerdown', this.onPointerDown);
    window.addEventListener('pointermove', this.onPointerMove);
    window.addEventListener('pointerup', this.onPointerUp);
    this.exportButton.addEventListener('click', this.handleExport);
    this.redraw();
  }

  disconnectedCallback() {
    this.loadButton.removeEventListener('click', this.handleLoadClick);
    this.zoom.removeEventListener('input', this.handleZoomInput);
    this.canvas.removeEventListener('pointerdown', this.onPointerDown);
    window.removeEventListener('pointermove', this.onPointerMove);
    window.removeEventListener('pointerup', this.onPointerUp);
    this.exportButton.removeEventListener('click', this.handleExport);
  }

  // --- UI handlers ---
  handleLoadClick() {
    const files = this.fileInput.files;
    if (!files || files.length === 0) return;
    const file = files[0];
    const reader = new FileReader();
    reader.onload = (e) => {
      this.image.onload = () => {
        this.imageLoaded = true;
        this.setInitialView();
      };
      this.image.src = e.target.result;
    };
    reader.readAsDataURL(file);
  }

  handleZoomInput(e) {
    const newScale = parseFloat(e.target.value);
    this.setScale(newScale);
  }

  onPointerDown(e) {
    if (!this.imageLoaded) return;
    this.isDragging = true;
    this.canvas.setPointerCapture(e.pointerId);
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
  }

  handleExport() {
    const dataURL = this.canvas.toDataURL('image/png');
    this.downloadLink.href = dataURL;
    this.downloadLink.download = 'banner.png';
    this.downloadLink.textContent = 'Download banner.png';
  }

  // --- View helpers ---
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
