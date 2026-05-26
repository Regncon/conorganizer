const indicatorPrefix = "toastMessage";
const toastDurationMs = 3000;
const toastExitMs = 180;
let activeHost = null;

function flattenSignals(signals, prefix = "") {
    if (!signals || typeof signals !== "object") {
        return [];
    }

    return Object.entries(signals).flatMap(([key, value]) => {
        const path = prefix ? `${prefix}.${key}` : key;
        if (value && typeof value === "object" && !Array.isArray(value)) {
            return flattenSignals(value, path);
        }
        return [[path, value]];
    });
}

class ToastMessage extends HTMLElement {
    constructor() {
        super();
        this.activeIndicators = new Set();
        this.onSignalPatch = (event) => this.handleSignalPatch(event);
        this.onToast = (event) => this.handleToast(event);
    }

    connectedCallback() {
        if (activeHost) {
            return;
        }

        activeHost = this;
        document.addEventListener("datastar-signal-patch", this.onSignalPatch);
        document.addEventListener("toast", this.onToast);
    }

    disconnectedCallback() {
        if (activeHost !== this) {
            return;
        }

        document.removeEventListener("datastar-signal-patch", this.onSignalPatch);
        document.removeEventListener("toast", this.onToast);
        activeHost = null;
    }

    handleSignalPatch(event) {
        for (const [path, value] of flattenSignals(event.detail)) {
            if (!path.startsWith(indicatorPrefix)) {
                continue;
            }

            if (value === true) {
                this.activeIndicators.add(path);
                continue;
            }

            if (value === false && this.activeIndicators.has(path)) {
                this.activeIndicators.delete(path);
                this.show("Lagret");
            }
        }
    }

    handleToast(event) {
        const message = event.detail?.message;
        if (typeof message !== "string" || message.trim() === "") {
            console.warn("toast event requires detail.message");
            return;
        }

        this.show(message);
    }

    show(message) {
        const template = this.querySelector("template");
        const toast = template?.content.firstElementChild?.cloneNode(true) ?? this.createFallbackToast();
        const text = toast.querySelector("[data-toast-message-text]");
        if (text) {
            text.textContent = message;
        }

        this.append(toast);
        window.setTimeout(() => this.removeToast(toast), toastDurationMs);
    }

    createFallbackToast() {
        const toast = document.createElement("div");
        toast.className = "toast";
        const text = document.createElement("span");
        text.setAttribute("data-toast-message-text", "");
        toast.append(text);
        return toast;
    }

    removeToast(toast) {
        toast.classList.add("is-removing");
        window.setTimeout(() => toast.remove(), toastExitMs);
    }
}

if (!customElements.get("toast-message")) {
    customElements.define("toast-message", ToastMessage);
}
