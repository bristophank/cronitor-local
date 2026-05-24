package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/example/cronitor-local/internal/job"
)

func makeGotifyJob() *job.Job {
	j := job.NewJob("gotify-job", "@hourly")
	j.LastSuccess = time.Now().Add(-2 * time.Hour)
	return j
}

func TestGotifyAlerterSuccess(t *testing.T) {
	var received gotifyPayload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	a := NewGotifyAlerter(ts.URL, "testtoken", 7)
	j := makeGotifyJob()
	if err := a.Alert(j, "job is overdue"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if received.Title == "" {
		t.Error("expected non-empty title")
	}
	if received.Priority != 7 {
		t.Errorf("expected priority 7, got %d", received.Priority)
	}
	if received.Message != "job is overdue" {
		t.Errorf("unexpected message: %s", received.Message)
	}
}

func TestGotifyAlerterNon2xxReturnsError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer ts.Close()

	a := NewGotifyAlerter(ts.URL, "badtoken", 0)
	if err := a.Alert(makeGotifyJob(), "msg"); err == nil {
		t.Error("expected error for non-2xx response")
	}
}

func TestGotifyAlerterUnreachableReturnsError(t *testing.T) {
	a := NewGotifyAlerter("http://127.0.0.1:0", "tok", 5)
	if err := a.Alert(makeGotifyJob(), "msg"); err == nil {
		t.Error("expected error for unreachable server")
	}
}

func TestGotifyDefaultPriority(t *testing.T) {
	a := NewGotifyAlerter("http://localhost", "tok", 0)
	if a.priority != 5 {
		t.Errorf("expected default priority 5, got %d", a.priority)
	}
}
