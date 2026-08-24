package stats

import (
	"example.com/orbit-notebook/internal/model"
	"sort"
)

type TagCount struct {
	Tag   string `json:"tag"`
	Count int    `json:"count"`
}

func ByTag(items []model.Observation) []TagCount {
	counts := map[string]int{}
	for _, item := range items {
		for _, tag := range item.Tags {
			counts[tag]++
		}
	}
	result := make([]TagCount, 0, len(counts))
	for tag, count := range counts {
		result = append(result, TagCount{Tag: tag, Count: count})
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].Count == result[j].Count {
			return result[i].Tag < result[j].Tag
		}
		return result[i].Count > result[j].Count
	})
	return result
}
func HasDuplicateIDs(items []model.Observation) bool {
	seen := map[string]bool{}
	for _, item := range items {
		if seen[item.ID] {
			return true
		}
		seen[item.ID] = true
	}
	return false
}
func CountReleased(items []model.Observation) int {
	result := 0
	for _, item := range items {
		if item.State == model.StateReleased {
			result++
		}
	}
	return result
}
