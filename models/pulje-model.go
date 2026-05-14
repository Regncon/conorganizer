package models

import (
	"time"
)

type Pulje string

const (
	PuljeFredagKveld  Pulje = "FredagKveld"
	PuljeLordagMorgen Pulje = "LordagMorgen"
	PuljeLordagKveld  Pulje = "LordagKveld"
	PuljeSondagMorgen Pulje = "SondagMorgen"
)

var validPuljes = map[string]Pulje{
	string(PuljeFredagKveld):  PuljeFredagKveld,
	string(PuljeLordagMorgen): PuljeLordagMorgen,
	string(PuljeLordagKveld):  PuljeLordagKveld,
	string(PuljeSondagMorgen): PuljeSondagMorgen,
}

func ParsePulje(s string) (Pulje, bool) {
	p, ok := validPuljes[s]
	return p, ok
}

func ValidPuljeValues() []string {
	return []string{
		string(PuljeFredagKveld),
		string(PuljeLordagMorgen),
		string(PuljeLordagKveld),
		string(PuljeSondagMorgen),
	}
}

type PuljeRow struct {
	ID      Pulje     `json:"id"`
	Name    string    `json:"name"`
	Status  string    `json:"status"`
	StartAt time.Time `json:"start_at"`
	EndAt   time.Time `json:"end_at"`
}

type EventPulje struct {
	EventID     string `json:"event_id"`
	PuljeID     Pulje  `json:"pulje_id"`
	IsInPulje   bool   `json:"isInPulje"`
	IsPublished bool   `json:"isPublished"`
	RoomID      string `json:"room_id"`
}
