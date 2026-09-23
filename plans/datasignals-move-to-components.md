# Move All `DataSignals` Call Sites To `components.DataSignals`

## Summary
- The Datastar signal helper now lives in `components/shared_partials.go` as `components.DataSignals`, next to `KVPairsAttrs`.
- `components/formsubmission/signals.go` is currently a thin delegation kept so existing call sites keep compiling.
- This plan removes that delegation and points all 15 remaining call sites at `components.DataSignals`, so a component never has to import `formsubmission` just to render signals.

## Current State
- `components/shared_partials.go` owns the implementation (`json.Marshal`, `"{}"` on error) and documents the key-naming rule: keys are used verbatim, unlike `data-signals:<name>` which converts kebab-case to camelCase.
- `components/formsubmission/signals.go` holds two wrappers: exported `DataSignals` and package-private `dataSignals`.
- `components/event_components/event_interests.templ` already calls `components.DataSignals` directly and is the reference for how a converted call site should look.

## Key Changes

1. `components/formsubmission/event_img_upload/event_img_upload.templ` (lines 435, 462).
- Replace `formsubmission.DataSignals(` with `components.DataSignals(`.
- No import change needed: this file already imports `github.com/Regncon/conorganizer/components`.
- Check whether the `formsubmission` import is still used for anything else in the file; drop it if it is not.

2. `components/formsubmission/about_event.templ` (lines 187, 226, 281).
- Replace `dataSignals(` with `components.DataSignals(`.
- Add `"github.com/Regncon/conorganizer/components"` to the import block.

3. `components/formsubmission/contact_info.templ` (lines 123, 129, 147).
- Same replacement and import addition.

4. `components/formsubmission/other_details.templ` (lines 250, 274, 355, 379).
- Same replacement and import addition.

5. `components/formsubmission/statusCard.templ` (line 54).
- Same replacement and import addition.

6. `components/formsubmission/who_is_interested.templ` (lines 688, 754).
- Same replacement and import addition.

7. Delete `components/formsubmission/signals.go`.
- Both `DataSignals` and `dataSignals` go away once the call sites above are converted.

Line numbers are from the state of the branch when this plan was written and will drift; the grep in the verification section is the source of truth.

## Import Cycle Check
`components` (root) imports only `components/icons`, `models`, `service/eventimage`, and `service/program`. It does not import `formsubmission` or any of its subpackages, so every package listed above can import `components` without creating a cycle. `components/formsubmission/event_img_upload` already does.

## Test Plan
1. Let the templ watcher regenerate, or run `go tool templ generate`.
2. `go build ./...`.
3. `go test ./components/... ./pages/event/`.
4. Guard against leftovers — this must print nothing:

```bash
grep -rn "formsubmission\.DataSignals\|[^.]\bdataSignals(" \
  --include=*.templ --include=*.go components pages layouts service \
  | grep -v _templ.go
```

5. Spot-check rendered markup for one converted form field (for example `title` in `about_event.templ`) and confirm the `data-signals` attribute is unchanged JSON, e.g. `{"title":"..."}`.

## Assumptions
1. Pure refactor: no rendered attribute value changes, so no behavior change in the browser.
2. `components.DataSignals` keeps its current contract, including returning `"{}"` when marshalling fails.
3. Generated `*_templ.go` files are gitignored and rebuilt by the watcher or CI, so they are not part of the diff.
4. No test currently calls `formsubmission.DataSignals`; if one is added before this work starts, convert it too.

## Out Of Scope
- Converting elements that still use several `data-signals:<name>` attributes into a single `DataSignals` map. That is a separate cleanup, and the kebab-case-to-camelCase rename makes it behavior-affecting rather than mechanical.
- `components/previous_next.templ` inline `style=` attributes and other unrelated cleanups.

## Issue Text (ready to paste)
**Refactor: use `components.DataSignals` everywhere**

`DataSignals` now lives in `components/shared_partials.go`. `components/formsubmission/signals.go` is a temporary delegation so old call sites keep working, which means components still import `formsubmission` only to render Datastar signals.

Convert the remaining 15 call sites to `components.DataSignals` and delete `components/formsubmission/signals.go`:

- `components/formsubmission/event_img_upload/event_img_upload.templ` (2)
- `components/formsubmission/about_event.templ` (3)
- `components/formsubmission/contact_info.templ` (3)
- `components/formsubmission/other_details.templ` (4)
- `components/formsubmission/statusCard.templ` (1)
- `components/formsubmission/who_is_interested.templ` (2)

Pure refactor, no rendered output change. Plan with details: `plans/datasignals-move-to-components.md`.
