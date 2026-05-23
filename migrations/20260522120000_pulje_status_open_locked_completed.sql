-- +goose Up
PRAGMA foreign_keys = OFF;
PRAGMA legacy_alter_table = ON;

INSERT INTO pulje_statuses(status)
VALUES
    ('open'),
    ('locked'),
    ('completed')
ON CONFLICT(status) DO NOTHING;

CREATE TABLE puljer_new(
    id TEXT NOT NULL PRIMARY KEY,
    name TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'open' CHECK(status IN('open', 'locked', 'completed')),
    start_at TEXT NOT NULL,
    end_at TEXT NOT NULL,
    FOREIGN KEY(status) REFERENCES pulje_statuses(status) ON UPDATE CASCADE
) STRICT;

INSERT INTO puljer_new(id, name, status, start_at, end_at)
SELECT
    id,
    name,
    CASE status
        WHEN 'locked' THEN 'locked'
        WHEN 'completed' THEN 'completed'
        ELSE 'open'
    END,
    start_at,
    end_at
FROM puljer;

DROP TABLE puljer;
ALTER TABLE puljer_new RENAME TO puljer;

DELETE FROM pulje_statuses
WHERE status IN('not_published', 'published');

PRAGMA foreign_keys = ON;
PRAGMA legacy_alter_table = OFF;

-- +goose Down
PRAGMA foreign_keys = OFF;
PRAGMA legacy_alter_table = ON;

INSERT INTO pulje_statuses(status)
VALUES
    ('not_published'),
    ('published'),
    ('locked'),
    ('completed')
ON CONFLICT(status) DO NOTHING;

CREATE TABLE puljer_new(
    id TEXT NOT NULL PRIMARY KEY,
    name TEXT NOT NULL,
    status TEXT NOT NULL CHECK(status IN('not_published', 'published', 'locked', 'completed')),
    start_at TEXT NOT NULL,
    end_at TEXT NOT NULL,
    FOREIGN KEY(status) REFERENCES pulje_statuses(status) ON UPDATE CASCADE
) STRICT;

INSERT INTO puljer_new(id, name, status, start_at, end_at)
SELECT
    id,
    name,
    CASE status
        WHEN 'locked' THEN 'locked'
        WHEN 'completed' THEN 'completed'
        ELSE 'published'
    END,
    start_at,
    end_at
FROM puljer;

DROP TABLE puljer;
ALTER TABLE puljer_new RENAME TO puljer;

DELETE FROM pulje_statuses
WHERE status = 'open';

PRAGMA foreign_keys = ON;
PRAGMA legacy_alter_table = OFF;
