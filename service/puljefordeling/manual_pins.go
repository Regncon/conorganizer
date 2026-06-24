package puljefordeling

import (
	"database/sql"
	"fmt"

	"github.com/Regncon/conorganizer/models"
)

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
