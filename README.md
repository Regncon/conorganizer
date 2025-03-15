# Con Organizer

## Description

This is a spike exploring Go, Data-Star and Templ using the Northstar template.

For more details, visit:
- [Templ Documentation](https://templ.guide)
- [Data-Star Documentation](https://data-star.dev/)

## Quick Start

Choose your preferred method to run the project:

1. **Docker (Recommended for Windows)**:
   ```bash
   docker compose up
   ```

> [!NOTE]
> The Docker setup handles all dependencies and environment configuration automatically.

2. **Direct Installation**: Follow the [Linux/Mac Setup Guide](#linuxmac-setup-guide) below.

## Database Issues: events.db Troubleshooting

> [!NOTE]
> If you encounter database issues (crashes, loading errors, or data retrieval problems):
> 1. Delete the `events.db` file
> 2. Restart the project to recreate the database
> 3. Seed the database:
>    ```bash
>    sqlite3 "./events.db" < seed_data.sql
>    ```

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
* [Basic Setup and Usage](https://www.youtube.com/watch?v=NhTPVXP8n7w)
* [Advanced Features](https://www.youtube.com/watch?v=ALGBuFLzDSA)

## Linux/Mac Setup Guide

> [!NOTE]
> Windows users should use [Docker Setup](#docker-setup) for consistency.

### Prerequisites

#### 1. Required Tools

| Tool | Description | Installation Command | Expected Version |
|------|-------------|---------------------|------------------|
| [Go](https://go.dev/doc/install) | Programming language | Follow installation guide | ≥ 1.21 |
| [Templ](https://templ.guide) | Template engine | `go install github.com/a-h/templ/cmd/templ@latest` | ≥ 0.2.0 |
| [Air](https://github.com/cosmtrek/air) | Live reload tool | `go install github.com/air-verse/air@latest` | Any |
| [Task](https://taskfile.dev/installation) | Task runner | Follow installation guide | ≥ 3.0 |

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

3. **Access the Application**:
   ```
   http://localhost:7331
   ```

### Troubleshooting

> [!TIP]
> Common issues and solutions:
> - **Tool not found**: Ensure `$HOME/go/bin` is in your PATH
> - **Port in use**: Check if another service is using port 7331 or 8080
> - **Database errors**: See [Database Issues](#database-issues-eventsdb-troubleshooting)
> - **Build errors**: Run `go mod tidy` to fix dependencies

## Additional Resources

- [Northstar Template Documentation](https://github.com/zangster300/northstar)
- [Go Documentation](https://go.dev/doc/)
- [Docker Documentation](https://docs.docker.com/)
