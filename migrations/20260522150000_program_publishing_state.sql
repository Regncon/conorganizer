-- +goose Up
CREATE TABLE program_publishing_state(
    id INTEGER NOT NULL PRIMARY KEY CHECK(id = 1),
    is_published INTEGER NOT NULL DEFAULT 0 CHECK(is_published IN(0, 1))
) STRICT;

INSERT INTO program_publishing_state(id, is_published)
VALUES (1, 0);

-- +goose Down
DROP TABLE program_publishing_state;
