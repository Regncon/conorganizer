package puljefordeling

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"

	"github.com/Regncon/conorganizer/models"
)

// SaveStatus describes how the emulated distribution for a pulje differs from
// what is saved, so "Lagre fordeling" can say what saving would change.
type SaveStatus struct {
	PuljeName    string
	FirstSave    bool          // the distribution has never been saved for this pulje
	Changes      []Consequence // From = saved seat, To = seat a save would store; sorted by name
	Confirmation string        // confirms exactly these changes; empty when there is nothing to save
}

// HasChanges reports whether saving would change anything.
func (s SaveStatus) HasChanges() bool {
	return len(s.Changes) > 0
}

// LoadSaveStatus compares the pulje's saved Player seats with the current
// emulated distribution.
func LoadSaveStatus(db *sql.DB, pulje models.Pulje) (SaveStatus, error) {
	em, err := emulateSeatings(db)
	if err != nil {
		return SaveStatus{}, fmt.Errorf("emulate for save status: %w", err)
	}
	return SaveStatusFor(db, em, pulje)
}

// SaveStatusFor is LoadSaveStatus for an emulation the caller already has.
func SaveStatusFor(db emulationQuerier, em Emulation, pulje models.Pulje) (SaveStatus, error) {
	saved, firstSave, err := loadSavedSeats(db, pulje)
	if err != nil {
		return SaveStatus{}, err
	}
	out := SaveStatus{FirstSave: firstSave}

	proposed := map[int]Seat{}
	names := map[int]string{}
	for id, seat := range saved {
		names[id] = seat.Name
	}
	for _, p := range em.Puljer {
		if p.PuljeID != pulje {
			continue
		}
		out.PuljeName = p.Name
		for _, ev := range p.Events {
			for _, pl := range ev.AssignedPlayers {
				if _, seen := proposed[pl.BillettholderID]; seen {
					continue
				}
				proposed[pl.BillettholderID] = Seat{EventID: ev.EventID, EventTitle: ev.Title, Level: pl.Level}
				names[pl.BillettholderID] = pl.Name
			}
		}
	}

	for id := range names {
		if saved[id].To == proposed[id] {
			continue
		}
		out.Changes = append(out.Changes, Consequence{BillettholderID: id, Name: names[id], From: saved[id].To, To: proposed[id]})
	}
	sortConsequences(out.Changes)
	if out.HasChanges() {
		out.Confirmation = saveConfirmationToken(pulje, out.Changes)
	}
	return out, nil
}

// loadSavedSeats returns each billettholder's saved Player seat in the pulje
// (in To), and whether the distribution has never been saved there.
func loadSavedSeats(db emulationQuerier, pulje models.Pulje) (map[int]Consequence, bool, error) {
	rows, err := db.Query(`
		SELECT rep.billettholder_id, rep.event_id, e.title, COALESCE(i.interest_level, ''),
		       TRIM(b.first_name || ' ' || b.last_name), rep.source
		FROM relation_events_players rep
		JOIN events e ON e.id = rep.event_id
		JOIN billettholdere b ON b.id = rep.billettholder_id
		LEFT JOIN interests i ON i.billettholder_id = rep.billettholder_id
			AND i.event_id = rep.event_id AND i.pulje_id = rep.pulje_id
		WHERE rep.pulje_id = ? AND rep.role = ?
		ORDER BY rep.billettholder_id, rep.event_id
	`, pulje, models.EventPlayerRolePlayer)
	if err != nil {
		return nil, false, fmt.Errorf("load saved seats in %s: %w", pulje, err)
	}
	defer rows.Close()

	saved := map[int]Consequence{}
	firstSave := true
	for rows.Next() {
		var id int
		var seat Seat
		var name, source string
		if err := rows.Scan(&id, &seat.EventID, &seat.EventTitle, &seat.Level, &name, &source); err != nil {
			return nil, false, fmt.Errorf("read saved seat in %s: %w", pulje, err)
		}
		if source == SourceSolver {
			firstSave = false
		}
		if _, seen := saved[id]; !seen {
			saved[id] = Consequence{BillettholderID: id, Name: name, To: seat}
		}
	}
	return saved, firstSave, rows.Err()
}

func saveConfirmationToken(pulje models.Pulje, changes []Consequence) string {
	h := sha256.New()
	skrivHashfelt(h, "save-v1", string(pulje))
	for _, change := range changes {
		skrivHashfelt(h, fmt.Sprint(change.BillettholderID), change.From.EventID, change.To.EventID)
	}
	return hex.EncodeToString(h.Sum(nil))
}

// SaveConfirmedDistribution saves the pulje's distribution only if confirmation
// still matches the changes the admin was shown. Otherwise nothing is saved and
// the current changes are returned so they can be confirmed again.
func SaveConfirmedDistribution(db *sql.DB, pulje models.Pulje, confirmation string) (*SaveStatus, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, fmt.Errorf("begin commit tx: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	var puljeStatus models.PuljeStatus
	if err := tx.QueryRow(`SELECT status FROM puljer WHERE id = ?`, pulje).Scan(&puljeStatus); err != nil {
		return nil, fmt.Errorf("read pulje status before commit: %w", err)
	}
	if puljeStatus == models.PuljeStatusCompleted {
		return nil, ErrPuljeCompleted
	}
	em, err := emulateSeatings(tx)
	if err != nil {
		return nil, fmt.Errorf("emulate before commit: %w", err)
	}
	status, err := SaveStatusFor(tx, em, pulje)
	if err != nil {
		return nil, err
	}
	if !status.HasChanges() || confirmation != status.Confirmation {
		return &status, nil
	}
	if err := saveDistribution(tx, em, pulje); err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit distribution for %s: %w", pulje, err)
	}
	return nil, nil
}
