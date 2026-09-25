package puljefordeling

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"sort"

	"github.com/Regncon/conorganizer/models"
)

// Lagringsendringer describes how the emulated distribution for a pulje differs
// from what is saved, so "Lagre fordeling" can say what saving would change.
type Lagringsendringer struct {
	PuljeNavn     string
	ForsteLagring bool         // the fordeling has never been saved for this pulje
	Endringer     []Konsekvens // Fra = saved seat, Til = seat a save would store; sorted by name
	Bekreftelse   string       // confirms exactly these changes; empty when there is nothing to save
}

// HarEndringer reports whether saving would change anything.
func (l Lagringsendringer) HarEndringer() bool {
	return len(l.Endringer) > 0
}

// Lagringsstatus compares the pulje's saved Player seats with the current
// emulated distribution.
func Lagringsstatus(db *sql.DB, pulje models.Pulje) (Lagringsendringer, error) {
	em, err := emulateSeatings(db)
	if err != nil {
		return Lagringsendringer{}, fmt.Errorf("emuler fordeling for lagringsstatus: %w", err)
	}
	return LagringsstatusFor(db, em, pulje)
}

// LagringsstatusFor is Lagringsstatus for an emulation the caller already has.
func LagringsstatusFor(db emulationQuerier, em Emulation, pulje models.Pulje) (Lagringsendringer, error) {
	lagret, forsteLagring, err := hentLagredePlasser(db, pulje)
	if err != nil {
		return Lagringsendringer{}, err
	}
	out := Lagringsendringer{ForsteLagring: forsteLagring}

	foreslatt := map[int]Plass{}
	navn := map[int]string{}
	for id, d := range lagret {
		navn[id] = d.Name
	}
	for _, p := range em.Puljer {
		if p.PuljeID != pulje {
			continue
		}
		out.PuljeNavn = p.Name
		for _, ev := range p.Events {
			for _, pl := range ev.AssignedPlayers {
				if _, seen := foreslatt[pl.BillettholderID]; seen {
					continue
				}
				foreslatt[pl.BillettholderID] = Plass{EventID: ev.EventID, EventTitle: ev.Title, Level: pl.Level}
				navn[pl.BillettholderID] = pl.Name
			}
		}
	}

	for id := range navn {
		if lagret[id].Til == foreslatt[id] {
			continue
		}
		out.Endringer = append(out.Endringer, Konsekvens{BillettholderID: id, Name: navn[id], Fra: lagret[id].Til, Til: foreslatt[id]})
	}
	sort.Slice(out.Endringer, func(i, j int) bool {
		if out.Endringer[i].Name != out.Endringer[j].Name {
			return out.Endringer[i].Name < out.Endringer[j].Name
		}
		return out.Endringer[i].BillettholderID < out.Endringer[j].BillettholderID
	})
	if out.HarEndringer() {
		out.Bekreftelse = lagringsbekreftelse(pulje, out.Endringer)
	}
	return out, nil
}

// hentLagredePlasser returns each billettholder's saved Player seat in the
// pulje (in Til), and whether the fordeling has never been saved there.
func hentLagredePlasser(db emulationQuerier, pulje models.Pulje) (map[int]Konsekvens, bool, error) {
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
		return nil, false, fmt.Errorf("hent lagrede plasser i %s: %w", pulje, err)
	}
	defer rows.Close()

	lagret := map[int]Konsekvens{}
	forsteLagring := true
	for rows.Next() {
		var id int
		var plass Plass
		var name, source string
		if err := rows.Scan(&id, &plass.EventID, &plass.EventTitle, &plass.Level, &name, &source); err != nil {
			return nil, false, fmt.Errorf("les lagret plass i %s: %w", pulje, err)
		}
		if source == SourceSolver {
			forsteLagring = false
		}
		if _, seen := lagret[id]; !seen {
			lagret[id] = Konsekvens{BillettholderID: id, Name: name, Til: plass}
		}
	}
	return lagret, forsteLagring, rows.Err()
}

func lagringsbekreftelse(pulje models.Pulje, endringer []Konsekvens) string {
	h := sha256.New()
	skrivHashfelt(h, "lagring-v1", string(pulje))
	for _, e := range endringer {
		skrivHashfelt(h, fmt.Sprint(e.BillettholderID), e.Fra.EventID, e.Til.EventID)
	}
	return hex.EncodeToString(h.Sum(nil))
}

// LagreBekreftetFordeling saves the pulje's distribution only if bekreftelse
// still matches the changes the admin was shown. Otherwise nothing is saved and
// the current changes are returned so they can be confirmed again.
func LagreBekreftetFordeling(db *sql.DB, pulje models.Pulje, bekreftelse string) (*Lagringsendringer, error) {
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
	status, err := LagringsstatusFor(tx, em, pulje)
	if err != nil {
		return nil, err
	}
	if !status.HarEndringer() || bekreftelse != status.Bekreftelse {
		return &status, nil
	}
	if err := lagreFordeling(tx, em, pulje); err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit distribution for %s: %w", pulje, err)
	}
	return nil, nil
}
