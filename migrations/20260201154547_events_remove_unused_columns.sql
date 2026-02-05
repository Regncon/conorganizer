-- +goose Up
PRAGMA foreign_keys = OFF;

CREATE TABLE IF NOT EXISTS events_new (
    id TEXT PRIMARY KEY NOT NULL DEFAULT ( lower(hex(randomblob(8))) ),
    title TEXT NOT NULL,
    intro TEXT NOT NULL,
    description TEXT NOT NULL,
    system TEXT DEFAULT '',
    event_type TEXT NOT NULL DEFAULT 'Other',
    age_group TEXT NOT NULL DEFAULT 'Default',
    event_runtime TEXT NOT NULL DEFAULT 'Normal',
    host_name TEXT NOT NULL,
    host INTEGER,
    email TEXT NOT NULL,
    phone_number TEXT NOT NULL,
    max_players INTEGER NOT NULL,
    beginner_friendly BOOLEAN NOT NULL,
    can_be_run_in_english BOOLEAN NOT NULL,
    notes TEXT DEFAULT '',
    status TEXT NOT NULL DEFAULT 'Kladd',
    inserted_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (host) REFERENCES users(id) ON DELETE SET NULL,
    FOREIGN KEY (status) REFERENCES event_statuses(status) ON UPDATE CASCADE,
    FOREIGN KEY (event_type) REFERENCES events_types(event_type) ON UPDATE CASCADE,
    FOREIGN KEY (age_group) REFERENCES age_groups(age_group) ON UPDATE CASCADE,
    FOREIGN KEY (event_runtime) REFERENCES event_runtimes(runtime) ON UPDATE CASCADE
);

INSERT INTO events_new (
    id,
    title,
    intro,
    description,
    system,
    event_type,
    age_group,
    event_runtime,
    host_name,
    host,
    email,
    phone_number,
    max_players,
    beginner_friendly,
    can_be_run_in_english,
    notes,
    status,
    inserted_time
)
SELECT
    id,
    title,
    intro,
    description,
    system,
    event_type,
    age_group,
    event_runtime,
    host_name,
    host,
    email,
    phone_number,
    max_players,
    beginner_friendly,
    can_be_run_in_english,
    notes,
    status,
    inserted_time
FROM events;

DROP TABLE events;
ALTER TABLE events_new RENAME TO events;

PRAGMA foreign_keys = ON;


-- +goose Down
PRAGMA foreign_keys = OFF;

CREATE TABLE IF NOT EXISTS events_old (
    id TEXT PRIMARY KEY NOT NULL DEFAULT ( lower(hex(randomblob(8))) ),
    title TEXT NOT NULL,
    intro TEXT NOT NULL,
    description TEXT NOT NULL,
    image_url TEXT DEFAULT '',
    system TEXT DEFAULT '',
    event_type TEXT NOT NULL DEFAULT 'Other',
    age_group TEXT NOT NULL DEFAULT 'Default',
    event_runtime TEXT NOT NULL DEFAULT 'Normal',
    host_name TEXT NOT NULL,
    host INTEGER,
    email TEXT NOT NULL,
    phone_number TEXT NOT NULL,
    pulje_name INTEGER,
    max_players INTEGER NOT NULL,
    beginner_friendly BOOLEAN NOT NULL,
    can_be_run_in_english BOOLEAN NOT NULL,
    notes TEXT DEFAULT '',
    status TEXT NOT NULL DEFAULT 'Kladd',
    inserted_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (host) REFERENCES users(id) ON DELETE SET NULL,
    FOREIGN KEY (pulje_name) REFERENCES puljer(name) ON UPDATE CASCADE,
    FOREIGN KEY (status) REFERENCES event_statuses(status) ON UPDATE CASCADE,
    FOREIGN KEY (event_type) REFERENCES events_types(event_type) ON UPDATE CASCADE,
    FOREIGN KEY (age_group) REFERENCES age_groups(age_group) ON UPDATE CASCADE,
    FOREIGN KEY (event_runtime) REFERENCES event_runtimes(runtime) ON UPDATE CASCADE
);

INSERT INTO events_old (
    id,
    title,
    intro,
    description,
    image_url,
    system,
    event_type,
    age_group,
    event_runtime,
    host_name,
    host,
    email,
    phone_number,
    pulje_name,
    max_players,
    beginner_friendly,
    can_be_run_in_english,
    notes,
    status,
    inserted_time
)
SELECT
    id,
    title,
    intro,
    description,
    '' AS image_url,
    system,
    event_type,
    age_group,
    event_runtime,
    host_name,
    host,
    email,
    phone_number,
    NULL AS pulje_name,
    max_players,
    beginner_friendly,
    can_be_run_in_english,
    notes,
    status,
    inserted_time
FROM events;

DROP TABLE events;
ALTER TABLE events_old RENAME TO events;

PRAGMA foreign_keys = ON;
