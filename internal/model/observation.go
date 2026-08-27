package model

import (
	"sort"
	"strings"
	"time"
)

type State string

const (
	StateDraft    State = "draft"
	StateReviewed State = "reviewed"
	StateReleased State = "released"
)

type Observation struct {
	ID          string    `json:"id"`
	Target      string    `json:"target"`
	Instrument  string    `json:"instrument"`
	ObservedAt  time.Time `json:"observed_at"`
	Description string    `json:"description"`
	Tags        []string  `json:"tags"`
	State       State     `json:"state"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Review      *Review   `json:"review,omitempty"`
	Bulletin    *Bulletin `json:"bulletin,omitempty"`
}

func (o Observation) HasDescription() bool { return strings.TrimSpace(o.Description) != "" }
func (o Observation) HasTag(value string) bool {
	needle := strings.ToLower(strings.TrimSpace(value))
	for _, tag := range o.Tags {
		if strings.ToLower(tag) == needle {
			return true
		}
	}
	return false
}
func (o Observation) SortedTags() []string {
	result := append([]string(nil), o.Tags...)
	sort.Strings(result)
	return result
}
func (o Observation) Summary() Summary {
	return Summary{ID: o.ID, Target: o.Target, Instrument: o.Instrument, ObservedAt: o.ObservedAt, State: o.State, Tags: o.SortedTags(), Description: o.Description}
}

type Summary struct {
	ID          string    `json:"id"`
	Target      string    `json:"target"`
	Instrument  string    `json:"instrument"`
	ObservedAt  time.Time `json:"observed_at"`
	State       State     `json:"state"`
	Tags        []string  `json:"tags"`
	Description string    `json:"description"`
}

func SortSummaries(items []Summary) {
	sort.SliceStable(items, func(i, j int) bool {
		if items[i].ObservedAt.Equal(items[j].ObservedAt) {
			return items[i].ID < items[j].ID
		}
		return items[i].ObservedAt.Before(items[j].ObservedAt)
	})
}
