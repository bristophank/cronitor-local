package notify

import (
	"errors"
	"os"
	"strconv"
)

// Config holds optional alerter configuration loaded from environment variables.
type Config struct {
	WebhookURL     string
	SlackURL       string
	EmailSMTP      string
	EmailFrom      string
	EmailTo        string
	PagerDutyKey   string
	OpsGenieKey    string
	VictorOpsURL   string
	TeamsURL       string
	DiscordURL     string
	GotifyURL      string
	GotifyToken    string
	GotifyPriority int
}

// LoadConfig reads alerter settings from environment variables.
func LoadConfig() Config {
	priority, _ := strconv.Atoi(os.Getenv("GOTIFY_PRIORITY"))
	return Config{
		WebhookURL:     os.Getenv("WEBHOOK_URL"),
		SlackURL:       os.Getenv("SLACK_WEBHOOK_URL"),
		EmailSMTP:      os.Getenv("EMAIL_SMTP"),
		EmailFrom:      os.Getenv("EMAIL_FROM"),
		EmailTo:        os.Getenv("EMAIL_TO"),
		PagerDutyKey:   os.Getenv("PAGERDUTY_INTEGRATION_KEY"),
		OpsGenieKey:    os.Getenv("OPSGENIE_API_KEY"),
		VictorOpsURL:   os.Getenv("VICTOROPS_WEBHOOK_URL"),
		TeamsURL:       os.Getenv("TEAMS_WEBHOOK_URL"),
		DiscordURL:     os.Getenv("DISCORD_WEBHOOK_URL"),
		GotifyURL:      os.Getenv("GOTIFY_URL"),
		GotifyToken:    os.Getenv("GOTIFY_TOKEN"),
		GotifyPriority: priority,
	}
}

// BuildAlerter constructs a MultiAlerter from the provided Config.
// Returns an error if no alerter is configured.
func BuildAlerter(cfg Config) (Alerter, error) {
	var alerters []Alerter

	if cfg.WebhookURL != "" {
		alerters = append(alerters, NewWebhookAlerter(cfg.WebhookURL))
	}
	if cfg.SlackURL != "" {
		alerters = append(alerters, NewSlackAlerter(cfg.SlackURL))
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
	if cfg.VictorOpsURL != "" {
		alerters = append(alerters, NewVictorOpsAlerter(cfg.VictorOpsURL))
	}
	if cfg.TeamsURL != "" {
		alerters = append(alerters, NewTeamsAlerter(cfg.TeamsURL))
	}
	if cfg.DiscordURL != "" {
		alerters = append(alerters, NewDiscordAlerter(cfg.DiscordURL))
	}
	if cfg.GotifyURL != "" && cfg.GotifyToken != "" {
		alerters = append(alerters, NewGotifyAlerter(cfg.GotifyURL, cfg.GotifyToken, cfg.GotifyPriority))
	}

	if len(alerters) == 0 {
		return nil, errors.New("no alerter configured: set at least one of WEBHOOK_URL, SLACK_WEBHOOK_URL, EMAIL_SMTP, PAGERDUTY_INTEGRATION_KEY, OPSGENIE_API_KEY, VICTOROPS_WEBHOOK_URL, TEAMS_WEBHOOK_URL, DISCORD_WEBHOOK_URL, GOTIFY_URL+GOTIFY_TOKEN")
	}
	return NewMultiAlerter(alerters...), nil
}
