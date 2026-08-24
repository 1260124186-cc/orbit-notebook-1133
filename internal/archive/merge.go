package archive

import (
	"example.com/orbit-notebook/internal/integrity"
	"example.com/orbit-notebook/internal/model"
	"fmt"
	"sort"
)

func Merge(left, right []model.Observation) ([]model.Observation, error) {
	combined := append(append([]model.Observation{}, left...), right...)
	byID := map[string]model.Observation{}
	for _, item := range combined {
		normalized := integrity.Normalize(item)
		if existing, ok := byID[normalized.ID]; ok && existing.UpdatedAt.After(normalized.UpdatedAt) {
			continue
		}
		byID[normalized.ID] = normalized
	}
	result := make([]model.Observation, 0, len(byID))
	for _, item := range byID {
		result = append(result, item)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
	if err := integrity.Collection(result); err != nil {
		return result, fmt.Errorf("merged collection invalid: %w", err)
	}
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
