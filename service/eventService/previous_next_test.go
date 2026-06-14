package eventservice

import (
	"context"
	"database/sql"
	"testing"

	"github.com/Regncon/conorganizer/components"
	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/bdd"
)

func TestGetPreviousNextInnsendtGodkjent_ReturnsNeighborsAmongSubmittedAndApprovedEvents(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given submitted, approved, and draft events ordered by creation time.",
		When:  "When previous/next navigation is requested for an event.",
		Then:  "Then only submitted and approved events are used as navigation neighbors.",
	})

	// Given
	expectedCases := []expectedPreviousNextCase{
		{
			name:      "middle event has both neighbors",
			currentID: "e2",
			expected:  expectedSubmittedPreviousNext{previousURL: "e3", previousTitle: "New", nextURL: "e1", nextTitle: "Old"},
		},
		{
			name:      "newest event has next only",
			currentID: "e3",
			expected:  expectedSubmittedPreviousNext{nextURL: "e2", nextTitle: "Mid"},
		},
		{
			name:      "oldest event has previous only",
			currentID: "e1",
			expected:  expectedSubmittedPreviousNext{previousURL: "e2", previousTitle: "Mid"},
		},
		{
			name:      "draft event is excluded",
			currentID: "e4",
			expected:  expectedSubmittedPreviousNext{},
		},
		{
			name:      "missing event has no neighbors",
			currentID: "does-not-exist",
			expected:  expectedSubmittedPreviousNext{},
		},
	}
	ctx := context.Background()
	imgDir := ""
	db := testutil.CreateTestDB(t, "previous-next")
	seedPreviousNextEvents(t, db)

	for _, tc := range expectedCases {
		t.Run(tc.name, func(t *testing.T) {
			// Given
			expected := tc.expected

			// When
			actual, err := GetPreviousNextInnsendtGodkjent(ctx, db, tc.currentID, &imgDir)

			// Then
			if err != nil {
				t.Fatalf("expected previous/next lookup to succeed: %v", err)
			}
			assertPreviousNext(t, expected, actual)
		})
	}
}

type expectedPreviousNextCase struct {
	name      string
	currentID string
	expected  expectedSubmittedPreviousNext
}

type expectedSubmittedPreviousNext struct {
	previousURL   string
	previousTitle string
	nextURL       string
	nextTitle     string
}

func seedPreviousNextEvents(t testing.TB, db *sql.DB) {
	t.Helper()

	testutil.MustExec(t, db, `DELETE FROM events`)
	testutil.MustExec(t, db, `INSERT OR IGNORE INTO event_statuses (status) VALUES (?), (?), (?)`, models.EventStatusDraft, models.EventStatusSubmitted, models.EventStatusApproved)
	testutil.MustExec(t, db, `INSERT OR IGNORE INTO events_types (event_type) VALUES (?)`, models.EventTypeOther)
	testutil.MustExec(t, db, `INSERT OR IGNORE INTO age_groups (age_group) VALUES (?)`, models.AgeGroupDefault)
	testutil.MustExec(t, db, `INSERT OR IGNORE INTO event_runtimes (runtime) VALUES (?)`, models.RunTimeNormal)
	testutil.MustExec(t, db, `
		INSERT INTO events (
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
		) VALUES
		('e1','Old','intro e1','desc e1',?,?,?,'Host One','one@test.test','11111111',4,1,1,?,'2025-10-01 10:00:00'),
		('e2','Mid','intro e2','desc e2',?,?,?,'Host Two','two@test.test','22222222',5,0,1,?,'2025-10-02 10:00:00'),
		('e3','New','intro e3','desc e3',?,?,?,'Host Tre','tre@test.test','33333333',6,1,0,?,'2025-10-03 10:00:00'),
		('e4','KladdRow','intro e4','desc e4',?,?,?,'Host Four','four@test.test','44444444',3,0,0,?,'2025-10-04 10:00:00')
	`,
		models.EventTypeOther, models.AgeGroupDefault, models.RunTimeNormal, models.EventStatusApproved,
		models.EventTypeOther, models.AgeGroupDefault, models.RunTimeNormal, models.EventStatusSubmitted,
		models.EventTypeOther, models.AgeGroupDefault, models.RunTimeNormal, models.EventStatusApproved,
		models.EventTypeOther, models.AgeGroupDefault, models.RunTimeNormal, models.EventStatusDraft,
	)
}

func assertPreviousNext(t testing.TB, expected expectedSubmittedPreviousNext, actual components.PreviousNext) {
	t.Helper()

	if actual.PreviousUrl != expected.previousURL {
		t.Fatalf("previous URL mismatch\nexpected: %q\nactual:   %q", expected.previousURL, actual.PreviousUrl)
	}
	if actual.PreviousTitle != expected.previousTitle {
		t.Fatalf("previous title mismatch\nexpected: %q\nactual:   %q", expected.previousTitle, actual.PreviousTitle)
	}
	if actual.NextUrl != expected.nextURL {
		t.Fatalf("next URL mismatch\nexpected: %q\nactual:   %q", expected.nextURL, actual.NextUrl)
	}
	if actual.NextTitle != expected.nextTitle {
		t.Fatalf("next title mismatch\nexpected: %q\nactual:   %q", expected.nextTitle, actual.NextTitle)
	}
	if actual.PreviousImageURL != "" {
		t.Fatalf("expected previous image URL to be empty, got %q", actual.PreviousImageURL)
	}
	if actual.NextImageURL != "" {
		t.Fatalf("expected next image URL to be empty, got %q", actual.NextImageURL)
	}
}
