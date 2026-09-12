package varsler

import (
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"time"

	webpush "github.com/SherClockHolmes/webpush-go"
	"github.com/go-chi/chi/v5"
)

const maxRequestBodyBytes = 16 << 10

type Service struct {
	db         *sql.DB
	logger     *slog.Logger
	config     Config
	httpClient webpush.HTTPClient
	now        func() time.Time
	wake       time.Duration
}

func New(db *sql.DB, logger *slog.Logger, config Config) (*Service, error) {
	if logger == nil {
		logger = slog.Default()
	}
	logger = logger.With("component", "varsler")
	service := &Service{
		db:     db,
		logger: logger,
		config: config,
		now:    time.Now,
		wake:   time.Second,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
			CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
	}
	if !config.Enabled {
		return service, nil
	}
	if db == nil {
		return nil, errors.New("varsler requires a database when enabled")
	}
	if config.PublicKey == "" || config.PrivateKey == "" || config.Subject == "" {
		return nil, errors.New("enabled varsler requires complete Web Push configuration")
	}
	if err := validateVAPIDKeys(config.PublicKey, config.PrivateKey); err != nil {
		return nil, err
	}
	if err := validateSubject(config.Subject); err != nil {
		return nil, err
	}
	return service, nil
}

func (service *Service) RegisterRoutes(router chi.Router) {
	router.Get("/api/varsler/config", service.getConfig)
	router.Get("/api/varsler/subscriptions", service.getSubscription)
	router.Post("/api/varsler/subscriptions", service.putSubscription)
	router.Delete("/api/varsler/subscriptions", service.deleteSubscription)
}
