# Datastar signals

A signal is a value in the browser that Datastar expressions can read. Signals
connect user input, server responses, and the elements that display those values.
For example, `data-text="$saveMessage"` updates the element when `saveMessage`
changes, and `data-show="$_detailsOpen"` controls whether an element is visible.

Signals and live buckets have different responsibilities:

- A live bucket tells subscribed server connections to render fresh HTML.
- A signal controls browser state. It can be initialized by that HTML or updated
  by an endpoint's signal response.

Signals are not stored in NATS buckets. A signal response goes to the browser
making that request. To notify other open pages, save the change and broadcast
the relevant bucket; those pages then receive their own HTML updates.

## Choose who owns the value

Use the `_` prefix for browser-only state, such as whether a details panel is
expanded. Datastar excludes these signals from its default request payload.
Use ordinary names for server-provided data and selections sent to an endpoint.
This is the project's naming convention; the prefix does not grant access or
replace validation on the server.

`data-signals__ifmissing` supplies defaults without overwriting an existing
selection or a newer endpoint response. Plain `data-signals` applies the value
when Datastar processes the attribute, including when a live HTML patch changes
its value. Choose plain initialization when the server should replace the value.

## Example: expand details in the browser

```html
<section data-signals__ifmissing:_details-open="false">
    <button data-on:click="$_detailsOpen = !$_detailsOpen">
        Show details
    </button>
    <p data-show="$_detailsOpen" style="display: none;">
        Additional information.
    </p>
</section>
```

This interaction needs no endpoint. The signal is browser-only, and `ifmissing`
lets the user's choice survive later HTML updates that retain this component.

## Example: report the result of saving a form

Initialize a server-provided message and bind it to the UI:

```html
<div data-signals__ifmissing:save-message="''">
    <button data-on:click="@post('/example/save')">Save</button>
    <p role="status" data-text="$saveMessage"></p>
</div>
```

The URL is illustrative. After the real save handler has validated and saved the
form, it can call a helper like this:

```go
func patchSaveMessage(w http.ResponseWriter, r *http.Request, message string) error {
    sse := datastar.NewSSE(w, r)
    return sse.MarshalAndPatchSignals(map[string]string{
        "saveMessage": message,
    })
}
```

The caller handles/logs the returned error. This updates the status text without
replacing the form or disturbing its input fields. Use the SDK's JSON helper so
quotes and other characters in messages are encoded correctly.

## Example: update a room count through ordinary live HTML

A live component can carry a server-owned signal in its markup:

```templ
templ RoomSummary(count int) {
    <section id="room-summary" data-signals:room-count={ strconv.Itoa(count) }>
        <span data-text="$roomCount">{ strconv.Itoa(count) }</span> rooms
    </section>
}
```

This example uses Go's `strconv` package. The parent live page loads the count
from the database and renders this component. After a room is added and the
handler broadcasts `live.BucketRooms`, subscribers render it again. A changed
`data-signals:room-count` value updates the browser signal. The initial count is
also rendered as text, so the page has content before Datastar starts.

No extra signal callback is needed in `live.Page`: its existing HTML patch
carries the initializer.

## Effects and requests

`data-effect` runs initially and when signals it reads change. For example:

```html
<div data-signals__ifmissing:search-text="''"
     data-effect="if ($searchText.length >= 3) { @get('/example/search') }">
    <input data-bind:search-text />
</div>
```

The illustrative search endpoint reads `searchText` from the request and returns
results. The bundled Datastar client collects request signals without making
the effect depend on every value in the payload. Keep the effect's explicit
dependencies limited to the values that should trigger a request, and avoid
having the response change those same values repeatedly.

A long-lived SSE request keeps the selection and cookies it received when it
opened. When results depend on a selection the user can change, send that
current selection in a fresh request. Validate it on the server.

Use normal morphing for server-rendered components. Reserve `data-ignore-morph`
for DOM that another widget owns: it also prevents ordinary text and markup
inside that element from receiving server updates.
