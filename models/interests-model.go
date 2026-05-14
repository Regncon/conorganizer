package models

import "time"

const (
	InterestLevelHigh   = "Veldig interessert"
	InterestLevelMedium = "Middels interessert"
	InterestLevelLow    = "Litt interessert"
)

type Interest struct {
	BillettholderId int           `json:"billettholder_id"`
	EventId         string        `json:"event_id"`
	PuljeId         string        `json:"pulje_id"`
	InterestLevel   string        `json:"interest_level"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`
	CreatedByID     sql.NullInt64 `json:"created_by_id"`
	UpdatedByID     sql.NullInt64 `json:"updated_by_id"`
}
