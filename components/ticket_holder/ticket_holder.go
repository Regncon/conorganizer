package ticketholder

import (
	"database/sql"
	"fmt"
	"log/slog"
	"slices"
	"strings"
	"unicode"

	"github.com/Regncon/conorganizer/service/requestctx"
)

type BillettHolder struct {
	Email string
	Name  string
	Id    int
}

func GetTicketHolders(userInfo requestctx.UserRequestInfo, db *sql.DB, logger *slog.Logger) ([]BillettHolder, error) {
	logger.Info("Fetching ticket holders from the database...")
	query := `SELECT email, billettholder_id, first_name, last_name
                FROM billettholder_emails [be]
                JOIN billettholdere [bh] ON [be].billettholder_id = [bh].id
                WHERE [be].email = ? `
	rows, ticketHolderQueryErr := db.Query(query, userInfo.Email)
	if ticketHolderQueryErr != nil {
		logger.Error("Failed to query ticket holders", "ticketHolderQueryErr", ticketHolderQueryErr)
		return []BillettHolder{}, ticketHolderQueryErr
	}
	defer rows.Close()

	var email, firstName, lastName string
	var associatedTicketholders []BillettHolder

	for rows.Next() {
		var billettHolderId int

		if ticketHolderScanErr := rows.Scan(&email, &billettHolderId, &firstName, &lastName); ticketHolderScanErr != nil {
			logger.Error("Failed to scan ticket holder row", "ticketHolderScanErr", ticketHolderScanErr)
			continue
		}
		associatedTicketholders = append(associatedTicketholders, BillettHolder{
			Email: email,
			Name:  fmt.Sprintf("%s %s", firstName, lastName),
			Id:    billettHolderId,
		})

		logger.Info("Ticket Holder", "email", email, "id", billettHolderId, "firstName", firstName, "lastName", lastName)
	}
	if ticketHolderRowsErr := rows.Err(); ticketHolderRowsErr != nil {
		logger.Error("Error iterating over ticket holder rows", "ticketHolderRowsErr", ticketHolderRowsErr)
	}

	associatedTicketholders = append(associatedTicketholders, BillettHolder{
		Email: "lo@najcuksuc.sn",
		Name:  "Leonard Moreno",
		Id:    1,
	})
	associatedTicketholders = append(associatedTicketholders, BillettHolder{
		Email: "lacbe@lecuc.my",
		Name:  "Olive Berry",
		Id:    2,
	})
	associatedTicketholders = append(associatedTicketholders, BillettHolder{
		Email: "mijinpu@posrik.cz",
		Name:  "Bobby Silva",
		Id:    3,
	})
	associatedTicketholders = append(associatedTicketholders, BillettHolder{
		Email: "igkir@mukpunuc.be",
		Name:  "Bertha Francis",
		Id:    5,
	})
	associatedTicketholders = append(associatedTicketholders, BillettHolder{
		Email: "ruidavuf@otavig.gy",
		Name:  "Mario Ross",
		Id:    6,
	})

	return associatedTicketholders, nil

}

func GetYourBillettHolderInfo(userInfo requestctx.UserRequestInfo, ticketHolders []BillettHolder) BillettHolder {
	idx := slices.IndexFunc(ticketHolders, func(th BillettHolder) bool {
		fmt.Printf("Comparing ticket holder email: %s with user email: %s\n", th.Email, userInfo.Email)
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
