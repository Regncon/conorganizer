package profilepage

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/service/userctx"
)

func GetEventsByUserId(userID string, db *sql.DB, logger *slog.Logger) []models.EventCardModel {
	baseLogger := logger
	logger = logger.With("component", "profile")
	var events []models.EventCardModel

	// Get events where event created id is the same as user
	userDbId, err := userctx.GetIdFromUserIdInDb(userID, db, baseLogger)
	if err != nil {
		logger.Error(fmt.Errorf("failed to get user database ID for user %q: %w", userID, err).Error())
		return events
	}

	// Query for events created by user
	eventsQuery := "SELECT id, title, intro, status, system, host_name, beginner_friendly, event_type, age_group, event_runtime, can_be_run_in_english FROM events WHERE host = ?"
	rows, eventsQueryErr := db.Query(eventsQuery, userDbId)
	if eventsQueryErr != nil {
		logger.Error(fmt.Errorf("failed to query events for user %q: %w", userID, eventsQueryErr).Error())
		return events
	}
	defer rows.Close()

	// Validate database query return
	for rows.Next() {
		var event models.EventCardModel
		if scanErr := rows.Scan(&event.Id, &event.Title, &event.Intro, &event.Status, &event.System, &event.HostName, &event.BeginnerFriendly, &event.EventType, &event.AgeGroup, &event.Runtime, &event.CanBeRunInEnglish); scanErr != nil {
			logger.Error(fmt.Errorf("failed to scan event row for user %q: %w", userID, scanErr).Error())
			return events
		}
		events = append(events, event)
	}

	return events
}
