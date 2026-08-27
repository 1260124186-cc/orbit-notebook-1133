package query

import (
	"example.com/orbit-notebook/internal/model"
	"strings"
)

type Request struct {
	State    model.State
	Target   string
	Tag      string
	Contains string
	Offset   int
	Limit    int
}
type Page struct {
	Items   []model.Summary `json:"items"`
	Offset  int             `json:"offset"`
	Limit   int             `json:"limit"`
	Total   int             `json:"total"`
	HasMore bool            `json:"has_more"`
}

func Match(item model.Observation, request Request) bool {
	if request.State != "" && item.State != request.State {
		return false
	}
	if request.Target != "" && !strings.EqualFold(item.Target, request.Target) {
		return false
	}
	if request.Tag != "" && !item.HasTag(request.Tag) {
		return false
	}
	if request.Contains != "" && !strings.Contains(strings.ToLower(item.Description), strings.ToLower(request.Contains)) {
		return false
	}
	return true
}
func Filter(items []model.Observation, request Request) []model.Observation {
	result := make([]model.Observation, 0, len(items))
	for _, item := range items {
		if Match(item, request) {
			result = append(result, item)
		}
	}
	return result
}
func Paginate(items []model.Observation, request Request) Page {
	if request.Offset < 0 {
		request.Offset = 0
	}
	if request.Limit <= 0 {
		request.Limit = 20
	}
	total := len(items)
	start := request.Offset
	if start > total {
		start = total
	}
	end := start + request.Limit
	if end > total {
		end = total
	}
	summaries := make([]model.Summary, 0, end-start)
	for _, item := range items[start:end] {
		summaries = append(summaries, item.Summary())
	}
	return Page{Items: summaries, Offset: request.Offset, Limit: request.Limit, Total: total, HasMore: end < total}
}
func IDs(items []model.Observation) []string {
	result := make([]string, 0, len(items))
	for _, item := range items {
		result = append(result, item.ID)
	}
	return result
}
