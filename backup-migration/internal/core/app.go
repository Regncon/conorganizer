package core

import (
	"context"

	"github.com/Regncon/conorganizer/backup-migration/internal/gui"
	"github.com/Regncon/conorganizer/backup-migration/services"
)

func NewApp(ctx context.Context, reg *services.Registry) error {
	// Initiate GUI
	gui.NewFyneApp(ctx, reg)

	return nil
}
