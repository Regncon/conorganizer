package ticketholder

import (
	"database/sql"
	"fmt"
	"slices"
	"strings"
	"unicode"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/service/requestctx"
)

type BillettHolder struct {
	Email string
	Name  string
	Id    int
	Color string
}

func GetTicketHolders(userInfo requestctx.UserRequestInfo, db *sql.DB) ([]BillettHolder, error) {
	// todo: use the correct way to get billettholders (billettholderservice.GetBilettholdere has a fallback to get all billettholders)
	query := `
    SELECT
        [be].email,
        [be].billettholder_id,
        [bh].first_name,
        [bh].last_name
    FROM
        billettholder_emails [be]
        LEFT JOIN billettholdere [bh] ON [be].billettholder_id = [bh].id
    WHERE
        [be].kind = 'Ticket'
        AND [be].billettholder_id IN (
            SELECT
                billettholder_id
            FROM
                billettholder_emails
            WHERE
                email = ?
        )
`
	rows, ticketHolderQueryErr := db.Query(query, userInfo.Email)
	if ticketHolderQueryErr != nil {
		return nil, fmt.Errorf("failed to query ticket holders for email %q: %w", userInfo.Email, ticketHolderQueryErr)
	}
	defer rows.Close()

	var email, firstName, lastName string
	var associatedTicketholders []BillettHolder

	for rows.Next() {
		var billettHolderId int

		if ticketHolderScanErr := rows.Scan(&email, &billettHolderId, &firstName, &lastName); ticketHolderScanErr != nil {
			return nil, fmt.Errorf("failed to scan ticket holder row: %w", ticketHolderScanErr)
		}
		associatedTicketholders = append(associatedTicketholders, BillettHolder{
			Email: email,
			Name:  fmt.Sprintf("%s %s", firstName, lastName),
			Id:    billettHolderId,
			Color: ColorForName(fmt.Sprintf("%s %s", firstName, lastName)),
		})

	}
	if ticketHolderRowsErr := rows.Err(); ticketHolderRowsErr != nil {
		return nil, fmt.Errorf("error iterating over ticket holder rows: %w", ticketHolderRowsErr)
	}

	return associatedTicketholders, nil

}

func GetPuljerFromEventId(eventId string, db *sql.DB) ([]models.Pulje, error) {
	puljerQuery := `SELECT pulje_id FROM event_puljer WHERE event_id = ? AND is_active = 1 AND is_published = 1`
	rows, puljerErr := db.Query(puljerQuery, eventId)
	if puljerErr != nil {
		return nil, fmt.Errorf("failed to query event puljer for event %s: %w", eventId, puljerErr)
	}
	defer rows.Close()

	var puljer []models.Pulje
	for rows.Next() {
		var puljeName models.Pulje
		if scanErr := rows.Scan(&puljeName); scanErr != nil {
			return nil, fmt.Errorf("failed to scan pulje row for event %s: %w", eventId, scanErr)
		}
		puljer = append(puljer, puljeName)
	}
	if rowsErr := rows.Err(); rowsErr != nil {
		return nil, fmt.Errorf("error iterating over pulje rows for event %s: %w", eventId, rowsErr)
	}

	return puljer, nil
}

func GetYourBillettHolderInfo(userInfo requestctx.UserRequestInfo, ticketHolders []BillettHolder) BillettHolder {
	idx := slices.IndexFunc(ticketHolders, func(th BillettHolder) bool {
		return th.Email == userInfo.Email
	})

	if idx == -1 {
		return BillettHolder{
			Email: "unknown@example.com",
			Name:  "Unknown Ticket Holder",
		}
	}

	return ticketHolders[idx]
}

func GetInitials(s string) string {
	var initialsBuilder strings.Builder
	words := strings.Fields(s)
	if len(words) == 0 {
		return "TT"
	}

	firstWordRunes := []rune(words[0])
	lastWordRunes := []rune(words[len(words)-1])
	if len(firstWordRunes) == 0 || len(lastWordRunes) == 0 {
		return "TT"
	}
	firstChar := firstWordRunes[0]
	lastChar := lastWordRunes[0]

	if unicode.IsLetter(firstChar) {
		initialsBuilder.WriteRune(unicode.ToUpper(firstChar))
	}
	if !unicode.IsLetter(firstChar) {
		initialsBuilder.WriteRune('T')
	}

	if unicode.IsLetter(lastChar) {
		initialsBuilder.WriteRune(unicode.ToUpper(lastChar))
	}
	if !unicode.IsLetter(lastChar) {
		initialsBuilder.WriteRune('T')
	}

	return initialsBuilder.String()
}

func GetFirstNameAndLastInitial(fullName string) string {
	nameParts := strings.Fields(fullName)
	if len(nameParts) == 0 {
		return ""
	}

	if len(nameParts) == 1 {
		return nameParts[0]
	}

	lastNameRunes := []rune(nameParts[len(nameParts)-1])
	if len(lastNameRunes) == 0 {
		return nameParts[0]
	}

	return fmt.Sprintf("%s %c", nameParts[0], unicode.ToUpper(lastNameRunes[0]))
}
