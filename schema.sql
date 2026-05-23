CREATE TABLE
    users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        external_id TEXT NOT NULL UNIQUE,
        email TEXT NOT NULL,
        is_admin INTEGER NOT NULL DEFAULT 0 CHECK (is_admin IN (0, 1)),
        inserted_at TEXT DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))
    ) STRICT;

CREATE TABLE
    relation_billettholdere_users (
        billettholder_id INTEGER NOT NULL,
        user_id INTEGER NOT NULL,
        inserted_at TEXT DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
        PRIMARY KEY (billettholder_id, user_id),
        FOREIGN KEY (billettholder_id) REFERENCES billettholdere (id),
        FOREIGN KEY (user_id) REFERENCES users (id)
    ) STRICT;

CREATE TABLE
    event_statuses (status TEXT PRIMARY KEY) STRICT;

CREATE TABLE
    events_types (event_type TEXT PRIMARY KEY) STRICT;

CREATE TABLE
    age_groups (age_group TEXT PRIMARY KEY) STRICT;

CREATE TABLE
    event_runtimes (runtime TEXT PRIMARY KEY) STRICT;

CREATE TABLE
    interest_levels (interest_level TEXT PRIMARY KEY) STRICT;

CREATE TABLE
    pulje_statuses (status TEXT PRIMARY KEY) STRICT;

CREATE TABLE
    goose_db_version (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        version_id INTEGER NOT NULL,
        is_applied INTEGER NOT NULL,
        tstamp TEXT DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))
    ) STRICT;

CREATE TABLE
    "events" (
        id TEXT PRIMARY KEY NOT NULL DEFAULT (lower(hex (randomblob (8)))),
        title TEXT NOT NULL,
        intro TEXT NOT NULL,
        description TEXT NOT NULL,
        system TEXT DEFAULT '',
        event_type TEXT NOT NULL DEFAULT 'Other',
        age_group TEXT NOT NULL DEFAULT 'Default',
        event_runtime TEXT NOT NULL DEFAULT 'Normal',
        host_name TEXT NOT NULL,
        user_id INTEGER,
        email TEXT NOT NULL,
        phone_number TEXT NOT NULL,
        max_players INTEGER NOT NULL,
        beginner_friendly INTEGER NOT NULL DEFAULT 0 CHECK (beginner_friendly IN (0, 1)),
        can_be_run_in_english INTEGER NOT NULL DEFAULT 0 CHECK (can_be_run_in_english IN (0, 1)),
        notes TEXT DEFAULT '',
        status TEXT NOT NULL DEFAULT 'Kladd',
        created_at TEXT DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
        updated_at TEXT DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
        created_by_id INTEGER,
        updated_by_id INTEGER,
        status_changed_by_id INTEGER,
        status_changed_at TEXT,
        status_changed_action TEXT,
        FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE SET NULL,
        FOREIGN KEY (updated_by_id) REFERENCES users (id) ON DELETE SET NULL,
        FOREIGN KEY (status_changed_by_id) REFERENCES users (id) ON DELETE SET NULL,
        FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE SET NULL,
        FOREIGN KEY (status) REFERENCES event_statuses (status) ON UPDATE CASCADE,
        FOREIGN KEY (event_type) REFERENCES events_types (event_type) ON UPDATE CASCADE,
        FOREIGN KEY (age_group) REFERENCES age_groups (age_group) ON UPDATE CASCADE,
        FOREIGN KEY (event_runtime) REFERENCES event_runtimes (runtime) ON UPDATE CASCADE
    ) STRICT;

CREATE TABLE
    "billettholdere" (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        first_name TEXT NOT NULL,
        last_name TEXT NOT NULL,
        ticket_type_id INTEGER NOT NULL,
        ticket_type TEXT NOT NULL,
        is_over_18 INTEGER NOT NULL DEFAULT 0 CHECK (is_over_18 IN (0, 1)),
        order_id INTEGER NOT NULL,
        ticket_id INTEGER NOT NULL UNIQUE,
        created_at TEXT DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
        updated_at TEXT DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
        created_by_id INTEGER,
        updated_by_id INTEGER,
        FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE SET NULL,
        FOREIGN KEY (updated_by_id) REFERENCES users (id) ON DELETE SET NULL
    ) STRICT;

CREATE TABLE
    relation_billettholder_emails (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        billettholder_id INTEGER NOT NULL,
        email TEXT NOT NULL COLLATE NOCASE,
        kind TEXT NOT NULL CHECK (kind IN ('Ticket', 'Associated', 'Manual')),
        created_at TEXT DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
        updated_at TEXT DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
        created_by_id INTEGER,
        updated_by_id INTEGER,
        FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE SET NULL,
        FOREIGN KEY (updated_by_id) REFERENCES users (id) ON DELETE SET NULL,
        FOREIGN KEY (billettholder_id) REFERENCES billettholdere (id)
    ) STRICT;

CREATE TABLE
    relation_event_puljer (
        event_id TEXT NOT NULL,
        pulje_id TEXT NOT NULL,
        is_in_pulje INTEGER NOT NULL DEFAULT 1 CHECK (is_in_pulje IN (0, 1)),
        is_published INTEGER NOT NULL DEFAULT 0 CHECK (is_published IN (0, 1)),
        room_id INTEGER,
        PRIMARY KEY (event_id, pulje_id),
        FOREIGN KEY (event_id) REFERENCES events (id) ON DELETE CASCADE,
        FOREIGN KEY (pulje_id) REFERENCES puljer (id) ON UPDATE CASCADE,
        FOREIGN KEY (room_id) REFERENCES rooms (id) ON DELETE SET NULL
    ) STRICT;

CREATE TABLE
    puljer (
        id TEXT NOT NULL PRIMARY KEY,
        name TEXT NOT NULL,
        status TEXT NOT NULL CHECK (
            status IN (
                'not_published',
                'published',
                'locked',
                'completed'
            )
        ),
        start_at TEXT NOT NULL,
        end_at TEXT NOT NULL,
        FOREIGN KEY (status) REFERENCES pulje_statuses (status) ON UPDATE CASCADE
    ) STRICT;

CREATE TABLE
    relation_events_players (
        event_id TEXT NOT NULL,
        pulje_id TEXT NOT NULL,
        billettholder_id INTEGER NOT NULL,
        role TEXT NOT NULL DEFAULT 'Player' CHECK (role IN ('Player', 'GM')),
        inserted_at TEXT DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
        PRIMARY KEY (billettholder_id, event_id, pulje_id),
        FOREIGN KEY (billettholder_id) REFERENCES billettholdere (id),
        FOREIGN KEY (event_id) REFERENCES events (id),
        FOREIGN KEY (pulje_id) REFERENCES puljer (id)
    ) STRICT;

CREATE TABLE
    "interests" (
        billettholder_id INTEGER NOT NULL,
        event_id TEXT NOT NULL,
        pulje_id TEXT NOT NULL,
        interest_level TEXT NOT NULL,
        created_at TEXT DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
        updated_at TEXT DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
        created_by_id INTEGER,
        updated_by_id INTEGER,
        FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE SET NULL,
        FOREIGN KEY (updated_by_id) REFERENCES users (id) ON DELETE SET NULL,
        PRIMARY KEY (billettholder_id, event_id, pulje_id),
        FOREIGN KEY (billettholder_id) REFERENCES billettholdere (id),
        FOREIGN KEY (event_id) REFERENCES events (id),
        FOREIGN KEY (pulje_id) REFERENCES puljer (id),
        FOREIGN KEY (interest_level) REFERENCES interest_levels (interest_level) ON UPDATE CASCADE
    ) STRICT;

CREATE TABLE
    rooms (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        room_number TEXT NOT NULL DEFAULT '',
        name TEXT NOT NULL,
        floor INTEGER NOT NULL,
        max_concurrent_games INTEGER NOT NULL,
        notes TEXT NOT NULL DEFAULT '',
        is_disabled INTEGER NOT NULL DEFAULT 0 CHECK (is_disabled IN (0, 1))
    ) STRICT;

INSERT INTO
    event_statuses (status)
VALUES
    ('Kladd'),
    ('Innsendt'),
    ('Godkjent'),
    ('Forkastet'),
    ('Publisert');

INSERT INTO
    events_types (event_type)
VALUES
    ('Roleplay'),
    ('Boardgame'),
    ('Cardgame'),
    ('Other');

INSERT INTO
    age_groups (age_group)
VALUES
    ('Default'),
    ('ChildFriendly'),
    ('AdultsOnly');

INSERT INTO
    event_runtimes (runtime)
VALUES
    ('Normal'),
    ('ShortRunning'),
    ('LongRunning');

INSERT INTO
    interest_levels (interest_level)
VALUES
    ('Veldig interessert'),
    ('Middels interessert'),
    ('Litt interessert');

INSERT INTO
    pulje_statuses (status)
VALUES
    ('not_published'),
    ('published'),
    ('locked'),
    ('completed');

INSERT INTO
    puljer (id, name, status, start_at, end_at)
VALUES
    (
        'FredagKveld',
        'Fredag kveld',
        'not_published',
        '2025-10-10T18:00:00Z',
        '2025-10-10T23:00:00Z'
    ),
    (
        'LordagMorgen',
        'Lørdag morgen',
        'not_published',
        '2025-10-11T10:00:00Z',
        '2025-10-11T15:00:00Z'
    ),
    (
        'LordagKveld',
        'Lørdag kveld',
        'not_published',
        '2025-10-11T18:00:00Z',
        '2025-10-11T23:00:00Z'
    ),
    (
        'SondagMorgen',
        'Søndag morgen',
        'not_published',
        '2025-10-12T10:00:00Z',
        '2025-10-12T15:00:00Z'
    );

-- Views information /database/views.md

CREATE VIEW
    v_get_user_billettholder AS
SELECT
    u.id AS user_db_id,
    u.external_id AS external_id,
    u.email AS user_email,
    u.is_admin AS user_is_admin,
    u.inserted_at AS user_inserted_at,
    bu.billettholder_id AS billettholder_id,
    bu.user_id AS billettholder_user_db_id,
    bu.inserted_at AS billettholder_user_inserted_at
FROM
    relation_billettholdere_users AS bu
    LEFT JOIN users AS u ON u.id = bu.user_id;

CREATE VIEW
    v_events_by_pulje_active AS
SELECT
    e.id AS id,
    e.title,
    e.intro,
    e.description,
    e.system,
    e.event_type,
    e.age_group,
    e.event_runtime,
    e.host_name,
    e.user_id,
    e.email,
    e.phone_number,
    e.max_players,
    e.beginner_friendly,
    e.can_be_run_in_english,
    e.notes,
    e.status,
    e.created_at,
    ep.is_published AS is_published,
    ep.pulje_id,
    ep.room_id,
    r.room_number,
    r.name AS room_name,
    r.floor AS room_floor,
    r.max_concurrent_games AS room_max_concurrent_games,
    r.notes AS room_notes,
    r.is_disabled AS room_is_disabled,
    p.name AS pulje_name,
    p.start_at AS pulje_start_at,
    p.end_at AS pulje_end_at
FROM
    events e
    INNER JOIN relation_event_puljer ep ON ep.event_id = e.id
    INNER JOIN puljer p ON p.id = ep.pulje_id
    LEFT JOIN rooms r ON r.id = ep.room_id
WHERE
    e.status = 'Godkjent'
    AND ep.is_in_pulje = 1;

CREATE VIEW
    v_event_summary AS
SELECT
    id,
    title,
    intro,
    status,
    system,
    host_name,
    beginner_friendly,
    event_type,
    age_group,
    event_runtime,
    can_be_run_in_english,
    created_at,
    updated_at
FROM
    events;

CREATE VIEW
    v_billettholder_emails AS
SELECT
    b.id AS billettholder_id,
    b.first_name,
    b.last_name,
    b.ticket_type_id,
    b.ticket_type,
    b.is_over_18,
    b.order_id,
    b.ticket_id,
    b.created_at AS billettholder_created_at,
    b.updated_at AS billettholder_updated_at,
    e.id AS email_id,
    e.email,
    e.kind,
    e.created_at AS email_created_at,
    e.updated_at AS email_updated_at
FROM
    billettholdere AS b
    LEFT JOIN relation_billettholder_emails AS e ON b.id = e.billettholder_id;

CREATE VIEW
    v_event_puljer_active AS
SELECT
    ep.event_id,
    ep.pulje_id,
    ep.room_id,
    r.room_number,
    r.name AS room_name,
    r.floor AS room_floor,
    r.max_concurrent_games AS room_max_concurrent_games,
    r.notes AS room_notes,
    r.is_disabled AS room_is_disabled,
    p.name AS pulje_name,
    p.start_at AS pulje_start_at,
    p.end_at AS pulje_end_at,
    ep.is_in_pulje,
    ep.is_published
FROM
    relation_event_puljer ep
    JOIN puljer p ON p.id = ep.pulje_id
    LEFT JOIN rooms r ON r.id = ep.room_id
WHERE
    ep.is_in_pulje = 1
    AND ep.is_published = 1;
