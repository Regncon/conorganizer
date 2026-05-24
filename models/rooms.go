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
