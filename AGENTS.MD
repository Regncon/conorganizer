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

## finishing instructions
When done reading this you will output "Done reading AGENTS.MD file."
