package admin

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/Regncon/conorganizer/testutil/templtest"
)

func TestPuljefordeling_CapacityWarningListsEveryPlayer(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{Given: "Et arrangement med fire plasser har fem manuelt tildelte spillere.", When: "En administrator åpner puljefordelingen.", Then: "Et utvidbart kapasitetsvarsel viser fem av fire plasser og alle spillerne."})
	// Given
	expectedPlayers := 5
	db, _ := tildelingsFixture(t)
	for id := 1; id <= expectedPlayers; id++ {
		if id > 1 {
			testutil.MustExec(t, db, `INSERT INTO billettholdere(id,first_name,last_name,ticket_type_id,ticket_type,order_id,ticket_id,is_over_18) VALUES (?,?,'Nordmann',0,'',0,?,1)`, id, fmt.Sprintf("Spiller %d", id), id)
		}
		testutil.MustExec(t, db, `INSERT INTO relation_events_players(event_id,pulje_id,billettholder_id,role,source) VALUES ('evA','FredagKveld',?,'Player','manual')`, id)
	}

	// When
	doc := templtest.Render(t, PuljefordelingTabContent(db, testutil.NewTestLogger(), models.PuljeFredagKveld, nil))

	// Then
	warning := doc.Find("details[data-kapasitetsvarsel='evA']")
	if warning.Length() != 1 {
		t.Fatalf("forventet ett kapasitetsvarsel, fikk %d", warning.Length())
	}
	for _, text := range []string{"Arrangement X", "5 / 4 spillerplasser"} {
		if !strings.Contains(warning.Find("summary").Text(), text) {
			t.Errorf("kapasitetsvarselet mangler %q", text)
		}
	}
	if got := warning.Find("li").Length(); got != expectedPlayers {
		t.Errorf("kapasitetsvarselet viser %d spillere, vil ha %d", got, expectedPlayers)
	}
	if preserve, _ := warning.Attr("data-preserve-attr"); preserve != "open" {
		t.Error("åpnet kapasitetsvarsel må beholdes under oppdatering")
	}
}

func TestPuljefordeling_OverlapWarningListsEveryAssignment(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{Given: "Kari is GM for three arrangementer and pinned as Player on one of them.", When: "An admin views the pulje distribution.", Then: "One expandable warning names Kari and lists all four roles."})
	// Given
	expectedAssignments := 4
	db, _ := tildelingsFixture(t)
	testutil.MustExec(t, db, `UPDATE billettholdere SET first_name='Kari', last_name='Nordmann' WHERE id=1`)
	for _, eventID := range []string{"evA", "evB", "evC"} {
		testutil.MustExec(t, db, `INSERT INTO relation_events_players(event_id,pulje_id,billettholder_id,role) VALUES (?,'FredagKveld',1,'GM')`, eventID)
	}
	testutil.MustExec(t, db, `INSERT INTO relation_events_players(event_id,pulje_id,billettholder_id,role,source) VALUES ('evB','FredagKveld',1,'Player','manual')`)

	// When
	doc := templtest.Render(t, PuljefordelingTabContent(db, testutil.NewTestLogger(), models.PuljeFredagKveld, nil))

	// Then
	warning := doc.Find("details[data-fordelingsvarsel='1']")
	if warning.Length() != 1 {
		t.Fatalf("expected one expandable warning for Kari, got %d", warning.Length())
	}
	for _, text := range []string{"Kari Nordmann", "3 spillederoppdrag", "1 spillerplass"} {
		if !strings.Contains(warning.Find("summary").Text(), text) {
			t.Errorf("warning summary missing %q", text)
		}
	}
	if got := warning.Find("li").Length(); got != expectedAssignments {
		t.Fatalf("expected %d assignments, got %d", expectedAssignments, got)
	}
	for _, text := range []string{"Arrangement X", "Arrangement Y", "Arrangement Z", "Spilleder", "Spiller"} {
		if !strings.Contains(warning.Find("ul").Text(), text) {
			t.Errorf("assignment list missing %q", text)
		}
	}
	if preserve, _ := warning.Attr("data-preserve-attr"); preserve != "open" {
		t.Fatal("expanded assignments should stay open during live refresh")
	}
}

func TestPuljefordeling_SingleAssignmentHasNoOverlapWarning(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{Given: "A billettholder has one GM assignment in the pulje.", When: "An admin views the distribution.", Then: "No overlap warning is displayed."})
	// Given
	expectedWarnings := 0
	db, _ := tildelingsFixture(t)
	testutil.MustExec(t, db, `INSERT INTO relation_events_players(event_id,pulje_id,billettholder_id,role) VALUES ('evA','FredagKveld',1,'GM')`)

	// When
	doc := templtest.Render(t, PuljefordelingTabContent(db, testutil.NewTestLogger(), models.PuljeFredagKveld, nil))

	// Then
	if got := doc.Find("[data-fordelingsvarsel]").Length(); got != expectedWarnings {
		t.Fatalf("expected no warning, got %d", got)
	}
}
