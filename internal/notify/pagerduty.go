package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const pagerDutyEventsURL = "https://events.pagerduty.com/v2/enqueue"

// PagerDutyAlerter sends alerts to PagerDuty via the Events API v2.
type PagerDutyAlerter struct {
	integrationKey string
	client         *http.Client
}

// NewPagerDutyAlerter creates a new PagerDutyAlerter with the given integration key.
func NewPagerDutyAlerter(integrationKey string) *PagerDutyAlerter {
	return &PagerDutyAlerter{
		integrationKey: integrationKey,
		client:         &http.Client{Timeout: 10 * time.Second},
	}
}

type pdPayload struct {
	RoutingKey  string    `json:"routing_key"`
	EventAction string    `json:"event_action"`
	Payload     pdDetails `json:"payload"`
}

type pdDetails struct {
	Summary  string `json:"summary"`
	Source   string `json:"source"`
	Severity string `json:"severity"`
}

// Alert sends a trigger event to PagerDuty for the given job name and message.
func (p *PagerDutyAlerter) Alert(jobName, message string) error {
	body := pdPayload{
		RoutingKey:  p.integrationKey,
		EventAction: "trigger",
		Payload: pdDetails{
			Summary:  fmt.Sprintf("cronitor-local: %s — %s", jobName, message),
			Source:   "cronitor-local",
			Severity: "error",
		},
	}

	data, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("pagerduty: marshal payload: %w", err)
	}

	resp, err := p.client.Post(pagerDutyEventsURL, "application/json", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("pagerduty: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("pagerduty: unexpected status %d", resp.StatusCode)
	}
	return nil
}
