package eventservice

import (
	"database/sql"
	"fmt"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"

	"github.com/Regncon/conorganizer/models"
	puljerService "github.com/Regncon/conorganizer/service/puljer"
)

func MdToHTML(md []byte) []byte {
	// create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}

func GetEventById(eventID string, db *sql.DB) (*models.Event, error) {
	query := `
            SELECT
                id,
                title,
                intro,
                description,
                system,
                event_type,
                age_group,
                event_runtime,
                host_name,
                user_id,
                email,
                phone_number,
                max_players,
                beginner_friendly,
                can_be_run_in_english,
                notes,
                status,
                created_at,
                updated_at,
                created_by_id,
                updated_by_id,
                status_changed_by_id,
                status_changed_at,
                status_changed_action
            FROM events WHERE id = ?
            `
	row := db.QueryRow(query, eventID)

	var event models.Event
	if err := row.Scan(
		&event.ID,
		&event.Title,
		&event.Intro,
		&event.Description,
		&event.System,
		&event.EventType,
		&event.AgeGroup,
		&event.Runtime,
		&event.HostName,
		&event.UserID,
		&event.Email,
		&event.PhoneNumber,
		&event.MaxPlayers,
		&event.BeginnerFriendly,
		&event.CanBeRunInEnglish,
		&event.Notes,
		&event.Status,
		&event.CreatedAt,
		&event.UpdatedAt,
		&event.CreatedByID,
		&event.UpdatedByID,
		&event.StatusChangedByID,
		&event.StatusChangedAt,
		&event.StatusChangedAction,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No event found
		}
		return nil, fmt.Errorf("failed to scan event row for event %q: %w", eventID, err)
	}

	return &event, nil
}

func GetPujerForEvent(
	eventID string,
	db *sql.DB,
) ([]models.PuljeRow, error) {
	puljer, err := puljerService.GetActivePuljeForEvent(eventID, db)
	if err != nil {
		return nil, fmt.Errorf("error querying puljer for event %q: %w", eventID, err)
	}
	return puljer, nil
}
