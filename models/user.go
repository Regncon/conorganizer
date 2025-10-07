package models

import "time"

type User struct {
	ID      int    `json:"id"`
	UserID  string `json:"user_id"`
	Email   string `json:"email"`
	IsAdmin bool   `json:"is_admin"`
	// InsertedTime is optional and not required
	InsertedTime *time.Time `json:"inserted_time,omitempty"`
}
