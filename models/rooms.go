package models

import (
	"database/sql"
)

type Room struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	RoomNumber string `json:"room_number"`
	Floor      int    `json:"floor"`
	Notes      string `json:"notes"`
}

// Normalized version of `Room` type for use when updating a room, or quering for a specific room with optional params
type RoomInput struct {
	ID         int
	Name       *string
	RoomNumber *string
	Floor      *int
	Notes      *string
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

// RoomByPulje is a snapshot of room assignments for a specific pulje.
type RoomByPulje struct {
	ID               int64
	Name             string
	RoomNumber       string
	Floor            int
	Notes            string
	AssignedEventsID []RoomEventPuljeSummary
}

// RoomStatusByPulje is a map of puljer containing room statuses, such as which games are assigned to that room
// You can access status by keys: [Pulje][RoomID]
type RoomStatusByPulje = map[Pulje]map[int64]RoomByPulje

type RoomStatusRow struct {
	PuljeID    Pulje
	RoomID     int64
	RoomName   string
	RoomNumber string
	Floor      int
	RoomNotes  string

	EventID         sql.NullString
	EventTitle      sql.NullString
	EventMaxPlayers sql.NullInt32
}

// RoomFormSignals is used in data-star input form bindings for sending signals to users
type RoomFormSignals struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	RoomNumber string `json:"room_number"`
	Floor      int    `json:"floor"`
	Notes      string `json:"notes"`

	Mode        string `json:"mode"`
	FormTitle   string `json:"form_title"`
	ButtonLabel string `json:"button_label"`
}

type RoomErrorKey string

const (
	RoomError           RoomErrorKey = "error"
	RoomErrorFloor      RoomErrorKey = "floor"
	RoomErrorName       RoomErrorKey = "name"
	RoomErrorNotes      RoomErrorKey = "notes"
	RoomErrorRoomNumber RoomErrorKey = "room_number"
)

// RoomFormErrors is used in validation and error handling when creating and updating rooms
type RoomFormErrors map[RoomErrorKey]string

// ResetErrors resets all the errors to empty strings
func (errors RoomFormErrors) ResetErrors() {
	errors[RoomError] = ""
	errors[RoomErrorFloor] = ""
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
