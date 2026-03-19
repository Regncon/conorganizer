package models

import "time"

type Billettholder struct {
	ID           int                  `json:"id"`
	FirstName    string               `json:"first_name"`
	LastName     string               `json:"last_name"`
	Emails       []BillettholderEmail `json:"emails,omitempty"`
	TicketTypeId int                  `json:"ticket_type_id"`
	TicketType   string               `json:"ticket_type"`
	IsOver18     bool                 `json:"is_over_18"`
	OrderID      int                  `json:"order_id"`
	TicketID     int                  `json:"ticket_id"`
	InsertedTime time.Time            `json:"inserted_time"`
}

type BillettholderEmail struct {
	ID              int       `json:"id"`
	BillettholderID int       `json:"billettholder_id"`
	Email           string    `json:"email"`
	Kind            string    `json:"kind"` // 'Ticket','Associated','Manual'
	InsertedTime    time.Time `json:"inserted_time"`
}

type BillettholderUsers struct {
	BillettholderID int       `json:"billettholder_id"`
	UserID          int       `json:"user_id"`
	InsertedTime    time.Time `json:"inserted_time"`
}

/*
CREATE TABLE events_players (

	    event_id TEXT NOT NULL,
	    pulje_id TEXT NOT NULL,
	    billettholder_id INTEGER NOT NULL,
	    isPlayer BOOLEAN NOT NULL DEFAULT TRUE,
	    isGm BOOLEAN NOT NULL DEFAULT FALSE,
	    inserted_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	    PRIMARY KEY (billettholder_id, event_id, pulje_id),
	    FOREIGN KEY (billettholder_id) REFERENCES billettholdere (id),
	    FOREIGN KEY (event_id) REFERENCES events (id),
	    FOREIGN KEY (pulje_id) REFERENCES puljer (id)
	);
*/
type EventPlayer struct {
	EventID         string    `json:"event_id"`
	PuljeID         string    `json:"pulje_id"`
	BillettholderID int       `json:"billettholder_id"`
	IsPlayer        bool      `json:"is_player"`
	IsGm            bool      `json:"is_gm"`
	InsertedTime    time.Time `json:"inserted_time"`
}
