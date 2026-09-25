package puljefordeling

import (
	"database/sql"
	"slices"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/bdd"
)

const (
	konsekvensX = 1
	konsekvensY = 2
	konsekvensZ = 3
)

// seedKonsekvensKjede sets up the cascade from the admin's point of view:
// placing X on Bravo pushes Y from Bravo (their førstevalg) to Charlie, and Z
// from Charlie to no seat.
func seedKonsekvensKjede(t *testing.T, db *sql.DB) {
	t.Helper()
	seedPulje(t, db, models.PuljeFredagKveld, "Fredag Kveld", "2026-09-04T18:00:00Z")
	seedEvent(t, db, "evB", "Bravo", 1, models.PuljeFredagKveld)
	seedEvent(t, db, "evC", "Charlie", 1, models.PuljeFredagKveld)
	seedParticipant(t, db, konsekvensX, "Xander", "Xu")
	seedParticipant(t, db, konsekvensY, "Yngve", "Yri")
	seedParticipant(t, db, konsekvensZ, "Zara", "Zahl")
	seedInterest(t, db, konsekvensY, "evB", models.PuljeFredagKveld, models.InterestLevelHigh)
	seedInterest(t, db, konsekvensY, "evC", models.PuljeFredagKveld, models.InterestLevelMedium)
	seedInterest(t, db, konsekvensZ, "evC", models.PuljeFredagKveld, models.InterestLevelLow)
}

func TestForhandsvisTildeling_ListsCascadeWithoutSaving(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at Y har Bravo og Z har Charlie, og en admin vil plassere X på Bravo.",
		When:  "Når tildelingen forhåndsvises.",
		Then:  "Så skal forhåndsvisningen vise at Y flyttes til Charlie uten førstevalg og at Z mister plassen, uten at noe lagres.",
	})

	// Given
	expected := []Konsekvens{
		{BillettholderID: konsekvensY, Name: "Yngve Yri",
			Fra: Plass{EventID: "evB", EventTitle: "Bravo", Level: models.InterestLevelHigh, Forstevalg: true},
			Til: Plass{EventID: "evC", EventTitle: "Charlie", Level: models.InterestLevelMedium}},
		{BillettholderID: konsekvensZ, Name: "Zara Zahl",
			Fra: Plass{EventID: "evC", EventTitle: "Charlie", Level: models.InterestLevelLow}},
	}
	db, _ := testutil.CreateTestDBAndLogger(t, "konsekvenser_cascade")
	seedKonsekvensKjede(t, db)
	valg := Tildelingsvalg{PuljeID: models.PuljeFredagKveld, EventID: "evB", BillettholderID: konsekvensX, Role: models.EventPlayerRolePlayer, FraLeggTil: true}

	// When
	varsel, err := ForhandsvisTildeling(db, valg)

	// Then
	if err != nil {
		t.Fatalf("ForhandsvisTildeling: %v", err)
	}
	if varsel == nil || varsel.Konsekvenser == nil || !varsel.KanBekrefte || varsel.Bekreftelse == "" {
		t.Fatalf("expected a confirmable preview with consequences, got %+v", varsel)
	}
	if !slices.Equal(varsel.Konsekvenser.Endringer, expected) {
		t.Fatalf("expected consequences\n%+v\ngot\n%+v", expected, varsel.Konsekvenser.Endringer)
	}
	if got := testutil.QueryInt(t, db, `SELECT COUNT(*) FROM relation_events_players`); got != 0 {
		t.Fatalf("preview must not save anything, found %d assignments", got)
	}
}

func TestForhandsvisTildeling_WithoutSideEffectsHasNoConsequences(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt et arrangement med ledige plasser og én interessert deltaker.",
		When:  "Når en annen billettholder forhåndsvises på arrangementet.",
		Then:  "Så skal forhåndsvisningen kunne bekreftes, uten at andre får endret plass.",
	})

	// Given
	db, _ := testutil.CreateTestDBAndLogger(t, "konsekvenser_none")
	seedPulje(t, db, models.PuljeFredagKveld, "Fredag Kveld", "2026-09-04T18:00:00Z")
	seedEvent(t, db, "evB", "Bravo", 4, models.PuljeFredagKveld)
	seedParticipant(t, db, konsekvensX, "Xander", "Xu")
	seedParticipant(t, db, konsekvensY, "Yngve", "Yri")
	seedInterest(t, db, konsekvensY, "evB", models.PuljeFredagKveld, models.InterestLevelHigh)
	valg := Tildelingsvalg{PuljeID: models.PuljeFredagKveld, EventID: "evB", BillettholderID: konsekvensX, Role: models.EventPlayerRolePlayer, FraLeggTil: true}

	// When
	varsel, err := ForhandsvisTildeling(db, valg)

	// Then
	if err != nil {
		t.Fatalf("ForhandsvisTildeling: %v", err)
	}
	if varsel == nil || varsel.Konsekvenser == nil || !varsel.KanBekrefte {
		t.Fatalf("expected a confirmable preview, got %+v", varsel)
	}
	if len(varsel.Konsekvenser.Endringer) != 0 {
		t.Fatalf("expected no consequences for others, got %+v", varsel.Konsekvenser.Endringer)
	}
}

func TestForhandsvisTildeling_ConfirmationSavesTheSameChange(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en forhåndsvist tildeling av X på Bravo.",
		When:  "Når admin bekrefter med bekreftelsen fra forhåndsvisningen.",
		Then:  "Så skal X lagres som manuelt plassert på Bravo.",
	})

	// Given
	db, _ := testutil.CreateTestDBAndLogger(t, "konsekvenser_confirm")
	seedKonsekvensKjede(t, db)
	valg := Tildelingsvalg{PuljeID: models.PuljeFredagKveld, EventID: "evB", BillettholderID: konsekvensX, Role: models.EventPlayerRolePlayer, FraLeggTil: true}
	preview, err := ForhandsvisTildeling(db, valg)
	if err != nil {
		t.Fatalf("ForhandsvisTildeling: %v", err)
	}
	valg.Bekreftelse = preview.Bekreftelse

	// When
	varsel, err := TildelBillettholder(db, valg)

	// Then
	if err != nil || varsel != nil {
		t.Fatalf("expected the confirmed change to be saved, got varsel %+v, err %v", varsel, err)
	}
	if got := assignmentEvents(t, db, models.PuljeFredagKveld, konsekvensX, models.EventPlayerRolePlayer); !slices.Equal(got, []string{"evB"}) {
		t.Fatalf("expected X pinned on Bravo, got %v", got)
	}
}

func TestForhandsvisFjerning_ListsWhoGetsTheFreedSeat(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at X er manuelt plassert på Bravo, slik at Y har Charlie og Z står uten plass.",
		When:  "Når fjerning av den manuelle plassen forhåndsvises.",
		Then:  "Så skal Y flyttes tilbake til Bravo og Z få Charlie, uten at plassen fjernes ennå.",
	})

	// Given
	expectedNames := []string{"Yngve Yri", "Zara Zahl"}
	db, _ := testutil.CreateTestDBAndLogger(t, "konsekvenser_remove")
	seedKonsekvensKjede(t, db)
	seedManualSeat(t, db, "evB", models.PuljeFredagKveld, konsekvensX)

	// When
	varsel, err := ForhandsvisFjerning(db, models.PuljeFredagKveld, "evB", konsekvensX, models.EventPlayerRolePlayer)

	// Then
	if err != nil {
		t.Fatalf("ForhandsvisFjerning: %v", err)
	}
	if varsel.BillettholderNavn != "Xander Xu" || varsel.EventTitle != "Bravo" {
		t.Fatalf("expected the preview to describe removing Xander from Bravo, got %+v", varsel)
	}
	var names []string
	for _, k := range varsel.Konsekvenser.Endringer {
		names = append(names, k.Name)
	}
	if !slices.Equal(names, expectedNames) {
		t.Fatalf("expected %v to change seats, got %+v", expectedNames, varsel.Konsekvenser.Endringer)
	}
	if got := assignmentEvents(t, db, models.PuljeFredagKveld, konsekvensX, models.EventPlayerRolePlayer); !slices.Equal(got, []string{"evB"}) {
		t.Fatalf("preview must keep the manual seat, got %v", got)
	}
}
