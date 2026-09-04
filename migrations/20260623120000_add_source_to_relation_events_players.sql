-- +goose Up
ALTER TABLE relation_events_players
    ADD COLUMN source TEXT NOT NULL DEFAULT 'manual' CHECK (source IN ('manual','solver'));

-- +goose Down
ALTER TABLE relation_events_players DROP COLUMN source;
