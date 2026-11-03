---
name: conorganizer-pr
description: Conorganizer's rules for pull request titles and descriptions. Use whenever you open or edit a pull request, run gh pr create or gh pr edit, write or rewrite a PR title or body, or write the commit subject of a one-commit branch, since a squash merge puts that subject on main as the title. Also use when the user says open a PR, make a PR, write this up, shorten or fix the description, or that a PR text is too long, vague or off, and when you reach the end of a branch and a PR is the next step, even if nobody mentioned the PR text. Every PR here needs a title of at most 65 characters and a body under 200 words, and without this skill both come out two to five times too long.
---

# Conorganizer PR

The repository merges with the PR title as the commit subject and the PR body as the commit message, so what you write here is what `git log` on main shows for good. People read the title in a list and the body once, at review time or months later when something breaks. A title that names the outcome and a body a reader finishes in a minute serve both. AGENTS.md "Git" has the base rules (imperative English, the user-visible outcome, Norwegian domain words, no Co-Authored-By, merge main rather than rebase, leave the CI preview block alone); this skill is the shape of the title and body on top of them. `references/examples.md` holds the merged PRs the maintainers picked as the model, and three to avoid.

## Before you write

- Read the whole branch, not the last commit: `git log --format='%s%n%b' main..HEAD`, `git diff main...HEAD --stat`, and the diff itself for anything you did not write in this session. Describe the diff against main: a change that a later commit on the branch undid never happened as far as the reader is concerned.
- Say in one sentence what a user or admin sees differently after the merge. That sentence is the title and the opening of the body. When it will not fit in one sentence, the branch does two things: tell the user and offer to split before writing a list title.
- Find the hand steps: a new file in `migrations/` (CI migrates only PR previews; main and demo are migrated by hand), a change under `configuration-as-code/`, `hetzner/` or `deploy/` that the server picks up only after a checkout move, restow or reload, a setup task everyone runs again, a dashboard to re-import. Each becomes a line in the rollout section, and nothing else does.

## Title

One line, at most 65 characters, so the merge commit subject with " (#NNN)" is at most 72. Count it: `printf '%s' "$title" | wc -m`.

- It follows the commit subject rule, because a merge turns it into the commit subject plus " (#NNN)": capitalized first word, imperative, verb first (Add, Let, Show, Fix, Stop, Keep, Require, Remove), no trailing period.
- It names the outcome someone sees, not the mechanism: "Let admins print a QR code to the feedback form", not "Add qrcode package and admin route".
- Norwegian domain words are ordinary words, lowercase mid-title (pulje, billettholdere, puljefordeling); page and product names keep their casing (Puljeoppsett, Adminverktøy, Min side, Datastar).
- No issue number, branch slug, "Fix:" or "feat:" prefix, Title Case, trailing period or emoji. GitHub adds "(#NNN)" at merge, and the issue goes in the body.
- Not a list. "Fix templ in place: v0.3.1070, Prettier formatting, fmt check in CI, once-rendered looped styles" is four outcomes; name the one the reader cares about most and let the bullets carry the rest, or split the PR.
- A one-commit PR that is squash-merged gets its commit subject on main, not the PR title, so that subject must say the same thing as the title.

| Before | After |
| --- | --- |
| 583 bug kartene som vises i kart modalene er i light mode | Show the room maps in dark mode in the map dialogs |
| move kontoistillinger down and some small refactoring | Move Kontoadministrasjon to the bottom of Min side |
| 615 migration i ci | Let CI run migrations for PR previews |
| Fix: Names with spaces | Put a space between first and last name on the ticket card |
| Last-Minute UI Polish and Cleanup | One PR per fix, or the biggest one named: Colour event status tags after the annonsert rename |

## Body

The body is the shortest text that tells a reviewer what changed and an admin what to run. Most PRs fit in 60 to 120 words; 200 is the ceiling, counted above the rollout section and the preview block, which are exempt. #653 and #611 in `references/examples.md` are the size to aim for.

It has these parts, in this order, and nothing else:

1. `Closes #NNN` alone on the first line when the PR closes an issue. A PR that only relates to one says so in the summary: "Follow-up to #699", "Stacked on #605".
2. The summary, one to three sentences, no heading: what is different after the merge and, when it is not obvious, why: the bug as someone saw it and its cause, in a sentence. For a small PR this is most of the body.
3. `## What changed`, only when the branch does several things in parallel: one bullet per change, one sentence each, naming where and what ("Romfordeling tab: a Publiser romfordeling switch per pulje"), at most six. A single change stays in the summary. Seven or more means the bullets are too fine or the branch too wide. No nested bullets, no bold lead-ins.
4. `## Tests`, one to three lines: the test names or behaviours that prove the change, what was checked by hand, and what was not checked. Only what you ran or read in the branch's tests: when you did not run something, say so in a few words ("not run on Windows"), because a manual check inferred from a commit message is a claim nobody made. "The full suite and lint pass" is what CI shows, so leave it out, and never a checklist of things still to do.
5. `## After merging`, only when someone has to act by hand; see Rollout steps.
6. The preview block, which CI adds; see below.

That is all. A known limit or follow-up that changes what the reviewer does gets one sentence at the end of the summary or one bullet ("Not in this PR: the image upload page GET still has no owner check."); the rest belongs in the docs, an issue or the review thread. The body carries no account of how the work was done (agents, review rounds, hours), no praise of the change, no emoji, no attribution line, and no walk through the diff file by file: the diff is one click away, the body says what it adds up to. Name a file only when the reader has to open it, such as the migration or the doc with the procedure; the reviewer sees the file list in the diff. Headings other than the three above are a sign a part has outgrown its sentence.

Language: English, with user-facing bokmål strings quoted as they appear ("Lagre fordeling", "Publiser romfordeling"). Screenshots are not required; the preview link CI adds covers UI changes.

## Rollout steps

`## After merging` (or `## On the server, after merging`) exists only when the merge alone does not finish the change:

- A migration: name the file and point at the procedure, in one line: "Run `20260926120000_feedback.sql` on main and demo by hand (documentation/migrations.md)."
- Server config under `configuration-as-code/`, `hetzner/` or `deploy/`: the commands as one fenced block, or as numbered steps when a step needs a word of context, in order, with the one line that shows they worked. Take them from the README section the PR updates instead of inventing a procedure: the checkout moves with `git fetch origin` and `git checkout --detach origin/main`, Caddy changes need `caddy validate` before `systemctl reload caddy`, stowed files need `stow --restow <package>`. Never name a server account, because the repository is public; when the README you copy from still names one, write that step as plain `git` from your own login shell, as "Repository location" in `configuration-as-code/README.md` does.
- Something everyone runs once, such as `go tool task setup:mcp`, or a dashboard to re-import.
- `## Before merging` only when a check must pass first, such as a sudoers line or a service already in place.

The rollout section does not count toward the 200 words, so a procedure that needs eight steps and a warning about their order gets them; the body above it still does. Nothing to run means no section: "No rollout steps needed" costs the reader a sentence for nothing.

## The preview block

After the first preview deploy, CI appends this to the body and skips it on later deploys when it is already there:

```
<!-- preview-deployment-url:start -->
## Preview deployment

_Added automatically by CI_

🔗 https://NNN-merge.lekeplassen.regncon.no
<!-- preview-deployment-url:end -->
```

`gh pr edit --body-file` replaces the whole body, so fetch the current one, change the text above the block, and keep the block byte for byte at the end:

```bash
gh pr view NNN --json body -q .body > body.md
```

## Writing it with gh

- Put the body in a file outside the repository and pass `--body-file`; backticks and `$` inside `--body "..."` are eaten by the shell, on Git Bash especially.
- Create: `gh pr create --base main --title "$title" --body-file body.md`. Add `--draft` when a server rollout has to be staged before the merge.
- Edit: `gh pr edit NNN --title "$title" --body-file body.md`, after the fetch above.
- Then read it back once with `gh pr view NNN --json title,body -q '.title, .body'`; what you see is what main's history will carry.

## Check before you run gh

1. Title at most 65 characters, capitalized verb first, the outcome, no number, slug or prefix.
2. Body under 200 words above the rollout section and the preview block, and aimed at about 100: `wc -w body.md`.
3. Only `## What changed`, `## Tests` and `## After merging` as headings, each present only when it has something to say.
4. `Closes #NNN` exactly when the PR closes an issue.
5. A rollout section exactly when a hand step exists, with migrations named by file.
6. The tests lines say what was checked and what was not.
7. Nothing about process, no emoji, no attribution, no praise, no file-by-file tour.
8. On an edit, the preview block is unchanged.

## References

- `references/examples.md`: the three merged PRs the maintainers picked as the model, three to avoid with what is wrong in each, and #671's body rewritten to this shape.
