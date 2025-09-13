package services

import (
	"fmt"

	"fyne.io/fyne/v2"
	"github.com/Regncon/conorganizer/backup-migration/config"
	"github.com/Regncon/conorganizer/backup-migration/services/s3"
)

type Registry struct {
	Config config.Config
	App    fyne.App
	Window fyne.Window

	S3 s3.S3Client
}

func NewRegistry() *Registry {
	// Init config
	config := config.NewConfig()

	// Try to load local .env file on runtime
	err := config.LoadEnv(".env")
	if err != nil {
		fmt.Printf(".env init failed: %f\n", err)
	}

	// Init S3 client
	s3client := s3.NewS3Client()

	return &Registry{
		Config: *config,
		S3:     *s3client,
	}
}
