// @ts-check


if (!customElements.get("ticket-holder-dropdown")) {
    const GLOBAL_STYLE_URLS = [
        "/static/index.css",
        "/static/buttons.css",
        "/static/web_components/ticket_holder_dropdown.css",
    ]

    const globalSheetsPromise = (async () => {
        // Feature check for constructable stylesheets
        const supportsConstructable =
            !!(Document.prototype && "adoptedStyleSheets" in Document.prototype) &&
            !!(CSSStyleSheet.prototype && "replace" in CSSStyleSheet.prototype)

        if (!supportsConstructable) return null

        const sheets = []
        for (const url of GLOBAL_STYLE_URLS) {
            const resp = await fetch(url, { credentials: "same-origin" })
            const cssText = await resp.text()
            const sheet = new CSSStyleSheet()
            await sheet.replace(cssText)
            sheets.push(sheet)
        }
        return sheets
    })()

    /**
     * Type for billettholder objects expected in the input JSON array.
     * @typedef {Object} BillettHolder
     * @property {number} Id
     * @property {string} Name
     * @property {string} Email
     * @property {string} Color
     */
    const DATA_BILLETTHOLDERE_ATTR = "data-billettholdere"
    /**
     * Ticket-holder dropdown custom element.
     *
     * Required input:
     * - `data-billettholdere`: JSON array of ticket holders.
     *
     * Optional templ-provided icon:
     * - Provide a child `<template data-arrow-icon>...</template>`.
     * - The component clones this template into the arrow slot during render.
     * - If omitted, it falls back to a plain text arrow.
     */
    class TicketHolderDropdown extends HTMLElement {
        static get observedAttributes() {
            return [DATA_BILLETTHOLDERE_ATTR]
        }

        constructor() {
            super()
            if (!this.shadowRoot) {
                this.attachShadow({ mode: "open" })
                globalSheetsPromise.then((sheets) => {
                    if (sheets && this.shadowRoot) {
                        this.shadowRoot.adoptedStyleSheets = [...this.shadowRoot.adoptedStyleSheets, ...sheets]
                    } else if (this.shadowRoot) {
                        // Fallback: inject <link> tags
                        for (const url of GLOBAL_STYLE_URLS) {
                            const link = document.createElement("link")
                            link.rel = "stylesheet"
                            link.href = url
                            this.shadowRoot.appendChild(link)
                        }
                    }
                })
            }
            /** @type {HTMLButtonElement | null} */
            this.selectButtonEle = null
            /** @type {HTMLUListElement | null} */
            this.dropdownEle = null
            /** @type {HTMLSpanElement | null} */
            this.selectedValueEle = null
            /** @type {number} */
            this.focusedIndex = -1
            /** @type {string} */
            this.LSKey = "selectedBillettHolder"
            /** @type {BillettHolder[]} */
            this.billettholdere = []
            /** @type {HTMLTemplateElement | null} */
            this.arrowIconTemplateEle = null

            this.onButtonClick = this.onButtonClick.bind(this)
            this.onButtonKeydown = this.onButtonKeydown.bind(this)
            this.onDropdownKeydown = this.onDropdownKeydown.bind(this)
            this.onDropdownClick = this.onDropdownClick.bind(this)
            this.onDocumentClick = this.onDocumentClick.bind(this)
        }

        /**
         * @param {string} name
         * @param {string | null} oldValue
         * @param {string | null} newValue
         * @returns {void}
         */
        attributeChangedCallback(name, oldValue, newValue) {
            if (name !== DATA_BILLETTHOLDERE_ATTR || !this.isConnected) {
                return
            }
            this.syncFromAttribute()
        }

        connectedCallback() {
            this.ensureShadowStyles()
            this.syncFromAttribute()
        }

        disconnectedCallback() {
            this.teardownInteractiveElements()
        }

        /**
         * @returns {void}
         */
        syncFromAttribute() {
            this.billettholdere = this.parseHolders()
            if (this.billettholdere.length === 0) {
                this.teardownInteractiveElements()
                this.shadowRoot?.replaceChildren()
                this.ensureShadowStyles()
                return
            }

            if (!this.arrowIconTemplateEle) {
                this.arrowIconTemplateEle = this.querySelector("template[data-arrow-icon]")
            }

            this.render()
            if (!this.setupInteractiveElements()) {
                return
            }
            this.hydrateSelection()
        }

        /**
         * @returns {void}
         */
        ensureShadowStyles() {
            // if (!this.shadowRoot || this.shadowRoot.getElementById(STYLE_ID)) {
            //     return
            // }
            // const styleEle = document.createElement("style")
            // styleEle.id = STYLE_ID
            // styleEle.textContent = STYLE_TEXT
            // this.shadowRoot.appendChild(styleEle)
        }

        /**
         * @returns {boolean}
         */
        setupInteractiveElements() {
            this.teardownInteractiveElements()

            this.selectButtonEle = this.shadowRoot?.querySelector(".select-button") || null
            this.dropdownEle = this.shadowRoot?.querySelector(".dropdown-list") || null
            this.selectedValueEle = this.shadowRoot?.querySelector(".selected-value") || null
            if (!this.selectButtonEle || !this.dropdownEle || !this.selectedValueEle) {
                return false
            }

            const controlId = `dropdown-list-${ Math.random().toString(36).slice(2, 10) }`
            const buttonId = `dropdown-button-${ Math.random().toString(36).slice(2, 10) }`
            this.dropdownEle.id = controlId
            this.selectButtonEle.id = buttonId
            this.selectButtonEle.setAttribute("aria-controls", controlId)
            this.dropdownEle.setAttribute("aria-labelledby", buttonId)

            this.selectButtonEle.addEventListener("click", this.onButtonClick)
            this.selectButtonEle.addEventListener("keydown", this.onButtonKeydown)
            this.dropdownEle.addEventListener("keydown", this.onDropdownKeydown)
            this.dropdownEle.addEventListener("click", this.onDropdownClick)
            document.addEventListener("click", this.onDocumentClick)
            return true
        }

        /**
         * @returns {void}
         */
        teardownInteractiveElements() {
            if (this.selectButtonEle) {
                this.selectButtonEle.removeEventListener("click", this.onButtonClick)
                this.selectButtonEle.removeEventListener("keydown", this.onButtonKeydown)
            }
            if (this.dropdownEle) {
                this.dropdownEle.removeEventListener("keydown", this.onDropdownKeydown)
                this.dropdownEle.removeEventListener("click", this.onDropdownClick)
            }
            document.removeEventListener("click", this.onDocumentClick)
        }

        /**
         * @returns {BillettHolder[]}
         */
        parseHolders() {
            const raw = this.getAttribute(DATA_BILLETTHOLDERE_ATTR)
            if (!raw) {
                return []
            }
            try {
                const parsed = JSON.parse(raw)
                if (!Array.isArray(parsed)) {
                    return []
                }
                return parsed.map((item) => ({
                    Id: Number(item?.Id || 0),
                    Name: String(item?.Name || ""),
                    Email: String(item?.Email || ""),
                    Color: String(item?.Color || ""),
                }))
            } catch {
                return []
            }
        }

        /**
         * @param {string} name
         * @returns {string}
         */
        getInitials(name) {
            return name
                .split(" ")
                .map((n) => n.trim())
                .filter((n) => n.length > 0)
                .map((n) => n[0])
                .join("")
                .toUpperCase()
        }

        /**
         * @param {BillettHolder} holder
         * @returns {HTMLDivElement}
         */
        createNameInitialsNode(holder) {
            const wrapperEle = document.createElement("div")
            wrapperEle.className = "name-initials"

            const initialsEle = document.createElement("span")
            initialsEle.className = "initials"
            if (holder.Color) {
                initialsEle.style.backgroundColor = `hsl(from ${ holder.Color } h s l / 0.5)`
                initialsEle.style.border = `1px solid ${ holder.Color }`
            }
            initialsEle.textContent = this.getInitials(holder.Name)

            const nameEle = document.createElement("p")
            nameEle.className = "name"
            nameEle.textContent = holder.Name

            wrapperEle.appendChild(initialsEle)
            wrapperEle.appendChild(nameEle)
            return wrapperEle
        }

        /**
         * @returns {void}
         */
        render() {
            if (!this.shadowRoot) {
                return
            }
            const wrapperEle = document.createElement("div")
            wrapperEle.className = "ticket-holder-dropdown-wrapper"

            const buttonEle = document.createElement("button")
            buttonEle.className = "select-button input no-marking"
            buttonEle.setAttribute("role", "combobox")
            buttonEle.setAttribute("aria-label", "select button")
            buttonEle.setAttribute("aria-haspopup", "listbox")
            buttonEle.setAttribute("aria-expanded", "false")
            buttonEle.type = "button"

            const selectedValueEle = document.createElement("span")
            selectedValueEle.className = "selected-value"

            const buttonEndEle = document.createElement("div")
            buttonEndEle.className = "select-button-end"
            const arrowEle = document.createElement("i")
            arrowEle.className = "arrow"
            arrowEle.setAttribute("aria-hidden", "true")
            if (this.arrowIconTemplateEle) {
                arrowEle.appendChild(this.arrowIconTemplateEle.content.cloneNode(true))
            } else {
                arrowEle.textContent = "▾"
            }
            buttonEndEle.appendChild(arrowEle)

            buttonEle.appendChild(selectedValueEle)
            buttonEle.appendChild(buttonEndEle)

            const listEle = document.createElement("ul")
            listEle.className = "dropdown-list hidden"
            listEle.setAttribute("role", "listbox")

            this.billettholdere.forEach((holder) => {
                const liEle = document.createElement("li")
                liEle.setAttribute("role", "option")
                liEle.dataset.billettHolderId = String(holder.Id)
                liEle.dataset.billettHolderName = holder.Name
                liEle.dataset.billettHolderEmail = holder.Email
                liEle.dataset.billettHolderColor = holder.Color
                liEle.setAttribute("data-bind", "billettHolderId")
                liEle.setAttribute("data-on:click", `$billettHolderId = ${ holder.Id }`)
                liEle.appendChild(this.createNameInitialsNode(holder))
                listEle.appendChild(liEle)
            })

            wrapperEle.appendChild(buttonEle)
            wrapperEle.appendChild(listEle)
            this.shadowRoot.replaceChildren()
            this.ensureShadowStyles()
            this.shadowRoot.appendChild(wrapperEle)
        }

        /**
         * @returns {HTMLLIElement[]}
         */
        getOptions() {
            return Array.from(this.shadowRoot?.querySelectorAll("li") || [])
        }

        /**
         * @param {HTMLLIElement} optionEle
         * @returns {BillettHolder}
         */
        toBillettHolder(optionEle) {
            return {
                Id: Number(optionEle.dataset.billettHolderId || "0"),
                Name: optionEle.dataset.billettHolderName || "",
                Email: optionEle.dataset.billettHolderEmail || "",
                Color: optionEle.dataset.billettHolderColor || "",
            }
        }

        /**
         * @param {HTMLLIElement} optionEle
         * @returns {void}
         */
        saveSelected(optionEle) {
            localStorage.setItem(this.LSKey, JSON.stringify(this.toBillettHolder(optionEle)))
        }

        /**
         * @param {HTMLLIElement[]} optionEles
         * @returns {void}
         */
        updateFocus(optionEles) {
            optionEles.forEach((optionEle, index) => {
                optionEle.setAttribute("tabindex", index === this.focusedIndex ? "0" : "-1")
                if (index === this.focusedIndex) {
                    optionEle.focus()
                }
            })
        }

        /**
         * @param {HTMLLIElement} optionEle
         * @returns {void}
         */
        renderSelected(optionEle) {
            if (!this.selectedValueEle) {
                return
            }
            this.getOptions().forEach((opt) => opt.classList.remove("selected"))
            optionEle.classList.add("selected")

            const holder = this.toBillettHolder(optionEle)
            this.selectedValueEle.replaceChildren(this.createNameInitialsNode(holder))
        }

        /**
         * @param {boolean | null} [expand]
         * @returns {void}
         */
        toggleDropdown(expand = null) {
            if (!this.dropdownEle || !this.selectButtonEle) {
                return
            }

            const optionEles = this.getOptions()
            const isOpen = expand !== null ? expand : this.dropdownEle.classList.contains("hidden")
            this.dropdownEle.classList.toggle("hidden", !isOpen)
            this.selectButtonEle.setAttribute("aria-expanded", String(isOpen))

            if (isOpen) {
                this.focusedIndex = optionEles.findIndex((optionEle) => optionEle.classList.contains("selected"))
                this.focusedIndex = this.focusedIndex === -1 ? 0 : this.focusedIndex
                this.updateFocus(optionEles)
                return
            }

            this.focusedIndex = -1
            this.selectButtonEle.focus()
        }

        /**
         * @returns {void}
         */
        hydrateSelection() {
            const optionEles = this.getOptions()
            const firstOptionEle = optionEles[0]
            if (!firstOptionEle) {
                return
            }

            const selectedBillettholderLS = localStorage.getItem(this.LSKey)
            if (!selectedBillettholderLS) {
                this.renderSelected(firstOptionEle)
                this.saveSelected(firstOptionEle)
                return
            }

            try {
                /** @type {BillettHolder} */
                const selectedBillettholder = JSON.parse(selectedBillettholderLS)
                const selectedOptionEle = optionEles.find(
                    (optionEle) => Number(optionEle.dataset.billettHolderId || "0") === Number(selectedBillettholder.Id),
                )
                if (!selectedOptionEle) {
                    this.renderSelected(firstOptionEle)
                    this.saveSelected(firstOptionEle)
                    return
                }
                this.renderSelected(selectedOptionEle)
            } catch {
                this.renderSelected(firstOptionEle)
                this.saveSelected(firstOptionEle)
            }
        }

        /**
         * @param {HTMLLIElement} optionEle
         * @returns {void}
         */
        handleOptionSelect(optionEle) {
            this.renderSelected(optionEle)
            this.saveSelected(optionEle)
        }

        /**
         * @returns {void}
         */
        onButtonClick() {
            this.toggleDropdown()
        }

        /**
         * @param {KeyboardEvent} event
         * @returns {void}
         */
        onButtonKeydown(event) {
            if (event.key === "ArrowDown") {
                event.preventDefault()
                this.toggleDropdown(true)
                return
            }
            if (event.key === "Escape") {
                this.toggleDropdown(false)
            }
        }

        /**
         * @param {KeyboardEvent} event
         * @returns {void}
         */
        onDropdownKeydown(event) {
            const optionEles = this.getOptions()
            if (optionEles.length === 0) {
                return
            }

            if (event.key === "ArrowDown") {
                event.preventDefault()
                this.focusedIndex = (this.focusedIndex + 1) % optionEles.length
                this.updateFocus(optionEles)
                return
            }
            if (event.key === "ArrowUp") {
                event.preventDefault()
                this.focusedIndex = (this.focusedIndex - 1 + optionEles.length) % optionEles.length
                this.updateFocus(optionEles)
                return
            }
            if (event.key === "Enter" || event.key === " ") {
                event.preventDefault()
                const optionEle = optionEles[this.focusedIndex]
                if (!optionEle) {
                    return
                }
                this.handleOptionSelect(optionEle)
                optionEle.click()
                this.toggleDropdown(false)
                return
            }
            if (event.key === "Escape") {
                this.toggleDropdown(false)
            }
        }

        /**
         * @param {MouseEvent} event
         * @returns {void}
         */
        onDropdownClick(event) {
            const targetEle = event.target
            if (!(targetEle instanceof Element)) {
                return
            }

            const optionEle = targetEle.closest("li")
            if (!(optionEle instanceof HTMLLIElement)) {
                return
            }

            this.handleOptionSelect(optionEle)
            this.toggleDropdown(false)
        }

        /**
         * @param {MouseEvent} event
         * @returns {void}
         */
        onDocumentClick(event) {
            const targetEle = event.target
            const isOutsideClick = !(targetEle instanceof Node) || !this.contains(targetEle)
            if (isOutsideClick) {
                this.toggleDropdown(false)
            }
        }
    }

    customElements.define("ticket-holder-dropdown", TicketHolderDropdown)
}
