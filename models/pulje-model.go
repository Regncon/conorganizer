package models

import (
	"database/sql"
	"fmt"
)

type Pulje string

const (
	PuljeFredagKveld  Pulje = "FredagKveld"
	PuljeLordagMorgen Pulje = "LordagMorgen"
	PuljeLordagKveld  Pulje = "LordagKveld"
	PuljeSondagMorgen Pulje = "SondagMorgen"
)

func AllPuljer() []Pulje {
	return []Pulje{
		PuljeFredagKveld,
		PuljeLordagMorgen,
		PuljeLordagKveld,
		PuljeSondagMorgen,
	}
}

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
	puljes := AllPuljer()
	values := make([]string, len(puljes))
	for i, pulje := range puljes {
		values[i] = string(pulje)
	}
	return values
}

type PuljeStatus string

const (
	PuljeStatusOpen         PuljeStatus = "open"
	PuljeStatusNotPublished PuljeStatus = "not_published"
	PuljeStatusPublished    PuljeStatus = "published"
	PuljeStatusLocked       PuljeStatus = "locked"
	PuljeStatusCompleted    PuljeStatus = "completed"
)

func (status PuljeStatus) Label() string {
	switch status {
	case PuljeStatusOpen:
		return "Åpen"
	case PuljeStatusNotPublished:
		return "Ikke publisert"
	case PuljeStatusPublished:
		return "Publisert"
	case PuljeStatusLocked:
		return "Låst"
	case PuljeStatusCompleted:
		return "Fullført"
	default:
		return string(status)
	}
}

type PuljeRow struct {
	ID      Pulje       `json:"id"`
	Name    string      `json:"name"`
	Status  PuljeStatus `json:"status"`
	StartAt DBDateTime  `json:"start_at"`
	EndAt   DBDateTime  `json:"end_at"`
}

func (pulje PuljeRow) TimeRange() string {
	return fmt.Sprintf("%s - %s", pulje.StartAt.Format("15:04"), pulje.EndAt.Format("15:04"))
}

type EventPulje struct {
	EventID     string        `json:"event_id"`
	PuljeID     Pulje         `json:"pulje_id"`
	IsInPulje   bool          `json:"isInPulje"`
	IsPublished bool          `json:"isPublished"`
	RoomID      sql.NullInt64 `json:"room_id"`
}
