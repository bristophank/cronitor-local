package notify

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestBuildAlerterNoConfigReturnsError(t *testing.T) {
	cfg := Config{}
	_, err := BuildAlerter(cfg)
	if err == nil {
		t.Fatal("expected error when no destinations configured")
	}
}

func TestBuildAlerterWithWebhook(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	cfg := Config{WebhookURL: ts.URL}
	a, err := BuildAlerter(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := a.Alert("test-job", "overdue"); err != nil {
		t.Fatalf("alert failed: %v", err)
	}
}

func TestBuildAlerterWithSlack(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	cfg := Config{SlackWebhook: ts.URL}
	a, err := BuildAlerter(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a == nil {
		t.Fatal("expected non-nil alerter")
	}
}

func TestLoadConfigReadsEnv(t *testing.T) {
	os.Setenv("ALERT_WEBHOOK_URL", "http://example.com/hook")
	os.Setenv("SLACK_WEBHOOK_URL", "http://hooks.slack.com/test")
	defer os.Unsetenv("ALERT_WEBHOOK_URL")
	defer os.Unsetenv("SLACK_WEBHOOK_URL")

	cfg := LoadConfig()
	if cfg.WebhookURL != "http://example.com/hook" {
		t.Errorf("expected webhook URL, got %q", cfg.WebhookURL)
	}
	if cfg.SlackWebhook != "http://hooks.slack.com/test" {
		t.Errorf("expected slack webhook, got %q", cfg.SlackWebhook)
	}
}
