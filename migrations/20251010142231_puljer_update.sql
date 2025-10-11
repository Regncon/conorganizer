-- +goose Up
PRAGMA foreign_keys = OFF;

DROP TABLE IF EXISTS puljer;
CREATE TABLE IF NOT EXISTS puljer (
    id TEXT NOT NULL PRIMARY KEY,
    name TEXT NOT NULL,
    is_closed   INTEGER NOT NULL DEFAULT FALSE,
    is_published INTEGER NOT NULL DEFAULT FALSE,
    start_time DATE NOT NULL,
    end_time DATE NOT NULL
);

INSERT INTO puljer (id, name, start_time, end_time) VALUES
('FredagKveld',  'Fredag kveld', '2025-10-10T18:00:00Z', '2025-10-10T23:00:00Z'),
('LordagMorgen', 'Lørdag morgen', '2025-10-11T10:00:00Z', '2025-10-11T15:00:00Z'),
('LordagKveld',  'Lørdag kveld', '2025-10-11T18:00:00Z', '2025-10-11T23:00:00Z'),
('SondagMorgen', 'Søndag morgen', '2025-10-12T10:00:00Z', '2025-10-12T15:00:00Z');


PRAGMA foreign_keys = ON;


-- +goose Down

