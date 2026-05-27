package models

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
}

// RoomByPulje is a snapshot of room delegation for a specific pulje, this is mainly used for figuring
// out what `max_concurrent_events` is based on a pulje, but also for the dropdown input component
// used in assigning rooms to an event in a pulje
type RoomByPulje struct {
	ID                 int
	Name               string
	RoomNumber         string
	AssignedEventsID   []RoomEventPuljeSummary
	MaxConcurrentGames int
	Notes              string
}

// RoomStatusByPulje is a map of puljer containing room statuses, such as which games are assigned to that room
// You can access status by keys: [Pulje][RoomID]
type RoomStatusByPulje = map[Pulje]map[int64]RoomByPulje
