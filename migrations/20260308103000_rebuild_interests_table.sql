-- +goose Up
-- Rebuild interests to repair a potentially corrupted PK index.

DROP TABLE IF EXISTS interests_rebuild_tmp;

CREATE TABLE interests_rebuild_tmp (
    billettholder_id INTEGER NOT NULL,
    event_id TEXT NOT NULL,
    pulje_id TEXT NOT NULL,
    interest_level TEXT NOT NULL,
    inserted_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (billettholder_id, event_id, pulje_id),
    FOREIGN KEY (billettholder_id) REFERENCES billettholdere(id),
    FOREIGN KEY (event_id) REFERENCES events(id),
    FOREIGN KEY (pulje_id) REFERENCES puljer(id),
    FOREIGN KEY (interest_level) REFERENCES interest_levels(interest_level) ON UPDATE CASCADE
);

INSERT OR IGNORE INTO interests_rebuild_tmp (
    billettholder_id,
    event_id,
    pulje_id,
    interest_level,
    inserted_time
)
SELECT
    billettholder_id,
    event_id,
    pulje_id,
    interest_level,
    inserted_time
FROM interests
-- Only include rows with non-empty pulje_id to avoid violating the new PK constraint.
WHERE TRIM(pulje_id) <> '';

DROP TABLE interests;
ALTER TABLE interests_rebuild_tmp RENAME TO interests;

-- Add back previous migration index
CREATE INDEX IF NOT EXISTS `idx_interests_event_lookup`
ON `interests` (event_id, pulje_id, interest_level, billettholder_id);

-- Non-transaction fallback (manual, if a runner cannot use tx):
-- 1) Take a backup first.
-- 2) Validate with: PRAGMA integrity_check;

-- Irreversible data repair migration.
