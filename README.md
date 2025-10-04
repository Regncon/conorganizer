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
8. [Additional Resources](#additional-resources)

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

````bash
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
````

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
    task live
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
> Make sure that you can do all of the following steps before you start. These actions require goose, flyctl and wget be installed.

Since our site and services are linked to S3 it's important that we update all the relevant secrets and variables to point to this new database before we run the manual deployment. These are the steps we need to do when running this manually

1. Check our S3 bucket for the current active `dbPrefix` -> https://fly.io/dashboard/regncon/tigris
2. Update the `DB_PREFIX` variable

- local `.env` files
- litestream config in `.flyio/litestream.yml`
- Regncon secret: https://fly.io/apps/regncon/secrets
- Backup-service secret: https://fly.io/apps/backup-service/secrets

3. Take down the website
4. Run `task download` and apply migrations with `goose up`
5. Start deployment with `flyctl deploy`
    ```console
    flyctl -o regncon -a regncon --dockerfile Dockerfile.migration
    ```
6. ???
7. Profit

## Additional Resources

- [Northstar Template Documentation](https://github.com/zangster300/northstar)
- [Go Documentation](https://go.dev/doc/)
- [Docker Documentation](https://docs.docker.com/)
  .
