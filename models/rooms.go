package models

type Room struct {
	ID                  int    `json:"id"`
	Name                string `json:"name"`
	RoomNumber          string `json:"room_number"`
	Floor               int    `json:"floor"`
	MaxConcurrentEvents int    `json:"max_concurrent_events"`
	Notes               string `json:"note"`
	IsDisabled          bool   `json:"is_disabled"`
}

type RoomByPulje struct {
	ID                  int
	Name                string
	RoomNumber          string
	Events              []EventPulje
	Notes               string
	MaxConcurrentEvents int
}
