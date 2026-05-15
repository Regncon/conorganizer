CREATE TABLE
    users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        external_id TEXT NOT NULL UNIQUE,
        email TEXT NOT NULL,
        is_admin BOOLEAN NOT NULL DEFAULT FALSE,
        inserted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

CREATE TABLE
    relation_billettholdere_users (
        billettholder_id INTEGER NOT NULL,
        user_id INTEGER NOT NULL,
        inserted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        PRIMARY KEY (billettholder_id, user_id),
        FOREIGN KEY (billettholder_id) REFERENCES billettholdere (id),
        FOREIGN KEY (user_id) REFERENCES users (id)
    );

CREATE TABLE
    event_statuses (status TEXT PRIMARY KEY);

CREATE TABLE
    events_types (event_type TEXT PRIMARY KEY);

CREATE TABLE
    age_groups (age_group TEXT PRIMARY KEY);

CREATE TABLE
    event_runtimes (runtime TEXT PRIMARY KEY);

CREATE TABLE
    interest_levels (interest_level TEXT PRIMARY KEY);

CREATE TABLE
    pulje_statuses (status TEXT PRIMARY KEY);

CREATE TABLE
    goose_db_version (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        version_id INTEGER NOT NULL,
        is_applied INTEGER NOT NULL,
        tstamp TIMESTAMP DEFAULT (datetime ('now'))
    );

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
        beginner_friendly BOOLEAN NOT NULL,
        can_be_run_in_english BOOLEAN NOT NULL,
        notes TEXT DEFAULT '',
        status TEXT NOT NULL DEFAULT 'Kladd',
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        created_by_id INTEGER,
        updated_by_id INTEGER,
        status_changed_by_id INTEGER,
        status_changed_at TIMESTAMP,
        status_changed_action TEXT,
        FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE SET NULL,
        FOREIGN KEY (updated_by_id) REFERENCES users (id) ON DELETE SET NULL,
        FOREIGN KEY (status_changed_by_id) REFERENCES users (id) ON DELETE SET NULL,
        FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE SET NULL,
        FOREIGN KEY (status) REFERENCES event_statuses (status) ON UPDATE CASCADE,
        FOREIGN KEY (event_type) REFERENCES events_types (event_type) ON UPDATE CASCADE,
        FOREIGN KEY (age_group) REFERENCES age_groups (age_group) ON UPDATE CASCADE,
        FOREIGN KEY (event_runtime) REFERENCES event_runtimes (runtime) ON UPDATE CASCADE
    );

CREATE TABLE
    "billettholdere" (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        first_name TEXT NOT NULL,
        last_name TEXT NOT NULL,
        ticket_type_id INTEGER NOT NULL,
        ticket_type TEXT NOT NULL,
        is_over_18 BOOLEAN NOT NULL,
        order_id INTEGER NOT NULL,
        ticket_id INTEGER NOT NULL UNIQUE,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        created_by_id INTEGER,
        updated_by_id INTEGER,
        FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE SET NULL,
        FOREIGN KEY (updated_by_id) REFERENCES users (id) ON DELETE SET NULL
    );

CREATE TABLE
    relation_billettholder_emails (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        billettholder_id INTEGER NOT NULL,
        email TEXT NOT NULL COLLATE NOCASE,
        kind TEXT NOT NULL CHECK (kind IN ('Ticket', 'Associated', 'Manual')),
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        created_by_id INTEGER,
        updated_by_id INTEGER,
        FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE SET NULL,
        FOREIGN KEY (updated_by_id) REFERENCES users (id) ON DELETE SET NULL,
        FOREIGN KEY (billettholder_id) REFERENCES billettholdere (id)
    );

CREATE TABLE
    relation_event_puljer (
        event_id TEXT NOT NULL,
        pulje_id TEXT NOT NULL,
        is_in_pulje BOOLEAN NOT NULL DEFAULT TRUE,
        is_published BOOLEAN NOT NULL DEFAULT FALSE,
        room_id INTEGER,
        PRIMARY KEY (event_id, pulje_id),
        FOREIGN KEY (event_id) REFERENCES events (id) ON DELETE CASCADE,
        FOREIGN KEY (pulje_id) REFERENCES puljer (id) ON UPDATE CASCADE,
        FOREIGN KEY (room_id) REFERENCES rooms (id) ON DELETE SET NULL
    );

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
        start_at DATE NOT NULL,
        end_at DATE NOT NULL,
        FOREIGN KEY (status) REFERENCES pulje_statuses (status) ON UPDATE CASCADE
    );

CREATE TABLE
    relation_events_players (
        event_id TEXT NOT NULL,
        pulje_id TEXT NOT NULL,
        billettholder_id INTEGER NOT NULL,
        role TEXT NOT NULL DEFAULT 'Player' CHECK (role IN ('Player', 'GM')),
        inserted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        PRIMARY KEY (billettholder_id, event_id, pulje_id),
        FOREIGN KEY (billettholder_id) REFERENCES billettholdere (id),
        FOREIGN KEY (event_id) REFERENCES events (id),
        FOREIGN KEY (pulje_id) REFERENCES puljer (id)
    );

CREATE TABLE
    "interests" (
        billettholder_id INTEGER NOT NULL,
        event_id TEXT NOT NULL,
        pulje_id TEXT NOT NULL,
        interest_level TEXT NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        created_by_id INTEGER,
        updated_by_id INTEGER,
        FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE SET NULL,
        FOREIGN KEY (updated_by_id) REFERENCES users (id) ON DELETE SET NULL,
        PRIMARY KEY (billettholder_id, event_id, pulje_id),
        FOREIGN KEY (billettholder_id) REFERENCES billettholdere (id),
        FOREIGN KEY (event_id) REFERENCES events (id),
        FOREIGN KEY (pulje_id) REFERENCES puljer (id),
        FOREIGN KEY (interest_level) REFERENCES interest_levels (interest_level) ON UPDATE CASCADE
    );

CREATE TABLE
    rooms (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        floor INTEGER NOT NULL,
        max_concurrent_games INTEGER NOT NULL
    );

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
    r.name AS room_name,
    r.floor AS room_floor,
    r.max_concurrent_games AS room_max_concurrent_games,
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
    r.name AS room_name,
    r.floor AS room_floor,
    r.max_concurrent_games AS room_max_concurrent_games,
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
