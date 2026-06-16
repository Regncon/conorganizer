package login

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Regncon/conorganizer/service/authctx"
	"github.com/descope/go-sdk/descope"
	"github.com/go-chi/chi/v5"
)

func TestSessionRoute_ValidTokensStoresAuthCookies(t *testing.T) {
	// Given valid Descope tokens,
	// when the login page establishes an application session,
	// then HttpOnly auth cookies are stored before redirecting to post-login.

	// Given
	expectedStatusCode := http.StatusNoContent
	expectedSessionJWT := "refreshed-session-jwt"
	expectedRefreshJWT := "browser-refresh-jwt"
	validator := &fakeSessionValidator{
		sessionOK:    true,
		sessionToken: &descope.Token{JWT: expectedSessionJWT},
	}
	restoreSessionValidator(t, validator, nil)
	router := authTestRouter(t)
	request := httptest.NewRequest(
		http.MethodPost,
		"/auth/session",
		strings.NewReader(`{"sessionJwt":"browser-session-jwt","refreshJwt":"browser-refresh-jwt"}`),
	)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Forwarded-Proto", "https")
	recorder := httptest.NewRecorder()

	// When
	router.ServeHTTP(recorder, request)

	// Then
	if recorder.Code != expectedStatusCode {
		t.Fatalf("expected status %d, got %d", expectedStatusCode, recorder.Code)
	}
	if validator.sessionCalls != 1 {
		t.Fatalf("expected session validator to be called once, got %d", validator.sessionCalls)
	}
	if validator.refreshCalls != 0 {
		t.Fatalf("expected refresh validator not to be called, got %d", validator.refreshCalls)
	}
	if validator.validatedSessionToken != "browser-session-jwt" {
		t.Fatalf("expected session token passed to validator, got %q", validator.validatedSessionToken)
	}
	assertAuthCookie(t, recorder.Result(), authctx.SessionCookieName, expectedSessionJWT, true)
	assertAuthCookie(t, recorder.Result(), authctx.RefreshCookieName, expectedRefreshJWT, true)
}

func TestSessionRoute_MissingTokensReturnsBadRequest(t *testing.T) {
	// Given missing login tokens,
	// when the login page tries to establish an application session,
	// then the request is rejected before validation.

	// Given
	expectedStatusCode := http.StatusBadRequest
	validator := &fakeSessionValidator{}
	restoreSessionValidator(t, validator, nil)
	router := authTestRouter(t)
	request := httptest.NewRequest(
		http.MethodPost,
		"/auth/session",
		strings.NewReader(`{"sessionJwt":"browser-session-jwt"}`),
	)
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	// When
	router.ServeHTTP(recorder, request)

	// Then
	if recorder.Code != expectedStatusCode {
		t.Fatalf("expected status %d, got %d", expectedStatusCode, recorder.Code)
	}
	if validator.sessionCalls != 0 {
		t.Fatalf("expected session validator not to be called, got %d calls", validator.sessionCalls)
	}
	if validator.refreshCalls != 0 {
		t.Fatalf("expected refresh validator not to be called, got %d calls", validator.refreshCalls)
	}
	assertNoAuthCookies(t, recorder.Result())
}

func TestSessionRoute_ExpiredSessionRefreshesBeforeStoringAuthCookies(t *testing.T) {
	// Given an expired session token and valid refresh token,
	// when the login page establishes an application session,
	// then the refreshed session is stored with the original refresh token.

	// Given
	expectedStatusCode := http.StatusNoContent
	expectedSessionJWT := "refreshed-session-jwt"
	expectedRefreshJWT := "browser-refresh-jwt"
	validator := &fakeSessionValidator{
		sessionOK:  false,
		sessionErr: errors.New("expired token"),
		refreshOK:  true,
		refreshToken: &descope.Token{
			JWT: expectedSessionJWT,
		},
	}
	restoreSessionValidator(t, validator, nil)
	router := authTestRouter(t)
	request := httptest.NewRequest(
		http.MethodPost,
		"/auth/session",
		strings.NewReader(`{"sessionJwt":"expired-session-jwt","refreshJwt":"browser-refresh-jwt"}`),
	)
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	// When
	router.ServeHTTP(recorder, request)

	// Then
	if recorder.Code != expectedStatusCode {
		t.Fatalf("expected status %d, got %d", expectedStatusCode, recorder.Code)
	}
	if validator.sessionCalls != 1 {
		t.Fatalf("expected session validator to be called once, got %d", validator.sessionCalls)
	}
	if validator.refreshCalls != 1 {
		t.Fatalf("expected refresh validator to be called once, got %d", validator.refreshCalls)
	}
	assertAuthCookie(t, recorder.Result(), authctx.SessionCookieName, expectedSessionJWT, false)
	assertAuthCookie(t, recorder.Result(), authctx.RefreshCookieName, expectedRefreshJWT, false)
}

func TestSessionRoute_InvalidTokensReturnsUnauthorized(t *testing.T) {
	// Given rejected Descope tokens,
	// when the login page tries to establish an application session,
	// then no auth cookies are stored.

	// Given
	expectedStatusCode := http.StatusUnauthorized
	validator := &fakeSessionValidator{
		sessionOK:  false,
		sessionErr: errors.New("invalid token"),
		refreshOK:  false,
		refreshErr: errors.New("invalid refresh token"),
	}
	restoreSessionValidator(t, validator, nil)
	router := authTestRouter(t)
	request := httptest.NewRequest(
		http.MethodPost,
		"/auth/session",
		strings.NewReader(`{"sessionJwt":"browser-session-jwt","refreshJwt":"browser-refresh-jwt"}`),
	)
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	// When
	router.ServeHTTP(recorder, request)

	// Then
	if recorder.Code != expectedStatusCode {
		t.Fatalf("expected status %d, got %d", expectedStatusCode, recorder.Code)
	}
	if validator.sessionCalls != 1 {
		t.Fatalf("expected session validator to be called once, got %d", validator.sessionCalls)
	}
	if validator.refreshCalls != 1 {
		t.Fatalf("expected refresh validator to be called once, got %d", validator.refreshCalls)
	}
	assertNoAuthCookies(t, recorder.Result())
}

type fakeSessionValidator struct {
	sessionOK             bool
	sessionToken          *descope.Token
	sessionErr            error
	refreshOK             bool
	refreshToken          *descope.Token
	refreshErr            error
	sessionCalls          int
	refreshCalls          int
	validatedSessionToken string
	validatedRefreshToken string
}

func (f *fakeSessionValidator) ValidateSessionWithToken(_ context.Context, sessionToken string) (bool, *descope.Token, error) {
	f.sessionCalls++
	f.validatedSessionToken = sessionToken
	return f.sessionOK, f.sessionToken, f.sessionErr
}

func (f *fakeSessionValidator) RefreshSessionWithToken(_ context.Context, refreshToken string) (bool, *descope.Token, error) {
	f.refreshCalls++
	f.validatedRefreshToken = refreshToken
	return f.refreshOK, f.refreshToken, f.refreshErr
}

func restoreSessionValidator(t *testing.T, validator authctx.SessionValidator, err error) {
	t.Helper()
	previous := newSessionValidator
	newSessionValidator = func() (authctx.SessionValidator, error) {
		return validator, err
	}
	t.Cleanup(func() {
		newSessionValidator = previous
	})
}

func authTestRouter(t *testing.T) chi.Router {
	t.Helper()
	router := chi.NewRouter()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	if err := SetupAuthRoute(router, nil, logger); err != nil {
		t.Fatalf("expected auth route setup to succeed: %v", err)
	}
	return router
}

func assertAuthCookie(t *testing.T, response *http.Response, name, expectedValue string, expectedSecure bool) {
	t.Helper()
	cookie := findCookie(response, name)
	if cookie == nil {
		t.Fatalf("expected %s cookie to be set", name)
	}
	if cookie.Value != expectedValue {
		t.Fatalf("expected %s cookie value %q, got %q", name, expectedValue, cookie.Value)
	}
	if cookie.Path != "/" {
		t.Fatalf("expected %s cookie path /, got %q", name, cookie.Path)
	}
	if cookie.MaxAge != 31536000 {
		t.Fatalf("expected %s cookie max age 31536000, got %d", name, cookie.MaxAge)
	}
	if !cookie.HttpOnly {
		t.Fatalf("expected %s cookie to be HttpOnly", name)
	}
	if cookie.Secure != expectedSecure {
		t.Fatalf("expected %s cookie secure=%v, got %v", name, expectedSecure, cookie.Secure)
	}
	if cookie.SameSite != http.SameSiteLaxMode {
		t.Fatalf("expected %s cookie SameSite=Lax, got %v", name, cookie.SameSite)
	}
}

func assertNoAuthCookies(t *testing.T, response *http.Response) {
	t.Helper()
	for _, name := range []string{authctx.SessionCookieName, authctx.RefreshCookieName} {
		if cookie := findCookie(response, name); cookie != nil {
			t.Fatalf("expected no %s cookie, got %q", name, cookie.Value)
		}
	}
}

func findCookie(response *http.Response, name string) *http.Cookie {
	for _, cookie := range response.Cookies() {
		if cookie.Name == name {
			return cookie
		}
	}
	return nil
}
