# Con Organizer

## Why

The main purpose of this project is to help the Regncon festival achieve its goals of having all the players get to play at least one game that they are very interested in during the festival.

## Quick Start

Choose your preferred method to run the project:
### Mac/Linux Setup
Just install Go.

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

## Get the Latest Database Backup and Images

> [!NOTE]
> Downloads require `DB_SSH_USER` in your `.env` file or shell environment.

```bash
go tool task download:main
go tool task download:demo
```

## Run a Restored Database Backup

`restored.lekeplassen.regncon.no` runs a selected production database backup.
It is a public, isolated environment: restoring a backup replaces its database
and refreshes its event images from main, without changing main or demo.

On the production server, choose one of the backup files in
`/mnt/HC_Volume_103911252/backups/sqlite` and run:

```bash
sudo conorganizer-sqlite-restore events-20260920T120007Z.db.zst
```

The command verifies the compressed SQLite backup before replacing the restored
environment, then starts `conorganizer-restored.service`. Verify the result at
[https://restored.lekeplassen.regncon.no/](https://restored.lekeplassen.regncon.no/).

For the first rollout, add the DNS record, let the first application deployment
install the restored binary, apply the configuration-as-code changes, then run
the restore command. The first application deployment only installs its binary
until a backup has been selected.

## Run Locally

```bash
go tool task start
go tool task start:demo
```

Then open your browser and navigate to: [http://localhost:8080](http://localhost:8080)

## Run tests

The fist time you run tests you need to create a new schema.sql

```bash
go tool task test
```

> [!TIP]
> Format the `schema.sql` using the Prettier plugin in your IDE to make it look nice.

After that, you can choose run the tests with go tool task or directly with go test:

```bash
go test ./...
```

## IDE Setup

See [Templ Guide: Developer Tools](https://templ.guide/developer-tools/ide-support/) for detailed IDE support information.

Neovim-specific setup lives in [documentation/neovim-setup.md](documentation/neovim-setup.md).

## Troubleshooting

Common issues and solutions:

- **Manual templ generation**: If you encounter issues with Templ, run:

```bash
go tool templ build
```

- **Docker HTTPS port in use**: Check if another service is using the port set
  by `HTTPS_PORT` in `.env` (or port `7331` when it is unset)
- **Build errors**: Run `go mod tidy` to fix dependencies

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

If `templ` was updated, make sure workflow pins match the new version where relevant, especially:

```text
.github/workflows/golangci-lint.yml
```

Look for hardcoded commands like:

```bash
go install github.com/a-h/templ/cmd/templ@v0.3.1020
```

and update the version to match `go.mod`.

## Agent Skills Path Compatibility

Some agents do not discover skills directly from `.agents/skills`.

If that happens, link each skill into that agent's own skills folder (create the folder first if needed).

If you need a true symlink instead (may require admin/dev mode):

```powershell
$agentSkillsFolder = ".codex\skills"  # replace with your agent's skills folder
New-Item -ItemType Directory -Force -Path $agentSkillsFolder | Out-Null
New-Item -ItemType SymbolicLink -Path "$agentSkillsFolder" -Target ".agents\skills"
```

## Additional Resources

- [Northstar Template Documentation](https://github.com/zangster300/northstar)
- [Go Documentation](https://go.dev/doc/)
- [Docker Documentation](https://docs.docker.com/)

