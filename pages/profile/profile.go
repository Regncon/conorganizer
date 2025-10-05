package profilepage

import (
	"database/sql"
	"log/slog"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/service/userctx"
)

func GetEventsByUserId(userID string, db *sql.DB, logger *slog.Logger) []models.EventCardModel {
	var events []models.EventCardModel

	// Get events where event created id is the same as user
	userDbId, err := userctx.GetIdFromUserIdInDb(userID, db, logger)
	if err != nil {
		logger.Error("userDbIdErr", "error", err)
		return events
	}

	// Query for events created by user
	eventsQuery := "SELECT id, title, intro, status, system, host_name, beginner_friendly, event_type, age_group, event_runtime, can_be_run_in_english FROM events WHERE host = ?"
	rows, eventsQueryErr := db.Query(eventsQuery, userDbId)
	if eventsQueryErr != nil {
		logger.Error("Failed to query events", "user_id", userID, "error", eventsQueryErr)
		return events
	}
	defer rows.Close()

	// Validate database query return
	for rows.Next() {
		var event models.EventCardModel
		if scanErr := rows.Scan(&event.Id, &event.Title, &event.Intro, &event.Status, &event.System, &event.HostName, &event.BeginnerFriendly, &event.EventType, &event.AgeGroup, &event.Runtime, &event.CanBeRunInEnglish); scanErr != nil {
			logger.Error("Failed to scan event row", "user_id", userID, "error", scanErr)
			return events
		}
		events = append(events, event)
	}

	return events
}
