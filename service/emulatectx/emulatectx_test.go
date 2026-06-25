package emulatectx_test

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/descope/go-sdk/descope"

	"github.com/Regncon/conorganizer/service/authctx"
	"github.com/Regncon/conorganizer/service/emulatectx"
	"github.com/Regncon/conorganizer/service/userctx"
	"github.com/Regncon/conorganizer/testutil"
)

const (
	targetUserID     = 42
	targetExternalID = "player-external-42"
	targetEmail      = "player@example.com"
	adminExternalID  = "admin-external-1"
	adminEmail       = "admin@example.com"
	playerFacingPath = "/"
	adminFacingPath  = "/admin/dashboard"
)

func adminToken() *descope.Token {
	return &descope.Token{
		ID:     adminExternalID,
		Claims: map[string]any{"roles": []any{"Admin"}, "email": adminEmail},
	}
}

func nonAdminToken() *descope.Token {
	return &descope.Token{
		ID:     "some-user",
		Claims: map[string]any{"roles": []any{"User"}, "email": "user@example.com"},
	}
}

func TestEmulateMiddleware_AdminWithCookie_SwapsIdentityToPlayer(t *testing.T) {
	db := testutil.CreateTestDB(t, "emulatectx_swap")
	testutil.MustExec(t, db, `INSERT INTO users (id, external_id, email, is_admin) VALUES (?, ?, ?, 0)`, targetUserID, targetExternalID, targetEmail)

	var info struct {
		id          string
		email       string
		isAdmin     bool
		emulating   bool
		hasReal     bool
		realIsAdmin bool
	}
	handler := emulatectx.EmulateMiddleware(db, testutil.NewTestLogger())(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userInfo := userctx.GetUserRequestInfo(ctx)
		info.id = userInfo.Id
		info.email = userInfo.Email
		info.isAdmin = userInfo.IsAdmin
		info.emulating = emulatectx.IsEmulating(ctx)
		real, ok := emulatectx.RealAdminTokenFromContext(ctx)
		info.hasReal = ok
		if ok {
			info.realIsAdmin = authctx.GetAdminFromUserToken(authctx.WithUserToken(ctx, real))
		}
		w.WriteHeader(http.StatusOK)
	}))

	request := httptest.NewRequest(http.MethodGet, playerFacingPath, nil)
	request = request.WithContext(authctx.WithUserToken(request.Context(), adminToken()))
	request.AddCookie(&http.Cookie{Name: emulatectx.CookieName, Value: strconv.Itoa(targetUserID)})
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status mismatch\nexpected: %d\nactual:   %d", http.StatusOK, recorder.Code)
	}
	if info.id != targetExternalID {
		t.Fatalf("emulated id mismatch\nexpected: %q\nactual:   %q", targetExternalID, info.id)
	}
	if info.email != targetEmail {
		t.Fatalf("emulated email mismatch\nexpected: %q\nactual:   %q", targetEmail, info.email)
	}
	if info.isAdmin {
		t.Fatalf("emulated player must not be admin")
	}
	if !info.emulating {
		t.Fatalf("expected IsEmulating to be true")
	}
	if !info.hasReal || !info.realIsAdmin {
		t.Fatalf("expected real admin token to be stashed and still admin")
	}
}

func TestEmulateMiddleware_NonAdminWithCookie_IsIgnoredAndCleared(t *testing.T) {
	db := testutil.CreateTestDB(t, "emulatectx_nonadmin")
	testutil.MustExec(t, db, `INSERT INTO users (id, external_id, email, is_admin) VALUES (?, ?, ?, 0)`, targetUserID, targetExternalID, targetEmail)

	var sawEmulation bool
	var sawID string
	handler := emulatectx.EmulateMiddleware(db, testutil.NewTestLogger())(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sawEmulation = emulatectx.IsEmulating(r.Context())
		sawID = userctx.GetUserRequestInfo(r.Context()).Id
		w.WriteHeader(http.StatusOK)
	}))

	request := httptest.NewRequest(http.MethodGet, playerFacingPath, nil)
	request = request.WithContext(authctx.WithUserToken(request.Context(), nonAdminToken()))
	request.AddCookie(&http.Cookie{Name: emulatectx.CookieName, Value: strconv.Itoa(targetUserID)})
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if sawEmulation {
		t.Fatalf("non-admin must never emulate")
	}
	if sawID != "some-user" {
		t.Fatalf("non-admin identity must be unchanged\nexpected: %q\nactual:   %q", "some-user", sawID)
	}
	if !cookieCleared(recorder) {
		t.Fatalf("expected emulate cookie to be cleared for non-admin")
	}
}

func TestEmulateMiddleware_AdminPathWhileEmulating_AutoExits(t *testing.T) {
	db := testutil.CreateTestDB(t, "emulatectx_autoexit")
	testutil.MustExec(t, db, `INSERT INTO users (id, external_id, email, is_admin) VALUES (?, ?, ?, 0)`, targetUserID, targetExternalID, targetEmail)

	var emulating bool
	var isAdmin bool
	handler := emulatectx.EmulateMiddleware(db, testutil.NewTestLogger())(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		emulating = emulatectx.IsEmulating(r.Context())
		isAdmin = userctx.GetUserRequestInfo(r.Context()).IsAdmin
		w.WriteHeader(http.StatusOK)
	}))

	request := httptest.NewRequest(http.MethodGet, adminFacingPath, nil)
	request = request.WithContext(authctx.WithUserToken(request.Context(), adminToken()))
	request.AddCookie(&http.Cookie{Name: emulatectx.CookieName, Value: strconv.Itoa(targetUserID)})
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if emulating {
		t.Fatalf("admin route must auto-exit emulation")
	}
	if !isAdmin {
		t.Fatalf("admin identity must be restored on admin route")
	}
	if !cookieCleared(recorder) {
		t.Fatalf("expected emulate cookie to be cleared on auto-exit")
	}
}

func TestEmulateMiddleware_UnsafeMethodWhileEmulating_IsBlocked(t *testing.T) {
	db := testutil.CreateTestDB(t, "emulatectx_readonly")
	testutil.MustExec(t, db, `INSERT INTO users (id, external_id, email, is_admin) VALUES (?, ?, ?, 0)`, targetUserID, targetExternalID, targetEmail)

	handlerCalled := false
	handler := emulatectx.EmulateMiddleware(db, testutil.NewTestLogger())(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.WriteHeader(http.StatusOK)
	}))

	request := httptest.NewRequest(http.MethodPost, playerFacingPath, nil)
	request = request.WithContext(authctx.WithUserToken(request.Context(), adminToken()))
	request.AddCookie(&http.Cookie{Name: emulatectx.CookieName, Value: strconv.Itoa(targetUserID)})
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if handlerCalled {
		t.Fatalf("write must not reach handler while emulating")
	}
	if recorder.Code != http.StatusForbidden {
		t.Fatalf("status mismatch\nexpected: %d\nactual:   %d", http.StatusForbidden, recorder.Code)
	}
}

func TestEmulateMiddleware_MissingTargetUser_FallsBackToAdmin(t *testing.T) {
	db := testutil.CreateTestDB(t, "emulatectx_missing")

	var emulating bool
	var isAdmin bool
	handler := emulatectx.EmulateMiddleware(db, testutil.NewTestLogger())(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		emulating = emulatectx.IsEmulating(r.Context())
		isAdmin = userctx.GetUserRequestInfo(r.Context()).IsAdmin
		w.WriteHeader(http.StatusOK)
	}))

	request := httptest.NewRequest(http.MethodGet, playerFacingPath, nil)
	request = request.WithContext(authctx.WithUserToken(request.Context(), adminToken()))
	request.AddCookie(&http.Cookie{Name: emulatectx.CookieName, Value: "9999"})
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if emulating {
		t.Fatalf("missing target must not emulate")
	}
	if !isAdmin {
		t.Fatalf("admin identity must remain when target is missing")
	}
	if !cookieCleared(recorder) {
		t.Fatalf("expected emulate cookie to be cleared when target missing")
	}
}

func cookieCleared(recorder *httptest.ResponseRecorder) bool {
	for _, c := range recorder.Result().Cookies() {
		if c.Name == emulatectx.CookieName && c.MaxAge < 0 {
			return true
		}
	}
	return false
}
