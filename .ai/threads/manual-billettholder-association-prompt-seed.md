# Prompt Seed: Manual Billettholder Email Should Create User Association

## Ticket Context

We have a bug where manually adding an email to an existing billettholder does not create the durable association row in `billettholdere_users`.

The application should not rely on email joins when querying billettholder data for a user. Email addresses are only a discovery/reconciliation mechanism. The durable relationship should be represented by:

- `billettholdere_users.billettholder_id`
- `billettholdere_users.user_id`

User-facing billettholder reads already depend on this association table. If the association is missing, the user can have a matching email in `billettholder_emails` but still not see or access the billettholder through code paths that query by `billettholdere_users`.

## Intended Behavior

When an email is manually added to an existing billettholder, and there is an existing `users` row with that same email address, the system should insert the matching `(billettholder_id, user_id)` pair into `billettholdere_users`.

The operation should be idempotent:

- adding the same association more than once should not create duplicates
- existing associations should be left intact
- email matching should preserve the current case-insensitive behavior used by `AssociateUserWithBillettholder`

The first implementation step should be to write a test that reproduces the current bug.

## Relevant Existing Code

`AssociateUserWithBillettholder`

- File: `service/checkIn/assign.go`
- Function: `AssociateUserWithBillettholder(userID string, db *sql.DB, logger *slog.Logger) error`
- Current responsibility: finds a user by Descope/user ID, looks up matching `billettholder_emails` rows by the user's email, and inserts rows into `billettholdere_users`.
- Current limitation: this function is only called from the self-service ticket flow, not from the manual email-add flows.

Self-service ticket flow

- File: `pages/profile/tickets/tickets_page.templ`
- Function: `getTicketsRouter`
- Current behavior:
  - fetches tickets from CheckIn
  - calls `checkIn.ConvertTicketToBillettholder(...)`
  - then calls `checkIn.AssociateUserWithBillettholder(user.Id, db, logger)`

Admin manual email-add flow

- File: `pages/admin/billettholder_admin/billettholder_card.templ`
- Function: `addEmailToBilettholderRoute`
- Current behavior:
  - validates the new email
  - checks for duplicate `billettholder_emails` rows
  - inserts `INSERT INTO billettholder_emails (billettholder_id, email, kind) VALUES (?, ?, 'Manual')`
- Current bug:
  - does not create or reconcile the matching `billettholdere_users` row

Profile manual email-add flow

- File: `pages/profile/tickets/billettholder_profile_card.templ`
- Function: `addEmailToBilettholderRoute`
- Current behavior:
  - same broad pattern as the admin route
  - inserts a `Manual` row into `billettholder_emails`
- Related issue:
  - also does not reconcile `billettholdere_users`
  - this may be handled in a separate ticket depending on ticket scope

User billettholder lookup

- File: `service/billettholder/billettholder.go`
- Function: `GetBilettholdere(userId string, db *sql.DB) ([]models.Billettholder, error)`
- Current behavior:
  - when `userId` is provided, joins through `billettholdere_users`
  - does not rely on `billettholder_emails` for the durable user-to-billettholder relationship

Database tables

- `billettholdere`
- `billettholder_emails`
- `users`
- `billettholdere_users`

Schema reference:

- `initialize.sql`
- `schema.sql`

Existing tests to inspect:

- `service/checkIn/assign_users_test.go`
- `service/checkIn/assign_users_generated_test.go`
- `service/checkIn/assign_billettholder_test.go`

Test helpers:

- `service/testdb.go`
- `testutil/createTmpDbLogger.go`

## Suggested First Test

Write a failing test before changing the implementation.

The test should arrange:

- a user in `users`, for example `user_id = "test-user"` and `email = "participant@example.com"`
- an existing billettholder in `billettholdere`
- no existing row in `billettholdere_users`

Then exercise the manual email-add behavior, or a small extracted helper if the route is too awkward to test directly:

- add `participant@example.com` to `billettholder_emails` as `kind = 'Manual'`
- reconcile the association

Assert:

- `billettholdere_users` contains exactly the expected `(billettholder_id, user_id)` pair
- running the association path again does not duplicate the row

## Likely Implementation Direction

Avoid duplicating association SQL inside the HTTP handlers.

Prefer extracting or reusing a service-level helper that can be called after successful manual email insert. Possible shapes:

- reuse `AssociateUserWithBillettholder(...)` if the handler can identify the matching user ID
- add a narrower helper that associates by email after a manual email insert, for example `AssociateUsersWithBillettholderEmail(billettholderID int, email string, db *sql.DB, logger *slog.Logger) error`

The narrower helper may fit this ticket better because the manual email-add route naturally has:

- `billettholderID`
- `newEmailAddress`

It does not naturally have the target user's Descope `user_id`.

The helper should:

- find matching rows in `users` by email, case-insensitively
- insert `(billettholder_id, user_id)` into `billettholdere_users`
- use `INSERT OR IGNORE` or equivalent idempotent behavior
- return useful errors with enough context for logs

## Related But Separate Concerns

Reverse flow:

- When a manual email is removed from a billettholder, the system should remove the corresponding `billettholdere_users` association if that user is no longer associated with the billettholder by any remaining email.
- Existing delete handlers already attempt some cleanup, but this should likely be handled in its own ticket to avoid expanding this change too far.
- Relevant files:
  - `pages/admin/billettholder_admin/billettholder_card.templ`
  - `pages/profile/tickets/billettholder_profile_card.templ`

CheckIn backfill:

- New ticket purchases happen in the external CheckIn ticketing system.
- We do not control that flow and do not want to call the CheckIn API too often.
- Backfill/reconciliation on `/profile/tickets` may be used to avoid stressing the CheckIn API.
- That is related but separate from the manual email-add association bug.

Possible future architecture:

- We may need a dedicated table for trusted user-email associations.
- This would separate raw billettholder emails from verified or user-owned emails.
- `billettholdere_users` should remain the durable authorization/query relation between users and billettholders.

## Prompt For Implementation Session

Investigate and fix the bug where manually adding an email to an existing billettholder does not create the expected row in `billettholdere_users`.

Start by writing a failing test that reproduces the bug. The test should prove that when a manual email is added to a billettholder and a user with that email already exists, the system creates the durable `(billettholder_id, user_id)` association.

Keep the implementation focused on the manual email-add flow for this ticket. Prefer a service-level helper over duplicating SQL in handlers. After the fix, run the targeted Go tests for the affected package(s) and report any remaining risks or follow-up tickets.
