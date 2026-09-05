-- +goose Up
ALTER TABLE events
    ADD COLUMN is_open_registration INTEGER NOT NULL DEFAULT 0
    CHECK (is_open_registration IN (0, 1));

-- +goose Down
ALTER TABLE events DROP COLUMN is_open_registration;
