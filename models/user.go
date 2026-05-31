package models

type User struct {
	ID         int        `json:"id"`
	ExternalID string     `json:"external_id"`
	Email      string     `json:"email"`
	IsAdmin    bool       `json:"is_admin"`
	InsertedAt DBDateTime `json:"inserted_at"`
}
