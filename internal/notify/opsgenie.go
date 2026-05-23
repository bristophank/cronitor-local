package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const opsGenieAPIURL = "https://api.opsgenie.com/v2/alerts"

// OpsGenieAlerter sends alerts to OpsGenie.
type OpsGenieAlerter struct {
	apiKey  string
	client  *http.Client
	apiURL  string
}

type opsGeniePayload struct {
	Message  string `json:"message"`
	Alias    string `json:"alias"`
	Priority string `json:"priority"`
}

// NewOpsGenieAlerter creates a new OpsGenieAlerter with the given API key.
func NewOpsGenieAlerter(apiKey string) *OpsGenieAlerter {
	return &OpsGenieAlerter{
		apiKey: apiKey,
		client: &http.Client{},
		apiURL: opsGenieAPIURL,
	}
}

// Alert sends an alert to OpsGenie for the given job name.
func (a *OpsGenieAlerter) Alert(jobName string) error {
	payload := opsGeniePayload{
		Message:  fmt.Sprintf("Cron job overdue: %s", jobName),
		Alias:    fmt.Sprintf("cronitor-local-%s", jobName),
		Priority: "P3",
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("opsgenie: marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, a.apiURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("opsgenie: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "GenieKey "+a.apiKey)

	resp, err := a.client.Do(req)
	if err != nil {
		return fmt.Errorf("opsgenie: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("opsgenie: unexpected status %d", resp.StatusCode)
	}
	return nil
}
