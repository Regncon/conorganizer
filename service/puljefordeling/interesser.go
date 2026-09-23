package puljefordeling

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/Regncon/conorganizer/models"
)

// SetInteresse records an administrator's interest choice for one
// billettholder and one event occurrence. Locked puljer remain editable by
// administrators; a completed pulje is immutable for everyone.
func SetInteresse(db *sql.DB, puljeID models.Pulje, eventID string, billettholderID int, level models.InterestLevel) error {
	if _, ok := models.ParsePulje(string(puljeID)); !ok {
		return fmt.Errorf("%w: ukjent pulje", ErrUgyldigTildeling)
	}
	if strings.TrimSpace(eventID) == "" || billettholderID <= 0 || !level.Valid() {
		return fmt.Errorf("%w: arrangement, billettholder og interesse må være gyldige", ErrUgyldigTildeling)
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("start interesseendring: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	var status models.PuljeStatus
	if err := tx.QueryRow(`SELECT status FROM puljer WHERE id = ?`, puljeID).Scan(&status); err != nil {
		return fmt.Errorf("hent pulje %s: %w", puljeID, err)
	}
	if status == models.PuljeStatusCompleted {
		return ErrPuljeCompleted
	}

	var exists bool
	if err := tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM billettholdere WHERE id = ?)`, billettholderID).Scan(&exists); err != nil {
		return fmt.Errorf("kontroller billettholder %d: %w", billettholderID, err)
	}
	if !exists {
		return sql.ErrNoRows
	}
	if err := tx.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM relation_event_puljer
			WHERE event_id = ? AND pulje_id = ? AND is_in_pulje = 1
		)
	`, eventID, puljeID).Scan(&exists); err != nil {
		return fmt.Errorf("kontroller arrangement %s i pulje %s: %w", eventID, puljeID, err)
	}
	if !exists {
		return fmt.Errorf("%w: arrangementet er ikke med i puljen", ErrUgyldigTildeling)
	}

	if level == models.InterestLevelNone {
		if _, err := tx.Exec(`DELETE FROM interests WHERE billettholder_id = ? AND event_id = ? AND pulje_id = ?`, billettholderID, eventID, puljeID); err != nil {
			return fmt.Errorf("fjern interesse: %w", err)
		}
	} else if _, err := tx.Exec(`
		INSERT INTO interests (billettholder_id, event_id, pulje_id, interest_level)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(billettholder_id, event_id, pulje_id) DO UPDATE SET
			interest_level = excluded.interest_level,
			updated_at = `+models.DBDateTimeNowSQL, billettholderID, eventID, puljeID, level); err != nil {
		return fmt.Errorf("lagre interesse: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("fullfør interesseendring: %w", err)
	}
	return nil
}
