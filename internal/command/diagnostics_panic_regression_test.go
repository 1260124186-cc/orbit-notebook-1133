package command

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"example.com/orbit-notebook/internal/model"
)

func TestDiagnosticsHandlesPersistedObservationWithoutCreationTime(t *testing.T) {
	home := t.TempDir()
	item := model.Observation{
		ID: "obs-000001", Target: "M42", Instrument: "scope-a",
		ObservedAt:  time.Date(2026, 8, 24, 10, 0, 0, 0, time.UTC),
		Description: "观测记录", State: model.StateDraft,
		UpdatedAt: time.Date(2026, 8, 24, 10, 1, 0, 0, time.UTC),
	}
	data, err := json.Marshal([]model.Observation{item})
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(home, "observations.json"), data, 0644); err != nil {
		t.Fatal(err)
	}
	if err := Run([]string{"diagnostics", "--home", home}); err != nil {
		t.Fatalf("diagnostics should report the damaged record instead of failing: %v", err)
	}
}
