package store

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/cronitor-local/internal/job"
)

// Store manages persistence of job state to disk.
type Store struct {
	mu       sync.RWMutex
	filePath string
	jobs     map[string]*job.Job
}

// New creates a new Store backed by the given file path.
func New(filePath string) (*Store, error) {
	s := &Store{
		filePath: filePath,
		jobs:     make(map[string]*job.Job),
	}
	if err := s.load(); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("store: loading data: %w", err)
	}
	return s, nil
}

// Save persists a job to the store.
func (s *Store) Save(j *job.Job) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.jobs[j.Name] = j
	return s.flush()
}

// Get retrieves a job by name. Returns nil if not found.
func (s *Store) Get(name string) (*job.Job, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	j, ok := s.jobs[name]
	return j, ok
}

// All returns a copy of all stored jobs.
func (s *Store) All() []*job.Job {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*job.Job, 0, len(s.jobs))
	for _, j := range s.jobs {
		list = append(list, j)
	}
	return list
}

// Delete removes a job from the store by name.
func (s *Store) Delete(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.jobs, name)
	return s.flush()
}

func (s *Store) load() error {
	f, err := os.Open(s.filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(&s.jobs)
}

func (s *Store) flush() error {
	f, err := os.Create(s.filePath)
	if err != nil {
		return fmt.Errorf("store: creating file: %w", err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(s.jobs)
}
