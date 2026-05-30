# SQLite In Production

## Human developer section

Conorganizer uses SQLite in production.

- Production DB: `/mnt/HC_Volume_103911252/environments/main/database/events.db`
- Production images: `/mnt/HC_Volume_103911252/environments/main/event-images`
- Branch/PR environment root: `/mnt/HC_Volume_103911252/environments/<environment-name>`
- Local server backups: `/mnt/HC_Volume_103911252/backups`

Operational rules:

- Migrations are manual only. The application must not run migrations on startup, `/healthz`, `/readyz`, or systemd startup.
- WAL mode is enabled and required for the app database.
- Foreign keys are enabled and required.
- Do not copy or download the live `events.db` file directly. In WAL mode, live state may also be in `events.db-wal` and `events.db-shm`.
- Use `task download:db` or `task download`; the task creates a temporary SQLite `.backup` snapshot on the server and downloads that snapshot.
- Scheduled SQLite backups use `/usr/local/bin/conorganizer-sqlite-backup`, which also uses SQLite `.backup`.

Health endpoints:

- `/healthz` means the process is alive.
- `/readyz` means startup checks passed and a cheap DB live check succeeds.
- Public health responses are generic. Look at service logs for details.

Basic restore/download notes:

- Restore only from a known-good SQLite backup snapshot, not from a raw live DB copy.
- Stop the app before replacing the production DB file.
- Keep ownership writable by `deploy:www-data`.
- After restore, verify `PRAGMA integrity_check;`, required tables, `/healthz`, and `/readyz`.

## LLM/context section

### Current decisions

The SQLite driver is still `_ "modernc.org/sqlite"`. It was kept because the driver already supports the needed per-connection PRAGMA setup through documented DSN `_pragma` parameters, and there was no strong reason to introduce CGO or switch drivers.

The app opens SQLite through `service.InitDB`, which delegates to `service.InitDBWithConfig`. The default config is production-oriented:

```go
BusyTimeoutMillis: 5000
Synchronous:       "NORMAL"
RequireWAL:        true
MaxOpenConns:      1
MaxIdleConns:      1
```

The DSN uses repeated `_pragma` values:

```text
_pragma=journal_mode(WAL)
_pragma=foreign_keys(ON)
_pragma=busy_timeout(5000)
_pragma=synchronous(NORMAL)
```

This matters because `foreign_keys`, `busy_timeout`, and `synchronous` are connection-scoped in practice. Do not replace this with a single startup `db.Exec("PRAGMA ...")` unless every future connection is configured reliably. `modernc.org/sqlite` documents `_pragma` as being run when each connection opens, which is the simple reliable option for this application.

### WAL

WAL mode is required for file-backed app databases. Startup verifies:

```sql
PRAGMA journal_mode;
```

and fails unless it returns `wal`.

WAL means live database state can involve three files:

```text
events.db
events.db-wal
events.db-shm
```

Raw-copying only `events.db` can produce an incomplete snapshot. Raw-copying all three files while the app is live can still race. Use SQLite `.backup` or the backup API.

### Foreign keys

SQLite foreign key enforcement must be enabled by the application. Startup verifies:

```sql
PRAGMA foreign_keys;
```

and fails unless it returns `1`.

Do not remove this verification. Without it, foreign key declarations can silently stop protecting data.

### Busy timeout

`busy_timeout` is set to 5000ms. This lets the app wait briefly when an external process, manual `sqlite3` session, backup, or another request has a SQLite lock. Do not replace this with application retry loops unless there is a measured problem.

### Synchronous mode

`synchronous=NORMAL` is deliberate. In WAL mode, it is the common small-web-app tradeoff: good durability behavior with less fsync pressure than `FULL`. Use `FULL` only if a future operator explicitly accepts the write-latency cost for stronger power-loss durability semantics.

### Connection pool

The app uses:

```go
db.SetMaxOpenConns(1)
db.SetMaxIdleConns(1)
```

Reasoning:

- The app is small, with roughly 200 total users.
- Current write transactions are short.
- Existing image upload work writes files outside DB transactions.
- Current query scans generally close `rows`.
- A single in-process DB connection serializes app writes and avoids multiple app connections competing for SQLite's single writer.
- External `sqlite3`, backup, and snapshot processes can still connect. WAL plus `busy_timeout` is the compatibility mechanism.

Do not increase this casually. If future load testing shows unacceptable latency, consider a deliberate reader/writer split or a small pool only after auditing every connection-scoped PRAGMA and every long-running query.

### No automatic migrations

Migrations are manual only.

Do not add migrations to:

- app startup
- `/healthz`
- `/readyz`
- systemd units
- deploy scripts without explicit human migration steps
- any startup or readiness check

Startup may validate that expected schema exists. It must not create or alter tables. Current startup validation checks core tables:

```text
users
events
billettholdere
puljer
```

If a table is missing, startup fails so the operator can run the correct manual migration or restore process.

### Startup checks

Startup DB checks happen in `service.InitDB`:

- DB path is non-empty.
- DB directory exists.
- DB file exists. Missing DB is fatal; the app must not silently create an empty production DB.
- SQLite opens and pings.
- WAL is enabled and verified.
- foreign keys are enabled and verified.
- busy timeout is configured and verified.
- synchronous mode is configured and verified.
- required core tables exist.

Image directory checks happen once at startup through `service.CheckWritableDirectory`:

- path is non-empty
- path exists
- path is a directory
- app user can create and remove a hidden `.conorganizer-write-check-*` temp file

Do not move image directory writability checks into `/readyz`; it is static startup state.

### Degraded state and readiness

The app still has degraded-mode routing for non-DB startup problems such as an unusable image directory or route setup failure.

DB open/configuration/schema failure is fatal. The process exits instead of serving a degraded page with no usable database.

When degraded mode is entered, `readinessState.MarkDegraded` caches the failure and logs it. Public HTTP responses stay generic.

`/healthz`:

- cheap
- no DB query
- returns `200 ok` when the process can serve HTTP

`/readyz`:

- returns `503 not ready` if startup marked the app degraded
- otherwise performs only `SELECT 1` with a short timeout
- does not run migrations
- does not run schema-changing PRAGMAs
- does not expose internal paths or error strings in the response body

### Safe database download

`task download:db` must not `scp` the live production DB file directly.

Current flow:

1. SSH to the server.
2. Create `/mnt/HC_Volume_103911252/backups/tmp`.
3. Run SQLite `.backup` from the live production DB to a uniquely named temporary snapshot.
4. Run `PRAGMA quick_check;` on the snapshot.
5. Download the snapshot to `database/events.db.download`.
6. Move it into `database/events.db`.
7. Remove the temporary server snapshot via a local shell trap.

This intentionally uses a separate temp directory from the scheduled backup script, so developer downloads do not conflict with retained local backups.

### Branch environment cloning

Branch/PR environments live under:

```text
/mnt/HC_Volume_103911252/environments/<environment-name>
```

When deploying a non-main environment, `deploy/deploy.sh` creates the branch DB with SQLite `.backup` from main if the branch DB file does not already exist. It does not copy the live main database directory.

### Future load testing

Do load testing as a separate task after these production SQLite settings are deployed.

Focus on:

- no 500s
- no panics
- no `database is locked` or `database is busy` errors
- acceptable latency for normal read paths
- acceptable latency for short write paths

Do not run write-path load tests against production data. Possible future tools: k6 or vegeta.

### Things not to change casually

- Do not switch away from `modernc.org/sqlite` without a documented reason.
- Do not add automatic migrations.
- Do not raw-copy a live WAL-mode DB.
- Do not remove WAL or foreign-key startup verification.
- Do not increase the connection pool without reviewing connection-scoped PRAGMAs.
- Do not expose readiness failure details in public HTTP response bodies.

### Primary references

- `modernc.org/sqlite` package docs: https://pkg.go.dev/modernc.org/sqlite
- SQLite PRAGMA docs: https://www.sqlite.org/pragma.html
- SQLite WAL docs: https://www.sqlite.org/wal.html
- SQLite Online Backup API docs: https://www.sqlite.org/backup.html
