package formsubmission

import (
	"bytes"
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func TestUpdateDescription_WhenEventUpdateFails_LogsContext(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given an event description update without a current user.",
		When:  "When the handler cannot write the audited event update.",
		Then:  "Then it returns an internal server error and logs request context.",
	})

	// Given
	expectedStatus := http.StatusInternalServerError
	expectedLogFragments := []string{
		`"component":"event_form"`,
		`get current user id for event audit`,
		`"event_id":"event-123"`,
		`"request_id":"request-123"`,
	}
	db := testutil.CreateTestDB(t, "update_description_logs")
	var logs bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logs, nil))
	router := chi.NewRouter()
	router.Route("/profile/api/new/{id}/description", func(descriptionRouter chi.Router) {
		UpdateDescription(descriptionRouter, db, nil, logger)
	})

	request := httptest.NewRequest(http.MethodPut, "/profile/api/new/event-123/description", strings.NewReader(`{"description":"New description"}`))
	request = request.WithContext(context.WithValue(request.Context(), middleware.RequestIDKey, "request-123"))
	recorder := httptest.NewRecorder()

	// When
	router.ServeHTTP(recorder, request)

	// Then
	if recorder.Code != expectedStatus {
		t.Fatalf("HTTP status mismatch\nexpected: %d\nactual:   %d", expectedStatus, recorder.Code)
	}

	logText := logs.String()
	for _, expected := range expectedLogFragments {
		if !strings.Contains(logText, expected) {
			t.Fatalf("expected log to contain %q\nactual log: %s", expected, logText)
		}
	}
}
