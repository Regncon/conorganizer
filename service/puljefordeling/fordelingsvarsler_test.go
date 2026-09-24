package puljefordeling

import (
	"database/sql"
	"fmt"
	"reflect"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/bdd"
)

func TestFinnKapasitetsvarsler_ManuellePlasserOgForhandsvisningTellesEnGang(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{Given: "Fem manuelle spillere er tildelt et arrangement med fire plasser, og fire vises også i forhåndsvisningen.", When: "Kapasitetsvarsler beregnes.", Then: "Ett varsel teller hver spiller én gang og viser fem av fire plasser."})
	// Given
	expectedPlayers := []string{"Spiller 1 Nordmann", "Spiller 2 Nordmann", "Spiller 3 Nordmann", "Spiller 4 Nordmann", "Spiller 5 Nordmann"}
	db := capacityWarningFixture(t)
	for id := 1; id <= 5; id++ {
		seedWarningAssignment(t, db, id, "ev-a", models.PuljeFredagKveld, models.EventPlayerRolePlayer, SourceManual)
	}
	pulje := capacityWarningPreview(1, 2, 3, 4)

	// When
	varsler, err := FinnKapasitetsvarsler(db, pulje)

	// Then
	if err != nil {
		t.Fatalf("finn kapasitetsvarsler: %v", err)
	}
	if len(varsler) != 1 {
		t.Fatalf("antall kapasitetsvarsler = %d, vil ha 1", len(varsler))
	}
	varsel := varsler[0]
	if varsel.EventID != "ev-a" || varsel.Kapasitet != 4 || varsel.AntallSpillerplasser != 5 {
		t.Errorf("kapasitetsvarsel = %+v, vil ha ev-a med 5 av 4 plasser", varsel)
	}
	if !reflect.DeepEqual(varsel.Spillere, expectedPlayers) {
		t.Errorf("spillere = %v, vil ha %v", varsel.Spillere, expectedPlayers)
	}
}

func TestFinnKapasitetsvarsler_FulltArrangementUtenUtdaterteSolverplasserVarslesIkke(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{Given: "Fire viste manuelle spillere fyller fire plasser, og en gammel solverplass og et GM-oppdrag er lagret på samme arrangement.", When: "Kapasitetsvarsler beregnes.", Then: "Arrangementet regnes som fullt, uten å varsle om overkapasitet."})
	// Given
	expectedWarnings := 0
	db := capacityWarningFixture(t)
	for id := 1; id <= 4; id++ {
		seedWarningAssignment(t, db, id, "ev-a", models.PuljeFredagKveld, models.EventPlayerRolePlayer, SourceManual)
	}
	seedWarningAssignment(t, db, 5, "ev-a", models.PuljeFredagKveld, models.EventPlayerRolePlayer, SourceSolver)
	seedWarningAssignment(t, db, 6, "ev-a", models.PuljeFredagKveld, models.EventPlayerRoleGM, SourceManual)
	pulje := capacityWarningPreview(1, 2, 3, 4)

	// When
	varsler, err := FinnKapasitetsvarsler(db, pulje)

	// Then
	if err != nil {
		t.Fatalf("finn kapasitetsvarsler: %v", err)
	}
	if len(varsler) != expectedWarnings {
		t.Errorf("antall kapasitetsvarsler = %d, vil ha %d: %+v", len(varsler), expectedWarnings, varsler)
	}
}

func TestFinnKapasitetsvarsler_ManuellePlasserUtenforForhandsvisningVarsles(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{Given: "Et arrangement utenfor den synlige listen har fem manuelle spillere og fire plasser.", When: "Kapasitetsvarsler beregnes.", Then: "De lagrede manuelle plassene gir ett varsel med arrangementets kapasitet."})
	// Given
	expectedEventID := "ev-b"
	db := capacityWarningFixture(t)
	for id := 1; id <= 5; id++ {
		seedWarningAssignment(t, db, id, "ev-b", models.PuljeFredagKveld, models.EventPlayerRolePlayer, SourceManual)
	}
	pulje := capacityWarningPreview()

	// When
	varsler, err := FinnKapasitetsvarsler(db, pulje)

	// Then
	if err != nil {
		t.Fatalf("finn kapasitetsvarsler: %v", err)
	}
	if len(varsler) != 1 || varsler[0].EventID != expectedEventID || varsler[0].AntallSpillerplasser != 5 || varsler[0].Kapasitet != 4 {
		t.Errorf("kapasitetsvarsler = %+v, vil ha ev-b med 5 av 4 plasser", varsler)
	}
}

func TestFinnKapasitetsvarsler_TildelingerIAndrePuljerTellesIkke(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{Given: "Et arrangement har fire spillere i én pulje og en manuell spillerplass i en annen pulje.", When: "Kapasitetsvarsler beregnes for den første puljen.", Then: "Plassen i den andre puljen gir ikke kapasitetsvarsel."})
	// Given
	expectedWarnings := 0
	db := capacityWarningFixture(t)
	seedPulje(t, db, models.PuljeLordagKveld, "Lørdag kveld", "2026-09-16T18:00:00Z")
	seedWarningAssignment(t, db, 5, "ev-a", models.PuljeLordagKveld, models.EventPlayerRolePlayer, SourceManual)
	pulje := capacityWarningPreview(1, 2, 3, 4)

	// When
	varsler, err := FinnKapasitetsvarsler(db, pulje)

	// Then
	if err != nil {
		t.Fatalf("finn kapasitetsvarsler: %v", err)
	}
	if len(varsler) != expectedWarnings {
		t.Errorf("antall kapasitetsvarsler = %d, vil ha %d: %+v", len(varsler), expectedWarnings, varsler)
	}
}

func capacityWarningFixture(t *testing.T) *sql.DB {
	t.Helper()
	db := testutil.CreateTestDB(t, "kapasitetsvarsler")
	seedWarningBase(t, db, models.PuljeFredagKveld)
	for id := 1; id <= 6; id++ {
		seedParticipant(t, db, id, fmt.Sprintf("Spiller %d", id), "Nordmann")
	}
	return db
}

func capacityWarningPreview(ids ...int) EmulatedPulje {
	event := EmulatedEvent{EventID: "ev-a", Title: "Alpha", Capacity: 4}
	for _, id := range ids {
		event.AssignedPlayers = append(event.AssignedPlayers, AssignedPlayer{BillettholderID: id, Name: fmt.Sprintf("Spiller %d Nordmann", id)})
	}
	return EmulatedPulje{PuljeID: models.PuljeFredagKveld, Events: []EmulatedEvent{event}}
}

func TestFinnFordelingsvarsler_FlereGMOppdragVarslesOgSorteres(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt to billettholdere med flere GM-oppdrag i samme pulje.",
		When:  "Når fordelingsvarsler beregnes.",
		Then:  "Så vises begge i navneorden med alle GM-oppdrag sortert etter arrangement.",
	})

	// Given
	expectedNames := []string{"Ada Aasen", "Zoe Zulu"}
	expectedEvents := []string{"Alpha", "Bravo"}
	db := testutil.CreateTestDB(t, "fordelingsvarsler_flere_gm")
	seedWarningBase(t, db, models.PuljeFredagKveld)
	seedParticipant(t, db, 1, "Zoe", "Zulu")
	seedParticipant(t, db, 2, "Ada", "Aasen")
	for _, billettholderID := range []int{1, 2} {
		seedWarningAssignment(t, db, billettholderID, "ev-b", models.PuljeFredagKveld, models.EventPlayerRoleGM, SourceManual)
		seedWarningAssignment(t, db, billettholderID, "ev-a", models.PuljeFredagKveld, models.EventPlayerRoleGM, SourceManual)
	}
	pulje := EmulatedPulje{PuljeID: models.PuljeFredagKveld}

	// When
	varsler, err := FinnFordelingsvarsler(db, pulje)

	// Then
	if err != nil {
		t.Fatalf("finn fordelingsvarsler: %v", err)
	}
	if len(varsler) != len(expectedNames) {
		t.Fatalf("antall varsler = %d, vil ha %d", len(varsler), len(expectedNames))
	}
	for i, expectedName := range expectedNames {
		if varsler[i].BillettholderNavn != expectedName {
			t.Fatalf("varsel %d gjelder %q, vil ha %q", i, varsler[i].BillettholderNavn, expectedName)
		}
		if varsler[i].AntallGMOppdrag != 2 || varsler[i].AntallSpillerplasser != 0 {
			t.Fatalf("antall roller for %s = GM %d, spiller %d; vil ha GM 2, spiller 0", expectedName, varsler[i].AntallGMOppdrag, varsler[i].AntallSpillerplasser)
		}
		for assignmentIndex, expectedEvent := range expectedEvents {
			if varsler[i].Tildelinger[assignmentIndex].EventTitle != expectedEvent {
				t.Fatalf("tildeling %d for %s = %q, vil ha %q", assignmentIndex, expectedName, varsler[i].Tildelinger[assignmentIndex].EventTitle, expectedEvent)
			}
		}
	}
}

func TestFinnFordelingsvarsler_FlereManuelleSpillerplasserVarsles(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt eldre data med flere manuelle spillerplasser for samme billettholder og pulje.",
		When:  "Når fordelingsvarsler beregnes.",
		Then:  "Så vises én advarsel med begge spillerplassene.",
	})

	// Given
	const expectedPlayerAssignments = 2
	db := testutil.CreateTestDB(t, "fordelingsvarsler_flere_spillerplasser")
	seedWarningBase(t, db, models.PuljeFredagKveld)
	seedParticipant(t, db, 1, "Kari", "Nordmann")
	seedWarningAssignment(t, db, 1, "ev-a", models.PuljeFredagKveld, models.EventPlayerRolePlayer, SourceManual)
	seedWarningAssignment(t, db, 1, "ev-b", models.PuljeFredagKveld, models.EventPlayerRolePlayer, SourceManual)
	pulje := EmulatedPulje{PuljeID: models.PuljeFredagKveld}

	// When
	varsler, err := FinnFordelingsvarsler(db, pulje)

	// Then
	if err != nil {
		t.Fatalf("finn fordelingsvarsler: %v", err)
	}
	varsel := requireSingleWarning(t, varsler)
	if varsel.AntallSpillerplasser != expectedPlayerAssignments || varsel.AntallGMOppdrag != 0 {
		t.Fatalf("antall roller = spiller %d, GM %d; vil ha spiller %d, GM 0", varsel.AntallSpillerplasser, varsel.AntallGMOppdrag, expectedPlayerAssignments)
	}
}

func TestFinnFordelingsvarsler_GMOgForhandsvistSpillerPaSammeArrangementVarsles(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en lagret GM og en forhåndsvist spillerplass på samme arrangement.",
		When:  "Når fordelingsvarsler beregnes.",
		Then:  "Så vises begge rollene som to selvstendige tildelinger.",
	})

	// Given
	expectedRoles := []models.EventPlayerRole{models.EventPlayerRoleGM, models.EventPlayerRolePlayer}
	db := testutil.CreateTestDB(t, "fordelingsvarsler_samme_arrangement")
	seedWarningBase(t, db, models.PuljeFredagKveld)
	seedParticipant(t, db, 1, "Kari", "Nordmann")
	seedWarningAssignment(t, db, 1, "ev-a", models.PuljeFredagKveld, models.EventPlayerRoleGM, SourceManual)
	pulje := warningPreview(models.PuljeFredagKveld, "ev-a", "Alpha", 1, "Kari Nordmann")

	// When
	varsler, err := FinnFordelingsvarsler(db, pulje)

	// Then
	if err != nil {
		t.Fatalf("finn fordelingsvarsler: %v", err)
	}
	varsel := requireSingleWarning(t, varsler)
	if varsel.AntallSpillerplasser != 1 || varsel.AntallGMOppdrag != 1 {
		t.Fatalf("antall roller = spiller %d, GM %d; vil ha 1 av hver", varsel.AntallSpillerplasser, varsel.AntallGMOppdrag)
	}
	if len(varsel.Tildelinger) != len(expectedRoles) {
		t.Fatalf("antall tildelinger = %d, vil ha %d", len(varsel.Tildelinger), len(expectedRoles))
	}
	for i, expectedRole := range expectedRoles {
		if varsel.Tildelinger[i].Role != expectedRole {
			t.Fatalf("rolle %d = %q, vil ha %q", i, varsel.Tildelinger[i].Role, expectedRole)
		}
	}
}

func TestFinnFordelingsvarsler_GMOgManuellSpillerPaUlikeArrangementVarslesUtenDuplikat(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en GM og en manuell spillerplass som også finnes i forhåndsvisningen.",
		When:  "Når fordelingsvarsler beregnes.",
		Then:  "Så vises de to arrangementene én gang hver.",
	})

	// Given
	const expectedAssignments = 2
	db := testutil.CreateTestDB(t, "fordelingsvarsler_ulike_arrangement")
	seedWarningBase(t, db, models.PuljeFredagKveld)
	seedParticipant(t, db, 1, "Kari", "Nordmann")
	seedWarningAssignment(t, db, 1, "ev-a", models.PuljeFredagKveld, models.EventPlayerRoleGM, SourceManual)
	seedWarningAssignment(t, db, 1, "ev-b", models.PuljeFredagKveld, models.EventPlayerRolePlayer, SourceManual)
	pulje := warningPreview(models.PuljeFredagKveld, "ev-b", "Bravo", 1, "Kari Nordmann")

	// When
	varsler, err := FinnFordelingsvarsler(db, pulje)

	// Then
	if err != nil {
		t.Fatalf("finn fordelingsvarsler: %v", err)
	}
	varsel := requireSingleWarning(t, varsler)
	if len(varsel.Tildelinger) != expectedAssignments {
		t.Fatalf("antall tildelinger = %d, vil ha %d: %+v", len(varsel.Tildelinger), expectedAssignments, varsel.Tildelinger)
	}
}

func TestFinnFordelingsvarsler_TildelingerIAndrePuljerVarslesIkke(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en forhåndsvist spiller i én pulje og et GM-oppdrag i en annen pulje.",
		When:  "Når fordelingsvarsler beregnes for spillerens pulje.",
		Then:  "Så regnes ikke GM-oppdraget som en overlapp.",
	})

	// Given
	expectedWarnings := 0
	db := testutil.CreateTestDB(t, "fordelingsvarsler_annen_pulje")
	seedWarningBase(t, db, models.PuljeFredagKveld)
	seedPulje(t, db, models.PuljeLordagKveld, "Lørdag kveld", "2026-09-16T18:00:00Z")
	seedEvent(t, db, "ev-d", "Delta", 4, models.PuljeLordagKveld)
	seedParticipant(t, db, 1, "Kari", "Nordmann")
	seedWarningAssignment(t, db, 1, "ev-d", models.PuljeLordagKveld, models.EventPlayerRoleGM, SourceManual)
	pulje := warningPreview(models.PuljeFredagKveld, "ev-a", "Alpha", 1, "Kari Nordmann")

	// When
	varsler, err := FinnFordelingsvarsler(db, pulje)

	// Then
	if err != nil {
		t.Fatalf("finn fordelingsvarsler: %v", err)
	}
	if len(varsler) != expectedWarnings {
		t.Fatalf("antall varsler = %d, vil ha %d: %+v", len(varsler), expectedWarnings, varsler)
	}
}

func TestFinnFordelingsvarsler_UtelaterLagretSolverplassSomIkkeVises(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en utdatert lagret solverplass, en annen forhåndsvist spillerplass og et GM-oppdrag.",
		When:  "Når fordelingsvarsler beregnes.",
		Then:  "Så inneholder varselet forhåndsvisningen og GM-oppdraget, men ikke den utdaterte solverplassen.",
	})

	// Given
	expectedEventIDs := []string{"ev-b", "ev-c"}
	db := testutil.CreateTestDB(t, "fordelingsvarsler_utdatert_solverplass")
	seedWarningBase(t, db, models.PuljeFredagKveld)
	seedParticipant(t, db, 1, "Kari", "Nordmann")
	seedWarningAssignment(t, db, 1, "ev-a", models.PuljeFredagKveld, models.EventPlayerRolePlayer, SourceSolver)
	seedWarningAssignment(t, db, 1, "ev-c", models.PuljeFredagKveld, models.EventPlayerRoleGM, SourceManual)
	pulje := warningPreview(models.PuljeFredagKveld, "ev-b", "Bravo", 1, "Kari Nordmann")

	// When
	varsler, err := FinnFordelingsvarsler(db, pulje)

	// Then
	if err != nil {
		t.Fatalf("finn fordelingsvarsler: %v", err)
	}
	varsel := requireSingleWarning(t, varsler)
	if len(varsel.Tildelinger) != len(expectedEventIDs) {
		t.Fatalf("antall tildelinger = %d, vil ha %d: %+v", len(varsel.Tildelinger), len(expectedEventIDs), varsel.Tildelinger)
	}
	for i, expectedEventID := range expectedEventIDs {
		if varsel.Tildelinger[i].EventID != expectedEventID {
			t.Fatalf("tildeling %d gjelder %q, vil ha %q", i, varsel.Tildelinger[i].EventID, expectedEventID)
		}
	}
}

func seedWarningBase(t *testing.T, db *sql.DB, pulje models.Pulje) {
	t.Helper()
	seedPulje(t, db, pulje, "Fredag kveld", "2026-09-15T18:00:00Z")
	seedEvent(t, db, "ev-a", "Alpha", 4, pulje)
	seedEvent(t, db, "ev-b", "Bravo", 4, pulje)
	seedEvent(t, db, "ev-c", "Charlie", 4, pulje)
}

func seedWarningAssignment(t *testing.T, db *sql.DB, billettholderID int, eventID string, pulje models.Pulje, role models.EventPlayerRole, source string) {
	t.Helper()
	if _, err := db.Exec(`
		INSERT INTO relation_events_players (billettholder_id, event_id, pulje_id, role, source)
		VALUES (?, ?, ?, ?, ?)`, billettholderID, eventID, pulje, role, source); err != nil {
		t.Fatalf("seed %s assignment for billettholder %d: %v", role, billettholderID, err)
	}
}

func warningPreview(pulje models.Pulje, eventID, eventTitle string, billettholderID int, name string) EmulatedPulje {
	return EmulatedPulje{
		PuljeID: pulje,
		Events: []EmulatedEvent{{
			EventID: eventID,
			Title:   eventTitle,
			AssignedPlayers: []AssignedPlayer{{
				BillettholderID: billettholderID,
				Name:            name,
			}},
		}},
	}
}

func requireSingleWarning(t *testing.T, varsler []Fordelingsvarsel) Fordelingsvarsel {
	t.Helper()
	if len(varsler) != 1 {
		t.Fatalf("antall varsler = %d, vil ha 1: %+v", len(varsler), varsler)
	}
	return varsler[0]
}
