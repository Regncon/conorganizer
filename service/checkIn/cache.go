package checkIn

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/sahilm/fuzzy"
)

type Cache struct {
	mu        sync.Mutex
	data      []CheckInTicket
	lastFetch time.Time
	ttl       time.Duration
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
	Data struct {
		EventTickets []eventTicket `json:"eventTickets"`
	} `json:"data"`
}

var ticketCache = &Cache{
	ttl: 5 * time.Minute, // Cache TTL set to 5 minutes
}

func (c *Cache) Get(logger *slog.Logger, searchTerm string) ([]CheckInTicket, error) {
	baseLogger := logger
	logger = logger.With("component", "checkin_cache")
	c.mu.Lock()
	defer c.mu.Unlock()

	// Check if cache is valid
	if time.Since(c.lastFetch) < c.ttl {
		logger.Debug("Returning cached tickets")
		return filterTickets(c.data, searchTerm), nil
	}

	// Fetch new data and update cache
	tickets, err := fetchTicketsFromCheckIn(baseLogger)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tickets from CheckIn: %w", err)
	}

	c.data = tickets
	c.lastFetch = time.Now()
	return filterTickets(tickets, searchTerm), nil
}

func filterTickets(tickets []CheckInTicket, searchTerm string) []CheckInTicket {
	if searchTerm == "" {
		return tickets
	}

	var ticketStrings []string
	for _, ticket := range tickets {
		combinedSearchString :=
			fmt.Sprintf("%s %s %s %s", strconv.Itoa(ticket.OrderID), ticket.Type, ticket.Email, ticket.FirstName+" "+ticket.LastName)
		ticketStrings = append(ticketStrings, combinedSearchString)
	}

	matches := fuzzy.Find(searchTerm, ticketStrings)
	var filteredTickets []CheckInTicket
	for _, match := range matches {
		filteredTickets = append(filteredTickets, tickets[match.Index])
	}

	return filteredTickets
}

func fetchTicketsFromCheckIn(logger *slog.Logger) ([]CheckInTicket, error) {
	logger = logger.With("component", "checkin")
	query := `{
		eventTickets(customer_id: 13446, id: 109715, onlyCompleted: true) {
			id
			category
			category_id
			crm {
				first_name
				last_name
				id
				email
				born
			}
			order_id
		}
	}`

	clientID := os.Getenv("CHECKIN_KEY")
	clientSecret := os.Getenv("CHECKIN_SECRET")
	if clientID == "" || clientSecret == "" {
		logger.Error("Missing CHECKIN_KEY or CHECKIN_SECRET")
		return nil, fmt.Errorf("missing CHECKIN_KEY or CHECKIN_SECRET environment variables")
	}

	reqBody, err := json.Marshal(map[string]string{
		"query": query,
	})
	if err != nil {
		marshalErr := fmt.Errorf("failed to marshal request body: %w", err)
		logger.Error(marshalErr.Error())
		return nil, marshalErr
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("https://app.checkin.no/graphql?client_id=%s&client_secret=%s", clientID, clientSecret), bytes.NewBuffer(reqBody))
	if err != nil {
		requestErr := fmt.Errorf("failed to create check-in request: %w", err)
		logger.Error(requestErr.Error())
		return nil, requestErr
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		doErr := fmt.Errorf("check-in request failed: %w", err)
		logger.Error(doErr.Error())
		return nil, doErr
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		readErr := fmt.Errorf("failed to read check-in response body: %w", err)
		logger.Error(readErr.Error())
		return nil, readErr
	}

	var result queryResult
	if err := json.Unmarshal(body, &result); err != nil {
		unmarshalErr := fmt.Errorf("failed to unmarshal check-in response: %w", err)
		logger.Error(unmarshalErr.Error())
		return nil, unmarshalErr
	}

	var tickets []CheckInTicket
	for _, et := range result.Data.EventTickets {
		tickets = append(tickets, CheckInTicket{
			ID:        et.ID,
			OrderID:   et.OrderID,
			TypeId:    et.CategoryID,
			Type:      et.Category,
			FirstName: et.Crm.FirstName,
			LastName:  et.Crm.LastName,
			Email:     et.Crm.Email,
			IsOver18:  isOver18(et.Crm.Born),
		})
	}

	return tickets, nil
}

func isOver18(born string) bool {
	//For some F***ing reason, go time.Parse cant take YYYY-MM-DD like a normal programming language.
	parseLayout := "2006-01-02"
	birthDate, err := time.Parse(parseLayout, born)
	if err != nil {
		return false
	}

	regnConDate := time.Date(2025, 10, 10, 0, 0, 0, 0, time.UTC)
	eighteenth := birthDate.AddDate(18, 0, 0)
	return !eighteenth.After(regnConDate)
}
