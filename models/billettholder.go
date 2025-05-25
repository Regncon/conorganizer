package models

import "time"

/*
CREATE TABLE IF NOT EXISTS billettholdere (

	id INTEGER PRIMARY KEY AUTOINCREMENT,
	first_name TEXT NOT NULL,
	last_name TEXT NOT NULL,
	ticket_type TEXT NOT NULL,
	is_over_18 BOOLEAN NOT NULL,
	order_id INTEGER NOT NULL,
	ticket_id INTEGER NOT NULL UNIQUE,
	ticket_email TEXT NOT NULL,
	order_email TEXT NOT NULL,
	ticket_category_id TEXT NOT NULL,
	inserted_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (ticket_type) REFERENCES ticket_types(name)

);
*/
type Billettholder struct {
	ID               int       `json:"id"`
	FirstName        string    `json:"first_name"`
	LastName         string    `json:"last_name"`
	TicketType       string    `json:"ticket_type"`
	IsOver18         bool      `json:"is_over_18"`
	OrderID          int       `json:"order_id"`
	TicketID         int       `json:"ticket_id"`
	TicketEmail      string    `json:"ticket_email"`
	OrderEmail       string    `json:"order_email"`
	TicketCategoryID string    `json:"ticket_category_id"`
	InsertedTime     time.Time `json:"inserted_time"`
}
