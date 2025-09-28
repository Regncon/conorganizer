-- +goose Up
PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS billettholdere_modified (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    ticket_type_id INTEGER NOT NULL,
    ticket_type TEXT NOT NULL,
    is_over_18 BOOLEAN NOT NULL,
    order_id INTEGER NOT NULL,
    ticket_id INTEGER NOT NULL UNIQUE,
    inserted_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose StatementBegin
INSERT INTO billettholdere_modified (
    id, first_name, last_name, ticket_type_id,
    ticket_type, is_over_18, order_id, ticket_id, inserted_time
)
SELECT
    id, first_name, last_name, CAST(ticket_category_id AS INTEGER),
    ticket_type, is_over_18, order_id, ticket_id, inserted_time
FROM billettholdere;
DROP TABLE billettholdere;
ALTER TABLE billettholdere_modified RENAME TO billettholdere;

PRAGMA foreign_keys=ON;
-- +goose StatementEnd


-- +goose Down
PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS billettholdere_modified (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    ticket_type TEXT NOT NULL,
    is_over_18 BOOLEAN NOT NULL,
    order_id INTEGER NOT NULL,
    ticket_id INTEGER NOT NULL UNIQUE,
    ticket_email TEXT NOT NULL,
    order_email TEXT NOT NULL,
    ticket_category_id TEXT NOT NULL,
    inserted_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (ticket_type) REFERENCES ticket_types(name)
);

-- +goose StatementBegin
INSERT INTO billettholdere_modified (
    id, first_name, last_name, ticket_category_id,
    ticket_type, is_over_18, order_id, ticket_id, inserted_time
)
SELECT
    id, first_name, last_name, CAST(ticket_type_id AS TEXT),
    ticket_type, is_over_18, order_id, ticket_id, inserted_time
FROM billettholdere;
DROP TABLE billettholdere;
ALTER TABLE billettholdere_modified RENAME TO billettholdere;

PRAGMA foreign_keys=ON;
-- +goose StatementEnd
