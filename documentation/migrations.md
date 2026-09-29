# Database migrations

Conorganizer database migrations use [Goose](https://pressly.github.io/goose/). They are applied manually; do not add automatic migrations to application startup, health checks, readiness checks, or systemd startup.

The only exception is PR preview environments. On every PR deploy, `deploy/deploy.sh` runs `goose up -allow-missing` against the preview database, which is cloned from main on the first deploy. The main, demo, and restored environments are never migrated by CI; main and demo are migrated with the procedure below. Goose does not re-run a migration that has already been applied, so if a PR edits one, delete that preview's database to make the next deploy clone it again.

## Create a migration

From the repository root, create a SQL migration with the same Goose version that CI uses:

```console
go run github.com/pressly/goose/v3/cmd/goose@v3.28.0 -dir migrations create <brief-description> sql
```

Without `-dir migrations` (or `GOOSE_MIGRATION_DIR` in a local, git-ignored `.env`), Goose writes the file to the current directory. Goose names the file with a timestamp version, like the existing files in `migrations/`.

See the [Goose annotations guide](https://pressly.github.io/goose/documentation/annotations/) for the migration-file format. Write both a `-- +goose Up` and a `-- +goose Down` section.

Rehearse a new migration on a copy of your local database before you run it against `database/events.db`. In the commands below and under Troubleshooting, `goose` means an installed Goose v3.28.0 or `go run github.com/pressly/goose/v3/cmd/goose@v3.28.0`:

```console
sqlite3 database/events.db ".backup '/tmp/events-copy.db'"
goose -env /dev/null -dir migrations sqlite3 /tmp/events-copy.db up
goose -env /dev/null -dir migrations sqlite3 /tmp/events-copy.db down
```

After migrating `database/events.db`, run `go tool task test`. It regenerates `schema.sql` from `database/events.db` before running the tests, and the tests build their databases from `schema.sql` (`service/testdb.go`). A migration with tricky data handling can also get its own Go test that runs the Up and Down SQL against a hand-built legacy schema, as `service/puljefordeling/migration_roles_test.go` does.

## Conventions

Most lookup and enum-like values stored in SQLite are capitalized strings, and some of them are Norwegian. Match the Go constants in `models/`.

| Column | Values |
| --- | --- |
| `event_statuses.status` / `events.status` | `Kladd`, `Innsendt`, `Godkjent`, `Forkastet`, `Annonsert` |
| `events_types.event_type` | `Roleplay`, `Boardgame`, `Cardgame`, `Other` |
| `age_groups.age_group` | `Default`, `ChildFriendly`, `AdultsOnly` |
| `event_runtimes.runtime` | `Normal`, `ShortRunning`, `LongRunning` |
| `interest_levels.interest_level` | `Veldig interessert`, `Middels interessert`, `Litt interessert` |
| `pulje_statuses.status` / `puljer.status` | `Open`, `Locked`, `Completed` |
| `relation_events_players.role` | `Player`, `GM` |
| `relation_billettholder_emails.kind` | `Ticket`, `Associated`, `Manual` |

The convention is not universal: `relation_events_players.source` uses lowercase `manual` and `solver`. Use capitalized values for new enum-like columns. For the meaning of the event statuses, see [domeneordbok.md](../domeneordbok.md#arrangementstatus).

## Changing CHECK constraints

SQLite cannot alter a `CHECK` constraint, so a migration that changes one must recreate the table: create `<table>_new` with the new constraint, copy the rows, drop the old table and rename `<table>_new` to the old name.

If views reference the table, the rename fails while the old table is gone, because SQLite validates the views during `ALTER TABLE ... RENAME`. `puljer`, for example, is used by `v_events_by_pulje_active` and `v_event_puljer_active`. Wrap the swap in these PRAGMAs and set them back afterwards, in both the Up and the Down section:

```sql
PRAGMA foreign_keys = OFF;
PRAGMA legacy_alter_table = ON;

-- CREATE TABLE puljer_new (...), INSERT INTO puljer_new SELECT ... FROM puljer,
-- DROP TABLE puljer, ALTER TABLE puljer_new RENAME TO puljer

PRAGMA foreign_keys = ON;
PRAGMA legacy_alter_table = OFF;
```

See `migrations/20260522120000_pulje_status_open_locked_completed.sql` and `migrations/20260524100000_capitalize_pulje_statuses.sql`. `legacy_alter_table` is what makes the rename work with the views. `PRAGMA foreign_keys` has no effect inside a transaction, and Goose runs each migration in one; it still matters when the SQL is run outside Goose, for example with `sqlite3` or from a Go test.

## Apply migrations on the server

Every push to main deploys Goose and `migrations/*.sql` next to each fixed environment's binary in `/opt/conorganizer/<environment>/`, but never runs them. The server scripts below use those deployed files, so no repository checkout is needed. The `restored` environment is not migrated.

1. Enable the maintenance page for `program.regncon.no`. `main.lekeplassen.regncon.no` and `demo.lekeplassen.regncon.no` do not show it. See [Maintenance mode](../configuration-as-code/README.md#maintenance-mode).

    ```console
    sudo conorganizer-maintenance-mode on
    ```

2. Merge the pull request and wait for CI to deploy main. Until the migration has run, the new binary runs against the old schema, which is why the public site is in maintenance mode.

3. Back up the main database.

    ```console
    sudo conorganizer-sqlite-backup
    ```

4. Migrate the demo and main databases. The script backs up demo to `events.db.pre-migrate-<timestamp>` next to its database, migrates demo, and then migrates main, running Goose as the `deploy` user with `up -allow-missing` and a `PRAGMA quick_check` after each. It stops at the first failure, so a failing demo migration never reaches main.

    ```console
    sudo conorganizer-sqlite-migrate
    ```

5. Verify the migration on `https://main.lekeplassen.regncon.no` and `https://demo.lekeplassen.regncon.no`, which are not affected by maintenance mode.

6. Disable the maintenance page and verify `https://program.regncon.no`.

    ```console
    sudo conorganizer-maintenance-mode off
    ```

Apart from PR previews, nothing else runs migrations. `deploy/deploy.sh` only runs Goose for preview databases, and `conorganizer-sqlite-restore` does not run it at all, so a backup from before the migration that is restored into the `restored` environment keeps the old schema.

## Troubleshooting

### `goose up` says "no next version found"

If `goose up` reports `no next version found` but the database is clearly behind the files in `migrations/`, check whether `goose_db_version` is missing its baseline row (`version_id` 0):

```console
sqlite3 database/events.db "SELECT id, version_id, is_applied, tstamp FROM goose_db_version ORDER BY id;"
```

Goose reports this error when the table exists but has no applied rows. A database created from `initialize.sql` has the table but no rows.

To repair:

1. Back up the database with `.backup` (for example to `/tmp`) and rehearse the repair on a copy first.
2. Insert the missing baseline row:

    ```console
    sqlite3 database/events.db "INSERT INTO goose_db_version (version_id, is_applied) VALUES (0, 1);"
    ```

3. Run the pending migrations in order, for example one at a time with `goose -env /dev/null -dir migrations sqlite3 database/events.db up-by-one`.
4. If a migration's tables or columns already exist because they were created outside Goose, do not run its SQL. Mark only that version as applied, then continue with the rest:

    ```console
    sqlite3 database/events.db "INSERT INTO goose_db_version (version_id, is_applied) VALUES (<migration-version>, 1);"
    ```

In this database `goose_db_version.tstamp` is `TEXT` with the default `strftime('%Y-%m-%dT%H:%M:%fZ', 'now')`. That can make `goose status` fail to scan the timestamp even when `goose up` works.
