package event

import (
	"database/sql"
	"fmt"
	"strings"

	ticketholder "github.com/Regncon/conorganizer/components/ticket_holder"
	"github.com/Regncon/conorganizer/models"
)

type ProgramInterestNotice struct {
	PuljeName   string
	EventTitles []string
}

func (notice ProgramInterestNotice) HasInterest() bool {
	return len(notice.EventTitles) > 0
}

func canShowInterestControls(programPublished bool, isInPuljefordeling bool, puljerForEvent []models.PuljeRow) bool {
	return programPublished && isInPuljefordeling && len(puljerForEvent) > 0
}

func getProgramInterestNotice(db *sql.DB, ticketHolders []ticketholder.BillettHolder, puljeID string) (ProgramInterestNotice, error) {
	notice := ProgramInterestNotice{}
	if len(ticketHolders) == 0 || puljeID == "" {
		return notice, nil
	}

	placeholders := strings.TrimSuffix(strings.Repeat("?,", len(ticketHolders)), ",")
	args := make([]any, 0, len(ticketHolders)+6)
	for _, ticketHolder := range ticketHolders {
		args = append(args, ticketHolder.Id)
	}
	args = append(args,
		puljeID,
		models.InterestLevelHigh,
		models.InterestLevelMedium,
		models.InterestLevelLow,
		models.EventStatusAnnounced,
		models.PuljeStatusCompleted,
	)

	query := fmt.Sprintf(`
		SELECT DISTINCT p.name, e.title
		FROM interests i
		JOIN events e ON e.id = i.event_id
		JOIN relation_event_puljer ep
			ON ep.event_id = i.event_id
			AND ep.pulje_id = i.pulje_id
		JOIN puljer p ON p.id = i.pulje_id
		WHERE i.billettholder_id IN (%s)
			AND i.pulje_id = ?
			AND i.interest_level IN (?, ?, ?)
			AND e.is_in_puljefordeling = 1
			AND e.status = ?
			AND ep.is_in_pulje = 1
			AND p.status != ?
		ORDER BY e.title ASC
	`, placeholders)

	rows, err := db.Query(query, args...)
	if err != nil {
		return notice, fmt.Errorf("query program interest notice for pulje %s: %w", puljeID, err)
	}
	defer rows.Close()

	for rows.Next() {
		var puljeName, eventTitle string
		if err := rows.Scan(&puljeName, &eventTitle); err != nil {
			return notice, fmt.Errorf("scan program interest notice for pulje %s: %w", puljeID, err)
		}

		if notice.PuljeName == "" {
			notice.PuljeName = puljeName
		}
		notice.EventTitles = append(notice.EventTitles, eventTitle)
	}
	if err := rows.Err(); err != nil {
		return notice, fmt.Errorf("iterate program interest notice for pulje %s: %w", puljeID, err)
	}

	return notice, nil
}
