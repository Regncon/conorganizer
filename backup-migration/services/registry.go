package services

import (
	"context"
	"io"
	"log/slog"

	"backup-migration/models"
	s3svc "backup-migration/services/s3"
)

type DBService interface {
	// Open(path string) error
	// ListTables(ctx context.Context) ([]string, error)
}

type S3Service interface {
	Browse(ctx context.Context, prefix string, max int32) ([]s3svc.Object, error)
	Download(ctx context.Context, key, outDir string) (string, error)
	Upload(ctx context.Context, key, localPath string) error
}

type EnvService interface {
	// Read(...), Write(...)
}

type FlyService interface {
	// SetSecrets(...)
}

type Registry struct {
	Logger        *slog.Logger
	UILogger      *slog.Logger // logs that should appear in the UI console
	ConsoleWriter io.Writer    // raw writer to the console (fmt.Fprintf, etc.)

	S3 S3Service
}

func NewRegistry(logger *slog.Logger) *Registry {
	config := models.Config(logger)
	registry := &Registry{
		Logger: logger,
	}

	// Init S3 client
	s3client, err := s3svc.NewClient(config, registry.Logger)
	if err == nil {
		registry.S3 = s3client
	}

	return registry
}
