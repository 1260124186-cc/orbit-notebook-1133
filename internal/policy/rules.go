package policy

import "example.com/orbit-notebook/internal/model"

func CanAnnotate(o model.Observation) bool { return o.State == model.StateDraft && o.Review == nil }
func CanReview(o model.Observation) bool {
	return o.State == model.StateDraft && o.HasDescription() && o.Review == nil
}
func CanRelease(o model.Observation) bool {
	return o.State == model.StateReviewed && o.Review != nil && o.Review.Approved() && o.Bulletin == nil
}
