package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Regncon/conorganizer/service/authctx"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/descope/go-sdk/descope"
	"github.com/go-chi/chi/v5"
)

func TestRoutes_InvalidPushConfigurationKeepsProfileNotificationsDisabled(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given expected protected notification routes with an invalid Web Push configuration.",
		When:  "When unauthenticated and authenticated clients request the notification configuration.",
		Then:  "Then authentication is required and the authenticated response reports notifications as disabled.",
	})

	// Given
	expectedStatus := http.StatusOK
	expectedUnauthenticatedStatus := http.StatusUnauthorized
	t.Setenv("WEB_PUSH_PUBLIC_KEY", "invalid")
	t.Setenv("WEB_PUSH_PRIVATE_KEY", "")
	t.Setenv("WEB_PUSH_SUBJECT", "")
	db, logger := testutil.CreateTestDBAndLogger(t, "notification_routes")
	router := chi.NewRouter()
	validator := notificationRouteSessionValidator{
		token: &descope.Token{
			ID:     "fixture-user",
			Claims: map[string]any{"email": "fixture@example.com"},
		},
	}
	authenticatedRouter := router.With(authctx.AuthMiddleware(validator, logger))
	imageDir := t.TempDir()
	cleanup, err := setupRoutes(context.Background(), logger, authenticatedRouter, router, db, &imageDir, t.TempDir(), validator)
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
	unauthenticatedRequest := httptest.NewRequest(http.MethodGet, "/api/varsler/config", nil)
	unauthenticatedResponse := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/varsler/config", nil)
	request.AddCookie(&http.Cookie{Name: authctx.SessionCookieName, Value: "fixture-session"})
	response := httptest.NewRecorder()

	// When
	router.ServeHTTP(unauthenticatedResponse, unauthenticatedRequest)
	router.ServeHTTP(response, request)

	// Then
	if unauthenticatedResponse.Code != expectedUnauthenticatedStatus {
		t.Fatalf("unauthenticated config returned %d: %s", unauthenticatedResponse.Code, unauthenticatedResponse.Body.String())
	}
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

type notificationRouteSessionValidator struct {
	token *descope.Token
}

func (validator notificationRouteSessionValidator) ValidateSessionWithToken(context.Context, string) (bool, *descope.Token, error) {
	return true, validator.token, nil
}

func (validator notificationRouteSessionValidator) RefreshSessionWithToken(context.Context, string) (bool, *descope.Token, error) {
	return false, nil, nil
}
