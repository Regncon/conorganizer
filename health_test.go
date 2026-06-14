package main

import (
	"database/sql"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	_ "modernc.org/sqlite"

	"github.com/Regncon/conorganizer/testutil/bdd"
)

func TestHealthzReturnsGenericOK(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given the process is serving HTTP.",
		When:  "When the health endpoint is requested.",
		Then:  "Then it returns a generic OK response.",
	})

	// Given
	expectedStatusCode := http.StatusOK
	expectedBody := "ok\n"

	router := chi.NewRouter()
	mountHealthRoutes(router, newReadinessState(nil, testLogger()), testLogger())

	// When
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	router.ServeHTTP(recorder, request)

	// Then
	assertHTTPStatusAndBody(t, recorder, expectedStatusCode, expectedBody)
}

func TestReadyzReturnsOKWhenStartupAndLiveCheckPass(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given startup checks have passed and the database answers a live check.",
		When:  "When the readiness endpoint is requested.",
		Then:  "Then it returns a generic OK response.",
	})

	// Given
	expectedStatusCode := http.StatusOK
	expectedBody := "ok\n"

	db := openMemoryDB(t)
	defer db.Close()

	router := chi.NewRouter()
	mountHealthRoutes(router, newReadinessState(db, testLogger()), testLogger())

	// When
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	router.ServeHTTP(recorder, request)

	// Then
	assertHTTPStatusAndBody(t, recorder, expectedStatusCode, expectedBody)
}

func TestReadyzReturnsSanitizedFailureReasonWhenDegraded(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given the app is degraded because image storage is unavailable.",
		When:  "When the readiness endpoint is requested.",
		Then:  "Then it returns a sanitized reason without exposing internal paths.",
	})

	// Given
	expectedStatusCode := http.StatusServiceUnavailable
	expectedBody := "not ready: image directory not writable\n"
	internalPath := "/secret/path"

	state := newReadinessState(nil, testLogger())
	state.MarkDegraded(notReadyImageReason, fmt.Errorf("event image directory %s is not writable", internalPath))

	router := chi.NewRouter()
	mountHealthRoutes(router, state, testLogger())

	// When
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	router.ServeHTTP(recorder, request)

	// Then
	assertHTTPStatusAndBody(t, recorder, expectedStatusCode, expectedBody)

	body := recorder.Body.String()
	if strings.Contains(body, internalPath) {
		t.Fatalf("/readyz exposed internal path: %q", body)
	}
}

func TestReadyzReturnsDatabaseReasonWhenLiveCheckFails(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given no database is available for the live readiness check.",
		When:  "When the readiness endpoint is requested.",
		Then:  "Then it returns a sanitized database unavailable reason.",
	})

	// Given
	expectedStatusCode := http.StatusServiceUnavailable
	expectedBody := "not ready: database not available\n"

	router := chi.NewRouter()
	mountHealthRoutes(router, newReadinessState(nil, testLogger()), testLogger())

	// When
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	router.ServeHTTP(recorder, request)

	// Then
	assertHTTPStatusAndBody(t, recorder, expectedStatusCode, expectedBody)
}

func openMemoryDB(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open memory db: %v", err)
	}
	if err := db.Ping(); err != nil {
		_ = db.Close()
		t.Fatalf("ping memory db: %v", err)
	}
	return db
}

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func assertHTTPStatusAndBody(t *testing.T, recorder *httptest.ResponseRecorder, expectedStatusCode int, expectedBody string) {
	t.Helper()

	if recorder.Code != expectedStatusCode {
		t.Fatalf("HTTP status mismatch\nexpected: %d\nactual:   %d", expectedStatusCode, recorder.Code)
	}
	if recorder.Body.String() != expectedBody {
		t.Fatalf("HTTP body mismatch\nexpected: %q\nactual:   %q", expectedBody, recorder.Body.String())
	}
}
