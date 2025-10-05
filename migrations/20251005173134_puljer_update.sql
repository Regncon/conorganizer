-- +goose Up
PRAGMA foreign_keys = ON;

DROP TABLE IF EXISTS puljer;

CREATE TABLE IF NOT EXISTS puljer(
  id TEXT NOT NULL PRIMARY KEY,
  name TEXT NOT NULL,
  start_time DATE NOT NULL,
  end_time DATE NOT NULL
);

-- +goose StatementBegin
INSERT INTO puljer (id, name, start_time, end_time) VALUES
('FredagKveld',  'Fredag kveld', '2025-10-10T18:00:00Z', '2025-10-10T22:00:00Z'),
('LordagMorgen', 'Lørdag morgen', '2025-10-11T10:00:00Z', '2025-10-11T15:00:00Z'),
('LordagKveld',  'Lørdag kveld', '2025-10-11T18:00:00Z', '2025-10-11T22:00:00Z'),
('SondagMorgen', 'Søndag morgen', '2025-10-12T10:00:00Z', '2025-10-12T15:00:00Z');
-- +goose StatementEnd

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


-- +goose Down
PRAGMA foreign_keys = ON;

DROP TABLE IF EXISTS event_puljer;

DROP TABLE IF EXISTS puljer;

CREATE TABLE IF NOT EXISTS puljer (
    name TEXT PRIMARY KEY,
    start_time DATE NOT NULL
);

-- +goose StatementBegin
INSERT INTO puljer (name, start_time) VALUES
('Fredag kveld', '2025-09-06T18:00:00Z'),
('Lørdag morgen', '2025-09-07T10:00:00Z'),
('Lørdag kveld', '2025-09-07T18:00:00Z'),
('Søndag morgen', '2025-09-08T10:00:00Z');
-- +goose StatementEnd
