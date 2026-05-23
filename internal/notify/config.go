package notify

import (
	"errors"
	"os"
)

// Config holds alerting configuration loaded from environment variables.
type Config struct {
	WebhookURL      string
	SlackWebhookURL string
	EmailSMTP       string
	EmailFrom       string
	EmailTo         string
	PagerDutyKey    string
	OpsGenieKey     string
}

// LoadConfig reads alerting configuration from environment variables.
func LoadConfig() Config {
	return Config{
		WebhookURL:      os.Getenv("ALERT_WEBHOOK_URL"),
		SlackWebhookURL: os.Getenv("ALERT_SLACK_WEBHOOK_URL"),
		EmailSMTP:       os.Getenv("ALERT_EMAIL_SMTP"),
		EmailFrom:       os.Getenv("ALERT_EMAIL_FROM"),
		EmailTo:         os.Getenv("ALERT_EMAIL_TO"),
		PagerDutyKey:    os.Getenv("ALERT_PAGERDUTY_KEY"),
		OpsGenieKey:     os.Getenv("ALERT_OPSGENIE_KEY"),
	}
}

// BuildAlerter constructs a MultiAlerter from the given Config.
// Returns an error if no alerting method is configured.
func BuildAlerter(cfg Config) (Alerter, error) {
	var alerters []Alerter

	if cfg.WebhookURL != "" {
		alerters = append(alerters, NewWebhookAlerter(cfg.WebhookURL))
	}
	if cfg.SlackWebhookURL != "" {
		alerters = append(alerters, NewSlackAlerter(cfg.SlackWebhookURL))
	}
	if cfg.EmailSMTP != "" && cfg.EmailFrom != "" && cfg.EmailTo != "" {
		alerters = append(alerters, NewEmailAlerter(cfg.EmailSMTP, cfg.EmailFrom, cfg.EmailTo))
	}
	if cfg.PagerDutyKey != "" {
		alerters = append(alerters, NewPagerDutyAlerter(cfg.PagerDutyKey))
	}
	if cfg.OpsGenieKey != "" {
		alerters = append(alerters, NewOpsGenieAlerter(cfg.OpsGenieKey))
	}

	if len(alerters) == 0 {
		return nil, errors.New("notify: no alerting method configured")
	}
	return NewMultiAlerter(alerters...), nil
}
