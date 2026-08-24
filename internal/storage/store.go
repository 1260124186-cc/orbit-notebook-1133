package storage

import (
	"example.com/orbit-notebook/internal/codec"
	"example.com/orbit-notebook/internal/model"
	"os"
	"path/filepath"
	"sync"
)

type Store struct {
	dir string
	mu  sync.Mutex
}

func New(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	return &Store{dir: dir}, nil
}
func (s *Store) path() string { return filepath.Join(s.dir, "observations.json") }
func (s *Store) Load() ([]model.Observation, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.loadUnlocked()
}
func (s *Store) loadUnlocked() ([]model.Observation, error) {
	var items []model.Observation
	err := codec.ReadJSON(s.path(), &items)
	if os.IsNotExist(err) {
		return []model.Observation{}, nil
	}
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (s *Store) Replace(fn func([]model.Observation) ([]model.Observation, error)) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	current, err := s.loadUnlocked()
	if err != nil {
		return err
	}
	next, err := fn(current)
	if err != nil {
		return err
	}
	return codec.AtomicWrite(s.path(), next)
}
