package feedback

import (
	"database/sql"
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/bdd"
)

func TestSubmit_StoredFeedbackIsListedWithTopics(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given an empty feedback table.",
		When:  "When a user submits feedback about the website with topics.",
		Then:  "Then the feedback is listed with its category, texts, topics and a creation time.",
	})

	// Given
	expected := Submission{
		Category:     CategoryWebsite,
		WentWell:     "Påmeldingen var enkel å finne.",
		CouldImprove: "Programmet er tregt på mobil.",
		Topics:       []Topic{TopicSignup, TopicMobile},
	}
	db, _ := testutil.CreateTestDBAndLogger(t, "feedback_round_trip")

	// When
	err := Submit(db, expected)

	// Then
	if err != nil {
		t.Fatalf("expected submit to succeed: %v", err)
	}
	entries := mustList(t, db, "")
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	entry := entries[0]
	if entry.Category != expected.Category || entry.WentWell != expected.WentWell || entry.CouldImprove != expected.CouldImprove {
		t.Fatalf("expected %+v, got %+v", expected, entry)
	}
	if !reflect.DeepEqual(entry.Topics, expected.Topics) {
		t.Fatalf("expected topics %v, got %v", expected.Topics, entry.Topics)
	}
	if entry.ID == 0 || entry.CreatedAt.IsZero() {
		t.Fatalf("expected id and created_at to be set, got %+v", entry)
	}
}

func TestSubmit_StoresNoTopicsAsEmptyList(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given feedback without any topics.",
		When:  "When the feedback is submitted and listed.",
		Then:  "Then it is stored with an empty topic list.",
	})

	// Given
	expectedTopicsJSON := "[]"
	db, _ := testutil.CreateTestDBAndLogger(t, "feedback_no_topics")

	// When
	err := Submit(db, Submission{Category: CategoryOther, CouldImprove: "Mer kaffe."})

	// Then
	if err != nil {
		t.Fatalf("expected submit to succeed: %v", err)
	}
	var topicsJSON string
	if err := db.QueryRow(`SELECT topics FROM feedback`).Scan(&topicsJSON); err != nil {
		t.Fatalf("expected to read topics: %v", err)
	}
	if topicsJSON != expectedTopicsJSON {
		t.Fatalf("expected topics %q, got %q", expectedTopicsJSON, topicsJSON)
	}
	if entries := mustList(t, db, ""); entries[0].Topics == nil || len(entries[0].Topics) != 0 {
		t.Fatalf("expected an empty non-nil topic list, got %#v", entries[0].Topics)
	}
}

func TestList_ReturnsNewestFirst(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given feedback stored at different times.",
		When:  "When the feedback is listed.",
		Then:  "Then the newest feedback comes first.",
	})

	// Given
	expectedTexts := []string{"nyest", "midten", "eldst"}
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
	assertWentWell(t, entries, expectedTexts)
}

func TestList_FiltersByCategory(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given feedback in several categories.",
		When:  "When the feedback is listed for one category.",
		Then:  "Then only feedback in that category is returned.",
	})

	// Given
	expectedTexts := []string{"om festivalen"}
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
	assertWentWell(t, entries, expectedTexts)
}

func TestSubmit_TrimsSurroundingWhitespace(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given texts with surrounding whitespace.",
		When:  "When the feedback is submitted.",
		Then:  "Then the texts are stored without the surrounding whitespace.",
	})

	// Given
	expectedWentWell := "Flotte lokaler."
	expectedCouldImprove := "Mer kaffe, takk."
	db, _ := testutil.CreateTestDBAndLogger(t, "feedback_trim")

	// When
	err := Submit(db, Submission{
		Category:     CategoryConvention,
		WentWell:     " \tFlotte lokaler.\n",
		CouldImprove: "  \n\tMer kaffe, takk. \n ",
	})

	// Then
	if err != nil {
		t.Fatalf("expected submit to succeed: %v", err)
	}
	entry := mustList(t, db, "")[0]
	if entry.WentWell != expectedWentWell || entry.CouldImprove != expectedCouldImprove {
		t.Fatalf("expected %q/%q, got %q/%q", expectedWentWell, expectedCouldImprove, entry.WentWell, entry.CouldImprove)
	}
}

func TestSubmit_AcceptsOnlyOneText(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given feedback with only one of the two texts filled in.",
		When:  "When the feedback is submitted.",
		Then:  "Then it is accepted.",
	})

	cases := []struct {
		name       string
		submission Submission
	}{
		{name: "only went well", submission: Submission{Category: CategoryOther, WentWell: "Bra stemning."}},
		{name: "only could improve", submission: Submission{Category: CategoryOther, CouldImprove: "Mer skilting."}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Given
			expectedCount := 1
			db, _ := testutil.CreateTestDBAndLogger(t, "feedback_one_text")

			// When
			err := Submit(db, tc.submission)

			// Then
			if err != nil {
				t.Fatalf("expected submit to succeed: %v", err)
			}
			if count := len(mustList(t, db, "")); count != expectedCount {
				t.Fatalf("expected %d entry, got %d", expectedCount, count)
			}
		})
	}
}

func TestValidate_RequiresAtLeastOneText(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given feedback where both texts are empty or whitespace.",
		When:  "When the submission is validated.",
		Then:  "Then it is rejected as required on the went well field.",
	})

	// Given
	expectedErrors := []FieldError{{Field: FieldWentWell, Code: CodeRequired}}
	submission := Submission{Category: CategoryWebsite, WentWell: " \n", CouldImprove: "\t "}

	// When
	errs := submission.Validate()

	// Then
	if !reflect.DeepEqual(errs, expectedErrors) {
		t.Fatalf("expected %v, got %v", expectedErrors, errs)
	}
}

func TestValidate_RejectsUnknownCategory(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given feedback with an unknown or empty category.",
		When:  "When the submission is validated.",
		Then:  "Then the category is reported as invalid.",
	})

	for _, category := range []Category{"", "food"} {
		// Given
		expectedErrors := []FieldError{{Field: FieldCategory, Code: CodeInvalid}}
		submission := Submission{Category: category, WentWell: "Hei"}

		// When
		errs := submission.Validate()

		// Then
		if !reflect.DeepEqual(errs, expectedErrors) {
			t.Fatalf("expected %v for %q, got %v", expectedErrors, category, errs)
		}
	}
}

func TestValidate_LimitsEachTextInRunes(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given texts of multi-byte characters at and over the maximum length.",
		When:  "When the submission is validated.",
		Then:  "Then the limit counts characters per field, and only the too long field is reported.",
	})

	cases := []struct {
		name       string
		submission Submission
		expected   []FieldError
	}{
		{
			name:       "both at max",
			submission: Submission{Category: CategoryOther, WentWell: strings.Repeat("ø", MaxTextLength), CouldImprove: strings.Repeat("å", MaxTextLength)},
			expected:   nil,
		},
		{
			name:       "went well too long",
			submission: Submission{Category: CategoryOther, WentWell: strings.Repeat("ø", MaxTextLength+1)},
			expected:   []FieldError{{Field: FieldWentWell, Code: CodeTooLong}},
		},
		{
			name:       "could improve too long",
			submission: Submission{Category: CategoryOther, WentWell: "Bra", CouldImprove: strings.Repeat("å", MaxTextLength+1)},
			expected:   []FieldError{{Field: FieldCouldImprove, Code: CodeTooLong}},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Given
			expectedErrors := tc.expected

			// When
			errs := tc.submission.Validate()

			// Then
			if !reflect.DeepEqual(errs, expectedErrors) {
				t.Fatalf("expected %v, got %v", expectedErrors, errs)
			}
		})
	}
}

func TestSubmit_StoresTextAtMaxLength(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given texts of exactly the maximum number of multi-byte characters.",
		When:  "When the feedback is submitted.",
		Then:  "Then the database accepts and stores them.",
	})

	// Given
	expectedText := strings.Repeat("ø", MaxTextLength)
	db, _ := testutil.CreateTestDBAndLogger(t, "feedback_max_length")

	// When
	err := Submit(db, Submission{Category: CategoryOther, WentWell: expectedText, CouldImprove: expectedText})

	// Then
	if err != nil {
		t.Fatalf("expected submit to succeed: %v", err)
	}
	assertWentWell(t, mustList(t, db, ""), []string{expectedText})
}

func TestValidate_TopicsMustBelongToCategoryAndBeUnique(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given topics from another category, unknown topics or duplicates.",
		When:  "When the submission is validated.",
		Then:  "Then the topics are reported as invalid.",
	})

	cases := []struct {
		name     string
		category Category
		topics   []Topic
	}{
		{name: "convention topic on website", category: CategoryWebsite, topics: []Topic{TopicFood}},
		{name: "website topic on convention", category: CategoryConvention, topics: []Topic{TopicEvents, TopicSignup}},
		{name: "any topic on other", category: CategoryOther, topics: []Topic{TopicProgram}},
		{name: "unknown topic", category: CategoryWebsite, topics: []Topic{"coffee"}},
		{name: "duplicate topic", category: CategoryWebsite, topics: []Topic{TopicProgram, TopicProgram}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Given
			expectedErrors := []FieldError{{Field: FieldTopics, Code: CodeInvalid}}
			submission := Submission{Category: tc.category, WentWell: "Hei", Topics: tc.topics}

			// When
			errs := submission.Validate()

			// Then
			if !reflect.DeepEqual(errs, expectedErrors) {
				t.Fatalf("expected %v, got %v", expectedErrors, errs)
			}
		})
	}
}

func TestValidate_RejectsPersonalInfoPerField(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given a phone number in one text and an email in the other.",
		When:  "When the submission is validated.",
		Then:  "Then both fields are reported as containing personal info.",
	})

	// Given
	expectedErrors := []FieldError{
		{Field: FieldWentWell, Code: CodePersonalInfo},
		{Field: FieldCouldImprove, Code: CodePersonalInfo},
	}
	submission := Submission{Category: CategoryOther, WentWell: "Ring meg på 123 45 678", CouldImprove: "kari@example.no"}

	// When
	errs := submission.Validate()

	// Then
	if !reflect.DeepEqual(errs, expectedErrors) {
		t.Fatalf("expected %v, got %v", expectedErrors, errs)
	}
}

func TestSubmit_RejectsInvalidFeedbackWithoutStoring(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given feedback that fails validation, including every personal info example.",
		When:  "When the feedback is submitted.",
		Then:  "Then it is rejected with ErrInvalidFeedback and nothing is stored.",
	})

	cases := []submitCase{
		{name: "no text", submission: Submission{Category: CategoryWebsite}},
		{name: "unknown category", submission: Submission{Category: "food", WentWell: "Hei"}},
		{name: "too long", submission: Submission{Category: CategoryWebsite, CouldImprove: strings.Repeat("a", MaxTextLength+1)}},
		{name: "foreign topic", submission: Submission{Category: CategoryWebsite, WentWell: "Hei", Topics: []Topic{TopicFood}}},
	}
	for _, example := range personalInfoExamples {
		cases = append(cases,
			submitCase{name: "went well " + example, submission: Submission{Category: CategoryOther, WentWell: "Hei " + example}},
			submitCase{name: "could improve " + example, submission: Submission{Category: CategoryOther, WentWell: "Hei", CouldImprove: example}},
		)
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Given
			expectedErr := ErrInvalidFeedback
			db, _ := testutil.CreateTestDBAndLogger(t, "feedback_invalid")

			// When
			err := Submit(db, tc.submission)

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

func TestTopicsFor_ListsTopicsPerCategory(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given every feedback category and an unknown one.",
		When:  "When the topics for a category are listed.",
		Then:  "Then each category gets its own labelled topics in display order, and other gets none.",
	})

	// Given
	expectedTopics := map[Category][]Topic{
		CategoryWebsite:    {TopicSignup, TopicProgram, TopicMyPage, TopicMobile},
		CategoryConvention: {TopicEvents, TopicVenue, TopicFood, TopicInformation},
		CategoryOther:      {},
		Category("food"):   {},
	}

	for category, expected := range expectedTopics {
		// When
		topics := TopicsFor(category)

		// Then
		if !reflect.DeepEqual(topics, expected) {
			t.Fatalf("expected topics %v for %q, got %v", expected, category, topics)
		}
		for _, topic := range topics {
			if topic.Label() == "" {
				t.Fatalf("expected topic %q to have a label", topic)
			}
		}
	}
}

func TestTopicsFor_ReturnsACopy(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given the topics listed for the website.",
		When:  "When the caller changes the returned slice.",
		Then:  "Then later listings are unaffected.",
	})

	// Given
	expectedFirst := TopicSignup
	topics := TopicsFor(CategoryWebsite)

	// When
	topics[0] = TopicFood

	// Then
	if first := TopicsFor(CategoryWebsite)[0]; first != expectedFirst {
		t.Fatalf("expected first topic %q, got %q", expectedFirst, first)
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

type submitCase struct {
	name       string
	submission Submission
}

func insertFeedback(t *testing.T, db *sql.DB, category Category, wentWell, createdAt string) {
	t.Helper()
	testutil.MustExec(t, db, `INSERT INTO feedback (category, went_well, created_at) VALUES (?, ?, ?)`, string(category), wentWell, createdAt)
}

func mustList(t *testing.T, db *sql.DB, category Category) []Entry {
	t.Helper()
	entries, err := List(db, category)
	if err != nil {
		t.Fatalf("expected list to succeed: %v", err)
	}
	return entries
}

func assertWentWell(t *testing.T, entries []Entry, expected []string) {
	t.Helper()
	if len(entries) != len(expected) {
		t.Fatalf("expected %d entries, got %d: %+v", len(expected), len(entries), entries)
	}
	for i, entry := range entries {
		if entry.WentWell != expected[i] {
			t.Fatalf("expected entry %d to be %q, got %q", i, expected[i], entry.WentWell)
		}
	}
}
