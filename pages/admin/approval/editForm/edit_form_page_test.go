package edit_form

import (
	"context"
	"database/sql"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/templtest"
)

func TestEditEventFormPageContent_DoesNotRenderPreviousNextNavigation(t *testing.T) {
	// Gitt en arrangementside i admin-godkjenning,
	// når siden rendres,
	// så rendres ikke forrige/neste-komponenter.

	// Given
	db := createEditFormNavigationTestDB(t)
	logger := testutil.NewSlogAdapter(&testutil.StubLogger{})
	seedEditFormNavigationLookups(t, db)
	seedEditFormNavigationEvent(t, db, "newer-approved", "Newer Approved", models.EventStatusApproved, "2026-01-03T10:00:00Z")
	seedEditFormNavigationEvent(t, db, "submitted-event", "Submitted Event", models.EventStatusSubmitted, "2026-01-02T10:00:00Z")
	seedEditFormNavigationEvent(t, db, "older-approved", "Older Approved", models.EventStatusApproved, "2026-01-01T10:00:00Z")

	// When
	doc := templtest.Render(t, EditEventFormPageContent(context.Background(), "submitted-event", db, nil, logger))

	// Then
	if templtest.HasSelector(doc, ".breadcrumb-prev-next-buttons") {
		t.Fatalf("expected breadcrumb previous/next navigation not to render")
	}
	if templtest.HasSelector(doc, ".previous-next-buttons") {
		t.Fatalf("expected event previous/next navigation not to render")
	}
}

func createEditFormNavigationTestDB(t *testing.T) *sql.DB {
	t.Helper()

	return testutil.CreateTestDB(t, "edit-form-navigation")
}

func seedEditFormNavigationLookups(t *testing.T, db *sql.DB) {
	t.Helper()

	for _, status := range []models.EventStatus{
		models.EventStatusDraft,
		models.EventStatusSubmitted,
		models.EventStatusApproved,
		models.EventStatusArchived,
		models.EventStatusAnnounced,
	} {
		mustExecEditFormNavigation(t, db, `INSERT INTO event_statuses(status) VALUES (?) ON CONFLICT(status) DO NOTHING`, status)
	}

	mustExecEditFormNavigation(t, db, `INSERT INTO events_types(event_type) VALUES (?) ON CONFLICT(event_type) DO NOTHING`, models.EventTypeOther)
	mustExecEditFormNavigation(t, db, `INSERT INTO age_groups(age_group) VALUES (?) ON CONFLICT(age_group) DO NOTHING`, models.AgeGroupDefault)
	mustExecEditFormNavigation(t, db, `INSERT INTO event_runtimes(runtime) VALUES (?) ON CONFLICT(runtime) DO NOTHING`, models.RunTimeNormal)
	mustExecEditFormNavigation(t, db, `INSERT INTO interest_levels(interest_level) VALUES (?), (?), (?)`, models.InterestLevelHigh, models.InterestLevelMedium, models.InterestLevelLow)
}

func seedEditFormNavigationEvent(t *testing.T, db *sql.DB, id string, title string, status models.EventStatus, createdAt string) {
	t.Helper()

	mustExecEditFormNavigation(t, db, `
		INSERT INTO events(
			id,
			title,
			intro,
			description,
			event_type,
			age_group,
			event_runtime,
			host_name,
			email,
			phone_number,
			max_players,
			beginner_friendly,
			can_be_run_in_english,
			status,
			created_at
		)
		VALUES(?, ?, 'Intro', 'Description', ?, ?, ?, 'Host', 'host@example.com', '12345678', 4, 1, 1, ?, ?)
	`, id, title, models.EventTypeOther, models.AgeGroupDefault, models.RunTimeNormal, status, createdAt)
}

func mustExecEditFormNavigation(t *testing.T, db *sql.DB, query string, args ...any) {
	t.Helper()

	if _, err := db.Exec(query, args...); err != nil {
		t.Fatalf("exec failed: %v\nquery:\n%s", err, query)
	}
}
