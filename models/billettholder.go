package models

import (
	"database/sql"
)

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
	CreatedAt    DBDateTime           `json:"created_at"`
	UpdatedAt    DBDateTime           `json:"updated_at"`
	CreatedByID  sql.NullInt64        `json:"created_by_id"`
	UpdatedByID  sql.NullInt64        `json:"updated_by_id"`
}

type BillettholderEmail struct {
	ID              int                    `json:"id"`
	BillettholderID int                    `json:"billettholder_id"`
	Email           string                 `json:"email"`
	Kind            BillettholderEmailKind `json:"kind"` // See BillettholderEmailKind constants.
	CreatedAt       DBDateTime             `json:"created_at"`
	UpdatedAt       DBDateTime             `json:"updated_at"`
	CreatedByID     sql.NullInt64          `json:"created_by_id"`
	UpdatedByID     sql.NullInt64          `json:"updated_by_id"`
}

type BillettholderEmailKind string

const (
	BillettholderEmailKindTicket     BillettholderEmailKind = "Ticket"
	BillettholderEmailKindAssociated BillettholderEmailKind = "Associated"
	BillettholderEmailKindManual     BillettholderEmailKind = "Manual"
)

func (kind BillettholderEmailKind) Label() string {
	switch kind {
	case BillettholderEmailKindTicket:
		return "Billett"
	case BillettholderEmailKindAssociated:
		return "Tilknyttet"
	case BillettholderEmailKindManual:
		return "Manuell"
	default:
		return string(kind)
	}
}

type BillettholderUsers struct {
	BillettholderID int        `json:"billettholder_id"`
	UserID          int        `json:"user_id"`
	InsertedAt      DBDateTime `json:"inserted_at"`
}

/*
CREATE TABLE relation_events_players (

	    event_id TEXT NOT NULL,
	    pulje_id TEXT NOT NULL,
	    billettholder_id INTEGER NOT NULL,
	    role TEXT NOT NULL DEFAULT 'Player' CHECK (role IN ('Player', 'GM')),
	    -- inserted_at uses the DBDateTimeNowSQL default expression.
	    inserted_at TEXT,
	    PRIMARY KEY (billettholder_id, event_id, pulje_id),
	    FOREIGN KEY (billettholder_id) REFERENCES billettholdere (id),
	    FOREIGN KEY (event_id) REFERENCES events (id),
	    FOREIGN KEY (pulje_id) REFERENCES puljer (id)
	);
*/
type EventPlayerRole string

const (
	EventPlayerRolePlayer EventPlayerRole = "Player"
	EventPlayerRoleGM     EventPlayerRole = "GM"
)

func (role EventPlayerRole) Label() string {
	switch role {
	case EventPlayerRolePlayer:
		return "spelar"
	case EventPlayerRoleGM:
		return "GM"
	default:
		return string(role)
	}
}

type EventPlayer struct {
	EventID         string          `json:"event_id"`
	PuljeID         string          `json:"pulje_id"`
	BillettholderID int             `json:"billettholder_id"`
	Role            EventPlayerRole `json:"role"`
	InsertedAt      DBDateTime      `json:"inserted_at"`
}
