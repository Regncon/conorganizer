---
name: conorganizer-worktrees
description: Conorganizer's rules for git worktrees and session branch names. Use whenever you create, enter, name, list or remove a worktree (git worktree add or remove, EnterWorktree, ExitWorktree, the superpowers using-git-worktrees skill, start_session with a branch_name), start work that needs isolation from the main checkout, or set up a fresh .claude/worktrees checkout to build, test or run under docker compose. Use it in the first turn of every session that finds itself on a generated branch such as claude/lucid-kirch-b765f7 or claude/worktree-skill-2cec9f, before the first commit, and when the user says worktree, branch name, EnterWorktree, rename the branch, fix the name, clean up the worktree or the PR merged. Use it even when the worktree is not the main task, since a worktree here has no templ watcher and task test rewrites schema.sql.
---

# Conorganizer Worktrees

One task, one worktree under `.claude/worktrees/`, one branch named after the task. The desktop app names a session's branch `claude/<adjective>-<surname>-<hex>` or `claude/<requested>-<hex>`, so the name GitHub and `git branch` show says nothing about the work; this skill fixes that, and makes every fresh worktree in this repo build the same way. AGENTS.md "Dev setup" and "Git" hold the short rules and are not repeated here. `references/app-behavior.md` is the evidence for every claim about the app, with the rename test.

## Branch names

`claude/<issue>-<slug>` when the task has a GitHub issue, `claude/<slug>` when it has none.

| Task | Branch |
| --- | --- |
| Issue #612, the new level colors | `claude/612-level-colors` |
| Issue #731, the pulje closing warning | `claude/731-pulje-closing-warning` |
| A skill for PR descriptions, no issue | `claude/pr-skill` |
| A sortable billettholder list on the admin page, no issue | `claude/sortable-billettholder-list` |

- The slug is two to five English words in kebab-case, taken from the issue title or the task as the user put it. Domain words from `domeneordbok.md` stay Norwegian without æøå (`pulje`, `billettholder`, `forstevalg`), as identifiers do. A name the code uses today stays as it is (`event`, not `arrangement`, until the planned rename lands), so the branch reads like the code it touches.
- The number is a GitHub issue number, never a PR number. When the user names the task but no issue, run `gh issue list --state all --search "<words>"` once; take the number when an issue matches, and go without one when none does, rather than inventing one. A closed issue the user points at still counts.
- Keep the `claude/` prefix. It is the app's Branch prefix setting (`ccd_settings get_settings`, key `branch_prefix`) and it groups every session branch under one namespace; change the setting only when the user asks.
- No hex suffix, no date, no author. The suffix is the app's collision guard, and a name that says what the task is rarely collides; check `git branch --list 'claude/<slug>*'` and `git ls-remote --heads origin 'claude/<slug>*'` anyway, and add a word rather than a number when the name is taken.

## A session on a generated branch

The app adds `-<6 hex>` to every session branch, also to a requested one, and picks a word pair when no name was requested, so `claude/lucid-kirch-b765f7` and `claude/worktree-skill-2cec9f` both need the rename. Do it in the first turn, before the first commit:

```bash
git branch --show-current
git branch -m claude/lucid-kirch-b765f7 claude/612-level-colors
git worktree list
```

- Before the first commit and push, because a rename after the push leaves a dead remote branch, and a branch with an open PR keeps its name (deleting the remote head closes the PR).
- `git branch -m` is safe: the worktree follows the branch, `git status` and `git worktree list` are right at once, and the push and `gh pr create` use the new name. The app had set no upstream, so there is nothing to re-point.
- The app's session record keeps the old name (the sidebar label, and the branch a resumed session checks out after a worktree cleanup). Say so in one line of the reply, so the label does not surprise the user.
- Never rename or move the folder (`.claude/worktrees/adoring-saha-f65bac`). It is the session's working directory and the app's handle on the worktree, and the app recycles folders across sessions, so its name was never meant to match the branch.
- When you start a session with `start_session`, pass `branch_name: "<issue>-<slug>"`. The app still adds the suffix, but the new session's rename then only drops it, and the label is readable from the start.

## Creating a worktree yourself

Create it from the main checkout: `EnterWorktree` refuses inside a worktree session, and a path given from the main checkout keeps every worktree under `.claude/worktrees/`.

1. `EnterWorktree` with `name: "<issue>-<slug>"` when the tool is available. It creates `.claude/worktrees/<issue>-<slug>` on branch `worktree-<issue>-<slug>` from `origin/main`. Then rename the branch: `git branch -m worktree-612-level-colors claude/612-level-colors`. `ExitWorktree` with `remove` is documented to delete the worktree and its branch; whether it finds the renamed branch is unverified, so after an exit run `git branch --list 'claude/*' 'worktree-*'` and delete what it left.
2. Otherwise:

```bash
git fetch origin
git worktree add --no-track .claude/worktrees/612-level-colors -b claude/612-level-colors origin/main
```

The folder name is the branch's slug, so `git worktree list` reads like the PR list. `.claude/worktrees/` because it is gitignored and where the app, this skill and the cleanup look; the `.worktrees/` default of the superpowers using-git-worktrees skill does not apply here. `origin/main` after a fetch, because the local `main` of a long-running checkout is often behind. `--no-track`, because git otherwise makes the new branch track `origin/main`, and a plain `git push` then offers to push to `main`; the first push is `git push -u origin claude/612-level-colors`. If the folder already exists, it is the empty shell of a removed worktree (see the cleanup section); clear it first.

Then set the worktree up before building or testing.

## Setting up a fresh worktree

From the worktree:

```bash
bash .agents/skills/conorganizer-worktrees/scripts/setup-worktree.sh
```

It generates the templ code, links `.claude/skills`, builds every package and says what is still missing. The same steps by hand, and why each exists:

1. `go tool templ generate -path <dir> -log-level error` for `components`, `pages`, `layouts` and `service`, as `build:templ` in the Taskfile does. The `*_templ.go` files are gitignored and the watcher runs only where `go tool task start` runs, so without this `go build` fails in every page package. Repeat it after each `.templ` edit here. Never run it in the main checkout while the dev server is up.
2. `go tool task setup:skills`. The `.claude/skills` link is gitignored; the app creates it in its own worktrees, a worktree you created has nothing, and without it this skill and the other repo skills do not load there. The app also copies `.claude/settings.local.json` (its permission allowlist) into its own worktrees; a worktree you create has none, so expect permission prompts the main checkout does not show.
3. `go build ./...` and `go vet ./...` for a baseline, `go test ./<pkg>/...` for the packages you touch, `go test ./...` for the whole suite (about a minute). Never `go tool task test`, `go tool task build`, `make test` or `make check` here: `task test` runs the host `sqlite3` on `database/events.db` and rewrites `schema.sql`, and in a worktree without a database it empties it. If `schema.sql` shows up in `git status`, `git checkout -- schema.sql`.
4. Tests need no database; `service/testdb.go` builds one from `schema.sql`.

Running the app is a separate step, only when the task needs the real pages:

- `.env` is gitignored. Copy it from the main checkout and never print it.
- `database/events.db` and `local-event-images/` are gitignored too. Copy the database with sqlite's `.backup`, never `cp`, because a WAL file left behind corrupts the copy. While Docker runs, do it inside the container, since the host `sqlite3` on that file makes the container fail with "unable to open database file (14)": `docker compose exec -T webserver sqlite3 database/events.db ".backup .claude/worktrees/<name>/database/events.db"`. With Docker stopped, the host `sqlite3` with the same `.backup` from the main checkout is fine. Copy the images folder as files.
- Native, side by side with the main checkout's server, from the worktree: `go build -tags dev -o bin/main . && PORT=8090 ./bin/main -dbp database/events.db -image-path local-event-images -nats-store-dir data/nats`. The `dev` tag serves `static/` from disk without caching, so CSS edits show on reload, and the port keeps the main checkout's 8080 free. The built-in browser opens plain `http://localhost:8090`.
- Docker: `docker compose -p conorganizer up -d --build` from the worktree replaces the main checkout's containers with ones that serve the worktree. The same project name reuses the Go caches and the trusted Caddy CA, so https://localhost:7331 keeps working, but the main checkout's server is gone until the same command runs there again. Say so before doing it. `--build` is needed whenever `go.mod` bumps templ, air or task, because the Dockerfile compiles them into the image.

## Cleaning up after the merge

From the main checkout, not from the worktree's own folder:

```bash
git fetch --prune origin
gh pr view <number> --json state,mergedAt,headRefName
git worktree remove .claude/worktrees/612-level-colors
git branch -d claude/612-level-colors
git push origin --delete claude/612-level-colors
```

- Confirm the merge first, with `gh pr view` or `git branch --merged origin/main`. A force-less `remove` and `-d` then succeed, and a refusal means there is work to look at, not a reason for `--force` or `-D`. When `gh` finds no PR for the branch but `git merge-base --is-ancestor <branch> origin/main` holds, nothing is lost by removing it: do it and say that no PR was found. When the branch has commits that are not on `origin/main`, stop and ask.
- Delete the remote branch only when `git ls-remote --heads origin <branch>` still lists it; GitHub often deletes it on merge.
- Check that the folder is gone with `ls .claude/worktrees/<name>`. `git worktree remove` exits 0 and still leaves `.claude/worktrees/<name>/.claude/skills` behind when the worktree had the skills junction: an empty shell with a dangling link, which blocks the next `git worktree add` of the same name, and which `test -e` does not see. Remove the junction with a plain `rm` on the link (never a recursive delete through it) and the two folders with `rmdir`.
- On Windows, when the worktree is the folder your own session runs in, `git worktree remove` unregisters it and empties it but cannot delete the folder, because it is the process's working directory. Say so, and that the empty folder goes once the session is closed; a new message in that chat makes the app check the branch out again.
- Keep the chat. Never archive, clear or delete a session to get rid of its worktree: in auto mode the app archives without a card, and the user keeps finished chats as reference.
- Leave other sessions' worktrees, their branches and the shared stash alone. `git worktree prune` and bulk cleanups are the user's call; the app's `clean_up_worktrees` does them with a card.
- `ExitWorktree` with `remove` is the same cleanup for a worktree `EnterWorktree` created in this session; use `keep` while the work is not merged.

## Pointers

- AGENTS.md "Dev setup" (templ, `task test`, other sessions' checkouts, searches scoped to `git ls-files`) and "Git" (merge main in, no rebase; commit subjects).
- `references/app-behavior.md`: what the app and Claude Code do with names, the rename test, the EnterWorktree and `worktree.*` facts, the cleanup facts from earlier sessions.
- `documentation/deployment.md`: preview deployments are per PR, so the branch name changes nothing there.
