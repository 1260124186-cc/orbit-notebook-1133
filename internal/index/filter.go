package index

import (
	"example.com/orbit-notebook/internal/model"
	"strings"
)

type Filter struct {
	State  model.State
	Target string
	Tag    string
}

func Apply(items []model.Observation, f Filter) []model.Summary {
	result := make([]model.Summary, 0, len(items))
	for _, item := range items {
		if f.State != "" && item.State != f.State {
			continue
		}
		if f.Target != "" && !strings.EqualFold(item.Target, f.Target) {
			continue
		}
		if f.Tag != "" && !item.HasTag(f.Tag) {
			continue
		}
		result = append(result, item.Summary())
	}
	model.SortSummaries(result)
	return result
}
