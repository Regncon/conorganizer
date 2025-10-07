package billettholderService

import (
	"database/sql"
	"log/slog"
	"sort"

	"github.com/Regncon/conorganizer/models"
)

func GetBilettholdere(userId string, db *sql.DB, logger *slog.Logger) ([]models.Billettholder, error) {
	var rows *sql.Rows
	var err error
	if userId == "" {
		queryAll := (`
        SELECT
            b.id, b.first_name, b.last_name, b.ticket_type_id, b.ticket_type,
            b.is_over_18, b.order_id, b.ticket_id, b.inserted_time,
            e.id, e.email, e.kind, e.inserted_time
        FROM billettholdere AS b
        LEFT JOIN billettholder_emails AS e
            ON b.id = e.billettholder_id
        ORDER BY b.id, e.id
	`)
		rows, err = db.Query(queryAll)
	} else {
		queryByUser := (`
        SELECT
            b.id, b.first_name, b.last_name, b.ticket_type_id, b.ticket_type,
            b.is_over_18, b.order_id, b.ticket_id, b.inserted_time,
            e.id, e.email, e.kind, e.inserted_time
        FROM billettholdere AS b
        JOIN billettholdere_users bu ON b.id = bu.billettholder_id
        JOIN users u ON bu.user_id = u.id
        LEFT JOIN billettholder_emails AS e
            ON b.id = e.billettholder_id
        WHERE u.user_id = ?
        ORDER BY b.id, e.id
    `)
		rows, err = db.Query(queryByUser, userId)
	}

	if err != nil {
		logger.Error("Failed to query billettholdere", "error", err)
		return nil, err
	}
	defer rows.Close()

	type emailRow struct {
		id           sql.NullInt64
		email        sql.NullString
		kind         sql.NullString
		insertedTime sql.NullTime
	}

	byID := make(map[int]*models.Billettholder)
	order := make([]int, 0, 512)

	for rows.Next() {
		var b models.Billettholder
		var er emailRow

		if err := rows.Scan(
			&b.ID, &b.FirstName, &b.LastName, &b.TicketTypeId, &b.TicketType,
			&b.IsOver18, &b.OrderID, &b.TicketID, &b.InsertedTime,
			&er.id, &er.email, &er.kind, &er.insertedTime,
		); err != nil {
			logger.Error("Failed to scan row", "error", err)
			return nil, err
		}

		holder, ok := byID[b.ID]
		if !ok {
			holder = &models.Billettholder{
				ID:           b.ID,
				FirstName:    b.FirstName,
				LastName:     b.LastName,
				TicketTypeId: b.TicketTypeId,
				TicketType:   b.TicketType,
				IsOver18:     b.IsOver18,
				OrderID:      b.OrderID,
				TicketID:     b.TicketID,
				InsertedTime: b.InsertedTime,
				Emails:       nil,
			}
			byID[b.ID] = holder
			order = append(order, b.ID)
		}

		if er.id.Valid {
			holder.Emails = append(holder.Emails, models.BillettholderEmail{
				ID:              int(er.id.Int64),
				BillettholderID: b.ID,
				Email:           er.email.String,
				Kind:            er.kind.String,
				InsertedTime:    er.insertedTime.Time,
			})
		}
	}
	if err := rows.Err(); err != nil {
		logger.Error("Row iteration error", "error", err)
		return nil, err
	}

	sort.Ints(order) // or remove this to keep INSERT order from the query
	out := make([]models.Billettholder, 0, len(order))
	for _, id := range order {
		out = append(out, *byID[id])
	}
	return out, nil
}

func GetBillettholderByUserId(db *sql.DB, logger *slog.Logger, userID string) (int, error) {
	var billettholderId int
	row := db.QueryRow(`
        SELECT id FROM billettholdere WHERE user_id = $1 `, userID)

	if err := row.Scan(&billettholderId); err != nil {
		logger.Error("Failed to scan row", "error", err)
		return 0, err
	}

	return billettholderId, nil
}
