//go:build dev

package main

import (
	"log/slog"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	datastar "github.com/starfederation/datastar-go/datastar"
)

func mountDevReloadRoutes(router chi.Router, logger *slog.Logger) {
	reloadChan := make(chan struct{}, 1)
	var firstConnection sync.Once

	router.Get("/reload", func(w http.ResponseWriter, r *http.Request) {
		sse := datastar.NewSSE(w, r)
		reload := func() {
			if err := sse.ExecuteScript("window.location.reload()"); err != nil && logger != nil {
				logger.Warn("failed to send dev reload script", "error", err)
			}
		}

		firstConnection.Do(reload)

		select {
		case <-reloadChan:
			reload()
		case <-r.Context().Done():
		}
	})

	router.Get("/hotreload", func(w http.ResponseWriter, _ *http.Request) {
		select {
		case reloadChan <- struct{}{}:
		default:
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})
}
