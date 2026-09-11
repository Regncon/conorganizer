package varsler

import (
	"bytes"
	"crypto/ecdh"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Regncon/conorganizer/service/authctx"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/go-chi/chi/v5"
)

func TestSubscriptionRoutes_CurrentUserCanRegisterInspectAndDelete(t *testing.T) {
	// Given
	db, logger := testutil.CreateTestDBAndLogger(t, "varsler_subscription_lifecycle")
	testutil.MustExec(t, db, `INSERT INTO users(id,external_id,email) VALUES (1,'external-1','one@example.com')`)
	service, err := New(db, logger, testConfig(t))
	if err != nil {
		t.Fatalf("create varsler service: %v", err)
	}
	router := chi.NewRouter()
	service.RegisterRoutes(router)
	endpoint := "https://fcm.googleapis.com/fcm/send/device-one"
	p256dh, auth := validSubscriptionKeys()
	body := fmt.Sprintf(`{"endpoint":%q,"expirationTime":null,"keys":{"p256dh":%q,"auth":%q}}`, endpoint, p256dh, auth)

	// When / Then: register.
	response := serveAuthenticated(router, http.MethodPost, "/api/varsler/subscriptions", body, "external-1")
	if response.Code != http.StatusNoContent {
		t.Fatalf("expected registration status 204, got %d: %s", response.Code, response.Body.String())
	}

	// When / Then: inspect only the current owner's endpoint.
	response = serveAuthenticated(router, http.MethodGet, "/api/varsler/subscriptions?endpoint="+endpoint, "", "external-1")
	if response.Code != http.StatusOK || response.Body.String() != "{\"subscribed\":true}\n" {
		t.Fatalf("expected subscribed response, got %d: %s", response.Code, response.Body.String())
	}

	// When / Then: delete.
	response = serveAuthenticated(router, http.MethodDelete, "/api/varsler/subscriptions", fmt.Sprintf(`{"endpoint":%q}`, endpoint), "external-1")
	if response.Code != http.StatusNoContent {
		t.Fatalf("expected deletion status 204, got %d: %s", response.Code, response.Body.String())
	}
	if got := testutil.QueryInt(t, db, `SELECT COUNT(*) FROM web_push_subscriptions`); got != 0 {
		t.Fatalf("expected subscription to be deleted, got %d", got)
	}
}

func TestSubscriptionRoutes_EndpointOwnedByAnotherUserReturnsConflict(t *testing.T) {
	// Given
	db, logger := testutil.CreateTestDBAndLogger(t, "varsler_subscription_conflict")
	testutil.MustExec(t, db, `INSERT INTO users(id,external_id,email) VALUES (1,'external-1','one@example.com'),(2,'external-2','two@example.com')`)
	endpoint := "https://fcm.googleapis.com/fcm/send/shared-device"
	testutil.MustExec(t, db, `INSERT INTO web_push_subscriptions(user_id,endpoint,p256dh,auth) VALUES (1,?,'old','old')`, endpoint)
	service, err := New(db, logger, testConfig(t))
	if err != nil {
		t.Fatalf("create varsler service: %v", err)
	}
	router := chi.NewRouter()
	service.RegisterRoutes(router)
	p256dh, auth := validSubscriptionKeys()
	body := fmt.Sprintf(`{"endpoint":%q,"keys":{"p256dh":%q,"auth":%q}}`, endpoint, p256dh, auth)

	// When
	response := serveAuthenticated(router, http.MethodPost, "/api/varsler/subscriptions", body, "external-2")

	// Then
	if response.Code != http.StatusConflict {
		t.Fatalf("expected endpoint ownership conflict, got %d: %s", response.Code, response.Body.String())
	}
	if got := testutil.QueryInt(t, db, `SELECT user_id FROM web_push_subscriptions WHERE endpoint = ?`, endpoint); got != 1 {
		t.Fatalf("expected endpoint to remain with original owner, got user %d", got)
	}
}

func TestSubscriptionRoutes_DifferentUserCannotInspectOrDeleteEndpoint(t *testing.T) {
	// Given
	db, logger := testutil.CreateTestDBAndLogger(t, "varsler_subscription_owner_scope")
	testutil.MustExec(t, db, `INSERT INTO users(id,external_id,email) VALUES (1,'external-1','one@example.com'),(2,'external-2','two@example.com')`)
	endpoint := "https://fcm.googleapis.com/fcm/send/private-device"
	testutil.MustExec(t, db, `INSERT INTO web_push_subscriptions(user_id,endpoint,p256dh,auth) VALUES (1,?,'old','old')`, endpoint)
	service, err := New(db, logger, testConfig(t))
	if err != nil {
		t.Fatalf("create varsler service: %v", err)
	}
	router := chi.NewRouter()
	service.RegisterRoutes(router)

	// When / Then: another user cannot inspect ownership.
	response := serveAuthenticated(router, http.MethodGet, "/api/varsler/subscriptions?endpoint="+endpoint, "", "external-2")
	if response.Code != http.StatusOK || response.Body.String() != "{\"subscribed\":false}\n" {
		t.Fatalf("expected owner-scoped lookup, got %d: %s", response.Code, response.Body.String())
	}

	// When / Then: another user's delete is a no-op.
	response = serveAuthenticated(router, http.MethodDelete, "/api/varsler/subscriptions", fmt.Sprintf(`{"endpoint":%q}`, endpoint), "external-2")
	if response.Code != http.StatusNoContent {
		t.Fatalf("expected idempotent owner-scoped deletion, got %d: %s", response.Code, response.Body.String())
	}
	if got := testutil.QueryInt(t, db, `SELECT COUNT(*) FROM web_push_subscriptions WHERE user_id = 1`); got != 1 {
		t.Fatalf("expected original owner's subscription to remain, got %d", got)
	}
}

func TestSubscriptionRoutes_RejectsUnknownPushHost(t *testing.T) {
	// Given
	db, logger := testutil.CreateTestDBAndLogger(t, "varsler_subscription_host")
	testutil.MustExec(t, db, `INSERT INTO users(id,external_id,email) VALUES (1,'external-1','one@example.com')`)
	service, err := New(db, logger, testConfig(t))
	if err != nil {
		t.Fatalf("create varsler service: %v", err)
	}
	router := chi.NewRouter()
	service.RegisterRoutes(router)
	p256dh, auth := validSubscriptionKeys()
	body := fmt.Sprintf(`{"endpoint":"https://example.com/collect","keys":{"p256dh":%q,"auth":%q}}`, p256dh, auth)

	// When
	response := serveAuthenticated(router, http.MethodPost, "/api/varsler/subscriptions", body, "external-1")

	// Then
	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected unknown host rejection, got %d: %s", response.Code, response.Body.String())
	}
}

func TestSubscriptionRoutes_DisabledServiceReportsConfigAndRejectsMutation(t *testing.T) {
	// Given
	db, logger := testutil.CreateTestDBAndLogger(t, "varsler_disabled")
	service, err := New(db, logger, Config{})
	if err != nil {
		t.Fatalf("create disabled varsler service: %v", err)
	}
	router := chi.NewRouter()
	service.RegisterRoutes(router)

	// When / Then: configuration remains readable.
	response := serveAuthenticated(router, http.MethodGet, "/api/varsler/config", "", "external-1")
	if response.Code != http.StatusOK || response.Body.String() != "{\"enabled\":false,\"publicKey\":\"\"}\n" {
		t.Fatalf("expected disabled configuration, got %d: %s", response.Code, response.Body.String())
	}

	// When / Then: mutations are unavailable.
	response = serveAuthenticated(router, http.MethodPost, "/api/varsler/subscriptions", `{}`, "external-1")
	if response.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected disabled mutation status 503, got %d", response.Code)
	}
}

func serveAuthenticated(handler http.Handler, method, target, body, externalID string) *httptest.ResponseRecorder {
	request := httptest.NewRequest(method, target, bytes.NewBufferString(body))
	request.Header.Set("Content-Type", "application/json")
	request = request.WithContext(authctx.WithUserToken(request.Context(), externalID, "user@example.com"))
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	return response
}

func testConfig(t *testing.T) Config {
	t.Helper()
	privateKey := make([]byte, 32)
	privateKey[31] = 1
	private, err := ecdh.P256().NewPrivateKey(privateKey)
	if err != nil {
		t.Fatalf("create test VAPID key: %v", err)
	}
	return Config{
		Enabled:    true,
		PublicKey:  base64.RawURLEncoding.EncodeToString(private.PublicKey().Bytes()),
		PrivateKey: base64.RawURLEncoding.EncodeToString(privateKey),
		Subject:    "mailto:varsler@example.com",
	}
}

func validSubscriptionKeys() (string, string) {
	privateKey := make([]byte, 32)
	privateKey[31] = 3
	private, err := ecdh.P256().NewPrivateKey(privateKey)
	if err != nil {
		panic(err)
	}
	return base64.RawURLEncoding.EncodeToString(private.PublicKey().Bytes()), base64.RawURLEncoding.EncodeToString([]byte("0123456789abcdef"))
}
