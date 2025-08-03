package services

import (
	"log/slog"
	"time"

	"github.com/go-co-op/gocron/v2"
)

// StartScheduler initializes and starts the backup scheduler service.
func StartScheduler(backupService *BackupService, logger *slog.Logger) error {
	logger.Info("Scheduler is starting registering jobs...")

	location, _ := time.LoadLocation("Europe/Oslo")
	scheduler, err := gocron.NewScheduler(
		gocron.WithLocation(location),
	)
	if err != nil {
		return err
	}

	// Register all backup jobs
	if err := registerJobs(scheduler, *backupService, *logger); err != nil {
		return err
	}

	// Start background job runner
	scheduler.Start()
	logger.Info("Scheduler finished registering all jobs successfully!")

	return nil
}

// Registers jobs for hourly, daily, weekly and yearly backups
func registerJobs(scheduler gocron.Scheduler, backupService BackupService, logger slog.Logger) error {
	jobs := []struct {
		name     string
		schedule string
		task     func()
	}{
		{"hourly", "0 * * * *", backupService.Hourly},
		{"daily", "15 0 * * *", backupService.Daily},
		{"weekly", "30 0 * * 0", backupService.Weekly},
		{"yearly", "45 0 1 1 *", backupService.Yearly},
	}

	for _, job := range jobs {
		_, err := scheduler.NewJob(
			gocron.CronJob(job.schedule, false),
			gocron.NewTask(job.task),
			gocron.WithName(job.name),
		)

		if err != nil {
			return err
		}

		logger.Info("Scheduled new job", "type", job.name, "shedule", job.schedule)
	}

	return nil
}
