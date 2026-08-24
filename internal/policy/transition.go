package policy

import "example.com/orbit-notebook/internal/model"

func StateName(state model.State) string {
	switch state {
	case model.StateDraft:
		return "draft"
	case model.StateReviewed:
		return "reviewed"
	case model.StateReleased:
		return "released"
	default:
		return "unknown"
	}
}
func CanMove(from, to model.State) bool {
	if from == model.StateDraft && (to == model.StateDraft || to == model.StateReviewed) {
		return true
	}
	if from == model.StateReviewed && (to == model.StateReviewed || to == model.StateReleased || to == model.StateDraft) {
		return true
	}
	return from == to
}
func Terminal(state model.State) bool { return state == model.StateReleased }
func EditableFields(state model.State) []string {
	if state != model.StateDraft {
		return []string{}
	}
	return []string{"description", "tags"}
}
func AllowedDecisions() []model.ReviewDecision {
	return []model.ReviewDecision{model.DecisionApprove, model.DecisionReturn}
}
