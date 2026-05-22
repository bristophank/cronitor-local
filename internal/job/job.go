package job

import (
	"time"
)

// Status represents the current state of a cron job
type Status string

const (
	StatusHealthy  Status = "healthy"
	StatusFailing  Status = "failing"
	StatusUnknown  Status = "unknown"
	StatusRunning  Status = "running"
)

// Job represents a monitored cron job
type Job struct {
	ID            string        `json:"id"`
	Name          string        `json:"name"`
	Schedule      string        `json:"schedule"`
	Status        Status        `json:"status"`
	LastPing      *time.Time    `json:"last_ping,omitempty"`
	LastSuccess   *time.Time    `json:"last_success,omitempty"`
	LastFailure   *time.Time    `json:"last_failure,omitempty"`
	GracePeriod   time.Duration `json:"grace_period"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}

// NewJob creates a new Job with default values
func NewJob(id, name, schedule string) *Job {
	now := time.Now().UTC()
	return &Job{
		ID:          id,
		Name:        name,
		Schedule:    schedule,
		Status:      StatusUnknown,
		GracePeriod: 5 * time.Minute,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// RecordSuccess marks the job as healthy and records the success time
func (j *Job) RecordSuccess() {
	now := time.Now().UTC()
	j.LastPing = &now
	j.LastSuccess = &now
	j.Status = StatusHealthy
	j.UpdatedAt = now
}

// RecordFailure marks the job as failing and records the failure time
func (j *Job) RecordFailure() {
	now := time.Now().UTC()
	j.LastPing = &now
	j.LastFailure = &now
	j.Status = StatusFailing
	j.UpdatedAt = now
}

// IsOverdue returns true if the job has not pinged within its expected window
func (j *Job) IsOverdue(expectedInterval time.Duration) bool {
	if j.LastPing == nil {
		return false
	}
	deadline := j.LastPing.Add(expectedInterval).Add(j.GracePeriod)
	return time.Now().UTC().After(deadline)
}
