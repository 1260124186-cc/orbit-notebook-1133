package storage

import (
	"example.com/orbit-notebook/internal/archive"
	"example.com/orbit-notebook/internal/codec"
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
	data, err := archive.Encode(archive.New(items, now))
	if err != nil {
		return err
	}
	if err = os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	// ExportLockPath reserves the destination for the duration of this write only.
	// It must be removed once the export completes (success or failure) so that
	// the same destination path can be reused by subsequent exports. O_EXCL makes
	// a leaked marker permanently block the destination, which is why cleanup
	// happens in a defer rather than only on the success path.
	lockPath := archive.ExportLockPath(path)
	lock, err := codec.AcquireExclusive(lockPath)
	if err != nil {
		return fmt.Errorf("reserve export destination: %w", err)
	}
	defer lock.Close()
	defer os.Remove(lockPath)
	if err = os.WriteFile(path, append(data, '\n'), 0644); err != nil {
		return err
	}
	return nil
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
