package store_test

import (
	"os"
	"testing"
	"time"

	"github.com/cronitor-local/internal/job"
	"github.com/cronitor-local/internal/store"
)

func tempFile(t *testing.T) string {
	t.Helper()
	f, err := os.CreateTemp("", "store-test-*.json")
	if err != nil {
		t.Fatalf("creating temp file: %v", err)
	}
	f.Close()
	os.Remove(f.Name()) // let store create it fresh
	return f.Name()
}

func TestSaveAndGet(t *testing.T) {
	s, err := store.New(tempFile(t))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	j := job.NewJob("backup", "0 2 * * *", 5*time.Minute)
	if err := s.Save(j); err != nil {
		t.Fatalf("Save: %v", err)
	}
	got, ok := s.Get("backup")
	if !ok {
		t.Fatal("expected job to be found")
	}
	if got.Name != "backup" {
		t.Errorf("name: got %q, want %q", got.Name, "backup")
	}
}

func TestAll(t *testing.T) {
	s, _ := store.New(tempFile(t))
	s.Save(job.NewJob("job1", "* * * * *", time.Minute))
	s.Save(job.NewJob("job2", "* * * * *", time.Minute))
	if len(s.All()) != 2 {
		t.Errorf("All: expected 2 jobs, got %d", len(s.All()))
	}
}

func TestDelete(t *testing.T) {
	s, _ := store.New(tempFile(t))
	s.Save(job.NewJob("temp", "* * * * *", time.Minute))
	if err := s.Delete("temp"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if _, ok := s.Get("temp"); ok {
		t.Error("expected job to be deleted")
	}
}

func TestPersistence(t *testing.T) {
	path := tempFile(t)
	s1, _ := store.New(path)
	s1.Save(job.NewJob("persist", "@hourly", 10*time.Minute))

	s2, err := store.New(path)
	if err != nil {
		t.Fatalf("reloading store: %v", err)
	}
	if _, ok := s2.Get("persist"); !ok {
		t.Error("expected persisted job to be reloaded")
	}
}
