package monitor

import (
	"sync/atomic"
	"testing"
	"time"

	"cronitor-local/internal/job"
	"cronitor-local/internal/store"
)

// countingAlerter counts how many alerts were fired.
type countingAlerter struct {
	count int64
}

func (c *countingAlerter) Alert(j *job.Job) error {
	atomic.AddInt64(&c.count, 1)
	return nil
}

func tempStore(t *testing.T) *store.Store {
	t.Helper()
	f := t.TempDir() + "/jobs.json"
	s, err := store.New(f)
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	return s
}

func TestMonitorDetectsOverdueJob(t *testing.T) {
	s := tempStore(t)
	j := job.NewJob("backup", 1*time.Millisecond)
	j.LastRunAt = time.Now().Add(-1 * time.Hour) // force overdue
	if err := s.Save(j); err != nil {
		t.Fatalf("save: %v", err)
	}

	alerter := &countingAlerter{}
	m := New(s, alerter, 20*time.Millisecond)
	m.Start()
	time.Sleep(60 * time.Millisecond)
	m.Stop()

	if atomic.LoadInt64(&alerter.count) == 0 {
		t.Error("expected at least one alert for overdue job")
	}
}

func TestMonitorNoAlertForFreshJob(t *testing.T) {
	s := tempStore(t)
	j := job.NewJob("heartbeat", 1*time.Hour)
	j.RecordSuccess()
	if err := s.Save(j); err != nil {
		t.Fatalf("save: %v", err)
	}

	alerter := &countingAlerter{}
	m := New(s, alerter, 20*time.Millisecond)
	m.Start()
	time.Sleep(60 * time.Millisecond)
	m.Stop()

	if atomic.LoadInt64(&alerter.count) != 0 {
		t.Errorf("expected no alerts for fresh job, got %d", alerter.count)
	}
}

func TestMultiAlerter(t *testing.T) {
	a1 := &countingAlerter{}
	a2 := &countingAlerter{}
	ma := NewMultiAlerter(a1, a2)
	j := job.NewJob("test", time.Minute)
	if err := ma.Alert(j); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a1.count != 1 || a2.count != 1 {
		t.Errorf("expected each alerter called once, got %d and %d", a1.count, a2.count)
	}
}
