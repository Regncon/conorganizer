-- +goose Up
PRAGMA foreign_keys = OFF;
CREATE TABLE puljer_modified (
    id TEXT NOT NULL PRIMARY KEY,
    name TEXT NOT NULL,
    is_closed BOOLEAN NOT NULL DEFAULT FALSE,
    is_published BOOLEAN NOT NULL DEFAULT FALSE,
    start_time DATE NOT NULL,
    end_time DATE NOT NULL
);

-- +goose StatementBegin
INSERT INTO puljer_modified (id, name, start_time, end_time)
SELECT id, name, start_time, end_time
FROM puljer;
DROP TABLE puljer;
ALTER TABLE puljer_modified RENAME TO puljer;
-- +goose StatementEnd

PRAGMA foreign_keys = ON;
PRAGMA foreign_key_check


-- +goose Down
PRAGMA foreign_keys = OFF;
CREATE TABLE puljer_modified (
    id TEXT NOT NULL PRIMARY KEY,
    name TEXT NOT NULL,
    start_time DATE NOT NULL,
    end_time DATE NOT NULL
);

-- +goose StatementBegin
INSERT INTO puljer_modified (id, name, start_time, end_time)
SELECT id, name, start_time, end_time
FROM puljer;
DROP TABLE puljer;
ALTER TABLE puljer_modified RENAME TO puljer;
-- +goose StatementEnd

PRAGMA foreign_keys = ON;
PRAGMA foreign_key_check
