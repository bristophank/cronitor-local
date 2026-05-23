package notify

import (
	"testing"
)

func TestBuildAlerterNoConfigReturnsError(t *testing.T) {
	_, err := BuildAlerter(Config{})
	if err == nil {
		t.Fatal("expected error when no alerters configured, got nil")
	}
}

func TestBuildAlerterWithWebhook(t *testing.T) {
	cfg := Config{WebhookURL: "https://example.com/hook"}
	a, err := BuildAlerter(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a == nil {
		t.Fatal("expected non-nil alerter")
	}
}

func TestBuildAlerterWithSlack(t *testing.T) {
	cfg := Config{SlackWebhookURL: "https://hooks.slack.com/services/TEST"}
	a, err := BuildAlerter(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a == nil {
		t.Fatal("expected non-nil alerter")
	}
}

func TestBuildAlerterWithEmail(t *testing.T) {
	cfg := Config{
		SMTPHost:  "smtp.example.com",
		SMTPPort:  "587",
		EmailFrom: "alert@example.com",
		EmailTo:   []string{"ops@example.com"},
	}
	a, err := BuildAlerter(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a == nil {
		t.Fatal("expected non-nil alerter")
	}
}

func TestLoadConfigReadsEnv(t *testing.T) {
	t.Setenv("ALERT_WEBHOOK_URL", "https://example.com/hook")
	t.Setenv("ALERT_EMAIL_TO", "a@example.com, b@example.com")

	cfg := LoadConfig()
	if cfg.WebhookURL != "https://example.com/hook" {
		t.Errorf("expected webhook URL, got %q", cfg.WebhookURL)
	}
	if len(cfg.EmailTo) != 2 {
		t.Errorf("expected 2 email recipients, got %d", len(cfg.EmailTo))
	}
}

func TestBuildAlerterEmailDefaultsPort(t *testing.T) {
	cfg := Config{
		SMTPHost:  "smtp.example.com",
		EmailFrom: "alert@example.com",
		EmailTo:   []string{"ops@example.com"},
	}
	// Should not error — port defaults to 587
	_, err := BuildAlerter(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
