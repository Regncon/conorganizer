package puljefordeling

import (
	"database/sql"
	"strings"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/bdd"
)

func TestTildelBillettholder_CapacityOverrideRequiresBoundConfirmation(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{Given: "four manual Player pins in a four-seat event", When: "an admin adds a fifth Player with only the age-confirmed flag", Then: "nothing changes until the capacity warning is explicitly confirmed"})

	// Given
	const expectedAssignments = 1
	db, _ := testutil.CreateTestDBAndLogger(t, "tildeling_capacity_confirmation")
	seedTildelingFixture(t, db, models.AgeGroupDefault, true)
	seedOtherPlayerSeats(t, db, 4, SourceManual)
	valg := Tildelingsvalg{PuljeID: models.PuljeFredagKveld, EventID: "evA", BillettholderID: 1, Role: models.EventPlayerRolePlayer, FraLeggTil: true, AlderBekreftet: true}

	// When
	varsel, err := TildelBillettholder(db, valg)

	// Then
	if err != nil || varsel == nil || !varsel.KanBekrefte || varsel.Bekreftelse == "" {
		t.Fatalf("expected confirmable capacity warning: warning=%+v err=%v", varsel, err)
	}
	if got := assignmentCountByRole(t, db, valg.PuljeID, valg.BillettholderID, valg.Role); got != 0 {
		t.Fatalf("unconfirmed capacity override wrote %d assignments", got)
	}
	valg.Bekreftelse = varsel.Bekreftelse
	varsel, err = TildelBillettholder(db, valg)
	if err != nil || varsel != nil {
		t.Fatalf("confirm capacity override: warning=%+v err=%v", varsel, err)
	}
	if got := assignmentCountByRole(t, db, valg.PuljeID, valg.BillettholderID, valg.Role); got != expectedAssignments {
		t.Fatalf("confirmed assignments: got %d, want %d", got, expectedAssignments)
	}
}

func TestTildelBillettholder_ChangedCapacityRequiresFreshConfirmation(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{Given: "a pending Player confirmation for a full event", When: "its capacity or manual occupancy changes", Then: "the old token writes nothing and returns a fresh warning"})

	for _, change := range []string{"capacity", "manual occupancy"} {
		t.Run(change, func(t *testing.T) {
			// Given
			const expectedNewAssignments = 0
			db, _ := testutil.CreateTestDBAndLogger(t, "tildeling_stale_capacity")
			seedTildelingFixture(t, db, models.AgeGroupDefault, true)
			seedOtherPlayerSeats(t, db, 4, SourceManual)
			seedEvent(t, db, "evB", "Bravo", 4, models.PuljeFredagKveld)
			seedGM(t, db, "evB", models.PuljeFredagKveld, 1)
			valg := Tildelingsvalg{PuljeID: models.PuljeFredagKveld, EventID: "evA", BillettholderID: 1, Role: models.EventPlayerRolePlayer, FraLeggTil: true}
			initial, err := TildelBillettholder(db, valg)
			if err != nil || initial == nil {
				t.Fatalf("initial warning: warning=%+v err=%v", initial, err)
			}
			valg.Bekreftelse = initial.Bekreftelse
			if change == "capacity" {
				testutil.MustExec(t, db, `UPDATE events SET max_players = 3 WHERE id = 'evA'`)
			} else {
				seedParticipant(t, db, 6, "Ola", "Nordmann")
				seedManualSeat(t, db, "evA", models.PuljeFredagKveld, 6)
			}

			// When
			fresh, err := TildelBillettholder(db, valg)

			// Then
			if err != nil || fresh == nil || fresh.Bekreftelse == initial.Bekreftelse {
				t.Fatalf("expected fresh capacity warning: warning=%+v err=%v", fresh, err)
			}
			if got := assignmentCountByRole(t, db, valg.PuljeID, valg.BillettholderID, valg.Role); got != expectedNewAssignments {
				t.Fatalf("stale confirmation wrote %d Player assignments", got)
			}
		})
	}
}

func TestTildelBillettholder_SolverSeatsDoNotRequireCapacityOverride(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{Given: "a four-seat event contains four replaceable solver seats", When: "an admin pins a new Player", Then: "the pin is saved without a manual capacity override"})

	// Given
	const expectedAssignments = 1
	db, _ := testutil.CreateTestDBAndLogger(t, "tildeling_solver_capacity")
	seedTildelingFixture(t, db, models.AgeGroupDefault, true)
	seedOtherPlayerSeats(t, db, 4, SourceSolver)
	valg := Tildelingsvalg{PuljeID: models.PuljeFredagKveld, EventID: "evA", BillettholderID: 1, Role: models.EventPlayerRolePlayer, FraLeggTil: true}

	// When
	varsel, err := TildelBillettholder(db, valg)

	// Then
	if err != nil || varsel != nil {
		t.Fatalf("replaceable solver seats triggered warning: warning=%+v err=%v", varsel, err)
	}
	if got := assignmentCountByRole(t, db, valg.PuljeID, valg.BillettholderID, valg.Role); got != expectedAssignments {
		t.Fatalf("manual assignments: got %d, want %d", got, expectedAssignments)
	}
}

func TestTildelBillettholder_ExistingPinDoesNotConsumeAnotherSeat(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{Given: "the Player already has one of four manual pins in a four-seat event", When: "the admin adds the same Player again", Then: "the same pin remains without a capacity warning"})

	// Given
	const expectedAssignments = 1
	db, _ := testutil.CreateTestDBAndLogger(t, "tildeling_existing_pin_capacity")
	seedTildelingFixture(t, db, models.AgeGroupDefault, true)
	seedOtherPlayerSeats(t, db, 3, SourceManual)
	seedManualSeat(t, db, "evA", models.PuljeFredagKveld, 1)
	valg := Tildelingsvalg{PuljeID: models.PuljeFredagKveld, EventID: "evA", BillettholderID: 1, Role: models.EventPlayerRolePlayer, FraLeggTil: true}

	// When
	varsel, err := TildelBillettholder(db, valg)

	// Then
	if err != nil || varsel != nil {
		t.Fatalf("existing pin triggered warning: warning=%+v err=%v", varsel, err)
	}
	if got := assignmentCountByRole(t, db, valg.PuljeID, valg.BillettholderID, valg.Role); got != expectedAssignments {
		t.Fatalf("assignments: got %d, want %d", got, expectedAssignments)
	}
}

func TestTildelBillettholder_PreservingSolverSeatWarnsAboutItsCapacity(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{Given: "an existing solver seat shares an event with four manual pins", When: "an admin requests a Player pin on another event", Then: "confirmation warns about preserving the fifth pin in the existing event"})

	// Given
	const expectedEventTitle = "Alpha"
	db, _ := testutil.CreateTestDBAndLogger(t, "tildeling_preserved_capacity")
	seedTildelingFixture(t, db, models.AgeGroupDefault, true)
	seedOtherPlayerSeats(t, db, 4, SourceManual)
	seedEvent(t, db, "evB", "Bravo", 4, models.PuljeFredagKveld)
	testutil.MustExec(t, db, `INSERT INTO relation_events_players(event_id,pulje_id,billettholder_id,role,source) VALUES ('evA',?,1,'Player','solver')`, models.PuljeFredagKveld)
	valg := Tildelingsvalg{PuljeID: models.PuljeFredagKveld, EventID: "evB", BillettholderID: 1, Role: models.EventPlayerRolePlayer, FraLeggTil: true}

	// When
	varsel, err := TildelBillettholder(db, valg)

	// Then
	if err != nil || varsel == nil || !strings.Contains(varsel.Kapasitetsvarsel, expectedEventTitle) {
		t.Fatalf("missing capacity warning for preserved seat: warning=%+v err=%v", varsel, err)
	}
	if got := testutil.QueryInt(t, db, `SELECT COUNT(*) FROM relation_events_players WHERE billettholder_id=1 AND source='manual'`); got != 0 {
		t.Fatalf("unconfirmed request pinned %d assignments", got)
	}
}

func TestTildelBillettholder_PreservedSeatCapacityChangeRequiresFreshConfirmation(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{Given: "confirmation will preserve a solver Player seat while adding a GM role", When: "another admin changes the existing Player event's capacity", Then: "the old confirmation cannot pin the seat or add the GM role"})

	// Given
	const expectedManualAssignments = 0
	db, _ := testutil.CreateTestDBAndLogger(t, "tildeling_preserved_stale_capacity")
	seedTildelingFixture(t, db, models.AgeGroupDefault, true)
	seedEvent(t, db, "evB", "Bravo", 4, models.PuljeFredagKveld)
	testutil.MustExec(t, db, `INSERT INTO relation_events_players(event_id,pulje_id,billettholder_id,role,source) VALUES ('evA',?,1,'Player','solver')`, models.PuljeFredagKveld)
	valg := Tildelingsvalg{PuljeID: models.PuljeFredagKveld, EventID: "evB", BillettholderID: 1, Role: models.EventPlayerRoleGM, FraLeggTil: true}
	initial, err := TildelBillettholder(db, valg)
	if err != nil || initial == nil {
		t.Fatalf("initial warning: warning=%+v err=%v", initial, err)
	}
	valg.Bekreftelse = initial.Bekreftelse
	testutil.MustExec(t, db, `UPDATE events SET max_players=0 WHERE id='evA'`)

	// When
	fresh, err := TildelBillettholder(db, valg)

	// Then
	if err != nil || fresh == nil || fresh.Bekreftelse == initial.Bekreftelse || fresh.Kapasitetsvarsel == "" {
		t.Fatalf("expected refreshed preserved-seat capacity warning: warning=%+v err=%v", fresh, err)
	}
	if got := testutil.QueryInt(t, db, `SELECT COUNT(*) FROM relation_events_players WHERE billettholder_id=1 AND source='manual'`); got != expectedManualAssignments {
		t.Fatalf("stale confirmation wrote %d manual assignments", got)
	}
}

func TestTildelBillettholder_PreservedSeatAgeRestrictionRequiresConfirmation(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{Given: "a minor has a saved solver seat that now requires age 18", When: "an admin confirms another assignment that would preserve that seat", Then: "the existing event's age restriction is shown and invalidates older confirmation"})

	// Given
	const expectedRestrictedEvent = "Alpha"
	db, _ := testutil.CreateTestDBAndLogger(t, "tildeling_preserved_age")
	seedTildelingFixture(t, db, models.AgeGroupDefault, false)
	seedEvent(t, db, "evB", "Bravo", 4, models.PuljeFredagKveld)
	testutil.MustExec(t, db, `INSERT INTO relation_events_players(event_id,pulje_id,billettholder_id,role,source) VALUES ('evA',?,1,'Player','solver')`, models.PuljeFredagKveld)
	valg := Tildelingsvalg{PuljeID: models.PuljeFredagKveld, EventID: "evB", BillettholderID: 1, Role: models.EventPlayerRoleGM, FraLeggTil: true}
	initial, err := TildelBillettholder(db, valg)
	if err != nil || initial == nil {
		t.Fatalf("initial warning: warning=%+v err=%v", initial, err)
	}
	valg.Bekreftelse = initial.Bekreftelse
	valg.AlderBekreftet = true
	testutil.MustExec(t, db, `UPDATE events SET age_group=? WHERE id='evA'`, models.AgeGroupAdultsOnly)

	// When
	fresh, err := TildelBillettholder(db, valg)

	// Then
	if err != nil || fresh == nil || !strings.Contains(fresh.Aldersvarsel, expectedRestrictedEvent) || fresh.Bekreftelse == initial.Bekreftelse {
		t.Fatalf("expected fresh age warning for preserved seat: warning=%+v err=%v", fresh, err)
	}
	if got := testutil.QueryInt(t, db, `SELECT COUNT(*) FROM relation_events_players WHERE billettholder_id=1 AND source='manual'`); got != 0 {
		t.Fatalf("stale age confirmation pinned %d assignments", got)
	}
	valg.Bekreftelse = fresh.Bekreftelse
	varsel, err := TildelBillettholder(db, valg)
	if err != nil || varsel != nil {
		t.Fatalf("confirm preserved age override: warning=%+v err=%v", varsel, err)
	}
	if got := testutil.QueryInt(t, db, `SELECT COUNT(*) FROM relation_events_players WHERE billettholder_id=1 AND source='manual'`); got != 2 {
		t.Fatalf("confirmed age override wrote %d assignments, want 2", got)
	}
}

func seedOtherPlayerSeats(t *testing.T, db *sql.DB, count int, source string) {
	t.Helper()
	for id := 2; id < count+2; id++ {
		seedParticipant(t, db, id, "Spiller", "Nordmann")
		testutil.MustExec(t, db, `INSERT INTO relation_events_players(event_id,pulje_id,billettholder_id,role,source) VALUES ('evA',?,?,'Player',?)`, models.PuljeFredagKveld, id, source)
	}
}
