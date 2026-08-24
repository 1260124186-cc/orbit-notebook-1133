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
