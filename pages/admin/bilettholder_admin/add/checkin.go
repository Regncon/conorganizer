package addbilettholder

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
)

type CheckInTicket struct {
	OrderID int
	Type    string
	Name    string
	Email   string
	IsAdult bool
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

func GetTicketsFromCheckIn(logger *slog.Logger) ([]CheckInTicket, error) {
	query := `{
		eventTickets(customer_id: 13446, id: 73685, onlyCompleted: true) {
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
			OrderID: et.OrderID,
			Type:    et.Category,
			Name:    fmt.Sprintf("%s %s", et.Crm.FirstName, et.Crm.LastName),
			Email:   et.Crm.Email,
			IsAdult: et.Crm.Born < "2007-01-01", // Example logic for determining if adult
		})
	}

	return tickets, nil
}
