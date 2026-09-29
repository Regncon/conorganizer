# Database migrations

Conorganizer database migrations use [Goose](https://pressly.github.io/goose/). They are applied manually; do not add automatic migrations to application startup, health checks, readiness checks, or systemd startup.

## Create a migration

Install the [Goose CLI](https://pressly.github.io/goose/installation/) and, from the repository root, create a SQL migration:

```console
goose -dir migrations create <brief-description> sql
```

Without `-dir migrations` (or `GOOSE_MIGRATION_DIR` in a local, git-ignored `.env`), Goose writes the file to the current directory. Goose names the file with a timestamp version, like the existing files in `migrations/`.

See the [Goose annotations guide](https://pressly.github.io/goose/documentation/annotations/) for the migration-file format. Write both a `-- +goose Up` and a `-- +goose Down` section.

Rehearse a new migration on a copy of your local database before you run it against `database/events.db`:

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

This is a manual maintenance-window procedure. Replace every placeholder with the correct value for the server and environment. Do not introduce shell variables or run the procedure as a copied script.

1. Enable the maintenance page in `/etc/caddy/Caddyfile`: in the `program.regncon.no` block, comment out `import conorganizer-main` and uncomment `import conorganizer-maintenance`. Then restart Caddy and confirm that `program.regncon.no` shows the maintenance page. See [Maintenance mode](../configuration-as-code/README.md#maintenance-mode) for details. `main.lekeplassen.regncon.no` and `demo.lekeplassen.regncon.no` do not show the maintenance page.

    ```console
    sudo vi /srv/configuration-as-code-repo/conorganizer/configuration-as-code/stow/caddy/etc/caddy/Caddyfile
    sudo systemctl restart caddy
    ```

2. Stop both application services.

    ```console
    sudo systemctl stop <main-service>
    sudo systemctl stop <demo-service>
    ```

3. Back up both databases. `conorganizer-sqlite-backup` backs up the main database; make a separate SQLite backup of the demo database.

    ```console
    sudo conorganizer-sqlite-backup
    sudo sqlite3 <demo-database-path> ".backup '<demo-backup-path>/pre-migration-<migration-id>-events.db'"
    ```

4. Update the checkout that contains the migrations.

    ```console
    cd <checkout-containing-migrations>
    git pull
    ```

5. Migrate the demo database. Move it into the checkout, where the Goose command expects `database/events.db`, then restore its service ownership before starting the service.

    The app runs SQLite in WAL mode and does not close the database cleanly when systemd stops it, so `events.db-wal` and `events.db-shm` files can be left next to the database. Moving only `events.db` would leave committed data behind in the WAL file. Checkpoint the database first and confirm that no `-wal` or `-shm` file is left in its directory (see [WAL: never copy the raw database file](sqlite-and-health.md#wal-never-copy-the-raw-database-file)):

    ```console
    sudo sqlite3 <demo-database-path> "PRAGMA wal_checkpoint(TRUNCATE);"
    ls -la <demo-database-directory>
    ```

    Then move and migrate it:

    ```console
    sudo mv <demo-database-path> database/events.db
    sudo chown <operator>:<operator> database/events.db
    goose -env /dev/null -dir migrations sqlite3 database/events.db up
    sqlite3 database/events.db "PRAGMA integrity_check;"
    sudo mv database/events.db <demo-database-path>
    sudo chown deploy:www-data <demo-database-path>
    sudo systemctl start <demo-service>
    sudo systemctl status <demo-service>
    ```

6. Migrate the main database in the same way, including the WAL checkpoint and check from step 5.

    ```console
    sudo sqlite3 <main-database-path> "PRAGMA wal_checkpoint(TRUNCATE);"
    ls -la <main-database-directory>
    sudo mv <main-database-path> database/events.db
    sudo chown <operator>:<operator> database/events.db
    goose -env /dev/null -dir migrations sqlite3 database/events.db up
    sqlite3 database/events.db "PRAGMA integrity_check;"
    sudo mv database/events.db <main-database-path>
    sudo chown deploy:www-data <main-database-path>
    sudo systemctl start <main-service>
    sudo systemctl status <main-service>
    ```

7. Disable the maintenance page by switching the `program.regncon.no` block back to `import conorganizer-main`, and restart Caddy. Verify the real public URL before considering the migration complete.

    ```console
    sudo vi /srv/configuration-as-code-repo/conorganizer/configuration-as-code/stow/caddy/etc/caddy/Caddyfile
    sudo systemctl restart caddy
    ```

    `/etc/caddy/Caddyfile` is a Stow symlink to `configuration-as-code/stow/caddy/etc/caddy/Caddyfile` in the server checkout, so the commands above edit the checked-in file directly. `sudoedit /etc/caddy/Caddyfile` does not work, because sudoedit refuses to follow symlinks by default. Keep the repository in sync with the edit.

Nothing else runs migrations. `deploy/deploy.sh` and `conorganizer-sqlite-restore` do not run Goose, so a backup from before the migration that is restored into the `restored` environment keeps the old schema.

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
