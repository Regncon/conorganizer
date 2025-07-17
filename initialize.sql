CREATE TABLE IF NOT EXISTS puljer (
    name TEXT PRIMARY KEY,
    start_time DATE NOT NULL
);

INSERT INTO puljer (name, start_time) VALUES
('Fredag kveld', '2025-09-06T18:00:00Z'),
('Lørdag morgen', '2025-09-07T10:00:00Z'),
('Lørdag kveld', '2025-09-07T18:00:00Z'),
('Søndag morgen', '2025-09-08T10:00:00Z');

CREATE TABLE IF NOT EXISTS ticket_types (
    name TEXT PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS billettholdere (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    ticket_type TEXT NOT NULL,
    is_over_18 BOOLEAN NOT NULL,
    order_id INTEGER NOT NULL,
    ticket_id INTEGER NOT NULL UNIQUE,
    ticket_email TEXT NOT NULL,
    order_email TEXT NOT NULL,
    ticket_category_id TEXT NOT NULL,
    inserted_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (ticket_type) REFERENCES ticket_types(name)
);

CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL,
    is_admin BOOLEAN NOT NULL DEFAULT FALSE,
    inserted_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS billettholdere_users (
    billettholder_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    inserted_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (billettholder_id, user_id),
    FOREIGN KEY (billettholder_id) REFERENCES billettholdere(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS event_statuses (
    status TEXT PRIMARY KEY
);

INSERT INTO event_statuses (status) VALUES
('Kladd'),
('Innsendt'),
('Publisert'),
('Godkjent'),
('Avvist');

CREATE TABLE IF NOT EXISTS events_types (
    event_type TEXT PRIMARY KEY
);

INSERT INTO events_types (event_type) VALUES
('Roleplay'),
('Boardgame'),
('Cardgame'),
('Other');

CREATE TABLE IF NOT EXISTS age_grups (
    age_group TEXT PRIMARY KEY
);

INSERT INTO age_grups (age_group) VALUES
('AllAges'),
('ChildFriendly'),
('TeenFriendly'),
('AdultsOnly');

CREATE TABLE IF NOT EXISTS event_runtimes (
    runtime TEXT PRIMARY KEY
);

INSERT INTO event_runtimes (runtime) VALUES
('Normal'),
('ShortRunning'),
('LongRunning');

CREATE TABLE IF NOT EXISTS events (
    id TEXT PRIMARY KEY NOT NULL DEFAULT ( lower(hex(randomblob(8))) ),
    title TEXT NOT NULL,
    intro TEXT NOT NULL,
    description TEXT NOT NULL,
    image_url TEXT DEFAULT '',
    system TEXT DEFAULT '',
    event_type TEXT NOT NULL DEFAULT 'Other',
    age_group TEXT NOT NULL DEFAULT 'AllAges',
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
    FOREIGN KEY (age_group) REFERENCES age_grups(age_group) ON UPDATE CASCADE,
    FOREIGN KEY (event_runtime) REFERENCES event_runtimes(runtime) ON UPDATE CASCADE,
);

CREATE TABLE IF NOT EXISTS interest_levels (
    interest_level TEXT PRIMARY KEY
);

INSERT INTO interest_levels (interest_level) VALUES
('Litt interessert'),
('Middels interessert'),
('Veldig interessert');

CREATE TABLE IF NOT EXISTS interests (
    billettholder_id INTEGER NOT NULL,
    event_id TEXT NOT NULL,
    interest_level TEXT NOT NULL,
    inserted_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (billettholder_id, event_id),
    FOREIGN KEY (billettholder_id) REFERENCES billettholdere(id),
    FOREIGN KEY (event_id) REFERENCES events(id),
    FOREIGN KEY (interest_level) REFERENCES interest_levels(interest_level) ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS events_players (
    event_id TEXT NOT NULL,
    billettholder_id INTEGER NOT NULL,
	interest_level TEXT NOT NULL,
    inserted_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (event_id, billettholder_id),
    FOREIGN KEY (event_id) REFERENCES events(id),
    FOREIGN KEY (billettholder_id) REFERENCES billettholdere(id),
	FOREIGN KEY (interest_level) REFERENCES interest_levels(interest_level) ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS events_puljes_exclusions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    pulje_name TEXT NOT NULL,
    event_id TEXT NOT NULL,
    FOREIGN KEY (pulje_name) REFERENCES puljer(name),
    FOREIGN KEY (event_id) REFERENCES events(id)
);

