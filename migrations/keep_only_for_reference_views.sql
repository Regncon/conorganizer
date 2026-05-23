-- +goose Up
DROP VIEW IF EXISTS v_get_user_billettholder;
DROP VIEW IF EXISTS v_events_by_pulje_active;
DROP VIEW IF EXISTS v_event_summary;
DROP VIEW IF EXISTS v_billettholder_emails;
DROP VIEW IF EXISTS v_event_puljer_active;

-- Used by: view for external user<->billettholder lookup.
-- Relevant code: service/userctx/userctx.go:55 (users lookup), pages/event/event.go:266-269 (join users<->relation_billettholdere_users)
CREATE VIEW IF NOT EXISTS
    v_get_user_billettholder AS
SELECT
    u.id AS user_db_id,
    u.external_id AS external_id,
    u.email AS user_email,
    u.is_admin AS user_is_admin,
    u.inserted_at AS user_inserted_at,
    bu.billettholder_id AS billettholder_id,
    bu.user_id AS billettholder_user_db_id,
    bu.inserted_at AS billettholder_user_inserted_at
FROM
    relation_billettholdere_users AS bu
    LEFT JOIN users AS u ON u.id = bu.user_id;

-- Used by: approved event listings grouped by active pulje.
-- Relevant code: pages/root/event_list.templ:40, pages/print-friendly/print-friendly-page.templ:93
CREATE VIEW IF NOT EXISTS
    v_events_by_pulje_active AS
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
FROM
    events e
    INNER JOIN relation_event_puljer ep ON ep.event_id = e.id
    INNER JOIN puljer p ON p.id = ep.pulje_id
    LEFT JOIN rooms r ON r.id = ep.room_id
WHERE
    e.status = 'Godkjent'
    AND ep.is_in_pulje = 1;

-- Used by: compact event summary queries.
-- Relevant code: pages/admin/approval/approval_page_templ.go:41, pages/myprofile/myevents/myevents_page_templ.go:29, pages/profile/profile.go:22
CREATE VIEW IF NOT EXISTS
    v_event_summary AS
SELECT
    id,
    title,
    intro,
    status,
    system,
    host_name,
    beginner_friendly,
    event_type,
    age_group,
    event_runtime,
    can_be_run_in_english,
    created_at,
    updated_at
FROM
    events;

-- Used by: billettholder + email lookups.
-- Relevant code: components/ticket_holder/ticket_holder.go:25, service/billettholder/billettholder.go:21, service/billettholder/billettholder.go:35, service/checkIn/assign.go:57, service/checkIn/assign.go:133
CREATE VIEW IF NOT EXISTS
    v_billettholder_emails AS
SELECT
    b.id AS billettholder_id,
    b.first_name,
    b.last_name,
    b.ticket_type_id,
    b.ticket_type,
    b.is_over_18,
    b.order_id,
    b.ticket_id,
    b.created_at AS billettholder_created_at,
    b.updated_at AS billettholder_updated_at,
    e.id AS email_id,
    e.email,
    e.kind,
    e.created_at AS email_created_at,
    e.updated_at AS email_updated_at
FROM
    billettholdere AS b
    LEFT JOIN relation_billettholder_emails AS e ON b.id = e.billettholder_id;

-- Used by: published pulje lookups for an event.
-- Relevant code: components/ticket_holder/ticket_holder.go:93, service/eventService/event_helpers.go:100, pages/event/event.go:258
CREATE VIEW IF NOT EXISTS
    v_event_puljer_active AS
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
FROM
    relation_event_puljer ep
    JOIN puljer p ON p.id = ep.pulje_id
    LEFT JOIN rooms r ON r.id = ep.room_id
WHERE
    ep.is_in_pulje = 1
    AND ep.is_published = 1;

-- +goose Down
DROP VIEW IF EXISTS v_get_user_billettholder;
DROP VIEW IF EXISTS v_events_by_pulje_active;
DROP VIEW IF EXISTS v_event_summary;
DROP VIEW IF EXISTS v_billettholder_emails;
DROP VIEW IF EXISTS v_event_puljer_active;
