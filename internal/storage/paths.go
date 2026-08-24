package storage

import (
	"os"
	"path/filepath"
)

func ResolveHome(value string) string {
	if value != "" {
		return value
	}
	if env := os.Getenv("ORBIT_NOTEBOOK_HOME"); env != "" {
		return env
	}
	return filepath.Join(".", ".orbit-notebook")
}
