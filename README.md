# Con Organizer

## Description

This is a spike exploring go and datastar using the northstar template.

## Ide septup

https://templ.guide/developer-tools/ide-support/
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

## Docker setup

To build the docker image run the following command:

```bash
 docker build --rm -t my-dev-environment .
```

To run the docker image use the run-docker bash or powershell script.

## Links

Se the [northstar](https://github.com/zangster300/northstar) README for installation instructions.
