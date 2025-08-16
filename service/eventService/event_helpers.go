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
                image_url,
                system,
                event_type,
                age_group,
                event_runtime,
                host_name,
                host,
                email,
                phone_number,
                pulje_name,
                max_players,
                beginner_friendly,
                can_be_run_in_english,
                notes,
                status
            FROM events WHERE id = ?
            `
	row := db.QueryRow(query, eventID)

	var event models.Event
	if err := row.Scan(
		&event.ID,
		&event.Title,
		&event.Intro,
		&event.Description,
		&event.ImageURL,
		&event.System,
		&event.EventType,
		&event.AgeGroup,
		&event.Runtime,
		&event.HostName,
		&event.Host,
		&event.Email,
		&event.PhoneNumber,
		&event.PuljeName,
		&event.MaxPlayers,
		&event.BeginnerFriendly,
		&event.CanBeRunInEnglish,
		&event.Notes,
		&event.Status,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No event found
		}
		fmt.Printf("Error scanning event: %v\n", err)
		return nil, err
	}
	return &event, nil
}

/*
func GetEventByID(id string, db *sql.DB) (*models.Event, error) {
	query := `
            SELECT
                id,
                title,
                description,
                image_url,
                system,
                host_name,
                host,
                email,
                phone_number,
                pulje_name,
                max_players,
                beginner_friendly,
                can_be_run_in_english,
                status
            FROM events WHERE id = ? AND status = ?
            `
	row := db.QueryRow(query, id, models.EventStatusPublished)

	var event models.Event
	if err := row.Scan(
		&event.ID,
		&event.Title,
		&event.Description,
		&event.ImageURL,
		&event.System,
		&event.HostName,
		&event.Host,
		&event.Email,
		&event.PhoneNumber,
		&event.PuljeName,
		&event.MaxPlayers,
		&event.BeginnerFriendly,
		&event.CanBeRunInEnglish,
		&event.Status,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No event found
		}
		return nil, err
	}

	return &event, nil
}
*/
