package checkIn

import "testing"

func TestIsOver18_WhenBirthdayIsOnConventionStart_ReturnsTrue(t *testing.T) {
	// Given a person who turns eighteen on the first day of Regncon,
	// when their age is checked,
	// then they count as over eighteen for the convention.

	// Given
	expectedOver18 := true
	born := "2007-10-10"

	// When
	actualOver18 := isOver18(born)

	// Then
	if actualOver18 != expectedOver18 {
		t.Fatalf("over-18 result mismatch\nexpected: %v\nactual:   %v", expectedOver18, actualOver18)
	}
}

func TestIsOver18_WhenBirthdayIsAfterConventionStart_ReturnsFalse(t *testing.T) {
	// Given a person who turns eighteen after the first day of Regncon,
	// when their age is checked,
	// then they do not count as over eighteen for the convention.

	// Given
	expectedOver18 := false
	born := "2007-10-11"

	// When
	actualOver18 := isOver18(born)

	// Then
	if actualOver18 != expectedOver18 {
		t.Fatalf("over-18 result mismatch\nexpected: %v\nactual:   %v", expectedOver18, actualOver18)
	}
}
