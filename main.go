package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"database/sql"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/sync/errgroup"
	"io/ioutil"
	_ "modernc.org/sqlite"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	db, err := initDB("events.db", "initialize.sql")
	if err != nil {
		logger.Error("Could not initialize DB: %v", err)
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
		router := chi.NewMux()

		router.Use(
			middleware.Logger,
			middleware.Recoverer,
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

func initDB(dbPath string, sqlFile string) (*sql.DB, error) {
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		db, err := sql.Open("sqlite", dbPath)
		if err != nil {
			return nil, fmt.Errorf("failed to open DB: %w", err)
		}

		if err = initializeDatabase(db, sqlFile); err != nil {
			return nil, fmt.Errorf("failed to initialize database: %w", err)
		}

		return db, nil
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open DB: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping DB: %w", err)
	}

	return db, nil
}

func initializeDatabase(db *sql.DB, filename string) error {
	sqlContent, err := loadSQLFile(filename)
	if err != nil {
		return fmt.Errorf("error loading SQL file: %w", err)
	}

	_, err = db.Exec(sqlContent)
	if err != nil {
		return fmt.Errorf("failed to execute SQL commands: %w", err)
	}

	return nil
}

func loadSQLFile(filename string) (string, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("failed to read file %s: %w", filename, err)
	}
	return string(data), nil
}
