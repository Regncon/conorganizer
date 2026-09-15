package puljefordeling

import (
	"database/sql"
	"errors"
	"slices"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/bdd"
)

func TestTildelBillettholder_PlacesPlayerWithoutConflict(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{Given: "an eligible event and an unassigned billettholder", When: "an admin assigns the billettholder as Player", Then: "one manual Player pin is stored without changing interests"})

	// Given
	const expectedPlayerAssignments = 1
	const expectedInterests = 0
	db, _ := testutil.CreateTestDBAndLogger(t, "tildeling_player_without_conflict")
	seedTildelingFixture(t, db, models.AgeGroupDefault, true)
	valg := Tildelingsvalg{
		PuljeID:         models.PuljeFredagKveld,
		EventID:         "evA",
		BillettholderID: 1,
		Role:            models.EventPlayerRolePlayer,
		FraLeggTil:      true,
	}

	// When
	varsel, err := TildelBillettholder(db, valg)

	// Then
	if err != nil {
		t.Fatalf("assign Player: %v", err)
	}
	if varsel != nil {
		t.Fatalf("unexpected warning: %+v", varsel)
	}
	if got := assignmentCountByRole(t, db, valg.PuljeID, valg.BillettholderID, valg.Role); got != expectedPlayerAssignments {
		t.Fatalf("Player assignments: got %d, want %d", got, expectedPlayerAssignments)
	}
	var source string
	if err := db.QueryRow(`SELECT source FROM relation_events_players WHERE event_id = ? AND pulje_id = ? AND billettholder_id = ? AND role = ?`, valg.EventID, valg.PuljeID, valg.BillettholderID, valg.Role).Scan(&source); err != nil {
		t.Fatalf("read Player assignment source: %v", err)
	}
	if source != SourceManual {
		t.Fatalf("assignment source: got %q, want %q", source, SourceManual)
	}
	if got := interestCount(t, db, valg.BillettholderID, valg.EventID, valg.PuljeID); got != expectedInterests {
		t.Fatalf("interests: got %d, want %d", got, expectedInterests)
	}
}

func TestTildelBillettholder_AddPlayerConfirmsAndPreservesAllAssignments(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{Given: "a billettholder has one GM assignment and one Player assignment in the pulje", When: "an admin confirms adding a Player to another event", Then: "the warning lists both assignments and confirmation preserves the existing Player and GM"})

	// Given
	expectedHandling := "Legg til som spelar på «Alpha»"
	expectedWarningEvents := []string{"evB", "evC"}
	expectedGMEvents := []string{"evB"}
	expectedPlayerEvents := []string{"evA", "evC"}
	db, _ := testutil.CreateTestDBAndLogger(t, "tildeling_confirm_player_add")
	seedTildelingFixture(t, db, models.AgeGroupDefault, true)
	seedEvent(t, db, "evB", "Bravo", 4, models.PuljeFredagKveld)
	seedEvent(t, db, "evC", "Charlie", 4, models.PuljeFredagKveld)
	seedGM(t, db, "evB", models.PuljeFredagKveld, 1)
	seedManualSeat(t, db, "evC", models.PuljeFredagKveld, 1)
	valg := Tildelingsvalg{
		PuljeID:         models.PuljeFredagKveld,
		EventID:         "evA",
		BillettholderID: 1,
		Role:            models.EventPlayerRolePlayer,
		FraLeggTil:      true,
	}

	// When
	varsel, err := TildelBillettholder(db, valg)

	// Then
	if err != nil {
		t.Fatalf("request Player placement confirmation: %v", err)
	}
	if varsel == nil {
		t.Fatal("expected an assignment warning")
	}
	if !varsel.KanBekrefte {
		t.Fatal("the add flow should permit explicit confirmation")
	}
	if varsel.Bekreftelse == "" {
		t.Fatal("warning must include a bound confirmation token")
	}
	if varsel.Handling != expectedHandling {
		t.Fatalf("warning action: got %q, want %q", varsel.Handling, expectedHandling)
	}
	var warningEvents []string
	for _, tildeling := range varsel.Tildelinger {
		warningEvents = append(warningEvents, tildeling.EventID)
		if tildeling.EventTitle == "" || tildeling.Source == "" || tildeling.InsertedAt == "" {
			t.Fatalf("warning assignment lacks display/state data: %+v", tildeling)
		}
	}
	if !slices.Equal(warningEvents, expectedWarningEvents) {
		t.Fatalf("warning events: got %v, want %v", warningEvents, expectedWarningEvents)
	}
	if got := assignmentEvents(t, db, valg.PuljeID, valg.BillettholderID, models.EventPlayerRolePlayer); !slices.Equal(got, []string{"evC"}) {
		t.Fatalf("unconfirmed request changed Player assignments: %v", got)
	}

	valg.Bekreftelse = varsel.Bekreftelse
	varsel, err = TildelBillettholder(db, valg)
	if err != nil {
		t.Fatalf("confirm Player addition: %v", err)
	}
	if varsel != nil {
		t.Fatalf("confirmed addition returned another warning: %+v", varsel)
	}
	if got := assignmentEvents(t, db, valg.PuljeID, valg.BillettholderID, models.EventPlayerRoleGM); !slices.Equal(got, expectedGMEvents) {
		t.Fatalf("GM assignments after Player addition: got %v, want %v", got, expectedGMEvents)
	}
	if got := assignmentEvents(t, db, valg.PuljeID, valg.BillettholderID, models.EventPlayerRolePlayer); !slices.Equal(got, expectedPlayerEvents) {
		t.Fatalf("Player assignments after addition: got %v, want %v", got, expectedPlayerEvents)
	}
}

func TestTildelBillettholder_DragWithGMOverlapCannotBeConfirmed(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{Given: "a billettholder is GM in the pulje", When: "an admin drags them to a Player seat with every confirmation field set", Then: "the service rejects the overlap without changing assignments"})

	// Given
	expectedGMEvents := []string{"evB"}
	db, _ := testutil.CreateTestDBAndLogger(t, "tildeling_drag_gm_overlap")
	seedTildelingFixture(t, db, models.AgeGroupDefault, true)
	seedEvent(t, db, "evB", "Bravo", 4, models.PuljeFredagKveld)
	seedGM(t, db, "evB", models.PuljeFredagKveld, 1)
	valg := Tildelingsvalg{
		PuljeID:         models.PuljeFredagKveld,
		EventID:         "evA",
		BillettholderID: 1,
		Role:            models.EventPlayerRolePlayer,
		FraLeggTil:      false,
		Bekreftelse:     "kan-ikkje-overstyre",
		AlderBekreftet:  true,
	}

	// When
	varsel, err := TildelBillettholder(db, valg)

	// Then
	if err != nil {
		t.Fatalf("drag GM overlap: %v", err)
	}
	if varsel == nil || varsel.KanBekrefte {
		t.Fatalf("drag GM overlap must return a non-confirmable warning: %+v", varsel)
	}
	if got := assignmentEvents(t, db, valg.PuljeID, valg.BillettholderID, models.EventPlayerRoleGM); !slices.Equal(got, expectedGMEvents) {
		t.Fatalf("GM assignments changed: got %v, want %v", got, expectedGMEvents)
	}
	if got := assignmentEvents(t, db, valg.PuljeID, valg.BillettholderID, models.EventPlayerRolePlayer); len(got) != 0 {
		t.Fatalf("blocked drag stored a Player assignment: %v", got)
	}
}

func TestTildelBillettholder_DragMovesOrdinaryPlayerWithoutAssignmentConfirmation(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{Given: "a billettholder has an ordinary Player assignment", When: "an admin drags them to another event", Then: "the Player assignment moves immediately"})

	// Given
	expectedPlayerEvents := []string{"evA"}
	db, _ := testutil.CreateTestDBAndLogger(t, "tildeling_drag_player_move")
	seedTildelingFixture(t, db, models.AgeGroupDefault, true)
	seedEvent(t, db, "evB", "Bravo", 4, models.PuljeFredagKveld)
	seedManualSeat(t, db, "evB", models.PuljeFredagKveld, 1)
	valg := Tildelingsvalg{
		PuljeID:         models.PuljeFredagKveld,
		EventID:         "evA",
		BillettholderID: 1,
		Role:            models.EventPlayerRolePlayer,
		FraLeggTil:      false,
	}

	// When
	varsel, err := TildelBillettholder(db, valg)

	// Then
	if err != nil {
		t.Fatalf("drag Player: %v", err)
	}
	if varsel != nil {
		t.Fatalf("ordinary Player drag should move without assignment confirmation: %+v", varsel)
	}
	if got := assignmentEvents(t, db, valg.PuljeID, valg.BillettholderID, valg.Role); !slices.Equal(got, expectedPlayerEvents) {
		t.Fatalf("Player assignments: got %v, want %v", got, expectedPlayerEvents)
	}
}

func TestTildelBillettholder_ChangedAssignmentsRequireFreshConfirmation(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{Given: "an admin has a confirmation token for the current assignments", When: "another assignment is added before confirmation", Then: "the stale token writes nothing and a fresh warning lists the changed state"})

	// Given
	const expectedWarningAssignments = 3
	db, _ := testutil.CreateTestDBAndLogger(t, "tildeling_stale_confirmation")
	seedTildelingFixture(t, db, models.AgeGroupDefault, true)
	for _, event := range []struct{ id, title string }{{"evB", "Bravo"}, {"evC", "Charlie"}, {"evD", "Delta"}} {
		seedEvent(t, db, event.id, event.title, 4, models.PuljeFredagKveld)
	}
	seedGM(t, db, "evB", models.PuljeFredagKveld, 1)
	seedManualSeat(t, db, "evC", models.PuljeFredagKveld, 1)
	valg := Tildelingsvalg{PuljeID: models.PuljeFredagKveld, EventID: "evA", BillettholderID: 1, Role: models.EventPlayerRolePlayer, FraLeggTil: true}
	firstWarning, err := TildelBillettholder(db, valg)
	if err != nil || firstWarning == nil {
		t.Fatalf("get initial warning: warning=%+v err=%v", firstWarning, err)
	}
	seedGM(t, db, "evD", models.PuljeFredagKveld, 1)
	valg.Bekreftelse = firstWarning.Bekreftelse

	// When
	freshWarning, err := TildelBillettholder(db, valg)

	// Then
	if err != nil {
		t.Fatalf("submit stale confirmation: %v", err)
	}
	if freshWarning == nil || len(freshWarning.Tildelinger) != expectedWarningAssignments {
		t.Fatalf("expected fresh warning with %d assignments: %+v", expectedWarningAssignments, freshWarning)
	}
	if freshWarning.Bekreftelse == firstWarning.Bekreftelse {
		t.Fatal("changed assignments produced the same confirmation token")
	}
	if got := assignmentEvents(t, db, valg.PuljeID, valg.BillettholderID, models.EventPlayerRolePlayer); !slices.Equal(got, []string{"evC"}) {
		t.Fatalf("stale confirmation changed Player assignment: %v", got)
	}
}

func TestTildelBillettholder_RemovedConflictStillRequiresFreshConfirmation(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{Given: "an admin has a confirmation token for an existing Player assignment", When: "that assignment is removed before confirmation", Then: "the stale token returns the now-empty assignment state before any mutation"})

	// Given
	const expectedCurrentAssignments = 0
	db, _ := testutil.CreateTestDBAndLogger(t, "tildeling_removed_conflict_confirmation")
	seedTildelingFixture(t, db, models.AgeGroupDefault, true)
	seedEvent(t, db, "evB", "Bravo", 4, models.PuljeFredagKveld)
	seedManualSeat(t, db, "evB", models.PuljeFredagKveld, 1)
	valg := Tildelingsvalg{PuljeID: models.PuljeFredagKveld, EventID: "evA", BillettholderID: 1, Role: models.EventPlayerRolePlayer, FraLeggTil: true}
	firstWarning, err := TildelBillettholder(db, valg)
	if err != nil || firstWarning == nil {
		t.Fatalf("get initial warning: warning=%+v err=%v", firstWarning, err)
	}
	if _, err := db.Exec(`DELETE FROM relation_events_players WHERE pulje_id = ? AND billettholder_id = ? AND role = ?`, valg.PuljeID, valg.BillettholderID, valg.Role); err != nil {
		t.Fatalf("remove original conflict: %v", err)
	}
	valg.Bekreftelse = firstWarning.Bekreftelse

	// When
	freshWarning, err := TildelBillettholder(db, valg)

	// Then
	if err != nil {
		t.Fatalf("submit stale confirmation: %v", err)
	}
	if freshWarning == nil || len(freshWarning.Tildelinger) != expectedCurrentAssignments {
		t.Fatalf("expected fresh warning with empty current state: %+v", freshWarning)
	}
	if freshWarning.Bekreftelse == firstWarning.Bekreftelse {
		t.Fatal("removed assignment produced the same confirmation token")
	}
	if got := assignmentEvents(t, db, valg.PuljeID, valg.BillettholderID, valg.Role); len(got) != 0 {
		t.Fatalf("stale confirmation stored a Player assignment: %v", got)
	}
}

func TestTildelBillettholder_AgeFlagCannotBypassStaleTokenAfterConflictRemoved(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{Given: "a minor has a Player assignment and a combined warning for an 18+ destination", When: "the old assignment disappears before the age-confirmed retry", Then: "the stale token returns a fresh age warning without moving the Player"})

	// Given
	const expectedCurrentAssignments = 0
	db, _ := testutil.CreateTestDBAndLogger(t, "tildeling_age_stale_after_removed_conflict")
	seedTildelingFixture(t, db, models.AgeGroupAdultsOnly, false)
	seedEvent(t, db, "evB", "Bravo", 4, models.PuljeFredagKveld)
	seedManualSeat(t, db, "evB", models.PuljeFredagKveld, 1)
	valg := Tildelingsvalg{PuljeID: models.PuljeFredagKveld, EventID: "evA", BillettholderID: 1, Role: models.EventPlayerRolePlayer, FraLeggTil: true}
	firstWarning, err := TildelBillettholder(db, valg)
	if err != nil || firstWarning == nil || firstWarning.Aldersvarsel == "" {
		t.Fatalf("get combined warning: warning=%+v err=%v", firstWarning, err)
	}
	if _, err := db.Exec(`DELETE FROM relation_events_players WHERE pulje_id = ? AND billettholder_id = ? AND role = ?`, valg.PuljeID, valg.BillettholderID, valg.Role); err != nil {
		t.Fatalf("remove original Player assignment: %v", err)
	}
	valg.Bekreftelse = firstWarning.Bekreftelse
	valg.AlderBekreftet = true

	// When
	freshWarning, err := TildelBillettholder(db, valg)

	// Then
	if err != nil {
		t.Fatalf("submit stale combined confirmation: %v", err)
	}
	if freshWarning == nil || freshWarning.Aldersvarsel == "" || len(freshWarning.Tildelinger) != expectedCurrentAssignments {
		t.Fatalf("expected fresh age warning with empty assignment state: %+v", freshWarning)
	}
	if freshWarning.Bekreftelse == firstWarning.Bekreftelse {
		t.Fatal("removed assignment produced the same confirmation token")
	}
	if got := assignmentEvents(t, db, valg.PuljeID, valg.BillettholderID, valg.Role); len(got) != 0 {
		t.Fatalf("stale combined confirmation moved the Player: %v", got)
	}
}

func TestTildelBillettholder_AgeConfirmationRemainsCompatibleWithoutAssignmentConflict(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{Given: "a minor and an unoccupied 18+ event", When: "the admin repeats the placement with the existing age-confirmed flag", Then: "the manual Player assignment is stored"})

	// Given
	const expectedPlayerAssignments = 1
	db, _ := testutil.CreateTestDBAndLogger(t, "tildeling_age_compatibility")
	seedTildelingFixture(t, db, models.AgeGroupAdultsOnly, false)
	valg := Tildelingsvalg{PuljeID: models.PuljeFredagKveld, EventID: "evA", BillettholderID: 1, Role: models.EventPlayerRolePlayer, FraLeggTil: true}
	varsel, err := TildelBillettholder(db, valg)
	if err != nil || varsel == nil || varsel.Aldersvarsel == "" {
		t.Fatalf("expected age warning: warning=%+v err=%v", varsel, err)
	}
	valg.AlderBekreftet = true

	// When
	varsel, err = TildelBillettholder(db, valg)

	// Then
	if err != nil {
		t.Fatalf("confirm age warning: %v", err)
	}
	if varsel != nil {
		t.Fatalf("age-confirmed placement returned warning: %+v", varsel)
	}
	if got := assignmentCountByRole(t, db, valg.PuljeID, valg.BillettholderID, valg.Role); got != expectedPlayerAssignments {
		t.Fatalf("Player assignments: got %d, want %d", got, expectedPlayerAssignments)
	}
}

func TestTildelBillettholder_AgeFlagCannotBypassCombinedGMConflict(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{Given: "a minor is GM in the pulje and the Player destination is 18+", When: "the add request has only the age-confirmed flag", Then: "one combined warning is returned and no Player row is written"})

	// Given
	const expectedWarningAssignments = 1
	db, _ := testutil.CreateTestDBAndLogger(t, "tildeling_combined_warning")
	seedTildelingFixture(t, db, models.AgeGroupAdultsOnly, false)
	seedEvent(t, db, "evB", "Bravo", 4, models.PuljeFredagKveld)
	seedGM(t, db, "evB", models.PuljeFredagKveld, 1)
	valg := Tildelingsvalg{
		PuljeID:         models.PuljeFredagKveld,
		EventID:         "evA",
		BillettholderID: 1,
		Role:            models.EventPlayerRolePlayer,
		FraLeggTil:      true,
		AlderBekreftet:  true,
	}

	// When
	varsel, err := TildelBillettholder(db, valg)

	// Then
	if err != nil {
		t.Fatalf("request combined confirmation: %v", err)
	}
	if varsel == nil || varsel.Aldersvarsel == "" || !varsel.KanBekrefte || len(varsel.Tildelinger) != expectedWarningAssignments {
		t.Fatalf("expected confirmable combined warning with current GM: %+v", varsel)
	}
	if got := assignmentCountByRole(t, db, valg.PuljeID, valg.BillettholderID, valg.Role); got != 0 {
		t.Fatalf("age-only confirmation bypassed GM conflict: %d Player rows", got)
	}

	valg.Bekreftelse = varsel.Bekreftelse
	varsel, err = TildelBillettholder(db, valg)
	if err != nil || varsel != nil {
		t.Fatalf("confirm combined warning: warning=%+v err=%v", varsel, err)
	}
}

func TestTildelBillettholder_ForstevalgStoresHighInterestInSameOperation(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{Given: "an eligible event and an unassigned billettholder", When: "the approval flow assigns the Player as first choice", Then: "the Player pin and high interest are both stored"})

	// Given
	expectedInterest := models.InterestLevelHigh
	db, _ := testutil.CreateTestDBAndLogger(t, "tildeling_forstevalg")
	seedTildelingFixture(t, db, models.AgeGroupDefault, true)
	valg := Tildelingsvalg{
		PuljeID:         models.PuljeFredagKveld,
		EventID:         "evA",
		BillettholderID: 1,
		Role:            models.EventPlayerRolePlayer,
		FraLeggTil:      true,
		Forstevalg:      true,
	}

	// When
	varsel, err := TildelBillettholder(db, valg)

	// Then
	if err != nil || varsel != nil {
		t.Fatalf("assign first choice: warning=%+v err=%v", varsel, err)
	}
	var actualInterest models.InterestLevel
	if err := db.QueryRow(`SELECT interest_level FROM interests WHERE billettholder_id = ? AND event_id = ? AND pulje_id = ?`, valg.BillettholderID, valg.EventID, valg.PuljeID).Scan(&actualInterest); err != nil {
		t.Fatalf("read first-choice interest: %v", err)
	}
	if actualInterest != expectedInterest {
		t.Fatalf("interest: got %q, want %q", actualInterest, expectedInterest)
	}
	if got := assignmentCountByRole(t, db, valg.PuljeID, valg.BillettholderID, valg.Role); got != 1 {
		t.Fatalf("first-choice Player assignments: got %d, want 1", got)
	}
}

func TestFjernTildeling_RemovesOnlySelectedRole(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{Given: "a billettholder is both GM and Player on the same event", When: "the Player assignment is removed", Then: "the GM assignment remains"})

	// Given
	expectedGMEvents := []string{"evA"}
	db, _ := testutil.CreateTestDBAndLogger(t, "fjern_tildeling_role")
	seedTildelingFixture(t, db, models.AgeGroupDefault, true)
	seedGM(t, db, "evA", models.PuljeFredagKveld, 1)
	seedManualSeat(t, db, "evA", models.PuljeFredagKveld, 1)

	// When
	err := FjernTildeling(db, models.PuljeFredagKveld, "evA", 1, models.EventPlayerRolePlayer)

	// Then
	if err != nil {
		t.Fatalf("remove Player assignment: %v", err)
	}
	if got := assignmentEvents(t, db, models.PuljeFredagKveld, 1, models.EventPlayerRoleGM); !slices.Equal(got, expectedGMEvents) {
		t.Fatalf("GM assignments: got %v, want %v", got, expectedGMEvents)
	}
	if got := assignmentEvents(t, db, models.PuljeFredagKveld, 1, models.EventPlayerRolePlayer); len(got) != 0 {
		t.Fatalf("Player assignment remains: %v", got)
	}
}

func TestTildelBillettholder_AddingPlayerAcrossGMEventsPreservesEveryAssignment(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{Given: "a billettholder is GM on three events and manually Player on one of them", When: "an admin confirms adding the Player role to another GM event", Then: "all three GM rows remain with Player pins at both events"})

	// Given
	expectedGMEvents := []string{"evX", "evY", "evZ"}
	expectedPlayerEvents := []string{"evY", "evZ"}
	db, _ := testutil.CreateTestDBAndLogger(t, "tildeling_add_across_gm_events")
	seedPulje(t, db, models.PuljeFredagKveld, "Fredag Kveld", "2026-09-04T18:00:00Z")
	for _, event := range []struct{ id, title string }{{"evX", "X"}, {"evY", "Y"}, {"evZ", "Z"}} {
		seedEvent(t, db, event.id, event.title, 4, models.PuljeFredagKveld)
	}
	seedParticipant(t, db, 1, "Kari", "Nordmann")
	for _, eventID := range expectedGMEvents {
		seedGM(t, db, eventID, models.PuljeFredagKveld, 1)
	}
	seedManualSeat(t, db, "evY", models.PuljeFredagKveld, 1)
	valg := Tildelingsvalg{PuljeID: models.PuljeFredagKveld, EventID: "evZ", BillettholderID: 1, Role: models.EventPlayerRolePlayer, FraLeggTil: true}
	varsel, err := TildelBillettholder(db, valg)
	if err != nil || varsel == nil {
		t.Fatalf("request addition confirmation: warning=%+v err=%v", varsel, err)
	}
	valg.Bekreftelse = varsel.Bekreftelse

	// When
	varsel, err = TildelBillettholder(db, valg)

	// Then
	if err != nil || varsel != nil {
		t.Fatalf("confirm Player addition: warning=%+v err=%v", varsel, err)
	}
	if got := assignmentEvents(t, db, valg.PuljeID, valg.BillettholderID, models.EventPlayerRoleGM); !slices.Equal(got, expectedGMEvents) {
		t.Fatalf("GM assignments: got %v, want %v", got, expectedGMEvents)
	}
	if got := assignmentEvents(t, db, valg.PuljeID, valg.BillettholderID, models.EventPlayerRolePlayer); !slices.Equal(got, expectedPlayerEvents) {
		t.Fatalf("Player assignments: got %v, want %v", got, expectedPlayerEvents)
	}
}

func TestTildelBillettholder_ValidationFailurePreservesExistingAssignments(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{Given: "a billettholder has an existing Player assignment", When: "a new placement fails input or entity validation", Then: "the existing assignment remains unchanged"})

	tests := []struct {
		name     string
		prepare  func(t *testing.T, db *sql.DB)
		valg     Tildelingsvalg
		expected error
	}{
		{
			name:     "invalid role",
			valg:     Tildelingsvalg{PuljeID: models.PuljeFredagKveld, EventID: "evA", BillettholderID: 1, Role: models.EventPlayerRole("Tilskodar"), FraLeggTil: true},
			expected: ErrUgyldigTildeling,
		},
		{
			name:     "missing event",
			valg:     Tildelingsvalg{PuljeID: models.PuljeFredagKveld, EventID: "missing", BillettholderID: 1, Role: models.EventPlayerRolePlayer, FraLeggTil: true},
			expected: sql.ErrNoRows,
		},
		{
			name: "event outside pulje",
			prepare: func(t *testing.T, db *sql.DB) {
				seedEvent(t, db, "outside", "Outside", 4, models.PuljeLordagMorgen)
			},
			valg:     Tildelingsvalg{PuljeID: models.PuljeFredagKveld, EventID: "outside", BillettholderID: 1, Role: models.EventPlayerRolePlayer, FraLeggTil: true},
			expected: ErrUgyldigTildeling,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Given
			expectedPlayerEvents := []string{"evB"}
			db, _ := testutil.CreateTestDBAndLogger(t, "tildeling_validation_"+test.name)
			seedTildelingFixture(t, db, models.AgeGroupDefault, true)
			seedEvent(t, db, "evB", "Bravo", 4, models.PuljeFredagKveld)
			seedManualSeat(t, db, "evB", models.PuljeFredagKveld, 1)
			if test.prepare != nil {
				seedPulje(t, db, models.PuljeLordagMorgen, "Lordag Morgen", "2026-09-05T10:00:00Z")
				test.prepare(t, db)
			}

			// When
			_, err := TildelBillettholder(db, test.valg)

			// Then
			if !errors.Is(err, test.expected) {
				t.Fatalf("error: got %v, want errors.Is(_, %v)", err, test.expected)
			}
			if got := assignmentEvents(t, db, models.PuljeFredagKveld, 1, models.EventPlayerRolePlayer); !slices.Equal(got, expectedPlayerEvents) {
				t.Fatalf("existing Player assignments changed: got %v, want %v", got, expectedPlayerEvents)
			}
		})
	}
}

func TestTildelBillettholder_PublishedPuljePreservesAssignments(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{Given: "a published pulje with an existing Player assignment", When: "an admin requests another assignment", Then: "ErrPuljeCompleted is returned and the existing assignment remains"})

	// Given
	expectedPlayerEvents := []string{"evB"}
	db, _ := testutil.CreateTestDBAndLogger(t, "tildeling_published")
	seedTildelingFixture(t, db, models.AgeGroupDefault, true)
	seedEvent(t, db, "evB", "Bravo", 4, models.PuljeFredagKveld)
	seedManualSeat(t, db, "evB", models.PuljeFredagKveld, 1)
	if _, err := db.Exec(`UPDATE puljer SET status = ? WHERE id = ?`, models.PuljeStatusCompleted, models.PuljeFredagKveld); err != nil {
		t.Fatalf("publish pulje: %v", err)
	}
	valg := Tildelingsvalg{PuljeID: models.PuljeFredagKveld, EventID: "evA", BillettholderID: 1, Role: models.EventPlayerRolePlayer, FraLeggTil: true}

	// When
	_, err := TildelBillettholder(db, valg)

	// Then
	if !errors.Is(err, ErrPuljeCompleted) {
		t.Fatalf("error: got %v, want ErrPuljeCompleted", err)
	}
	if got := assignmentEvents(t, db, valg.PuljeID, valg.BillettholderID, valg.Role); !slices.Equal(got, expectedPlayerEvents) {
		t.Fatalf("published rejection changed Player assignments: got %v, want %v", got, expectedPlayerEvents)
	}
}

func TestGetTildelinger_AbsentBillettholderReturnsNoRows(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{Given: "an existing pulje without the requested billettholder", When: "their assignments are requested", Then: "the service returns sql.ErrNoRows"})

	// Given
	expectedError := sql.ErrNoRows
	db, _ := testutil.CreateTestDBAndLogger(t, "get_tildelinger_missing_holder")
	seedPulje(t, db, models.PuljeFredagKveld, "Fredag Kveld", "2026-09-04T18:00:00Z")

	// When
	_, err := GetTildelinger(db, models.PuljeFredagKveld, 999)

	// Then
	if !errors.Is(err, expectedError) {
		t.Fatalf("error: got %v, want errors.Is(_, sql.ErrNoRows)", err)
	}
}

func seedTildelingFixture(t *testing.T, db *sql.DB, ageGroup models.AgeGroup, isOver18 bool) {
	t.Helper()
	seedPulje(t, db, models.PuljeFredagKveld, "Fredag Kveld", "2026-09-04T18:00:00Z")
	seedEvent(t, db, "evA", "Alpha", 4, models.PuljeFredagKveld)
	if _, err := db.Exec(`UPDATE events SET age_group = ? WHERE id = 'evA'`, ageGroup); err != nil {
		t.Fatalf("set event age group: %v", err)
	}
	seedParticipant(t, db, 1, "Kari", "Nordmann")
	if _, err := db.Exec(`UPDATE billettholdere SET is_over_18 = ? WHERE id = 1`, isOver18); err != nil {
		t.Fatalf("set billettholder age: %v", err)
	}
}

func interestCount(t *testing.T, db *sql.DB, billettholderID int, eventID string, pulje models.Pulje) int {
	t.Helper()
	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM interests WHERE billettholder_id = ? AND event_id = ? AND pulje_id = ?`, billettholderID, eventID, pulje).Scan(&count); err != nil {
		t.Fatalf("count interests: %v", err)
	}
	return count
}

func assignmentEvents(t *testing.T, db *sql.DB, pulje models.Pulje, billettholderID int, role models.EventPlayerRole) []string {
	t.Helper()
	rows, err := db.Query(`SELECT event_id FROM relation_events_players WHERE pulje_id = ? AND billettholder_id = ? AND role = ? ORDER BY event_id`, pulje, billettholderID, role)
	if err != nil {
		t.Fatalf("query assignment events: %v", err)
	}
	defer rows.Close()
	var eventIDs []string
	for rows.Next() {
		var eventID string
		if err := rows.Scan(&eventID); err != nil {
			t.Fatalf("scan assignment event: %v", err)
		}
		eventIDs = append(eventIDs, eventID)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate assignment events: %v", err)
	}
	return eventIDs
}
