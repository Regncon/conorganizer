package formsubmission

import (
	"testing"

	"github.com/Regncon/conorganizer/testutil/bdd"
)

func TestNormalizeTextareaSubmission_NormalizesOuterWhitespaceAndNewlines(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given textarea submissions with surrounding whitespace and mixed newline styles.",
		When:  "When each submission is normalized.",
		Then:  "Then only outer whitespace and Windows newlines are normalized.",
	})

	// Given
	expectedCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "trims leading and trailing whitespace",
			input:    "\n  hello world  \n",
			expected: "hello world",
		},
		{
			name:     "preserves interior newlines",
			input:    "first line\nsecond line\nthird line",
			expected: "first line\nsecond line\nthird line",
		},
		{
			name:     "normalizes windows newlines",
			input:    "\r\nfirst line\r\nsecond line\r\n",
			expected: "first line\nsecond line",
		},
		{
			name:     "whitespace only becomes empty",
			input:    " \n\t\r\n ",
			expected: "",
		},
	}

	for _, tc := range expectedCases {
		t.Run(tc.name, func(t *testing.T) {
			// Given
			expected := tc.expected

			// When
			actual := normalizeTextareaSubmission(tc.input)

			// Then
			if actual != expected {
				t.Fatalf("normalized submission mismatch\nexpected: %q\nactual:   %q", expected, actual)
			}
		})
	}
}
