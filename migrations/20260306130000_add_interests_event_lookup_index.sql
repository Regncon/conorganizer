-- +goose Up
CREATE INDEX IF NOT EXISTS `idx_interests_event_lookup`
ON `interests` (event_id, pulje_id, interest_level, billettholder_id);

-- +goose Down
DROP INDEX IF EXISTS `idx_interests_event_lookup`;
