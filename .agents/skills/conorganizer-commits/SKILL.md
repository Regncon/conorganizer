---
name: conorganizer-commits
description: Conorganizer's rules for commit messages, with a checker script: the subject and body shape, where the issue reference goes, attribution lines, and how to rewrite history without an editor. Use whenever you are about to run git commit, amend, squash, split, reword or rebase commits, write or rewrite a commit message, or the user says "commit this", "commit these changes", "clean up the commits", "fix the commit message", "squash" or "amend"; also when reviewing the commit messages on a branch or in a PR before it is opened or merged, and when writing a PR title, since GitHub turns it into the merge commit subject. Use it even when the commit is a small step inside a larger task, since every commit lands in main's history and the checker runs on each one.
---

# Conorganizer Commits

A commit message is read later by someone with the diff open and no memory of the task: in `git log`, in `git blame`, in a bisect, in the PR's merge commit. The subject tells them what changed for a user or admin; the body tells them why it was done this way. AGENTS.md has the one-line rule under "Git"; this is the depth, and `scripts/check-commit-msg.sh` checks the mechanical part of it.

## The shape of a message

```
Show pulje status alerts in a bar at the top of the front page

The "lukkes snart" and "ikke publisert" warnings sat on every event
card, so a billettholder scrolling the program saw the same warning
repeated down the page and missed it when no card was in view. One
bar at the top says it once, for every pulje that needs it.

Closes #580
```

1. One subject line: capitalized imperative, the outcome, 72 characters or fewer, no period.
2. One blank line.
3. The body, when there is one: why, in one to three short paragraphs wrapped at 72.
4. `Closes #NNN` on its own line, last, when the commit closes an issue.

## Subject

- Imperative, as if completing "This commit will ...": Add, Fix, Show, Let, Stop, Keep, Make, Move. "Updated" or "Fixes" reads as a changelog entry, and the history reads as instructions.
- First word capitalized, no trailing period: the subject is a headline in `git log`, on GitHub and as the PR title, and that is how the history already reads.
- Say what a user or admin gets, not the mechanism: "Start the Grafana MCP server in about a second instead of timing out", not "Install mcp-grafana with go install instead of go run". The mechanism is in the diff; the outcome is what someone searching the log needs.
- 72 characters at most, 50 to 65 is the comfortable range: `git log --oneline` and GitHub cut longer subjects, and a PR title gets " (#NNN)" appended when it is merged.
- Norwegian domain words are ordinary words: lowercase inside the subject (billettholder, pulje, puljefordeling), capitalized only as the first word. Product names keep their casing (Regncon, CheckIn, Datastar, Grafana, NATS).
- No "NNN | " prefix: the issue goes in the body, so the subject spends its characters on the change.
- A subject that needs "and also" describes two commits. Split them when the code is yours to rearrange (`references/rewriting-history.md`); when the commit has to stay as it is, name the main change in the subject and the others in the body.

## Body

Write a body when the why is not in the diff: a symptom the diff does not show, a decision between alternatives, a workaround for something outside the repo, a reversal of an earlier decision, a fact the next reader would otherwise get wrong. Leave it out when the subject says it all: a typo fix, a rename, a doc link, a version bump with no story.

- One blank line after the subject: git treats the first paragraph as the subject, so without the blank line the body becomes part of it.
- Why, not what: what was wrong before, why this fix and not another, what changes for the people who use it. A paragraph that walks through the diff ("now loops three times with a 500 ms sleep", "also fixes the README typo") is a code comment in the wrong place. The test for each sentence: would it still tell the reader something with the diff open?
- Take the reason from the person who asked for the commit when they gave one. When you inferred the why from a test name or an error string, say so in your reply, outside the message, so the user can correct it before it is pushed.
- Wrap at 72: `git log` indents the body by four spaces, and an 80-column terminal cuts the rest.
- No attribution lines (Co-Authored-By, "Generated with", a robot emoji): the history names people only in git's author field.
- English, with the Norwegian domain words kept, like the subject.

## Issue references

- `Closes #NNN` on its own line, last in the body, when the commit closes the issue. One keyword, one place: GitHub closes the issue when the commit reaches main, and a reader finds it without scanning the text. "Fixes" and "Resolves" also work on GitHub, but the history uses "Closes".
- A plain `#NNN` inside a sentence when the commit only relates to an issue: "Follow-up to #699", "reverses the decision in #524".
- Never in the subject, and no commit on main carries a "NNN | " prefix.

## Before you commit

1. Stage only the files of this change. `git add -A` in a tree with other work in it commits that work too.
2. Write the message to a file (`tmp/` is gitignored) and run the checker on it, then commit with `-F`:

```bash
bash ".agents/skills/conorganizer-commits/scripts/check-commit-msg.sh" tmp/commit-msg.txt
git commit -F tmp/commit-msg.txt
```

The repo path contains a space, so the script path is quoted. A file keeps the wrapping; `git commit -m subject -m body` works for a one-sentence body. The checker catches the mechanical rules (casing, period, length, prefix, blank line, attribution) and prints each problem; the why is yours to judge.

3. To have git refuse a bad message on every commit, install the hook once per clone: `bash ".agents/skills/conorganizer-commits/scripts/install-commit-msg-hook.sh"`. It writes `.git/hooks/commit-msg`, which calls the checker inside the checkout, so a later change to the script applies without reinstalling. Nothing is committed; each developer decides for their own clone.

## Rewriting commits

Amending, rewording, squashing and splitting change commit hashes, so do it only on commits that are not on main, and on a pushed branch only when it is the user's own; tell them the push will need `--force-with-lease`. Interactive rebase needs an editor, which is not available here, so use the recipes in `references/rewriting-history.md`: amend the last commit with `-F`, reword older commits with `exec` lines in a rebase todo, squash with `git reset --soft`, split with `git reset` and selective `git add`. After any rewrite, `git diff <old-tip> HEAD` is empty when only messages changed, and the old tip stays in the reflog. Merging main into the branch is still the house way to update it; rewording your own unpushed commits does not conflict with that rule.

## Reviewing messages on a branch or PR

Run every message through the checker, then read the bodies.

- Report what breaks a rule, each with the rewritten message, most important first: a subject that says nothing or names the mechanism, a body that restates the diff, an attribution line, an issue reference in the wrong place.
- A body that states something the diff contradicts is a defect: say which claim and what the diff shows.
- A message that follows the rules is left alone, even when you would word it differently. A review of messages is not a copy edit; naming a passing message as "mostly good" with a reworded version costs the user a decision for nothing.
- Propose a split when a subject needs "and also", and say the code stays the same.
- Do not rewrite unless asked; the user may want to do it in their own order.

## PR titles

A PR title follows the subject rules, because GitHub makes it the merge commit's subject and appends " (#NNN)". Keep it at 65 characters or fewer, so the seven-character suffix fits in 72. The description has its own rules in the `conorganizer-pr` skill.

## Checklist

1. Subject: capitalized imperative, outcome not mechanism, 72 or fewer, no period, no "NNN | " prefix, domain words lowercase mid-subject.
2. Blank line after the subject.
3. Body only when the why is not in the diff, and then why, not what; wrapped at 72; no attribution lines.
4. `Closes #NNN` last when the commit closes an issue; `#NNN` in a sentence otherwise.
5. `bash ".agents/skills/conorganizer-commits/scripts/check-commit-msg.sh"` passes.
6. Only this change's files are staged.

## References

- `references/examples.md`: subjects and bodies from the real history, before and after, with the reason for each change. Read it when a subject or body does not come out right on the first try, or when reviewing.
- `references/rewriting-history.md`: amend, reword, squash and split without an editor, with the checks that the code is unchanged and the force-push etiquette. Read it only when you change existing commits.
