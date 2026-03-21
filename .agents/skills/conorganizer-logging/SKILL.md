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
- Preferred one-use inline example:
```go
logger.Error(fmt.Errorf("failed to stop watcher: %w", err).Error())
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

6. Wrap returned errors deliberately.
- Run this pass after logging edits and before final validation.
- Wrap returned errors with `fmt.Errorf("...: %w", err)` when adding useful local context or crossing a package/system boundary.
- Prefer direct returns like `return fmt.Errorf("...: %w", err)` over `wrappedErr := ...; return wrappedErr` when the error is only returned once.
- If the same wrapped error is used twice or more in the same block, prefer a local error variable.
- A log-plus-return pair is enough reason to create a temporary error variable.
- Keep `return err` when there is no meaningful context to add.
- Avoid `%w` at HTTP handlers, `main`, or other top-level response/logging boundaries.
- Use `%v` or translate to your own exported error when callers should not depend on an underlying dependency error type.
- Do not double-wrap errors that already carry the needed local context.
- At log-only boundaries, prefer contextualizing the error with `fmt.Errorf(...)` and logging `logger.Error(err.Error())`.

7. Validate before finishing.
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
7. Returned errors are wrapped with `%w` only when useful context is added.
8. No broad refactors unrelated to logging.

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
