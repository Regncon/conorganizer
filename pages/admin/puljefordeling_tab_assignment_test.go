package admin

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/service/live"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/templtest"
	"github.com/go-chi/chi/v5"
)

func TestPuljefordelingRemoveManualSeatRoute_DeletesPin(t *testing.T) {
	db, logger := testutil.CreateTestDBAndLogger(t, "puljefordeling_remove_route")
	router := chi.NewRouter()
	puljefordelingRoute(router, db, &live.Manager{}, logger)

	const fredag = models.PuljeFredagKveld
	seedTabPulje(t, db, fredag, "Fredag Kveld", models.PuljeStatusOpen, "2026-01-01 18:00")
	testutil.MustExec(t, db, `INSERT INTO events (id, title, intro, description, host_name, email, phone_number, max_players)
		VALUES ('evA','Alpha','','','','','',4)`)
	testutil.MustExec(t, db, `INSERT INTO relation_event_puljer (event_id, pulje_id, is_in_pulje) VALUES ('evA',?,1)`, string(fredag))
	testutil.MustExec(t, db, `INSERT INTO billettholdere (id, first_name, last_name, ticket_type_id, ticket_type, order_id, ticket_id)
		VALUES (1,'Kari','Nordmann',0,'',0,1)`)
	testutil.MustExec(t, db, `INSERT INTO relation_events_players (event_id, pulje_id, billettholder_id, role, source)
		VALUES ('evA',?,1,'Player','manual')`, string(fredag))

	req := httptest.NewRequest(http.MethodDelete, "/api/puljefordeling/FredagKveld/evA/1", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("want 204 No Content, got %d (%s)", rec.Code, rec.Body.String())
	}

	var n int
	if err := db.QueryRow(
		`SELECT COUNT(*) FROM relation_events_players WHERE event_id='evA' AND billettholder_id=1 AND source='manual'`,
	).Scan(&n); err != nil {
		t.Fatalf("count manual seats: %v", err)
	}
	if n != 0 {
		t.Fatalf("manual seat should be deleted, still found %d", n)
	}
}

func TestPuljefordelingTabContent_RendersAddPickerAndManualRemove(t *testing.T) {
	db, logger := testutil.CreateTestDBAndLogger(t, "puljefordeling_tab_interactive")

	const fredag = models.PuljeFredagKveld
	seedTabPulje(t, db, fredag, "Fredag Kveld", models.PuljeStatusOpen, "2026-01-01 18:00")
	testutil.MustExec(t, db, `INSERT INTO events (id, title, intro, description, host_name, email, phone_number, max_players)
		VALUES ('evA','Alpha','','','','','',4)`)
	testutil.MustExec(t, db, `INSERT INTO relation_event_puljer (event_id, pulje_id, is_in_pulje) VALUES ('evA',?,1)`, string(fredag))
	testutil.MustExec(t, db, `INSERT INTO billettholdere (id, first_name, last_name, ticket_type_id, ticket_type, order_id, ticket_id)
		VALUES (1,'Kari','Nordmann',0,'',0,1)`)
	// Kari is a manual pin → her tile must offer a × that removes her seat.
	testutil.MustExec(t, db, `INSERT INTO relation_events_players (event_id, pulje_id, billettholder_id, role, source)
		VALUES ('evA',?,1,'Player','manual')`, string(fredag))

	doc := templtest.Render(t, PuljefordelingTabContent(db, logger, fredag))
	html, err := doc.Html()
	if err != nil {
		t.Fatalf("render html: %v", err)
	}

	wants := []string{
		// The shared picker dialog and custom search element are present.
		"puljefordeling-assign-dialog",
		"admin-billettholder-search",
		// The × on Kari's manual tile deletes her manual seat.
		"/admin/api/puljefordeling/FredagKveld/evA/1",
	}
	for _, want := range wants {
		if !strings.Contains(html, want) {
			t.Errorf("rendered tab should contain %q", want)
		}
	}

	// The + button opens the dialog scoped to this event (attribute decoded).
	addClick := doc.Find(".pulje-add").AttrOr("data-on:click", "")
	if !strings.Contains(addClick, "$assignmentEventId = 'evA'") {
		t.Errorf("+ button should set assignmentEventId to the event; got %q", addClick)
	}
}

func TestPuljefordelingRemoveManualSeatRoute_RejectsInvalidPulje(t *testing.T) {
	db, logger := testutil.CreateTestDBAndLogger(t, "puljefordeling_remove_route_invalid")
	router := chi.NewRouter()
	puljefordelingRoute(router, db, &live.Manager{}, logger)

	req := httptest.NewRequest(http.MethodDelete, "/api/puljefordeling/NotAPulje/evA/1", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("invalid pulje should be rejected with 400, got %d", rec.Code)
	}
}
