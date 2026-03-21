-- +goose Up
PRAGMA foreign_keys = OFF;

DROP TABLE IF EXISTS _email_anonymize_map;
CREATE TABLE _email_anonymize_map (
    email TEXT PRIMARY KEY,
    anon_email TEXT NOT NULL
);

INSERT INTO _email_anonymize_map (email, anon_email)
WITH all_emails AS (
    SELECT email FROM users
    UNION
    SELECT email FROM billettholder_emails
    UNION
    SELECT email FROM events
),
dedup AS (
    SELECT email
    FROM all_emails
    WHERE email IS NOT NULL AND email <> ''
    GROUP BY email
),
numbered AS (
    SELECT
        email,
        row_number() OVER (ORDER BY email) AS rn
    FROM dedup
)
SELECT
    email,
    printf('user_%05d@example.invalid', rn) AS anon_email
FROM numbered;

UPDATE users
SET email = (
    SELECT anon_email
    FROM _email_anonymize_map m
    WHERE m.email = users.email
);

UPDATE billettholder_emails
SET email = (
    SELECT anon_email
    FROM _email_anonymize_map m
    WHERE m.email = billettholder_emails.email
);

UPDATE events
SET email = (
    SELECT anon_email
    FROM _email_anonymize_map m
    WHERE m.email = events.email
);

UPDATE billettholdere
SET
    first_name = 'User',
    last_name = printf('%06d', id);

UPDATE events
SET host_name = CASE
    WHEN host IS NOT NULL THEN 'Host ' || printf('%06d', host)
    ELSE 'Host'
END;

UPDATE events
SET phone_number = '00000000';

DROP TABLE IF EXISTS _email_anonymize_map;

PRAGMA foreign_keys = ON;


-- +goose Down
-- Irreversible: original PII is not recoverable after anonymization.
SELECT 'irreversible';
