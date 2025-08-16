package formsubmission

import (
	"database/sql"
	"fmt"
	"github.com/Regncon/conorganizer/models"
	"log/slog"
)

func shouldShowStringValue(value string) string {
	if value != "" {
		return value
	}
	return ""
}

func shouldShowNumberValue(value int64) string {
	if value != 0 {
		return fmt.Sprintf("%d", value)
	}
	return ""
}

func GetEventById(userId string, eventID string, db *sql.DB, logger *slog.Logger) (*models.Event, error) {
	if userId == "" {
		logger.Error("Unauthorized", "User is not logged in")
		return nil, fmt.Errorf("unauthorized access")
	}

	/*userDbId, userDbIdErr := userctx.GetIdFromUserIdInDb(userId, db, logger)
	if userDbIdErr != nil {
		logger.Error("Failed to get user ID from database", "error", userDbIdErr)
		return nil, fmt.Errorf("failed to get user ID from database: %w", userDbIdErr)
	}
	*/

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
