# Deployment

Conorganizer runs on one Hetzner server. CI builds one binary and deploys it to the fixed environments on pushes to `main`, and to a short-lived preview environment for each open pull request.

The pipeline is `.github/workflows/buildAndTest.yml`. CI uploads only build artefacts to the server and runs two server-side scripts through sudo: `conorganizer-deploy` and `conorganizer-cleanup`. They, the preview templates they render, and the fixed environments' systemd units, Caddy routes and backup scripts are in `configuration-as-code/` (see [configuration-as-code/README.md](../configuration-as-code/README.md#deploy-scripts)).

## Environments

There are two kinds of environments:

- **Fixed environments**: `main`, `demo` and `restored`. Every push to `main` deploys the same binary to all three. Their systemd units and Caddy routes are checked in under `configuration-as-code/stow/` and CI never generates them. `conorganizer-deploy` installs no config for them and clones no data into them.
- **PR previews**: short-lived, one per pull request. On every push to the PR, `conorganizer-deploy` renders the preview's systemd unit and Caddy site on the server from `configuration-as-code/stow/scripts/usr/local/share/conorganizer/preview.service.tmpl` and `preview.caddy.tmpl`, with the port and host derived from the preview's name, and installs them. CI never renders them.

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

`restored` is special: `conorganizer-deploy` installs its binary but does not start the service until a backup has been restored into it with `conorganizer-sqlite-restore` (see [Run a Restored Database Backup](../README.md#run-a-restored-database-backup)).

## CI/CD jobs

- **`build`** runs on every push to `main` and on every PR `opened`, `reopened`, `synchronize` and `ready_for_review` event. It runs `go tool templ generate`, then the tests with the behavior report (`go run ./cmd/testreport`), then builds the `conorganizer` binary and uploads it together with a Goose binary as an artifact. The Goose version is `GOOSE_VERSION` in the job's `env` (v3.28.0, the same as in [migrations.md](migrations.md)). The finished `goose` binary is restored from the Actions cache with a key made of the OS, the arch, `GOOSE_VERSION` and the Go version, and `go install github.com/pressly/goose/v3/cmd/goose@$GOOSE_VERSION` only runs, followed by a cache save, when that key has no cache. The Go module and build caches are saved before the tests and never contain Goose, so bumping Goose or Go is what triggers a rebuild.
- **`deploy-fixed`** runs on pushes to `main`. It is a matrix over `[main, demo, restored]` with `fail-fast: false`, so the three environments deploy in parallel and a failure in one does not cancel the others. This is safe because each environment has its own app dir, binary, service and data dir, and fixed deploys install no generated Caddy or systemd files, so they share no files. Each matrix entry uploads the binary, `goose` and `migrations/*.sql` to `/opt/conorganizer/<name>/` and runs `sudo -n /usr/local/bin/conorganizer-deploy <name>`. CI never runs Goose against a fixed environment; the deployed files are used by the manual `conorganizer-sqlite-migrate` (see [migrations.md](migrations.md)).
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

CI computes the preview's name and host from `github.event.pull_request.number`; the server derives the port from the name:

- name: `<PR_NUMBER>-merge`, for example `482-merge`
- host: `https://<PR_NUMBER>-merge.lekeplassen.regncon.no`
- port: `20000 + PR_NUMBER`, for example `20482`, computed by `conorganizer-deploy`

On each push to the PR, CI uploads the binary as `conorganizer.new`, `goose` and `migrations/*.sql` to the preview's app dir and runs `sudo -n /usr/local/bin/conorganizer-deploy <PR_NUMBER>-merge`. The script then:

1. When the preview has no `database/events.db` yet (or an empty one), makes a SQLite `.backup` snapshot of main's database and checks it with `PRAGMA quick_check`. A stale `events.db-wal` and `-shm` are removed before the clone, so deleting `events.db` by hand is enough to make the next deploy clone again. Separately, when the preview has no `event-images/` directory, copies main's. Data that already exists is kept.
2. Runs `goose up -allow-missing` with the deployed migrations against the preview database, on every deploy, so later pushes to the PR get their new migrations.
3. Renders the unit and the Caddy site from the templates and writes them, as root and atomically, to `/etc/systemd/system/conorganizer-<name>.service` and `/etc/caddy/sites-enabled/conorganizer-<name>.caddy`, then reloads Caddy. If Caddy refuses the new config, the previous site file is put back and the deploy fails.
4. Promotes the new binary (the previous one stays as `conorganizer-<name>.old`), restarts the service and fails the job (printing the last 50 journal lines) if the service does not become active.

Every step inside the app dir and the data dir runs as the service user (see [Service user and ownership](#service-user-and-ownership)). Files it creates in the data dir are `deploy:www-data`. The files CI uploads (`conorganizer.new`, `goose`, `migrations/`), the promoted binary and a new preview's app dir are `deploy:deploy`, the deploy user's own group, which is enough because the unit runs as `User=deploy`. Older app dirs can still show `www-data` from the previous deploy script's recursive `chown`. A PR cannot change its own preview's unit or Caddy site: a change to the templates or to `conorganizer-deploy` is merged to `main` and then installed on the server (see [Deploy scripts](../configuration-as-code/README.md#deploy-scripts)).

After a successful deploy, CI adds a `Preview deployment` block with the URL to the PR description. It only does this once: the block starts with the HTML comment `<!-- preview-deployment-url:start -->`, and CI skips the step when the description already contains it.

## Caddy routing

The root Caddyfile is `configuration-as-code/stow/caddy/etc/caddy/Caddyfile`. Routing for `program.regncon.no`, `main.lekeplassen.regncon.no`, `demo.lekeplassen.regncon.no`, `restored.lekeplassen.regncon.no` and `grafana.regncon.no` is always defined there. `demo` and `restored` proxy straight to `127.0.0.1:19081` and `127.0.0.1:19082`.

Only PR preview sites live in `/etc/caddy/sites-enabled/conorganizer-<name>.caddy`. `conorganizer-deploy` installs them and `conorganizer-cleanup` removes them. The root Caddyfile loads them with:

```caddyfile
import /etc/caddy/sites-enabled/*.caddy
```

### Maintenance page

The Caddyfile defines two snippets: `conorganizer-main` (a reverse proxy to `127.0.0.1:19080`) and `conorganizer-maintenance` (a bilingual English/Norwegian maintenance page served with status 503 and `Cache-Control: no-store`). While the con is over, `program.regncon.no` only redirects (302) to `https://regncon.no/gamle-regncon-program/`, so the maintenance flag has no effect until the previous block is restored from git history. That block serves the maintenance snippet while `/var/lib/conorganizer/maintenance.on` exists and `conorganizer-main` otherwise. Switch with `sudo conorganizer-maintenance-mode on|off`; Caddy checks the flag file on every request, so no reload is needed. `main.lekeplassen.regncon.no` always imports `conorganizer-main`, so main can still be reached there while `program.regncon.no` shows the maintenance page. The migration procedure in [migrations.md](migrations.md) uses this toggle. See [configuration-as-code/README.md](../configuration-as-code/README.md#maintenance-mode) for details.

## Service user and ownership

All Conorganizer services run as `User=deploy` and `Group=www-data`. This covers the main, demo and restored units and the preview template. Database files you put in place by hand, for example during a restore or migration, must keep that ownership, or the service cannot open them.

The `deploy` user is also the SSH user CI logs in as, and it owns `/opt/conorganizer`. Anything it can write is untrusted to root: a preview is PR code running as `deploy`. Its sudo rights are therefore limited to `/usr/local/bin/conorganizer-deploy` and `/usr/local/bin/conorganizer-cleanup` (`/etc/sudoers.d/deploy-conorganizer`, created by `hetzner/scripts/setup_deploy_user.sh` on a fresh server and edited with `sudo visudo -f` after that; check it with `sudo -l -U deploy`). Both are Stow links into the admin-owned server checkout `/srv/configuration-as-code-repo/conorganizer`, which root already trusts for the backup units, and the deploy user cannot write there. The scripts validate the environment name (`main`, `demo`, `restored` or `<PR_NUMBER>-merge`; cleanup only the last) and derive every path from it, so the deploy user can trigger a deploy or a cleanup but cannot change what either does.

Inside `conorganizer-deploy`, root does only what needs root: it checks that `/mnt/HC_Volume_103911252/environments` is a root-owned directory that only root can write (its parent directories must stay that way too; the script checks only the last component, so check them with `namei -l`), creates the preview's top-level data dir there, writes the rendered unit and Caddy site into `/etc`, and calls `systemctl`. Everything inside the app dir and the data dir, from `chmod +x goose` through the database clone, goose and the binary promotion, runs as the service user through `runuser -u deploy -g www-data`. Root never writes to, opens or executes anything under a path the deploy user controls. It only checks with file tests whether the uploaded binary, goose and migrations, the preview's database and `event-images`, and the current binary exist, and those tests read metadata only. No recursive `chown` is needed, because the service user creates or moves every file it later uses. `conorganizer-cleanup` only removes fixed paths built from the validated name, and `rm -rf` does not follow symlinks.

This does not cover everything. Every environment, including main, runs as the same `deploy` user, and the three CI jobs share one SSH key in the `Production` GitHub environment. A same-repo PR can therefore still, as `deploy`, overwrite main's binary or database or read `/etc/conorganizer/env`, from its preview process or by changing its own workflow. That is a production compromise without root. Caddy's admin API on `localhost:2019` is reachable by the `deploy` user as well, so a preview can change Caddy's running config without root. Closing both needs a separate service user and CI identity for previews and the admin API on a Unix socket, which is a planned follow-up.

## Backups

Scheduled backups of main are written to `/mnt/HC_Volume_103911252/backups`:

- `sqlite/` holds compressed database backups (`events-<timestamp>.db.zst`) plus an uncompressed `events-latest.db`, written by `conorganizer-sqlite-backup` every 15 minutes.
- `images/` holds image backups (`event-images-<timestamp>.tar.zst`), written by `conorganizer-images-backup` daily at 03:30.

Both keep 14 days of backups and are run by systemd timers from `configuration-as-code`. See [configuration-as-code/README.md](../configuration-as-code/README.md#backups) for the timers, and [sqlite-and-health.md](sqlite-and-health.md) for why backups use SQLite `.backup`.

After each successful run the scripts write `conorganizer_*` metrics for the Grafana backup panels to node_exporter's textfile directory. See [Backup Metrics](../configuration-as-code/grafana-dashboards/README.md#backup-metrics).

## Preview cleanup

`cleanup-preview` runs `sudo -n /usr/local/bin/conorganizer-cleanup <PR_NUMBER>-merge` on the server; it uploads nothing and checks out nothing. The script removes everything a preview deploy creates, in this order:

1. Stops and disables the service. If the service is still active after that, the script exits without removing anything, so a running preview never loses its files.
2. Removes the unit and its `multi-user.target.wants` link. `systemctl disable` only removes the link while the unit file exists, so the script removes both itself.
3. Removes the Caddy site and the data dir (`database/`, `event-images/`, everything the app wrote there).
4. Reloads systemd, resets the unit's failed state so it stops showing as `not-found`, and reloads Caddy. A failed Caddy reload fails the job, because Caddy then still serves the old route.
5. Removes the app dir.

`/run/conorganizer-<name>/` is the unit's `RuntimeDirectory` and holds the NATS store; systemd removes it when the service stops. Two things stay on purpose: Caddy's certificate for the host, which Caddy deletes from its storage after it expires, and the journal entries of the unit, which follow the journal's own retention.

Every removal tolerates a target that is already gone, so running it twice, after a failed deploy, or for a draft that never got a preview, is safe.

The job runs on `pull_request_target`, not `pull_request`, because GitHub runs no `pull_request` workflow for a PR with merge conflicts, not even when it is closed. `pull_request_target` runs the workflow of the PR's base branch with the secrets, so the job never runs PR code; for a PR into `main` that is main's workflow. A change to `cleanup-preview` therefore only takes effect once it is merged to the base branch, and a change to `conorganizer-cleanup` once the server checkout is updated.

The job needs the `deploy` user's sudoers entry for `/usr/local/bin/conorganizer-cleanup`, which `hetzner/scripts/setup_deploy_user.sh` writes on a fresh server. Check it on the server with `sudo -l -U deploy`. Without it the job fails at the `sudo` line and the preview stays.

To see what the script would remove without changing anything, run it with `--dry-run`. That needs no root and works for any `<PR_NUMBER>-merge` name, whether or not that preview exists:

```bash
conorganizer-cleanup --dry-run <PR_NUMBER>-merge
```

A real run refuses to start unless it runs as root, so it always goes through sudo.

Safety notes:

- The job takes the target name from the PR number (`<PR_NUMBER>-merge`), not from `github.ref_name`. On the close event of a merged PR, `github.ref_name` can resolve to `main`, so never use it to pick the cleanup target.
- `conorganizer-cleanup` refuses every name that is not `<PR_NUMBER>-merge` (one to five digits, no leading zero), so it cannot remove `main`, `demo`, `restored` or anything outside the preview paths.

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

`conorganizer-cleanup` is already installed on the server. Run it for each stale preview, first as a dry run, then for real:

```bash
for name in $stale; do conorganizer-cleanup --dry-run "$name"; done
for name in $stale; do sudo /usr/local/bin/conorganizer-cleanup "$name"; done
```

Each real run reloads Caddy, so don't update the configuration-as-code checkout in the middle of a sweep unless you want its Caddyfile changes live. When the last site file is gone, Caddy's `import /etc/caddy/sites-enabled/*.caddy` only logs a warning, so the reload still succeeds.
