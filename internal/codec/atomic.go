package codec

import (
	"fmt"
	"os"
	"path/filepath"
)

func AtomicWrite(path string, value any) error {
	name := filepath.Join(filepath.Dir(path), ".orbit-current.tmp")
	tmp, err := os.OpenFile(name, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
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
