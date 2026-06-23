package puljefordeling

import (
	"database/sql"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/Regncon/conorganizer/models"
)

// CommitPuljeAssignments solves the target pulje — seeded by the actual seats of
// earlier frozen puljer and pinning the pulje's existing manual placements — and
// persists the solver-chosen seats. In one transaction it removes any prior
// source='solver' Player rows for the pulje, inserts the fresh ones, and sets the
// pulje status to Locked. Manual and GM rows are left untouched. Idempotent.
func CommitPuljeAssignments(db *sql.DB, target models.Pulje, logger *slog.Logger) error {
	d, err := loadSeatingData(db)
	if err != nil {
		return fmt.Errorf("load seating data: %w", err)
	}
	idx := -1
	for i := range d.puljer {
		if d.puljer[i].ID == target {
			idx = i
			break
		}
	}
	if idx == -1 {
		return fmt.Errorf("pulje %s not found", target)
	}

	_, results := d.solveChronological(idx)
	res := results[idx]
	manual := d.manualFixed[string(target)]

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("begin commit tx for %s: %w", target, err)
	}
	defer tx.Rollback()

	if _, err := tx.Exec(
		`DELETE FROM relation_events_players WHERE pulje_id = ? AND role = ? AND source = 'solver'`,
		string(target), models.EventPlayerRolePlayer,
	); err != nil {
		return fmt.Errorf("clear solver seats for %s: %w", target, err)
	}

	const insert = `
		INSERT INTO relation_events_players (event_id, pulje_id, billettholder_id, role, source)
		VALUES (?, ?, ?, ?, 'solver')
		ON CONFLICT(billettholder_id, event_id, pulje_id) DO NOTHING
	`
	inserted := 0
	for evID, pids := range res.Assignments {
		for _, pid := range pids {
			if _, isManual := manual[pid]; isManual {
				continue // already persisted as a manual row
			}
			bh, convErr := strconv.Atoi(pid)
			if convErr != nil {
				continue
			}
			if _, err := tx.Exec(insert, evID, string(target), bh, models.EventPlayerRolePlayer); err != nil {
				return fmt.Errorf("insert solver seat (pulje=%s event=%s bh=%d): %w", target, evID, bh, err)
			}
			inserted++
		}
	}

	if _, err := tx.Exec(`UPDATE puljer SET status = ? WHERE id = ?`, string(models.PuljeStatusLocked), string(target)); err != nil {
		return fmt.Errorf("lock pulje %s: %w", target, err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit seats for %s: %w", target, err)
	}
	if logger != nil {
		logger.Info("committed puljefordeling", "pulje_id", target, "solver_seats", inserted)
	}
	return nil
}

// RevertPuljeAssignments removes the solver-written seats for a pulje (leaving manual
// and GM rows intact) and reopens it, in one transaction.
func RevertPuljeAssignments(db *sql.DB, target models.Pulje) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("begin revert tx for %s: %w", target, err)
	}
	defer tx.Rollback()

	if _, err := tx.Exec(
		`DELETE FROM relation_events_players WHERE pulje_id = ? AND role = ? AND source = 'solver'`,
		string(target), models.EventPlayerRolePlayer,
	); err != nil {
		return fmt.Errorf("revert solver seats for %s: %w", target, err)
	}
	if _, err := tx.Exec(`UPDATE puljer SET status = ? WHERE id = ?`, string(models.PuljeStatusOpen), string(target)); err != nil {
		return fmt.Errorf("reopen pulje %s: %w", target, err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit revert for %s: %w", target, err)
	}
	return nil
}
