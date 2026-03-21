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

## Safety and Privacy Rules

Never log:

1. Session/refresh tokens.
2. Cookies or authorization headers.
3. Full request/response bodies by default.
4. Sensitive personal data unless explicitly required.

When context is needed, prefer internal IDs over raw personal fields.

## Migration Rules for Existing Logs

When updating existing logging:

1. Replace `"err"` fields with `"error"`.
2. Convert `fmt.Print*` diagnostic logs in production paths to structured `logger.*` calls.
3. Rename `componentLogger` and `<scope>Logger` locals to `logger` when touching that scope.
4. Keep migration scope focused to touched areas.
5. Avoid unrelated refactors.

## Fast Verification Queries

```powershell
rg -n 'logger\.With\(\"component\",' -S .
rg -n 'componentLogger|[a-z][A-Za-z0-9]*Logger\s*:=' -S . --glob '*.go' --glob '*.templ'
rg -n 'logger\.(Debug|Info|Warn|Error)\(\"[^\"]*\"\s*,\s*\"err\"' -S . --glob '!backup-*'
rg -n 'fmt\.Println|fmt\.Printf|log\.Print|log\.Fatal' -S . --glob '!backup-*'
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
