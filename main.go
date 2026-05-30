package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Regncon/conorganizer/service"
	"github.com/Regncon/conorganizer/service/applog"
	"github.com/Regncon/conorganizer/service/authctx"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"golang.org/x/sync/errgroup"
	_ "modernc.org/sqlite"
)

func main() {
	logger := applog.NewJSONLogger()
	slog.SetDefault(logger)
	baseLogger := logger
	logger = logger.With("component", "main")

	if err := godotenv.Load(); err != nil {
		logger.Debug("No .env file found")
	}

	dsn := flag.String("dbp", "database/events.db", "absolute path to database file")
	eventImageDir := flag.String("image-path", "local-event-images", "directory to store event images")
	flag.Parse()

	db, dbErr := service.InitDB(*dsn)
	if dbErr != nil {
		logger.Error(fmt.Errorf("could not initialize DB at %q: %w", *dsn, dbErr).Error())
		os.Exit(1)
	}
	logger.Info("SQLite database initialized", "path", *dsn, "max_open_connections", db.Stats().MaxOpenConnections)
	defer func() {
		if db != nil {
			db.Close()
		}
	}()

	getPort := func() string {
		if p, ok := os.LookupEnv("PORT"); ok {
			return p
		}
		return "8080"
	}

	logger.Info("Starting server", "address", "0.0.0.0:"+getPort())
	defer logger.Info("Stopping server")

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := run(ctx, baseLogger, getPort(), eventImageDir, db); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

func run(ctx context.Context, logger *slog.Logger, port string, eventImageDir *string, db *sql.DB) error {
	g, ctx := errgroup.WithContext(ctx)
	g.Go(startServer(ctx, logger, port, eventImageDir, db))
	if err := g.Wait(); err != nil {
		return fmt.Errorf("error running server: %w", err)
	}
	return nil
}
func startServer(ctx context.Context, logger *slog.Logger, port string, eventImageDir *string, db *sql.DB) func() error {
	return func() error {
		baseLogger := logger
		logger = logger.With("component", "http_server")
		router := chi.NewRouter()
		readiness := newReadinessState(db, baseLogger)

		router.Use(
			middleware.RequestID,
			RequestLoggingMiddleware(baseLogger.With("component", "http")),
			middleware.Recoverer,
		)

		mountHealthRoutes(router, readiness, baseLogger)

		var imgErr error
		if eventImageDir != nil && *eventImageDir != "" {
			if err := service.CheckWritableDirectory(*eventImageDir); err != nil {
				imgErr = fmt.Errorf("event image directory startup check failed: %w", err)
			}
		} else {
			imgErr = fmt.Errorf("event image directory path is empty")
		}

		degradedErr := imgErr
		fullMode := degradedErr == nil && db != nil
		if degradedErr != nil {
			readiness.MarkDegraded(degradedErr)
		}

		if fullMode {
			router.Use(authctx.AuthMiddleware(baseLogger))
		}

		if eventImageDir != nil && *eventImageDir != "" {
			router.Handle("/event-images/*", http.StripPrefix("/event-images/", http.FileServer(http.Dir(*eventImageDir))))
		}
		router.Handle("/static/*", http.StripPrefix("/static/", static(baseLogger)))

		if fullMode {
			cleanup, err := setupRoutes(ctx, baseLogger, router, db, eventImageDir)
			if err != nil {
				logger.Error(fmt.Errorf("error setting up routes; falling back to degraded mode: %w", err).Error())
				readiness.MarkDegraded(err)
				mountDBErrorRoutes(router, err)
			} else if cleanup != nil {
				defer func() {
					if err := cleanup(); err != nil {
						logger.Error(fmt.Errorf("failed to clean up routes: %w", err).Error())
					}
				}()
			}
		} else {
			// Show a single degraded page without exposing operational details.
			mountDBErrorRoutes(router, degradedErr)
		}

		srv := &http.Server{
			Addr:    "0.0.0.0:" + port,
			Handler: router,
		}

		go func() {
			<-ctx.Done()
			_ = srv.Shutdown(context.Background())
		}()

		return srv.ListenAndServe()
	}
}

func mountDBErrorRoutes(r chi.Router, _ error) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprint(w, `<!doctype html>
<html lang="en">
<meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>Application unavailable</title>
<style>
  :root { color-scheme: light dark; }
  body { font-family: ui-sans-serif, system-ui, -apple-system, Segoe UI, Roboto, Ubuntu, Cantarell, Noto Sans, Helvetica, Arial, Apple Color Emoji, Segoe UI Emoji; margin:0; padding:2rem; }
  .card { max-width: 56ch; margin: 5vh auto; border: 1px solid rgba(127,127,127,.35); border-radius: .75rem; padding: 1.5rem; }
  h1 { margin: 0 0 .5rem; }
  code { padding: .15rem .35rem; border-radius: .35rem; background: rgba(127,127,127,.15); }
  .muted { opacity: .8; }
</style>
<body>
  <div class="card">
    <h1>Application unavailable</h1>
    <p>The server is running, but the application is not ready. Please try again later.</p>
    <p class="muted">Operational details are available in the service logs.</p>
  </div>
</body>
</html>`)
	})

	r.Get("/", handler)
	r.NotFound(handler.ServeHTTP)
	r.MethodNotAllowed(handler.ServeHTTP)
}
