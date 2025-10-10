-- +goose Up
PRAGMA foreign_keys = OFF;

DROP TABLE IF EXISTS puljer_new;

CREATE TABLE puljer_new (
    id          TEXT    NOT NULL PRIMARY KEY,
    name        TEXT    NOT NULL,
    is_closed   INTEGER NOT NULL DEFAULT FALSE,
    is_published INTEGER NOT NULL DEFAULT FALSE,
    start_time  TEXT    NOT NULL,
    end_time    TEXT    NOT NULL
);

-- +goose StatementBegin
INSERT INTO puljer_new (id, name, is_closed, is_published, start_time, end_time)
SELECT id, name, 0, 0, start_time, end_time
FROM puljer;

DROP TABLE puljer;
ALTER TABLE puljer_new RENAME TO puljer;
-- +goose StatementEnd

PRAGMA foreign_keys = ON;
PRAGMA foreign_key_check;


-- +goose Down
PRAGMA foreign_keys = OFF;

DROP TABLE IF EXISTS puljer_old;

CREATE TABLE puljer_old (
    id         TEXT NOT NULL PRIMARY KEY,
    name       TEXT NOT NULL,
    start_time TEXT NOT NULL,
    end_time   TEXT NOT NULL
);

-- +goose StatementBegin
INSERT INTO puljer_old (id, name, start_time, end_time)
SELECT id, name, start_time, end_time
FROM puljer;

DROP TABLE puljer;
ALTER TABLE puljer_old RENAME TO puljer;
-- +goose StatementEnd

PRAGMA foreign_keys = ON;
PRAGMA foreign_key_check;
