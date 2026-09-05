package formsubmission

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/service/authctx"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/Regncon/conorganizer/testutil/templtest"
	"github.com/go-chi/chi/v5"
)

func TestFormBody_OpenRegistrationSettingIsAdminOnly(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at skjemaet kan vises både til en administrator og en vanlig arrangør.",
		When:  "Når arrangementsskjemaet rendres.",
		Then:  "Så skal bare administratoren se innstillingen for åpen påmelding.",
	})

	// Given
	expectedAdminControls := 1
	expectedRegularControls := 0
	db := testutil.CreateTestDB(t, "open_registration_visibility")
	event := &models.Event{ID: "event-1"}

	// When
	adminDocument := templtest.Render(t, FormBody(event, nil, true, db))
	regularDocument := templtest.Render(t, FormBody(event, nil, false, db))
	actualAdminControls := adminDocument.Find(`[data-testid="open-registration"]`).Length()
	actualRegularControls := regularDocument.Find(`[data-testid="open-registration"]`).Length()

	// Then
	if actualAdminControls != expectedAdminControls {
		t.Fatalf("admin control count mismatch\nexpected: %d\nactual:   %d", expectedAdminControls, actualAdminControls)
	}
	if actualRegularControls != expectedRegularControls {
		t.Fatalf("regular control count mismatch\nexpected: %d\nactual:   %d", expectedRegularControls, actualRegularControls)
	}
}

func TestSetEventOpenRegistration_WhenEnabled_PersistsConfiguration(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt et arrangement der åpen påmelding er slått av.",
		When:  "Når en administrator slår på åpen påmelding.",
		Then:  "Så skal innstillingen og hvem som endret den lagres.",
	})

	// Given
	expectedOpenRegistration := 1
	expectedUpdatedByID := 42
	db, logger := testutil.CreateTestDBAndLogger(t, "set_open_registration")
	testutil.MustExec(t, db, `INSERT INTO users (id, external_id, email, is_admin) VALUES (42, 'ext-42', 'admin@x.no', 1)`)
	testutil.MustExec(t, db, `
		INSERT INTO events (id, title, intro, description, host_name, email, phone_number, max_players)
		VALUES ('event-1', 'Spel', '', '', 'Ola', 'ola@x.no', '', 4)`)
	ctx := authctx.WithUserToken(context.Background(), "ext-42", "admin@x.no")

	// When
	if err := setEventOpenRegistration(ctx, db, logger, "event-1", true); err != nil {
		t.Fatalf("expected open registration update to succeed: %v", err)
	}
	actualOpenRegistration := testutil.QueryInt(t, db, `SELECT is_open_registration FROM events WHERE id = 'event-1'`)
	actualUpdatedByID := testutil.QueryInt(t, db, `SELECT updated_by_id FROM events WHERE id = 'event-1'`)

	// Then
	if actualOpenRegistration != expectedOpenRegistration {
		t.Fatalf("open-registration mismatch\nexpected: %d\nactual:   %d", expectedOpenRegistration, actualOpenRegistration)
	}
	if actualUpdatedByID != expectedUpdatedByID {
		t.Fatalf("updated-by mismatch\nexpected: %d\nactual:   %d", expectedUpdatedByID, actualUpdatedByID)
	}
}

func TestUpdateOpenRegistration_WhenUserIsNotAdmin_ReturnsForbidden(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at en vanlig arrangør forsøker å endre åpen påmelding direkte.",
		When:  "Når endepunktet mottar forespørselen.",
		Then:  "Så skal forespørselen avvises uten å endre arrangementet.",
	})

	// Given
	expectedStatusCode := http.StatusForbidden
	expectedOpenRegistration := 0
	db, logger := testutil.CreateTestDBAndLogger(t, "open_registration_forbidden")
	testutil.MustExec(t, db, `INSERT INTO users (id, external_id, email, is_admin) VALUES (42, 'ext-42', 'user@x.no', 0)`)
	testutil.MustExec(t, db, `
		INSERT INTO events (id, title, intro, description, host_name, email, phone_number, max_players)
		VALUES ('event-1', 'Spel', '', '', 'Ola', 'ola@x.no', '', 4)`)
	router := chi.NewRouter()
	router.Route("/profile/api/new/{id}/open-registration", func(openRegistrationRouter chi.Router) {
		UpdateOpenRegistration(openRegistrationRouter, db, nil, logger)
	})
	request := httptest.NewRequest(
		http.MethodPut,
		"/profile/api/new/event-1/open-registration",
		strings.NewReader(`{"isOpenRegistration":true}`),
	)
	request = request.WithContext(authctx.WithUserToken(request.Context(), "ext-42", "user@x.no"))
	recorder := httptest.NewRecorder()

	// When
	router.ServeHTTP(recorder, request)
	actualOpenRegistration := testutil.QueryInt(t, db, `SELECT is_open_registration FROM events WHERE id = 'event-1'`)

	// Then
	if recorder.Code != expectedStatusCode {
		t.Fatalf("HTTP status mismatch\nexpected: %d\nactual:   %d", expectedStatusCode, recorder.Code)
	}
	if actualOpenRegistration != expectedOpenRegistration {
		t.Fatalf("open-registration mismatch\nexpected: %d\nactual:   %d", expectedOpenRegistration, actualOpenRegistration)
	}
}
