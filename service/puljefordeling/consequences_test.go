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
	consequenceX = 1
	consequenceY = 2
	consequenceZ = 3
)

// seedConsequenceChain sets up the cascade from the admin's point of view:
// placing X on Bravo pushes Y from Bravo (their førstevalg) to Charlie, and Z
// from Charlie to no seat.
func seedConsequenceChain(t *testing.T, db *sql.DB) {
	t.Helper()
	seedPulje(t, db, models.PuljeFredagKveld, "Fredag Kveld", "2026-09-04T18:00:00Z")
	seedEvent(t, db, "evB", "Bravo", 1, models.PuljeFredagKveld)
	seedEvent(t, db, "evC", "Charlie", 1, models.PuljeFredagKveld)
	seedParticipant(t, db, consequenceX, "Xander", "Xu")
	seedParticipant(t, db, consequenceY, "Yngve", "Yri")
	seedParticipant(t, db, consequenceZ, "Zara", "Zahl")
	seedInterest(t, db, consequenceY, "evB", models.PuljeFredagKveld, models.InterestLevelHigh)
	seedInterest(t, db, consequenceY, "evC", models.PuljeFredagKveld, models.InterestLevelMedium)
	seedInterest(t, db, consequenceZ, "evC", models.PuljeFredagKveld, models.InterestLevelLow)
}

func TestPreviewAssignment_ListsCascadeWithoutSaving(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at Y har Bravo og Z har Charlie, og en admin vil plassere X på Bravo.",
		When:  "Når tildelingen forhåndsvises.",
		Then:  "Så skal forhåndsvisningen vise at X får Bravo, at Y flyttes til Charlie uten førstevalg og at Z mister plassen, uten at noe lagres.",
	})

	// Given
	expected := []Consequence{
		{BillettholderID: consequenceY, Name: "Yngve Yri",
			From: Seat{EventID: "evB", EventTitle: "Bravo", Level: models.InterestLevelHigh, Forstevalg: true},
			To:   Seat{EventID: "evC", EventTitle: "Charlie", Level: models.InterestLevelMedium}},
		{BillettholderID: consequenceZ, Name: "Zara Zahl",
			From: Seat{EventID: "evC", EventTitle: "Charlie", Level: models.InterestLevelLow}},
	}
	expectedOwn := Consequence{BillettholderID: consequenceX, Name: "Xander Xu",
		To: Seat{EventID: "evB", EventTitle: "Bravo", Level: models.InterestLevelNone}}
	db, _ := testutil.CreateTestDBAndLogger(t, "consequences_cascade")
	seedConsequenceChain(t, db)
	assignment := Tildelingsvalg{PuljeID: models.PuljeFredagKveld, EventID: "evB", BillettholderID: consequenceX, Role: models.EventPlayerRolePlayer, FraLeggTil: true}

	// When
	warning, err := PreviewAssignment(db, assignment)

	// Then
	if err != nil {
		t.Fatalf("PreviewAssignment: %v", err)
	}
	if warning == nil || warning.Consequences == nil || !warning.KanBekrefte || warning.Bekreftelse == "" {
		t.Fatalf("expected a confirmable preview with consequences, got %+v", warning)
	}
	if !slices.Equal(warning.Consequences.Changes, expected) {
		t.Fatalf("expected consequences\n%+v\ngot\n%+v", expected, warning.Consequences.Changes)
	}
	if warning.Consequences.Own != expectedOwn {
		t.Fatalf("expected Xander's own change %+v, got %+v", expectedOwn, warning.Consequences.Own)
	}
	if got := testutil.QueryInt(t, db, `SELECT COUNT(*) FROM relation_events_players`); got != 0 {
		t.Fatalf("preview must not save anything, found %d assignments", got)
	}
}

func TestPreviewAssignment_WithoutSideEffectsHasNoConsequences(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt et arrangement med ledige plasser og én interessert deltaker.",
		When:  "Når en annen billettholder forhåndsvises på arrangementet.",
		Then:  "Så skal forhåndsvisningen kunne bekreftes, uten at andre får endret plass.",
	})

	// Given
	db, _ := testutil.CreateTestDBAndLogger(t, "consequences_none")
	seedPulje(t, db, models.PuljeFredagKveld, "Fredag Kveld", "2026-09-04T18:00:00Z")
	seedEvent(t, db, "evB", "Bravo", 4, models.PuljeFredagKveld)
	seedParticipant(t, db, consequenceX, "Xander", "Xu")
	seedParticipant(t, db, consequenceY, "Yngve", "Yri")
	seedInterest(t, db, consequenceY, "evB", models.PuljeFredagKveld, models.InterestLevelHigh)
	assignment := Tildelingsvalg{PuljeID: models.PuljeFredagKveld, EventID: "evB", BillettholderID: consequenceX, Role: models.EventPlayerRolePlayer, FraLeggTil: true}

	// When
	warning, err := PreviewAssignment(db, assignment)

	// Then
	if err != nil {
		t.Fatalf("PreviewAssignment: %v", err)
	}
	if warning == nil || warning.Consequences == nil || !warning.KanBekrefte {
		t.Fatalf("expected a confirmable preview, got %+v", warning)
	}
	if len(warning.Consequences.Changes) != 0 {
		t.Fatalf("expected no consequences for others, got %+v", warning.Consequences.Changes)
	}
}

func TestPreviewAssignment_ConfirmationSavesTheSameChange(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en forhåndsvist tildeling av X på Bravo.",
		When:  "Når admin bekrefter med bekreftelsen fra forhåndsvisningen.",
		Then:  "Så skal X lagres som manuelt plassert på Bravo.",
	})

	// Given
	db, _ := testutil.CreateTestDBAndLogger(t, "consequences_confirm")
	seedConsequenceChain(t, db)
	assignment := Tildelingsvalg{PuljeID: models.PuljeFredagKveld, EventID: "evB", BillettholderID: consequenceX, Role: models.EventPlayerRolePlayer, FraLeggTil: true}
	preview, err := PreviewAssignment(db, assignment)
	if err != nil {
		t.Fatalf("PreviewAssignment: %v", err)
	}
	assignment.Bekreftelse = preview.Bekreftelse

	// When
	warning, err := TildelBillettholder(db, assignment)

	// Then
	if err != nil || warning != nil {
		t.Fatalf("expected the confirmed change to be saved, got warning %+v, err %v", warning, err)
	}
	if got := assignmentEvents(t, db, models.PuljeFredagKveld, consequenceX, models.EventPlayerRolePlayer); !slices.Equal(got, []string{"evB"}) {
		t.Fatalf("expected X pinned on Bravo, got %v", got)
	}
}

func TestPreviewRemoval_ListsWhoGetsTheFreedSeat(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at X er manuelt plassert på Bravo, slik at Y har Charlie og Z står uten plass.",
		When:  "Når fjerning av den manuelle plassen forhåndsvises.",
		Then:  "Så skal Y flyttes tilbake til Bravo og Z få Charlie, uten at plassen fjernes ennå.",
	})

	// Given
	expectedNames := []string{"Yngve Yri", "Zara Zahl"}
	db, _ := testutil.CreateTestDBAndLogger(t, "consequences_remove")
	seedConsequenceChain(t, db)
	seedManualSeat(t, db, "evB", models.PuljeFredagKveld, consequenceX)

	// When
	warning, err := PreviewRemoval(db, models.PuljeFredagKveld, "evB", consequenceX, models.EventPlayerRolePlayer)

	// Then
	if err != nil {
		t.Fatalf("PreviewRemoval: %v", err)
	}
	if warning.BillettholderName != "Xander Xu" || warning.EventTitle != "Bravo" {
		t.Fatalf("expected the preview to describe removing Xander from Bravo, got %+v", warning)
	}
	var names []string
	for _, k := range warning.Consequences.Changes {
		names = append(names, k.Name)
	}
	if !slices.Equal(names, expectedNames) {
		t.Fatalf("expected %v to change seats, got %+v", expectedNames, warning.Consequences.Changes)
	}
	if got := assignmentEvents(t, db, models.PuljeFredagKveld, consequenceX, models.EventPlayerRolePlayer); !slices.Equal(got, []string{"evB"}) {
		t.Fatalf("preview must keep the manual seat, got %v", got)
	}
}

func TestTildeling_DropOnOwnEventDoesNothing(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at fordelingen har gitt Yngve plass på Bravo.",
		When:  "Når admin drar Yngve og slipper ham på Bravo igjen.",
		Then:  "Så skal det verken spørres eller lagres noe.",
	})

	// Given
	db, _ := testutil.CreateTestDBAndLogger(t, "consequences_drop_same_event")
	seedConsequenceChain(t, db)
	assignment := Tildelingsvalg{PuljeID: models.PuljeFredagKveld, EventID: "evB", BillettholderID: consequenceY, Role: models.EventPlayerRolePlayer, FraEventID: "evB"}

	// When
	preview, previewErr := PreviewAssignment(db, assignment)
	warning, saveErr := TildelBillettholder(db, assignment)

	// Then
	if previewErr != nil || preview != nil {
		t.Fatalf("expected no preview for a drop on the same event, got %+v, %v", preview, previewErr)
	}
	if saveErr != nil || warning != nil {
		t.Fatalf("expected the drop to be a no-op, got %+v, %v", warning, saveErr)
	}
	if got := testutil.QueryInt(t, db, `SELECT COUNT(*) FROM relation_events_players`); got != 0 {
		t.Fatalf("drop on the same event saved %d assignments", got)
	}
}
