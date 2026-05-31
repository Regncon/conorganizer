package billettholderadmin

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/Regncon/conorganizer/models"
)

type billettholderInterestEventRow struct {
	EventID       string
	EventTitle    string
	EventStatus   models.EventStatus
	IsPublished   bool
	InterestLevel models.InterestLevel
	AssignedRole  models.EventPlayerRole
}

type billettholderInterestPuljeSection struct {
	PuljeID  models.Pulje
	Name     string
	Assigned []billettholderInterestEventRow
	High     []billettholderInterestEventRow
	Medium   []billettholderInterestEventRow
	Low      []billettholderInterestEventRow
}

type billettholderInterestSectionKey struct {
	BillettholderID int
	PuljeID         models.Pulje
}

func billettholderIDs(billettholdere []models.Billettholder) []int {
	ids := make([]int, 0, len(billettholdere))
	for _, billettholder := range billettholdere {
		ids = append(ids, billettholder.ID)
	}
	return ids
}

func getBillettholderInterestSectionsByBillettholderID(
	db *sql.DB,
	billettholderIDs []int,
) (map[int][]billettholderInterestPuljeSection, error) {
	ids := uniquePositiveBillettholderIDs(billettholderIDs)
	result := make(map[int][]billettholderInterestPuljeSection, len(ids))
	for _, id := range ids {
		result[id] = nil
	}
	if len(ids) == 0 {
		return result, nil
	}

	args := make([]any, 0, len(ids))
	for _, id := range ids {
		args = append(args, id)
	}

	query := fmt.Sprintf(`
		WITH requested_billettholdere(billettholder_id) AS (
			VALUES %s
		),
		billettholder_interest_rows AS (
			SELECT
				rep.billettholder_id,
				p.id AS pulje_id,
				p.name AS pulje_name,
				p.start_at AS pulje_start_at,
				e.id AS event_id,
				e.title AS event_title,
				e.status AS event_status,
				COALESCE(ep.is_published, 0) AS is_published,
				COALESCE(i.interest_level, '') AS interest_level,
				rep.role AS assigned_role,
				1 AS is_assigned
			FROM relation_events_players AS rep
			INNER JOIN requested_billettholdere AS rb
				ON rb.billettholder_id = rep.billettholder_id
			INNER JOIN events AS e
				ON e.id = rep.event_id
			INNER JOIN puljer AS p
				ON p.id = rep.pulje_id
			LEFT JOIN relation_event_puljer AS ep
				ON ep.event_id = rep.event_id
				AND ep.pulje_id = rep.pulje_id
			LEFT JOIN interests AS i
				ON i.billettholder_id = rep.billettholder_id
				AND i.event_id = rep.event_id
				AND i.pulje_id = rep.pulje_id

			UNION ALL

			SELECT
				i.billettholder_id,
				p.id AS pulje_id,
				p.name AS pulje_name,
				p.start_at AS pulje_start_at,
				e.id AS event_id,
				e.title AS event_title,
				e.status AS event_status,
				COALESCE(ep.is_published, 0) AS is_published,
				i.interest_level AS interest_level,
				'' AS assigned_role,
				0 AS is_assigned
			FROM interests AS i
			INNER JOIN requested_billettholdere AS rb
				ON rb.billettholder_id = i.billettholder_id
			INNER JOIN events AS e
				ON e.id = i.event_id
			INNER JOIN puljer AS p
				ON p.id = i.pulje_id
			LEFT JOIN relation_event_puljer AS ep
				ON ep.event_id = i.event_id
				AND ep.pulje_id = i.pulje_id
			LEFT JOIN relation_events_players AS rep
				ON rep.billettholder_id = i.billettholder_id
				AND rep.event_id = i.event_id
				AND rep.pulje_id = i.pulje_id
			WHERE rep.billettholder_id IS NULL
		)
		SELECT
			billettholder_id,
			pulje_id,
			pulje_name,
			event_id,
			event_title,
			event_status,
			is_published,
			interest_level,
			assigned_role,
			is_assigned
		FROM billettholder_interest_rows
		ORDER BY
			pulje_start_at,
			pulje_name COLLATE NOCASE,
			pulje_id,
			is_assigned DESC,
			CASE WHEN is_assigned = 1 THEN event_title ELSE '' END COLLATE NOCASE,
			CASE interest_level
				WHEN ? THEN 1
				WHEN ? THEN 2
				WHEN ? THEN 3
				ELSE 4
			END,
			event_title COLLATE NOCASE,
			event_id
	`, valuesPlaceholders(len(ids)))

	args = append(args, models.InterestLevelHigh, models.InterestLevelMedium, models.InterestLevelLow)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("query billettholder interest sections: %w", err)
	}
	defer rows.Close()

	sectionIndex := map[billettholderInterestSectionKey]int{}
	for rows.Next() {
		var (
			billettholderID int
			puljeID         string
			puljeName       string
			eventID         string
			eventTitle      string
			eventStatus     string
			isPublished     int
			interestLevel   string
			assignedRole    string
			isAssigned      int
		)
		if err := rows.Scan(
			&billettholderID,
			&puljeID,
			&puljeName,
			&eventID,
			&eventTitle,
			&eventStatus,
			&isPublished,
			&interestLevel,
			&assignedRole,
			&isAssigned,
		); err != nil {
			return nil, fmt.Errorf("scan billettholder interest row: %w", err)
		}

		key := billettholderInterestSectionKey{
			BillettholderID: billettholderID,
			PuljeID:         models.Pulje(puljeID),
		}
		sectionIndexForKey, ok := sectionIndex[key]
		if !ok {
			result[billettholderID] = append(result[billettholderID], billettholderInterestPuljeSection{
				PuljeID: key.PuljeID,
				Name:    puljeName,
			})
			sectionIndexForKey = len(result[billettholderID]) - 1
			sectionIndex[key] = sectionIndexForKey
		}

		row := billettholderInterestEventRow{
			EventID:       eventID,
			EventTitle:    eventTitle,
			EventStatus:   models.EventStatus(eventStatus),
			IsPublished:   isPublished == 1,
			InterestLevel: models.InterestLevel(interestLevel),
			AssignedRole:  models.EventPlayerRole(assignedRole),
		}
		appendBillettholderInterestRow(result[billettholderID], sectionIndexForKey, row, isAssigned == 1)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate billettholder interest rows: %w", err)
	}

	return result, nil
}

func uniquePositiveBillettholderIDs(ids []int) []int {
	uniqueIDs := make([]int, 0, len(ids))
	seen := map[int]struct{}{}
	for _, id := range ids {
		if id <= 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		uniqueIDs = append(uniqueIDs, id)
	}
	return uniqueIDs
}

func valuesPlaceholders(count int) string {
	placeholders := make([]string, 0, count)
	for range count {
		placeholders = append(placeholders, "(?)")
	}
	return strings.Join(placeholders, ", ")
}

func appendBillettholderInterestRow(
	sections []billettholderInterestPuljeSection,
	sectionIndex int,
	row billettholderInterestEventRow,
	isAssigned bool,
) {
	if isAssigned {
		sections[sectionIndex].Assigned = append(sections[sectionIndex].Assigned, row)
		return
	}

	switch row.InterestLevel {
	case models.InterestLevelHigh:
		sections[sectionIndex].High = append(sections[sectionIndex].High, row)
	case models.InterestLevelMedium:
		sections[sectionIndex].Medium = append(sections[sectionIndex].Medium, row)
	case models.InterestLevelLow:
		sections[sectionIndex].Low = append(sections[sectionIndex].Low, row)
	}
}
