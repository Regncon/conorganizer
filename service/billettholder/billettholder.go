package billettholderService

import (
	"database/sql"
	"log/slog"

	"github.com/Regncon/conorganizer/models"
)

func GetBilettholdere(db *sql.DB, logger *slog.Logger) ([]models.Billettholder, error) {
	rows, err := db.Query("SELECT id, first_name,  last_name, ticket_type, ticket_type_id, is_over_18, order_id, ticket_id FROM billettholdere")
	if err != nil {
		logger.Error("Failed to query billettholdere", "error", err)
		return nil, err
	}
	defer rows.Close()

	var billettholdere []models.Billettholder
	for rows.Next() {
		var b models.Billettholder
		if err := rows.Scan(&b.ID, &b.FirstName, &b.LastName, &b.TicketType, &b.TicketTypeId, &b.IsOver18, &b.OrderID, &b.TicketID); err != nil {
			logger.Error("Failed to scan billettholder", "error", err)
			return nil, err
		}
		billettholdere = append(billettholdere, b)
	}
	return billettholdere, nil
}
