-- +goose Up
CREATE TABLE IF NOT EXISTS billettholdere (
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

-- +goose Down
CREATE TABLE IF NOT EXISTS billettholdere (
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
