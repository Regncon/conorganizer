package services

import (
	"context"
	"io"
	"log/slog"

	"backup-migration/models"
	db "backup-migration/services/db"
	s3svc "backup-migration/services/s3"
)

type DBService interface {
	Open()
	// Open(path string) error
	// ListTables(ctx context.Context) ([]string, error)
}

type S3Service interface {
	Browse(ctx context.Context, prefix string, max int32) ([]s3svc.Object, error)
	Download(ctx context.Context, key, outDir string) (string, error)
	Upload(ctx context.Context, key, localPath string) error
	GetExistingPrefixes(ctx context.Context)
}

type EnvService interface {
	// Read(...), Write(...)
}

type FlyService interface {
	// SetSecrets(...)
}

type Registry struct {
	AppState      models.AppState
	Logger        *slog.Logger
	UILogger      *slog.Logger
	ConsoleWriter io.Writer
	S3            S3Service
	DB            DBService
}

func NewRegistry(logger *slog.Logger) *Registry {
	config := models.Config(logger)
	appState := models.NewAppState()

	registry := &Registry{
		AppState: *appState,
		Logger:   logger,
	}

	dbClient, err := db.NewDBClient(appState, logger)
	if err == nil {
		registry.DB = dbClient
	}

	// Init S3 client
	s3client, err := s3svc.NewClient(config)
	if err == nil {
		registry.S3 = s3client
	}

	return registry
}
