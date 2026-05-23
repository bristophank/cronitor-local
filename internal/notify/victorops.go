package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// VictorOpsAlerter sends alerts to VictorOps (Splunk On-Call) via REST endpoint.
type VictorOpsAlerter struct {
	webhookURL string
	client     *http.Client
}

type victorOpsPayload struct {
	MessageType       string `json:"message_type"`
	EntityID          string `json:"entity_id"`
	EntityDisplayName string `json:"entity_display_name"`
	StateMessage      string `json:"state_message"`
	Timestamp         int64  `json:"timestamp"`
}

// NewVictorOpsAlerter creates a new VictorOpsAlerter with the given REST endpoint URL.
func NewVictorOpsAlerter(webhookURL string) *VictorOpsAlerter {
	return &VictorOpsAlerter{
		webhookURL: webhookURL,
		client:     &http.Client{Timeout: 10 * time.Second},
	}
}

// Alert sends a CRITICAL alert to VictorOps for the given job name and message.
func (v *VictorOpsAlerter) Alert(jobName, message string) error {
	payload := victorOpsPayload{
		MessageType:       "CRITICAL",
		EntityID:          jobName,
		EntityDisplayName: fmt.Sprintf("Cron job overdue: %s", jobName),
		StateMessage:      message,
		Timestamp:         time.Now().Unix(),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("victorops: marshal payload: %w", err)
	}

	resp, err := v.client.Post(v.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("victorops: post alert: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("victorops: unexpected status %d", resp.StatusCode)
	}

	return nil
}
