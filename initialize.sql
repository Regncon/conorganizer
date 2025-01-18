CREATE TABLE IF NOT EXISTS rooms (
    name TEXT PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS rooms_map_data (
	name TEXT PRIMARY KEY,
	x INTEGER NOT NULL,
	y INTEGER NOT NULL,
	FOREIGN KEY (name) REFERENCES rooms(name) ON DELETE CASCADE
);

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

CREATE TABLE IF NOT EXISTS ticketholders (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    ticket_type TEXT NOT NULL,
    is_over_18 BOOLEAN NOT NULL,
    order_id INTEGER NOT NULL,
    ticket_id INTEGER NOT NULL,
    ticket_email TEXT NOT NULL,
    order_email TEXT NOT NULL,
    ticket_category_id TEXT NOT NULL,
    inserted_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (ticket_type) REFERENCES ticket_types(name)
);

CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT NOT NULL,
    is_admin BOOLEAN NOT NULL DEFAULT FALSE,
    inserted_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS ticketholders_users (
    ticketholder_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    inserted_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (ticketholder_id, user_id),
    FOREIGN KEY (ticketholder_id) REFERENCES ticketholders(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS suggested_event (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	host INTEGER NOT NULL,
	title TEXT,
	description TEXT,
	image_url TEXT,
	system TEXT,
	max_players INTEGER,
    child_friendly BOOLEAN NOT NULL,
    adults_only BOOLEAN NOT NULL,
    beginner_friendly BOOLEAN NOT NULL,
    experienced_only BOOLEAN NOT NULL,
    can_be_run_in_english BOOLEAN NOT NULL,
    long_running BOOLEAN NOT NULL,
    short_running BOOLEAN NOT NULL,
	submitted_time TIMESTAMP,
    inserted_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (host) REFERENCES users(id) ON DELETE CASCADE,
	-- Ensure some flags are mutually exclusive
    CHECK (child_friendly + adults_only <= 1),
    CHECK (beginner_friendly + experienced_only <= 1),
    CHECK (long_running + short_running <= 1)
);

CREATE TABLE IF NOT EXISTS events (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
	suggested_event_id INTEGER,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
	image_url TEXT,
	system TEXT,
    host_name TEXT NOT NULL,
    host INTEGER,
    room_name INTEGER,
    pulje_name INTEGER,
    max_players INTEGER NOT NULL,
    child_friendly BOOLEAN NOT NULL,
    adults_only BOOLEAN NOT NULL,
    beginner_friendly BOOLEAN NOT NULL,
    experienced_only BOOLEAN NOT NULL,
    can_be_run_in_english BOOLEAN NOT NULL,
    long_running BOOLEAN NOT NULL,
    short_running BOOLEAN NOT NULL,
    inserted_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (suggested_event_id) REFERENCES suggested_event(id) ON DELETE SET NULL,
    FOREIGN KEY (host) REFERENCES users(id) ON DELETE SET NULL,
    FOREIGN KEY (room_name) REFERENCES rooms(name) ON UPDATE CASCADE,
    FOREIGN KEY (pulje_name) REFERENCES puljer(name) ON UPDATE CASCADE,
	-- Ensure some flags are mutually exclusive
    CHECK (child_friendly + adults_only <= 1),
    CHECK (beginner_friendly + experienced_only <= 1),
    CHECK (long_running + short_running <= 1)
);

CREATE TABLE IF NOT EXISTS interest_levels (
    interest_level TEXT PRIMARY KEY
);

INSERT INTO interest_levels (interest_level) VALUES
('Litt interessert'),
('Middels interessert'),
('Veldig interessert');

CREATE TABLE IF NOT EXISTS interests (
    ticketholder_id INTEGER NOT NULL,
    event_id INTEGER NOT NULL,
    interest_level TEXT NOT NULL,
    inserted_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (ticketholder_id, event_id),
    FOREIGN KEY (ticketholder_id) REFERENCES ticketholders(id),
    FOREIGN KEY (event_id) REFERENCES events(id),
    FOREIGN KEY (interest_level) REFERENCES interest_levels(interest_level) ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS events_players (
    event_id INTEGER NOT NULL,
    ticketholder_id INTEGER NOT NULL,
	interest_level TEXT NOT NULL,
    inserted_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (event_id, ticketholder_id),
    FOREIGN KEY (event_id) REFERENCES events(id),
    FOREIGN KEY (ticketholder_id) REFERENCES ticketholders(id),
	FOREIGN KEY (interest_level) REFERENCES interest_levels(interest_level) ON UPDATE CASCADE
);
