package puljefordeling

import (
	"slices"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/bdd"
)

func TestSaveStatus_UnsavedPuljeIsFirstSave(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en pulje der fordelingen aldri er lagret.",
		When:  "Når lagringsstatusen hentes.",
		Then:  "Så skal den vise første lagring med alle som får plass som endringer.",
	})

	// Given
	db, _ := testutil.CreateTestDBAndLogger(t, "save_status_first_save")
	seedConsequenceChain(t, db)

	// When
	status, err := LoadSaveStatus(db, models.PuljeFredagKveld)

	// Then
	if err != nil {
		t.Fatalf("LoadSaveStatus: %v", err)
	}
	if !status.FirstSave || !status.HasChanges() || len(status.Changes) != 2 {
		t.Fatalf("expected a first save with two new seats, got %+v", status)
	}
	if status.Confirmation == "" {
		t.Fatal("expected a confirmation token for the changes")
	}
}

func TestSaveStatus_SavedDistributionHasNoChanges(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en pulje der fordelingen nettopp er lagret.",
		When:  "Når lagringsstatusen hentes.",
		Then:  "Så skal det ikke være noe å lagre.",
	})

	// Given
	db, _ := testutil.CreateTestDBAndLogger(t, "save_status_saved")
	seedConsequenceChain(t, db)
	if err := CommitDistribution(db, models.PuljeFredagKveld); err != nil {
		t.Fatalf("CommitDistribution: %v", err)
	}

	// When
	status, err := LoadSaveStatus(db, models.PuljeFredagKveld)

	// Then
	if err != nil {
		t.Fatalf("LoadSaveStatus: %v", err)
	}
	if status.FirstSave || status.HasChanges() {
		t.Fatalf("expected nothing to save, got %+v", status)
	}
}

func TestSaveStatus_ListsChangesSinceLastSave(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en lagret fordeling der Yngve har Bravo og Zara har Charlie.",
		When:  "Når Xander deretter plasseres manuelt på Bravo.",
		Then:  "Så skal lagringsstatusen vise hvem som får en annen plass enn den lagrede; den manuelle plassen er allerede lagret.",
	})

	// Given
	expected := []Consequence{
		{BillettholderID: consequenceY, Name: "Yngve Yri",
			From: Seat{EventID: "evB", EventTitle: "Bravo", Level: models.InterestLevelHigh},
			To:   Seat{EventID: "evC", EventTitle: "Charlie", Level: models.InterestLevelMedium}},
		{BillettholderID: consequenceZ, Name: "Zara Zahl",
			From: Seat{EventID: "evC", EventTitle: "Charlie", Level: models.InterestLevelLow}},
	}
	db, _ := testutil.CreateTestDBAndLogger(t, "save_status_changes")
	seedConsequenceChain(t, db)
	if err := CommitDistribution(db, models.PuljeFredagKveld); err != nil {
		t.Fatalf("CommitDistribution: %v", err)
	}
	seedManualSeat(t, db, "evB", models.PuljeFredagKveld, consequenceX)

	// When
	status, err := LoadSaveStatus(db, models.PuljeFredagKveld)

	// Then
	if err != nil {
		t.Fatalf("LoadSaveStatus: %v", err)
	}
	if status.FirstSave {
		t.Fatal("expected a later save, not the first")
	}
	if len(status.Changes) != len(expected) {
		t.Fatalf("expected %d changes, got %+v", len(expected), status.Changes)
	}
	for i := range expected {
		if status.Changes[i] != expected[i] {
			t.Errorf("change %d: expected %+v, got %+v", i, expected[i], status.Changes[i])
		}
	}
}

func TestSaveStatus_UndoingAManualMoveLeavesNothingToSave(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en lagret fordeling der Yngve har Bravo.",
		When:  "Når admin flytter Yngve manuelt til Charlie, og så fjerner den manuelle plassen igjen.",
		Then:  "Så skal Yngve ha fått tilbake plassen sin i den lagrede fordelingen, og det skal ikke være noe å lagre.",
	})

	// Given
	expectedSeat := "evB"
	db, _ := testutil.CreateTestDBAndLogger(t, "save_status_undo_manual_move")
	seedConsequenceChain(t, db)
	if err := CommitDistribution(db, models.PuljeFredagKveld); err != nil {
		t.Fatalf("CommitDistribution: %v", err)
	}
	move := Tildelingsvalg{PuljeID: models.PuljeFredagKveld, EventID: "evC", BillettholderID: consequenceY, Role: models.EventPlayerRolePlayer, FraEventID: "evB"}
	if warning, err := TildelBillettholder(db, move); err != nil || warning != nil {
		t.Fatalf("move Yngve to Charlie: %+v, %v", warning, err)
	}

	// When
	err := RemoveManualSeat(db, models.PuljeFredagKveld, "evC", consequenceY)

	// Then
	if err != nil {
		t.Fatalf("RemoveManualSeat: %v", err)
	}
	status, err := LoadSaveStatus(db, models.PuljeFredagKveld)
	if err != nil {
		t.Fatalf("LoadSaveStatus: %v", err)
	}
	if status.HasChanges() {
		t.Fatalf("expected nothing to save after undoing the move, got %+v", status.Changes)
	}
	if got := assignmentEvents(t, db, models.PuljeFredagKveld, consequenceY, models.EventPlayerRolePlayer); !slices.Equal(got, []string{expectedSeat}) {
		t.Fatalf("expected Yngve's saved seat to be back on %s, got %v", expectedSeat, got)
	}
}

func TestRemoveManualSeat_BeforeFirstSaveKeepsItUnsaved(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en pulje der fordelingen aldri er lagret, og en manuell plass.",
		When:  "Når den manuelle plassen fjernes.",
		Then:  "Så skal ingen plass lagres, slik at første lagring fortsatt er første lagring.",
	})

	// Given
	db, _ := testutil.CreateTestDBAndLogger(t, "save_status_remove_before_first_save")
	seedConsequenceChain(t, db)
	seedManualSeat(t, db, "evB", models.PuljeFredagKveld, consequenceX)

	// When
	err := RemoveManualSeat(db, models.PuljeFredagKveld, "evB", consequenceX)

	// Then
	if err != nil {
		t.Fatalf("RemoveManualSeat: %v", err)
	}
	if got := testutil.QueryInt(t, db, `SELECT COUNT(*) FROM relation_events_players`); got != 0 {
		t.Fatalf("expected no saved seats before the first save, got %d", got)
	}
	status, err := LoadSaveStatus(db, models.PuljeFredagKveld)
	if err != nil {
		t.Fatalf("LoadSaveStatus: %v", err)
	}
	if !status.FirstSave {
		t.Fatal("expected the pulje to still need its first save")
	}
}

func TestSaveConfirmedDistribution_ClearsManualSeatOnEventNoLongerInPulje(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en lagret fordeling der Xander er plassert manuelt på Bravo.",
		When:  "Når Bravo tas ut av puljen og admin lagrer fordelingen.",
		Then:  "Så skal den manuelle plassen på Bravo fjernes, og det skal ikke være noe mer å lagre.",
	})

	// Given
	db, _ := testutil.CreateTestDBAndLogger(t, "save_status_manual_seat_removed_event")
	seedConsequenceChain(t, db)
	seedManualSeat(t, db, "evB", models.PuljeFredagKveld, consequenceX)
	if err := CommitDistribution(db, models.PuljeFredagKveld); err != nil {
		t.Fatalf("CommitDistribution: %v", err)
	}
	testutil.MustExec(t, db,
		`UPDATE relation_event_puljer SET is_in_pulje = 0 WHERE event_id = 'evB' AND pulje_id = ?`,
		string(models.PuljeFredagKveld))
	status, err := LoadSaveStatus(db, models.PuljeFredagKveld)
	if err != nil || !status.HasChanges() {
		t.Fatalf("expected changes after removing Bravo, got %+v, %v", status, err)
	}

	// When
	retry, err := SaveConfirmedDistribution(db, models.PuljeFredagKveld, status.Confirmation)

	// Then
	if err != nil || retry != nil {
		t.Fatalf("SaveConfirmedDistribution: %+v, %v", retry, err)
	}
	if got := assignmentEvents(t, db, models.PuljeFredagKveld, consequenceX, models.EventPlayerRolePlayer); slices.Contains(got, "evB") {
		t.Fatalf("expected Xander's seat on Bravo to be gone, got %v", got)
	}
	after, err := LoadSaveStatus(db, models.PuljeFredagKveld)
	if err != nil {
		t.Fatalf("LoadSaveStatus: %v", err)
	}
	if after.HasChanges() {
		t.Fatalf("expected nothing to save, got %+v", after.Changes)
	}
}
