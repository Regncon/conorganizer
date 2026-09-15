package admin

import (
	"strings"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/Regncon/conorganizer/testutil/templtest"
)

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
