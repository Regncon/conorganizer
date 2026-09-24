package admin

import (
	"strings"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/service/puljefordeling"
	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/Regncon/conorganizer/testutil/templtest"
)

func TestPuljeForstevalgStats_TilesOpenListsOfParticipants(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en pulje der én fikk førstevalg, to mangler det og én er uten plass.",
		When:  "Når førstevalgsoversikten rendres.",
		Then:  "Så skal hver flis vise antallet og åpne en liste med deltakerne, uten interessepoeng.",
	})

	// Given
	expectedTiles := map[string]string{"fikk-forstevalg": "1", "uten-forstevalg": "2", "uten-plass": "1"}
	result := puljefordeling.EmulatedPulje{
		TotalScore:     90,
		Unassigned:     []string{"Bjørn Berg"},
		FikkForstevalg: []puljefordeling.PuljeDeltaker{{Name: "Anne Aas", EventTitle: "Drager", Level: models.InterestLevelHigh}},
		UtenForstevalg: []puljefordeling.PuljeDeltaker{
			{Name: "Bjørn Berg", WantedForstevalg: true},
			{Name: "Cato Carlsen", EventTitle: "Brett", Level: models.InterestLevelLow},
		},
	}

	// When
	doc := templtest.Render(t, puljeForstevalgStats(result))

	// Then
	for list, expected := range expectedTiles {
		tile := doc.Find(`.pulje-stat[aria-controls="pulje-liste-` + list + `"]`)
		if got := strings.TrimSpace(tile.Find(".pulje-stat-value").Text()); got != expected {
			t.Errorf("tile %s: expected %s, got %q", list, expected, got)
		}
		if action := tile.AttrOr("data-on:click", ""); action != "$_puljeListe = '"+list+"'" {
			t.Errorf("tile %s: expected it to open its list, got %q", list, action)
		}
		if doc.Find("dialog#pulje-liste-"+list).Length() != 1 {
			t.Errorf("expected a dialog for %s", list)
		}
	}
	uten := strings.Join(templtest.CollectTexts(doc, "#pulje-liste-uten-forstevalg li"), " | ")
	for _, part := range []string{"Bjørn Berg", "Ville ha førstevalg her", "Uten plass", "Cato Carlsen", "Brett"} {
		if !strings.Contains(uten, part) {
			t.Errorf("expected uten førstevalg list to contain %q, got %q", part, uten)
		}
	}
	if strings.Contains(doc.Text(), "interessepoeng") {
		t.Error("expected interessepoeng to no longer be shown")
	}
}
