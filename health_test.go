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
)

func TestHealthzReturnsGenericOK(t *testing.T) {
	router := chi.NewRouter()
	mountHealthRoutes(router, newReadinessState(nil, testLogger()), testLogger())

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("/healthz status = %d, want %d", recorder.Code, http.StatusOK)
	}
	if recorder.Body.String() != "ok\n" {
		t.Fatalf("/healthz body = %q, want generic ok", recorder.Body.String())
	}
}

func TestReadyzReturnsOKWhenStartupAndLiveCheckPass(t *testing.T) {
	db := openMemoryDB(t)
	defer db.Close()

	router := chi.NewRouter()
	mountHealthRoutes(router, newReadinessState(db, testLogger()), testLogger())

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("/readyz status = %d, want %d", recorder.Code, http.StatusOK)
	}
	if recorder.Body.String() != "ok\n" {
		t.Fatalf("/readyz body = %q, want generic ok", recorder.Body.String())
	}
}

func TestReadyzReturnsSanitizedFailureReasonWhenDegraded(t *testing.T) {
	state := newReadinessState(nil, testLogger())
	state.MarkDegraded(notReadyImageReason, fmt.Errorf("event image directory /secret/path is not writable"))

	router := chi.NewRouter()
	mountHealthRoutes(router, state, testLogger())

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusServiceUnavailable {
		t.Fatalf("/readyz status = %d, want %d", recorder.Code, http.StatusServiceUnavailable)
	}
	body := recorder.Body.String()
	if body != "not ready: image directory not writable\n" {
		t.Fatalf("/readyz body = %q, want sanitized not ready reason", body)
	}
	if strings.Contains(body, "/secret/path") {
		t.Fatalf("/readyz exposed internal path: %q", body)
	}
}

func TestReadyzReturnsDatabaseReasonWhenLiveCheckFails(t *testing.T) {
	router := chi.NewRouter()
	mountHealthRoutes(router, newReadinessState(nil, testLogger()), testLogger())

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusServiceUnavailable {
		t.Fatalf("/readyz status = %d, want %d", recorder.Code, http.StatusServiceUnavailable)
	}
	if recorder.Body.String() != "not ready: database not available\n" {
		t.Fatalf("/readyz body = %q, want database not available reason", recorder.Body.String())
	}
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
