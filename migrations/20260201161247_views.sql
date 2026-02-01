-- +goose Up
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

-- +goose Down
DROP VIEW IF EXISTS v_get_user_billettholder;

