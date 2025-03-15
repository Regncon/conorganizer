# Con Organizer

## Description

This is a spike exploring Go, Data-Star and Templ using the Northstar template.

this project runs

For more details about Templ, please visit [Templ](https://templ.guide).

For more details about Data-Star, please visit [Data-Star](https://data-star.dev/).

# Database Issues: events.db Crash or Data Retrieval Errors

> [!NOTE]
>If you encounter issues with the events.db database—such as crashes, errors during loading, or problems retrieving data from tables—follow these steps:
> 1. Delete the `events.db` file
> 2. Run the project using the [Run Project](#run-project) to recreate the database
> 3. Open terminal and run `sqlite3 "<<path to repo>>\conorganizer\events.db" < seed_data.sql`


## IDE Setup
For more information on IDE support, see [Templ Guide: Developer Tools](https://templ.guide/developer-tools/ide-support/).

**Tip:** We recommend installing the suggested extensions.

### NeoVim

#### Templ

> [!WARNING]
> Don't install joerdav/templ.vim.

#### SQL Support
For SQL support in NeoVim, use Dadbod. Add the following plugins to your configuration:

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
For video tutorials on using Dadbod, check out:

* [video 1](https://www.youtube.com/watch?v=ALGBuFLzDSA)
* [video 2](https://www.youtube.com/watch?v=NhTPVXP8n7w&t=219s)


# Run project
There are two ways to run this project:

1. Directly on Linux/Mac: Follow the instructions in the [Linux/Mac Setup](#linux-setup) section to run the project via the terminal.
2. Via Docker: Follow the instructions in the [Docker Setup](#docker-setup) section to manage the environment.

> [!NOTE]
>On Windows, the executable is named `main.exe `while on Linux it is simply `main`. To keep things consistent with a single file (`main`), we recommend running Docker on Windows.



## Linux/Mac setup
1. install [Go](https://go.dev/doc/install)
2. Install Templ by running: `go install github.com/a-h/templ/cmd/templ@latest`
3. Install Air by running: `go install github.com/air-verse/air@latest`
4. Follow the [task install instructions](https://taskfile.dev/installation)

5. Open your terminal and run `task live`

## Docker Setup

To build the Docker image and start a container, run:
```console
 docker compose up
```

## Additional Resources

For further installation instructions, see the [northstar](https://github.com/zangster300/northstar)
