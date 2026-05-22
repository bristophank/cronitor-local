package monitor

import (
	"log"
	"time"

	"cronitor-local/internal/job"
	"cronitor-local/internal/store"
)

// Monitor periodically checks all jobs for overdue status and triggers alerts.
type Monitor struct {
	store    *store.Store
	alerter  Alerter
	interval time.Duration
	stop     chan struct{}
}

// Alerter is the interface for sending alerts.
type Alerter interface {
	Alert(j *job.Job) error
}

// New creates a new Monitor with the given store, alerter, and check interval.
func New(s *store.Store, a Alerter, interval time.Duration) *Monitor {
	return &Monitor{
		store:    s,
		alerter:  a,
		interval: interval,
		stop:     make(chan struct{}),
	}
}

// Start begins the monitoring loop in a background goroutine.
func (m *Monitor) Start() {
	go func() {
		ticker := time.NewTicker(m.interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				m.checkAll()
			case <-m.stop:
				return
			}
		}
	}()
}

// Stop halts the monitoring loop.
func (m *Monitor) Stop() {
	close(m.stop)
}

// checkAll iterates all jobs and alerts on any that are overdue.
func (m *Monitor) checkAll() {
	jobs, err := m.store.All()
	if err != nil {
		log.Printf("monitor: failed to load jobs: %v", err)
		return
	}
	for _, j := range jobs {
		if j.IsOverdue() {
			log.Printf("monitor: job %q is overdue, sending alert", j.Name)
			if err := m.alerter.Alert(j); err != nil {
				log.Printf("monitor: alert failed for job %q: %v", j.Name, err)
			}
		}
	}
}
