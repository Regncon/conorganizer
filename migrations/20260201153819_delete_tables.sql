-- +goose Up
PRAGMA foreign_keys = OFF;

DROP TABLE IF EXISTS events_puljes_exclusions;
DROP TABLE IF EXISTS _litestream_lock;
DROP TABLE IF EXISTS _litestream_seq;

PRAGMA foreign_keys = ON;

-- +goose Down
PRAGMA foreign_keys = OFF;

CREATE TABLE IF NOT EXISTS events_puljes_exclusions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    pulje_name TEXT NOT NULL,
    event_id TEXT NOT NULL,
    FOREIGN KEY (pulje_name) REFERENCES puljer(name),
    FOREIGN KEY (event_id) REFERENCES events(id)
);

CREATE TABLE IF NOT EXISTS _litestream_seq (
    id INTEGER PRIMARY KEY,
    seq INTEGER
);

CREATE TABLE IF NOT EXISTS _litestream_lock (
    id INTEGER
);

PRAGMA foreign_keys = ON;
