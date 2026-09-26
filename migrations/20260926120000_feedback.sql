-- +goose Up
CREATE TABLE feedback(
  id INTEGER PRIMARY KEY,
  category TEXT NOT NULL CHECK(category IN('website', 'convention', 'other')),
  message TEXT NOT NULL CHECK(length(message) BETWEEN 1 AND 2000),
  topics TEXT NOT NULL DEFAULT '[]' CHECK(json_valid(topics)),
  created_at TEXT NOT NULL DEFAULT(strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))
) STRICT;
CREATE INDEX idx_feedback_created_at ON feedback(created_at);

-- +goose Down
DROP INDEX idx_feedback_created_at;
DROP TABLE feedback;
