# Configuration as Code

Server configuration for the Conorganizer production server.

This directory uses GNU Stow to symlink configuration files into the root filesystem.

For environments, deployment and CI, see [documentation/deployment.md](../documentation/deployment.md).
For the Grafana dashboards and observability notes, see [grafana-dashboards/README.md](grafana-dashboards/README.md).

## Production stack overview

Production runs on one Hetzner Ubuntu VPS. Caddy is the reverse proxy in front of several app instances, each with its own systemd unit, SQLite database, event-images directory and embedded NATS store:

| Instance | systemd unit | Port | Data directory | NATS store |
| --- | --- | --- | --- | --- |
| main | `conorganizer-main.service` | `19080` | `/mnt/HC_Volume_103911252/environments/main/` | `/run/conorganizer-main/nats` |
| demo | `conorganizer-demo.service` | `19081` | `/mnt/HC_Volume_103911252/environments/demo/` | `/run/conorganizer-demo/nats` |
| restored | `conorganizer-restored.service` | `19082` | `/mnt/HC_Volume_103911252/environments/restored/` | `/run/conorganizer-restored/nats` |

Each data directory contains `database/events.db` and `event-images/`. The units run as `deploy:www-data` from `/opt/conorganizer/<env>/`. Main and restored read `/etc/conorganizer/env`; demo reads `/etc/conorganizer/demo-env`. These environment files are not in the repository.

Local backups are written to `/mnt/HC_Volume_103911252/backups/sqlite` and `/mnt/HC_Volume_103911252/backups/images` (see [Backups](#backups)).

Observability services:

| Service | Address | Notes |
| --- | --- | --- |
| Grafana | `127.0.0.1:3400` | Public at `grafana.regncon.no` through Caddy |
| Loki | `127.0.0.1:3500` | Not the Loki default `3100`. Data in `/var/lib/loki`, kept for 120 days (see [Loki storage](#loki-storage)) |
| Promtail | `:9080` | No bind address set. Tails `/var/log/messages` and pushes to Loki |
| Prometheus | `127.0.0.1:9090` | OS package, config from this repo |
| node_exporter | `127.0.0.1:9100` | OS package, config from this repo |
| blackbox_exporter | `127.0.0.1:9115` | OS package, config from this repo |

Grafana Alloy may be installed on the server, but no Alloy config is in the repository. Prometheus only scrapes Alloy's self-metrics on `127.0.0.1:12345`.

## Repository location

The production repository is stored at:

```text
/srv/configuration-as-code-repo/conorganizer/configuration-as-code
```

## Stow packages

| Package | Installs |
| --- | --- |
| `scripts` | `/usr/local/bin/conorganizer-*` and `/usr/local/share/conorganizer/anonymize-export.sql` |
| `caddy` | `/etc/caddy/Caddyfile` |
| `grafana` | `/etc/grafana/grafana.ini` |
| `loki` | `/etc/loki/config.yml` |
| `prometheus` | `/etc/prometheus/prometheus.yml`, `/etc/prometheus/blackbox.yml` and `/etc/default/prometheus-node-exporter` |
| `promtail` | `/etc/promtail/config.yml` |
| `systemd` | `/etc/systemd/system/` units for the app instances, the backup services and timers, `loki.service`, its `loki.service.d` drop-in and `promtail.service` |

Prometheus, `prometheus-node-exporter` and `prometheus-blackbox-exporter` are installed as OS packages; only their config comes from the `prometheus` package. `/etc/default/prometheus-node-exporter` is the only `/etc/default` override in the repository. It binds node_exporter to `127.0.0.1:9100`, enables `--collector.systemd`, and uses a filesystem mount-point exclude that does not exclude `/mnt`, so the mounted volume is visible.

## Install

From the repository root:

```bash
./configuration-as-code/install.sh
```

The script runs `sudo stow --dir=stow --target=/ --restow <package>` for these packages, in order: `scripts`, `caddy`, `grafana`, `loki`, `prometheus`, `promtail`, `systemd`. A missing package directory prints `Skipping missing Stow package: <name>` instead of failing. It then runs `sudo systemctl daemon-reload`.

The script does not enable, start or restart any service. Restart the affected services yourself after a config change, and enable new units or timers with `sudo systemctl enable --now <unit>`.

## Fix permissions

```bash
sudo find configuration-as-code/stow -type d -exec chmod 755 {} \;
sudo find configuration-as-code/stow -type f -exec chmod 644 {} \;
sudo find configuration-as-code/stow/scripts/usr/local/bin -type f -exec chmod 755 {} \;
```

## Database scripts

The `scripts` Stow package installs the database maintenance commands in
`/usr/local/bin`.

- `conorganizer-sqlite-backup` creates the scheduled compressed main-database backups.
- `conorganizer-images-backup` creates the scheduled main event-images archive.
- `conorganizer-sqlite-restore events-YYYYMMDDTHHMMSSZ.db.zst` installs a selected backup into the public `restored` environment and refreshes its event images from main.
- `conorganizer-export-db main|demo` streams a tar with a SQLite copy of that environment's database to stdout. The main export is anonymized with `/usr/local/share/conorganizer/anonymize-export.sql`. It is used by `scripts/download-environment` (`go tool task download:main` / `download:demo`).
- `conorganizer-sqlite-migrate` backs up the demo database, then runs the deployed Goose migrations against demo and main. See [documentation/migrations.md](../documentation/migrations.md).
- `conorganizer-maintenance-mode on|off` shows or hides the maintenance page for `program.regncon.no`. Without an argument it prints the current state.

## Backups

Backups of the main environment run from systemd timers:

| Timer | Schedule | Output |
| --- | --- | --- |
| `conorganizer-sqlite-backup.timer` | `OnCalendar=*:0/15` (every 15 minutes) | `/mnt/HC_Volume_103911252/backups/sqlite/events-<timestamp>.db.zst`, plus an uncompressed `events-latest.db` |
| `conorganizer-images-backup.timer` | `OnCalendar=*-*-* 03:30:00` (daily) | `/mnt/HC_Volume_103911252/backups/images/event-images-<timestamp>.tar.zst` |

Both timers use `Persistent=true`, so a run missed while the server was down happens when the timer starts again. Each timer starts the matching oneshot `.service`, which runs the script with the same name in `/usr/local/bin`. Timestamps are UTC (`YYYYMMDDTHHMMSSZ`). Compressed backups older than 14 days are deleted.

The SQLite backup runs `PRAGMA integrity_check` and checks that the `users` and `events` tables exist before it keeps the file.

Log lines start with `conorganizer-sqlite-backup:` or `conorganizer-images-backup:` and end with `completed successfully` on success. The Grafana backup panels depend on these strings.

Before logging `completed successfully`, each script writes Prometheus metrics to `/var/lib/prometheus/node-exporter/conorganizer-sqlite-backup.prom` or `conorganizer-images-backup.prom`, which node_exporter's textfile collector serves. They hold the time of the last successful backup and, from the SQLite backup, the database, WAL, backup and backup directory sizes. See [Backup Metrics](grafana-dashboards/README.md#backup-metrics).

## Loki storage

Loki stores its chunks, index and write-ahead log in `/var/lib/loki` on the root disk. The drop-in `loki.service.d/state-directory.conf` creates the folder for the `loki` user, with mode 0700 (`StateDirectory=loki`). It is a drop-in because the Loki package overwrites `/etc/systemd/system/loki.service` on upgrade, which would undo changes to the unit itself.

The compactor deletes logs older than `retention_period` (120 days). Without `compactor.retention_enabled`, Loki ignores `retention_period` and keeps everything. Enabling retention also enables Loki's log-deletion API, without authentication. `deletion_mode: disabled` turns that API off, and retention does not need it.

On `/tmp`, Loki grew about 15 MB a day over its first 30 hours. It grows less on quiet days and more while people run queries, because about half of what Loki stores is its own query logs. At 15 MB a day, 120 days comes to roughly 2 GB. That takes `/` from about 84% to about 89% used, just under the 90% red line of the "Root filesystem usage" panel, and no alert watches it. Run `sudo du -sh /var/lib/loki` a week after the move, and lower `retention_period` if Loki grows faster than 15 MB a day.

Loki used to keep this data in `/tmp/loki`, because the repository's config started from the example file the Loki package installs (`cmd/loki/loki-local-config.yaml`). On this server `/tmp` is a 1 GB tmpfs in RAM, so every reboot wiped all logs. The originals in the systemd journal and `/var/log/syslog*` were never affected.

### Moving from /tmp/loki

To move without losing the logs collected since the last reboot, run all of this in one go, from `/srv/configuration-as-code-repo/conorganizer`. That checkout belongs to `gerhard` and is kept on a detached `origin/main`, so git runs as `sudo -u gerhard`. Updating the checkout switches `/etc/loki/config.yml` to `/var/lib/loki` at once, because it is a stow symlink into the checkout. If Loki restarted before systemd knew about the drop-in, it could not create its folder and would stay down. So stop Loki before you update the checkout, and don't update it past this change without running these steps.

First check what the update brings, and that the installed Loki accepts the new config:

```bash
cd /srv/configuration-as-code-repo/conorganizer
sudo -u gerhard git fetch origin
sudo -u gerhard git diff --stat HEAD origin/main -- configuration-as-code/stow
sudo -u gerhard git show origin/main:configuration-as-code/stow/loki/etc/loki/config.yml > /tmp/loki-config-new.yml
/usr/bin/loki -config.file=/tmp/loki-config-new.yml -verify-config
```

The diff lists every stowed file the update changes. Changes to the Caddyfile or to other units take effect the next time those services reload.

Then stop Loki, update the checkout, copy the data and start Loki again:

```bash
sudo systemctl stop loki
sudo -u gerhard git checkout --detach origin/main
test -e /var/lib/loki && echo "/var/lib/loki already exists: do not copy, check it first" || sudo cp -a /tmp/loki /var/lib/loki
sudo stow --dir=configuration-as-code/stow --target=/ --restow loki systemd
sudo systemctl daemon-reload
sudo systemctl start loki
```

`cp -a` creates `/var/lib/loki` with the owner of `/tmp/loki`. The Loki package creates no `loki` group, so don't create the folder with `install -g loki`.

If Loki does not start, `/tmp/loki` is untouched. To go back to the old path, run `sudo -u gerhard git checkout HEAD@{1} -- configuration-as-code/stow/loki/etc/loki/config.yml`, then `sudo systemctl start loki`. The drop-in is harmless with the old path and can stay.

Alloy retries each batch 10 times, with backoff growing from 0.5 s to 5 min, so nothing is dropped as long as Loki is back within about 4 minutes. The journal keeps every line either way. Check the result:

```bash
systemctl status loki --no-pager
for i in $(seq 30); do curl -fsS http://127.0.0.1:3500/ready && break; sleep 2; done
curl -fsS http://127.0.0.1:3500/config | grep -E 'path_prefix|working_directory|active_index_directory|retention_enabled|deletion_mode'
sudo du -sh /var/lib/loki
sudo journalctl -u loki -n 50 --no-pager -o cat | grep -E "level=(error|warn)"
curl -fsS http://127.0.0.1:12345/metrics | grep '^loki_write_dropped_entries_total'
```

The last command should show only 0 values.

In Grafana, logs from before the move should still be there. Re-import `grafana-dashboards/conorganizer-production-health.json` (see [Manual Import](grafana-dashboards/README.md#manual-import)), so the "Root filesystem usage" description no longer says Loki's data is in `/tmp/loki`. `/tmp/loki` can then be deleted, or left for the next reboot to clear.

## Caddy hostnames

Hostnames in the checked-in Caddyfile:

| Hostname | Target |
| --- | --- |
| `main.lekeplassen.regncon.no` | `import conorganizer-main` (`127.0.0.1:19080`) |
| `program.regncon.no` | Temporary (302) redirect of every request to `https://regncon.no/gamle-regncon-program/` while the con is over, see [Maintenance mode](#maintenance-mode) |
| `demo.lekeplassen.regncon.no` | `127.0.0.1:19081` |
| `restored.lekeplassen.regncon.no` | `127.0.0.1:19082` |
| `grafana.regncon.no` | `127.0.0.1:3400` |
| `meetup-january-2026.lekeplassen.regncon.no` | `127.0.0.1:8080` |

The Caddyfile ends with `import /etc/caddy/sites-enabled/*.caddy`. That directory is not in the repository. `deploy/deploy.sh` writes one site file per PR preview there. The fixed environments (main, demo, restored) do not use it.

Blackbox probes and the TLS-expiry panels only check `https://main.lekeplassen.regncon.no/` and `https://grafana.regncon.no/`.

## Maintenance mode

The maintenance page is defined entirely in the Caddyfile as the snippet `(conorganizer-maintenance)`. It returns a static bilingual English/Norwegian HTML page with:

- HTTP `503`, which tells browsers and crawlers that the outage is temporary.
- `Cache-Control: no-store`, so the maintenance page is not cached and still shown after the site is back.

While the con is over, `program.regncon.no` redirects every request to `https://regncon.no/gamle-regncon-program/` and does not import either snippet, so the flag file has no effect. To bring the site back, restore the previous `program.regncon.no` block from git history. That block serves the maintenance snippet while the flag file `/var/lib/conorganizer/maintenance.on` exists, and `import conorganizer-main` otherwise. Caddy checks the file on every request, so switching needs no Caddy reload or restart and no edit to the Caddyfile:

```bash
sudo conorganizer-maintenance-mode on    # create the flag file
sudo conorganizer-maintenance-mode off   # remove it
conorganizer-maintenance-mode            # print the current state
```

Only `program.regncon.no` is affected. `main.lekeplassen.regncon.no` and `demo.lekeplassen.regncon.no` keep serving the app. Nothing switches maintenance mode automatically. [documentation/migrations.md](../documentation/migrations.md) uses it for the migration procedure.

## Find all stowed files

List symlinks below `/etc` and `/usr/local` that resolve into this repository:

```bash
sudo find /etc /usr/local -type l -print0 |
while IFS= read -r -d '' symlink_path; do
    resolved_target="$(readlink -f -- "$symlink_path" 2>/dev/null || true)"

    case "$resolved_target" in
        /srv/configuration-as-code-repo/conorganizer/configuration-as-code/stow/*)
            printf '%s -> %s\n' "$symlink_path" "$resolved_target"
            ;;
    esac
done
```
