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
You have two options when it comes to running this poject, you can either use the bash or powershell script or docker compose cli.

### Script method
First build an image with the following command, then execute either [run-docker.sh](run-docker.sh) or [run-docker.ps1](run-docker.ps1) script depeding on your OS.
```console
 docker build --rm -t my-dev-environment .
```

### Compose method
To build the docker image and mount a container, run the following command:
```console
 docker compose up
```

## Links

Se the [northstar](https://github.com/zangster300/northstar) README for installation instructions.
