package models

// Room is the internal struct of a room
type Room struct {
	ID                  int
	Name                string
	RoomNumber          string
	Floor               int
	MaxConcurrentEvents int
	Notes               string
	IsDisabled          bool
}

// RoomJSON is the JSON representation of `Room` type for use with db query and front-end
type RoomJSON struct {
	ID                  int    `json:"id"`
	Name                string `json:"name"`
	RoomNumber          string `json:"room_number"`
	Floor               int    `json:"floor"`
	MaxConcurrentEvents int    `json:"max_concurrent_events"`
	Notes               string `json:"notes"`
	IsDisabled          bool   `json:"is_disabled"`
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
}

// RoomByPulje is a snapshot of room delegation for a specific pulje, this is mainly used for figuring
// out what `max_concurrent_events` is based on a pulje, but also for the dropdown input component
// used in assigning rooms to an event in a pulje
type RoomByPulje struct {
	ID                  int
	Name                string
	RoomNumber          string
	AssignedEventsID    []RoomEventPuljeSummary
	MaxConcurrentEvents int
	Notes               string
}
