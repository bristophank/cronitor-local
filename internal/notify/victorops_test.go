package notify

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestVictorOpsAlerterSuccess(t *testing.T) {
	var received victorOpsPayload

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	alerter := NewVictorOpsAlerter(server.URL)
	err := alerter.Alert("backup-job", "job has not run in 2 hours")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if received.MessageType != "CRITICAL" {
		t.Errorf("expected message_type CRITICAL, got %s", received.MessageType)
	}
	if received.EntityID != "backup-job" {
		t.Errorf("expected entity_id backup-job, got %s", received.EntityID)
	}
	if received.StateMessage != "job has not run in 2 hours" {
		t.Errorf("unexpected state_message: %s", received.StateMessage)
	}
	if received.Timestamp == 0 {
		t.Error("expected non-zero timestamp")
	}
}

func TestVictorOpsAlerterNon2xxReturnsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	alerter := NewVictorOpsAlerter(server.URL)
	err := alerter.Alert("backup-job", "overdue")
	if err == nil {
		t.Fatal("expected error for non-2xx status, got nil")
	}
}

func TestVictorOpsAlerterUnreachableReturnsError(t *testing.T) {
	alerter := NewVictorOpsAlerter("http://127.0.0.1:0/nonexistent")
	err := alerter.Alert("backup-job", "overdue")
	if err == nil {
		t.Fatal("expected error for unreachable server, got nil")
	}
}
