CREATE TABLE users(
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  external_id TEXT NOT NULL UNIQUE,
  email TEXT NOT NULL,
  is_admin INTEGER NOT NULL DEFAULT 0 CHECK(is_admin IN(0, 1)),
  inserted_at TEXT DEFAULT(strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))
) STRICT;
CREATE TABLE relation_billettholdere_users(
  billettholder_id INTEGER NOT NULL,
  user_id INTEGER NOT NULL,
  inserted_at TEXT DEFAULT(strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
  PRIMARY KEY(billettholder_id, user_id),
  FOREIGN KEY(billettholder_id) REFERENCES billettholdere(id),
  FOREIGN KEY(user_id) REFERENCES users(id)
) STRICT;
CREATE TABLE event_statuses(status TEXT PRIMARY KEY) STRICT;
CREATE TABLE events_types(event_type TEXT PRIMARY KEY) STRICT;
CREATE TABLE age_groups(age_group TEXT PRIMARY KEY) STRICT;
CREATE TABLE event_runtimes(runtime TEXT PRIMARY KEY) STRICT;
CREATE TABLE interest_levels(interest_level TEXT PRIMARY KEY) STRICT;
CREATE TABLE pulje_statuses(status TEXT PRIMARY KEY) STRICT;
CREATE TABLE goose_db_version(
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  version_id INTEGER NOT NULL,
  is_applied INTEGER NOT NULL,
  tstamp TEXT DEFAULT(strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))
) STRICT;
CREATE TABLE "billettholdere"(
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  first_name TEXT NOT NULL,
  last_name TEXT NOT NULL,
  ticket_type_id INTEGER NOT NULL,
  ticket_type TEXT NOT NULL,
  is_over_18 INTEGER NOT NULL DEFAULT 0 CHECK(is_over_18 IN(0, 1)),
  order_id INTEGER NOT NULL,
  ticket_id INTEGER NOT NULL UNIQUE,
  created_at TEXT DEFAULT(strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
  updated_at TEXT DEFAULT(strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
  created_by_id INTEGER,
  updated_by_id INTEGER,
  FOREIGN KEY(created_by_id) REFERENCES users(id) ON DELETE SET NULL,
  FOREIGN KEY(updated_by_id) REFERENCES users(id) ON DELETE SET NULL
) STRICT;
CREATE TABLE relation_billettholder_emails(
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  billettholder_id INTEGER NOT NULL,
  email TEXT NOT NULL COLLATE NOCASE,
  kind TEXT NOT NULL CHECK(kind IN('Ticket', 'Associated', 'Manual')),
  created_at TEXT DEFAULT(strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
  updated_at TEXT DEFAULT(strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
  created_by_id INTEGER,
  updated_by_id INTEGER,
  FOREIGN KEY(created_by_id) REFERENCES users(id) ON DELETE SET NULL,
  FOREIGN KEY(updated_by_id) REFERENCES users(id) ON DELETE SET NULL,
  FOREIGN KEY(billettholder_id) REFERENCES billettholdere(id)
) STRICT;
CREATE TABLE relation_event_puljer(
  event_id TEXT NOT NULL,
  pulje_id TEXT NOT NULL,
  is_in_pulje INTEGER NOT NULL DEFAULT 1 CHECK(is_in_pulje IN(0, 1)),
  is_published INTEGER NOT NULL DEFAULT 0 CHECK(is_published IN(0, 1)),
  room_id INTEGER,
  PRIMARY KEY(event_id, pulje_id),
  FOREIGN KEY(event_id) REFERENCES events(id) ON DELETE CASCADE,
  FOREIGN KEY(pulje_id) REFERENCES puljer(id) ON UPDATE CASCADE,
  FOREIGN KEY(room_id) REFERENCES rooms(id) ON DELETE SET NULL
) STRICT;
CREATE TABLE relation_events_players(
  event_id TEXT NOT NULL,
  pulje_id TEXT NOT NULL,
  billettholder_id INTEGER NOT NULL,
  role TEXT NOT NULL DEFAULT 'Player' CHECK(role IN('Player', 'GM')),
  source TEXT NOT NULL DEFAULT 'manual' CHECK (source IN ('manual','solver')),
  inserted_at TEXT DEFAULT(strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
  PRIMARY KEY(billettholder_id, event_id, pulje_id),
  FOREIGN KEY(billettholder_id) REFERENCES billettholdere(id),
  FOREIGN KEY(event_id) REFERENCES events(id),
  FOREIGN KEY(pulje_id) REFERENCES puljer(id)
) STRICT;
CREATE TABLE "interests"(
  billettholder_id INTEGER NOT NULL,
  event_id TEXT NOT NULL,
  pulje_id TEXT NOT NULL,
  interest_level TEXT NOT NULL,
  created_at TEXT DEFAULT(strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
  updated_at TEXT DEFAULT(strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
  created_by_id INTEGER,
  updated_by_id INTEGER,
  FOREIGN KEY(created_by_id) REFERENCES users(id) ON DELETE SET NULL,
  FOREIGN KEY(updated_by_id) REFERENCES users(id) ON DELETE SET NULL,
  PRIMARY KEY(billettholder_id, event_id, pulje_id),
  FOREIGN KEY(billettholder_id) REFERENCES billettholdere(id),
  FOREIGN KEY(event_id) REFERENCES events(id),
  FOREIGN KEY(pulje_id) REFERENCES puljer(id),
  FOREIGN KEY(interest_level) REFERENCES interest_levels(interest_level) ON UPDATE CASCADE
) STRICT;
CREATE TABLE "rooms"(
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  room_number TEXT NOT NULL DEFAULT '',
  name TEXT NOT NULL,
  floor INTEGER NOT NULL,
  max_concurrent_games INTEGER NOT NULL,
  notes TEXT NOT NULL DEFAULT '',
  is_disabled INTEGER NOT NULL DEFAULT 0 CHECK(is_disabled IN(0, 1))
) STRICT;
CREATE TABLE "events"(
  id TEXT PRIMARY KEY NOT NULL DEFAULT(lower(hex(randomblob(8)))),
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
  beginner_friendly INTEGER NOT NULL DEFAULT 0 CHECK(beginner_friendly IN(0, 1)),
  can_be_run_in_english INTEGER NOT NULL DEFAULT 0 CHECK(can_be_run_in_english IN(0, 1)),
  notes TEXT DEFAULT '',
  status TEXT NOT NULL DEFAULT 'Kladd',
  created_at TEXT DEFAULT(strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
  updated_at TEXT DEFAULT(strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
  created_by_id INTEGER,
  updated_by_id INTEGER,
  status_changed_by_id INTEGER,
  status_changed_at TEXT,
  status_changed_action TEXT,
  FOREIGN KEY(created_by_id) REFERENCES users(id) ON DELETE SET NULL,
  FOREIGN KEY(updated_by_id) REFERENCES users(id) ON DELETE SET NULL,
  FOREIGN KEY(status_changed_by_id) REFERENCES users(id) ON DELETE SET NULL,
  FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE SET NULL,
  FOREIGN KEY(status) REFERENCES event_statuses(status) ON UPDATE CASCADE,
  FOREIGN KEY(event_type) REFERENCES events_types(event_type) ON UPDATE CASCADE,
  FOREIGN KEY(age_group) REFERENCES age_groups(age_group) ON UPDATE CASCADE,
  FOREIGN KEY(event_runtime) REFERENCES event_runtimes(runtime) ON UPDATE CASCADE
) STRICT;
CREATE VIEW v_get_user_billettholder AS
SELECT
    u.id AS user_id,
    u.external_id AS external_id,
    u.email AS user_email,
    u.is_admin AS user_is_admin,
    u.inserted_at AS user_inserted_at,
    bu.billettholder_id AS billettholder_id,
    bu.user_id AS billettholder_user_id,
    bu.inserted_at AS billettholder_user_inserted_at
FROM
    relation_billettholdere_users AS bu
    LEFT JOIN users AS u ON u.id = bu.user_id
/* v_get_user_billettholder(user_id,external_id,user_email,user_is_admin,user_inserted_at,billettholder_id,billettholder_user_id,billettholder_user_inserted_at) */;
CREATE VIEW v_event_summary AS
SELECT
    e.id,
    e.title,
    e.intro,
    e.status,
    e.system,
    e.host_name,
    e.beginner_friendly,
    e.event_type,
    e.age_group,
    e.event_runtime,
    e.can_be_run_in_english,
    e.created_at,
    e.updated_at
FROM
    events e
/* v_event_summary(id,title,intro,status,system,host_name,beginner_friendly,event_type,age_group,event_runtime,can_be_run_in_english,created_at,updated_at) */;
CREATE VIEW v_billettholder_emails AS
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
    LEFT JOIN relation_billettholder_emails AS e ON b.id = e.billettholder_id
/* v_billettholder_emails(billettholder_id,first_name,last_name,ticket_type_id,ticket_type,is_over_18,order_id,ticket_id,billettholder_created_at,billettholder_updated_at,email_id,email,kind,email_created_at,email_updated_at) */;
CREATE VIEW v_event_puljer_active AS
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
    AND ep.is_published = 1
/* v_event_puljer_active(event_id,pulje_id,room_id,room_number,room_name,room_floor,room_max_concurrent_games,room_notes,room_is_disabled,pulje_name,pulje_start_at,pulje_end_at,is_in_pulje,is_published) */;
CREATE TABLE program_publishing_state(
  id INTEGER NOT NULL PRIMARY KEY CHECK(id = 1),
  is_published INTEGER NOT NULL DEFAULT 0 CHECK(is_published IN(0, 1))
) STRICT;
CREATE TABLE "puljer"(
  id TEXT NOT NULL PRIMARY KEY,
  name TEXT NOT NULL,
  status TEXT NOT NULL DEFAULT 'Open' CHECK(status IN('Open', 'Locked', 'Completed')),
  start_at TEXT NOT NULL,
  end_at TEXT NOT NULL,
  FOREIGN KEY(status) REFERENCES pulje_statuses(status) ON UPDATE CASCADE
) STRICT;
CREATE VIEW v_events_by_pulje_active AS
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
FROM events e
INNER JOIN relation_event_puljer ep ON ep.event_id = e.id
INNER JOIN puljer p ON p.id = ep.pulje_id
LEFT JOIN rooms r ON r.id = ep.room_id
WHERE
    e.status = 'Annonsert'
    AND ep.is_in_pulje = 1
/* v_events_by_pulje_active(id,title,intro,description,system,event_type,age_group,event_runtime,host_name,user_id,email,phone_number,max_players,beginner_friendly,can_be_run_in_english,notes,status,created_at,is_published,pulje_id,room_id,room_number,room_name,room_floor,room_max_concurrent_games,room_notes,room_is_disabled,pulje_name,pulje_start_at,pulje_end_at) */;
