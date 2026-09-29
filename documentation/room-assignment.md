# Romadministrasjon

Romadministrasjon has two admin pages, both server-rendered with templ and driven by Datastar signals:

- `/admin/rooms` lists rooms per floor and lets admins create, edit and delete them (`pages/admin/rooms/rooms_page.templ`).
- `/admin/rooms/assignment/{pulje}` assigns events to rooms for one pulje (`pages/admin/rooms/rooms_assignment_page.templ`).

The API routes live under `/admin/rooms/api` in `pages/admin/admin.go`. The manual test checklist is in [testing/admin-rooms.md](testing/admin-rooms.md).

## Room model

`models.Room` has only `ID`, `Name`, `RoomNumber`, `Floor` and `Notes`. The per-pulje snapshot used by the assignment page is `models.RoomByPulje`, which adds `AssignedEventsID []RoomEventPuljeSummary`.

A room is linked to an event in a pulje through `relation_event_puljer.room_id`. An assignment only counts while the relation has `is_in_pulje = 1`.

### Legacy columns

`rooms.max_concurrent_games` and `rooms.is_disabled` are legacy columns. The application has no capacity control and no way to disable a room, and no application code reads either column.

- Both are still in `schema.sql` and `initialize.sql`. `max_concurrent_games` is `NOT NULL` with no default, and `is_disabled` is `NOT NULL DEFAULT 0`.
- The views `v_events_by_pulje_active` and `v_event_puljer_active` still select them as `room_max_concurrent_games` and `room_is_disabled`.
- `rooms.CreateRoom` (`service/rooms/rooms.go`) writes `0` to both.
- Test fixtures that insert rooms with raw SQL must set `max_concurrent_games` explicitly, because the column is `NOT NULL` (see `pages/admin/rooms_assignment_route_test.go`).

## Room maps

Room maps are matched to rooms by **room number only**. No database column stores a map path. `service/rooms/maps.go` holds two allowlists:

| Allowlist | Key | Value |
| --- | --- | --- |
| `roomMaps` | room number | per-room wayfinding SVG, `/static/rooms/terminus-<floor>-etasje-<number>.svg` |
| `floorMaps` | floor | floor overview SVG: `0` → `/static/rooms/terminus-0-etasje.svg`, `7` → `/static/rooms/terminus-7-etasje.svg` |

`roomMaps` covers `001`, `002`, `003` and `004` on floor 0, and `705`, `706`, `707` and `709`–`714` on floor 7. There is no `708`.

Use the helpers, and never build a path from user input:

- `MapPathForRoom(roomNumber string) (string, bool)`
- `MapPathForFloor(floor int) (string, bool)`

Both return `("", false)` for an unknown key. They are used by the admin room list, the assignment page, the event details page (`pages/event/event_rooms.go`) and the print-friendly page (`pages/print-friendly/print-friendly-page.templ`). On `/admin/rooms`, a room without a map shows "Kart er ikke tilgjengelig for dette rommet." instead of a broken image or link.

`service/rooms/maps_test.go` checks that every allowlisted path points to an existing SVG in `static/rooms/`:

- Floor 7 per-room maps must contain `id="room-<number>-destination-fill"`.
- The other per-room maps must contain `data-room-number="<number>"`.
- Floor overview maps must contain `id="room-targets"`.

To add a map, put the SVG in `static/rooms/` and add the entry to the matching allowlist.

## The `room-map` web component

`static/web_components/room_map.js` defines `<room-map>`. The assignment page renders one for each floor that has an entry in `floorMaps`. The SVG supplies geometry, and the server-rendered room cards (slotted in from light DOM) own the content and the Datastar actions.

Attributes:

- `map-path`: the floor overview SVG to fetch. The component parses it and inserts it with `aria-hidden="true"`.
- `rooms`: JSON of the floor's `[]RoomByPulje`. The component reads `ID` and `RoomNumber`.

For each room, the component looks for SVG geometry by room number:

1. It prefers an element with `data-room-target="<number>"`, which is a `<rect>` read from its `x`, `y`, `width` and `height`.
2. If there is none, it falls back to an element with `data-room-number="<number>"`, for example `<g id="room-705" data-room-number="705">`, measured with `getBBox()`.

Each matched room becomes an overlay positioned in percentages of the SVG `viewBox`, containing `<slot name="room-<ID>">`. A room with no target, or whose bounds can't be read (the measurement throws or has zero size), is listed under "Rom utenfor kartet" below the map. If the map fails to load, the component shows "Kartet kunne ikke lastes. Bruk romlisten nedenfor, eller last siden på nytt." and lists all its slotted rooms under "Rom".

Room coordinates live only in the SVG files, not in Go or JS.

The server sets `slot="room-<ID>"` on a room card only when **both** `MapPathForRoom` and `MapPathForFloor` succeed for it. The component's shadow root has no default slot, so a room on a mapped floor that has no `roomMaps` entry (for example `716`, which has a `data-room-number` group in the floor 7 SVG) is rendered in the HTML but not shown by the component. Rooms on floors without a floor map are rendered in the ordinary `.rooms-container` list instead.

## Room number vs. database ID

The room number (`rooms.room_number`, used in the SVGs) and the room's database ID (`rooms.id`) are different identifiers:

- The SVG is matched by `RoomNumber` only.
- Every action uses the database ID:
  - the room's "+" button sets `$room = <room.ID>` (and `$_roomNumber` for the dialog title),
  - drop handlers post to `/admin/rooms/api/assignment/{pulje}/{event}/{room.ID}`,
  - map slots are named `room-<ID>`,
  - the remove button (×) calls `DELETE /admin/rooms/api/assignment/{pulje}/{event}/{room.ID}`.
- The room edit button on `/admin/rooms` also targets `/admin/rooms/api/<ID>`.

## Assigning events

Admins can drag an event card onto a room or use the room's "+" button. Both post to the same endpoint.

### Event picker

The "+" button opens `#assignment-dialog`, which lists **every** event with status approved or announced, whatever its pulje. Each option is labelled:

- "Ikke i denne puljen": the event has no active relation in this pulje,
- "Mangler rom": it is in the pulje without a room,
- "Har rom – flyttes ved valg": it already has a room in this pulje and will be moved.

When there are no such events, the dialog shows "Ingen godkjente arrangementer.".

### `POST /admin/rooms/api/assignment/{pulje}/{event}/{room}`

1. It validates the pulje and parses the room ID. It returns 400 if either is invalid or the room doesn't exist.
2. It checks that the event exists and has status approved or announced. It returns 409 otherwise.
3. It upserts `relation_event_puljer`:

   ```sql
   INSERT INTO relation_event_puljer (event_id, pulje_id, is_in_pulje, room_id)
   VALUES (?, ?, 1, ?)
   ON CONFLICT(event_id, pulje_id) DO UPDATE SET
       is_in_pulje = 1,
       room_id = excluded.room_id
   ```

   This creates or reactivates the event's relation in the pulje and sets the room. Relations in other puljer are not touched.
4. It broadcasts the `rooms` and `events` live buckets.
5. It responds with a script that closes `#assignment-dialog` and dispatches `room-saved` on `window`.

The endpoint does **not** change the legacy `relation_event_puljer.is_published` flag. A new row keeps the database default `0`, and an existing row keeps its value (checked in `pages/admin/rooms_assignment_route_test.go`). The dialog text says events from other puljer are "lagt til og publisert i denne puljen", but what the endpoint actually sets is `is_in_pulje = 1`. Whether the program is published is a separate global flag in `program_publishing_state`.

### `DELETE /admin/rooms/api/assignment/{pulje}/{event}/{room}`

This sets `room_id = NULL` only when the relation is active and still points to that room. A stale card therefore can't undo another admin's move, and gets a 409 instead. An invalid pulje or a room ID that is not a positive number gets a 400. On success it broadcasts the same buckets and dispatches `room-saved`.

### Page signals

`#room-assignment-page` sets up `draggedEventId`, `dragOverRoom`, `room`, `_roomNumber`, `_roomSaving` (the request indicator that disables the buttons) and `_roomAssignmentStatus`. The `room-saved` window listener clears the drag state, resets `_roomSaving` and shows "Romtildelingen er lagret.".

## Live updates and the dialog

The live endpoint `GET /admin/rooms/api/assignment/{pulje}` streams `RoomsAssignmentPageContent`, which renders only the `#room-assignment` section. `#assignment-dialog` is rendered once by `RoomsAssignmentPage`, **outside** that live fragment, so a live patch can't replace an open modal. `TestRoomsAssignmentPage_DialogRendersOutsideLiveRegion` enforces this. See [live-update-lifecycle.md](live-update-lifecycle.md) for the general rules.

## Event notes

Event notes (`events.notes`) appear on every mini card on the assignment page: unassigned events, events in mapped rooms and events in rooms outside the map. They also appear on each event option in the dialog. All of these use the shared `eventNotesAccordion` component:

```html
<details class="event-notes" data-preserve-attr="open" draggable="false">
  <summary>Notater</summary>
  <p>…</p>
</details>
```

- It is collapsed by default.
- It is not rendered when the notes are blank after trimming.
- It keeps line breaks (`white-space: pre-wrap`), and long notes are capped at `max-block-size: 8rem` with vertical scroll.
- `data-preserve-attr="open"` stops live Datastar morphs from collapsing an accordion an admin has opened.

Room notes (`rooms.notes`) are separate. On the assignment page they are shown only on room cards that are not on a map.

## Data queries

The cards get their data from two queries. Both select event notes, so there is no lookup per card.

- **Assigned cards:** `rooms.GetAllRoomStatusesByPulje` (`service/rooms/rooms.go`). This is a `puljer × rooms` cross join, left-joined to active `relation_event_puljer` rows and `events`. It fills `RoomEventPuljeSummary.Notes` from `e.notes`, and the page picks the selected pulje from the result.
- **Unassigned cards and dialog options:** the page-local `getApprovedEventsForPulje` in `rooms_assignment_page.templ`. It selects approved and announced events with `COALESCE(e.notes, '')` into `RoomEventPuljeSummaryJson.Notes`, left-joined to the active relation in the selected pulje. The list is loaded once per render and passed to both the live content (which keeps events in the pulje with no room) and the dialog.
