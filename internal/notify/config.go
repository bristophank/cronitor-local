package notify

import (
	"errors"
	"os"
	"strings"
)

// Config holds notification settings loaded from environment variables.
type Config struct {
	// Webhook
	WebhookURL string

	// Slack
	SlackWebhookURL string

	// Email (SMTP)
	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
	EmailFrom    string
	EmailTo      []string
}

// LoadConfig reads notification configuration from environment variables.
func LoadConfig() Config {
	var to []string
	if raw := os.Getenv("ALERT_EMAIL_TO"); raw != "" {
		for _, addr := range strings.Split(raw, ",") {
			if t := strings.TrimSpace(addr); t != "" {
				to = append(to, t)
			}
		}
	}
	return Config{
		WebhookURL:      os.Getenv("ALERT_WEBHOOK_URL"),
		SlackWebhookURL: os.Getenv("ALERT_SLACK_WEBHOOK_URL"),
		SMTPHost:        os.Getenv("ALERT_SMTP_HOST"),
		SMTPPort:        os.Getenv("ALERT_SMTP_PORT"),
		SMTPUsername:    os.Getenv("ALERT_SMTP_USERNAME"),
		SMTPPassword:    os.Getenv("ALERT_SMTP_PASSWORD"),
		EmailFrom:       os.Getenv("ALERT_EMAIL_FROM"),
		EmailTo:         to,
	}
}

// BuildAlerter constructs a MultiAlerter from the provided Config.
// Returns an error if no alerters are configured.
func BuildAlerter(cfg Config) (Alerter, error) {
	var alerters []Alerter

	if cfg.WebhookURL != "" {
		alerters = append(alerters, NewWebhookAlerter(cfg.WebhookURL))
	}
	if cfg.SlackWebhookURL != "" {
		alerters = append(alerters, NewSlackAlerter(cfg.SlackWebhookURL))
	}
	if cfg.SMTPHost != "" && cfg.EmailFrom != "" && len(cfg.EmailTo) > 0 {
		port := cfg.SMTPPort
		if port == "" {
			port = "587"
		}
		alerters = append(alerters, NewEmailAlerter(
			cfg.SMTPHost, port, cfg.SMTPUsername, cfg.SMTPPassword,
			cfg.EmailFrom, cfg.EmailTo,
		))
	}

	if len(alerters) == 0 {
		return nil, errors.New("no alerters configured: set at least one of ALERT_WEBHOOK_URL, ALERT_SLACK_WEBHOOK_URL, or SMTP settings")
	}
	return NewMultiAlerter(alerters...), nil
}
