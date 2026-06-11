package login

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Regncon/conorganizer/service/authctx"
	"github.com/go-chi/chi/v5"
)

func TestLogout_ClearsSessionAndRefreshCookies(t *testing.T) {
	// Gitt at en innlogget bruker har sesjons- og refresh-cookie,
	// når brukeren logger ut,
	// så skal begge cookies slettes.

	// Given
	expectedStatusCode := http.StatusOK
	expectedExpiredCookieNames := []string{
		authctx.SessionCookieName,
		authctx.RefreshCookieName,
	}
	expectedCookieMaxAge := -1
	expectedCookieValue := ""
	expectedCookiePath := "/"

	router := chi.NewRouter()
	if err := SetupAuthRoute(router, nil, discardLogger()); err != nil {
		t.Fatalf("setup auth route: %v", err)
	}
	request := httptest.NewRequest(http.MethodGet, "/auth/logout", nil)
	request.AddCookie(&http.Cookie{Name: authctx.SessionCookieName, Value: "session"})
	request.AddCookie(&http.Cookie{Name: authctx.RefreshCookieName, Value: "refresh"})
	recorder := httptest.NewRecorder()

	// When
	router.ServeHTTP(recorder, request)

	// Then
	if recorder.Code != expectedStatusCode {
		t.Fatalf("HTTP status mismatch\nexpected: %d\nactual:   %d", expectedStatusCode, recorder.Code)
	}
	for _, cookieName := range expectedExpiredCookieNames {
		cookie := responseCookie(t, recorder, cookieName)
		if cookie.Value != expectedCookieValue {
			t.Fatalf("%s value mismatch\nexpected: %q\nactual:   %q", cookieName, expectedCookieValue, cookie.Value)
		}
		if cookie.MaxAge != expectedCookieMaxAge {
			t.Fatalf("%s MaxAge mismatch\nexpected: %d\nactual:   %d", cookieName, expectedCookieMaxAge, cookie.MaxAge)
		}
		if cookie.Path != expectedCookiePath {
			t.Fatalf("%s path mismatch\nexpected: %q\nactual:   %q", cookieName, expectedCookiePath, cookie.Path)
		}
		if !cookie.HttpOnly {
			t.Fatalf("expected %s to be HttpOnly", cookieName)
		}
		if cookie.SameSite != http.SameSiteLaxMode {
			t.Fatalf("%s SameSite mismatch\nexpected: %v\nactual:   %v", cookieName, http.SameSiteLaxMode, cookie.SameSite)
		}
	}
}

func responseCookie(t *testing.T, recorder *httptest.ResponseRecorder, name string) *http.Cookie {
	t.Helper()

	for _, cookie := range recorder.Result().Cookies() {
		if cookie.Name == name {
			return cookie
		}
	}
	t.Fatalf("expected response cookie %q to exist", name)
	return nil
}

func discardLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}
