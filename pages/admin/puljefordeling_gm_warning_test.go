package admin

import (
	"net/http"
	"strings"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/service/live"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/go-chi/chi/v5"
)

func TestPuljefordelingAssignRoute_DragGMReceivesWarningEvenAfterAgeConfirmation(t *testing.T) {
	for _, confirmed := range []bool{false, true} {
		name := "ordinary"
		if confirmed {
			name = "age_confirmed"
		}
		t.Run(name, func(t *testing.T) {
			// Given
			const expectedPins = 0
			db, logger := testutil.CreateTestDBAndLogger(t, "gm_warning_"+name)
			router := chi.NewRouter()
			puljefordelingRoute(router, db, &live.Manager{}, logger, nil)
			age := models.AgeGroupDefault
			if confirmed {
				age = models.AgeGroupAdultsOnly
			}
			seedAssignFixture(t, db, models.PuljeFredagKveld, age, false)
			testutil.MustExec(t, db, `INSERT INTO relation_events_players (event_id, pulje_id, billettholder_id, role) VALUES ('evA', 'FredagKveld', 1, 'GM')`)
			// When
			post := postAssignSignals
			if confirmed {
				post = postAssignSignalsConfirmed
			}
			rec := post(t, router, 1, "evA", "FredagKveld")
			// Then
			body := rec.Body.String()
			if rec.Code != http.StatusOK || !strings.Contains(body, "datastar-patch-signals") || !strings.Contains(body, "tildelingOpen") || !strings.Contains(body, "Spilleder") || strings.Contains(body, "data-confirmation=") {
				t.Errorf("expected GM warning signal patch, got %d: %s", rec.Code, body)
			}
			var pins int
			if err := db.QueryRow(`SELECT COUNT(*) FROM relation_events_players WHERE role = 'Player'`).Scan(&pins); err != nil {
				t.Fatal(err)
			}
			if pins != expectedPins {
				t.Errorf("want %d pins, got %d", expectedPins, pins)
			}
		})
	}
}
