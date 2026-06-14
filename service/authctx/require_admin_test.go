package authctx

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/descope/go-sdk/descope"

	"github.com/Regncon/conorganizer/testutil/bdd"
)

func TestRequireAdmin_WhenUserIsNotAdmin_ReturnsForbidden(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at en innlogget bruker mangler adminrolle.",
		When:  "Når brukeren åpner en adminbeskyttet rute.",
		Then:  "Så skal tilgangen avvises før handleren kjøres.",
	})

	// Given
	expectedStatusCode := http.StatusForbidden
	expectedBody := "You are not an admin\n"
	expectedHandlerCalled := false

	handlerCalled := false
	handler := RequireAdmin(discardLogger())(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.WriteHeader(http.StatusOK)
	}))
	request := httptest.NewRequest(http.MethodGet, "/admin", nil)
	request = request.WithContext(context.WithValue(request.Context(), ctxUserToken, userTokenWithRoles("User")))
	recorder := httptest.NewRecorder()

	// When
	handler.ServeHTTP(recorder, request)

	// Then
	if recorder.Code != expectedStatusCode {
		t.Fatalf("HTTP status mismatch\nexpected: %d\nactual:   %d", expectedStatusCode, recorder.Code)
	}
	if recorder.Body.String() != expectedBody {
		t.Fatalf("HTTP body mismatch\nexpected: %q\nactual:   %q", expectedBody, recorder.Body.String())
	}
	if handlerCalled != expectedHandlerCalled {
		t.Fatalf("handler call mismatch\nexpected: %v\nactual:   %v", expectedHandlerCalled, handlerCalled)
	}
}

func TestRequireAdmin_WhenUserIsAdmin_AllowsRequest(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at en innlogget bruker har adminrolle.",
		When:  "Når brukeren åpner en adminbeskyttet rute.",
		Then:  "Så skal handleren få behandle forespørselen.",
	})

	// Given
	expectedStatusCode := http.StatusOK
	expectedBody := "admin ok"
	expectedHandlerCalled := true

	handlerCalled := false
	handler := RequireAdmin(discardLogger())(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		_, _ = w.Write([]byte(expectedBody))
	}))
	request := httptest.NewRequest(http.MethodGet, "/admin", nil)
	request = request.WithContext(context.WithValue(request.Context(), ctxUserToken, userTokenWithRoles("Admin")))
	recorder := httptest.NewRecorder()

	// When
	handler.ServeHTTP(recorder, request)

	// Then
	if recorder.Code != expectedStatusCode {
		t.Fatalf("HTTP status mismatch\nexpected: %d\nactual:   %d", expectedStatusCode, recorder.Code)
	}
	if recorder.Body.String() != expectedBody {
		t.Fatalf("HTTP body mismatch\nexpected: %q\nactual:   %q", expectedBody, recorder.Body.String())
	}
	if handlerCalled != expectedHandlerCalled {
		t.Fatalf("handler call mismatch\nexpected: %v\nactual:   %v", expectedHandlerCalled, handlerCalled)
	}
}

func userTokenWithRoles(roles ...string) *descope.Token {
	roleClaims := make([]any, 0, len(roles))
	for _, role := range roles {
		roleClaims = append(roleClaims, role)
	}
	return &descope.Token{Claims: map[string]any{"roles": roleClaims}}
}

func discardLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}
