package notify

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTeamsAlerterSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("expected Content-Type application/json, got %s", ct)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	alerter := NewTeamsAlerter(server.URL)
	if err := alerter.Alert("my-job", "job is overdue"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestTeamsAlerterNon2xxReturnsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer server.Close()

	alerter := NewTeamsAlerter(server.URL)
	err := alerter.Alert("my-job", "job is overdue")
	if err == nil {
		t.Fatal("expected error for non-2xx status, got nil")
	}
}

func TestTeamsAlerterUnreachableReturnsError(t *testing.T) {
	alerter := NewTeamsAlerter("http://127.0.0.1:0/webhook")
	err := alerter.Alert("my-job", "job is overdue")
	if err == nil {
		t.Fatal("expected error for unreachable server, got nil")
	}
}

func TestTeamsAlerterPayloadFields(t *testing.T) {
	var capturedBody []byte
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf := make([]byte, r.ContentLength)
		r.Body.Read(buf)
		capturedBody = buf
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	alerter := NewTeamsAlerter(server.URL)
	if err := alerter.Alert("backup-job", "missed schedule"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	body := string(capturedBody)
	if len(body) == 0 {
		t.Fatal("expected non-empty request body")
	}
	for _, want := range []string{"MessageCard", "backup-job", "missed schedule", "FF0000"} {
		if !containsString(body, want) {
			t.Errorf("expected body to contain %q, got: %s", want, body)
		}
	}
}

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && stringContains(s, substr))
}

func stringContains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
