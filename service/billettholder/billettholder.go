package billettholderService

import (
	"database/sql"
	"fmt"
	"sort"

	"github.com/Regncon/conorganizer/models"
)

type BillettholderFilters struct {
	WithoutFirstChoice bool
	GMOrDM             bool
}

func GetBillettholdere(userId string, db *sql.DB) ([]models.Billettholder, error) {
	var rows *sql.Rows
	var err error
	if userId == "" {
		queryAll := (`
        SELECT
            b.id, b.first_name, b.last_name, b.ticket_type_id, b.ticket_type,
            b.is_over_18, b.order_id, b.ticket_id, b.created_at, b.updated_at,
            e.id, e.email, e.kind, e.created_at, e.updated_at
        FROM billettholdere AS b
        LEFT JOIN relation_billettholder_emails AS e
            ON b.id = e.billettholder_id
        ORDER BY b.id, e.id
	`)
		rows, err = db.Query(queryAll)
	} else {
		queryByUser := (`
        SELECT
            b.id, b.first_name, b.last_name, b.ticket_type_id, b.ticket_type,
            b.is_over_18, b.order_id, b.ticket_id, b.created_at, b.updated_at,
            e.id, e.email, e.kind, e.created_at, e.updated_at
        FROM billettholdere AS b
        JOIN relation_billettholdere_users bu ON b.id = bu.billettholder_id
        JOIN users u ON bu.user_id = u.id
        LEFT JOIN relation_billettholder_emails AS e
            ON b.id = e.billettholder_id
        WHERE u.external_id = ?
        ORDER BY b.id, e.id
    `)
		rows, err = db.Query(queryByUser, userId)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to query billettholdere: %w", err)
	}
	defer rows.Close()

	return scanBillettholdere(rows)
}

func GetBillettholdereWithFilters(userId string, db *sql.DB, filters BillettholderFilters) ([]models.Billettholder, error) {
	query := `
		WITH published_assignments AS (
			SELECT
				rep.billettholder_id,
				rep.role,
				i.interest_level
			FROM relation_events_players AS rep
			INNER JOIN relation_event_puljer AS ep
				ON ep.event_id = rep.event_id
				AND ep.pulje_id = rep.pulje_id
				AND ep.is_published = 1
			LEFT JOIN interests AS i
				ON i.billettholder_id = rep.billettholder_id
				AND i.event_id = rep.event_id
				AND i.pulje_id = rep.pulje_id
		),
		first_choice_billettholdere AS (
			SELECT DISTINCT billettholder_id
			FROM published_assignments
			WHERE role = ? AND interest_level = ?
		),
		gm_billettholdere AS (
			SELECT DISTINCT billettholder_id
			FROM published_assignments
			WHERE role = ?
		)
		SELECT
			b.id, b.first_name, b.last_name, b.ticket_type_id, b.ticket_type,
			b.is_over_18, b.order_id, b.ticket_id, b.created_at, b.updated_at,
			e.id, e.email, e.kind, e.created_at, e.updated_at
		FROM billettholdere AS b
		LEFT JOIN relation_billettholder_emails AS e
			ON b.id = e.billettholder_id
	`
	args := []any{
		models.EventPlayerRolePlayer,
		models.InterestLevelHigh,
		models.EventPlayerRoleGM,
	}

	if userId != "" {
		query += `
			JOIN relation_billettholdere_users bu ON b.id = bu.billettholder_id
			JOIN users u ON bu.user_id = u.id
		`
	}

	query += `
		WHERE (? = 0 OR NOT EXISTS (
			SELECT 1
			FROM first_choice_billettholdere AS fc
			WHERE fc.billettholder_id = b.id
		))
		AND (? = 0 OR EXISTS (
			SELECT 1
			FROM gm_billettholdere AS gm
			WHERE gm.billettholder_id = b.id
		))
	`
	args = append(args, boolToInt(filters.WithoutFirstChoice), boolToInt(filters.GMOrDM))

	if userId != "" {
		query += ` AND u.external_id = ?`
		args = append(args, userId)
	}

	query += ` ORDER BY b.id, e.id`

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query filtered billettholdere: %w", err)
	}
	defer rows.Close()

	return scanBillettholdere(rows)
}

func scanBillettholdere(rows *sql.Rows) ([]models.Billettholder, error) {
	type emailRow struct {
		id        sql.NullInt64
		email     sql.NullString
		kind      sql.NullString
		createdAt models.DBDateTime
		updatedAt models.DBDateTime
	}

	byID := make(map[int]*models.Billettholder)
	order := make([]int, 0, 512)

	for rows.Next() {
		var b models.Billettholder
		var er emailRow

		if err := rows.Scan(
			&b.ID, &b.FirstName, &b.LastName, &b.TicketTypeId, &b.TicketType,
			&b.IsOver18, &b.OrderID, &b.TicketID, &b.CreatedAt, &b.UpdatedAt,
			&er.id, &er.email, &er.kind, &er.createdAt, &er.updatedAt,
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
				CreatedAt:    b.CreatedAt,
				UpdatedAt:    b.UpdatedAt,
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
				Kind:            models.BillettholderEmailKind(er.kind.String),
				CreatedAt:       er.createdAt,
				UpdatedAt:       er.updatedAt,
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

func boolToInt(value bool) int {
	if value {
		return 1
	}
	return 0
}

func GetBillettholderByUserId(db *sql.DB, userID string) (int, error) {
	var billettholderId int
	row := db.QueryRow(`
        SELECT bu.billettholder_id
        FROM relation_billettholdere_users bu
        JOIN users u ON u.id = bu.user_id
        WHERE u.external_id = $1
        LIMIT 1`, userID)

	if err := row.Scan(&billettholderId); err != nil {
		return 0, fmt.Errorf("failed to scan billettholder row for user %q: %w", userID, err)
	}

	return billettholderId, nil
}
