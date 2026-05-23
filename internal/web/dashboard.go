package web

import (
	"net/http"
	"time"

	"github.com/user/cronitor-local/internal/job"
	"github.com/user/cronitor-local/internal/store"
)

type dashboardData struct {
	Jobs      []*job.Job
	CheckedAt time.Time
}

func (h *Handler) handleDashboard(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	jobs, err := h.store.All()
	if err != nil {
		http.Error(w, "failed to load jobs", http.StatusInternalServerError)
		return
	}

	data := dashboardData{
		Jobs:      jobs,
		CheckedAt: time.Now().UTC(),
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := dashboardTmpl.Execute(w, data); err != nil {
		http.Error(w, "template error: "+err.Error(), http.StatusInternalServerError)
	}
}

// Ensure Handler has a store field — this file relies on the existing Handler type.
// The store field is expected to be of type *store.Store as defined in handler.go.
var _ *store.Store // import guard
