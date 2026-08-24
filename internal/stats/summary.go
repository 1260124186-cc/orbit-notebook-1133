package stats

import (
	"example.com/orbit-notebook/internal/model"
	"sort"
	"time"
)

type Overview struct {
	Total         int        `json:"total"`
	Draft         int        `json:"draft"`
	Reviewed      int        `json:"reviewed"`
	Released      int        `json:"released"`
	FirstObserved *time.Time `json:"first_observed,omitempty"`
	LastObserved  *time.Time `json:"last_observed,omitempty"`
	Targets       int        `json:"targets"`
	Tags          int        `json:"tags"`
}

type TargetCount struct {
	Target string `json:"target"`
	Count  int    `json:"count"`
}

func Build(items []model.Observation) Overview {
	result := Overview{}
	targets := map[string]bool{}
	tags := map[string]bool{}
	var first, last time.Time
	for _, item := range items {
		result.Total++
		targets[item.Target] = true
		for _, tag := range item.Tags {
			tags[tag] = true
		}
		switch item.State {
		case model.StateDraft:
			result.Draft++
		case model.StateReviewed:
			result.Reviewed++
		case model.StateReleased:
			result.Released++
		}
		if first.IsZero() || item.ObservedAt.Before(first) {
			first = item.ObservedAt
		}
		if last.IsZero() || item.ObservedAt.After(last) {
			last = item.ObservedAt
		}
	}
	result.Targets, result.Tags = len(targets), len(tags)
	if !first.IsZero() {
		result.FirstObserved = &first
	}
	if !last.IsZero() {
		result.LastObserved = &last
	}
	return result
}
func ByTarget(items []model.Observation) []TargetCount {
	counts := map[string]int{}
	for _, item := range items {
		counts[item.Target]++
	}
	result := make([]TargetCount, 0, len(counts))
	for target, count := range counts {
		result = append(result, TargetCount{Target: target, Count: count})
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].Count == result[j].Count {
			return result[i].Target < result[j].Target
		}
		return result[i].Count > result[j].Count
	})
	return result
}
func StateCounts(items []model.Observation) map[model.State]int {
	result := map[model.State]int{}
	for _, item := range items {
		result[item.State]++
	}
	return result
}
func Recent(items []model.Observation, limit int) []model.Summary {
	if limit < 1 {
		return nil
	}
	sorted := append([]model.Observation{}, items...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].ObservedAt.After(sorted[j].ObservedAt) })
	if limit > len(sorted) {
		limit = len(sorted)
	}
	out := make([]model.Summary, 0, limit)
	for _, item := range sorted[:limit] {
		out = append(out, item.Summary())
	}
	return out
}
