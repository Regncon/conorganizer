package models

type BackupLogStatus string

const (
	Pending BackupLogStatus = "pending"
	Success BackupLogStatus = "success"
	Error   BackupLogStatus = "error"
)

type BackupLogInput struct {
	ID      int64
	Status  BackupLogStatus
	Message error
}

type BackupLogMessage struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Intro       string `json:"intro"`
	Description string `json:"description"`
}
