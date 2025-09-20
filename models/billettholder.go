package models

import "time"

type Billettholder struct {
	ID           int       `json:"id"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	TicketTypeId int       `json:"ticket_type_id"`
	TicketType   string    `json:"ticket_type"`
	IsOver18     bool      `json:"is_over_18"`
	OrderID      int       `json:"order_id"`
	TicketID     int       `json:"ticket_id"`
	InsertedTime time.Time `json:"inserted_time"`
}

type BillettholderEmail struct {
	ID              int       `json:"id"`
	BillettholderID int       `json:"billettholder_id"`
	Email           string    `json:"email"`
	Kind            string    `json:"kind"` // 'Ticket','Associated','Manual'
	InsertedTime    time.Time `json:"inserted_time"`
}
