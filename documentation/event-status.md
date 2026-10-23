# Event status

This document describes the five statuses an event (arrangement) can have, who can change them and where, and what each status allows. For who can see an event in each status, see [access-control-and-error-pages.md](access-control-and-error-pages.md). For what "publisert" means, see [pulje-status-and-publishing.md](pulje-status-and-publishing.md#what-publisert-means).

## The five statuses

`events.status` is `TEXT NOT NULL DEFAULT 'Kladd'` with a foreign key to the lookup table `event_statuses` (`ON UPDATE CASCADE`). In Go the values are the `models.EventStatus` constants in `models/event-model.go`. An event always has exactly one of them.

| Status | Go constant | Meaning |
| --- | --- | --- |
| `Kladd` | `EventStatusDraft` | The organizer is still working on the event. |
| `Innsendt` | `EventStatusSubmitted` | Sent to the board. |
| `Godkjent` | `EventStatusApproved` | Approved by the board, but not public. Can be placed in a room on the room-assignment page. |
| `Annonsert` | `EventStatusAnnounced` | Public. The only status shown on the front page and in the previous/next navigation. Can be placed in a room on the room-assignment page. |
| `Forkastet` | `EventStatusArchived` | Rejected or withdrawn. |

- There is no `Avvist` status. A rejected event is `Forkastet`, and its Go constant is `EventStatusArchived`.
- `Annonsert` replaced the old status `Publisert` (`EventStatusPublished`), which no longer exists. Migration `20260531080000_rename_published_event_status_to_announced.sql` converted every `Publisert` and `Godkjent` event to `Annonsert`, removed `Publisert` from `event_statuses` and changed `v_events_by_pulje_active` to filter on `Annonsert`. `Godkjent` is still a valid status. An event is "annonsert", never "publisert".
- `EventStatus.Label()` returns the value itself, so the UI shows `Kladd`, `Innsendt`, `Godkjent`, `Annonsert` and `Forkastet`. `button.EventStatus` (`components/ui/button/eventStatus.templ`) renders the label with a progress icon.

## Who changes the status, and where

### The event form

The status card (`statusCard` in `components/formsubmission/statusCard.templ`) is a `<select>` that calls `PUT /profile/api/new/{id}/status` on change. The options depend on the viewer:

- An admin gets all five statuses. `Godkjent` and `Annonsert` carry the suffix "(du er administrator)".
- A non-admin owner gets `Kladd` and `Forkastet`, plus a disabled `Innsendt` "(bruk knappen nederst i skjemaet)", and only while the event is `Kladd`, `Innsendt` or `Forkastet`. In any other status the select has no options.

`FormBody` (`components/formsubmission/form_body.templ`) renders the status card for everyone and the "Puljefordeling" card (`is_in_puljefordeling` and the puljer) only for admins.

Owners edit at `/profile/new/{id}`. That page (`pages/profile/newevent/new_page.templ`) refuses a non-admin while the event is `Godkjent`, showing "Arrangementet er allerede godkjent eller annonsert" and a note to contact the board. The check looks only at `Godkjent`, not at `Annonsert`. Admins edit the same form at `/admin/approval/edit/{id}` (`pages/admin/approval/editForm/edit_form_page.templ`), where `isAdmin` comes from `authctx.GetAdminFromUserToken`.

### "Send inn"

A non-admin submits with the "Send inn" button in `SubmitSection` (`components/formsubmission/submit_section.templ`). `POST /profile/api/new/{id}/submit` sets the status to `Innsendt` (`submitEventWithAudit`) and redirects to `/profile/`.

### The server side

Every route under `/profile/api/new/{id}` runs `newevent.RequireEventOwnerOrAdmin` (`pages/profile/newevent/new_index.templ`), so only the owner or an admin reaches the status handler. The handler itself, `UpdateStatus` in `statusCard.templ`, writes whatever status the `status` signal holds. It does not limit a non-admin to `Kladd` and `Forkastet`; only the foreign key to `event_statuses` rejects unknown values. The options hidden in the template are not a security boundary.

Every status change goes through `updateEventStatusWithAudit` (`components/formsubmission/form_body.templ`), which also sets `status_changed_at`, `status_changed_by_id` (`users.id`) and `status_changed_action` (the new status) on the `events` row.

### Admin approval page

`/admin/approval` (`pages/admin/approval/approval_page.templ`) lists every event grouped by status, in the order `Innsendt`, `Godkjent`, `Annonsert`, `Kladd`, `Forkastet`. Each card links to `/admin/approval/edit/{id}`. The page has no approve or announce buttons of its own; the admin changes the status in the form's status card.

## What each status allows

| Area | Rule |
| --- | --- |
| Event page visibility | `Annonsert` is public. The other statuses are shown only to the owner and admins, and `Forkastet` answers 410 Gone to everyone else. See [access-control-and-error-pages.md](access-control-and-error-pages.md#event-visibility). |
| Front page | Only `Annonsert`: `program.GetAnnouncedEvents` while the program is unpublished, and the view `v_events_by_pulje_active` (`e.status = 'Annonsert'` and `is_in_pulje = 1`) once it is published. |
| Previous/next on the event page | `GetPreviousNextForRootEventList` (`service/eventService/previous_next_root_list.go`) walks the same lists as the front page, so only `Annonsert` events appear. `GetPreviousNextInnsendtGodkjent` in `getPreviousNext.go` filters on `Innsendt` and `Godkjent` instead, but only its test calls it. |
| Interest | `updateInterest` requires `events.status = 'Annonsert'`. See [pulje-status-and-publishing.md](pulje-status-and-publishing.md#server-side-enforcement-in-updateinterest). |
| Mitt festivalprogram | Shows assignments and interests only for `Annonsert` events. |
| Room assignment | The room-assignment page accepts only `Godkjent` and `Annonsert` events: its dialog lists those two statuses, and `POST /admin/rooms/api/assignment/{pulje}/{event}/{room}` answers 409 "Arrangementet er ikke godkjent." for anything else (`pages/admin/admin.go`). The "Legg til rom" select in the event form's Puljefordeling card (`UpdateRoomInPulje` in `components/formsubmission/puljefordeling.templ`, `PUT /profile/api/new/{id}/assign-room/{puljeId}`) has no status check, so it can give a room to an event in any status. See [room-assignment.md](room-assignment.md). |
| Puljefordeling | The solver loads events by `is_in_pulje = 1` and `is_in_puljefordeling = 1` (`loadEligibleEvents` in `service/puljefordeling/emulate.go`), with no status filter. |
| Mine arrangementer on the profile | The card links to `/event/{id}` for `Innsendt`, `Godkjent` and `Annonsert`, and to the form `/profile/new/{id}` for `Kladd` and `Forkastet` (`components/profile/my_events.templ`). |
