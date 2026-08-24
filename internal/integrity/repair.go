package integrity

import (
	"example.com/orbit-notebook/internal/model"
	"example.com/orbit-notebook/internal/normalize"
	"sort"
	"strings"
)

func Normalize(item model.Observation) model.Observation {
	item.ID = strings.TrimSpace(item.ID)
	item.Target = normalize.Target(item.Target)
	item.Instrument = normalize.Instrument(item.Instrument)
	item.Description = normalize.Text(item.Description)
	item.Tags = normalize.Tags(item.Tags)
	if item.Review != nil {
		item.Review.Reviewer = normalize.Person(item.Review.Reviewer)
		item.Review.Note = normalize.Text(item.Review.Note)
	}
	return item
}
func Sort(items []model.Observation) []model.Observation {
	result := append([]model.Observation{}, items...)
	sort.SliceStable(result, func(i, j int) bool { return result[i].ID < result[j].ID })
	return result
}
func MissingDescription(items []model.Observation) []string {
	result := []string{}
	for _, item := range items {
		if !item.HasDescription() {
			result = append(result, item.ID)
		}
	}
	return result
}
func MissingReview(items []model.Observation) []string {
	result := []string{}
	for _, item := range items {
		if item.Review == nil {
			result = append(result, item.ID)
		}
	}
	return result
}
