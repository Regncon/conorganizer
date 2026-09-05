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
		 VALUES (?, ?, ?, ?, ?)`,
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
