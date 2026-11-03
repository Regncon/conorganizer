# Examples from the history

The "before" messages are real commits on main. The "after" messages were written from their diffs and PR bodies to show the rules; where the reason behind a change is not on record, the example says so instead of inventing one.

## Subjects

### Casing and period

Before, `3fd8f64c`:

```
fix meta panel colors
```

After:

```
Give the puljefordeling meta panels their Level backgrounds
```

Capitalized, and the subject names the page and the rule the diff applies (`--level-2-background` becomes `--level-3-background` and so on). "fix ... colors" could be any of a hundred commits in this history.

Before, `7208349b`:

```
fix background on dropdown
```

After:

```
Give the header menu dropdown a Level 3 background
```

### Outcome, not mechanism

A mechanism subject for `5c3b41ac` would have been:

```
Install mcp-grafana with go install instead of running it with go run
```

As committed:

```
Start the Grafana MCP server in about a second instead of timing out
```

The install command is in the diff; what a developer gets is an MCP server that starts before Claude Code gives up on it.

### Length

Before, `06abad27`, 89 characters:

```
Keep a SQLite backup successful when the WAL or a backup file vanishes while it is measured
```

After, 61 characters, with the detail in the body:

```
Keep a SQLite backup successful when a measured file vanishes

du ignores files that disappear during its walk, which happens when
the image backup deletes an old archive at 03:30 while the SQLite
backup is running, and a WAL file can vanish just before it is
measured. Both made the size measurement fail a backup that had
already succeeded; now a vanished file counts as 0.
```

### "and also"

Before, `2dd0bf88`, 82 characters:

```
Show backup age from the last successful run and add a backup directory size panel
```

After, the main change in the subject, the second one in a sentence of the body:

```
Show backup age from the last successful run
```

When the two changes can be separate commits, split them instead (`references/rewriting-history.md`).

Before, `32cafaaf`:

```
move kontoistillinger down and some small refactoring
```

After:

```
Move kontoinnstillinger below the ticket cards on the profile page
```

"some small refactoring" has no place in a subject. When the refactoring is worth a word, it goes in the body (here: `LoggedInAccount` now renders its own pane, so `profile_index.templ` no longer wraps it); when it is not, it goes nowhere.

### Issue prefix

Before, `f379c2c9`, with the PR body "Fix some bugs, update times. Fixes #651":

```
651 | Update contents of time schedule
```

After:

```
Update the times and day headers in the time schedule

Closes #651
```

The issue reference moves to its own last line, and the subject says which contents changed.

### Conventional-commit prefix and attribution

Before, `805f428e`, with the two changes of the PR joined by a semicolon and a session link at the end of the body:

```
fix: never auto-seat participants under 18 in 18+ games; confirm manual placements

...
Claude-Session: https://claude.ai/code/session_...
```

After, one commit per change, no prefix, no session line:

```
Never seat participants under 18 in 18+ games automatically

The solver had no notion of age: any interesse became an assignment
edge, so a minor who rated an 18+ game highly could be seated there.
Manual pins stay as the admin's deliberate override.
```

```
Confirm before manually placing a minor in an 18+ game
```

### Title Case and a grab-bag

Before, `c47c4900`, the merge commit of a PR with twelve commits:

```
Last-Minute UI Polish and Cleanup
```

After, the main change of the PR as its title, so the merge commit says what landed:

```
Link program puljer to their day and tidy the program cards
```

Title Case is for headlines in print; a subject is a sentence. A PR title that names no change leaves a merge commit that names no change.

### Already fine

`c3b95d55`:

```
Add banner to show during festival
```

Articles would make it read as a sentence ("Add a banner to show during the festival"), but it is capitalized, imperative, short and says what a visitor gets. A review leaves it alone.

### Merge commits

Before, `7c0180ea`:

```
merge main in and solve the merge boy
```

After: git's own subject, which the checker skips:

```
Merge branch 'main' into 612-impliment-new-color-scheme
```

## Bodies

### Why, with the evidence

`01176b95`, as committed:

```
Reset a preview unit's failed state only when it is failed

A cleanly stopped unit is not loaded after the unit file is removed, so
reset-failed printed a 'not loaded' error on every normal cleanup. Seen
on the first real run against 676-merge.
```

Two sentences: what went wrong and where it was seen. Nothing about the `if` the diff adds.

### A reversal

`8027909e`, trimmed to the decision and the reference to what it reverses:

```
Let the solver seat participants under 18 in 18+ games

The board decided that age limits are handled by hand, so the solver
no longer knows about age. The confirmation before manually placing a
minor in an 18+ game stays, and is now the only guard. This reverses
#524.
```

The reader of `git blame` on the solver finds the decision and the earlier commit without searching the issues.

### What turned into why

Before, the body of `2dd0bf88`:

```
Production Health's SQLite and image backup age panels now read the
conorganizer_*_last_success_timestamp_seconds metrics instead of the
systemd timer, and a Backup directory size panel joins them. Panel
descriptions, the dashboard notes and the README no longer call these
metrics optional.
```

After, with the reason from the PR (#671):

```
The age panels from #670 measured when the systemd timer last fired,
not when a backup last succeeded, so a failing backup looked fresh.
The panels now read the scripts' own last-success metrics. A Backup
directory size panel joins them, since the same scripts publish it.
```

Which panels read which metric is in the diff. What the reader needs is why the timer was the wrong source.

### A constraint outside the repo

`5c3b41ac`, trimmed:

```
Start the Grafana MCP server in about a second instead of timing out

The agent configs ran mcp-grafana through go run, which spends about
15 seconds checking its build on every start and sometimes ran past
Claude Code's 30-second MCP startup limit. The server is now installed
once by go tool task setup:mcp and started from PATH.
```

The limit belongs to another tool and is nowhere in the diff, so the body is the only place a reader learns why `go run` had to go.

### No body

```
Remove the duplicate --color-text-soft on the print page
```

```
Unquote the generic monospace fallback in --font-monospace
```

The subject is the whole story; a body would repeat it.
