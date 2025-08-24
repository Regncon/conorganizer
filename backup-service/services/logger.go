package services

import (
	"database/sql"

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
