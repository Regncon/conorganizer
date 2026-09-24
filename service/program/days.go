package program

import (
	"database/sql"
	"fmt"
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
}

func (day Day) QueryValue() string {
	return day.Date.Format(programDateLayout)
}

func (day Day) Label() string {
	return fmt.Sprintf("%s %d.%d", norwegianWeekdays[day.Date.Weekday()], day.Date.Day(), day.Date.Month())
}

func GetDays(db *sql.DB) ([]Day, error) {
	location, err := time.LoadLocation(programTimeZone)
	if err != nil {
		return nil, fmt.Errorf("load program time zone %q: %w", programTimeZone, err)
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
func buildDays(puljer []models.PuljeRow, eventsByPulje map[models.Pulje][]models.EventCardModel, location *time.Location) []Day {
	days := make([]Day, 0)
	var seenProgramEvents map[string]struct{}
	for _, pulje := range puljer {
		date := programDate(pulje.StartAt.TimeOrZero(), location)
		if len(days) == 0 || !days[len(days)-1].Date.Equal(date) {
			days = append(days, Day{Date: date})
			seenProgramEvents = make(map[string]struct{})
		}

		// Keep the day available even when this pulje has no announced events.
		events := eventsByPulje[pulje.ID]
		if len(events) == 0 {
			continue
		}
		day := &days[len(days)-1]
		day.Blocks = append(day.Blocks, PuljeBlock{Pulje: pulje, Events: events})
		for _, event := range events {
			if event.IsInPuljefordeling {
				continue
			}
			if _, seen := seenProgramEvents[event.Id]; seen {
				continue
			}
			// A program event links to its first occurrence on this day.
			day.ProgramEvents = append(day.ProgramEvents, EventOccurrence{Event: event, PuljeID: pulje.ID})
			seenProgramEvents[event.Id] = struct{}{}
		}
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
