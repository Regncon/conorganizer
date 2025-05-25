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
('Publisert'),
('Godkjent'),
('Avvist');

CREATE TABLE IF NOT EXISTS events (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    image_url TEXT,
    system TEXT,
    host_name TEXT NOT NULL,
    host INTEGER,
    email TEXT NOT NULL,
    phone_number INTEGER NOT NULL,
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
    status TEXT NOT NULL DEFAULT 'Kladd',
    inserted_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (host) REFERENCES users(id) ON DELETE SET NULL,
    FOREIGN KEY (room_name) REFERENCES rooms(name) ON UPDATE CASCADE,
    FOREIGN KEY (pulje_name) REFERENCES puljer(name) ON UPDATE CASCADE,
    FOREIGN KEY (status) REFERENCES event_statuses(status) ON UPDATE CASCADE,
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
    billettholder_id INTEGER NOT NULL,
    event_id INTEGER NOT NULL,
    interest_level TEXT NOT NULL,
    inserted_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (billettholder_id, event_id),
    FOREIGN KEY (billettholder_id) REFERENCES billettholdere(id),
    FOREIGN KEY (event_id) REFERENCES events(id),
    FOREIGN KEY (interest_level) REFERENCES interest_levels(interest_level) ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS events_players (
    event_id INTEGER NOT NULL,
    billettholder_id INTEGER NOT NULL,
	interest_level TEXT NOT NULL,
    inserted_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (event_id, billettholder_id),
    FOREIGN KEY (event_id) REFERENCES events(id),
    FOREIGN KEY (billettholder_id) REFERENCES billettholdere(id),
	FOREIGN KEY (interest_level) REFERENCES interest_levels(interest_level) ON UPDATE CASCADE
);

-- Add 2 test events and a test user to to database

-- Add a test admin user
INSERT INTO users (email, is_admin) VALUES
('test.admin@example.com', true);

-- Add two example events
INSERT INTO events (
    title,
    description,
    image_url,
    system,
    host_name,
    email,
    phone_number,
    host,
    pulje_name,
    max_players,
    child_friendly,
    adults_only,
    beginner_friendly,
    experienced_only,
    can_be_run_in_english,
    long_running,
    short_running,
    status
) VALUES (
    'Dungeons & Dragons: Den Tapte Minen',
    'Bli med på et spennende eventyr i den klassiske D&D-modulen "Den Tapte Minen av Phandelver". Perfekt for nye spillere!',
    'https://imgur.com/example1',
    'D&D 5e',
    'Erik Spilleder',
    'test.admin@example.com',
    12345678,
    1, -- refererer til test.admin@example.com
    'Lørdag morgen',
    6,
    false,
    false,
    true,
    false,
    true,
    false,
    true,
    'Publisert'
),
(
    'Vampire: Nattens Barn',
    'En intens fortelling om intriger og makt i Oslos vampyrsamfunn. Kun for erfarne rollespillere.',
    'https://imgur.com/example2',
    'Vampire: The Masquerade 5th Edition',
    'Maria Storyteller',
    'test.admin@example.com',
    12345678,
    1, -- refererer til test.admin@example.com
    'Lørdag kveld',
    4,
    false,
    true,
    false,
    true,
    false,
    true,
    false,
    'Publisert'
);
