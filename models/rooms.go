package models

import (
	"database/sql"
)

type Room struct {
	ID                 int    `json:"id"`
	Name               string `json:"name"`
	RoomNumber         string `json:"room_number"`
	Floor              int    `json:"floor"`
	MaxConcurrentGames int    `json:"max_concurrent_games"`
	Notes              string `json:"notes"`
	IsDisabled         bool   `json:"is_disabled"`
}

// Normalized version of `Room` type for use when updating a room, or quering for a specific room with optional params
type RoomInput struct {
	ID                 int
	Name               *string
	RoomNumber         *string
	Floor              *int
	MaxConcurrentGames *int
	Notes              *string
	IsDisabled         *bool
}

/*
RoomEventPuljeSummary is the summary of an event in `relation_event_puljer` and used in `RoomByPulje` struct
  - `EventPuljeID` is the ID of the unique event in a pulje
  - `EventID`      is the ID of the pulje the unique event is in
  - `Title`        is the title of the event
*/
type RoomEventPuljeSummary struct {
	EventPuljeID string
	EventID      string
	Title        string
	MaxPlayers   int
	RoomID       int64
}
type RoomEventPuljeSummaryJson struct {
	EventPuljeID string        `json:"pulje_id"`
	EventID      string        `json:"event_id"`
	Title        string        `json:"title"`
	MaxPlayers   int           `json:"max_players"`
	RoomID       sql.NullInt64 `json:"room_id"`
}

// RoomByPulje is a snapshot of room delegation for a specific pulje, this is mainly used for figuring
// out what `max_concurrent_events` is based on a pulje, but also for the dropdown input component
// used in assigning rooms to an event in a pulje
type RoomByPulje struct {
	ID                 int64
	Name               string
	RoomNumber         string
	Floor              int
	MaxConcurrentGames int
	IsDisabled         bool
	Notes              string
	AssignedEventsID   []RoomEventPuljeSummary
}

// Helper function for getting currently assigned events to a room in a pulje
func (r RoomByPulje) CurrentOccupancy() int {
	return len(r.AssignedEventsID)
}

// Helper function for getting available free slots for a room in a pulje
func (r RoomByPulje) RemainingCapacity() int {
	return r.MaxConcurrentGames - len(r.AssignedEventsID)
}

// Helper function for quickly checking if a room is full in a pulje
func (r RoomByPulje) IsFull() bool {
	return len(r.AssignedEventsID) >= r.MaxConcurrentGames
}

// RoomStatusByPulje is a map of puljer containing room statuses, such as which games are assigned to that room
// You can access status by keys: [Pulje][RoomID]
type RoomStatusByPulje = map[Pulje]map[int64]RoomByPulje

type RoomStatusRow struct {
	PuljeID Pulje

	RoomID             int64
	RoomName           string
	RoomNumber         string
	Floor              int
	MaxConcurrentGames int
	IsDisabled         bool
	RoomNotes          string

	EventID         sql.NullString
	EventTitle      sql.NullString
	EventMaxPlayers sql.NullInt32
}

// RoomFormSignals is used in data-star input form bindings for sending signals to users
type RoomFormSignals struct {
	ID                 int    `json:"id"`
	Name               string `json:"name"`
	RoomNumber         string `json:"room_number"`
	Floor              int    `json:"floor"`
	MaxConcurrentGames int    `json:"max_concurrent_games"`
	Notes              string `json:"notes"`
	IsDisabled         bool   `json:"is_disabled"`

	Mode        string `json:"mode"`
	FormTitle   string `json:"form_title"`
	ButtonLabel string `json:"button_label"`
}

type RoomErrorKey string

const (
	RoomError              RoomErrorKey = "error"
	RoomErrorFloor         RoomErrorKey = "floor"
	RoomErrorIsDisabled    RoomErrorKey = "is_disabled"
	RoomErrorMaxConcurrent RoomErrorKey = "max_concurrent_games"
	RoomErrorName          RoomErrorKey = "name"
	RoomErrorNotes         RoomErrorKey = "notes"
	RoomErrorRoomNumber    RoomErrorKey = "room_number"
)

// RoomFormErrors is used in validation and error handling when creating and updating rooms
type RoomFormErrors map[RoomErrorKey]string

// ResetErrors resets all the errors to empty strings
func (errors RoomFormErrors) ResetErrors() {
	errors[RoomError] = ""
	errors[RoomErrorFloor] = ""
	errors[RoomErrorIsDisabled] = ""
	errors[RoomErrorMaxConcurrent] = ""
	errors[RoomErrorName] = ""
	errors[RoomErrorNotes] = ""
	errors[RoomErrorRoomNumber] = ""
}

// AddError is a helper function for adding an error message
func (errors RoomFormErrors) AddError(errorKey RoomErrorKey, errorMessage string) {
	errors[errorKey] = errorMessage
}

// HasErrors is a hepler function for quickly checking if a certain error exists
func (errors RoomFormErrors) HasError(errorKey RoomErrorKey) bool {
	return errors[errorKey] != ""
}

// HasErrors is a hepler function for quickly checking if any errors exists from validation
func (errors RoomFormErrors) HasErrors() bool {
	for _, msg := range errors {
		if msg != "" {
			return true
		}
	}
	return false
}

func (errors RoomFormErrors) GetKeys() []string {
	keys := make([]string, 0, len(errors))

	for key := range errors {
		keys = append(keys, string(key))
	}

	return keys
}
