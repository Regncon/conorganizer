package testfailure

import (
	"testing"

	"github.com/Regncon/conorganizer/testutil/bdd"
)

func TestReporterDemo_OrdinaryTestFailure(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given an intentional reporter demonstration failure.",
		When:  "When the demonstration test runs in the pipeline.",
		Then:  "Then the behavior report displays its assertion output.",
	})

	// Given
	expected := "passing value"
	actual := "intentional failing value"

	// When
	valuesMatch := actual == expected

	// Then
	if !valuesMatch {
		t.Fatalf("intentional test failure for reporter demo: expected %q, got %q", expected, actual)
	}
}
