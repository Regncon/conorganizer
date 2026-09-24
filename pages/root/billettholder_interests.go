package root

import (
	"cmp"
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"slices"
	"strings"

	"github.com/Regncon/conorganizer/components"
	ticketholder "github.com/Regncon/conorganizer/components/ticket_holder"
	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/service/requestctx"
	"github.com/Regncon/conorganizer/service/userctx"
)

type billettholderInterestsByEvent map[string][]components.BillettholderInterest

func loadPuljeInterests(ctx context.Context, db *sql.DB, puljeID models.Pulje) billettholderInterestsByEvent {
	logger := slog.Default().With("component", "root")
	userInfo := userctx.GetUserRequestInfo(ctx)
	interests, err := loadBillettholderInterests(userInfo, requestctx.SelectedBillettholderID(ctx), puljeID, db)
	if err != nil {
		logger.Error(err.Error(), "user_id", userInfo.Id, "pulje_id", puljeID)
		return nil
	}
	return interests
}

func loadBillettholderInterests(userInfo requestctx.UserRequestInfo, selectedBillettholderHint int, puljeID models.Pulje, db *sql.DB) (billettholderInterestsByEvent, error) {
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
	args := []any{puljeID}
	for id := range billettholdere {
		placeholders = append(placeholders, "?")
		args = append(args, id)
	}

	query := fmt.Sprintf(`
		SELECT billettholder_id, event_id, interest_level
		FROM interests
		WHERE pulje_id = ? AND billettholder_id IN (%s)
	`, strings.Join(placeholders, ", "))
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query billettholder interests: %w", err)
	}
	defer rows.Close()

	interests := make(billettholderInterestsByEvent)
	for rows.Next() {
		var billettholderID int
		var eventID string
		var level models.InterestLevel
		if err := rows.Scan(&billettholderID, &eventID, &level); err != nil {
			return nil, fmt.Errorf("failed to scan billettholder interest: %w", err)
		}
		if level.Score() == 0 {
			continue
		}

		interests[eventID] = append(interests[eventID], components.BillettholderInterest{
			BillettholderID:   billettholderID,
			BillettholderName: billettholdere[billettholderID].Name,
			InterestLevel:     level,
			IsSelected:        billettholderID == selectedID,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate billettholder interests: %w", err)
	}

	for _, eventInterests := range interests {
		slices.SortFunc(eventInterests, compareBillettholderInterests)
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
