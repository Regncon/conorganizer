-- +goose Up
CREATE TABLE feedback(
  id INTEGER PRIMARY KEY,
  category TEXT NOT NULL CHECK(category IN('website', 'convention', 'other')),
  went_well TEXT NOT NULL DEFAULT '' CHECK(length(went_well) <= 1000),
  could_improve TEXT NOT NULL DEFAULT '' CHECK(length(could_improve) <= 1000),
  topics TEXT NOT NULL DEFAULT '[]' CHECK(json_valid(topics)),
  created_at TEXT NOT NULL DEFAULT(strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
  CHECK(went_well <> '' OR could_improve <> '')
) STRICT;
CREATE INDEX idx_feedback_created_at ON feedback(created_at);

-- +goose Down
DROP INDEX idx_feedback_created_at;
DROP TABLE feedback;
