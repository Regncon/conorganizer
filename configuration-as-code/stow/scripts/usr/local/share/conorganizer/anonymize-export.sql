.bail on

PRAGMA foreign_keys = OFF;

BEGIN IMMEDIATE;

CREATE TEMP TABLE _admin_billettholdere (
    billettholder_id INTEGER PRIMARY KEY
);

CREATE TEMP TABLE _email_anonymize_keep (
    email_key TEXT PRIMARY KEY
);

CREATE TEMP TABLE _email_anonymize_map (
    email_key TEXT PRIMARY KEY,
    anon_email TEXT NOT NULL
);

INSERT INTO _admin_billettholdere (billettholder_id)
SELECT DISTINCT r.billettholder_id
FROM relation_billettholdere_users r
JOIN users u ON u.id = r.user_id
WHERE u.is_admin = 1
ON CONFLICT(billettholder_id) DO NOTHING;

INSERT INTO _email_anonymize_keep (email_key)
SELECT lower(trim(email))
FROM users
WHERE is_admin = 1
  AND email IS NOT NULL
  AND trim(email) <> ''
ON CONFLICT(email_key) DO NOTHING;

INSERT INTO _email_anonymize_keep (email_key)
SELECT lower(trim(e.email))
FROM relation_billettholder_emails e
JOIN _admin_billettholdere a ON a.billettholder_id = e.billettholder_id
WHERE e.email IS NOT NULL
  AND trim(e.email) <> ''
ON CONFLICT(email_key) DO NOTHING;

INSERT INTO _email_anonymize_keep (email_key)
SELECT lower(trim(e.email))
FROM events e
JOIN users u ON u.id = e.user_id
WHERE u.is_admin = 1
  AND e.email IS NOT NULL
  AND trim(e.email) <> ''
ON CONFLICT(email_key) DO NOTHING;

INSERT INTO _email_anonymize_map (email_key, anon_email)
WITH all_emails AS (
    SELECT email FROM users
    UNION
    SELECT email FROM relation_billettholder_emails
    UNION
    SELECT email FROM events
),
dedup AS (
    SELECT lower(trim(email)) AS email_key
    FROM all_emails
    WHERE email IS NOT NULL
      AND trim(email) <> ''
      AND lower(trim(email)) NOT IN (
          SELECT email_key FROM _email_anonymize_keep
      )
    GROUP BY lower(trim(email))
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
    WHERE m.email_key = lower(trim(users.email))
)
WHERE lower(trim(email)) IN (
    SELECT email_key FROM _email_anonymize_map
);

UPDATE relation_billettholder_emails
SET email = (
    SELECT anon_email
    FROM _email_anonymize_map m
    WHERE m.email_key = lower(trim(relation_billettholder_emails.email))
)
WHERE lower(trim(email)) IN (
    SELECT email_key FROM _email_anonymize_map
);

UPDATE events
SET email = (
    SELECT anon_email
    FROM _email_anonymize_map m
    WHERE m.email_key = lower(trim(events.email))
)
WHERE lower(trim(email)) IN (
    SELECT email_key FROM _email_anonymize_map
);

UPDATE billettholdere
SET
    first_name = 'User',
    last_name = printf('%06d', ((id * 2654435761) % 900000) + 100000)
WHERE id NOT IN (
    SELECT billettholder_id FROM _admin_billettholdere
);

UPDATE events
SET host_name = CASE
    WHEN user_id IS NOT NULL THEN 'Host ' || printf('%06d', ((user_id * 2654435761) % 900000) + 100000)
    ELSE 'Host'
END;

UPDATE events
SET phone_number = '00000000';

COMMIT;

PRAGMA foreign_keys = ON;
