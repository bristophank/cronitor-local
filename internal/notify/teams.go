package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// TeamsAlerter sends alerts to a Microsoft Teams channel via an incoming webhook.
type TeamsAlerter struct {
	webhookURL string
	client     *http.Client
}

// NewTeamsAlerter creates a new TeamsAlerter with the given webhook URL.
func NewTeamsAlerter(webhookURL string) *TeamsAlerter {
	return &TeamsAlerter{
		webhookURL: webhookURL,
		client:     &http.Client{},
	}
}

type teamsPayload struct {
	Type       string         `json:"@type"`
	Context    string         `json:"@context"`
	ThemeColor string         `json:"themeColor"`
	Summary    string         `json:"summary"`
	Sections   []teamsSection `json:"sections"`
}

type teamsSection struct {
	ActivityTitle string `json:"activityTitle"`
	ActivityText  string `json:"activityText"`
}

// Alert sends an alert message to the configured Microsoft Teams channel.
func (a *TeamsAlerter) Alert(jobName, message string) error {
	payload := teamsPayload{
		Type:       "MessageCard",
		Context:    "http://schema.org/extensions",
		ThemeColor: "FF0000",
		Summary:    fmt.Sprintf("cronitor-local: job %q is overdue", jobName),
		Sections: []teamsSection{
			{
				ActivityTitle: fmt.Sprintf("Job Alert: %s", jobName),
				ActivityText:  message,
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("teams: marshal payload: %w", err)
	}

	resp, err := a.client.Post(a.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("teams: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("teams: unexpected status code %d", resp.StatusCode)
	}

	return nil
}
