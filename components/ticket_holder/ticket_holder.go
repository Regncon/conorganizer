package ticketholder

import (
	"database/sql"
	"fmt"
	"log/slog"
)

func GetTicketHolders(db *sql.DB, logger *slog.Logger) {

	fmt.Println("Fetching ticket holders from the database...")

	email := "__"
	query := `SELECT email, billettholder_id, first_name, last_name
                FROM billettholder_emails [be]
                JOIN billettholdere [bh] ON [be].billettholder_id = [bh].id
                WHERE [be].email = ? `
	rows, queryErr := db.Query(query, email)
	if queryErr != nil {
		logger.Error("Failed to query ticket holders", "error", queryErr)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var email, firstName, lastName string
		var ticketHolderID int
		if err := rows.Scan(&email, &ticketHolderID, &firstName, &lastName); err != nil {
			logger.Error("Failed to scan ticket holder row", "error", err)
			continue
		}
		logger.Info("Ticket Holder", "email", email, "id", ticketHolderID, "firstName", firstName, "lastName", lastName)
	}
	if err := rows.Err(); err != nil {
		logger.Error("Error iterating over ticket holder rows", "error", err)
	}

}
