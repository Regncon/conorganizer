-- +goose Up
ALTER TABLE puljer
ADD COLUMN rooms_published INTEGER NOT NULL DEFAULT 0 CHECK(rooms_published IN(0, 1));

-- +goose Down
ALTER TABLE puljer DROP COLUMN rooms_published;
