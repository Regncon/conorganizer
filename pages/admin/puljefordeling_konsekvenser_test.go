package admin

import (
	"database/sql"
	"net/http"
	"strings"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/service/puljefordeling"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/Regncon/conorganizer/testutil/templtest"
)

// seedKonsekvensRoute: placing Xander on Bravo pushes Yngve from Bravo (his
// førstevalg) to Charlie, and Zara from Charlie to no seat.
func seedKonsekvensRoute(t *testing.T) (*sql.DB, http.Handler) {
	t.Helper()
	db, _ := testutil.CreateTestDBAndLogger(t, t.Name())
	seedTabPulje(t, db, models.PuljeFredagKveld, "Fredag Kveld", models.PuljeStatusOpen, "2026-01-01 18:00")
	for _, ev := range []struct{ id, title string }{{"evB", "Bravo"}, {"evC", "Charlie"}} {
		testutil.MustExec(t, db, `INSERT INTO events (id, title, intro, description, host_name, email, phone_number, max_players, is_in_puljefordeling) VALUES (?, ?, '', '', '', '', '', 1, 1)`, ev.id, ev.title)
		testutil.MustExec(t, db, `INSERT INTO relation_event_puljer (event_id, pulje_id, is_in_pulje) VALUES (?, 'FredagKveld', 1)`, ev.id)
	}
	for _, p := range []struct {
		id          int
		first, last string
	}{{1, "Xander", "Xu"}, {2, "Yngve", "Yri"}, {3, "Zara", "Zahl"}} {
		testutil.MustExec(t, db, `INSERT INTO billettholdere (id, first_name, last_name, ticket_type_id, ticket_type, order_id, ticket_id, is_over_18) VALUES (?, ?, ?, 0, '', 0, ?, 1)`, p.id, p.first, p.last, p.id)
	}
	testutil.MustExec(t, db, `INSERT INTO interests (billettholder_id, event_id, pulje_id, interest_level) VALUES
		(2, 'evB', 'FredagKveld', ?), (2, 'evC', 'FredagKveld', ?), (3, 'evC', 'FredagKveld', ?)`,
		string(models.InterestLevelHigh), string(models.InterestLevelMedium), string(models.InterestLevelLow))
	return db, assignmentRouterFor(t, db)
}

func TestPuljefordelingAssignPreview_AsksWithConsequencesBeforeSaving(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at Yngve har Bravo og Zara har Charlie.",
		When:  "Når en admin drar Xander til Bravo.",
		Then:  "Så skal admin få «Er du sikker?» med hvem som flyttes, og ingenting lagres før det bekreftes.",
	})

	// Given
	expectedParts := []string{"Er du sikker?", "For Xander Xu", "Yngve Yri", "Bravo", "Charlie", "Får ikke førstevalg i denne puljen", "Zara Zahl", "Uten plass"}
	db, router := seedKonsekvensRoute(t)

	// When
	preview := postAssignmentSignals(t, router, http.MethodPost, "/api/puljefordeling/assign/preview", 1, "evB", "FredagKveld", `,"assignmentRole":"Player","assignmentFromAddMenu":true`)

	// Then
	if preview.Code != http.StatusOK {
		t.Fatalf("expected a confirmation dialog, got %d: %s", preview.Code, preview.Body.String())
	}
	for _, part := range expectedParts {
		if !strings.Contains(preview.Body.String(), part) {
			t.Errorf("expected the dialog to mention %q", part)
		}
	}
	if got := testutil.QueryInt(t, db, `SELECT COUNT(*) FROM relation_events_players`); got != 0 {
		t.Fatalf("preview wrote %d assignments", got)
	}

	confirmed := postAssignmentSignals(t, router, http.MethodPost, "/api/puljefordeling/assign", 1, "evB", "FredagKveld", `,"assignmentRole":"Player","assignmentFromAddMenu":true,"assignmentConfirmation":"`+confirmationFromResponse(t, preview)+`"`)
	if confirmed.Code != http.StatusNoContent {
		t.Fatalf("expected the confirmed change to save, got %d: %s", confirmed.Code, confirmed.Body.String())
	}
	if got := testutil.QueryInt(t, db, `SELECT COUNT(*) FROM relation_events_players WHERE billettholder_id = 1 AND event_id = 'evB' AND source = 'manual'`); got != 1 {
		t.Fatal("expected Xander to be pinned on Bravo after confirming")
	}
}

func TestPuljefordelingRemovePreview_AsksWithConsequencesBeforeRemoving(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at Xander er manuelt plassert på Bravo.",
		When:  "Når en admin trykker × på plassen hans.",
		Then:  "Så skal admin få «Er du sikker?» med hvem som får plassen, og bekreftelsen fjerner den manuelle plassen.",
	})

	// Given
	expectedParts := []string{"Er du sikker?", "Fjern den manuelle plassen på «Bravo»", "Yngve Yri", "Zara Zahl", "@delete(&#34;/admin/api/puljefordeling/FredagKveld/evB/1&#34;)"}
	db, router := seedKonsekvensRoute(t)
	testutil.MustExec(t, db, `INSERT INTO relation_events_players (event_id, pulje_id, billettholder_id, role, source) VALUES ('evB', 'FredagKveld', 1, 'Player', 'manual')`)

	// When
	preview := postAssignmentSignals(t, router, http.MethodPost, "/api/puljefordeling/FredagKveld/evB/1/preview", 0, "", "", "")

	// Then
	if preview.Code != http.StatusOK {
		t.Fatalf("expected a confirmation dialog, got %d: %s", preview.Code, preview.Body.String())
	}
	for _, part := range expectedParts {
		if !strings.Contains(preview.Body.String(), part) {
			t.Errorf("expected the dialog to contain %q", part)
		}
	}
	if got := testutil.QueryInt(t, db, `SELECT COUNT(*) FROM relation_events_players WHERE billettholder_id = 1`); got != 1 {
		t.Fatal("preview must keep the manual seat")
	}
}

func TestPuljeKonsekvenser_DescribesSeatChangesNeutrally(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en endring som flytter én deltaker og tar plassen fra en annen.",
		When:  "Når konsekvensene rendres.",
		Then:  "Så skal hver flytting vises fra og til, med førstevalg nevnt nøytralt.",
	})

	// Given
	konsekvenser := puljefordeling.Konsekvenser{PuljeNavn: "Fredag Kveld", Endringer: []puljefordeling.Konsekvens{
		{Name: "Yngve Yri",
			Fra: puljefordeling.Plass{EventID: "evB", EventTitle: "Bravo", Level: models.InterestLevelHigh, Forstevalg: true},
			Til: puljefordeling.Plass{EventID: "evC", EventTitle: "Charlie", Level: models.InterestLevelMedium}},
		{Name: "Zara Zahl", Fra: puljefordeling.Plass{EventID: "evC", EventTitle: "Charlie", Level: models.InterestLevelLow}},
	}}

	// When
	doc := templtest.Render(t, puljeKonsekvenser(konsekvenser))

	// Then
	rows := templtest.CollectTexts(doc, ".tildeling-konsekvens")
	if len(rows) != 2 {
		t.Fatalf("expected two changes, got %v", rows)
	}
	for _, part := range []string{"Yngve Yri", "🤩 Bravo", "😁 Charlie", "Får ikke førstevalg i denne puljen"} {
		if !strings.Contains(rows[0], part) {
			t.Errorf("expected %q in %q", part, rows[0])
		}
	}
	if !strings.Contains(rows[1], "Uten plass") {
		t.Errorf("expected Zara to end up without a seat, got %q", rows[1])
	}
	if empty := templtest.Render(t, puljeKonsekvenser(puljefordeling.Konsekvenser{PuljeNavn: "Fredag Kveld"})); !strings.Contains(empty.Text(), "Ingen andre deltakere får endret plass.") {
		t.Error("expected a note when nobody else changes seats")
	}
}

func TestPuljeEgenKonsekvens_ShowsTheMovedPlayersOwnChange(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at Xander flyttes fra Charlie til Bravo, der han får førstevalget sitt.",
		When:  "Når hans egen konsekvens rendres.",
		Then:  "Så skal flyttingen hans vises fra og til, med førstevalg nevnt.",
	})

	// Given
	expectedParts := []string{"For Xander Xu", "😁 Charlie", "🤩 Bravo", "Får førstevalg i denne puljen"}
	konsekvenser := puljefordeling.Konsekvenser{PuljeNavn: "Fredag Kveld", Egen: puljefordeling.Konsekvens{
		Fra: puljefordeling.Plass{EventID: "evC", EventTitle: "Charlie", Level: models.InterestLevelMedium},
		Til: puljefordeling.Plass{EventID: "evB", EventTitle: "Bravo", Level: models.InterestLevelHigh, Forstevalg: true},
	}}

	// When
	doc := templtest.Render(t, puljeEgenKonsekvens(konsekvenser, "Xander Xu"))

	// Then
	text := doc.Text()
	for _, part := range expectedParts {
		if !strings.Contains(text, part) {
			t.Errorf("expected %q in %q", part, text)
		}
	}
}
