// Package model defines the core data types for the puljefordeling algorithm.
package model

// Score is a player's interest level in an event (1–5).
type Score int

// MinScore and MaxScore are the inclusive bounds for preference scores.
const (
	MinScore Score = 1
	MaxScore Score = 5
)

// Event is a game session within a slot with a fixed seat count.
// DMIDs are the player IDs running this event; none can also be
// assigned as a participant in any event during the same slot.
type Event struct {
	ID       string
	Name     string
	Capacity int
	DMIDs    []string

	// AdultsOnly marks an 18+ game: a player who is not over 18 is never
	// seated here by the solver. Only an admin pin can place a minor here.
	AdultsOnly bool
}

// Slot is a time block containing one or more events.
type Slot struct {
	ID     string
	Name   string
	Events []Event
}

// Player is a convention attendee.
// Prefs[slotID][eventID] = score. A missing entry means no interest.
type Player struct {
	ID    string
	Name  string
	Prefs map[string]map[string]Score

	// IsOver18 gates access to AdultsOnly events.
	IsOver18 bool
}

// Weekend holds all slots and players for a convention.
type Weekend struct {
	Slots   []Slot
	Players []Player
}

// ScoreBreakdown is how the solver valued one player's seat: the interest band
// plus the fairness bumps. Total is the edge weight the solver maximised.
type ScoreBreakdown struct {
	Score           Score // the player's interest in the event
	Satisfied       bool  // already had a top choice before this slot
	Band            int
	Misses          int // prior puljer with a missed top choice
	MissBonus       int
	NeverSeatedBump int
	DMBump          int
	Total           int
}

// SlotResult is the assignment output for a single slot.
type SlotResult struct {
	SlotID                string
	Assignments           map[string][]string       // eventID -> assigned playerIDs
	UndersubscribedEvents []string                  // eventIDs assigned fewer players than MinPlayers — flagged for organiser review (not cancelled)
	Unassigned            []string                  // playerIDs with interest but no seat
	NewlySatisfied        []string                  // playerIDs satisfied for the first time this slot
	MovedPlayers          []string                  // playerIDs bumped down to a strictly lower-interest event to make room (lateral, equal-interest swaps are excluded)
	TotalScore            int                       // sum of actual (unadjusted) scores for all assignments
	Scores                map[string]ScoreBreakdown // playerID -> how the solver valued their seat, pins included when the player has interest in the event; empty for replayed puljer
	Seed                  int64                     // seed used for tie-breaking shuffle this slot
}
