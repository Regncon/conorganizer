-- +goose Up
PRAGMA foreign_keys = ON;

DROP TABLE ticket_types;

-- +goose Down
CREATE TABLE IF NOT EXISTS ticket_types (
    name TEXT PRIMARY KEY
);
