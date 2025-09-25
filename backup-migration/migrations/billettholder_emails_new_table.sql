-- +goose Up
CREATE TABLE IF NOT EXISTS billettholder_emails (
    id                INTEGER PRIMARY KEY AUTOINCREMENT,
    billettholder_id  INTEGER NOT NULL,
    email             TEXT    NOT NULL COLLATE NOCASE,
    kind              TEXT    NOT NULL CHECK (kind IN ('Ticket','Associated','Manual')),
    inserted_time     TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE billettholder_emails
