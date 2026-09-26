package feedback

import (
	"database/sql"
	"errors"
	"strings"
	"testing"

	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/bdd"
)

func TestSubmit_StoredFeedbackIsListed(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given an empty feedback table.",
		When:  "When a user submits feedback about the website.",
		Then:  "Then the feedback is listed with its category, message and a creation time.",
	})

	// Given
	expectedCategory := CategoryWebsite
	expectedMessage := "Påmeldingen var enkel å finne."
	db, _ := testutil.CreateTestDBAndLogger(t, "feedback_round_trip")

	// When
	err := Submit(db, expectedCategory, expectedMessage)

	// Then
	if err != nil {
		t.Fatalf("expected submit to succeed: %v", err)
	}
	entries := mustList(t, db, "")
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	entry := entries[0]
	if entry.Category != expectedCategory || entry.Message != expectedMessage {
		t.Fatalf("expected %q/%q, got %q/%q", expectedCategory, expectedMessage, entry.Category, entry.Message)
	}
	if entry.ID == 0 || entry.CreatedAt.IsZero() {
		t.Fatalf("expected id and created_at to be set, got %+v", entry)
	}
}

func TestList_ReturnsNewestFirst(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given feedback stored at different times.",
		When:  "When the feedback is listed.",
		Then:  "Then the newest feedback comes first.",
	})

	// Given
	expectedMessages := []string{"nyest", "midten", "eldst"}
	db, _ := testutil.CreateTestDBAndLogger(t, "feedback_newest_first")
	insertFeedback(t, db, CategoryOther, "midten", "2026-09-20T12:00:00.000Z")
	insertFeedback(t, db, CategoryOther, "nyest", "2026-09-21T12:00:00.000Z")
	insertFeedback(t, db, CategoryOther, "eldst", "2026-09-19T12:00:00.000Z")

	// When
	entries, err := List(db, "")

	// Then
	if err != nil {
		t.Fatalf("expected list to succeed: %v", err)
	}
	assertMessages(t, entries, expectedMessages)
}

func TestList_FiltersByCategory(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given feedback in several categories.",
		When:  "When the feedback is listed for one category.",
		Then:  "Then only feedback in that category is returned.",
	})

	// Given
	expectedMessages := []string{"om festivalen"}
	db, _ := testutil.CreateTestDBAndLogger(t, "feedback_filter")
	insertFeedback(t, db, CategoryWebsite, "om nettsiden", "2026-09-20T12:00:00.000Z")
	insertFeedback(t, db, CategoryConvention, "om festivalen", "2026-09-20T13:00:00.000Z")
	insertFeedback(t, db, CategoryOther, "om noe annet", "2026-09-20T14:00:00.000Z")

	// When
	entries, err := List(db, CategoryConvention)

	// Then
	if err != nil {
		t.Fatalf("expected list to succeed: %v", err)
	}
	assertMessages(t, entries, expectedMessages)
}

func TestSubmit_TrimsSurroundingWhitespace(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given a message with surrounding whitespace.",
		When:  "When the feedback is submitted.",
		Then:  "Then the message is stored without the surrounding whitespace.",
	})

	// Given
	expectedMessages := []string{"Mer kaffe, takk."}
	db, _ := testutil.CreateTestDBAndLogger(t, "feedback_trim")

	// When
	err := Submit(db, CategoryConvention, "  \n\tMer kaffe, takk. \n ")

	// Then
	if err != nil {
		t.Fatalf("expected submit to succeed: %v", err)
	}
	assertMessages(t, mustList(t, db, ""), expectedMessages)
}

func TestSubmit_AcceptsMessageAtMaxLength(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given a message of exactly the maximum number of multi-byte characters.",
		When:  "When the feedback is submitted.",
		Then:  "Then it is accepted, since the limit counts characters and not bytes.",
	})

	// Given
	expectedMessage := strings.Repeat("ø", MaxMessageLength)
	db, _ := testutil.CreateTestDBAndLogger(t, "feedback_max_length")

	// When
	err := Submit(db, CategoryOther, expectedMessage)

	// Then
	if err != nil {
		t.Fatalf("expected submit to succeed: %v", err)
	}
	assertMessages(t, mustList(t, db, ""), []string{expectedMessage})
}

func TestSubmit_RejectsInvalidFeedback(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given feedback with an empty, too long or uncategorized message.",
		When:  "When the feedback is submitted.",
		Then:  "Then it is rejected with ErrInvalidFeedback and nothing is stored.",
	})

	cases := []struct {
		name     string
		category Category
		message  string
	}{
		{name: "empty message", category: CategoryWebsite, message: ""},
		{name: "whitespace only message", category: CategoryWebsite, message: " \n\t "},
		{name: "too long message", category: CategoryWebsite, message: strings.Repeat("a", MaxMessageLength+1)},
		{name: "unknown category", category: Category("food"), message: "Hei"},
		{name: "empty category", category: "", message: "Hei"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Given
			expectedErr := ErrInvalidFeedback
			db, _ := testutil.CreateTestDBAndLogger(t, "feedback_invalid")

			// When
			err := Submit(db, tc.category, tc.message)

			// Then
			if !errors.Is(err, expectedErr) {
				t.Fatalf("expected ErrInvalidFeedback, got %v", err)
			}
			if count := testutil.QueryInt(t, db, `SELECT COUNT(*) FROM feedback`); count != 0 {
				t.Fatalf("expected no stored feedback, got %d", count)
			}
		})
	}
}

func TestCategory_SlugRoundTrip(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given every feedback category.",
		When:  "When a category is turned into its URL slug and back.",
		Then:  "Then the Norwegian slug maps back to the same category.",
	})

	// Given
	expectedSlugs := map[Category]string{
		CategoryWebsite:    "nettsiden",
		CategoryConvention: "festivalen",
		CategoryOther:      "annet",
	}

	for _, category := range Categories {
		// When
		slug := category.Slug()
		parsed, ok := CategoryFromSlug(slug)

		// Then
		if slug != expectedSlugs[category] {
			t.Fatalf("expected slug %q for %q, got %q", expectedSlugs[category], category, slug)
		}
		if !ok || parsed != category {
			t.Fatalf("expected slug %q to map back to %q, got %q (ok=%v)", slug, category, parsed, ok)
		}
		if !category.Valid() || category.Label() == "" {
			t.Fatalf("expected %q to be valid with a label", category)
		}
	}
}

func TestCategoryFromSlug_RejectsUnknownSlug(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given slugs that are not Norwegian category words.",
		When:  "When they are mapped to a category.",
		Then:  "Then no category is found.",
	})

	// Given
	unknownSlugs := []string{"", "website", "Nettsiden", "mat"}

	for _, slug := range unknownSlugs {
		// When
		category, ok := CategoryFromSlug(slug)

		// Then
		if ok || category != "" {
			t.Fatalf("expected slug %q to be unknown, got %q (ok=%v)", slug, category, ok)
		}
	}
}

func insertFeedback(t *testing.T, db *sql.DB, category Category, message, createdAt string) {
	t.Helper()
	testutil.MustExec(t, db, `INSERT INTO feedback (category, message, created_at) VALUES (?, ?, ?)`, string(category), message, createdAt)
}

func mustList(t *testing.T, db *sql.DB, category Category) []Entry {
	t.Helper()
	entries, err := List(db, category)
	if err != nil {
		t.Fatalf("expected list to succeed: %v", err)
	}
	return entries
}

func assertMessages(t *testing.T, entries []Entry, expected []string) {
	t.Helper()
	if len(entries) != len(expected) {
		t.Fatalf("expected %d entries, got %d: %+v", len(expected), len(entries), entries)
	}
	for i, entry := range entries {
		if entry.Message != expected[i] {
			t.Fatalf("expected entry %d to be %q, got %q", i, expected[i], entry.Message)
		}
	}
}
