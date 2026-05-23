package web

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/user/cronitor-local/internal/job"
)

func TestHandleDashboardEmpty(t *testing.T) {
	s := tempStore(t)
	h := NewHandler(s)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	body := rr.Body.String()
	if !strings.Contains(body, "Cronitor Local") {
		t.Error("expected dashboard title in response")
	}
	if !strings.Contains(body, "No jobs registered") {
		t.Error("expected empty-state message")
	}
}

func TestHandleDashboardWithJobs(t *testing.T) {
	s := tempStore(t)

	j := job.NewJob("backup", "0 2 * * *", 25*time.Hour)
	j.RecordSuccess()
	if err := s.Save(j); err != nil {
		t.Fatal(err)
	}

	h := NewHandler(s)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	body := rr.Body.String()
	if !strings.Contains(body, "backup") {
		t.Error("expected job name in dashboard")
	}
	if !strings.Contains(body, "0 2 * * *") {
		t.Error("expected schedule in dashboard")
	}
	if !strings.Contains(body, "OK") {
		t.Error("expected OK status for fresh job")
	}
}

func TestHandleDashboardOverdueJob(t *testing.T) {
	s := tempStore(t)

	j := job.NewJob("stale", "* * * * *", time.Minute)
	// last success was 10 minutes ago — overdue
	j.LastSuccess = time.Now().Add(-10 * time.Minute)
	if err := s.Save(j); err != nil {
		t.Fatal(err)
	}

	h := NewHandler(s)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	body := rr.Body.String()
	if !strings.Contains(body, "OVERDUE") {
		t.Error("expected OVERDUE status for late job")
	}
}

func TestHandleDashboardNotFoundOnSubpath(t *testing.T) {
	s := tempStore(t)
	h := NewHandler(s)

	req := httptest.NewRequest(http.MethodGet, "/unknown-path", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rr.Code)
	}
}
