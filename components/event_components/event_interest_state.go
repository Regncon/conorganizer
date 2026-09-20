package event_components

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/service/eventimage"
)

type AssignedEvent struct {
	Title    string
	System   string
	Role     string
	Source   string
	EventID  string
	ImageURL string
}

type InterestNoticeState struct {
	ShowAssignedEvent bool
	ShowUnder18       bool
	CanChooseInterest bool
	AssignedEvent     AssignedEvent
}

func LoadInterestNoticeState(billettholderID int, puljeID string, ageGroup models.AgeGroup, eventImageDir *string, db *sql.DB) (InterestNoticeState, error) {
	if billettholderID <= 0 {
		return InterestNoticeState{}, nil
	}
	var isOver18 bool
	if err := db.QueryRow("SELECT is_over_18 FROM billettholdere WHERE id = ?1", billettholderID).Scan(&isOver18); err != nil {
		return InterestNoticeState{}, fmt.Errorf("failed to fetch billettholder age: %w", err)
	}
	assignedEvent, err := getAssignedEventForBillettholder(billettholderID, puljeID, db)
	if err != nil {
		return InterestNoticeState{}, err
	}
	if assignedEvent.EventID != "" {
		assignedEvent.ImageURL = eventimage.GetEventImageUrl(assignedEvent.EventID, "banner", eventImageDir)
	}
	return InterestNoticeState{
		ShowAssignedEvent: assignedEvent.EventID != "",
		ShowUnder18:       !isOver18 && ageGroup == models.AgeGroupAdultsOnly,
		CanChooseInterest: isOver18 && assignedEvent.EventID == "",
		AssignedEvent:     assignedEvent,
	}, nil
}

func getAssignedEventForBillettholder(billettholderID int, puljeID string, db *sql.DB) (AssignedEvent, error) {
	query := `
		SELECT e.title, e.system, bp.role, bp.source, e.id
		FROM relation_events_players bp
		JOIN events e ON e.id = bp.event_id
		WHERE bp.billettholder_id = ?1 AND bp.pulje_id = ?2
			AND (bp.source = 'manual' OR bp.role = 'GM')
		ORDER BY CASE WHEN bp.role = 'GM' THEN 0 ELSE 1 END, e.id
		LIMIT 1
	`
	var assignedEvent AssignedEvent
	err := db.QueryRow(query, billettholderID, puljeID).Scan(
		&assignedEvent.Title, &assignedEvent.System, &assignedEvent.Role, &assignedEvent.Source, &assignedEvent.EventID,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return AssignedEvent{}, nil
	}
	if err != nil {
		return AssignedEvent{}, fmt.Errorf("failed to fetch assigned event for billettholder: %w", err)
	}
	return assignedEvent, nil
}

func assignedEventHeading(role string) string {
	if role == "GM" {
		return "Du er arrangør for et arrangement i denne puljen."
	}
	return "Du har allerede blitt tildelt denne puljen."
}
