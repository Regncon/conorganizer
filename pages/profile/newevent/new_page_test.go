package newevent

import (
	"context"
	"database/sql"
	"strings"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/Regncon/conorganizer/testutil/templtest"
)

func TestValidateEventExists_WhenEventBelongsToUser_Succeeds(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at arrangementet tilhører brukeren.",
		When:  "Når skjema-ruten validerer tilgang.",
		Then:  "Så skal arrangementet kunne åpnes.",
	})

	// Given
	expectedStatusCode := 200
	db := testutil.CreateTestDB(t, "new_event_validate_owner")
	seedNewEventPageLookups(t, db)
	insertNewEventPageEvent(t, db, "owned-event", 501, models.EventStatusDraft)

	// When
	err, actualStatusCode := validateEventExists(db, "owned-event", "501")

	// Then
	if err != nil {
		t.Fatalf("expected owned event validation to succeed: %v", err)
	}
	if actualStatusCode != expectedStatusCode {
		t.Fatalf("status code mismatch\nexpected: %d\nactual:   %d", expectedStatusCode, actualStatusCode)
	}
}

func TestValidateEventExists_WhenEventBelongsToAnotherUser_ReturnsNotFound(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at arrangementet tilhører en annen bruker.",
		When:  "Når skjema-ruten validerer tilgang.",
		Then:  "Så skal vanlig bruker ikke få åpne redigering.",
	})

	// Given
	expectedStatusCode := 404
	db := testutil.CreateTestDB(t, "new_event_validate_foreign")
	seedNewEventPageLookups(t, db)
	insertNewEventPageEvent(t, db, "foreign-event", 501, models.EventStatusDraft)

	// When
	err, actualStatusCode := validateEventExists(db, "foreign-event", "999")

	// Then
	if err == nil {
		t.Fatalf("expected foreign event validation to fail")
	}
	if actualStatusCode != expectedStatusCode {
		t.Fatalf("status code mismatch\nexpected: %d\nactual:   %d", expectedStatusCode, actualStatusCode)
	}
}

func TestNewEventFormPageContent_WhenApprovedEventIsOpenedByNonAdmin_RendersLockedMessage(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at et arrangement allerede er godkjent og brukeren ikke er admin.",
		When:  "Når skjemaet rendres.",
		Then:  "Så skal brukeren se beskjed om at arrangementet ikke kan redigeres videre.",
	})

	// Given
	expectedTextPart := "Arrangementet er allerede godkjent eller publisert"
	db, logger := testutil.CreateTestDBAndLogger(t, "new_event_approved")
	seedNewEventPageLookups(t, db)
	insertNewEventPageEvent(t, db, "approved-event", 501, models.EventStatusApproved)

	// When
	doc := templtest.Render(t, NewEventFormPageContent("approved-event", "501", context.Background(), db, nil, logger))
	actualText := strings.Join(templtest.CollectTexts(doc, "body, h1, h2"), " ")

	// Then
	if !strings.Contains(actualText, expectedTextPart) {
		t.Fatalf("expected approved event message %q\nactual text: %s", expectedTextPart, actualText)
	}
}

func seedNewEventPageLookups(t *testing.T, db *sql.DB) {
	t.Helper()

	for _, status := range []models.EventStatus{
		models.EventStatusDraft,
		models.EventStatusApproved,
	} {
		testutil.MustExec(t, db, `INSERT INTO event_statuses(status) VALUES (?) ON CONFLICT(status) DO NOTHING`, status)
	}
	testutil.MustExec(t, db, `INSERT INTO events_types(event_type) VALUES (?) ON CONFLICT(event_type) DO NOTHING`, models.EventTypeOther)
	testutil.MustExec(t, db, `INSERT INTO age_groups(age_group) VALUES (?) ON CONFLICT(age_group) DO NOTHING`, models.AgeGroupDefault)
	testutil.MustExec(t, db, `INSERT INTO event_runtimes(runtime) VALUES (?) ON CONFLICT(runtime) DO NOTHING`, models.RunTimeNormal)
}

func insertNewEventPageEvent(t *testing.T, db *sql.DB, eventID string, userID int64, status models.EventStatus) {
	t.Helper()

	testutil.MustExec(t, db, `
		INSERT INTO events(
			id,
			title,
			intro,
			description,
			system,
			event_type,
			age_group,
			event_runtime,
			host_name,
			user_id,
			email,
			phone_number,
			max_players,
			beginner_friendly,
			can_be_run_in_english,
			notes,
			status
		)
		VALUES(?, 'New Event', 'Intro', 'Description', 'System', ?, ?, ?, 'Host', ?, 'host@example.com', '12345678', 4, 1, 1, '', ?)
	`, eventID, models.EventTypeOther, models.AgeGroupDefault, models.RunTimeNormal, sql.NullInt64{Int64: userID, Valid: true}, status)
}
