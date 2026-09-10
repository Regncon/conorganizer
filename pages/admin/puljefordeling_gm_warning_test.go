package admin

import (
	"net/http"
	"strings"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/service/live"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/templtest"
	"github.com/go-chi/chi/v5"
)

func TestPuljefordelingAssignRoute_GMReceivesWarningWithoutPin(t *testing.T) {
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
			if rec.Code != http.StatusOK || !strings.Contains(body, "datastar-patch-signals") || !strings.Contains(body, "gmWarningText") || !strings.Contains(body, "GM") {
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

func TestPuljefordelingIndex_GMWarningDialogCanOnlyBeDismissed(t *testing.T) {
	// Given
	const selector = "#puljefordeling-gm-dialog"
	db, logger := testutil.CreateTestDBAndLogger(t, "gm_warning_dialog")
	seedTabPulje(t, db, models.PuljeFredagKveld, "Fredag Kveld", models.PuljeStatusOpen, "2026-01-01 18:00")
	// When
	doc := templtest.Render(t, puljefordelingIndex(db, logger, models.PuljeFredagKveld, nil))
	// Then
	dialog := doc.Find(selector)
	if dialog.Length() != 1 {
		t.Fatalf("expected one GM warning dialog, got %d", dialog.Length())
	}
	if doc.Find("#puljefordeling-tab "+selector).Length() != 0 {
		t.Error("warning must survive live region updates")
	}
	if !strings.Contains(dialog.AttrOr("data-effect", ""), "$gmWarningText") || !strings.Contains(dialog.AttrOr("data-effect", ""), "showModal()") {
		t.Error("warning signal must open popup")
	}
	if dialog.Find("p").AttrOr("data-text", "") != "$gmWarningText" {
		t.Error("popup must show server warning")
	}
	if !strings.Contains(dialog.AttrOr("data-on:close", ""), "$gmWarningText = ''") {
		t.Error("closing must clear warning so it can open again")
	}
	if dialog.Find("button").Length() != 1 {
		t.Fatal("warning must offer only dismissal")
	}
	action := dialog.Find("button").AttrOr("data-on:click", "")
	if !strings.Contains(action, ".close()") || strings.Contains(action, "@post") {
		t.Errorf("expected dismissal without override: %s", action)
	}
}
