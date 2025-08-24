package web

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(ctx context.Context, logger *slog.Logger, db *sql.DB) http.Handler {
	router := chi.NewRouter()

	router.Use(
		middleware.Logger,
		middleware.Recoverer,
	)

	// Routes
	router.Get("/", IndexHandler)
	router.Get("/{interval}", IntervalHandler)

	return router
}
