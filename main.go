package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"database/sql"

	"github.com/Regncon/conorganizer/service"
	"github.com/Regncon/conorganizer/service/authctx"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"golang.org/x/sync/errgroup"
	_ "modernc.org/sqlite"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}

	dsn := flag.String("dbp", "database/events.db", "absolute path to database file")
	flag.Parse()

	db, dbErr := service.InitDB(*dsn)
	if dbErr != nil {
		logger.Error("Could not initialize DB; starting in degraded mode", "err", dbErr, "dsn", *dsn)
		db = nil
	}
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

	logger.Info("Starting Server 0.0.0.0:" + getPort())
	defer logger.Info("Stopping Server")

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := run(ctx, logger, getPort(), db, dbErr); err != nil {
		logger.Error("Error running server", slog.Any("err", err))
		os.Exit(1)
	}
}

func run(ctx context.Context, logger *slog.Logger, port string, db *sql.DB, dbErr error) error {
	g, ctx := errgroup.WithContext(ctx)
	g.Go(startServer(ctx, logger, port, db, dbErr))
	if err := g.Wait(); err != nil {
		return fmt.Errorf("error running server: %w", err)
	}
	return nil
}

func startServer(ctx context.Context, logger *slog.Logger, port string, db *sql.DB, dbErr error) func() error {
	return func() error {
		router := chi.NewRouter()

		router.Use(
			middleware.Logger,
			middleware.Recoverer,
		)

		if dbErr == nil && db != nil {
			router.Use(authctx.AuthMiddleware(logger))
		}

		router.Handle("/static/*", http.StripPrefix("/static/", static(logger)))

		if dbErr == nil && db != nil {
			cleanup, err := setupRoutes(ctx, logger, router, db)
			if err != nil {
				logger.Error("error setting up routes; falling back to degraded mode", "err", err)
				mountDBErrorRoutes(router, err)
			} else if cleanup != nil {
				defer cleanup()
			}
		} else {
			mountDBErrorRoutes(router, dbErr)
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

func mountDBErrorRoutes(r chi.Router, cause error) {
	errMsg := "The application database could not be opened."
	if cause != nil {
		errMsg = cause.Error()
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintf(w, `<!doctype html>
<html lang="en">
<meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>Database unavailable</title>
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
    <h1>Database unavailable</h1>
    <p>The server is running, but the database is not available. Please check that the file exists, the directory path is correct, and the process has permission to access it.</p>
    <p class="muted"><strong>Reason:</strong> <code>%s</code></p>
  </div>
</body>
</html>`, errMsg)
	})

	r.Get("/", handler)
	r.NotFound(handler.ServeHTTP)
	r.MethodNotAllowed(handler.ServeHTTP)
}
