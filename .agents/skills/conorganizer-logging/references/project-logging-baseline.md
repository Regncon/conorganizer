# Project Logging Baseline

## Purpose

Capture the current logging architecture and non-negotiable rules for conorganizer so logging edits stay consistent.

## Architecture Snapshot

1. Logger bootstrap:
- `service/applog/logger.go` builds a JSON `slog` logger.
- Default level is `INFO`.
- `LOG_LEVEL` supports `DEBUG|INFO|WARN|WARNING|ERROR`.

2. App wiring:
- `main.go` initializes logger via `applog.NewJSONLogger()`.
- `slog.SetDefault(logger)` is called at startup.
- Component-scoped loggers are created with `logger.With("component", "...")`.
- Bind component-scoped loggers to the local variable name `logger`, not `componentLogger` or `<scope>Logger`.

3. HTTP middleware:
- `http_logging_middleware.go` emits one log per completed request.
- Current request fields:
- `method`
- `path`
- `status_code`
- `duration_ms`
- `request_id` (when available)
- Level mapping:
- `INFO`: status < 400
- `WARN`: 400-499
- `ERROR`: >= 500

4. Request IDs:
- `main.go` uses `chi` `middleware.RequestID`.
- Middleware and auth/user code read request IDs from context.

## Canonical Field Conventions

Prefer these common keys:

1. Core:
- `component`
- `error`
- `request_id`
- `path`

2. Domain examples:
- `event_id`
- `user_id`
- `pulje_id`
- `ticket_id`
- `billettholder_id`

3. Message style:
- Short and action-focused.
- No redundant data already available in structured fields.
- If a log site creates a contextualized error with `fmt.Errorf(...)`, prefer `logger.Error(fmt.Errorf(...).Error())` instead of repeating the same context in both the message and structured fields.
- For one-off log sites, keep the `fmt.Errorf(...)` inline instead of assigning it to a variable first.
- If the same exact wrapped error is used twice or more in the same scope, prefer assigning it to a variable.
- A log-plus-return pair counts as reuse and should normally use a named error variable.
- Log a failure once at the boundary that decides the outcome.
- Lower-level helpers and services should usually wrap and return errors without logging when a caller will decide the response, retry, or degraded-mode behavior.
- Route handlers, background loop tops, stream consumers, and similar handling boundaries should log the final failure before returning a response or taking an operational decision.
- If a lower layer already logs because it fully handled the failure locally, higher layers should not log the same error again; only log the new decision when it adds value.
- Preferred one-use inline pattern:
```go
logger.Error(fmt.Errorf("failed to stop watcher: %w", err).Error())
```
- Preferred reuse pattern:
```go
saveErr := fmt.Errorf("error saving event form submission: %w", insertError)
logger.Error(saveErr.Error())
http.Error(w, fmt.Sprintf("Error updating event: %v", insertError), http.StatusBadRequest)
return 0, saveErr
```
- Use `"error"` as the field key only when keeping a separate structured error field is still useful.

## Safety and Privacy Rules

Never log:

1. Session/refresh tokens.
2. Cookies or authorization headers.
3. Full request/response bodies by default.
4. Sensitive personal data unless explicitly required.

When context is needed, prefer internal IDs over raw personal fields.

## Advanced Runtime Rules

1. Expected failures should usually not be logged at `Error`.
2. Treat `sql.ErrNoRows`, validation failures, canceled requests, client disconnects, and similar expected outcomes as `Debug`, `Info`, `Warn`, or no log depending on the handling boundary and user impact.
3. Repeated watcher, polling, stream, and retry-loop logs should default to `Debug` unless they indicate a state transition or actionable problem.
4. Avoid one log per steady-state loop iteration; prefer rate-limiting, aggregation, or summary logs for noisy paths.
5. Async and background flows should carry correlation fields that help reconstruct the sequence.
6. Prefer stable, actionable identifiers that operators can actually use: `request_id` when it still exists in context, otherwise domain IDs like `event_id`, `pulje_id`, `user_id`, or `billettholder_id`.
7. Do not add ephemeral internal values like in-memory session IDs, transient stream subjects, or similar correlation fields unless they are persisted, searchable, or already used operationally in this repo.
8. Log lifecycle transitions once at the operational boundary: startup complete, shutdown start/complete, entering degraded mode, retry scheduled, retries exhausted, and recovery after failure.
9. Do not emit success logs for every normal heartbeat, poll, or stream tick.
10. If a path emits high-volume state-change logs, consider metrics/tracing instead of per-item logs; use logs for transitions, anomalies, and decisions.

## Error Wrapping Rules

When updating returned errors:

1. Wrap with `fmt.Errorf("...: %w", err)` when the current function can add useful local context.
2. Wrap when an error crosses a package, subsystem, or abstraction boundary and the higher layer can explain the failed operation.
3. Prefer direct returns like `return fmt.Errorf("...: %w", err)` when the wrapped error is only returned once.
4. If the same exact wrapped error is used twice or more in the same block, prefer a local error variable.
5. A log-plus-return pair counts as reuse and should normally use a named error variable.
6. Keep `return err` when the current function adds no meaningful context.
7. Avoid `%w` at HTTP handlers, `main`, and similar top-level response/logging boundaries; log the accumulated context and return a safe error or response instead.
8. Use `%v` or translate to an exported repo error when callers should not depend on an underlying dependency error type.
9. Do not double-wrap errors that already contain the needed local context.
10. At log-only boundaries, prefer contextualizing the error with `fmt.Errorf(...)` and logging `logger.Error(err.Error())`.
11. Do not introduce temporary `wrappedErr`-style variables for one-use errors.
12. Avoid double-logging across layers. Helpers usually wrap and return; the handling boundary usually emits the log.

## Migration Rules for Existing Logs

When updating existing logging:

1. Replace `"err"` fields with `"error"`.
2. Convert `fmt.Print*` diagnostic logs in production paths to structured `logger.*` calls.
3. Rename `componentLogger` and `<scope>Logger` locals to `logger` when touching that scope.
4. Run a final pass over returned errors and wrap them with `fmt.Errorf(...: %w)` only when the function can add useful context.
5. At log-only boundaries, prefer `logger.Error(fmt.Errorf(...).Error())` over duplicating the same context in both the message and structured fields.
6. Prefer inline wrapping for one-use errors, but create a temporary error variable when the same exact error is used twice or more.
7. Remove duplicate logging of the same failure across helper and handler layers when touching that path.
8. Demote expected failures and noisy loop logs when touching watchers, streams, polling, or request-boundary code.
9. Add missing correlation fields to async/background logs when the surrounding context makes that practical.
10. Keep migration scope focused to touched areas.
11. Avoid unrelated refactors.

## Fast Verification Queries

```powershell
rg -n 'logger\.With\(\"component\",' -S .
rg -n 'componentLogger|[a-z][A-Za-z0-9]*Logger\s*:=' -S . --glob '*.go' --glob '*.templ'
rg -n 'logger\.(Debug|Info|Warn|Error)\(\"[^\"]*\"\s*,\s*\"err\"' -S . --glob '!backup-*'
rg -n 'fmt\.Println|fmt\.Printf|log\.Print|log\.Fatal' -S . --glob '!backup-*'
rg -n 'return .*err$' -S service pages components --glob '*.go' --glob '*.templ'
rg -n 'request_id|status_code|duration_ms|http request completed' -S main.go http_logging_middleware.go service pages components
```

## Source Context Used

This baseline is derived from:

1. `.ai/threads/new-log-structure.md`
2. `service/applog/logger.go`
3. `main.go`
4. `http_logging_middleware.go`
5. `service/authctx/authctx.go`
6. `service/userctx/userctx.go`
