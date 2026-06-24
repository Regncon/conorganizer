package puljefordeling

import (
	"database/sql"
	"fmt"

	"github.com/Regncon/conorganizer/models"
)

// CommitDistribution persists the current emulated distribution for a pulje to
// relation_events_players so it becomes the actual seating shown in other views
// (and, once the pulje is published, to participants).
//
// Solver-placed players are written as source='solver'; manually pinned players
// are already persisted as source='manual' and are left untouched. Any previous
// solver-committed seats for the pulje are cleared first, so re-committing always
// reflects the latest distribution. GM rows are not touched.
func CommitDistribution(db *sql.DB, pulje models.Pulje) error {
	em, err := EmulateSeatings(db)
	if err != nil {
		return fmt.Errorf("emulate before commit: %w", err)
	}

	var target EmulatedPulje
	found := false
	for _, p := range em.Puljer {
		if p.PuljeID == pulje {
			target = p
			found = true
			break
		}
	}
	if !found {
		return nil
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("begin commit tx: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	// Clear the previous solver-committed seats; manual pins and GM rows stay.
	if _, err := tx.Exec(
		`DELETE FROM relation_events_players WHERE pulje_id = ? AND source = ? AND role = ?`,
		string(pulje), SourceSolver, models.EventPlayerRolePlayer,
	); err != nil {
		return fmt.Errorf("clear solver seats for %s: %w", pulje, err)
	}

	const upsert = `
		INSERT INTO relation_events_players (event_id, pulje_id, billettholder_id, role, source)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT(billettholder_id, event_id, pulje_id) DO UPDATE SET
			role = EXCLUDED.role,
			source = EXCLUDED.source
	`
	for _, ev := range target.Events {
		for _, pl := range ev.AssignedPlayers {
			if pl.Manual {
				continue // already persisted as source='manual'
			}
			if _, err := tx.Exec(upsert, ev.EventID, string(pulje), pl.BillettholderID, models.EventPlayerRolePlayer, SourceSolver); err != nil {
				return fmt.Errorf("commit solver seat (event=%s bh=%d): %w", ev.EventID, pl.BillettholderID, err)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit distribution for %s: %w", pulje, err)
	}
	return nil
}
