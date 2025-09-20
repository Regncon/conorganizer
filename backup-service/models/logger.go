package models

import "time"

type BackupLogInput struct {
	ID      int64
	Status  BackupLogStatus
	Message string
}

type BackupLogMessage struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Intro       string `json:"intro"`
	Description string `json:"description"`
}

type BackupLog struct {
	ID         int
	BackupType string
	Stage      string
	Status     string
	FilePath   string
	FileSize   int64
	DBPrefix   string
	Message    string
	CreatedAt  time.Time
}
