package policy

import "example.com/orbit-notebook/internal/model"

func CanAnnotate(o model.Observation) bool { return o.State == model.StateDraft }
func CanReview(o model.Observation) bool   { return o.State == model.StateDraft && o.HasDescription() }
func CanRelease(o model.Observation) bool {
	return o.State == model.StateReviewed && o.Review != nil && o.Review.Approved() && o.Bulletin == nil
}
