package notify

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds all notification backend configuration loaded from environment.
type Config struct {
	SlackWebhookURL string
	Email           EmailConfig
}

// LoadConfig reads notification settings from environment variables.
func LoadConfig() (*Config, error) {
	cfg := &Config{}
	cfg.SlackWebhookURL = os.Getenv("CRONITOR_SLACK_WEBHOOK")

	cfg.Email.Host = os.Getenv("CRONITOR_SMTP_HOST")
	cfg.Email.Username = os.Getenv("CRONITOR_SMTP_USER")
	cfg.Email.Password = os.Getenv("CRONITOR_SMTP_PASS")
	cfg.Email.From = os.Getenv("CRONITOR_ALERT_FROM")
	cfg.Email.To = os.Getenv("CRONITOR_ALERT_TO")

	portStr := os.Getenv("CRONITOR_SMTP_PORT")
	if portStr == "" {
		portStr = "587"
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("notify: invalid CRONITOR_SMTP_PORT %q: %w", portStr, err)
	}
	cfg.Email.Port = port
	return cfg, nil
}

// BuildAlerter constructs a MultiAlerter from the loaded config,
// including only backends that are fully configured.
func BuildAlerter(cfg *Config) Alerter {
	var alerters []Alerter
	if cfg.SlackWebhookURL != "" {
		alerters = append(alerters, NewSlackAlerter(cfg.SlackWebhookURL))
	}
	if cfg.Email.Host != "" && cfg.Email.To != "" {
		alerters = append(alerters, NewEmailAlerter(cfg.Email))
	}
	if len(alerters) == 0 {
		return &NoopAlerter{}
	}
	return NewMultiAlerter(alerters...)
}
