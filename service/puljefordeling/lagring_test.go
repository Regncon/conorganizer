package puljefordeling

import (
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/bdd"
)

func TestLagringsstatus_UnsavedPuljeIsFirstSave(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en pulje der fordelingen aldri er lagret.",
		When:  "Når lagringsstatusen hentes.",
		Then:  "Så skal den vise første lagring med alle som får plass som endringer.",
	})

	// Given
	db, _ := testutil.CreateTestDBAndLogger(t, "lagring_first_save")
	seedKonsekvensKjede(t, db)

	// When
	status, err := Lagringsstatus(db, models.PuljeFredagKveld)

	// Then
	if err != nil {
		t.Fatalf("Lagringsstatus: %v", err)
	}
	if !status.ForsteLagring || !status.HarEndringer() || len(status.Endringer) != 2 {
		t.Fatalf("expected a first save with two new seats, got %+v", status)
	}
	if status.Bekreftelse == "" {
		t.Fatal("expected a confirmation token for the changes")
	}
}

func TestLagringsstatus_SavedDistributionHasNoChanges(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en pulje der fordelingen nettopp er lagret.",
		When:  "Når lagringsstatusen hentes.",
		Then:  "Så skal det ikke være noe å lagre.",
	})

	// Given
	db, _ := testutil.CreateTestDBAndLogger(t, "lagring_saved")
	seedKonsekvensKjede(t, db)
	if err := CommitDistribution(db, models.PuljeFredagKveld); err != nil {
		t.Fatalf("CommitDistribution: %v", err)
	}

	// When
	status, err := Lagringsstatus(db, models.PuljeFredagKveld)

	// Then
	if err != nil {
		t.Fatalf("Lagringsstatus: %v", err)
	}
	if status.ForsteLagring || status.HarEndringer() {
		t.Fatalf("expected nothing to save, got %+v", status)
	}
}

func TestLagringsstatus_ListsChangesSinceLastSave(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en lagret fordeling der Yngve har Bravo og Zara har Charlie.",
		When:  "Når Xander deretter plasseres manuelt på Bravo.",
		Then:  "Så skal lagringsstatusen vise hvem som får en annen plass enn den lagrede; den manuelle plassen er allerede lagret.",
	})

	// Given
	expected := []Konsekvens{
		{BillettholderID: konsekvensY, Name: "Yngve Yri",
			Fra: Plass{EventID: "evB", EventTitle: "Bravo", Level: models.InterestLevelHigh},
			Til: Plass{EventID: "evC", EventTitle: "Charlie", Level: models.InterestLevelMedium}},
		{BillettholderID: konsekvensZ, Name: "Zara Zahl",
			Fra: Plass{EventID: "evC", EventTitle: "Charlie", Level: models.InterestLevelLow}},
	}
	db, _ := testutil.CreateTestDBAndLogger(t, "lagring_changes")
	seedKonsekvensKjede(t, db)
	if err := CommitDistribution(db, models.PuljeFredagKveld); err != nil {
		t.Fatalf("CommitDistribution: %v", err)
	}
	seedManualSeat(t, db, "evB", models.PuljeFredagKveld, konsekvensX)

	// When
	status, err := Lagringsstatus(db, models.PuljeFredagKveld)

	// Then
	if err != nil {
		t.Fatalf("Lagringsstatus: %v", err)
	}
	if status.ForsteLagring {
		t.Fatal("expected a later save, not the first")
	}
	if len(status.Endringer) != len(expected) {
		t.Fatalf("expected %d changes, got %+v", len(expected), status.Endringer)
	}
	for i := range expected {
		if status.Endringer[i] != expected[i] {
			t.Errorf("change %d: expected %+v, got %+v", i, expected[i], status.Endringer[i])
		}
	}
}
