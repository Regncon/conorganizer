package services

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/Regncon/conorganizer/backup-service/models"
	"github.com/Regncon/conorganizer/backup-service/utils"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// BackupService handles scheduled backup operations and file rotation logic.
type BackupService struct {
	Config   models.Config
	Db       *sql.DB
	S3Client *s3.Client
	Logger   *slog.Logger
}

// NewBackupService creates a new instance of BackupService with dependencies injected.
func NewBackupService(cfg models.Config, db *sql.DB, s3Client *s3.Client, logger *slog.Logger) *BackupService {
	return &BackupService{
		Config:   cfg,
		Db:       db,
		S3Client: s3Client,
		Logger:   logger,
	}
}

// run is the internal function that performs backup + rotation,
// retention is how many files we need of that interval
func (b *BackupService) run(ctx context.Context, interval models.BackupInterval, retention int) {
	// Signal that backup service has started a job
	b.Logger.Info("Scheduled backup job triggered", "interval", interval)

	// Create db entry for logging
	logID, err := NewLogBackup(b.Db, interval)
	if err != nil {
		b.Logger.Error("Failed to log backup", "err", err)
		return
	}

	// Create status object for tracking
	output := models.BackupHandlerOptions{
		DB:       b.Db,
		Logger:   b.Logger,
		Cfg:      b.Config,
		Id:       logID,
		Interval: interval,
	}

	// download snapshot
	snapshotPath, err := DownloadLatestSnapshot(ctx, b.S3Client, b.Config.BUCKET_NAME, b.Config.DB_PREFIX)
	if err != nil {
		output.Status = models.Error
		output.Stage = models.Downloading
		output.Error = err.Error()
		HandleBackupResult(output)
		return
	}
	b.Logger.Info("Snapshot downloaded", "snapshot", snapshotPath)

	// decompress snapshot
	dbPath, err := utils.DecompressLZ4(snapshotPath)
	if err != nil {
		b.Logger.Error("Decompression failed", "err", err)
		return
	}
	b.Logger.Info("Decompression complete", "db", dbPath)

	// validate snapshot
	if err := utils.ValidateSnapshot(dbPath); err != nil {
		b.Logger.Error("Invalid SQLite snapshot", "err", err)
		return
	}

	// handle storing db backup, overwrite or delete existing as required
	backupDir := filepath.Join("/data/regncon/backup", string(interval))
	finalPath, err := utils.RotateBackups(dbPath, backupDir, retention)
	if err != nil {
		b.Logger.Error("Failed to finalize backup", "err", err)
		return
	}

	// Cleanup temp files after successful backup
	if err := os.Remove(snapshotPath); err != nil {
		b.Logger.Warn("Failed to remove snapshot file", "path", snapshotPath, "err", err)
	}

	// Backup successful
	b.Logger.Info("Backup stored successfully", "path", finalPath)
}

// Hourly triggers a backup task for the hourly interval.
func (b *BackupService) Hourly() {
	ctx := context.Background()
	b.run(ctx, models.Hourly, 24)
}

// Daily triggers a backup task for the daily interval.
func (b *BackupService) Daily() {
	ctx := context.Background()
	b.run(ctx, models.Daily, 7)
}

// Weekly triggers a backup task for the weekly interval.
func (b *BackupService) Weekly() {
	ctx := context.Background()
	b.run(ctx, models.Weekly, 4)
}

// Yearly triggers a backup task for the yearly interval.
func (b *BackupService) Yearly() {
	ctx := context.Background()
	b.run(ctx, models.Yearly, 99)
}

// Manual trigger which is not used with gocron. Mostly used for testing
func (b *BackupService) Manual() {
	ctx := context.Background()
	b.run(ctx, models.Manually, 20)
}

func HandleBackupResult(outcome models.BackupHandlerOptions) {
	// Update log in DB
	err := UpdateLogBackup(outcome.DB, models.BackupLogInput{
		ID:      outcome.Id,
		Status:  outcome.Status,
		Message: outcome.Error,
	})
	if err != nil {
		outcome.Logger.Error("Failed to write to database", "stage", outcome.Stage, "error", err)
	}

	// Log result
	if outcome.Status == models.Success {
		outcome.Logger.Info("Backup stored successfully", "type", outcome.Interval)
	} else {
		outcome.Logger.Error("Backup failed", "stage", outcome.Stage, "type", outcome.Interval, "error", outcome.Error)
	}

	// Send discord notification
	err = SendDiscordMessage(outcome)
	if err != nil {
		outcome.Logger.Error("Discord notification failed", "stage", outcome.Stage, "error", outcome.Error)
	}
}
