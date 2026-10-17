# Rename translated domain words to the Norwegian domain words

Status: plan, not started. Written 2026-10-18 against main at 8033ed17.

The code will have changed by the time this plan runs. Do [Before you start: rescan](#before-you-start-rescan) first, so the plan covers every name in the code as it is then.

## Why

On 2026-10-18 the naming rule was decided (AGENTS.md "Names", added in #663). Code is English, but words from [domeneordbok.md](../domeneordbok.md) stay Norwegian, written without æøå (`getBillettholder`, `Forstevalg`). Some existing names translate a domain word. This plan renames them:

| Translated word | Domain word |
| --- | --- |
| event / events | arrangement / arrangementer |
| interest / interests | interesse / interesser |
| room / rooms | rom (same in plural) |
| ticket holder (`components/ticket_holder`, package `ticketholder`) | billettholder / billettholdere |

Casing follows the existing name: `EventID` → `ArrangementID`, `eventID` → `arrangementID`, `event_id` → `arrangement_id`, `.event-room-dialog` → `.arrangement-rom-dialog`.

## Decisions (answered 2026-10-18)

| Question | Answer |
| --- | --- |
| Rename database tables, columns and views? | Yes, as the last phase. |
| Change URLs? | No. Every URL stays: `/event/{id}`, `/event/api/...`, `/admin/rooms/...`, `/profile/api/new/{id}/event-in-pulje/...`, `/profile/api/new/{id}/assign-room/...`, `/admin/api/puljefordeling/interest`, `/event-images/*`, `/static/...`. Only code names change. Route parameter names like `{event}` and `{room}` are code names and can change, because they are not part of the URL. |
| Is the `ticketholder` package in scope? | Yes. |
| CSS classes, HTML ids, Datastar signals, JSON tags? | Yes. |
| Fix the plural `billettholders` → `billettholdere`? | Yes. |
| Rename log field keys and component names? | Yes. |
| Compound words, like `InterestLevel` (Interessenivå)? | Swap only the translated word: `InterestLevel` → `InteresseLevel`, `interest_level` → `interesse_level`, `RoomAssignment` → `RomAssignment`, `rooms_published` → `rom_published`. |

## Name mapping

### Go packages

Directory and package names keep their current style, so each change is a pure rename (see open item 3).

| Today (package) | After (package) | Files | Importers |
| --- | --- | --- | --- |
| `service/eventService` (`eventservice`) | `service/arrangementService` (`arrangementservice`) | 7 | 8 |
| `service/eventimage` | `service/arrangementimage` | 1 | 10 |
| `service/rooms` (`rooms`) | `service/rom` (`rom`) | 12 | 8 |
| `pages/event` (`event`) | `pages/arrangement` (`arrangement`) | 22 | 4 |
| `pages/admin/rooms` (`rooms`) | `pages/admin/rom` (`rom`) | 5 | 1 |
| `pages/profile/newevent` | `pages/profile/newarrangement` | 3 | 1 |
| `components/event_components` | `components/arrangement_components` | 8 | 7 |
| `components/formsubmission/event_img_upload` (`eventimgupload`) | `.../arrangement_img_upload` (`arrangementimgupload`) | 2 | 1 |
| `components/ticket_holder` (`ticketholder`) | `components/billettholder` (`billettholder`) | 11 | 8 |

`service/billettholder` is package `billettholderService`, so a new `billettholder` package does not clash with it.

### Database

| Today | After |
| --- | --- |
| `events` | `arrangementer` |
| `interests` | `interesser` |
| `rooms` | `rom` |
| `event_statuses` | `arrangement_statuses` |
| `events_types` | `arrangementer_types` (see open item 4) |
| `event_runtimes` | `arrangement_runtimes` |
| `interest_levels` | `interesse_levels` |
| `relation_event_puljer` | `relation_arrangement_puljer` |
| `relation_events_players` | `relation_arrangementer_players` |
| column `event_id` (in `relation_event_puljer`, `interests`, `relation_events_players`) | `arrangement_id` |
| column `relation_event_puljer.room_id` | `rom_id` |
| column `interest_level` (in `interests`, `interest_levels`) | `interesse_level` |
| column `event_type` (in `events`, `events_types`) | `arrangement_type` |
| column `events.event_runtime` | `arrangement_runtime` |
| column `rooms.room_number` | `rom_number` |
| column `puljer.rooms_published` | `rom_published` |
| view `v_event_summary` | `v_arrangement_summary` |
| view `v_events_by_pulje_active` | `v_arrangementer_by_pulje_active` |
| view `v_event_puljer_active` | `v_arrangement_puljer_active` |
| view aliases `room_id`, `room_number`, `room_name`, `room_floor`, `room_max_concurrent_games`, `room_notes`, `room_is_disabled` | `rom_*` |

That is 9 tables, 11 column renames and 3 views.

### Live buckets

| Today | After |
| --- | --- |
| `BucketEvents` = `"events"` | `BucketArrangementer` = `"arrangementer"` |
| `BucketInterests` = `"interests"` | `BucketInteresser` = `"interesser"` |
| `BucketRooms` = `"rooms"` | `BucketRom` = `"rom"` |
| `BucketBillettholders` = `"billettholders"` | `BucketBillettholdere` = `"billettholdere"` |

## Not in scope

- Every URL, including static file paths: the 15 maps in `/static/rooms/`, `/static/rooms-thumbnail.webp`, `/static/web_components/room_map.js`, and `ticket_holder_dropdown.js` and `.css`.
- Database files, server paths and config: `database/events.db`, `events-demo.db`, the `events-<timestamp>.db.zst` backups, `event-images/` and `local-event-images/`, `APP_DB_PATH`, `APP_IMAGE_PATH` and the `-image-path` flag. The variable `eventImageDir` is a code name and does change.
- Words that are not the domain word: DOM events (`EmulatedEvent`, `addEventListener`, `CustomEvent`, `removeEventListener`, `ServerSentEventGenerator`: 35 uses in Go and templ), Server-Sent Events, `github.event` in CI, and event handling in `static/**.js`.
- "Ticket" words that mean the CheckIn ticket rather than billettholder: `ticket_id`, `ticket_type`, `my_tickets`, `profile_tickets`, `post_login_ticket_sync`.
- Norwegian UI text, which already says arrangement.
- Applied migrations in `migrations/2026*.sql`. They are never edited.
- `dev.db` (tracked binary, last changed in 535fa239) and old `plans/` files.

## Inventory

Counts come from `git ls-files` and `git grep` over the 589 tracked files, excluding `.claude/`, at 8033ed17. The commands are in the appendix.

### Files per area that mention a word

Case-insensitive, files/lines. Tracked files per area: Go 99, Go tests 153, templ 98, JS 10, CSS 7, SQL 15, Markdown 44, config 55.

| Area | event | interest | room | ticket holder |
| --- | --- | --- | --- | --- |
| Go (non-test) | 53/1062 | 23/290 | 15/366 | 7/31 |
| Go tests | 98/2651 | 49/918 | 20/793 | 8/56 |
| templ | 67/1545 | 22/472 | 12/436 | 13/84 |
| static JS | 9/182 (mostly DOM events) | 0/0 | 1/24 | 1/4 |
| CSS | 3/71 | 1/2 | 0/0 | 0/0 |
| SQL (migrations, schema.sql, initialize.sql) | 7/137 | 3/10 | 6/99 | 0/0 |
| Markdown | 27/329 | 13/90 | 12/121 | 5/27 |
| config, CI, scripts | 24/118 (mostly `github.event`) | 0/0 | 0/0 | 1/2 |

Paths that contain the word: event 73, interest 24, room 44, ticket_holder 13. The room count includes the 17 static files that stay (15 maps, `rooms-thumbnail.webp` and `room_map.js`). The event count includes 8 icons (`components/icons/assets/icon-event-*.svg`). They are embedded with `go:embed`, so they are code names and do change.

### Database

Uses of each identifier as a whole word, uses/files:

| Identifier | Go | Go tests | templ | migrations (uses) | Markdown (uses) |
| --- | --- | --- | --- | --- | --- |
| `events` | 105/28 | 161/47 | 96/18 | 15 | 113 |
| `interests` | 28/11 | 84/24 | 69/12 | 0 | 26 |
| `rooms` | 56/8 | 92/14 | 95/9 | 10 | 94 |
| `event_id` | 93/24 | 195/47 | 65/19 | 16 | 7 |
| `relation_event_puljer` | 17/12 | 73/29 | 14/7 | 6 | 14 |
| `relation_events_players` | 25/9 | 140/32 | 5/4 | 8 | 7 |
| `interest_level` | 24/10 | 31/17 | 7/3 | 0 | 2 |
| `room_id` | 12/5 | 20/6 | 4/4 | 12 | 5 |
| `room_number` | 15/3 | 11/5 | 4/2 | 6 | 1 |
| `event_type` | 10/6 | 38/17 | 5/4 | 4 | 1 |
| `event_runtime` | 6/5 | 16/15 | 3/3 | 4 | 0 |
| `rooms_published` | 3/2 | 6/3 | 1/1 | 3 | 4 |
| `event_statuses` | 0/0 | 15/15 | 0/0 | 4 | 2 |
| `events_types` | 0/0 | 14/14 | 0/0 | 0 | 2 |
| `event_runtimes` | 0/0 | 14/14 | 0/0 | 0 | 1 |
| `interest_levels` | 0/0 | 7/7 | 0/0 | 0 | 1 |
| `v_events_by_pulje_active` | 1/1 | 0/0 | 1/1 | 8 | 6 |
| `v_event_puljer_active` | 0/0 | 0/0 | 0/0 | 4 | 4 |
| `v_event_summary` | 0/0 | 0/0 | 0/0 | 0 | 0 |

`events`, `interests` and `rooms` also match variables and URL segments. SQL statements that name the table (`FROM`/`JOIN`/`INTO`/`UPDATE`/`TABLE`/`REFERENCES`/`EXISTS`), lines/files:

| Table | Go | Go tests | templ |
| --- | --- | --- | --- |
| `events` | 20/15 | 81/40 | 12/7 |
| `interests` | 13/9 | 27/16 | 3/3 |
| `rooms` | 10/4 | 18/6 | 0/0 |

Files that name any renamed table, column or view: Go 29, Go tests 66, templ 23, Markdown 9.

Outside the app:

- `configuration-as-code/stow/scripts/usr/local/bin/conorganizer-sqlite-backup` checks that the `users` and `events` tables exist and runs `SELECT COUNT(*) FROM events`. It runs every 15 minutes.
- `configuration-as-code/stow/scripts/usr/local/share/conorganizer/anonymize-export.sql` (6 lines) and `migrations/permanent-migrations/anonymize_user_data.sql` (5 lines) use `events`. The first one is used by `conorganizer-export-db`, and therefore by `go tool task download:main`.
- `configuration-as-code/README.md:102` describes the backup check.

### Go identifiers

All Go and templ files:

| Pattern | Distinct identifiers | Uses | Files |
| --- | --- | --- | --- |
| `*Event*` | 422 | 3186 | 195 |
| `*Interest*` | 229 | 1352 | 87 |
| `*Room*` | 173 | 863 | 36 |
| `*TicketHolder*` / `*Ticketholder*` | 27 / 1 | 71 / 4 | 22 / 1 |
| starts with `event` | 97 | 4159 | |
| starts with `interest` | 30 | 919 | |
| starts with `room` | 63 | 1329 | |
| starts with `ticketHolder` / `ticketholder` / `ticket_holder` | 11 | 93 | |

The lowercase counts include SQL text, CSS class pieces and URL segments, because `git grep -w` splits words at `-` and `/`. The most used: `event` 1899, `eventID` 505, `events` 362, `event_id` 353, `eventImageDir` 211, `eventId` 153.

Non-test declarations: top-level `type` lines with Event 30, Interest 17, Room 13, TicketHolder 1. Functions and methods: 56, 34, 24 and 1.

`rom` is the same in singular and plural, so 7 identifier pairs collide (`rooms`/`room`, `Rooms`/`Room`, `actualRooms`, `expectedRooms`, `wantRooms`, `roomsRouter`, `roomsApiRouter`), and 25 Go and templ files use both `room` and `rooms`. See open item 1.

### templ, CSS and signals

- 236 templ components. Names with the word: Event 48, Interest 27, Room 20, TicketHolder 6.
- 3 CSS custom properties: `--event-banner-height`, `--event-status-color`, `--nav-event-height`.
- 106 distinct CSS selectors with the word in templ `<style>` blocks and `.css` files (one is the false positive `.prevent`). Examples: `.program-event-row`, `.event-room-dialog`, `.pulje-interests-summary`.
- 11 HTML `id` values, 23 Datastar signals (`$draggedEventId`, `$currentInterestLevelChoice`, `$_roomNumber`, ...) and 13 JSON tags (`json:"event_id"`, `json:"roomsPublished"`, ...).
- The custom element `room-map` in `static/web_components/room_map.js`. `ticket_holder_dropdown.js` already defines `billettholder-dropdown`.
- The 15 SVG maps in `static/rooms/` hold 560 distinct names with `room` (839 uses): ids like `room-705` and `room-705-door-706`, `id="room-targets"`, and the attributes `data-room-number`, `data-room-target` and `data-room-name`. Those attributes are a contract with `room_map.js`, `service/rooms/maps_test.go`, `pages/admin/rooms/rooms_assignment_page.templ` and `documentation/room-assignment.md`. See open item 7.

### Live buckets

`BucketEvents` 63 uses in 21 files, `BucketInterests` 14/11, `BucketRooms` 19/11. The word `billettholders` appears 21 times. The buckets use `jetstream.MemoryStorage` (`service/live/live.go:380`), so nothing is persisted and renaming them needs no migration. The "Buckets" and page tables in [documentation/live-update-lifecycle.md](../documentation/live-update-lifecycle.md) name them (32 matching lines in that file).

### Routes

All paths stay. The handlers and route setup functions named after them change (`SetupEventRoute`, `eventLayoutRoute`, `roomsRouter`, ...). Distinct path strings found: `/event/{id}`, `/event/api/{id}`, `/event/api/{id}/interest/selected-interest`, `/event/api/{id}/interest/update/interest`, `/admin/rooms/`, `/admin/rooms/api/...` (8 forms), `/admin/rooms/assignment/{pulje}`, `/admin/api/puljefordeling/interest`, `/profile/api/new/{id}/event-in-pulje/{pulje}`, `/profile/api/new/{id}/assign-room/{pulje}`, `/event-images/*`, `/static/rooms/*.svg` and `/static/web_components/room_map.js`.

### Tests

- Test functions with the word in their name: Event 151, Interest 83, Room 64, TicketHolder 9. They show up in the behavior report (`go run ./cmd/testreport`).
- Lowercase test helpers with the word: 67, 33, 14 and 5. `testutil/` has none.
- Fixtures that insert lookup rows: `event_statuses` in 15 test files, `events_types` 14, `event_runtimes` 14, `interest_levels` 7.

### Documentation and logs

- 29 Markdown files mention a word. The most lines: `documentation/room-assignment.md` 85, `documentation/live-update-lifecycle.md` 32, `documentation/billettholdere.md` 27, `README.md` 24, `AGENTS.md` 11, `domeneordbok.md` 10.
- `domeneordbok.md` names the tables and columns (`events.is_in_puljefordeling`, `interests`, `interest_level`, `v_events_by_pulje_active`, `puljer.rooms_published`). `documentation/migrations.md` has the enum-value table.
- The logging skill (`.agents/skills/conorganizer-logging/SKILL.md`, `references/project-logging-baseline.md`) uses `event_id` and `event_service` as examples.
- `.vscode/settings.json` has the spelling words `eventimage`, `eventimgupload`, `eventservice`, `myevents` and `newevent`.
- Log keys in non-test Go: `"event_id"` 60, `"room_id"` 3, `"interest_level"` 1. Component names: `event`, `event_form`, `my_events`, `new_event`, `room_assignment`, `rooms`. No dashboard in `configuration-as-code/grafana-dashboards` filters on them.
- Log message text is not a name and stays. One dashboard depends on it: `conorganizer-main-service-debugging.json` line 586 filters on the message `event image directory`. If that message changes, update the dashboard.

### Baseline scan

`plans/2026-10-18-norwegian-domain-rename-scan.txt` lists every distinct name that contains event, interest, room or ticket holder, from all tracked text files except `.claude/` and `plans/`, plus every such path. It has 2178 lines: 2039 names (14743 uses) and 139 paths. The rescan diffs against it.

It was reviewed by kind on 2026-10-18, not line by line. Besides the areas above, it holds names that stay: paths and file names (`events.db`, `event-images`, `local-event-images`) in `.gitignore`, `.dockerignore`, `Makefile`, `Taskfile.yml`, `deploy/`, `scripts/download-environment`, the systemd units and the images-backup script; `eventstream` in `go.mod`/`go.sum`; "event" in `LICENSE`, `grafana.ini` and `static/datastar.js.map`; and the names of applied migrations.

## How the migration renames tables that views depend on

Tested on 2026-10-18 with sqlite3 3.52.0, loading `schema.sql` into an in-memory database:

- With `PRAGMA legacy_alter_table = OFF` (the default), `ALTER TABLE ... RENAME TO` and `RENAME COLUMN` rewrite the foreign keys in other tables and the views. After `events` → `arrangementer`, the foreign key in `interests` read `REFERENCES "arrangementer"(id)`. Renaming `interest_levels.interest_level` also rewrote the foreign key in `interests`.
- With `legacy_alter_table = ON`, the foreign key in `interests` still said `REFERENCES events(id)` and `v_event_summary` still read `FROM events`. Selecting from the view failed with `no such table: main.events`. The PRAGMA takes effect inside a transaction too.
- **So the pattern under "Changing CHECK constraints" in [documentation/migrations.md](../documentation/migrations.md) must not be used for this rename.** It is for recreating a table under the same name. A rename needs the default mode.
- The views need new names and new `rom_*` aliases anyway. The migration drops them first, renames the tables and columns, and then creates the views again under the new names. With the views gone, nothing else in the schema depends on the renamed objects.

The Up section below passed `PRAGMA integrity_check` (ok) and `PRAGMA foreign_key_check` (no rows) in that test. Afterwards no schema SQL mentioned event, room or interest, and all three new views could be queried.

```sql
-- +goose Up
PRAGMA legacy_alter_table = OFF;

DROP VIEW v_event_summary;
DROP VIEW v_events_by_pulje_active;
DROP VIEW v_event_puljer_active;

ALTER TABLE events RENAME TO arrangementer;
ALTER TABLE interests RENAME TO interesser;
ALTER TABLE rooms RENAME TO rom;
ALTER TABLE event_statuses RENAME TO arrangement_statuses;
ALTER TABLE events_types RENAME TO arrangementer_types;
ALTER TABLE event_runtimes RENAME TO arrangement_runtimes;
ALTER TABLE interest_levels RENAME TO interesse_levels;
ALTER TABLE relation_event_puljer RENAME TO relation_arrangement_puljer;
ALTER TABLE relation_events_players RENAME TO relation_arrangementer_players;

ALTER TABLE relation_arrangement_puljer RENAME COLUMN event_id TO arrangement_id;
ALTER TABLE relation_arrangement_puljer RENAME COLUMN room_id TO rom_id;
ALTER TABLE interesser RENAME COLUMN event_id TO arrangement_id;
ALTER TABLE interesser RENAME COLUMN interest_level TO interesse_level;
ALTER TABLE interesse_levels RENAME COLUMN interest_level TO interesse_level;
ALTER TABLE relation_arrangementer_players RENAME COLUMN event_id TO arrangement_id;
ALTER TABLE arrangementer RENAME COLUMN event_type TO arrangement_type;
ALTER TABLE arrangementer_types RENAME COLUMN event_type TO arrangement_type;
ALTER TABLE arrangementer RENAME COLUMN event_runtime TO arrangement_runtime;
ALTER TABLE rom RENAME COLUMN room_number TO rom_number;
ALTER TABLE puljer RENAME COLUMN rooms_published TO rom_published;

-- CREATE VIEW v_arrangement_summary, v_arrangementer_by_pulje_active and
-- v_arrangement_puljer_active: today's definitions with the new table,
-- column and rom_* alias names.

-- +goose Down
-- The mirror: PRAGMA legacy_alter_table = OFF, drop the three new views,
-- rename the columns back, then the tables, and create the views exactly as
-- they are in schema.sql today.
```

The same PR adds a short "Renaming tables and columns" section to `documentation/migrations.md` with these findings.

## Before you start: rescan

Do this before phase 0, in its own PR that changes only this plan and the baseline. No rename starts until it is merged. The commands are in the appendix.

1. **Start from current main.** Write the commit you scan into the status line at the top.
2. **Diff the names against the baseline.** Run the scan command into `tmp/rename-scan.txt` and diff it with `plans/2026-10-18-norwegian-domain-rename-scan.txt`. Lines starting with `>` are names or paths added since 2026-10-18; lines starting with `<` are gone. Put every `>` line in one of three places: the name mapping and a phase, "Not in scope", or "Open items". Then replace the baseline file with the new scan.
3. **Check the database.** Run `git diff 8033ed17..main -- schema.sql migrations/` and `grep -n -i -E 'event|interest|room' schema.sql`. Add every new table, column, view, index or trigger to the database mapping and to the migration. Then rerun the in-memory rehearsal against the current `schema.sql`, with the full Up section and the `CREATE VIEW` statements written out. `integrity_check` must return ok, `foreign_key_check` must return no rows, and the leftover query must return nothing.
4. **Check packages, files, buckets and routes.** Look for new folders and files in `git diff --name-status 8033ed17..main`, new constants in `service/live/live.go`, and new paths with the route command.
5. **Check outside the app.** Run the server-script command. Anything that reads a renamed table must handle both schemas in phase 0 if it runs on the server, or change in phase 6.
6. **Recount.** Rerun the counting commands and update the Inventory, so the phase sizes are right.
7. **Check open PRs.** Run `gh pr list`. A rename PR conflicts with almost everything, so merge or wait for large open PRs that touch the same files.
8. **Check the rules.** Look for new domain words in `domeneordbok.md` and changes to "Names" in `AGENTS.md`. Go through "Decisions" and "Open items" again.

Done when every line of the new scan is in the mapping, "Not in scope" or "Open items", and the counts, the migration and the baseline match the code you scanned.

## Phases

Each phase is one PR, or a few for the large ones. Each one builds, passes `go tool task test`, and can ship alone. The Go renames come before the database because the compiler checks them and they need no rollout. One word per PR keeps the diffs reviewable. Each PR is a pure rename: no behavior changes and no refactors.

How to do a Go phase:

- Move files and folders with `git mv` so history follows.
- Rename identifiers with `gopls rename` where it works, then run `go tool templ generate`. In a worktree you must run it yourself.
- Find what is left with the inventory commands in the appendix. Only the not-in-scope uses should remain.
- Bulk edits from the main checkout also rewrite `.claude/worktrees/`, so stage only your own files. Multi-line `sed`/`perl` edits silently miss on CRLF files; check with `git diff --stat`.

### Phase 0: groundwork

- Finish [Before you start: rescan](#before-you-start-rescan).
- Merge #663. It adds the "Names" rule and deletes `initialize.sql`, so that file needs no update.
- Make `conorganizer-sqlite-backup` accept either `events` or `arrangementer` in its table check and its `COUNT`. Pull it into the server checkout at `/srv/configuration-as-code-repo/conorganizer`, whose files the scripts are symlinked to. Without this, every backup fails from the migration until someone updates the checkout.
- Settle the open items below.

### Phase 1: live buckets (1 small PR)

Rename the four constants and bucket names, and update the two tables in `documentation/live-update-lifecycle.md`. The buckets are in memory, so a restart starts with the new names and no migration is needed.

### Phase 2: ticketholder → billettholder (1 PR)

`components/ticket_holder` → `components/billettholder` (11 files, 8 importers), `ticket_holder_*.templ` → `billettholder_*.templ`, and the 27 `TicketHolder` identifiers (71 uses) with their CSS classes and signals. The static files `ticket_holder_dropdown.js` and `.css` keep their paths, because those are URLs.

### Phase 3: room → rom (1–2 PRs)

`service/rooms` → `service/rom`, `pages/admin/rooms` → `pages/admin/rom`, `models/rooms.go` → `models/rom.go` (`Room` → `Rom`, `RoomInput` → `RomInput`, `RoomEventPuljeSummary` → `RomEventPuljeSummary`, ...), plus `pages/event/event_rooms.*`. Also the CSS classes, `$_room*` signals, log key `room_id`, and components `rooms` → `rom` and `room_assignment` → `rom_assignment`. Paths under `/admin/rooms/` and `/static/rooms/` stay. The SVG map contents follow open item 7. Size: 173 distinct `*Room*` identifiers with 863 uses in 36 files, plus 63 lowercase `room*` names.

### Phase 4: interest → interesse (1–2 PRs)

`models/interests-model.go` → `models/interesser-model.go` (`Interest` → `Interesse`, `InterestLevel` → `InteresseLevel` and its `High`/`Medium`/`Low`/`None` constants), `event_interests.templ`, `programpulje_interests.templ`, `interest_indicator.templ`, `billettholder_interest_*.templ`, `pages/root/billettholder_interests.go`, the signals, and the log key `interest_level`. The paths under `/event/api/{id}/interest/` stay. Size: 229 distinct identifiers with 1352 uses in 87 files, plus 30 lowercase names.

### Phase 5: event → arrangement (3 PRs)

- **5a models and services**: `models/event-model.go` → `arrangement-model.go` (`Event`, `EventStatus*`, `EventType*`, `EventPlayerRole*`, `EventCardModel`), `service/eventService`, `service/eventimage`, `service/program/events.go`.
- **5b pages and components**: `pages/event` (22 files), `components/event_components`, `pages/profile/newevent`, `event_img_upload`, `event_card.templ`, `my_events.templ`, `banner_more_events.templ`, `eventStatus.templ`, `eventBar.templ`, `pages/root/event_list.templ`, the `icon-event-*` icons, CSS classes and signals.
- **5c the rest**: test names and helpers (151 test functions, 67 helpers), the log key `event_id` (60 uses) and component names, `eventImageDir` (211 uses), and the logging skill and `.vscode/settings.json`.

Size: 422 distinct `*Event*` identifiers with 3186 uses in 195 files, plus 97 lowercase `event*` names. `/event/...`, `/event-images/` and `local-event-images/` stay.

### Phase 6: database (1 PR, then a migration by hand)

The PR contains:

- The migration above.
- SQL strings in 29 Go files, 66 test files and 23 templ files.
- `schema.sql`, regenerated by `go tool task test` after you migrate your local database.
- `anonymize-export.sql` and `migrations/permanent-migrations/anonymize_user_data.sql`.
- Docs: `documentation/migrations.md` (enum table and the new section), `domeneordbok.md`, `documentation/room-assignment.md`, `documentation/pulje-status-and-publishing.md`, `documentation/billettholdere.md`, `AGENTS.md`, `configuration-as-code/README.md:102`, and a note in `migrations.md` and `configuration-as-code/README.md` that older backups need `goose up` (see below).

Rehearse before merging:

1. Run `go tool task download:main`, make a `.backup` copy, and run `goose up`, `down` and `up` on the copy.
2. Run `PRAGMA integrity_check` and `PRAGMA foreign_key_check` on the copy.
3. Run `go tool task test`.
4. Click through the event page, admin rooms, room assignment, puljefordeling, puljeoppsett and profile locally.

Rollout follows "Apply migrations on the server" in `documentation/migrations.md`:

1. `sudo conorganizer-maintenance-mode on`. It has no effect while `program.regncon.no` redirects after the con, but it matters if that block has been restored.
2. Merge and wait for `deploy-fixed`. Main, demo and restored now run the new binary on the old schema, and every query fails until they are migrated.
3. `sudo conorganizer-sqlite-backup`. This works on either schema because of phase 0.
4. `sudo conorganizer-sqlite-migrate`. It migrates demo, then main, with a `quick_check` after each.
5. Update the server checkout `/srv/configuration-as-code-repo/conorganizer` to main, so `anonymize-export.sql` uses the new names. Until then, `download:main` fails.
6. Check `main.lekeplassen.regncon.no` and `demo.lekeplassen.regncon.no`, then turn maintenance mode off.
7. `restored` is never migrated by the script. It stays broken until a backup taken after step 4 is restored into it with `sudo conorganizer-sqlite-restore events-<timestamp>.db.zst`.
8. Locally, run `goose up` on `database/events.db` (in Docker, through `docker compose exec`), or download a fresh copy after step 5.

Lasting consequence: every backup from before step 4 has the old schema. Restoring one into main or demo needs `goose up` afterwards, and so does restoring one into `restored`.

Rollback: restore the step 3 backup and revert the PR. The data written since the migration is lost; maintenance mode keeps that small. `goose down` is also possible, but only while the deployed `migrations/` folder still holds the migration file, so run it before the revert is deployed.

PR previews:

- The phase 6 PR's own preview clones main (old schema) on its first deploy, then `goose up -allow-missing` migrates it, so the preview doubles as a rehearsal. If the migration changes after the first deploy, delete the preview's database so it clones again, because Goose does not re-run an applied migration.
- Existing previews of other PRs keep their old cloned database and keep working until the branch merges main. The next deploy then migrates them.
- A preview created after main is migrated clones the new schema but runs the PR's old code, so it is broken until the branch merges main. Merge main into every open PR after the rollout. Today #663 is the only open PR, and phase 0 merges it.
- Phases 1–5 do not touch the database, so previews are not affected. Open PRs only get merge conflicts.

## Open items

These are suggestions to approve or change before the phase that needs them.

1. **Plural `rom`** (phase 3). Suggestion: when both forms are in scope, the plural gets the suffix `List`: `rooms` → `romList`, `Rooms` → `RomList`, `expectedRooms` → `expectedRomList`. That covers 7 pairs in 25 files. Usikker.
2. **The `room-map` element** (phase 3). The plan renames the tag to `rom-map`, because a tag is an HTML name and not a URL, but keeps the file path `static/web_components/room_map.js`. Usikker: should the tag stay too, to match the file?
3. **Package directory style** (phases 2–5). The plan keeps today's style (`eventService` → `arrangementService`, `event_components` → `arrangement_components`), so the diff stays a pure rename. Cleaning up the style is separate work.
4. **`events_types`** (phase 6). A pure rename gives `arrangementer_types`, which keeps today's odd plural. `arrangement_types` would match `arrangement_statuses` and `arrangement_runtimes`. Usikker.
5. **`v_event_summary` is unused** (0 uses in code). The plan renames it. Dropping it would be a separate change.
6. **Saved Loki queries** outside the repo that filter on `event_id` or `component=event` stop matching new logs after phase 5c.
7. **SVG map contents** (phase 3). The maps hold 560 room names, and `data-room-number`/`data-room-target` are read by `room_map.js` and `maps_test.go`. Suggestion: leave the SVG files and these attribute names as they are, because they are map data under a URL that stays, and keep the code that reads them pointing at `data-room-*`. The other choice is to rename them all in one PR together with `room_map.js`, `maps_test.go`, `rooms_assignment_page.templ` and `room-assignment.md`. Usikker.

## Appendix: commands

Run from the repository root (Git Bash) with `set -f` so that pathspecs are not expanded as globs.

```bash
# full scan: every distinct name and path with a word, compared with the baseline
mkdir -p tmp
{ git grep -I -o -h -i -E '[A-Za-z0-9_-]*(event|interest|room|ticket.?holder)[A-Za-z0-9_-]*' -- ':!.claude' ':!plans' | sed 's/^/name /'
  git ls-files -- ':!.claude' ':!plans' | grep -i -E 'event|interest|room|ticket.?holder' | sed 's/^/path /'
} | LC_ALL=C sort -u > tmp/rename-scan.txt
diff plans/2026-10-18-norwegian-domain-rename-scan.txt tmp/rename-scan.txt

# schema objects that still mention a word (run after the Up section in the rehearsal)
#   SELECT type, name FROM sqlite_master
#   WHERE sql LIKE '%event%' OR sql LIKE '%room%' OR sql LIKE '%interest%';

# server scripts, deploy and CI that name a table or column
git grep -n -w -E 'events|interests|rooms|event_id|room_id|interest_level' -- configuration-as-code deploy scripts Taskfile.yml .github

# URL path strings with a word
git grep -o -h -E "[\"'\`(]/[A-Za-z0-9_{}./%-]*(event|room|interest)[A-Za-z0-9_{}./%-]*" -- '*.go' '*.templ' '*.js' ':!.claude' ':!*_test.go' | sort | uniq -c

# paths that contain a word
git ls-files -- ':!.claude' | grep -i -c event

# files/lines per area (example: Go non-test, word "room")
git grep -l -i -E room -- '*.go' ':!*_test.go' ':!.claude' | wc -l
git grep -h -i -E room -- '*.go' ':!*_test.go' ':!.claude' | wc -l

# database identifier, uses/files (example: event_id in templ)
git grep -o -w event_id -- '*.templ' ':!.claude' | wc -l
git grep -l -w event_id -- '*.templ' ':!.claude' | wc -l

# SQL statements that name a table
git grep -h -i -E '(from|join|into|update|table|references|exists) +"?events"?([^A-Za-z0-9_]|$)' -- '*.go' ':!*_test.go' ':!.claude' | wc -l

# CamelCase identifiers: distinct / uses / files
git grep -o -h -E '[A-Za-z0-9_]*Event[A-Za-z0-9_]*' -- '*.go' '*.templ' ':!.claude' | sort -u | wc -l
git grep -o -h -E '[A-Za-z0-9_]*Event[A-Za-z0-9_]*' -- '*.go' '*.templ' ':!.claude' | wc -l
git grep -l -E '[A-Za-z0-9_]*Event[A-Za-z0-9_]*' -- '*.go' '*.templ' ':!.claude' | wc -l

# templ components, test functions and types with the word
git grep -h -E '^templ [A-Za-z0-9_]*[Ee]vent[A-Za-z0-9_]*\(' -- '*.templ' ':!.claude' | wc -l
git grep -h -E '^func Test[A-Za-z0-9_]*Event' -- '*_test.go' ':!.claude' | wc -l
git grep -h -E '^type [A-Za-z0-9_]*[Ee]vent[A-Za-z0-9_]* ' -- '*.go' '*.templ' ':!*_test.go' ':!.claude' | wc -l

# live buckets
git grep -o -w BucketEvents -- ':!.claude' | wc -l

# rename rehearsal in memory (writes no file)
printf '%s\n' '.read schema.sql' 'ALTER TABLE events RENAME TO arrangementer;' \
  "SELECT sql FROM sqlite_master WHERE name='interests';" | sqlite3 :memory:
```
