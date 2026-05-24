package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/cronitor-local/internal/job"
)

// DiscordAlerter sends alert notifications to a Discord channel via webhook.
type DiscordAlerter struct {
	webhookURL string
	client     *http.Client
}

type discordPayload struct {
	Content  string         `json:"content,omitempty"`
	Embeds   []discordEmbed `json:"embeds,omitempty"`
}

type discordEmbed struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Color       int    `json:"color"`
}

// NewDiscordAlerter creates a new DiscordAlerter with the given webhook URL.
func NewDiscordAlerter(webhookURL string) *DiscordAlerter {
	return &DiscordAlerter{
		webhookURL: webhookURL,
		client:     &http.Client{},
	}
}

// Alert sends a notification to Discord about an overdue job.
func (d *DiscordAlerter) Alert(j *job.Job, message string) error {
	payload := discordPayload{
		Embeds: []discordEmbed{
			{
				Title:       fmt.Sprintf("⚠️ Cron Alert: %s", j.Name),
				Description: message,
				Color:       15158332, // red
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("discord: marshal payload: %w", err)
	}

	resp, err := d.client.Post(d.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("discord: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("discord: unexpected status %d", resp.StatusCode)
	}

	return nil
}
