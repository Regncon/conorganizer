package admin

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/service/live"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/Regncon/conorganizer/testutil/templtest"
	"github.com/go-chi/chi/v5"
)

func seedGMRemoval(t *testing.T, db *sql.DB, status models.PuljeStatus) {
	t.Helper()
	seedTabPulje(t, db, models.PuljeFredagKveld, "Fredag", status, "2026-01-01 18:00")
	seedTabPulje(t, db, models.PuljeLordagMorgen, "Lørdag", models.PuljeStatusOpen, "2026-01-02 10:00")
	seedTabEventWithInterest(t, db, "evA", "Alpha", models.PuljeFredagKveld)
	testutil.MustExec(t, db, `INSERT INTO billettholdere (id, first_name, last_name, ticket_type_id, ticket_type, order_id, ticket_id) VALUES (2,'Anne','Nordmann',0,'',0,2)`)
	testutil.MustExec(t, db, `INSERT INTO relation_events_players (event_id, pulje_id, billettholder_id, role, source) VALUES ('evA','FredagKveld',1,'GM','manual'), ('evA','FredagKveld',1,'Player','manual'), ('evA','FredagKveld',2,'GM','manual'), ('evA','LordagMorgen',1,'GM','manual'), ('evA','LordagMorgen',2,'Player','manual')`)
}

func TestPuljefordelingTabContent_IndividualGMRemoval(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{Given: "an event has two GMs", When: "the editable pulje is displayed", Then: "each GM has an accessible removal control targeting their own assignment"})
	// Given
	expectedNames := []string{"Anne Nordmann", "Kari Nordmann"}
	expectedIDs := []int{2, 1}
	db, logger := testutil.CreateTestDBAndLogger(t, "individual_gm_removal")
	seedGMRemoval(t, db, models.PuljeStatusOpen)
	// When
	doc := templtest.Render(t, PuljefordelingTabContent(db, logger, models.PuljeFredagKveld, nil))
	// Then
	buttons := doc.Find(".pulje-gm .pulje-remove")
	if buttons.Length() != 2 {
		t.Fatalf("expected two GM removal controls, got %d", buttons.Length())
	}
	for i, name := range expectedNames {
		if got := buttons.Eq(i).AttrOr("aria-label", ""); got != "Fjern spilleder "+name {
			t.Errorf("accessible name: %q", got)
		}
		expected := fmt.Sprintf("@post('/admin/api/puljefordeling/FredagKveld/evA/%d/gm/preview')", expectedIDs[i])
		if got := buttons.Eq(i).AttrOr("data-on:click", ""); got != expected {
			t.Errorf("action: got %q, want %q", got, expected)
		}
	}
}

func TestPuljefordelingRemoveGMRoute_PreservesOtherAssignments(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{Given: "two GMs and assignments in another pulje", When: "one GM is removed", Then: "only that GM assignment is deleted"})
	// Given
	expectedRemaining := 4
	db, logger := testutil.CreateTestDBAndLogger(t, "remove_one_gm")
	seedGMRemoval(t, db, models.PuljeStatusOpen)
	router := chi.NewRouter()
	puljefordelingRoute(router, db, &live.Manager{}, logger, nil)
	// When
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodDelete, "/api/puljefordeling/FredagKveld/evA/1/gm", nil))
	// Then
	if rec.Code != http.StatusNoContent {
		t.Fatalf("status %d: %s", rec.Code, rec.Body.String())
	}
	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM relation_events_players`).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != expectedRemaining {
		t.Fatalf("remaining assignments %d, want %d", count, expectedRemaining)
	}
	if err := db.QueryRow(`SELECT COUNT(*) FROM relation_events_players WHERE pulje_id='FredagKveld' AND billettholder_id=1 AND role='GM'`).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatal("selected GM still assigned")
	}
	if err := db.QueryRow(`SELECT COUNT(*) FROM relation_events_players WHERE pulje_id='FredagKveld' AND event_id='evA' AND billettholder_id=1 AND role='Player'`).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatal("removing GM also removed their Player assignment")
	}
}

func TestPuljefordelingRemoveGMRoute_PublishedIsReadOnly(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{Given: "a published pulje with two GMs", When: "removal is attempted", Then: "controls are hidden and the server preserves all assignments"})
	// Given
	expectedStatus := http.StatusConflict
	db, logger := testutil.CreateTestDBAndLogger(t, "remove_published_gm")
	seedGMRemoval(t, db, models.PuljeStatusCompleted)
	router := chi.NewRouter()
	puljefordelingRoute(router, db, &live.Manager{}, logger, nil)
	// When
	doc := templtest.Render(t, PuljefordelingTabContent(db, logger, models.PuljeFredagKveld, nil))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodDelete, "/api/puljefordeling/FredagKveld/evA/1/gm", nil))
	// Then
	if doc.Find(".pulje-gm .pulje-remove").Length() != 0 {
		t.Fatal("published GM controls visible")
	}
	if rec.Code != expectedStatus {
		t.Fatalf("status %d, want %d", rec.Code, expectedStatus)
	}
	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM relation_events_players`).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 5 {
		t.Fatalf("published assignments changed: %d", count)
	}
}
