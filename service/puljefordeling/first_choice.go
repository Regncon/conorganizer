package puljefordeling

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Regncon/conorganizer/models"
)

var (
	ErrFirstChoiceInvalidInput          = errors.New("first-choice invalid input")
	ErrFirstChoiceMissingAssignment     = errors.New("first-choice missing assignment")
	ErrFirstChoiceGMAssignment          = errors.New("first-choice GM assignment")
	ErrFirstChoiceOtherPuljeFirstChoice = errors.New("first-choice already exists")
)

type FirstChoiceKey struct {
	BillettholderID int
	EventID         string
	PuljeID         string
}

type FirstChoiceStatus struct {
	HasCurrentPuljeFirstChoice bool
	HasOtherPuljeFirstChoice   bool
}

func GetFirstChoiceStatusesForEvent(db *sql.DB, eventID string) (map[FirstChoiceKey]FirstChoiceStatus, error) {
	const query = `
		WITH current_rows AS (
			SELECT billettholder_id, event_id, pulje_id
			FROM interests
			WHERE event_id = ?

			UNION

			SELECT billettholder_id, event_id, pulje_id
			FROM relation_events_players
			WHERE event_id = ?
		),
		qualifying_first_choices AS (
			SELECT ep.billettholder_id, ep.event_id, ep.pulje_id
			FROM relation_events_players ep
			JOIN interests i
				ON i.billettholder_id = ep.billettholder_id
				AND i.event_id = ep.event_id
				AND i.pulje_id = ep.pulje_id
			WHERE ep.role = ?
				AND i.interest_level = ?
		)
		SELECT
			cr.billettholder_id,
			cr.event_id,
			cr.pulje_id,
			EXISTS (
				SELECT 1
				FROM qualifying_first_choices q
				WHERE q.billettholder_id = cr.billettholder_id
					AND q.event_id = cr.event_id
					AND q.pulje_id = cr.pulje_id
			) AS has_current_pulje_first_choice,
			EXISTS (
				SELECT 1
				FROM qualifying_first_choices q
				WHERE q.billettholder_id = cr.billettholder_id
					AND NOT (
						q.event_id = cr.event_id
						AND q.pulje_id = cr.pulje_id
					)
			) AS has_other_pulje_first_choice
		FROM current_rows cr
	`

	rows, err := db.Query(query, eventID, eventID, models.EventPlayerRolePlayer, models.InterestLevelHigh)
	if err != nil {
		return nil, fmt.Errorf("query first-choice statuses for event %s: %w", eventID, err)
	}
	defer rows.Close()

	statuses := make(map[FirstChoiceKey]FirstChoiceStatus)
	for rows.Next() {
		var key FirstChoiceKey
		var current, other int
		if err := rows.Scan(&key.BillettholderID, &key.EventID, &key.PuljeID, &current, &other); err != nil {
			return nil, fmt.Errorf("scan first-choice status for event %s: %w", eventID, err)
		}
		statuses[key] = FirstChoiceStatus{
			HasCurrentPuljeFirstChoice: current != 0,
			HasOtherPuljeFirstChoice:   other != 0,
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate first-choice statuses for event %s: %w", eventID, err)
	}

	return statuses, nil
}

func SetAssignmentFirstChoice(db *sql.DB, eventID string, puljeID string, billettholderID int, enabled bool) error {
	if eventID == "" {
		return fmt.Errorf("%w: eventID is required", ErrFirstChoiceInvalidInput)
	}
	if puljeID == "" {
		return fmt.Errorf("%w: puljeID is required", ErrFirstChoiceInvalidInput)
	}
	if billettholderID <= 0 {
		return fmt.Errorf("%w: billettholderID must be positive", ErrFirstChoiceInvalidInput)
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("begin first-choice transaction: %w", err)
	}
	defer tx.Rollback()

	var role models.EventPlayerRole
	if err := tx.QueryRow(`
		SELECT role
		FROM relation_events_players
		WHERE event_id = ?
			AND pulje_id = ?
			AND billettholder_id = ?
	`, eventID, puljeID, billettholderID).Scan(&role); err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("%w: cannot update first-choice without an assignment", ErrFirstChoiceMissingAssignment)
		}
		return fmt.Errorf("query first-choice assignment role for billettholder_id=%d event_id=%s pulje_id=%s: %w", billettholderID, eventID, puljeID, err)
	}

	if role == models.EventPlayerRoleGM {
		return fmt.Errorf("%w: cannot update first-choice for GM assignment", ErrFirstChoiceGMAssignment)
	}

	if enabled {
		var hasOtherFirstChoice int
		if err := tx.QueryRow(`
			SELECT EXISTS (
				SELECT 1
				FROM relation_events_players ep
				JOIN interests i
					ON i.billettholder_id = ep.billettholder_id
					AND i.event_id = ep.event_id
					AND i.pulje_id = ep.pulje_id
				WHERE ep.billettholder_id = ?
					AND ep.role = ?
					AND i.interest_level = ?
					AND NOT (
						ep.event_id = ?
						AND ep.pulje_id = ?
					)
			)
		`, billettholderID, models.EventPlayerRolePlayer, models.InterestLevelHigh, eventID, puljeID).Scan(&hasOtherFirstChoice); err != nil {
			return fmt.Errorf("query other first-choice assignment for billettholder_id=%d pulje_id=%s: %w", billettholderID, puljeID, err)
		}
		if hasOtherFirstChoice != 0 {
			return fmt.Errorf("%w: cannot set first-choice while another first-choice already exists", ErrFirstChoiceOtherPuljeFirstChoice)
		}
	}

	level := models.InterestLevelMedium
	if enabled {
		level = models.InterestLevelHigh
	}

	upsertInterestQuery := `
		INSERT INTO interests(billettholder_id, event_id, pulje_id, interest_level)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(billettholder_id, event_id, pulje_id) DO UPDATE SET
			interest_level = excluded.interest_level,
			updated_at = ` + models.DBDateTimeNowSQL + `
	`
	if _, err := tx.Exec(upsertInterestQuery, billettholderID, eventID, puljeID, level); err != nil {
		return fmt.Errorf("upsert first-choice interest for billettholder_id=%d event_id=%s pulje_id=%s: %w", billettholderID, eventID, puljeID, err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit first-choice transaction: %w", err)
	}

	return nil
}
