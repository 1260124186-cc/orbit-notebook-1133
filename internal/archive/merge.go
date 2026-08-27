package archive

import (
	"example.com/orbit-notebook/internal/integrity"
	"example.com/orbit-notebook/internal/model"
	"sort"
)

func Merge(left, right []model.Observation) ([]model.Observation, error) {
	byID := map[string]model.Observation{}
	for _, item := range left {
		byID[item.ID] = item
	}
	for _, item := range right {
		normalized := integrity.Normalize(item)
		if existing, ok := byID[normalized.ID]; ok && !existing.UpdatedAt.Before(normalized.UpdatedAt) {
			continue
		}
		byID[normalized.ID] = normalized
	}
	result := make([]model.Observation, 0, len(byID))
	for _, item := range byID {
		result = append(result, item)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
	return result, nil
}
func Difference(base, current []model.Observation) []model.Observation {
	known := map[string]model.Observation{}
	for _, item := range base {
		known[item.ID] = item
	}
	result := []model.Observation{}
	for _, item := range current {
		old, ok := known[item.ID]
		if !ok || item.UpdatedAt.After(old.UpdatedAt) {
			result = append(result, item)
		}
	}
	return result
}
