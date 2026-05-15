package models

import (
	"database/sql"
	"time"
)

type EventStatus string

const (
	EventStatusDraft     EventStatus = "Kladd"
	EventStatusSubmitted EventStatus = "Innsendt"
	EventStatusApproved  EventStatus = "Godkjent"
	EventStatusArchived  EventStatus = "Forkastet"
)

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

type Runtime string

const (
	RunTimeNormal       Runtime = "Normal"
	RunTimeShortRunning Runtime = "ShortRunning"
	RunTimeLongRunning  Runtime = "LongRunning"
)

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
	CreatedAt           time.Time      `json:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at"`
	CreatedByID         sql.NullInt64  `json:"created_by_id"`
	UpdatedByID         sql.NullInt64  `json:"updated_by_id"`
	StatusChangedByID   sql.NullInt64  `json:"status_changed_by_id"`
	StatusChangedAt     sql.NullTime   `json:"status_changed_at"`
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
