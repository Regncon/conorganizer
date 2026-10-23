package profilepage

import (
	"context"
	"database/sql"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Regncon/conorganizer/service/authctx"
	"github.com/Regncon/conorganizer/service/live"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/delaneyj/toolbelt/embeddednats"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
	natsserver "github.com/nats-io/nats-server/v2/server"
)

const (
	newEventAccessEventID       = "arr-eier"
	newEventAccessOriginalTitle = "Gammelt navn"
	newEventAccessOwnerExtID    = "ekstern-kari-eier"
	newEventAccessOtherExtID    = "ekstern-per-annen"
	newEventAccessAdminExtID    = "ekstern-ada-admin"
	newEventLiveStreamDuration  = 200 * time.Millisecond
)

func TestNewEventLiveRoute_WhenOwnerOpensIt_StreamsTheForm(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt et arrangement som tilhører den innloggede brukeren.",
		When:  "Når eieren åpner live-strømmen for skjemaet.",
		Then:  "Så skal skjemaet strømmes til eieren.",
	})

	// Given
	expectedStatus := http.StatusOK
	expectedBodyPart := "formsubmission-container"
	_, router := newNewEventAccessFixture(t)
	request := newNewEventLiveRequest(t, authctx.WithUserToken(context.Background(), newEventAccessOwnerExtID, "kari.eier@example.com"))
	recorder := httptest.NewRecorder()

	// When
	router.ServeHTTP(recorder, request)

	// Then
	assertNewEventAccessResponse(t, recorder, expectedStatus, expectedBodyPart)
}

func TestNewEventLiveRoute_WhenAnotherUserOpensIt_RespondsNotFound(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt et arrangement som tilhører en annen bruker.",
		When:  "Når en annen innlogget bruker åpner live-strømmen for skjemaet.",
		Then:  "Så skal svaret være 404 uten skjemaet.",
	})

	// Given
	expectedStatus := http.StatusNotFound
	forbiddenBodyPart := "formsubmission-container"
	_, router := newNewEventAccessFixture(t)
	request := newNewEventLiveRequest(t, authctx.WithUserToken(context.Background(), newEventAccessOtherExtID, "per.annen@example.com"))
	recorder := httptest.NewRecorder()

	// When
	router.ServeHTTP(recorder, request)

	// Then
	if recorder.Code != expectedStatus {
		t.Fatalf("expected status %d, got %d: %s", expectedStatus, recorder.Code, recorder.Body.String())
	}
	if strings.Contains(recorder.Body.String(), forbiddenBodyPart) {
		t.Fatalf("expected no form in the response, got: %s", recorder.Body.String())
	}
}

func TestNewEventLiveRoute_WhenAdminOpensIt_StreamsTheForm(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt et arrangement som tilhører en annen bruker.",
		When:  "Når en admin åpner live-strømmen for skjemaet.",
		Then:  "Så skal skjemaet strømmes til admin.",
	})

	// Given
	expectedStatus := http.StatusOK
	expectedBodyPart := "formsubmission-container"
	_, router := newNewEventAccessFixture(t)
	request := newNewEventLiveRequest(t, authctx.WithAdminUserToken(context.Background(), newEventAccessAdminExtID, "ada.admin@example.com"))
	recorder := httptest.NewRecorder()

	// When
	router.ServeHTTP(recorder, request)

	// Then
	assertNewEventAccessResponse(t, recorder, expectedStatus, expectedBodyPart)
}

func TestNewEventTitleRoute_WhenOwnerUpdatesIt_SavesTheTitle(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt et arrangement som tilhører den innloggede brukeren.",
		When:  "Når eieren endrer tittelen via skjema-API-et.",
		Then:  "Så skal den nye tittelen lagres.",
	})

	// Given
	expectedStatus := http.StatusOK
	expectedTitle := "Nytt navn"
	db, router := newNewEventAccessFixture(t)
	request := newNewEventTitleRequest(t, expectedTitle, authctx.WithUserToken(context.Background(), newEventAccessOwnerExtID, "kari.eier@example.com"))
	recorder := httptest.NewRecorder()

	// When
	router.ServeHTTP(recorder, request)

	// Then
	if recorder.Code != expectedStatus {
		t.Fatalf("expected status %d, got %d: %s", expectedStatus, recorder.Code, recorder.Body.String())
	}
	assertNewEventTitle(t, db, expectedTitle)
}

func TestNewEventTitleRoute_WhenAnotherUserUpdatesIt_RespondsNotFoundAndKeepsTitle(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt et arrangement som tilhører en annen bruker.",
		When:  "Når en annen innlogget bruker prøver å endre tittelen via skjema-API-et.",
		Then:  "Så skal svaret være 404 og tittelen forbli uendret.",
	})

	// Given
	expectedStatus := http.StatusNotFound
	expectedTitle := newEventAccessOriginalTitle
	db, router := newNewEventAccessFixture(t)
	request := newNewEventTitleRequest(t, "Kapret navn", authctx.WithUserToken(context.Background(), newEventAccessOtherExtID, "per.annen@example.com"))
	recorder := httptest.NewRecorder()

	// When
	router.ServeHTTP(recorder, request)

	// Then
	if recorder.Code != expectedStatus {
		t.Fatalf("expected status %d, got %d: %s", expectedStatus, recorder.Code, recorder.Body.String())
	}
	assertNewEventTitle(t, db, expectedTitle)
}

func TestNewEventTitleRoute_WhenAdminUpdatesIt_SavesTheTitle(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt et arrangement som tilhører en annen bruker.",
		When:  "Når en admin endrer tittelen via skjema-API-et.",
		Then:  "Så skal den nye tittelen lagres.",
	})

	// Given
	expectedStatus := http.StatusOK
	expectedTitle := "Navn satt av admin"
	db, router := newNewEventAccessFixture(t)
	request := newNewEventTitleRequest(t, expectedTitle, authctx.WithAdminUserToken(context.Background(), newEventAccessAdminExtID, "ada.admin@example.com"))
	recorder := httptest.NewRecorder()

	// When
	router.ServeHTTP(recorder, request)

	// Then
	if recorder.Code != expectedStatus {
		t.Fatalf("expected status %d, got %d: %s", expectedStatus, recorder.Code, recorder.Body.String())
	}
	assertNewEventTitle(t, db, expectedTitle)
}

func newNewEventAccessFixture(t *testing.T) (*sql.DB, http.Handler) {
	t.Helper()

	ns, err := embeddednats.New(context.Background(), embeddednats.WithNATSServerOptions(&natsserver.Options{
		Host: "127.0.0.1", Port: -1, JetStream: true, StoreDir: t.TempDir(),
	}))
	if err != nil {
		t.Fatalf("failed to start embedded nats: %v", err)
	}
	t.Cleanup(func() { ns.Close() })
	ns.WaitForServer()
	manager, err := live.NewManager(context.Background(), ns, sessions.NewCookieStore([]byte("new-event-access-test")))
	if err != nil {
		t.Fatalf("failed to create live manager: %v", err)
	}

	db, logger := testutil.CreateTestDBAndLogger(t, "new_event_access")
	testutil.MustExec(t, db, `INSERT INTO users (id, external_id, email) VALUES (501, ?, 'kari.eier@example.com'), (502, ?, 'per.annen@example.com'), (503, ?, 'ada.admin@example.com')`,
		newEventAccessOwnerExtID, newEventAccessOtherExtID, newEventAccessAdminExtID)
	testutil.MustExec(t, db, `INSERT INTO events (id, user_id, title, intro, description, host_name, email, phone_number, max_players)
		VALUES (?, 501, ?, '', '', 'Kari Eier', 'kari.eier@example.com', '', 4)`, newEventAccessEventID, newEventAccessOriginalTitle)

	imageDir := t.TempDir()
	router := chi.NewRouter()
	if err := SetupProfileRoute(router, manager, db, &imageDir, logger); err != nil {
		t.Fatalf("failed to set up profile routes: %v", err)
	}
	return db, router
}

func newNewEventLiveRequest(t *testing.T, ctx context.Context) *http.Request {
	t.Helper()

	ctx, cancel := context.WithTimeout(ctx, newEventLiveStreamDuration)
	t.Cleanup(cancel)
	return httptest.NewRequest(http.MethodGet, "/profile/api/new/"+newEventAccessEventID+"/", nil).WithContext(ctx)
}

func newNewEventTitleRequest(t *testing.T, title string, ctx context.Context) *http.Request {
	t.Helper()

	var body io.Reader = strings.NewReader(`{"title":"` + title + `"}`)
	request := httptest.NewRequest(http.MethodPut, "/profile/api/new/"+newEventAccessEventID+"/title", body).WithContext(ctx)
	request.Header.Set("Content-Type", "application/json")
	return request
}

func assertNewEventAccessResponse(t *testing.T, recorder *httptest.ResponseRecorder, expectedStatus int, expectedBodyPart string) {
	t.Helper()

	if recorder.Code != expectedStatus {
		t.Fatalf("expected status %d, got %d: %s", expectedStatus, recorder.Code, recorder.Body.String())
	}
	if !strings.Contains(recorder.Body.String(), expectedBodyPart) {
		t.Fatalf("expected response to contain %q, got: %s", expectedBodyPart, recorder.Body.String())
	}
}

func assertNewEventTitle(t *testing.T, db *sql.DB, expectedTitle string) {
	t.Helper()

	var actualTitle string
	if err := db.QueryRow(`SELECT title FROM events WHERE id = ?`, newEventAccessEventID).Scan(&actualTitle); err != nil {
		t.Fatalf("failed to read event title: %v", err)
	}
	if actualTitle != expectedTitle {
		t.Fatalf("title mismatch\nexpected: %q\nactual:   %q", expectedTitle, actualTitle)
	}
}
