-- +goose Up
-- Existing room notes were only shown to admins, so they become admin notes.
ALTER TABLE rooms RENAME COLUMN notes TO admin_notes;
ALTER TABLE rooms ADD COLUMN public_notes TEXT NOT NULL DEFAULT '';

-- +goose Down
ALTER TABLE rooms DROP COLUMN public_notes;
ALTER TABLE rooms RENAME COLUMN admin_notes TO notes;
