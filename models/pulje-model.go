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

type PuljeRow struct {
	ID        Pulje     `json:"id"`
	Name      string    `json:"name"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

type EventPulje struct {
	EventID     string `json:"event_id"`
	PuljeID     Pulje  `json:"pulje_id"`
	IsPublished bool   `json:"isPublished"`
	Room        string `json:"room"`
}
