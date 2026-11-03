# What the desktop app and Claude Code do with worktree names

Observed on 2026-11-03 in the session that wrote this skill. It was started from another session with `start_session` and `branch_name: "worktree-skill"`, so it shows what a requested name becomes. Tool names are the Claude desktop app's (`ccd_settings`, `ccd_session_mgmt`, `ccd_pr`). Re-check the facts marked *unverified* before relying on them; the rest were run and read.

## 1. Branch names the desktop app gives a session

| How the session started | Branch the app created | Worktree folder |
| --- | --- | --- |
| `start_session` with `branch_name: "worktree-skill"` | `claude/worktree-skill-2cec9f` | `.claude/worktrees/adoring-saha-f65bac` |
| `start_session` with `branch_name: "commit-skill"` | `claude/commit-skill-803af8` | `.claude/worktrees/blissful-herschel-82377c` |
| `start_session` without a name (sibling "pr-skill") | `claude/keen-kowalevski-3ae2ec` | `.claude/worktrees/sharp-shtern-55b2d3` |
| Earlier sessions (September and October) | `claude/lucid-kirch-b765f7`, `claude/goofy-lalande-d17893`, `claude/hopeful-shamir-1a2903` | the same word pair as the branch, or no worktree at all |

So the pattern is `<prefix>/<requested name or adjective-surname>-<6 hex>`. The hex suffix is always added, also to a requested name. A session that runs in the main checkout instead of a worktree gets a `claude/…` branch too ("Fix templ in place" ran in the main checkout on `claude/infallible-turing-611030`).

The prefix is the app setting `branch_prefix` (Settings → Code → Branch prefix), read with `ccd_settings get_settings`:

- value `"claude"`, meaning "The prefix added to the branch names of new sessions, local and cloud (default "claude", giving branches like claude/<name>). An account setting."
- allowed: 1-40 letters, digits, `.`, `_` or `-`, starting with a letter or digit; not main, master, staging, warm, prod-*, release*.
- changing it is a `set_setting` call the user approves on a card. It was read, not changed, for this skill.

The prefix is the only naming control the app offers. Nothing sets the rest of the name or drops the suffix.

## 2. Worktree folders are the app's, and recycled

- Location: the read-only setting "Worktree location" is "Inside the project (<repo>/.claude/worktrees)". `/.claude/worktrees/` is in `.gitignore` (line 27).
- The folder name is generated and says nothing about the task. It is also recycled: this session's folder `adoring-saha-f65bac` was the worktree of an archived session from 2026-10-05 (its branch `claude/adoring-saha-f65bac` still exists, and the folder's reflog still holds that session's commits). `blissful-herschel-82377c` served an archived session with PR #671 before the commit-skill session, and `sharp-shtern-55b2d3` served one with PR #696 before the pr-skill session.
- The app registers a worktree by path in `git-worktrees.json` (`gitRoot`, `commonDir`, `gitDir`), not by branch, and the session file holds `worktreePath` and `worktreeName`. The folder is the session's cwd. Do not rename or move it.
- The worktree's git dir carries `config.worktree` with `core.longpaths = true` and an empty marker file `claude-desktop-worktree`.
- The app also links `.claude/skills` and copies `.claude/settings.local.json` (the permission allowlist, ignored through the user's global git ignore) into its worktrees at session start. A worktree made with `git worktree add` or `EnterWorktree` gets neither.

## 3. Renaming an app-generated branch (tested before any commit)

```bash
git branch -m claude/worktree-skill-2cec9f claude/worktree-skill
```

Afterwards, in git: `git branch --show-current` prints the new name, the worktree's `HEAD` file points at `refs/heads/claude/worktree-skill`, `git worktree list` shows the folder on the new branch, `git status` is clean and the old ref is gone. The app had set no upstream (it sets none until the first push), so nothing had to be re-pointed.

In the app:

- `ccd_pr get_status` was "no pull request bound" before and after, so it says nothing about the rename; the app binds a PR when the session runs `gh pr create` or pushes a branch with an open PR. Section 7 records what happened when this skill's PR was opened from the renamed branch.
- `ccd_session_mgmt get_session self` still reports `"branch": "claude/worktree-skill-2cec9f"`. The app writes the branch name into the session file (`claude-code-sessions/…/local_<id>.json`, keys `branch` and `sourceBranch`) when it creates the session and does not read it back from git. So the sidebar label and anything the app does with the stored name keep the old one. The description of `clean_up_worktrees` says a resumed session "checks its branch out again"; what that does when the stored branch no longer exists is *unverified*.

Conclusion: the rename is safe for git and for the work. Commits, the push and the PR use the new name. The one cost is the stale name in the app's session record.

## 4. Claude Code's own worktrees

These are the CLI's, not the desktop app's: `claude --worktree <name>`, the `EnterWorktree` tool, subagents with `isolation: "worktree"` and Workflow runs.

- Directory `.claude/worktrees/<name>`, branch `worktree-<name>`. Observed: a subagent with worktree isolation got `.claude/worktrees/agent-a6d2ff1162bfbd8e8` on `worktree-agent-a6d2ff1162bfbd8e8` (locked), and a Workflow run `wf_7b8f7f7e-af1-5` sits on `worktree-wf_7b8f7f7e-af1-5`. The docs (code.claude.com/docs/en/worktrees) give the same pattern for `--worktree`.
- `EnterWorktree` takes `name` (each "/"-separated segment letters, digits, dots, underscores and dashes, at most 64 characters) or `path` to enter an existing worktree listed by `git worktree list`. It cannot create a worktree from inside a worktree session, and the docs do not say which branch name the `name` form produces beyond the `worktree-<name>` default.
- No setting names the branch or its prefix. The `worktree.*` settings are `baseRef` (`"fresh"`, the default, branches from `origin/<default branch>`; `"head"` from the local HEAD), `bgIsolation`, `location` (SSH only), `sparsePaths` and `symlinkDirectories`.
- Reopening a worktree by name checks that it is still on the branch Claude Code created for it; otherwise it reopens at the old tip. `ExitWorktree` with `remove` deletes the directory and its branch and refuses when there are uncommitted files or unmerged commits unless `discard_changes` is true. Renaming the branch is documented neither as allowed nor as a problem.
- The `superpowers:using-git-worktrees` skill prefers a native tool and otherwise runs `git worktree add "$LOCATION/$BRANCH_NAME" -b "$BRANCH_NAME"`, so its directory name is the branch name and its default location is `.worktrees/`, not this repo's `.claude/worktrees/`. It takes the branch name from the caller and does not say how to choose it. Its baseline step runs the whole test suite, which in this repo needs `go tool templ generate` first.

## 5. `git worktree add` in this repository

No tracked file runs it. README, Taskfile.yml and documentation/ never mention worktrees; AGENTS.md "Dev setup" has the three worktree rules (no watcher, no `task test`, other sessions' checkouts live under `.claude/worktrees/`). The user's notes describe `git worktree add --detach <scratch>/main-baseline main` for a baseline server on another port.

Branch names in the merged history are mixed (`task/…`, `bug/651-…`, `feature/…`, `NNN-slug`, `claude/…`). CI runs on every branch (`branches: ["**"]`) and reads the branch name only to detect `main`, so a `claude/<name>` branch changes nothing in CI or the preview deployment.

## 6. Cleanup facts from earlier sessions

- On Windows, `git worktree remove` cannot delete the folder of the session you are running in, because it is the process's working directory: git unregisters the worktree and empties the folder, and the empty folder goes once that session is closed (2026-11-01). Seen again on 2026-11-03: `.claude/worktrees/lucid-kirch-b765f7` is an empty folder that `git worktree list` no longer shows, left by the session whose PR #699 merged the day before.
- Sending a new message in a chat whose worktree was removed makes the app check the branch out again.
- `git worktree remove` on a worktree that has the `.claude/skills` junction exits 0 and still leaves `<folder>/.claude/skills` (now dangling) and the two folders behind. Reproduced twice on 2026-11-03 with `.claude/worktrees/612-level-colors`; the shell blocked the next `git worktree add` of that name. `rm` on the link and `rmdir` on the folders clear it. `git worktree add -b <branch> <path> origin/main` also sets `origin/main` as the branch's upstream unless `--no-track` is given.
- Archiving a session is not a way to remove its worktree. In auto mode the app archived without a card, and the user keeps finished chats as reference (2026-11-01).

## 7. PR binding after the rename

Recorded when this skill's PR was opened from `claude/worktree-skill`.
