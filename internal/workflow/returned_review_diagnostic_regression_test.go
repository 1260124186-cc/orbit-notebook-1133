package workflow

import (
	"example.com/orbit-notebook/internal/clock"
	"example.com/orbit-notebook/internal/model"
	"example.com/orbit-notebook/internal/storage"
	"testing"
	"time"
)

func TestDiagnosticsFlagsReturnedReviewAsNotReleaseReady(t *testing.T) {
	st, err := storage.New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, 8, 24, 10, 0, 0, 0, time.UTC)
	item := model.Observation{ID: "obs-000081", Target: "M101", Instrument: "scope-y", ObservedAt: now.Add(-time.Hour), Description: "需补充记录", State: model.StateReviewed, CreatedAt: now.Add(-time.Hour), UpdatedAt: now, Review: &model.Review{Reviewer: "alice", Decision: model.DecisionReturn, Note: "补充细节", ReviewedAt: now}}
	if err = st.Replace(func(items []model.Observation) ([]model.Observation, error) { return append(items, item), nil }); err != nil {
		t.Fatal(err)
	}
	out, err := New(st, clock.Fixed{Value: now}).Diagnostics(item.ID)
	if err != nil {
		t.Fatal(err)
	}
	if out.Consistent {
		t.Fatalf("returned review was reported consistent: %#v", out)
	}
	for _, id := range out.ReleaseReady {
		if id == item.ID {
			t.Fatalf("returned review was reported release-ready: %#v", out.ReleaseReady)
		}
	}
}
