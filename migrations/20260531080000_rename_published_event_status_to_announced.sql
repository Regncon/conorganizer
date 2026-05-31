-- +goose Up
INSERT INTO event_statuses(status)
VALUES ('Annonsert')
ON CONFLICT(status) DO NOTHING;

UPDATE events
SET status = 'Annonsert'
WHERE status IN ('Publisert', 'Godkjent');

DELETE FROM event_statuses
WHERE status = 'Publisert';

DROP VIEW IF EXISTS v_events_by_pulje_active;

CREATE VIEW v_events_by_pulje_active AS
SELECT
    e.id AS id,
    e.title,
    e.intro,
    e.description,
    e.system,
    e.event_type,
    e.age_group,
    e.event_runtime,
    e.host_name,
    e.user_id,
    e.email,
    e.phone_number,
    e.max_players,
    e.beginner_friendly,
    e.can_be_run_in_english,
    e.notes,
    e.status,
    e.created_at,
    ep.is_published AS is_published,
    ep.pulje_id,
    ep.room_id,
    r.room_number,
    r.name AS room_name,
    r.floor AS room_floor,
    r.max_concurrent_games AS room_max_concurrent_games,
    r.notes AS room_notes,
    r.is_disabled AS room_is_disabled,
    p.name AS pulje_name,
    p.start_at AS pulje_start_at,
    p.end_at AS pulje_end_at
FROM events e
INNER JOIN relation_event_puljer ep ON ep.event_id = e.id
INNER JOIN puljer p ON p.id = ep.pulje_id
LEFT JOIN rooms r ON r.id = ep.room_id
WHERE
    e.status = 'Annonsert'
    AND ep.is_in_pulje = 1;

-- +goose Down
INSERT INTO event_statuses(status)
VALUES
    ('Godkjent'),
    ('Publisert')
ON CONFLICT(status) DO NOTHING;

-- The Up migration folds both old public statuses into Annonsert.
-- On rollback, map them to Godkjent to preserve the old front-page behavior.
UPDATE events
SET status = 'Godkjent'
WHERE status = 'Annonsert';

DELETE FROM event_statuses
WHERE status = 'Annonsert';

DROP VIEW IF EXISTS v_events_by_pulje_active;

CREATE VIEW v_events_by_pulje_active AS
SELECT
    e.id AS id,
    e.title,
    e.intro,
    e.description,
    e.system,
    e.event_type,
    e.age_group,
    e.event_runtime,
    e.host_name,
    e.user_id,
    e.email,
    e.phone_number,
    e.max_players,
    e.beginner_friendly,
    e.can_be_run_in_english,
    e.notes,
    e.status,
    e.created_at,
    ep.is_published AS is_published,
    ep.pulje_id,
    ep.room_id,
    r.room_number,
    r.name AS room_name,
    r.floor AS room_floor,
    r.max_concurrent_games AS room_max_concurrent_games,
    r.notes AS room_notes,
    r.is_disabled AS room_is_disabled,
    p.name AS pulje_name,
    p.start_at AS pulje_start_at,
    p.end_at AS pulje_end_at
FROM events e
INNER JOIN relation_event_puljer ep ON ep.event_id = e.id
INNER JOIN puljer p ON p.id = ep.pulje_id
LEFT JOIN rooms r ON r.id = ep.room_id
WHERE
    e.status = 'Godkjent'
    AND ep.is_in_pulje = 1;
