package formsubmission

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/service/authctx"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/go-chi/chi/v5"
)

func TestUpdateProgramTime_StoresTheTimeForOnlyThatPulje(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt et programarrangement i to lørdagspuljer.",
		When:  "Når admin skriver « Hele dagen » som tidspunkt for lørdag morgen.",
		Then:  "Så lagres «Hele dagen» uten mellomrom rundt for lørdag morgen, og lørdag kveld er uendret.",
	})

	// Given
	expectedMorning := "Hele dagen"
	expectedEvening := ""
	db, logger := testutil.CreateTestDBAndLogger(t, "update_program_time")
	insertProgramTimeFixture(t, db)
	ctx := authctx.WithUserToken(context.Background(), "ext-42", "admin@x.no")

	// When
	err := updateProgramTimeInRelationEventPuljeWithAudit(ctx, db, logger, "e1", string(models.PuljeLordagMorgen), "  Hele dagen ")

	// Then
	if err != nil {
		t.Fatalf("expected the program time to be saved: %v", err)
	}
	assertProgramTime(t, db, models.PuljeLordagMorgen, expectedMorning)
	assertProgramTime(t, db, models.PuljeLordagKveld, expectedEvening)
}

func TestUpdateProgramTimeRoute_MissingSignalKeepsTheSavedTime(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt et programarrangement med et lagret tidspunkt for lørdag morgen.",
		When:  "Når en lagring kommer uten tidspunkt for den puljen.",
		Then:  "Så avvises den, og det lagrede tidspunktet blir ikke tømt.",
	})

	// Given
	expectedStatus := http.StatusBadRequest
	expectedMorning := "Hele dagen"
	db, logger := testutil.CreateTestDBAndLogger(t, "update_program_time_missing_signal")
	insertProgramTimeFixture(t, db)
	testutil.MustExec(t, db, `UPDATE relation_event_puljer SET program_time = ? WHERE pulje_id = ?`, expectedMorning, models.PuljeLordagMorgen)
	router := chi.NewRouter()
	router.Route("/profile/api/new/{id}/program-time", func(programTimeRouter chi.Router) {
		UpdateProgramTimeInPulje(programTimeRouter, db, nil, logger)
	})
	request := httptest.NewRequest(http.MethodPut, "/profile/api/new/e1/program-time/LordagMorgen", strings.NewReader(`{"programTime":{"LordagKveld":""}}`))
	recorder := httptest.NewRecorder()

	// When
	router.ServeHTTP(recorder, request)

	// Then
	if recorder.Code != expectedStatus {
		t.Fatalf("HTTP status = %d, want %d", recorder.Code, expectedStatus)
	}
	assertProgramTime(t, db, models.PuljeLordagMorgen, expectedMorning)
}

// insertProgramTimeFixture seeds an admin and a program event in both Saturday puljer.
func insertProgramTimeFixture(t *testing.T, db *sql.DB) {
	t.Helper()
	testutil.MustExec(t, db, `INSERT INTO users (id, external_id, email, is_admin) VALUES (42, 'ext-42', 'admin@x.no', 1)`)
	testutil.MustExec(t, db, `
		INSERT INTO puljer (id, name, status, start_at, end_at) VALUES
			(?, 'Lørdag morgen', 'Open', '2026-10-03T10:00:00+02:00', '2026-10-03T15:00:00+02:00'),
			(?, 'Lørdag kveld', 'Open', '2026-10-03T18:00:00+02:00', '2026-10-03T23:00:00+02:00')
	`, models.PuljeLordagMorgen, models.PuljeLordagKveld)
	testutil.MustExec(t, db, `
		INSERT INTO events (id, title, intro, description, host_name, email, phone_number, max_players, is_in_puljefordeling)
		VALUES ('e1', 'Infodisk', '', '', 'Styret', 'styret@x.no', '', 0, 0)
	`)
	testutil.MustExec(t, db, `
		INSERT INTO relation_event_puljer (event_id, pulje_id, is_in_pulje) VALUES ('e1', ?, 1), ('e1', ?, 1)
	`, models.PuljeLordagMorgen, models.PuljeLordagKveld)
}

func assertProgramTime(t *testing.T, db *sql.DB, puljeID models.Pulje, expected string) {
	t.Helper()
	var actual string
	if err := db.QueryRow(`SELECT program_time FROM relation_event_puljer WHERE event_id = 'e1' AND pulje_id = ?`, puljeID).Scan(&actual); err != nil {
		t.Fatalf("read program time for %s: %v", puljeID, err)
	}
	if actual != expected {
		t.Fatalf("program time for %s = %q, want %q", puljeID, actual, expected)
	}
}
