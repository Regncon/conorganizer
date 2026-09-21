package program

import (
	"database/sql"
	"fmt"

	"github.com/Regncon/conorganizer/models"
)

func GetAnnouncedEvents(db *sql.DB) ([]models.EventCardModel, error) {
	const q = `
		SELECT
			e.id,
			e.title,
			e.intro,
			e.status,
			e.system,
			e.host_name,
			e.event_type,
			e.age_group,
			e.event_runtime,
			e.beginner_friendly,
			e.can_be_run_in_english,
			e.is_in_puljefordeling
		FROM events e
		WHERE e.status = ?
		ORDER BY e.title COLLATE NOCASE ASC, e.id ASC
	`

	rows, err := db.Query(q, models.EventStatusAnnounced)
	if err != nil {
		return nil, fmt.Errorf("query announced events alphabetically: %w", err)
	}
	defer rows.Close()

	events := make([]models.EventCardModel, 0)
	for rows.Next() {
		var event models.EventCardModel
		if err := rows.Scan(
			&event.Id,
			&event.Title,
			&event.Intro,
			&event.Status,
			&event.System,
			&event.HostName,
			&event.EventType,
			&event.AgeGroup,
			&event.Runtime,
			&event.BeginnerFriendly,
			&event.CanBeRunInEnglish,
			&event.IsInPuljefordeling,
		); err != nil {
			return nil, fmt.Errorf("scan announced event row: %w", err)
		}
		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate announced event rows: %w", err)
	}

	return events, nil
}

func getEventsByPulje(db *sql.DB) (map[models.Pulje][]models.EventCardModel, error) {
	const q = `
		SELECT
			e.id,
			e.title,
			e.intro,
			e.status,
			e.system,
			e.host_name,
			e.event_type,
			e.age_group,
			e.event_runtime,
			e.beginner_friendly,
			e.can_be_run_in_english,
			e.is_in_puljefordeling,

			e.pulje_id
		FROM v_events_by_pulje_active e
		ORDER BY e.pulje_start_at ASC, e.title COLLATE NOCASE ASC, e.id ASC
	`

	rows, err := db.Query(q)
	if err != nil {
		return nil, fmt.Errorf("querying events by pulje: %w", err)
	}
	defer rows.Close()

	out := make(map[models.Pulje][]models.EventCardModel)

	for rows.Next() {
		var ev models.EventCardModel
		var puljeID models.Pulje

		if err := rows.Scan(
			&ev.Id,
			&ev.Title,
			&ev.Intro,
			&ev.Status,
			&ev.System,
			&ev.HostName,
			&ev.EventType,
			&ev.AgeGroup,
			&ev.Runtime,
			&ev.BeginnerFriendly,
			&ev.CanBeRunInEnglish,
			&ev.IsInPuljefordeling,

			&puljeID,
		); err != nil {
			return nil, fmt.Errorf("scan event by pulje: %w", err)
		}

		out[puljeID] = append(out[puljeID], ev)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate events by pulje: %w", err)
	}
	return out, nil
}
