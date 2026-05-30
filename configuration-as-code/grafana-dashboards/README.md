# Conorganizer Grafana Dashboards

This directory contains importable Grafana dashboard JSON for manual review:

- `conorganizer-production-health.json`
- `conorganizer-main-service-debugging.json`

The files are intentionally outside `configuration-as-code/stow/`. They are not server-side Grafana provisioning and should not be treated as provisioned dashboards until that is explicitly added later.

## Manual Import

1. Open Grafana.
2. Go to **Dashboards** -> **New** -> **Import**.
3. Upload one JSON file from this directory, or paste its JSON.
4. Pick the Prometheus datasource for `DS_PROMETHEUS`.
5. Pick the Loki datasource for `DS_LOKI`.
6. Import the dashboard.
7. Adjust variables if the server labels differ from the defaults.

The JSON defaults datasource variable values to `Prometheus` and `Loki` for normal import/overwrite workflows. It does not contain production Grafana database IDs or secrets.

## Dashboards

`Conorganizer Production Health` answers platform and restore-confidence questions:

- Is the public site reachable?
- Is Grafana reachable?
- Are Caddy and `conorganizer-main.service` active?
- Are recent backups visible?
- Is root and mounted-volume capacity safe?
- Are TLS certificates close to expiry?
- Are there system, backup, Caddy, Loki, Promtail, Alloy, SQLite, or critical log errors?

`Conorganizer Main Service Debugging` is focused on `conorganizer-main.service`:

- Is the service active and publicly reachable?
- Is the current app error signal zero?
- Are there panics, 5xx request logs, or database lock/busy errors?
- What are request rate, status mix, route usage, and latency when JSON request logs are parseable?
- Which recent app logs, failed requests, or slow requests should be inspected next?

## Variables

- `DS_PROMETHEUS`: Prometheus datasource selected at import time.
- `DS_LOKI`: Loki datasource selected at import time.
- `node_instance`: Prometheus query variable using `label_values(node_uname_info, instance)`.
- `loki_job`: Loki query variable using `label_values(job)`. The all value is broad by design; narrow it after import if exact labels exist.
- `main_host`: defaults to `main.lekeplassen.regncon.no`.
- `grafana_host`: defaults to `grafana.regncon.no`.
- `main_service`: defaults to `conorganizer-main.service`.
- `caddy_service`: defaults to `caddy.service`.
- `volume_mountpoint`: defaults to `/mnt/HC_Volume_103911252`.

## Repo Observations

- `configuration-as-code/stow/prometheus/etc/prometheus/prometheus.yml` exists and scrapes local `prometheus`, `node`, `alloy`, and `blackbox_https` jobs.
- The Prometheus blackbox job probes `https://main.lekeplassen.regncon.no/` and `https://grafana.regncon.no/` through `127.0.0.1:9115`.
- `configuration-as-code/stow/prometheus/etc/prometheus/blackbox.yml` defines an `http_2xx` module with `fail_if_not_ssl: true`, redirects enabled, IPv4 preferred, and a 5 second timeout.
- `configuration-as-code/stow/prometheus/etc/default/prometheus-node-exporter` enables `--collector.systemd` and keeps `/mnt` from being excluded by the filesystem collector.
- `configuration-as-code/stow/loki/etc/loki/config.yml` configures Loki on `127.0.0.1:3500` with filesystem storage and `retention_period: 120d`.
- `configuration-as-code/stow/promtail/etc/promtail/config.yml` pushes to `http://localhost:3500/loki/api/v1/push`, labels logs with `job=varlogs`, and scrapes `/var/log/messages`.
- The repo contains Loki and Promtail systemd unit files. It does not contain an Alloy config file, although Prometheus is configured to scrape Alloy self-metrics on `127.0.0.1:12345` if Alloy exists on the server.
- The Caddyfile routes `grafana.regncon.no` to `127.0.0.1:3400`, includes `meetup-january-2026.lekeplassen.regncon.no`, and imports `/etc/caddy/sites-enabled/*.caddy`. The main host config is not visible in the checked-in root Caddyfile and may be in an imported server-local file.
- No Caddy native metrics endpoint or Caddy access-log format is configured in the checked-in Caddyfile.
- `conorganizer-main.service` runs `/opt/conorganizer/main/conorganizer-main` on `PORT=18856`, with the SQLite DB at `/mnt/HC_Volume_103911252/environments/main/database/events.db` and images at `/mnt/HC_Volume_103911252/environments/main/event-images`.
- The app uses JSON `slog` logs to stdout. `LOG_LEVEL` supports `DEBUG`, `INFO`, `WARN`/`WARNING`, and `ERROR`.
- Request logs are emitted with `msg="http request completed"`, `component="http"`, `method`, `path`, `status_code`, `duration_ms`, and optional `request_id`.
- No app Prometheus instrumentation, `/metrics`, `/healthz`, or `/readyz` route was found in the Go code. Public health panels therefore use blackbox probes rather than an app-native health endpoint.
- SQLite backups run every 15 minutes from `conorganizer-sqlite-backup.timer`.
- Image backups run daily at `03:30:00` from `conorganizer-images-backup.timer`; systemd calendar times use the server's local timezone unless configured otherwise.
- Backup scripts emit stable prefixes:
  - `conorganizer-sqlite-backup:`
  - `conorganizer-images-backup:`
- SQLite backup success logs include `conorganizer-sqlite-backup: completed successfully`.
- Image backup success logs include `conorganizer-images-backup: completed successfully`.
- Backup failure strings include `database does not exist`, `image directory does not exist`, `integrity check failed`, and `sanity check failed`.
- Grafana config exists under `configuration-as-code/stow/grafana/etc/grafana/grafana.ini`, but no Grafana datasource or dashboard provisioning is added for these JSON files.
- `configuration-as-code/install.sh` stows `scripts`, `conorganizer` if present, `caddy`, `grafana`, and `prometheus`. It does not import these dashboard JSON files.

## Historical Server Notes

The previous README contained manual server notes dated 2026-05-28. Preserve them as historical context only; verify current state on the server before relying on them:

- Loki was ready on `http://127.0.0.1:3500`.
- Promtail was active and tailed `/var/log/messages` with `job=varlogs`.
- Alloy v1.16.1 was active on `127.0.0.1:12345` with Loki journal collection.
- Prometheus, node_exporter, and blackbox_exporter were initially absent and later installed.
- The mounted volume was hidden by package-default node_exporter filesystem excludes until the repo's `/etc/default/prometheus-node-exporter` override was applied.

## Known Assumptions and Unknowns

- The dashboards assume Grafana can query a Prometheus-compatible datasource and a Loki datasource.
- The dashboards assume blackbox `instance` labels contain the full probed URL, matching the checked-in Prometheus relabeling.
- Loki JSON parsing panels require request logs to arrive as parseable JSON lines. If Promtail scrapes syslog-prefixed lines, text-only log panels may work while `| json` panels show no data.
- Current Loki labels are server-dependent. Start broad with `loki_job=All`, then narrow the variable or edit selectors after inspecting labels.
- The actual Caddy site file for `main.lekeplassen.regncon.no` is not visible in the repo if it lives under server-local `sites-enabled`.
- Server runtime state, installed packages, open ports, and active services require manual server verification.

## Panel Dependencies

Prometheus panels:

- Public probe status and HTTP/TLS probe panels.
- `caddy.service` and `conorganizer-main.service` systemd status panels.
- Root and mounted-volume filesystem usage, inode usage, and available bytes.
- CPU, memory, load, and network panels.
- Optional SQLite DB/WAL/backup size and precise backup-age panels.

Loki panels:

- Backup success, backup failure, and backup log panels.
- Main service logs, warning/error logs, panic logs, and database symptom logs.
- Log-derived app error signal.
- Text fallback panels for request, Caddy, systemd, Loki, Promtail, Alloy, and critical logs.

node_exporter panels:

- `node_filesystem_*` for `/` and `/mnt/HC_Volume_103911252`.
- `node_filesystem_files*` for inode usage.
- `node_cpu_seconds_total`, `node_memory_*`, `node_load1`, and `node_network_*`.

blackbox_exporter panels:

- `probe_success`.
- `probe_duration_seconds`.
- `probe_http_status_code`.
- `probe_ssl_earliest_cert_expiry`.

systemd metrics panels:

- `node_systemd_unit_state{name="caddy.service", state="active"}`.
- `node_systemd_unit_state{name="conorganizer-main.service", state="active"}`.
- `changes(node_systemd_unit_state{name="conorganizer-main.service", state="active"}[1h])`.

Caddy native metrics:

- No current dashboard panel requires Caddy native Prometheus metrics. Caddy is represented by systemd status, blackbox symptoms, and log searches.

Application metrics:

- No current app Prometheus metrics were found in the repo.
- The debugging dashboard intentionally uses Loki request logs for request rate, status-code mix, route usage, and latency.
- Do not add common HTTP metric names such as `http_requests_total` unless the app actually exposes them later.

Log-derived data:

- Backup freshness in the top row is log-derived from backup success prefixes.
- App request rate can work with text-only request completion logs.
- Status-code, route, and latency panels require parseable JSON request logs with `status_code`, `path`, `method`, and `duration_ms`.

## Representative Queries

Prometheus:

```promql
max(probe_success{instance=~".*${main_host:regex}.*"})
max(node_systemd_unit_state{instance=~"$node_instance",name="$main_service",state="active"})
100 * (1 - (node_filesystem_avail_bytes{instance=~"$node_instance",mountpoint="$volume_mountpoint",fstype!~"tmpfs|overlay|squashfs|ramfs"} / node_filesystem_size_bytes{instance=~"$node_instance",mountpoint="$volume_mountpoint",fstype!~"tmpfs|overlay|squashfs|ramfs"}))
(probe_ssl_earliest_cert_expiry{instance=~".*${main_host:regex}.*"} - time()) / 86400
```

Loki:

```logql
{job=~"$loki_job"} |= "conorganizer-sqlite-backup:"
{job=~"$loki_job"} |= "conorganizer-images-backup:"
sum(count_over_time({job=~"$loki_job"} |= "conorganizer-sqlite-backup: completed successfully" [2h]))
sum(count_over_time({job=~"$loki_job"} |~ "conorganizer-main.service|conorganizer-main|http request completed|level.:.ERROR|database is locked|database is busy|panic" |~ "ERROR|error|failed|panic|status_code.:5[0-9][0-9]|database is locked|database is busy" [5m]))
quantile_over_time(0.95, {job=~"$loki_job"} |= "http request completed" | json | unwrap duration_ms | __error__="" [5m])
```

## Optional Future Metrics

These would make the dashboards more precise if added through node_exporter textfile collector, a custom exporter, or app instrumentation:

- `conorganizer_sqlite_db_size_bytes`
- `conorganizer_sqlite_wal_size_bytes`
- `conorganizer_sqlite_backup_size_bytes`
- `conorganizer_sqlite_backup_last_success_timestamp_seconds`
- `conorganizer_images_backup_last_success_timestamp_seconds`
- `conorganizer_backup_directory_size_bytes`

## Validating Labels and Metrics

In Grafana Explore:

- Loki: start with `{job=~".+"}` and inspect available labels.
- Loki: try `{job=~".+"} |= "http request completed"` for app request logs.
- Loki: try `{job=~".+"} |= "conorganizer-sqlite-backup:"` for backup logs.
- Loki: if labels such as `unit`, `systemd_unit`, or `__journal__systemd_unit` exist, narrow app panels to `conorganizer-main.service`.
- Prometheus: try `up`, `node_uname_info`, `node_filesystem_avail_bytes`, `probe_success`, and `node_systemd_unit_state`.

If JSON parsing panels are empty but text log panels work, inspect whether Loki lines begin with raw JSON or include a syslog prefix before the JSON payload.

## Alert Candidates

These are documented only. No alert rules are provisioned.

- DB backup success missing for more than 2 hours.
- Image backup success missing for more than 26 hours.
- Mounted volume disk usage over 80%.
- Mounted volume disk usage over 90%.
- WAL file larger than expected after future WAL metrics exist.
- `conorganizer-main.service` not active.
- `caddy.service` not active.
- Caddy reload/start failure.
- Public HTTP probe failure for 2-5 minutes.
- TLS certificate expires in less than 14 days.
- Backup failure log seen in the last 24 hours.
- Main service error signal greater than 0 in the last 5 minutes.
- Panic detected.
- 5xx responses detected.
- Database locked/busy errors detected.
- Repeated service restarts or active-state changes.

## Commands for the developer to run on the server

These commands are read-only checks. Do not run them from a workstation against production. Run them manually on the server, and adjust ports if the actual config differs.

```bash
# Service status checks.
systemctl status loki.service --no-pager
systemctl status promtail.service --no-pager
systemctl status alloy.service --no-pager
systemctl status prometheus.service --no-pager
systemctl status prometheus-node-exporter.service --no-pager
systemctl status prometheus-blackbox-exporter.service --no-pager
systemctl status caddy.service --no-pager
systemctl status conorganizer-main.service --no-pager

# Loki readiness. The checked-in Loki config uses port 3500; some installs use 3100.
LOKI_URL=http://127.0.0.1:3500
curl -fsS "$LOKI_URL/ready"

# Promtail and Alloy recent logs.
journalctl -u promtail.service -n 80 --no-pager
journalctl -u alloy.service -n 80 --no-pager
curl -fsS http://127.0.0.1:12345/-/ready
curl -fsS http://127.0.0.1:12345/metrics | head

# Prometheus and exporter readiness.
curl -fsS http://127.0.0.1:9090/-/ready
curl -fsS http://127.0.0.1:9090/api/v1/targets | jq .
curl -fsS http://127.0.0.1:9100/metrics | head
curl -fsS http://127.0.0.1:9115/metrics | head

# Caddy native metrics, only relevant if Caddy metrics are intentionally enabled later.
curl -fsS http://127.0.0.1:2019/metrics | head

# Systemd metrics in Prometheus.
curl -G http://127.0.0.1:9090/api/v1/query \
  --data-urlencode 'query=node_systemd_unit_state{name="conorganizer-main.service",state="active"}' | jq .
curl -G http://127.0.0.1:9090/api/v1/query \
  --data-urlencode 'query=node_systemd_unit_state{name="caddy.service",state="active"}' | jq .

# Node filesystem metrics for root and the mounted volume.
curl -G http://127.0.0.1:9090/api/v1/query \
  --data-urlencode 'query=node_filesystem_avail_bytes{mountpoint="/"}' | jq .
curl -G http://127.0.0.1:9090/api/v1/query \
  --data-urlencode 'query=node_filesystem_avail_bytes{mountpoint="/mnt/HC_Volume_103911252"}' | jq .
curl -G http://127.0.0.1:9090/api/v1/query \
  --data-urlencode 'query=node_filesystem_files_free{mountpoint="/mnt/HC_Volume_103911252"}' | jq .

# Blackbox HTTP and TLS metrics.
curl -G http://127.0.0.1:9090/api/v1/query \
  --data-urlencode 'query=probe_success{instance=~".*main\\.lekeplassen\\.regncon\\.no.*"}' | jq .
curl -G http://127.0.0.1:9090/api/v1/query \
  --data-urlencode 'query=probe_http_status_code{instance=~".*main\\.lekeplassen\\.regncon\\.no.*"}' | jq .
curl -G http://127.0.0.1:9090/api/v1/query \
  --data-urlencode 'query=(probe_ssl_earliest_cert_expiry{instance=~".*main\\.lekeplassen\\.regncon\\.no.*"} - time()) / 86400' | jq .

# Loki labels and label values.
curl -fsS "$LOKI_URL/loki/api/v1/labels" | jq .
curl -fsS "$LOKI_URL/loki/api/v1/label/job/values" | jq .

# Loki backup log checks.
curl -G "$LOKI_URL/loki/api/v1/query_range" \
  --data-urlencode 'query={job=~".+"} |= "conorganizer-sqlite-backup:"' \
  --data-urlencode 'since=24h' \
  --data-urlencode 'limit=20' | jq .
curl -G "$LOKI_URL/loki/api/v1/query_range" \
  --data-urlencode 'query={job=~".+"} |= "conorganizer-images-backup:"' \
  --data-urlencode 'since=48h' \
  --data-urlencode 'limit=20' | jq .

# Backup timers, backup service logs, and backup files.
systemctl list-timers 'conorganizer-*backup*' --no-pager
journalctl -u conorganizer-sqlite-backup.service -n 80 --no-pager
journalctl -u conorganizer-images-backup.service -n 80 --no-pager
find /mnt/HC_Volume_103911252/backups/sqlite -maxdepth 1 -type f -printf '%TY-%Tm-%Td %TH:%TM %s %p\n' | sort | tail
find /mnt/HC_Volume_103911252/backups/images -maxdepth 1 -type f -printf '%TY-%Tm-%Td %TH:%TM %s %p\n' | sort | tail

# Main service logs and request/error samples.
journalctl -u conorganizer-main.service -n 200 --no-pager
journalctl -u conorganizer-main.service --since '1 hour ago' --no-pager | grep 'http request completed'
journalctl -u conorganizer-main.service --since '1 hour ago' --no-pager | grep -Ei 'error|failed|panic|status_code":5[0-9][0-9]'
journalctl -u conorganizer-main.service --since '24 hours ago' --no-pager | grep -Ei 'database is locked|database is busy'
```
