package puljefordeling

import (
	"database/sql"
	"fmt"
	"sort"

	"github.com/Regncon/conorganizer/models"
)

// Plass is where a participant sits in one pulje, as shown in a consequence.
type Plass struct {
	EventID    string // empty when the participant has no seat
	EventTitle string
	Level      models.InterestLevel // their interest in a player seat
	IsGM       bool
	Forstevalg bool // the seat gives them their førstevalg in this pulje
}

// Konsekvens is another participant whose seat changes along with a manual change.
type Konsekvens struct {
	BillettholderID int
	Name            string
	Fra, Til        Plass
}

// Konsekvenser lists how a manual change moves other participants in the same
// pulje. Later puljer are left out on purpose: they are likely to change again
// before they are published.
type Konsekvenser struct {
	PuljeNavn string
	Endringer []Konsekvens // sorted by name
}

// Fjerningsvarsel describes a manual seat or GM removal and its consequences,
// so an admin can confirm it.
type Fjerningsvarsel struct {
	BillettholderNavn string
	EventTitle        string
	Role              models.EventPlayerRole
	Konsekvenser      Konsekvenser
}

// ForhandsvisFjerning shows what removing a manual Player seat or a GM would
// change in the pulje, without removing anything.
func ForhandsvisFjerning(db *sql.DB, pulje models.Pulje, eventID string, billettholderID int, role models.EventPlayerRole) (Fjerningsvarsel, error) {
	valg := Tildelingsvalg{PuljeID: pulje, EventID: eventID, BillettholderID: billettholderID, Role: role, FraLeggTil: true}
	if err := validerTildelingsvalg(valg); err != nil {
		return Fjerningsvarsel{}, err
	}
	tx, err := db.Begin()
	if err != nil {
		return Fjerningsvarsel{}, fmt.Errorf("start forhåndsvisning av fjerning: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	grunnlag, err := hentTildelingsgrunnlag(tx, valg)
	if err != nil {
		return Fjerningsvarsel{}, err
	}
	konsekvenser, err := forhandsvisKonsekvenser(tx, pulje, billettholderID, func(tx *sql.Tx) error {
		if role == models.EventPlayerRoleGM {
			return slettTildeling(tx, pulje, eventID, billettholderID, role)
		}
		return fjernManuellPlass(tx, pulje, eventID, billettholderID)
	})
	if err != nil {
		return Fjerningsvarsel{}, err
	}
	return Fjerningsvarsel{
		BillettholderNavn: grunnlag.billettholderNavn,
		EventTitle:        grunnlag.eventTitle,
		Role:              role,
		Konsekvenser:      konsekvenser,
	}, nil
}

// forhandsvisKonsekvenser applies endring inside tx and compares the pulje's
// emulated seating before and after. The caller must roll tx back.
func forhandsvisKonsekvenser(tx *sql.Tx, pulje models.Pulje, billettholderID int, endring func(*sql.Tx) error) (Konsekvenser, error) {
	before, err := emulateSeatings(tx)
	if err != nil {
		return Konsekvenser{}, fmt.Errorf("emuler fordeling før endring: %w", err)
	}
	if err := endring(tx); err != nil {
		return Konsekvenser{}, err
	}
	after, err := emulateSeatings(tx)
	if err != nil {
		return Konsekvenser{}, fmt.Errorf("emuler fordeling etter endring: %w", err)
	}
	return sammenlignPlasser(before, after, pulje, billettholderID), nil
}

// sammenlignPlasser lists everyone except the changed billettholder whose seat
// in pulje differs between the two emulations.
func sammenlignPlasser(before, after Emulation, pulje models.Pulje, billettholderID int) Konsekvenser {
	fra, navn, puljeNavn := puljePlasser(before, pulje)
	til, navnEtter, _ := puljePlasser(after, pulje)
	for id, n := range navnEtter {
		navn[id] = n
	}

	out := Konsekvenser{PuljeNavn: puljeNavn}
	for id := range navn {
		if id == billettholderID || fra[id] == til[id] {
			continue
		}
		out.Endringer = append(out.Endringer, Konsekvens{BillettholderID: id, Name: navn[id], Fra: fra[id], Til: til[id]})
	}
	sort.Slice(out.Endringer, func(i, j int) bool {
		if out.Endringer[i].Name != out.Endringer[j].Name {
			return out.Endringer[i].Name < out.Endringer[j].Name
		}
		return out.Endringer[i].BillettholderID < out.Endringer[j].BillettholderID
	})
	return out
}

// puljePlasser returns each seated participant's seat and name in pulje.
func puljePlasser(em Emulation, pulje models.Pulje) (map[int]Plass, map[int]string, string) {
	plasser := map[int]Plass{}
	navn := map[int]string{}
	for _, p := range em.Puljer {
		if p.PuljeID != pulje {
			continue
		}
		forstevalg := make(map[int]bool, len(p.FikkForstevalg))
		for _, d := range p.FikkForstevalg {
			forstevalg[d.BillettholderID] = true
		}
		for _, ev := range p.Events {
			for _, pl := range ev.AssignedPlayers {
				plasser[pl.BillettholderID] = Plass{EventID: ev.EventID, EventTitle: ev.Title, Level: pl.Level, Forstevalg: forstevalg[pl.BillettholderID]}
				navn[pl.BillettholderID] = pl.Name
			}
			for _, gm := range ev.AssignedGMs {
				plasser[gm.BillettholderID] = Plass{EventID: ev.EventID, EventTitle: ev.Title, IsGM: true}
				navn[gm.BillettholderID] = gm.Name
			}
		}
		return plasser, navn, p.Name
	}
	return plasser, navn, ""
}
