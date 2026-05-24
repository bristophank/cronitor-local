package notify

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/user/cronitor-local/internal/job"
)

func TestDiscordAlerterSuccess(t *testing.T) {
	var received discordPayload

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	alerter := NewDiscordAlerter(ts.URL)
	j := &job.Job{Name: "backup", Schedule: "@daily", LastSuccess: time.Now().Add(-25 * time.Hour)}

	err := alerter.Alert(j, "job is overdue")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(received.Embeds) == 0 {
		t.Fatal("expected at least one embed")
	}
	if !strings.Contains(received.Embeds[0].Title, "backup") {
		t.Errorf("expected title to contain job name, got %q", received.Embeds[0].Title)
	}
	if received.Embeds[0].Description != "job is overdue" {
		t.Errorf("expected description %q, got %q", "job is overdue", received.Embeds[0].Description)
	}
}

func TestDiscordAlerterNon2xxReturnsError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	alerter := NewDiscordAlerter(ts.URL)
	j := &job.Job{Name: "sync", Schedule: "@hourly"}

	err := alerter.Alert(j, "overdue")
	if err == nil {
		t.Fatal("expected error for non-2xx response")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Errorf("expected error to mention status code, got %v", err)
	}
}

func TestDiscordAlerterUnreachableReturnsError(t *testing.T) {
	alerter := NewDiscordAlerter("http://127.0.0.1:1")
	j := &job.Job{Name: "cleanup", Schedule: "@weekly"}

	err := alerter.Alert(j, "overdue")
	if err == nil {
		t.Fatal("expected error for unreachable server")
	}
}
