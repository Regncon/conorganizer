package rooms

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Regncon/conorganizer/models"
)

// RoomsPublished reports whether the room assignment for a pulje is published to non-admin users.
func RoomsPublished(db *sql.DB, pulje models.Pulje) (bool, error) {
	var published bool

	err := db.QueryRow(`SELECT rooms_published FROM puljer WHERE id = ?`, pulje).Scan(&published)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("pulje %s not found", pulje)
		}
		return false, fmt.Errorf("query rooms published for pulje %s: %w", pulje, err)
	}

	return published, nil
}

// SetRoomsPublished sets whether the room assignment for a pulje is published to non-admin users.
func SetRoomsPublished(db *sql.DB, pulje models.Pulje, published bool) error {
	result, err := db.Exec(`UPDATE puljer SET rooms_published = ? WHERE id = ?`, published, pulje)
	if err != nil {
		return fmt.Errorf("update rooms published for pulje %s: %w", pulje, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected for pulje %s: %w", pulje, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("pulje %s not found", pulje)
	}

	return nil
}
