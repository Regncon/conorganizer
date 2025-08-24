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

type BackupHandlerOptions struct {
	DB       *sql.DB
	Logger   *slog.Logger
	Cfg      Config
	FilePath string
	Id       int64
	Stage    BackupStage
	Status   BackupLogStatus
	Error    string
	Interval BackupInterval
}
