package integrity

import (
	"example.com/orbit-notebook/internal/model"
	"fmt"
	"strings"
)

func Observation(item model.Observation) error {
	if strings.TrimSpace(item.ID) == "" {
		return fmt.Errorf("observation id is required")
	}
	if strings.TrimSpace(item.Target) == "" {
		return fmt.Errorf("observation target is required")
	}
	if strings.TrimSpace(item.Instrument) == "" {
		return fmt.Errorf("observation instrument is required")
	}
	if item.ObservedAt.IsZero() {
		return fmt.Errorf("observation time is required")
	}
	if item.State != model.StateDraft && item.State != model.StateReviewed && item.State != model.StateReleased {
		return fmt.Errorf("unknown observation state %q", item.State)
	}
	if item.Review != nil && !item.Review.Valid() {
		return fmt.Errorf("invalid review")
	}
	if item.Bulletin != nil && !item.Bulletin.Valid() {
		return fmt.Errorf("invalid bulletin")
	}
	return nil
}
func Collection(items []model.Observation) error {
	seen := map[string]bool{}
	for _, item := range items {
		if seen[item.ID] {
			return fmt.Errorf("duplicate observation id %q", item.ID)
		}
		seen[item.ID] = true
		if err := Observation(item); err != nil {
			return err
		}
	}
	return nil
}
func Consistent(item model.Observation) error {
	if err := Observation(item); err != nil {
		return err
	}
	switch item.State {
	case model.StateDraft:
		if item.Review != nil {
			return fmt.Errorf("draft cannot retain review")
		}
		if item.Bulletin != nil {
			return fmt.Errorf("draft cannot have bulletin")
		}
	case model.StateReviewed:
		if item.Review == nil || !item.Review.Approved() {
			return fmt.Errorf("reviewed item must have approval")
		}
		if item.Bulletin != nil {
			return fmt.Errorf("reviewed item cannot have bulletin")
		}
	case model.StateReleased:
		if item.Review == nil || !item.Review.Approved() || item.Bulletin == nil {
			return fmt.Errorf("released item requires approval and bulletin")
		}
	}
	return nil
}
