package ticketholder

import (
	"database/sql"
	"fmt"
	"log/slog"
	"strings"
	"unicode"

	"github.com/Regncon/conorganizer/service/requestctx"
)

type TicketHolderResponse struct {
	Email string
	Name  string
}

func GetTicketHolders(userInfo requestctx.UserRequestInfo, db *sql.DB, logger *slog.Logger) (TicketHolderResponse, error) {
	logger.Info("Fetching ticket holders from the database...")
	query := `SELECT email, billettholder_id, first_name, last_name
                FROM billettholder_emails [be]
                JOIN billettholdere [bh] ON [be].billettholder_id = [bh].id
                WHERE [be].email = ? `
	rows, ticketHolderQueryErr := db.Query(query, userInfo.Email)
	if ticketHolderQueryErr != nil {
		logger.Error("Failed to query ticket holders", "ticketHolderQueryErr", ticketHolderQueryErr)
		return TicketHolderResponse{
			Email: "",
			Name:  "",
		}, ticketHolderQueryErr
	}
	defer rows.Close()

	var email, firstName, lastName string

	for rows.Next() {
		var ticketHolderID int
		if ticketHolderScanErr := rows.Scan(&email, &ticketHolderID, &firstName, &lastName); ticketHolderScanErr != nil {
			logger.Error("Failed to scan ticket holder row", "ticketHolderScanErr", ticketHolderScanErr)
			continue
		}
		logger.Info("Ticket Holder", "email", email, "id", ticketHolderID, "firstName", firstName, "lastName", lastName)
	}
	if ticketHolderRowsErr := rows.Err(); ticketHolderRowsErr != nil {
		logger.Error("Error iterating over ticket holder rows", "ticketHolderRowsErr", ticketHolderRowsErr)
	}

	return TicketHolderResponse{
			Email: email,
			Name:  fmt.Sprintf("%s %s", firstName, lastName),
		},
		nil

}

func GetInitials(s string) string {
	var initialsBuilder strings.Builder
	words := strings.Fields(s)
	if len(words) == 0 {
		return "TT"
	}

	firstChar := rune(words[0][0])
	lastChar := rune(words[len(words)-1][0])

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
