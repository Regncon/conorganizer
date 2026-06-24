package admin

import (
	"database/sql"
	"strings"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/Regncon/conorganizer/testutil/templtest"
)

func seedTabPulje(t *testing.T, db *sql.DB, id models.Pulje, name string, status models.PuljeStatus, startAt string) {
	t.Helper()
	_, err := db.Exec(
		`INSERT INTO puljer (id, name, status, start_at, end_at) VALUES (?, ?, ?, ?, ?)`,
		string(id), name, string(status), startAt, startAt,
	)
	if err != nil {
		t.Fatalf("seed pulje %s: %v", id, err)
	}
}

func seedTabEventWithInterest(t *testing.T, db *sql.DB, eventID, title string, pulje models.Pulje) {
	t.Helper()
	if _, err := db.Exec(
		`INSERT INTO events (id, title, intro, description, host_name, email, phone_number, max_players)
		 VALUES (?, ?, '', '', '', '', '', 4)`,
		eventID, title,
	); err != nil {
		t.Fatalf("seed event %s: %v", eventID, err)
	}
	if _, err := db.Exec(
		`INSERT INTO relation_event_puljer (event_id, pulje_id, is_in_pulje) VALUES (?, ?, 1)`,
		eventID, string(pulje),
	); err != nil {
		t.Fatalf("place event %s in %s: %v", eventID, pulje, err)
	}
	if _, err := db.Exec(
		`INSERT INTO billettholdere (id, first_name, last_name, ticket_type_id, ticket_type, order_id, ticket_id)
		 VALUES (1, 'Kari', 'Nordmann', 0, '', 0, 1)`,
	); err != nil {
		t.Fatalf("seed participant: %v", err)
	}
	if _, err := db.Exec(
		`INSERT INTO interests (billettholder_id, event_id, pulje_id, interest_level) VALUES (1, ?, ?, ?)`,
		eventID, string(pulje), string(models.InterestLevelHigh),
	); err != nil {
		t.Fatalf("seed interest: %v", err)
	}
}

func TestPuljefordelingTabContent_RendersPuljeEventsAndStatusToggles(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en pulje med et arrangement og en interessert deltaker.",
		When:  "Når puljefordeling-fanen rendres.",
		Then:  "Så skal arrangementet, deltakeren og status-bryterne vises.",
	})

	// Given
	expectedTextParts := []string{
		"Fredag Kveld",
		"Drager og Fangehull",
		"Kari Nordmann",
		"Puljefordeling lukket",
		"Puljefordeling publisert",
	}
	db, logger := testutil.CreateTestDBAndLogger(t, "puljefordeling_tab_content")
	seedTabPulje(t, db, models.PuljeFredagKveld, "Fredag Kveld", models.PuljeStatusOpen, "2026-01-01 18:00")
	seedTabEventWithInterest(t, db, "drager-og-fangehull", "Drager og Fangehull", models.PuljeFredagKveld)

	// When
	doc := templtest.Render(t, PuljefordelingTabContent(db, logger, models.PuljeFredagKveld))
	actualText := strings.Join(templtest.CollectTexts(doc, "#puljefordeling-tab"), " ")

	// Then
	for _, expectedTextPart := range expectedTextParts {
		if !strings.Contains(actualText, expectedTextPart) {
			t.Fatalf("expected tab text to contain %q\nactual text: %s", expectedTextPart, actualText)
		}
	}
}

func TestPuljefordelingTabContent_ShowsRunningUnsatisfiedCount(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en deltaker hvis eneste førstevalg ligger i en senere pulje.",
		When:  "Når en tidligere og den senere puljen rendres.",
		Then:  "Så skal den tidligere puljen telle deltakeren som uten førstevalg, og den senere som fornøyd.",
	})

	// Given: one participant whose only first choice is in the later pulje.
	db, logger := testutil.CreateTestDBAndLogger(t, "puljefordeling_tab_unsatisfied")
	seedTabPulje(t, db, models.PuljeFredagKveld, "Fredag Kveld", models.PuljeStatusOpen, "2026-01-01 18:00")
	seedTabPulje(t, db, models.PuljeLordagKveld, "Lørdag Kveld", models.PuljeStatusOpen, "2026-01-02 18:00")
	seedTabEventWithInterest(t, db, "kveldsspill", "Kveldsspill", models.PuljeLordagKveld)

	// When / Then: earlier pulje — participant not yet satisfied.
	fredag := strings.Join(
		templtest.CollectTexts(templtest.Render(t, PuljefordelingTabContent(db, logger, models.PuljeFredagKveld)), "#puljefordeling-tab"),
		" ",
	)
	if !strings.Contains(fredag, "1 uten førstevalg så langt") {
		t.Fatalf("expected Fredag tab to report 1 still without first choice\nactual text: %s", fredag)
	}

	// When / Then: later pulje — participant gets their first choice.
	lordag := strings.Join(
		templtest.CollectTexts(templtest.Render(t, PuljefordelingTabContent(db, logger, models.PuljeLordagKveld)), "#puljefordeling-tab"),
		" ",
	)
	if !strings.Contains(lordag, "0 uten førstevalg så langt") {
		t.Fatalf("expected Lørdag tab to report 0 still without first choice\nactual text: %s", lordag)
	}
}

func TestPuljefordelingTabContent_WarnsWhenEventHasNoDm(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt et arrangement i en pulje uten tildelt spilleder.",
		When:  "Når puljefordeling-fanen rendres.",
		Then:  "Så skal arrangementet markeres som at det mangler spilleder.",
	})

	// Given: an event in the pulje with no GM assigned.
	db, logger := testutil.CreateTestDBAndLogger(t, "puljefordeling_missing_dm")
	seedTabPulje(t, db, models.PuljeFredagKveld, "Fredag Kveld", models.PuljeStatusOpen, "2026-01-01 18:00")
	testutil.MustExec(t, db, `INSERT INTO events (id, title, intro, description, host_name, email, phone_number, max_players)
		VALUES ('evA','Alpha','','','','','',4)`)
	testutil.MustExec(t, db, `INSERT INTO relation_event_puljer (event_id, pulje_id, is_in_pulje) VALUES ('evA',?,1)`, string(models.PuljeFredagKveld))

	// When
	doc := templtest.Render(t, PuljefordelingTabContent(db, logger, models.PuljeFredagKveld))
	text := strings.Join(templtest.CollectTexts(doc, "#puljefordeling-tab"), " ")

	// Then
	if !strings.Contains(text, "Mangler spilleder") {
		t.Fatalf("event without a DM should be flagged 'Mangler spilleder'\nactual text: %s", text)
	}
}

func TestPuljefordelingTabContent_PublishedHidesAssignmentControls(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en publisert (Completed) pulje med en manuell plassering.",
		When:  "Når puljefordeling-fanen rendres.",
		Then:  "Så skal verken «legg til» eller «fjern»-kontrollene vises.",
	})

	db, logger := testutil.CreateTestDBAndLogger(t, "puljefordeling_published_readonly")
	seedTabPulje(t, db, models.PuljeFredagKveld, "Fredag Kveld", models.PuljeStatusCompleted, "2026-01-01 18:00")
	testutil.MustExec(t, db, `INSERT INTO events (id, title, intro, description, host_name, email, phone_number, max_players)
		VALUES ('evA','Alpha','','','','','',4)`)
	testutil.MustExec(t, db, `INSERT INTO relation_event_puljer (event_id, pulje_id, is_in_pulje) VALUES ('evA',?,1)`, string(models.PuljeFredagKveld))
	testutil.MustExec(t, db, `INSERT INTO billettholdere (id, first_name, last_name, ticket_type_id, ticket_type, order_id, ticket_id)
		VALUES (1,'Kari','Nordmann',0,'',0,1)`)
	testutil.MustExec(t, db, `INSERT INTO relation_events_players (event_id, pulje_id, billettholder_id, role, source)
		VALUES ('evA',?,1,'Player','manual')`, string(models.PuljeFredagKveld))

	doc := templtest.Render(t, PuljefordelingTabContent(db, logger, models.PuljeFredagKveld))

	if n := doc.Find(".pulje-add").Length(); n != 0 {
		t.Errorf("published pulje must not show the add button, found %d", n)
	}
	if n := doc.Find(".pulje-remove").Length(); n != 0 {
		t.Errorf("published pulje must not show remove controls, found %d", n)
	}
}

func TestPuljefordelingTabContent_PlayerTilesDragToEventBoxes(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en pulje med et seatet arrangement.",
		When:  "Når fanen rendres.",
		Then:  "Så skal spillerflisene kunne dras og arrangementsboksene være slippmål.",
	})

	db, logger := testutil.CreateTestDBAndLogger(t, "puljefordeling_dragdrop")
	seedTabPulje(t, db, models.PuljeFredagKveld, "Fredag Kveld", models.PuljeStatusOpen, "2026-01-01 18:00")
	seedTabEventWithInterest(t, db, "evA", "Alpha", models.PuljeFredagKveld)

	doc := templtest.Render(t, PuljefordelingTabContent(db, logger, models.PuljeFredagKveld))

	tile := doc.Find(".pulje-players li[draggable='true']")
	if tile.Length() == 0 {
		t.Fatal("expected a draggable player tile")
	}
	if got := tile.AttrOr("data-on:dragstart", ""); !strings.Contains(got, "$draggedBillettholderId = 1") {
		t.Errorf("tile dragstart should set the dragged billettholder id; got %q", got)
	}

	drop := doc.Find(".pulje-event").AttrOr("data-on:drop__prevent", "")
	if !strings.Contains(drop, "/admin/api/puljefordeling/assign") {
		t.Errorf("event box should be a drop target posting to the assign endpoint; got %q", drop)
	}
}

func TestPuljefordelingTabContent_PublishedTilesNotDraggable(t *testing.T) {
	db, logger := testutil.CreateTestDBAndLogger(t, "puljefordeling_dragdrop_published")
	seedTabPulje(t, db, models.PuljeFredagKveld, "Fredag Kveld", models.PuljeStatusCompleted, "2026-01-01 18:00")
	seedTabEventWithInterest(t, db, "evA", "Alpha", models.PuljeFredagKveld)

	doc := templtest.Render(t, PuljefordelingTabContent(db, logger, models.PuljeFredagKveld))

	if n := doc.Find(".pulje-players li[draggable='true']").Length(); n != 0 {
		t.Errorf("published pulje must not allow dragging, found %d draggable tiles", n)
	}
}

func TestPuljeStatusToggles_ReflectLockedAndCompletedState(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en pulje som er publisert (Completed).",
		When:  "Når status-bryterne rendres.",
		Then:  "Så skal begge bryterne være avkrysset.",
	})

	// Given
	row := models.PuljeRow{ID: models.PuljeFredagKveld, Name: "Fredag Kveld", Status: models.PuljeStatusCompleted}

	// When
	doc := templtest.Render(t, puljeStatusToggles(row))

	// Then
	checked := doc.Find("input[type=checkbox][checked]")
	if checked.Length() != 2 {
		t.Fatalf("expected both toggles checked for Completed pulje, got %d checked", checked.Length())
	}
}
