package admin

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/service/live"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/delaneyj/toolbelt/embeddednats"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
	natsserver "github.com/nats-io/nats-server/v2/server"
)

func TestRoomAssignment_ChangesOnlyActiveEventInSelectedPulje(t *testing.T) {
	// Given
	ns, err := embeddednats.New(context.Background(), embeddednats.WithNATSServerOptions(&natsserver.Options{
		Host: "127.0.0.1", Port: -1, JetStream: true, StoreDir: t.TempDir(),
	}))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { ns.Close() })
	ns.WaitForServer()
	manager, err := live.NewManager(context.Background(), ns, sessions.NewCookieStore([]byte("room-assignment-test")))
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name   string
		active int
		event  string
		status int
		roomID int
		method string
		target string
	}{
		{name: "assign active", active: 1, event: "event-1", status: http.StatusOK, roomID: 42, method: http.MethodPost, target: "42"},
		{name: "assign inactive", active: 0, event: "event-1", status: http.StatusOK, roomID: 42, method: http.MethodPost, target: "42"},
		{name: "assign missing", active: 1, event: "missing", status: http.StatusConflict, roomID: 41, method: http.MethodPost, target: "42"},
		{name: "remove active", active: 1, event: "event-1", status: http.StatusOK, roomID: 0, method: http.MethodDelete, target: "41"},
		{name: "remove stale", active: 1, event: "event-1", status: http.StatusConflict, roomID: 41, method: http.MethodDelete, target: "42"},
		{name: "remove inactive", active: 0, event: "event-1", status: http.StatusConflict, roomID: 41, method: http.MethodDelete, target: "41"},
		{name: "remove missing", active: 1, event: "missing", status: http.StatusConflict, roomID: 41, method: http.MethodDelete, target: "41"},
		{name: "remove invalid room", active: 1, event: "event-1", status: http.StatusBadRequest, roomID: 41, method: http.MethodDelete, target: "invalid"},
		{name: "remove nonpositive room", active: 1, event: "event-1", status: http.StatusBadRequest, roomID: 41, method: http.MethodDelete, target: "0"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			// Given
			db, logger := testutil.CreateTestDBAndLogger(t, "room_assignment_route")
			for _, pulje := range []models.Pulje{models.PuljeFredagKveld, models.PuljeLordagKveld} {
				testutil.MustExec(t, db, `INSERT INTO puljer(id,name,status,start_at,end_at) VALUES(?,?,'Open','2026-10-09T18:00:00Z','2026-10-09T23:00:00Z')`, pulje, pulje)
			}
			testutil.MustExec(t, db, `INSERT INTO rooms(id,name,room_number,floor,max_concurrent_games) VALUES(41,'Første','706',7,0),(42,'Andre','705',7,0)`)
			insertBoardEvent(t, db, "event-1", "Spill", "Godkjent", "Default", 0, 0, "host@example.com", "Host")
			testutil.MustExec(t, db, `INSERT INTO relation_event_puljer(event_id,pulje_id,is_in_pulje,is_published,room_id) VALUES('event-1',?,?,1,41)`, models.PuljeFredagKveld, tc.active)
			testutil.MustExec(t, db, `INSERT INTO relation_event_puljer(event_id,pulje_id,is_in_pulje,is_published,room_id) VALUES('event-1',?,1,1,41)`, models.PuljeLordagKveld)
			router := chi.NewRouter()
			if err := SetupAdminRoute(router, logger, manager, db, nil); err != nil {
				t.Fatal(err)
			}

			// When
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, httptest.NewRequest(tc.method, fmt.Sprintf("/admin/rooms/api/assignment/FredagKveld/%s/%s", tc.event, tc.target), nil))

			// Then
			if recorder.Code != tc.status {
				t.Fatalf("expected status %d, got %d: %s", tc.status, recorder.Code, recorder.Body.String())
			}
			if got := testutil.QueryInt(t, db, `SELECT COALESCE(room_id, 0) FROM relation_event_puljer WHERE event_id='event-1' AND pulje_id=?`, models.PuljeFredagKveld); got != tc.roomID {
				t.Fatalf("Friday room = %d, want %d", got, tc.roomID)
			}
			if got := testutil.QueryInt(t, db, `SELECT room_id FROM relation_event_puljer WHERE event_id='event-1' AND pulje_id=?`, models.PuljeLordagKveld); got != 41 {
				t.Fatalf("Saturday assignment changed to %d", got)
			}
			if got := testutil.QueryInt(t, db, `SELECT is_published FROM relation_event_puljer WHERE event_id='event-1' AND pulje_id=?`, models.PuljeFredagKveld); got != 1 {
				t.Fatal("room assignment changed publishing state")
			}
			if tc.status == http.StatusOK && !strings.Contains(recorder.Body.String(), "room-saved") {
				t.Fatal("successful assignment should notify the page")
			}
		})
	}
}

func TestRoomAssignment_CreatesAndPublishesRelationForApprovedEvent(t *testing.T) {
	ns, err := embeddednats.New(context.Background(), embeddednats.WithNATSServerOptions(&natsserver.Options{
		Host: "127.0.0.1", Port: -1, JetStream: true, StoreDir: t.TempDir(),
	}))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { ns.Close() })
	ns.WaitForServer()
	manager, err := live.NewManager(context.Background(), ns, sessions.NewCookieStore([]byte("room-assignment-create-test")))
	if err != nil {
		t.Fatal(err)
	}

	db, logger := testutil.CreateTestDBAndLogger(t, "room_assignment_create_relation")
	pulje := models.PuljeFredagKveld
	testutil.MustExec(t, db, `INSERT INTO puljer(id,name,status,start_at,end_at) VALUES(?,?,'Open','2026-10-09T18:00:00Z','2026-10-09T23:00:00Z')`, pulje, pulje)
	testutil.MustExec(t, db, `INSERT INTO rooms(id,name,room_number,floor,max_concurrent_games) VALUES(42,'Andre','705',7,0)`)
	insertBoardEvent(t, db, "event-without-pulje", "Spill", "Godkjent", "Default", 0, 0, "host@example.com", "Host")

	router := chi.NewRouter()
	if err := SetupAdminRoute(router, logger, manager, db, nil); err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodPost, "/admin/rooms/api/assignment/FredagKveld/event-without-pulje/42", nil))

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, recorder.Code, recorder.Body.String())
	}
	if got := testutil.QueryInt(t, db, `SELECT COUNT(*) FROM relation_event_puljer WHERE event_id=? AND pulje_id=?`, "event-without-pulje", pulje); got != 1 {
		t.Fatalf("expected a new pulje relation, got %d rows", got)
	}
	if got := testutil.QueryInt(t, db, `SELECT room_id FROM relation_event_puljer WHERE event_id=? AND pulje_id=?`, "event-without-pulje", pulje); got != 42 {
		t.Fatalf("room assignment = %d, want 42", got)
	}
	if got := testutil.QueryInt(t, db, `SELECT is_in_pulje FROM relation_event_puljer WHERE event_id=? AND pulje_id=?`, "event-without-pulje", pulje); got != 1 {
		t.Fatal("new relation should be active")
	}
	if got := testutil.QueryInt(t, db, `SELECT is_published FROM relation_event_puljer WHERE event_id=? AND pulje_id=?`, "event-without-pulje", pulje); got != 0 {
		t.Fatal("new relation should retain the database default for the legacy publication flag")
	}
	if !strings.Contains(recorder.Body.String(), "room-saved") {
		t.Fatal("successful assignment should notify the page")
	}
}
