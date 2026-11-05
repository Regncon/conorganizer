## Summary

Loki lost all its logs on every reboot. Its config was the example file the Loki package installs (`cmd/loki/loki-local-config.yaml`), which keeps everything in `/tmp/loki`. On the server, `/tmp` is a 1 GB tmpfs in RAM. That's why Loki had nothing from before 19 October. The originals in the systemd journal and `/var/log/syslog*` were never affected.

There was a second fault: `retention_period: 120d` never deleted anything, because retention was off. On a real disk, Loki would have grown without limit.

### Changes

- **`config.yml`:**
  - `path_prefix`, `chunks_directory` and `rules_directory` move to `/var/lib/loki` on the root disk. WAL, index and compactor follow `path_prefix`.
  - `compactor.retention_enabled: true` with `delete_request_store: filesystem`, so the 120-day retention runs.
  - `limits_config.deletion_mode: disabled`. Turning on retention also turns on Loki's delete API, which has no authentication here (`auth_enabled: false`). Retention doesn't need it.
- **New drop-in `loki.service.d/state-directory.conf`:** `StateDirectory=loki` with mode 0700, so systemd creates the folder for the `loki` user.
  - It's a drop-in because the Loki package ships `/etc/systemd/system/loki.service` as an ordinary file and overwrites it on upgrade.
- **`configuration-as-code/README.md`:** a "Loki storage" section explaining the storage and the size estimate, with the one-off server steps. The dashboards README and the "Root filesystem usage" panel description are updated too.

Storing Loki's data on the local disk is the setup Loki's docs describe as common for a single binary on one server. Object storage, such as Hetzner Object Storage, would be the next step if logs have to survive losing the server.

## Server steps (after merging)

Follow "Moving from /tmp/loki" in `configuration-as-code/README.md`, in one go:
1. Check what the update brings, and run `loki -verify-config` on the new config.
2. Stop Loki.
3. Update the checkout with `git checkout --detach origin/main`.
4. Copy `/tmp/loki` with `cp -a`.
5. Restow `loki` and `systemd`, run `daemon-reload` and start Loki.
6. Check `/ready`, `/config`, the size, and that Alloy dropped nothing.

Then re-import Production Health into the "Configuration-as-code" folder.

**Order matters.** `/etc/loki/config.yml` is a stow symlink into the checkout, so updating the checkout switches Loki to `/var/lib/loki` straight away. If Loki restarts before `daemon-reload` and the drop-in are in place, it can't create its folder and stays down. Stop Loki first, and don't update the server's checkout past this commit without running the steps.

## Size

On `/tmp`, Loki grew about 15 MB a day. About half of that is its own query logs, so the rate depends on how much people query. 120 days comes to roughly 2 GB, taking `/` from about 84% to about 89%, just under the 90% red line. Check `du -sh /var/lib/loki` after a week, and lower `retention_period` if Loki grows faster.

## Testing

- Every config key was checked against the Loki v3.7.8 source, the version on the server (`pkg/compactor/config.go`, `pkg/loki/config_wrapper.go`, `docs/sources/operations/storage/retention.md`). Nothing still points at `/tmp`.
- A review workflow covered three angles: config, migration safety and documentation. A fourth agent then tried to refute each finding. It confirmed 10 problems and fixed them all, including:
  - The step order: stop Loki before updating the checkout.
  - The missing `loki` group, which ruled out `install -g loki`.
  - The delete API.
  - The drop-in instead of an edit to the unit itself.
  - Checks that confirm the switch actually worked.
- **Not tested here:** `loki -verify-config` runs on the server as the first step, because Docker wasn't available locally.

## Follow-up

The "Check in the Alloy config and cap Loki memory" task covers Loki memory limits. On 2026-10-20, about 50 heavy 7-day queries briefly ran the 2 GB server out of memory.

The merge commit from main also changes the server steps to the admin-owned, detached checkout (`git …`).


<!-- preview-deployment-url:start -->
## Preview deployment

_Added automatically by CI_

🔗 https://673-merge.lekeplassen.regncon.no
<!-- preview-deployment-url:end -->

