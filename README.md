# Con Organizer

## Table of Contents

1. [Why](#why)
2. [Quick Start](#quick-start)
    - [Mac/Linux Setup](#maclinux-setup)
    - [Docker Setup](#docker-setup-recommended-for-windows)
    - [Direct Installation](#direct-installation)
3. [Access the Application](#access-the-application)
4. [Database Issues: events.db Troubleshooting](#database-issues-eventsdb-troubleshooting)
5. [IDE Setup](#ide-setup)
    - [Troubleshooting](#troubleshooting)
6. [Migrations](#migrations)
7. [Agent Skills Path Compatibility](#agent-skills-path-compatibility)
8. [Update Dependencies](#update-dependencies)
9. [Additional Resources](#additional-resources)

## Why

The main purpose of this project is to help the Regncon festival achieve its goals of having all the players get to play at least one game that they are very interested in during the festival.

## Quick Start

Choose your preferred method to run the project:
### Mac/Linux Setup
Just install Go.

### Docker Setup (Recommended for Windows)

1. **Start the application using Docker Compose**

```bash
docker compose up --build
```

2. Once the containers are up and running, head over to [Access the Application](#access-the-application) to view the app in your browser.

> [!NOTE]
> Docker is a platform that allows you to run applications in containers. It handles all dependencies and environment configuration automatically, making it ideal for Windows users.

### Direct Installation


## Access the Application:

```bash
go tool task start
```

Then open your browser and navigate to: [http://localhost:8080](http://localhost:8080)

## Database Issues: events.db Troubleshooting

> [!NOTE]
> To get the latest database backup and all images from production, run:
> The database download requires `DB_SSH_USER` in your `.env` file or shell environment.
> The database task creates a temporary SQLite backup snapshot on the server; it does not copy the live WAL-mode database file directly.

```bash
go tool task download
```

Production SQLite operational notes are in [documentation/sqlite-production.md](documentation/sqlite-production.md).

To get the latest schema of the database, run:

```bash
sqlite3 database/events.db ".schema --indent" > schema.sql
```

> [!TIP]
> Format the `schema.sql` using the Prettier plugin in your IDE to make it look nice.

## IDE Setup

See [Templ Guide: Developer Tools](https://templ.guide/developer-tools/ide-support/) for detailed IDE support information.

Neovim-specific setup lives in [documentation/neovim-setup.md](documentation/neovim-setup.md).

### Troubleshooting

Common issues and solutions:

- **Manual templ generation**: If you encounter issues with Templ, run:

```bash
go tool templ build
```

- **Port in use**: Check if another service is using port 8080
- **Database errors**: See [Database Issues](#database-issues-eventsdb-troubleshooting)
- **Build errors**: Run `go mod tidy` to fix dependencies

## Migrations

Migration notes and the production database update runbook live in [documentation/migrations.md](documentation/migrations.md).

Migrations are manual only. Do not add automatic migrations to application startup, health checks, readiness checks, or systemd startup.

## Agent Skills Path Compatibility

Some agents do not discover skills directly from `.agents/skills`.

If that happens, link each skill into that agent's own skills folder (create the folder first if needed).

If you need a true symlink instead (may require admin/dev mode):

```powershell
$agentSkillsFolder = ".codex\skills"  # replace with your agent's skills folder
New-Item -ItemType Directory -Force -Path $agentSkillsFolder | Out-Null
New-Item -ItemType SymbolicLink -Path "$agentSkillsFolder" -Target ".agents\skills"
```


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

## Additional Resources

- [Northstar Template Documentation](https://github.com/zangster300/northstar)
- [Go Documentation](https://go.dev/doc/)
- [Docker Documentation](https://docs.docker.com/)
