# Form Textarea Trim Normalization (`notes`, `intro`, `description`)

## Summary
Standardize Datastar textarea behavior so all three fields trim whitespace on initial signal value and on user input, while preserving existing SSE `PUT` on change behavior.

## Key Changes
1. Update `other-notes` in `components/formsubmission/other_details.templ`.
- Keep `data-signals:notes={ fmt.Sprintf("'%s'", strings.TrimSpace(notes)) }`.
- Keep `data-bind:notes`.
- Keep `data-on-input="$notes = $notes.trim()"`.
- Keep existing `data-on-change={ datastar.PutSSE("/my-events/api/new/%s/notes", eventId) }`.

2. Update `intro` in `components/formsubmission/about_event.templ`.
- Add `data-signals:intro={ fmt.Sprintf("'%s'", strings.TrimSpace(intro)) }`.
- Change `data-bind="intro"` to `data-bind:intro`.
- Add `data-on-input="$intro = $intro.trim()"`.
- Keep existing `data-on-change={ datastar.PutSSE("/my-events/api/new/%s/intro", eventId) }`.

3. Update `description` in `components/formsubmission/about_event.templ`.
- Add `data-signals:description={ fmt.Sprintf("'%s'", strings.TrimSpace(description)) }`.
- Change `data-bind="description"` to `data-bind:description`.
- Add `data-on-input="$description = $description.trim()"`.
- Keep existing `data-on-change={ datastar.PutSSE("/my-events/api/new/%s/description", eventId) }`.

4. Imports.
- Ensure `about_event.templ` includes `strings`.
- Keep `strings` import in `other_details.templ`.

## Test Plan
1. For `notes`, `intro`, and `description`, type leading/trailing spaces and verify they are trimmed during input.
2. Verify internal spaces remain unchanged.
3. Verify whitespace-only input becomes empty string.
4. Blur/change each field and verify trimmed values are sent to:
- `/my-events/api/new/{id}/notes`
- `/my-events/api/new/{id}/intro`
- `/my-events/api/new/{id}/description`
5. Reload and verify all three textareas initialize with trimmed values.

## Assumptions
1. Keep the exact `data-signals:*` string approach (`fmt.Sprintf("'%s'", strings.TrimSpace(...))`) without additional escaping strategy.
2. Trimming on every input is intended for all three fields.
3. No backend API, schema, or type changes are required.
