package program

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/Regncon/conorganizer/models"
	puljerService "github.com/Regncon/conorganizer/service/puljer"
)

const (
	programDateLayout = "2006-01-02"
	programTimeZone   = "Europe/Oslo"
)

var norwegianWeekdays = [...]string{
	"Søndag",
	"Mandag",
	"Tirsdag",
	"Onsdag",
	"Torsdag",
	"Fredag",
	"Lørdag",
}

type Day struct {
	Date          time.Time
	ProgramEvents []EventOccurrence
	Blocks        []PuljeBlock
}

type EventOccurrence struct {
	Event   models.EventCardModel
	PuljeID models.Pulje
	// Times is what the program event card shows for the day, one line each.
	Times []string
}

func (day Day) QueryValue() string {
	return day.Date.Format(programDateLayout)
}

func (day Day) Label() string {
	return fmt.Sprintf("%s %d.%d", norwegianWeekdays[day.Date.Weekday()], day.Date.Day(), day.Date.Month())
}

var loadLocation = sync.OnceValues(func() (*time.Location, error) {
	location, err := time.LoadLocation(programTimeZone)
	if err != nil {
		return nil, fmt.Errorf("load program time zone %q: %w", programTimeZone, err)
	}
	return location, nil
})

// Location is the convention's time zone (Europe/Oslo), used for every time
// shown to users. It is loaded once.
func Location() (*time.Location, error) {
	return loadLocation()
}

func GetDays(db *sql.DB) ([]Day, error) {
	location, err := Location()
	if err != nil {
		return nil, err
	}

	puljer, err := puljerService.GetAllPuljer(db)
	if err != nil {
		return nil, err
	}

	eventsByPulje, err := getEventsByPulje(db)
	if err != nil {
		return nil, err
	}

	return buildDays(puljer, eventsByPulje, location), nil
}

// buildDays keeps puljer in their query order (start time ascending).
func buildDays(puljer []models.PuljeRow, eventsByPulje map[models.Pulje][]puljeEvent, location *time.Location) []Day {
	days := make([]Day, 0)
	var programEventIndex map[string]int
	for _, pulje := range puljer {
		date := programDate(pulje.StartAt.TimeOrZero(), location)
		if len(days) == 0 || !days[len(days)-1].Date.Equal(date) {
			days = append(days, Day{Date: date})
			programEventIndex = make(map[string]int)
		}

		// Keep the day available even when this pulje has no announced events.
		puljeEvents := eventsByPulje[pulje.ID]
		if len(puljeEvents) == 0 {
			continue
		}
		day := &days[len(days)-1]
		block := PuljeBlock{Pulje: pulje, Events: make([]models.EventCardModel, 0, len(puljeEvents))}
		for _, puljeEvent := range puljeEvents {
			block.Events = append(block.Events, puljeEvent.Event)
			if puljeEvent.Event.IsInPuljefordeling {
				continue
			}
			index, seen := programEventIndex[puljeEvent.Event.Id]
			if !seen {
				// A program event links to its first occurrence on this day.
				index = len(day.ProgramEvents)
				programEventIndex[puljeEvent.Event.Id] = index
				day.ProgramEvents = append(day.ProgramEvents, EventOccurrence{Event: puljeEvent.Event, PuljeID: pulje.ID})
			}
			occurrence := &day.ProgramEvents[index]
			occurrence.Times = appendProgramTime(occurrence.Times, puljeEvent.ProgramTime)
		}
		day.Blocks = append(day.Blocks, block)
	}
	return days
}

func SelectDayIndex(days []Day, requestedDate string, now time.Time) int {
	if len(days) == 0 {
		return -1
	}

	for index, day := range days {
		if day.QueryValue() == requestedDate {
			return index
		}
	}

	location := days[0].Date.Location()
	today := programDate(now, location)
	if !today.After(days[0].Date) {
		return 0
	}
	if !today.Before(days[len(days)-1].Date) {
		return len(days) - 1
	}

	for index := 1; index < len(days); index++ {
		if today.Before(days[index].Date) {
			before := days[index-1].Date
			after := days[index].Date
			if today.Sub(before) <= after.Sub(today) {
				return index - 1
			}
			return index
		}
	}

	return len(days) - 1
}

func programDate(value time.Time, location *time.Location) time.Time {
	local := value.In(location)
	return time.Date(local.Year(), local.Month(), local.Day(), 0, 0, 0, 0, location)
}

type PuljeBlock struct {
	Pulje  models.PuljeRow
	Events []models.EventCardModel
}

func (block PuljeBlock) RaffleEvents() []models.EventCardModel {
	events := make([]models.EventCardModel, 0, len(block.Events))
	for _, event := range block.Events {
		if event.IsInPuljefordeling {
			events = append(events, event)
		}
	}
	return events
}
