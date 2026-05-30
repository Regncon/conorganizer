# Con Organizer

## Table of Contents

1. [Description](#description)
2. [Quick Start](#quick-start)
    - [Docker Setup](#docker-setup-recommended-for-windows)
    - [Direct Installation](#direct-installation)
3. [Database Issues: events.db Troubleshooting](#database-issues-eventsdb-troubleshooting)
4. [IDE Setup](#ide-setup)
    - [NeoVim Configuration](#neovim-configuration)
5. [Troubleshooting](#troubleshooting)
6. [Linux/Mac Setup Guide](#linuxmac-setup-guide)
    - [Prerequisites](#prerequisites)
    - [Verification and Startup](#verification-and-startup)
7. [Migrations](#-Migrations-with-goose)
    - [Running Goose manually](#running-goose-manually)
    - [Pushing migrations to prod](#pushing-migrations-to-prod-and-services)
    - [Step by step](#step-by-step-to-update-db)
8. [Agent Skills Path Compatibility](#agent-skills-path-compatibility)
9. [Update go dependencies](#update-dependencies)
10. [Additional Resources](#additional-resources)

## Description

This is a spike exploring Go, Data-Star and Templ using the Northstar template.

For more details, visit:

- [Templ Documentation](https://templ.guide)
- [Data-Star Documentation](https://data-star.dev/)

## Quick Start

Choose your preferred method to run the project:

### Docker Setup (Recommended for Windows)

1. **Start the application using Docker Compose**

```bash
docker compose up --build
```

2. Once the containers are up and running, head over to [Access the Application](#access-the-application) to view the app in your browser.

> [!NOTE]
> Docker is a platform that allows you to run applications in containers. It handles all dependencies and environment configuration automatically, making it ideal for Windows users.

### Direct Installation

Follow the [Linux/Mac Setup Guide](#linuxmac-setup-guide) below.

## Database Issues: events.db Troubleshooting

> [!NOTE]
> To get the latest backup of the database and all the images from prod run:
> The database download requires `DB_SSH_USER` in your `.env` file or shell environment.

```bash
task download
```

To get the latest schema of the database, run:

```bash
sqlite3 database/events.db ".schema --indent" > schema.sql
```

> [!TIP]
> Format the `schema.sql` using the prettier plugin in your IDE to make it look nice.

## IDE Setup

See [Templ Guide: Developer Tools](https://templ.guide/developer-tools/ide-support/) for detailed IDE support information.

### NeoVim Configuration

#### Templ Support

> [!WARNING]
> Do not install `joerdav/templ.vim` - it is deprecated.

#### SQL Support with Dadbod

Add these plugins to your NeoVim configuration:

```lua
{
  "tpope/vim-dadbod",
  "kristijanhusak/vim-dadbod-completion",
  {
    "kristijanhusak/vim-dadbod-ui",
    config = function()
      vim.keymap.set("n", "<leader>td", ":DBUIToggle<CR>", { desc = "Toggle Dadbod UI" })
    end,
  },
}
```

Helpful Dadbod tutorials:

- [Basic Setup and Usage](https://www.youtube.com/watch?v=NhTPVXP8n7w)
- [Advanced Features](https://www.youtube.com/watch?v=ALGBuFLzDSA)

### Troubleshooting

Common issues and solutions:

- **Manual generate templ**: If you encounter issues with Templ, run:

```bash
templ generate && go build -buildvcs=false -o tmp/main .
```

- **Tool not found**: Ensure `$HOME/go/bin` is in your PATH
- **Port in use**: Check if another service is using port 7331 or 8080
- **Database errors**: See [Database Issues](#database-issues-eventsdb-troubleshooting)
- **Build errors**: Run `go mod tidy` to fix dependencies

## Linux/Mac Setup Guide

> [!NOTE]
> Windows users should use [Docker Setup](#docker-setup-recommended-for-windows) for consistency.

### Prerequisites

#### 1. Required Tools

| Tool                                      | Description          | Installation Command                               |
| ----------------------------------------- | -------------------- | -------------------------------------------------- |
| [Go](https://go.dev/doc/install)          | Programming language | Follow installation guide                          |
| [Templ](https://templ.guide)              | Template engine      | `go install github.com/a-h/templ/cmd/templ@latest` |
| [Air](https://github.com/cosmtrek/air)    | Live reload tool     | `go install github.com/air-verse/air@latest`       |
| [Task](https://taskfile.dev/installation) | Task runner          | Follow installation guide                          |

#### 2. Shell Configuration

<details>
<summary>Bash Setup (Linux/macOS)</summary>

```bash
# Add to ~/.bashrc (Linux) or ~/.bash_profile (macOS)
echo 'export PATH=$PATH:$HOME/go/bin' >> ~/.bashrc  # or ~/.bash_profile for macOS

# Apply changes
source ~/.bashrc  # or source ~/.bash_profile for macOS
```

</details>

<details>
<summary>Zsh Setup</summary>

```bash
# Add Go binaries to PATH
echo 'export PATH=$PATH:$HOME/go/bin' >> ~/.zshrc

# Apply changes
source ~/.zshrc
```

</details>

### Update templ

```bash
go install github.com/a-h/templ/cmd/templ@latest
```

### Update Go Modules

```bash
go get -u
go mod tidy
```

### Verification and Startup

1. **Verify Tool Installation**:

    Check Go installation:

    ```bash
    go version
    ```

    Check Templ installation:

    ```bash
    templ version
    ```

    Check Air installation:

    ```bash
    air -v
    ```

    Check Task installation:

    ```bash
    task --version
    ```

> [!TIP]
> Each command should return a version number. If any command fails:
>
> 1. Ensure the tool is installed correctly
> 2. Verify your PATH includes Go binaries
> 3. Try reopening your terminal

2. **Start Development Server**:
    ```bash
    go tool task start
    ```

> [!NOTE]
> This will start the server with hot-reload enabled.
> Any code changes will automatically trigger a rebuild.

### Access the Application:

```
http://localhost:7331
```

## Migrations with goose

> [!NOTE]
> Goose will try to read some basic variables from `.env`, make sure that this file is updated with the most recent version from discord before running any commands.

We're using [Goose](https://pressly.github.io/goose/) in our migration process for its simplicity and reliability. While Goose is available as a go dependency for programatically migrating databases, we're mostly using its CLI tool for manual updates.

### Running Goose manually

> [!WARNING]
> Before running Goose, run `task download` to fetch the newest version of the database!
> You can install Goose CLI tool from [here](https://pressly.github.io/goose/installation/), afterwards you should have `goose` globally available in your terminal.

To create a new migration file you can run the following command, read [here](https://pressly.github.io/goose/documentation/annotations/) for more annotation examples.

```console
goose create <briefly describe changes> sql
```

After you've added your migration files you can use the keywords `up` or `down` to handle the migrations

```console
goose up
```

### Pushing migrations to prod and services

> [!CAUTION]
> Make sure that you can do all of the following steps before you start. These actions require goose, account on server.

1. Run goose on local db (preferably a copy)
2. make a backup on server
3. upload to server.
5. profit

### Step By Step To Update DB

#### 1. Find The Correct Service Name

```bash
systemctl list-units --type=service | grep -i conorganizer
systemctl list-unit-files | grep -i conorganizer
```

#### 2. Inspect The Service

Set the service name:

```bash
SERVICE=conorganizer-{service-name}.service
```

Check the service command and find the mounted database path:

```bash
systemctl show "$SERVICE" -p ExecStart --value | fold -s -w 120
```

Look for the host path that contains `events.db` or maps the database folder into the app.

Example:

```text
/mnt/DC_Hardware_13377331/308-merge/database:/app/database
```

In that case the database path is:

```bash
/mnt/DC_Hardware_13377331/1337-merge/database/events.db
```

#### 3. Replace The Database

Set the paths:

```bash
DB_PATH="{mounted path to db}/events.db"
UPLOADED_DB="/home/{account-name}/events.db"
```

Stop the service:

```bash
sudo systemctl stop "$SERVICE"
```

Back up the current database:

```bash
sudo cp -a "$DB_PATH" "$DB_PATH.bak.$(date +%Y%m%d-%H%M%S)"
```

Move the uploaded database into place:

```bash
sudo mv "$UPLOADED_DB" "$DB_PATH"
```

Fix ownership and permissions:

```bash
sudo chown deploy:deploy "$DB_PATH"
sudo chmod 644 "$DB_PATH"
```

Start the service again:

```bash
sudo systemctl start "$SERVICE"
sudo systemctl status "$SERVICE"
```

Check logs:

```bash
journalctl -u "$SERVICE" -n 100 --no-pager
```

<!--
!!!Dont run this unless you know all caveats, this can affect prod negatively!!!
docker compose -f compose-restore.yaml down && docker image rm regncon-migration
docker compose -f compose-restore.yaml up
-->

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
go get -u ./...
go mod tidy
```

Check what changed:

```bash
git diff go.mod go.sum
```

Verify tool versions used by the repo:

```bash
go tool templ version
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
