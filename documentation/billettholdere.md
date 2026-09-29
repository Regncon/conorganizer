# Billettholdere

A billettholder is a CheckIn ticket that is not of type "Middag" (see
[domeneordbok.md](../domeneordbok.md#billettholder)). This page covers how
billettholdere are imported from CheckIn, how they are linked to logged-in users
through e-post, and how the admin billettholder page reads them.

## Data model

| Table | Purpose |
| --- | --- |
| `billettholdere` | One row per imported CheckIn ticket. `ticket_id` is `UNIQUE`. |
| `relation_billettholder_emails` | E-post addresses on a billettholder. `email` is `COLLATE NOCASE`; `kind` is `Ticket`, `Associated` or `Manual`. |
| `relation_billettholdere_users` | Ownership link between a user and a billettholder. `PRIMARY KEY (billettholder_id, user_id)`, `user_id` references `users.id`. |

- `relation_billettholdere_users` is the durable link. All user-scoped reads use
  it. Do not use `relation_billettholder_emails` for ownership or access; it is
  only used to discover and reconcile links.
- User-scoped reads go through `billettholderService.GetBillettholdere(userId, db)`
  (`service/billettholder/billettholder.go`). `userId` is the auth provider id
  stored in `users.external_id`. When it is non-empty the query joins
  `relation_billettholdere_users` and `users`; when it is empty all
  billettholdere are returned. `GetBillettholdereWithFilters` and
  `GetBillettholderByUserId` use the same join.

E-post kinds (`models.BillettholderEmailKind*`):

| Kind | Created by | Deletable |
| --- | --- | --- |
| `Ticket` | Ticket import: the e-post on the ticket itself. | No |
| `Associated` | Ticket import: other e-post addresses on tickets in the same CheckIn order. | No |
| `Manual` | Admin or the user on Min Side, via "Legg til epost". | Yes |

## CheckIn integration

Code lives in `service/checkIn`.

- `GetTicketsFromCheckIn(ctx, logger, searchTerm)` reads tickets through an
  in-memory cache (`cache.go`) in front of the CheckIn GraphQL API: 5 minute TTL,
  30 second failure cooldown (no repeated upstream calls while CheckIn is down),
  10 second request timeout. When a refresh fails but earlier data exists, it
  returns the stale tickets together with the error and sets
  `TicketFetchResult.UsedStaleCache`. Callers decide whether stale data is good
  enough; both the ticket fetch and admin conversion accept it.
- `TicketTypeMiddag = 251934`. Middag tickets are never imported.
- `billettholdere.is_over_18` comes from the ticket's birth date (`isOver18` in
  `cache.go`): true when the billettholder is 18 or older on 2026-10-04. A
  missing or unparseable birth date counts as under 18. The puljefordeling
  solver does not use the flag, so it can seat a billettholder under 18 on an
  18+ arrangement (`AgeGroupAdultsOnly`). The admin puljefordeling page uses it
  to warn about such seats, and to ask for confirmation before a manual pin. See
  [domeneordbok.md](../domeneordbok.md#aldersgrense-18).
- `ConvertTicketToBillettholder(ctx, ticketId, db, logger)` fetches all tickets and
  calls the unexported `converTicketIdToNewBillettholder`, which:
  - refuses Middag tickets,
  - reuses an existing billettholder with the same `ticket_id`, or inserts one
    (reporting `TicketConversionResult.CreatedBillettholders = 1`),
  - adds the ticket's e-post as `Ticket` and other e-post addresses from the same
    order as `Associated`, skipping pairs that already exist on that billettholder.
- `AssociateTicketsWithEmail` returns the non-Middag tickets whose e-post equals
  the given address, ignoring case (`strings.EqualFold`).
- `AssociateTicketsWithBillettholder(tickets, email, db, logger)` imports every
  non-Middag ticket from each order that contains a ticket registered to `email`,
  so tickets bought for others on the same order are included. It sums created
  billettholdere into `TicketAssociationResult.CreatedBillettholders`. No match is
  not an error.

The ticket purchase link is the constant `checkInTicketsURL` in
`pages/profile/tickets/tickets_page.templ` (currently
`https://event.checkin.no/221572/regncon-xxxiv-2026`). Update it for each new
convention year.

## "Hent billetter" flow

`/profile/tickets` only reads `GetBillettholdere(userId, db)`; loading the page does
not import or link anything. Reconciliation starts when the user presses
"Hent billetter" (`POST /profile/tickets/api/get-tickets`):

1. `checkIn.GetTicketsFromCheckIn(ctx, logger, "")` fetches all tickets once.
2. `AssociateTicketsWithBillettholder(tickets, user.Email, db, logger)` imports the
   user's orders.
3. `AssociateUserWithBillettholder(user.Id, db, logger)` looks up the user by
   `users.external_id`, finds every `relation_billettholder_emails` row matching
   the user's e-post (`COLLATE NOCASE`, any kind), and inserts the pairs with
   `INSERT OR IGNORE` into `relation_billettholdere_users`. It returns the number
   of new links (`RowsAffected()`), `0` when nothing matches, and is idempotent.
   It never imports tickets on its own.

The route patches `getTicketsSuccessMessage`, `getTicketsInfoMessage` and
`getTicketsErrorMessage`:

| Situation | Message |
| --- | --- |
| 1 new link | Success: "1 billett hentet!" |
| N new links | Success: "N billetter hentet!" |
| No new links | Info: "Ingen nye billetter funnet." |
| CheckIn refresh failed, stale cache used | Info: stale-cache warning (replaces "Ingen nye billetter funnet."; a success message can still show) |
| CheckIn unavailable and no cache | Error: CheckIn unavailable message |
| Import or linking failed | Error: "Vi klarte ikke å oppdatere billettene dine akkurat nå. Prøv igjen senere." |

N is only the number of new `relation_billettholdere_users` rows for the current
user in this request. It is not the total ticket count and does not include
billettholdere created for other people on the order.

The route broadcasts `live.BucketBillettholders` only when
`shouldBroadcastTicketFetch` is true: at least one billettholder was created (for
anyone) or at least one new link was created for the user. No previous count is
stored or compared.

The `/profile/tickets` empty state offers both "Hent billetter" and
"Kjøp billetter" (`checkInTicketsURL`). The Min Side overview
(`components/profile/my_tickets.templ`) links only to `/profile/tickets`
("Hent billetter"), because no tickets usually means they have not been fetched
or linked yet, not that none were bought.

## Association paths and helpers

The only ways a user gets linked to, or unlinked from, a billettholder:

| Path | Route | Helper |
| --- | --- | --- |
| Self-service ticket fetch | `POST /profile/tickets/api/get-tickets` | `AssociateTicketsWithBillettholder`, then `AssociateUserWithBillettholder` |
| Add manual e-post (admin) | `POST /admin/billettholder/api/new-email/{id}/` | `AssociateUsersWithBillettholderEmail` |
| Add manual e-post (Min Side) | `POST /profile/tickets/api/new-email/{id}/` | `AssociateUsersWithBillettholderEmail` |
| Delete manual e-post (admin) | `POST /admin/billettholder/api/delete-email/{id}/{emailID}/` | `DisassociateUsersFromBillettholderEmail` |
| Delete manual e-post (Min Side) | `POST /profile/tickets/api/delete-email/{id}/{emailID}/` | `DisassociateUsersFromBillettholderEmail` |

Admin ticket conversion (`GET /admin/billettholder/add/api/convert/{ticketID}/`
→ `checkIn.ConvertTicketToBillettholder`) creates no user links. A billettholder
converted by an admin shows up for the user only after the user presses
"Hent billetter" or an admin adds a matching manual e-post.

Add manual e-post:

- The handler reads the card-scoped signal `newEmail-{id}`, rejects an empty value
  ("Tomt felt for epostadresse") and an address already on that billettholder,
  then inserts the row with kind `Manual`.
- `AssociateUsersWithBillettholderEmail(billettholderID, email, db, logger)` runs
  `INSERT OR IGNORE INTO relation_billettholdere_users ... SELECT ?, id FROM users
  WHERE email = ? COLLATE NOCASE`; the primary key prevents duplicates. It takes
  the billettholder id and e-post, not a user id, because the handler does not
  know which user (if any) owns the address.

Delete manual e-post:

- The handler looks up the e-post row scoped to that billettholder, blocks any
  kind other than `Manual` ("Kun manuelle epostadresser kan slettes"), deletes the
  row, and then calls `DisassociateUsersFromBillettholderEmail`.
- The helper runs a single `DELETE ... AND NOT EXISTS`: it removes links for users
  whose e-post matches the removed address (`COLLATE NOCASE`), but only when no
  remaining e-post on the same billettholder still matches. Deleting one of
  several matching addresses keeps the link. Handlers must not do their own
  unconditional link cleanup, or they would disagree with the helper.

Caveat: add and delete are not transactional. The `relation_billettholder_emails`
insert/delete and the following link change are separate statements. If the
second step fails, the route reports an error but the e-post change stays
committed.

Logging: the helpers in `service/checkIn/assign.go` use a logger with
`component=checkin_assign`. They log at Info only when rows change —
"Created billettholder user associations" (`created_associations`) and
"Removed billettholder user associations" (`removed_associations`), both with
`billettholder_id` and `association_flow=billettholder_email` — and at Debug
otherwise. They never log e-post addresses. Errors are returned and logged at the
route boundary by `handleError`, with `billettholder_id`. Exception:
`AssociateUserWithBillettholder` prints to stdout with `fmt.Printf` when the
insert fails, including the matching e-post rows.

## Admin billettholder page

`/admin/billettholder` (`pages/admin/billettholder_admin`). The live stream at
`/admin/billettholder/api/` subscribes to the `billettholders` and `interests`
buckets.

### Interest data

`BillettholderAdminPage` loads interests and assignments for all shown
billettholdere in one query per render:
`getBillettholderInterestSectionsByBillettholderID(db, billettholderIDs(...))`
(`billettholder_interest_data.go`) returns a map keyed by billettholder id, and
each card gets its `[]billettholderInterestPuljeSection`.

The query is a `UNION ALL` of:

- assigned rows from `relation_events_players`, grouped per
  billettholder/event/pulje with `MAX(role = Player)` / `MAX(role = GM)`, so one
  event with both roles gives one row; `LEFT JOIN interests` supplies the
  interessenivå,
- interest rows from `interests` without a matching `relation_events_players` row.

Both parts join `events` and `puljer`. Rows are ordered by pulje `start_at`, pulje
name, pulje id, then assigned first (by title), then interessenivå
(Veldig, Middels, Litt), then title. The loader does not look at
`relation_event_puljer` (no `is_in_pulje` or `is_published` filter).

### Interest dialog

Each card has an "Interesser (n) / Tildelt (n)" button that opens a native
`<dialog>` rendered inside the card with `showModal()`, so opening it does not
change the card height or the grid layout. Which dialog is open is the page-level
signal `$billettholderInterestOpenDialogId`; see
[datastar-signals.md](datastar-signals.md) for the signal-ownership pattern.

- One tab per pulje, in `start_at` order. Tabs are client-only: the page-wide
  signal `$billettholderInterestActivePulje` (declared with
  `data-signals__ifmissing:billettholder-interest-active-pulje`) drives
  `data-class:selected`, `aria-selected`, `tabindex` and which panel is hidden.
  Switching tabs needs no server round trip. The open button sets the signal to
  that billettholder's first pulje before opening.
- Each pulje panel shows a "Tildelt (n)" group first, with the role label
  ("Spiller", "Spilleder (GM/DM)" or "Spiller og Spilleder (GM/DM)") and the
  interessenivå, then groups per interessenivå: "🤩 Veldig interessert",
  "😁 Interessert" (stored as `Middels interessert`), "🙂 Litt interessert".
- Each row links to `/admin/approval/edit/{eventId}` and closes the dialog on click.

### Counters, badges and filters

The overview counters (Totalt, "Uten førstevalg", "GM/DM") and the filters count
only assignments in active puljer: `relation_events_players` joined to
`relation_event_puljer` with `is_in_pulje = 1` (the `active_assignments` CTE).
They do not use `is_published`. See [domeneordbok.md](../domeneordbok.md#førstevalg)
for what counts as førstevalg.

The per-card badges and the GM icon on dialog tabs (`billettholderHasFirstChoice`,
`billettholderHasGMOrDM`, `billettholderSectionHasGMOrDM`) only look at the
assigned rows from the interest query, which does not check `is_in_pulje`. If an
assignment exists for an event that is no longer in that pulje, a card badge can
differ from the counters and filters.

The "Uten førstevalg" and "GM/DM" filters run on the server:

1. The button flips `$filterWithoutFirstChoice` / `$filterGmOrDm` and calls
   `GET /admin/billettholder/api/search/`.
2. The handler parses the `datastar` query payload into `BillettholderFilters`
   and calls `liveManager.EnsureConnection(w, r, live.BucketBillettholders)`,
   which re-renders only the current connection (see
   [live-update-lifecycle.md](live-update-lifecycle.md)).
3. `BillettholderAdminPage` calls
   `billettholderService.GetBillettholdereWithFilters("", db, filters)`, which
   combines the filters with AND on the `active_assignments` CTE.
   `GetBillettholdere` is unchanged for profile pages; both share
   `scanBillettholdere`.

Gotcha: `searchTerm` and `billettholderFilters` are package-level variables in
`billettholder_admin_page.templ`. The last search/filter request applies to every
admin's next render, not just the requester's. (The search input itself is
currently commented out, so each filter request also resets `searchTerm` to "".)

## Tests

- `service/checkIn/assign_users_test.go`: `AssociateUserWithBillettholder`,
  `AssociateUsersWithBillettholderEmail` and
  `DisassociateUsersFromBillettholderEmail`.
- `service/checkIn/assign_billettholder_test.go` and `assign_ticket_test.go`:
  ticket import and e-post matching, including whole-order import.
- `service/checkIn/cache_test.go`: CheckIn cache, stale fallback and cooldown.
- `pages/profile/tickets/tickets_page_test.go`: "Hent billetter" feedback messages
  and the empty state.
- `pages/admin/billettholder_admin/billettholder_admin_overview_test.go` and
  `service/billettholder/billettholder_admin_filter_test.go`: active-pulje counters
  and filters.
- `pages/admin/billettholder_admin/billettholder_interest_dialog_test.go`: the
  interest query.

E-post route suites:
`pages/admin/billettholder_admin/billettholder_email_routes_test.go` (admin) and
`pages/profile/tickets/billettholder_email_routes_test.go` (Min Side). They mount
the real handlers on a chi router, POST Datastar JSON signal bodies through
`httptest`, and assert on the card-scoped signals (`successMessage-{id}`,
`errorMessage-{id}`, `newEmail-{id}`) and on database rows. Both suites cover the
same cases:

- adding an e-post lands on the requested card only and links a matching user,
- an empty e-post is rejected with a card-scoped error and nothing is stored,
- a duplicate e-post on the same billettholder is rejected without a new row,
- deleting removes only the requested manual e-post,
- deleting removes the user link when no remaining e-post matches the user,
- deleting keeps the link when another matching e-post remains.

Manual test checklists: [testing/profile-tickets.md](testing/profile-tickets.md)
and [testing/admin-billettholders.md](testing/admin-billettholders.md).
