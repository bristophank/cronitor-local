package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/example/cronitor-local/internal/job"
)

// GotifyAlerter sends alerts to a self-hosted Gotify server.
type GotifyAlerter struct {
	baseURL  string
	token    string
	priority int
	client   *http.Client
}

type gotifyPayload struct {
	Title    string `json:"title"`
	Message  string `json:"message"`
	Priority int    `json:"priority"`
}

// NewGotifyAlerter creates a new GotifyAlerter.
func NewGotifyAlerter(baseURL, token string, priority int) *GotifyAlerter {
	if priority <= 0 {
		priority = 5
	}
	return &GotifyAlerter{
		baseURL:  baseURL,
		token:    token,
		priority: priority,
		client:   &http.Client{},
	}
}

// Alert sends a notification to Gotify.
func (g *GotifyAlerter) Alert(j *job.Job, message string) error {
	payload := gotifyPayload{
		Title:    fmt.Sprintf("cronitor-local: %s overdue", j.Name),
		Message:  message,
		Priority: g.priority,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("gotify: marshal payload: %w", err)
	}

	url := fmt.Sprintf("%s/message?token=%s", g.baseURL, g.token)
	resp, err := g.client.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("gotify: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("gotify: unexpected status %d", resp.StatusCode)
	}
	return nil
}
