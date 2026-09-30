# Con Organizer

## Why

The main purpose of this project is to help the Regncon festival achieve its goals of having all the players get to play at least one game that they are very interested in during the festival.

## Quick Start

Task, templ and Air are tool dependencies in `go.mod`, so you do not install
them separately. Run every Taskfile command as `go tool task <name>`, and the
other tools as `go tool templ` and `go tool air`.

The app needs no `.env` file to start, but it does need a database and a
writable image directory. `go tool task start` opens `database/events.db` and
refuses to create it when it is missing. If `local-event-images/` (present in
the checkout) is missing or not writable, the app starts in degraded mode. A
fresh checkout has no database, so download data first:

1. Set `DB_SSH_USER` in `.env` (see [`.env`](#env)).
2. `go tool task download:main` (see
   [Get the Latest Database Backup and Images](#get-the-latest-database-backup-and-images)).
3. `go tool task start` (see [Run Locally](#run-locally)).

Choose your preferred method to run the project:
### Mac/Linux Setup
Install Go. Downloads also need `ssh`, `sqlite3` and `tar`, and
`go tool task test` needs `sqlite3`.

### Docker Setup (Recommended for Windows)

Start the application using Docker Compose

```bash
docker compose up --build
```

Docker serves the application through Caddy using HTTPS and HTTP/2. The Go
server stays internal on port `7332`. Caddy uses `HTTPS_PORT` from `.env` as
its public HTTPS port, falling back to `7331` when `HTTPS_PORT` is not set.
With `HTTPS_PORT=7331`, open
[https://localhost:7331](https://localhost:7331).

Caddy starts once the Go server reports healthy, so the first start waits
for templates to generate and the server to compile.

The first time Caddy starts, trust its local certificate authority for your
user account. Keep Docker Compose running, open another terminal, and run the
command for your operating system.

Windows PowerShell:

```powershell
.\scripts\trust-docker-ca.ps1
```

If PowerShell refuses with `running scripts is disabled on this system`, run
the script with a bypass that only applies to that one process:

```powershell
powershell -ExecutionPolicy Bypass -File .\scripts\trust-docker-ca.ps1
```

Linux or macOS:

```bash
sh scripts/trust-docker-ca.sh
```

On Linux the script also adds the certificate to the Chrome and Firefox
certificate databases, which needs `certutil` (`libnss3-tools` on
Debian/Ubuntu, `nss-tools` on Fedora, `nss` on Arch).

Restart the browser after trusting the certificate, then open the HTTPS URL
for the configured `HTTPS_PORT`. The certificate is retained in the
`caddy_data` Docker volume across container rebuilds. Removing that volume,
for example with `docker compose down -v`, creates a new certificate
authority. Run the trust script again afterwards; the old certificate stays
trusted until you remove it from your certificate store.

### Testing from a phone

Set `DEV_LAN_IP` in `.env` to your computer's local network IP address, for
example `DEV_LAN_IP=192.168.1.20`, and restart Docker Compose. Then:

1. Run the trust script once so the certificate is saved to
   `tmp/caddy/conorganizer-caddy-root.crt`.
2. Copy that file to the phone and install it as a trusted certificate
   authority.
   - iOS: open the file, install the profile under Settings → General → VPN &
     Device Management, then enable it under Settings → General → About →
     Certificate Trust Settings.
   - Android: Settings → Security → Encryption & credentials → Install a
     certificate → CA certificate.
3. Open `https://<DEV_LAN_IP>:7331` on the phone, using your `HTTPS_PORT`.

While `DEV_LAN_IP` is set, use `https://localhost:7331` on the computer
itself: connections to an IP address carry no hostname, so Caddy answers
`127.0.0.1` with the LAN address certificate. If the phone cannot connect,
check that the firewall allows inbound connections on the HTTPS port.

### `.env`

`.env` is gitignored and holds only local, per-developer values and secrets.
Task loads it (`dotenv: .env` in `Taskfile.yml`), and the app loads it on a
best-effort basis with godotenv. Known keys:

| Key | Used by |
| --- | --- |
| `DB_SSH_USER` | Required by `go tool task download:*` |
| `HTTPS_PORT`, `DEV_LAN_IP` | Docker Compose and Caddy dev setup |
| `CHECKIN_KEY`, `CHECKIN_SECRET` | Optional. Checkin API credentials, read by `service/checkIn` |
| `GOOSE_DRIVER`, `GOOSE_MIGRATION_DIR` | Optional Goose CLI convenience |
| `PORT` | Optional. HTTP port, defaults to `8080` |
| `LOG_LEVEL` | Optional. `DEBUG`, `INFO`, `WARN` or `ERROR`, defaults to `INFO` |

The Descope project ID is not configuration: it is the hardcoded constant
`authctx.DescopeProjectID` in `service/authctx/authctx.go`. It is not a secret,
since the login page shows it to every browser. Any `DESCOPE_PROJECT_ID` left
in an old `.env` is ignored.

## Get the Latest Database Backup and Images

> [!NOTE]
> Downloads require `DB_SSH_USER` in your `.env` file or shell environment.

```bash
go tool task download:main
go tool task download:demo
```

Add `:db` or `:images` to download only one part, for example
`go tool task download:main:images`. Plain `go tool task download` only prints
usage and exits with status 1, so nothing is downloaded without naming an
environment.

Each task runs `bash scripts/download-environment <main|demo> <all|db|images>`:

- It opens one SSH ControlMaster connection to the production server as
  `DB_SSH_USER`, so you authenticate once per run.
- The database is exported on the server by `conorganizer-export-db <env>`
  and streamed back as a tar. The main export is anonymized. See
  [documentation/sqlite-and-health.md](documentation/sqlite-and-health.md) for
  how the export works.
- The export is extracted into a fresh temporary directory and checked with
  `PRAGMA quick_check`. Only when the check passes is the local database (and
  its `-wal`/`-shm` files) replaced: `database/events.db` for main,
  `database/events-demo.db` for demo. A failed or interrupted download leaves
  your working database alone.
- Images are streamed as a tar from the environment's `event-images`
  directory into the shared `local-event-images/` directory. Existing files are
  never pruned. Image files are named after the random event id, so main and
  demo images do not collide; an event present in both gets the same file
  name, and the last download wins.

`DB_SSH_USER` must be allowed to run `conorganizer-export-db` and to read
`/mnt/HC_Volume_103911252/environments/{main,demo}/database/events.db` and the
matching `event-images` directories. The export only writes to the user's
default temporary directory, not to the backup area.

- **Mac/Linux:** the download runs on your machine and needs `ssh`, `sqlite3`
  and `tar`.
- **Windows:** the download runs inside the Docker container, so keep
  `docker compose up` running. Without Go installed on Windows, run the task in
  the container instead: `docker compose exec webserver go tool task download:main`.
  The container has no SSH key, so ssh asks for your password, once per run.
  After the container has been recreated, it also asks you to confirm the
  server's fingerprint.

Restart the application after downloading a database so it uses the new file.

## Run a Restored Database Backup

`restored.lekeplassen.regncon.no` runs a selected production database backup.
It is a public, isolated environment: restoring a backup replaces its database
and refreshes its event images from main, without changing main or demo.

On the production server, choose one of the backup files in
`/mnt/HC_Volume_103911252/backups/sqlite` and run:

```bash
sudo conorganizer-sqlite-restore events-20260920T120007Z.db.zst
```

The command decompresses the backup and runs `PRAGMA integrity_check` and a
core-table check before replacing the database and event images in
`/mnt/HC_Volume_103911252/environments/restored/`, then starts
`conorganizer-restored.service` (port `19082`, running as `deploy:www-data`).
Verify the result at
[https://restored.lekeplassen.regncon.no/](https://restored.lekeplassen.regncon.no/).

For the first rollout, add the DNS record, let the first application deployment
install the restored binary, apply the configuration-as-code changes, then run
the restore command. The first application deployment only installs its binary
until a backup has been selected.

See [documentation/deployment.md](documentation/deployment.md) for the
environments, ports and CI/CD.

## Run Locally

```bash
go tool task start
go tool task start:demo
```

Then open your browser and navigate to: [http://localhost:8080](http://localhost:8080)

Main and demo data use separate local SQLite files, so there is no command for
switching environments:

| Command | Database | Written by |
| --- | --- | --- |
| `go tool task start` | `database/events.db` | `go tool task download:main` |
| `go tool task start:demo` | `database/events-demo.db` | `go tool task download:demo` |

Both use the shared `local-event-images/` directory. Air passes the paths to
the app as `-dbp <path> -image-path local-event-images`.

## Run tests

```bash
go tool task test
```

`go tool task test` first regenerates `schema.sql` from your local
`database/events.db` (`sqlite3 database/events.db ".schema --indent --nosys"`)
and then runs `go test ./...`. Test databases are built from `schema.sql`
(`service/testdb.go`).

`schema.sql` is committed, so the task overwrites the committed file with your
local database's schema. If your local database is missing a migration,
`schema.sql` loses those tables or constraints, and tests that use them fail.
To fix it:

- Never edit `schema.sql` by hand.
- Apply the pending migrations to the local database, for example
  `goose -dir migrations sqlite3 database/events.db up`, or download a fresh
  database with `go tool task download:main:db`. Then run
  `go tool task test` again.
- Check the `schema.sql` diff before committing.

Plain `go test ./...` uses the committed `schema.sql` as it is:

```bash
go test ./...
```

`go tool task test:report` runs the tests and prints the BDD behavior report
used in CI. See [documentation/automated-tests.md](documentation/automated-tests.md)
for how the automated tests are written.

## Templ

Generated `*_templ.go` and `*_templ.txt` files are gitignored; commit only the
`.templ` sources. They never show up in `git status`, which is expected. They
are generated by the Taskfile (`build:templ`, and `start:templ` while
`go tool task start` runs), during Docker startup (which runs
`go tool task start`), and in CI before tests, build and lint. golangci-lint
also excludes `_templ.go` files.

After editing `.templ` files outside `go tool task start`, run:

```bash
go tool task build
```

Its `build:templ` dependency runs `go tool templ generate` for `components`,
`pages`, `layouts` and `service` before `go build`. Compile errors in the
generated code, such as unused variables or a shadowed `err` in `{{ }}` Go
blocks, only show up after generation.

A bare `return` line in template markup is rendered as text, not as a Go early
exit, so rendering continues. For guard and early-exit logic inside a
component, use a Go code block, as `pages/event/event_page.templ` does:

```templ
if event == nil {
	<p>Event not found</p>
	{{ return nil }}
}
```

## IDE Setup

See [Templ Guide: Developer Tools](https://templ.guide/developer-tools/ide-support/) for detailed IDE support information.

Neovim-specific setup lives in [documentation/neovim-setup.md](documentation/neovim-setup.md).

## Troubleshooting

Common issues and solutions:

- **Manual templ generation**: If you encounter issues with Templ, run:

```bash
go tool templ generate
```

- **Docker HTTPS port in use**: Check if another service is using the port set
  by `HTTPS_PORT` in `.env` (or port `7331` when it is unset)
- **Build errors**: Run `go mod tidy` to fix dependencies
- **Admin route returns 404**: every admin route is registered in
  `SetupAdminRoute` in `pages/admin/admin.go`. Feature handlers such as
  `programPublishingRoute` (`pages/admin/publiser_program.templ`) and
  `puljefordelingStatusRoute` (`pages/admin/puljefordeling.go`) live next to
  their feature and are only reachable when `SetupAdminRoute` calls them. A
  handler and a frontend `@put` call are not enough. Merges can drop these
  calls, so check that every `*Route(adminRouter, ...)` call is still there.

## Migrations

[documentation/migrations.md](documentation/migrations.md)

## Update Dependencies

Update all Go dependencies:

```bash
go get -u
go mod tidy
```

Check what changed:

```bash
git diff go.mod go.sum
```

Verify tool versions used by the repo:

```bash
go tool templ --version
go tool task --version
go tool air -v
```

templ, Task and Air are tool dependencies in `go.mod`, so updating `go.mod`
updates them everywhere: CI and the Dockerfile use the `go.mod` versions.

When changing the Go version in `go.mod`, also update `go-version` in
`.github/workflows/buildAndTest.yml` and `.github/workflows/golangci-lint.yml`,
and the `golang:` image in `Dockerfile`.

If golangci-lint fails in CI with "the Go language version ... used to build
golangci-lint is lower than the targeted Go version", bump the golangci-lint
`version:` pin (and `go-version` if needed) in
`.github/workflows/golangci-lint.yml`. Do not downgrade the project's Go
version.

## Agent Skills Path Compatibility

Some agents do not discover skills directly from `.agents/skills`.

For Claude Code, run `go tool task setup:skills` once per clone. It links
`.claude/skills` to `.agents/skills` (a junction on Windows, so no admin rights
are needed).

For other agents, link the skills into that agent's own skills folder (create the folder first if needed).

If you need a true symlink instead (may require admin/dev mode):

```powershell
$agentSkillsFolder = ".codex\skills"  # replace with your agent's skills folder
New-Item -ItemType Directory -Force -Path $agentSkillsFolder | Out-Null
New-Item -ItemType SymbolicLink -Path "$agentSkillsFolder" -Target ".agents\skills"
```

## Additional Resources

Project documentation in `documentation/`:

- [Deployment and environments](documentation/deployment.md)
- [SQLite and health checks](documentation/sqlite-and-health.md)
- [Migrations](documentation/migrations.md)
- [Live update lifecycle](documentation/live-update-lifecycle.md)
- [Datastar signals](documentation/datastar-signals.md)
- [Pulje status and publishing](documentation/pulje-status-and-publishing.md)
- [Billettholdere](documentation/billettholdere.md)
- [Room assignment](documentation/room-assignment.md)
- [Access control and error pages](documentation/access-control-and-error-pages.md)
- [UI conventions](documentation/ui-conventions.md)
- [Automated tests](documentation/automated-tests.md)
- [Manual tests](documentation/testing/index.md)
- [Neovim setup](documentation/neovim-setup.md)
- [Hetzner admin](hetzner/hetzner-admin.md)

External:

- [Northstar Template Documentation](https://github.com/zangster300/northstar)
- [Go Documentation](https://go.dev/doc/)
- [Docker Documentation](https://docs.docker.com/)
