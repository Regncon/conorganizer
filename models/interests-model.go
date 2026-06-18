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
		return "Ikkje interessert"
	default:
		return string(level)
	}
}

// Score maps an interest level onto the puljefordeling solver's 1–5 preference
// scale. "Veldig interessert" is the top choice (5) and is the threshold the
// solver uses for participant satisfaction; "Middels" is 3 and "Litt" is 1.
// No/invalid interest returns 0, meaning no edge in the assignment graph.
func (level InterestLevel) Score() int {
	switch level {
	case InterestLevelHigh:
		return 5
	case InterestLevelMedium:
		return 3
	case InterestLevelLow:
		return 1
	default:
		return 0
	}
}

// Emoji returns a short glyph for an interest level, used to show at a glance
// how much a seated participant wanted the game they got.
func (level InterestLevel) Emoji() string {
	switch level {
	case InterestLevelHigh:
		return "🔥"
	case InterestLevelMedium:
		return "👍"
	case InterestLevelLow:
		return "🤷"
	default:
		return ""
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
