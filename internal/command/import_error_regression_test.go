package command

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"example.com/orbit-notebook/internal/model"
	"example.com/orbit-notebook/internal/storage"
)

func TestImportRejectsInvalidArchiveAndPreservesState(t *testing.T) {
	home := t.TempDir()
	store, err := storage.New(home)
	if err != nil {
		t.Fatal(err)
	}
	before := model.Observation{
		ID: "obs-000001", Target: "M42", Instrument: "scope-a",
		ObservedAt:  time.Date(2026, 8, 24, 10, 0, 0, 0, time.UTC),
		Description: "原始记录", State: model.StateDraft,
		CreatedAt: time.Date(2026, 8, 24, 10, 1, 0, 0, time.UTC),
		UpdatedAt: time.Date(2026, 8, 24, 10, 1, 0, 0, time.UTC),
	}
	if err := store.Replace(func([]model.Observation) ([]model.Observation, error) {
		return []model.Observation{before}, nil
	}); err != nil {
		t.Fatal(err)
	}

	archivePath := filepath.Join(t.TempDir(), "invalid.json")
	document := map[string]any{
		"format":       "orbit-notebook/v1",
		"generated_at": "2026-08-24T10:00:00Z",
		"items": []model.Observation{{
			Target: "M31", Instrument: "scope-b",
			ObservedAt:  time.Date(2026, 8, 24, 11, 0, 0, 0, time.UTC),
			Description: "缺少标识", State: model.StateDraft,
			CreatedAt: time.Date(2026, 8, 24, 11, 1, 0, 0, time.UTC),
			UpdatedAt: time.Date(2026, 8, 24, 11, 1, 0, 0, time.UTC),
		}},
	}
	data, err := json.Marshal(document)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(archivePath, data, 0644); err != nil {
		t.Fatal(err)
	}

	if err := Run([]string{"import", "--home", home, "--path", archivePath}); err == nil {
		t.Fatal("expected invalid archive import to return an error")
	}
	items, err := store.Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 || items[0].ID != before.ID || items[0].Description != before.Description {
		t.Fatalf("invalid import changed persisted state: %#v", items)
	}
}
