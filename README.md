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

Then open your browser and navigate to: [http://localhost:8080](http://localhost:8080)

## Get the Latest Database Backup and Images

> [!NOTE]
> Downloads require `DB_SSH_USER` in your `.env` file or shell environment.

```bash
go tool task download:main
go tool task download:demo
```

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

- **Port in use**: Check if another service is using port 8080
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
