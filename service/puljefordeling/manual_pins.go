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
// pure emulation based on their real interests. If a seat already exists it is
// reclaimed as a manual player seat.
func AddManualSeat(db *sql.DB, pulje models.Pulje, eventID string, billettholderID int) error {
	const query = `
		INSERT INTO relation_events_players (event_id, pulje_id, billettholder_id, role, source)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT(billettholder_id, event_id, pulje_id) DO UPDATE SET
			role = EXCLUDED.role,
			source = EXCLUDED.source
	`
	if _, err := db.Exec(query, eventID, string(pulje), billettholderID, models.EventPlayerRolePlayer, SourceManual); err != nil {
		return fmt.Errorf("add manual seat (pulje=%s event=%s bh=%d): %w", pulje, eventID, billettholderID, err)
	}
	return nil
}

// RemoveManualSeat deletes an admin-pinned player seat (source='manual',
// role='Player') for the given pulje/event/participant. It only removes manual
// player pins — solver seats and GM rows are left untouched. Removing the pin
// does not touch the player's interest, so a later emulation may still seat them
// in the same event by simulation (now as a non-manual placement).
func RemoveManualSeat(db *sql.DB, pulje models.Pulje, eventID string, billettholderID int) error {
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
