package models

import (
	"database/sql"
	"time"
)

type EventStatus string

const (
	EventStatusDraft     EventStatus = "Kladd"
	EventStatusSubmitted EventStatus = "Innsendt"
	EventStatusPublished EventStatus = "Publisert"
	EventStatusApproved  EventStatus = "Godkjent"
	EventStatusArchived  EventStatus = "Avist"
)

type EventType string

const (
	EventTypeRoleplay  EventType = "roleplay"
	EventTypeBoardGame EventType = "boardgame"
	EventTypeCardGame  EventType = "cardgame"
	EventTypeOther     EventType = "other"
)

type AgeGroup string

const (
	AgeGroupAllAges       AgeGroup = "AllAges"
	AgeGroupChildFriendly AgeGroup = "ChildFriendly"
	AgeGroupTeenFriendly  AgeGroup = "TeenFriendly"
	AgeGroupAdultsOnly    AgeGroup = "AdultsOnly"
)

type Runtime string

const (
	RunTimeNormal       Runtime = "Normal"
	RunTimeShortRunning Runtime = "ShortRunning"
	RunTimeLongRunning  Runtime = "LongRunning"
)

type Event struct {
	ID                string         `json:"id"`
	Title             string         `json:"title"`
	Intro             string         `json:"intro"`
	Description       string         `json:"description"`
	ImageURL          sql.NullString `json:"image_url"`
	System            string         `json:"system"`
	EventType         EventType      `json:"event_type"`
	AgeGroup          AgeGroup       `json:"age_group"`
	Runtime           Runtime        `json:"runtime"`
	HostName          string         `json:"host_name"`
	Host              sql.NullInt64  `json:"host"`
	Email             string         `json:"email"`
	PhoneNumber       string         `json:"phone_number"`
	PuljeName         sql.NullString `json:"pulje_name"`
	MaxPlayers        int            `json:"max_players"`
	BeginnerFriendly  bool           `json:"beginner_friendly"`
	CanBeRunInEnglish bool           `json:"can_be_run_in_english"`
	Notes             string         `json:"notes"`
	Status            EventStatus    `json:"status"`
	InsertedTime      time.Time      `json:"inserted_time"`
}
