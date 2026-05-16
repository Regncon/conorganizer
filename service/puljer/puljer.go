package puljerService

import (
	"database/sql"
	"fmt"

	"github.com/Regncon/conorganizer/models"
)

func GetActivePuljeForEvent(eventID string, db *sql.DB) ([]models.PuljeRow, error) {
	const query = `
		SELECT p.id, p.name, p.status, p.start_at, p.end_at
		FROM puljer p
		JOIN relation_event_puljer ep ON p.id = ep.pulje_id
		WHERE ep.event_id = ? AND ep.is_in_pulje = TRUE AND ep.is_published = TRUE
		ORDER BY p.start_at ASC
	`

	rows, err := db.Query(query, eventID)
	if err != nil {
		return nil, fmt.Errorf("query active puljer for event %s: %w", eventID, err)
	}
	defer rows.Close()

	return scanPulje(rows)
}

func GetAllPuljer(db *sql.DB) ([]models.PuljeRow, error) {
	const query = `
		SELECT id, name, status, start_at, end_at
		FROM puljer
		ORDER BY start_at ASC
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("query puljer: %w", err)
	}
	defer rows.Close()

	return scanPulje(rows)
}

func BuildPulje(id models.Pulje, name string, status models.PuljeStatus, startAtRaw string, endAtRaw string) (models.PuljeRow, error) {
	startAt, err := models.ParsePuljeTime(startAtRaw)
	if err != nil {
		return models.PuljeRow{}, fmt.Errorf("parse pulje %s start_at %q: %w", id, startAtRaw, err)
	}

	endAt, err := models.ParsePuljeTime(endAtRaw)
	if err != nil {
		return models.PuljeRow{}, fmt.Errorf("parse pulje %s end_at %q: %w", id, endAtRaw, err)
	}

	return models.PuljeRow{
		ID:      id,
		Name:    name,
		Status:  status,
		StartAt: startAt,
		EndAt:   endAt,
	}, nil
}

func scanPulje(rows *sql.Rows) ([]models.PuljeRow, error) {
	puljer := make([]models.PuljeRow, 0)
	for rows.Next() {
		var (
			id             models.Pulje
			name           string
			status         models.PuljeStatus
			startAt, endAt string
		)

		if err := rows.Scan(&id, &name, &status, &startAt, &endAt); err != nil {
			return nil, fmt.Errorf("scan pulje row: %w", err)
		}

		pulje, err := BuildPulje(id, name, status, startAt, endAt)
		if err != nil {
			return nil, err
		}

		puljer = append(puljer, pulje)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate pulje rows: %w", err)
	}

	return puljer, nil
}
