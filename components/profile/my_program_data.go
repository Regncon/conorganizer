package profilecomponent

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/Regncon/conorganizer/models"
)

type UserEvent struct {
	EventID      string
	Title        string
	Intro        string
	Description  string
	UserID       int
	EventType    models.EventType
	PuljeID      models.Pulje
	PuljeStatus  models.PuljeStatus
	IsGM         bool
	RoomName     string
	RoomNumber   string
	GMNames      []string
	GroupMembers []ProgramGroupMember
}

type ProgramGroupMember struct {
	FirstName string
	Role      models.EventPlayerRole
	IsPlayer  bool
	IsGM      bool
}

type UserProgramPulje struct {
	Pulje     models.PuljeRow
	Events    []UserEvent
	Interests []UserInterest
}

type UserInterest struct {
	EventID       string
	EventName     string
	InterestLevel models.InterestLevel
	PuljeID       models.Pulje
}

func (event UserEvent) RoomLabel() string {
	switch {
	case event.RoomNumber != "" && event.RoomName != "":
		return fmt.Sprintf("Rom %s · %s", event.RoomNumber, event.RoomName)
	case event.RoomNumber != "":
		return "Rom " + event.RoomNumber
	case event.RoomName != "":
		return event.RoomName
	default:
		return "Rom kommer"
	}
}

func (event UserEvent) GMLabel() string {
	return strings.Join(event.GMNames, ", ")
}

func (member ProgramGroupMember) RoleLabel() string {
	if member.IsPlayer && member.IsGM {
		return "Spiller og GM"
	}
	if member.IsGM || member.Role == models.EventPlayerRoleGM {
		return "GM"
	}
	return ""
}

func enrichUserProgramEvents(db *sql.DB, events []UserEvent) error {
	for i := range events {
		if err := loadUserProgramEventDetails(db, &events[i]); err != nil {
			return err
		}
	}

	return nil
}

func loadUserProgramEventDetails(db *sql.DB, event *UserEvent) error {
	const roomQuery = `
		SELECT COALESCE(r.room_number, ''), COALESCE(r.name, '')
		FROM relation_event_puljer event_pulje
		LEFT JOIN rooms r ON r.id = event_pulje.room_id
		WHERE event_pulje.event_id = ? AND event_pulje.pulje_id = ?
	`
	if err := db.QueryRow(roomQuery, event.EventID, event.PuljeID).Scan(&event.RoomNumber, &event.RoomName); err != nil {
		return fmt.Errorf("query room for event %s in pulje %s: %w", event.EventID, event.PuljeID, err)
	}

	const gmQuery = `
		SELECT b.first_name || ' ' || b.last_name
		FROM relation_events_players player
		JOIN billettholdere b ON b.id = player.billettholder_id
		WHERE player.event_id = ?
			AND player.pulje_id = ?
			AND player.role = ?
		ORDER BY b.first_name, b.last_name, b.id
	`
	gmRows, err := db.Query(gmQuery, event.EventID, event.PuljeID, models.EventPlayerRoleGM)
	if err != nil {
		return fmt.Errorf("query GMs for event %s in pulje %s: %w", event.EventID, event.PuljeID, err)
	}
	for gmRows.Next() {
		var gmName string
		if err := gmRows.Scan(&gmName); err != nil {
			gmRows.Close()
			return fmt.Errorf("scan GM for event %s in pulje %s: %w", event.EventID, event.PuljeID, err)
		}
		event.GMNames = append(event.GMNames, gmName)
	}
	if err := gmRows.Err(); err != nil {
		gmRows.Close()
		return fmt.Errorf("iterate GMs for event %s in pulje %s: %w", event.EventID, event.PuljeID, err)
	}
	if err := gmRows.Close(); err != nil {
		return fmt.Errorf("close GMs for event %s in pulje %s: %w", event.EventID, event.PuljeID, err)
	}

	if event.PuljeStatus != models.PuljeStatusCompleted {
		return nil
	}

	const groupQuery = `
		SELECT
			b.first_name,
			MAX(player.role = ?) AS is_player,
			MAX(player.role = ?) AS is_gm
		FROM relation_events_players player
		JOIN billettholdere b ON b.id = player.billettholder_id
		WHERE player.event_id = ? AND player.pulje_id = ?
		GROUP BY b.id, b.first_name, b.last_name
		ORDER BY
			MAX(player.role = ?) DESC,
			b.first_name,
			b.last_name,
			b.id
	`
	groupRows, err := db.Query(
		groupQuery,
		models.EventPlayerRolePlayer,
		models.EventPlayerRoleGM,
		event.EventID,
		event.PuljeID,
		models.EventPlayerRoleGM,
	)
	if err != nil {
		return fmt.Errorf("query group for event %s in pulje %s: %w", event.EventID, event.PuljeID, err)
	}
	defer groupRows.Close()

	for groupRows.Next() {
		var member ProgramGroupMember
		if err := groupRows.Scan(&member.FirstName, &member.IsPlayer, &member.IsGM); err != nil {
			return fmt.Errorf("scan group member for event %s in pulje %s: %w", event.EventID, event.PuljeID, err)
		}
		if member.IsGM {
			member.Role = models.EventPlayerRoleGM
		} else if member.IsPlayer {
			member.Role = models.EventPlayerRolePlayer
		}
		event.GroupMembers = append(event.GroupMembers, member)
	}
	if err := groupRows.Err(); err != nil {
		return fmt.Errorf("iterate group for event %s in pulje %s: %w", event.EventID, event.PuljeID, err)
	}

	return nil
}
