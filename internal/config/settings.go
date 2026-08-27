package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Settings struct {
	Home              string
	DefaultInstrument string
	DefaultPublisher  string
	DefaultReviewer   string
	MaxDescription    int
}

func Default() Settings {
	return Settings{Home: ".orbit-notebook", DefaultInstrument: "unknown", DefaultPublisher: "orbit-editorial", DefaultReviewer: "orbit-review", MaxDescription: 2000}
}

func FromEnvironment() Settings {
	settings := Default()
	if value := strings.TrimSpace(os.Getenv("ORBIT_NOTEBOOK_HOME")); value != "" {
		settings.Home = value
	}
	if value := strings.TrimSpace(os.Getenv("ORBIT_NOTEBOOK_INSTRUMENT")); value != "" {
		settings.DefaultInstrument = value
	}
	if value := strings.TrimSpace(os.Getenv("ORBIT_NOTEBOOK_PUBLISHER")); value != "" {
		settings.DefaultPublisher = value
	}
	if value := strings.TrimSpace(os.Getenv("ORBIT_NOTEBOOK_REVIEWER")); value != "" {
		settings.DefaultReviewer = value
	}
	return settings
}

func (s Settings) Validate() error {
	if strings.TrimSpace(s.Home) == "" {
		return fmt.Errorf("home is required")
	}
	if s.MaxDescription < 64 {
		return fmt.Errorf("max description is too small")
	}
	return nil
}

func (s Settings) AbsoluteHome(base string) string {
	if filepath.IsAbs(s.Home) {
		return filepath.Clean(s.Home)
	}
	if base == "" {
		base = "."
	}
	return filepath.Clean(filepath.Join(base, s.Home))
}
