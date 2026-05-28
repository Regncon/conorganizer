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
5. Import the dashboard.
6. Use the dashboard variables at the top:
   - `node_instance`: the node_exporter instance to inspect.
   - `loki_job`: the Loki job containing service/app logs.
   - `main_host`: defaults to `main.lekeplassen.regncon.no`.
   - `grafana_host`: defaults to `grafana.regncon.no`.

## Expected Datasources

- Prometheus datasource for node_exporter, blackbox_exporter, optional systemd metrics, and optional future app/custom metrics.
- Loki datasource for app logs, backup logs, Caddy/systemd error logs, and log-derived HTTP behavior.

The checked-in repo contains Grafana, Loki, Promtail, Caddy, backup scripts, and systemd unit files. It does not contain Prometheus configuration, blackbox_exporter configuration, node_exporter configuration, app Prometheus instrumentation, `/metrics`, `/healthz`, or `/readyz`.

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
systemctl status node_exporter.service --no-pager

# blackbox_exporter, if installed.
curl -fsS http://127.0.0.1:9115/metrics | head
systemctl status blackbox_exporter.service --no-pager

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
