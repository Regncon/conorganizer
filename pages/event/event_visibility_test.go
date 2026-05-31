package event

import (
	"database/sql"
	"net/http/httptest"
	"slices"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/service/requestctx"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/templtest"
)

func TestCanViewEvent_AllowsPublicAccessToAnnouncedEvent(t *testing.T) {
	// Gitt at arrangementet er annonsert,
	// når visningstilgang sjekkes,
	// så skal alle kunne se arrangementet.

	// Given
	expectedCanView := true

	db := createEventVisibilityTestDB(t)
	event := &models.Event{Status: models.EventStatusAnnounced}
	userInfo := requestctx.UserRequestInfo{}

	// When
	actualCanView, err := canViewEvent(event, userInfo, db)

	// Then
	if err != nil {
		t.Fatalf("expected visibility check to succeed: %v", err)
	}
	if actualCanView != expectedCanView {
		t.Fatalf("visibility mismatch\nexpected: %v\nactual:   %v", expectedCanView, actualCanView)
	}
}

func TestCanViewEvent_AllowsAdminAccessToUnannouncedEvent(t *testing.T) {
	// Gitt at arrangementet ikke er annonsert og brukeren er admin,
	// når visningstilgang sjekkes,
	// så skal admin kunne se arrangementet.

	// Given
	expectedCanView := true

	db := createEventVisibilityTestDB(t)
	event := &models.Event{Status: models.EventStatusApproved}
	userInfo := requestctx.UserRequestInfo{IsLoggedIn: true, IsAdmin: true}

	// When
	actualCanView, err := canViewEvent(event, userInfo, db)

	// Then
	if err != nil {
		t.Fatalf("expected visibility check to succeed: %v", err)
	}
	if actualCanView != expectedCanView {
		t.Fatalf("visibility mismatch\nexpected: %v\nactual:   %v", expectedCanView, actualCanView)
	}
}

func TestCanViewEvent_AllowsCreatorAccessToUnannouncedEvent(t *testing.T) {
	// Gitt at arrangementet ikke er annonsert og brukeren er eieren,
	// når visningstilgang sjekkes,
	// så skal eieren kunne se arrangementet.

	// Given
	expectedCanView := true

	db := createEventVisibilityTestDB(t)
	seedEventVisibilityUser(t, db, 101, "creator-user")
	event := &models.Event{
		Status: models.EventStatusApproved,
		UserID: sql.NullInt64{
			Int64: 101,
			Valid: true,
		},
	}
	userInfo := requestctx.UserRequestInfo{IsLoggedIn: true, Id: "creator-user"}

	// When
	actualCanView, err := canViewEvent(event, userInfo, db)

	// Then
	if err != nil {
		t.Fatalf("expected visibility check to succeed: %v", err)
	}
	if actualCanView != expectedCanView {
		t.Fatalf("visibility mismatch\nexpected: %v\nactual:   %v", expectedCanView, actualCanView)
	}
}

func TestCanViewEvent_DeniesAnonymousAccessToUnannouncedEvent(t *testing.T) {
	// Gitt at arrangementet ikke er annonsert og brukeren ikke er innlogget,
	// når visningstilgang sjekkes,
	// så skal brukeren ikke kunne se arrangementet.

	// Given
	expectedCanView := false

	db := createEventVisibilityTestDB(t)
	event := &models.Event{Status: models.EventStatusApproved}
	userInfo := requestctx.UserRequestInfo{}

	// When
	actualCanView, err := canViewEvent(event, userInfo, db)

	// Then
	if err != nil {
		t.Fatalf("expected visibility check to succeed: %v", err)
	}
	if actualCanView != expectedCanView {
		t.Fatalf("visibility mismatch\nexpected: %v\nactual:   %v", expectedCanView, actualCanView)
	}
}

func TestEventPageContent_WhenUnannouncedEventIsHidden_RendersFriendlyMessage(t *testing.T) {
	// Gitt at arrangementet ikke er annonsert og brukeren ikke har tilgang,
	// når arrangementssiden vises,
	// så skal brukeren se en vennlig melding om at arrangementet ikke er annonsert ennå.

	// Given
	expectedMessages := []string{"Dette arrangementet er ikke annonsert ennå. Kom tilbake senere, så får du se hva som venter."}

	db := createEventVisibilityTestDB(t)
	logger := testutil.NewSlogAdapter(&testutil.StubLogger{})
	seedEventVisibilityEvent(t, db, "hidden-event", "Hidden Event", models.EventStatusApproved, sql.NullInt64{})
	request := httptest.NewRequest("GET", "/event/hidden-event", nil)

	// When
	doc := templtest.Render(t, event_page_content("hidden-event", false, logger, db, nil, request))
	actualMessages := templtest.CollectTexts(doc, ".event-not-announced-message")

	// Then
	if !slices.Equal(expectedMessages, actualMessages) {
		t.Fatalf("hidden event message mismatch\nexpected: %v\nactual:   %v", expectedMessages, actualMessages)
	}
}

func createEventVisibilityTestDB(t *testing.T) *sql.DB {
	t.Helper()

	db, _, err := testutil.CreateTemporaryDBAndLogger("event_visibility", t)
	if err != nil {
		t.Fatalf("failed to create test database: %v", err)
	}
	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Fatalf("failed to close test database: %v", err)
		}
	})

	mustExecEventVisibilityTest(t, db, `INSERT OR IGNORE INTO event_statuses(status) VALUES (?), (?), (?)`, models.EventStatusApproved, models.EventStatusAnnounced, models.EventStatusDraft)
	mustExecEventVisibilityTest(t, db, `INSERT OR IGNORE INTO events_types(event_type) VALUES (?)`, models.EventTypeOther)
	mustExecEventVisibilityTest(t, db, `INSERT OR IGNORE INTO age_groups(age_group) VALUES (?)`, models.AgeGroupDefault)
	mustExecEventVisibilityTest(t, db, `INSERT OR IGNORE INTO event_runtimes(runtime) VALUES (?)`, models.RunTimeNormal)

	return db
}

func seedEventVisibilityUser(t *testing.T, db *sql.DB, userID int, externalID string) {
	t.Helper()

	mustExecEventVisibilityTest(t, db, `
		INSERT INTO users (id, external_id, email, is_admin)
		VALUES (?, ?, ?, 0)
	`, userID, externalID, externalID+"@example.com")
}

func seedEventVisibilityEvent(t *testing.T, db *sql.DB, eventID string, title string, status models.EventStatus, userID sql.NullInt64) {
	t.Helper()

	mustExecEventVisibilityTest(t, db, `
		INSERT INTO events (
			id, title, intro, description, system, event_type,
			age_group, event_runtime, host_name, user_id, email, phone_number,
			max_players, beginner_friendly, can_be_run_in_english,
			status
		) VALUES (?, ?, 'intro', 'description', '', ?, ?, ?, 'Host', ?, 'host@example.com', '11111111', 4, 1, 1, ?)
	`, eventID, title, models.EventTypeOther, models.AgeGroupDefault, models.RunTimeNormal, userID, status)
}

func mustExecEventVisibilityTest(t *testing.T, db *sql.DB, query string, args ...any) {
	t.Helper()

	if _, err := db.Exec(query, args...); err != nil {
		t.Fatalf("exec failed: %v\nquery:\n%s", err, query)
	}
}
