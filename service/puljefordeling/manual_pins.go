package puljefordeling

import (
	"database/sql"
	"fmt"

	"github.com/Regncon/conorganizer/models"
)

// AddManualSeat force-pins a participant into an event for the given pulje by
// writing a player seat tagged source='manual'. It removes the interest for the
// assigned event and pulje, while preserving other interests and open-registration
// seats. A participant holds at most one non-registration player seat per pulje,
// so moving them between ordinary events leaves a single pin.
func AddManualSeat(db *sql.DB, pulje models.Pulje, eventID string, billettholderID int) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("begin add manual seat tx: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	// Clear any prior ordinary player seat in this pulje first, so a move leaves a
	// single pin. Open-registration seats are independent and must survive.
	if _, err := tx.Exec(
		`DELETE FROM relation_events_players
		 WHERE billettholder_id = ? AND pulje_id = ? AND role = ? AND source IN (?, ?)`,
		billettholderID,
		string(pulje),
		models.EventPlayerRolePlayer,
		models.EventPlayerSourceManual,
		models.EventPlayerSourceSolver,
	); err != nil {
		return fmt.Errorf("clear prior seat (pulje=%s bh=%d): %w", pulje, billettholderID, err)
	}
	if _, err := tx.Exec(
		`INSERT INTO relation_events_players (event_id, pulje_id, billettholder_id, role, source)
		 VALUES (?, ?, ?, ?, ?)
		 ON CONFLICT(billettholder_id, event_id, pulje_id) DO UPDATE SET
			role = EXCLUDED.role,
			source = EXCLUDED.source,
			inserted_at = `+models.DBDateTimeNowSQL,
		eventID, string(pulje), billettholderID, models.EventPlayerRolePlayer, models.EventPlayerSourceManual,
	); err != nil {
		return fmt.Errorf("add manual seat (pulje=%s event=%s bh=%d): %w", pulje, eventID, billettholderID, err)
	}
	if _, err := tx.Exec(
		`DELETE FROM interests WHERE event_id = ? AND pulje_id = ? AND billettholder_id = ?`,
		eventID, string(pulje), billettholderID,
	); err != nil {
		return fmt.Errorf("remove interest for manual seat (pulje=%s event=%s bh=%d): %w", pulje, eventID, billettholderID, err)
	}
	return tx.Commit()
}

// RemoveManualSeat deletes an admin-pinned player seat (source='manual',
// role='Player') for the given pulje/event/participant. It only removes manual
// player pins — solver, registration, and GM rows are left untouched. The
// interest removed by the original manual assignment is not restored.
func RemoveManualSeat(db *sql.DB, pulje models.Pulje, eventID string, billettholderID int) error {
	const query = `
		DELETE FROM relation_events_players
		WHERE event_id = ? AND pulje_id = ? AND billettholder_id = ?
		  AND source = ? AND role = ?
	`
	if _, err := db.Exec(query, eventID, string(pulje), billettholderID, models.EventPlayerSourceManual, models.EventPlayerRolePlayer); err != nil {
		return fmt.Errorf("remove manual seat (pulje=%s event=%s bh=%d): %w", pulje, eventID, billettholderID, err)
	}
	return nil
}

// AddManualGM assigns a billettholder as the GM for one event and pulje. It
// deliberately leaves interests untouched; GM assignment is independent from
// the participant's stated preferences.
func AddManualGM(db *sql.DB, pulje models.Pulje, eventID string, billettholderID int) error {
	const query = `
		INSERT INTO relation_events_players (event_id, pulje_id, billettholder_id, role, source)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT(billettholder_id, event_id, pulje_id) DO UPDATE SET
			role = EXCLUDED.role,
			source = EXCLUDED.source,
			inserted_at = ` + models.DBDateTimeNowSQL
	if _, err := db.Exec(
		query,
		eventID,
		string(pulje),
		billettholderID,
		models.EventPlayerRoleGM,
		models.EventPlayerSourceManual,
	); err != nil {
		return fmt.Errorf("add manual GM (pulje=%s event=%s bh=%d): %w", pulje, eventID, billettholderID, err)
	}
	return nil
}

// RemoveManualGM removes exactly one admin-assigned GM. Player seats and
// interests are independent and remain untouched.
func RemoveManualGM(db *sql.DB, pulje models.Pulje, eventID string, billettholderID int) error {
	const query = `
		DELETE FROM relation_events_players
		WHERE event_id = ? AND pulje_id = ? AND billettholder_id = ?
		  AND source = ? AND role = ?
	`
	if _, err := db.Exec(
		query,
		eventID,
		string(pulje),
		billettholderID,
		models.EventPlayerSourceManual,
		models.EventPlayerRoleGM,
	); err != nil {
		return fmt.Errorf("remove manual GM (pulje=%s event=%s bh=%d): %w", pulje, eventID, billettholderID, err)
	}
	return nil
}

// AddFirstChoiceSeat assigns a manual player seat and records the event as the
// billettholder's first choice in that pulje. Independent registrations and
// interests for other events remain untouched.
func AddFirstChoiceSeat(db *sql.DB, pulje models.Pulje, eventID string, billettholderID int) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("begin add first-choice seat tx: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	if _, err := tx.Exec(
		`DELETE FROM relation_events_players
		 WHERE billettholder_id = ? AND pulje_id = ? AND role = ? AND source IN (?, ?)`,
		billettholderID,
		string(pulje),
		models.EventPlayerRolePlayer,
		models.EventPlayerSourceManual,
		models.EventPlayerSourceSolver,
	); err != nil {
		return fmt.Errorf("clear prior seat before first-choice assignment (pulje=%s bh=%d): %w", pulje, billettholderID, err)
	}

	if _, err := tx.Exec(
		`INSERT INTO relation_events_players (event_id, pulje_id, billettholder_id, role, source)
		 VALUES (?, ?, ?, ?, ?)
		 ON CONFLICT(billettholder_id, event_id, pulje_id) DO UPDATE SET
			role = EXCLUDED.role,
			source = EXCLUDED.source,
			inserted_at = `+models.DBDateTimeNowSQL,
		eventID,
		string(pulje),
		billettholderID,
		models.EventPlayerRolePlayer,
		models.EventPlayerSourceManual,
	); err != nil {
		return fmt.Errorf("add first-choice seat (pulje=%s event=%s bh=%d): %w", pulje, eventID, billettholderID, err)
	}

	if _, err := tx.Exec(
		`INSERT INTO interests (billettholder_id, event_id, pulje_id, interest_level)
		 VALUES (?, ?, ?, ?)
		 ON CONFLICT(billettholder_id, event_id, pulje_id) DO UPDATE SET
			interest_level = EXCLUDED.interest_level,
			updated_at = `+models.DBDateTimeNowSQL,
		billettholderID,
		eventID,
		string(pulje),
		models.InterestLevelHigh,
	); err != nil {
		return fmt.Errorf("record first-choice interest (pulje=%s event=%s bh=%d): %w", pulje, eventID, billettholderID, err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit first-choice assignment (pulje=%s event=%s bh=%d): %w", pulje, eventID, billettholderID, err)
	}
	return nil
}

// RemoveFirstChoiceSeat removes both sides of an admin-created first choice:
// the manual player seat and its matching high-interest row. Other interests
// and independent registrations remain untouched.
func RemoveFirstChoiceSeat(db *sql.DB, pulje models.Pulje, eventID string, billettholderID int) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("begin remove first-choice seat tx: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	result, err := tx.Exec(
		`DELETE FROM relation_events_players
		 WHERE event_id = ? AND pulje_id = ? AND billettholder_id = ?
		   AND source = ? AND role = ?`,
		eventID,
		string(pulje),
		billettholderID,
		models.EventPlayerSourceManual,
		models.EventPlayerRolePlayer,
	)
	if err != nil {
		return fmt.Errorf("remove first-choice seat (pulje=%s event=%s bh=%d): %w", pulje, eventID, billettholderID, err)
	}
	removedSeats, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("count removed first-choice seats (pulje=%s event=%s bh=%d): %w", pulje, eventID, billettholderID, err)
	}

	if removedSeats > 0 {
		if _, err := tx.Exec(
			`DELETE FROM interests
			 WHERE event_id = ? AND pulje_id = ? AND billettholder_id = ? AND interest_level = ?`,
			eventID,
			string(pulje),
			billettholderID,
			models.InterestLevelHigh,
		); err != nil {
			return fmt.Errorf("remove first-choice interest (pulje=%s event=%s bh=%d): %w", pulje, eventID, billettholderID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit first-choice removal (pulje=%s event=%s bh=%d): %w", pulje, eventID, billettholderID, err)
	}
	return nil
}
