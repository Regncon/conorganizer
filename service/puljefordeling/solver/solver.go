package solver // solver is defined in flow.go

import (
	"math/rand/v2"
	"sort"

	"github.com/Regncon/conorganizer/service/puljefordeling/solver/model"
)

// Weight bands for the assignment graph. The solver maximises total weight
// among the assignments it makes, so these constants encode the priority order
// we want. Bands are spaced by 200, and every within-band bump
// (miss / never-seated / DM) sums to less than that — so the declared order
// below can never be broken by a bump.
//
// Declared priority, highest first (bumps break ties within a band):
//
//	1. unsatisfied + top choice (Veldig)   — the satisfaction goal
//	2. satisfied   + top choice
//	3. medium interest (Middels)           — satisfied or not, same band
//	4. low interest    (Litt)              — satisfied or not, same band
//
// The unsatisfied advantage exists ONLY on the top choice (the satisfaction
// goal); Middels/Litt are valued the same whether or not the player is
// satisfied, which shrinks the surface for "stay unsatisfied" gaming.
const (
	bandUnsatVeldig = 800
	bandSatVeldig   = 600
	bandMiddels     = 400
	bandLitt        = 200
)

// Within-band bumps (all small relative to the 200 band gap).
const (
	// dmBump rewards a player who runs at least one event in the weekend. It
	// stays within the band, so a regular player's unmet top choice (≥800)
	// always outranks a DM's medium/low interest (≤470).
	dmBump = 60

	// neverSeatedBump nudges a player who has not been given any seat yet this
	// weekend. Small and self-limiting (to stay never-seated you must forgo
	// seats you'd have wanted), so it is not worth gaming.
	neverSeatedBump = 10

	// missStep / missCap implement the (un-gameable) scarcity bonus: priority
	// for a top choice grows with the number of prior puljer in which a player
	// wanted a top choice but did not get one. It is backward-looking on locked
	// results, so it cannot be farmed by concentrating declarations.
	missStep = 20
	missCap  = 60
)

// participationBonus prices "one more attendee seated" in the same units as the
// weights. It is added to every assignment edge; the flow stops filling chairs
// once the next seat would cost more than this in total preference (see
// flowGraph.minCostFlow). Raise toward infinity to recover "fill every chair";
// lower toward zero for pure welfare-maximisation.
const participationBonus = 300

// minViablePlayers is the threshold below which an event is flagged for
// organiser review (not cancelled). Three players is enough for most
// tabletop games to be worth running; below that the organiser should look
// at the assignment and decide.
const minViablePlayers = 3

// State carries fairness data forward across slots:
//   - satisfied: has received a top-choice (score-5) assignment
//   - seated:    has received any assignment at all this weekend
//   - misses:    number of prior puljer the player wanted a top choice and
//     missed (drives the scarcity bonus)
type State struct {
	satisfied map[string]struct{}
	seated    map[string]struct{}
	misses    map[string]int
	year      int
	slotIndex int

	// dmSlots[playerID] is the set of slot IDs in which playerID is DMing
	// an event. Players in this set are excluded from being assigned during
	// those slots.
	dmSlots map[string]map[string]struct{}

	// isDM[playerID] is true if playerID DMs at least one event anywhere
	// in the weekend. Such players receive a fixed bump on every edge.
	isDM map[string]struct{}
}

// NewState returns a fresh State for the start of a weekend.
//
// year is used to derive per-slot tie-breaking seeds (seed = year×1000 + slotIndex),
// making results deterministic within a year but different across years.
//
// weekend provides the full slot schedule so the solver can recognise which
// players DM events in which slots and apply the DM bump across the weekend.
func NewState(year int, weekend model.Weekend) *State {
	dmSlots := make(map[string]map[string]struct{})
	isDM := make(map[string]struct{})
	for _, sl := range weekend.Slots {
		for _, ev := range sl.Events {
			if ev.DMID == "" {
				continue
			}
			isDM[ev.DMID] = struct{}{}
			if dmSlots[ev.DMID] == nil {
				dmSlots[ev.DMID] = make(map[string]struct{})
			}
			dmSlots[ev.DMID][sl.ID] = struct{}{}
		}
	}
	return &State{
		satisfied: make(map[string]struct{}),
		seated:    make(map[string]struct{}),
		misses:    make(map[string]int),
		year:      year,
		dmSlots:   dmSlots,
		isDM:      isDM,
	}
}

// IsSatisfied reports whether playerID has received a score-5 assignment.
func (s *State) IsSatisfied(playerID string) bool {
	_, ok := s.satisfied[playerID]
	return ok
}

// SatisfiedCount returns the number of satisfied players.
func (s *State) SatisfiedCount() int {
	return len(s.satisfied)
}

// IsDM reports whether playerID runs at least one event in the weekend.
func (s *State) IsDM(playerID string) bool {
	_, ok := s.isDM[playerID]
	return ok
}

// SolveSlot assigns players to events for one slot, updates the fairness
// state, and returns the result.
//
// Players who are DMing any event in this slot are excluded from the player
// pool — they cannot also be assigned as participants in the same slot.
//
// Participation is priced (see participationBonus): a chair is filled only when
// doing so is worth the welfare it costs, so an interested player may be left
// unseated rather than dragged off a strong preference to fill a seat.
//
// Events whose minViablePlayers threshold is not met are NOT cancelled — they
// are reported in UndersubscribedEvents for the organiser to review manually.
func (s *State) SolveSlot(slot model.Slot, players []model.Player) model.SlotResult {
	currentIndex := s.slotIndex
	seed := int64(s.year)*1000 + int64(currentIndex)
	s.slotIndex++

	result := model.SlotResult{
		SlotID:      slot.ID,
		Assignments: make(map[string][]string),
		Seed:        seed,
	}

	// Players DMing in this slot are unavailable as players.
	dmingHere := make(map[string]bool)
	for _, ev := range slot.Events {
		if ev.DMID != "" {
			dmingHere[ev.DMID] = true
		}
	}

	// Only consider players who expressed interest and are not DMing here.
	interested := make([]model.Player, 0, len(players))
	for _, p := range players {
		if dmingHere[p.ID] {
			continue
		}
		if len(p.Prefs[slot.ID]) > 0 {
			interested = append(interested, p)
		}
	}
	if len(interested) == 0 {
		return result
	}

	// Shuffle players before building the graph so that equal-weight ties are
	// broken randomly but reproducibly.
	rng := rand.New(rand.NewPCG(uint64(seed), 0)) //nolint:gosec
	rng.Shuffle(len(interested), func(i, j int) {
		interested[i], interested[j] = interested[j], interested[i]
	})

	// Index players by ID so score lookups during the totals pass are O(1)
	// rather than a linear scan per assignment.
	playerByID := make(map[string]model.Player, len(interested))
	for _, p := range interested {
		playerByID[p.ID] = p
	}

	assignments, moved := s.runMCMF(slot.ID, slot.Events, interested)
	result.Assignments = assignments

	// Flag events with fewer than minViablePlayers for human review.
	for _, ev := range slot.Events {
		if len(assignments[ev.ID]) < minViablePlayers {
			result.UndersubscribedEvents = append(result.UndersubscribedEvents, ev.ID)
		}
	}

	// Update satisfaction/seated state and collect totals.
	assigned := make(map[string]bool, len(interested))
	for evID, playerIDs := range assignments {
		for _, pid := range playerIDs {
			score := playerByID[pid].Prefs[slot.ID][evID]
			result.TotalScore += int(score)
			s.seated[pid] = struct{}{}
			if _, ok := s.satisfied[pid]; score == model.MaxScore && !ok {
				s.satisfied[pid] = struct{}{}
				result.NewlySatisfied = append(result.NewlySatisfied, pid)
			}
			assigned[pid] = true
		}
	}

	// Players bumped off a higher-scoring event by a residual-edge augmentation
	// and re-seated elsewhere to improve the global result. Restricted to players
	// who still hold a seat, so it always describes a relocation, never a player
	// who merely lost contention without ever being tentatively seated.
	for pid := range moved {
		if assigned[pid] {
			result.MovedPlayers = append(result.MovedPlayers, pid)
		}
	}

	// Record a "miss" for every still-unsatisfied player who wanted a top
	// choice this slot but did not get one — this drives the scarcity bonus in
	// future slots. It is backward-looking on the locked result, so it cannot
	// be gamed by how a player declares future interests.
	for _, p := range interested {
		if _, ok := s.satisfied[p.ID]; ok {
			continue
		}
		if wantedTopChoice(p, slot.ID) {
			s.misses[p.ID]++
		}
	}

	// Collect unassigned players (had interest but no seat).
	for _, p := range interested {
		if !assigned[p.ID] {
			result.Unassigned = append(result.Unassigned, p.ID)
		}
	}

	// Sort all output slices for deterministic results.
	sort.Strings(result.NewlySatisfied)
	sort.Strings(result.Unassigned)
	sort.Strings(result.UndersubscribedEvents)
	sort.Strings(result.MovedPlayers)
	for evID := range result.Assignments {
		sort.Strings(result.Assignments[evID])
	}

	return result
}

// wantedTopChoice reports whether the player rated any event in this slot as a
// top choice (score 5).
func wantedTopChoice(p model.Player, slotID string) bool {
	for _, score := range p.Prefs[slotID] {
		if score == model.MaxScore {
			return true
		}
	}
	return false
}

// runMCMF builds and solves the flow network for the given events and players,
// returning the raw assignment map (eventID -> []playerID, unsorted) and the set
// of player IDs that were bumped off an event by a residual-edge augmentation.
func (s *State) runMCMF(
	slotID string,
	events []model.Event,
	players []model.Player,
) (map[string][]string, map[string]bool) {
	assignments := make(map[string][]string)
	if len(events) == 0 {
		return assignments, nil
	}

	// Node layout:
	//   0           → source
	//   1 .. P      → one node per player
	//   P+1 .. P+E  → one node per event
	//   P+E+1       → sink
	P := len(players)
	E := len(events)
	source := 0
	sink := P + E + 1
	g := newFlowGraph(sink + 1)

	for i := range players {
		g.addEdge(source, i+1, 1, 0)
	}

	for j, ev := range events {
		g.addEdge(P+1+j, sink, ev.Capacity, 0)
	}

	// Iterate events (slice, deterministic) for each player rather than the
	// player's preference map so the edge addition order is identical
	// run-to-run.
	for i, p := range players {
		for j, ev := range events {
			score, ok := p.Prefs[slotID][ev.ID]
			if !ok {
				continue
			}
			_, satisfied := s.satisfied[p.ID]
			_, seated := s.seated[p.ID]
			_, isDM := s.isDM[p.ID]
			w := adjustScore(
				score,
				satisfied,
				!seated,
				isDM,
				s.misses[p.ID],
			)
			// Cost is negated (we minimise cost = maximise weight). The
			// participation bonus is folded into every assignment edge so the
			// flow stops once a new seat would cost more than it is worth.
			g.addEdge(i+1, P+1+j, 1, -(w + participationBonus))
		}
	}

	_, _, reduced := g.minCostFlow(source, sink)

	// A reduced forward edge ran player→event (forward at fe, its reverse at
	// fe^1 runs event→player, so its .to is the player node). Flow pushed back
	// along it means that player was bumped off the event.
	moved := make(map[string]bool)
	for _, fe := range reduced {
		playerNode := g.edges[fe^1].to
		eventNode := g.edges[fe].to
		if playerNode < 1 || playerNode > P || eventNode < P+1 || eventNode > P+E {
			continue
		}
		moved[players[playerNode-1].ID] = true
	}

	for i, p := range players {
		for _, eid := range g.adj[i+1] {
			e := g.edges[eid]
			if e.flow != 1 || e.to < P+1 || e.to > P+E {
				continue
			}
			evID := events[e.to-P-1].ID
			assignments[evID] = append(assignments[evID], p.ID)
		}
	}

	return assignments, moved
}

// adjustScore returns the priority weight for a (player, event) edge. Larger =
// higher priority for that seat. See the band constants for the declared order.
//
//   - The unsatisfied advantage applies only to the top choice (Veldig).
//   - The scarcity (miss) bonus and never-seated bump apply only while the
//     player is unsatisfied.
//   - The DM bump applies to every edge but stays within its band.
func adjustScore(score model.Score, satisfied, neverSeated, isDM bool, misses int) int {
	w := bandBase(score, satisfied)

	if !satisfied && score == model.MaxScore {
		bonus := misses * missStep
		if bonus > missCap {
			bonus = missCap
		}
		w += bonus
	}
	if !satisfied && neverSeated {
		w += neverSeatedBump
	}
	if isDM {
		w += dmBump
	}
	return w
}

// bandBase returns the category base weight for an edge.
func bandBase(score model.Score, satisfied bool) int {
	switch {
	case score == model.MaxScore && !satisfied:
		return bandUnsatVeldig
	case score == model.MaxScore:
		return bandSatVeldig
	case score >= 3:
		return bandMiddels
	default:
		return bandLitt
	}
}
