// @ts-check
const STYLE_ID = "ticket-holder-dropdown-styles"
const STYLE_TEXT = `
ticket-holder-dropdown.custom-select {
    position: relative;
    display: inline-block;
    width: 100%;

    .name-initials {
        display: inline-flex;
        place-items: center;
        gap: var(--spacing-2x);
        min-width: 0;
        overflow: clip;

        .name {
            white-space: nowrap;
            overflow: clip;
        }

        .initials {
            display: flex;
            place-content: center;
            place-items: center;
            color: var(--color-text-strong);
            font-size: 12px;
            border-radius: 50%;
            min-inline-size: 1.5rem;
            min-block-size: 1.5rem;
        }
    }

    .select-button {
        display: flex;
        place-content: space-between;
        place-items: center;
        width: 100%;
        cursor: pointer;
        padding-inline: var(--spacing-4x);
        padding-block: var(--spacing-3x);

        .select-button-end {
            .arrow {
                display: flex;
                transition: transform ease-in-out 0.3s;
            }
        }

        .selected-value {
            .name-initials {
                max-width: 12.8rem;
            }
        }

        &[aria-expanded="true"] {
            .select-button-end {
                .arrow {
                    transform: rotate(180deg);
                }
            }
        }
    }

    .dropdown-list {
        position: absolute;
        top: 100%;
        left: 0;
        width: 100%;
        background-color: var(--bg-surface);
        border: 1px solid var(--bg-item-border);
        border-radius: 0.25rem;
        list-style: none;
        padding: 10px;
        margin: 10px 0 0;
        box-shadow: var(--shadow-dialog);
        max-height: 200px;
        overflow-y: auto;
        scrollbar-width: thin;
        display: flex;
        flex-direction: column;
        gap: 0.5rem;

        &.hidden {
            display: none;
        }

        li {
            background-color: var(--bg-item);
            border: 1px solid var(--bg-item-border);
            border-radius: var(--border-radius-2x);
            padding: 10px;
            cursor: pointer;
            display: flex;
            gap: 0.5rem;
            align-items: center;

            &.selected {
                border-color: var(--color-primary);
                color: var(--color-primary);
                font-weight: bold;
            }

            &:hover,
            &:focus-visible {
                border-color: var(--color-primary);
                color: var(--color-primary);
            }
        }
    }
}
`

/**
 * @typedef {Object} BillettHolder
 * @property {number} Id
 * @property {string} Name
 * @property {string} Email
 * @property {string} Color
 */

/**
 * @returns {void}
 */
function ensureStyles() {
    if (document.getElementById(STYLE_ID)) {
        return
    }
    const styleEle = document.createElement("style")
    styleEle.id = STYLE_ID
    styleEle.textContent = STYLE_TEXT
    document.head.appendChild(styleEle)
}

if (!customElements.get("ticket-holder-dropdown")) {
    /**
     * Ticket-holder dropdown custom element.
     *
     * Required input:
     * - `data-billettholders`: JSON array of ticket holders.
     *
     * Optional templ-provided icon:
     * - Provide a child `<template data-arrow-icon>...</template>`.
     * - The component clones this template into the arrow slot during render.
     * - If omitted, it falls back to a plain text arrow.
     */
    class TicketHolderDropdown extends HTMLElement {
        constructor() {
            super()
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

        connectedCallback() {
            if (this.dataset.initialized === "true") {
                return
            }

            ensureStyles()
            this.billettholdere = this.parseHolders()
            if (this.billettholdere.length === 0) {
                return
            }
            this.arrowIconTemplateEle = this.querySelector("template[data-arrow-icon]")

            this.render()

            this.selectButtonEle = this.querySelector(".select-button")
            this.dropdownEle = this.querySelector(".dropdown-list")
            this.selectedValueEle = this.querySelector(".selected-value")
            if (!this.selectButtonEle || !this.dropdownEle || !this.selectedValueEle) {
                return
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

            this.hydrateSelection()
            this.dataset.initialized = "true"
        }

        disconnectedCallback() {
            if (this.selectButtonEle) {
                this.selectButtonEle.removeEventListener("click", this.onButtonClick)
                this.selectButtonEle.removeEventListener("keydown", this.onButtonKeydown)
            }
            if (this.dropdownEle) {
                this.dropdownEle.removeEventListener("keydown", this.onDropdownKeydown)
                this.dropdownEle.removeEventListener("click", this.onDropdownClick)
            }
            document.removeEventListener("click", this.onDocumentClick)
            this.dataset.initialized = "false"
        }

        /**
         * @returns {BillettHolder[]}
         */
        parseHolders() {
            const raw = this.getAttribute("data-billettholders")
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

            this.replaceChildren(buttonEle, listEle)
        }

        /**
         * @returns {HTMLLIElement[]}
         */
        getOptions() {
            return Array.from(this.querySelectorAll("li"))
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
