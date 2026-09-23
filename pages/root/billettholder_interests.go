package root

import (
	"cmp"
	"database/sql"
	"fmt"
	"slices"
	"strings"

	"github.com/Regncon/conorganizer/components"
	ticketholder "github.com/Regncon/conorganizer/components/ticket_holder"
	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/service/requestctx"
)

type billettholderInterestsByEvent map[models.Pulje]map[string][]components.BillettholderInterest

func (interests billettholderInterestsByEvent) forEvent(puljeID models.Pulje, eventID string) []components.BillettholderInterest {
	return interests[puljeID][eventID]
}

func loadBillettholderInterests(userInfo requestctx.UserRequestInfo, selectedBillettholderHint int, db *sql.DB) (billettholderInterestsByEvent, error) {
	if userInfo.Email == "" {
		return nil, nil
	}

	associated, err := ticketholder.GetTicketHolders(userInfo, db)
	if err != nil {
		return nil, fmt.Errorf("failed to load billettholdere for interest indicator: %w", err)
	}
	if len(associated) == 0 {
		return nil, nil
	}

	selectedID := ticketholder.ResolveSelectedBillettholderID(userInfo, associated, selectedBillettholderHint)
	billettholdere := make(map[int]ticketholder.BillettHolder, len(associated))
	for _, billettholder := range associated {
		billettholdere[billettholder.Id] = billettholder
	}

	placeholders := make([]string, 0, len(billettholdere))
	args := make([]any, 0, len(billettholdere))
	for id := range billettholdere {
		placeholders = append(placeholders, "?")
		args = append(args, id)
	}

	query := fmt.Sprintf(`
		SELECT billettholder_id, event_id, pulje_id, interest_level
		FROM interests
		WHERE billettholder_id IN (%s)
	`, strings.Join(placeholders, ", "))
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query billettholder interests: %w", err)
	}
	defer rows.Close()

	interests := make(billettholderInterestsByEvent)
	for rows.Next() {
		var billettholderID int
		var eventID, puljeID string
		var level models.InterestLevel
		if err := rows.Scan(&billettholderID, &eventID, &puljeID, &level); err != nil {
			return nil, fmt.Errorf("failed to scan billettholder interest: %w", err)
		}
		if level.Score() == 0 {
			continue
		}

		pulje := models.Pulje(puljeID)
		if interests[pulje] == nil {
			interests[pulje] = make(map[string][]components.BillettholderInterest)
		}
		interests[pulje][eventID] = append(interests[pulje][eventID], components.BillettholderInterest{
			BillettholderID:   billettholderID,
			BillettholderName: billettholdere[billettholderID].Name,
			InterestLevel:     level,
			IsSelected:        billettholderID == selectedID,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate billettholder interests: %w", err)
	}

	for _, interestsByEvent := range interests {
		for _, eventInterests := range interestsByEvent {
			slices.SortFunc(eventInterests, compareBillettholderInterests)
		}
	}

	return interests, nil
}

func compareBillettholderInterests(a, b components.BillettholderInterest) int {
	if a.IsSelected != b.IsSelected {
		if a.IsSelected {
			return -1
		}
		return 1
	}
	return cmp.Or(
		cmp.Compare(b.InterestLevel.Score(), a.InterestLevel.Score()),
		cmp.Compare(a.BillettholderName, b.BillettholderName),
	)
}
