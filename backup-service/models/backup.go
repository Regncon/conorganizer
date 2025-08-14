package models

type BackupInterval string

const (
	Hourly   BackupInterval = "hourly"
	Daily    BackupInterval = "daily"
	Weekly   BackupInterval = "weekly"
	Yearly   BackupInterval = "yearly"
	Manually BackupInterval = "manually"
)
