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

func TestEventPageContent_WhenProgramPublishingIsOff_DoesNotRenderInterestDialog(t *testing.T) {
	// Gitt at programmet ikke er publisert, men arrangementet er publisert i en pulje,
	// når arrangementssiden vises,
	// så skal dialogen ikke rendres og brukeren fortsatt se lenken for å hente billett.

	// Given
	expectedDialogVisible := false
	expectedHrefs := []string{"/", "/profile/tickets"}

	db := createEventVisibilityTestDB(t)
	logger := testutil.NewSlogAdapter(&testutil.StubLogger{})
	seedEventVisibilityEvent(t, db, "program-off-event", "Program Off Event", models.EventStatusAnnounced, sql.NullInt64{})
	seedEventVisibilityPulje(t, db, models.PuljeFredagKveld)
	seedEventVisibilityEventPulje(t, db, "program-off-event", models.PuljeFredagKveld, true)
	setEventVisibilityProgramPublishing(t, db, false)
	request := httptest.NewRequest("GET", "/event/program-off-event?pulje=fredag_kveld", nil)

	// When
	doc := templtest.Render(t, event_page_content("program-off-event", false, logger, db, nil, request))
	actualDialogVisible := templtest.HasSelector(doc, ".interest-dialog")
	actualHrefs := templtest.CollectUniqueHrefs(doc)

	// Then
	if actualDialogVisible != expectedDialogVisible {
		t.Fatalf("interest dialog visibility mismatch\nexpected: %v\nactual:   %v", expectedDialogVisible, actualDialogVisible)
	}
	templtest.AssertSameHrefs(t, expectedHrefs, actualHrefs)
}

func TestEventPageContent_WhenEventIsNotPublishedInPulje_DoesNotRenderInterestDialog(t *testing.T) {
	// Gitt at programmet er publisert, men arrangementet ikke er publisert i puljen,
	// når arrangementssiden vises,
	// så skal dialogen ikke rendres og brukeren fortsatt se lenken for å hente billett.

	// Given
	expectedDialogVisible := false
	expectedHrefs := []string{"/", "/profile/tickets"}

	db := createEventVisibilityTestDB(t)
	logger := testutil.NewSlogAdapter(&testutil.StubLogger{})
	seedEventVisibilityEvent(t, db, "unpublished-pulje-event", "Unpublished Pulje Event", models.EventStatusAnnounced, sql.NullInt64{})
	seedEventVisibilityPulje(t, db, models.PuljeFredagKveld)
	seedEventVisibilityEventPulje(t, db, "unpublished-pulje-event", models.PuljeFredagKveld, false)
	setEventVisibilityProgramPublishing(t, db, true)
	request := httptest.NewRequest("GET", "/event/unpublished-pulje-event?pulje=fredag_kveld", nil)

	// When
	doc := templtest.Render(t, event_page_content("unpublished-pulje-event", false, logger, db, nil, request))
	actualDialogVisible := templtest.HasSelector(doc, ".interest-dialog")
	actualHrefs := templtest.CollectUniqueHrefs(doc)

	// Then
	if actualDialogVisible != expectedDialogVisible {
		t.Fatalf("interest dialog visibility mismatch\nexpected: %v\nactual:   %v", expectedDialogVisible, actualDialogVisible)
	}
	templtest.AssertSameHrefs(t, expectedHrefs, actualHrefs)
}

func TestEventPageContent_WhenProgramAndPuljeArePublished_RendersInterestDialog(t *testing.T) {
	// Gitt at programmet er publisert og arrangementet er publisert i en pulje,
	// når arrangementssiden vises,
	// så skal dialogen for interessevalg rendres.

	// Given
	expectedDialogVisible := true

	db := createEventVisibilityTestDB(t)
	logger := testutil.NewSlogAdapter(&testutil.StubLogger{})
	seedEventVisibilityEvent(t, db, "published-interest-event", "Published Interest Event", models.EventStatusAnnounced, sql.NullInt64{})
	seedEventVisibilityPulje(t, db, models.PuljeFredagKveld)
	seedEventVisibilityEventPulje(t, db, "published-interest-event", models.PuljeFredagKveld, true)
	setEventVisibilityProgramPublishing(t, db, true)
	request := httptest.NewRequest("GET", "/event/published-interest-event?pulje=fredag_kveld", nil)

	// When
	doc := templtest.Render(t, event_page_content("published-interest-event", false, logger, db, nil, request))
	actualDialogVisible := templtest.HasSelector(doc, ".interest-dialog")

	// Then
	if actualDialogVisible != expectedDialogVisible {
		t.Fatalf("interest dialog visibility mismatch\nexpected: %v\nactual:   %v", expectedDialogVisible, actualDialogVisible)
	}
}

func createEventVisibilityTestDB(t *testing.T) *sql.DB {
	t.Helper()

	db := testutil.CreateTestDB(t, "event_visibility")
	mustExecEventVisibilityTest(t, db, `INSERT OR IGNORE INTO event_statuses(status) VALUES (?), (?), (?)`, models.EventStatusApproved, models.EventStatusAnnounced, models.EventStatusDraft)
	mustExecEventVisibilityTest(t, db, `INSERT OR IGNORE INTO events_types(event_type) VALUES (?)`, models.EventTypeOther)
	mustExecEventVisibilityTest(t, db, `INSERT OR IGNORE INTO age_groups(age_group) VALUES (?)`, models.AgeGroupDefault)
	mustExecEventVisibilityTest(t, db, `INSERT OR IGNORE INTO event_runtimes(runtime) VALUES (?)`, models.RunTimeNormal)
	mustExecEventVisibilityTest(t, db, `INSERT OR IGNORE INTO pulje_statuses(status) VALUES (?)`, models.PuljeStatusOpen)
	mustExecEventVisibilityTest(t, db, `INSERT OR IGNORE INTO interest_levels(interest_level) VALUES (?), (?), (?)`, models.InterestLevelHigh, models.InterestLevelMedium, models.InterestLevelLow)
	setEventVisibilityProgramPublishing(t, db, true)

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

func seedEventVisibilityPulje(t *testing.T, db *sql.DB, puljeID models.Pulje) {
	t.Helper()

	mustExecEventVisibilityTest(t, db, `
		INSERT INTO puljer (id, name, status, start_at, end_at)
		VALUES (?, 'Fredag kveld', ?, '2026-10-09T18:30:00+02:00', '2026-10-09T23:00:00+02:00')
	`, puljeID, models.PuljeStatusOpen)
}

func seedEventVisibilityEventPulje(t *testing.T, db *sql.DB, eventID string, puljeID models.Pulje, published bool) {
	t.Helper()

	isPublished := 0
	if published {
		isPublished = 1
	}

	mustExecEventVisibilityTest(t, db, `
		INSERT INTO relation_event_puljer (event_id, pulje_id, is_in_pulje, is_published)
		VALUES (?, ?, 1, ?)
	`, eventID, puljeID, isPublished)
}

func setEventVisibilityProgramPublishing(t *testing.T, db *sql.DB, programPublished bool) {
	t.Helper()

	isPublished := 0
	if programPublished {
		isPublished = 1
	}

	mustExecEventVisibilityTest(t, db, `
		INSERT INTO program_publishing_state(id, is_published)
		VALUES (1, ?)
		ON CONFLICT(id) DO UPDATE SET is_published = excluded.is_published
	`, isPublished)
}

func mustExecEventVisibilityTest(t *testing.T, db *sql.DB, query string, args ...any) {
	t.Helper()

	if _, err := db.Exec(query, args...); err != nil {
		t.Fatalf("exec failed: %v\nquery:\n%s", err, query)
	}
}
