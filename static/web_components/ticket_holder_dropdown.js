// @ts-check

/**
 * @typedef {Object} BillettHolder
 * @property {number} Id
 * @property {string} Name
 * @property {string} Email
 * @property {string} Color
 */

if (!customElements.get("ticket-holder-dropdown")) {
    class TicketHolderDropdown extends HTMLElement {
        constructor() {
            super()
            /** @type {HTMLButtonElement | null} */
            this.selectButton = null
            /** @type {HTMLUListElement | null} */
            this.dropdown = null
            /** @type {HTMLSpanElement | null} */
            this.selectedValue = null
            /** @type {number} */
            this.focusedIndex = -1
            /** @type {string} */
            this.storageKey = "selectedBillettHolder"

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

            this.selectButton = this.querySelector(".select-button")
            this.dropdown = this.querySelector(".dropdown-list")
            this.selectedValue = this.querySelector(".selected-value")

            if (!this.selectButton || !this.dropdown || !this.selectedValue) {
                return
            }

            const controlId = `dropdown-list-${ Math.random().toString(36).slice(2, 10) }`
            const buttonId = `dropdown-button-${ Math.random().toString(36).slice(2, 10) }`
            this.dropdown.id = controlId
            this.selectButton.id = buttonId
            this.selectButton.setAttribute("aria-controls", controlId)
            this.dropdown.setAttribute("aria-labelledby", buttonId)

            /** @type {{ SelectedBilletHolder?: string } | undefined} */
            const lsEnum = /** @type {any} */ (window).LSEnum
            this.storageKey = lsEnum?.SelectedBilletHolder || "selectedBillettHolder"

            this.selectButton.addEventListener("click", this.onButtonClick)
            this.selectButton.addEventListener("keydown", this.onButtonKeydown)
            this.dropdown.addEventListener("keydown", this.onDropdownKeydown)
            this.dropdown.addEventListener("click", this.onDropdownClick)
            document.addEventListener("click", this.onDocumentClick)

            this.hydrateSelection()
            this.dataset.initialized = "true"
        }

        disconnectedCallback() {
            if (this.selectButton) {
                this.selectButton.removeEventListener("click", this.onButtonClick)
                this.selectButton.removeEventListener("keydown", this.onButtonKeydown)
            }
            if (this.dropdown) {
                this.dropdown.removeEventListener("keydown", this.onDropdownKeydown)
                this.dropdown.removeEventListener("click", this.onDropdownClick)
            }
            document.removeEventListener("click", this.onDocumentClick)
            this.dataset.initialized = "false"
        }

        /**
         * @returns {HTMLLIElement[]}
         */
        getOptions() {
            return Array.from(this.querySelectorAll("li"))
        }

        /**
         * @param {HTMLLIElement} option
         * @returns {BillettHolder}
         */
        toBillettHolder(option) {
            return {
                Id: Number(option.dataset.billettHolderId || "0"),
                Name: option.dataset.billettHolderName || "",
                Email: option.dataset.billettHolderEmail || "",
                Color: option.dataset.billettHolderColor || "",
            }
        }

        /**
         * @param {HTMLLIElement} option
         * @returns {void}
         */
        saveSelected(option) {
            localStorage.setItem(this.storageKey, JSON.stringify(this.toBillettHolder(option)))
        }

        /**
         * @param {HTMLLIElement[]} options
         * @returns {void}
         */
        updateFocus(options) {
            options.forEach((option, index) => {
                option.setAttribute("tabindex", index === this.focusedIndex ? "0" : "-1")
                if (index === this.focusedIndex) {
                    option.focus()
                }
            })
        }

        /**
         * @param {HTMLLIElement} option
         * @returns {void}
         */
        renderSelected(option) {
            if (!this.selectedValue) {
                return
            }

            this.getOptions().forEach((opt) => opt.classList.remove("selected"))
            option.classList.add("selected")

            const fragment = document.createDocumentFragment()
            option.childNodes.forEach((node) => {
                fragment.appendChild(node.cloneNode(true))
            })
            this.selectedValue.replaceChildren(fragment)
        }

        /**
         * @param {boolean | null} [expand]
         * @returns {void}
         */
        toggleDropdown(expand = null) {
            if (!this.dropdown || !this.selectButton) {
                return
            }

            const options = this.getOptions()
            const isOpen = expand !== null ? expand : this.dropdown.classList.contains("hidden")
            this.dropdown.classList.toggle("hidden", !isOpen)
            this.selectButton.setAttribute("aria-expanded", String(isOpen))

            if (isOpen) {
                this.focusedIndex = options.findIndex((option) => option.classList.contains("selected"))
                this.focusedIndex = this.focusedIndex === -1 ? 0 : this.focusedIndex
                this.updateFocus(options)
                return
            }

            this.focusedIndex = -1
            this.selectButton.focus()
        }

        /**
         * @returns {void}
         */
        hydrateSelection() {
            const options = this.getOptions()
            const firstOption = options[0]
            if (!firstOption) {
                return
            }

            const selectedBillettholderLS = localStorage.getItem(this.storageKey)
            if (!selectedBillettholderLS) {
                this.renderSelected(firstOption)
                this.saveSelected(firstOption)
                return
            }

            try {
                /** @type {BillettHolder} */
                const selectedBillettholder = JSON.parse(selectedBillettholderLS)
                const selectedOption = options.find(
                    (option) => Number(option.dataset.billettHolderId || "0") === Number(selectedBillettholder.Id),
                )
                if (!selectedOption) {
                    this.renderSelected(firstOption)
                    this.saveSelected(firstOption)
                    return
                }
                this.renderSelected(selectedOption)
            } catch {
                this.renderSelected(firstOption)
                this.saveSelected(firstOption)
            }
        }

        /**
         * @param {HTMLLIElement} option
         * @returns {void}
         */
        handleOptionSelect(option) {
            this.renderSelected(option)
            this.saveSelected(option)
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
            const options = this.getOptions()
            if (options.length === 0) {
                return
            }

            if (event.key === "ArrowDown") {
                event.preventDefault()
                this.focusedIndex = (this.focusedIndex + 1) % options.length
                this.updateFocus(options)
                return
            }
            if (event.key === "ArrowUp") {
                event.preventDefault()
                this.focusedIndex = (this.focusedIndex - 1 + options.length) % options.length
                this.updateFocus(options)
                return
            }
            if (event.key === "Enter" || event.key === " ") {
                event.preventDefault()
                const option = options[this.focusedIndex]
                if (!option) {
                    return
                }
                this.handleOptionSelect(option)
                option.click()
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
            const target = event.target
            if (!(target instanceof Element)) {
                return
            }

            const option = target.closest("li")
            if (!(option instanceof HTMLLIElement)) {
                return
            }

            this.handleOptionSelect(option)
            this.toggleDropdown(false)
        }

        /**
         * @param {MouseEvent} event
         * @returns {void}
         */
        onDocumentClick(event) {
            const target = event.target
            const isOutsideClick = !(target instanceof Node) || !this.contains(target)
            if (isOutsideClick) {
                this.toggleDropdown(false)
            }
        }
    }

    customElements.define("ticket-holder-dropdown", TicketHolderDropdown)
}
