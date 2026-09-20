-- +goose Up
ALTER TABLE events
ADD COLUMN is_in_puljefordeling INTEGER NOT NULL DEFAULT 0 CHECK (is_in_puljefordeling IN (0, 1));

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
    e.is_in_puljefordeling AS is_in_puljefordeling,
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

DROP VIEW IF EXISTS v_event_puljer_active;

CREATE VIEW v_event_puljer_active AS
SELECT
    ep.event_id,
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
    p.end_at AS pulje_end_at,
    ep.is_in_pulje,
    ep.is_published
FROM relation_event_puljer ep
JOIN puljer p ON p.id = ep.pulje_id
LEFT JOIN rooms r ON r.id = ep.room_id
WHERE ep.is_in_pulje = 1;

-- +goose Down
DROP VIEW IF EXISTS v_events_by_pulje_active;

ALTER TABLE events DROP COLUMN is_in_puljefordeling;

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

DROP VIEW IF EXISTS v_event_puljer_active;

CREATE VIEW v_event_puljer_active AS
SELECT
    ep.event_id,
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
    p.end_at AS pulje_end_at,
    ep.is_in_pulje,
    ep.is_published
FROM relation_event_puljer ep
JOIN puljer p ON p.id = ep.pulje_id
LEFT JOIN rooms r ON r.id = ep.room_id
WHERE ep.is_in_pulje = 1
  AND ep.is_published = 1;
