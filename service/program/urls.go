package program

import "net/url"

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
