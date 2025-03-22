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
	"io/ioutil"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/supertokens/supertokens-golang/recipe/dashboard"
	"github.com/supertokens/supertokens-golang/recipe/dashboard/dashboardmodels"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/recipe/userroles"
	"github.com/supertokens/supertokens-golang/supertokens"
	"golang.org/x/sync/errgroup"
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
		initSupertokens()
		router := chi.NewMux()

		router.Use(cors.Handler(cors.Options{
			AllowedOrigins: []string{"http://localhost:8080"},
			AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders: append([]string{"Content-Type"},
				supertokens.GetAllCORSHeaders()...),
			AllowCredentials: true,
		}))

		router.Use(
			middleware.Logger,
			middleware.Recoverer,
			supertokens.Middleware,
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

func setupThirdParty() supertokens.Recipe {
	return thirdparty.Init(&tpmodels.TypeInput{
		SignInAndUpFeature: tpmodels.TypeInputSignInAndUp{
			Providers: []tpmodels.ProviderInput{
				{
					Config: tpmodels.ProviderConfig{
						ThirdPartyId: "google",
						Clients: []tpmodels.ProviderClientConfig{
							{
								ClientID:     "1060725074195-kmeum4crr01uirfl2op9kd5acmi9jutn.apps.googleusercontent.com",
								ClientSecret: "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",
							},
						},
					},
				},
				{
					Config: tpmodels.ProviderConfig{
						ThirdPartyId: "discord",
						Clients: []tpmodels.ProviderClientConfig{
							{
								ClientID:     "TODO:",
								ClientSecret: "TODO:",
							},
						},
					},
				},
				{
					Config: tpmodels.ProviderConfig{
						ThirdPartyId: "facebook",
						Clients: []tpmodels.ProviderClientConfig{
							{
								ClientID:     "TODO:",
								ClientSecret: "TODO:",
							},
						},
					},
				},
				{
					Config: tpmodels.ProviderConfig{
						ThirdPartyId: "twitter",
						Clients: []tpmodels.ProviderClientConfig{
							{
								ClientID:     "4398792-WXpqVXRiazdRMGNJdEZIa3RVQXc6MTpjaQ",
								ClientSecret: "BivMbtwmcygbRLNQ0zk45yxvW246tnYnTFFq-LH39NwZMxFpdC",
							},
						},
					},
				},
				{
					Config: tpmodels.ProviderConfig{
						ThirdPartyId: "active-directory",
						Clients: []tpmodels.ProviderClientConfig{
							{
								ClientID:     "TODO:",
								ClientSecret: "TODO:",
							},
						},
						OIDCDiscoveryEndpoint: "https://login.microsoftonline.com/<directoryId>/v2.0/.well-known/openid-configuration",
					},
				},
				{
					Config: tpmodels.ProviderConfig{
						ThirdPartyId: "apple",
						Clients: []tpmodels.ProviderClientConfig{
							{
								ClientID: "4398792-io.supertokens.example.service",
								AdditionalConfig: map[string]interface{}{
									"keyId":      "7M48Y4RYDL",
									"privateKey": "-----BEGIN PRIVATE KEY-----\nMIGTAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBHkwdwIBAQQgu8gXs+XYkqXD6Ala9Sf/iJXzhbwcoG5dMh1OonpdJUmgCgYIKoZIzj0DAQehRANCAASfrvlFbFCYqn3I2zeknYXLwtH30JuOKestDbSfZYxZNMqhF/OzdZFTV0zc5u5s3eN+oCWbnvl0hM+9IW0UlkdA\n-----END PRIVATE KEY-----",
									"teamId":     "YWQCXGJRJL",
								},
							},
						},
					},
				},
			},
		},
	})
}

func initSupertokens() {
	apiBasePath := "/auth"
	websiteBasePath := "/auth"
	err := supertokens.Init(supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			// We use try.supertokens for demo purposes.
			// At the end of the tutorial we will show you how to create
			// your own SuperTokens core instance and then update your config.
			ConnectionURI: "https://st-dev-2eb6cde0-0683-11f0-9a7f-5d160be73652.aws.supertokens.io",
			APIKey:        "01w-Ccsm64UrzH1pjiAU9x=H1V",
		},
		AppInfo: supertokens.AppInfo{
			AppName:         "Regncon 2025",
			APIDomain:       "http://localhost:8080",
			WebsiteDomain:   "http://localhost:8080",
			APIBasePath:     &apiBasePath,
			WebsiteBasePath: &websiteBasePath,
		},
		RecipeList: []supertokens.Recipe{
			emailpassword.Init(nil),
			session.Init(nil),
			setupThirdParty(),
			dashboard.Init(&dashboardmodels.TypeInput{
				Admins: &[]string{
					"coldasice_t@hotmail.com",
				},
			}),
			userroles.Init(nil),
		},
	})

	if err != nil {
		panic(err.Error())
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
