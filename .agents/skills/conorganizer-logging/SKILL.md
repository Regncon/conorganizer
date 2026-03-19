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
3. `main.go` for wiring (`slog.SetDefault`, request ID middleware, component logger usage).
4. `references/project-logging-baseline.md` for rules and examples.

## Workflow

1. Identify scope.
- Decide whether the request is:
- Add logging in new code.
- Migrate existing logs.
- Review PR logging quality only.

2. Pick logger ownership and component.
- Reuse passed `*slog.Logger` when available.
- Create a local logger near route/service entry.
- Never use the variable name `componentLogger`.
- Always use scoped naming: `<scope>Logger := logger.With("component", "<component_name>")`.
- Examples: `authLogger`, `eventServiceLogger`, `profileTicketsLogger`, `eventImageUploadLogger`.
- Keep component names short and stable (`main`, `http`, `auth`, `event_service`, `profile_tickets`).

3. Apply structured log rules.
- Use `logger.Info/Warn/Error/Debug` with message plus key/value pairs.
- Use `"error"` as the error field key (not `"err"`).
- Prefer stable context keys: `event_id`, `user_id`, `pulje_id`, `ticket_id`, `billettholder_id`, `path`, `request_id`.
- Keep `msg` short and action-oriented.

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

6. Validate before finishing.
- Run targeted scans for anti-patterns.
- Run Go tests for changed packages if possible.
- Report any pre-existing failures separately from logging edits.

## PR Review Checklist

When the user asks to "check logging in this PR", review for:

1. Structured fields present, no concatenated string formatting for context.
2. `"error"` key used consistently.
3. Reasonable level choice (`Info/Warn/Error/Debug`).
4. `component` context exists at service/route boundaries.
5. Request ID propagated where middleware context exists.
6. No secret or PII leakage in new logs.
7. No broad refactors unrelated to logging.

## Useful Commands

Run from repo root:

```powershell
rg -n 'logger\.With\(\"component\",' -S .
rg -n 'logger\.(Debug|Info|Warn|Error)\(\"[^\"]*\"\s*,\s*\"err\"' -S . --glob '!backup-*'
rg -n 'fmt\.Println|fmt\.Printf|log\.Print|log\.Fatal' -S . --glob '!backup-*'
rg -n 'request_id|status_code|duration_ms|http request completed' -S main.go http_logging_middleware.go service pages components
```

Treat findings in test files separately from production code.
