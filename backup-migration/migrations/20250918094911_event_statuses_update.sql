-- +goose Up
INSERT OR IGNORE INTO event_statuses (status) VALUES ('Forkastet');

UPDATE events
SET status = 'Forkastet'
WHERE status = 'Avvist';

DELETE FROM event_statuses
WHERE status = 'Avvist';

-- +goose Down
INSERT OR IGNORE INTO event_statuses (status) VALUES ('Avvist');

UPDATE events
SET status = 'Avvist'
WHERE status = 'Forkastet';

DELETE FROM event_statuses
WHERE status = 'Forkastet';
