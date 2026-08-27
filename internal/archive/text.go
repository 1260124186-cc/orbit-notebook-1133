package archive

import (
	"example.com/orbit-notebook/internal/model"
	"sort"
	"strings"
)

func Headline(item model.Observation) string {
	parts := []string{item.Target, item.Instrument}
	if item.Description != "" {
		parts = append(parts, item.Description)
	}
	return strings.Join(parts, " — ")
}
func Lines(items []model.Observation) []string {
	sorted := append([]model.Observation{}, items...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].ObservedAt.Before(sorted[j].ObservedAt) })
	result := make([]string, 0, len(sorted))
	for _, item := range sorted {
		result = append(result, item.ID+"\t"+Headline(item))
	}
	return result
}
func Markdown(items []model.Observation) string {
	lines := []string{"# Orbit Notebook", ""}
	for _, line := range Lines(items) {
		lines = append(lines, "- "+line)
	}
	return strings.Join(lines, "\n") + "\n"
}
func CountByState(items []model.Observation) map[model.State]int {
	result := map[model.State]int{}
	for _, item := range items {
		result[item.State]++
	}
	return result
}
