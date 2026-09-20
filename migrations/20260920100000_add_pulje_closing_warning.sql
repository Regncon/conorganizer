-- +goose Up
ALTER TABLE puljer
ADD COLUMN closing_warning_active INTEGER NOT NULL DEFAULT 0 CHECK(closing_warning_active IN(0, 1));

-- +goose Down
ALTER TABLE puljer DROP COLUMN closing_warning_active;
