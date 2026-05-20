package models

import (
	"database/sql"
)

type InterestLevel string

const (
	InterestLevelNone   InterestLevel = ""
	InterestLevelHigh   InterestLevel = "Veldig interessert"
	InterestLevelMedium InterestLevel = "Middels interessert"
	InterestLevelLow    InterestLevel = "Litt interessert"
)

func (level InterestLevel) Label() string {
	switch level {
	case InterestLevelHigh:
		return "Veldig interessert"
	case InterestLevelMedium:
		return "Interessert"
	case InterestLevelLow:
		return "Litt interessert"
	case InterestLevelNone:
		return "Ikke interessert"
	default:
		return string(level)
	}
}

func (level InterestLevel) Valid() bool {
	switch level {
	case InterestLevelHigh, InterestLevelMedium, InterestLevelLow, InterestLevelNone:
		return true
	default:
		return false
	}
}

type Interest struct {
	BillettholderId int           `json:"billettholder_id"`
	EventId         string        `json:"event_id"`
	PuljeId         string        `json:"pulje_id"`
	InterestLevel   InterestLevel `json:"interest_level"`
	CreatedAt       DBDateTime    `json:"created_at"`
	UpdatedAt       DBDateTime    `json:"updated_at"`
	CreatedByID     sql.NullInt64 `json:"created_by_id"`
	UpdatedByID     sql.NullInt64 `json:"updated_by_id"`
}
