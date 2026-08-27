package timeline

import (
	"example.com/orbit-notebook/internal/model"
	"fmt"
	"sort"
	"time"
)

type EventType string

const (
	Captured  EventType = "captured"
	Annotated EventType = "annotated"
	Reviewed  EventType = "reviewed"
	Released  EventType = "released"
)

type Event struct {
	Type  EventType `json:"type"`
	At    time.Time `json:"at"`
	Actor string    `json:"actor"`
	Note  string    `json:"note"`
}

func ForObservation(item model.Observation) []Event {
	events := []Event{}
	if !item.CreatedAt.IsZero() {
		events = append(events, Event{Type: Captured, At: item.CreatedAt, Note: "observation created"})
	}
	if item.UpdatedAt.After(item.CreatedAt) {
		events = append(events, Event{Type: Annotated, At: item.UpdatedAt, Note: "observation changed"})
	}
	if item.Review != nil {
		events = append(events, Event{Type: Reviewed, At: item.Review.ReviewedAt, Actor: item.Review.Reviewer, Note: item.Review.Note})
	}
	if item.State == model.StateReleased {
		events = append(events, Event{Type: Released, At: item.Bulletin.ReleasedAt, Actor: item.Bulletin.Publisher, Note: item.Bulletin.Headline})
	}
	sort.SliceStable(events, func(i, j int) bool { return events[i].At.Before(events[j].At) })
	return events
}
func Describe(event Event) string {
	if event.Actor == "" {
		return fmt.Sprintf("%s at %s", event.Type, event.At.Format(time.RFC3339))
	}
	return fmt.Sprintf("%s by %s at %s", event.Type, event.Actor, event.At.Format(time.RFC3339))
}
func Latest(events []Event) (Event, bool) {
	if len(events) == 0 {
		return Event{}, false
	}
	result := events[0]
	for _, event := range events[1:] {
		if event.At.After(result.At) {
			result = event
		}
	}
	return result, true
}
