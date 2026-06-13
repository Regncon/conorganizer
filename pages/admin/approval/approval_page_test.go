package approval

import (
	"database/sql"
	"slices"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/templtest"
)

func TestSplitEventsByStatusMap_ReturnsSubmittedAndApprovedOnly(t *testing.T) {
	// Gitt arrangementer med kladd, innsendt, godkjent og annonsert status,
	// når godkjenningssiden grupperer arrangementer,
	// så skal bare innsendte og godkjente arrangementer vises i sine seksjoner.

	// Given
	expectedSubmittedIDs := []string{"submitted-event"}
	expectedApprovedIDs := []string{"approved-event"}
	events := []models.EventCardModel{
		{Id: "draft-event", Status: models.EventStatusDraft},
		{Id: "submitted-event", Status: models.EventStatusSubmitted},
		{Id: "approved-event", Status: models.EventStatusApproved},
		{Id: "announced-event", Status: models.EventStatusAnnounced},
	}

	// When
	submitted, approved := splitEventsByStatusMap(events)
	actualSubmittedIDs := approvalEventIDs(submitted)
	actualApprovedIDs := approvalEventIDs(approved)

	// Then
	if !slices.Equal(expectedSubmittedIDs, actualSubmittedIDs) {
		t.Fatalf("submitted events mismatch\nexpected: %v\nactual:   %v", expectedSubmittedIDs, actualSubmittedIDs)
	}
	if !slices.Equal(expectedApprovedIDs, actualApprovedIDs) {
		t.Fatalf("approved events mismatch\nexpected: %v\nactual:   %v", expectedApprovedIDs, actualApprovedIDs)
	}
}

func TestApprovalPageContent_RendersSectionsAndEditLinks(t *testing.T) {
	// Gitt at det finnes innsendte og godkjente arrangementer,
	// når godkjenningssiden rendres,
	// så skal hvert arrangement vises i riktig seksjon med lenke til adminredigering.

	// Given
	expectedSectionTitles := []string{
		"Arrangementer til Godkjenning",
		"Godkjente Arrangementer",
	}
	expectedHrefs := []string{
		"/admin/approval/edit/approved-event",
		"/admin/approval/edit/submitted-event",
	}
	db, logger := testutil.CreateTestDBAndLogger(t, "approval_page")
	seedApprovalPageLookups(t, db)
	insertApprovalPageEvent(t, db, "draft-event", "Draft Event", models.EventStatusDraft)
	insertApprovalPageEvent(t, db, "submitted-event", "Submitted Event", models.EventStatusSubmitted)
	insertApprovalPageEvent(t, db, "approved-event", "Approved Event", models.EventStatusApproved)
	insertApprovalPageEvent(t, db, "announced-event", "Announced Event", models.EventStatusAnnounced)

	// When
	doc := templtest.Render(t, ApprovalPageContent(db, logger))
	actualSectionTitles := templtest.CollectTexts(doc, ".event-approval-section-header")
	actualHrefs := templtest.CollectUniqueHrefs(doc)

	// Then
	if !slices.Equal(expectedSectionTitles, actualSectionTitles) {
		t.Fatalf("approval section titles mismatch\nexpected: %v\nactual:   %v", expectedSectionTitles, actualSectionTitles)
	}
	for _, expectedHref := range expectedHrefs {
		if !slices.Contains(actualHrefs, expectedHref) {
			t.Fatalf("expected approval edit href %q in %v", expectedHref, actualHrefs)
		}
	}
	for _, unexpectedHref := range []string{"/admin/approval/edit/draft-event", "/admin/approval/edit/announced-event"} {
		if slices.Contains(actualHrefs, unexpectedHref) {
			t.Fatalf("did not expect approval edit href %q in %v", unexpectedHref, actualHrefs)
		}
	}
}

func approvalEventIDs(events []models.EventCardModel) []string {
	ids := make([]string, 0, len(events))
	for _, event := range events {
		ids = append(ids, event.Id)
	}
	return ids
}

func seedApprovalPageLookups(t *testing.T, db *sql.DB) {
	t.Helper()

	for _, status := range []models.EventStatus{
		models.EventStatusDraft,
		models.EventStatusSubmitted,
		models.EventStatusApproved,
		models.EventStatusAnnounced,
	} {
		testutil.MustExec(t, db, `INSERT INTO event_statuses(status) VALUES (?) ON CONFLICT(status) DO NOTHING`, status)
	}
	testutil.MustExec(t, db, `INSERT INTO events_types(event_type) VALUES (?) ON CONFLICT(event_type) DO NOTHING`, models.EventTypeOther)
	testutil.MustExec(t, db, `INSERT INTO age_groups(age_group) VALUES (?) ON CONFLICT(age_group) DO NOTHING`, models.AgeGroupDefault)
	testutil.MustExec(t, db, `INSERT INTO event_runtimes(runtime) VALUES (?) ON CONFLICT(runtime) DO NOTHING`, models.RunTimeNormal)
}

func insertApprovalPageEvent(t *testing.T, db *sql.DB, id string, title string, status models.EventStatus) {
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
			email,
			phone_number,
			max_players,
			beginner_friendly,
			can_be_run_in_english,
			status
		)
		VALUES(?, ?, 'Intro', 'Description', 'System', ?, ?, ?, 'Host', 'host@example.com', '12345678', 4, 1, 1, ?)
	`, id, title, models.EventTypeOther, models.AgeGroupDefault, models.RunTimeNormal, status)
}
