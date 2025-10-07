-- +goose Up
PRAGMA foreign_keys = ON;

DROP TABLE interests;

CREATE TABLE IF NOT EXISTS interests_new (
    billettholder_id INTEGER NOT NULL,
    event_id TEXT NOT NULL,
    pulje_id TEXT NOT NULL,
    interest_level TEXT NOT NULL,
    inserted_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (billettholder_id, event_id, pulje_id),
    FOREIGN KEY (billettholder_id) REFERENCES billettholdere(id),
    FOREIGN KEY (event_id) REFERENCES events(id),
    FOREIGN KEY (pulje_id) REFERENCES puljer(id),
    FOREIGN KEY (interest_level) REFERENCES interest_levels(interest_level) ON UPDATE CASCADE
);

-- +goose Down
