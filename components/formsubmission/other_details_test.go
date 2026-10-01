package formsubmission

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/service/authctx"
	"github.com/Regncon/conorganizer/service/live"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/Regncon/conorganizer/testutil/templtest"
	"github.com/delaneyj/toolbelt/embeddednats"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
	natsserver "github.com/nats-io/nats-server/v2/server"
)

func TestOtherDetails_RegularUserMaxPlayersMinimumIsFour(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at en vanlig bruker fyller ut arrangementsskjemaet.",
		When:  "Når feltet for maks antall spillere rendres.",
		Then:  "Så er laveste tillatte verdi 4.",
	})

	// Given
	expectedMin := "4"

	// When
	doc := templtest.Render(t, otherDetails("event-1", models.AgeGroupDefault, models.RunTimeNormal, true, false, 4, "", false))

	// Then
	assertMaxPlayersMin(t, doc.Find(`input[name="max-players"]`).AttrOr("min", ""), expectedMin)
}

func TestOtherDetails_AdminMaxPlayersMinimumIsZero(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at en admin redigerer et arrangement.",
		When:  "Når feltet for maks antall spillere rendres.",
		Then:  "Så er laveste tillatte verdi 0.",
	})

	// Given
	expectedMin := "0"

	// When
	doc := templtest.Render(t, otherDetails("event-1", models.AgeGroupDefault, models.RunTimeNormal, true, false, 4, "", true))

	// Then
	assertMaxPlayersMin(t, doc.Find(`input[name="max-players"]`).AttrOr("min", ""), expectedMin)
}

func TestUpdateMaxPlayers_RegularUserBelowFourIsRejected(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt et arrangement med 4 spillere og en vanlig bruker.",
		When:  "Når brukeren prøver å sette maks antall spillere til 3.",
		Then:  "Så blir endringen avvist og antallet forblir 4.",
	})

	// Given
	expectedStatus := http.StatusBadRequest
	expectedMaxPlayers := 4
	db := createMaxPlayersTestDB(t, "update_max_players_user_rejected")
	ctx := authctx.WithUserToken(context.Background(), "ext-42", "host@x.no")

	// When
	recorder := putMaxPlayers(t, db, nil, ctx, 3)

	// Then
	if recorder.Code != expectedStatus {
		t.Fatalf("HTTP status mismatch\nexpected: %d\nactual:   %d", expectedStatus, recorder.Code)
	}
	assertStoredMaxPlayers(t, db, expectedMaxPlayers)
}

func TestUpdateMaxPlayers_AdminCanSetZero(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt et arrangement med 4 spillere og en admin.",
		When:  "Når admin setter maks antall spillere til 0.",
		Then:  "Så blir 0 lagret på arrangementet.",
	})

	// Given
	expectedStatus := http.StatusOK
	expectedMaxPlayers := 0
	db := createMaxPlayersTestDB(t, "update_max_players_admin_zero")
	ctx := authctx.WithAdminUserToken(context.Background(), "ext-42", "admin@x.no")

	// When
	recorder := putMaxPlayers(t, db, newMaxPlayersTestLiveManager(t), ctx, expectedMaxPlayers)

	// Then
	if recorder.Code != expectedStatus {
		t.Fatalf("HTTP status mismatch\nexpected: %d\nactual:   %d\nbody: %s", expectedStatus, recorder.Code, recorder.Body.String())
	}
	assertStoredMaxPlayers(t, db, expectedMaxPlayers)
}

func TestUpdateMaxPlayers_AdminBelowZeroIsRejected(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt et arrangement med 4 spillere og en admin.",
		When:  "Når admin prøver å sette maks antall spillere til -1.",
		Then:  "Så blir endringen avvist og antallet forblir 4.",
	})

	// Given
	expectedStatus := http.StatusBadRequest
	expectedMaxPlayers := 4
	db := createMaxPlayersTestDB(t, "update_max_players_admin_negative")
	ctx := authctx.WithAdminUserToken(context.Background(), "ext-42", "admin@x.no")

	// When
	recorder := putMaxPlayers(t, db, nil, ctx, -1)

	// Then
	if recorder.Code != expectedStatus {
		t.Fatalf("HTTP status mismatch\nexpected: %d\nactual:   %d", expectedStatus, recorder.Code)
	}
	assertStoredMaxPlayers(t, db, expectedMaxPlayers)
}

func createMaxPlayersTestDB(t *testing.T, name string) *sql.DB {
	t.Helper()

	db := testutil.CreateTestDB(t, name)
	testutil.MustExec(t, db, `INSERT INTO users (id, external_id, email) VALUES (42, 'ext-42', 'user@x.no')`)
	testutil.MustExec(t, db,
		`INSERT INTO events (id, title, intro, description, host_name, email, phone_number, max_players)
		 VALUES ('e1', 'Spill', '', '', 'Ola', 'ola@x.no', '', 4)`)
	return db
}

func newMaxPlayersTestLiveManager(t *testing.T) *live.Manager {
	t.Helper()

	ns, err := embeddednats.New(context.Background(), embeddednats.WithNATSServerOptions(&natsserver.Options{
		Host: "127.0.0.1", Port: -1, JetStream: true, StoreDir: t.TempDir(),
	}))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { ns.Close() })
	ns.WaitForServer()
	manager, err := live.NewManager(context.Background(), ns, sessions.NewCookieStore([]byte("max-players-test")))
	if err != nil {
		t.Fatal(err)
	}
	return manager
}

func putMaxPlayers(t *testing.T, db *sql.DB, liveManager *live.Manager, ctx context.Context, maxPlayers int) *httptest.ResponseRecorder {
	t.Helper()

	router := chi.NewRouter()
	router.Route("/profile/api/new/{id}/max-players", func(maxPlayersRouter chi.Router) {
		UpdateMaxPlayers(maxPlayersRouter, db, liveManager, testutil.NewTestLogger())
	})
	body := strings.NewReader(`{"maxPlayers":` + strconv.Itoa(maxPlayers) + `}`)
	request := httptest.NewRequest(http.MethodPut, "/profile/api/new/e1/max-players", body).WithContext(ctx)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	return recorder
}

func assertMaxPlayersMin(t *testing.T, actualMin string, expectedMin string) {
	t.Helper()

	if actualMin != expectedMin {
		t.Fatalf("max players min mismatch\nexpected: %q\nactual:   %q", expectedMin, actualMin)
	}
}

func assertStoredMaxPlayers(t *testing.T, db *sql.DB, expectedMaxPlayers int) {
	t.Helper()

	actualMaxPlayers := testutil.QueryInt(t, db, `SELECT max_players FROM events WHERE id = 'e1'`)
	if actualMaxPlayers != expectedMaxPlayers {
		t.Fatalf("stored max_players mismatch\nexpected: %d\nactual:   %d", expectedMaxPlayers, actualMaxPlayers)
	}
}
