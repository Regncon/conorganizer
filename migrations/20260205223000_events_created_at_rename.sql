-- +goose Up
DROP VIEW IF EXISTS v_events_by_pulje_active;

PRAGMA foreign_keys = OFF;
ALTER TABLE events RENAME COLUMN inserted_time TO created_at;
PRAGMA foreign_keys = ON;

-- Used by: v_events_by_pulje_active view consumers
-- Related code updates: pages/root/event_list.templ, pages/print-friendly/print-friendly-page.templ,
-- pages/admin/approval/approval_page.templ, service/eventService/getPreviousNext.go,
-- service/eventService/previous_next_test.go, models/event-model.go, schema.sql, initialize.sql
CREATE VIEW IF NOT EXISTS
    v_events_by_pulje_active AS
SELECT
    e.id AS id,
    e.title,
    e.intro,
    e.description,
    e.image_url,
    e.system,
    e.event_type,
    e.age_group,
    e.event_runtime,
    e.host_name,
    e.host,
    e.email,
    e.phone_number,
    e.pulje_name AS event_pulje_name,
    e.max_players,
    e.beginner_friendly,
    e.can_be_run_in_english,
    e.notes,
    e.status,
    e.created_at,
    ep.is_published AS is_published,
    ep.pulje_id,
    p.name AS pulje_name,
    p.start_time AS pulje_start_time,
    p.end_time AS pulje_end_time
FROM
    events e
    INNER JOIN event_puljer ep ON ep.event_id = e.id
    INNER JOIN puljer p ON p.id = ep.pulje_id
WHERE
    e.status = 'Godkjent'
    AND ep.is_active = 1;

-- +goose Down
DROP VIEW IF EXISTS v_events_by_pulje_active;

PRAGMA foreign_keys = OFF;
ALTER TABLE events RENAME COLUMN created_at TO inserted_time;
PRAGMA foreign_keys = ON;

-- Used by: v_events_by_pulje_active view consumers
-- Related code updates: pages/root/event_list.templ, pages/print-friendly/print-friendly-page.templ,
-- pages/admin/approval/approval_page.templ, service/eventService/getPreviousNext.go,
-- service/eventService/previous_next_test.go, models/event-model.go, schema.sql, initialize.sql
CREATE VIEW IF NOT EXISTS
    v_events_by_pulje_active AS
SELECT
    e.id AS id,
    e.title,
    e.intro,
    e.description,
    e.image_url,
    e.system,
    e.event_type,
    e.age_group,
    e.event_runtime,
    e.host_name,
    e.host,
    e.email,
    e.phone_number,
    e.pulje_name AS event_pulje_name,
    e.max_players,
    e.beginner_friendly,
    e.can_be_run_in_english,
    e.notes,
    e.status,
    e.inserted_time,
    ep.is_published AS is_published,
    ep.pulje_id,
    p.name AS pulje_name,
    p.start_time AS pulje_start_time,
    p.end_time AS pulje_end_time
FROM
    events e
    INNER JOIN event_puljer ep ON ep.event_id = e.id
    INNER JOIN puljer p ON p.id = ep.pulje_id
WHERE
    e.status = 'Godkjent'
    AND ep.is_active = 1;
