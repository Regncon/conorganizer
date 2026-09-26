package puljefordeling

import (
	"database/sql"
	"fmt"
	"sort"

	"github.com/Regncon/conorganizer/models"
)

// Seat is where a participant sits in one pulje, as shown in a consequence.
type Seat struct {
	EventID    string // empty when the participant has no seat
	EventTitle string
	Level      models.InterestLevel // their interest in a player seat
	IsGM       bool
	Forstevalg bool // the seat gives them their førstevalg in this pulje
}

// Consequence is a participant whose seat changes along with a manual change.
type Consequence struct {
	BillettholderID int
	Name            string
	From, To        Seat
}

// Consequences lists how a manual change moves participants in the same pulje.
// Later puljer are left out on purpose: they are likely to change again before
// they are published.
type Consequences struct {
	PuljeName string
	Own       Consequence   // the changed billettholder's own seat before and after
	Changes   []Consequence // everyone else whose seat changes, sorted by name
}

// RemovalPreview describes a manual seat or GM removal and its consequences,
// so an admin can confirm it.
type RemovalPreview struct {
	BillettholderName string
	EventTitle        string
	Role              models.EventPlayerRole
	Consequences      Consequences
}

// PreviewRemoval shows what removing a manual Player seat or a GM would change
// in the pulje, without removing anything.
func PreviewRemoval(db *sql.DB, pulje models.Pulje, eventID string, billettholderID int, role models.EventPlayerRole) (RemovalPreview, error) {
	valg := Tildelingsvalg{PuljeID: pulje, EventID: eventID, BillettholderID: billettholderID, Role: role, FraLeggTil: true}
	if err := validerTildelingsvalg(valg); err != nil {
		return RemovalPreview{}, err
	}
	tx, err := db.Begin()
	if err != nil {
		return RemovalPreview{}, fmt.Errorf("begin removal preview: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	grunnlag, err := hentTildelingsgrunnlag(tx, valg)
	if err != nil {
		return RemovalPreview{}, err
	}
	consequences, err := previewConsequences(tx, pulje, billettholderID, func(tx *sql.Tx) error {
		if role == models.EventPlayerRoleGM {
			return deleteAssignment(tx, pulje, eventID, billettholderID, role)
		}
		return removeManualSeat(tx, pulje, eventID, billettholderID)
	})
	if err != nil {
		return RemovalPreview{}, err
	}
	return RemovalPreview{
		BillettholderName: grunnlag.billettholderNavn,
		EventTitle:        grunnlag.eventTitle,
		Role:              role,
		Consequences:      consequences,
	}, nil
}

// previewConsequences applies change inside tx and compares the pulje's
// emulated seating before and after. The caller must roll tx back.
func previewConsequences(tx *sql.Tx, pulje models.Pulje, billettholderID int, change func(*sql.Tx) error) (Consequences, error) {
	before, err := emulateSeatings(tx)
	if err != nil {
		return Consequences{}, fmt.Errorf("emulate before change: %w", err)
	}
	if err := change(tx); err != nil {
		return Consequences{}, err
	}
	after, err := emulateSeatings(tx)
	if err != nil {
		return Consequences{}, fmt.Errorf("emulate after change: %w", err)
	}
	return compareSeats(before, after, pulje, billettholderID), nil
}

// compareSeats lists the changed billettholder's own seat before and after, and
// everyone else whose seat in pulje differs between the two emulations.
func compareSeats(before, after Emulation, pulje models.Pulje, billettholderID int) Consequences {
	from, names, puljeName := puljeSeats(before, pulje)
	to, namesAfter, _ := puljeSeats(after, pulje)
	for id, name := range namesAfter {
		names[id] = name
	}

	out := Consequences{
		PuljeName: puljeName,
		Own:       Consequence{BillettholderID: billettholderID, Name: names[billettholderID], From: from[billettholderID], To: to[billettholderID]},
	}
	for id := range names {
		if id == billettholderID || from[id] == to[id] {
			continue
		}
		out.Changes = append(out.Changes, Consequence{BillettholderID: id, Name: names[id], From: from[id], To: to[id]})
	}
	sortConsequences(out.Changes)
	return out
}

func sortConsequences(consequences []Consequence) {
	sort.Slice(consequences, func(i, j int) bool {
		if consequences[i].Name != consequences[j].Name {
			return consequences[i].Name < consequences[j].Name
		}
		return consequences[i].BillettholderID < consequences[j].BillettholderID
	})
}

// puljeSeats returns each seated participant's seat and name in pulje, and the
// pulje's name.
func puljeSeats(em Emulation, pulje models.Pulje) (map[int]Seat, map[int]string, string) {
	seats := map[int]Seat{}
	names := map[int]string{}
	for _, p := range em.Puljer {
		if p.PuljeID != pulje {
			continue
		}
		gotForstevalg := make(map[int]bool, len(p.GotForstevalg))
		for _, participant := range p.GotForstevalg {
			gotForstevalg[participant.BillettholderID] = true
		}
		for _, ev := range p.Events {
			for _, pl := range ev.AssignedPlayers {
				seats[pl.BillettholderID] = Seat{EventID: ev.EventID, EventTitle: ev.Title, Level: pl.Level, Forstevalg: gotForstevalg[pl.BillettholderID]}
				names[pl.BillettholderID] = pl.Name
			}
			for _, gm := range ev.AssignedGMs {
				seats[gm.BillettholderID] = Seat{EventID: ev.EventID, EventTitle: ev.Title, IsGM: true}
				names[gm.BillettholderID] = gm.Name
			}
		}
		return seats, names, p.Name
	}
	return seats, names, ""
}
