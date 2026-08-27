package workflow

import (
	"testing"
	"time"

	"example.com/orbit-notebook/internal/clock"
	"example.com/orbit-notebook/internal/model"
	"example.com/orbit-notebook/internal/storage"
)

func TestDiagnosticsSurvivesReleasedRecordWithoutBulletin(t *testing.T) {
	st, err := storage.New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, 8, 24, 10, 0, 0, 0, time.UTC)
	bad := model.Observation{ID: "obs-000071", Target: "M87", Instrument: "scope-z", ObservedAt: now.Add(-time.Hour), Description: "旧记录", State: model.StateReleased, CreatedAt: now.Add(-time.Hour), UpdatedAt: now, Review: &model.Review{Reviewer: "alice", Decision: model.DecisionApprove, Note: "旧数据", ReviewedAt: now}}
	if err = st.Replace(func(items []model.Observation) ([]model.Observation, error) { return append(items, bad), nil }); err != nil {
		t.Fatal(err)
	}
	svc := New(st, clock.Fixed{Value: now})
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("diagnostics panicked for malformed released record: %v", r)
		}
	}()
	out, err := svc.Diagnostics(bad.ID)
	if err != nil {
		t.Fatal(err)
	}
	if out.Records != 1 {
		t.Fatalf("records=%d", out.Records)
	}
}

func TestDiagnosticsSurvivesRecordWithoutCreationTime(t *testing.T) {
	st, err := storage.New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, 8, 24, 10, 0, 0, 0, time.UTC)
	bad := model.Observation{ID: "obs-000072", Target: "M89", Instrument: "scope-z", ObservedAt: now, Description: "旧记录", State: model.StateDraft}
	if err = st.Replace(func(items []model.Observation) ([]model.Observation, error) { return append(items, bad), nil }); err != nil {
		t.Fatal(err)
	}
	svc := New(st, clock.Fixed{Value: now})
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("diagnostics panicked for record without creation time: %v", r)
		}
	}()
	out, err := svc.Diagnostics(bad.ID)
	if err != nil {
		t.Fatal(err)
	}
	if out.Records != 1 {
		t.Fatalf("records=%d", out.Records)
	}
}
