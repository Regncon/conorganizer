package checkIn

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/sahilm/fuzzy"
)

const (
	checkInCacheTTL        = 5 * time.Minute
	checkInFailureCooldown = 30 * time.Second // Avoid repeat upstream calls while Checkin is unavailable.
	checkInRequestTimeout  = 10 * time.Second
	checkInEndpoint        = "https://app.checkin.no/graphql"
)

type Cache struct {
	mu              sync.Mutex
	data            []CheckInTicket
	lastFetch       time.Time
	lastFailure     time.Time
	lastFailureErr  error
	ttl             time.Duration
	failureCooldown time.Duration
	client          *http.Client
	endpoint        string
	now             func() time.Time
}

type crm struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	ID        int    `json:"id"`
	Email     string `json:"email"`
	Born      string `json:"born"`
}

type eventTicket struct {
	ID         int    `json:"id"`
	Category   string `json:"category"`
	CategoryID int    `json:"category_id"`
	Crm        crm    `json:"crm"`
	OrderID    int    `json:"order_id"`
}

type queryResult struct {
	Data *struct {
		EventTickets *[]eventTicket `json:"eventTickets"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

type checkInFetchError struct {
	outcome string
	status  int
	err     error
}

func (e *checkInFetchError) Error() string { return e.err.Error() }
func (e *checkInFetchError) Unwrap() error { return e.err }

func newTicketCache() *Cache {
	return &Cache{
		ttl:             checkInCacheTTL,
		failureCooldown: checkInFailureCooldown,
		client:          &http.Client{Timeout: checkInRequestTimeout},
		now:             time.Now,
	}
}

var ticketCache = newTicketCache()

func (c *Cache) Get(ctx context.Context, logger *slog.Logger, searchTerm string) (TicketFetchResult, error) {
	logger = logger.With("component", "checkin_cache", "integration", "checkin")
	c.mu.Lock()
	defer c.mu.Unlock()

	now := c.currentTime()
	cacheAvailable := !c.lastFetch.IsZero()
	cacheAge := c.cacheAge(now)
	if cacheAvailable && now.Sub(c.lastFetch) < c.ttl {
		logger.Debug("checkin cache hit", c.logFields(ctx, cacheAvailable, cacheAge, false)...)
		return TicketFetchResult{Tickets: filterTickets(c.data, searchTerm)}, nil
	}

	if !c.lastFailure.IsZero() && now.Sub(c.lastFailure) < c.cooldown() {
		return c.cachedFailureResult(searchTerm, cacheAvailable)
	}

	if err := ctx.Err(); err != nil {
		logger.Debug("checkin refresh canceled", append(c.logFields(ctx, cacheAvailable, cacheAge, cacheAvailable), "outcome", "canceled")...)
		return c.failureResult(searchTerm, cacheAvailable, err)
	}

	logger.Info("checkin refresh started", c.logFields(ctx, cacheAvailable, cacheAge, false)...)
	start := c.currentTime()
	endpoint, err := c.requestEndpoint()
	var tickets []CheckInTicket
	if err == nil {
		tickets, err = fetchTicketsFromCheckIn(ctx, c.httpClient(), endpoint)
	}
	finishedAt := c.currentTime()
	duration := finishedAt.Sub(start)
	if err == nil {
		c.data = tickets
		c.lastFetch = finishedAt
		c.lastFailure = time.Time{}
		c.lastFailureErr = nil
		logger.Info("checkin refresh completed", append(c.logFields(ctx, cacheAvailable, cacheAge, false), "duration_ms", duration.Milliseconds(), "outcome", "success", "ticket_count", len(tickets))...)
		return TicketFetchResult{Tickets: filterTickets(tickets, searchTerm)}, nil
	}

	if outcome, _ := fetchFailureDetails(err); outcome != "canceled" {
		c.lastFailure = finishedAt
		c.lastFailureErr = err
	}
	return c.logRefreshFailure(ctx, logger, searchTerm, cacheAvailable, cacheAge, duration, err)
}

func (c *Cache) logRefreshFailure(ctx context.Context, logger *slog.Logger, searchTerm string, cacheAvailable bool, cacheAge, duration time.Duration, err error) (TicketFetchResult, error) {
	outcome, status := fetchFailureDetails(err)
	fields := append(c.logFields(ctx, cacheAvailable, cacheAge, cacheAvailable), "duration_ms", duration.Milliseconds(), "outcome", outcome)
	if status != 0 {
		fields = append(fields, "upstream_status_code", status)
	}
	if outcome == "canceled" {
		logger.Debug("checkin refresh canceled", fields...)
	} else if cacheAvailable {
		logger.Warn("checkin refresh failed; using stale cache", fields...)
	} else {
		logger.Error("checkin refresh failed", fields...)
	}
	return c.failureResult(searchTerm, cacheAvailable, err)
}

func (c *Cache) cachedFailureResult(searchTerm string, cacheAvailable bool) (TicketFetchResult, error) {
	err := c.lastFailureErr
	if err == nil {
		err = errors.New("checkin refresh unavailable")
	}
	return c.failureResult(searchTerm, cacheAvailable, err)
}

func (c *Cache) failureResult(searchTerm string, cacheAvailable bool, err error) (TicketFetchResult, error) {
	if cacheAvailable {
		return TicketFetchResult{Tickets: filterTickets(c.data, searchTerm), UsedStaleCache: true}, err
	}
	return TicketFetchResult{}, err
}

func (c *Cache) currentTime() time.Time {
	if c.now != nil {
		return c.now()
	}
	return time.Now()
}

func (c *Cache) cacheAge(now time.Time) time.Duration {
	if c.lastFetch.IsZero() {
		return 0
	}
	return now.Sub(c.lastFetch)
}

func (c *Cache) cooldown() time.Duration {
	if c.failureCooldown > 0 {
		return c.failureCooldown
	}
	return checkInFailureCooldown
}

func (c *Cache) httpClient() *http.Client {
	if c.client != nil {
		return c.client
	}
	return &http.Client{Timeout: checkInRequestTimeout}
}

func (c *Cache) requestEndpoint() (string, error) {
	if c.endpoint != "" {
		return c.endpoint, nil
	}

	clientID := os.Getenv("CHECKIN_KEY")
	clientSecret := os.Getenv("CHECKIN_SECRET")
	if clientID == "" || clientSecret == "" {
		return "", &checkInFetchError{outcome: "http_error", err: errors.New("missing Checkin credentials")}
	}

	endpoint, err := url.Parse(checkInEndpoint)
	if err != nil {
		return "", &checkInFetchError{outcome: "http_error", err: errors.New("invalid Checkin endpoint")}
	}
	query := endpoint.Query()
	query.Set("client_id", clientID)
	query.Set("client_secret", clientSecret)
	endpoint.RawQuery = query.Encode()
	return endpoint.String(), nil
}

func (c *Cache) logFields(ctx context.Context, cacheAvailable bool, cacheAge time.Duration, usedStaleCache bool) []any {
	fields := []any{"cache_available", cacheAvailable, "used_stale_cache", usedStaleCache}
	if cacheAvailable {
		fields = append(fields, "cache_age_ms", cacheAge.Milliseconds())
	}
	if requestID := middleware.GetReqID(ctx); requestID != "" {
		fields = append(fields, "request_id", requestID)
	}
	return fields
}

func filterTickets(tickets []CheckInTicket, searchTerm string) []CheckInTicket {
	if searchTerm == "" {
		return tickets
	}

	var ticketStrings []string
	for _, ticket := range tickets {
		ticketStrings = append(ticketStrings, fmt.Sprintf("%s %s %s %s", strconv.Itoa(ticket.OrderID), ticket.Type, ticket.Email, ticket.FirstName+" "+ticket.LastName))
	}

	matches := fuzzy.Find(searchTerm, ticketStrings)
	var filteredTickets []CheckInTicket
	for _, match := range matches {
		filteredTickets = append(filteredTickets, tickets[match.Index])
	}
	return filteredTickets
}

func fetchTicketsFromCheckIn(ctx context.Context, client *http.Client, endpoint string) ([]CheckInTicket, error) {
	query := `{
		eventTickets(customer_id: 13446, id: 221572, onlyCompleted: true) {
			id category category_id
			crm { first_name last_name id email born }
			order_id
		}
	}`
	reqBody, err := json.Marshal(map[string]string{"query": query})
	if err != nil {
		return nil, &checkInFetchError{outcome: "decode_error", err: errors.New("failed to prepare Checkin request")}
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, &checkInFetchError{outcome: "http_error", err: errors.New("failed to prepare Checkin request")}
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, &checkInFetchError{outcome: requestOutcome(err), err: sanitizeNetworkError(err)}
	}
	defer resp.Body.Close()
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, &checkInFetchError{outcome: "http_error", status: resp.StatusCode, err: fmt.Errorf("Checkin returned HTTP status %d", resp.StatusCode)}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &checkInFetchError{outcome: requestOutcome(err), err: sanitizeNetworkError(err)}
	}
	var result queryResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, &checkInFetchError{outcome: "decode_error", err: fmt.Errorf("failed to decode Checkin response: %w", err)}
	}
	if len(result.Errors) > 0 {
		return nil, &checkInFetchError{outcome: "graphql_error", err: errors.New("Checkin returned GraphQL errors")}
	}
	if result.Data == nil || result.Data.EventTickets == nil {
		return nil, &checkInFetchError{outcome: "decode_error", err: errors.New("Checkin response did not contain ticket data")}
	}

	tickets := make([]CheckInTicket, 0, len(*result.Data.EventTickets))
	for _, et := range *result.Data.EventTickets {
		tickets = append(tickets, CheckInTicket{ID: et.ID, OrderID: et.OrderID, TypeId: et.CategoryID, Type: et.Category, FirstName: et.Crm.FirstName, LastName: et.Crm.LastName, Email: et.Crm.Email, IsOver18: isOver18(et.Crm.Born)})
	}
	return tickets, nil
}

func fetchFailureDetails(err error) (string, int) {
	var fetchErr *checkInFetchError
	if errors.As(err, &fetchErr) {
		return fetchErr.outcome, fetchErr.status
	}
	return requestOutcome(err), 0
}

func requestOutcome(err error) string {
	if errors.Is(err, context.Canceled) {
		return "canceled"
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return "timeout"
	}
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return "timeout"
	}
	return "http_error"
}

func sanitizeNetworkError(err error) error {
	var urlErr *url.Error
	if errors.As(err, &urlErr) && urlErr.Err != nil {
		return urlErr.Err
	}
	return err
}

func isOver18(born string) bool {
	birthDate, err := time.Parse("2006-01-02", born)
	if err != nil {
		return false
	}
	regnConDate := time.Date(2026, 10, 4, 0, 0, 0, 0, time.UTC)
	return !birthDate.AddDate(18, 0, 0).After(regnConDate)
}
