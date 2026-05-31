package event

import (
	"database/sql"
	"strings"
	"testing"
	"time"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/templtest"
)

func TestEventInterestPanel_WhenScheduledWarningHasFired_RendersWarningState(t *testing.T) {
	// Gitt at den planlagte varselmeldingen for en åpen pulje har blitt sendt,
	// når interessepanelet rendres på nytt,
	// så skal billettholderen se varselstatus ved knappen.

	// Given
	expectedHelperVisible := true
	expectedHelperClass := "pulje-interest-state--warning"
	expectedMessagePart := "låses snart"

	now := time.Now()
	puljer := []models.PuljeRow{
		buildEventInterestTestPulje(
			models.PuljeFredagKveld,
			"Fredag kveld",
			models.PuljeStatusOpen,
			now.Add(2*time.Hour),
		),
	}

	// When
	doc := templtest.Render(t, EventInterestPanel(true, puljer))
	helper := doc.Find(".event-interest-helper")
	actualHelperVisible := helper.Length() > 0
	actualMessage := strings.Join(strings.Fields(helper.Text()), " ")
	actualHasExpectedClass := helper.HasClass(expectedHelperClass)

	// Then
	if actualHelperVisible != expectedHelperVisible {
		t.Fatalf("helper visibility mismatch\nexpected: %v\nactual:   %v", expectedHelperVisible, actualHelperVisible)
	}
	if !actualHasExpectedClass {
		t.Fatalf("helper class mismatch\nexpected helper to have class: %s", expectedHelperClass)
	}
	if !strings.Contains(actualMessage, expectedMessagePart) {
		t.Fatalf("helper message mismatch\nexpected to contain: %q\nactual:              %q", expectedMessagePart, actualMessage)
	}
}

func TestEventInterestPanel_WhenCurrentTimeIsBeforeWarningThreshold_RendersNoWarningState(t *testing.T) {
	// Gitt at en åpen pulje ikke nærmer seg låsing,
	// når interessepanelet rendres,
	// så skal ingen låseadvarsel vises ved knappen.

	// Given
	expectedHelperVisible := false

	now := time.Now()
	puljer := []models.PuljeRow{
		buildEventInterestTestPulje(
			models.PuljeFredagKveld,
			"Fredag kveld",
			models.PuljeStatusOpen,
			now.Add(4*time.Hour),
		),
	}

	// When
	doc := templtest.Render(t, EventInterestPanel(true, puljer))
	actualHelperVisible := doc.Find(".event-interest-helper").Length() > 0

	// Then
	if actualHelperVisible != expectedHelperVisible {
		t.Fatalf("helper visibility mismatch\nexpected: %v\nactual:   %v", expectedHelperVisible, actualHelperVisible)
	}
}

func TestEventInterestPanel_WhenScheduledUrgentWarningHasFired_RendersUrgentWarningState(t *testing.T) {
	// Gitt at den planlagte hastevarselmeldingen for en åpen pulje har blitt sendt,
	// når interessepanelet rendres,
	// så skal den mest presserende puljemeldingen vises ved knappen.

	// Given
	expectedHelperVisible := true
	expectedHelperClass := "pulje-interest-state--urgent-warning"
	expectedMessagePart := "låses straks"

	now := time.Now()
	puljer := []models.PuljeRow{
		buildEventInterestTestPulje(
			models.PuljeFredagKveld,
			"Fredag kveld",
			models.PuljeStatusLocked,
			now.Add(-1*time.Hour),
		),
		buildEventInterestTestPulje(
			models.PuljeLordagMorgen,
			"Lørdag morgen",
			models.PuljeStatusOpen,
			now.Add(45*time.Minute),
		),
	}

	// When
	doc := templtest.Render(t, EventInterestPanel(true, puljer))
	helper := doc.Find(".event-interest-helper")
	actualHelperVisible := helper.Length() > 0
	actualMessage := strings.Join(strings.Fields(helper.Text()), " ")
	actualHasExpectedClass := helper.HasClass(expectedHelperClass)

	// Then
	if actualHelperVisible != expectedHelperVisible {
		t.Fatalf("helper visibility mismatch\nexpected: %v\nactual:   %v", expectedHelperVisible, actualHelperVisible)
	}
	if !actualHasExpectedClass {
		t.Fatalf("helper class mismatch\nexpected helper to have class: %s", expectedHelperClass)
	}
	if !strings.Contains(actualMessage, expectedMessagePart) {
		t.Fatalf("helper message mismatch\nexpected to contain: %q\nactual:              %q", expectedMessagePart, actualMessage)
	}
}

func TestUpdateInterest_WhenPuljeIsOpen_UpdatesInterest(t *testing.T) {
	// Gitt at en billettholder har meldt interesse i en åpen pulje,
	// når interessen endres,
	// så skal den nye interessen lagres.

	// Given
	expectedInterest := models.InterestLevelLow

	db := createEventInterestTestDB(t)
	fixture := seedEventInterestUpdateFixture(t, db, models.PuljeStatusOpen, models.InterestLevelHigh)

	// When
	err := updateInterest(
		fixture.userExternalID,
		fixture.billettholderID,
		fixture.eventID,
		expectedInterest,
		string(fixture.puljeID),
		db,
	)
	actualInterest := getEventInterestTestInterest(t, db, fixture.eventID, fixture.billettholderID, fixture.puljeID)

	// Then
	if err != nil {
		t.Fatalf("expected open pulje interest update to succeed: %v", err)
	}
	if actualInterest != expectedInterest {
		t.Fatalf("interest level mismatch\nexpected: %s\nactual:   %s", expectedInterest, actualInterest)
	}
}

func TestUpdateInterest_WhenPuljeIsLocked_RejectsInterestChangeAndKeepsExistingInterest(t *testing.T) {
	// Gitt at en billettholder allerede har meldt interesse i en låst pulje,
	// når interessen forsøkes endret,
	// så skal endringen avvises og eksisterende interesse beholdes.

	// Given
	expectedInterest := models.InterestLevelHigh
	expectedErrorText := "locked"

	db := createEventInterestTestDB(t)
	fixture := seedEventInterestUpdateFixture(t, db, models.PuljeStatusLocked, expectedInterest)

	// When
	err := updateInterest(
		fixture.userExternalID,
		fixture.billettholderID,
		fixture.eventID,
		models.InterestLevelLow,
		string(fixture.puljeID),
		db,
	)
	actualInterest := getEventInterestTestInterest(t, db, fixture.eventID, fixture.billettholderID, fixture.puljeID)

	// Then
	if err == nil {
		t.Errorf("expected locked pulje to reject interest update")
	} else if !strings.Contains(strings.ToLower(err.Error()), expectedErrorText) {
		t.Errorf("error mismatch\nexpected to contain: %q\nactual:              %v", expectedErrorText, err)
	}
	if actualInterest != expectedInterest {
		t.Fatalf("interest level mismatch\nexpected: %s\nactual:   %s", expectedInterest, actualInterest)
	}
}

func TestUpdateInterest_WhenPuljeIsCompleted_RejectsInterestChangeAndKeepsExistingInterest(t *testing.T) {
	// Gitt at en billettholder allerede har meldt interesse i en fullført pulje,
	// når interessen forsøkes endret,
	// så skal endringen avvises og eksisterende interesse beholdes.

	// Given
	expectedInterest := models.InterestLevelHigh
	expectedErrorText := "completed"

	db := createEventInterestTestDB(t)
	fixture := seedEventInterestUpdateFixture(t, db, models.PuljeStatusCompleted, expectedInterest)

	// When
	err := updateInterest(
		fixture.userExternalID,
		fixture.billettholderID,
		fixture.eventID,
		models.InterestLevelLow,
		string(fixture.puljeID),
		db,
	)
	actualInterest := getEventInterestTestInterest(t, db, fixture.eventID, fixture.billettholderID, fixture.puljeID)

	// Then
	if err == nil {
		t.Errorf("expected completed pulje to reject interest update")
	} else if !strings.Contains(strings.ToLower(err.Error()), expectedErrorText) {
		t.Errorf("error mismatch\nexpected to contain: %q\nactual:              %v", expectedErrorText, err)
	}
	if actualInterest != expectedInterest {
		t.Fatalf("interest level mismatch\nexpected: %s\nactual:   %s", expectedInterest, actualInterest)
	}
}

type eventInterestUpdateFixture struct {
	userExternalID  string
	billettholderID int
	eventID         string
	puljeID         models.Pulje
}

func createEventInterestTestDB(t *testing.T) *sql.DB {
	t.Helper()

	db, _, err := testutil.CreateTemporaryDBAndLogger("event_interest", t)
	if err != nil {
		t.Fatalf("failed to create test database: %v", err)
	}
	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Fatalf("failed to close test database: %v", err)
		}
	})

	seedEventInterestLookups(t, db)
	return db
}

func seedEventInterestLookups(t *testing.T, db *sql.DB) {
	t.Helper()

	mustExecEventInterestTest(t, db, `INSERT OR IGNORE INTO event_statuses(status) VALUES (?)`, models.EventStatusApproved)
	mustExecEventInterestTest(t, db, `INSERT OR IGNORE INTO events_types(event_type) VALUES (?)`, models.EventTypeOther)
	mustExecEventInterestTest(t, db, `INSERT OR IGNORE INTO age_groups(age_group) VALUES (?)`, models.AgeGroupDefault)
	mustExecEventInterestTest(t, db, `INSERT OR IGNORE INTO event_runtimes(runtime) VALUES (?)`, models.RunTimeNormal)
	mustExecEventInterestTest(t, db, `INSERT OR IGNORE INTO interest_levels(interest_level) VALUES (?), (?), (?)`, models.InterestLevelHigh, models.InterestLevelMedium, models.InterestLevelLow)
	mustExecEventInterestTest(t, db, `INSERT OR IGNORE INTO pulje_statuses(status) VALUES (?), (?), (?)`, models.PuljeStatusOpen, models.PuljeStatusLocked, models.PuljeStatusCompleted)
}

func seedEventInterestUpdateFixture(
	t *testing.T,
	db *sql.DB,
	puljeStatus models.PuljeStatus,
	existingInterest models.InterestLevel,
) eventInterestUpdateFixture {
	t.Helper()

	fixture := eventInterestUpdateFixture{
		userExternalID:  "event-interest-user",
		billettholderID: 901,
		eventID:         "event-interest-event",
		puljeID:         models.PuljeFredagKveld,
	}

	mustExecEventInterestTest(t, db, `
		INSERT INTO users (id, external_id, email, is_admin)
		VALUES (501, ?, 'event-interest-user@example.com', 0)
	`, fixture.userExternalID)
	mustExecEventInterestTest(t, db, `
		INSERT INTO billettholdere (
			id, first_name, last_name, ticket_type_id, ticket_type, is_over_18, order_id, ticket_id
		) VALUES (?, 'Event', 'Interest', 1, 'Ticket', 1, 7001, 8001)
	`, fixture.billettholderID)
	mustExecEventInterestTest(t, db, `
		INSERT INTO relation_billettholdere_users (billettholder_id, user_id)
		VALUES (?, 501)
	`, fixture.billettholderID)
	mustExecEventInterestTest(t, db, `
		INSERT INTO puljer (id, name, status, start_at, end_at)
		VALUES (?, 'Fredag kveld', ?, '2026-10-09T18:30:00+02:00', '2026-10-09T23:00:00+02:00')
	`, fixture.puljeID, puljeStatus)
	mustExecEventInterestTest(t, db, `
		INSERT INTO events (
			id, title, intro, description, system, event_type,
			age_group, event_runtime, host_name, email, phone_number,
			max_players, beginner_friendly, can_be_run_in_english,
			status
		) VALUES (?, 'Interest Event', 'intro', 'description', '', ?, ?, ?, 'Host', 'host@example.com', '11111111', 4, 1, 1, ?)
	`, fixture.eventID, models.EventTypeOther, models.AgeGroupDefault, models.RunTimeNormal, models.EventStatusApproved)
	mustExecEventInterestTest(t, db, `
		INSERT INTO relation_event_puljer (event_id, pulje_id, is_in_pulje, is_published)
		VALUES (?, ?, 1, 1)
	`, fixture.eventID, fixture.puljeID)
	mustExecEventInterestTest(t, db, `
		INSERT INTO interests (billettholder_id, event_id, pulje_id, interest_level)
		VALUES (?, ?, ?, ?)
	`, fixture.billettholderID, fixture.eventID, fixture.puljeID, existingInterest)

	return fixture
}

func getEventInterestTestInterest(t *testing.T, db *sql.DB, eventID string, billettholderID int, puljeID models.Pulje) models.InterestLevel {
	t.Helper()

	interest, err := getSelectedInterest(eventID, billettholderID, string(puljeID), db)
	if err != nil {
		t.Fatalf("failed to get selected interest: %v", err)
	}
	return interest
}

func buildEventInterestTestPulje(id models.Pulje, name string, status models.PuljeStatus, startAt time.Time) models.PuljeRow {
	return models.PuljeRow{
		ID:      id,
		Name:    name,
		Status:  status,
		StartAt: models.NewDBDateTime(startAt),
		EndAt:   models.NewDBDateTime(startAt.Add(4 * time.Hour)),
	}
}

func mustExecEventInterestTest(t *testing.T, db *sql.DB, query string, args ...any) {
	t.Helper()

	if _, err := db.Exec(query, args...); err != nil {
		t.Fatalf("exec failed: %v\nquery:\n%s", err, query)
	}
}
