package services

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/Regncon/conorganizer/backup-service/config"
	"github.com/Regncon/conorganizer/backup-service/models"
	"github.com/Regncon/conorganizer/backup-service/utils"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// BackupService handles scheduled backup operations and file rotation logic.
type BackupService struct {
	Config   config.Config
	S3Client *s3.Client
	Logger   *slog.Logger
}

// NewBackupService creates a new instance of BackupService with dependencies injected.
func NewBackupService(cfg config.Config, s3Client *s3.Client, logger *slog.Logger) *BackupService {
	return &BackupService{
		Config:   cfg,
		S3Client: s3Client,
		Logger:   logger,
	}
}

// run is the internal function that performs backup + rotation,
// retention is how many files we need of that interval
func (b *BackupService) run(ctx context.Context, interval models.BackupInterval, retention int) {
	b.Logger.Info("Scheduled backup job triggered", "interval", interval)

	// download snapshot
	snapshotPath, err := DownloadLatestSnapshot(ctx, b.S3Client, b.Config.BUCKET_NAME, b.Config.DB_PREFIX)
	if err != nil {
		b.Logger.Error("Downloading snapshot failed", "err", err)
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
