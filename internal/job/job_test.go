package job_test

import (
	"testing"
	"time"

	"github.com/cronitor-local/internal/job"
)

func TestNewJob(t *testing.T) {
	j := job.NewJob("test", "* * * * *", 5*time.Minute)
	if j.Name != "test" {
		t.Errorf("Name: got %q, want %q", j.Name, "test")
	}
	if j.Schedule != "* * * * *" {
		t.Errorf("Schedule: got %q, want %q", j.Schedule, "* * * * *")
	}
	if j.Status != job.StatusUnknown {
		t.Errorf("Status: got %q, want %q", j.Status, job.StatusUnknown)
	}
}

func TestRecordSuccess(t *testing.T) {
	j := job.NewJob("test", "* * * * *", time.Minute)
	j.RecordSuccess()
	if j.Status != job.StatusOK {
		t.Errorf("Status: got %q, want %q", j.Status, job.StatusOK)
	}
	if j.LastSuccess == nil {
		t.Error("LastSuccess should not be nil")
	}
	if j.ConsecFails != 0 {
		t.Errorf("ConsecFails: got %d, want 0", j.ConsecFails)
	}
}

func TestRecordFailure(t *testing.T) {
	j := job.NewJob("test", "* * * * *", time.Minute)
	j.RecordFailure()
	j.RecordFailure()
	if j.Status != job.StatusFailed {
		t.Errorf("Status: got %q, want %q", j.Status, job.StatusFailed)
	}
	if j.ConsecFails != 2 {
		t.Errorf("ConsecFails: got %d, want 2", j.ConsecFails)
	}
	j.RecordSuccess()
	if j.ConsecFails != 0 {
		t.Errorf("ConsecFails after success: got %d, want 0", j.ConsecFails)
	}
}

func TestIsOverdue(t *testing.T) {
	j := job.NewJob("test", "* * * * *", 5*time.Minute)

	// No last success — never overdue.
	if j.IsOverdue(time.Now().Add(-10 * time.Minute)) {
		t.Error("job with no success record should not be overdue")
	}

	// Successful run before the scheduled time — should be overdue.
	j.RecordSuccess()
	past := time.Now().Add(-20 * time.Minute)
	if !j.IsOverdue(past) {
		t.Error("job should be overdue when last success is before scheduled run + grace period")
	}

	// Scheduled run in the future — not yet overdue.
	future := time.Now().Add(10 * time.Minute)
	if j.IsOverdue(future) {
		t.Error("job should not be overdue when next run is in the future")
	}
}
