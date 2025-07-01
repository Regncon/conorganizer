package models

import (
	"database/sql"
	"time"
)

type EventStatus string

const (
	EventStatusDraft     EventStatus = "Kladd"
	EventStatusPublished EventStatus = "Publisert"
	EventStatusClosed    EventStatus = "Godkjent"
	EventStatusArchived  EventStatus = "Avist"
)

type AgeGroup struct {
	ChildFriendly bool
	AdultsOnly    bool
}
type Duration struct {
	LongRunning  bool
	ShortRunning bool
}

type Event struct {
	ID                int64          `json:"id"`
	Title             string         `json:"title"`
	Description       string         `json:"description"`
	ImageURL          sql.NullString `json:"image_url"`
	System            sql.NullString `json:"system"`
	HostName          string         `json:"host_name"`
	Host              sql.NullInt64  `json:"host"`
	Email             string         `json:"email"`
	PhoneNumber       int64          `json:"phone_number"`
	RoomName          sql.NullString `json:"room_name"`
	PuljeName         sql.NullString `json:"pulje_name"`
	MaxPlayers        int64          `json:"max_players"`
	ChildFriendly     bool           `json:"child_friendly"`
	AdultsOnly        bool           `json:"adults_only"`
	BeginnerFriendly  bool           `json:"beginner_friendly"`
	ExperiencedOnly   bool           `json:"experienced_only"`
	CanBeRunInEnglish bool           `json:"can_be_run_in_english"`
	LongRunning       bool           `json:"long_running"`
	ShortRunning      bool           `json:"short_running"`
	Status            EventStatus    `json:"status"`
	InsertedTime      time.Time      `json:"inserted_time"`
	AgeGroup          AgeGroup
	Duration          Duration
}
