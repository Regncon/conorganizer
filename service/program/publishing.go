package program

import (
	"database/sql"
	"errors"
	"fmt"
)

// IsPublished reports the global program publication state.
func IsPublished(db *sql.DB) (bool, error) {
	const query = `
		SELECT is_published
		FROM program_publishing_state
		WHERE id = 1
	`

	var isPublished int
	if err := db.QueryRow(query).Scan(&isPublished); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("query program publishing state: %w", err)
	}

	return isPublished == 1, nil
}
