package admin

import (
	"database/sql"
	"fmt"
	"log/slog"
	"sort"
	"strings"

	"github.com/Regncon/conorganizer/models"
	puljerService "github.com/Regncon/conorganizer/service/puljer"
)

// BoardGame is one game card on the puljeoppsett board.
type BoardGame struct {
	EventID  string
	Title    string
	Status   models.EventStatus
	AgeGroup models.AgeGroup
	Beginner bool
	OwnerKey string // "user:<id>" or "email:<lowercased>" — DM identity
	HostName string
}

// SlotStats are the per-column counters shown in the header.
type SlotStats struct {
	Games    int
	Adults   int // age_group == AdultsOnly
	Beginner int // beginner_friendly
}

// DMCollision is one owner who runs >=2 games in the same pulje.
type DMCollision struct {
	HostName string
	Count    int
}

// SlotColumn is one pulje column: its games, stats and collisions.
type SlotColumn struct {
	Pulje      models.PuljeRow
	Games      []BoardGame
	Stats      SlotStats
	Collisions []DMCollision
}

// ScheduleBoard is the whole board: pool + one column per pulje.
type ScheduleBoard struct {
	Pool           []BoardGame
	Columns        []SlotColumn
	CollisionCount int
}

func ownerKey(userID sql.NullInt64, email string) string {
	if userID.Valid {
		return fmt.Sprintf("user:%d", userID.Int64)
	}
	return "email:" + strings.ToLower(strings.TrimSpace(email))
}

// buildScheduleBoard loads eligible games (Godkjent/Annonsert), shapes them into
// a pool (games in no pulje) plus one column per pulje (ordered by start_at), and
// computes per-slot stats and DM double-booking collisions.
func buildScheduleBoard(db *sql.DB, logger *slog.Logger) (ScheduleBoard, error) {
	puljer, err := puljerService.GetAllPuljer(db)
	if err != nil {
		return ScheduleBoard{}, fmt.Errorf("load puljer: %w", err)
	}

	const query = `
		SELECT e.id, e.title, e.status, e.age_group, e.beginner_friendly,
		       e.user_id, e.email, e.host_name, ep.pulje_id
		FROM events e
		LEFT JOIN relation_event_puljer ep
		       ON ep.event_id = e.id AND ep.is_in_pulje = 1
		WHERE e.status IN ('Godkjent', 'Annonsert')
		ORDER BY e.title ASC
	`
	rows, err := db.Query(query)
	if err != nil {
		return ScheduleBoard{}, fmt.Errorf("query board games: %w", err)
	}
	defer rows.Close()

	gamesByPulje := make(map[models.Pulje][]BoardGame)
	poolByID := make(map[string]BoardGame)
	slotted := make(map[string]bool)

	for rows.Next() {
		var (
			g       BoardGame
			userID  sql.NullInt64
			email   string
			puljeID sql.NullString
			beginI  int
		)
		if err := rows.Scan(&g.EventID, &g.Title, &g.Status, &g.AgeGroup, &beginI,
			&userID, &email, &g.HostName, &puljeID); err != nil {
			return ScheduleBoard{}, fmt.Errorf("scan board game: %w", err)
		}
		g.Beginner = beginI == 1
		g.OwnerKey = ownerKey(userID, email)

		poolByID[g.EventID] = g
		if puljeID.Valid {
			p := models.Pulje(puljeID.String)
			gamesByPulje[p] = append(gamesByPulje[p], g)
			slotted[g.EventID] = true
		}
	}
	if err := rows.Err(); err != nil {
		return ScheduleBoard{}, fmt.Errorf("iterate board games: %w", err)
	}

	board := ScheduleBoard{}
	for _, p := range puljer {
		col := SlotColumn{Pulje: p, Games: gamesByPulje[p.ID]}
		col.Stats = computeStats(col.Games)
		col.Collisions = computeCollisions(col.Games)
		board.CollisionCount += len(col.Collisions)
		board.Columns = append(board.Columns, col)
	}
	board.Pool = poolGames(poolByID, slotted)

	return board, nil
}

func computeStats(games []BoardGame) SlotStats {
	s := SlotStats{Games: len(games)}
	for _, g := range games {
		if g.AgeGroup == models.AgeGroupAdultsOnly {
			s.Adults++
		}
		if g.Beginner {
			s.Beginner++
		}
	}
	return s
}

func computeCollisions(games []BoardGame) []DMCollision {
	order := []string{}
	counts := map[string]int{}
	host := map[string]string{}
	for _, g := range games {
		if _, seen := counts[g.OwnerKey]; !seen {
			order = append(order, g.OwnerKey)
			host[g.OwnerKey] = g.HostName
		}
		counts[g.OwnerKey]++
	}
	var out []DMCollision
	for _, key := range order {
		if counts[key] >= 2 {
			out = append(out, DMCollision{HostName: host[key], Count: counts[key]})
		}
	}
	return out
}

func poolGames(byID map[string]BoardGame, slotted map[string]bool) []BoardGame {
	var out []BoardGame
	for id, g := range byID {
		if !slotted[id] {
			out = append(out, g)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Title < out[j].Title })
	return out
}
