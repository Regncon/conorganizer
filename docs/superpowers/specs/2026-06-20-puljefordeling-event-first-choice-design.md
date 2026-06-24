# Puljefordeling Event First-Choice Design

Date: 2026-06-20

## Goal

Complete first-choice handling in the event puljefordeling UI by separating player assignment from first-choice interest changes, deriving first-choice status from existing assignment and interest data, and verifying that successful admin mutations broadcast interest updates.

## Current Context

The event UI is rendered from `components/formsubmission/who_is_interested.templ`. It currently mixes presentation, interest queries, assignment queries, and first-choice SQL. The existing `AddPlayersFirstChoice` helper couples two writes: it sets interest to `Veldig interessert` and assigns the billettholder as a player.

The database already has the necessary source data:

- `relation_events_players` stores event/pulje assignments and role (`Player` or `GM`).
- `interests` stores interest level per billettholder/event/pulje.

No new persisted first-choice column is needed.

## First-Choice Rule

A qualifying first-choice row exists only when:

- `relation_events_players.role = 'Player'`
- a matching `interests` row exists for the same `billettholder_id`, `event_id`, and `pulje_id`
- the matching interest level is `Veldig interessert`

GM rows never count as first-choice evidence for that event/pulje, even if the billettholder has `Veldig interessert` interest for the same event. A billettholder can still be GM for one event and have first-choice as a player on another event.

First-choice status for a row is:

- `HasCurrentPuljeFirstChoice`: this exact billettholder/event/pulje row is a qualifying first-choice row.
- `HasOtherPuljeFirstChoice`: the billettholder has any qualifying first-choice row outside the exact current event/pulje row. The name was kept from the initial design, but the rule is festival-wide: a same-pulje, different-event first-choice also counts as already used.

The UI should show already-used first-choice status before current-pulje first-choice status. Because the conflict can be in another pulje or another event in the same pulje, the label should be generic rather than saying "previous pulje".

The service should use persisted assignment and interest rows as the source of truth and should not depend on pulje chronology.

## Architecture

Create a small first-choice service in `service/puljefordeling/first_choice.go`.

Public API:

```go
type FirstChoiceKey struct {
	BillettholderID int
	EventID         string
	PuljeID         string
}

type FirstChoiceStatus struct {
	HasCurrentPuljeFirstChoice bool
	HasOtherPuljeFirstChoice   bool
}

func GetFirstChoiceStatusesForEvent(db *sql.DB, eventID string) (map[FirstChoiceKey]FirstChoiceStatus, error)
```

The function should batch-load statuses for the event instead of querying per row. It should compare the event's interest and assignment rows against all qualifying first-choice rows in the festival.

`components/formsubmission/who_is_interested.templ` should keep rendering and row assembly responsibilities, but should consume this service result instead of embedding the full first-choice rule in component SQL. The existing `queryFirstChoice` SQL should be removed or replaced by service-derived status.

This design intentionally does not move all interest and assignment loading into `service/puljefordeling`; that refactor is larger than the current feature needs.

## Mutation Design

Assignment mutations and first-choice mutations must be independent.

Assignment operations write only `relation_events_players`:

- assign as player
- assign as GM
- remove assignment
- switch between player and GM

First-choice operations write only `interests`:

- set first-choice: upsert/update the matching interest to `Veldig interessert`
- remove first-choice: update the matching interest to `Middels interessert`

The existing combined `AddPlayersFirstChoice` behavior should be split or replaced so assigning a player does not automatically set first-choice.

## UI Behavior

In `who_is_interested.templ`:

- keep player and GM controls as assignment controls
- add a separate first-choice control for assigned player rows
- show `Set førstevalg` when the current row does not have current-pulje first-choice
- show `Fjern førstevalg` when the current row has current-pulje first-choice
- disable the first-choice control for GM rows
- disable setting first-choice when `HasOtherPuljeFirstChoice` is true, because first-choice should only happen once per festival
- replace the search action `Legg til som førsteval` with `Legg til som spelar`; first-choice can then be set from the assigned row

Interest rows that are not assigned as a player do not count as first-choice, even when their interest level is `Veldig interessert`.

## Routes And Broadcasts

Keep the existing assignment route for assignment state:

- `PUT /admin/approval/api/event-players/update_status`

Add a separate endpoint for first-choice interest state:

- `PUT /admin/approval/api/event-players/first-choice`

The request should include the existing assignment signals:

- `assignmentEventId`
- `assignmentPuljeId`
- `assignmentBillettholderId`
- desired first-choice boolean

Successful assignment and first-choice mutations must broadcast `live.BucketInterests`.

## Error Handling

Mutation handlers should validate billettholder ID, event ID, and pulje ID before writing. Invalid input should return a client error. Database write failures should return a server error and should not broadcast a successful update.

If the database write succeeds but broadcasting fails, the route should return an error so the UI does not silently miss the live refresh failure.

## Testing Plan

Service tests in `service/puljefordeling/first_choice_test.go`:

- player assignment plus `Veldig interessert` in the current pulje gives current first-choice
- GM assignment plus `Veldig interessert` does not count
- GM in one pulje does not block player first-choice in another pulje
- player first-choice in another pulje sets `HasOtherPuljeFirstChoice`
- player first-choice in another event in the same pulje also sets `HasOtherPuljeFirstChoice`
- other-pulje first-choice takes display precedence over current row
- assigned player with `Middels interessert` does not count as first-choice

Mutation tests:

- assigning player does not change interest level
- setting first-choice changes interest to `Veldig interessert` without changing assignment
- removing first-choice changes interest to `Middels interessert` without removing assignment
- first-choice mutation rejects or performs no first-choice change for GM rows

Route or thin handler tests:

- assignment mutation broadcasts `live.BucketInterests`
- first-choice mutation broadcasts `live.BucketInterests`
- failed database mutation does not report a successful broadcast
- broadcast failure returns an error

The live manager already has coverage for NATS bucket behavior. Route tests should verify that the admin mutation path calls the interests broadcast after successful mutations rather than re-testing browser/SSE behavior.

## Non-Goals

- Add a persisted `first_choice` column.
- Refactor all event interest and assignment loading out of the templ component.
- Add chronology-based first-choice logic.
- Re-test the full Datastar/SSE live stream in event-player route tests.
