package root

import (
	"database/sql"
	"fmt"
	"net/url"
	"time"

	"github.com/Regncon/conorganizer/models"
	puljerService "github.com/Regncon/conorganizer/service/puljer"
)

const (
	programDateQueryParam = "date"
	programDateLayout     = "2006-01-02"
	programTimeZone       = "Europe/Oslo"
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

type ProgramDay struct {
	Date   time.Time
	Blocks []PuljeBlock
}

func (day ProgramDay) QueryValue() string {
	return day.Date.Format(programDateLayout)
}

func (day ProgramDay) Label() string {
	return fmt.Sprintf("%s %d.%d", norwegianWeekdays[day.Date.Weekday()], day.Date.Day(), day.Date.Month())
}

func GetPublishedProgramDays(db *sql.DB) ([]ProgramDay, error) {
	location, err := time.LoadLocation(programTimeZone)
	if err != nil {
		return nil, fmt.Errorf("load program time zone %q: %w", programTimeZone, err)
	}

	puljer, err := puljerService.GetAllPuljer(db)
	if err != nil {
		return nil, err
	}

	blocks, err := getPublishedEventPuljeBlocks(db, puljer)
	if err != nil {
		return nil, err
	}

	return buildProgramDays(puljer, blocks, location), nil
}

func buildProgramDays(puljer []models.PuljeRow, blocks []PuljeBlock, location *time.Location) []ProgramDay {
	blocksByPulje := make(map[models.Pulje]PuljeBlock, len(blocks))
	for _, block := range blocks {
		blocksByPulje[block.Pulje.ID] = block
	}

	days := make([]ProgramDay, 0)
	dayIndexes := make(map[string]int)
	for _, pulje := range puljer {
		date := programDate(pulje.StartAt.TimeOrZero(), location)
		dateKey := date.Format(programDateLayout)

		dayIndex, ok := dayIndexes[dateKey]
		if !ok {
			dayIndex = len(days)
			dayIndexes[dateKey] = dayIndex
			days = append(days, ProgramDay{Date: date})
		}

		if block, ok := blocksByPulje[pulje.ID]; ok {
			days[dayIndex].Blocks = append(days[dayIndex].Blocks, block)
		}
	}

	return days
}

func SelectProgramDayIndex(days []ProgramDay, requestedDate string, now time.Time) int {
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

func programDateURL(date string) string {
	query := url.Values{}
	query.Set(programDateQueryParam, date)
	return "/?" + query.Encode()
}

func rootAPIURL(requestedDate string) string {
	if requestedDate == "" {
		return "/root/api"
	}

	query := url.Values{}
	query.Set(programDateQueryParam, requestedDate)
	return "/root/api?" + query.Encode()
}
