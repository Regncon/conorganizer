# Database migrations

Conorganizer database migrations use [Goose](https://pressly.github.io/goose/). They are applied manually; do not add automatic migrations to application startup, health checks, readiness checks, or systemd startup.

The only exception is PR preview environments. On every PR deploy, `deploy/deploy.sh` runs `goose up -allow-missing` against the preview database, which is cloned from main on the first deploy. The main, demo, and restored environments are never migrated by CI; main and demo are migrated with the procedure below. Goose does not re-run a migration that has already been applied, so if a PR edits one, delete that preview's database to make the next deploy clone it again.

## Create a migration

From the repository root, create a SQL migration with the same Goose version that CI uses:

```console
go run github.com/pressly/goose/v3/cmd/goose@v3.28.0 -dir migrations create <brief-description> sql
```

See the [Goose annotations guide](https://pressly.github.io/goose/documentation/annotations/) for the migration-file format.

## Apply migrations on the server

Every push to main deploys Goose and `migrations/*.sql` next to each fixed environment's binary in `/opt/conorganizer/<environment>/`, but never runs them. The server scripts below use those deployed files, so no repository checkout is needed. The `restored` environment is not migrated.

1. Enable the maintenance page for `program.regncon.no`.

    ```console
    sudo conorganizer-maintenance-mode on
    ```

2. Merge the pull request and wait for CI to deploy main. Until the migration has run, the new binary runs against the old schema, which is why the public site is in maintenance mode.

3. Back up the main database.

    ```console
    sudo conorganizer-sqlite-backup
    ```

4. Migrate the demo and main databases. The script backs up demo to `events.db.pre-migrate-<timestamp>` next to its database, migrates demo, and then migrates main. It stops at the first failure, so a failing demo migration never reaches main.

    ```console
    sudo conorganizer-sqlite-migrate
    ```

5. Verify the migration on `https://main.lekeplassen.regncon.no` and `https://demo.lekeplassen.regncon.no`, which are not affected by maintenance mode.

6. Disable the maintenance page and verify `https://program.regncon.no`.

    ```console
    sudo conorganizer-maintenance-mode off
    ```
