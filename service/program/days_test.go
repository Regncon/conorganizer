package program

import (
	"slices"
	"testing"
	"time"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil/bdd"
)

func TestSelectDayIndex_UsesRequestedDateBeforeDefault(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at programdagene er tilgjengelige og en dato er valgt i URL-en.",
		When:  "Når valgt programdag bestemmes.",
		Then:  "Så skal datoen fra URL-en brukes foran standardvalget.",
	})

	// Given
	expectedIndex := 2
	days := testDays()
	requestedDate := "2026-10-04"
	now := time.Date(2026, 10, 1, 12, 0, 0, 0, time.UTC)

	// When
	actualIndex := SelectDayIndex(days, requestedDate, now)

	// Then
	if actualIndex != expectedIndex {
		t.Fatalf("selected day index = %d, want %d", actualIndex, expectedIndex)
	}
}

func TestSelectDayIndex_ClampsDefaultToProgramBoundaries(t *testing.T) {
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
	days := testDays()

	// When
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Then
			actualIndex := SelectDayIndex(days, "", tc.now)
			if actualIndex != expectedIndexes[tc.name] {
				t.Fatalf("selected day index = %d, want %d", actualIndex, expectedIndexes[tc.name])
			}
		})
	}
}

func testDays() []Day {
	location := time.UTC
	return []Day{
		{Date: time.Date(2026, 10, 2, 0, 0, 0, 0, location)},
		{Date: time.Date(2026, 10, 3, 0, 0, 0, 0, location)},
		{Date: time.Date(2026, 10, 4, 0, 0, 0, 0, location)},
	}
}

func TestSelectDayIndex_HandlesMissingDatesAndOsloMidnight(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt programdager i Oslo og en valgfri dato fra URL-en.",
		When:  "Når valgt dag bestemmes rundt midnatt eller med manglende data.",
		Then:  "Så brukes lokal dato som standard, mens tomme programmer ikke får en valgt dag.",
	})

	location, err := time.LoadLocation("Europe/Oslo")
	if err != nil {
		t.Fatal(err)
	}
	days := []Day{
		{Date: time.Date(2026, 10, 2, 0, 0, 0, 0, location)},
		{Date: time.Date(2026, 10, 3, 0, 0, 0, 0, location)},
		{Date: time.Date(2026, 10, 4, 0, 0, 0, 0, location)},
	}
	beforeMidnight := time.Date(2026, 10, 2, 21, 59, 0, 0, time.UTC)
	midnight := time.Date(2026, 10, 2, 22, 0, 0, 0, time.UTC)
	cases := []struct {
		name string
		days []Day
		date string
		now  time.Time
		want int
	}{
		{name: "no program days", now: midnight, want: -1},
		{name: "before local midnight", days: days, now: beforeMidnight, want: 0},
		{name: "at local midnight", days: days, now: midnight, want: 1},
		{name: "invalid date falls back", days: days, date: "invalid", now: midnight, want: 1},
		{name: "unknown date falls back", days: days, date: "2026-11-01", now: midnight, want: 1},
		{name: "requested day wins", days: days, date: "2026-10-04", now: midnight, want: 2},
		{name: "gap ties choose earlier day", days: []Day{days[0], days[2]}, now: midnight, want: 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := SelectDayIndex(tc.days, tc.date, tc.now); got != tc.want {
				t.Fatalf("selected day = %d, want %d", got, tc.want)
			}
		})
	}
}

func TestBuildDays_PreservesEmptyDaysAndDeduplicatesProgramWithinEachLocalDay(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en tom fredag og et programarrangement i flere puljer lørdag og søndag.",
		When:  "Når programdagene bygges.",
		Then:  "Så beholdes fredag, programarrangementet vises én gang per lokal dag, og rafflearrangementer beholdes i hver pulje.",
	})

	location, err := time.LoadLocation("Europe/Oslo")
	if err != nil {
		t.Fatal(err)
	}
	puljer := []models.PuljeRow{
		{ID: models.PuljeFredagKveld, StartAt: models.NewDBDateTime(time.Date(2026, 10, 2, 18, 0, 0, 0, time.UTC))},
		{ID: models.PuljeLordagMorgen, StartAt: models.NewDBDateTime(time.Date(2026, 10, 2, 22, 30, 0, 0, time.UTC))},
		{ID: models.PuljeLordagKveld, StartAt: models.NewDBDateTime(time.Date(2026, 10, 3, 18, 0, 0, 0, time.UTC))},
		{ID: models.PuljeSondagMorgen, StartAt: models.NewDBDateTime(time.Date(2026, 10, 4, 8, 0, 0, 0, time.UTC))},
	}
	programEvent := models.EventCardModel{Id: "shared-program"}
	raffleEvent := models.EventCardModel{Id: "shared-raffle", IsInPuljefordeling: true}
	events := map[models.Pulje][]models.EventCardModel{
		models.PuljeLordagMorgen: {programEvent, raffleEvent},
		models.PuljeLordagKveld:  {programEvent, raffleEvent},
		models.PuljeSondagMorgen: {programEvent},
	}

	days := buildDays(puljer, events, location)

	if len(days) != 3 {
		t.Fatalf("day count = %d, want 3", len(days))
	}
	dates := []string{days[0].QueryValue(), days[1].QueryValue(), days[2].QueryValue()}
	if want := []string{"2026-10-02", "2026-10-03", "2026-10-04"}; !slices.Equal(dates, want) {
		t.Fatalf("dates = %v, want %v", dates, want)
	}
	if len(days[0].Blocks) != 0 || len(days[0].ProgramEvents) != 0 {
		t.Fatal("empty Friday contains events")
	}
	for _, index := range []int{1, 2} {
		if len(days[index].ProgramEvents) != 1 || days[index].ProgramEvents[0].Event.Id != programEvent.Id {
			t.Fatalf("program events on %s = %v, want one shared program event", days[index].QueryValue(), days[index].ProgramEvents)
		}
	}
	if got := days[1].ProgramEvents[0].PuljeID; got != models.PuljeLordagMorgen {
		t.Fatalf("Saturday program link uses %s, want first occurrence", got)
	}
	if got := days[2].ProgramEvents[0].PuljeID; got != models.PuljeSondagMorgen {
		t.Fatalf("Sunday program link uses %s, want Sunday occurrence", got)
	}
	if len(days[1].Blocks) != 2 {
		t.Fatalf("Saturday pulje count = %d, want 2", len(days[1].Blocks))
	}
	for _, block := range days[1].Blocks {
		raffleEvents := block.RaffleEvents()
		if len(raffleEvents) != 1 || raffleEvents[0].Id != raffleEvent.Id {
			t.Fatalf("raffle events in %s = %v, want shared raffle event", block.Pulje.ID, raffleEvents)
		}
	}
}
