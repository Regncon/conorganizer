package root

import (
	"net/url"

	"github.com/Regncon/conorganizer/service/program"
)

const programDateQueryParam = "date"

type programEventRow struct {
	Events   []program.EventOccurrence
	Reversed bool
}

func programEventRows(events []program.EventOccurrence) []programEventRow {
	rows := make([]programEventRow, 0, (len(events)+1)/2)
	for index := 0; index < len(events); index += 2 {
		end := index + 2
		if end > len(events) {
			end = len(events)
		}
		rows = append(rows, programEventRow{
			Events:   events[index:end],
			Reversed: len(rows)%2 == 1,
		})
	}
	return rows
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
