// SVG supplies geometry; slotted server-rendered rooms own cards and Datastar actions.
class RoomMap extends HTMLElement {
    static observedAttributes = ['rooms']

    #svg
    #controller

    constructor() {
        super()
        this.attachShadow({ mode: 'open' })
        this.shadowRoot.innerHTML = `
            <style>
                :host { display: block; min-inline-size: 0; container-type: inline-size; }
                .scroll { overflow-x: auto; border-radius: 8px; }
                .stage { position: relative; min-inline-size: 1100px; }
                svg { display: block; inline-size: 100%; block-size: auto; background: white; }
                [data-room-target] { pointer-events: none; }
                .map-room { position: absolute; padding: 3px; box-sizing: border-box; }
                .map-room slot { display: block; block-size: 100%; }
                .map-room ::slotted(.room) { block-size: 100%; box-sizing: border-box; }
                .outside { display: grid; grid-template-columns: 1fr; gap: 16px; }
                .outside h3 { grid-column: 1/-1; margin-block: 16px 0; }
                @container (width > 700px) {
                    .outside { grid-template-columns: repeat(2, minmax(0, 1fr)); }
                }
                p { margin-block: 0.5rem; }
            </style>
            <p role="status">Laster kart …</p>
            <div class="scroll" tabindex="0" role="region" aria-label="Romkart – Terminus, 7. etasje">
                <div class="stage"></div>
            </div>
            <div class="outside"></div>`
    }

    connectedCallback() {
        this.#renderRooms()
        if (!this.#svg) this.#load()
    }

    disconnectedCallback() {
        this.#controller?.abort()
    }

    attributeChangedCallback() {
        this.#renderRooms()
    }

    async #load() {
        this.#controller?.abort()
        const controller = new AbortController()
        this.#controller = controller
        try {
            const response = await fetch('/static/rooms/terminus-7-etasje.svg', { signal: controller.signal })
            if (!response.ok) throw new Error('Map unavailable')
            const parsed = new DOMParser().parseFromString(await response.text(), 'image/svg+xml')
            if (controller.signal.aborted) return
            const svg = parsed.documentElement
            if (svg.localName !== 'svg' || parsed.querySelector('parsererror')) throw new Error('Invalid map')
            svg.setAttribute('aria-hidden', 'true')
            this.#svg = svg
            this.shadowRoot.querySelector('.stage').prepend(svg)
            this.#renderRooms()
            this.#status('')
        } catch (error) {
            if (error.name !== 'AbortError') {
                this.#status('Kartet kunne ikke lastes. Bruk romlisten nedenfor, eller last siden på nytt.')
            }
        }
    }

    #renderRooms() {
        const rooms = JSON.parse(this.getAttribute('rooms') || '[]')
        const stage = this.shadowRoot.querySelector('.stage')
        const outside = this.shadowRoot.querySelector('.outside')
        stage.querySelectorAll('.map-room').forEach(room => room.remove())
        outside.replaceChildren()
        const targets = new Map(Array.from(this.#svg?.querySelectorAll('[data-room-target]') || [],
            target => [target.dataset.roomTarget, target]))
        const viewBox = this.#svg?.viewBox.baseVal
        let hasOutside = false

        for (const room of rooms) {
            const slot = document.createElement('slot')
            slot.name = `room-${room.ID}`
            const target = targets.get(room.RoomNumber)
            if (target) {
                const overlay = document.createElement('div')
                overlay.className = 'map-room'
                overlay.style.left = `${(target.x.baseVal.value - viewBox.x) / viewBox.width * 100}%`
                overlay.style.top = `${(target.y.baseVal.value - viewBox.y) / viewBox.height * 100}%`
                overlay.style.width = `${target.width.baseVal.value / viewBox.width * 100}%`
                overlay.style.height = `${target.height.baseVal.value / viewBox.height * 100}%`
                overlay.append(slot)
                stage.append(overlay)
            } else {
                if (!hasOutside) {
                    const heading = document.createElement('h3')
                    heading.textContent = this.#svg ? 'Rom utenfor kartet' : 'Rom'
                    outside.append(heading)
                    hasOutside = true
                }
                outside.append(slot)
            }
        }
        this.shadowRoot.querySelector('.scroll').hidden = !this.#svg
    }

    #status(message) {
        const status = this.shadowRoot.querySelector('[role="status"]')
        status.textContent = message
        status.hidden = !message
    }
}

if (!customElements.get('room-map')) customElements.define('room-map', RoomMap)
