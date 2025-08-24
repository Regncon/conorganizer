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

	// Middleware
	router.Use(
		middleware.Logger,
		middleware.Recoverer,
	)

	// Static file server
	fs := http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static")))
	router.Handle("/static/*", fs)

	// Routes
	handlers := &Handlers{DB: db, Logger: logger}
	router.Get("/", handlers.IndexHandler)
	/* router.Get("/{interval}", handlers.IntervalHandler) */

	return router
}
