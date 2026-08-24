package workflow

import (
	"example.com/orbit-notebook/internal/clock"
	"example.com/orbit-notebook/internal/storage"
	"example.com/orbit-notebook/internal/validation"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestConcurrentCapturePreservesCollectionAndIDs(t *testing.T) {
	dir := t.TempDir()
	store, err := storage.New(dir)
	if err != nil {
		t.Fatal(err)
	}
	svc := New(store, clock.Fixed{Value: time.Date(2026, 8, 24, 10, 0, 0, 0, time.UTC)})
	const workers = 32
	var wg sync.WaitGroup
	errs := make(chan error, workers)
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_, err := svc.Capture(validation.CaptureInput{
				Target: fmt.Sprintf("M-%02d", i), Instrument: "scope-a",
				ObservedAt: time.Date(2026, 8, 24, 10, i, 0, 0, time.UTC),
				Description: "并发观测",
			})
			if err != nil {
				errs <- err
			}
		}(i)
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		t.Errorf("capture failed: %v", err)
	}
	items, err := store.Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != workers {
		t.Fatalf("concurrent capture lost records: got %d want %d", len(items), workers)
	}
	seen := map[string]bool{}
	for _, item := range items {
		if seen[item.ID] {
			t.Fatalf("duplicate observation id %q", item.ID)
		}
		seen[item.ID] = true
	}
	for i := 1; i <= workers; i++ {
		id := fmt.Sprintf("obs-%06d", i)
		if !seen[id] {
			t.Fatalf("missing contiguous observation id %q", id)
		}
	}
}

func TestCaptureStillSupportsSequentialIDs(t *testing.T) {
	dir := t.TempDir()
	store, err := storage.New(dir)
	if err != nil {
		t.Fatal(err)
	}
	svc := New(store, clock.Fixed{Value: time.Date(2026, 8, 24, 10, 0, 0, 0, time.UTC)})
	for i := 0; i < 3; i++ {
		o, err := svc.Capture(validation.CaptureInput{Target: fmt.Sprintf("M-%d", i), Instrument: "scope-a", ObservedAt: time.Date(2026, 8, 24, 10, i, 0, 0, time.UTC), Description: "单次观测"})
		if err != nil {
			t.Fatal(err)
		}
		want := fmt.Sprintf("obs-%06d", i+1)
		if o.ID != want {
			t.Fatalf("got id %q want %q", o.ID, want)
		}
	}
}
