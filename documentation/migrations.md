# Database migrations

Conorganizer database migrations use [Goose](https://pressly.github.io/goose/). They are applied manually; do not add automatic migrations to application startup, health checks, readiness checks, or systemd startup.

## Create a migration

Install the [Goose CLI](https://pressly.github.io/goose/installation/) and, from the repository root, create a SQL migration:

```console
goose create <brief-description> sql
```

See the [Goose annotations guide](https://pressly.github.io/goose/documentation/annotations/) for the migration-file format.

## Apply migrations on the server

This is a manual maintenance-window procedure. Replace every placeholder with the correct value for the server and environment. Do not introduce shell variables or run the procedure as a copied script.

1. Enable the maintenance page in `/etc/caddy/Caddyfile`, then restart Caddy. Confirm that the public URL shows the maintenance page.

    ```console
    sudoedit /etc/caddy/Caddyfile
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

6. Migrate the main database in the same way.

    ```console
    sudo mv <main-database-path> database/events.db
    sudo chown <operator>:<operator> database/events.db
    goose -env /dev/null -dir migrations sqlite3 database/events.db up
    sqlite3 database/events.db "PRAGMA integrity_check;"
    sudo mv database/events.db <main-database-path>
    sudo chown deploy:www-data <main-database-path>
    sudo systemctl start <main-service>
    sudo systemctl status <main-service>
    ```

7. Disable the maintenance page and restart Caddy. Verify the real public URL before considering the migration complete.

    ```console
    sudoedit /etc/caddy/Caddyfile
    sudo systemctl restart caddy
    ```
