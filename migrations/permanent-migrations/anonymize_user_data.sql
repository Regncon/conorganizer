-- +goose Up
PRAGMA foreign_keys = OFF;

DROP TABLE IF EXISTS _email_anonymize_map;
CREATE TABLE _email_anonymize_map (
    email_key TEXT PRIMARY KEY,
    anon_email TEXT NOT NULL
);

DROP TABLE IF EXISTS _email_anonymize_keep;
CREATE TABLE _email_anonymize_keep (
    email_key TEXT PRIMARY KEY
);

-- Emails listed in ANONYMIZE_PRESERVE_EMAILS are intentionally kept for testing/admin access.
-- Set ANONYMIZE_PRESERVE_EMAILS in .env as a comma-separated list before running the migration.
-- Example: ANONYMIZE_PRESERVE_EMAILS=admin@example.com,test@example.com
-- +goose ENVSUB ON
WITH RECURSIVE keep_emails(value) AS (
    SELECT '${ANONYMIZE_PRESERVE_EMAILS:-}'
),
split(email_key, rest) AS (
    SELECT
        lower(trim(CASE
            WHEN instr(value, ',') = 0 THEN value
            ELSE substr(value, 1, instr(value, ',') - 1)
        END)),
        CASE
            WHEN instr(value, ',') = 0 THEN ''
            ELSE substr(value, instr(value, ',') + 1)
        END
    FROM keep_emails

    UNION ALL

    SELECT
        lower(trim(CASE
            WHEN instr(rest, ',') = 0 THEN rest
            ELSE substr(rest, 1, instr(rest, ',') - 1)
        END)),
        CASE
            WHEN instr(rest, ',') = 0 THEN ''
            ELSE substr(rest, instr(rest, ',') + 1)
        END
    FROM split
    WHERE rest <> ''
)
INSERT INTO _email_anonymize_keep (email_key)
SELECT email_key
FROM split
WHERE email_key <> ''
ON CONFLICT(email_key) DO NOTHING;
-- +goose ENVSUB OFF

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
      AND lower(email) NOT IN (
          SELECT lower(trim(email_key)) FROM _email_anonymize_keep
      )
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
    printf('user_%05d@example.invalid', ((rn * 7919) % 90000) + 10000) AS anon_email
FROM numbered;

UPDATE users
SET email = (
    SELECT anon_email
    FROM _email_anonymize_map m
    WHERE m.email_key = lower(users.email)
)
WHERE lower(email) IN (
    SELECT email_key FROM _email_anonymize_map
);

UPDATE relation_billettholder_emails
SET email = (
    SELECT anon_email
    FROM _email_anonymize_map m
    WHERE m.email_key = lower(relation_billettholder_emails.email)
)
WHERE lower(email) IN (
    SELECT email_key FROM _email_anonymize_map
);

UPDATE events
SET email = (
    SELECT anon_email
    FROM _email_anonymize_map m
    WHERE m.email_key = lower(events.email)
)
WHERE lower(email) IN (
    SELECT email_key FROM _email_anonymize_map
);

UPDATE billettholdere
SET
    first_name = 'User',
    last_name = printf('%06d', ((id * 2654435761) % 900000) + 100000)
WHERE NOT EXISTS (
    SELECT 1
    FROM relation_billettholder_emails e
    JOIN _email_anonymize_keep k ON lower(trim(k.email_key)) = lower(e.email)
    WHERE e.billettholder_id = billettholdere.id
);

UPDATE events
SET host_name = CASE
    WHEN user_id IS NOT NULL THEN 'Host ' || printf('%06d', ((user_id * 2654435761) % 900000) + 100000)
    ELSE 'Host'
END;

UPDATE events
SET phone_number = '00000000';

DROP TABLE IF EXISTS _email_anonymize_map;
DROP TABLE IF EXISTS _email_anonymize_keep;

PRAGMA foreign_keys = ON;


-- +goose Down
-- Irreversible: original PII is not recoverable after anonymization.
SELECT 'irreversible';
