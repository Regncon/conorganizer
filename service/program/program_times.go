package program

import (
	"fmt"
	"slices"

	"github.com/Regncon/conorganizer/models"
)

// PuljeTime is one pulje an event is placed in, with the "Tidspunkt" an admin
// wrote for it (relation_event_puljer.program_time). Only program events show it.
//
// A program event does not follow the pulje times. The pulje only decides which
// day it is shown on and which room it has, so the time comes from this text.
type PuljeTime struct {
	Pulje       models.PuljeRow
	ProgramTime string
}

// EventSchedule returns the schedule lines an event shows for the given puljer,
// which must be in pulje start order. A pulje event shows the pulje times, for
// example "Lørdag kveld · 18:00 - 23:00". A program event does not follow the
// pulje times and shows its own time one day at a time instead, for example
// "Lørdag 3.10 · Hele dagen". A day where no time is written shows only the day,
// so the event still says when it happens.
func EventSchedule(puljeTimes []PuljeTime, isProgramEvent bool) ([]string, error) {
	if !isProgramEvent {
		lines := make([]string, 0, len(puljeTimes))
		for _, puljeTime := range puljeTimes {
			lines = append(lines, fmt.Sprintf("%s · %s", puljeTime.Pulje.Name, puljeTime.Pulje.TimeRange()))
		}
		return lines, nil
	}
	return programEventSchedule(puljeTimes)
}

func programEventSchedule(puljeTimes []PuljeTime) ([]string, error) {
	location, err := Location()
	if err != nil {
		return nil, err
	}

	type dayTimes struct {
		day   Day
		times []string
	}
	var days []dayTimes
	for _, puljeTime := range puljeTimes {
		date := programDate(puljeTime.Pulje.StartAt.TimeOrZero(), location)
		if len(days) == 0 || !days[len(days)-1].day.Date.Equal(date) {
			days = append(days, dayTimes{day: Day{Date: date}})
		}
		current := &days[len(days)-1]
		current.times = appendProgramTime(current.times, puljeTime.ProgramTime)
	}

	lines := make([]string, 0, len(days))
	for _, day := range days {
		if len(day.times) == 0 {
			lines = append(lines, day.day.Label())
			continue
		}
		for _, programTime := range day.times {
			lines = append(lines, fmt.Sprintf("%s · %s", day.day.Label(), programTime))
		}
	}
	return lines, nil
}

// appendProgramTime adds a program event's time for one pulje to the times it
// shows for that day. Each distinct text is shown once, in pulje order, and an
// empty text shows nothing. An all-day event is placed in both of a day's
// puljer to keep its room all day, and should still say "Hele dagen" only once.
func appendProgramTime(times []string, programTime string) []string {
	if programTime == "" || slices.Contains(times, programTime) {
		return times
	}
	return append(times, programTime)
}
