// Package puljefordeling runs the seating-distribution solver against live
// conorganizer data to produce a read-only preview of how participants would
// be assigned to events in each pulje. It never writes to the database.
package puljefordeling

import (
	"database/sql"
	"fmt"
	"sort"
	"strconv"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/service/puljefordeling/solver"
	smodel "github.com/Regncon/conorganizer/service/puljefordeling/solver/model"
)

// AssignedPlayer is one seated participant, with enough context for the UI to
// show how much they wanted this game, whether they carry the DM bump, and
// whether they were bumped off a stronger preference to make room for others.
type AssignedPlayer struct {
	Name   string
	IsDM   bool                 // runs at least one game in the weekend (DM bump)
	Level  models.InterestLevel // their interest in the game they got
	Moved  bool                 // relocated off a higher-scoring event by the solver to make room for others
	Pinned bool                 // manually placed (source=manual); honored by the solver, not chosen by it
}

// EmulatedEvent is the proposed seating for a single event within a pulje.
type EmulatedEvent struct {
	EventID         string
	Title           string
	Capacity        int
	GMName          string           // empty if the event has no GM assigned
	AssignedPlayers []AssignedPlayer // sorted by name
	Undersubscribed bool             // fewer than the solver's viable-player threshold
}

// EmulatedPulje is the proposed seating for one pulje (time slot).
type EmulatedPulje struct {
	PuljeID        models.Pulje
	Name           string
	Events         []EmulatedEvent
	Unassigned     []string // names of interested participants who got no seat
	NewlySatisfied int      // participants who got a top-choice seat this pulje
	TotalScore     int      // sum of actual (unadjusted) interest scores
}

// Emulation is the full preview across all puljer.
type Emulation struct {
	Year           int
	Puljer         []EmulatedPulje
	PlayerCount    int // distinct participants with at least one interest
	SatisfiedTotal int // distinct participants satisfied across the weekend
}

type eligibleEvent struct {
	title    string
	capacity int
}

// EmulateSeatings builds the solver model from the database, runs the
// distribution for every pulje in chronological order, and returns the
// proposed seating. It performs only reads.
func EmulateSeatings(db *sql.DB) (Emulation, error) {
	d, err := loadSeatingData(db)
	if err != nil {
		return Emulation{}, err
	}
	if len(d.puljer) == 0 {
		return Emulation{}, nil
	}

	state, results := d.solveChronological(len(d.puljer) - 1)

	// PlayerCount is distinct participants with at least one interest (unchanged
	// semantics — excludes manual-only placements with no interest).
	emulation := Emulation{Year: d.year, PlayerCount: len(d.prefs)}
	for i := range d.puljer {
		pid := string(d.puljer[i].ID)
		emulation.Puljer = append(emulation.Puljer, shapePulje(
			d.puljer[i], d.weekend.Slots[i], results[i],
			d.gms, d.names, d.prefs, d.dmSet, d.pinnedSet[pid],
		))
	}
	emulation.SatisfiedTotal = state.SatisfiedCount()
	return emulation, nil
}

// shapePulje resolves the solver's ID-based result into display-ready data.
func shapePulje(
	pulje models.PuljeRow,
	slot smodel.Slot,
	res smodel.SlotResult,
	gms map[string]int,
	names map[int]string,
	prefs map[int]map[string]map[string]smodel.Score,
	dmSet map[int]bool,
	pinned map[string]bool,
) EmulatedPulje {
	under := make(map[string]bool, len(res.UndersubscribedEvents))
	for _, eid := range res.UndersubscribedEvents {
		under[eid] = true
	}

	moved := make(map[string]bool, len(res.MovedPlayers))
	for _, pid := range res.MovedPlayers {
		moved[pid] = true
	}

	out := EmulatedPulje{
		PuljeID:        pulje.ID,
		Name:           pulje.Name,
		NewlySatisfied: len(res.NewlySatisfied),
		TotalScore:     res.TotalScore,
		Unassigned:     playerIDsToNames(res.Unassigned, names),
	}

	for _, ev := range slot.Events {
		emEv := EmulatedEvent{
			EventID:         ev.ID,
			Title:           ev.Name,
			Capacity:        ev.Capacity,
			AssignedPlayers: assignedPlayers(res.Assignments[ev.ID], ev.ID, string(pulje.ID), names, prefs, dmSet, moved, pinned),
			Undersubscribed: under[ev.ID],
		}
		if gmID, ok := gms[eventPuljeKey(ev.ID, pulje.ID)]; ok {
			emEv.GMName = names[gmID]
		}
		out.Events = append(out.Events, emEv)
	}

	return out
}

// assignedPlayers turns solver player IDs into display rows: name, DM flag, the
// interest level the player had for the game they were seated in, whether
// the solver relocated them off a higher-scoring event (the moved set), and
// whether the seat was manually placed (pinned).
func assignedPlayers(
	ids []string,
	eventID, puljeID string,
	names map[int]string,
	prefs map[int]map[string]map[string]smodel.Score,
	dmSet map[int]bool,
	moved map[string]bool,
	pinned map[string]bool,
) []AssignedPlayer {
	if len(ids) == 0 {
		return nil
	}
	out := make([]AssignedPlayer, 0, len(ids))
	for _, id := range ids {
		bh, err := strconv.Atoi(id)
		if err != nil {
			out = append(out, AssignedPlayer{Name: id})
			continue
		}
		ap := AssignedPlayer{
			Name:   names[bh],
			IsDM:   dmSet[bh],
			Moved:  moved[id],
			Pinned: pinned[seatKey(eventID, id)],
		}
		if byPulje, ok := prefs[bh]; ok {
			got := byPulje[puljeID][eventID]
			ap.Level = models.InterestLevelFromScore(int(got))
		}
		out = append(out, ap)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

// --- seating data -----------------------------------------------------------

// puljeFrozen reports whether a pulje's seats are committed and no longer
// auto-solved (Locked or Completed).
func puljeFrozen(status models.PuljeStatus) bool {
	return status == models.PuljeStatusLocked || status == models.PuljeStatusCompleted
}

func seatKey(eventID, playerID string) string {
	return eventID + "\x00" + playerID
}

type persistedSeat struct {
	pulje    string
	event    string
	playerID string
	source   string
}

// loadPersistedPlayerSeats returns every role=Player row (manual and solver).
func loadPersistedPlayerSeats(db *sql.DB) ([]persistedSeat, error) {
	const query = `
		SELECT event_id, pulje_id, billettholder_id, source
		FROM relation_events_players
		WHERE role = ?
	`
	rows, err := db.Query(query, models.EventPlayerRolePlayer)
	if err != nil {
		return nil, fmt.Errorf("query player seats: %w", err)
	}
	defer rows.Close()

	var seats []persistedSeat
	for rows.Next() {
		var s persistedSeat
		var bhID int
		if err := rows.Scan(&s.event, &s.pulje, &bhID, &s.source); err != nil {
			return nil, fmt.Errorf("scan player seat: %w", err)
		}
		s.playerID = strconv.Itoa(bhID)
		seats = append(seats, s)
	}
	return seats, rows.Err()
}

// seatingData holds everything needed to run the seeded chronological solve.
type seatingData struct {
	puljer      []models.PuljeRow
	weekend     smodel.Weekend
	players     []smodel.Player
	gms         map[string]int
	names       map[int]string
	prefs       map[int]map[string]map[string]smodel.Score
	dmSet       map[int]bool
	actual      map[string]map[string][]string // puljeID -> eventID -> []playerID (all Player rows)
	manualFixed map[string]map[string]string   // puljeID -> playerID -> eventID (source=manual)
	pinnedSet   map[string]map[string]bool      // puljeID -> seatKey(event,player) (source=manual)
	year        int
}

// loadSeatingData loads puljer, events, GMs, names, prefs, and persisted Player
// seats, and assembles the solver model. Players include every participant with an
// interest plus anyone holding a manual placement (even without an interest).
func loadSeatingData(db *sql.DB) (*seatingData, error) {
	puljer, err := loadPuljer(db)
	if err != nil {
		return nil, err
	}
	if len(puljer) == 0 {
		return &seatingData{}, nil
	}
	events, err := loadEligibleEvents(db)
	if err != nil {
		return nil, err
	}
	gms, err := loadGMs(db)
	if err != nil {
		return nil, err
	}
	names, err := loadParticipantNames(db)
	if err != nil {
		return nil, err
	}
	prefs, err := loadPrefs(db, events)
	if err != nil {
		return nil, err
	}
	seats, err := loadPersistedPlayerSeats(db)
	if err != nil {
		return nil, err
	}

	d := &seatingData{
		puljer:      puljer,
		gms:         gms,
		names:       names,
		prefs:       prefs,
		actual:      make(map[string]map[string][]string),
		manualFixed: make(map[string]map[string]string),
		pinnedSet:   make(map[string]map[string]bool),
	}

	for _, s := range seats {
		if d.actual[s.pulje] == nil {
			d.actual[s.pulje] = make(map[string][]string)
		}
		d.actual[s.pulje][s.event] = append(d.actual[s.pulje][s.event], s.playerID)
		if s.source == models.EventPlayerSourceManual {
			if d.manualFixed[s.pulje] == nil {
				d.manualFixed[s.pulje] = make(map[string]string)
			}
			d.manualFixed[s.pulje][s.playerID] = s.event
			if d.pinnedSet[s.pulje] == nil {
				d.pinnedSet[s.pulje] = make(map[string]bool)
			}
			d.pinnedSet[s.pulje][seatKey(s.event, s.playerID)] = true
		}
	}

	// Build the solver's Weekend in chronological pulje order.
	d.weekend = smodel.Weekend{Slots: make([]smodel.Slot, 0, len(puljer))}
	for _, p := range puljer {
		slot := smodel.Slot{ID: string(p.ID), Name: p.Name}
		for _, eid := range sortedEventIDs(events[p.ID]) {
			e := events[p.ID][eid]
			ev := smodel.Event{ID: eid, Name: e.title, Capacity: e.capacity}
			if gmID, ok := gms[eventPuljeKey(eid, p.ID)]; ok {
				ev.DMID = strconv.Itoa(gmID)
			}
			slot.Events = append(slot.Events, ev)
		}
		d.weekend.Slots = append(d.weekend.Slots, slot)
	}

	// Players: everyone with an interest, plus anyone holding a manual seat.
	playerIDs := make(map[int]bool, len(prefs))
	for bh := range prefs {
		playerIDs[bh] = true
	}
	for _, s := range seats {
		if s.source != models.EventPlayerSourceManual {
			continue
		}
		if bh, err := strconv.Atoi(s.playerID); err == nil {
			playerIDs[bh] = true
		}
	}
	for _, bh := range sortedIntKeys(playerIDs) { // sortedIntKeys is generic over map[int]V
		d.players = append(d.players, smodel.Player{
			ID:    strconv.Itoa(bh),
			Name:  names[bh],
			Prefs: prefs[bh],
		})
	}

	d.dmSet = make(map[int]bool, len(gms))
	for _, bhID := range gms {
		d.dmSet[bhID] = true
	}

	d.year = puljer[0].StartAt.TimeOrZero().Year()
	return d, nil
}

// solveChronological threads one State across slots[0..upTo] inclusive: frozen
// puljer are replayed from their persisted seats (ApplyActual); open puljer are
// solved with their manual placements pinned (SolveSlotFixed). Returns the State
// and per-slot results index-aligned with d.puljer[0..upTo].
func (d *seatingData) solveChronological(upTo int) (*solver.State, []smodel.SlotResult) {
	state := solver.NewState(d.year, d.weekend)
	results := make([]smodel.SlotResult, 0, upTo+1)
	for i := 0; i <= upTo; i++ {
		slot := d.weekend.Slots[i]
		pid := string(d.puljer[i].ID)
		var res smodel.SlotResult
		if puljeFrozen(d.puljer[i].Status) {
			res = state.ApplyActual(slot, d.players, d.actual[pid])
		} else {
			res = state.SolveSlotFixed(slot, d.players, d.manualFixed[pid])
		}
		results = append(results, res)
	}
	return state, results
}

// --- data loading -----------------------------------------------------------

func loadPuljer(db *sql.DB) ([]models.PuljeRow, error) {
	const query = `
		SELECT id, name, status, start_at, end_at
		FROM puljer
		ORDER BY start_at ASC
	`
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("query puljer: %w", err)
	}
	defer rows.Close()

	var puljer []models.PuljeRow
	for rows.Next() {
		var p models.PuljeRow
		if err := rows.Scan(&p.ID, &p.Name, &p.Status, &p.StartAt, &p.EndAt); err != nil {
			return nil, fmt.Errorf("scan pulje row: %w", err)
		}
		puljer = append(puljer, p)
	}
	return puljer, rows.Err()
}

func loadEligibleEvents(db *sql.DB) (map[models.Pulje]map[string]eligibleEvent, error) {
	const query = `
		SELECT ep.pulje_id, e.id, e.title, e.max_players
		FROM relation_event_puljer ep
		JOIN events e ON e.id = ep.event_id
		WHERE ep.is_in_pulje = 1
	`
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("query eligible events: %w", err)
	}
	defer rows.Close()

	out := make(map[models.Pulje]map[string]eligibleEvent)
	for rows.Next() {
		var pulje models.Pulje
		var eventID, title string
		var maxPlayers int
		if err := rows.Scan(&pulje, &eventID, &title, &maxPlayers); err != nil {
			return nil, fmt.Errorf("scan event row: %w", err)
		}
		if out[pulje] == nil {
			out[pulje] = make(map[string]eligibleEvent)
		}
		out[pulje][eventID] = eligibleEvent{title: title, capacity: maxPlayers}
	}
	return out, rows.Err()
}

func loadGMs(db *sql.DB) (map[string]int, error) {
	const query = `
		SELECT event_id, pulje_id, billettholder_id
		FROM relation_events_players
		WHERE role = ?
	`
	rows, err := db.Query(query, models.EventPlayerRoleGM)
	if err != nil {
		return nil, fmt.Errorf("query GMs: %w", err)
	}
	defer rows.Close()

	out := make(map[string]int)
	for rows.Next() {
		var eventID string
		var pulje models.Pulje
		var bhID int
		if err := rows.Scan(&eventID, &pulje, &bhID); err != nil {
			return nil, fmt.Errorf("scan GM row: %w", err)
		}
		out[eventPuljeKey(eventID, pulje)] = bhID
	}
	return out, rows.Err()
}

func loadParticipantNames(db *sql.DB) (map[int]string, error) {
	const query = `SELECT id, first_name, last_name FROM billettholdere`
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("query participants: %w", err)
	}
	defer rows.Close()

	out := make(map[int]string)
	for rows.Next() {
		var id int
		var first, last string
		if err := rows.Scan(&id, &first, &last); err != nil {
			return nil, fmt.Errorf("scan participant row: %w", err)
		}
		out[id] = first + " " + last
	}
	return out, rows.Err()
}

// loadPrefs builds billettholderID -> puljeID -> eventID -> score, keeping only
// interests in events that are actually placed in that pulje and only positive
// scores (an edge in the assignment graph).
func loadPrefs(
	db *sql.DB,
	events map[models.Pulje]map[string]eligibleEvent,
) (map[int]map[string]map[string]smodel.Score, error) {
	const query = `SELECT billettholder_id, event_id, pulje_id, interest_level FROM interests`
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("query interests: %w", err)
	}
	defer rows.Close()

	out := make(map[int]map[string]map[string]smodel.Score)
	for rows.Next() {
		var bhID int
		var eventID string
		var pulje models.Pulje
		var level models.InterestLevel
		if err := rows.Scan(&bhID, &eventID, &pulje, &level); err != nil {
			return nil, fmt.Errorf("scan interest row: %w", err)
		}

		// Skip interests for events not placed in this pulje, so the solver
		// only sees real, seatable choices.
		if _, ok := events[pulje][eventID]; !ok {
			continue
		}
		score := level.Score()
		if score == 0 {
			continue
		}

		byPulje := out[bhID]
		if byPulje == nil {
			byPulje = make(map[string]map[string]smodel.Score)
			out[bhID] = byPulje
		}
		byEvent := byPulje[string(pulje)]
		if byEvent == nil {
			byEvent = make(map[string]smodel.Score)
			byPulje[string(pulje)] = byEvent
		}
		byEvent[eventID] = smodel.Score(score)
	}
	return out, rows.Err()
}

// --- helpers -----------------------------------------------------------------

func eventPuljeKey(eventID string, pulje models.Pulje) string {
	return eventID + "\x00" + string(pulje)
}

func playerIDsToNames(ids []string, names map[int]string) []string {
	if len(ids) == 0 {
		return nil
	}
	out := make([]string, 0, len(ids))
	for _, id := range ids {
		if n, err := strconv.Atoi(id); err == nil {
			if name, ok := names[n]; ok {
				out = append(out, name)
				continue
			}
		}
		out = append(out, id)
	}
	sort.Strings(out)
	return out
}

func sortedEventIDs(m map[string]eligibleEvent) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Slice(out, func(i, j int) bool {
		// Sort by title for stable, readable output; fall back to ID.
		if m[out[i]].title != m[out[j]].title {
			return m[out[i]].title < m[out[j]].title
		}
		return out[i] < out[j]
	})
	return out
}

func sortedIntKeys[V any](m map[int]V) []int {
	out := make([]int, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Ints(out)
	return out
}
