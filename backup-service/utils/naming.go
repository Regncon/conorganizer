package utils

import (
	"fmt"
	"time"

	"github.com/Regncon/conorganizer/backup-service/models"
)

// GenerateBackupFilename returns a string like "hourly-14.bak" based on the provided interval string for a backup filename.
func GenerateBackupFilename(interval models.BackupInterval) string {
	t := time.Now()

	switch interval {
	case models.Hourly:
		return fmt.Sprintf("hourly-%02d.bak", t.Hour())
	case models.Daily:
		return fmt.Sprintf("daily-%02d.bak", t.Day())
	case models.Weekly:
		_, week := t.ISOWeek()
		return fmt.Sprintf("weekly-%02d.bak", week)
	case models.Yearly:
		return fmt.Sprintf("yearly-%d.bak", t.Year())
	default:
		return fmt.Sprintf("backup-%d.bak", t.Unix())
	}
}
