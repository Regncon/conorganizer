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
	Cfg      Config
	DB       *sql.DB
	DBPrefix string
	Error    string
	FilePath string
	FileSize int64
	Id       int64
	Interval BackupInterval
	Logger   *slog.Logger
	Stage    BackupStage
	Status   BackupLogStatus
}

type BackupLogOutput struct {
	CreatedAt string `json:"createdAt"`
	DBPrefix  string `json:"db_prefix"`
	FilePath  string `json:"filePath"`
	FileSize  int64  `json:"file_size"`
	Id        int64  `json:"title"`
	Interval  string `json:"interval"`
	Message   string `json:"message"`
	Stage     string `json:"stage"`
	Status    string `json:"status"`
}
