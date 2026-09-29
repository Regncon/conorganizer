# Neovim Setup

See [Templ Guide: Developer Tools](https://templ.guide/developer-tools/ide-support/) for general IDE setup.

## Templ Support

> [!WARNING]
> Do not install `joerdav/templ.vim` - it is deprecated.

The repository pins templ as a Go tool in `go.mod`. To make the templ language server match that version, configure it to run `go tool templ lsp` from the repository root instead of a globally installed `templ`.

## SQL Support with Dadbod

Add these plugins to your Neovim configuration:

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

The local development database is `database/events.db`. Add it in Dadbod UI with the connection URL `sqlite:database/events.db` (relative to the repository root).

Helpful Dadbod tutorials:

- [Basic Setup and Usage](https://www.youtube.com/watch?v=NhTPVXP8n7w)
- [Advanced Features](https://www.youtube.com/watch?v=ALGBuFLzDSA)
