package monitor

import (
	"fmt"
	"log"

	"cronitor-local/internal/job"
)

// LogAlerter is a simple Alerter that writes alerts to the standard logger.
type LogAlerter struct{}

// Alert logs an overdue alert for the given job.
func (l *LogAlerter) Alert(j *job.Job) error {
	log.Printf("ALERT: job %q has not run since %v (interval: %v)",
		j.Name, j.LastRunAt, j.Interval)
	return nil
}

// MultiAlerter fans out alerts to multiple Alerter implementations.
type MultiAlerter struct {
	alerters []Alerter
}

// NewMultiAlerter creates a MultiAlerter wrapping the provided alerters.
func NewMultiAlerter(alerters ...Alerter) *MultiAlerter {
	return &MultiAlerter{alerters: alerters}
}

// Alert sends the alert to all registered alerters, collecting any errors.
func (m *MultiAlerter) Alert(j *job.Job) error {
	var errs []string
	for _, a := range m.alerters {
		if err := a.Alert(j); err != nil {
			errs = append(errs, err.Error())
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("multi-alerter errors: %v", errs)
	}
	return nil
}
