package varsler

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/go-chi/chi/v5"
)

func TestPushRoutes_RejectUnauthenticatedRequests(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given notification routes with no authenticated user in the request context.",
		When:  "When configuration or subscription endpoints are requested.",
		Then:  "Then every route rejects the request as unauthorized.",
	})

	// Given
	expectedStatus := http.StatusUnauthorized
	db, logger := testutil.CreateTestDBAndLogger(t, "varsler_unauthenticated_routes")
	service, err := New(db, logger, testConfig(t))
	if err != nil {
		t.Fatalf("create varsler service: %v", err)
	}
	router := chi.NewRouter()
	service.RegisterRoutes(router)
	p256dh, auth := validSubscriptionKeys()
	endpoint := "https://fcm.googleapis.com/fcm/send/unauthenticated-device"
	subscriptionBody := fmt.Sprintf(`{"endpoint":%q,"keys":{"p256dh":%q,"auth":%q}}`, endpoint, p256dh, auth)
	requests := []struct {
		name   string
		method string
		target string
		body   string
	}{
		{name: "configuration GET", method: http.MethodGet, target: "/api/varsler/config"},
		{name: "subscription GET", method: http.MethodGet, target: "/api/varsler/subscriptions?endpoint=" + endpoint},
		{name: "subscription POST", method: http.MethodPost, target: "/api/varsler/subscriptions", body: subscriptionBody},
		{name: "subscription DELETE", method: http.MethodDelete, target: "/api/varsler/subscriptions", body: fmt.Sprintf(`{"endpoint":%q}`, endpoint)},
	}

	for _, request := range requests {
		t.Run(request.name, func(t *testing.T) {
			// When
			response := serveUnauthenticatedPushRequest(router, request.method, request.target, request.body)

			// Then
			if response.Code != expectedStatus {
				t.Fatalf("unauthenticated %s %s returned %d, want %d: %s", request.method, request.target, response.Code, expectedStatus, response.Body.String())
			}
		})
	}
}

func TestSubscriptionRoutes_RejectInvalidSubscriptionKeys(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given an authenticated user registering an allowed push endpoint with malformed subscription keys.",
		When:  "When the subscription is submitted.",
		Then:  "Then the route rejects it without storing a subscription.",
	})

	// Given
	expectedStatus := http.StatusBadRequest
	db, logger := testutil.CreateTestDBAndLogger(t, "varsler_invalid_subscription_keys")
	testutil.MustExec(t, db, `INSERT INTO users(id,external_id,email) VALUES (1,'external-1','one@example.com')`)
	service, err := New(db, logger, testConfig(t))
	if err != nil {
		t.Fatalf("create varsler service: %v", err)
	}
	router := chi.NewRouter()
	service.RegisterRoutes(router)
	validP256dh, validAuth := validSubscriptionKeys()
	endpoint := "https://fcm.googleapis.com/fcm/send/invalid-key-device"
	tests := []struct {
		name   string
		p256dh string
		auth   string
	}{
		{name: "malformed public key", p256dh: "not-a-p256-key", auth: validAuth},
		{name: "short auth secret", p256dh: validP256dh, auth: "c2hvcnQ"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// When
			body := fmt.Sprintf(`{"endpoint":%q,"keys":{"p256dh":%q,"auth":%q}}`, endpoint, test.p256dh, test.auth)
			response := serveAuthenticated(router, http.MethodPost, "/api/varsler/subscriptions", body, "external-1")

			// Then
			if response.Code != expectedStatus {
				t.Fatalf("invalid subscription keys returned %d, want %d: %s", response.Code, expectedStatus, response.Body.String())
			}
			if got := testutil.QueryInt(t, db, `SELECT COUNT(*) FROM web_push_subscriptions`); got != 0 {
				t.Fatalf("invalid subscription stored %d rows, want 0", got)
			}
		})
	}
}

func serveUnauthenticatedPushRequest(handler http.Handler, method string, target string, body string) *httptest.ResponseRecorder {
	request := httptest.NewRequest(method, target, bytes.NewBufferString(body))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	return response
}
