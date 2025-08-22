package models

import (
	"database/sql"
	"log/slog"
)

type BackupInterval string
type BackupStage string

const (
	Hourly        BackupInterval = "hourly"
	Daily         BackupInterval = "daily"
	Weekly        BackupInterval = "weekly"
	Yearly        BackupInterval = "yearly"
	Manually      BackupInterval = "manually"
	Initializing  BackupStage    = "starting"
	Downloading   BackupStage    = "downloading"
	Decompressing BackupStage    = "decompressing"
	Moving        BackupStage    = "moving"
	Finalizing    BackupStage    = "finalizing"
)

type BackupOutcome struct {
	DB         *sql.DB
	LogID      int64
	Status     BackupLogStatus
	Error      string
	Stage      BackupStage
	Interval   BackupInterval
	Logger     *slog.Logger
	WebhookURL string
}
