# Editor setup for templ

The repository pins templ in `go.mod`, and CI, the Taskfile and the Dockerfile run that version as `go tool templ`. Editors start their own templ binary for the language server, so the editor can drift from the pinned version. This page covers VS Code. Neovim lives in [neovim-setup.md](neovim-setup.md).

The machine that runs your editor needs [Go](https://go.dev/dl/) and [Node.js](https://nodejs.org/en/download) with `npm`, on every platform. Node is only there to install Prettier; the app itself does not use it.

## Keep the templ binary on the go.mod version

The VS Code extension (`a-h.templ`) picks the templ binary in this order: the `templ.executablePath` setting, then `go tool -n templ` run from the workspace, then `templ` on `PATH`, then a few well-known paths such as `~/go/bin`.

On Windows the `go tool -n templ` step fails when Go is installed under `C:\Program Files\Go`. The extension runs the unquoted path through the shell, so `C:\Program` is not recognized as a command and the extension silently falls back to the `templ.exe` on `PATH`. If that binary is older than `go.mod`, formatting and diagnostics in VS Code differ from what `go tool templ` and CI do.

Install the pinned version globally so the fallback is the right version. The task reads the version from `go.mod`, installs it, installs the Prettier version CI uses, and then compares both with what is on `PATH`:

```bash
go tool task setup:editor
```

Run it again whenever `go.mod` bumps templ. To see whether anything is off without installing, run:

```bash
go tool task check:editor
```

Setting `templ.executablePath` in your user settings is the alternative when you want to point VS Code at one specific binary.

## Linux and macOS

The same task does everything there. templ runs Prettier through `sh`, where its default command works as is, so the Windows notes below do not apply.

## Format on save

`.vscode/settings.json` makes `a-h.templ` the formatter for `.templ` files and turns on format on save, so saved files come out the way `go tool templ fmt` writes them.

CI and `make check` run `go tool task check:templ-fmt`, which fails when a `.templ` file is not formatted. Both need Prettier on `PATH`. The version is `PRETTIER_VERSION` in the Taskfile; `go tool task setup:prettier` installs it, locally and in CI, so bump it in that one place. `go tool task fmt:templ` formats everything in place.

## Prettier for `<style>` and `<script>` blocks

`templ fmt` formats the templ and Go parts by itself. The bodies of `<style>` and `<script>` elements, and constant attribute values, are formatted by Prettier when `prettier` or `prettierd` is on `PATH`. Without Prettier those parts are left as they are, which is how the CSS indentation inside `.templ` files drifted. Install the version CI uses:

```bash
go tool task setup:prettier
```

templ runs `prettier --use-tabs --stdin-filepath $TEMPL_PRETTIER_FILENAME`. On Windows that command goes through `cmd.exe`, where `$TEMPL_PRETTIER_FILENAME` does not expand, so Prettier is asked to format a file literally named `$TEMPL_PRETTIER_FILENAME` and fails with "No parser could be inferred". `.prettierrc` has an override that maps that literal file name to the HTML parser, which makes the default command work on Windows as long as Prettier finds the repository's `.prettierrc`. If you run `templ fmt` from outside the repository, pass the Windows form of the variable instead:

```bash
go tool templ fmt -prettier-command "prettier --use-tabs --stdin-filepath %TEMPL_PRETTIER_FILENAME%" .
```

The VS Code extension cannot pass `-prettier-command` to the language server, so on Windows the language server depends on the `.prettierrc` override. The language server starts in the workspace folder, so the override is found as long as the checkout you have open contains it. When Prettier is installed but its command fails, the language server logs "prettier failure" and format on save does nothing at all. Set `templ.log` to a file path to see that log.
