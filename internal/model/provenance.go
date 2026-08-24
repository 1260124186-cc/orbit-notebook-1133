package model

import (
	"sort"
	"strings"
	"time"
)

// Provenance records the human and session context attached to a local
// observation import. It is deliberately small so it can be displayed in
// diagnostics without exposing the full storage representation.
type Provenance struct {
	Observer  string    `json:"observer"`
	Session   string    `json:"session"`
	Site      string    `json:"site"`
	StartedAt time.Time `json:"started_at"`
	EndedAt   time.Time `json:"ended_at"`
}

func NewProvenance(observer, session, site string, started time.Time) Provenance {
	return Provenance{Observer: strings.TrimSpace(observer), Session: strings.TrimSpace(session), Site: strings.TrimSpace(site), StartedAt: started.UTC()}
}

func (p Provenance) Completed(at time.Time) Provenance {
	p.EndedAt = at.UTC()
	return p
}

func (p Provenance) Valid() bool {
	return p.Observer != "" && p.Session != "" && p.Site != "" && !p.StartedAt.IsZero() && !p.EndedAt.IsZero() && !p.EndedAt.Before(p.StartedAt)
}

func (p Provenance) Duration() time.Duration {
	if p.EndedAt.Before(p.StartedAt) {
		return 0
	}
	return p.EndedAt.Sub(p.StartedAt)
}

func (p Provenance) Label() string {
	parts := []string{p.Observer, p.Session, p.Site}
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return strings.Join(parts, " / ")
}

func (p Provenance) SameSession(other Provenance) bool {
	return p.Session != "" && p.Session == other.Session
}

func CombineProvenance(values ...Provenance) (Provenance, bool) {
	if len(values) == 0 {
		return Provenance{}, false
	}
	result := values[0]
	for _, value := range values[1:] {
		if value.StartedAt.Before(result.StartedAt) {
			result.StartedAt = value.StartedAt
		}
		if value.EndedAt.After(result.EndedAt) {
			result.EndedAt = value.EndedAt
		}
		if result.Observer == "" {
			result.Observer = value.Observer
		}
		if result.Site == "" {
			result.Site = value.Site
		}
	}
	return result, result.Valid()
}

func SortProvenance(values []Provenance) []Provenance {
	result := append([]Provenance{}, values...)
	sort.SliceStable(result, func(i, j int) bool {
		if result[i].StartedAt.Equal(result[j].StartedAt) {
			return result[i].Session < result[j].Session
		}
		return result[i].StartedAt.Before(result[j].StartedAt)
	})
	return result
}

func SessionIDs(values []Provenance) []string {
	result := make([]string, 0, len(values))
	seen := map[string]bool{}
	for _, value := range values {
		if value.Session == "" || seen[value.Session] {
			continue
		}
		seen[value.Session] = true
		result = append(result, value.Session)
	}
	sort.Strings(result)
	return result
}
