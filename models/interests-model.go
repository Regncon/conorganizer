package models

import (
	"time"
)

type Interest struct {
	BillettholderId int       `json:"billettholder_id"`
	EventId         string    `json:"event_id"`
	PuljeId         string    `json:"pulje_id"`
	InterestLevel   string    `json:"interest_level"`
	InsertedTime    time.Time `json:"inserted_time"`
}
