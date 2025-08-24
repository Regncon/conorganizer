package models

import (
	"database/sql"
	"log/slog"
)

type BackupInterval string
type BackupStage string
type BackupLogStatus string

const (
	Hourly        BackupInterval  = "hourly"
	Daily         BackupInterval  = "daily"
	Weekly        BackupInterval  = "weekly"
	Yearly        BackupInterval  = "yearly"
	Manually      BackupInterval  = "manually"
	Initializing  BackupStage     = "starting"
	Downloading   BackupStage     = "downloading"
	Decompressing BackupStage     = "decompressing"
	Validating    BackupStage     = "validating"
	Moving        BackupStage     = "moving"
	Finalizing    BackupStage     = "completed"
	Pending       BackupLogStatus = "pending"
	Success       BackupLogStatus = "success"
	Error         BackupLogStatus = "error"
)

type BackupHandlerOptions struct {
	DB       *sql.DB
	Logger   *slog.Logger
	Cfg      Config
	FilePath string
	Id       int64
	Interval BackupInterval
	Stage    BackupStage
	Status   BackupLogStatus
	Error    string
}
