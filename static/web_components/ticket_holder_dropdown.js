// @ts-check

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
    const style = document.createElement("style")
    style.id = STYLE_ID
    style.textContent = STYLE_TEXT
    document.head.appendChild(style)
}

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
            /** @type {BillettHolder[]} */
            this.billettholdere = []

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

            this.render()

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
            const wrapper = document.createElement("div")
            wrapper.className = "name-initials"

            const initials = document.createElement("span")
            initials.className = "initials"
            if (holder.Color) {
                initials.style.backgroundColor = `hsl(from ${ holder.Color } h s l / 0.5)`
                initials.style.border = `1px solid ${ holder.Color }`
            }
            initials.textContent = this.getInitials(holder.Name)

            const name = document.createElement("p")
            name.className = "name"
            name.textContent = holder.Name

            wrapper.appendChild(initials)
            wrapper.appendChild(name)
            return wrapper
        }

        /**
         * @returns {void}
         */
        render() {
            const button = document.createElement("button")
            button.className = "select-button input no-marking"
            button.setAttribute("role", "combobox")
            button.setAttribute("aria-label", "select button")
            button.setAttribute("aria-haspopup", "listbox")
            button.setAttribute("aria-expanded", "false")
            button.type = "button"

            const selectedValue = document.createElement("span")
            selectedValue.className = "selected-value"

            const buttonEnd = document.createElement("div")
            buttonEnd.className = "select-button-end"
            const arrow = document.createElement("i")
            arrow.className = "arrow"
            arrow.textContent = "▾"
            buttonEnd.appendChild(arrow)

            button.appendChild(selectedValue)
            button.appendChild(buttonEnd)

            const list = document.createElement("ul")
            list.className = "dropdown-list hidden"
            list.setAttribute("role", "listbox")

            this.billettholdere.forEach((holder) => {
                const li = document.createElement("li")
                li.setAttribute("role", "option")
                li.dataset.billettHolderId = String(holder.Id)
                li.dataset.billettHolderName = holder.Name
                li.dataset.billettHolderEmail = holder.Email
                li.dataset.billettHolderColor = holder.Color
                li.setAttribute("data-bind", "billettHolderId")
                li.setAttribute("data-on:click", `$billettHolderId = ${ holder.Id }`)
                li.appendChild(this.createNameInitialsNode(holder))
                list.appendChild(li)
            })

            this.replaceChildren(button, list)
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

            const holder = this.toBillettHolder(option)
            this.selectedValue.replaceChildren(this.createNameInitialsNode(holder))
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


const STYLE_ID = "ticket-holder-dropdown-styles"
const STYLE_TEXT = `
ticket-holder-dropdown.custom-select {
    position: relative;
    display: inline-block;
    width: 100%;
}

ticket-holder-dropdown .select-button {
    display: flex;
    place-content: space-between;
    place-items: center;
    width: 100%;
    cursor: pointer;
    padding-inline: var(--spacing-4x);
    padding-block: var(--spacing-3x);
}

ticket-holder-dropdown .select-button .select-button-end .arrow {
    display: flex;
    transition: transform ease-in-out 0.3s;
}

ticket-holder-dropdown .select-button[aria-expanded="true"] .select-button-end .arrow {
    transform: rotate(180deg);
}

ticket-holder-dropdown .selected-value .name-initials {
    max-width: 12.8rem;
}

ticket-holder-dropdown .dropdown-list {
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
}

ticket-holder-dropdown .dropdown-list.hidden {
    display: none;
}

ticket-holder-dropdown .dropdown-list li {
    background-color: var(--bg-item);
    border: 1px solid var(--bg-item-border);
    border-radius: var(--border-radius-2x);
    padding: 10px;
    cursor: pointer;
    display: flex;
    gap: 0.5rem;
    align-items: center;
}

ticket-holder-dropdown .dropdown-list li.selected,
ticket-holder-dropdown .dropdown-list li:hover,
ticket-holder-dropdown .dropdown-list li:focus-visible {
    border-color: var(--color-primary);
    color: var(--color-primary);
}

ticket-holder-dropdown .dropdown-list li.selected {
    font-weight: bold;
}
`
