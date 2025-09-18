package models

import "time"

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
