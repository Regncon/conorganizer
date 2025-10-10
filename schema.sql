CREATE TABLE sqlite_sequence(name,seq);
CREATE TABLE users(
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id TEXT NOT NULL UNIQUE,
  email TEXT NOT NULL,
  is_admin BOOLEAN NOT NULL DEFAULT FALSE,
  inserted_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE billettholdere_users(
  billettholder_id INTEGER NOT NULL,
  user_id INTEGER NOT NULL,
  inserted_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY(billettholder_id, user_id),
  FOREIGN KEY(billettholder_id) REFERENCES billettholdere(id),
  FOREIGN KEY(user_id) REFERENCES users(id)
);
CREATE TABLE event_statuses(status TEXT PRIMARY KEY);
CREATE TABLE events_types(event_type TEXT PRIMARY KEY);
CREATE TABLE IF NOT EXISTS "age_groups"(age_group TEXT PRIMARY KEY);
CREATE TABLE event_runtimes(runtime TEXT PRIMARY KEY);
CREATE TABLE interest_levels(interest_level TEXT PRIMARY KEY);
CREATE TABLE events_players(
  event_id TEXT NOT NULL,
  billettholder_id INTEGER NOT NULL,
  interest_level TEXT NOT NULL,
  inserted_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY(event_id, billettholder_id),
  FOREIGN KEY(event_id) REFERENCES events(id),
  FOREIGN KEY(billettholder_id) REFERENCES billettholdere(id),
  FOREIGN KEY(interest_level) REFERENCES interest_levels(interest_level) ON UPDATE CASCADE
);
CREATE TABLE events_puljes_exclusions(
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  pulje_name TEXT NOT NULL,
  event_id TEXT NOT NULL,
  FOREIGN KEY(pulje_name) REFERENCES puljer(name),
  FOREIGN KEY(event_id) REFERENCES events(id)
);
CREATE TABLE _litestream_seq(id INTEGER PRIMARY KEY, seq INTEGER);
CREATE TABLE _litestream_lock(id INTEGER);
CREATE TABLE goose_db_version(
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  version_id INTEGER NOT NULL,
  is_applied INTEGER NOT NULL,
  tstamp TIMESTAMP DEFAULT(datetime('now'))
);
CREATE TABLE IF NOT EXISTS "events"(
  id TEXT PRIMARY KEY NOT NULL DEFAULT(lower(hex(randomblob(8)))),
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
  FOREIGN KEY(host) REFERENCES users(id) ON DELETE SET NULL,
  FOREIGN KEY(pulje_name) REFERENCES puljer(name) ON UPDATE CASCADE,
  FOREIGN KEY(status) REFERENCES event_statuses(status) ON UPDATE CASCADE,
  FOREIGN KEY(event_type) REFERENCES events_types(event_type) ON UPDATE CASCADE,
  FOREIGN KEY(age_group) REFERENCES age_groups(age_group) ON UPDATE CASCADE,
  FOREIGN KEY(event_runtime) REFERENCES event_runtimes(runtime) ON UPDATE CASCADE
);
CREATE TABLE IF NOT EXISTS "billettholdere"(
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  first_name TEXT NOT NULL,
  last_name TEXT NOT NULL,
  ticket_type_id INTEGER NOT NULL,
  ticket_type TEXT NOT NULL,
  is_over_18 BOOLEAN NOT NULL,
  order_id INTEGER NOT NULL,
  ticket_id INTEGER NOT NULL UNIQUE,
  inserted_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE billettholder_emails(
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  billettholder_id INTEGER NOT NULL,
  email TEXT NOT NULL COLLATE NOCASE,
  kind TEXT NOT NULL CHECK(kind IN('Ticket','Associated','Manual')),
  inserted_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY(billettholder_id) REFERENCES billettholdere(id)
);
CREATE TABLE puljer(
  id TEXT NOT NULL PRIMARY KEY,
  name TEXT NOT NULL,
  start_time DATE NOT NULL,
  end_time DATE NOT NULL
);
CREATE TABLE event_puljer(
  event_id TEXT NOT NULL,
  pulje_id TEXT NOT NULL,
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  is_published BOOLEAN NOT NULL DEFAULT FALSE,
  room TEXT DEFAULT '',
  PRIMARY KEY(event_id, pulje_id),
  FOREIGN KEY(event_id) REFERENCES events(id) ON DELETE CASCADE,
  FOREIGN KEY(pulje_id) REFERENCES puljer(id) ON UPDATE CASCADE
);
CREATE TABLE interests(
  billettholder_id INTEGER NOT NULL,
  event_id TEXT NOT NULL,
  pulje_id TEXT NOT NULL,
  interest_level TEXT NOT NULL,
  inserted_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY(billettholder_id, event_id, pulje_id),
  FOREIGN KEY(billettholder_id) REFERENCES billettholdere(id),
  FOREIGN KEY(event_id) REFERENCES events(id),
  FOREIGN KEY(pulje_id) REFERENCES puljer(id),
  FOREIGN KEY(interest_level) REFERENCES interest_levels(interest_level) ON UPDATE CASCADE
);
