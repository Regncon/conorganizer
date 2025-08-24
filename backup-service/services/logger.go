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

func UpdateLogBackup(db *sql.DB, input models.BackupLogInput) error {
	_, err := db.Exec(`
        UPDATE backup_logs
        SET status = ?, message = ?
        WHERE id = ?
    `, input.Status, input.Message, input.ID)
	return err
}
