package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
)

const readinessLiveCheckTimeout = 250 * time.Millisecond

type readinessState struct {
	db     *sql.DB
	logger *slog.Logger

	mu          sync.RWMutex
	degradedErr error
}

func newReadinessState(db *sql.DB, logger *slog.Logger) *readinessState {
	return &readinessState{
		db:     db,
		logger: logger.With("component", "readiness"),
	}
}

func (state *readinessState) MarkDegraded(err error) {
	if err == nil {
		return
	}

	state.mu.Lock()
	defer state.mu.Unlock()

	if state.degradedErr != nil {
		state.degradedErr = fmt.Errorf("%v; %w", state.degradedErr, err)
		return
	}

	state.degradedErr = err
	state.logger.Warn("application marked not ready", "error", err.Error())
}

func (state *readinessState) DegradedError() error {
	state.mu.RLock()
	defer state.mu.RUnlock()
	return state.degradedErr
}

func (state *readinessState) CheckLive(ctx context.Context) error {
	state.mu.RLock()
	db := state.db
	state.mu.RUnlock()

	if db == nil {
		return fmt.Errorf("database is unavailable")
	}

	liveCheckCtx, cancel := context.WithTimeout(ctx, readinessLiveCheckTimeout)
	defer cancel()

	var one int
	if err := db.QueryRowContext(liveCheckCtx, "SELECT 1;").Scan(&one); err != nil {
		return fmt.Errorf("readiness database check failed: %w", err)
	}
	if one != 1 {
		return fmt.Errorf("readiness database check returned %d", one)
	}

	return nil
}

func mountHealthRoutes(router chi.Router, readiness *readinessState, logger *slog.Logger) {
	logger = logger.With("component", "health")

	router.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok\n"))
	})

	router.Get("/readyz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")

		if err := readiness.DegradedError(); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = w.Write([]byte("not ready\n"))
			return
		}

		if err := readiness.CheckLive(r.Context()); err != nil {
			logger.Warn("readiness check failed", "error", err.Error())
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = w.Write([]byte("not ready\n"))
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok\n"))
	})
}
