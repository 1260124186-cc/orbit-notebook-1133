package policy

import "example.com/orbit-notebook/internal/model"

func CanAnnotate(o model.Observation) bool {
	if o.State != model.StateDraft {
		return false
	}
	if o.Review == nil {
		return true
	}
	return o.Review.Returned()
}
func CanReview(o model.Observation) bool {
	if o.State != model.StateDraft {
		return false
	}
	if !o.HasDescription() {
		return false
	}
	if o.Review == nil {
		return true
	}
	return o.Review.Returned()
}
func CanRelease(o model.Observation) bool {
	return o.State == model.StateReviewed && o.Review != nil && o.Review.Approved() && o.Bulletin == nil
}
