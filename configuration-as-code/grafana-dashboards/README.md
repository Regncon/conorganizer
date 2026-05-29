# Conorganizer Grafana Dashboards

These dashboards are importable Grafana JSON files for manual review/import:

- `conorganizer-production-health.json`
- `conorganizer-main-service-debugging.json`

They are intentionally kept outside `configuration-as-code/stow/` for now. Do not treat them as server-side provisioning until that is explicitly added later.

## Manual Import

1. Open Grafana.
2. Go to **Dashboards** -> **New** -> **Import**.
3. Upload one JSON file from this directory, or paste its contents.
4. When Grafana asks for variables, choose:
   - `DS_PROMETHEUS`: the Prometheus datasource.
   - `DS_LOKI`: the Loki datasource.
   - The JSON defaults these to `Prometheus` and `Loki` for overwrite/import workflows.
5. Import the dashboard.
6. Use the dashboard variables at the top:
   - `node_instance`: the node_exporter instance to inspect.
   - `loki_job`: the Loki job containing service/app logs.
   - `main_host`: defaults to `main.lekeplassen.regncon.no`.
   - `grafana_host`: defaults to `grafana.regncon.no`.

## Expected Datasources

- Prometheus datasource for node_exporter, blackbox_exporter, optional systemd metrics, and optional future app/custom metrics.
- Loki datasource for app logs, backup logs, Caddy/systemd error logs, and log-derived HTTP behavior.

The checked-in repo contains Grafana, Loki, Promtail, Caddy, backup scripts, and systemd unit files. It does not contain Prometheus configuration, blackbox_exporter configuration, node_exporter configuration, app Prometheus instrumentation, `/metrics`, `/healthz`, or `/readyz`. As of the manual server check on 2026-05-28, Prometheus, node_exporter, and blackbox_exporter were not running on the server.

## Repo Observations

- `conorganizer-main.service` runs `/opt/conorganizer/main/conorganizer-main` with the database at `/mnt/HC_Volume_103911252/environments/main/database/events.db` and images at `/mnt/HC_Volume_103911252/environments/main/event-images`.
- The app emits JSON `slog` logs to stdout.
- Request logs use `msg="http request completed"` with `component="http"`, `method`, `path`, `status_code`, `duration_ms`, and optional `request_id`.
- Backup scripts log with stable prefixes:
  - `conorganizer-sqlite-backup:`
  - `conorganizer-images-backup:`
- The SQLite backup timer runs every 15 minutes.
- The image backup timer runs daily at 03:30.
- The checked-in Promtail config labels logs as `job=varlogs` and scrapes `/var/log/messages`.
- The checked-in Loki config listens on HTTP port `3500`, while the checked-in Promtail client points to `http://localhost:3100/loki/api/v1/push`; verify the actual server config before relying on Loki panels.

## Manual Server Check Results

From the developer-run checks on 2026-05-28:

- Loki is ready on `http://127.0.0.1:3500`.
- Promtail is active and currently tails `/var/log/messages` with `job=varlogs`.
- Promtail has old `no space left on device` position-file errors from 2026-05-16, but the latest 2026-05-28 startup logs are clean.
- `prometheus.service`, `node_exporter.service`, and `blackbox_exporter.service` were not found.
- Ports `9090`, `9100`, and `9115` were not listening.
- `jq` was installed after the first check, so command examples use `jq`.
- SQLite backup is running every 15 minutes and recent runs completed successfully.
- Image backup has run successfully and the next timer is scheduled for the following 03:30 UTC window.
- Alloy v1.16.1 is active and exposes its own readiness and self-metrics on `127.0.0.1:12345`.
- The current Alloy config includes `loki.write` and `loki.source.journal`.
- The current Alloy config does not include `prometheus.scrape`, `prometheus.remote_write`, `prometheus.exporter.unix`, or `prometheus.exporter.blackbox`, so it currently helps with logs, not the Prometheus panels in these dashboards.
- Prometheus, node_exporter, and blackbox_exporter were later installed and are scrapeable.
- The package default node_exporter filesystem settings did not expose `/mnt/HC_Volume_103911252`, so the stow package includes `/etc/default/prometheus-node-exporter` to keep `/mnt` visible and enable systemd metrics.
- After stowing `/etc/default/prometheus-node-exporter`, Prometheus sees `/mnt/HC_Volume_103911252` and returns `node_filesystem_avail_bytes` / `node_filesystem_size_bytes` for `/dev/sdb`.

## Alloy Note

Grafana Alloy can replace Prometheus for scrape configuration and collection pipelines. Grafana's migration guide converts Prometheus `scrape_configs` into `prometheus.scrape` components and converts `remote_write` into `prometheus.remote_write` components. It does not make Alloy a drop-in replacement for Prometheus storage, `remote_read`, rules, alerting, or the Prometheus query API that Grafana uses for historical PromQL dashboard panels.

For these dashboards, Alloy only compensates for the missing Prometheus service if it forwards metrics to a Prometheus-compatible datasource that Grafana can query, such as Grafana Cloud Metrics, Mimir, or a Prometheus instance.

Relevant Alloy roles:

- It can replace standalone Promtail for logs if configured with `loki.source.journal` or `loki.source.file` and a `loki.write` destination.
- It can replace standalone node_exporter collection if configured with `prometheus.exporter.unix`, `prometheus.scrape`, and a metrics write destination.
- It can replace standalone blackbox_exporter collection if configured with `prometheus.exporter.blackbox`, `prometheus.scrape`, and a metrics write destination.
- It still needs a metrics backend such as Prometheus, Grafana Mimir, Grafana Cloud Metrics, or another Prometheus-compatible remote-write target for Grafana dashboards to query historical metrics with PromQL.

Useful references:

- https://grafana.com/docs/alloy/latest/reference/components/prometheus/prometheus.scrape/
- https://grafana.com/docs/alloy/latest/reference/components/prometheus/prometheus.remote_write/
- https://grafana.com/docs/alloy/latest/reference/components/prometheus/prometheus.exporter.unix/
- https://grafana.com/docs/alloy/latest/reference/components/prometheus/prometheus.exporter.blackbox/
- https://grafana.com/docs/alloy/latest/reference/components/loki/loki.source.journal/

## Getting Prometheus Panels Working

The dashboards need a Grafana datasource that can answer PromQL queries. There are two practical paths:

1. Local Prometheus on the VPS.
   - Install and run Prometheus.
   - Use the stow-owned Prometheus config files in `configuration-as-code/stow/prometheus/etc/prometheus/`.
   - Configure Prometheus as a Grafana datasource.
   - Scrape node/system metrics and blackbox probes either from standalone exporters or from Alloy-generated exporter targets.
   - Use the stow-owned `/etc/default/prometheus-node-exporter` override because the package default filesystem exclude hid the mounted Hetzner volume.

2. Alloy as collector plus a Prometheus-compatible backend.
   - Configure Alloy `prometheus.exporter.unix` for node filesystem, CPU, memory, load, network, and systemd metrics.
   - Enable the systemd collector in the Unix exporter so `node_systemd_unit_state` exists.
   - Configure Alloy `prometheus.exporter.blackbox` for `main.lekeplassen.regncon.no` and `grafana.regncon.no`.
   - Configure Alloy `prometheus.scrape` to scrape those exporter targets.
   - Configure Alloy `prometheus.remote_write` to send metrics to Grafana Cloud Metrics, Mimir, or another Prometheus-compatible remote-write backend.
   - Configure that backend as the Grafana `DS_PROMETHEUS` datasource.

Minimum metrics needed by the current dashboards:

- `node_filesystem_*`
- `node_cpu_seconds_total`
- `node_memory_*`
- `node_load1`
- `node_network_*`
- `node_systemd_unit_state`
- `probe_success`
- `probe_duration_seconds`
- `probe_http_status_code`
- `probe_ssl_earliest_cert_expiry`

Optional future custom metrics:

- `conorganizer_sqlite_db_size_bytes`
- `conorganizer_sqlite_wal_size_bytes`
- `conorganizer_sqlite_backup_size_bytes`
- `conorganizer_sqlite_backup_last_success_timestamp_seconds`
- `conorganizer_images_backup_last_success_timestamp_seconds`
- `conorganizer_backup_directory_size_bytes`

## Panel Dependencies

Prometheus panels:

- Website/Grafana probe status.
- Caddy and `conorganizer-main.service` status.
- Disk, inode, CPU, memory, load, and network panels.
- TLS expiry and HTTP probe panels.
- Optional SQLite DB/WAL/backup size and backup age panels.

Loki panels:

- Backup logs and backup failure logs.
- SQLite/database symptom logs.
- Caddy/systemd/Loki/Promtail error logs.
- Main service logs, warning/error logs, panic logs, failed request logs.
- Log-derived request rate, status-code rate, latency, and route usage.

node_exporter panels:

- Root filesystem and `/mnt/HC_Volume_103911252` usage.
- Inode usage.
- Available volume bytes.
- CPU, memory, load, and network throughput.

blackbox_exporter panels:

- `probe_success` for `main.lekeplassen.regncon.no` and `grafana.regncon.no`.
- `probe_duration_seconds`.
- `probe_http_status_code`.
- `probe_ssl_earliest_cert_expiry`.

systemd metrics panels:

- `node_systemd_unit_state{name="caddy.service", state="active"}`.
- `node_systemd_unit_state{name="conorganizer-main.service", state="active"}`.
- Active-state changes for `conorganizer-main.service`.

Application metrics panels:

- No current app Prometheus metrics were found in the repo.
- The debugging dashboard uses Loki request logs for HTTP behavior and latency.
- Optional future metrics are clearly labeled, including `conorganizer_sqlite_db_size_bytes`, `conorganizer_sqlite_wal_size_bytes`, backup size, and backup last-success timestamp metrics.

## Alert Candidates

These are documented in panel descriptions only. No alert rules are provisioned.

- DB backup older than 2 hours.
- Image backup older than 26 hours.
- Mounted volume disk usage over 80%.
- WAL file larger than expected.
- `conorganizer-main.service` not active.
- Caddy reload/start failure.
- Public HTTP probe failure for 2-5 minutes.
- TLS certificate expires in less than 14 days.
- Main service error signal greater than 0 in the last 5 minutes.
- Panic detected.
- 5xx responses detected.
- Database locked errors detected.

## Validating Labels and Metrics

In Grafana Explore:

- Loki: start with `{job=~".+"}` and inspect available labels.
- Loki: try `{job=~".+"} |= "http request completed"` for app request logs.
- Loki: try `{job=~".+"} |= "conorganizer-sqlite-backup:"` for backup logs.
- Prometheus: try `up`, `node_uname_info`, `node_filesystem_avail_bytes`, `probe_success`, and `node_systemd_unit_state`.

If Loki has systemd labels such as `unit`, `systemd_unit`, or `__journal__systemd_unit`, narrow the dashboard selectors after import for cleaner app and service panels.

## Commands for the Developer to Run on the Server

Do not run all of these blindly; use the relevant checks for the exporter or panel you are validating.

```bash
# Install local Prometheus and standalone exporters.
sudo apt-get update
sudo apt-get install -y stow prometheus prometheus-node-exporter prometheus-blackbox-exporter

# Back up package-created config before replacing it with stow symlinks.
sudo cp -a /etc/prometheus/prometheus.yml "/etc/prometheus/prometheus.yml.$(date -u +%Y%m%dT%H%M%SZ).bak" 2>/dev/null || true
sudo cp -a /etc/prometheus/blackbox.yml "/etc/prometheus/blackbox.yml.$(date -u +%Y%m%dT%H%M%SZ).bak" 2>/dev/null || true
sudo cp -a /etc/default/prometheus-node-exporter "/etc/default/prometheus-node-exporter.$(date -u +%Y%m%dT%H%M%SZ).bak" 2>/dev/null || true
sudo rm -f /etc/prometheus/prometheus.yml /etc/prometheus/blackbox.yml
sudo rm -f /etc/default/prometheus-node-exporter

# Stow only the Prometheus config files from this repo.
cd /home/cinmay/Documents/conorganizer/configuration-as-code/stow
sudo stow --target=/ prometheus

# Start/restart services with the stowed config.
sudo systemctl daemon-reload
sudo systemctl enable --now prometheus prometheus-node-exporter prometheus-blackbox-exporter
sudo systemctl restart prometheus prometheus-node-exporter prometheus-blackbox-exporter

# Loki readiness. Use 3100 instead if that is the actual configured HTTP port.
LOKI_URL=http://127.0.0.1:3500
curl -fsS "$LOKI_URL/ready"

# Promtail service and recent errors.
systemctl status promtail.service --no-pager
journalctl -u promtail.service -n 80 --no-pager

# Prometheus readiness.
curl -fsS http://127.0.0.1:9090/-/ready
systemctl status prometheus.service --no-pager

# node_exporter.
curl -fsS http://127.0.0.1:9100/metrics | head
systemctl status prometheus-node-exporter.service --no-pager

# blackbox_exporter, if installed.
curl -fsS http://127.0.0.1:9115/metrics | head
systemctl status prometheus-blackbox-exporter.service --no-pager

# Systemd metrics in Prometheus.
curl -G http://127.0.0.1:9090/api/v1/query \
  --data-urlencode 'query=node_systemd_unit_state{name="conorganizer-main.service",state="active"}' | jq .

# Node filesystem metrics for the mounted volume.
curl -G http://127.0.0.1:9090/api/v1/query \
  --data-urlencode 'query=node_filesystem_avail_bytes{mountpoint="/mnt/HC_Volume_103911252"}' | jq .

# Blackbox TLS metrics.
curl -G http://127.0.0.1:9090/api/v1/query \
  --data-urlencode 'query=(probe_ssl_earliest_cert_expiry{instance=~".*main\\.lekeplassen\\.regncon\\.no.*"} - time()) / 86400' | jq .

# Loki labels and job values.
curl -fsS "$LOKI_URL/loki/api/v1/labels" | jq .
curl -fsS "$LOKI_URL/loki/api/v1/label/job/values" | jq .

# Loki backup log checks.
curl -G "$LOKI_URL/loki/api/v1/query_range" \
  --data-urlencode 'query={job=~".+"} |= "conorganizer-sqlite-backup:"' \
  --data-urlencode 'limit=20' | jq .

curl -G "$LOKI_URL/loki/api/v1/query_range" \
  --data-urlencode 'query={job=~".+"} |= "conorganizer-images-backup:"' \
  --data-urlencode 'limit=20' | jq .

# Backup timers, service logs, and files.
systemctl list-timers 'conorganizer-*backup*' --no-pager
journalctl -u conorganizer-sqlite-backup.service -n 80 --no-pager
journalctl -u conorganizer-images-backup.service -n 80 --no-pager
ls -lh /mnt/HC_Volume_103911252/backups/sqlite | tail
ls -lh /mnt/HC_Volume_103911252/backups/images | tail
```
