package services

import (
	"io"
	"log/slog"
)

type DBService interface {
	// Open(path string) error
	// ListTables(ctx context.Context) ([]string, error)
}

type S3Service interface {
	// List(...), Upload(...), Download(...)
}

type EnvService interface {
	// Read(...), Write(...)
}

type FlyService interface {
	// SetSecrets(...)
}

type Registry struct {
	Logger        *slog.Logger
	UILogger      *slog.Logger
	ConsoleWriter io.Writer
}

func NewRegistry() *Registry {
	return &Registry{}
}
