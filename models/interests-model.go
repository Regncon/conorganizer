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

// interestScores is the single source of truth mapping interest levels onto the
// puljefordeling solver's 1–5 preference scale. "Veldig interessert" is the top
// choice (5) and is the threshold the solver uses for participant satisfaction;
// "Middels" is 3 and "Litt" is 1. Both Score and InterestLevelFromScore read
// this table, so the two directions can never drift apart.
var interestScores = []struct {
	level InterestLevel
	score int
}{
	{InterestLevelHigh, 5},
	{InterestLevelMedium, 3},
	{InterestLevelLow, 1},
}

// Score maps an interest level onto the solver's 1–5 preference scale.
// No/invalid interest returns 0, meaning no edge in the assignment graph.
func (level InterestLevel) Score() int {
	for _, m := range interestScores {
		if m.level == level {
			return m.score
		}
	}
	return 0
}

// InterestLevelFromScore reverses Score: it maps a solver preference score back
// to the interest level it came from, returning InterestLevelNone for any score
// that is not one of the defined levels (including 0).
func InterestLevelFromScore(score int) InterestLevel {
	for _, m := range interestScores {
		if m.score == score {
			return m.level
		}
	}
	return InterestLevelNone
}

// Emoji returns a short glyph for an interest level, used to show at a glance
// how much a seated participant wanted the game they got. These are the single
// source of truth for the interest glyphs and match the buttons in the interest
// picker (TicketHolderInterestPicker), so the two views never drift apart.
// InterestLevelNone deliberately returns "" — the puljefordeling box uses the
// empty string as the signal to substitute a 📌 pin for a manual seat that has
// no real interest behind it.
func (level InterestLevel) Emoji() string {
	switch level {
	case InterestLevelHigh:
		return "🤩"
	case InterestLevelMedium:
		return "🙂"
	case InterestLevelLow:
		return "🤨"
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
