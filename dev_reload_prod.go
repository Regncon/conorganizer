//go:build !dev

package main

import (
	"log/slog"

	"github.com/go-chi/chi/v5"
)

func mountDevReloadRoutes(_ chi.Router, _ *slog.Logger) {}
