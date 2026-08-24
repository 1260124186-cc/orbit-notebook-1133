package workflow

import (
	"example.com/orbit-notebook/internal/clock"
	"example.com/orbit-notebook/internal/index"
	"example.com/orbit-notebook/internal/storage"
	"example.com/orbit-notebook/internal/validation"
	"testing"
	"time"
)

func TestLifecycle(t *testing.T) {
	dir := t.TempDir()
	store, err := storage.New(dir)
	if err != nil {
		t.Fatal(err)
	}
	svc := New(store, clock.Fixed{Value: time.Date(2026, 8, 24, 10, 0, 0, 0, time.UTC)})
	o, err := svc.Capture(validation.CaptureInput{Target: "M42", Instrument: "scope-a", ObservedAt: time.Now(), Description: "清晰"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err = svc.Annotate(o.ID, "细节", []string{"nebula"}); err != nil {
		t.Fatal(err)
	}
	if _, err = svc.Review(o.ID, "alice", "approve", "完整"); err != nil {
		t.Fatal(err)
	}
	if _, err = svc.Release(o.ID, "bob"); err != nil {
		t.Fatal(err)
	}
	items, err := svc.List(index.Filter{State: "released"})
	if err != nil || len(items) != 1 {
		t.Fatalf("items=%v err=%v", items, err)
	}
}
func TestRejectsEmptyCapture(t *testing.T) {
	if err := validation.Capture(validation.CaptureInput{}); err == nil {
		t.Fatal("expected validation error")
	}
}
