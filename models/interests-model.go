package models

import (
	"time"
)

type Interest struct {
	BillettholderId int       `json:"billettholder_id"`
	EventId         string    `json:"event_id"`
	PuljeId         string    `json:"pulje_id"`
	InterestLevel   string    `json:"interest_level"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
