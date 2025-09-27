package main

import (
	"context"
	"log"

	"github.com/Regncon/conorganizer/backup-migration/internal/core"
	"github.com/Regncon/conorganizer/backup-migration/services"
)

func main() {
	// Create registry
	reg := services.NewRegistry()

	// Run app entrypoint
	if err := core.NewApp(context.Background(), reg); err != nil {
		log.Fatal(err)
	}
}
