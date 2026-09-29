# SQLite, startup checks and health endpoints

How the app opens SQLite, what it checks at startup, what `/healthz` and `/readyz` report, and how to copy a live database safely. For where each environment's database lives, see [deployment.md](deployment.md).

## SQLite driver and PRAGMAs

The app uses the pure-Go `modernc.org/sqlite` driver. It needs no CGO and supports setting PRAGMAs on each connection through the DSN.

`service.InitDB` builds a `file:` URI DSN with repeated `_pragma` parameters:

- `journal_mode(WAL)`
- `foreign_keys(ON)`
- `busy_timeout(5000)`
- `synchronous(NORMAL)`

The driver applies them to every new pooled connection. Do not replace this with a one-time `db.Exec("PRAGMA ...")`. `foreign_keys`, `busy_timeout` and `synchronous` apply per connection, so a one-time exec would only configure whichever connection happened to run it.

After opening, `InitDB` reads the PRAGMAs back and fails if `foreign_keys` is not `1`, `journal_mode` is not `wal`, `busy_timeout` is below 5000 ms, or `synchronous` is not `NORMAL`. `TestInitDBAppliesProductionSQLiteSettings` in `service/database_test.go` covers this.

## Connection pool

`DefaultSQLiteConfig` sets `MaxOpenConns` and `MaxIdleConns` to 1, and `normalizeSQLiteConfig` never lets idle exceed open. This is deliberate: the write load is small, and one connection means there is only one SQLite writer inside the process.

Other processes can still use the database alongside the app, for example a `sqlite3` CLI session, the scheduled backup, `conorganizer-export-db` or `deploy.sh` cloning a preview. They are separate OS processes, not pool connections, and WAL mode plus `busy_timeout` let them work safely next to the app.

## Startup checks

Checks whose result cannot change while the app runs are done once at startup.

### Database: fail fast

`main.go` calls `service.InitDB` and exits with `os.Exit(1)` on any error. The app stops if:

- the database path is empty, or the database file or its directory does not exist (no empty database is created)
- the ping fails
- `foreign_keys` is not enabled, `journal_mode` is not WAL, or `busy_timeout` or `synchronous` do not match (`RequireWAL` is on in the default config, so this applies to every environment, including local runs)
- one of the required tables `users`, `events`, `billettholdere` or `puljer` is missing

The goal is to catch a broken SQLite setup early, instead of running quietly with a bad setup. Nothing is cached for these checks, because the process does not start if they fail.

Authentication setup (`authctx.NewSessionValidator`) also stops the process if it fails.

### Image directory and route setup: degraded mode

Other startup problems do not stop the process. The app starts in degraded mode instead:

- **Image directory.** `service.CheckWritableDirectory` checks that the `-image-path` directory is set, exists and is a directory. It then creates a hidden temp file with `os.CreateTemp(path, ".conorganizer-write-check-*")`, closes it and removes it. The check fails if cleanup fails, so it never leaves a file behind. If the check fails, route setup is skipped.
- **Route setup.** If `setupRoutes` fails, the app switches to degraded mode. `setupRoutes` also starts embedded NATS and creates the live buckets.

Each failure is recorded once in `readinessState` with `MarkDegraded`. The full error, including paths, goes to the log. In degraded mode `/` and every unmatched route return a generic 503 page ("Conorganizer is temporarily unavailable"). `/healthz`, `/readyz` and the public assets are still served.

## `/healthz` and `/readyz`

Both routes are mounted on the root router and do not go through auth middleware. Both return `text/plain; charset=utf-8`.

`/healthz` always returns `200` with body `ok\n`. It only shows that the process is alive and checks nothing else.

`/readyz` first reads the cached startup state, then runs one cheap live query, `SELECT 1`, with a 250 ms timeout. It returns `200 ok\n` when both pass. Otherwise it returns `503` with `not ready: <reason>\n`:

| Reason | Cause |
| --- | --- |
| `image directory not writable` | The image directory startup check failed. |
| `application startup incomplete` | `setupRoutes` failed. |
| `multiple startup checks failed` | Startup recorded more than one different reason. |
| `database not available` | The live `SELECT 1` failed or timed out. |

The body only ever contains these fixed reasons. Internal errors and paths are logged, never written to the response.

`/readyz` does not check NATS or JetStream while the app runs. It makes no connection check, account-info call or KV write. If NATS fails after startup, JetStream KV writes can fail while `/readyz` still returns `200`. There is currently no NATS-aware readiness check. For how NATS is used, see [live-update-lifecycle.md](live-update-lifecycle.md). The embedded NATS store dir is set with `-nats-store-dir` (default `data/nats`; `/run/conorganizer-<name>/nats` on the server).

The tests are in `health_test.go` and `service/startup_checks_test.go`.

## Router middleware order

chi panics if `router.Use(...)` is called after routes have been mounted on that router. `startServer` in `main.go` is therefore ordered like this:

1. `router.Use(...)` once, first, with the shared middleware: `middleware.RequestID`, request logging, `middleware.Recoverer` and `requestctx.BillettholderSelectionMiddleware`.
2. Mount `/healthz` and `/readyz`, the dev-reload route (only in `dev` builds), then `/static/*` and `/event-images/*` (`mountPublicAssetRoutes`) on the root router.
3. Create `appRouter := router.With(authMiddleware)` and pass it to `setupRoutes` for the app routes.

This means health checks and public assets never go through Descope session validation or refresh. The shared middleware from step 1 still applies to them. Never add auth or app middleware with `router.Use` after mounting. Use `router.With` or a sub-router instead. `TestPublicAssetRoutesBypassAppMiddleware` in `health_test.go` checks that assets bypass app middleware while app routes still use it.

Assets used to go through auth. A slow auth refresh could then hold up `/static/datastar.js`, and the admin page's placeholder text ("If you are seeing this message, please clear your cookies and refresh the page.") stayed on screen for several seconds until Datastar loaded.

## WAL: never copy the raw database file

In WAL mode the live state can be split across `events.db`, `events.db-wal` and `events.db-shm`. Copying only `events.db` can silently give you an incomplete or stale database. Never copy the live file directly. Always take a SQLite `.backup` snapshot:

```bash
sqlite3 /path/to/events.db ".backup '/path/to/snapshot.db'"
```

Everything that copies a database already does this. The scheduled backups, `deploy.sh` when it clones main into a new PR preview, and `conorganizer-export-db` all use `.backup`, and they check the copy with `PRAGMA quick_check` (the scheduled backup uses `PRAGMA integrity_check`).

## Downloading a database: `conorganizer-export-db`

The local download commands (`go tool task download:main`, `download:demo`, and the `:db`/`:images` variants) are described in [README.md](../README.md#get-the-latest-database-backup-and-images). They run `scripts/download-environment`, which uses `conorganizer-export-db` on the server.

`conorganizer-export-db main|demo` is installed into `/usr/local/bin` by the `scripts` Stow package. It:

1. creates a temp dir with `mktemp -d` and makes a `.backup` snapshot of `/mnt/HC_Volume_103911252/environments/<env>/database/events.db` there
2. for **main only**, anonymizes the snapshot with `/usr/local/share/conorganizer/anonymize-export.sql` and fails if that file is missing (demo is not anonymized)
3. runs `PRAGMA quick_check`
4. writes `events.db` as a tar stream to stdout

Logs go to stderr, and the temp dir is removed on exit.

Locally, `scripts/download-environment` opens a single SSH ControlMaster connection (`ControlPersist=60`), so you are asked for a password at most once per run. It extracts the tar, runs `quick_check` again, removes the old local database and its `-wal` and `-shm` files, and moves the new file into place. Main is saved as `database/events.db` and demo as `database/events-demo.db`.

### What the main anonymization keeps and rewrites

Admin data is kept. An admin is a user with `users.is_admin = 1`. The script keeps:

- admin users' emails
- names and emails of billettholdere linked to an admin
- emails on events owned by an admin

Everything else is rewritten:

- Every other email in `users`, `relation_billettholder_emails` and `events` becomes a deterministic `user_NNNNN@example.invalid`. The same original email gets the same replacement in every table.
- Billettholdere not linked to an admin get `first_name` `User` and a deterministic 6-digit `last_name`.
- `host_name` on all events becomes `Host <6 digits>`, derived from the event's `user_id`, or `Host` when there is no user.
- `phone_number` on all events becomes `00000000`.

`users.external_id` is never changed, not even for non-admins. It is the auth identity, and changing it would break login and user lookup.
