-- +goose Up
PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS billettholder_emails (
    id                INTEGER PRIMARY KEY AUTOINCREMENT,
    billettholder_id  INTEGER NOT NULL,
    email             TEXT    NOT NULL COLLATE NOCASE,
    kind              TEXT    NOT NULL CHECK (kind IN ('Ticket','Associated','Manual')),
    inserted_time     TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (billettholder_id) REFERENCES billettholdere(id)
);

-- +goose Down
PRAGMA foreign_keys = ON;

DROP TABLE billettholder_emails
