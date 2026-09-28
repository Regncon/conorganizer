package admin

import (
	"strings"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/Regncon/conorganizer/testutil/templtest"
)

func TestPuljefordeling_AldersvarselForUnder18SattPaaVoksenarrangementAvSolveren(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en billettholder under 18 som er veldig interessert i et 18+-arrangement, uten manuell tildeling.",
		When:  "Når en administrator åpner puljefordelingen.",
		Then:  "Så setter solveren hen på arrangementet, og et aldersvarsel nevner hen og arrangementet.",
	})

	// Given
	expectedName := "Kari Nordmann"
	db, logger := testutil.CreateTestDBAndLogger(t, "puljefordeling_aldersvarsel_solver")
	seedAssignFixture(t, db, models.PuljeFredagKveld, models.AgeGroupAdultsOnly, false)
	testutil.MustExec(t, db, `INSERT INTO interests(billettholder_id,event_id,pulje_id,interest_level) VALUES (1,'evA','FredagKveld',?)`, models.InterestLevelHigh)

	// When
	doc := templtest.Render(t, PuljefordelingTabContent(db, logger, models.PuljeFredagKveld, nil))

	// Then
	warning := doc.Find("details[data-aldersvarsel='evA']")
	if warning.Length() != 1 {
		t.Fatalf("forventet ett aldersvarsel for evA, fikk %d", warning.Length())
	}
	if !strings.Contains(warning.Find("summary").Text(), "Voksenspill") {
		t.Errorf("aldersvarselet mangler arrangementet: %q", warning.Find("summary").Text())
	}
	deltaker := warning.Find("li.aldersvarsel-deltaker")
	if got := deltaker.Find(".aldersvarsel-navn").Text(); got != expectedName {
		t.Errorf("uthevet navn = %q, vil ha %q", got, expectedName)
	}
	merke := deltaker.Find(`.under-18-merke[role="img"][aria-label="Under 18 år på 18+-arrangement"]`)
	if merke.Length() != 1 {
		t.Fatalf("forventet ett under 18-merke ved deltakeren, fikk %d", merke.Length())
	}
	if got := merke.Find("svg").Length(); got != 2 {
		t.Errorf("under 18-merket skal vise 18+-ikonet med forbudsskilt over, fikk %d ikoner", got)
	}
	if preserve, _ := warning.Attr("data-preserve-attr"); preserve != "open" {
		t.Error("åpnet aldersvarsel må beholdes under oppdatering")
	}
}

func TestPuljefordeling_IngenAldersvarselForVoksenPaaVoksenarrangement(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en voksen billettholder som er veldig interessert i et 18+-arrangement.",
		When:  "Når en administrator åpner puljefordelingen.",
		Then:  "Så vises ingen aldersvarsler.",
	})

	// Given
	expectedWarnings := 0
	db, logger := testutil.CreateTestDBAndLogger(t, "puljefordeling_aldersvarsel_voksen")
	seedAssignFixture(t, db, models.PuljeFredagKveld, models.AgeGroupAdultsOnly, true)
	testutil.MustExec(t, db, `INSERT INTO interests(billettholder_id,event_id,pulje_id,interest_level) VALUES (1,'evA','FredagKveld',?)`, models.InterestLevelHigh)

	// When
	doc := templtest.Render(t, PuljefordelingTabContent(db, logger, models.PuljeFredagKveld, nil))

	// Then
	if got := doc.Find("details[data-aldersvarsel]").Length(); got != expectedWarnings {
		t.Errorf("antall aldersvarsler = %d, vil ha %d", got, expectedWarnings)
	}
}
