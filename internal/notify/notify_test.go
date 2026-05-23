package notify_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/example/cronitor-local/internal/notify"
)

// recordingAlerter captures calls to Alert for assertion.
type recordingAlerter struct {
	calls []string
	err   error
}

func (r *recordingAlerter) Alert(jobName, reason string) error {
	r.calls = append(r.calls, jobName+":"+reason)
	return r.err
}

func TestNoopAlerter(t *testing.T) {
	a := &notify.NoopAlerter{}
	if err := a.Alert("myjob", "overdue"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestMultiAlerterFansOut(t *testing.T) {
	a1 := &recordingAlerter{}
	a2 := &recordingAlerter{}
	m := notify.NewMultiAlerter(a1, a2)

	if err := m.Alert("job1", "overdue"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(a1.calls) != 1 || len(a2.calls) != 1 {
		t.Fatalf("expected 1 call each, got %d and %d", len(a1.calls), len(a2.calls))
	}
}

func TestMultiAlerterReturnsFirstError(t *testing.T) {
	expected := errors.New("send failed")
	a1 := &recordingAlerter{err: expected}
	a2 := &recordingAlerter{}
	m := notify.NewMultiAlerter(a1, a2)

	if err := m.Alert("job1", "overdue"); !errors.Is(err, expected) {
		t.Fatalf("expected %v, got %v", expected, err)
	}
	// a2 should still have been called
	if len(a2.calls) != 1 {
		t.Fatalf("expected a2 to be called despite a1 error")
	}
}

func TestSlackAlerterSuccess(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	a := notify.NewSlackAlerter(ts.URL)
	if err := a.Alert("backup", "missed schedule"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSlackAlerterNonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	a := notify.NewSlackAlerter(ts.URL)
	if err := a.Alert("backup", "missed schedule"); err == nil {
		t.Fatal("expected error for non-200 status")
	}
}
