package puljefordeling

import (
	"database/sql"
	"fmt"
	"sort"
	"strings"

	"github.com/Regncon/conorganizer/models"
)

type Fordelingstildeling struct {
	EventID    string
	EventTitle string
	Role       models.EventPlayerRole
}

type Fordelingsvarsel struct {
	BillettholderID      int
	BillettholderNavn    string
	Tildelinger          []Fordelingstildeling
	AntallSpillerplasser int
	AntallGMOppdrag      int
}

type fordelingstildelingKey struct {
	eventID string
	role    models.EventPlayerRole
}

type fordelingsvarselgrunnlag struct {
	billettholderID   int
	billettholderNavn string
	tildelinger       map[fordelingstildelingKey]Fordelingstildeling
}

// FinnFordelingsvarsler finds overlapping assignments in one displayed pulje.
// Player assignments come from the current preview, supplemented by persisted
// manual Player rows. Every persisted GM row is included. Persisted solver
// Player rows are deliberately ignored because they may describe an older
// saved distribution than the preview being shown.
func FinnFordelingsvarsler(db *sql.DB, pulje EmulatedPulje) ([]Fordelingsvarsel, error) {
	grunnlag := make(map[int]*fordelingsvarselgrunnlag)
	leggTil := func(billettholderID int, navn string, tildeling Fordelingstildeling) {
		if billettholderID <= 0 {
			return
		}
		person := grunnlag[billettholderID]
		if person == nil {
			person = &fordelingsvarselgrunnlag{
				billettholderID: billettholderID,
				tildelinger:     make(map[fordelingstildelingKey]Fordelingstildeling),
			}
			grunnlag[billettholderID] = person
		}
		if strings.TrimSpace(navn) != "" {
			person.billettholderNavn = strings.TrimSpace(navn)
		}
		person.tildelinger[fordelingstildelingKey{eventID: tildeling.EventID, role: tildeling.Role}] = tildeling
	}

	for _, event := range pulje.Events {
		for _, player := range event.AssignedPlayers {
			leggTil(player.BillettholderID, player.Name, Fordelingstildeling{
				EventID:    event.EventID,
				EventTitle: event.Title,
				Role:       models.EventPlayerRolePlayer,
			})
		}
	}

	const query = `
		SELECT
			ep.billettholder_id,
			b.first_name,
			b.last_name,
			ep.event_id,
			e.title,
			ep.role
		FROM relation_events_players ep
		JOIN billettholdere b ON b.id = ep.billettholder_id
		JOIN events e ON e.id = ep.event_id
		WHERE ep.pulje_id = ?
			AND (
				ep.role = ?
				OR (ep.role = ? AND ep.source = ?)
			)
	`
	rows, err := db.Query(
		query,
		pulje.PuljeID,
		models.EventPlayerRoleGM,
		models.EventPlayerRolePlayer,
		SourceManual,
	)
	if err != nil {
		return nil, fmt.Errorf("hent lagrede tildelinger for varsler i %s: %w", pulje.PuljeID, err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			billettholderID int
			firstName       string
			lastName        string
			tildeling       Fordelingstildeling
		)
		if err := rows.Scan(
			&billettholderID,
			&firstName,
			&lastName,
			&tildeling.EventID,
			&tildeling.EventTitle,
			&tildeling.Role,
		); err != nil {
			return nil, fmt.Errorf("les lagret tildeling for varsel i %s: %w", pulje.PuljeID, err)
		}
		leggTil(billettholderID, firstName+" "+lastName, tildeling)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("gå gjennom lagrede tildelinger for varsler i %s: %w", pulje.PuljeID, err)
	}

	varsler := make([]Fordelingsvarsel, 0, len(grunnlag))
	for _, person := range grunnlag {
		varsel := Fordelingsvarsel{
			BillettholderID:   person.billettholderID,
			BillettholderNavn: person.billettholderNavn,
			Tildelinger:       make([]Fordelingstildeling, 0, len(person.tildelinger)),
		}
		for _, tildeling := range person.tildelinger {
			varsel.Tildelinger = append(varsel.Tildelinger, tildeling)
			switch tildeling.Role {
			case models.EventPlayerRolePlayer:
				varsel.AntallSpillerplasser++
			case models.EventPlayerRoleGM:
				varsel.AntallGMOppdrag++
			}
		}
		harFleireSpillerplasser := varsel.AntallSpillerplasser > 1
		harFleireGMOppdrag := varsel.AntallGMOppdrag > 1
		harBeggeRoller := varsel.AntallSpillerplasser > 0 && varsel.AntallGMOppdrag > 0
		if !harFleireSpillerplasser && !harFleireGMOppdrag && !harBeggeRoller {
			continue
		}
		sort.Slice(varsel.Tildelinger, func(i, j int) bool {
			left := varsel.Tildelinger[i]
			right := varsel.Tildelinger[j]
			if left.EventTitle != right.EventTitle {
				return left.EventTitle < right.EventTitle
			}
			if left.EventID != right.EventID {
				return left.EventID < right.EventID
			}
			return left.Role < right.Role
		})
		varsler = append(varsler, varsel)
	}

	sort.Slice(varsler, func(i, j int) bool {
		if varsler[i].BillettholderNavn != varsler[j].BillettholderNavn {
			return varsler[i].BillettholderNavn < varsler[j].BillettholderNavn
		}
		return varsler[i].BillettholderID < varsler[j].BillettholderID
	})
	return varsler, nil
}
