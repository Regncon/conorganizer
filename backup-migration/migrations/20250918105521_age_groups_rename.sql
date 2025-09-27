-- +goose Up
PRAGMA foreign_keys=ON;
PRAGMA legacy_alter_table=OFF;

ALTER TABLE age_grups RENAME TO age_groups;

-- +goose Down
PRAGMA foreign_keys=ON;
PRAGMA legacy_alter_table=OFF;

ALTER TABLE age_groups RENAME TO age_grups;
