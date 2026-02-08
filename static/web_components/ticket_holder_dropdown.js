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
     * @typedef {Object} Billettholder
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
            /** @type {Billettholder[]} */
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
         * Runs when a watched attribute changes.
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
            this.syncFromAttribute()
        }

        disconnectedCallback() {
            this.teardownInteractiveElements()
        }

        /**
         * Reads input data, renders dropdown, and restores selection.
         * @returns {void}
         */
        syncFromAttribute() {
            this.billettholdere = this.parseBillettholdere()
            if (this.billettholdere.length === 0) {
                this.teardownInteractiveElements()
                this.shadowRoot?.replaceChildren()
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
         * Finds key DOM nodes and attaches event listeners.
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
         * Removes previously attached event listeners.
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
         * Parses `data-billettholdere` JSON into a normalized list.
         * @returns {Billettholder[]}
         */
        parseBillettholdere() {
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
         * Builds initials from a full name.
         * @param {string} name
         * @returns {string}
         */
        getInitials(name) {
            const parts = name
                .split(" ")
                .map((part) => part.trim())
                .filter((part) => part.length > 0)

            if (parts.length === 0) {
                return ""
            }

            const firstName = parts[0]
            const lastName = parts[parts.length - 1]
            const firstInitial = firstName[0] || ""
            const lastInitial = lastName[0] || ""

            if (parts.length === 1) {
                return firstInitial.toUpperCase()
            }

            return `${ firstInitial }${ lastInitial }`.toUpperCase()
        }

        /**
         * Creates the visual row for a billettholder.
         * @param {Billettholder} billettholder
         * @returns {HTMLDivElement}
         */
        createNameInitialsNode(billettholder) {
            const wrapperEle = document.createElement("div")
            wrapperEle.className = "name-initials"

            const initialsEle = document.createElement("span")
            initialsEle.className = "initials"
            if (billettholder.Color) {
                initialsEle.style.backgroundColor = `hsl(from ${ billettholder.Color } h s l / 0.5)`
                initialsEle.style.border = `1px solid ${ billettholder.Color }`
            }
            initialsEle.textContent = this.getInitials(billettholder.Name)

            const nameEle = document.createElement("p")
            nameEle.className = "name"
            nameEle.textContent = billettholder.Name

            wrapperEle.appendChild(initialsEle)
            wrapperEle.appendChild(nameEle)
            return wrapperEle
        }

        /**
         * Renders button + dropdown options inside shadow DOM.
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

            this.billettholdere.forEach((billettholder) => {
                const liEle = document.createElement("li")
                liEle.setAttribute("role", "option")
                liEle.dataset.Id = String(billettholder.Id)
                liEle.dataset.Name = billettholder.Name
                liEle.dataset.Email = billettholder.Email
                liEle.dataset.Color = billettholder.Color
                liEle.onclick = () => this.emitBillettholderSelected(billettholder.Id)
                liEle.appendChild(this.createNameInitialsNode(billettholder))
                listEle.appendChild(liEle)
            })

            wrapperEle.appendChild(buttonEle)
            wrapperEle.appendChild(listEle)
            this.shadowRoot.replaceChildren()
            this.shadowRoot.appendChild(wrapperEle)
        }

        /**
         * Returns all list option elements.
         * @returns {HTMLLIElement[]}
         */
        getOptionElements() {
            return Array.from(this.shadowRoot?.querySelectorAll("li") || [])
        }

        /**
         * Converts a list element dataset into a billettholder object.
         * @param {HTMLLIElement} optionEle
         * @returns {Billettholder}
         */
        toBillettHolder(optionEle) {
            return {
                Id: Number(optionEle.dataset.Id || "0"),
                Name: optionEle.dataset.Name || "",
                Email: optionEle.dataset.Email || "",
                Color: optionEle.dataset.Color || "",
            }
        }

        /**
         * Saves selected billettholder to localStorage.
         * @param {HTMLLIElement} optionEle
         * @returns {void}
         */
        saveSelectedToLocalStorage(optionEle) {
            localStorage.setItem(this.LSKey, JSON.stringify(this.toBillettHolder(optionEle)))
        }

        /**
         * Updates keyboard focus on dropdown options.
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
         * Updates selected UI state and button content.
         * @param {HTMLLIElement} optionEle
         * @returns {void}
         */
        renderSelected(optionEle) {
            if (!this.selectedValueEle) {
                return
            }
            this.getOptionElements().forEach((opt) => opt.classList.remove("selected"))
            optionEle.classList.add("selected")

            const billettholder = this.toBillettHolder(optionEle)
            this.selectedValueEle.replaceChildren(this.createNameInitialsNode(billettholder))
        }

        /**
         * Opens or closes the dropdown and manages focus.
         * @param {boolean | null} [expand]
         * @returns {void}
         */
        toggleDropdown(expand = null) {
            if (!this.dropdownEle || !this.selectButtonEle) {
                return
            }

            const optionEles = this.getOptionElements()
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
         * Restores selection from localStorage or picks first option.
         * @returns {void}
         */
        hydrateSelection() {
            const optionEles = this.getOptionElements()
            const firstOptionEle = optionEles[0]
            if (!firstOptionEle) {
                return
            }

            const selectedBillettholderLS = localStorage.getItem(this.LSKey)
            if (!selectedBillettholderLS) {
                this.renderSelected(firstOptionEle)
                this.saveSelectedToLocalStorage(firstOptionEle)
                return
            }

            try {
                /** @type {Billettholder} */
                const selectedBillettholder = JSON.parse(selectedBillettholderLS)
                const selectedOptionEle = optionEles.find(
                    (optionEle) => Number(optionEle.dataset.Id || "0") === Number(selectedBillettholder.Id),
                )
                if (!selectedOptionEle) {
                    this.renderSelected(firstOptionEle)
                    this.saveSelectedToLocalStorage(firstOptionEle)
                    return
                }
                this.renderSelected(selectedOptionEle)
            } catch {
                this.renderSelected(firstOptionEle)
                this.saveSelectedToLocalStorage(firstOptionEle)
            }
        }

        /**
         * Handles option selection side effects.
         * @param {HTMLLIElement} optionEle
         * @returns {void}
         */
        handleOptionSelect(optionEle) {
            this.renderSelected(optionEle)
            this.saveSelectedToLocalStorage(optionEle)
        }

        /**
         * Toggles dropdown on button click.
         * @returns {void}
         */
        onButtonClick() {
            this.toggleDropdown()
        }

        /**
         * Handles keyboard input while focus is on button.
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
         * Handles keyboard navigation and selection in the list.
         * @param {KeyboardEvent} event
         * @returns {void}
         */
        onDropdownKeydown(event) {
            const optionEles = this.getOptionElements()
            if (optionEles.length === 0) {
                return
            }

            switch (event.key) {
                case "ArrowDown":
                    event.preventDefault()
                    this.focusedIndex = (this.focusedIndex + 1) % optionEles.length
                    this.updateFocus(optionEles)
                    return
                case "ArrowUp":
                    event.preventDefault()
                    this.focusedIndex = (this.focusedIndex - 1 + optionEles.length) % optionEles.length
                    this.updateFocus(optionEles)
                    return
                case "Enter":
                case " ":
                    event.preventDefault()
                    {
                        const optionEle = optionEles[this.focusedIndex]
                        if (!optionEle) {
                            return
                        }

                        this.handleOptionSelect(optionEle)
                        this.emitBillettholderSelected(this.toBillettHolder(optionEle).Id)
                        this.toggleDropdown(false)
                    }
                    return
                case "Escape":
                    this.toggleDropdown(false)
                    return
                default:
                    return
            }
        }

        /**
         * Handles mouse selection on dropdown options.
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
         * Closes dropdown when clicking outside the component.
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


        /**
         * Emits selected billettholder id as a custom event.
         * @param {number} billettholderId
         * @returns {void}
         */
        emitBillettholderSelected(billettholderId) {
            this.dispatchEvent(
                new CustomEvent("billett-holder-selected", {
                    detail: billettholderId,
                    bubbles: true,
                    composed: true,
                }),
            )
        }

    }

    customElements.define("ticket-holder-dropdown", TicketHolderDropdown)
}
