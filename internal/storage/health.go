package storage

import (
	"example.com/orbit-notebook/internal/integrity"
	"fmt"
	"os"
)

type Health struct {
	Exists  bool  `json:"exists"`
	Bytes   int64 `json:"bytes"`
	Records int   `json:"records"`
}

func (s *Store) Health() (Health, error) {
	size, err := s.Size()
	if err != nil {
		return Health{}, err
	}
	items, err := s.Load()
	if err != nil {
		return Health{}, err
	}
	if err = integrity.Collection(items); err != nil {
		return Health{}, err
	}
	return Health{Exists: s.Exists(), Bytes: size, Records: len(items)}, nil
}
func (s *Store) Remove() error {
	if err := os.Remove(s.path()); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("remove store: %w", err)
	}
	return nil
}
