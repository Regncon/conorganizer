package userProgram

import (
	"database/sql"
	"log/slog"

	billettholderService "github.com/Regncon/conorganizer/service/billettholder"
	"github.com/Regncon/conorganizer/service/requestctx"
)

func getAllEventsForUser(userInfo requestctx.UserRequestInfo, db *sql.DB, logger *slog.Logger) ([]UserEvent, error) {
	logger.Info("Fetching events for user", "userId", userInfo.Id)

	billettholdere, billettholderErr := billettholderService.GetBilettholdere(userInfo.Id, db, logger)
	if billettholderErr != nil {
		logger.Error("Failed to get billettholdere", "error", billettholderErr)
		return nil, billettholderErr
	}

	if len(billettholdere) == 0 {
		logger.Info("User has no billettholdere")
		return []UserEvent{}, nil
	}

	billettholderID := billettholdere[0].ID

	query := `
		SELECT
			ep.event_id,
			e.title,
			COALESCE(e.intro, '') as intro,
			COALESCE(e.description, '') as description,
			COALESCE(e.image_url, '') as image_url,
			COALESCE(e.host, 0) as host,
			e.event_type,
			ep.pulje_id,
			p.name as pulje_name,
			p.start_time,
			p.end_time
		FROM events_players ep
		JOIN events e ON ep.event_id = e.id
		JOIN puljer p ON ep.pulje_id = p.id
		WHERE ep.billettholder_id = ?
			AND ep.is_player = 1
		ORDER BY p.start_time ASC
	`

	rows, queryErr := db.Query(query, billettholderID)
	if queryErr != nil {
		logger.Error("Failed to query events", "error", queryErr)
		return nil, queryErr
	}
	defer rows.Close()

	var events []UserEvent

	for rows.Next() {
		var event UserEvent

		scanErr := rows.Scan(
			&event.EventID,
			&event.Title,
			&event.Intro,
			&event.Description,
			&event.ImageURL,
			&event.Host,
			&event.EventType,
			&event.PuljeID,
			&event.PuljeName,
			&event.StartTime,
			&event.EndTime,
		)

		if scanErr != nil {
			logger.Error("Failed to scan event", "error", scanErr)
			continue
		}
		logger.Info("Event scanned", "eventID", event.EventID, "title", event.Title, "pulje", event.PuljeName, "puljeID", event.PuljeID, "imageURL", event.ImageURL, "eventType", event.EventType)
		events = append(events, event)
	}

	if rowsErr := rows.Err(); rowsErr != nil {
		logger.Error("Error iterating rows", "error", rowsErr)
		return nil, rowsErr
	}

	logger.Info("Found events", "count", len(events))

	return events, nil
}
