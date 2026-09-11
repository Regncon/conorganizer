-- +goose Up
CREATE TABLE web_push_subscriptions(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    endpoint TEXT NOT NULL UNIQUE,
    p256dh TEXT NOT NULL,
    auth TEXT NOT NULL,
    created_at TEXT NOT NULL DEFAULT(strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT(strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
) STRICT;

CREATE INDEX web_push_subscriptions_user_idx ON web_push_subscriptions(user_id);

CREATE TABLE pulje_varsel_publications(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    pulje_id TEXT NOT NULL,
    revision INTEGER NOT NULL,
    created_at TEXT NOT NULL DEFAULT(strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    UNIQUE(pulje_id, revision),
    FOREIGN KEY(pulje_id) REFERENCES puljer(id) ON DELETE CASCADE ON UPDATE CASCADE
) STRICT;

CREATE TABLE pulje_varsel_results(
    publication_id INTEGER NOT NULL,
    billettholder_id INTEGER NOT NULL,
    result_fingerprint TEXT NOT NULL,
    result_json TEXT NOT NULL,
    payload_json TEXT NOT NULL,
    PRIMARY KEY(publication_id, billettholder_id),
    FOREIGN KEY(publication_id) REFERENCES pulje_varsel_publications(id) ON DELETE CASCADE,
    FOREIGN KEY(billettholder_id) REFERENCES billettholdere(id) ON DELETE CASCADE
) STRICT;

CREATE INDEX pulje_varsel_results_billettholder_idx ON pulje_varsel_results(billettholder_id, publication_id);

CREATE TABLE web_push_jobs(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    publication_id INTEGER NOT NULL,
    billettholder_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    subscription_id INTEGER NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending' CHECK(status IN('pending','processing','retrying','sent','canceled','exhausted')),
    attempts INTEGER NOT NULL DEFAULT 0,
    next_attempt_at TEXT NOT NULL DEFAULT(strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    lease_until TEXT,
    lease_token TEXT,
    last_error TEXT NOT NULL DEFAULT '',
    created_at TEXT NOT NULL DEFAULT(strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    sent_at TEXT,
    UNIQUE(publication_id, billettholder_id, subscription_id),
    FOREIGN KEY(publication_id, billettholder_id) REFERENCES pulje_varsel_results(publication_id, billettholder_id) ON DELETE CASCADE,
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
) STRICT;

CREATE INDEX web_push_jobs_ready_idx ON web_push_jobs(status, next_attempt_at, lease_until);

-- +goose Down
DROP INDEX web_push_jobs_ready_idx;
DROP TABLE web_push_jobs;
DROP INDEX pulje_varsel_results_billettholder_idx;
DROP TABLE pulje_varsel_results;
DROP TABLE pulje_varsel_publications;
DROP INDEX web_push_subscriptions_user_idx;
DROP TABLE web_push_subscriptions;
