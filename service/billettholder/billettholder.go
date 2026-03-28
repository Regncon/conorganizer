package billettholderService

import (
	"database/sql"
	"fmt"
	"sort"

	"github.com/Regncon/conorganizer/models"
)

func GetBilettholdere(userId string, db *sql.DB) ([]models.Billettholder, error) {
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
        WITH current_user AS (
            SELECT id, email
            FROM users
            WHERE user_id = ?
        ),
        linked_orders AS (
            SELECT DISTINCT b.order_id
            FROM billettholdere AS b
            JOIN billettholdere_users AS bu
                ON b.id = bu.billettholder_id
            JOIN current_user AS u
                ON bu.user_id = u.id
            UNION
            SELECT DISTINCT b.order_id
            FROM billettholdere AS b
            JOIN billettholder_emails AS e
                ON b.id = e.billettholder_id
            JOIN current_user AS u
                ON e.email = u.email
        )
        SELECT
            b.id, b.first_name, b.last_name, b.ticket_type_id, b.ticket_type,
            b.is_over_18, b.order_id, b.ticket_id, b.inserted_time,
            e.id, e.email, e.kind, e.inserted_time
        FROM billettholdere AS b
        LEFT JOIN billettholder_emails AS e
            ON b.id = e.billettholder_id
        WHERE b.order_id IN (SELECT order_id FROM linked_orders)
        ORDER BY b.id, e.id
    `)
		rows, err = db.Query(queryByUser, userId)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to query billettholdere: %w", err)
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
			return nil, fmt.Errorf("failed to scan billettholder row: %w", err)
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
		return nil, fmt.Errorf("row iteration error for billettholdere: %w", err)
	}

	sort.Ints(order) // or remove this to keep INSERT order from the query
	out := make([]models.Billettholder, 0, len(order))
	for _, id := range order {
		out = append(out, *byID[id])
	}
	return out, nil
}

func GetBillettholderByUserId(db *sql.DB, userID string) (int, error) {
	var billettholderId int
	row := db.QueryRow(`
        SELECT id FROM billettholdere WHERE user_id = $1 `, userID)

	if err := row.Scan(&billettholderId); err != nil {
		return 0, fmt.Errorf("failed to scan billettholder row for user %q: %w", userID, err)
	}

	return billettholderId, nil
}
