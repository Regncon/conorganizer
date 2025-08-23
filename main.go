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

	db, err := service.InitDB(*dsn)
	if err != nil {
		logger.Error("Could not initialize DB", "initialize database", err)
	}
	defer db.Close()

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

	if err := run(ctx, logger, getPort(), db); err != nil {
		logger.Error("Error running server", slog.Any("err", err))
		os.Exit(1)
	}
}

func run(ctx context.Context, logger *slog.Logger, port string, db *sql.DB) error {
	g, ctx := errgroup.WithContext(ctx)

	g.Go(startServer(ctx, logger, port, db))

	if err := g.Wait(); err != nil {
		return fmt.Errorf("error running server: %w", err)
	}

	return nil
}

func startServer(ctx context.Context, logger *slog.Logger, port string, db *sql.DB) func() error {
	return func() error {
		router := chi.NewRouter()

		router.Use(
			middleware.Logger,
			middleware.Recoverer,
			authctx.AuthMiddleware(logger),
		)

		router.Handle("/static/*", http.StripPrefix("/static/", static(logger)))

		cleanup, err := setupRoutes(ctx, logger, router, db)
		defer cleanup()
		if err != nil {
			return fmt.Errorf("error setting up routes: %w", err)
		}

		srv := &http.Server{
			Addr:    "0.0.0.0:" + port,
			Handler: router,
		}

		go func() {
			<-ctx.Done()
			srv.Shutdown(context.Background())
		}()

		return srv.ListenAndServe()
	}
}
