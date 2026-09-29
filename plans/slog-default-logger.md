# Replace Logger Parameters With slog.Default

**Date:** 2026-09-25
**Status:** Proposed, discuss "Before execution" first

## Scope

Stop passing `logger *slog.Logger` through function and template parameters. Each function that logs gets its own logger from the process default:

```go
logger := slog.Default().With("component", "<component_name>")
```

The log output stays the same: JSON, same levels, same `component` values, same fields.

## Background

`main.go` builds the app logger with `applog.NewJSONLogger()` and makes it the process default:

```go
slog.SetDefault(logger)
```

After that, `slog.Default()` returns the same JSON logger anywhere in the process. This was verified in the running app on 2026-09-25: a log line from `pages/root/billettholder_interests.go` using `slog.Default().With("component", "root")` showed up in `docker compose logs webserver` as normal JSON with `"component":"root"`.

`service/live` and `pages/notfound` already use `slog.Default()`, and `pages/root` does since the interest indicator (#313).

## Problem

- 169 signatures in 73 files take `logger *slog.Logger` (production code only, excluding generated `_templ.go` files and tests).
- Most of them only forward the logger to the next function or template. For example, `layouts.Base(title, userInfo, db, logger, children)` takes a logger so the header menu can log.
- Templates need a logger parameter to log, even though they already have `ctx`.
- No test reads log output. 49 test files create a logger (`testutil.NewTestLogger()` or `testutil.NewSlogAdapter(&testutil.StubLogger{})`) only because the function signature requires one. `StubLogger` records calls in a private field that nothing reads.

## Decision

- Use `logger := slog.Default().With("component", "<component_name>")` at the top of each function or template block that logs.
- Keep the existing component names. The name that was bound with `logger.With("component", ...)` near a route or service entry moves into the function that logs.
- Remove the `logger` parameter once nothing in the function uses it, and update every call site in the same commit.
- Keep `main.go` as the only place that builds the logger and calls `slog.SetDefault`. It must run before anything logs.
- Keep adding `request_id`, `user_id`, `event_id` and similar fields per log call, as today.

### Not part of this plan

- Putting a request-scoped logger on the context (for example `applog.FromContext(ctx)` with `request_id` bound by middleware). It can be added later on top of this, since the fallback would be `slog.Default()`.
- Changing log messages, levels or fields.

## Before execution

Discuss and decide whether there is a shorter way to get a logger, so there is little to remember when writing code. `slog.Default().With("component", "...")` in every function that logs is long, and it is easy to forget the component or spell it differently each time.

Ideas to compare:

1. **Small helper in `service/applog`:** `logger := applog.For("root")`, which returns `slog.Default().With("component", "root")`. Short and hard to get wrong. It is one extra import.
2. **One function per package:** for example `func rootLogger() *slog.Logger { return slog.Default().With("component", "root") }` in each package, so the component name is written once per package.
3. **Log directly with `slog`:** `slog.Error(err.Error(), "component", "root", "user_id", ...)`. No logger variable at all, but the component must be repeated in every call.
4. **Logger on the context:** `applog.FromContext(ctx)`, with `component` and `request_id` bound by middleware or at route setup. Shortest in handlers and templates, and adds `request_id` automatically, but needs middleware and a way to set the component.

Watch out for a package-level variable like `var logger = slog.Default().With(...)`: it is evaluated when the package loads, before `main.go` calls `slog.SetDefault`, so it would keep Go's default text logger instead of the app's JSON logger. Any shortcut must call `slog.Default()` when logging, not at package load.

Whatever is chosen replaces `slog.Default().With("component", "<component_name>")` in the rest of this plan and in the logging skill.

### Enforcing the chosen way

Make CI fail when code gets a logger any other way, so nobody has to remember the rule.

- **Source test (main check):** a Go test, for example `service/applog/usage_test.go`, that reads all `.go` and `.templ` files (skipping `_templ.go`, tests, `tmp/` and `.superpowers/`) and fails on:
  - `slog.Default(` or `slog.New(` outside `service/applog` and `main.go`
  - `logger *slog.Logger` parameters

  The failure message says what to write instead. It runs in the existing `go test ./...`, locally and in `.github/workflows/buildAndTest.yml`, and needs no new tools.
- **Docs:** add the rule to `AGENTS.md` and the `conorganizer-logging` skill, so agents write it correctly from the start.
- **No lint rules:** do not add `forbidigo` or other golangci-lint rules for this. `.golangci.yml` skips generated `_templ.go` files, so lint rules would not cover templates anyway.

Add the source test in the first commit of the migration, with an allow list of the files not migrated yet. Remove files from the allow list as each step lands, and delete the allow list when the migration is done.

## Inventory

Production files with `logger *slog.Logger` parameters, grouped by directory (files / signatures):

| Directory | Files | Signatures |
|---|---:|---:|
| `components/formsubmission` | 8 | 30 |
| `components/formsubmission/event_img_upload` | 1 | 4 |
| `components/event_components` | 1 | 1 |
| `components/header` | 1 | 3 |
| `components/profile` | 1 | 3 |
| `components/ticket_holder` | 2 | 2 |
| `layouts` | 1 | 1 |
| `pages/admin` | 10 | 21 |
| `pages/admin/approval` | 4 | 9 |
| `pages/admin/billettholder_admin` | 8 | 10 |
| `pages/admin/rooms` | 3 | 9 |
| `pages/event` | 3 | 7 |
| `pages/login` | 1 | 4 |
| `pages/notfound` | 1 | 2 |
| `pages/print-friendly` | 1 | 4 |
| `pages/profile` | 2 | 8 |
| `pages/profile/newevent` | 2 | 4 |
| `pages/profile/tickets` | 3 | 7 |
| `pages/root` | 2 | 2 |
| `service/authctx` | 2 | 4 |
| `service/checkIn` | 4 | 9 |
| `service/live` | 1 | 1 |
| `service/userctx` | 1 | 4 |
| root package (`main.go`, `router.go`, `health.go`, `dev_reload.go`, `http_logging_middleware.go`, `static_dev.go`, `static_prod.go`) | 7 | 12 |

`tmp/` and `.superpowers/` also match but are not app code and are out of scope.

Regenerate the inventory before starting:

```bash
grep -rEn "logger \*slog\.Logger" --include=*.go --include=*.templ . | grep -v _templ.go | grep -v _test.go | grep -v -E "^\./(tmp|\.superpowers)/"
```

## Migration order

Work bottom-up so every commit builds and passes tests on its own. A function can only lose its parameter when all its callers are updated in the same commit, so start with the functions that are called by others and have no logger callees left.

1. **Services:** `service/checkIn`, `service/authctx`, `service/userctx`, `service/live` (the `WithLogger` option can go if nothing needs it).
2. **Shared components:** `components/ticket_holder`, `components/event_components`, `components/profile`, `components/formsubmission` (split into several commits, it has 30 signatures), `components/formsubmission/event_img_upload`.
3. **Header and layout:** `components/header`, then `layouts.Base`. `layouts.Base` has many callers, so do it as its own commit.
4. **Pages, one package per commit:** `pages/root`, `pages/event`, `pages/login`, `pages/notfound`, `pages/print-friendly`, `pages/profile/...`, `pages/admin/...`.
5. **Route setup and the root package:** `router.go`, `health.go`, `dev_reload.go`, `static_*.go`, `http_logging_middleware.go`. Keep `main.go` building the logger and calling `slog.SetDefault`.
6. **Test helpers:** remove `testutil.NewTestLogger`, `testutil.StubLogger`, `testutil.NewSlogAdapter`, the `testutil.Logger` interface and `CreateTestDBAndLogger` once no test uses them.

Suggested commit message per step: `Use slog.Default for logging in <package>`.

## Per-function steps

1. Add `logger := slog.Default().With("component", "<component_name>")` where the function logs. Reuse the component name the caller used to bind.
2. Remove the `logger *slog.Logger` parameter.
3. Update all call sites, including generated `_templ.go` callers by running `go tool templ generate`.
4. Remove now-unused `logger` arguments from tests.
5. Run `go build ./...` and `go test ./...`.

## Tests

- Tests that only passed a logger to satisfy a signature drop the argument.
- Code under test now logs through `slog.Default()`. In tests that is Go's default text logger, which writes to stderr. `go test` only shows that output for failing tests or with `-v`, so it does not change normal test runs.
- If a test later needs to check a log line, it can swap the default for the test with `slog.SetDefault(testLogger)` and restore it with `t.Cleanup`. Such a test must not use `t.Parallel()`.

## Docs to update

- `.claude/skills/conorganizer-logging/SKILL.md` and `references/project-logging-baseline.md`: replace "reuse the passed `*slog.Logger`" with "get the logger with `slog.Default().With("component", ...)`", and drop the rules about removing unused logger parameters.

## Verification

- `grep -rEn "logger \*slog\.Logger" ...` (see Inventory) returns nothing outside `main.go`.
- `go build ./...` and `go test ./...` pass after every step.
- In the running app, `docker compose logs webserver` still shows JSON log lines with the expected `component` values for the touched packages.

## Risks

- **Logging before `slog.SetDefault`:** anything that logs during package init or before `main.go` sets the default would use Go's text logger. Keep `slog.SetDefault` as the first thing `main.go` does after building the logger.
- **Lost component names:** when a caller bound `component` and passed the logger down, the callee must bind the same name itself, or log lines lose their `component`. Review the diff for each step with that in mind.
- **Large diff:** 73 files. Keep one package per commit so each step is easy to review and revert.
