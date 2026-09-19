package puljefordeling

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"sort"
	"strings"

	"github.com/Regncon/conorganizer/models"
)

var ErrUgyldigTildeling = errors.New("ugyldig tildeling")

type Tildeling struct {
	EventID    string
	EventTitle string
	Role       models.EventPlayerRole
	Source     string
	InsertedAt string
}

type Tildelingsvalg struct {
	PuljeID         models.Pulje
	EventID         string
	BillettholderID int
	Role            models.EventPlayerRole
	FraLeggTil      bool
	FraEventID      string
	FraManuellPlass bool
	Bekreftelse     string
	AlderBekreftet  bool
	Forstevalg      bool
}

type Tildelingsvarsel struct {
	BillettholderID   int
	BillettholderNavn string
	PuljeID           models.Pulje
	PuljeNavn         string
	EventID           string
	EventTitle        string
	Role              models.EventPlayerRole
	Tildelinger       []Tildeling
	Aldersvarsel      string
	Kapasitetsvarsel  string
	Bekreftelse       string
	Handling          string
	KanBekrefte       bool
}

func TildelBillettholder(db *sql.DB, valg Tildelingsvalg) (*Tildelingsvarsel, error) {
	if err := validerTildelingsvalg(valg); err != nil {
		return nil, err
	}

	tx, err := db.Begin()
	if err != nil {
		return nil, fmt.Errorf("start tildeling: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	grunnlag, err := hentTildelingsgrunnlag(tx, valg)
	if err != nil {
		return nil, err
	}
	if !valg.FraLeggTil {
		fraEventID, err := finnFlyttetSpillerplass(valg.FraEventID, valg.FraManuellPlass, grunnlag.tildelinger)
		if err != nil {
			return nil, err
		}
		valg.FraEventID = fraEventID
		for _, tildeling := range grunnlag.tildelinger {
			if tildeling.Role == models.EventPlayerRolePlayer && tildeling.Source == SourceManual && tildeling.EventID == valg.EventID {
				return nil, nil
			}
		}
	}
	if varsel := lagTildelingsvarsel(valg, grunnlag); varsel != nil {
		if !varsel.KanBekrefte {
			return varsel, nil
		}
		if valg.Bekreftelse != "" {
			if valg.Bekreftelse != varsel.Bekreftelse {
				return varsel, nil
			}
		} else if tildelingerKreverBekreftelse(valg, grunnlag.tildelinger) || varsel.Kapasitetsvarsel != "" || (varsel.Aldersvarsel != "" && !valg.AlderBekreftet) {
			return varsel, nil
		}
	}
	if err := lagreTildeling(tx, valg); err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("fullfør tildeling: %w", err)
	}
	return nil, nil
}

func finnFlyttetSpillerplass(fraEventID string, fraManuellPlass bool, tildelinger []Tildeling) (string, error) {
	var manuellePlasser []string
	for _, tildeling := range tildelinger {
		if tildeling.Role == models.EventPlayerRolePlayer && tildeling.Source == SourceManual {
			manuellePlasser = append(manuellePlasser, tildeling.EventID)
		}
	}
	// En automatisk forhåndsvisning trenger ikke å ha en lagret spillerplass.
	if len(manuellePlasser) == 0 && !fraManuellPlass {
		return fraEventID, nil
	}
	// Eldre klienter sendte ikke kilden. Én manuell plass er entydig.
	if fraEventID == "" && len(manuellePlasser) == 1 {
		return manuellePlasser[0], nil
	}
	for _, eventID := range manuellePlasser {
		if eventID == fraEventID {
			return fraEventID, nil
		}
	}
	return "", fmt.Errorf("%w: velg spillerplassen som skal flyttes", ErrUgyldigTildeling)
}

func GetTildelinger(db *sql.DB, pulje models.Pulje, billettholderID int) ([]Tildeling, error) {
	if _, ok := models.ParsePulje(string(pulje)); !ok || billettholderID <= 0 {
		return nil, ErrUgyldigTildeling
	}
	var exists int
	if err := db.QueryRow(`SELECT 1 FROM puljer WHERE id = ?`, pulje).Scan(&exists); err != nil {
		return nil, fmt.Errorf("hent pulje %s: %w", pulje, err)
	}
	if err := db.QueryRow(`SELECT 1 FROM billettholdere WHERE id = ?`, billettholderID).Scan(&exists); err != nil {
		return nil, fmt.Errorf("hent billettholder %d: %w", billettholderID, err)
	}
	return hentTildelinger(db, pulje, billettholderID)
}

func FjernTildeling(db *sql.DB, pulje models.Pulje, eventID string, billettholderID int, role models.EventPlayerRole) error {
	valg := Tildelingsvalg{
		PuljeID:         pulje,
		EventID:         eventID,
		BillettholderID: billettholderID,
		Role:            role,
		FraLeggTil:      true,
	}
	if err := validerTildelingsvalg(valg); err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("start fjerning av tildeling: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	if _, err := hentTildelingsgrunnlag(tx, valg); err != nil {
		return err
	}
	result, err := tx.Exec(
		`DELETE FROM relation_events_players WHERE event_id = ? AND pulje_id = ? AND billettholder_id = ? AND role = ?`,
		eventID, pulje, billettholderID, role,
	)
	if err != nil {
		return fmt.Errorf("fjern %s-tildeling for billettholder %d frå %s i %s: %w", role, billettholderID, eventID, pulje, err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("les fjerna %s-tildeling: %w", role, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("finn %s-tildeling for billettholder %d på %s i %s: %w", role, billettholderID, eventID, pulje, sql.ErrNoRows)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("fullfør fjerning av tildeling: %w", err)
	}
	return nil
}

type queryer interface {
	Query(string, ...any) (*sql.Rows, error)
}

func hentTildelinger(db queryer, pulje models.Pulje, billettholderID int) ([]Tildeling, error) {
	const query = `
		SELECT ep.event_id, e.title, ep.role, ep.source, COALESCE(ep.inserted_at, '')
		FROM relation_events_players ep
		JOIN events e ON e.id = ep.event_id
		WHERE ep.pulje_id = ? AND ep.billettholder_id = ?
		ORDER BY ep.event_id, ep.role, ep.source, ep.inserted_at
	`
	rows, err := db.Query(query, pulje, billettholderID)
	if err != nil {
		return nil, fmt.Errorf("hent tildelinger for billettholder %d i %s: %w", billettholderID, pulje, err)
	}
	defer rows.Close()

	var tildelinger []Tildeling
	for rows.Next() {
		var tildeling Tildeling
		if err := rows.Scan(&tildeling.EventID, &tildeling.EventTitle, &tildeling.Role, &tildeling.Source, &tildeling.InsertedAt); err != nil {
			return nil, fmt.Errorf("les tildeling for billettholder %d i %s: %w", billettholderID, pulje, err)
		}
		tildelinger = append(tildelinger, tildeling)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("gå gjennom tildelinger for billettholder %d i %s: %w", billettholderID, pulje, err)
	}
	return tildelinger, nil
}

type tildelingsgrunnlag struct {
	puljeNavn         string
	puljeStatus       models.PuljeStatus
	eventTitle        string
	ageGroup          models.AgeGroup
	billettholderNavn string
	isOver18          bool
	tildelinger       []Tildeling
	kapasiteter       []tildelingskapasitet
}

type tildelingskapasitet struct {
	ageGroup               models.AgeGroup
	eventID                string
	eventTitle             string
	maksSpillere           int
	manuelleSpillerplasser int
}

func hentTildelingsgrunnlag(tx *sql.Tx, valg Tildelingsvalg) (tildelingsgrunnlag, error) {
	var grunnlag tildelingsgrunnlag
	if err := tx.QueryRow(`SELECT name, status FROM puljer WHERE id = ?`, valg.PuljeID).Scan(&grunnlag.puljeNavn, &grunnlag.puljeStatus); err != nil {
		return grunnlag, fmt.Errorf("hent pulje %s: %w", valg.PuljeID, err)
	}
	if grunnlag.puljeStatus == models.PuljeStatusCompleted {
		return grunnlag, ErrPuljeCompleted
	}
	var firstName, lastName string
	if err := tx.QueryRow(
		`SELECT first_name, last_name, is_over_18 FROM billettholdere WHERE id = ?`,
		valg.BillettholderID,
	).Scan(&firstName, &lastName, &grunnlag.isOver18); err != nil {
		return grunnlag, fmt.Errorf("hent billettholder %d: %w", valg.BillettholderID, err)
	}
	grunnlag.billettholderNavn = strings.TrimSpace(firstName + " " + lastName)
	if err := tx.QueryRow(
		`SELECT title, age_group FROM events WHERE id = ?`,
		valg.EventID,
	).Scan(&grunnlag.eventTitle, &grunnlag.ageGroup); err != nil {
		return grunnlag, fmt.Errorf("hent arrangement %s: %w", valg.EventID, err)
	}
	var finnes bool
	if err := tx.QueryRow(
		`SELECT EXISTS(SELECT 1 FROM relation_event_puljer WHERE event_id = ? AND pulje_id = ? AND is_in_pulje = 1)`,
		valg.EventID, valg.PuljeID,
	).Scan(&finnes); err != nil {
		return grunnlag, fmt.Errorf("kontroller arrangement %s i pulje %s: %w", valg.EventID, valg.PuljeID, err)
	}
	if !finnes {
		return grunnlag, fmt.Errorf("%w: arrangementet er ikke med i puljen", ErrUgyldigTildeling)
	}
	tildelinger, err := hentTildelinger(tx, valg.PuljeID, valg.BillettholderID)
	if err != nil {
		return grunnlag, err
	}
	grunnlag.tildelinger = tildelinger
	grunnlag.kapasiteter, err = hentTildelingskapasiteter(tx, valg, tildelinger)
	return grunnlag, err
}

func hentTildelingskapasiteter(tx *sql.Tx, valg Tildelingsvalg, tildelinger []Tildeling) ([]tildelingskapasitet, error) {
	var eventIDs []string
	if valg.Role == models.EventPlayerRolePlayer {
		eventIDs = append(eventIDs, valg.EventID)
	}
	if valg.FraLeggTil {
		for _, tildeling := range tildelinger {
			if tildeling.Role != models.EventPlayerRolePlayer || tildeling.Source != SourceSolver {
				continue
			}
			if valg.Role == models.EventPlayerRolePlayer && tildeling.EventID == valg.EventID {
				continue
			}
			eventIDs = append(eventIDs, tildeling.EventID)
		}
	}
	sort.Strings(eventIDs)
	var kapasiteter []tildelingskapasitet
	for _, eventID := range eventIDs {
		kapasitet := tildelingskapasitet{eventID: eventID}
		if err := tx.QueryRow(
			`SELECT e.title, e.age_group, e.max_players,
			 (SELECT COUNT(*) + 1 FROM relation_events_players
			  WHERE event_id = e.id AND pulje_id = ? AND role = ? AND source = ? AND billettholder_id != ?)
			 FROM events e WHERE e.id = ?`,
			valg.PuljeID, models.EventPlayerRolePlayer, SourceManual, valg.BillettholderID, eventID,
		).Scan(&kapasitet.eventTitle, &kapasitet.ageGroup, &kapasitet.maksSpillere, &kapasitet.manuelleSpillerplasser); err != nil {
			return nil, fmt.Errorf("hent kapasitet på %s i %s: %w", eventID, valg.PuljeID, err)
		}
		kapasiteter = append(kapasiteter, kapasitet)
	}
	return kapasiteter, nil
}

func validerTildelingsvalg(valg Tildelingsvalg) error {
	if _, ok := models.ParsePulje(string(valg.PuljeID)); !ok {
		return fmt.Errorf("%w: ukjent pulje", ErrUgyldigTildeling)
	}
	if strings.TrimSpace(valg.EventID) == "" || valg.BillettholderID <= 0 {
		return fmt.Errorf("%w: arrangement og billettholder må være valgt", ErrUgyldigTildeling)
	}
	if valg.Role != models.EventPlayerRolePlayer && valg.Role != models.EventPlayerRoleGM {
		return fmt.Errorf("%w: ukjent rolle", ErrUgyldigTildeling)
	}
	if !valg.FraLeggTil && valg.Role != models.EventPlayerRolePlayer {
		return fmt.Errorf("%w: dra-og-slipp kan berre flytte ein spelar", ErrUgyldigTildeling)
	}
	return nil
}

func lagTildelingsvarsel(valg Tildelingsvalg, grunnlag tildelingsgrunnlag) *Tildelingsvarsel {
	bekreftelse := tildelingsbekreftelse(valg, grunnlag)
	harGM := false
	for _, tildeling := range grunnlag.tildelinger {
		if tildeling.Role == models.EventPlayerRoleGM {
			harGM = true
		}
	}
	var aldersvarsler []string
	if !grunnlag.isOver18 && grunnlag.ageGroup == models.AgeGroupAdultsOnly {
		aldersvarsler = append(aldersvarsler, fmt.Sprintf("%s er under 18 år, og «%s» er 18+.", grunnlag.billettholderNavn, grunnlag.eventTitle))
	}
	var kapasitetsvarsler []string
	for _, kapasitet := range grunnlag.kapasiteter {
		if !grunnlag.isOver18 && kapasitet.eventID != valg.EventID && kapasitet.ageGroup == models.AgeGroupAdultsOnly {
			aldersvarsler = append(aldersvarsler, fmt.Sprintf("%s er under 18 år, og «%s» er 18+.", grunnlag.billettholderNavn, kapasitet.eventTitle))
		}
		if kapasitet.manuelleSpillerplasser > kapasitet.maksSpillere {
			kapasitetsvarsler = append(kapasitetsvarsler, fmt.Sprintf("«%s» får %d manuelt festede spillerplasser, men har kapasitet til %d.", kapasitet.eventTitle, kapasitet.manuelleSpillerplasser, kapasitet.maksSpillere))
		}
	}
	aldersvarsel := strings.Join(aldersvarsler, " ")
	kapasitetsvarsel := strings.Join(kapasitetsvarsler, " ")
	bekreftelseKreves := tildelingerKreverBekreftelse(valg, grunnlag.tildelinger)
	if !bekreftelseKreves && aldersvarsel == "" && kapasitetsvarsel == "" && (valg.Bekreftelse == "" || valg.Bekreftelse == bekreftelse) {
		return nil
	}

	varsel := &Tildelingsvarsel{
		BillettholderID:   valg.BillettholderID,
		BillettholderNavn: grunnlag.billettholderNavn,
		PuljeID:           valg.PuljeID,
		PuljeNavn:         grunnlag.puljeNavn,
		EventID:           valg.EventID,
		EventTitle:        grunnlag.eventTitle,
		Role:              valg.Role,
		Tildelinger:       grunnlag.tildelinger,
		Aldersvarsel:      aldersvarsel,
		Kapasitetsvarsel:  kapasitetsvarsel,
		Handling:          tildelingshandling(valg, grunnlag),
		KanBekrefte:       true,
	}
	if !valg.FraLeggTil && harGM {
		varsel.KanBekrefte = false
		varsel.Handling = "Bruk Legg til for å plassere ein GM som spelar."
	}
	varsel.Bekreftelse = bekreftelse
	return varsel
}

func tildelingerKreverBekreftelse(valg Tildelingsvalg, tildelinger []Tildeling) bool {
	if !valg.FraLeggTil {
		for _, tildeling := range tildelinger {
			if tildeling.Role == models.EventPlayerRoleGM {
				return true
			}
		}
		return false
	}
	for _, tildeling := range tildelinger {
		if tildeling.EventID == valg.EventID && tildeling.Role == valg.Role {
			continue
		}
		return true
	}
	return false
}

func tildelingshandling(valg Tildelingsvalg, grunnlag tildelingsgrunnlag) string {
	if valg.Role == models.EventPlayerRoleGM {
		return fmt.Sprintf("Legg til som GM på «%s»", grunnlag.eventTitle)
	}
	if valg.FraLeggTil {
		return fmt.Sprintf("Legg til som spelar på «%s»", grunnlag.eventTitle)
	}
	for _, tildeling := range grunnlag.tildelinger {
		if tildeling.Role == models.EventPlayerRolePlayer && tildeling.EventID == valg.FraEventID && tildeling.EventID != valg.EventID {
			return fmt.Sprintf("Flytt spillerplassen fra «%s» til «%s»", tildeling.EventTitle, grunnlag.eventTitle)
		}
	}
	return fmt.Sprintf("Legg til som spelar på «%s»", grunnlag.eventTitle)
}

func tildelingsbekreftelse(valg Tildelingsvalg, grunnlag tildelingsgrunnlag) string {
	h := sha256.New()
	skrivHashfelt(h, "tildeling-v3")
	skrivHashfelt(h, string(valg.PuljeID), valg.EventID, fmt.Sprint(valg.BillettholderID), string(valg.Role))
	skrivHashfelt(h, fmt.Sprint(valg.FraLeggTil), valg.FraEventID, fmt.Sprint(valg.FraManuellPlass), fmt.Sprint(valg.Forstevalg))
	skrivHashfelt(h, grunnlag.puljeNavn, string(grunnlag.puljeStatus), grunnlag.eventTitle, string(grunnlag.ageGroup))
	skrivHashfelt(h, grunnlag.billettholderNavn, fmt.Sprint(grunnlag.isOver18))
	for _, kapasitet := range grunnlag.kapasiteter {
		skrivHashfelt(h, kapasitet.eventID, kapasitet.eventTitle, string(kapasitet.ageGroup), fmt.Sprint(kapasitet.maksSpillere), fmt.Sprint(kapasitet.manuelleSpillerplasser))
	}
	for _, tildeling := range grunnlag.tildelinger {
		skrivHashfelt(h, tildeling.EventID, tildeling.EventTitle, string(tildeling.Role), tildeling.Source, tildeling.InsertedAt)
	}
	return hex.EncodeToString(h.Sum(nil))
}

func skrivHashfelt(h hash.Hash, felter ...string) {
	for _, felt := range felter {
		_, _ = fmt.Fprintf(h, "%d:%s|", len(felt), felt)
	}
}

func lagreTildeling(tx *sql.Tx, valg Tildelingsvalg) error {
	if valg.FraLeggTil {
		if _, err := tx.Exec(
			`UPDATE relation_events_players SET source = ?
			 WHERE pulje_id = ? AND billettholder_id = ? AND role = ? AND source = ?`,
			SourceManual, valg.PuljeID, valg.BillettholderID, models.EventPlayerRolePlayer, SourceSolver,
		); err != nil {
			return fmt.Errorf("fest eksisterende spillerplasser i %s: %w", valg.PuljeID, err)
		}
	}
	if valg.Role == models.EventPlayerRolePlayer && !valg.FraLeggTil {
		if _, err := tx.Exec(
			`DELETE FROM relation_events_players WHERE pulje_id = ? AND billettholder_id = ? AND role = ?
			 AND (event_id = ? OR source = ? OR ? = '')`,
			valg.PuljeID, valg.BillettholderID, models.EventPlayerRolePlayer, valg.FraEventID, SourceSolver, valg.FraEventID,
		); err != nil {
			return fmt.Errorf("fjern tidligere Player-tildeling i %s: %w", valg.PuljeID, err)
		}
	}

	result, err := tx.Exec(
		`UPDATE relation_events_players
		 SET source = ?, inserted_at = `+models.DBDateTimeNowSQL+`
		 WHERE event_id = ? AND pulje_id = ? AND billettholder_id = ? AND role = ?`,
		SourceManual, valg.EventID, valg.PuljeID, valg.BillettholderID, valg.Role,
	)
	if err != nil {
		return fmt.Errorf("oppdater tildeling: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("les oppdatert tildeling: %w", err)
	}
	if rowsAffected == 0 {
		if _, err := tx.Exec(
			`INSERT INTO relation_events_players (event_id, pulje_id, billettholder_id, role, source) VALUES (?, ?, ?, ?, ?)`,
			valg.EventID, valg.PuljeID, valg.BillettholderID, valg.Role, SourceManual,
		); err != nil {
			return fmt.Errorf("lagre tildeling: %w", err)
		}
	}

	if valg.Forstevalg {
		const query = `
			INSERT INTO interests (billettholder_id, event_id, pulje_id, interest_level)
			VALUES (?, ?, ?, ?)
			ON CONFLICT(billettholder_id, event_id, pulje_id) DO UPDATE SET
				interest_level = excluded.interest_level,
				updated_at = ` + models.DBDateTimeNowSQL
		if _, err := tx.Exec(query, valg.BillettholderID, valg.EventID, valg.PuljeID, models.InterestLevelHigh); err != nil {
			return fmt.Errorf("lagre førstevalg: %w", err)
		}
	}
	return nil
}
