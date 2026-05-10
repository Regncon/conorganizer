package models

import "time"

type User struct {
	ID         int        `json:"id"`
	ExternalID string     `json:"external_id"`
	Email      string     `json:"email"`
	IsAdmin    bool       `json:"is_admin"`
	InsertedAt *time.Time `json:"inserted_at,omitempty"`
}
