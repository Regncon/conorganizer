package services

import (
	"database/sql"

	"github.com/Regncon/conorganizer/backup-service/models"
)

func NewLogBackup(db *sql.DB, intervalType models.BackupInterval) (int64, error) {
	res, err := db.Exec(`
        INSERT INTO backup_logs (backup_type, status, message)
        VALUES (?, 'pending', '')
    `, intervalType)

	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func UpdateLogBackup(db *sql.DB, input models.BackupLogInput) error {
	msg := ""
	if input.Message != nil {
		msg = input.Message.Error()
	}

	_, err := db.Exec(`
        UPDATE backup_logs
        SET status = ?, message = ?
        WHERE id = ?
    `, input.Status, msg, input.ID)
	return err
}
