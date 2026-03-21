---
name: conorganizer-logging
description: Create, migrate, and review structured logging in the conorganizer Go codebase using log/slog. Use when adding new logs, updating old logs after refactors, or checking PRs for logging quality in .go and .templ handlers, middleware, services, and route setup.
---

# Conorganizer Logging

Use this skill to make logging changes consistent with the repo's current JSON `slog` setup and field conventions.

## Quick Baseline

Read these first:

1. `service/applog/logger.go` for logger initialization and `LOG_LEVEL`.
2. `http_logging_middleware.go` for request log format and level mapping.
3. `main.go` for wiring (`slog.SetDefault`, request ID middleware, component-scoped logger usage).
4. `references/project-logging-baseline.md` for rules, examples, and error-wrapping guidance.

## Workflow

1. Identify scope.
- Decide whether the request is:
- Add logging in new code.
- Migrate existing logs.
- Review PR logging quality only.

2. Pick logger ownership and component.
- Reuse passed `*slog.Logger` when available.
- Create or rebind a scoped logger near route/service entry.
- Always bind the component-scoped logger to the variable name `logger`.
- If the current scope already has `logger`, reassign it: `logger = logger.With("component", "<component_name>")`.
- If the parent logger has a different name, derive `logger` from it: `logger := baseLogger.With("component", "<component_name>")`.
- Never introduce `componentLogger` or `<scope>Logger` variable names.
- Keep component names short and stable (`main`, `http`, `auth`, `event_service`, `profile_tickets`).

3. Apply structured log rules.
- Use `logger.Info/Warn/Error/Debug` with message plus key/value pairs.
- When you do use a structured error field, use `"error"` as the key (not `"err"`).
- Prefer stable context keys: `event_id`, `user_id`, `pulje_id`, `ticket_id`, `billettholder_id`, `path`, `request_id`.
- Keep `msg` short and action-oriented.
- Prefer building a contextualized error with `fmt.Errorf(...)` and logging `logger.Error(fmt.Errorf(...).Error())` over repeating the same context in both the message and structured fields.
- Prefer inline wrapping for one-off logs: `logger.Error(fmt.Errorf("...: %w", err).Error())`.
- Do not introduce variables like `wrappedErr`, `queryErr`, or `insertErr` when the wrapped error is only used once.
- If the same wrapped error is used twice or more in the same scope, prefer assigning it to a variable.
- A log-plus-return pair counts as reuse and should normally use a named error variable.
- Log a failure once at the boundary that decides the outcome.
- Lower-level helpers and services should usually wrap and return errors without logging when the caller will decide the response, retry, or degraded-mode behavior.
- If a helper or service no longer logs after a migration, remove the unused `logger` parameter from its signature unless it is still needed for non-error logs in that function.
- When you change a function signature, update all call sites in touched source files. If the function is called from generated `*_templ.go` files and you are not running `templ generate`, update those generated call sites too so the tree stays consistent.
- Route handlers, background loop tops, stream consumers, and similar handling boundaries should log the final failure before returning a response or taking an operational decision.
- If a lower layer already logs because it fully handled the failure locally, higher layers should not log the same error again; only log the new decision if that adds value.
- If you know a lower layer already returned a useful wrapped `fmt.Errorf(...)`, prefer logging `logger.Error(err.Error())` at the handling boundary instead of wrapping again just to restate the same failure.
- If the callee is in the touched repo and you can see that its returned errors are already wrapped with `fmt.Errorf(...)`, treat it as known-wrapped and use `logger.Error(err.Error())` at the boundary.
- If you do not know whether the returned error already carries useful wrapped context, prefer wrapping it at the boundary with `fmt.Errorf(...)` before logging so the boundary still contributes a clear operation-specific error.
- Only skip the extra `fmt.Errorf(...)` at the boundary when the existing returned error is already known to carry the needed context.
- Treat the boundary log shape and the returned error shape as separate decisions.
- Normalizing a boundary log to `logger.Error(err.Error())` does not imply `return err`.
- If the current function adds meaningful caller-local context for its own caller, keep wrapping on `return` even when the log uses `err.Error()`.
- Only collapse to `return err` when this function adds no useful context beyond what the callee already returned.
- Preferred one-use inline example:
```go
logger.Error(fmt.Errorf("failed to stop watcher: %w", err).Error())
```
- Preferred boundary log plus wrapped return example:
```go
logger.Error(err.Error(), "user_id", userInfo.Id)
return nil, fmt.Errorf("unable to get billettholder for user %s: %w", userInfo.Id, err)
```
- Preferred reuse example:
```go
saveErr := fmt.Errorf("error saving event form submission: %w", insertError)
logger.Error(saveErr.Error())
http.Error(w, fmt.Sprintf("Error updating event: %v", insertError), http.StatusBadRequest)
return 0, saveErr
```
- Example: `logger.Error(fmt.Errorf("event image directory %q does not exist: %w Create it and run task start again", *eventImageDir, statErr).Error())`.

4. Handle HTTP requests consistently.
- Keep one completion log per request in middleware.
- Ensure fields: `method`, `path`, `status_code`, `duration_ms`, and optional `request_id`.
- Keep level mapping:
- `INFO` for <400.
- `WARN` for 4xx.
- `ERROR` for 5xx.

5. Protect sensitive data.
- Do not log tokens, cookies, authorization headers, or raw request/response bodies.
- Avoid personal data when possible (especially email in info/debug paths).
- If user identifiers are needed, prefer internal IDs.

6. Apply advanced runtime logging rules.
- Expected failures should usually not be logged at `Error`.
- Treat `sql.ErrNoRows`, validation failures, canceled requests, client disconnects, and similar expected outcomes as `Debug`, `Info`, `Warn`, or no log depending on the boundary and user impact.
- Repeated watcher, polling, stream, and retry-loop logs should default to `Debug` unless they represent a state transition or actionable problem.
- Avoid one log per steady-state loop iteration; prefer rate-limiting, aggregation, or summary logs for noisy paths.
- In async and background flows, include correlation fields that help stitch events together.
- Prefer stable, actionable identifiers that operators can actually use: `request_id` when it still exists, otherwise domain IDs like `event_id`, `pulje_id`, `user_id`, or `billettholder_id`.
- Do not add ephemeral internal values like in-memory session IDs, transient stream subjects, or similar correlation fields unless they are persisted, searchable, or already used operationally in this repo.
- Log lifecycle transitions once: startup complete, shutdown start/complete, entering degraded mode, retry scheduled, retries exhausted, and recovery after failure.
- Do not emit success logs for every normal heartbeat, poll, or stream tick.
- If a path emits high-volume state-change logs, consider whether the signal belongs in metrics/tracing instead of per-item logs.

7. Wrap returned errors deliberately.
- Run this pass after logging edits and before final validation.
- Wrap returned errors with `fmt.Errorf("...: %w", err)` when adding useful local context or crossing a package/system boundary.
- Prefer direct returns like `return fmt.Errorf("...: %w", err)` over `wrappedErr := ...; return wrappedErr` when the error is only returned once.
- If the same wrapped error is used twice or more in the same block, prefer a local error variable.
- A log-plus-return pair is enough reason to create a temporary error variable.
- Keep `return err` when there is no meaningful context to add.
- Decide the return shape independently from the log shape.
- Changing `logger.Error(fmt.Errorf(...).Error())` to `logger.Error(err.Error())` is not a reason by itself to change `return fmt.Errorf(...: %w, err)` into `return err`.
- Keep the boundary wrap on return when the current function contributes caller-local context such as `user_id`, route intent, or operation name.
- Avoid `%w` at HTTP handlers, `main`, or other top-level response/logging boundaries.
- Use `%v` or translate to your own exported error when callers should not depend on an underlying dependency error type.
- Do not double-wrap errors that already carry the needed local context.
- At log-only boundaries, prefer `logger.Error(err.Error())` only when the returned error is already known to carry useful wrapped context.
- If that is not known, prefer `logger.Error(fmt.Errorf("...: %w", err).Error())`.
- Avoid double-logging across layers. Helpers usually wrap and return; the handling boundary usually emits the log.

8. Validate before finishing.
- Run targeted scans for anti-patterns.
- Run Go tests for changed packages if possible.
- Report any pre-existing failures separately from logging edits.

## PR Review Checklist

When the user asks to "check logging in this PR", review for:

1. Structured fields present, no concatenated string formatting for context.
2. Structured error fields use `"error"` when present; log-only wrapped errors may be emitted as `logger.Error(err.Error())`.
3. Reasonable level choice (`Info/Warn/Error/Debug`).
4. `component` context exists at service/route boundaries.
5. Request ID propagated where middleware context exists.
6. No secret or PII leakage in new logs.
7. Expected failures (`sql.ErrNoRows`, validation, canceled requests, client disconnects) are not promoted to `Error` without a clear reason.
8. Repetitive watcher/polling/stream logs are `Debug`-level, rate-limited, or summarized.
9. Async/background logs carry useful correlation fields (`request_id`, `session_id`, IDs, stream subject) when available.
10. Lifecycle and degraded-mode transitions are logged once at the operational boundary.
11. Boundary logs prefer `logger.Error(err.Error())` only when the returned error is known to already carry useful wrapped context.
12. If that is not known, boundary logs should wrap with `fmt.Errorf(...)` before logging.
13. Boundary log normalization does not accidentally flatten useful wrapped returns into `return err`.
14. Returned errors are wrapped with `%w` only when useful context is added.
15. Unused logger parameters are removed from helpers/services that no longer log.
16. No duplicate logging of the same failure across helper and handler layers.
17. No broad refactors unrelated to logging.

## Useful Commands

Run from repo root:

```powershell
rg -n 'logger\.With\(\"component\",' -S .
rg -n 'componentLogger|[a-z][A-Za-z0-9]*Logger\s*:=' -S . --glob '*.go' --glob '*.templ'
rg -n 'logger\.(Debug|Info|Warn|Error)\(\"[^\"]*\"\s*,\s*\"err\"' -S . --glob '!backup-*'
rg -n 'fmt\.Println|fmt\.Printf|log\.Print|log\.Fatal' -S . --glob '!backup-*'
rg -n 'return .*err$' -S service pages components --glob '*.go' --glob '*.templ'
rg -n 'request_id|status_code|duration_ms|http request completed' -S main.go http_logging_middleware.go service pages components
```

Treat findings in test files separately from production code.
