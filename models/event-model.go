package models

import (
	"database/sql"
)

type EventStatus string

const (
	EventStatusDraft     EventStatus = "Kladd"
	EventStatusSubmitted EventStatus = "Innsendt"
	EventStatusApproved  EventStatus = "Godkjent"
	EventStatusArchived  EventStatus = "Forkastet"
	EventStatusPublished EventStatus = "Publisert"
)

func (status EventStatus) Label() string {
	return string(status)
}

type EventType string

const (
	EventTypeRoleplay  EventType = "Roleplay"
	EventTypeBoardGame EventType = "Boardgame"
	EventTypeCardGame  EventType = "Cardgame"
	EventTypeOther     EventType = "Other"
)

func (eventType EventType) Label() string {
	switch eventType {
	case EventTypeRoleplay:
		return "Rollespill"
	case EventTypeBoardGame:
		return "Brettspill"
	case EventTypeCardGame:
		return "Kortspill"
	case EventTypeOther:
		return "Annet"
	default:
		return string(eventType)
	}
}

type AgeGroup string

const (
	AgeGroupDefault       AgeGroup = "Default"
	AgeGroupChildFriendly AgeGroup = "ChildFriendly"
	AgeGroupAdultsOnly    AgeGroup = "AdultsOnly"
)

func (ageGroup AgeGroup) Label() string {
	switch ageGroup {
	case AgeGroupDefault:
		return "Standard"
	case AgeGroupChildFriendly:
		return "Barn (under 12 år)"
	case AgeGroupAdultsOnly:
		return "Voksne (18+ år)"
	default:
		return string(ageGroup)
	}
}

func (ageGroup AgeGroup) BadgeLabel() string {
	switch ageGroup {
	case AgeGroupChildFriendly:
		return "Barnevennlig"
	case AgeGroupAdultsOnly:
		return "Egnet for voksne"
	default:
		return ""
	}
}

func (ageGroup AgeGroup) Valid() bool {
	switch ageGroup {
	case AgeGroupDefault, AgeGroupChildFriendly, AgeGroupAdultsOnly:
		return true
	default:
		return false
	}
}

type Runtime string

const (
	RunTimeNormal       Runtime = "Normal"
	RunTimeShortRunning Runtime = "ShortRunning"
	RunTimeLongRunning  Runtime = "LongRunning"
)

func (runtime Runtime) Label() string {
	switch runtime {
	case RunTimeNormal:
		return "Vanlig pulje"
	case RunTimeShortRunning:
		return "Kortere (2-3 timer)"
	case RunTimeLongRunning:
		return "Lengre (6+ timer)"
	default:
		return string(runtime)
	}
}

func (runtime Runtime) BadgeLabel() string {
	switch runtime {
	case RunTimeShortRunning:
		return "Varer under 3 timer"
	case RunTimeLongRunning:
		return "Varer over 6 timer"
	default:
		return ""
	}
}

func (runtime Runtime) Valid() bool {
	switch runtime {
	case RunTimeNormal, RunTimeShortRunning, RunTimeLongRunning:
		return true
	default:
		return false
	}
}

type Event struct {
	ID                  string         `json:"id"`
	Title               string         `json:"title"`
	Intro               string         `json:"intro"`
	Description         string         `json:"description"`
	System              string         `json:"system"`
	EventType           EventType      `json:"event_type"`
	AgeGroup            AgeGroup       `json:"age_group"`
	Runtime             Runtime        `json:"runtime"`
	HostName            string         `json:"host_name"`
	UserID              sql.NullInt64  `json:"user_id"`
	Email               string         `json:"email"`
	PhoneNumber         string         `json:"phone_number"`
	MaxPlayers          int            `json:"max_players"`
	BeginnerFriendly    bool           `json:"beginner_friendly"`
	CanBeRunInEnglish   bool           `json:"can_be_run_in_english"`
	Notes               string         `json:"notes"`
	Status              EventStatus    `json:"status"`
	CreatedAt           DBDateTime     `json:"created_at"`
	UpdatedAt           DBDateTime     `json:"updated_at"`
	CreatedByID         sql.NullInt64  `json:"created_by_id"`
	UpdatedByID         sql.NullInt64  `json:"updated_by_id"`
	StatusChangedByID   sql.NullInt64  `json:"status_changed_by_id"`
	StatusChangedAt     DBDateTime     `json:"status_changed_at"`
	StatusChangedAction sql.NullString `json:"status_changed_action"`
}

type EventCardModel struct {
	Id                string      `json:"id"`
	IsPublished       bool        `json:"is_published"`
	Title             string      `json:"title"`
	Intro             string      `json:"intro"`
	Status            EventStatus `json:"status"`
	System            string      `json:"system"`
	HostName          string      `json:"host_name"`
	EventType         EventType   `json:"event_type"`
	AgeGroup          AgeGroup    `json:"age_group"`
	Runtime           Runtime     `json:"runtime"`
	BeginnerFriendly  bool        `json:"beginner_friendly"`
	CanBeRunInEnglish bool        `json:"can_be_run_in_english"`
}
