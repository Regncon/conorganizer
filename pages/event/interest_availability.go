package event

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Regncon/conorganizer/models"
)

func getProgramPublished(db *sql.DB) (bool, error) {
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

func canShowInterestControls(programPublished bool, puljerForEvent []models.PuljeRow) bool {
	return programPublished && len(puljerForEvent) > 0
}
