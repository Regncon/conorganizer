// Package feedback stores anonymous feedback from users. Nothing identifying
// the sender (user id, email) is stored, on purpose.
package feedback

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/Regncon/conorganizer/models"
)

// MaxTextLength is the maximum number of characters (runes) in each text field
// after surrounding whitespace is trimmed.
const MaxTextLength = 1000

// Field names used in FieldError.
const (
	FieldCategory     = "category"
	FieldWentWell     = "wentWell"
	FieldCouldImprove = "couldImprove"
	FieldTopics       = "topics"
)

// Error codes used in FieldError.
const (
	CodeRequired     = "required"
	CodeTooLong      = "too_long"
	CodePersonalInfo = "personal_info"
	CodeInvalid      = "invalid"
)

// ErrInvalidFeedback is wrapped by Submit when the input is rejected.
var ErrInvalidFeedback = errors.New("invalid feedback")

// Submission is feedback as sent by a user.
type Submission struct {
	Category     Category
	WentWell     string
	CouldImprove string
	Topics       []Topic
}

// FieldError describes why one field of a Submission was rejected.
type FieldError struct {
	Field string
	Code  string
}

// Entry is one stored piece of feedback.
type Entry struct {
	ID           int
	Category     Category
	WentWell     string
	CouldImprove string
	Topics       []Topic
	CreatedAt    models.DBDateTime
}

// Validate returns every problem with the submission. Texts are trimmed before
// they are checked. An empty result means the submission can be stored.
func (s Submission) Validate() []FieldError {
	var errs []FieldError
	if !s.Category.Valid() {
		errs = append(errs, FieldError{Field: FieldCategory, Code: CodeInvalid})
	}

	wentWell := strings.TrimSpace(s.WentWell)
	couldImprove := strings.TrimSpace(s.CouldImprove)
	if wentWell == "" && couldImprove == "" {
		errs = append(errs, FieldError{Field: FieldWentWell, Code: CodeRequired})
	}
	errs = appendTextErrors(errs, FieldWentWell, wentWell)
	errs = appendTextErrors(errs, FieldCouldImprove, couldImprove)

	if !validTopics(s.Category, s.Topics) {
		errs = append(errs, FieldError{Field: FieldTopics, Code: CodeInvalid})
	}
	return errs
}

func appendTextErrors(errs []FieldError, field, text string) []FieldError {
	if utf8.RuneCountInString(text) > MaxTextLength {
		errs = append(errs, FieldError{Field: field, Code: CodeTooLong})
	}
	if ContainsPersonalInfo(text) {
		errs = append(errs, FieldError{Field: field, Code: CodePersonalInfo})
	}
	return errs
}

func validTopics(category Category, topics []Topic) bool {
	seen := make(map[Topic]bool, len(topics))
	for _, topic := range topics {
		if seen[topic] || !topic.belongsTo(category) {
			return false
		}
		seen[topic] = true
	}
	return true
}

// Submit validates and stores a piece of feedback with trimmed texts. Invalid
// input returns an error wrapping ErrInvalidFeedback and nothing is stored.
func Submit(db *sql.DB, s Submission) error {
	if errs := s.Validate(); len(errs) > 0 {
		return fmt.Errorf("%w: %v", ErrInvalidFeedback, errs)
	}

	topics := s.Topics
	if topics == nil {
		topics = []Topic{}
	}
	topicsJSON, err := json.Marshal(topics)
	if err != nil {
		return fmt.Errorf("encode feedback topics: %w", err)
	}

	if _, err := db.Exec(
		`INSERT INTO feedback (category, went_well, could_improve, topics) VALUES (?, ?, ?, ?)`,
		string(s.Category), strings.TrimSpace(s.WentWell), strings.TrimSpace(s.CouldImprove), string(topicsJSON),
	); err != nil {
		return fmt.Errorf("insert feedback: %w", err)
	}
	return nil
}

// List returns stored feedback, newest first. An empty category returns all.
func List(db *sql.DB, category Category) ([]Entry, error) {
	rows, err := db.Query(`
		SELECT id, category, went_well, could_improve, topics, created_at
		FROM feedback
		WHERE ? = '' OR category = ?
		ORDER BY created_at DESC, id DESC`,
		string(category), string(category),
	)
	if err != nil {
		return nil, fmt.Errorf("query feedback: %w", err)
	}
	defer rows.Close()

	entries := []Entry{}
	for rows.Next() {
		var entry Entry
		var topicsJSON string
		if err := rows.Scan(&entry.ID, &entry.Category, &entry.WentWell, &entry.CouldImprove, &topicsJSON, &entry.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan feedback: %w", err)
		}
		entry.Topics = []Topic{}
		if err := json.Unmarshal([]byte(topicsJSON), &entry.Topics); err != nil {
			return nil, fmt.Errorf("decode feedback topics for %d: %w", entry.ID, err)
		}
		entries = append(entries, entry)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate feedback: %w", err)
	}
	return entries, nil
}
