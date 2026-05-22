package job

import (
	"testing"
	"time"
)

func TestNewJob(t *testing.T) {
	j := NewJob("abc123", "backup-db", "0 2 * * *")

	if j.ID != "abc123" {
		t.Errorf("expected ID abc123, got %s", j.ID)
	}
	if j.Name != "backup-db" {
		t.Errorf("expected Name backup-db, got %s", j.Name)
	}
	if j.Status != StatusUnknown {
		t.Errorf("expected status unknown, got %s", j.Status)
	}
	if j.GracePeriod != 5*time.Minute {
		t.Errorf("expected grace period 5m, got %v", j.GracePeriod)
	}
	if j.LastPing != nil {
		t.Error("expected LastPing to be nil for new job")
	}
}

func TestRecordSuccess(t *testing.T) {
	j := NewJob("abc123", "backup-db", "0 2 * * *")
	j.RecordSuccess()

	if j.Status != StatusHealthy {
		t.Errorf("expected status healthy, got %s", j.Status)
	}
	if j.LastPing == nil {
		t.Error("expected LastPing to be set after success")
	}
	if j.LastSuccess == nil {
		t.Error("expected LastSuccess to be set after success")
	}
	if j.LastFailure != nil {
		t.Error("expected LastFailure to remain nil after success")
	}
}

func TestRecordFailure(t *testing.T) {
	j := NewJob("abc123", "backup-db", "0 2 * * *")
	j.RecordFailure()

	if j.Status != StatusFailing {
		t.Errorf("expected status failing, got %s", j.Status)
	}
	if j.LastFailure == nil {
		t.Error("expected LastFailure to be set after failure")
	}
}

func TestIsOverdue(t *testing.T) {
	j := NewJob("abc123", "backup-db", "0 2 * * *")

	// No ping yet — should not be overdue
	if j.IsOverdue(time.Hour) {
		t.Error("expected job with no ping to not be overdue")
	}

	// Recent ping — should not be overdue
	j.RecordSuccess()
	if j.IsOverdue(time.Hour) {
		t.Error("expected recently pinged job to not be overdue")
	}

	// Simulate old ping by backdating LastPing
	past := time.Now().UTC().Add(-2 * time.Hour)
	j.LastPing = &past
	j.GracePeriod = 0
	if !j.IsOverdue(time.Hour) {
		t.Error("expected old job to be overdue")
	}
}
