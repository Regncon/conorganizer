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
        SET stage = ?, status = ?, db_prefix = ?, file_path = ?, file_size = ?, message = ?
        WHERE id = ?
    `, options.Stage, options.Status, options.DBPrefix, options.FilePath, options.FileSize, options.Error, options.Id)
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

func (b *FetchLogRes) Logs(interval models.BackupInterval, status models.BackupLogStatus, limit int) ([]models.BackupLog, error) {
	var logs []models.BackupLog
	var args []interface{}

	// Base query
	query := `SELECT id, backup_type, stage, status, file_path, message, created_at FROM backup_logs WHERE 1=1`

	// Optional: filter by interval (backup_type)
	if interval != "" {
		query += ` AND backup_type = ?`
		args = append(args, interval)
	}

	// Optional: filter by status
	if status != "" {
		query += ` AND status = ?`
		args = append(args, status)
	}

	// Sort newest first
	query += ` ORDER BY created_at DESC`

	// Optional: limit number of rows
	if limit > 0 {
		query += ` LIMIT ?`
		args = append(args, limit)
	}

	// Run query
	rows, err := b.DB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query logs: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var log models.BackupLog
		if err := rows.Scan(
			&log.ID,
			&log.BackupType,
			&log.Stage,
			&log.Status,
			&log.FilePath,
			&log.Message,
			&log.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		logs = append(logs, log)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row error: %w", err)
	}

	return logs, nil
}
