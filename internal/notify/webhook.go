package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// WebhookAlerter sends alert payloads to a configurable HTTP endpoint.
type WebhookAlerter struct {
	URL    string
	client *http.Client
}

type webhookPayload struct {
	JobName   string    `json:"job_name"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

// NewWebhookAlerter creates a WebhookAlerter that posts to the given URL.
func NewWebhookAlerter(url string) *WebhookAlerter {
	return &WebhookAlerter{
		URL: url,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

// Alert sends a JSON POST request to the configured webhook URL.
func (w *WebhookAlerter) Alert(jobName, message string) error {
	payload := webhookPayload{
		JobName:   jobName,
		Message:   message,
		Timestamp: time.Now().UTC(),
	}
	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("webhook: marshal payload: %w", err)
	}
	resp, err := w.client.Post(w.URL, "application/json", bytes.NewReader(bodyBytes))
	if err != nil {
		return fmt.Errorf("webhook: post: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("webhook: unexpected status %d from %s", resp.StatusCode, w.URL)
	}
	return nil
}
