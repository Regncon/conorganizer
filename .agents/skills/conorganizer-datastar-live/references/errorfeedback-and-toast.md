# errorfeedback and toast

Two small conventions carry nearly all user feedback for Datastar requests. `components/errorfeedback` renders inline field errors from one shared signal, `feedbackErrors`. `components/toast` shows a "Lagret" toast when a save request finishes. `static/js/error_feedback.js` and `static/web_components/toast.js` glue both to Datastar's `datastar-fetch` and `datastar-signal-patch` events, and they know about each other: a backend field error for the trigger's own feedback key suppresses the toast for that request.

## errorfeedback: server side

```go
sse := datastar.NewSSE(w, r)
feedback := errorfeedback.New("name", "email")   // every key this response owns
if name == "" {
	feedback.Set("name", "Name must be filled in.")
}
if feedback.HasErrors() {
	return feedback.Patch(sse)                    // sends {"feedbackErrors": {"name": "...", "email": ""}}
}
```

- `New(keys...)` initializes each key to an empty string. `Patch` sends all of them, so keys you listed are cleared on the page when they have no new error. A key you did not list is left as it was, which is how a stale message survives a successful save.
- `Set`, `Clear`, `HasErrors`, `Patch(sse)` and `Signals()` are the whole API. `Signals()` returns the `{"feedbackErrors": {...}}` map so you can merge it into one `MarshalAndPatchSignals` call together with other signals.
- Because `Patch` needs the SSE, create it before the validation branch that reports field errors, and keep `http.Error` for the checks that come before `NewSSE` (malformed signals, missing ids, access).
- On the feedback page, `/tilbakemelding/reset` patches a fresh form and then `newFeedbackErrors().Patch(sse)` explicitly, because the message elements use `__ifmissing` and would keep showing old errors otherwise.

## errorfeedback: markup

```templ
<form { errorfeedback.DefaultMessage("Kunne ikke lagre. Prøv igjen.")... }>
	<input
		data-bind="title"
		data-indicator="toastTitle"
		{ errorfeedback.Attrs("title")... }
		data-on:change={ fmt.Sprintf("@put('/profile/api/new/%s/title')", eventId) }
	/>
	@errorfeedback.RenderMessage("title")
</form>
```

- `Attrs(key)` adds `data-feedback-key` to the element that starts the request; `Attrs(key, msg)` adds a one-off `data-feedback-message` fallback as well.
- `DefaultMessage(msg)` on the closest `form` or `dialog` sets `data-feedback-default-message` for all triggers inside.
- `RenderMessage(key)` renders a `<span class="feedback-error">` bound to `$feedbackErrors["key"]`, with `data-signals__ifmissing` initializing that key to `""`, `hidden` toggled by `data-attr:hidden`, and `aria-live="polite"`.
- `errorfeedback.Clear(key)` and `ClearMany(keys...)` return Datastar expressions for buttons that should clear a message locally without a request.
- Value signals for the inputs are initialized with `components.DataSignals` or `data-signals`, never through this package.

## errorfeedback: what the browser does

`static/js/error_feedback.js` listens to `datastar-fetch` for elements that carry `data-feedback-key`:

| Fetch event | Effect on `$feedbackErrors[key]` |
| --- | --- |
| `started` | cleared |
| `retrying` | set to the fallback message |
| `error`, `retries-failed` | set to the fallback message, trigger remembered as failed |
| `finished` | cleared, unless the trigger failed or the backend patched a non-empty message for that key during the request |

What a non-200 answer looks like in the browser (checked against the bundled v1.0.4 with a 404 and a 501 on 2026-10-18): `started`, `error`, `finished`, and no retry with the default `retry: 'auto'`. The response body is dropped. So an `http.Error(w, "Epostadressen finnes allerede", 400)` before the SSE shows the trigger's fallback message, never that text, and suppresses the toast. `retrying` and `retries-failed` come from a fetch that throws, such as a server that is down. Validation whose exact wording matters goes back as a 200 with a `feedbackErrors` patch.

Fallback message precedence: `data-feedback-message` on the trigger, then `data-feedback-default-message` on the closest `form` or `dialog`, then "Klarte ikke å lagre endringen. Prøv igjen.".

A backend `feedbackErrors` patch received while the request is active marks the trigger as failed, so the later `finished` event does not wipe the message. The script's own local patches are wrapped in `feedback-errors-local-patch-start/end` events so it and the toast can tell them apart from backend patches.

## toast

```templ
@toast.Render("Lagret", icons.ProgressComplete)
```

- Render it once per page, in the `_index` templ next to the live wrapper rather than inside it; `pages/profile/newevent/new_index.templ` and `pages/admin/approval/editForm/edit_form_index.templ` do this. The `<app-toast>` element moves itself to `document.body`, and only the first active host listens to global events.
- Saved toast: put `data-indicator="toastSomething"` on the trigger (prefix `toast`, overridable with `indicator-prefix` on the host). Use the value form: `toast.js` reads `el.dataset.indicator`, so the colon form `data-indicator:toast-something` never shows a toast. When that request reaches `finished` without `error` or `retries-failed`, the host shows its `indicator-message` ("Lagret"). A backend `feedbackErrors` patch suppresses it only when it carries a non-empty message for the trigger's own `data-feedback-key`; a trigger without `errorfeedback.Attrs` still shows "Lagret" next to the error.
- Manual toast from JS or a web component: dispatch a bubbling, composed `CustomEvent("toast", { detail: { message } })`. The event image upload does this from `data-init` on a hidden span after a successful upload.
- A custom message from Datastar: `data-on:click="evt.target.dispatchEvent(new CustomEvent('toast', { bubbles: true, composed: true, detail: { message: 'Kopiert' } }))"`.

## Choosing between the mechanisms

| Need | Use |
| --- | --- |
| Validation message next to a control, cleared on the next successful save | `errorfeedback` key per control |
| "Saved" confirmation for an auto-saving field | `data-indicator="toast<Field>"` plus `errorfeedback.Attrs` with the same key the handler patches |
| One-off status text for an action button, such as success, info or error | plain signals patched from the handler, as `getTicketsSuccessMessage` and friends in `pages/profile/tickets/tickets_page.templ`. Declare them kebab-case, `data-signals:get-tickets-success-message="''"`; that file's camelCase declarations are lowercased by HTML and declare other signals |
| Server error after the SSE already started | `errorfeedback` if there is a key, otherwise a plain signal; `sse.ConsoleError(err)` for the console |
| Navigation after a mutation | `sse.Redirect(url)` when `Datastar-Request: true`, else `http.Redirect` with 303 |

The per-card `errorMessage-<id>` and `successMessage-<id>` signals in the billettholder cards predate `errorfeedback`. They still work; prefer `errorfeedback` for new field validation so the toast suppression and the fallback messages come for free.
