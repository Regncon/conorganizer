# Rewriting history without an editor

Every recipe here changes commit hashes. Use them on commits that are not on main, and on a pushed branch only when it belongs to the user; say in your reply that the branch needs `git push --force-with-lease`. The old tip stays in `git reflog` for a while, so a mistake is recoverable: `git reset --hard <old-tip>`.

All commands run from the branch's checkout. If the branch is checked out in another worktree, git refuses to rebase or reset it from here; work in that worktree. The repo path contains a space, so quote every path that includes it.

Write each new message to a file first (`tmp/` is gitignored) and run the checker on it, so the rewrite never introduces a message the hook would refuse.

## Amend the last commit's message

```bash
git commit --amend -F tmp/commit-msg.txt
```

To also add forgotten files, stage them first; `--amend` takes the index as the new tree.

## Reword older commits

`git rebase -i` opens an editor for the todo list and again for every `reword`. Neither is available here, so write the todo list yourself and hand it to the rebase with `cp` as the sequence editor. After each `pick` that needs a new message, an `exec` line amends it with the prepared file:

```bash
git log --reverse --format='pick %h %s' HEAD~3..HEAD > tmp/todo.txt
```

Edit `tmp/todo.txt` so it reads, with absolute paths to the message files:

```
pick f6b9fab fix the thing.
exec git commit --amend -q -F "/full/path/to/tmp/msg-f6b9fab.txt"
pick 4849e19 651 | Updated the pulje closing warning text and also ...
exec git commit --amend -q -F "/full/path/to/tmp/msg-4849e19.txt"
pick ebb9e19 Add retry to the pulje status broadcast
```

Then run it:

```bash
GIT_SEQUENCE_EDITOR="cp tmp/todo.txt" git rebase -i HEAD~3 --quiet
```

The `exec` line runs right after its commit is replayed, so the amend hits the right one. A `pick` without an `exec` line keeps its message. The paths are absolute because `exec` runs from the repo root whatever the current directory was. The author name and date stay; the committer date becomes now, as with any rebase. Writing the todo as a file avoids quoting a `sed` script with `\n` in it, which some shells and tool harnesses mangle.

The same can be done in one line with `sed` as the sequence editor when the shell passes backslashes through unchanged (`-c` is an option of `git` itself, so it goes before `rebase`):

```bash
git -c sequence.editor="sed -i \
    -e 's|^pick f6b9fab .*|&\nexec git commit --amend -q -F \"$PWD/tmp/msg-f6b9fab.txt\"|'" \
  rebase -i HEAD~3 --quiet
```

An alternative that needs no rebase: `git commit-tree` builds a new commit from an existing tree, which is how to reword a commit whose parent chain you then rebuild by hand. The rebase recipe is shorter for one branch.

## Squash several commits into one

```bash
git reset --soft HEAD~3
git commit -F tmp/commit-msg.txt
```

`--soft` keeps the index and working tree, so the new commit holds exactly the sum of the three. Write the message fresh; the old subjects were steps, the new one is the outcome.

## Split one commit into several

```bash
git reset HEAD~1
git add service/puljer/closing.go service/puljer/closing_test.go
git commit -F tmp/msg-closing.txt
git add static/css/index.css
git commit -F tmp/msg-banner.txt
git add README.md
git commit -F tmp/msg-readme.txt
```

A mixed reset leaves the changes unstaged, so each `git add` picks a commit's files. `git add -p` splits a file when two changes share it. To split a commit further back, reword-style: rebase onto its parent with `exec git reset HEAD~1` after its pick, then the adds and commits, then `git rebase --continue` when the sequence stops.

## Check that only the messages changed

```bash
git diff <old-tip> HEAD --stat
```

Empty output means the tree at the tip is the same. For a reword, `git log --format='%h %T %s' <old-tip>` and the same on `HEAD` show the tree hashes; they match commit for commit. For a squash or split, the tip tree is what matters.

Then run the checker over the branch in one go:

```bash
for sha in $(git log --format=%h main..HEAD); do
  git log -1 --format=%B "$sha" | bash ".agents/skills/conorganizer-commits/scripts/check-commit-msg.sh" || echo "$sha"
done
```

A hash printed after a problem line is the commit to fix.

## Pushing afterwards

```bash
git push --force-with-lease
```

`--force-with-lease` refuses when someone else pushed to the branch in the meantime, which a plain `--force` would overwrite. Never force-push main. A branch that is already in an open PR keeps its PR; GitHub shows the rewritten commits and drops the old ones from the PR view.
