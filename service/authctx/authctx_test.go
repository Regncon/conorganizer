package authctx

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/descope/go-sdk/descope"
)

func TestAuthMiddleware_AnonymousRequestsDoNotCallDescope(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "A request has no authentication tokens.",
		When:  "The authentication middleware handles the request.",
		Then:  "The handler runs anonymously without validating or refreshing a token.",
	})

	// Given
	expectedStatus := http.StatusNoContent
	validator := &fakeSessionValidator{}
	handler := AuthMiddleware(validator, discardLogger())(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if token, err := GetUserTokenFromContext(r.Context()); token != nil || err == nil {
			t.Errorf("expected anonymous request, got token %v and error %v", token, err)
		}
		w.WriteHeader(expectedStatus)
	}))
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/", nil)

	// When
	handler.ServeHTTP(recorder, request)

	// Then
	if recorder.Code != expectedStatus {
		t.Fatalf("expected status %d, got %d", expectedStatus, recorder.Code)
	}
	if validator.sessionCalls != 0 || validator.refreshCalls != 0 {
		t.Fatalf("anonymous request called Descope: validation=%d refresh=%d", validator.sessionCalls, validator.refreshCalls)
	}
}

func TestAuthMiddleware_ValidSessionUsesInjectedValidatorWithoutRefreshing(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "A request has a valid session and a refresh cookie.",
		When:  "The authentication middleware handles the request.",
		Then:  "The injected validator supplies the user without refreshing or replacing cookies.",
	})

	// Given
	expectedToken := &descope.Token{ID: "user-123", JWT: "valid-session"}
	validator := &fakeSessionValidator{sessionOK: true, sessionToken: expectedToken}
	handler := AuthMiddleware(validator, discardLogger())(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := GetUserTokenFromContext(r.Context())
		if err != nil || token != expectedToken {
			t.Errorf("expected validated user token, got token %v and error %v", token, err)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	request := httptest.NewRequest(http.MethodGet, "/profile", nil)
	request.AddCookie(&http.Cookie{Name: SessionCookieName, Value: expectedToken.JWT})
	request.AddCookie(&http.Cookie{Name: RefreshCookieName, Value: "refresh-token"})
	recorder := httptest.NewRecorder()

	// When
	handler.ServeHTTP(recorder, request)

	// Then
	if recorder.Code != http.StatusNoContent {
		t.Fatalf("expected handler to run, got status %d", recorder.Code)
	}
	if validator.sessionCalls != 1 || validator.refreshCalls != 0 {
		t.Fatalf("expected one validation and no refresh, got validation=%d refresh=%d", validator.sessionCalls, validator.refreshCalls)
	}
	if len(recorder.Result().Cookies()) != 0 {
		t.Fatal("valid session should not replace cookies")
	}
}

func TestAuthMiddleware_ExpiredSessionSetsRefreshedCookieBeforeHandler(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "A request has an expired session and a valid refresh token.",
		When:  "The authentication middleware refreshes the session.",
		Then:  "The handler receives the refreshed identity and the response already contains the new session cookie.",
	})

	// Given
	expectedToken := &descope.Token{ID: "user-123", JWT: "new-session"}
	validator := &fakeSessionValidator{
		sessionErr:   errors.New("expired session"),
		refreshOK:    true,
		refreshToken: expectedToken,
	}
	handler := AuthMiddleware(validator, discardLogger())(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := GetUserTokenFromContext(r.Context())
		if err != nil || token != expectedToken {
			t.Errorf("expected refreshed user token, got token %v and error %v", token, err)
		}
		if len(w.Header().Values("Set-Cookie")) != 1 {
			t.Error("expected refreshed cookie before handler writes headers")
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	request := httptest.NewRequest(http.MethodGet, "/profile/api", nil)
	request.AddCookie(&http.Cookie{Name: SessionCookieName, Value: "expired-session"})
	request.AddCookie(&http.Cookie{Name: RefreshCookieName, Value: "refresh-token"})
	recorder := httptest.NewRecorder()

	// When
	handler.ServeHTTP(recorder, request)

	// Then
	if recorder.Code != http.StatusNoContent {
		t.Fatalf("expected handler to run, got status %d", recorder.Code)
	}
	if validator.sessionCalls != 1 || validator.refreshCalls != 1 {
		t.Fatalf("expected one validation and one refresh, got validation=%d refresh=%d", validator.sessionCalls, validator.refreshCalls)
	}
	cookies := recorder.Result().Cookies()
	if len(cookies) != 1 || cookies[0].Name != SessionCookieName || cookies[0].Value != expectedToken.JWT {
		t.Fatalf("expected only the refreshed session cookie, got %v", cookies)
	}
}

func TestAuthMiddleware_FailedRefreshPreservesExistingCookies(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Session validation and refresh both fail.",
		When:  "The authentication middleware handles the request.",
		Then:  "It preserves existing cookies and continues with the existing unauthenticated behavior.",
	})

	// Given
	expectedStatus := http.StatusUnauthorized
	validator := &fakeSessionValidator{
		sessionErr: errors.New("expired session"),
		refreshErr: errors.New("refresh unavailable"),
	}
	handler := AuthMiddleware(validator, discardLogger())(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if token, err := GetUserTokenFromContext(r.Context()); token != nil || err == nil {
			t.Errorf("expected no authenticated identity, got token %v and error %v", token, err)
		}
		w.WriteHeader(expectedStatus)
	}))
	request := httptest.NewRequest(http.MethodPost, "/profile/api/create", nil)
	request.AddCookie(&http.Cookie{Name: SessionCookieName, Value: "expired-session"})
	request.AddCookie(&http.Cookie{Name: RefreshCookieName, Value: "refresh-token"})
	recorder := httptest.NewRecorder()

	// When
	handler.ServeHTTP(recorder, request)

	// Then
	if recorder.Code != expectedStatus {
		t.Fatalf("expected status %d, got %d", expectedStatus, recorder.Code)
	}
	if len(recorder.Result().Cookies()) != 0 {
		t.Fatal("authentication failure must not clear or replace cookies")
	}
}

type fakeSessionValidator struct {
	sessionOK      bool
	sessionToken   *descope.Token
	sessionErr     error
	refreshOK      bool
	refreshToken   *descope.Token
	refreshErr     error
	sessionCalls   int
	refreshCalls   int
	sessionContext context.Context
	refreshContext context.Context
}

func (v *fakeSessionValidator) ValidateSessionWithToken(ctx context.Context, _ string) (bool, *descope.Token, error) {
	v.sessionCalls++
	v.sessionContext = ctx
	return v.sessionOK, v.sessionToken, v.sessionErr
}

func (v *fakeSessionValidator) RefreshSessionWithToken(ctx context.Context, _ string) (bool, *descope.Token, error) {
	v.refreshCalls++
	v.refreshContext = ctx
	return v.refreshOK, v.refreshToken, v.refreshErr
}
