package notify

import (
	"testing"
)

func TestBuildAlerterNoConfigReturnsError(t *testing.T) {
	_, err := BuildAlerter(Config{})
	if err == nil {
		t.Fatal("expected error when no alerter configured")
	}
}

func TestBuildAlerterWithWebhook(t *testing.T) {
	a, err := BuildAlerter(Config{WebhookURL: "http://example.com/hook"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a == nil {
		t.Fatal("expected non-nil alerter")
	}
}

func TestBuildAlerterWithSlack(t *testing.T) {
	a, err := BuildAlerter(Config{SlackURL: "http://hooks.slack.com/test"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a == nil {
		t.Fatal("expected non-nil alerter")
	}
}

func TestBuildAlerterWithEmail(t *testing.T) {
	cfg := Config{EmailSMTP: "smtp://localhost:25", EmailFrom: "a@b.com", EmailTo: "c@d.com"}
	a, err := BuildAlerter(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a == nil {
		t.Fatal("expected non-nil alerter")
	}
}

func TestBuildAlerterWithGotify(t *testing.T) {
	cfg := Config{GotifyURL: "http://gotify.local", GotifyToken: "abc123", GotifyPriority: 8}
	a, err := BuildAlerter(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a == nil {
		t.Fatal("expected non-nil alerter")
	}
}

func TestBuildAlerterGotifyMissingTokenIgnored(t *testing.T) {
	// Gotify requires both URL and token; only URL should not register it.
	_, err := BuildAlerter(Config{GotifyURL: "http://gotify.local"})
	if err == nil {
		t.Fatal("expected error when only GotifyURL is set without token")
	}
}

func TestLoadConfigReadsEnv(t *testing.T) {
	t.Setenv("GOTIFY_URL", "http://gotify.example.com")
	t.Setenv("GOTIFY_TOKEN", "mytoken")
	t.Setenv("GOTIFY_PRIORITY", "9")

	cfg := LoadConfig()
	if cfg.GotifyURL != "http://gotify.example.com" {
		t.Errorf("unexpected GotifyURL: %s", cfg.GotifyURL)
	}
	if cfg.GotifyToken != "mytoken" {
		t.Errorf("unexpected GotifyToken: %s", cfg.GotifyToken)
	}
	if cfg.GotifyPriority != 9 {
		t.Errorf("unexpected GotifyPriority: %d", cfg.GotifyPriority)
	}
}
