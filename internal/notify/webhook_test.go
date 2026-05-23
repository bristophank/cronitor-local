package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWebhookAlerterSuccess(t *testing.T) {
	var received webhookPayload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("expected Content-Type application/json, got %s", ct)
		}
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	a := NewWebhookAlerter(ts.URL)
	if err := a.Alert("backup-job", "overdue by 5m"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.JobName != "backup-job" {
		t.Errorf("expected job_name backup-job, got %s", received.JobName)
	}
	if received.Message != "overdue by 5m" {
		t.Errorf("expected message 'overdue by 5m', got %s", received.Message)
	}
	if received.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestWebhookAlerterNon2xxReturnsError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	a := NewWebhookAlerter(ts.URL)
	if err := a.Alert("job", "msg"); err == nil {
		t.Fatal("expected error for 500 response")
	}
}

func TestWebhookAlerterUnreachableReturnsError(t *testing.T) {
	a := NewWebhookAlerter("http://127.0.0.1:19999/no-server")
	if err := a.Alert("job", "msg"); err == nil {
		t.Fatal("expected error for unreachable server")
	}
}
