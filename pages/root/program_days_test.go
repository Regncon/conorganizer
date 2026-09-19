package root

import (
	"testing"
	"time"

	"github.com/Regncon/conorganizer/testutil/bdd"
)

func TestSelectProgramDayIndex_UsesRequestedDateBeforeDefault(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at programdagene er tilgjengelige og en dato er valgt i URL-en.",
		When:  "Når valgt programdag bestemmes.",
		Then:  "Så skal datoen fra URL-en brukes foran standardvalget.",
	})

	// Given
	expectedIndex := 2
	days := testProgramDays()
	requestedDate := "2026-10-04"
	now := time.Date(2026, 10, 1, 12, 0, 0, 0, time.UTC)

	// When
	actualIndex := SelectProgramDayIndex(days, requestedDate, now)

	// Then
	if actualIndex != expectedIndex {
		t.Fatalf("selected day index = %d, want %d", actualIndex, expectedIndex)
	}
}

func TestSelectProgramDayIndex_ClampsDefaultToProgramBoundaries(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at programdagene er tilgjengelige uten en valgt dato i URL-en.",
		When:  "Når valgt programdag bestemmes før, under eller etter programmet.",
		Then:  "Så skal standardvalget bruke nærmeste programdag og holde seg innenfor grensene.",
	})

	// Given
	expectedIndexes := map[string]int{
		"before program": 0,
		"friday":         0,
		"saturday":       1,
		"sunday":         2,
		"after program":  2,
	}
	cases := []struct {
		name string
		now  time.Time
	}{
		{name: "before program", now: time.Date(2026, 10, 1, 12, 0, 0, 0, time.UTC)},
		{name: "friday", now: time.Date(2026, 10, 2, 12, 0, 0, 0, time.UTC)},
		{name: "saturday", now: time.Date(2026, 10, 3, 12, 0, 0, 0, time.UTC)},
		{name: "sunday", now: time.Date(2026, 10, 4, 12, 0, 0, 0, time.UTC)},
		{name: "after program", now: time.Date(2026, 10, 5, 12, 0, 0, 0, time.UTC)},
	}
	days := testProgramDays()

	// When
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Then
			actualIndex := SelectProgramDayIndex(days, "", tc.now)
			if actualIndex != expectedIndexes[tc.name] {
				t.Fatalf("selected day index = %d, want %d", actualIndex, expectedIndexes[tc.name])
			}
		})
	}
}

func testProgramDays() []ProgramDay {
	location := time.UTC
	return []ProgramDay{
		{Date: time.Date(2026, 10, 2, 0, 0, 0, 0, location)},
		{Date: time.Date(2026, 10, 3, 0, 0, 0, 0, location)},
		{Date: time.Date(2026, 10, 4, 0, 0, 0, 0, location)},
	}
}
