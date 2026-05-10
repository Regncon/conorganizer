package eventservice

import (
	"database/sql"
	"fmt"
	"github.com/Regncon/conorganizer/models"
	"log/slog"
)

func GetEventById(eventID string, db *sql.DB, logger *slog.Logger) (*models.Event, error) {
	query := `
            SELECT
                id,
                title,
                intro,
                description,
                system,
                event_type,
                age_group,
                event_runtime,
                host_name,
                user_id,
                email,
                phone_number,
                max_players,
                beginner_friendly,
                can_be_run_in_english,
                notes,
                status,
                created_at,
                updated_at,
                created_by_id,
                updated_by_id,
                status_changed_by_id,
                status_changed_at,
                status_changed_action
            FROM events WHERE id = ?
            `
	row := db.QueryRow(query, eventID)

	var event models.Event
	if err := row.Scan(
		&event.ID,
		&event.Title,
		&event.Intro,
		&event.Description,
		&event.System,
		&event.EventType,
		&event.AgeGroup,
		&event.Runtime,
		&event.HostName,
		&event.UserID,
		&event.Email,
		&event.PhoneNumber,
		&event.MaxPlayers,
		&event.BeginnerFriendly,
		&event.CanBeRunInEnglish,
		&event.Notes,
		&event.Status,
		&event.CreatedAt,
		&event.UpdatedAt,
		&event.CreatedByID,
		&event.UpdatedByID,
		&event.StatusChangedByID,
		&event.StatusChangedAt,
		&event.StatusChangedAction,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No event found
		}
		fmt.Printf("Error scanning event: %v\n", err)
		return nil, err
	}
	return &event, nil
}

func GetPujerForEvent(
	eventID string,
	db *sql.DB,
	logger *slog.Logger) ([]models.PuljeRow, error) {
	/*
			   CREATE TABLE
			       relation_event_puljer (
			           event_id TEXT NOT NULL,
			           pulje_id TEXT NOT NULL,
			           is_in_pulje BOOLEAN NOT NULL DEFAULT TRUE,
			           is_published BOOLEAN NOT NULL DEFAULT FALSE,
			           room TEXT DEFAULT '',
			           PRIMARY KEY (event_id, pulje_id),
			           FOREIGN KEY (event_id) REFERENCES events (id) ON DELETE CASCADE,
			           FOREIGN KEY (pulje_id) REFERENCES puljer (id) ON UPDATE CASCADE
			       );
		CREATE TABLE
		    puljer (
		        id TEXT NOT NULL PRIMARY KEY,
		        name TEXT NOT NULL,
		        start_at DATE NOT NULL,
		        end_at DATE NOT NULL
		    );

		type PuljeRow struct {
			ID      Pulje     `json:"id"`
			Name    string    `json:"name"`
			StartAt time.Time `json:"start_at"`
			EndAt   time.Time `json:"end_at"`
		}
	*/
	query := `SELECT p.id, p.name, p.start_at, p.end_at
            FROM puljer p
            JOIN relation_event_puljer ep ON p.id = ep.pulje_id
            WHERE ep.event_id = ? AND ep.is_in_pulje = TRUE AND ep.is_published = TRUE
            ORDER BY p.start_at ASC
            `

	rows, err := db.Query(query, eventID)
	if err != nil {
		logger.Error("Error querying puljer for event", slog.String("eventID", eventID), slog.String("error", err.Error()))
		return nil, err
	}
	defer rows.Close()
	var puljer []models.PuljeRow
	for rows.Next() {
		var pulje models.PuljeRow
		if err := rows.Scan(&pulje.ID, &pulje.Name, &pulje.StartAt, &pulje.EndAt); err != nil {
			logger.Error("Error scanning pulje row", slog.String("eventID", eventID), slog.String("error", err.Error()))
			return nil, err
		}
		puljer = append(puljer, pulje)
	}
	if err := rows.Err(); err != nil {
		logger.Error("Error iterating over pulje rows", slog.String("eventID", eventID), slog.String("error", err.Error()))
		return nil, err
	}
	return puljer, nil
}
