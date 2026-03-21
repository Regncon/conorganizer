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
- Use `"error"` as the field key only when keeping a separate structured error field is still useful.

## Safety and Privacy Rules

Never log:

1. Session/refresh tokens.
2. Cookies or authorization headers.
3. Full request/response bodies by default.
4. Sensitive personal data unless explicitly required.

When context is needed, prefer internal IDs over raw personal fields.

## Error Wrapping Rules

When updating returned errors:

1. Wrap with `fmt.Errorf("...: %w", err)` when the current function can add useful local context.
2. Wrap when an error crosses a package, subsystem, or abstraction boundary and the higher layer can explain the failed operation.
3. Keep `return err` when the current function adds no meaningful context.
4. Avoid `%w` at HTTP handlers, `main`, and similar top-level response/logging boundaries; log the accumulated context and return a safe error or response instead.
5. Use `%v` or translate to an exported repo error when callers should not depend on an underlying dependency error type.
6. Do not double-wrap errors that already contain the needed local context.
7. At log-only boundaries, prefer contextualizing the error with `fmt.Errorf(...)` and logging `logger.Error(err.Error())`.

## Migration Rules for Existing Logs

When updating existing logging:

1. Replace `"err"` fields with `"error"`.
2. Convert `fmt.Print*` diagnostic logs in production paths to structured `logger.*` calls.
3. Rename `componentLogger` and `<scope>Logger` locals to `logger` when touching that scope.
4. Run a final pass over returned errors and wrap them with `fmt.Errorf(...: %w)` only when the function can add useful context.
5. At log-only boundaries, prefer `logger.Error(fmt.Errorf(...).Error())` over duplicating the same context in both the message and structured fields.
6. Keep migration scope focused to touched areas.
7. Avoid unrelated refactors.

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
