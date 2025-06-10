package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"database/sql"

	"github.com/Regncon/conorganizer/service"
	"github.com/Regncon/conorganizer/service/userctx"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"golang.org/x/sync/errgroup"
	_ "modernc.org/sqlite"
)

func main() {
	// Set up logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}

	// Parse cli flag for setting db path
	dsn := flag.String("dbp", "", "absolute path to database file")
	flag.Parse()
	if *dsn == "" {
		flag.Usage()
		logger.Error("required arg: use -dbp to specify database absolute file path")
	}

	// Validate dbp directory exists
	dir := filepath.Dir(*dsn)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		logger.Error("Database directory does not exist: %s", dir)
	}

	// Initialize database
	db, err := service.InitDB(*dsn, "initialize.sql")
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

	logger.Info(fmt.Sprintf("Starting Server 0.0.0.0:" + getPort()))
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
		publicRouter := chi.NewMux()

		publicRouter.Use(
			middleware.Logger,
			middleware.Recoverer,
			userctx.IsLoggedInMiddleware(logger),
		)

		publicRouter.Handle("/static/*", http.StripPrefix("/static/", static(logger)))

		cleanup, err := setupRoutes(ctx, logger, publicRouter, db)
		defer cleanup()
		if err != nil {
			return fmt.Errorf("error setting up routes: %w", err)
		}

		srv := &http.Server{
			Addr:    "0.0.0.0:" + port,
			Handler: publicRouter,
		}

		go func() {
			<-ctx.Done()
			srv.Shutdown(context.Background())
		}()

		return srv.ListenAndServe()
	}
}
