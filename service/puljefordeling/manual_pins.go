package puljefordeling

import (
	"database/sql"
	"fmt"

	"github.com/Regncon/conorganizer/models"
)

// AddManualSeat force-pins a participant into an event for the given pulje by
// writing a player seat tagged source='manual'. It deliberately does NOT touch
// the participant's interests: the pin forces and locks the placement on its own
// (the solver honours manual seats), and removing the pin reverts the player to
// pure emulation based on their real interests. A participant holds at most one
// player seat per pulje, so moving them between events leaves a single pin.
func AddManualSeat(db *sql.DB, pulje models.Pulje, eventID string, billettholderID int) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("begin add manual seat tx: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	if err := lagreTildeling(tx, Tildelingsvalg{
		PuljeID:         pulje,
		EventID:         eventID,
		BillettholderID: billettholderID,
		Role:            models.EventPlayerRolePlayer,
	}); err != nil {
		return fmt.Errorf("add manual seat (pulje=%s event=%s bh=%d): %w", pulje, eventID, billettholderID, err)
	}
	return tx.Commit()
}

// RemoveManualSeat deletes an admin-pinned player seat (source='manual',
// role='Player') for the given pulje/event/participant. It only removes manual
// player pins — other solver seats and GM rows are left untouched. Removing the
// pin does not touch the player's interest, so the emulation seats them again by
// simulation.
//
// Once the pulje's distribution has been saved, the player's new solver seat is
// saved right away as well. A manual move already replaced their saved seat, so
// without this, undoing the move would leave them without a saved seat until the
// next "Lagre fordeling". Seat changes it causes for others stay unsaved.
func RemoveManualSeat(db *sql.DB, pulje models.Pulje, eventID string, billettholderID int) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("begin remove manual seat: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	if err := removeManualSeat(tx, pulje, eventID, billettholderID); err != nil {
		return err
	}
	if err := restoreSolverSeat(tx, pulje, billettholderID); err != nil {
		return err
	}
	return tx.Commit()
}

// restoreSolverSeat saves the billettholder's emulated seat in pulje as a solver
// seat, if the pulje's distribution has been saved before and they have no
// other saved Player seat there.
func restoreSolverSeat(tx *sql.Tx, pulje models.Pulje, billettholderID int) error {
	var saved, ownSeats int
	if err := tx.QueryRow(
		`SELECT
			COALESCE(SUM(CASE WHEN source = ? THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN billettholder_id = ? THEN 1 ELSE 0 END), 0)
		 FROM relation_events_players WHERE pulje_id = ? AND role = ?`,
		SourceSolver, billettholderID, pulje, models.EventPlayerRolePlayer,
	).Scan(&saved, &ownSeats); err != nil {
		return fmt.Errorf("check saved seats in %s: %w", pulje, err)
	}
	if saved == 0 || ownSeats > 0 {
		return nil
	}
	em, err := emulateSeatings(tx)
	if err != nil {
		return fmt.Errorf("emulate after removing manual seat: %w", err)
	}
	seats, _, _ := puljeSeats(em, pulje)
	seat, ok := seats[billettholderID]
	if !ok || seat.IsGM {
		return nil
	}
	if _, err := tx.Exec(
		`INSERT INTO relation_events_players (event_id, pulje_id, billettholder_id, role, source) VALUES (?, ?, ?, ?, ?)`,
		seat.EventID, pulje, billettholderID, models.EventPlayerRolePlayer, SourceSolver,
	); err != nil {
		return fmt.Errorf("restore solver seat for billettholder %d in %s: %w", billettholderID, pulje, err)
	}
	return nil
}

type execer interface {
	Exec(query string, args ...any) (sql.Result, error)
}

func removeManualSeat(db execer, pulje models.Pulje, eventID string, billettholderID int) error {
	const query = `
		DELETE FROM relation_events_players
		WHERE event_id = ? AND pulje_id = ? AND billettholder_id = ?
		  AND source = ? AND role = ?
	`
	if _, err := db.Exec(query, eventID, string(pulje), billettholderID, SourceManual, models.EventPlayerRolePlayer); err != nil {
		return fmt.Errorf("remove manual seat (pulje=%s event=%s bh=%d): %w", pulje, eventID, billettholderID, err)
	}
	return nil
}
