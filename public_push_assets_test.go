package main

import (
	"net/http"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestPublicAssetRoutes_ServesRootScopedNotificationWorker(t *testing.T) {
	// Given
	expectedStatus := http.StatusOK
	router := chi.NewRouter()
	mountPublicAssetRoutes(router, nil, testLogger())

	// When
	rec := performTestRequest(router, "/varsler-sw.js")

	// Then
	if rec.Code != expectedStatus {
		t.Fatalf("worker returned %d, want %d", rec.Code, expectedStatus)
	}
	if !strings.Contains(rec.Header().Get("Content-Type"), "javascript") {
		t.Errorf("worker content type = %q", rec.Header().Get("Content-Type"))
	}
	if rec.Header().Get("Service-Worker-Allowed") != "/" {
		t.Errorf("worker must be allowed to control the profile route")
	}
	if cacheControl := rec.Header().Get("Cache-Control"); cacheControl != "no-cache" && cacheControl != "no-store" {
		t.Errorf("worker updates must be revalidated: %q", rec.Header().Get("Cache-Control"))
	}
}
