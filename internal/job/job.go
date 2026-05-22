package job

import (
	"time"
)

// Status represents the last known state of a job.
type Status string

const (
	StatusUnknown Status = "unknown"
	StatusOK      Status = "ok"
	StatusFailed  Status = "failed"
)

// Job holds the configuration and runtime state of a monitored cron job.
type Job struct {
	Name        string        `json:"name"`
	Schedule    string        `json:"schedule"`
	GracePeriod time.Duration `json:"grace_period"`
	LastRunAt   *time.Time    `json:"last_run_at,omitempty"`
	LastSuccess *time.Time    `json:"last_success,omitempty"`
	LastFailure *time.Time    `json:"last_failure,omitempty"`
	Status      Status        `json:"status"`
	ConsecFails int           `json:"consec_fails"`
}

// NewJob creates a new Job with the given name, cron schedule, and grace period.
func NewJob(name, schedule string, gracePeriod time.Duration) *Job {
	return &Job{
		Name:        name,
		Schedule:    schedule,
		GracePeriod: gracePeriod,
		Status:      StatusUnknown,
	}
}

// RecordSuccess marks the job as having completed successfully.
func (j *Job) RecordSuccess() {
	now := time.Now().UTC()
	j.LastRunAt = &now
	j.LastSuccess = &now
	j.Status = StatusOK
	j.ConsecFails = 0
}

// RecordFailure marks the job as having failed.
func (j *Job) RecordFailure() {
	now := time.Now().UTC()
	j.LastRunAt = &now
	j.LastFailure = &now
	j.Status = StatusFailed
	j.ConsecFails++
}

// IsOverdue returns true when the job has not succeeded within its grace period
// past the expected next run time. If no successful run has been recorded the
// job is never considered overdue.
func (j *Job) IsOverdue(nextRunAt time.Time) bool {
	if j.LastSuccess == nil {
		return false
	}
	deadline := nextRunAt.Add(j.GracePeriod)
	return time.Now().UTC().After(deadline) && j.LastSuccess.Before(nextRunAt)
}
