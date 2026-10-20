# Pulje status and program publishing

This document describes how pulje status and program publishing control interest, the event page, puljefordeling and Mitt festivalprogram, and how the per-pulje publish switch for the room assignment fits in. For the meaning of domain words such as pulje, interesse and billettholder, see [domeneordbok.md](../domeneordbok.md). For event statuses, including `Annonsert`, see [event-status.md](event-status.md).

## What "publisert" means

"Publisert" means several different things in the code. Keep them apart:

1. **The event is announced**: `events.status = 'Annonsert'` (`models.EventStatusAnnounced`). The event is public. Call this "annonsert", never "publisert". See [event-status.md](event-status.md).
2. **The program is published**: `program_publishing_state.is_published`, one row with `id = 1`, read by `program.IsPublished`. It covers the whole program: the front page with days and puljer, interest, Mitt festivalprogram, and the puljer, times and rooms on event pages. Toggled with "Publiser program" on `/admin`.
3. **The puljefordeling is published**: `puljer.status = 'Completed'`, shown as "Puljefordeling publisert" in admin. The assignments in that pulje are published, and the commit and tildeling paths refuse changes (`ErrPuljeCompleted`) while the pulje stays `Completed`. An admin can untick step 3 to move the pulje back to `Locked` and allow changes again.
4. **The room assignment is published**: `puljer.rooms_published`, per pulje, toggled with the "Publiser romfordeling" switch on `/admin/rooms/assignment/{pulje}` ("Synlig for alle" / "Ikke publisert"). Only then are the room and the room map for that pulje shown on the event page and the print-friendly page, to admins too, and the event page also shows the room's public notes. The event page also requires the program to be published. The flag is not tied to `puljer.status`.
5. **`relation_event_puljer.is_published`**: a legacy column per event and pulje. No production code reads it, and `SetEventInPulje` deliberately leaves it untouched. Do not build new behavior on it.

These fields are not about publishing, but are easy to confuse with it:

- `relation_event_puljer.is_in_pulje`: whether the event is set up in the pulje. `v_event_puljer_active` filters on `is_in_pulje = 1`, and `v_events_by_pulje_active` also filters on status `Annonsert`.
- `events.is_in_puljefordeling`: whether the event takes part in the puljefordeling, or is a [program event](#program-events-outside-the-puljefordeling).
- `puljer.status` (`Open`, `Locked`, `Completed`) publishes only the puljefordeling (item 3). It does not control which events are shown publicly or in the previous/next navigation.

## Independent settings

Three separate settings control pulje behavior. They are independent domain concepts and are kept apart on purpose.

| Setting | Scope | Controls |
| --- | --- | --- |
| `puljer.status` (`Open` / `Locked` / `Completed`) | Per pulje | Whether billettholdere can add or change interest in that pulje, and whether solver assignments are shown as results in Mitt festivalprogram. |
| `program_publishing_state.is_published` (0/1) | Global | Front page layout, whether interest controls are shown and interest updates are accepted on event pages, whether Mitt festivalprogram is shown, and whether the schedule (puljer, times and rooms) is shown on event pages. |
| `puljer.rooms_published` (0/1) | Per pulje | Whether that pulje's room name, public room note and room map are shown on event pages, and its room and map on the print-friendly page. It adds to the program flag: event pages show a room only when both are set. It is not tied to `puljer.status`. |

`puljer.rooms_published` defaults to `0` and is toggled with the "Publiser romfordeling" switch on `/admin/rooms/assignment/{pulje}`. The event page and the print-friendly page hide an unpublished pulje's room from admins too. See [room-assignment.md](room-assignment.md#publishing-the-room-assignment).

Do not encode program publishing or room publishing in `puljer.status`, and do not use the legacy `relation_event_puljer.is_published` column for any of the rules below. None of the code paths described here read it.

### `program_publishing_state`

A single-row global table:

```sql
CREATE TABLE program_publishing_state(
  id INTEGER NOT NULL PRIMARY KEY CHECK(id = 1),
  is_published INTEGER NOT NULL DEFAULT 0 CHECK(is_published IN(0, 1))
) STRICT;
```

- Migration `20260522150000_program_publishing_state.sql` seeds it with `(1, 0)`, so the program starts unpublished. `schema.sql` has the table but no row, so tests that need one insert it themselves.
- `program.IsPublished` (`service/program/publishing.go`) treats a missing row as not published. Fresh databases still need the table and its seed row.
- Admins toggle it with the "Publiser program" card on `/admin`, which calls `PUT /admin/api/program-publishing` (`pages/admin/publiser_program.templ`). The handler upserts the row with `id = 1` and broadcasts the `events` bucket (see [live-update-lifecycle.md](./live-update-lifecycle.md)).

Front page (`pages/root/root_page.templ`):

- Published: a day selector with pulje blocks, built from `v_events_by_pulje_active` (`service/program/days.go`, `service/program/events.go`).
- Not published: a plain alphabetical list of events with status `Annonsert` (`program.GetAnnouncedEvents`).

## Pulje status values

`puljer.status` is `TEXT NOT NULL DEFAULT 'Open'` with `CHECK(status IN('Open', 'Locked', 'Completed'))` and a foreign key to `pulje_statuses`. In Go the values are `models.PuljeStatusOpen`, `models.PuljeStatusLocked` and `models.PuljeStatusCompleted` (`models/pulje-model.go`). `PuljeStatus.Label()` shows them as Åpen, Låst and Fullført.

| Status | Meaning |
| --- | --- |
| `Open` | Billettholdere can add or change interest. |
| `Locked` | Existing interest is shown read-only and cannot be changed while players are being assigned. |
| `Completed` | The puljefordeling is published. Interest stays frozen, and solver assignments become visible in Mitt festivalprogram. |

The table also has `closing_warning_active` (`INTEGER NOT NULL DEFAULT 0`, added in `20260920100000_add_pulje_closing_warning.sql`). It is the manual "closes soon" warning described below.

Only `Open`, `Locked` and `Completed` exist. Older values (`not_published`, `published` and lowercase `open`/`locked`/`completed`) were migrated away by `20260522120000_pulje_status_open_locked_completed.sql` and `20260524100000_capitalize_pulje_statuses.sql` and removed from `pulje_statuses`. Do not bring back `PuljeStatusNotPublished` or `PuljeStatusPublished`. The capitalized values match the other lookup tables, such as `event_statuses` (`Kladd`, `Innsendt`, `Godkjent`, `Forkastet`, `Annonsert`) and `events_types`.

## Admin workflow on `/admin/puljefordeling/{pulje}`

Only admins change pulje status. They do it in three ordered steps on the puljefordeling page for each pulje. The same page also handles the distribution itself: assigning players, saving (committing) the distribution and the preview dialogs.

| Step | Label | Checked | Unchecked |
| --- | --- | --- | --- |
| 1 | Vis advarsel om at puljefordeling stenger snart | `closing_warning_active = 1` | `closing_warning_active = 0` |
| 2 | Puljefordeling lukket | `Locked` | `Open` |
| 3 | Puljefordeling publisert | `Completed` | `Locked` |

Rules:

- Step 1 can only change while the pulje is `Open`.
- Step 2 (lock) needs step 1 done. Going back to `Open` is rejected once the pulje is `Completed`.
- Step 3 (publish) needs the pulje to be `Locked` and the distribution to have no unsaved changes. Going back from `Completed` to `Locked` is allowed.
- Moving to `Locked` or `Completed` clears `closing_warning_active`. A pulje that is no longer `Open` counts step 1 as done.
- Every toggle asks for confirmation with `confirm()`. Cancelling reverts the checkbox.

In the UI, a step is disabled until the previous step is done, and it cannot be undone once the next step is done. Each step shows a caption:

- `Neste steg`: the step can be done now.
- `Fullført`: the step is done.
- `Krever steg 1` / `Krever steg 2`: the previous step is not done yet.
- `Lagre fordelingen først`: step 3 while the pulje is `Locked` and the distribution has unsaved changes.

The server enforces the same order in `puljeStatusStepAllowed`:

- A change out of order returns `409` with `Forrige steg må være fullført først`.
- Publishing with unsaved changes returns `409` with `Lagre fordelingen før den publiseres`.
- An unknown pulje returns `404`, and an invalid status returns `400`.

Both status handlers broadcast the `events` bucket.

Code locations:

- UI: `puljeStatusToggles` and `puljeStatusStep` in `pages/admin/puljefordeling_tab.templ`.
- Handlers `PUT /admin/api/puljer/{puljeId}/status` and `PUT /admin/api/puljer/{puljeId}/closing-warning`, together with the admin DB helpers (`getPuljer`, `updatePuljeStatus`, `updatePuljeClosingWarning`, `puljeStatusStepAllowed`), in `pages/admin/puljefordeling.go`. They are admin-local, not in the service layer.
- `puljefordelingStatusRoute` registers them in `SetupAdminRoute` (`pages/admin/admin.go`).
- Tests: `pages/admin/puljefordeling_status_steps_test.go` and `TestPuljeStatusRoute_PublishingRequiresSavedDistribution` in `pages/admin/puljefordeling_commit_completed_test.go`.

`components/formsubmission/puljefordeling.templ` is unrelated. It is the part of the event form where you choose which puljer an event belongs to.

## Pulje times and Oslo wall-clock time

- `puljer.start_at` and `puljer.end_at` are RFC3339 `TEXT`. No application code writes them.
- `DBDateTime` parses them and keeps the stored offset. `PuljeRow.TimeRange()` formats the stored clock time as is (`Format("15:04")`, no timezone conversion). The stored value must therefore hold Oslo wall-clock time.
- Test fixtures use explicit Oslo offsets, for example `2026-10-09T18:30:00+02:00`.
- The front page decides which day a pulje belongs to by converting `start_at` to `Europe/Oslo` (`service/program/days.go`).
- No time-based logic uses these timestamps for warnings or locking.

## Interest messages and states

`BuildPuljeInterestState` (`components/ticket_holder/ticket_holder.go`) derives the state only from `puljer.status` and `puljer.closing_warning_active`. Its `now` argument is ignored, so messages never depend on the clock.

| State | Condition | Message | Editing | Color | Picker label |
| --- | --- | --- | --- | --- | --- |
| open | `Open`, no warning | none | allowed | none | none |
| warning | `Open`, `closing_warning_active = 1` | Viktig: Puljefordelingen stenger snart. Gjør endringer nå hvis du vil endre interessene dine. | allowed | `--color-error` | Låses snart |
| locked | `Locked` | Puljen er låst. Du kan ikke melde eller endre interesse lenger. Vi jobber med å fordele spillere. | disabled | `--color-accent-blue` | Låst |
| completed | `Completed` | Puljefordelingen er klar. Se hva du fikk på profilen din. | disabled, profile link shown | `--color-success` | Klar |

The warning is informational. An admin turns it on manually (step 1 above), it never blocks editing, and locking or completing clears it.

## Program events outside the puljefordeling

There is no separate "påmelding" (sign-up) feature. Long-running events that anyone can join, such as "Blood on the Clocktower" or "Cosplay", are entered as program events (programarrangementer): events with `events.is_in_puljefordeling = 0`. Admins set this with the "Med i puljefordeling" checkbox in the "Puljefordeling" card of the event form (`puljefordeling` in `components/formsubmission/puljefordeling.templ`, rendered only for admins, saved by `PUT /profile/api/new/{id}/is-in-puljefordeling`). A program event is still placed in one or more puljer with `is_in_pulje = 1`, so that it appears on the right day.

A program event behaves differently from an event in the puljefordeling:

- The event page shows `ProgramEventInterestPanel` (`pages/event/event_interest_panel.templ`) instead of the interest panel: "Dette er et programarrangement som er åpent for alle. Du trenger ikke melde interesse eller reservere plass." When the viewer's billettholdere have interests in the selected pulje, the panel also warns that they may be picked for one of those events, with a link to `/profile`. The interest dialog is never rendered (`canShowInterestControls` requires `is_in_puljefordeling = 1`), and `updateInterest` rejects interest for it.
- The front page lists each day's program events once, under the heading "Programoversikt" above the pulje blocks (`pages/root/event_list.templ`). `program.GetDays` (`service/program/days.go`) links a program event to its first pulje on that day, and `PuljeBlock.RaffleEvents` leaves program events out of the pulje blocks.
- Times: program events run on their own schedule, given in the description. The event page rooms and the print-friendly page show `ProgramEventTimeHint` ("Se beskrivelse for tidspunkt.") instead of the pulje time. In Mitt festivalprogram the program event's card (`ProgramPuljeEvent`) shows the hint as an extra line, while the pulje heading keeps the pulje time.
- The solver only loads events with `is_in_puljefordeling = 1` (`loadEligibleEvents` in `service/puljefordeling/emulate.go`).

The UI still uses the word "påmelding" for showing interest in three places: the "Publiser program" card on `/admin` says "Publiser programmet for festivalen, og åpne for ønsker og påmeldinger.", "Påmelding" is a feedback topic (`service/feedback/topic.go`), and the feedback form's placeholder (`pages/feedback/feedback_page.templ`) begins "F.eks. «Påmeldingen var enkel, …»". All of them mean interesse, not a separate sign-up.

## Event page interest panel and dialog

For an event that is in puljefordeling (`events.is_in_puljefordeling = 1`), the event page (`pages/event/event_page.templ`) has two separate parts: the outer panel and the "Meld interesse" dialog. Events that are not in puljefordeling show `ProgramEventInterestPanel` instead, which says that no interest is needed. See [Program events outside the puljefordeling](#program-events-outside-the-puljefordeling). For which users can see an event at all, see [access-control-and-error-pages.md](./access-control-and-error-pages.md).

### Outer panel (`EventInterestPanel`)

In `pages/event/event_interest_panel.templ`:

- Program not published: everyone sees `Interessevalget åpner når programmet er publisert.`, with no `Hent billett` button.
- Program published, user has no tickets: a `Hent billett` button linking to `/profile/tickets`.
- Program published, user has tickets: the `Meld interesse` button is always shown. It is never replaced by status text, because an event can be in several puljer and a user can have several billettholdere. A status line may appear above it (see below). If interest controls are unavailable, the panel shows `Interessevalg er ikke åpnet for dette arrangementet ennå.`

The status line is strictly tied to the `?pulje=` query parameter. It uses `BuildQueryPuljeInterestState`, which returns nothing when:

- the parameter is missing,
- the value fails `models.ParsePulje`,
- the value names a pulje the event is not in (`is_in_pulje`), or
- that pulje's state has no message.

It never falls back to another pulje.

### Dialog (`EventInterests`)

The dialog is rendered only when `canShowInterestControls` (`pages/event/interest_availability.go`) is true. That requires all of:

- `program_publishing_state.is_published = 1`
- `events.is_in_puljefordeling = 1`
- at least one pulje with `relation_event_puljer.is_in_pulje = 1`

`GetPujerForEvent` and `GetActivePuljeForEvent` filter only on `is_in_pulje`.

The dialog picks its pulje differently from the outer panel, on purpose. `event_page_content` replaces a missing, invalid or unattached `?pulje=` value with the event's first pulje (ordered by `start_at`), so the dialog always has a usable default. Do not unify this with the outer panel.

Inside the dialog:

- The pulje picker is shown only when the event has more than one pulje. Every pulje stays selectable, and each non-open state is marked with an icon and label: warning `Låses snart`, locked `Låst` (lock icon), completed `Klar` (complete icon).
- The interest-level buttons become read-only (`$puljeCanEdit = false`) only when the selected pulje is `Locked` or `Completed`. The saved interest stays visible and the status message is shown. For `Completed`, a `Se profilen din` button links to `/profile`.

### Two helpers that must stay separate

`components/ticket_holder/ticket_holder.go` has two helpers for choosing a pulje:

- `BuildQueryPuljeInterestState(puljer, puljeQuery, now)` is strict. It returns a state only for the pulje named in the query value and never falls back. Use it for status shown from the URL, as the outer panel does.
- `BuildSelectedPuljeInterestState(puljer, puljeID, now)` falls back to the first pulje. Use it where a default is needed, as the `EventInterests` dialog does.

## Server-side enforcement in `updateInterest`

Hiding or disabling the interest picker in the UI is not a security boundary. `updateInterest` in `pages/event/event.go` checks everything again on the server before writing, so direct API calls are rejected too. It requires:

- `program_publishing_state.is_published = 1`
- a `relation_event_puljer` row for the event and pulje with `is_in_pulje = 1`
- `events.is_in_puljefordeling = 1`
- `events.status = 'Annonsert'`
- `puljer.status` that is neither `Locked` nor `Completed`
- the logged-in user linked to the billettholder through `relation_billettholdere_users`

## Mitt festivalprogram visibility

Mitt festivalprogram is `components/profile/my_program.templ`.

- Until `program_publishing_state.is_published = 1`, it shows only `Programmet for Regncon er ikke publisert ennå`.
- Once the program is published, an assignment or interest appears only if `events.status = 'Annonsert'` and `relation_event_puljer.is_in_pulje = 1` for that event and pulje.

Assignments from `relation_events_players` are shown as follows:

| Assignment | Shown when |
| --- | --- |
| Player assigned by the solver (`source = 'solver'`) | Only when the pulje is `Completed`. Until then, the pulje shows the billettholder's interests instead. |
| Player assigned manually (`source = 'manual'`) | Straight away, in any pulje status. This matches the assignment notice in the interest dialog. |
| GM (`role = 'GM'`) | In every pulje status. |

- GM means `relation_events_players.role = 'GM'`, not `events.user_id`, which is the event's creator.
- If the billettholder is both Player and GM on the same event, one card marked as GM is shown (`MAX(role = 'GM')`).
- When a pulje has any visible assignment, the interest list for that pulje is hidden.

Tests: `components/profile/my_program_test.go` and `components/profile/my_program_render_test.go`.
