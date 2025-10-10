-- +goose Up
PRAGMA foreign_keys = ON;
DROP TABLE events_players;

CREATE TABLE IF NOT EXISTS events_players (
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

-- +goose Down
