package query

import (
	"example.com/orbit-notebook/internal/model"
	"sort"
	"strings"
)

type SortMode string

const (
	SortObserved SortMode = "observed"
	SortUpdated  SortMode = "updated"
	SortTarget   SortMode = "target"
)

func Search(items []model.Observation, request Request, mode SortMode) Page {
	selected := Filter(items, request)
	sort.SliceStable(selected, func(i, j int) bool {
		switch mode {
		case SortTarget:
			if selected[i].Target == selected[j].Target {
				return selected[i].ID < selected[j].ID
			}
			return selected[i].Target < selected[j].Target
		case SortUpdated:
			return selected[i].UpdatedAt.Before(selected[j].UpdatedAt)
		default:
			if selected[i].ObservedAt.Equal(selected[j].ObservedAt) {
				return selected[i].ID < selected[j].ID
			}
			return selected[i].ObservedAt.Before(selected[j].ObservedAt)
		}
	})
	return Paginate(selected, request)
}
func Terms(value string) []string {
	words := strings.Fields(strings.ToLower(value))
	unique := make([]string, 0, len(words))
	seen := map[string]bool{}
	for _, word := range words {
		if !seen[word] {
			seen[word] = true
			unique = append(unique, word)
		}
	}
	return unique
}
func ContainsAll(item model.Observation, terms []string) bool {
	haystack := strings.ToLower(item.Target + " " + item.Instrument + " " + item.Description + " " + strings.Join(item.Tags, " "))
	for _, term := range terms {
		if !strings.Contains(haystack, strings.ToLower(term)) {
			return false
		}
	}
	return true
}
