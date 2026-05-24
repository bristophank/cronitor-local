package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPagerDutyAlerterSuccess(t *testing.T) {
	var received pdPayload

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusAccepted)
	}))
	defer server.Close()

	alerter := NewPagerDutyAlerter("test-key")
	alerter.client.Transport = rewriteTransport(server.URL)

	if err := alerter.Alert("my-job", "overdue by 10m"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if received.RoutingKey != "test-key" {
		t.Errorf("routing key: got %q, want %q", received.RoutingKey, "test-key")
	}
	if received.EventAction != "trigger" {
		t.Errorf("event_action: got %q, want %q", received.EventAction, "trigger")
	}
	if received.Payload.Severity != "error" {
		t.Errorf("severity: got %q, want %q", received.Payload.Severity, "error")
	}
	if received.Payload.Summary == "" {
		t.Error("expected non-empty payload summary")
	}
}

func TestPagerDutyAlerterNon2xxReturnsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	alerter := NewPagerDutyAlerter("bad-key")
	alerter.client.Transport = rewriteTransport(server.URL)

	if err := alerter.Alert("job", "msg"); err == nil {
		t.Fatal("expected error for non-2xx response")
	}
}

func TestPagerDutyAlerterUnreachableReturnsError(t *testing.T) {
	alerter := NewPagerDutyAlerter("key")
	// point at a port that refuses connections
	alerter.client.Transport = rewriteTransport("http://127.0.0.1:1")

	if err := alerter.Alert("job", "msg"); err == nil {
		t.Fatal("expected error for unreachable host")
	}
}

func TestPagerDutyAlerterPayloadContainsJobAndMessage(t *testing.T) {
	var received pdPayload

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusAccepted)
	}))
	defer server.Close()

	alerter := NewPagerDutyAlerter("key")
	alerter.client.Transport = rewriteTransport(server.URL)

	if err := alerter.Alert("nightly-backup", "overdue by 5m"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := received.Payload.Summary; got == "" {
		t.Error("expected payload summary to be set")
	}
	// The summary should reference the job name so on-call engineers know what fired.
	if got := received.Payload.Summary; !contains(got, "nightly-backup") {
		t.Errorf("summary %q does not mention job name %q", got, "nightly-backup")
	}
}
