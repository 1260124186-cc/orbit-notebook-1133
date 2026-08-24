package storage_test

import (
	"example.com/orbit-notebook/internal/archive"
	"example.com/orbit-notebook/internal/clock"
	"example.com/orbit-notebook/internal/storage"
	"example.com/orbit-notebook/internal/validation"
	"example.com/orbit-notebook/internal/workflow"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestExportedArchiveRoundTripsWithCurrentContract(t *testing.T) {
	now := time.Date(2026, 8, 24, 10, 0, 0, 0, time.UTC)
	source, err := storage.New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	svc := workflow.New(source, clock.Fixed{Value: now})
	if _, err = svc.Capture(validation.CaptureInput{Target: "M13", Instrument: "scope-c", ObservedAt: now.Add(-time.Hour), Description: "圆球状星团"}); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "archive.json")
	if err = svc.Export(path); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	doc, err := archive.Decode(data)
	if err != nil {
		t.Fatal(err)
	}
	if doc.Format != "orbit-notebook/v1" {
		t.Fatalf("format=%q want current v1", doc.Format)
	}
	if !doc.GeneratedAt.Equal(now) {
		t.Fatalf("generated_at=%s want %s", doc.GeneratedAt, now)
	}
	target, err := storage.New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if err = target.Import(path); err != nil {
		t.Fatalf("current exported archive cannot be imported: %v", err)
	}
	items, err := target.Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 || items[0].Target != "M13" {
		t.Fatalf("round trip items=%#v", items)
	}
}
