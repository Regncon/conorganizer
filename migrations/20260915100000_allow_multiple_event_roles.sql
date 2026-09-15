-- +goose Up
-- Existing multiple Player assignments in one pulje are preserved intentionally.
-- Application validation reports these legacy conflicts for manual resolution.
CREATE TABLE relation_events_players_new (
    event_id TEXT NOT NULL,
    pulje_id TEXT NOT NULL,
    billettholder_id INTEGER NOT NULL,
    role TEXT NOT NULL DEFAULT 'Player' CHECK (role IN ('Player', 'GM')),
    inserted_at TEXT DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    source TEXT NOT NULL DEFAULT 'manual' CHECK (source IN ('manual', 'solver')),
    PRIMARY KEY (billettholder_id, event_id, pulje_id, role),
    FOREIGN KEY (billettholder_id) REFERENCES billettholdere (id),
    FOREIGN KEY (event_id) REFERENCES events (id),
    FOREIGN KEY (pulje_id) REFERENCES puljer (id)
) STRICT;

INSERT INTO relation_events_players_new (
    event_id,
    pulje_id,
    billettholder_id,
    role,
    inserted_at,
    source
)
SELECT
    event_id,
    pulje_id,
    billettholder_id,
    role,
    inserted_at,
    source
FROM relation_events_players;

DROP TABLE relation_events_players;
ALTER TABLE relation_events_players_new RENAME TO relation_events_players;

-- +goose Down
-- Copying into the old key fails before the source table is dropped when dual roles exist.
CREATE TABLE relation_events_players_old (
    event_id TEXT NOT NULL,
    pulje_id TEXT NOT NULL,
    billettholder_id INTEGER NOT NULL,
    role TEXT NOT NULL DEFAULT 'Player' CHECK (role IN ('Player', 'GM')),
    inserted_at TEXT DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    source TEXT NOT NULL DEFAULT 'manual' CHECK (source IN ('manual', 'solver')),
    PRIMARY KEY (billettholder_id, event_id, pulje_id),
    FOREIGN KEY (billettholder_id) REFERENCES billettholdere (id),
    FOREIGN KEY (event_id) REFERENCES events (id),
    FOREIGN KEY (pulje_id) REFERENCES puljer (id)
) STRICT;

INSERT INTO relation_events_players_old (
    event_id,
    pulje_id,
    billettholder_id,
    role,
    inserted_at,
    source
)
SELECT
    event_id,
    pulje_id,
    billettholder_id,
    role,
    inserted_at,
    source
FROM relation_events_players;

DROP TABLE relation_events_players;
ALTER TABLE relation_events_players_old RENAME TO relation_events_players;
