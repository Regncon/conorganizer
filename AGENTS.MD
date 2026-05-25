## Check domain domeneordbok.md
when adding names check domeneordbok.md instead of translating to english

# CSS rules

### Required CSS style
- Use native CSS nesting.
- Use container queries for responsive component styling.
- Prefer component-based responsiveness over viewport-based media queries.
- Use modern CSS features where appropriate:
  - custom properties
  - logical properties
  - `:has()`
  - `clamp()`
  - `min()`
  - `max()`
  - `color-mix()`
  - container query units like `cqw`, `cqh`, `cqi`, `cqb`
- Keep selectors shallow and readable.
- Avoid Sass-style CSS that does not work natively in browsers.
- Avoid viewport media queries for components. Use `@container` instead.
- Use viewport media queries only for page-level/global layout concerns.


## Go test structure
When writing or cleaning up Go tests, prefer behavior-focused tests that read from intent to verification.

- Add a short BDD-style comment at the top of each test when it improves readability:
  - Given some initial context,
  - When an event occurs,
  - Then ensure the expected outcome.
- Split each test into three visible sections:
  - `// Given`
  - `// When`
  - `// Then`
- Start the `// Given` section with the expected output or expected result. This makes the test's purpose clear before setup details.
- Keep test names specific to one behavior. Prefer splitting creation, duplicate/idempotency, validation, and error cases into separate tests.
- Move noisy setup SQL and repeated assertions into local test helpers when that makes the test read like behavior.

Example shape:

```go
func TestService_EventProducesExpectedResult(t *testing.T) {
	// Given an existing domain object,
	// when the service reconciles related data,
	// then the expected durable relation is stored.

	// Given
	expectedResult := models.SomeResult{ID: 123}
	db := createTestDB(t)
	insertRequiredFixture(t, db, expectedResult.ID)

	// When
	err := ServiceFunction(expectedResult.ID, db)

	// Then
	if err != nil {
		t.Fatalf("expected service call to succeed: %v", err)
	}
	assertExpectedResult(t, db, expectedResult)
}
```

### Production helpers vs test helpers
Reuse production functions for the behavior under test. Do not add production service functions only to make test setup nicer.

- Use production helpers when they represent real application behavior that the test is exercising.
- Use test helpers for arranging prerequisite database state, such as inserting fixture users, billettholdere, emails, or association rows.
- If setup helpers become useful across packages, move them to `testutil`, not to production services.
- Avoid arranging a test through the exact code path being verified unless the test is explicitly an integration or route test.
- Route-level tests can exercise handlers and their production wiring, but keep service-level regression tests focused on the service behavior.

## When creating custom components that use signals

Datastar interacts with custom web components using a "props down, events up" model:
- Data flows into the component via HTML attributes.
- Changes inside the component are sent back via custom events.

### Interaction methods

#### Passing data in (Datastar -> Web Component)
- Use the `data-attr` attribute to bind a Datastar signal to a custom element attribute.
- The web component must observe that attribute (via `attributeChangedCallback`).

Example:
```html
<my-component data-attr:name="$_name"></my-component>
```

#### Getting data out (Web Component -> Datastar)
- Dispatch a custom event from the web component when internal state changes.
- Datastar listens with `data-on` and updates a signal.
- `evt.detail` or `$_eventDetails` provides the payload.

Example:
```html
<my-component data-on:custom-change="$_signal = evt.detail.value"></my-component>
```

### Key attributes and variables
- `data-attr:[attribute-name]`: Binds a Datastar signal to a component attribute.
- `data-on:[event-name]`: Listens for a (custom) event from the component.
- `evt`: Event object available in Datastar expressions.
- `evt.detail` or `$_eventDetails`: Event payload.

### Example flow
1. A user interacts with a `<reverse-component>` custom element.
2. The component reverses a string internally.
3. The component dispatches a `reverse` event with the result in `event.detail`.

JavaScript inside the component:
```javascript
// Inside your web component class
const event = new CustomEvent('reverse', {
  detail: { value: reversedValue }
});
this.dispatchEvent(event);
```

Datastar binding:
```html
<reverse-component
  data-on:reverse="$_reversed = evt.detail.value"
  data-attr:name="$_name">
</reverse-component>
```

This sets the `$_reversed` signal to the value from the event payload.

For more details and live examples, refer to the Datastar documentation on web components.

### Sub Agents
Use subagents when this task would be faster or more reliable with them.

## finishing instructions
When done reading this you will output "Done reading AGENTS.MD file."
