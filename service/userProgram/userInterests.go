package userProgram

import (
	"database/sql"
	"log/slog"

	billettholderService "github.com/Regncon/conorganizer/service/billettholder"
	"github.com/Regncon/conorganizer/service/requestctx"
)

type UserInterest struct {
	EventID       string
	EventName     string
	InterestLevel string
	PuljeID       string
}

func GetAllInterestsForUser(userInfo requestctx.UserRequestInfo, db *sql.DB, logger *slog.Logger) ([]UserInterest, error) {
	logger.Info("Fetching interests for user", "userId", userInfo.Id)

	billettholdere, billettholderErr := billettholderService.GetBilettholdere(userInfo.Id, db, logger)
	if billettholderErr != nil {
		logger.Error("Failed to get billettholdere", "error", billettholderErr)
		return nil, billettholderErr
	}

	if len(billettholdere) == 0 {
		logger.Info("User has no billettholdere")
		return []UserInterest{}, nil
	}

	billettholderID := billettholdere[0].ID

	query := `
		SELECT
			i.event_id,
			e.title,
			i.interest_level,
			i.pulje_id
		FROM interests i
		JOIN events e ON i.event_id = e.id
		WHERE i.billettholder_id = ?
		ORDER BY i.pulje_id ASC
	`

	rows, queryErr := db.Query(query, billettholderID)
	if queryErr != nil {
		logger.Error("Failed to query interests", "error", queryErr)
		return nil, queryErr
	}
	defer rows.Close()

	var interests []UserInterest

	for rows.Next() {
		var interest UserInterest

		scanErr := rows.Scan(
			&interest.EventID,
			&interest.EventName,
			&interest.InterestLevel,
			&interest.PuljeID,
		)

		if scanErr != nil {
			logger.Error("Failed to scan interest", "error", scanErr)
			continue
		}

		logger.Info("Interest scanned", "eventID", interest.EventID, "eventName", interest.EventName, "interestLevel", interest.InterestLevel, "puljeID", interest.PuljeID)
		interests = append(interests, interest)
	}

	if rowsErr := rows.Err(); rowsErr != nil {
		logger.Error("Error iterating rows", "error", rowsErr)
		return nil, rowsErr
	}

	logger.Info("Found interests", "count", len(interests))

	return interests, nil
}
