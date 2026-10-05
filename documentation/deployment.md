# Deployment

Conorganizer runs on one Hetzner server. CI builds one binary and deploys it to the fixed environments on pushes to `main`, and to a short-lived preview environment for each open pull request.

The pipeline is `.github/workflows/buildAndTest.yml`. The server-side scripts are in `deploy/`. The fixed environments' systemd units, Caddy routes and backup scripts are in `configuration-as-code/` (see [configuration-as-code/README.md](../configuration-as-code/README.md)).

## Environments

There are two kinds of environments:

- **Fixed environments**: `main`, `demo` and `restored`. Every push to `main` deploys the same binary to all three. Their systemd units and Caddy routes are checked in under `configuration-as-code/stow/` and CI never generates them. `deploy.sh` does not install config for them and does not clone data into them.
- **PR previews**: short-lived, one per pull request. CI renders a systemd unit and a Caddy site from `deploy/conorganizer.service.tmpl` and `deploy/caddy-site.tmpl` on every push to the PR, and `deploy.sh` installs them.

| Name | Host | Port | systemd unit | App dir | Data dir |
| --- | --- | --- | --- | --- | --- |
| main | `main.lekeplassen.regncon.no` (`program.regncon.no` redirects to `https://regncon.no/gamle-regncon-program/` while the con is over) | 19080 | `conorganizer-main.service` | `/opt/conorganizer/main/` | `/mnt/HC_Volume_103911252/environments/main/` |
| demo | `demo.lekeplassen.regncon.no` | 19081 | `conorganizer-demo.service` | `/opt/conorganizer/demo/` | `/mnt/HC_Volume_103911252/environments/demo/` |
| restored | `restored.lekeplassen.regncon.no` | 19082 | `conorganizer-restored.service` | `/opt/conorganizer/restored/` | `/mnt/HC_Volume_103911252/environments/restored/` |
| PR preview | `<PR_NUMBER>-merge.lekeplassen.regncon.no` | 20000 + PR_NUMBER | `conorganizer-<PR_NUMBER>-merge.service` | `/opt/conorganizer/<PR_NUMBER>-merge/` | `/mnt/HC_Volume_103911252/environments/<PR_NUMBER>-merge/` |

The fixed ports are set in the checked-in units with `Environment=PORT=...`. The app listens on `0.0.0.0:$PORT` and falls back to 8080 when `PORT` is unset. The fixed units read secrets from `EnvironmentFile=/etc/conorganizer/env`, except demo, which uses `/etc/conorganizer/demo-env`. Previews use `/etc/conorganizer/env`.

Each environment has the same layout:

- The data dir holds `database/events.db` and `event-images/`.
- The app dir holds the running binary `conorganizer-<name>`, plus the `goose` binary and `migrations/*.sql` that CI deploys with it. On deploy the previous binary is kept as `conorganizer-<name>.old`.
- The unit starts the binary with `-dbp <data dir>/database/events.db`, `-image-path <data dir>/event-images` and `-nats-store-dir /run/conorganizer-<name>/nats`. `/run/conorganizer-<name>` is the unit's `RuntimeDirectory`, so the embedded NATS data is thrown away when the service stops.

`restored` is special: `deploy.sh` installs its binary but does not start the service until a backup has been restored into it with `conorganizer-sqlite-restore` (see [Run a Restored Database Backup](../README.md#run-a-restored-database-backup)).

## CI/CD jobs

- **`build`** runs on every push to `main` and on every PR event except `closed`. It runs `go tool templ generate`, then the tests with the behavior report (`go run ./cmd/testreport`), then builds the `conorganizer` binary and a Goose v3.28.0 binary (`go install github.com/pressly/goose/v3/cmd/goose@v3.28.0`) and uploads both as an artifact.
- **`deploy-fixed`** runs on pushes to `main`. It is a matrix over `[main, demo, restored]` with `fail-fast: false`, so the three environments deploy in parallel and a failure in one does not cancel the others. This is safe because each environment has its own app dir, binary, service and data dir, and fixed deploys install no generated Caddy or systemd files, so they share no files. Each matrix entry uploads the binary, `goose`, `migrations/*.sql` and `deploy.sh` to `/opt/conorganizer/<name>/` and runs `sudo deploy.sh <name>`. CI never runs Goose against a fixed environment; the deployed files are used by the manual `conorganizer-sqlite-migrate` (see [migrations.md](migrations.md)).
- **`deploy`** builds a PR preview. It only runs for PRs that are open, not drafts, and come from this repository (fork PRs have no secrets). A draft gets a preview when it is marked ready for review.
- **`cleanup-preview`** runs when a same-repo PR is closed, merged or not. See [Preview cleanup](#preview-cleanup).

Runs are serialized per ref with a `concurrency` group on `${{ github.workflow }}-${{ github.ref }}` and `cancel-in-progress: false`. A deploy that has started is never interrupted halfway through moving the binary or restarting the service. GitHub keeps only the newest pending run, so the latest commit is always deployed last.

## Production environment and SSH secrets

`HETZNER_SSH_KEY`, `HETZNER_HOST` and `HETZNER_USER` are scoped to the `Production` GitHub environment. Every job that uses SSH (`deploy-fixed`, `deploy` and `cleanup-preview`) must declare:

```yaml
environment:
    name: Production
```

Without it the secrets are empty. Each of these jobs has a `Check that HETZNER_SSH_KEY is set` step that fails early with `HETZNER_SSH_KEY is empty or not available in this context`. The `build` job does not use SSH and does not declare the environment.

## PR previews

CI computes the preview's name, port and host from `github.event.pull_request.number`:

- name: `<PR_NUMBER>-merge`, for example `482-merge`
- port: `20000 + PR_NUMBER`, for example `20482`
- host: `https://<PR_NUMBER>-merge.lekeplassen.regncon.no`

On each push to the PR, CI renders the unit and Caddy site into a bundle together with the binary, `goose`, `migrations/*.sql` and `deploy.sh`. `deploy.sh` then:

1. On the first deploy only, makes a SQLite `.backup` snapshot of main's database, checks it with `PRAGMA quick_check`, and copies main's `event-images/`. Later deploys keep the preview's existing data.
2. Runs `goose up -allow-missing` with the deployed migrations against the preview database, on every deploy, so later pushes to the PR get their new migrations.
3. Changes ownership of the preview's data dir to `deploy:www-data`.
4. Installs the unit into `/etc/systemd/system/` and the Caddy site into `/etc/caddy/sites-enabled/conorganizer-<name>.caddy`, then reloads Caddy.
5. Promotes the new binary, restarts the service and fails the job (printing the last 50 journal lines) if the service does not become active.

After a successful deploy, CI adds a `Preview deployment` block with the URL to the PR description. It only does this once: the block starts with the HTML comment `<!-- preview-deployment-url:start -->`, and CI skips the step when the description already contains it.

## Caddy routing

The root Caddyfile is `configuration-as-code/stow/caddy/etc/caddy/Caddyfile`. Routing for `program.regncon.no`, `main.lekeplassen.regncon.no`, `demo.lekeplassen.regncon.no`, `restored.lekeplassen.regncon.no` and `grafana.regncon.no` is always defined there. `demo` and `restored` proxy straight to `127.0.0.1:19081` and `127.0.0.1:19082`.

Only PR preview sites live in `/etc/caddy/sites-enabled/conorganizer-<name>.caddy`. `deploy.sh` installs them and `cleanup.sh` removes them. The root Caddyfile loads them with:

```caddyfile
import /etc/caddy/sites-enabled/*.caddy
```

### Maintenance page

The Caddyfile defines two snippets: `conorganizer-main` (a reverse proxy to `127.0.0.1:19080`) and `conorganizer-maintenance` (a bilingual English/Norwegian maintenance page served with status 503 and `Cache-Control: no-store`). While the con is over, `program.regncon.no` only redirects (302) to `https://regncon.no/gamle-regncon-program/`, so the maintenance flag has no effect until the previous block is restored from git history. That block serves the maintenance snippet while `/var/lib/conorganizer/maintenance.on` exists and `conorganizer-main` otherwise. Switch with `sudo conorganizer-maintenance-mode on|off`; Caddy checks the flag file on every request, so no reload is needed. `main.lekeplassen.regncon.no` always imports `conorganizer-main`, so main can still be reached there while `program.regncon.no` shows the maintenance page. The migration procedure in [migrations.md](migrations.md) uses this toggle. See [configuration-as-code/README.md](../configuration-as-code/README.md#maintenance-mode) for details.

## Service user and ownership

All Conorganizer services run as `User=deploy` and `Group=www-data`. This covers the main, demo and restored units and the preview template. `deploy.sh` changes ownership of the app dir, the promoted binary and new preview data dirs to `deploy:www-data`. Database files you put in place by hand, for example during a restore or migration, must keep that ownership, or the service cannot open them.

## Backups

Scheduled backups of main are written to `/mnt/HC_Volume_103911252/backups`:

- `sqlite/` holds compressed database backups (`events-<timestamp>.db.zst`) plus an uncompressed `events-latest.db`, written by `conorganizer-sqlite-backup` every 15 minutes.
- `images/` holds image backups (`event-images-<timestamp>.tar.zst`), written by `conorganizer-images-backup` daily at 03:30.

Both keep 14 days of backups and are run by systemd timers from `configuration-as-code`. See [configuration-as-code/README.md](../configuration-as-code/README.md#backups) for the timers, and [sqlite-and-health.md](sqlite-and-health.md) for why backups use SQLite `.backup`.

## Preview cleanup

`cleanup-preview` uploads `deploy/cleanup.sh` to the preview's app dir and runs it with the preview name. That script stops and disables the service, removes the unit, the Caddy site, the data dir and the app dir, and reloads systemd and Caddy. Every step tolerates a target that is already gone, so running it twice, or after a failed deploy, is safe.

**Cleanup is currently a dry run.** `cleanup.sh` defaults to `DRY_RUN=true`, and the workflow also passes `DRY_RUN=true` explicitly. Closing a PR therefore only logs what would be removed. Preview services keep running, and their units, Caddy sites and data dirs pile up on the server until someone removes them by hand:

```bash
sudo DRY_RUN=false bash /opt/conorganizer/<PR_NUMBER>-merge/cleanup.sh <PR_NUMBER>-merge
```

Put `DRY_RUN=false` after `sudo`. `DRY_RUN=false sudo ...` normally still dry-runs, because sudo resets the environment. Real teardown needs root. To turn it on in CI, first add a sudoers entry for `cleanup.sh` on the server, then change the default in the script and remove `DRY_RUN=true` from the workflow.

Safety notes:

- The job takes the target name from the PR number (`<PR_NUMBER>-merge`), not from `github.ref_name`. On the close event of a merged PR, `github.ref_name` can resolve to `main`, so never use it to pick the cleanup target.
- `cleanup.sh` refuses to run for `main` and `demo`. It does **not** refuse `restored`. `restored` is only safe because the workflow always passes `<PR_NUMBER>-merge`. Never run the script by hand with a fixed environment's name.
