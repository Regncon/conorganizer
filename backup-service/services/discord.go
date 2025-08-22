package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Regncon/conorganizer/backup-service/models"
)

type DiscordWebhookPayload struct {
	Content string `json:"content"`
}

func SendDiscordMessage(cfg models.Config, message string) error {
	payload := DiscordWebhookPayload{Content: message}
	webhookUrl := cfg.DISCORD_WEBHOOK_URL

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal webhook payload: %w", err)
	}

	res, err := http.Post(webhookUrl, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to send webhook: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return fmt.Errorf("webhook returned non-successful status: %s", res.Status)
	}

	return nil
}
