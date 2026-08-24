package normalize

import (
	"strings"
	"unicode"
)

func Text(value string) string       { return strings.Join(strings.Fields(strings.TrimSpace(value)), " ") }
func Target(value string) string     { return strings.ToUpper(Text(value)) }
func Instrument(value string) string { return strings.ToLower(Text(value)) }
func Person(value string) string     { return Text(value) }
func Token(value string) string      { return strings.ToLower(Text(value)) }
func HasVisibleRune(value string) bool {
	for _, r := range value {
		if !unicode.IsSpace(r) {
			return true
		}
	}
	return false
}
func Words(value string) []string { return strings.Fields(Text(value)) }
func JoinSentences(values ...string) string {
	result := make([]string, 0, len(values))
	for _, value := range values {
		if value = Text(value); value != "" {
			result = append(result, value)
		}
	}
	return strings.Join(result, ". ")
}
