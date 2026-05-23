package notify

// Alerter is the interface implemented by all notification backends.
type Alerter interface {
	Alert(jobName, reason string) error
}

// NoopAlerter is an Alerter that does nothing. Useful for testing.
type NoopAlerter struct{}

func (n *NoopAlerter) Alert(_, _ string) error { return nil }

// MultiAlerter fans out alerts to multiple Alerter implementations.
type MultiAlerter struct {
	alerters []Alerter
}

// NewMultiAlerter creates a MultiAlerter that notifies all provided alerters.
func NewMultiAlerter(alerters ...Alerter) *MultiAlerter {
	return &MultiAlerter{alerters: alerters}
}

// Alert sends the alert to every registered alerter, collecting any errors.
func (m *MultiAlerter) Alert(jobName, reason string) error {
	var first error
	for _, a := range m.alerters {
		if err := a.Alert(jobName, reason); err != nil && first == nil {
			first = err
		}
	}
	return first
}
