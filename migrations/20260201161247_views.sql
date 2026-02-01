-- +goose Up
CREATE VIEW IF NOT EXISTS
    get_billettholder_id AS
SELECT
    u.*,
    bu.*
FROM
    billettholdere_users AS bu
    LEFT JOIN users AS u ON u.id = bu.user_id;

-- +goose Down
DROP VIEW IF EXISTS get_billettholder_id;
