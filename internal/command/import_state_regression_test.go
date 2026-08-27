package command

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"example.com/orbit-notebook/internal/archive"
	"example.com/orbit-notebook/internal/model"
	"example.com/orbit-notebook/internal/storage"
)

func TestImportAppliesNewerObservationState(t *testing.T) {
	home := t.TempDir()
	store, err := storage.New(home)
	if err != nil {
		t.Fatal(err)
	}
	oldTime := time.Date(2026, 8, 24, 10, 0, 0, 0, time.UTC)
	old := model.Observation{
		ID: "obs-000001", Target: "M42", Instrument: "scope-a",
		ObservedAt: oldTime, Description: "旧描述", State: model.StateDraft,
		CreatedAt: oldTime, UpdatedAt: oldTime,
	}
	if err := store.Replace(func([]model.Observation) ([]model.Observation, error) {
		return []model.Observation{old}, nil
	}); err != nil {
		t.Fatal(err)
	}

	newTime := oldTime.Add(time.Hour)
	newer := old
	newer.Description = "导入的新描述"
	newer.State = model.StateReviewed
	newer.UpdatedAt = newTime
	archivePath := filepath.Join(t.TempDir(), "newer.json")
	data, err := archive.Encode(archive.New([]model.Observation{newer}, newTime.Add(time.Minute)))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(archivePath, append(data, '\n'), 0644); err != nil {
		t.Fatal(err)
	}

	if err := Run([]string{"import", "--home", home, "--path", archivePath}); err != nil {
		t.Fatalf("import failed: %v", err)
	}
	items, err := store.Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 || items[0].Description != newer.Description || items[0].State != newer.State || !items[0].UpdatedAt.Equal(newTime) {
		t.Fatalf("newer imported observation was not applied: %#v", items)
	}
}
