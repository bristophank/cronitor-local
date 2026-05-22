package web

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/cronitor-local/internal/store"
)

// Handler holds dependencies for HTTP handlers.
type Handler struct {
	store *store.Store
}

// NewHandler creates a new Handler with the given store.
func NewHandler(s *store.Store) *Handler {
	return &Handler{store: s}
}

// RegisterRoutes wires up all HTTP routes on the given mux.
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", h.handleIndex)
	mux.HandleFunc("/api/jobs", h.handleJobs)
	mux.HandleFunc("/healthz", h.handleHealth)
}

// handleIndex serves a minimal HTML dashboard.
func (h *Handler) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	jobs, err := h.store.All()
	if err != nil {
		http.Error(w, "failed to load jobs", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_ = dashboardTmpl.Execute(w, map[string]interface{}{
		"Jobs":      jobs,
		"Timestamp": time.Now().Format(time.RFC1123),
	})
}

// handleJobs returns all jobs as JSON.
func (h *Handler) handleJobs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	jobs, err := h.store.All()
	if err != nil {
		http.Error(w, "failed to load jobs", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(jobs)
}

// handleHealth returns a simple 200 OK for liveness checks.
func (h *Handler) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
