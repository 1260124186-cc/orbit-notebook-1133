package command

import (
	"path/filepath"
	"testing"
)

func TestRepeatedExportReleasesDestination(t *testing.T) {
	home := t.TempDir()
	destination := filepath.Join(t.TempDir(), "archives", "observations.json")
	args := []string{"export", "--home", home, "--path", destination}
	if err := Run(args); err != nil {
		t.Fatalf("first export failed: %v", err)
	}
	if err := Run(args); err != nil {
		t.Fatalf("repeated export should succeed after the first call: %v", err)
	}
}
