package workflow

import (
	"example.com/orbit-notebook/internal/clock"
	"example.com/orbit-notebook/internal/integrity"
	"example.com/orbit-notebook/internal/model"
	"example.com/orbit-notebook/internal/storage"
	"example.com/orbit-notebook/internal/validation"
	"testing"
	"time"
)

func TestInvalidReviewsDoNotPersistOrAdvanceState(t *testing.T) {
	store, err := storage.New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	svc := New(store, clock.Fixed{Value: time.Date(2026, 8, 24, 10, 0, 0, 0, time.UTC)})
	good, err := svc.Capture(validation.CaptureInput{Target: "M31", Instrument: "scope-b", ObservedAt: time.Now(), Description: "可见旋臂"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err = svc.Review(good.ID, "alice", "approve", ""); err == nil {
		t.Fatal("blank review note was accepted")
	}
	after, err := svc.Show(good.ID)
	if err != nil {
		t.Fatal(err)
	}
	if after.State != model.StateDraft || after.Review != nil {
		t.Fatalf("invalid review changed record %#v", after)
	}
	incomplete := model.Observation{ID: "obs-900000", Target: "M51", Instrument: "scope-d", ObservedAt: time.Now(), State: model.StateDraft, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	if err = store.Replace(func(items []model.Observation) ([]model.Observation, error) { return append(items, incomplete), nil }); err != nil {
		t.Fatal(err)
	}
	if _, err = svc.Review(incomplete.ID, "bob", "approve", "信息完整"); err == nil {
		t.Fatal("description-less record was reviewed")
	}
	if err = integrity.Observation(model.Observation{ID: "obs-900001", Target: "M1", Instrument: "scope-a", ObservedAt: time.Now(), State: model.StateDraft, Review: &model.Review{Reviewer: "alice", Decision: model.DecisionApprove}}); err == nil {
		t.Fatal("blank review note passed integrity")
	}
	returned, err := svc.Review(good.ID, "alice", "return", "补充细节")
	if err != nil {
		t.Fatal(err)
	}
	if returned.State != model.StateDraft {
		t.Fatalf("returned review advanced state to %q", returned.State)
	}
}
