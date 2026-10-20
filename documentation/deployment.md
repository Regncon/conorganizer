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

- **`build`** runs on every push to `main` and on every PR `opened`, `reopened`, `synchronize` and `ready_for_review` event. It runs `go tool templ generate`, then the tests with the behavior report (`go run ./cmd/testreport`), then builds the `conorganizer` binary and a Goose v3.28.0 binary (`go install github.com/pressly/goose/v3/cmd/goose@v3.28.0`) and uploads both as an artifact.
- **`deploy-fixed`** runs on pushes to `main`. It is a matrix over `[main, demo, restored]` with `fail-fast: false`, so the three environments deploy in parallel and a failure in one does not cancel the others. This is safe because each environment has its own app dir, binary, service and data dir, and fixed deploys install no generated Caddy or systemd files, so they share no files. Each matrix entry uploads the binary, `goose`, `migrations/*.sql` and `deploy.sh` to `/opt/conorganizer/<name>/` and runs `sudo deploy.sh <name>`. CI never runs Goose against a fixed environment; the deployed files are used by the manual `conorganizer-sqlite-migrate` (see [migrations.md](migrations.md)).
- **`deploy`** builds a PR preview. It only runs for PRs that are open, not drafts, and come from this repository (fork PRs have no secrets). A draft gets a preview when it is marked ready for review.
- **`cleanup-preview`** runs when a same-repo PR is closed, merged or not. It is the only job that runs on the `pull_request_target` event. See [Preview cleanup](#preview-cleanup).

Runs are serialized per PR, and pushes per ref, with a `concurrency` group on `${{ github.workflow }}-${{ github.event.pull_request.number || github.ref }}` and `cancel-in-progress: false`. A deploy that has started is never interrupted halfway through moving the binary or restarting the service, and a PR's cleanup waits for that PR's running deploy. GitHub keeps only the newest pending run, so the latest commit is always deployed last. The group uses the PR number, not `github.ref`, because a merged PR's close event has `github.ref` `refs/heads/main`. In main's group, the next merge cancelled the pending cleanup.

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

After each successful run the scripts write `conorganizer_*` metrics for the Grafana backup panels to node_exporter's textfile directory. See [Backup Metrics](../configuration-as-code/grafana-dashboards/README.md#backup-metrics).

## Preview cleanup

`cleanup-preview` uploads `deploy/cleanup.sh` to the preview's app dir and runs `sudo cleanup.sh <PR_NUMBER>-merge`. The script removes everything a preview deploy creates, in this order:

1. Stops and disables the service. If the service is still active after that, the script exits without removing anything, so a running preview never loses its files.
2. Removes the unit and its `multi-user.target.wants` link. `systemctl disable` only removes the link while the unit file exists, so the script removes both itself.
3. Removes the Caddy site and the data dir (`database/`, `event-images/`, everything the app wrote there).
4. Reloads systemd, resets the unit's failed state so it stops showing as `not-found`, and reloads Caddy. A failed Caddy reload fails the job, because Caddy then still serves the old route.
5. Removes the app dir last, because it holds the script itself.

`/run/conorganizer-<name>/` is the unit's `RuntimeDirectory` and holds the NATS store; systemd removes it when the service stops. Two things stay on purpose: Caddy's certificate for the host, which Caddy deletes from its storage after it expires, and the journal entries of the unit, which follow the journal's own retention.

Every removal tolerates a target that is already gone, so running it twice, after a failed deploy, or for a draft that never got a preview, is safe.

The job runs on `pull_request_target`, not `pull_request`, because GitHub runs no `pull_request` workflow for a PR with merge conflicts, not even when it is closed. `pull_request_target` runs main's workflow with the secrets, whatever the PR's base branch, so the job only checks out main and never runs PR code. A change to `cleanup-preview` or `cleanup.sh` therefore only takes effect once it is merged to main.

The job needs the `deploy` user's sudoers entry for `/opt/conorganizer/*/cleanup.sh`, which `hetzner/scripts/setup_deploy_user.sh` writes. Check it on the server with `sudo -l -U deploy`. Without it the job fails at the `sudo` line and the preview stays.

To see what the script would remove without changing anything, run a copy of `deploy/cleanup.sh` as a dry run. A live preview's app dir has no copy, because CI only uploads the script when the PR closes. A dry run needs no root:

```bash
DRY_RUN=true bash cleanup.sh <PR_NUMBER>-merge
```

Never write `DRY_RUN=true sudo ...`: sudo drops `DRY_RUN`, so the script removes the preview.

Safety notes:

- The job takes the target name from the PR number (`<PR_NUMBER>-merge`), not from `github.ref_name`. On the close event of a merged PR, `github.ref_name` can resolve to `main`, so never use it to pick the cleanup target.
- `cleanup.sh` refuses every name that is not `<digits>-merge`, so it cannot remove `main`, `demo`, `restored` or anything outside the preview paths.

### Remove stale previews by hand

A preview whose cleanup job failed, was cancelled or never ran stays on the server: Caddy keeps its route and certificate and logs a `502` with `connection refused` for every request once its service is gone. To find and remove such previews, first list the open PRs on your own machine:

```bash
gh pr list --repo Regncon/conorganizer --state open --limit 200 --json number --jq '[.[].number] | join(" ")'
```

Then, on the server, collect every preview that still has a Caddy site, a unit, an app dir or a data dir, and leave out the open PRs. Paste the numbers into `OPEN_PRS` (any separator works). Leave it empty only when no PR is open, because every preview left out of it is removed:

```bash
OPEN_PRS=""
open=" $(echo "$OPEN_PRS" | grep -oE '[0-9]+' | tr '\n' ' ')"
stale=$(
  {
    sudo ls /etc/caddy/sites-enabled | sed -nE 's/^conorganizer-([0-9]+-merge)\.caddy$/\1/p'
    sudo ls /etc/systemd/system | sed -nE 's/^conorganizer-([0-9]+-merge)\.service$/\1/p'
    sudo ls /opt/conorganizer /mnt/HC_Volume_103911252/environments | grep -E '^[0-9]+-merge$'
  } | sort -un | while read -r name; do
    case "$open" in *" ${name%-merge} "*) ;; *) echo "$name" ;; esac
  done
)
echo $stale
```

Copy `deploy/cleanup.sh` from an up-to-date `main` to the server, for example with `scp deploy/cleanup.sh <user>@<server>:/tmp/cleanup.sh`. Don't `git pull` the configuration-as-code checkout for this: the script reloads Caddy, which would also put any pulled Caddyfile change live. Then run it for each stale preview, first as a dry run, then for real:

```bash
for name in $stale; do DRY_RUN=true bash /tmp/cleanup.sh "$name"; done
for name in $stale; do sudo DRY_RUN=false bash /tmp/cleanup.sh "$name"; done
```

Put `DRY_RUN=false` after `sudo`, because sudo resets the environment. When the last site file is gone, Caddy's `import /etc/caddy/sites-enabled/*.caddy` only logs a warning, so the reload still succeeds.
