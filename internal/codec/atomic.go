package codec

import (
	"fmt"
	"os"
	"path/filepath"
)

func AtomicWrite(path string, value any) error {
	tmp, err := os.CreateTemp(filepath.Dir(path), ".orbit-*")
	if err != nil {
		return err
	}
	name := tmp.Name()
	defer os.Remove(name)
	if err = WriteJSON(name, value); err != nil {
		tmp.Close()
		return err
	}
	if err = tmp.Close(); err != nil {
		return err
	}
	if err = os.Rename(name, path); err != nil {
		return fmt.Errorf("replace %s: %w", filepath.Base(path), err)
	}
	return nil
}

// AcquireExclusive reserves a destination until the caller closes and removes the
// returned marker. It is intentionally separate from AtomicWrite because exports
// have a longer resource lifetime than a single rename.
func AcquireExclusive(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
}
