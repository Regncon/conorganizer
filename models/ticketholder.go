package models

import "time"

type Billettholder struct {
	ID               int64     `json:"id"`
	FirstName        string    `json:"firstName"`
	LastName         string    `json:"lastName"`
	Over18           bool      `json:"over18"`
	TicketEmail      string    `json:"ticketEmail"`
	OrderEmails      []string  `json:"orderEmails"`
	TicketID         int       `json:"ticketId"`
	OrderID          int       `json:"orderId"`
	TicketCategory   string    `json:"ticketCategory"`
	TicketCategoryID int       `json:"ticketCategoryId"`
	CreatedAt        time.Time `json:"createdAt"`
	CreatedBy        string    `json:"createdBy"`
	UpdatedAt        time.Time `json:"updatedAt"`
	UpdatedBy        string    `json:"updatedBy"`
	ConnectedEmails  []string  `json:"connectedEmails"`
}
