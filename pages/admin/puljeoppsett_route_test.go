package admin

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/service/authctx"
	"github.com/Regncon/conorganizer/service/live"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/go-chi/chi/v5"
)

func authedRequest(method, target string) *http.Request {
	req := httptest.NewRequest(method, target, nil)
	return req.WithContext(authctx.WithUserToken(req.Context(), "ext-42", "admin@x.no"))
}

func TestPuljeoppsettRoute_AddAndRemoveMembership(t *testing.T) {
	db, logger := testutil.CreateTestDBAndLogger(t, "puljeoppsett_route")
	testutil.MustExec(t, db, `INSERT INTO users (id, external_id, email, is_admin) VALUES (42, 'ext-42', 'admin@x.no', 1)`)
	testutil.MustExec(t, db, `INSERT INTO puljer (id, name, status, start_at, end_at) VALUES (?,?,?,?,?)`,
		string(models.PuljeFredagKveld), "Fredag", "Open", "2026-01-01 18:00", "2026-01-01 22:00")
	insertBoardEvent(t, db, "e1", "Spel", "Godkjent", "Default", 0, 0, "ola@x.no", "Ola")

	router := chi.NewRouter()
	puljeoppsettRoute(router, db, &live.Manager{}, logger, nil)

	addRec := httptest.NewRecorder()
	router.ServeHTTP(addRec, authedRequest(http.MethodPut,
		"/api/puljeoppsett/e1/"+string(models.PuljeFredagKveld)+"/add"))
	if addRec.Code != http.StatusNoContent {
		t.Fatalf("add status = %d, want 204\nbody: %s", addRec.Code, addRec.Body.String())
	}
	if got := testutil.QueryInt(t, db,
		`SELECT is_in_pulje FROM relation_event_puljer WHERE event_id='e1' AND pulje_id=?`,
		string(models.PuljeFredagKveld)); got != 1 {
		t.Fatalf("after add is_in_pulje = %d, want 1", got)
	}

	delRec := httptest.NewRecorder()
	router.ServeHTTP(delRec, authedRequest(http.MethodDelete,
		"/api/puljeoppsett/e1/"+string(models.PuljeFredagKveld)))
	if delRec.Code != http.StatusNoContent {
		t.Fatalf("remove status = %d, want 204\nbody: %s", delRec.Code, delRec.Body.String())
	}
	if got := testutil.QueryInt(t, db,
		`SELECT is_in_pulje FROM relation_event_puljer WHERE event_id='e1' AND pulje_id=?`,
		string(models.PuljeFredagKveld)); got != 0 {
		t.Fatalf("after remove is_in_pulje = %d, want 0", got)
	}
}

func TestPuljeoppsettRoute_RejectsBadPulje(t *testing.T) {
	db, logger := testutil.CreateTestDBAndLogger(t, "puljeoppsett_route_bad")
	router := chi.NewRouter()
	puljeoppsettRoute(router, db, &live.Manager{}, logger, nil)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, authedRequest(http.MethodPut, "/api/puljeoppsett/e1/NotAPulje/add"))
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("bad pulje status = %d, want 400", rec.Code)
	}
}
