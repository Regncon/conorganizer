package billettholderadmin

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"github.com/Regncon/conorganizer/service/emulatectx"
	"github.com/Regncon/conorganizer/testutil"
)

func setupEmulateRouteTest(t testing.TB) (chi.Router, *testutil.StubLogger) {
	t.Helper()

	db, logger := testutil.CreateTestDBAndLogger(t, "admin_emulate_route")
	insertAdminRouteTestBillettholder(t, db, 7)
	insertAdminRouteTestUser(t, db, 42, "player@example.com")
	insertAdminRouteTestBillettholderUserAssociation(t, db, 7, 42)

	// A billettholder with no linked user account.
	insertAdminRouteTestBillettholder(t, db, 8)

	router := chi.NewRouter()
	EmulatePlayerRoute(router, db, logger)
	return router, nil
}

func TestEmulatePlayerRoute_WithLinkedUser_SetsCookieAndRedirects(t *testing.T) {
	router, _ := setupEmulateRouteTest(t)

	request := httptest.NewRequest(http.MethodPost, "/admin/emulate/7/", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusSeeOther {
		t.Fatalf("status mismatch\nexpected: %d\nactual:   %d", http.StatusSeeOther, recorder.Code)
	}
	if location := recorder.Header().Get("Location"); location != "/" {
		t.Fatalf("redirect mismatch\nexpected: %q\nactual:   %q", "/", location)
	}
	var found *http.Cookie
	for _, c := range recorder.Result().Cookies() {
		if c.Name == emulatectx.CookieName {
			found = c
		}
	}
	if found == nil {
		t.Fatalf("expected %s cookie to be set", emulatectx.CookieName)
	}
	if found.Value != "42" {
		t.Fatalf("cookie value mismatch\nexpected: %q\nactual:   %q", "42", found.Value)
	}
}

func TestEmulatePlayerRoute_NoLinkedUser_ReturnsFriendlyErrorWithoutCookie(t *testing.T) {
	router, _ := setupEmulateRouteTest(t)

	request := httptest.NewRequest(http.MethodPost, "/admin/emulate/8/", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code == http.StatusSeeOther {
		t.Fatalf("must not redirect when billettholder has no account")
	}
	for _, c := range recorder.Result().Cookies() {
		if c.Name == emulatectx.CookieName && c.Value != "" {
			t.Fatalf("must not set emulate cookie when no linked user")
		}
	}
}
