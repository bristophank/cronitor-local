package notify

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOpsGenieAlerterSuccess(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.Header.Get("Authorization") != "GenieKey test-key" {
			t.Errorf("missing or wrong Authorization header")
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("missing Content-Type header")
		}
		w.WriteHeader(http.StatusAccepted)
	}))
	defer ts.Close()

	a := NewOpsGenieAlerter("test-key")
	a.apiURL = ts.URL

	if err := a.Alert("my-job"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestOpsGenieAlerterNon2xxReturnsError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	a := NewOpsGenieAlerter("bad-key")
	a.apiURL = ts.URL

	if err := a.Alert("my-job"); err == nil {
		t.Fatal("expected error for non-2xx response")
	}
}

func TestOpsGenieAlerterUnreachableReturnsError(t *testing.T) {
	a := NewOpsGenieAlerter("test-key")
	a.apiURL = "http://127.0.0.1:0"

	if err := a.Alert("my-job"); err == nil {
		t.Fatal("expected error for unreachable server")
	}
}
