package program

import "testing"

func TestEventURL_PreservesOptionalContextAndEscapesIDs(t *testing.T) {
	cases := []struct {
		name, eventID, pulje, date, want string
	}{
		{name: "no context", eventID: "event", want: "/event/event"},
		{name: "pulje only", eventID: "event", pulje: "FredagKveld", want: "/event/event?pulje=FredagKveld"},
		{name: "date only", eventID: "event", date: "2026-10-03", want: "/event/event?date=2026-10-03"},
		{name: "full context", eventID: "event", pulje: "LordagMorgen", date: "2026-10-03", want: "/event/event?date=2026-10-03&pulje=LordagMorgen"},
		{name: "escaped ID", eventID: "event/with?spaces here", want: "/event/event%2Fwith%3Fspaces%20here"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := EventURL(tc.eventID, tc.pulje, tc.date); got != tc.want {
				t.Fatalf("URL = %q, want %q", got, tc.want)
			}
		})
	}
}
