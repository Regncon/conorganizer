# Con Organizer

## Description

This is a spike exploring Go, Data-Star and Templ using the Northstar template.

this project runs

For further details about Templ please visit [Templ](https://templ.guide).

For further details about Data-Star, please visit [Data-Star](https://data-star.dev/).

# Database Issues: events.db Crash or Data Retrieval Errors

> [!NOTE]
>If you encounter issues with the events.db database—such as crashes, errors during loading, or problems retrieving data from tables—follow these steps:
> 1. Delete the `events.db` file
> 2. Run the project using the [Run Project](#run-project) to recreate the database

## IDE Setup
For more information on IDE support, see [Templ Guide: Developer Tools](https://templ.guide/developer-tools/ide-support/).

**Tip:** Recommend installing the suggested extensions.

### NeoVim
#### Templ
> [!WARNING]
> Don't install joerdav/templ.vim.

#### Sql
Use Dadbod for sql support.

```lua
  "tpope/vim-dadbod",
  "kristijanhusak/vim-dadbod-completion",
  {
    "kristijanhusak/vim-dadbod-ui",
    config = function()
      vim.keymap.set("n", "<leader>td", ":DBUIToggle<CR>", { desc = "Toggle dbod" })
    end,
  },
```
https://www.youtube.com/watch?v=ALGBuFLzDSA
https://www.youtube.com/watch?v=NhTPVXP8n7w&t=219s

# Run project
There are two ways to run this project:

1. Directly on Linux/Mac: [Run the project](#linux-setup) using a command in your terminal.
2. Via Docker: Use [Docker](#docker-setup)  to manage the environment.

Note: On Windows, the executable is named main.exe while on Linux it is simply main. To keep things consistent with a single file (main), we run Docker on Windows.

## Linux/Mac setup
1. install [Go](https://go.dev/doc/install)
2. `go install github.com/a-h/templ/cmd/templ@latest`
3. `go install github.com/air-verse/air@latest`
4. Follow the [task install instructions](https://taskfile.dev/installation)

5. Open your terminal and run `task live`

## Docker setup

To build the Docker image and start a container, run:
```console
 docker compose up
```

## Additional Resources

See the  [northstar](https://github.com/zangster300/northstar) for further installation instructions.
