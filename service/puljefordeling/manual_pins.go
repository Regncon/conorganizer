package puljefordeling

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Regncon/conorganizer/models"
)

// ErrInterestNotFound means the quick assignment was based on an interest that
// disappeared before the admin action reached the server.
var ErrInterestNotFound = errors.New("matching interest does not exist")

// AddManualSeat force-pins a participant into an event for the given pulje by
// writing a player seat tagged source='manual'. Interests and open-registration
// seats remain untouched. A participant holds at most one non-registration
// player seat per pulje, so moving them between ordinary events leaves a single
// pin.
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
	return tx.Commit()
}

// AddManualSeatFromInterest force-pins a participant into an event while
// preserving the interest that led the admin to choose them. This is the quick
// backup-assignment workflow in Puljefordeling: a retained high interest can
// satisfy the participant's first choice, while medium and low interests keep
// their ordinary solver value. As with other manual seats, moving a participant
// replaces their previous non-registration player seat in the same pulje.
func AddManualSeatFromInterest(db *sql.DB, pulje models.Pulje, eventID string, billettholderID int) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("begin add interest-backed manual seat tx: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	var interestExists bool
	if err := tx.QueryRow(
		`SELECT EXISTS (
			SELECT 1
			FROM interests
			WHERE billettholder_id = ? AND event_id = ? AND pulje_id = ?
		)`,
		billettholderID,
		eventID,
		string(pulje),
	).Scan(&interestExists); err != nil {
		return fmt.Errorf("check interest for manual seat (pulje=%s event=%s bh=%d): %w", pulje, eventID, billettholderID, err)
	}
	if !interestExists {
		return fmt.Errorf("add interest-backed manual seat (pulje=%s event=%s bh=%d): %w", pulje, eventID, billettholderID, ErrInterestNotFound)
	}

	if _, err := tx.Exec(
		`DELETE FROM relation_events_players
		 WHERE billettholder_id = ? AND pulje_id = ? AND role = ? AND source IN (?, ?)`,
		billettholderID,
		string(pulje),
		models.EventPlayerRolePlayer,
		models.EventPlayerSourceManual,
		models.EventPlayerSourceSolver,
	); err != nil {
		return fmt.Errorf("clear prior seat for interest-backed assignment (pulje=%s bh=%d): %w", pulje, billettholderID, err)
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
		return fmt.Errorf("add interest-backed manual seat (pulje=%s event=%s bh=%d): %w", pulje, eventID, billettholderID, err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit interest-backed manual seat (pulje=%s event=%s bh=%d): %w", pulje, eventID, billettholderID, err)
	}
	return nil
}

// RemoveManualSeat deletes an admin-pinned player seat (source='manual',
// role='Player') for the given pulje/event/participant. It only removes manual
// player pins — solver, registration, GM, and interest rows are left untouched.
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
