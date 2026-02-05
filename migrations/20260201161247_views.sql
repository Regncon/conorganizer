-- +goose Up
-- Used by: view for user<->billettholder lookup
-- Relevant code: service/userctx/userctx.go:55 (users lookup), pages/event/event.go:266-269 (join users<->billettholdere_users)
CREATE VIEW IF NOT EXISTS
    v_get_user_billettholder AS
SELECT
    u.id AS user_db_id,
    u.user_id AS user_id,
    u.email AS user_email,
    u.is_admin AS user_is_admin,
    u.inserted_time AS user_inserted_time,
    bu.billettholder_id AS billettholder_id,
    bu.user_id AS billettholder_user_db_id,
    bu.inserted_time AS billettholder_user_inserted_time
FROM
    billettholdere_users AS bu
    LEFT JOIN users AS u ON u.id = bu.user_id;

-- Used by: pages/root/event_list.templ, pages/print-friendly/print-friendly-page.templ
-- Relevant code: pages/root/event_list.templ:40, pages/print-friendly/print-friendly-page.templ:93
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

-- Potential replacement for event summary queries
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
    can_be_run_in_english
FROM
    events;

-- Used by: billettholder + email lookups
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
    b.inserted_time AS billettholder_inserted_time,
    e.id AS email_id,
    e.email,
    e.kind,
    e.inserted_time AS email_inserted_time
FROM
    billettholdere AS b
    LEFT JOIN billettholder_emails AS e
        ON b.id = e.billettholder_id;

-- Used by: event pulje lookups
-- Relevant code: components/ticket_holder/ticket_holder.go:93, service/eventService/event_helpers.go:100, pages/event/event.go:258
CREATE VIEW IF NOT EXISTS
    v_event_puljer_active AS
SELECT
    ep.event_id,
    ep.pulje_id,
    p.name AS pulje_name,
    p.start_time AS pulje_start_time,
    p.end_time AS pulje_end_time,
    ep.is_active,
    ep.is_published
FROM
    event_puljer ep
    JOIN puljer p ON p.id = ep.pulje_id
WHERE
    ep.is_active = 1
    AND ep.is_published = 1;

-- +goose Down
DROP VIEW IF EXISTS v_get_user_billettholder;
DROP VIEW IF EXISTS v_events_by_pulje_active;
DROP VIEW IF EXISTS v_event_summary;
DROP VIEW IF EXISTS v_billettholder_emails;
DROP VIEW IF EXISTS v_event_puljer_active;
