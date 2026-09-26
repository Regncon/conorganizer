// Package feedback stores anonymous feedback from users. Nothing identifying
// the sender (user id, email) is stored, on purpose.
package feedback

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/Regncon/conorganizer/models"
)

// MaxMessageLength is the maximum number of characters (runes) in a message
// after surrounding whitespace is trimmed.
const MaxMessageLength = 2000

// ErrInvalidFeedback is wrapped by Submit when the input is rejected.
var ErrInvalidFeedback = errors.New("invalid feedback")

// Entry is one stored piece of feedback.
type Entry struct {
	ID        int
	Category  Category
	Message   string
	CreatedAt models.DBDateTime
}

// Submit validates and stores a piece of feedback. The message is trimmed
// before validation. Invalid input returns an error wrapping ErrInvalidFeedback.
func Submit(db *sql.DB, category Category, message string) error {
	if !category.Valid() {
		return fmt.Errorf("%w: unknown category %q", ErrInvalidFeedback, category)
	}
	message = strings.TrimSpace(message)
	if message == "" {
		return fmt.Errorf("%w: message is empty", ErrInvalidFeedback)
	}
	if length := utf8.RuneCountInString(message); length > MaxMessageLength {
		return fmt.Errorf("%w: message is %d characters, max is %d", ErrInvalidFeedback, length, MaxMessageLength)
	}

	if _, err := db.Exec(`INSERT INTO feedback (category, message) VALUES (?, ?)`, string(category), message); err != nil {
		return fmt.Errorf("insert feedback: %w", err)
	}
	return nil
}

// List returns stored feedback, newest first. An empty category returns all.
func List(db *sql.DB, category Category) ([]Entry, error) {
	rows, err := db.Query(`
		SELECT id, category, message, created_at
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
		if err := rows.Scan(&entry.ID, &entry.Category, &entry.Message, &entry.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan feedback: %w", err)
		}
		entries = append(entries, entry)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate feedback: %w", err)
	}
	return entries, nil
}
