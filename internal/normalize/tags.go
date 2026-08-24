package normalize

import (
	"sort"
	"strings"
)

func Tags(values []string) []string {
	seen := make(map[string]bool, len(values))
	result := make([]string, 0, len(values))
	for _, value := range values {
		value = Token(value)
		if value != "" && !seen[value] {
			seen[value] = true
			result = append(result, value)
		}
	}
	sort.Strings(result)
	return result
}
func AddTags(current []string, additions ...string) []string {
	values := append(append([]string{}, current...), additions...)
	return Tags(values)
}
func RemoveTags(current []string, removals ...string) []string {
	excluded := map[string]bool{}
	for _, value := range removals {
		excluded[Token(value)] = true
	}
	kept := make([]string, 0, len(current))
	for _, value := range current {
		if !excluded[Token(value)] {
			kept = append(kept, value)
		}
	}
	return Tags(kept)
}
func TagString(values []string) string     { return strings.Join(Tags(values), ",") }
func ParseTagString(value string) []string { return Tags(strings.Split(value, ",")) }
