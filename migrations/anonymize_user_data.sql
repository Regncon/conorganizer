-- +goose Up
PRAGMA foreign_keys = OFF;

DROP TABLE IF EXISTS _email_anonymize_map;
CREATE TABLE _email_anonymize_map (
    email_key TEXT PRIMARY KEY,
    anon_email TEXT NOT NULL
);

INSERT INTO _email_anonymize_map (email_key, anon_email)
WITH all_emails AS (
    SELECT email FROM users
    UNION
    SELECT email FROM relation_billettholder_emails
    UNION
    SELECT email FROM events
),
dedup AS (
    SELECT lower(email) AS email_key
    FROM all_emails
    WHERE email IS NOT NULL AND email <> ''
    GROUP BY lower(email)
),
numbered AS (
    SELECT
        email_key,
        row_number() OVER (ORDER BY email_key) AS rn
    FROM dedup
)
SELECT
    email_key,
    printf('user_%05d@example.invalid', rn) AS anon_email
FROM numbered;

UPDATE users
SET email = (
    SELECT anon_email
    FROM _email_anonymize_map m
    WHERE m.email_key = lower(users.email)
);

UPDATE relation_billettholder_emails
SET email = (
    SELECT anon_email
    FROM _email_anonymize_map m
    WHERE m.email_key = lower(relation_billettholder_emails.email)
);

UPDATE events
SET email = (
    SELECT anon_email
    FROM _email_anonymize_map m
    WHERE m.email_key = lower(events.email)
);

UPDATE billettholdere
SET
    first_name = 'User',
    last_name = printf('%06d', id);

UPDATE events
SET host_name = CASE
    WHEN user_id IS NOT NULL THEN 'Host ' || printf('%06d', user_id)
    ELSE 'Host'
END;

UPDATE events
SET phone_number = '00000000';

DROP TABLE IF EXISTS _email_anonymize_map;

PRAGMA foreign_keys = ON;


-- +goose Down
-- Irreversible: original PII is not recoverable after anonymization.
SELECT 'irreversible';
