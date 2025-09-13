package services

import (
	"fmt"

	"fyne.io/fyne/v2"
	"github.com/Regncon/conorganizer/backup-migration/config"
)

type Registry struct {
	Config config.Config
	App    fyne.App
}

func NewRegistry() *Registry {
	// Load application config
	config := config.NewConfig()

	// Try to load local .env file on runtime
	err := config.LoadEnv("../.env")
	if err != nil {
		fmt.Printf(".env init failed: %f", err)
	}

	return &Registry{
		Config: *config,
	}
}
