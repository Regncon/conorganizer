package checkIn

import (
	"testing"

	"github.com/Regncon/conorganizer/testutil/bdd"
)

func TestIsOver18_WhenBirthdayIsOnConventionEnd_ReturnsTrue(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given a person who turns eighteen on the last day of Regncon.",
		When:  "When their age is checked.",
		Then:  "Then they count as over eighteen for the convention.",
	})

	// Given
	expectedOver18 := true
	born := "2008-10-04"

	// When
	actualOver18 := isOver18(born)

	// Then
	if actualOver18 != expectedOver18 {
		t.Fatalf("over-18 result mismatch\nexpected: %v\nactual:   %v", expectedOver18, actualOver18)
	}
}

func TestIsOver18_WhenBirthdayIsAfterConventionEnd_ReturnsFalse(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given a person who turns eighteen after the last day of Regncon.",
		When:  "When their age is checked.",
		Then:  "Then they do not count as over eighteen for the convention.",
	})

	// Given
	expectedOver18 := false
	born := "2008-10-05"

	// When
	actualOver18 := isOver18(born)

	// Then
	if actualOver18 != expectedOver18 {
		t.Fatalf("over-18 result mismatch\nexpected: %v\nactual:   %v", expectedOver18, actualOver18)
	}
}
