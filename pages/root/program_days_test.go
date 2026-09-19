package root

import (
	"testing"
	"time"
)

func TestSelectProgramDayIndex_UsesRequestedDateBeforeDefault(t *testing.T) {
	days := testProgramDays()

	if got := SelectProgramDayIndex(days, "2026-10-04", time.Date(2026, 10, 1, 12, 0, 0, 0, time.UTC)); got != 2 {
		t.Fatalf("selected day index = %d, want %d", got, 2)
	}
}

func TestSelectProgramDayIndex_ClampsDefaultToProgramBoundaries(t *testing.T) {
	days := testProgramDays()
	cases := []struct {
		name string
		now  time.Time
		want int
	}{
		{name: "before program", now: time.Date(2026, 10, 1, 12, 0, 0, 0, time.UTC), want: 0},
		{name: "friday", now: time.Date(2026, 10, 2, 12, 0, 0, 0, time.UTC), want: 0},
		{name: "saturday", now: time.Date(2026, 10, 3, 12, 0, 0, 0, time.UTC), want: 1},
		{name: "sunday", now: time.Date(2026, 10, 4, 12, 0, 0, 0, time.UTC), want: 2},
		{name: "after program", now: time.Date(2026, 10, 5, 12, 0, 0, 0, time.UTC), want: 2},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := SelectProgramDayIndex(days, "", tc.now); got != tc.want {
				t.Fatalf("selected day index = %d, want %d", got, tc.want)
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
