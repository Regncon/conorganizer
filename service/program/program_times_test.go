package program

import (
	"slices"
	"testing"
	"time"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil/bdd"
)

func TestBuildDays_AllDayProgramEventInBothPuljerShowsItsTimeOnce(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt et programarrangement som går hele lørdagen og derfor er i begge lørdagspuljene med tidspunktet «Hele dagen».",
		When:  "Når programdagene bygges.",
		Then:  "Så viser programkortet «Hele dagen» én gang.",
	})

	// Given
	expectedTimes := []string{"Hele dagen"}
	events := map[models.Pulje][]puljeEvent{
		models.PuljeLordagMorgen: {programEventWithTime("info-desk", "Hele dagen")},
		models.PuljeLordagKveld:  {programEventWithTime("info-desk", "Hele dagen")},
	}

	// When
	days := buildDays(saturdayPuljer(), events, osloLocation(t))

	// Then
	assertProgramCardTimes(t, days, expectedTimes)
}

func TestBuildDays_ProgramEventShowsEachPuljeTimeOfTheDayInPuljeOrder(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt et programarrangement med ulikt tidspunkt i de to lørdagspuljene.",
		When:  "Når programdagene bygges.",
		Then:  "Så viser programkortet hvert tidspunkt på egen linje i puljerekkefølge.",
	})

	// Given
	expectedTimes := []string{"10:00–12:00", "19:00–21:00"}
	events := map[models.Pulje][]puljeEvent{
		models.PuljeLordagMorgen: {programEventWithTime("workshop", "10:00–12:00")},
		models.PuljeLordagKveld:  {programEventWithTime("workshop", "19:00–21:00")},
	}

	// When
	days := buildDays(saturdayPuljer(), events, osloLocation(t))

	// Then
	assertProgramCardTimes(t, days, expectedTimes)
}

func TestBuildDays_ProgramEventWithoutTimeShowsNoTime(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt et programarrangement der ingen tidspunkt er skrevet inn.",
		When:  "Når programdagene bygges.",
		Then:  "Så viser programkortet ingen tid.",
	})

	// Given
	var expectedTimes []string
	events := map[models.Pulje][]puljeEvent{
		models.PuljeLordagMorgen: {programEventWithTime("cosplay", "")},
		models.PuljeLordagKveld:  {programEventWithTime("cosplay", "")},
	}

	// When
	days := buildDays(saturdayPuljer(), events, osloLocation(t))

	// Then
	assertProgramCardTimes(t, days, expectedTimes)
}

func TestEventSchedule_ProgramEventShowsEachDayWithItsTimes(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt et programarrangement hele lørdagen, og søndag uten tidspunkt.",
		When:  "Når tidsplanen for programarrangementet lages.",
		Then:  "Så står lørdag med «Hele dagen» én gang, og søndag bare med dagen.",
	})

	// Given
	expectedLines := []string{"Lørdag 3.10 · Hele dagen", "Søndag 4.10"}
	saturday := saturdayPuljer()
	puljeTimes := []PuljeTime{
		{Pulje: saturday[0], ProgramTime: "Hele dagen"},
		{Pulje: saturday[1], ProgramTime: "Hele dagen"},
		{Pulje: models.PuljeRow{ID: models.PuljeSondagMorgen, StartAt: models.NewDBDateTime(time.Date(2026, 10, 4, 8, 0, 0, 0, time.UTC))}},
	}

	// When
	lines, err := EventSchedule(puljeTimes, true)

	// Then
	if err != nil {
		t.Fatalf("expected the schedule to be built: %v", err)
	}
	if !slices.Equal(lines, expectedLines) {
		t.Fatalf("schedule = %q, want %q", lines, expectedLines)
	}
}

func saturdayPuljer() []models.PuljeRow {
	return []models.PuljeRow{
		{ID: models.PuljeLordagMorgen, StartAt: models.NewDBDateTime(time.Date(2026, 10, 3, 8, 0, 0, 0, time.UTC))},
		{ID: models.PuljeLordagKveld, StartAt: models.NewDBDateTime(time.Date(2026, 10, 3, 16, 0, 0, 0, time.UTC))},
	}
}

func programEventWithTime(id string, programTime string) puljeEvent {
	return puljeEvent{Event: models.EventCardModel{Id: id}, ProgramTime: programTime}
}

func osloLocation(t *testing.T) *time.Location {
	t.Helper()
	location, err := Location()
	if err != nil {
		t.Fatal(err)
	}
	return location
}

// assertProgramCardTimes checks the only program event on the only day.
func assertProgramCardTimes(t *testing.T, days []Day, expectedTimes []string) {
	t.Helper()
	if len(days) != 1 || len(days[0].ProgramEvents) != 1 {
		t.Fatalf("expected one day with one program event, got %+v", days)
	}
	if got := days[0].ProgramEvents[0].Times; !slices.Equal(got, expectedTimes) {
		t.Fatalf("program card times = %q, want %q", got, expectedTimes)
	}
}
