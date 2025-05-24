package checkIn

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sahilm/fuzzy"
	"io"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"
)

type Cache struct {
	mu        sync.Mutex
	data      []CheckInTicket
	lastFetch time.Time
	ttl       time.Duration
}

var ticketCache = &Cache{
	ttl: 5 * time.Minute, // Cache TTL set to 5 minutes
}

func (c *Cache) Get(logger *slog.Logger, searchTerm string) ([]CheckInTicket, error) {
	fmt.Printf("search term in cache %q\n", searchTerm)
	c.mu.Lock()
	defer c.mu.Unlock()

	// Check if cache is valid
	if time.Since(c.lastFetch) < c.ttl {
		logger.Info("Returning cached tickets")
		return filterTickets(c.data, searchTerm), nil
	}

	// Fetch new data and update cache
	tickets, err := fetchTicketsFromCheckIn(logger)
	if err != nil {
		return nil, err
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
		combinedSearchString := fmt.Sprintf("%s %s %s %s", ticket.OrderID, ticket.Type, ticket.Email, ticket.Name)
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
		logger.Error("missing CHECKIN_KEY or CHECKIN_SECRET")
		return nil, fmt.Errorf("missing CHECKIN_KEY or CHECKIN_SECRET")
	}

	reqBody, err := json.Marshal(map[string]string{
		"query": query,
	})
	if err != nil {
		logger.Error("Error", "message", err.Error())
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("https://app.checkin.no/graphql?client_id=%s&client_secret=%s", clientID, clientSecret), bytes.NewBuffer(reqBody))
	if err != nil {
		logger.Error("Error", "message", err.Error())
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("Error", "message", err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Error", "message", err.Error())
		return nil, err
	}

	var result queryResult
	if err := json.Unmarshal(body, &result); err != nil {
		logger.Error("Error", "message", err.Error())
		return nil, err
	}

	var tickets []CheckInTicket
	for _, et := range result.Data.EventTickets {
		tickets = append(tickets, CheckInTicket{
			ID:      et.ID,
			OrderID: et.OrderID,
			Type:    et.Category,
			Name:    fmt.Sprintf("%s %s", et.Crm.FirstName, et.Crm.LastName),
			Email:   et.Crm.Email,
			IsAdult: et.Crm.Born < "2007-01-01", // Example logic for determining if adult
		})
	}

	return tickets, nil
}
