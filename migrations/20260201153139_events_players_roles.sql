-- +goose Up
PRAGMA foreign_keys = OFF;

CREATE TABLE IF NOT EXISTS events_players_new (
    event_id TEXT NOT NULL,
    pulje_id TEXT NOT NULL,
    billettholder_id INTEGER NOT NULL,
    role TEXT NOT NULL DEFAULT 'Player' CHECK (role IN ('Player','GM')),
    inserted_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (billettholder_id, event_id, pulje_id),
    FOREIGN KEY (billettholder_id) REFERENCES billettholdere (id),
    FOREIGN KEY (event_id) REFERENCES events (id),
    FOREIGN KEY (pulje_id) REFERENCES puljer (id)
);

INSERT INTO events_players_new (event_id, pulje_id, billettholder_id, role, inserted_time)
SELECT
    event_id,
    pulje_id,
    billettholder_id,
    CASE
        WHEN is_gm = 1 THEN 'GM'
        WHEN is_player = 1 THEN 'Player'
        ELSE 'Player'
    END AS role,
    inserted_time
FROM events_players;

DROP TABLE events_players;
ALTER TABLE events_players_new RENAME TO events_players;

PRAGMA foreign_keys = ON;

-- +goose Down
PRAGMA foreign_keys = OFF;

CREATE TABLE IF NOT EXISTS events_players_old (
    event_id TEXT NOT NULL,
    pulje_id TEXT NOT NULL,
    billettholder_id INTEGER NOT NULL,
    is_player BOOLEAN NOT NULL DEFAULT TRUE,
    is_gm BOOLEAN NOT NULL DEFAULT FALSE,
    inserted_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (billettholder_id, event_id, pulje_id),
    FOREIGN KEY (billettholder_id) REFERENCES billettholdere (id),
    FOREIGN KEY (event_id) REFERENCES events (id),
    FOREIGN KEY (pulje_id) REFERENCES puljer (id)
);

INSERT INTO events_players_old (event_id, pulje_id, billettholder_id, is_player, is_gm, inserted_time)
SELECT
    event_id,
    pulje_id,
    billettholder_id,
    CASE WHEN role = 'Player' THEN 1 ELSE 0 END AS is_player,
    CASE WHEN role = 'GM' THEN 1 ELSE 0 END AS is_gm,
    inserted_time
FROM events_players;

DROP TABLE events_players;
ALTER TABLE events_players_old RENAME TO events_players;

PRAGMA foreign_keys = ON;
