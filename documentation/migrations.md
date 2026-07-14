# Migrations with Goose

> [!NOTE]
> Goose reads variables from `.env`. Make sure this file is updated with the most recent version from Discord before running any commands.

We're using [Goose](https://pressly.github.io/goose/) in our migration process for its simplicity and reliability. While Goose is available as a Go dependency for programmatic database migrations, we're mostly using its CLI tool for manual updates.

## Running Goose manually

> [!WARNING]
> Before running Goose, run `go tool task download` to fetch the newest version of the database.
> Install the Goose CLI tool from the [official installation guide](https://pressly.github.io/goose/installation/). Afterward, `goose` should be globally available in your terminal.
> Migrations are manual only. Do not add automatic migrations to application startup, health checks, readiness checks, or systemd startup.

To create a new migration file, run this command. See the [Goose annotations guide](https://pressly.github.io/goose/documentation/annotations/) for more annotation examples.

```console
goose create <briefly describe changes> sql
```

After adding migration files, use `up` or `down` to run migrations.

```console
goose up
```

## Pushing migrations to production

> [!CAUTION]
> Make sure that you can do all of the following steps before you start. These actions require Goose and server access.

1. Run Goose on the local database (preferably a copy).
2. Make a backup on the server.
3. Upload the database to the server.

## Step-by-step database update

```bash
systemctl list-units --type=service | grep -i conorganizer
systemctl list-unit-files | grep -i conorganizer
```

Check the service command and find the mounted database path:

```bash
systemctl show INSERT_SERVICE_NAME -p ExecStart --value | fold -s -w 120
```

Look for the host path that contains `events.db` or maps the database folder into the app.

```bash
ls -lh /mnt/HC_Volume_103911252/environments
```

Example: `/mnt/HC_Volume_103911252/environments/1337-merge/database/events.db`

Stop the service:

```bash
sudo systemctl stop INSERT_SERVICE_NAME
```

Back up the current database if needed:

```bash
sqlite3 PATH_TO_DB ".backup 'PATH_TO_BACKUP/events.db.bak'"
```

Move the uploaded database into place:

```bash
mv /path/to/uploaded/events.db /mnt/HC_Volume_103911252/environments/1337-merge/database/events.db
cd /mnt/HC_Volume_103911252/environments/1337-merge/database
sudo chown deploy:deploy events.db
sudo chmod 644 events.db
```

Start the service again:

```bash
sudo systemctl start INSERT_SERVICE_NAME
sudo systemctl status INSERT_SERVICE_NAME
```

Check logs:

```bash
journalctl -u INSERT_SERVICE_NAME -n 100 --no-pager
```

## Restore compose commands

> [!CAUTION]
> Do not run this unless you know all caveats; this can affect production negatively.

```bash
docker compose -f compose-restore.yaml down && docker image rm regncon-migration
docker compose -f compose-restore.yaml up
```
