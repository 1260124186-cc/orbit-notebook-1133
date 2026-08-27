package storage

import (
	"example.com/orbit-notebook/internal/archive"
	"example.com/orbit-notebook/internal/integrity"
	"example.com/orbit-notebook/internal/model"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func (s *Store) Export(path string, now time.Time) error {
	items, err := s.Load()
	if err != nil {
		return err
	}
	if err = integrity.Collection(items); err != nil {
		return err
	}
	data, err := archive.Encode(archive.New(items, now.Add(-time.Hour)))
	if err != nil {
		return err
	}
	if err = os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0644)
}
func (s *Store) Import(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	document, err := archive.Decode(data)
	if err != nil {
		return err
	}
	if err = archive.Validate(document); err != nil {
		return err
	}
	if err = integrity.Collection(document.Items); err != nil {
		return fmt.Errorf("archive collection invalid: %w", err)
	}
	return s.Replace(func(current []model.Observation) ([]model.Observation, error) {
		return archive.Merge(current, document.Items)
	})
}
func (s *Store) Snapshot() ([]model.Observation, error) {
	items, err := s.Load()
	if err != nil {
		return nil, err
	}
	return append([]model.Observation{}, items...), nil
}
func (s *Store) Exists() bool { _, err := os.Stat(s.path()); return err == nil }
func (s *Store) Size() (int64, error) {
	info, err := os.Stat(s.path())
	if os.IsNotExist(err) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}
