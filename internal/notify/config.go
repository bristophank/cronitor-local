package notify

import (
	"fmt"
	"os"
)

// Config holds optional alerting destinations loaded from environment variables.
type Config struct {
	SMTPHost     string
	SMTPPort     string
	SMTPFrom     string
	SMTPTo       string
	SlackWebhook string
	WebhookURL   string
}

// LoadConfig reads alerting configuration from environment variables.
func LoadConfig() Config {
	return Config{
		SMTPHost:     os.Getenv("SMTP_HOST"),
		SMTPPort:     os.Getenv("SMTP_PORT"),
		SMTPFrom:     os.Getenv("SMTP_FROM"),
		SMTPTo:       os.Getenv("SMTP_TO"),
		SlackWebhook: os.Getenv("SLACK_WEBHOOK_URL"),
		WebhookURL:   os.Getenv("ALERT_WEBHOOK_URL"),
	}
}

// BuildAlerter constructs a MultiAlerter from the provided Config.
// At least one destination must be configured or an error is returned.
func BuildAlerter(cfg Config) (Alerter, error) {
	var alerters []Alerter

	if cfg.SlackWebhook != "" {
		alerters = append(alerters, NewSlackAlerter(cfg.SlackWebhook))
	}
	if cfg.WebhookURL != "" {
		alerters = append(alerters, NewWebhookAlerter(cfg.WebhookURL))
	}
	if cfg.SMTPHost != "" && cfg.SMTPFrom != "" && cfg.SMTPTo != "" {
		port := cfg.SMTPPort
		if port == "" {
			port = "25"
		}
		alerters = append(alerters, NewEmailAlerter(cfg.SMTPHost+":"+port, cfg.SMTPFrom, cfg.SMTPTo))
	}
	if len(alerters) == 0 {
		return nil, fmt.Errorf("notify: no alerting destinations configured")
	}
	return NewMultiAlerter(alerters...), nil
}
