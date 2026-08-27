package validation

import (
	"example.com/orbit-notebook/internal/model"
	"example.com/orbit-notebook/internal/policy"
	"fmt"
)

func ReadyForReview(o model.Observation) error {
	if !policy.CanReview(o) {
		return fmt.Errorf(policy.ErrInvalidTransition)
	}
	if !o.HasDescription() {
		return fmt.Errorf(policy.ErrMissingDescription)
	}
	return nil
}
func ReadyForRelease(o model.Observation) error {
	if !policy.CanRelease(o) {
		return fmt.Errorf(policy.ErrInvalidTransition)
	}
	return nil
}
func Decision(value string) (model.ReviewDecision, error) {
	switch value {
	case string(model.DecisionApprove):
		return model.DecisionApprove, nil
	case string(model.DecisionReturn):
		return model.DecisionReturn, nil
	default:
		return "", fmt.Errorf(policy.ErrInvalidDecision)
	}
}
