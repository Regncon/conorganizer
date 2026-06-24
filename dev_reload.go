//go:build dev

package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	datastar "github.com/starfederation/datastar-go/datastar"
)

func mountDevReloadRoutes(router chi.Router, logger *slog.Logger) {
	reloadHub := newDevReloadHub()
	bootID := fmt.Sprintf("%d", time.Now().UnixNano())

	router.Get("/reload", func(w http.ResponseWriter, r *http.Request) {
		sse := datastar.NewSSE(w, r)
		if err := sse.ExecuteScript(devReloadOnNewServerScript(bootID)); err != nil {
			logDevReloadError(logger, err)
			return
		}

		reloadSignal, unsubscribe := reloadHub.subscribe()
		defer unsubscribe()

		for {
			select {
			case <-reloadSignal:
				if err := sse.ExecuteScript("window.location.reload()"); err != nil {
					logDevReloadError(logger, err)
					return
				}
			case <-r.Context().Done():
				return
			}
		}
	})

	router.Get("/hotreload", func(w http.ResponseWriter, _ *http.Request) {
		reloadHub.broadcast()
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})
}

type devReloadHub struct {
	mu          sync.Mutex
	subscribers map[chan struct{}]struct{}
}

func newDevReloadHub() *devReloadHub {
	return &devReloadHub{
		subscribers: make(map[chan struct{}]struct{}),
	}
}

func (h *devReloadHub) subscribe() (<-chan struct{}, func()) {
	reloadSignal := make(chan struct{}, 1)

	h.mu.Lock()
	h.subscribers[reloadSignal] = struct{}{}
	h.mu.Unlock()

	unsubscribe := func() {
		h.mu.Lock()
		if _, ok := h.subscribers[reloadSignal]; ok {
			delete(h.subscribers, reloadSignal)
			close(reloadSignal)
		}
		h.mu.Unlock()
	}

	return reloadSignal, unsubscribe
}

func (h *devReloadHub) broadcast() {
	h.mu.Lock()
	defer h.mu.Unlock()

	for reloadSignal := range h.subscribers {
		select {
		case reloadSignal <- struct{}{}:
		default:
		}
	}
}

func devReloadOnNewServerScript(bootID string) string {
	return fmt.Sprintf(`(() => {
	const key = "conorganizer:dev-reload:%s";
	if (sessionStorage.getItem(key) === "reloaded") {
		return;
	}
	sessionStorage.setItem(key, "reloaded");
	window.location.reload();
})()`, bootID)
}

func logDevReloadError(logger *slog.Logger, err error) {
	if logger != nil {
		logger.Warn("failed to send dev reload script", "error", err)
	}
}
