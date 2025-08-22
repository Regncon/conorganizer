package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Regncon/conorganizer/backup-service/models"
)

func SendDiscordMessage(outcome models.BackupOutcome) error {
	payload := models.DiscordWebhookPayload{
		Content: ":face_with_peeking_eye:",
		Embeds: []models.DiscordEmbed{
			{
				Title: "Regncon backup-service error!",
				Color: 16711680,
				Fields: []models.DiscordEmbedField{
					{
						Name:   "Stage",
						Value:  string(outcome.Stage),
						Inline: true,
					},
					{
						Name:   "Type",
						Value:  string(outcome.Interval),
						Inline: true,
					},
					{},
					{
						Name:   "When",
						Value:  time.Now().Format("Mon, 02 Jan 15:04:05"),
						Inline: true,
					}, {
						Name:   "Link",
						Value:  "[Coming soon!](https://google.com)",
						Inline: true,
					},
					{
						Name:  "Error message",
						Value: outcome.Error,
					},
				},
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal webhook payload: %w", err)
	}

	res, err := http.Post(outcome.WebhookURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to send webhook: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return fmt.Errorf("webhook returned non-successful status: %s", res.Status)
	}

	return nil
}
