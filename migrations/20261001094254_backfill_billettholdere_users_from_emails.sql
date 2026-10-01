-- +goose Up
-- Link every existing user to every billettholder that carries their email.
-- Users were only linked when they pressed "Hent billetter" themselves or when
-- a manual email was added after they already had a user. Everyone else could
-- see billettholdere in the email-based picker but was refused when saving
-- interest. Login, "Hent billetter", ticket conversion and "Legg til epost"
-- now create these links, and every reader uses them, so this repairs the
-- users who were missed before.
INSERT OR IGNORE INTO relation_billettholdere_users (billettholder_id, user_id)
SELECT DISTINCT e.billettholder_id, u.id
FROM relation_billettholder_emails AS e
JOIN users AS u ON u.email = e.email COLLATE NOCASE;

-- +goose Down
-- Intentionally empty: links that existed before cannot be told apart from the
-- ones added here, and removing them would lock users out again.
SELECT 1;
