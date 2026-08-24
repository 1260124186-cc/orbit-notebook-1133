package config

import "strings"

type LabelSet struct{ Values []string }

func NewLabels(values ...string) LabelSet {
	result := LabelSet{Values: make([]string, 0, len(values))}
	seen := map[string]struct{}{}
	for _, value := range values {
		value = strings.ToLower(strings.TrimSpace(value))
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		result.Values = append(result.Values, value)
	}
	return result
}

func (s LabelSet) Contains(value string) bool {
	needle := strings.ToLower(strings.TrimSpace(value))
	for _, item := range s.Values {
		if item == needle {
			return true
		}
	}
	return false
}

func (s LabelSet) Merge(other LabelSet) LabelSet {
	return NewLabels(append(append([]string{}, s.Values...), other.Values...)...)
}
func (s LabelSet) Empty() bool { return len(s.Values) == 0 }
