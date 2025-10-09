package formsubmission

import (
	"database/sql"
	"fmt"
)

func shouldShowStringValue(value string) string {
	if value != "" {
		return value
	}
	return ""
}

func shouldShowNumberValue(value int64) string {
	if value != 0 {
		return fmt.Sprintf("%d", value)
	}
	return ""
}

func test(eventId string, db *sql.DB) {
	interessertePlayersQuery := `SELECT [E].id, [E].host, first_name, last_name
    FROM events [E]
    JOIN event_puljer [EP] ON [E].id = [EP].event_id
    JOIN billettholdere_users [BHU] ON [BHU].user_id = [E].host
    JOIN billettholdere [BH] ON [BH].Id = [BHU].billettholder_id
    WHERE event_id = ?`
	rows, err := db.Query(interessertePlayersQuery, eventId)
	if err != nil {
		fmt.Println("Error querying interested players:", err)
		return
	}

	fmt.Printf("Event ID: %s\n", eventId)
	fmt.Printf("Event ID: %+v\n", rows)
	defer rows.Close()

	for rows.Next() {
		var id, host, firstName, lastName string
		if err := rows.Scan(&id, &host, &firstName, &lastName); err != nil {
			fmt.Println("Error scanning row:", err)
			continue
		}
		fmt.Printf("ID: %s, Host: %s, First Name: %s, Last Name: %s\n",
			shouldShowStringValue(id),
			shouldShowStringValue(host),
			shouldShowStringValue(firstName),
			shouldShowStringValue(lastName),
		)
	}

	if err := rows.Err(); err != nil {
		fmt.Println("Error iterating over rows:", err)
	}
}
