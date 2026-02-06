// @ts-check

/**
 * @typedef {Object} BillettHolder
 * @property {number} Id
 * @property {string} Name
 * @property {string} Email
 * @property {string} Color
 */

/**
 * Custom dropdown for selecting a ticket holder.
 *
 * Responsibilities:
 * - Manages open/close state and keyboard navigation for the listbox UI.
 * - Mirrors the selected `<li>` content into `.selected-value`.
 * - Persists selected ticket holder in `localStorage` using `window.LSEnum.SelectedBilletHolder`
 *   (falls back to `"selectedBillettHolder"`).
 * - Keeps Datastar integration intact by letting option click handlers run (`option.click()`).
 *
 * Expected markup inside the element:
 * - `.select-button` (`<button>`) for combobox trigger.
 * - `.dropdown-list` (`<ul>`) containing `<li role="option">` items.
 * - `.selected-value` (`<span>`) where the active option preview is rendered.
 * @extends {HTMLElement}
 */
class TicketHolderDropdown extends HTMLElement {
    constructor() {
        super()
        /** @type {HTMLButtonElement | null} */
        this.selectButton = null
        /** @type {HTMLUListElement | null} */
        this.dropdown = null
        /** @type {HTMLSpanElement | null} */
        this.selectedValue = null
        /** @type {HTMLLIElement[]} */
        this.options = []
    }

    /**
     * Initializes the component once per element instance.
     * @returns {void}
     */
    connectedCallback() {
        if (this.dataset.initialized === "true") {
            return
        }
        this.dataset.initialized = "true"

        this.selectButton = this.querySelector(".select-button")
        this.dropdown = this.querySelector(".dropdown-list")
        this.selectedValue = this.querySelector(".selected-value")
        this.options = Array.from(this.querySelectorAll("li"))

        const selectButton = this.selectButton
        const dropdown = this.dropdown
        const selectedValue = this.selectedValue
        const options = this.options

        if (!selectButton || !dropdown || !selectedValue || options.length === 0) {
            return
        }

        const controlId = `dropdown-list-${ Math.random().toString(36).slice(2, 10) }`
        const buttonId = `dropdown-button-${ Math.random().toString(36).slice(2, 10) }`
        dropdown.id = controlId
        selectButton.id = buttonId
        selectButton.setAttribute("aria-controls", controlId)
        dropdown.setAttribute("aria-labelledby", buttonId)

        /** @type {{ SelectedBilletHolder?: string } | undefined} */
        const lsEnum = /** @type {any} */ (window).LSEnum
        const storageKey = lsEnum?.SelectedBilletHolder || "selectedBillettHolder"
        let focusedIndex = -1

        /**
         * Converts an option element's dataset into a persisted ticket-holder shape.
         * @param {HTMLLIElement} option
         * @returns {BillettHolder}
         */
        const toBillettHolder = (option) => ({
            Id: Number(option.dataset.billettHolderId || "0"),
            Name: option.dataset.billettHolderName || "",
            Email: option.dataset.billettHolderEmail || "",
            Color: option.dataset.billettHolderColor || "",
        })

        /**
         * Persists the selected option to localStorage.
         * @param {HTMLLIElement} option
         * @returns {void}
         */
        const saveSelected = (option) => {
            localStorage.setItem(storageKey, JSON.stringify(toBillettHolder(option)))
        }

        /**
         * Updates keyboard focus state among all option elements.
         * @returns {void}
         */
        const updateFocus = () => {
            options.forEach((option, index) => {
                option.setAttribute("tabindex", index === focusedIndex ? "0" : "-1")
                if (index === focusedIndex) {
                    option.focus()
                }
            })
        }

        /**
         * Updates visual selected state and selected value preview content.
         * @param {HTMLLIElement} option
         * @returns {void}
         */
        const renderSelected = (option) => {
            options.forEach((opt) => opt.classList.remove("selected"))
            option.classList.add("selected")
            const fragment = document.createDocumentFragment()
            option.childNodes.forEach((node) => {
                fragment.appendChild(node.cloneNode(true))
            })
            selectedValue.replaceChildren(fragment)
        }

        /**
         * Opens or closes the dropdown. If `expand` is null, toggles current state.
         * @param {boolean | null} [expand]
         * @returns {void}
         */
        const toggleDropdown = (expand = null) => {
            const isOpen = expand !== null ? expand : dropdown.classList.contains("hidden")
            dropdown.classList.toggle("hidden", !isOpen)
            selectButton.setAttribute("aria-expanded", String(isOpen))

            if (isOpen) {
                focusedIndex = options.findIndex((option) => option.classList.contains("selected"))
                focusedIndex = focusedIndex === -1 ? 0 : focusedIndex
                updateFocus()
                return
            }

            focusedIndex = -1
            selectButton.focus()
        }

        /**
         * Restores selection from localStorage, or falls back to first option.
         * @returns {void}
         */
        const hydrateSelection = () => {
            const firstOption = options[0]
            if (!firstOption) {
                return
            }

            const selectedBillettholderLS = localStorage.getItem(storageKey)
            if (!selectedBillettholderLS) {
                renderSelected(firstOption)
                saveSelected(firstOption)
                return
            }

            try {
                /** @type {BillettHolder} */
                const selectedBillettholder = JSON.parse(selectedBillettholderLS)
                const selectedOption = options.find(
                    (option) => Number(option.dataset.billettHolderId || "0") === Number(selectedBillettholder.Id),
                )
                if (!selectedOption) {
                    renderSelected(firstOption)
                    saveSelected(firstOption)
                    return
                }
                renderSelected(selectedOption)
            } catch {
                renderSelected(firstOption)
                saveSelected(firstOption)
            }
        }

        /**
         * Handles a single option selection event.
         * @param {HTMLLIElement} option
         * @returns {void}
         */
        const handleOptionSelect = (option) => {
            renderSelected(option)
            saveSelected(option)
        }

        selectButton.addEventListener("click", () => {
            toggleDropdown()
        })

        selectButton.addEventListener("keydown", (event) => {
            if (event.key === "ArrowDown") {
                event.preventDefault()
                toggleDropdown(true)
                return
            }
            if (event.key === "Escape") {
                toggleDropdown(false)
            }
        })

        dropdown.addEventListener("keydown", (event) => {
            if (event.key === "ArrowDown") {
                event.preventDefault()
                focusedIndex = (focusedIndex + 1) % options.length
                updateFocus()
                return
            }
            if (event.key === "ArrowUp") {
                event.preventDefault()
                focusedIndex = (focusedIndex - 1 + options.length) % options.length
                updateFocus()
                return
            }
            if (event.key === "Enter" || event.key === " ") {
                event.preventDefault()
                const option = options[focusedIndex]
                if (!option) {
                    return
                }
                handleOptionSelect(option)
                option.click()
                toggleDropdown(false)
                return
            }
            if (event.key === "Escape") {
                toggleDropdown(false)
            }
        })

        options.forEach((option) => {
            option.addEventListener("click", () => {
                handleOptionSelect(option)
                toggleDropdown(false)
            })
        })

        document.addEventListener("click", (event) => {
            const target = event.target
            const isOutsideClick = !(target instanceof Node) || !this.contains(target)
            if (isOutsideClick) {
                toggleDropdown(false)
            }
        })

        hydrateSelection()
    }
}

customElements.define("ticket-holder-dropdown", TicketHolderDropdown)
