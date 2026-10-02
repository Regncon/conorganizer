package program

import (
	"net/url"

	"github.com/Regncon/conorganizer/models"
)

// ProgramDateURL builds a homepage URL that selects the requested program day.
func ProgramDateURL(date string) string {
	query := url.Values{}
	query.Set("date", date)
	return "/?" + query.Encode()
}

// ProgramPuljeURL selects a program date and anchors to a pulje on the homepage.
func ProgramPuljeURL(date string, puljeID models.Pulje) string {
	return ProgramDateURL(date) + "#pulje-" + string(puljeID)
}

// EventURL preserves the selected pulje and day when linking to an event.
func EventURL(eventID string, pulje string, date string) string {
	query := url.Values{}
	if date != "" {
		query.Set("date", date)
	}
	if pulje != "" {
		query.Set("pulje", pulje)
	}

	href := "/event/" + url.PathEscape(eventID)
	if encoded := query.Encode(); encoded != "" {
		href += "?" + encoded
	}
	return href
}
