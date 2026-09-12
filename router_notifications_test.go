package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Regncon/conorganizer/service/authctx"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/go-chi/chi/v5"
)

func TestRoutes_InvalidPushConfigurationKeepsProfileNotificationsDisabled(t *testing.T) {
	// Given
	expectedStatus := http.StatusOK
	t.Setenv("WEB_PUSH_PUBLIC_KEY", "invalid")
	t.Setenv("WEB_PUSH_PRIVATE_KEY", "")
	t.Setenv("WEB_PUSH_SUBJECT", "")
	db, logger := testutil.CreateTestDBAndLogger(t, "notification_routes")
	router := chi.NewRouter()
	imageDir := t.TempDir()
	cleanup, err := setupRoutes(context.Background(), logger, router, db, &imageDir, t.TempDir())
	if cleanup != nil {
		t.Cleanup(func() {
			if err := cleanup(); err != nil {
				t.Errorf("cleanup routes: %v", err)
			}
		})
	}
	if err != nil {
		t.Fatalf("setup routes: %v", err)
	}
	request := httptest.NewRequest(http.MethodGet, "/api/varsler/config", nil)
	request = request.WithContext(authctx.WithUserToken(request.Context(), "fixture-user", "fixture@example.com"))
	response := httptest.NewRecorder()

	// When
	router.ServeHTTP(response, request)

	// Then
	if response.Code != expectedStatus {
		t.Fatalf("config returned %d: %s", response.Code, response.Body.String())
	}
	var config struct {
		Enabled   bool
		PublicKey string
	}
	if err := json.Unmarshal(response.Body.Bytes(), &config); err != nil {
		t.Fatal(err)
	}
	if config.Enabled || config.PublicKey != "" {
		t.Fatalf("invalid configuration must disable notifications: %+v", config)
	}
}
