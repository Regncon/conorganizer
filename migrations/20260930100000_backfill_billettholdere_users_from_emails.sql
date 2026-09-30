-- +goose Up
-- Backfill relation_billettholdere_users for secondary/manual billettholder
-- emails that were added before the person they belong to ever logged in.
--
-- pages/login.syncPostLoginUser created/updated a user's `users` row on every
-- login but never called service/checkIn.AssociateUserWithBillettholder, so a
-- user created at login after their email was already present in
-- relation_billettholder_emails never got the matching
-- relation_billettholdere_users row. That left them able to see a
-- billettholder in the "Meld interesse" picker
-- (components/ticket_holder.GetTicketHolders, keyed on
-- relation_billettholder_emails) but blocked from registering interest for
-- it (pages/event.updateInterest's access check, keyed on
-- relation_billettholdere_users) with "does not have access to this
-- billettholder interest". Login now calls AssociateUserWithBillettholder on
-- every login to keep the two tables in sync going forward; this migration
-- repairs users who were already affected.
INSERT OR IGNORE INTO relation_billettholdere_users (billettholder_id, user_id)
SELECT e.billettholder_id, u.id
FROM relation_billettholder_emails e
JOIN users u ON u.email = e.email COLLATE NOCASE;

-- +goose Down
-- Not reversible: there is no way to tell which of these associations
-- already existed before this migration ran versus which rows it just
-- created, and deleting them would re-break access for users who already
-- relied on them. Login (pages/login.syncPostLoginUser) and the "add
-- secondary email" routes recreate the same rows going forward regardless,
-- so leaving them in place on a Down is intentional.
SELECT 1;
