package validation

import (
	"example.com/orbit-notebook/internal/policy"
	"fmt"
	"strings"
	"time"
)

type CaptureInput struct {
	Target, Instrument, Description string
	ObservedAt                      time.Time
}

func Capture(in CaptureInput) error {
	if strings.TrimSpace(in.Target) == "" {
		return fmt.Errorf(policy.ErrMissingTarget)
	}
	if strings.TrimSpace(in.Instrument) == "" {
		return fmt.Errorf(policy.ErrMissingInstrument)
	}
	if in.ObservedAt.IsZero() {
		return fmt.Errorf("observed-at is required")
	}
	if strings.TrimSpace(in.Description) == "" {
		return fmt.Errorf(policy.ErrMissingDescription)
	}
	return nil
}
func Annotation(description string) error {
	if strings.TrimSpace(description) == "" {
		return fmt.Errorf(policy.ErrMissingDescription)
	}
	return nil
}
func Reviewer(name, note string) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("reviewer is required")
	}
	if strings.TrimSpace(note) == "" {
		return fmt.Errorf("review note is required")
	}
	return nil
}
