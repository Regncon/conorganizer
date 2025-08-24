package services

import (
	"database/sql"
	"fmt"
	"math"

	"github.com/Regncon/conorganizer/backup-service/models"
)

func NewLogBackup(db *sql.DB, intervalType models.BackupInterval) (int64, error) {
	res, err := db.Exec(`
        INSERT INTO backup_logs (backup_type, stage, status)
        VALUES (?, 'starting', 'pending')
    `, intervalType)

	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func UpdateLogBackup(options models.BackupHandlerOptions) error {
	_, err := options.DB.Exec(`
        UPDATE backup_logs
        SET stage = ?, status = ?, file_path = ?, message = ?
        WHERE id = ?
    `, options.Stage, options.Status, options.FilePath, options.Error, options.Id)
	return err
}

type FetchLogRes struct {
	DB *sql.DB
}

func FetchLog(db *sql.DB) *FetchLogRes {
	return &FetchLogRes{
		DB: db,
	}
}

type BackupStats struct {
	Total       int
	Success     int
	Failed      int
	SuccessRate float64
}

func (b *FetchLogRes) Stats() (BackupStats, error) {
	var stats BackupStats

	// Total backups
	err := b.DB.QueryRow(`SELECT COUNT(*) FROM backup_logs`).Scan(&stats.Total)
	if err != nil {
		return stats, fmt.Errorf("failed to get total backups: %w", err)
	}

	// Successful backups
	err = b.DB.QueryRow(`SELECT COUNT(*) FROM backup_logs WHERE status = 'success'`).Scan(&stats.Success)
	if err != nil {
		return stats, fmt.Errorf("failed to get successful backups: %w", err)
	}

	// Failed backups
	err = b.DB.QueryRow(`SELECT COUNT(*) FROM backup_logs WHERE status = 'error'`).Scan(&stats.Failed)
	if err != nil {
		return stats, fmt.Errorf("failed to get failed backups: %w", err)
	}

	// Calculate success rate safely
	if stats.Total > 0 {
		stats.SuccessRate = math.Round(float64(stats.Success) / float64(stats.Total) * 100)
	}

	return stats, nil
}
