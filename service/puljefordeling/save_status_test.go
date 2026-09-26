package puljefordeling

import (
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
