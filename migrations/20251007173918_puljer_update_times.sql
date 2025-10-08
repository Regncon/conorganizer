-- +goose Up
PRAGMA foreign_keys = ON;

-- +goose StatementBegin
UPDATE puljer
SET end_time = '2025-10-10T23:00:00Z'
WHERE id = 'FredagKveld' AND end_time = '2025-10-10T22:00:00Z';

UPDATE puljer
SET end_time = '2025-10-11T23:00:00Z'
WHERE id = 'LordagKveld' AND end_time = '2025-10-11T22:00:00Z';
-- +goose StatementEnd


-- +goose Down
PRAGMA foreign_keys = ON;

-- +goose StatementBegin
UPDATE puljer
SET end_time = '2025-10-10T22:00:00Z'
WHERE id = 'FredagKveld' AND end_time = '2025-10-10T23:00:00Z';

UPDATE puljer
SET end_time = '2025-10-11T22:00:00Z'
WHERE id = 'LordagKveld' AND end_time = '2025-10-11T23:00:00Z';
-- +goose StatementEnd
