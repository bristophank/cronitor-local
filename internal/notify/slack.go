package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// SlackAlerter sends alert notifications to a Slack webhook URL.
type SlackAlerter struct {
	webhookURL string
	client     *http.Client
}

// NewSlackAlerter creates a new SlackAlerter with the given webhook URL.
func NewSlackAlerter(webhookURL string) *SlackAlerter {
	return &SlackAlerter{
		webhookURL: webhookURL,
		client:     &http.Client{},
	}
}

type slackPayload struct {
	Text string `json:"text"`
}

// Alert sends a Slack message for the given job name and reason.
func (s *SlackAlerter) Alert(jobName, reason string) error {
	payload := slackPayload{
		Text: fmt.Sprintf(":warning: *cronitor-local alert*\nJob: `%s`\nReason: %s", jobName, reason),
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("slack: marshal payload: %w", err)
	}

	resp, err := s.client.Post(s.webhookURL, "application/json", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("slack: post webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("slack: unexpected status %d", resp.StatusCode)
	}
	return nil
}
